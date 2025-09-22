package repositories

import (
	"context"

	"github.com/suck-seed/yapp/internal/models"
)

type IHallRepository interface {
	CreateHall(ctx context.Context, hall *models.Hall) (*models.Hall, error)
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

	INSERT INTO halls (id, name, icon_url, description, owner)
    VALUES ($1, $2, $3, $4, $5)
    RETURNING id, name, icon_url, banner_color, description, owner, created_at, updated_at
	`

	row := r.db.QueryRow(ctx, query,
		hall.ID,
		hall.Name,
		hall.IconURL,
		hall.Description,
		hall.Owner,
	)

	saved := &models.Hall{}

	err := row.Scan(
		&saved.ID,
		&saved.Name,
		&saved.IconURL,
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
