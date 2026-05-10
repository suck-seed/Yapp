package repositories

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/suck-seed/yapp/internal/database"
	dto "github.com/suck-seed/yapp/internal/dto/room"
	"github.com/suck-seed/yapp/internal/models"
)

type IRoomRepository interface {
	// existing
	CreateRoom(ctx context.Context, db database.DBRunner, room *models.Room) (*models.Room, error)
	GetRoomByID(ctx context.Context, db database.DBRunner, roomID uuid.UUID) (*models.Room, error)

	DoesRoomExists(ctx context.Context, db database.DBRunner, roomID uuid.UUID) (bool, error)
	// new
	GetRoomsByHallID(ctx context.Context, db database.DBRunner, hallID uuid.UUID) ([]*models.Room, error)
	GetRoomsIDandPrivateInfoByHallID(ctx context.Context, db database.DBRunner, hallID uuid.UUID) ([]*dto.RoomIDandPrivate, error)
	UpdateRoom(ctx context.Context, db database.DBRunner, roomID uuid.UUID, name *string, isPrivate *bool) (*models.Room, error)
	DeleteRoom(ctx context.Context, db database.DBRunner, roomID uuid.UUID) error
	GetMaxPositionInContainer(ctx context.Context, db database.DBRunner, hallID uuid.UUID, floorID *uuid.UUID) (float64, error)
	ReorderRooms(ctx context.Context, db database.DBRunner, hallID uuid.UUID, floorID *uuid.UUID, orderedIDs []uuid.UUID) error
	MoveRoom(ctx context.Context, db database.DBRunner, roomID uuid.UUID, newFloorID *uuid.UUID, newPosition float64) (*models.Room, error)

	IsUserRoomMember(ctx context.Context, db database.DBRunner, roomID uuid.UUID, userID uuid.UUID) (bool, error)
	AddRoomMember(ctx context.Context, db database.DBRunner, roomID uuid.UUID, userID uuid.UUID) error
	SyncRoomMembersFromFloor(ctx context.Context, db database.DBRunner, roomID uuid.UUID, floorID uuid.UUID) error
	ClearRoomMembers(ctx context.Context, db database.DBRunner, roomID uuid.UUID) error

	GetRoomPositionBounds(ctx context.Context, db database.DBRunner, hallID uuid.UUID, floorID *uuid.UUID, afterID *uuid.UUID) (lower float64, upper *float64, err error)
}

type roomRepository struct{}

func NewRoomRepository() IRoomRepository {
	return &roomRepository{}
}

