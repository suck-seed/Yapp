package services

import (
	"context"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/suck-seed/yapp/internal/auth"
	"github.com/suck-seed/yapp/internal/database"
	dto "github.com/suck-seed/yapp/internal/dto/message"
	"github.com/suck-seed/yapp/internal/models"
	"github.com/suck-seed/yapp/internal/repositories"
	"github.com/suck-seed/yapp/internal/utils"
)

type IMessageService interface {
	CreateMessage(c context.Context, req *dto.CreateMessageReq) (*dto.CreateMessageRes, error)
	FetchMessages(c context.Context, req *dto.MessageQueryParams) (*dto.MessageListResponse, error)
}

type messageService struct {
	repositories.IHallRepository
	repositories.IRoomRepository
	repositories.IMessageRepository
	repositories.IUserRepository
	pool *pgxpool.Pool

	timeout time.Duration
	mu      sync.RWMutex
}

func NewMessageService(hallRepo repositories.IHallRepository, roomRepo repositories.IRoomRepository, messageRepo repositories.IMessageRepository, userRepo repositories.IUserRepository, pool *pgxpool.Pool) IMessageService {
	return &messageService{
		hallRepo,
		roomRepo,
		messageRepo,
		userRepo,
		pool,
		time.Duration(2) * time.Second,
		sync.RWMutex{},
	}
}

// TODO Properly fill in these methods

func (s *messageService) CreateMessage(c context.Context, req *dto.CreateMessageReq) (*dto.CreateMessageRes, error) {
	ctx, cancel := context.WithTimeout(c, s.timeout)
	defer cancel()

	// --------------- TRANSACTION INIT
	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return nil, utils.ErrorInternal
	}
	runner := database.NewTxWrapper(tx)
	defer runner.Rollback(ctx)

	// validate content
	normalizedContent := utils.SanitizeMessageContent(req.Content)

	// create a uuid
	messageId, err := uuid.NewV7()
	if err != nil {
		return nil, utils.ErrorInternal
	}

	// create a message object
	message := &models.Message{
		ID:       messageId,
		RoomId:   req.RoomID,
		AuthorId: req.AuthorID,
		Content:  normalizedContent,
		SentAt:   req.SentAt,
	}

	message.MentionEveryone = false
	if req.MentionEveryone != nil && *req.MentionEveryone == true {
		message.MentionEveryone = true
	}

	// call CreateMessage
	messageCRES, err := s.IMessageRepository.CreateMessage(ctx, runner, message)
	if err != nil {
		return nil, utils.ErrorWritingMessage
	}

	// else , create entry on message_mentions
	var mentions []dto.UserBasic
	if *req.MentionEveryone == false && req.Mentions != nil {

		// validate if the userId acc exists or not
		for _, mentionedUserID := range *req.Mentions {

			exists, err := s.IUserRepository.DoesUserExists(ctx, runner, &mentionedUserID)
			if err != nil {
				return nil, utils.ErrorInternal
			}

			if exists {

				// Add to table message_mentions
				if err := s.IMessageRepository.AddMessageMention(ctx, runner, messageCRES.ID, mentionedUserID); err != nil {
					return nil, utils.ErrorWritingMentions
				}

				// get userBasic information for generation
				userCRES, err := s.IUserRepository.GetUserById(ctx, runner, &mentionedUserID)
				if err != nil {
					return nil, utils.ErrorFetchingUser
				}

				mentions = append(mentions, dto.UserBasic{
					ID:        userCRES.ID,
					Username:  userCRES.Username,
					Email:     userCRES.Email,
					AvatarURL: userCRES.AvatarURL,
				})
			}

		}

	}

	// attachments check
	var attachments []models.Attachment
	if req.Attachments != nil && *req.Attachments != nil && len(*req.Attachments) > 0 {

		for _, currentAttachment := range *req.Attachments {

			//			validate fileName
			canonFileName, err := utils.ValidateFileName(currentAttachment.FileName)
			if err != nil {
				return nil, err
			}

			validatedFileType, err := utils.ValidateFileType(currentAttachment.FileType, currentAttachment.URL)
			if err != nil {
				return nil, err
			}

			//			validate file size
			if *currentAttachment.FileSize >= utils.FileSize {
				return nil, utils.ErrorLargeFileSize

			}

			//			attachment_id
			attachmentID, err := uuid.NewV7()
			if err != nil {
				return nil, utils.ErrorInternal
			}

			//			repo call
			attachmentCRES, err := s.IMessageRepository.AddAttachment(ctx, runner, &models.Attachment{
				ID:        attachmentID,
				MessageID: messageId,
				FileName:  canonFileName,
				URL:       currentAttachment.URL,
				FileType:  validatedFileType,
				FileSize:  currentAttachment.FileSize,
			})

			if err != nil {
				return nil, utils.ErrorCreatingAttachment
			}

			//	append to attachments
			attachments = append(attachments, *attachmentCRES)
		}

	}

	// --------------- COMMIT
	if err := runner.Commit(ctx); err != nil {
		return nil, utils.ErrorInternal
	}

	return &dto.CreateMessageRes{
		ID:               messageCRES.ID,
		RoomID:           messageCRES.RoomId,
		AuthorID:         messageCRES.AuthorId,
		Content:          messageCRES.Content,
		SentAt:           messageCRES.SentAt,
		MentionsEveryone: messageCRES.MentionEveryone,
		Mentions:         mentions,
		Attachments:      attachments,

		CreatedAt: messageCRES.CreatedAt,
		EditedAt:  messageCRES.EditedAt,
		DeletedAt: messageCRES.DeletedAt,
		UpdatedAt: messageCRES.UpdatedAt,
	}, nil

}

