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
	dto "github.com/suck-seed/yapp/internal/dto/floor"
	"github.com/suck-seed/yapp/internal/models"
	"github.com/suck-seed/yapp/internal/repositories"
	"github.com/suck-seed/yapp/internal/utils"
)

type IFloorService interface {
	CreateFloor(c context.Context, userInfo *auth.UserInfo, hallID uuid.UUID, req *dto.CreateFloorReq) (*dto.CreateFloorRes, error)
	GetFloor(c context.Context, userInfo *auth.UserInfo, hallID uuid.UUID, floorID uuid.UUID) (*dto.GetFloorRes, error)
	GetFloors(c context.Context, userInfo *auth.UserInfo, hallID uuid.UUID) (*dto.GetFloorsRes, error)
	UpdateFloor(c context.Context, userInfo *auth.UserInfo, hallID uuid.UUID, floorID uuid.UUID, req *dto.UpdateFloorReq) (*dto.UpdateFloorRes, error)
	DeleteFloor(c context.Context, userInfo *auth.UserInfo, hallID uuid.UUID, floorID uuid.UUID) error
	MoveFloor(c context.Context, userInfo *auth.UserInfo, hallID uuid.UUID, floorID uuid.UUID, req *dto.MoveFloorReq) (*dto.GetFloorRes, error)
}

type floorService struct {
	repositories.IHallRepository
	repositories.IFloorRepository
	repositories.IRoomRepository
	repositories.IBanRepsitory

	IPermissionCheckerService

	pool    *pgxpool.Pool
	timeout time.Duration
	mu      sync.RWMutex
}

func NewFloorService(
	hallRepo repositories.IHallRepository,
	floorRepo repositories.IFloorRepository,
	roomRepo repositories.IRoomRepository,
	banRepo repositories.IBanRepsitory,
	permissionChecker IPermissionCheckerService,
	pool *pgxpool.Pool,
) IFloorService {
	return &floorService{
		hallRepo,
		floorRepo,
		roomRepo,
		banRepo,
		permissionChecker,
		pool,
		time.Duration(2) * time.Second,
		sync.RWMutex{},
	}
}

// ── helpers ───────────────────────────────────────────────────────────────────

func floorToGetRes(f *models.Floor) dto.GetFloorRes {
	return dto.GetFloorRes{
		ID:        f.ID,
		HallID:    f.HallID,
		Name:      f.Name,
		Position:  f.Position,
		IsPrivate: f.IsPrivate,
		CreatedAt: f.CreatedAt,
		UpdatedAt: f.UpdatedAt,
	}
}

func (s *floorService) requireManageServers(ctx context.Context, runner database.DBRunner, userID, hallID uuid.UUID) error {
	ok, err := s.IPermissionCheckerService.CanManageServers(ctx, runner, userID, hallID)
	if err != nil {
		return err
	}
	if !ok {
		return utils.ErrorUserCannotManageServer
	}
	return nil
}

// add this alongside the other helpers in floor_service.go, below floorToGetRes

func isDeadline(err error) bool {
	return errors.Is(err, context.DeadlineExceeded) || errors.Is(err, context.Canceled)
}

// ── CreateFloor ───────────────────────────────────────────────────────────────

