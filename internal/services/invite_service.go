package services

import (
	"context"
	"crypto/rand"
	"errors"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/suck-seed/yapp/internal/auth"
	"github.com/suck-seed/yapp/internal/database"
	dto "github.com/suck-seed/yapp/internal/dto/hall"
	"github.com/suck-seed/yapp/internal/models"
	"github.com/suck-seed/yapp/internal/repositories"
	"github.com/suck-seed/yapp/internal/utils"
)

const (
	inviteCodeLength = 10
	baseInviteURL    = "https://yourdomain.com/invites/" // move to config
)

var expireAfterDurations = map[dto.ExpireAfterOption]time.Duration{
	dto.Expire30Min: 30 * time.Minute,
	dto.Expire1Hr:   time.Hour,
	dto.Expire6Hr:   6 * time.Hour,
	dto.Expire12Hr:  12 * time.Hour,
	dto.Expire1Day:  24 * time.Hour,
	dto.Expire7Days: 7 * 24 * time.Hour,
}

type IInviteService interface {
	// Management (requires manage_invites permission)
	CreateInviteLink(ctx context.Context, userInfo *auth.UserInfo, hallID uuid.UUID, req *dto.CreateInviteLinkReq) (*dto.InviteLinkRes, error)
	ListInviteLinks(ctx context.Context, userInfo *auth.UserInfo, hallID uuid.UUID) ([]*dto.InviteLinkRes, error)
	RevokeInviteLink(ctx context.Context, userInfo *auth.UserInfo, hallID uuid.UUID, inviteID uuid.UUID) (*dto.InviteLinkRes, error)
	// Public / joining
	GetInviteLinkInfo(ctx context.Context, code string) (*dto.InviteInfoRes, error)
	AcceptInviteLink(ctx context.Context, userInfo *auth.UserInfo, code string) (*dto.AcceptInviteLinkRes, error)
}

type inviteService struct {
	repositories.IInviteRepository
	repositories.IHallRepository
	repositories.IRoleRepository

	IPermissionCheckerService

	pool    *pgxpool.Pool
	timeout time.Duration
	mu      sync.RWMutex
}

func NewInviteService(
	inviteRepo repositories.IInviteRepository,
	hallRepo repositories.IHallRepository,
	roleRepo repositories.IRoleRepository,
	permSvc IPermissionCheckerService,
	pool *pgxpool.Pool,
) IInviteService {
	return &inviteService{
		IInviteRepository:         inviteRepo,
		IHallRepository:           hallRepo,
		IRoleRepository:           roleRepo,
		IPermissionCheckerService: permSvc,
		pool:                      pool,
		timeout:                   2 * time.Second,
		mu:                        sync.RWMutex{},
	}
}

// ---------- helpers ----------

func generateInviteCode(length int) (string, error) {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, length)
	if _, err := rand.Read(b); err != nil {
		return "", utils.ErrorGeneratingInviteCode
	}
	for i := range b {
		b[i] = charset[int(b[i])%len(charset)]
	}
	return string(b), nil
}

func inviteToRes(inv *models.HallInvite) *dto.InviteLinkRes {
	return &dto.InviteLinkRes{
		ID:        inv.ID,
		HallID:    inv.HallID,
		CreatedBy: inv.CreatedBy,
		Code:      inv.Code,
		URL:       baseInviteURL + inv.Code,
		RoleID:    inv.RoleID,
		MaxUses:   inv.MaxUses,
		UsedCount: inv.UsedCount,
		ExpiresAt: inv.ExpiresAt,
		CreatedAt: inv.CreatedAt,
		IsValid:   inv.IsValid(),
	}
}

// ---------- management ----------

