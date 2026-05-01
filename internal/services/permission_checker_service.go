package services

import (
	"context"
	"errors"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/suck-seed/yapp/internal/constants"
	"github.com/suck-seed/yapp/internal/database"
	"github.com/suck-seed/yapp/internal/repositories"
	"github.com/suck-seed/yapp/internal/utils"
)

type IPermissionCheckerService interface {
	// Permission checkers
	CanManageRoles(ctx context.Context, runner database.DBRunner, userID, hallID uuid.UUID) (bool, error)
	CanBanMembers(ctx context.Context, runner database.DBRunner, userID, hallID uuid.UUID) (bool, error)
	CanKickMembers(ctx context.Context, runner database.DBRunner, userID, hallID uuid.UUID) (bool, error)
	CanChangeNickname(ctx context.Context, runner database.DBRunner, userID, hallID uuid.UUID) (bool, error)
	CanManageNicknames(ctx context.Context, runner database.DBRunner, userID, hallID uuid.UUID) (bool, error)
	CanManageInvites(ctx context.Context, runner database.DBRunner, userID, hallID uuid.UUID) (bool, error)
	CanManageRequests(ctx context.Context, runner database.DBRunner, userID, hallID uuid.UUID) (bool, error)
	CanManageServers(ctx context.Context, runner database.DBRunner, userID, hallID uuid.UUID) (bool, error)

	checkPermission(ctx context.Context, runner database.DBRunner, userID uuid.UUID, hallID uuid.UUID, permColumn string) (bool, error)
}

type permissionCheckerService struct {
	repositories.IRoleRepository
	repositories.IUserRepository
	repositories.IHallRepository
	repositories.IBanRepsitory

	pool    *pgxpool.Pool
	timeout time.Duration
	mu      sync.RWMutex
}

func NewPermissionCheckerService(roleRepo repositories.IRoleRepository, userRepo repositories.IUserRepository, hallRepo repositories.IHallRepository, banRepo repositories.IBanRepsitory, pool *pgxpool.Pool) IPermissionCheckerService {
	return &permissionCheckerService{
		roleRepo,
		userRepo,
		hallRepo,
		banRepo,
		pool,
		time.Duration(2) * time.Second,
		sync.RWMutex{},
	}
}

// INTERNAL GENERIC FUNCTION
func (s *permissionCheckerService) checkPermission(ctx context.Context, runner database.DBRunner, userID uuid.UUID, hallID uuid.UUID, permColumn string) (bool, error) {

	// Checking if hall exists
	// Checking if hallOwner is userID
	// to rule out both conditions
	ownerID, err := s.IHallRepository.GetHallOwnerID(ctx, runner, hallID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return false, utils.ErrorHallNotFound
		}

		if errors.Is(err, context.DeadlineExceeded) || errors.Is(err, context.Canceled) {
			return false, utils.ErrorRequestTimeout
		}

		return false, utils.ErrorFetchingHall
	}

	if ownerID == userID {
		return true, nil
	}

	// giving access to admins
	userRole, err := s.IRoleRepository.GetUsersRoleInHall(ctx, runner, hallID, userID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return false, utils.ErrorUserDoesntBelongHall
		}
		if errors.Is(err, context.DeadlineExceeded) || errors.Is(err, context.Canceled) {
			return false, utils.ErrorRequestTimeout
		}
		return false, utils.ErrorFetchingRole
	}

	if userRole.IsAdmin {
		return true, nil
	}

	// Validating column
	if _, ok := constants.ValidPermissionColumns[permColumn]; !ok {
		return false, utils.ErrorPermissionsNotFound
	}

	// for any other role, check
	allowded, err := s.IRoleRepository.CheckUserPermission(ctx, runner, hallID, userID, permColumn)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return false, utils.ErrorUserDoesntBelongHall
		}
		if errors.Is(err, context.DeadlineExceeded) || errors.Is(err, context.Canceled) {
			return false, utils.ErrorRequestTimeout
		}
		return false, utils.ErrorFetchingPermission
	}

	return allowded, nil
}

// CanManageRoles - Return bool representing if the current user has appropriate permission to Manage other Roles from the corresponding hall
func (s *permissionCheckerService) CanManageRoles(ctx context.Context, runner database.DBRunner, userID, hallID uuid.UUID) (bool, error) {

	return s.checkPermission(ctx, runner, userID, hallID, constants.PermManageRoles)
}

// CanBanMembers - Return bool representing if the current user has appropriate permission to Ban other Users from the corresponding hall
func (s *permissionCheckerService) CanBanMembers(ctx context.Context, runner database.DBRunner, userID, hallID uuid.UUID) (bool, error) {

	return s.checkPermission(ctx, runner, userID, hallID, constants.PermBanMembers)

}

func (s *permissionCheckerService) CanKickMembers(ctx context.Context, runner database.DBRunner, userID, hallID uuid.UUID) (bool, error) {
	return s.checkPermission(ctx, runner, userID, hallID, constants.PermKickMembers)

}

func (s *permissionCheckerService) CanChangeNickname(ctx context.Context, runner database.DBRunner, userID, hallID uuid.UUID) (bool, error) {
	return s.checkPermission(ctx, runner, userID, hallID, constants.PermChangeNickname)
}

func (s *permissionCheckerService) CanManageNicknames(ctx context.Context, runner database.DBRunner, userID, hallID uuid.UUID) (bool, error) {
	return s.checkPermission(ctx, runner, userID, hallID, constants.PermManageNicknames)
}

func (s *permissionCheckerService) CanManageInvites(ctx context.Context, runner database.DBRunner, userID, hallID uuid.UUID) (bool, error) {
	return s.checkPermission(ctx, runner, userID, hallID, constants.PermManageInvites)
}

func (s *permissionCheckerService) CanManageRequests(ctx context.Context, runner database.DBRunner, userID, hallID uuid.UUID) (bool, error) {
	return s.checkPermission(ctx, runner, userID, hallID, constants.PermManageRequests)
}

// Implementation:
func (s *permissionCheckerService) CanManageServers(ctx context.Context, runner database.DBRunner, userID, hallID uuid.UUID) (bool, error) {
	return s.checkPermission(ctx, runner, userID, hallID, constants.PermManageServers)
}