func (s *floorService) CreateFloor(c context.Context, userInfo *auth.UserInfo, hallID uuid.UUID, req *dto.CreateFloorReq) (*dto.CreateFloorRes, error) {
	ctx, cancel := context.WithTimeout(c, s.timeout)
	defer cancel()

	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return nil, utils.ErrorInternal
	}
	runner := database.NewTxWrapper(tx)
	defer runner.Rollback(ctx)

	if err := s.requireManageServers(ctx, runner, userInfo.ID, hallID); err != nil {
		return nil, err
	}

	canonName, err := utils.SanitizeFloorname(req.Name)
	if err != nil {
		return nil, err
	}

	exists, err := s.IHallRepository.DoesHallExist(ctx, runner, hallID)
	if err != nil || !exists {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, utils.ErrorHallNotFound
		}
		if isDeadline(err) {
			return nil, utils.ErrorRequestTimeout
		}
		return nil, utils.ErrorFetchingHall
	}

	maxPos, err := s.IFloorRepository.GetMaxPosition(ctx, runner, hallID)
	if err != nil {
		if isDeadline(err) {
			return nil, utils.ErrorRequestTimeout
		}
		return nil, utils.ErrorFetchingFloor
	}

	floorID, err := uuid.NewV7()
	if err != nil {
		return nil, utils.ErrorInternal
	}

	floor := &models.Floor{
		ID:        floorID,
		HallID:    hallID,
		Name:      canonName,
		Position:  maxPos + 1000.0,
		IsPrivate: *req.IsPrivate,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	created, err := s.IFloorRepository.CreateFloor(ctx, runner, floor)
	if err != nil {
		if isDeadline(err) {
			return nil, utils.ErrorRequestTimeout
		}
		return nil, utils.ErrorCreatingFloor
	}

	if err := runner.Commit(ctx); err != nil {
		return nil, utils.ErrorInternal
	}

	return &dto.CreateFloorRes{
		ID:        created.ID,
		HallID:    created.HallID,
		Name:      created.Name,
		Position:  created.Position,
		IsPrivate: created.IsPrivate,
		CreatedAt: created.CreatedAt,
		UpdatedAt: created.UpdatedAt,
	}, nil
}

// ── GetFloors ─────────────────────────────────────────────────────────────────

func (s *floorService) GetFloors(c context.Context, userInfo *auth.UserInfo, hallID uuid.UUID) (*dto.GetFloorsRes, error) {
	ctx, cancel := context.WithTimeout(c, s.timeout)
	defer cancel()

	conn, err := s.pool.Acquire(ctx)
	if err != nil {
		return nil, utils.ErrorInternal
	}
	defer conn.Release()
	runner := database.NewConnWrapper(conn)

	ok, err := s.IHallRepository.IsUserHallMember(ctx, runner, hallID, userInfo.ID)
	if err != nil {
		return nil, utils.ErrorInternal
	}
	if !ok {
		return nil, utils.ErrorUserDoesntBelongHall
	}

	floors, err := s.IFloorRepository.GetFloorsByHallID(ctx, runner, hallID)
	if err != nil {
		if isDeadline(err) {
			return nil, utils.ErrorRequestTimeout
		}
		return nil, utils.ErrorFetchingFloor
	}

	res := make([]dto.GetFloorRes, len(floors))
	for i, f := range floors {
		res[i] = floorToGetRes(f)
	}
	return &dto.GetFloorsRes{Floors: res}, nil
}

// ── GetFloor ──────────────────────────────────────────────────────────────────

func (s *floorService) GetFloor(c context.Context, userInfo *auth.UserInfo, hallID uuid.UUID, floorID uuid.UUID) (*dto.GetFloorRes, error) {
	ctx, cancel := context.WithTimeout(c, s.timeout)
	defer cancel()

	conn, err := s.pool.Acquire(ctx)
	if err != nil {
		return nil, utils.ErrorInternal
	}
	defer conn.Release()
	runner := database.NewConnWrapper(conn)

	ok, err := s.IHallRepository.IsUserHallMember(ctx, runner, hallID, userInfo.ID)
	if err != nil {
		return nil, utils.ErrorInternal
	}
	if !ok {
		return nil, utils.ErrorUserDoesntBelongHall
	}

	floor, err := s.IFloorRepository.GetFloorByID(ctx, runner, floorID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, utils.ErrorFloorNotFound
		}
		if isDeadline(err) {
			return nil, utils.ErrorRequestTimeout
		}
		return nil, utils.ErrorFetchingFloor
	}

	// Ensure floor belongs to this hall
	if floor.HallID != hallID {
		return nil, utils.ErrorFloorNotFound
	}

	res := floorToGetRes(floor)
	return &res, nil
}

// ── UpdateFloor ───────────────────────────────────────────────────────────────