func (s *inviteService) CreateInviteLink(ctx context.Context, userInfo *auth.UserInfo, hallID uuid.UUID, req *dto.CreateInviteLinkReq) (*dto.InviteLinkRes, error) {
	ctx, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()

	// 1. validate expire_after
	if !dto.ValidExpireOptions[req.ExpireAfter] {
		return nil, utils.ErrorInvalidExpireAfter
	}

	// 2. validate max_uses — nil means unlimited, otherwise must be an allowed preset
	if req.MaxUses != nil && !dto.ValidMaxUses[*req.MaxUses] {
		return nil, utils.ErrorInvalidMaxUses
	}

	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return nil, utils.ErrorInternal
	}
	runner := database.NewTxWrapper(tx)
	defer runner.Rollback(ctx)

	// 3. permission check
	canManage, err := s.CanManageInvites(ctx, runner, userInfo.ID, hallID)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) || errors.Is(err, context.Canceled) {
			return nil, utils.ErrorRequestTimeout
		}
		return nil, utils.ErrorInternal
	}
	if !canManage {
		return nil, utils.ErrorUserCannotManageInvites
	}

	// 4. validate role belongs to this hall (if provided)
	if req.RoleID != nil {
		exists, err := s.IRoleRepository.DoesRoleExist(ctx, runner, *req.RoleID, hallID)
		if err != nil {
			if errors.Is(err, context.DeadlineExceeded) || errors.Is(err, context.Canceled) {
				return nil, utils.ErrorRequestTimeout
			}
			return nil, utils.ErrorFetchingRole
		}
		if !exists {
			return nil, utils.ErrorRoleDoesntBelongInThisHall
		}
	}

	// 5. build expiry timestamp
	var expiresAt *time.Time
	if req.ExpireAfter != dto.ExpireNever {
		t := time.Now().Add(expireAfterDurations[req.ExpireAfter])
		expiresAt = &t
	}

	// 6. generate unique code
	code, err := generateInviteCode(inviteCodeLength)
	if err != nil {
		return nil, err // already utils.ErrorGeneratingInviteCode
	}

	// 7. persist
	inv, err := s.IInviteRepository.CreateInvite(ctx, runner, &models.HallInvite{
		HallID:    hallID,
		CreatedBy: userInfo.ID,
		Code:      code,
		RoleID:    req.RoleID,
		MaxUses:   req.MaxUses,
		ExpiresAt: expiresAt,
	})
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) || errors.Is(err, context.Canceled) {
			return nil, utils.ErrorRequestTimeout
		}
		return nil, utils.ErrorCreatingInvite
	}

	if err := runner.Commit(ctx); err != nil {
		return nil, utils.ErrorInternal
	}

	return inviteToRes(inv), nil
}

func (s *inviteService) ListInviteLinks(ctx context.Context, userInfo *auth.UserInfo, hallID uuid.UUID) ([]*dto.InviteLinkRes, error) {
	ctx, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()

	conn, err := s.pool.Acquire(ctx)
	if err != nil {
		return nil, utils.ErrorInternal
	}
	defer conn.Release()
	runner := database.NewConnWrapper(conn)

	canManage, err := s.CanManageInvites(ctx, runner, userInfo.ID, hallID)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) || errors.Is(err, context.Canceled) {
			return nil, utils.ErrorRequestTimeout
		}
		return nil, utils.ErrorInternal
	}
	if !canManage {
		return nil, utils.ErrorUserCannotManageInvites
	}

	invites, err := s.IInviteRepository.ListHallInvites(ctx, runner, hallID)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) || errors.Is(err, context.Canceled) {
			return nil, utils.ErrorRequestTimeout
		}
		return nil, utils.ErrorFetchingInvite
	}

	res := make([]*dto.InviteLinkRes, len(invites))
	for i, inv := range invites {
		res[i] = inviteToRes(inv)
	}
	return res, nil
}

func (s *inviteService) RevokeInviteLink(ctx context.Context, userInfo *auth.UserInfo, hallID uuid.UUID, inviteID uuid.UUID) (*dto.InviteLinkRes, error) {
	ctx, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()

	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return nil, utils.ErrorInternal
	}
	runner := database.NewTxWrapper(tx)
	defer runner.Rollback(ctx)

	canManage, err := s.CanManageInvites(ctx, runner, userInfo.ID, hallID)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) || errors.Is(err, context.Canceled) {
			return nil, utils.ErrorRequestTimeout
		}
		return nil, utils.ErrorInternal
	}
	if !canManage {
		return nil, utils.ErrorUserCannotManageInvites
	}

	// ensure invite belongs to this hall before deleting
	existing, err := s.IInviteRepository.GetInviteByID(ctx, runner, inviteID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, utils.ErrorInviteNotFound
		}
		if errors.Is(err, context.DeadlineExceeded) || errors.Is(err, context.Canceled) {
			return nil, utils.ErrorRequestTimeout
		}
		return nil, utils.ErrorFetchingInvite
	}
	if existing.HallID != hallID {
		return nil, utils.ErrorInviteDoesntBelongToHall
	}

	deleted, err := s.IInviteRepository.DeleteInvite(ctx, runner, inviteID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, utils.ErrorInviteNotFound
		}
		if errors.Is(err, context.DeadlineExceeded) || errors.Is(err, context.Canceled) {
			return nil, utils.ErrorRequestTimeout
		}
		return nil, utils.ErrorDeletingInvite
	}

	if err := runner.Commit(ctx); err != nil {
		return nil, utils.ErrorInternal
	}

	return inviteToRes(deleted), nil
}

