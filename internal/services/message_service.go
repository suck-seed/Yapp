package services

import (
	"context"
	"errors"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
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
	FetchMessages(c context.Context, userInfo *auth.UserInfo, roomID uuid.UUID, params *dto.FetchMessagesQuery) (*dto.MessageListResponse, error)
	GetMessage(c context.Context, userInfo *auth.UserInfo, roomID uuid.UUID, messageID uuid.UUID) (*dto.MessageDetailed, error)
	UpdateMessage(c context.Context, userInfo *auth.UserInfo, roomID uuid.UUID, messageID uuid.UUID, req *dto.UpdateMessageReq) (*dto.MessageDetailed, error)
	DeleteMessage(c context.Context, userInfo *auth.UserInfo, roomID uuid.UUID, messageID uuid.UUID) error
	AddReaction(c context.Context, userInfo *auth.UserInfo, roomID uuid.UUID, messageID uuid.UUID, emoji string) (*dto.ReactionRes, error)
	RemoveReaction(c context.Context, userInfo *auth.UserInfo, roomID uuid.UUID, messageID uuid.UUID, emoji string) error
}

type messageService struct {
	repositories.IHallRepository
	repositories.IRoomRepository
	repositories.IMessageRepository
	repositories.IUserRepository

	IPermissionCheckerService

	pool    *pgxpool.Pool
	timeout time.Duration
	mu      sync.RWMutex
}

func NewMessageService(
	hallRepo repositories.IHallRepository,
	roomRepo repositories.IRoomRepository,
	messageRepo repositories.IMessageRepository,
	userRepo repositories.IUserRepository,
	permissionChecker IPermissionCheckerService,
	pool *pgxpool.Pool,
) IMessageService {
	return &messageService{
		hallRepo,
		roomRepo,
		messageRepo,
		userRepo,
		permissionChecker,
		pool,
		time.Duration(2) * time.Second,
		sync.RWMutex{},
	}
}

// ── helpers ───────────────────────────────────────────────────────────────────

// resolveRoom fetches the room and verifies the user is a hall member.
// Returns the room so callers can access room.HallID and room.IsPrivate.
func (s *messageService) resolveRoom(ctx context.Context, runner database.DBRunner, roomID uuid.UUID, userID uuid.UUID) (*models.Room, error) {
	room, err := s.IRoomRepository.GetRoomByID(ctx, runner, roomID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, utils.ErrorRoomNotFound
		}
		if isDeadline(err) {
			return nil, utils.ErrorRequestTimeout
		}
		return nil, utils.ErrorFetchingRoom
	}

	ok, err := s.IHallRepository.IsUserHallMember(ctx, runner, room.HallID, userID)
	if err != nil {
		return nil, utils.ErrorInternal
	}
	if !ok {
		return nil, utils.ErrorUserDoesntBelongHall
	}

	return room, nil
}

// resolveRoomWithPrivateCheck also enforces room membership for private rooms.
func (s *messageService) resolveRoomWithPrivateCheck(ctx context.Context, runner database.DBRunner, roomID uuid.UUID, userID uuid.UUID) (*models.Room, error) {
	room, err := s.resolveRoom(ctx, runner, roomID, userID)
	if err != nil {
		return nil, err
	}

	if room.IsPrivate {
		inRoom, err := s.IRoomRepository.IsUserRoomMember(ctx, runner, room.ID, userID)
		if err != nil {
			return nil, utils.ErrorInternal
		}
		if !inRoom {
			return nil, utils.ErrorUserDoesntBelongRoom
		}
	}

	return room, nil
}

// ── CreateMessage ─────────────────────────────────────────────────────────────
// Called internally by the WebSocket hub, not directly from HTTP.

