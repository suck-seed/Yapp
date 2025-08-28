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

type IUserService interface {
	CreateUser(c context.Context, user *dto.CreateUserReq) (*dto.CreateUserRes, error)
}

// userService : Behaves like a class, and implements IUserService's methods
type userService struct {
	repositories.IUserRepository
	timeout time.Duration
	mu      sync.RWMutex
}

// NewUserService : Constructor to return a new IUserService with all the user service methods
func NewUserService(repository repositories.IUserRepository) IUserService {
	return &userService{
		repository,
		time.Duration(2) * time.Second,
		sync.RWMutex{},
	}
}

// Methods
func (s *userService) CreateUser(c context.Context, req *dto.CreateUserReq) (*dto.CreateUserRes, error) {

	// interface that provides a way to control lifecycle, cancellation and prppaagation of requests
	ctx, cancel := context.WithTimeout(c, s.timeout)
	defer cancel()

	// generate id
	id, err := uuid.NewV7()
	if err != nil {
		return nil, err
	}

	// hash password
	password_hash, err := utils.HashPassword(req.Password)

	if err != nil {
		return &dto.CreateUserRes{}, err
	}

	//
	// a model.User{
	//
	user := &models.User{
		ID:           id,
		Username:     req.Username,
		Email:        req.Email,
		PhoneNumber:  &req.PhoneNumber,
		PasswordHash: password_hash,
		DisplayName:  &req.DisplayName,
	}

	// calling the repo
	r, err := s.IUserRepository.CreateUser(ctx, user)
	if err != nil {
		return &dto.CreateUserRes{}, err
	}

	// create a response
	res := &dto.CreateUserRes{
		ID:       r.ID.String(),
		Username: r.Username,
	}

	return res, nil

}
