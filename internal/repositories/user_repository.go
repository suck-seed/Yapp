package repositories

import (
	"context"

	"github.com/suck-seed/yapp/internal/models"
)

type IUserRepository interface {
	CreateUser(ctx context.Context, user *models.User) (*models.User, error)
}

type userRepository struct {
	db PGXTX
}

func NewUserRepository(db PGXTX) IUserRepository {

	return &userRepository{
		db: db,
	}

}

func (userRepository *userRepository) CreateUser(ctx context.Context, user *models.User) (*models.User, error) {

	query := "INSERT INTO users(id, username,display_name,email,password_hash,phone_number) VALUES($1,$2,$3,$4,$5,$6)"

	// userRepository.db.QueryRow(ctx, query, args ...any)

}
