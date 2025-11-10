package services

import (
	"context"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/suck-seed/yapp/internal/dto"
	"github.com/suck-seed/yapp/internal/models"
	"github.com/suck-seed/yapp/internal/repositories"
	"github.com/suck-seed/yapp/internal/utils"
)

type IMessageService interface {
	CreateMessage(c context.Context, req *dto.CreateMessageReq) (*dto.CreateMessageRes, error)
}

type messageService struct {
	repositories.IRoomRepository
	repositories.IMessageRepository
	repositories.IUserRepository
	timeout time.Duration
	mu      sync.RWMutex
}

func NewMessageService(roomRepo repositories.IRoomRepository, messageRepo repositories.IMessageRepository, userRepo repositories.IUserRepository) IMessageService {
	return &messageService{
		roomRepo,
		messageRepo,
		userRepo,
		time.Duration(2) * time.Second,
		sync.RWMutex{},
	}
}

// TODO Properly fill in these methods

func (s *messageService) CreateMessage(c context.Context, req *dto.CreateMessageReq) (*dto.CreateMessageRes, error) {
	ctx, cancel := context.WithTimeout(c, s.timeout)
	defer cancel()

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
	messageCRES, err := s.IMessageRepository.CreateMessage(ctx, message)
	if err != nil {
		return nil, utils.ErrorWritingMessage
	}

	// else , create entry on message_mentions
	var mentions []dto.UserBasic
	if *req.MentionEveryone == false && req.Mentions != nil {

		// validate if the userId acc exists or not
		for _, mentionedUserID := range *req.Mentions {

			// convert string into UUID

			exists, err := s.IUserRepository.DoesUserExists(ctx, mentionedUserID)
			if err != nil {
				return nil, utils.ErrorInternal
			}

			if exists {

				// Add to table message_mentions
				if err := s.IMessageRepository.AddMessageMention(ctx, messageCRES.ID, mentionedUserID); err != nil {
					return nil, utils.ErrorWritingMentions
				}

				// get userBasic information for generation
				userCRES, err := s.IUserRepository.GetUserById(ctx, mentionedUserID)
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
	var attachments []dto.AttachmentResponseMinimal
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
			attachmentCRES, err := s.IMessageRepository.AddAttachment(ctx, &models.Attachment{
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
			attachments = append(attachments, dto.AttachmentResponseMinimal{
				ID:       attachmentCRES.ID,
				URL:      attachmentCRES.URL,
				FileName: attachmentCRES.FileName,
				FileType: attachmentCRES.FileType,
			})
		}

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