func (s *messageService) CreateMessage(c context.Context, req *dto.CreateMessageReq) (*dto.CreateMessageRes, error) {
	ctx, cancel := context.WithTimeout(c, s.timeout)
	defer cancel()

	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return nil, utils.ErrorInternal
	}
	runner := database.NewTxWrapper(tx)
	defer runner.Rollback(ctx)

	normalizedContent := utils.SanitizeMessageContent(req.Content)

	messageID, err := uuid.NewV7()
	if err != nil {
		return nil, utils.ErrorInternal
	}

	message := &models.Message{
		ID:       messageID,
		RoomID:   req.RoomID,
		AuthorID: req.AuthorID,
		Content:  normalizedContent,
		SentAt:   req.SentAt,
	}

	message.MentionEveryone = false
	if req.MentionEveryone != nil && *req.MentionEveryone {
		message.MentionEveryone = true
	}

	messageCRES, err := s.IMessageRepository.CreateMessage(ctx, runner, message)
	if err != nil {
		return nil, utils.ErrorWritingMessage
	}

	var mentions []dto.UserBasic
	if req.MentionEveryone != nil && !*req.MentionEveryone && req.Mentions != nil {
		for _, mentionedUserID := range *req.Mentions {
			exists, err := s.IUserRepository.DoesUserExists(ctx, runner, mentionedUserID)
			if err != nil {
				return nil, utils.ErrorInternal
			}
			if !exists {
				continue
			}

			if err := s.IMessageRepository.AddMessageMention(ctx, runner, messageCRES.ID, mentionedUserID); err != nil {
				return nil, utils.ErrorWritingMentions
			}

			userCRES, err := s.IUserRepository.GetUserById(ctx, runner, mentionedUserID)
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

	var attachments []models.Attachment
	if req.Attachments != nil && len(*req.Attachments) > 0 {
		for _, currentAttachment := range *req.Attachments {
			canonFileName, err := utils.ValidateFileName(currentAttachment.FileName)
			if err != nil {
				return nil, err
			}

			validatedFileType, err := utils.ValidateFileType(currentAttachment.FileType, currentAttachment.URL)
			if err != nil {
				return nil, err
			}

			if currentAttachment.FileSize != nil && *currentAttachment.FileSize >= utils.FileSize {
				return nil, utils.ErrorLargeFileSize
			}

			attachmentID, err := uuid.NewV7()
			if err != nil {
				return nil, utils.ErrorInternal
			}

			attachmentCRES, err := s.IMessageRepository.AddAttachment(ctx, runner, &models.Attachment{
				ID:        attachmentID,
				MessageID: messageID,
				FileName:  canonFileName,
				URL:       currentAttachment.URL,
				FileType:  validatedFileType,
				FileSize:  currentAttachment.FileSize,
			})
			if err != nil {
				return nil, utils.ErrorCreatingAttachment
			}

			attachments = append(attachments, *attachmentCRES)
		}
	}

	if err := runner.Commit(ctx); err != nil {
		return nil, utils.ErrorInternal
	}

	return &dto.CreateMessageRes{
		ID:               messageCRES.ID,
		RoomID:           messageCRES.RoomID,
		AuthorID:         messageCRES.AuthorID,
		Content:          messageCRES.Content,
		SentAt:           messageCRES.SentAt,
		MentionsEveryone: messageCRES.MentionEveryone,
		Mentions:         mentions,
		Attachments:      attachments,
		CreatedAt:        messageCRES.CreatedAt,
		EditedAt:         messageCRES.EditedAt,
		DeletedAt:        messageCRES.DeletedAt,
		UpdatedAt:        messageCRES.UpdatedAt,
	}, nil
}

// ── FetchMessages ─────────────────────────────────────────────────────────────

func (s *messageService) FetchMessages(c context.Context, userInfo *auth.UserInfo, roomID uuid.UUID, params *dto.FetchMessagesQuery) (*dto.MessageListResponse, error) {
	ctx, cancel := context.WithTimeout(c, s.timeout)
	defer cancel()

	conn, err := s.pool.Acquire(ctx)
	if err != nil {
		return nil, utils.ErrorInternal
	}
	defer conn.Release()
	runner := database.NewConnWrapper(conn)

	// Validate at most one cursor is used.
	// Zero cursor is valid for the initial page fetch.
	cursorCount := 0
	if params.Before != nil {
		cursorCount++
	}
	if params.After != nil {
		cursorCount++
	}
	if params.Around != nil {
		cursorCount++
	}
	if cursorCount > 1 {
		return nil, utils.ErrorInvalidCursorCombination
	}

	if _, err := s.resolveRoomWithPrivateCheck(ctx, runner, roomID, userInfo.ID); err != nil {
		return nil, err
	}

	messages, err := s.IMessageRepository.GetMessages(ctx, runner, &dto.MessageQueryParams{
		RoomID: roomID,
		Before: params.Before,
		After:  params.After,
		Around: params.Around,
		Limit:  params.Limit + 1, // fetch one extra to determine hasMore
	})
	if err != nil {
		return nil, err
	}

	hasMore := len(messages) > params.Limit
	if hasMore {
		messages = messages[:params.Limit]
	}

	return &dto.MessageListResponse{
		Messages: messages,
		HasMore:  hasMore,
	}, nil
}

// ── GetMessage ────────────────────────────────────────────────────────────────

func (s *messageService) GetMessage(c context.Context, userInfo *auth.UserInfo, roomID uuid.UUID, messageID uuid.UUID) (*dto.MessageDetailed, error) {
	ctx, cancel := context.WithTimeout(c, s.timeout)
	defer cancel()

	conn, err := s.pool.Acquire(ctx)
	if err != nil {
		return nil, utils.ErrorInternal
	}
	defer conn.Release()
	runner := database.NewConnWrapper(conn)

	if _, err := s.resolveRoomWithPrivateCheck(ctx, runner, roomID, userInfo.ID); err != nil {
		return nil, err
	}

	message, err := s.IMessageRepository.GetMessageDetailed(ctx, runner, messageID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, utils.ErrorMessageNotFound
		}
		if isDeadline(err) {
			return nil, utils.ErrorRequestTimeout
		}
		return nil, utils.ErrorFetchingMessages
	}

	// Ensure message belongs to this room
	if message.RoomID != roomID {
		return nil, utils.ErrorMessageNotFound
	}

	return message, nil
}