func (r *roomRepository) CreateRoom(ctx context.Context, db database.DBRunner, room *models.Room) (*models.Room, error) {
	query := `
        INSERT INTO rooms (id, hall_id, floor_id, name, room_type, position, is_private, created_at, updated_at)
        VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
        RETURNING id, hall_id, floor_id, name, room_type, position, is_private, created_at, updated_at
    `
	out := &models.Room{}
	err := db.QueryRow(ctx, query,
		room.ID, room.HallID, room.FloorID, room.Name,
		room.RoomType, room.Position, room.IsPrivate,
		room.CreatedAt, room.UpdatedAt,
	).Scan(&out.ID, &out.HallID, &out.FloorID, &out.Name,
		&out.RoomType, &out.Position, &out.IsPrivate,
		&out.CreatedAt, &out.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (r *roomRepository) GetRoomByID(ctx context.Context, db database.DBRunner, roomID uuid.UUID) (*models.Room, error) {
	query := `
        SELECT id, hall_id, floor_id, name, room_type, position, is_private, created_at, updated_at
        FROM rooms WHERE id = $1
    `
	out := &models.Room{}
	err := db.QueryRow(ctx, query, roomID).Scan(
		&out.ID, &out.HallID, &out.FloorID, &out.Name,
		&out.RoomType, &out.Position, &out.IsPrivate,
		&out.CreatedAt, &out.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// GetRoomsByHallID returns all rooms for a hall.
// Ordered: NULL floor_id first (top-level), then by floor_id, then by position within each group.
// The service layer groups these into the structured response.
func (r *roomRepository) GetRoomsByHallID(ctx context.Context, db database.DBRunner, hallID uuid.UUID) ([]*models.Room, error) {
	query := `
        SELECT id, hall_id, floor_id, name, room_type, position, is_private, created_at, updated_at
        FROM rooms
        WHERE hall_id = $1
        ORDER BY
            floor_id IS NOT NULL,  -- NULLs (top-level) first
            floor_id,
            position ASC
    `
	rows, err := db.Query(ctx, query, hallID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var rooms []*models.Room
	for rows.Next() {
		rm := &models.Room{}
		if err := rows.Scan(&rm.ID, &rm.HallID, &rm.FloorID, &rm.Name,
			&rm.RoomType, &rm.Position, &rm.IsPrivate,
			&rm.CreatedAt, &rm.UpdatedAt); err != nil {
			return nil, err
		}
		rooms = append(rooms, rm)
	}
	return rooms, rows.Err()
}

func (r *roomRepository) GetRoomsIDandPrivateInfoByHallID(ctx context.Context, db database.DBRunner, hallID uuid.UUID) ([]*dto.RoomIDandPrivate, error) {
	query := `
        SELECT id, is_private
        FROM rooms
        WHERE hall_id = $1
        ORDER BY
            floor_id IS NOT NULL,  -- NULLs (top-level) first
            floor_id,
            position ASC
    `
	rows, err := db.Query(ctx, query, hallID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var roomIDandPrivate []*dto.RoomIDandPrivate

	for rows.Next() {
		rm := &dto.RoomIDandPrivate{}
		if err := rows.Scan(&rm.RoomID, &rm.IsPrivate); err != nil {
			return nil, err
		}
		roomIDandPrivate = append(roomIDandPrivate, rm)
	}
	return nil, rows.Err()
}

// func (r *roomRepository) UpdateRoom(ctx context.Context, db database.DBRunner, roomID uuid.UUID, name *string, isPrivate *bool) (*models.Room, error) {
// 	query := `
//         UPDATE rooms
//         SET
//             name       = COALESCE($1, name),
//             is_private = COALESCE($2, is_private),
//             updated_at = $3
//         WHERE id = $4
//         RETURNING id, hall_id, floor_id, name, room_type, position, is_private, created_at, updated_at
//     `
// 	out := &models.Room{}
// 	err := db.QueryRow(ctx, query, name, isPrivate, time.Now(), roomID).Scan(
// 		&out.ID, &out.HallID, &out.FloorID, &out.Name,
// 		&out.RoomType, &out.Position, &out.IsPrivate,
// 		&out.CreatedAt, &out.UpdatedAt,
// 	)
// 	if err != nil {
// 		return nil, err
// 	}
// 	return out, nil
// }

func (r *roomRepository) UpdateRoom(ctx context.Context, db database.DBRunner, roomID uuid.UUID, fields map[string]any) (*models.Room, error) {
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

	args = append(args, roomID)

	query := fmt.Sprintf(`
		UPDATE rooms
		SET %s
		WHERE id = $%d
		RETURNING id, hall_id, floor_id, name, room_type, position, is_private, created_at, updated_at
	`, strings.Join(setClauses, ", "), i)

	out := &models.Room{}
	err := db.QueryRow(ctx, query, args...).Scan(
		&out.ID,
		&out.HallID,
		&out.FloorID,
		&out.Name,
		&out.RoomType,
		&out.Position,
		&out.IsPrivate,
		&out.CreatedAt,
		&out.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return out, nil
}

func (r *roomRepository) DeleteRoom(ctx context.Context, db database.DBRunner, roomID uuid.UUID) error {
	_, err := db.Exec(ctx, `DELETE FROM rooms WHERE id = $1`, roomID)
	return err
}

// GetMaxPositionInContainer returns the highest position value among rooms
// in the same container (floor_id = nil means top-level).
func (r *roomRepository) GetMaxPositionInContainer(ctx context.Context, db database.DBRunner, hallID uuid.UUID, floorID *uuid.UUID) (float64, error) {
	query := `
        SELECT COALESCE(MAX(position), 0)
        FROM rooms
        WHERE hall_id = $1
          AND ($2::uuid IS NULL AND floor_id IS NULL OR floor_id = $2)
    `
	var max float64
	if err := db.QueryRow(ctx, query, hallID, floorID).Scan(&max); err != nil {
		return 0, err
	}
	return max, nil
}

// ReorderRooms reassigns positions 1000, 2000, 3000 … within a container.
// floorID = nil targets the top-level (floor_id IS NULL) group.
func (r *roomRepository) ReorderRooms(ctx context.Context, db database.DBRunner, hallID uuid.UUID, floorID *uuid.UUID, orderedIDs []uuid.UUID) error {
	ids := make([]string, len(orderedIDs))
	positions := make([]float64, len(orderedIDs))
	for i, id := range orderedIDs {
		ids[i] = id.String()
		positions[i] = float64((i + 1) * 1000)
	}
	query := `
        UPDATE rooms
        SET    position   = new_order.pos,
               updated_at = now()
        FROM (
            SELECT UNNEST($1::uuid[])   AS id,
                   UNNEST($2::float8[]) AS pos
        ) AS new_order
        WHERE rooms.id      = new_order.id
          AND rooms.hall_id = $3
          AND ($4::uuid IS NULL AND rooms.floor_id IS NULL OR rooms.floor_id = $4)
    `
	_, err := db.Exec(ctx, query, ids, positions, hallID, floorID)
	return err
}

// MoveRoom updates the room's container and position in one query.
func (r *roomRepository) MoveRoom(ctx context.Context, db database.DBRunner, roomID uuid.UUID, newFloorID *uuid.UUID, newPosition float64) (*models.Room, error) {
	query := `
        UPDATE rooms
        SET    floor_id   = $1,
               position   = $2,
               updated_at = now()
        WHERE  id = $3
        RETURNING id, hall_id, floor_id, name, room_type, position, is_private, created_at, updated_at
    `
	out := &models.Room{}
	err := db.QueryRow(ctx, query, newFloorID, newPosition, roomID).Scan(
		&out.ID, &out.HallID, &out.FloorID, &out.Name,
		&out.RoomType, &out.Position, &out.IsPrivate,
		&out.CreatedAt, &out.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (r *roomRepository) IsUserRoomMember(ctx context.Context, db database.DBRunner, roomID uuid.UUID, userID uuid.UUID) (bool, error) {
	var exists bool
	err := db.QueryRow(ctx,
		`SELECT EXISTS (SELECT 1 FROM room_members WHERE room_id = $1 AND user_id = $2)`,
		roomID, userID,
	).Scan(&exists)
	return exists, err
}

func (r *roomRepository) AddRoomMember(ctx context.Context, db database.DBRunner, roomID uuid.UUID, userID uuid.UUID) error {

	query := `
		INSERT INTO room_members (room_id, user_id)
		VALUES ($1, $2)
		ON CONFLICT (room_id, user_id) DO NOTHING
	`

	_, err := db.Exec(ctx, query, roomID, userID)
	return err
}

func (r *roomRepository) SyncRoomMembersFromFloor(ctx context.Context, db database.DBRunner, roomID uuid.UUID, floorID uuid.UUID) error {
	query := `
		INSERT INTO room_members (room_id, user_id)
		SELECT $1, fm.user_id
		FROM floor_members fm
		WHERE fm.floor_id = $2
		ON CONFLICT (room_id, user_id) DO NOTHING
	`

	_, err := db.Exec(ctx, query, roomID, floorID)
	return err
}

func (r *roomRepository) ClearRoomMembers(ctx context.Context, db database.DBRunner, roomID uuid.UUID) error {
	query := `
		DELETE FROM room_members
		WHERE room_id = $1
	`

	_, err := db.Exec(ctx, query, roomID)
	return err
}

func (r *roomRepository) DoesRoomExists(ctx context.Context, db database.DBRunner, roomID uuid.UUID) (bool, error) {
	var exists bool
	err := db.QueryRow(ctx,
		`SELECT EXISTS (SELECT 1 FROM rooms WHERE id = $1)`,
		roomID,
	).Scan(&exists)
	return exists, err
}

// GetRoomPositionBounds returns position bounds for placing a room within a container.
//
// Container is identified by (hallID, floorID):
//
//	floorID = nil  → top-level container
//	floorID = uuid → inside that floor
//
//	afterID = nil  → lower = 0, upper = current minimum in container (or nil)
//	afterID = uuid → lower = that room's position, upper = next room's position (or nil)
func (r *roomRepository) GetRoomPositionBounds(
	ctx context.Context,
	db database.DBRunner,
	hallID uuid.UUID,
	floorID *uuid.UUID,
	afterID *uuid.UUID,
) (lower float64, upper *float64, err error) {
	if afterID == nil {
		query := `
            SELECT MIN(position) FROM rooms
            WHERE hall_id = $1
              AND ($2::uuid IS NULL AND floor_id IS NULL OR floor_id = $2)
        `
		var min *float64
		if err := db.QueryRow(ctx, query, hallID, floorID).Scan(&min); err != nil {
			return 0, nil, err
		}
		return 0, min, nil
	}

	query := `
        WITH anchor AS (
            SELECT position FROM rooms
            WHERE id = $1 AND hall_id = $2
        )
        SELECT
            anchor.position,
            (
                SELECT position FROM rooms
                WHERE  hall_id  = $2
                  AND  ($3::uuid IS NULL AND floor_id IS NULL OR floor_id = $3)
                  AND  position > anchor.position
                ORDER  BY position ASC
                LIMIT  1
            )
        FROM anchor
    `
	var up *float64
	if err := db.QueryRow(ctx, query, afterID, hallID, floorID).Scan(&lower, &up); err != nil {
		return 0, nil, err
	}
	return lower, up, nil
}
