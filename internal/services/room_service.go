package services

import (
	"sync"
	"time"

	"github.com/suck-seed/yapp/internal/repositories"
)

type IRoomService interface {
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
