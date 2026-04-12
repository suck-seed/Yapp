package repositories

import (
	"context"

	"github.com/google/uuid"
	"github.com/suck-seed/yapp/internal/database"
	"github.com/suck-seed/yapp/internal/models"
)

type IInviteRepository interface {
	CreateInvite(ctx context.Context, db database.DBRunner, invite *models.HallInvite) (*models.HallInvite, error)
	GetInviteByID(ctx context.Context, db database.DBRunner, inviteID uuid.UUID) (*models.HallInvite, error)
	GetInviteByCode(ctx context.Context, db database.DBRunner, code string) (*models.HallInvite, error)
	ListHallInvites(ctx context.Context, db database.DBRunner, hallID uuid.UUID) ([]*models.HallInvite, error)
	AtomicIncrementUsedCount(ctx context.Context, db database.DBRunner, inviteID uuid.UUID) (*models.HallInvite, error)
	DeleteInvite(ctx context.Context, db database.DBRunner, inviteID uuid.UUID) (*models.HallInvite, error)
}

type inviteRepository struct{}

func NewInviteRepository() IInviteRepository {
	return &inviteRepository{}
}

func scanInvite(row interface{ Scan(...any) error }, out *models.HallInvite) error {
	return row.Scan(
		&out.ID, &out.HallID, &out.CreatedBy, &out.Code,
		&out.RoleID, &out.MaxUses, &out.UsedCount, &out.ExpiresAt, &out.CreatedAt,
	)
}

func (r *inviteRepository) CreateInvite(ctx context.Context, db database.DBRunner, invite *models.HallInvite) (*models.HallInvite, error) {
	query := `
		INSERT INTO hall_invites (hall_id, created_by, code, role_id, max_uses, expires_at)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id, hall_id, created_by, code, role_id, max_uses, used_count, expires_at, created_at`

	out := &models.HallInvite{}
	if err := scanInvite(db.QueryRow(ctx, query,
		invite.HallID, invite.CreatedBy, invite.Code,
		invite.RoleID, invite.MaxUses, invite.ExpiresAt,
	), out); err != nil {
		return nil, err
	}
	return out, nil
}

func (r *inviteRepository) GetInviteByID(ctx context.Context, db database.DBRunner, inviteID uuid.UUID) (*models.HallInvite, error) {
	query := `
		SELECT id, hall_id, created_by, code, role_id, max_uses, used_count, expires_at, created_at
		FROM hall_invites
		WHERE id = $1`

	out := &models.HallInvite{}
	if err := scanInvite(db.QueryRow(ctx, query, inviteID), out); err != nil {
		return nil, err
	}
	return out, nil
}

func (r *inviteRepository) GetInviteByCode(ctx context.Context, db database.DBRunner, code string) (*models.HallInvite, error) {
	query := `
		SELECT id, hall_id, created_by, code, role_id, max_uses, used_count, expires_at, created_at
		FROM hall_invites
		WHERE code = $1`

	out := &models.HallInvite{}
	if err := scanInvite(db.QueryRow(ctx, query, code), out); err != nil {
		return nil, err
	}
	return out, nil
}

func (r *inviteRepository) ListHallInvites(ctx context.Context, db database.DBRunner, hallID uuid.UUID) ([]*models.HallInvite, error) {
	query := `
		SELECT id, hall_id, created_by, code, role_id, max_uses, used_count, expires_at, created_at
		FROM hall_invites
		WHERE hall_id = $1
		ORDER BY created_at DESC`

	rows, err := db.Query(ctx, query, hallID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var invites []*models.HallInvite
	for rows.Next() {
		inv := &models.HallInvite{}
		if err := scanInvite(rows, inv); err != nil {
			return nil, err
		}
		invites = append(invites, inv)
	}
	return invites, rows.Err()
}

func (r *inviteRepository) AtomicIncrementUsedCount(ctx context.Context, db database.DBRunner, inviteID uuid.UUID) (*models.HallInvite, error) {
	query := `
		UPDATE hall_invites
		SET used_count = used_count + 1
		WHERE id = $1
		  AND (max_uses IS NULL OR used_count < max_uses)
		RETURNING id, hall_id, created_by, code, role_id, max_uses, used_count, expires_at, created_at`

	out := &models.HallInvite{}
	if err := scanInvite(db.QueryRow(ctx, query, inviteID), out); err != nil {
		return nil, err // pgx.ErrNoRows → cap was hit, service maps this
	}
	return out, nil
}

func (r *inviteRepository) DeleteInvite(ctx context.Context, db database.DBRunner, inviteID uuid.UUID) (*models.HallInvite, error) {
	query := `
		DELETE FROM hall_invites
		WHERE id = $1
		RETURNING id, hall_id, created_by, code, role_id, max_uses, used_count, expires_at, created_at`

	out := &models.HallInvite{}
	if err := scanInvite(db.QueryRow(ctx, query, inviteID), out); err != nil {
		return nil, err
	}
	return out, nil
}
