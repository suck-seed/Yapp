package repositories

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/suck-seed/yapp/internal/database"
	"github.com/suck-seed/yapp/internal/models"
)

type IFloorRepository interface {
	CreateFloor(ctx context.Context, db database.DBRunner, floor *models.Floor) (*models.Floor, error)
	GetFloorByID(ctx context.Context, db database.DBRunner, floorID uuid.UUID) (*models.Floor, error)
	GetFloorsByHallID(ctx context.Context, db database.DBRunner, hallID uuid.UUID) ([]*models.Floor, error)
	UpdateFloor(ctx context.Context, db database.DBRunner, floorID uuid.UUID, name *string, isPrivate *bool) (*models.Floor, error)
	DeleteFloor(ctx context.Context, db database.DBRunner, floorID uuid.UUID) error
	GetMaxPosition(ctx context.Context, db database.DBRunner, hallID uuid.UUID) (float64, error)
	ReorderFloors(ctx context.Context, db database.DBRunner, hallID uuid.UUID, orderedIDs []uuid.UUID) error
	DoesFloorExistInHall(ctx context.Context, db database.DBRunner, floorID uuid.UUID, hallID uuid.UUID) (bool, error)

	GetFloorPositionBounds(ctx context.Context, db database.DBRunner, hallID uuid.UUID, afterID *uuid.UUID) (lower float64, upper *float64, err error)
	UpdateFloorPosition(ctx context.Context, db database.DBRunner, floorID uuid.UUID, position float64) (*models.Floor, error)

	// Helper functions
	IsFloorPrivate(ctx context.Context, db database.DBRunner, floorID uuid.UUID) (bool, error)
}

type floorRepository struct{}

func NewFloorRepository() IFloorRepository {
	return &floorRepository{}
}

