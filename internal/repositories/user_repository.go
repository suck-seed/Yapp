package repositories

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/suck-seed/yapp/internal/database"
	"github.com/suck-seed/yapp/internal/models"
)

type IUserRepository interface {
	// ---------------- USER CORE
	CreateUser(ctx context.Context, db database.DBRunner, user *models.User) (*models.User, error)
	GetUserWithPasswordHashByEmail(ctx context.Context, db database.DBRunner, email string) (*models.User, error)
	GetUserByEmail(ctx context.Context, db database.DBRunner, email string) (*models.User, error)
	GetUserByUsername(ctx context.Context, db database.DBRunner, username string) (*models.User, error)
	GetUserByNumber(ctx context.Context, db database.DBRunner, number string) (*models.User, error)
	GetUserById(ctx context.Context, db database.DBRunner, userID uuid.UUID) (*models.User, error)
	DoesUserExists(ctx context.Context, db database.DBRunner, userID uuid.UUID) (bool, error)

	UpdateUserById(ctx context.Context, db database.DBRunner, userID uuid.UUID, fields map[string]any) (*models.User, error)
	UpdateUsername(ctx context.Context, db database.DBRunner, userID uuid.UUID, username string) (*models.User, error)
	UpdateEmail(ctx context.Context, db database.DBRunner, userID uuid.UUID, email string) (*models.User, error)
	DeleteUserById(ctx context.Context, db database.DBRunner, userID uuid.UUID) error

	// ---------------- FRIEND REQUESTS
	CreateFriendRequest(ctx context.Context, db database.DBRunner, req *models.FriendRequest) (*models.FriendRequest, error)
	GetFriendRequestByID(ctx context.Context, db database.DBRunner, requestID uuid.UUID) (*models.FriendRequest, error)
	GetFriendRequestByUsers(ctx context.Context, db database.DBRunner, senderID uuid.UUID, receiverID uuid.UUID) (*models.FriendRequest, error)
	DeleteFriendRequestByID(ctx context.Context, db database.DBRunner, requestID uuid.UUID) (*models.FriendRequest, error)
	DeleteFriendRequestByUsers(ctx context.Context, db database.DBRunner, senderID uuid.UUID, receiverID uuid.UUID) error
	ListIncomingFriendRequests(ctx context.Context, db database.DBRunner, userID uuid.UUID) ([]*models.FriendRequest, error)
	ListOutgoingFriendRequests(ctx context.Context, db database.DBRunner, userID uuid.UUID) ([]*models.FriendRequest, error)
	DoesFriendRequestExist(ctx context.Context, db database.DBRunner, senderID uuid.UUID, receiverID uuid.UUID) (bool, error)

	// ---------------- FRIENDS
	CreateFriendship(ctx context.Context, db database.DBRunner, userID1 uuid.UUID, userID2 uuid.UUID) (*models.Friend, error)
	DeleteFriendship(ctx context.Context, db database.DBRunner, userID1 uuid.UUID, userID2 uuid.UUID) error
	AreFriends(ctx context.Context, db database.DBRunner, userID1 uuid.UUID, userID2 uuid.UUID) (bool, error)
	ListFriends(ctx context.Context, db database.DBRunner, userID uuid.UUID) ([]*models.User, error)
	CountFriends(ctx context.Context, db database.DBRunner, userID uuid.UUID) (int, error)
	GetMutualFriends(ctx context.Context, db database.DBRunner, currentUserID uuid.UUID, targetUserID uuid.UUID) ([]*models.User, error)
	CountMutualFriends(ctx context.Context, db database.DBRunner, currentUserID uuid.UUID, targetUserID uuid.UUID) (int, error)

	// ---------------- APP LINKS
	UpsertAppLink(ctx context.Context, db database.DBRunner, link *models.UserAppLink) (*models.UserAppLink, error)
	DeleteAppLink(ctx context.Context, db database.DBRunner, userID uuid.UUID, provider models.AppProvider) error
	GetUserAppLinks(ctx context.Context, db database.DBRunner, userID uuid.UUID, onlyVisible bool) ([]*models.UserAppLink, error)
	GetUserAppLinkByProvider(ctx context.Context, db database.DBRunner, userID uuid.UUID, provider models.AppProvider) (*models.UserAppLink, error)
}

type userRepository struct{}

func NewUserRepository() IUserRepository {
	return &userRepository{}
}

func normalizedFriendPair(a, b uuid.UUID) (uuid.UUID, uuid.UUID) {
	if strings.Compare(a.String(), b.String()) < 0 {
		return a, b
	}
	return b, a
}

