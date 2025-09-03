package services

import (
	"context"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/suck-seed/yapp/internal/repositories"
)

type IHallService interface {
	IsMember(c context.Context, hallID *uuid.UUID, userId *uuid.UUID) (bool, error)
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

// ctx, cancel := context.WithTimeout(c, s.timeout)
// 	defer cancel()

func (s *hallService) IsMember(c context.Context, hallID *uuid.UUID, userId *uuid.UUID) (bool, error) {

	ctx, cancel := context.WithTimeout(c, s.timeout)
	defer cancel()

	print(ctx)

	return true, nil
}
