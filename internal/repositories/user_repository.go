package repositories

import (
	"context"

	"github.com/google/uuid"
	"github.com/suck-seed/yapp/internal/dto"
	"github.com/suck-seed/yapp/internal/models"
)

type IUserRepository interface {
	CreateUser(ctx context.Context, user *models.User) (*models.User, error)
	UpdateUserByUsername(ctx context.Context, username string, req *dto.UpdateProfileReq) (*models.User, error)

	GetUserByEmail(ctx context.Context, email string) (*models.User, error)
	GetUserByUsername(ctx context.Context, username string) (*models.User, error)
	GetUserByNumber(ctx context.Context, number *string) (*models.User, error)
	GetUserById(ctx context.Context, userId uuid.UUID) (*models.User, error)
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
  				INSERT INTO users (id, username, display_name, email, password_hash)
      			VALUES ($1, $2, $3, $4, $5)
         		RETURNING id, username, display_name, email, phone_number, avatar_url, friend_policy, created_at, updated_at

   			`

	row := userRepository.db.QueryRow(ctx, query,
		user.ID,
		user.Username,
		user.DisplayName,
		user.Email,
		user.PasswordHash,
	)

	saved := &models.User{}

	err := row.Scan(
		&saved.ID,
		&saved.Username,
		&saved.DisplayName,
		&saved.Email,
		&saved.PhoneNumber,
		&saved.AvatarURL,
		&saved.FriendPolicy,
		&saved.CreatedAt,
		&saved.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return saved, nil
}

func (userRepository *userRepository) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {

	user := &models.User{}

	query := `
				SELECT id, username, display_name, email, phone_number, avatar_url, password_hash
				FROM users
				WHERE lower(email) = lower($1)
			`

	row := userRepository.db.QueryRow(ctx, query, email)

	err := row.Scan(
		&user.ID,
		&user.Username,
		&user.DisplayName,
		&user.Email,
		&user.PhoneNumber,
		&user.AvatarURL,
		&user.PasswordHash,
	)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (userRepository *userRepository) GetUserById(ctx context.Context, userId uuid.UUID) (*models.User, error) {

	query := `
				SELECT username, display_name, email, avatar_url, active
				FROM users
				WHERE id = $1
			`

	row := userRepository.db.QueryRow(ctx, query, userId)

	user := &models.User{}
	err := row.Scan(
		&user.Username,
		&user.DisplayName,
		&user.Email,
		&user.AvatarURL,
		&user.Active,
	)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (userRepository *userRepository) GetUserByUsername(ctx context.Context, username string) (*models.User, error) {

	user := &models.User{}

	query := `
				SELECT id, username, display_name, email, phone_number, avatar_url, password_hash
				FROM users
				WHERE lower(username) = lower($1)
			`

	row := userRepository.db.QueryRow(ctx, query, username)

	err := row.Scan(
		&user.ID,
		&user.Username,
		&user.DisplayName,
		&user.Email,
		&user.PhoneNumber,
		&user.AvatarURL,
		&user.PasswordHash,
	)
	if err != nil {
		return nil, err
	}

	return user, nil

}

func (userRepository *userRepository) GetUserByNumber(ctx context.Context, number *string) (*models.User, error) {

	user := &models.User{}

	query := `
				SELECT id, username, display_name, email, phone_number, avatar_url, password_hash
				FROM users
				WHERE phone_number = $1
			`

	row := userRepository.db.QueryRow(ctx, query, number)

	err := row.Scan(
		&user.ID,
		&user.Username,
		&user.DisplayName,
		&user.Email,
		&user.PhoneNumber,
		&user.AvatarURL,
		&user.PasswordHash,
	)
	if err != nil {
		return nil, err
	}

	return user, nil

}

func (userRepository *userRepository) UpdateUserByUsername(ctx context.Context, username string, req *dto.UpdateProfileReq) (*models.User, error) {
	user := &models.User{}
	query := `
        UPDATE users
        SET display_name = $1, avatar_url = $2
        WHERE username = $3
        RETURNING username, display_name, email, avatar_url, active
    `
	err := userRepository.db.QueryRow(ctx, query, req.DisplayName, req.AvatarURL, username).Scan(
		&user.Username,
		&user.DisplayName,
		&user.Email,
		&user.AvatarURL,
		&user.Active,
	)
	if err != nil {
		return nil, err
	}
	return user, nil
}
