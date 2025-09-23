package repositories

import (
	"context"

	"github.com/suck-seed/yapp/internal/models"
)

type IHallRepository interface {
	CreateHall(ctx context.Context, hall *models.Hall) (*models.Hall, error)
	GetHallByName(ctx context.Context, hallName string) (*models.Hall, error)
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
	INSERT INTO halls (id, name, icon_url, banner_color, description, created_by_id)
    VALUES ($1, $2, $3, $4, $5, $6)
    RETURNING id, name, icon_url, banner_color, description, created_at, updated_at, created_by_id
	`

	row := r.db.QueryRow(ctx, query,
		hall.ID,
		hall.Name,
		hall.IconURL,
		hall.BannerColor,
		hall.Description,
		hall.CreatedBy,
	)

	saved := &models.Hall{}

	err := row.Scan(
		&saved.ID,
		&saved.Name,
		&saved.IconURL,
		&saved.BannerColor,
		&saved.Description,
		&saved.CreatedAt,
		&saved.UpdatedAt,
		&saved.CreatedBy,
	)

	if err != nil {
		return nil, err
	}

	return saved, nil
}

func (r *hallRepository) GetHallByName(ctx context.Context, hallName string) (*models.Hall, error) {
	hall := &models.Hall{}

	query := `SELECT id, name, icon_url, banner_color, description, created_at, updated_at, created_by_id
              FROM halls WHERE name = $1`

	err := r.db.QueryRow(ctx, query, hallName).Scan(
		&hall.ID,
		&hall.Name,
		&hall.IconURL,
		&hall.BannerColor,
		&hall.Description,
		&hall.CreatedAt,
		&hall.UpdatedAt,
		&hall.CreatedBy,
	)

	if err != nil {
		return nil, err
	}

	return hall, nil
}
