package repositories

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/suck-seed/yapp/internal/dto"
	"github.com/suck-seed/yapp/internal/models"
)

type IUserRepository interface {
	CreateUser(ctx context.Context, user *models.User) (*models.User, error)
	UpdateUserById(ctx context.Context, userID uuid.UUID, req *dto.UpdateUserMeReq) (*models.User, error)

	GetUserWithPasswordHashByEmail(ctx context.Context, email string) (*models.User, error)
	GetUserByEmail(ctx context.Context, email string) (*models.User, error)
	GetUserByUsername(ctx context.Context, username string) (*models.User, error)
	GetUserByNumber(ctx context.Context, number *string) (*models.User, error)
	GetUserById(ctx context.Context, userID uuid.UUID) (*models.User, error)

	UserExists(ctx context.Context, userID uuid.UUID) (bool, error)
	AddMessageMention(ctx context.Context, messageId uuid.UUID, userID uuid.UUID) error
}

type userRepository struct {
	db PGXTX
}

func NewUserRepository(db PGXTX) IUserRepository {

	return &userRepository{
		db: db,
	}
}

func (r *userRepository) CreateUser(ctx context.Context, user *models.User) (*models.User, error) {

	saved := &models.User{}

	query := `
  				INSERT INTO users (id, username, display_name, email, password_hash)
      			VALUES ($1, $2, $3, $4, $5)
         		RETURNING id, username, display_name, email, phone_number, avatar_url, avatar_thumbnail_url, friend_policy, created_at, updated_at

   			`

	row := r.db.QueryRow(ctx, query,
		user.ID,
		user.Username,
		user.DisplayName,
		user.Email,
		user.PasswordHash,
	)

	err := row.Scan(
		&saved.ID,
		&saved.Username,
		&saved.DisplayName,
		&saved.Email,
		&saved.PhoneNumber,
		&saved.AvatarURL,
		&saved.AvatarThumbnailURL,
		&saved.FriendPolicy,
		&saved.CreatedAt,
		&saved.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return saved, nil
}

func (r *userRepository) GetUserByUsername(ctx context.Context, username string) (*models.User, error) {

	user := &models.User{}

	query := `
				SELECT id, username, display_name, email, phone_number, avatar_url, avatar_thumbnail_url, friend_policy, created_at, updated_at, active
				FROM users
				WHERE lower(username) = lower($1)
			`

	row := r.db.QueryRow(ctx, query, username)

	err := row.Scan(
		&user.ID,
		&user.Username,
		&user.DisplayName,
		&user.Email,
		&user.PhoneNumber,
		&user.AvatarURL,
		&user.AvatarThumbnailURL,
		&user.FriendPolicy,
		&user.CreatedAt,
		&user.UpdatedAt,
		&user.Active,
	)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (r *userRepository) GetUserByNumber(ctx context.Context, number *string) (*models.User, error) {

	user := &models.User{}

	query := `
				SELECT id, username, display_name, email, phone_number, avatar_url, avatar_thumbnail_url, friend_policy, created_at, updated_at, active
				FROM users
				WHERE phone_number = $1
			`

	row := r.db.QueryRow(ctx, query, number)

	err := row.Scan(
		&user.ID,
		&user.Username,
		&user.DisplayName,
		&user.Email,
		&user.PhoneNumber,
		&user.AvatarURL,
		&user.AvatarThumbnailURL,
		&user.FriendPolicy,
		&user.CreatedAt,
		&user.UpdatedAt,
		&user.Active,
	)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (r *userRepository) GetUserWithPasswordHashByEmail(ctx context.Context, email string) (*models.User, error) {

	user := &models.User{}

	query := `
				SELECT id, username, display_name, email, password_hash, phone_number, avatar_url, avatar_thumbnail_url, friend_policy,
				created_at, updated_at, active
				FROM users
				WHERE lower(email) = lower($1)
			`

	row := r.db.QueryRow(ctx, query, email)

	err := row.Scan(
		&user.ID,
		&user.Username,
		&user.DisplayName,
		&user.Email,
		&user.PasswordHash,
		&user.PhoneNumber,
		&user.AvatarURL,
		&user.AvatarThumbnailURL,
		&user.FriendPolicy,
		&user.CreatedAt,
		&user.UpdatedAt,
		&user.Active,
	)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (r *userRepository) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {

	user := &models.User{}

	query := `
				SELECT id, username, display_name, email, phone_number, avatar_url, avatar_thumbnail_url, friend_policy,
				created_at, updated_at, active
				FROM users
				WHERE lower(email) = lower($1)
			`

	row := r.db.QueryRow(ctx, query, email)

	err := row.Scan(
		&user.ID,
		&user.Username,
		&user.DisplayName,
		&user.Email,
		&user.PhoneNumber,
		&user.AvatarURL,
		&user.AvatarThumbnailURL,
		&user.FriendPolicy,
		&user.CreatedAt,
		&user.UpdatedAt,
		&user.Active,
	)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (r *userRepository) GetUserById(ctx context.Context, userId uuid.UUID) (*models.User, error) {

	user := &models.User{}

	query := `
				SELECT id, username, display_name, email, phone_number, avatar_url, avatar_thumbnail_url, friend_policy, created_at,
				updated_at, active
				FROM users
				WHERE id = $1
			`

	row := r.db.QueryRow(ctx, query, userId)

	err := row.Scan(
		&user.ID,
		&user.Username,
		&user.DisplayName,
		&user.Email,
		&user.PhoneNumber,
		&user.AvatarURL,
		&user.AvatarThumbnailURL,
		&user.FriendPolicy,
		&user.CreatedAt,
		&user.UpdatedAt,
		&user.Active,
	)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (r *userRepository) UpdateUserById(ctx context.Context, userId uuid.UUID, req *dto.UpdateUserMeReq) (*models.User, error) {

	user := &models.User{}

	query := `
        UPDATE users
        SET display_name = $1, avatar_url = $2, avatar_thumbnail_url = $3, updated_at = $4
        WHERE id = $5
        RETURNING username, display_name, email, avatar_url, avatar_thumbnail_url, active
    `
	err := r.db.QueryRow(ctx, query, req.DisplayName, req.AvatarURL, req.AvatarThumbnailURL, time.Now(), userId).Scan(
		&user.Username,
		&user.DisplayName,
		&user.Email,
		&user.AvatarURL,
		&user.AvatarThumbnailURL,
		&user.Active,
	)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (r *userRepository) UserExists(ctx context.Context, userId uuid.UUID) (bool, error) {

	query := `

		SELECT EXISTS (SELECT 1 FROM users WHERE id = $1)

	`

	var exists bool

	err := r.db.QueryRow(ctx, query, userId).Scan(&exists)
	if err != nil {
		return false, err
	}

	return true, nil
}

func (r *userRepository) AddMessageMention(ctx context.Context, messageId uuid.UUID, userID uuid.UUID) error {

	query := `
  				INSERT INTO message_mentions (message_id,user_id)
      			VALUES ($1, $2)
        		ON CONFLICT (message_id, user_id) DO NOTHING

   			`

	_, err := r.db.Exec(ctx, query)

	return err
}