func scanUser(row pgx.Row) (*models.User, error) {
	user := &models.User{}
	err := row.Scan(
		&user.ID,
		&user.Username,
		&user.DisplayName,
		&user.Email,
		&user.PasswordHash,
		&user.Description,
		&user.PhoneNumber,
		&user.AvatarURL,
		&user.AvatarThumbnailURL,
		&user.FriendPolicy,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (r *userRepository) CreateUser(ctx context.Context, db database.DBRunner, user *models.User) (*models.User, error) {
	query := `
		INSERT INTO users (id, username, display_name, email, password_hash)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, username, display_name, email, password_hash, description, phone_number,
		          avatar_url, avatar_thumbnail_url, friend_policy, created_at, updated_at
	`

	return scanUser(db.QueryRow(ctx, query,
		user.ID,
		user.Username,
		user.DisplayName,
		user.Email,
		user.PasswordHash,
	))
}

func (r *userRepository) GetUserWithPasswordHashByEmail(ctx context.Context, db database.DBRunner, email string) (*models.User, error) {
	query := `
		SELECT id, username, display_name, email, password_hash, description, phone_number,
		       avatar_url, avatar_thumbnail_url, friend_policy, created_at, updated_at
		FROM users
		WHERE lower(email) = lower($1)
	`
	return scanUser(db.QueryRow(ctx, query, email))
}

func (r *userRepository) GetUserByEmail(ctx context.Context, db database.DBRunner, email string) (*models.User, error) {
	query := `
		SELECT id, username, display_name, email, password_hash, description, phone_number,
		       avatar_url, avatar_thumbnail_url, friend_policy, created_at, updated_at
		FROM users
		WHERE lower(email) = lower($1)
	`
	return scanUser(db.QueryRow(ctx, query, email))
}

func (r *userRepository) GetUserByUsername(ctx context.Context, db database.DBRunner, username string) (*models.User, error) {
	query := `
		SELECT id, username, display_name, email, password_hash, description, phone_number,
		       avatar_url, avatar_thumbnail_url, friend_policy, created_at, updated_at
		FROM users
		WHERE lower(username) = lower($1)
	`
	return scanUser(db.QueryRow(ctx, query, username))
}

func (r *userRepository) GetUserByNumber(ctx context.Context, db database.DBRunner, number string) (*models.User, error) {
	query := `
		SELECT id, username, display_name, email, password_hash, description, phone_number,
		       avatar_url, avatar_thumbnail_url, friend_policy, created_at, updated_at
		FROM users
		WHERE phone_number = $1
	`
	return scanUser(db.QueryRow(ctx, query, number))
}

func (r *userRepository) GetUserById(ctx context.Context, db database.DBRunner, userID uuid.UUID) (*models.User, error) {
	query := `
		SELECT id, username, display_name, email, password_hash, description, phone_number,
		       avatar_url, avatar_thumbnail_url, friend_policy, created_at, updated_at
		FROM users
		WHERE id = $1
	`
	return scanUser(db.QueryRow(ctx, query, userID))
}

func (r *userRepository) DoesUserExists(ctx context.Context, db database.DBRunner, userID uuid.UUID) (bool, error) {
	query := `SELECT EXISTS (SELECT 1 FROM users WHERE id = $1)`
	var exists bool
	err := db.QueryRow(ctx, query, userID).Scan(&exists)
	if err != nil {
		return false, err
	}
	return exists, nil
}

func (r *userRepository) UpdateUserById(ctx context.Context, db database.DBRunner, userID uuid.UUID, fields map[string]any) (*models.User, error) {
	setClauses := make([]string, 0, len(fields)+1)
	args := make([]any, 0, len(fields)+2)

	i := 1
	for col, val := range fields {
		setClauses = append(setClauses, fmt.Sprintf("%s = $%d", col, i))
		args = append(args, val)
		i++
	}

	setClauses = append(setClauses, fmt.Sprintf("updated_at = $%d", i))
	args = append(args, time.Now())
	i++

	args = append(args, userID)

	query := fmt.Sprintf(`
		UPDATE users
		SET %s
		WHERE id = $%d
		RETURNING id, username, display_name, email, password_hash, description, phone_number,
		          avatar_url, avatar_thumbnail_url, friend_policy, created_at, updated_at
	`, strings.Join(setClauses, ", "), i)

	return scanUser(db.QueryRow(ctx, query, args...))
}

func (r *userRepository) UpdateUsername(ctx context.Context, db database.DBRunner, userID uuid.UUID, username string) (*models.User, error) {
	return r.UpdateUserById(ctx, db, userID, map[string]any{
		"username": username,
	})
}

func (r *userRepository) UpdateEmail(ctx context.Context, db database.DBRunner, userID uuid.UUID, email string) (*models.User, error) {
	return r.UpdateUserById(ctx, db, userID, map[string]any{
		"email": email,
	})
}

func (r *userRepository) DeleteUserById(ctx context.Context, db database.DBRunner, userID uuid.UUID) error {
	tag, err := db.Exec(ctx, `DELETE FROM users WHERE id = $1`, userID)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return pgx.ErrNoRows
	}
	return nil
}

func (r *userRepository) CreateFriendRequest(ctx context.Context, db database.DBRunner, req *models.FriendRequest) (*models.FriendRequest, error) {
	query := `
		INSERT INTO friend_requests (id, sender_id, receiver_id)
		VALUES ($1, $2, $3)
		RETURNING id, sender_id, receiver_id, created_at
	`

	saved := &models.FriendRequest{}
	err := db.QueryRow(ctx, query, req.ID, req.SenderID, req.ReceiverID).Scan(
		&saved.ID,
		&saved.SenderID,
		&saved.ReceiverID,
		&saved.CreatedAt,
	)
	if err != nil {
		return nil, err
	}
	return saved, nil
}

func (r *userRepository) GetFriendRequestByID(ctx context.Context, db database.DBRunner, requestID uuid.UUID) (*models.FriendRequest, error) {
	query := `
		SELECT id, sender_id, receiver_id, created_at
		FROM friend_requests
		WHERE id = $1
	`
	req := &models.FriendRequest{}
	err := db.QueryRow(ctx, query, requestID).Scan(
		&req.ID,
		&req.SenderID,
		&req.ReceiverID,
		&req.CreatedAt,
	)
	if err != nil {
		return nil, err
	}
	return req, nil
}

func (r *userRepository) GetFriendRequestByUsers(ctx context.Context, db database.DBRunner, senderID uuid.UUID, receiverID uuid.UUID) (*models.FriendRequest, error) {
	query := `
		SELECT id, sender_id, receiver_id, created_at
		FROM friend_requests
		WHERE sender_id = $1 AND receiver_id = $2
	`
	req := &models.FriendRequest{}
	err := db.QueryRow(ctx, query, senderID, receiverID).Scan(
		&req.ID,
		&req.SenderID,
		&req.ReceiverID,
		&req.CreatedAt,
	)
	if err != nil {
		return nil, err
	}
	return req, nil
}

func (r *userRepository) DeleteFriendRequestByID(ctx context.Context, db database.DBRunner, requestID uuid.UUID) (*models.FriendRequest, error) {
	query := `
		DELETE FROM friend_requests
		WHERE id = $1
		RETURNING id, sender_id, receiver_id, created_at
	`
	req := &models.FriendRequest{}
	err := db.QueryRow(ctx, query, requestID).Scan(
		&req.ID,
		&req.SenderID,
		&req.ReceiverID,
		&req.CreatedAt,
	)
	if err != nil {
		return nil, err
	}
	return req, nil
}

func (r *userRepository) DeleteFriendRequestByUsers(ctx context.Context, db database.DBRunner, senderID uuid.UUID, receiverID uuid.UUID) error {
	tag, err := db.Exec(ctx, `
		DELETE FROM friend_requests
		WHERE sender_id = $1 AND receiver_id = $2
	`, senderID, receiverID)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return pgx.ErrNoRows
	}
	return nil
}

