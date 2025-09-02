package services

import (
	"context"
	"sync"
	"time"

	"github.com/suck-seed/yapp/internal/dto"
	"github.com/suck-seed/yapp/internal/models"
	"github.com/suck-seed/yapp/internal/repositories"
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
	timeout time.Duration
	mu      sync.RWMutex
}

func NewMessageService(roomRepo repositories.IRoomRepository, messageRepo repositories.IMessageRepository) IMessageService {
	return &messageService{
		roomRepo,
		messageRepo,
		time.Duration(2) * time.Second,
		sync.RWMutex{},
	}
}

// TODO Properly fill in these methods

func (r *messageService) CreateMessage(c context.Context, req *dto.CreateMessageReq) (*models.Message, error) {

	return &models.Message{}, nil

}

func (r *messageService) GetMessageByID(c context.Context, req *dto.FetchMessageByIdReq) (*models.Message, error) {

	return &models.Message{}, nil

}

func (r *messageService) GetMessagesByRoomID(c context.Context, req *dto.FetchMessageByRoomIDReq) ([]*models.Message, error) {

	return []*models.Message{}, nil

}
func (r *messageService) GetRoomMessages(c context.Context, req *dto.FetchRoomMessageReq) ([]*models.Message, error) {

	return []*models.Message{}, nil

}

func (r *messageService) UpdateMessage(c context.Context, req *dto.UpdateMessageReq) (*models.Message, error) {

	return &models.Message{}, nil

}

func (r *messageService) DeleteMessage(c context.Context, req *dto.DeleteMessageReq) error {

	return nil

}
