package repositories

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/suck-seed/yapp/internal/models"
)

type IHallRepository interface {
	CreateHall(ctx context.Context, hall *models.Hall) (*models.Hall, error)
	CreateHallRole(ctx context.Context, hallRole *models.Role) (*models.Role, error)
	CreateHallMember(ctx context.Context, hallMember *models.HallMember) error

	GetUserHallIDs(ctx context.Context, userID uuid.UUID) ([]uuid.UUID, error)
	GetHallByID(ctx context.Context, hallID uuid.UUID) (*models.Hall, error)

	DoesHallExists(ctx context.Context, hallID uuid.UUID) (*bool, error)
	IsUserHallMember(ctx context.Context, hallID uuid.UUID, userID uuid.UUID) (*bool, error)
}

type hallRepository struct {
	db PGXTX
}

func NewHallRepository(db PGXTX) IHallRepository {

	return &hallRepository{
		db: db,
	}
}

func (r *hallRepository) CreateHall(ctx context.Context, hall *models.Hall) (*models.Hall, error) {

	query := `

	INSERT INTO halls (id, name, icon_url, icon_thumbnail_url, banner_color, description, owner)
    VALUES ($1, $2, $3, $4, $5, $6, $7)
    RETURNING id, name, icon_url, icon_thumbnail_url, banner_color, description, owner, created_at, updated_at
	`

	row := r.db.QueryRow(ctx, query,
		hall.ID,
		hall.Name,
		hall.IconURL,
		hall.IconThumbnailURL,
		hall.BannerColor,
		hall.Description,
		hall.Owner,
	)

	saved := &models.Hall{}

	err := row.Scan(
		&saved.ID,
		&saved.Name,
		&saved.IconURL,
		&saved.IconThumbnailURL,
		&saved.BannerColor,
		&saved.Description,
		&saved.Owner,
		&saved.CreatedAt,
		&saved.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return saved, nil
}

func (r *hallRepository) CreateHallRole(ctx context.Context, hallRole *models.Role) (*models.Role, error) {
	query := `
    INSERT INTO roles (id, hall_id, name, color, icon_url, is_default, is_admin)
    VALUES ($1, $2, $3, $4, $5, $6, $7)
    RETURNING id, hall_id, name, color, icon_url, is_default, is_admin, created_at, updated_at
    `

	row := r.db.QueryRow(ctx, query,
		hallRole.ID,
		hallRole.HallID,
		hallRole.Name,
		hallRole.Color,
		hallRole.IconURL,
		hallRole.IsDefault,
		hallRole.IsAdmin,
	)

	saved := &models.Role{}
	err := row.Scan(
		&saved.ID,
		&saved.HallID,
		&saved.Name,
		&saved.Color,
		&saved.IconURL,
		&saved.IsDefault,
		&saved.IsAdmin,
		&saved.CreatedAt,
		&saved.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return saved, nil
}

func (r *hallRepository) CreateHallMember(ctx context.Context, hallMember *models.HallMember) error {
	query := `
    INSERT INTO hall_members (id, hall_id, user_id,role_id)
    VALUES ($1, $2, $3, $4)
    `

	_, err := r.db.Exec(ctx, query,
		hallMember.ID,
		hallMember.HallID,
		hallMember.UserID,
		hallMember.RoleID,
	)

	if err != nil {
		return err
	}

	return nil
}

func (r *hallRepository) GetUserHallIDs(ctx context.Context, userID uuid.UUID) ([]uuid.UUID, error) {
	query := `SELECT hall_id FROM hall_members WHERE user_id = $1`

	rows, err := r.db.Query(ctx, query, userID) // Fixed: use userID, not hallID
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

func (r *hallRepository) GetHallByID(ctx context.Context, hallID uuid.UUID) (*models.Hall, error) {
	hall := &models.Hall{}

	query := `SELECT id, name, icon_url, icon_thumbnail_url, banner_color, description, created_at, updated_at, created_by_id
              FROM halls WHERE id = $1`

	err := r.db.QueryRow(ctx, query, hallID).Scan(
		&hall.ID,
		&hall.Name,
		&hall.IconURL,
		&hall.IconThumbnailURL,
		&hall.BannerColor,
		&hall.Description,
		&hall.CreatedAt,
		&hall.UpdatedAt,
		&hall.Owner,
	)

	if err != nil {
		return nil, err
	}

	return hall, nil
}

func (r *hallRepository) DoesHallExists(ctx context.Context, hallID uuid.UUID) (*bool, error) {

	query := `

		SELECT EXISTS (SELECT 1 FROM halls WHERE id = $1)

	`

	var exists bool

	err := r.db.QueryRow(ctx, query, hallID).Scan(&exists)
	if err != nil {
		return nil, err
	}

	return &exists, nil
}

func (r *hallRepository) IsUserHallMember(ctx context.Context, hallID uuid.UUID, userID uuid.UUID) (*bool, error) {

	query := `

	SELECT EXISTS (SELECT 1 FROM hall_members WHERE hall_id = $1 and user_id = $2)
`

	var exists bool

	if err := r.db.QueryRow(ctx, query, hallID, userID).Scan(&exists); err != nil {
		return nil, err
	}

	return &exists, nil

}