// ── UpdateMessage ─────────────────────────────────────────────────────────────
// Only the author can update their own message.

func (s *messageService) UpdateMessage(c context.Context, userInfo *auth.UserInfo, roomID uuid.UUID, messageID uuid.UUID, req *dto.UpdateMessageReq) (*dto.MessageDetailed, error) {
	ctx, cancel := context.WithTimeout(c, s.timeout)
	defer cancel()

	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return nil, utils.ErrorInternal
	}
	runner := database.NewTxWrapper(tx)
	defer runner.Rollback(ctx)

	if _, err := s.resolveRoomWithPrivateCheck(ctx, runner, roomID, userInfo.ID); err != nil {
		return nil, err
	}

	message, err := s.IMessageRepository.GetMessageByID(ctx, runner, messageID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, utils.ErrorMessageNotFound
		}
		if isDeadline(err) {
			return nil, utils.ErrorRequestTimeout
		}
		return nil, utils.ErrorFetchingMessages
	}

	if message.RoomID != roomID {
		return nil, utils.ErrorMessageNotFound
	}

	if message.AuthorID != userInfo.ID {
		return nil, utils.ErrorForbidden
	}

	content := utils.SanitizeMessageContent(&req.Content)
	if content == nil || *content == "" {
		return nil, utils.ErrorInvalidInput
	}

	if _, err := s.IMessageRepository.UpdateMessageContent(ctx, runner, messageID, *content); err != nil {
		if isDeadline(err) {
			return nil, utils.ErrorRequestTimeout
		}
		return nil, utils.ErrorInternal
	}

	if err := runner.Commit(ctx); err != nil {
		return nil, utils.ErrorInternal
	}

	// Fetch the full updated message with author/attachments/reactions
	conn, err := s.pool.Acquire(ctx)
	if err != nil {
		return nil, utils.ErrorInternal
	}
	defer conn.Release()

	updated, err := s.IMessageRepository.GetMessageDetailed(ctx, database.NewConnWrapper(conn), messageID)
	if err != nil {
		return nil, utils.ErrorFetchingMessages
	}

	return updated, nil
}

// ── DeleteMessage ─────────────────────────────────────────────────────────────
// Author can delete their own message.
// Anyone with manage_servers (admin, owner, manage_servers permission) can also delete.

