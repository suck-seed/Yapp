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
	CreateMessage(c context.Context, req *dto.CreateMessageReq) (*models.Message, error)

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

func (s *messageService) CreateMessage(c context.Context, req *dto.CreateMessageReq) (*models.Message, error) {
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
		RoomId:   req.RoomId,
		AuthorId: req.UserId,
		Content:  *normalizedContent,
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
	if *req.MentionEveryone == false && req.Mentions != nil {

		// validate if the userId acc exists or not
		for _, mentionedUserIDString := range *req.Mentions {

			// convert string into UUID

			mentionedUserID, err := uuid.Parse(mentionedUserIDString)
			if err != nil {
				return nil, utils.ErrorInternal
			}

			exists, err := s.IUserRepository.UserExists(ctx, mentionedUserID)
			if err != nil {
				return nil, utils.ErrorInternal
			}

			if exists {

				// Add to table message_mentions
				if err := s.IUserRepository.AddMessageMention(ctx, messageCRES.ID, mentionedUserID); err != nil {
					return nil, utils.ErrorWritingMentions
				}

			}

		}

	}

	// attachments check
	if *req.Attachments != nil {

		for _, currentAttachment := range *req.Attachments {

		}
	}

	print(ctx)

	return &models.Message{}, nil

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
