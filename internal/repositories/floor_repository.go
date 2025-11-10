package repositories

import (
	"context"

	"github.com/google/uuid"
	"github.com/suck-seed/yapp/internal/models"
)

type IFloorRepository interface {
	CreateFloor(ctx context.Context, floor *models.Floor) (*models.Floor, error)
	DoesFloorExistsInRoom(ctx context.Context, floorID *uuid.UUID, hallID *uuid.UUID) (bool, error)
}

type floorRepository struct {
	db PGXTX
}

func NewFloorRepository(db PGXTX) IFloorRepository {

	return &floorRepository{
		db: db,
	}

}

func (r *floorRepository) CreateFloor(ctx context.Context, floor *models.Floor) (*models.Floor, error) {

	query := `

	INSERT INTO floors (id, hall_id, name, is_private, created_at, updated_at)
	VALUES ($1, $2, $3, $4, $5, $6)
	RETURNING  id, hall_id, name, is_private, created_at, updated_at

	`

	row := r.db.QueryRow(ctx, query,

		floor.ID,
		floor.HallID,
		floor.Name,
		floor.IsPrivate,
		floor.CreatedAt,
		floor.UpdatedAt,
	)

	floorCRES := &models.Floor{}

	err := row.Scan(
		&floorCRES.ID,
		&floorCRES.HallID,
		&floorCRES.Name,
		&floorCRES.IsPrivate,
		&floorCRES.CreatedAt,
		&floorCRES.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return floorCRES, nil
}

func (r *floorRepository) DoesFloorExistsInRoom(ctx context.Context, floorID *uuid.UUID, hallID *uuid.UUID) (bool, error) {

	query := `

		SELECT EXISTS (SELECT 1 FROM floors WHERE id = $1 and hall_id = $2)

	`

	var exists bool

	err := r.db.QueryRow(ctx, query, floorID, hallID).Scan(&exists)
	if err != nil {
		return false, err
	}

	return exists, nil
}