func (s *messageService) DeleteMessage(c context.Context, userInfo *auth.UserInfo, roomID uuid.UUID, messageID uuid.UUID) error {
	ctx, cancel := context.WithTimeout(c, s.timeout)
	defer cancel()

	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return utils.ErrorInternal
	}
	runner := database.NewTxWrapper(tx)
	defer runner.Rollback(ctx)

	room, err := s.resolveRoomWithPrivateCheck(ctx, runner, roomID, userInfo.ID)
	if err != nil {
		return err
	}

	message, err := s.IMessageRepository.GetMessageByID(ctx, runner, messageID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return utils.ErrorMessageNotFound
		}
		if isDeadline(err) {
			return utils.ErrorRequestTimeout
		}
		return utils.ErrorFetchingMessages
	}

	if message.RoomID != roomID {
		return utils.ErrorMessageNotFound
	}

	// Allow if author, otherwise require manage_servers
	if message.AuthorID != userInfo.ID {
		ok, err := s.IPermissionCheckerService.CanManageServers(ctx, runner, userInfo.ID, room.HallID)
		if err != nil {
			return err
		}
		if !ok {
			return utils.ErrorForbidden
		}
	}

	if err := s.IMessageRepository.SoftDeleteMessage(ctx, runner, messageID); err != nil {
		if isDeadline(err) {
			return utils.ErrorRequestTimeout
		}
		return utils.ErrorInternal
	}

	return runner.Commit(ctx)
}

// ── AddReaction ───────────────────────────────────────────────────────────────
// Any hall member can react. Duplicate reactions (same user+emoji) are silently ignored.

func (s *messageService) AddReaction(c context.Context, userInfo *auth.UserInfo, roomID uuid.UUID, messageID uuid.UUID, emoji string) (*dto.ReactionRes, error) {
	ctx, cancel := context.WithTimeout(c, s.timeout)
	defer cancel()

	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return nil, utils.ErrorInternal
	}
	runner := database.NewTxWrapper(tx)
	defer runner.Rollback(ctx)

	if _, err := s.resolveRoomWithPrivateCheck(ctx, runner, roomID, userInfo.ID); err != nil {
		return nil, err
	}

	// Verify message exists and belongs to this room
	message, err := s.IMessageRepository.GetMessageByID(ctx, runner, messageID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, utils.ErrorMessageNotFound
		}
		if isDeadline(err) {
			return nil, utils.ErrorRequestTimeout
		}
		return nil, utils.ErrorFetchingMessages
	}
	if message.RoomID != roomID {
		return nil, utils.ErrorMessageNotFound
	}

	reactionID, err := uuid.NewV7()
	if err != nil {
		return nil, utils.ErrorInternal
	}

	// ON CONFLICT DO NOTHING — duplicate reactions are silently ignored
	if err := s.IMessageRepository.AddReaction(ctx, runner, reactionID, messageID, userInfo.ID, emoji); err != nil {
		if isDeadline(err) {
			return nil, utils.ErrorRequestTimeout
		}
		return nil, utils.ErrorInternal
	}

	if err := runner.Commit(ctx); err != nil {
		return nil, utils.ErrorInternal
	}

	return &dto.ReactionRes{
		MessageID: messageID,
		UserID:    userInfo.ID,
		Emoji:     emoji,
	}, nil
}

// ── RemoveReaction ────────────────────────────────────────────────────────────

func (s *messageService) RemoveReaction(c context.Context, userInfo *auth.UserInfo, roomID uuid.UUID, messageID uuid.UUID, emoji string) error {
	ctx, cancel := context.WithTimeout(c, s.timeout)
	defer cancel()

	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return utils.ErrorInternal
	}
	runner := database.NewTxWrapper(tx)
	defer runner.Rollback(ctx)

	if _, err := s.resolveRoomWithPrivateCheck(ctx, runner, roomID, userInfo.ID); err != nil {
		return err
	}

	message, err := s.IMessageRepository.GetMessageByID(ctx, runner, messageID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return utils.ErrorMessageNotFound
		}
		if isDeadline(err) {
			return utils.ErrorRequestTimeout
		}
		return utils.ErrorFetchingMessages
	}
	if message.RoomID != roomID {
		return utils.ErrorMessageNotFound
	}

	deleted, err := s.IMessageRepository.RemoveReaction(ctx, runner, messageID, userInfo.ID, emoji)
	if err != nil {
		if isDeadline(err) {
			return utils.ErrorRequestTimeout
		}
		return utils.ErrorInternal
	}
	if !deleted {
		// Reaction didn't exist — treat as 404
		return utils.ErrorReactionNotFound
	}

	return runner.Commit(ctx)
}
