package services

import (
	"context"
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
	"github.com/suck-seed/yapp/internal/realtime"
	"github.com/suck-seed/yapp/internal/repositories"
	"github.com/suck-seed/yapp/internal/utils"
)

type IBanService interface {
	GetAllHallBans(ctx context.Context, userInfo *auth.UserInfo, hallID uuid.UUID) (*dto.AllBannedUserRes, error)
	BanUser(ctx context.Context, userInfo *auth.UserInfo, hallID uuid.UUID, req *dto.BanUserReq) (*dto.BanUserRes, error)
	UnbanUser(ctx context.Context, userInfo *auth.UserInfo, hallID uuid.UUID, banID uuid.UUID) (*dto.UnbanRes, error)
}

type banService struct {
	repositories.IBanRepsitory
	repositories.IUserRepository
	repositories.IHallRepository
	IPermissionCheckerService

	EventPublisher realtime.Publisher

	pool    *pgxpool.Pool
	timeout time.Duration
	mu      sync.RWMutex
}

func NewBanService(
	banRepo repositories.IBanRepsitory,
	userRepo repositories.IUserRepository,
	hallRepo repositories.IHallRepository,
	permissionChecker IPermissionCheckerService,
	eventPublisher realtime.Publisher,
	pool *pgxpool.Pool,
) IBanService {
	return &banService{
		banRepo,
		userRepo,
		hallRepo,
		permissionChecker,
		eventPublisher,
		pool,
		time.Duration(2) * time.Second,
		sync.RWMutex{},
	}
}

func (s *banService) GetAllHallBans(ctx context.Context, userInfo *auth.UserInfo, hallID uuid.UUID) (*dto.AllBannedUserRes, error) {
	ctx, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()

	conn, err := s.pool.Acquire(ctx)
	if err != nil {
		return nil, utils.ErrorInternal
	}
	defer conn.Release()
	runner := database.NewConnWrapper(conn)

	isMember, err := s.IHallRepository.IsUserHallMember(ctx, runner, hallID, userInfo.ID)
	if err != nil {
		return nil, utils.ErrorInternal
	}
	if !isMember {
		return nil, utils.ErrorUserDoesntBelongHall
	}

	canBan, err := s.CanBanMembers(ctx, runner, userInfo.ID, hallID)
	if err != nil {
		return nil, err
	}
	if !canBan {
		return nil, utils.ErrorUserCannotBanMembers
	}

	bans, err := s.IBanRepsitory.GetAllHallBans(ctx, runner, hallID)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) || errors.Is(err, context.Canceled) {
			return nil, utils.ErrorRequestTimeout
		}
		return nil, utils.ErrorFetchingBan
	}

	out := make([]dto.BanSummaryRes, 0, len(bans))
	for _, b := range bans {
		u, err := s.IUserRepository.GetUserById(ctx, runner, b.UserID)
		if err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				continue
			}
			return nil, utils.ErrorFetchingUser
		}
		var reason *string
		if b.Reason != "" {
			r := b.Reason
			reason = &r
		}
		out = append(out, dto.BanSummaryRes{
			ID:        b.ID,
			UserID:    b.UserID,
			Username:  u.Username,
			AvatarURL: u.AvatarURL,
			Reason:    reason,
			CreatedAt: b.CreatedAt,
			UpdatedAt: b.UpdatedAt,
		})
	}

	return &dto.AllBannedUserRes{Bans: out}, nil
}

