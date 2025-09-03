package services

import (
	"context"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/suck-seed/yapp/internal/models"
	"github.com/suck-seed/yapp/internal/repositories"
)

type IRoomService interface {
	GetRoomByID(c context.Context, roomId *uuid.UUID) (*models.Room, error)

	IsMember(c context.Context, roomId *uuid.UUID, userId *uuid.UUID) (bool, error)
}

type roomService struct {
	repositories.IHallRepository
	repositories.IFloorRepository
	repositories.IRoomRepository
	timeout time.Duration
	mu      sync.RWMutex
}

func NewRoomService(hallRepo repositories.IHallRepository, floorRepo repositories.IFloorRepository, roomRepo repositories.IRoomRepository) IRoomService {
	return &roomService{
		hallRepo,
		floorRepo,
		roomRepo,
		time.Duration(2) * time.Second,
		sync.RWMutex{},
	}
}

// ctx, cancel := context.WithTimeout(c, s.timeout)
// 	defer cancel()

func (s *roomService) GetRoomByID(c context.Context, rooomId *uuid.UUID) (*models.Room, error) {

	ctx, cancel := context.WithTimeout(c, s.timeout)
	defer cancel()

	print(ctx)

	return &models.Room{}, nil

}

func (s *roomService) IsMember(c context.Context, roomId *uuid.UUID, userId *uuid.UUID) (bool, error) {

	ctx, cancel := context.WithTimeout(c, s.timeout)
	defer cancel()

	print(ctx)

	return true, nil
}
