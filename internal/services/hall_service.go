package services

import (
	"sync"
	"time"

	"github.com/suck-seed/yapp/internal/repositories"
)

type IHallService interface {
}

type hallService struct {
	repositories.IHallRepository
	timeout time.Duration
	mu      sync.RWMutex
}

func NewHallService(hallRepo repositories.IHallRepository) IHallService {
	return &hallService{
		hallRepo,
		time.Duration(2) * time.Second,
		sync.RWMutex{},
	}
}
