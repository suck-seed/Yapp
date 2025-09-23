package repositories

import (
	"context"

	"github.com/suck-seed/yapp/internal/dto"
	"github.com/suck-seed/yapp/internal/models"
)

type IHallRepository interface {
	CreateHall(ctx context.Context, hall *dto.CreateHallReq) (*models.Hall, error)
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

func (r *hallRepository) CreateHall(ctx context.Context, hall *dto.CreateHallReq) (*models.Hall, error) {

	query := `
	INSERT INTO halls (id, name, icon_url, description, created_by_id)
    VALUES ($1, $2, $3, $4, $5)
    RETURNING id, name, icon_url, banner_color, description, created_at, updated_at, created_by_id
	`

	row := r.db.QueryRow(ctx, query,
		hall.ID,
		hall.Name,
		hall.IconURL,
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

	return &models.Hall{}, nil
}