// ---------- public / joining ----------

func (s *inviteService) GetInviteLinkInfo(ctx context.Context, code string) (*dto.InviteInfoRes, error) {
	ctx, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()

	conn, err := s.pool.Acquire(ctx)
	if err != nil {
		return nil, utils.ErrorInternal
	}
	defer conn.Release()
	runner := database.NewConnWrapper(conn)

	inv, err := s.IInviteRepository.GetInviteByCode(ctx, runner, code)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, utils.ErrorInviteNotFound
		}
		if errors.Is(err, context.DeadlineExceeded) || errors.Is(err, context.Canceled) {
			return nil, utils.ErrorRequestTimeout
		}
		return nil, utils.ErrorFetchingInvite
	}

	hall, err := s.IHallRepository.GetHallByID(ctx, runner, inv.HallID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, utils.ErrorHallNotFound
		}
		if errors.Is(err, context.DeadlineExceeded) || errors.Is(err, context.Canceled) {
			return nil, utils.ErrorRequestTimeout
		}
		return nil, utils.ErrorFetchingHall
	}

	// role name is best-effort — role may have been deleted (ON DELETE SET NULL)
	roleName := ""
	if inv.RoleID != nil {
		role, err := s.IRoleRepository.GetRole(ctx, runner, *inv.RoleID)
		if err == nil {
			roleName = role.Name
		}
	}

	return &dto.InviteInfoRes{
		Code:      inv.Code,
		HallID:    inv.HallID,
		HallName:  hall.Name,
		HallImage: *hall.IconThumbnailURL,
		RoleID:    inv.RoleID,
		RoleName:  roleName,
		MaxUses:   inv.MaxUses,
		UsedCount: inv.UsedCount,
		ExpiresAt: inv.ExpiresAt,
		IsValid:   inv.IsValid(),
	}, nil
}

func (s *inviteService) AcceptInviteLink(ctx context.Context, userInfo *auth.UserInfo, code string) (*dto.AcceptInviteLinkRes, error) {
	ctx, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()

	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return nil, utils.ErrorInternal
	}
	runner := database.NewTxWrapper(tx)
	defer runner.Rollback(ctx)

	inv, err := s.IInviteRepository.GetInviteByCode(ctx, runner, code)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, utils.ErrorInviteNotFound
		}
		if errors.Is(err, context.DeadlineExceeded) || errors.Is(err, context.Canceled) {
			return nil, utils.ErrorRequestTimeout
		}
		return nil, utils.ErrorFetchingInvite
	}

	// validity checks before touching anything
	if inv.IsExpired() {
		return nil, utils.ErrorInviteExpired
	}
	if inv.IsExhausted() {
		return nil, utils.ErrorInviteExhausted
	}

	// already a member?
	isMember, err := s.IHallRepository.IsUserHallMember(ctx, runner, inv.HallID, userInfo.ID)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) || errors.Is(err, context.Canceled) {
			return nil, utils.ErrorRequestTimeout
		}
		return nil, utils.ErrorTest1
	}
	if isMember {
		return nil, utils.ErrorAlreadyHallMember
	}

	// atomic increment — guards against the concurrent-last-slot race
	updated, err := s.IInviteRepository.AtomicIncrementUsedCount(ctx, runner, inv.ID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			// 0 rows matched the WHERE clause — cap was hit between our check and update
			return nil, utils.ErrorInviteExhausted
		}
		if errors.Is(err, context.DeadlineExceeded) || errors.Is(err, context.Canceled) {
			return nil, utils.ErrorRequestTimeout
		}
		return nil, utils.ErrorTest2
	}

	memberID, err := uuid.NewV7()
	if err != nil {
		return nil, utils.ErrorTest3
	}

	member, err := s.IHallRepository.CreateHallMember(ctx, runner, &models.HallMember{
		ID:     memberID,
		HallID: updated.HallID,
		UserID: userInfo.ID,
		RoleID: *updated.RoleID, // nil if no role was set on the invite, repo handles it
	})
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) || errors.Is(err, context.Canceled) {
			return nil, utils.ErrorRequestTimeout
		}
		return nil, utils.ErrorCreatingHallMember
	}

	if err := runner.Commit(ctx); err != nil {
		return nil, utils.ErrorTest4
	}

	return &dto.AcceptInviteLinkRes{
		HallID:   member.HallID,
		MemberID: member.ID,
		RoleID:   &member.RoleID,
		JoinedAt: member.JoinedAt,
	}, nil
}
