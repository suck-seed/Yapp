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

	GetMessageByID(c context.Context, req *dto.FetchMessageByIdReq) (*models.Message, error)

	GetMessagesByRoomID(c context.Context, req *dto.FetchMessageByRoomIDReq) ([]*models.Message, error)
	GetRoomMessages(c context.Context, req *dto.FetchRoomMessageReq) ([]*models.Message, error)

	UpdateMessage(c context.Context, req *dto.UpdateMessageReq) (*models.Message, error)
	DeleteMessage(c context.Context, req *dto.DeleteMessageReq) error
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
	var mentions []dto.MentionResponseMinimal
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

				mentions = append(mentions, dto.MentionResponseMinimal{
					ID: mentionedUserID,
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
				AttachmentID: attachmentID,
				MessageID:    messageId,
				FileName:     canonFileName,
				URL:          currentAttachment.URL,
				FileType:     validatedFileType,
				FileSize:     currentAttachment.FileSize,
			})

			if err != nil {
				return nil, utils.ErrorCreatingAttachment
			}

			//	append to attachments
			attachments = append(attachments, dto.AttachmentResponseMinimal{
				ID:       attachmentCRES.AttachmentID,
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
	}, nil

}

func (s *messageService) GetMessageByID(c context.Context, req *dto.FetchMessageByIdReq) (*models.Message, error) {
	ctx, cancel := context.WithTimeout(c, s.timeout)
	defer cancel()

	print(ctx)

	return &models.Message{}, nil

}

func (s *messageService) GetMessagesByRoomID(c context.Context, req *dto.FetchMessageByRoomIDReq) ([]*models.Message, error) {
	ctx, cancel := context.WithTimeout(c, s.timeout)
	defer cancel()

	print(ctx)

	return []*models.Message{}, nil

}
func (s *messageService) GetRoomMessages(c context.Context, req *dto.FetchRoomMessageReq) ([]*models.Message, error) {
	ctx, cancel := context.WithTimeout(c, s.timeout)
	defer cancel()

	print(ctx)

	return []*models.Message{}, nil

}

func (s *messageService) UpdateMessage(c context.Context, req *dto.UpdateMessageReq) (*models.Message, error) {
	ctx, cancel := context.WithTimeout(c, s.timeout)
	defer cancel()

	print(ctx)

	return &models.Message{}, nil

}

func (s *messageService) DeleteMessage(c context.Context, req *dto.DeleteMessageReq) error {
	ctx, cancel := context.WithTimeout(c, s.timeout)
	defer cancel()

	print(ctx)

	return nil

}