func (r *userRepository) ListIncomingFriendRequests(ctx context.Context, db database.DBRunner, userID uuid.UUID) ([]*models.FriendRequest, error) {
	query := `
		SELECT id, sender_id, receiver_id, created_at
		FROM friend_requests
		WHERE receiver_id = $1
		ORDER BY created_at DESC
	`
	rows, err := db.Query(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var requests []*models.FriendRequest
	for rows.Next() {
		current := &models.FriendRequest{}
		if err := rows.Scan(&current.ID, &current.SenderID, &current.ReceiverID, &current.CreatedAt); err != nil {
			return nil, err
		}
		requests = append(requests, current)
	}
	return requests, rows.Err()
}

func (r *userRepository) ListOutgoingFriendRequests(ctx context.Context, db database.DBRunner, userID uuid.UUID) ([]*models.FriendRequest, error) {
	query := `
		SELECT id, sender_id, receiver_id, created_at
		FROM friend_requests
		WHERE sender_id = $1
		ORDER BY created_at DESC
	`
	rows, err := db.Query(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var requests []*models.FriendRequest
	for rows.Next() {
		current := &models.FriendRequest{}
		if err := rows.Scan(&current.ID, &current.SenderID, &current.ReceiverID, &current.CreatedAt); err != nil {
			return nil, err
		}
		requests = append(requests, current)
	}
	return requests, rows.Err()
}

func (r *userRepository) DoesFriendRequestExist(ctx context.Context, db database.DBRunner, senderID uuid.UUID, receiverID uuid.UUID) (bool, error) {
	query := `
		SELECT EXISTS (
			SELECT 1 FROM friend_requests
			WHERE sender_id = $1 AND receiver_id = $2
		)
	`
	var exists bool
	err := db.QueryRow(ctx, query, senderID, receiverID).Scan(&exists)
	if err != nil {
		return false, err
	}
	return exists, nil
}

func (r *userRepository) CreateFriendship(ctx context.Context, db database.DBRunner, userID1 uuid.UUID, userID2 uuid.UUID) (*models.Friend, error) {
	userID1, userID2 = normalizedFriendPair(userID1, userID2)

	query := `
		INSERT INTO friends (user_id_1, user_id_2)
		VALUES ($1, $2)
		RETURNING user_id_1, user_id_2, created_at
	`

	friend := &models.Friend{}
	err := db.QueryRow(ctx, query, userID1, userID2).Scan(
		&friend.UserID1,
		&friend.UserID2,
		&friend.CreatedAt,
	)
	if err != nil {
		return nil, err
	}
	return friend, nil
}

func (r *userRepository) DeleteFriendship(ctx context.Context, db database.DBRunner, userID1 uuid.UUID, userID2 uuid.UUID) error {
	userID1, userID2 = normalizedFriendPair(userID1, userID2)

	tag, err := db.Exec(ctx, `
		DELETE FROM friends
		WHERE user_id_1 = $1 AND user_id_2 = $2
	`, userID1, userID2)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return pgx.ErrNoRows
	}
	return nil
}

func (r *userRepository) AreFriends(ctx context.Context, db database.DBRunner, userID1 uuid.UUID, userID2 uuid.UUID) (bool, error) {
	userID1, userID2 = normalizedFriendPair(userID1, userID2)

	query := `
		SELECT EXISTS (
			SELECT 1 FROM friends
			WHERE user_id_1 = $1 AND user_id_2 = $2
		)
	`
	var exists bool
	err := db.QueryRow(ctx, query, userID1, userID2).Scan(&exists)
	if err != nil {
		return false, err
	}
	return exists, nil
}

func (r *userRepository) ListFriends(ctx context.Context, db database.DBRunner, userID uuid.UUID) ([]*models.User, error) {
	query := `
		SELECT u.id, u.username, u.display_name, u.email, u.password_hash, u.description,
		       u.phone_number, u.avatar_url, u.avatar_thumbnail_url, u.friend_policy,
		       u.created_at, u.updated_at
		FROM friends f
		INNER JOIN users u
			ON u.id = CASE
				WHEN f.user_id_1 = $1 THEN f.user_id_2
				ELSE f.user_id_1
			END
		WHERE f.user_id_1 = $1 OR f.user_id_2 = $1
		ORDER BY lower(u.username) ASC
	`
	rows, err := db.Query(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []*models.User
	for rows.Next() {
		u, err := scanUser(rows)
		if err != nil {
			return nil, err
		}
		users = append(users, u)
	}
	return users, rows.Err()
}

func (r *userRepository) CountFriends(ctx context.Context, db database.DBRunner, userID uuid.UUID) (int, error) {
	query := `
		SELECT COUNT(*)
		FROM friends
		WHERE user_id_1 = $1 OR user_id_2 = $1
	`
	var count int
	err := db.QueryRow(ctx, query, userID).Scan(&count)
	if err != nil {
		return 0, err
	}
	return count, nil
}

func (r *userRepository) GetMutualFriends(ctx context.Context, db database.DBRunner, currentUserID uuid.UUID, targetUserID uuid.UUID) ([]*models.User, error) {
	query := `
		WITH current_user_friends AS (
			SELECT CASE
				WHEN user_id_1 = $1 THEN user_id_2
				ELSE user_id_1
			END AS friend_id
			FROM friends
			WHERE user_id_1 = $1 OR user_id_2 = $1
		),
		target_user_friends AS (
			SELECT CASE
				WHEN user_id_1 = $2 THEN user_id_2
				ELSE user_id_1
			END AS friend_id
			FROM friends
			WHERE user_id_1 = $2 OR user_id_2 = $2
		)
		SELECT u.id, u.username, u.display_name, u.email, u.password_hash, u.description,
		       u.phone_number, u.avatar_url, u.avatar_thumbnail_url, u.friend_policy,
		       u.created_at, u.updated_at
		FROM users u
		INNER JOIN current_user_friends cuf ON cuf.friend_id = u.id
		INNER JOIN target_user_friends tuf ON tuf.friend_id = u.id
		ORDER BY lower(u.username) ASC
	`
	rows, err := db.Query(ctx, query, currentUserID, targetUserID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []*models.User
	for rows.Next() {
		u, err := scanUser(rows)
		if err != nil {
			return nil, err
		}
		users = append(users, u)
	}
	return users, rows.Err()
}

func (r *userRepository) CountMutualFriends(ctx context.Context, db database.DBRunner, currentUserID uuid.UUID, targetUserID uuid.UUID) (int, error) {
	query := `
		WITH current_user_friends AS (
			SELECT CASE
				WHEN user_id_1 = $1 THEN user_id_2
				ELSE user_id_1
			END AS friend_id
			FROM friends
			WHERE user_id_1 = $1 OR user_id_2 = $1
		),
		target_user_friends AS (
			SELECT CASE
				WHEN user_id_1 = $2 THEN user_id_2
				ELSE user_id_1
			END AS friend_id
			FROM friends
			WHERE user_id_1 = $2 OR user_id_2 = $2
		)
		SELECT COUNT(*)
		FROM current_user_friends cuf
		INNER JOIN target_user_friends tuf ON tuf.friend_id = cuf.friend_id
	`
	var count int
	err := db.QueryRow(ctx, query, currentUserID, targetUserID).Scan(&count)
	if err != nil {
		return 0, err
	}
	return count, nil
}

func (r *userRepository) UpsertAppLink(ctx context.Context, db database.DBRunner, link *models.UserAppLink) (*models.UserAppLink, error) {
	query := `
		INSERT INTO user_app_links (id, user_id, provider, url, show_on_profile)
		VALUES ($1, $2, $3, $4, $5)
		ON CONFLICT (user_id, provider)
		DO UPDATE SET
			url = EXCLUDED.url,
			show_on_profile = EXCLUDED.show_on_profile,
			updated_at = now()
		RETURNING id, user_id, provider, url, show_on_profile, created_at, updated_at
	`

	out := &models.UserAppLink{}
	err := db.QueryRow(ctx, query,
		link.ID,
		link.UserID,
		link.Provider,
		link.URL,
		link.ShowOnProfile,
	).Scan(
		&out.ID,
		&out.UserID,
		&out.Provider,
		&out.URL,
		&out.ShowOnProfile,
		&out.CreatedAt,
		&out.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return out, nil
}

func (r *userRepository) DeleteAppLink(ctx context.Context, db database.DBRunner, userID uuid.UUID, provider models.AppProvider) error {
	tag, err := db.Exec(ctx, `
		DELETE FROM user_app_links
		WHERE user_id = $1 AND provider = $2
	`, userID, provider)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return pgx.ErrNoRows
	}
	return nil
}

func (r *userRepository) GetUserAppLinks(ctx context.Context, db database.DBRunner, userID uuid.UUID, onlyVisible bool) ([]*models.UserAppLink, error) {
	query := `
		SELECT id, user_id, provider, url, show_on_profile, created_at, updated_at
		FROM user_app_links
		WHERE user_id = $1
	`
	args := []any{userID}

	if onlyVisible {
		query += ` AND show_on_profile = true`
	}

	query += ` ORDER BY provider ASC`

	rows, err := db.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var links []*models.UserAppLink
	for rows.Next() {
		current := &models.UserAppLink{}
		if err := rows.Scan(
			&current.ID,
			&current.UserID,
			&current.Provider,
			&current.URL,
			&current.ShowOnProfile,
			&current.CreatedAt,
			&current.UpdatedAt,
		); err != nil {
			return nil, err
		}
		links = append(links, current)
	}

	return links, rows.Err()
}

func (r *userRepository) GetUserAppLinkByProvider(ctx context.Context, db database.DBRunner, userID uuid.UUID, provider models.AppProvider) (*models.UserAppLink, error) {
	query := `
		SELECT id, user_id, provider, url, show_on_profile, created_at, updated_at
		FROM user_app_links
		WHERE user_id = $1 AND provider = $2
	`

	link := &models.UserAppLink{}
	err := db.QueryRow(ctx, query, userID, provider).Scan(
		&link.ID,
		&link.UserID,
		&link.Provider,
		&link.URL,
		&link.ShowOnProfile,
		&link.CreatedAt,
		&link.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return link, nil
}
