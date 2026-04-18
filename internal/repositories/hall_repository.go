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

type IHallRepository interface {
	// -------------------------- HALL
	// core cud
	CreateHall(ctx context.Context, db database.DBRunner, hall *models.Hall) (*models.Hall, error)                         // C
	CreateHallMember(ctx context.Context, db database.DBRunner, hallMember *models.HallMember) (*models.HallMember, error) // C
	DeleteHall(ctx context.Context, db database.DBRunner, hallID uuid.UUID) (*models.Hall, error)                          // D

	// list operation
	GetUserHallIDs(ctx context.Context, db database.DBRunner, userID uuid.UUID) ([]uuid.UUID, error) // R
	GetHallByID(ctx context.Context, db database.DBRunner, hallID uuid.UUID) (*models.Hall, error)   // R
	GetHallOwnerID(ctx context.Context, db database.DBRunner, hallID uuid.UUID) (uuid.UUID, error)   // R
	GetHallMemberByUserID(ctx context.Context, db database.DBRunner, hallID uuid.UUID, userID uuid.UUID) (*models.HallMember, error)
	GetHallMemberByID(ctx context.Context, db database.DBRunner, hallID uuid.UUID, memberID uuid.UUID) (*models.HallMember, error)
	ListHallMembers(ctx context.Context, db database.DBRunner, hallID uuid.UUID) ([]*models.HallMember, error)
	UpdateHallMember(ctx context.Context, db database.DBRunner, hallID uuid.UUID, userID uuid.UUID, fields map[string]any) (*models.HallMember, error)

	// ---------------- USER MANAGEMENT
	KickHallMember(ctx context.Context, db database.DBRunner, hallID uuid.UUID, userID uuid.UUID) error

	// -------------- HALL PROFILE
	UpdateHallProfile(ctx context.Context, db database.DBRunner, hallID uuid.UUID, fields map[string]any) (*models.Hall, error)

	// ------------- CHECK OPERATION
	DoesHallExist(ctx context.Context, db database.DBRunner, hallID uuid.UUID) (bool, error)
	IsUserHallMember(ctx context.Context, db database.DBRunner, hallID uuid.UUID, userID uuid.UUID) (bool, error)

	// ---------------- JOIN OPERATIONS
	CreateJoinRequest(ctx context.Context, db database.DBRunner, request *models.HallRequest) (*models.HallRequest, error)
	GetJoinRequestByID(ctx context.Context, db database.DBRunner, requestID uuid.UUID) (*models.HallRequest, error)
	GetJoinRequestByHallAndUser(ctx context.Context, db database.DBRunner, hallID, userID uuid.UUID) (*models.HallRequest, error)
	GetAllHallRequests(ctx context.Context, db database.DBRunner, hallID uuid.UUID) ([]*models.HallRequest, error)
	DeleteJoinRequest(ctx context.Context, db database.DBRunner, requestID uuid.UUID) (*models.HallRequest, error)
	DoesPendingJoinRequestExist(ctx context.Context, db database.DBRunner, hallID, userID uuid.UUID) (bool, error)
}

type hallRepository struct {
}

func NewHallRepository() IHallRepository {

	return &hallRepository{}
}

