package repositories

import (
	"context"

	"github.com/google/uuid"
	"github.com/suck-seed/yapp/internal/database"
	"github.com/suck-seed/yapp/internal/models"
	"github.com/suck-seed/yapp/internal/utils"
)

type IBanRepsitory interface {

	// ------------------------------- BANS
	// core cud
	BanUser(ctx context.Context, db database.DBRunner, ban *models.HallBan) (*models.HallBan, error)
	UnBanUser(ctx context.Context, db database.DBRunner, banID uuid.UUID) (*models.HallBan, error)
	GetBanByID(ctx context.Context, db database.DBRunner, banID uuid.UUID) (*models.HallBan, error)

	// list operation
	GetAllHallBans(ctx context.Context, db database.DBRunner, hallID uuid.UUID) ([]*models.HallBan, error)
	GetAllHallBansPaginated(ctx context.Context, db database.DBRunner, hallID uuid.UUID, offset int, limit int) ([]*models.HallBan, error)

	// check operation
	IsUserBanned(ctx context.Context, db database.DBRunner, hallID uuid.UUID, userID uuid.UUID) (bool, error)
	GetBanByUserAndHall(ctx context.Context, db database.DBRunner, hallID uuid.UUID, userID uuid.UUID) (*models.HallBan, error)

	// Statistics
	GetBanCount(ctx context.Context, db database.DBRunner, hallID uuid.UUID) (int, error)
}

type banRepository struct{}

func NewBanRepository() IBanRepsitory {
	return &banRepository{}
}

// ------------------------------ BANS
func (r *banRepository) BanUser(ctx context.Context, db database.DBRunner, ban *models.HallBan) (*models.HallBan, error) {

	query := `

	INSERT INTO hall_bans (id, reason, user_id, hall_id)
    VALUES ($1, $2, $3, $4)
    RETURNING id, reason, user_id, hall_id, created_at, updated_at
	`

	row := db.QueryRow(ctx, query, ban.ID, ban.Reason, ban.UserID, ban.HallID)

	saved := &models.HallBan{}
	err := row.Scan(
		&saved.ID,
		&saved.Reason,
		&saved.UserID,
		&saved.HallID,
		&saved.CreatedAt,
		&saved.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return saved, nil
}

func (r *banRepository) UnBanUser(ctx context.Context, db database.DBRunner, banID uuid.UUID) (*models.HallBan, error) {

	// Fetching ban before deleting it
	ban, err := r.GetBanByID(ctx, db, banID)
	if err != nil {
		return nil, utils.ErrorFetchingBan
	}

	// delete the ban from registery
	query := `	DELETE
			  	FROM hall_bans
				WHERE id = $1
			`
	result, err := db.Exec(ctx, query, banID)
	if err != nil {
		return nil, err
	}

	rowsAffected := result.RowsAffected()
	if rowsAffected == 0 {
		return nil, utils.ErrorBanNotFound
	}

	return ban, nil
}

func (r *banRepository) GetBanByID(ctx context.Context, db database.DBRunner, banID uuid.UUID) (*models.HallBan, error) {

	query := `
	SELECT id, reason, user_id, hall_id, created_at, updated_at
	FROM hall_bans
	WHERE id = $
	`

	row := db.QueryRow(ctx, query, banID)

	saved := &models.HallBan{}
	err := row.Scan(
		&saved.ID,
		&saved.Reason,
		&saved.UserID,
		&saved.HallID,
		&saved.CreatedAt,
		&saved.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return saved, nil
}

func (r *banRepository) GetAllHallBans(ctx context.Context, db database.DBRunner, hallID uuid.UUID) ([]*models.HallBan, error) {

	query := `	SELECT id, reason, user_id, hall_id, created_at, updated_at
				FROM hall_bans
				WHERE hall_id = $1
				ORDER BY created_at DESC
	`

	rows, err := db.Query(ctx, query, hallID)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	bans := []*models.HallBan{}
	for rows.Next() {
		currentBan := &models.HallBan{}
		err := rows.Scan(
			&currentBan.ID,
			&currentBan.Reason,
			&currentBan.UserID,
			&currentBan.HallID,
			&currentBan.CreatedAt,
			&currentBan.UpdatedAt,
		)

		if err != nil {
			return nil, err
		}

		bans = append(bans, currentBan)
	}

	// check if error iterating bans
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return bans, nil
}

func (r *banRepository) GetAllHallBansPaginated(ctx context.Context, db database.DBRunner, hallID uuid.UUID, offset int, limit int) ([]*models.HallBan, error)

// check operation
func (r *banRepository) IsUserBanned(ctx context.Context, db database.DBRunner, hallID uuid.UUID, userID uuid.UUID) (bool, error) {
	query := `	SELECT EXISTS (
					SELECT 1
					FROM hall_bans
					WHERE hall_id = $1 AND user_id = $2
				)
	`
	var exists bool
	row := db.QueryRow(ctx, query, hallID, userID)

	if err := row.Scan(&exists); err != nil {
		return false, err
	}
	return exists, nil
}

func (r *banRepository) GetBanByUserAndHall(ctx context.Context, db database.DBRunner, hallID uuid.UUID, userID uuid.UUID) (*models.HallBan, error) {
	query := `
				SELECT id, reason, user_id, hall_id, created_at, updated_at
				FROM hall_bans
				WHERE hall_id = $1 AND user_id = $2

	`
	saved := &models.HallBan{}
	row := db.QueryRow(ctx, query, hallID, userID)

	err := row.Scan(
		&saved.ID,
		&saved.Reason,
		&saved.UserID,
		&saved.HallID,
		&saved.CreatedAt,
		&saved.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return saved, nil
}

// Statistics
func (r *banRepository) GetBanCount(ctx context.Context, db database.DBRunner, hallID uuid.UUID) (int, error) {

	query := `
				SELECT COUNT(*)
				FROM hall_bans
				WHERE hall_id = $1
	`

	var count int
	row := db.QueryRow(ctx, query, hallID)

	if err := row.Scan(&count); err != nil {
		return 0, err
	}
	return count, nil
}
