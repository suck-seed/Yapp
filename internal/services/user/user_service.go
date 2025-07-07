package user

import (
	"errors"
	"fmt"
	"sync"

	"github.com/suck-seed/yapp/internal/dto"
	"github.com/suck-seed/yapp/internal/models"
)

// ! INTERFACE
func NewUserService() IUserService {
	return &userService{
		users: make(map[string]*models.User),
	}
}

// ! CLASS
type userService struct {
	mu    sync.RWMutex
	users map[string]*models.User
}

// ! CLASS METHODS

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
