package user

import "github.com/suck-seed/yapp/internal/models"

type IUserService interface {

	//TODO implement functions that has to be implemented and accessed
	GetUserByID(id string) (*models.User, error)
}