func (s *banService) BanUser(ctx context.Context, userInfo *auth.UserInfo, hallID uuid.UUID, req *dto.BanUserReq) (*dto.BanUserRes, error) {
	ctx, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()

	reasonSanitized, err := utils.SanitizeText(&req.Reason)
	if err != nil {
		return nil, err
	}
	if reasonSanitized == nil || *reasonSanitized == "" {
		return nil, utils.ErrorInvalidInput
	}

	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return nil, utils.ErrorInternal
	}
	runner := database.NewTxWrapper(tx)
	defer runner.Rollback(ctx)

	isMember, err := s.IHallRepository.IsUserHallMember(ctx, runner, hallID, userInfo.ID)
	if err != nil {
		return nil, utils.ErrorInternal
	}
	if !isMember {
		return nil, utils.ErrorUserDoesntBelongHall
	}

	canBan, err := s.CanBanMembers(ctx, runner, userInfo.ID, hallID)
	if err != nil {
		return nil, err
	}
	if !canBan {
		return nil, utils.ErrorUserCannotBanMembers
	}

	ownerID, err := s.IHallRepository.GetHallOwnerID(ctx, runner, hallID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, utils.ErrorHallNotFound
		}
		return nil, utils.ErrorFetchingHall
	}

	if req.UserID == ownerID {
		return nil, utils.ErrorCannotBanHallOwner
	}
	if req.UserID == userInfo.ID {
		return nil, utils.ErrorCannotBanYourself
	}

	already, err := s.IBanRepsitory.IsUserBanned(ctx, runner, hallID, req.UserID)
	if err != nil {
		return nil, utils.ErrorFetchingBan
	}
	if already {
		return nil, utils.ErrorUserAlreadyBanned
	}

	bannedUser, err := s.IUserRepository.GetUserById(ctx, runner, req.UserID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, utils.ErrorUserNotFound
		}
		return nil, utils.ErrorFetchingUser
	}

	banID, err := uuid.NewV7()
	if err != nil {
		return nil, utils.ErrorInternal
	}

	newBan := &models.HallBan{
		ID:     banID,
		Reason: *reasonSanitized,
		UserID: req.UserID,
		HallID: hallID,
	}

	saved, err := s.IBanRepsitory.BanUser(ctx, runner, newBan)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) || errors.Is(err, context.Canceled) {
			return nil, utils.ErrorRequestTimeout
		}
		return nil, utils.ErrorInternal
	}

	isTargetMember, err := s.IHallRepository.IsUserHallMember(ctx, runner, hallID, req.UserID)
	if err != nil {
		return nil, utils.ErrorInternal
	}
	if isTargetMember {
		if err := s.IHallRepository.KickHallMember(ctx, runner, hallID, req.UserID); err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				// No longer a member before kick (concurrent leave); ban row still applies.
			} else if errors.Is(err, context.DeadlineExceeded) || errors.Is(err, context.Canceled) {
				return nil, utils.ErrorRequestTimeout
			} else {
				return nil, utils.ErrorInternal
			}
		}
	}

	if err := runner.Commit(ctx); err != nil {
		return nil, utils.ErrorInternal
	}

	// PUBLISH EVENT
	publishHubEvent(s.EventPublisher, realtime.HubEvent{
		Type:   realtime.HubEventUserBannedFromHall,
		HallID: hallID,
		UserID: req.UserID,
	})

	return &dto.BanUserRes{
		ID:     saved.ID,
		Reason: saved.Reason,
		UserID: saved.UserID,
		User: dto.BannedUserInfo{
			ID:       bannedUser.ID,
			Username: bannedUser.Username,
			Avatar:   bannedUser.AvatarURL,
		},
		CreatedAt: saved.CreatedAt,
		UpdatedAt: saved.UpdatedAt,
	}, nil
}

func (s *banService) UnbanUser(ctx context.Context, userInfo *auth.UserInfo, hallID uuid.UUID, banID uuid.UUID) (*dto.UnbanRes, error) {
	ctx, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()

	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return nil, utils.ErrorInternal
	}
	runner := database.NewTxWrapper(tx)
	defer runner.Rollback(ctx)

	isMember, err := s.IHallRepository.IsUserHallMember(ctx, runner, hallID, userInfo.ID)
	if err != nil {
		return nil, utils.ErrorInternal
	}
	if !isMember {
		return nil, utils.ErrorUserDoesntBelongHall
	}

	canBan, err := s.CanBanMembers(ctx, runner, userInfo.ID, hallID)
	if err != nil {
		return nil, err
	}
	if !canBan {
		return nil, utils.ErrorUserCannotBanMembers
	}

	ban, err := s.IBanRepsitory.GetBanByID(ctx, runner, banID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, utils.ErrorBanNotFound
		}
		return nil, utils.ErrorFetchingBan
	}

	if ban.HallID != hallID {
		return nil, utils.ErrorBanNotFound
	}

	if _, err := s.IBanRepsitory.UnBanUser(ctx, runner, banID); err != nil {
		if errors.Is(err, context.DeadlineExceeded) || errors.Is(err, context.Canceled) {
			return nil, utils.ErrorRequestTimeout
		}
		return nil, err
	}

	u, err := s.IUserRepository.GetUserById(ctx, runner, ban.UserID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, utils.ErrorUserNotFound
		}
		return nil, utils.ErrorFetchingUser
	}

	if err := runner.Commit(ctx); err != nil {
		return nil, utils.ErrorInternal
	}

	return &dto.UnbanRes{
		UserID:   u.ID,
		Username: u.Username,
		Message:  "User unbanned successfully",
	}, nil
}
