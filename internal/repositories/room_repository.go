package repositories

import (
	"context"
	"github.com/google/uuid"
	"github.com/suck-seed/yapp/internal/models"
)

type IRoomRepository interface {
	CreateRoom(ctx context.Context, room *models.Room) (*models.Room, error)
	GetRoomByID(ctx context.Context, roomID *uuid.UUID) (*models.Room, error)
	IsUserRoomMember(c context.Context, roomId *uuid.UUID, userId *uuid.UUID) (*bool, error)
}

type roomRepository struct {
	db PGXTX
}

func NewRoomRepository(db PGXTX) IRoomRepository {

	return &roomRepository{
		db: db,
	}
}

func (r *roomRepository) CreateRoom(ctx context.Context, room *models.Room) (*models.Room, error) {

	query := `


	INSERT INTO rooms (id, hall_id, floor_id, name, room_type, is_private, created_at, updated_at)
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	RETURNING  id, hall_id, floor_id, name, room_type, is_private, created_at, updated_at


`

	row := r.db.QueryRow(
		ctx,
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

func (r *roomRepository) GetRoomByID(c context.Context, roomID *uuid.UUID) (*models.Room, error) {

	room := &models.Room{}

	query := `
				SELECT id, hall_id, floor_id, name, room_type, is_private, created_at, updated_at
				FROM rooms
				WHERE id = $1
			`

	row := r.db.QueryRow(c, query, roomID)

	err := row.Scan(
		&room.ID,
		&room.HallId,
		&room.FloorId,
		&room.Name,
		&room.RoomType,
		&room.IsPrivate,
		&room.CreatedAt,
		&room.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return room, nil

}

func (r *roomRepository) IsUserRoomMember(c context.Context, roomID *uuid.UUID, userID *uuid.UUID) (*bool, error) {

	query := `
		SELECT EXISTS (SELECT 1 FROM room_members WHERE room_id = $1 AND user_id = $2)
`
	var exists bool

	err := r.db.QueryRow(c, query, roomID, userID).Scan(&exists)
	if err != nil {
		return nil, err
	}

	return &exists, err

}
