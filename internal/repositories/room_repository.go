package repositories

import (
	"context"
	"github.com/suck-seed/yapp/internal/models"
)

type IRoomRepository interface {
	CreateRoom(c context.Context, room *models.Room) (*models.Room, error)
}

type roomRepository struct {
	db PGXTX
}

func NewRoomRepository(db PGXTX) IRoomRepository {

	return &roomRepository{
		db: db,
	}
}

func (r *roomRepository) CreateRoom(c context.Context, room *models.Room) (*models.Room, error) {

	query := `


	INSERT INTO rooms (id, hall_id, floor_id, name, room_type, is_private, created_at, updated_at)
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	RETURNING  id, hall_id, floor_id, name, room_type, is_private, created_at, updated_at


`

	row := r.db.QueryRow(
		c,
		query,
		room.ID,
		room.HallId,
		room.FloorId,
		room.Name,
		room.RoomType,
		room.IsPrivate,
		room.CreatedAt,
		room.UpdatedAt,
	)

	roomCRES := &models.Room{}

	err := row.Scan(
		&roomCRES.ID,
		&roomCRES.HallId,
		&roomCRES.FloorId,
		&roomCRES.Name,
		&roomCRES.RoomType,
		&roomCRES.IsPrivate,
		&roomCRES.CreatedAt,
		&roomCRES.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return roomCRES, nil
}
