package services

import (
	"sync"
	"time"

	"github.com/suck-seed/yapp/internal/repositories"
)

type IFloorService interface {
}

type floorService struct {
	repositories.IHallRepository
	repositories.IFloorRepository
	timeout time.Duration
	mu      sync.RWMutex
}

func NewFloorService(hallRepo repositories.IHallRepository, floorRepo repositories.IFloorRepository) IFloorService {
	return &floorService{
		hallRepo,
		floorRepo,
		time.Duration(2) * time.Second,
		sync.RWMutex{},
	}
}
