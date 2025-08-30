package repositories

import (
	"context"
	"fmt"

	"github.com/suck-seed/yapp/internal/models"
)

type IUserRepository interface {
	CreateUser(ctx context.Context, user *models.User) (*models.User, error)
	GetUserByEmail(ctx context.Context, email string) (*models.User, error)
	GetUserByUsername(ctx context.Context, username string) (*models.User, error)
	GetUserByNumber(ctx context.Context, number *string) (*models.User, error)
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

	query := `
  				INSERT INTO users (id, username, display_name, email, password_hash, phone_number)
      			VALUES ($1, $2, $3, $4, $5, $6)
   			`

	tag, err := userRepository.db.Exec(ctx, query,
		user.ID,
		user.Username,
		user.DisplayName,
		user.Email,
		user.PasswordHash,
		user.PhoneNumber,
	)

	if err != nil {
		return nil, err
	}

	fmt.Println("Rows effected : ", tag.RowsAffected())

	return user, nil
}

func (userRepository *userRepository) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {

	// create a user model
	user := &models.User{}

	query := `
				SELECT id,username,display_name, email, phone_number, avatar_url, password_hash
				FROM users
				WHERE email = $1
			`

	row := userRepository.db.QueryRow(ctx, query, email)

	err := row.Scan(&user.ID, &user.Username, &user.DisplayName, &user.Email, &user.PhoneNumber, &user.AvatarURL, &user.PasswordHash)
	if err != nil {
		return nil, err
	}

	return user, nil

}

func (userRepository *userRepository) GetUserByUsername(ctx context.Context, username string) (*models.User, error) {

	// create a user model
	user := &models.User{}

	query := `
				SELECT id,username,display_name, email, phone_number, avatar_url, password_hash
				FROM users
				WHERE username = $1
			`

	row := userRepository.db.QueryRow(ctx, query, username)

	err := row.Scan(&user.ID, &user.Username, &user.DisplayName, &user.Email, &user.PhoneNumber, &user.AvatarURL, &user.PasswordHash)
	if err != nil {
		return nil, err
	}

	return user, nil

}

func (userRepository *userRepository) GetUserByNumber(ctx context.Context, number *string) (*models.User, error) {

	// create a user model
	user := &models.User{}

	query := `
				SELECT id,username,display_name, email, phone_number, avatar_url, password_hash
				FROM users
				WHERE phone_number = $1
			`

	row := userRepository.db.QueryRow(ctx, query, number)

	err := row.Scan(&user.ID, &user.Username, &user.DisplayName, &user.Email, &user.PhoneNumber, &user.AvatarURL, &user.PasswordHash)
	if err != nil {
		return nil, err
	}

	return user, nil

}
