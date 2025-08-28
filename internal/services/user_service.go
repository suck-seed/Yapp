package services

import (
	"errors"
	"fmt"
	"sync"

	"github.com/suck-seed/yapp/internal/dto"
	"github.com/suck-seed/yapp/internal/models"
)

type IUserService interface {
	RegisterUser(user dto.UserSignup) (string, error)
	GetUserByID(id string) (*models.User, error)
}

// userService : Behaves like a class, and implements IUserService's methods
type userService struct {
	mu    sync.RWMutex
	users map[string]*models.User
}

// NewUserService : Constructor to return a new IUserService with all the user service methods
func NewUserService() IUserService {
	return &userService{
		users: make(map[string]*models.User),
	}
}

// Methods
func (s *userService) RegisterUser(user dto.UserSignup) (string, error) {

	fmt.Println("User created under : ", user)

	return "User created sucesfully", nil

}

func (s *userService) GetUserByID(id string) (*models.User, error) {
	//TODO implement me
	s.mu.RLock()
	defer s.mu.RUnlock()

	u, exists := s.users[id]
	if !exists {
		return nil, errors.New("user not found")
	}
	return u, nil
}