func (r *hallRepository) CreateHall(ctx context.Context, db database.DBRunner, hall *models.Hall) (*models.Hall, error) {

	query := `

	INSERT INTO halls (id, name, is_private, icon_url, icon_thumbnail_url, banner_color, description, owner_id)
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	RETURNING id, name, is_private, icon_url, icon_thumbnail_url, banner_color, description, owner_id, created_at, updated_at

	`

	row := db.QueryRow(ctx, query,
		hall.ID,
		hall.Name,
		hall.IsPrivate,
		hall.IconURL,
		hall.IconThumbnailURL,
		hall.BannerColor,
		hall.Description,
		hall.OwnerID,
	)

	saved := &models.Hall{}

	err := row.Scan(
		&saved.ID,
		&saved.Name,
		&saved.IsPrivate,
		&saved.IconURL,
		&saved.IconThumbnailURL,
		&saved.BannerColor,
		&saved.Description,
		&saved.OwnerID,
		&saved.CreatedAt,
		&saved.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return saved, nil
}

// Updated implementation
func (r *hallRepository) CreateHallMember(ctx context.Context, db database.DBRunner, hallMember *models.HallMember) (*models.HallMember, error) {
	query := `
		INSERT INTO hall_members (id, hall_id, user_id, role_id, nickname)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, hall_id, user_id, role_id, nickname, joined_at, created_at, updated_at
	`

	saved := &models.HallMember{}
	err := db.QueryRow(ctx, query,
		hallMember.ID,
		hallMember.HallID,
		hallMember.UserID,
		hallMember.RoleID,
		hallMember.Nickname,
	).Scan(
		&saved.ID,
		&saved.HallID,
		&saved.UserID,
		&saved.RoleID,
		&saved.Nickname,
		&saved.JoinedAt,
		&saved.CreatedAt,
		&saved.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return saved, nil
}

func (r *hallRepository) DeleteHall(ctx context.Context, db database.DBRunner, hallID uuid.UUID) (*models.Hall, error) {
	query := `
		DELETE FROM halls
		WHERE id = $1
		RETURNING id, name, is_private, icon_url, icon_thumbnail_url, banner_color, description, created_at, updated_at, owner_id
	`

	hall := &models.Hall{}
	err := db.QueryRow(ctx, query, hallID).Scan(
		&hall.ID,
		&hall.Name,
		&hall.IsPrivate,
		&hall.IconURL,
		&hall.IconThumbnailURL,
		&hall.BannerColor,
		&hall.Description,
		&hall.CreatedAt,
		&hall.UpdatedAt,
		&hall.OwnerID,
	)
	if err != nil {
		return nil, err
	}

	return hall, nil
}

func (r *hallRepository) GetUserHallIDs(ctx context.Context, db database.DBRunner, userID uuid.UUID) ([]uuid.UUID, error) {

	query := `
	SELECT hall_id FROM hall_members WHERE user_id = $1
	`

	rows, err := db.Query(ctx, query, userID) // Fixed: use userID, not hallID
	if err != nil {
		return nil, fmt.Errorf("failed to query user hall IDs: %w", err)
	}
	defer rows.Close()

	var hallIDs []uuid.UUID

	for rows.Next() {
		var hallID uuid.UUID
		err := rows.Scan(&hallID)
		if err != nil {
			return nil, fmt.Errorf("failed to scan hall ID: %w", err)
		}
		hallIDs = append(hallIDs, hallID)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error during rows iteration: %w", err)
	}

	return hallIDs, nil
}

func (r *hallRepository) GetHallByID(ctx context.Context, db database.DBRunner, hallID uuid.UUID) (*models.Hall, error) {
	hall := &models.Hall{}

	query := `SELECT id, name, is_private, icon_url, icon_thumbnail_url, banner_color, description, created_at, updated_at, owner_id
              FROM halls WHERE id = $1
              `

	err := db.QueryRow(ctx, query, hallID).Scan(
		&hall.ID,
		&hall.Name,
		&hall.IsPrivate,
		&hall.IconURL,
		&hall.IconThumbnailURL,
		&hall.BannerColor,
		&hall.Description,
		&hall.CreatedAt,
		&hall.UpdatedAt,
		&hall.OwnerID,
	)

	if err != nil {
		return nil, err
	}

	return hall, nil
}

func (r *hallRepository) GetHallOwnerID(ctx context.Context, db database.DBRunner, hallID uuid.UUID) (uuid.UUID, error) {

	var ownerID uuid.UUID

	query := `SELECT owner_id
              FROM halls WHERE id = $1`

	err := db.QueryRow(ctx, query, hallID).Scan(
		&ownerID,
	)

	if err != nil {
		return uuid.Nil, err
	}

	return ownerID, nil
}

func (r *hallRepository) GetHallMemberByUserID(ctx context.Context, db database.DBRunner, hallID uuid.UUID, userID uuid.UUID) (*models.HallMember, error) {
	query := `
		SELECT id, hall_id, user_id, role_id, nickname, joined_at, created_at, updated_at
		FROM hall_members
		WHERE hall_id = $1 AND user_id = $2
	`

	m := &models.HallMember{}
	err := db.QueryRow(ctx, query, hallID, userID).Scan(
		&m.ID,
		&m.HallID,
		&m.UserID,
		&m.RoleID,
		&m.Nickname,
		&m.JoinedAt,
		&m.CreatedAt,
		&m.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return m, nil
}

func (r *hallRepository) GetHallMemberByID(ctx context.Context, db database.DBRunner, hallID uuid.UUID, memberID uuid.UUID) (*models.HallMember, error) {
	query := `
		SELECT id, hall_id, user_id, role_id, nickname, joined_at, created_at, updated_at
		FROM hall_members
		WHERE hall_id = $1 AND id = $2
	`

	m := &models.HallMember{}
	err := db.QueryRow(ctx, query, hallID, memberID).Scan(
		&m.ID,
		&m.HallID,
		&m.UserID,
		&m.RoleID,
		&m.Nickname,
		&m.JoinedAt,
		&m.CreatedAt,
		&m.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return m, nil
}

func (r *hallRepository) ListHallMembers(ctx context.Context, db database.DBRunner, hallID uuid.UUID) ([]*models.HallMember, error) {
	query := `
		SELECT id, hall_id, user_id, role_id, nickname, joined_at, created_at, updated_at
		FROM hall_members
		WHERE hall_id = $1
		ORDER BY joined_at ASC
	`

	rows, err := db.Query(ctx, query, hallID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var members []*models.HallMember
	for rows.Next() {
		m := &models.HallMember{}
		if err := rows.Scan(
			&m.ID,
			&m.HallID,
			&m.UserID,
			&m.RoleID,
			&m.Nickname,
			&m.JoinedAt,
			&m.CreatedAt,
			&m.UpdatedAt,
		); err != nil {
			return nil, err
		}
		members = append(members, m)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return members, nil
}

// UpdateHallMember — swap the WHERE clause to use user_id instead of id
func (r *hallRepository) UpdateHallMember(ctx context.Context, db database.DBRunner, hallID uuid.UUID, userID uuid.UUID, fields map[string]any) (*models.HallMember, error) {
	setClauses := make([]string, 0, len(fields)+1)
	args := make([]any, 0, len(fields)+3)

	i := 1
	for col, val := range fields {
		setClauses = append(setClauses, fmt.Sprintf("%s = $%d", col, i))
		args = append(args, val)
		i++
	}
	setClauses = append(setClauses, fmt.Sprintf("updated_at = $%d", i))
	args = append(args, time.Now())
	i++

	args = append(args, hallID, userID)

	query := fmt.Sprintf(`
		UPDATE hall_members
		SET %s
		WHERE hall_id = $%d AND user_id = $%d
		RETURNING id, hall_id, user_id, role_id, nickname, joined_at, created_at, updated_at
	`, strings.Join(setClauses, ", "), i, i+1)

	m := &models.HallMember{}
	err := db.QueryRow(ctx, query, args...).Scan(
		&m.ID,
		&m.HallID,
		&m.UserID,
		&m.RoleID,
		&m.Nickname,
		&m.JoinedAt,
		&m.CreatedAt,
		&m.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return m, nil
}

// DeleteHallMember — same swap
func (r *hallRepository) KickHallMember(ctx context.Context, db database.DBRunner, hallID uuid.UUID, userID uuid.UUID) error {
	query := `
		DELETE FROM hall_members
		WHERE hall_id = $1 AND user_id = $2
	`

	ct, err := db.Exec(ctx, query, hallID, userID)
	if err != nil {
		return err
	}

	if ct.RowsAffected() == 0 {
		return pgx.ErrNoRows
	}

	return nil
}

func (r *hallRepository) UpdateHallProfile(ctx context.Context, db database.DBRunner, hallID uuid.UUID, fields map[string]any) (*models.Hall, error) {

	// Allowed columns and their sanitized values are built by the service.
	// We build a dynamic SET clause: SET name = $1, description = $2 ... WHERE id = $N
	setClauses := make([]string, 0, len(fields))
	args := make([]any, 0, len(fields)+1)

	i := 1
	for col, val := range fields {
		setClauses = append(setClauses, fmt.Sprintf("%s = $%d", col, i))
		args = append(args, val)
		i++
	}
	// always bump updated_at
	setClauses = append(setClauses, fmt.Sprintf("updated_at = $%d", i))
	args = append(args, time.Now())
	i++

	args = append(args, hallID) // WHERE id = $N

	query := fmt.Sprintf(`
			UPDATE halls
			SET %s
			WHERE id = $%d
			RETURNING id, name, is_private, icon_url, icon_thumbnail_url, banner_color, description, owner_id, created_at, updated_at
		`, strings.Join(setClauses, ", "), i)

	hall := &models.Hall{}
	err := db.QueryRow(ctx, query, args...).Scan(
		&hall.ID,
		&hall.Name,
		&hall.IsPrivate,
		&hall.IconURL,
		&hall.IconThumbnailURL,
		&hall.BannerColor,
		&hall.Description,
		&hall.OwnerID,
		&hall.CreatedAt,
		&hall.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return hall, nil
}

func (r *hallRepository) DoesHallExist(ctx context.Context, db database.DBRunner, hallID uuid.UUID) (bool, error) {

	query := `

		SELECT EXISTS (SELECT 1 FROM halls WHERE id = $1)

	`

	var exists bool

	err := db.QueryRow(ctx, query, hallID).Scan(&exists)
	if err != nil {
		return false, err
	}

	return exists, nil
}

func (r *hallRepository) IsUserHallMember(ctx context.Context, db database.DBRunner, hallID uuid.UUID, userID uuid.UUID) (bool, error) {

	query := `

	SELECT EXISTS (SELECT 1 FROM hall_members WHERE hall_id = $1 and user_id = $2)
`

	var exists bool

	if err := db.QueryRow(ctx, query, hallID, userID).Scan(&exists); err != nil {
		return false, err
	}

	return exists, nil

}

// ------------------- JOIN REQUEST OPERATIONS

func (r *hallRepository) CreateJoinRequest(ctx context.Context, db database.DBRunner, request *models.HallRequest) (*models.HallRequest, error) {
	query := `
		INSERT INTO hall_requests (id, hall_id, user_id)
		VALUES ($1, $2, $3)
		RETURNING id, hall_id, user_id, created_at, updated_at
	`

	saved := &models.HallRequest{}
	err := db.QueryRow(ctx, query,
		request.ID,
		request.HallID,
		request.UserID,
	).Scan(
		&saved.ID,
		&saved.HallID,
		&saved.UserID,
		&saved.CreatedAt,
		&saved.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return saved, nil
}

func (r *hallRepository) GetJoinRequestByID(ctx context.Context, db database.DBRunner, requestID uuid.UUID) (*models.HallRequest, error) {
	query := `
		SELECT id, hall_id, user_id, created_at, updated_at
		FROM hall_requests
		WHERE id = $1
	`

	saved := &models.HallRequest{}
	err := db.QueryRow(ctx, query, requestID).Scan(
		&saved.ID,
		&saved.HallID,
		&saved.UserID,
		&saved.CreatedAt,
		&saved.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return saved, nil
}

func (r *hallRepository) GetJoinRequestByHallAndUser(ctx context.Context, db database.DBRunner, hallID, userID uuid.UUID) (*models.HallRequest, error) {
	query := `
		SELECT id, hall_id, user_id, created_at, updated_at
		FROM hall_requests
		WHERE hall_id = $1 AND user_id = $2
	`

	saved := &models.HallRequest{}
	err := db.QueryRow(ctx, query, hallID, userID).Scan(
		&saved.ID,
		&saved.HallID,
		&saved.UserID,
		&saved.CreatedAt,
		&saved.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return saved, nil
}

func (r *hallRepository) GetAllHallRequests(ctx context.Context, db database.DBRunner, hallID uuid.UUID) ([]*models.HallRequest, error) {
	query := `
		SELECT id, hall_id, user_id, created_at, updated_at
		FROM hall_requests
		WHERE hall_id = $1
		ORDER BY created_at ASC
	`

	rows, err := db.Query(ctx, query, hallID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var requests []*models.HallRequest
	for rows.Next() {
		current := &models.HallRequest{}
		if err := rows.Scan(
			&current.ID,
			&current.HallID,
			&current.UserID,
			&current.CreatedAt,
			&current.UpdatedAt,
		); err != nil {
			return nil, err
		}
		requests = append(requests, current)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return requests, nil
}

func (r *hallRepository) DeleteJoinRequest(ctx context.Context, db database.DBRunner, requestID uuid.UUID) (*models.HallRequest, error) {
	query := `
		DELETE FROM hall_requests
		WHERE id = $1
		RETURNING id, hall_id, user_id, created_at, updated_at
	`

	deleted := &models.HallRequest{}
	err := db.QueryRow(ctx, query, requestID).Scan(
		&deleted.ID,
		&deleted.HallID,
		&deleted.UserID,
		&deleted.CreatedAt,
		&deleted.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return deleted, nil
}

func (r *hallRepository) DoesPendingJoinRequestExist(ctx context.Context, db database.DBRunner, hallID, userID uuid.UUID) (bool, error) {
	query := `
		SELECT EXISTS (
			SELECT 1
			FROM hall_requests
			WHERE hall_id = $1 AND user_id = $2
		)
	`

	var exists bool
	if err := db.QueryRow(ctx, query, hallID, userID).Scan(&exists); err != nil {
		return false, err
	}

	return exists, nil
}
