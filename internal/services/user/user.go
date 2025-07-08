package user

import (
	"github.com/suck-seed/yapp/internal/dto"
	"github.com/suck-seed/yapp/internal/models"
)

type IUserService interface {

	//TODO implement functions that has to be implemented and accessed
	RegisterUser(user dto.UserSignup) (string, error)
	GetUserByID(id string) (*models.User, error)
}