func (s *floorService) UpdateFloor(c context.Context, userInfo *auth.UserInfo, hallID uuid.UUID, floorID uuid.UUID, req *dto.UpdateFloorReq) (*dto.UpdateFloorRes, error) {
	ctx, cancel := context.WithTimeout(c, s.timeout)
	defer cancel()

	if req.Name == nil && req.IsPrivate == nil {
		return nil, utils.ErrorNoFieldsToUpdate
	}

	if req.Name != nil {
		canon, err := utils.SanitizeFloorname(*req.Name)
		if err != nil {
			return nil, err
		}
		req.Name = &canon
	}

	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return nil, utils.ErrorInternal
	}
	runner := database.NewTxWrapper(tx)
	defer runner.Rollback(ctx)

	if err := s.requireManageServers(ctx, runner, userInfo.ID, hallID); err != nil {
		return nil, err
	}

	// Verify floor belongs to this hall before updating
	exists, err := s.IFloorRepository.DoesFloorExistInHall(ctx, runner, floorID, hallID)
	if err != nil || !exists {
		return nil, utils.ErrorFloorNotFound
	}

	updated, err := s.IFloorRepository.UpdateFloor(ctx, runner, floorID, req.Name, req.IsPrivate)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, utils.ErrorFloorNotFound
		}
		if isDeadline(err) {
			return nil, utils.ErrorRequestTimeout
		}
		return nil, utils.ErrorFetchingFloor
	}

	if err := runner.Commit(ctx); err != nil {
		return nil, utils.ErrorInternal
	}

	return &dto.UpdateFloorRes{
		ID:        updated.ID,
		HallID:    updated.HallID,
		Name:      updated.Name,
		Position:  updated.Position,
		IsPrivate: updated.IsPrivate,
		UpdatedAt: updated.UpdatedAt,
	}, nil
}

// ── DeleteFloor ───────────────────────────────────────────────────────────────

func (s *floorService) DeleteFloor(c context.Context, userInfo *auth.UserInfo, hallID uuid.UUID, floorID uuid.UUID) error {
	ctx, cancel := context.WithTimeout(c, s.timeout)
	defer cancel()

	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return utils.ErrorInternal
	}
	runner := database.NewTxWrapper(tx)
	defer runner.Rollback(ctx)

	if err := s.requireManageServers(ctx, runner, userInfo.ID, hallID); err != nil {
		return err
	}

	// Verify floor belongs to this hall before deleting
	exists, err := s.IFloorRepository.DoesFloorExistInHall(ctx, runner, floorID, hallID)
	if err != nil || !exists {
		return utils.ErrorFloorNotFound
	}

	if err := s.IFloorRepository.DeleteFloor(ctx, runner, floorID); err != nil {
		if isDeadline(err) {
			return utils.ErrorRequestTimeout
		}
		return utils.ErrorInternal
	}

	if err := runner.Commit(ctx); err != nil {
		return utils.ErrorInternal
	}
	return nil
}

// ── MoveFloor ─────────────────────────────────────────────────────────────────

func (s *floorService) MoveFloor(c context.Context, userInfo *auth.UserInfo, hallID uuid.UUID, floorID uuid.UUID, req *dto.MoveFloorReq) (*dto.GetFloorRes, error) {
	ctx, cancel := context.WithTimeout(c, s.timeout)
	defer cancel()

	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return nil, utils.ErrorInternal
	}
	runner := database.NewTxWrapper(tx)
	defer runner.Rollback(ctx)

	if err := s.requireManageServers(ctx, runner, userInfo.ID, hallID); err != nil {
		return nil, err
	}

	// Verify floor belongs to this hall
	exists, err := s.IFloorRepository.DoesFloorExistInHall(ctx, runner, floorID, hallID)
	if err != nil || !exists {
		if isDeadline(err) {
			return nil, utils.ErrorRequestTimeout
		}
		return nil, utils.ErrorFloorNotFound
	}

	lower, upper, err := s.IFloorRepository.GetFloorPositionBounds(ctx, runner, hallID, req.AfterID)
	if err != nil {
		if isDeadline(err) {
			return nil, utils.ErrorRequestTimeout
		}
		return nil, utils.ErrorFetchingFloor
	}

	newPos := utils.CalcPosition(lower, upper)

	updated, err := s.IFloorRepository.UpdateFloorPosition(ctx, runner, floorID, newPos)
	if err != nil {
		if isDeadline(err) {
			return nil, utils.ErrorRequestTimeout
		}
		return nil, utils.ErrorFetchingFloor
	}

	if err := runner.Commit(ctx); err != nil {
		return nil, utils.ErrorInternal
	}

	res := floorToGetRes(updated)
	return &res, nil
}