func (s *messageService) FetchMessages(c context.Context, req *dto.MessageQueryParams) (*dto.MessageListResponse, error) {

	ctx, cancel := context.WithTimeout(c, s.timeout)
	defer cancel()

	// --------------- CONNECTION INIT
	conn, err := s.pool.Acquire(ctx)
	if err != nil {
		return nil, utils.ErrorInternal
	}
	defer conn.Release()
	runner := database.NewConnWrapper(conn)

	// validate only single cursor is used
	cursorCount := 0

	if req.Before != nil {
		cursorCount++
	}
	if req.After != nil {
		cursorCount++
	}
	if req.Around != nil {
		cursorCount++
	}

	if cursorCount > 1 || cursorCount < 1 {
		return nil, utils.ErrorInvalidCursorCombination
	}

	// Validate Limit
	if req.Limit <= 0 {
		return nil, utils.ErrorInvalidCursorLimit
	}

	// Check if room exists
	roomExist, err := s.IRoomRepository.DoesRoomExists(ctx, runner, &req.RoomID)
	if err != nil {
		return nil, utils.ErrorInternal
	}

	if !*roomExist {
		return nil, utils.ErrorRoomDoesntExist
	}

	// Get userId from context
	userId, _, err := auth.CurrentUserFromContext(ctx)
	if err != nil {
		return nil, utils.ErrorInvalidUserUUID
	}

	// get the room
	room, err := s.IRoomRepository.GetRoomByID(ctx, runner, &req.RoomID)
	if err != nil {
		return nil, utils.ErrorFetchingRoom
	}

	// does hall exist
	hallExist, err := s.IHallRepository.DoesHallExist(ctx, runner, room.HallId)
	if err != nil {
		return nil, utils.ErrorFetchingHall
	}

	if !*hallExist {
		return nil, utils.ErrorHallDoesntExist
	}

	// does user belong to hall and room
	userBelongsToHall, err := s.IHallRepository.IsUserHallMember(ctx, runner, room.HallId, *userId)
	if err != nil {
		return nil, utils.ErrorFetchingHall
	}

	if !*userBelongsToHall {
		return nil, utils.ErrorUserDoesntBelongHall
	}

	// if room = private , check if user belongs to the room
	if room.IsPrivate {
		userBelongsToRoom, err := s.IRoomRepository.IsUserRoomMember(ctx, runner, &room.ID, userId)

		if err != nil {
			return nil, utils.ErrorFetchingRoom
		}

		if !*userBelongsToRoom {
			return nil, utils.ErrorUserDoesntBelongRoom
		}
	}

	// call repository
	messages, err := s.IMessageRepository.GetMessages(ctx, runner, &dto.MessageQueryParams{
		RoomID: req.RoomID,
		Before: req.Before,
		After:  req.After,
		Around: req.Around,
		Limit:  req.Limit + 1,
	})
	if err != nil {
		return nil, err
	}

	// Check if we got +1 more message
	hasMore := len(messages) > req.Limit

	if hasMore {
		// trim messages to limit
		messages = messages[:req.Limit]
	}

	messageListResponse := &dto.MessageListResponse{
		Messages: messages,
		HasMore:  hasMore,
	}

	return messageListResponse, nil
}