func (r *floorRepository) CreateFloor(ctx context.Context, db database.DBRunner, floor *models.Floor) (*models.Floor, error) {
	query := `
		INSERT INTO floors (id, hall_id, name, position, is_private, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id, hall_id, name, position, is_private, created_at, updated_at
	`
	row := db.QueryRow(ctx, query,
		floor.ID, floor.HallID, floor.Name,
		floor.Position, floor.IsPrivate,
		floor.CreatedAt, floor.UpdatedAt,
	)
	out := &models.Floor{}
	err := row.Scan(&out.ID, &out.HallID, &out.Name, &out.Position, &out.IsPrivate, &out.CreatedAt, &out.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (r *floorRepository) GetFloorByID(ctx context.Context, db database.DBRunner, floorID uuid.UUID) (*models.Floor, error) {
	query := `
		SELECT id, hall_id, name, position, is_private, created_at, updated_at
		FROM floors
		WHERE id = $1
	`
	out := &models.Floor{}
	err := db.QueryRow(ctx, query, floorID).Scan(
		&out.ID, &out.HallID, &out.Name, &out.Position,
		&out.IsPrivate, &out.CreatedAt, &out.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (r *floorRepository) GetFloorsByHallID(ctx context.Context, db database.DBRunner, hallID uuid.UUID) ([]*models.Floor, error) {
	query := `
		SELECT id, hall_id, name, position, is_private, created_at, updated_at
		FROM floors
		WHERE hall_id = $1
		ORDER BY position ASC
	`
	rows, err := db.Query(ctx, query, hallID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var floors []*models.Floor
	for rows.Next() {
		f := &models.Floor{}
		if err := rows.Scan(&f.ID, &f.HallID, &f.Name, &f.Position, &f.IsPrivate, &f.CreatedAt, &f.UpdatedAt); err != nil {
			return nil, err
		}
		floors = append(floors, f)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return floors, nil
}

// UpdateFloor only sets columns that are non-nil. Caller guarantees at least one is non-nil.
func (r *floorRepository) UpdateFloor(ctx context.Context, db database.DBRunner, floorID uuid.UUID, name *string, isPrivate *bool) (*models.Floor, error) {
	query := `
		UPDATE floors
		SET
			name       = COALESCE($1, name),
			is_private = COALESCE($2, is_private),
			updated_at = $3
		WHERE id = $4
		RETURNING id, hall_id, name, position, is_private, created_at, updated_at
	`
	out := &models.Floor{}
	err := db.QueryRow(ctx, query, name, isPrivate, time.Now(), floorID).Scan(
		&out.ID, &out.HallID, &out.Name, &out.Position,
		&out.IsPrivate, &out.CreatedAt, &out.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (r *floorRepository) DeleteFloor(ctx context.Context, db database.DBRunner, floorID uuid.UUID) error {
	query := `DELETE FROM floors WHERE id = $1`
	_, err := db.Exec(ctx, query, floorID)
	return err
}

func (r *floorRepository) GetMaxPosition(ctx context.Context, db database.DBRunner, hallID uuid.UUID) (float64, error) {
	query := `SELECT COALESCE(MAX(position), 0) FROM floors WHERE hall_id = $1`
	var max float64
	if err := db.QueryRow(ctx, query, hallID).Scan(&max); err != nil {
		return 0, err
	}
	return max, nil
}

// ReorderFloors reassigns positions as 1000, 2000, 3000 … for each ID in orderedIDs.
func (r *floorRepository) ReorderFloors(ctx context.Context, db database.DBRunner, hallID uuid.UUID, orderedIDs []uuid.UUID) error {
	ids := make([]string, len(orderedIDs))
	positions := make([]float64, len(orderedIDs))
	for i, id := range orderedIDs {
		ids[i] = id.String()
		positions[i] = float64((i + 1) * 1000)
	}
	query := `
		UPDATE floors
		SET    position   = new_order.pos,
		       updated_at = now()
		FROM (
			SELECT UNNEST($1::uuid[])   AS id,
			       UNNEST($2::float8[]) AS pos
		) AS new_order
		WHERE floors.id      = new_order.id
		  AND floors.hall_id = $3
	`
	_, err := db.Exec(ctx, query, ids, positions, hallID)
	return err
}

func (r *floorRepository) DoesFloorExistInHall(ctx context.Context, db database.DBRunner, floorID uuid.UUID, hallID uuid.UUID) (bool, error) {
	query := `SELECT EXISTS (SELECT 1 FROM floors WHERE id = $1 AND hall_id = $2)`
	var exists bool
	if err := db.QueryRow(ctx, query, floorID, hallID).Scan(&exists); err != nil {
		return false, err
	}
	return exists, nil
}

// GetFloorPositionBounds returns the position of afterID and the position
// of the next floor after it within the hall.
//
//	afterID = nil  → lower = 0, upper = position of first floor (or nil)
//	afterID = uuid → lower = that floor's position, upper = next floor's position (or nil)
//
// If upper is nil, caller should use lower + 1000 (append at end).
func (r *floorRepository) GetFloorPositionBounds(
	ctx context.Context,
	db database.DBRunner,
	hallID uuid.UUID,
	afterID *uuid.UUID,
) (lower float64, upper *float64, err error) {
	if afterID == nil {
		// inserting at top: lower = 0, upper = current minimum position
		query := `
            SELECT MIN(position) FROM floors WHERE hall_id = $1
        `
		var min *float64
		if err := db.QueryRow(ctx, query, hallID).Scan(&min); err != nil {
			return 0, nil, err
		}
		return 0, min, nil
	}

	// lower = afterID's position
	// upper = smallest position in the hall that is strictly greater than lower
	query := `
        WITH anchor AS (
            SELECT position FROM floors
            WHERE id = $1 AND hall_id = $2
        )
        SELECT
            anchor.position,
            (
                SELECT position FROM floors
                WHERE  hall_id  = $2
                  AND  position > anchor.position
                ORDER  BY position ASC
                LIMIT  1
            )
        FROM anchor
    `
	var up *float64
	if err := db.QueryRow(ctx, query, afterID, hallID).Scan(&lower, &up); err != nil {
		return 0, nil, err
	}
	return lower, up, nil
}

func (r *floorRepository) UpdateFloorPosition(ctx context.Context, db database.DBRunner, floorID uuid.UUID, position float64) (*models.Floor, error) {
	query := `
        UPDATE floors SET position = $1, updated_at = now()
        WHERE id = $2
        RETURNING id, hall_id, name, position, is_private, created_at, updated_at
    `
	out := &models.Floor{}
	err := db.QueryRow(ctx, query, position, floorID).Scan(
		&out.ID, &out.HallID, &out.Name, &out.Position,
		&out.IsPrivate, &out.CreatedAt, &out.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// HELPER FUNCTIONS
func (r *floorRepository) IsFloorPrivate(ctx context.Context, db database.DBRunner, floorID uuid.UUID) (bool, error) {
	query := `
		SELECT is_private
		FROM floors
		WHERE id = $1
	`

	var isPrivate bool

	err := db.QueryRow(ctx, query, floorID).Scan(&isPrivate)

	if err != nil {
		return false, err
	}

	return isPrivate, nil
}
