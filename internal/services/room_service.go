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
	dto "github.com/suck-seed/yapp/internal/dto/room"
	"github.com/suck-seed/yapp/internal/models"
	"github.com/suck-seed/yapp/internal/repositories"
	"github.com/suck-seed/yapp/internal/utils"
)

type IRoomService interface {
	CreateRoom(c context.Context, userInfo *auth.UserInfo, hallID uuid.UUID, req *dto.CreateRoomReq) (*dto.CreateRoomRes, error)
	GetRoom(c context.Context, userInfo *auth.UserInfo, hallID uuid.UUID, roomID uuid.UUID) (*dto.GetRoomRes, error)
	GetHallRooms(c context.Context, userInfo *auth.UserInfo, hallID uuid.UUID) (*dto.GetHallRoomsRes, error)
	UpdateRoom(c context.Context, userInfo *auth.UserInfo, hallID uuid.UUID, roomID uuid.UUID, req *dto.UpdateRoomReq) (*dto.UpdateRoomRes, error)
	DeleteRoom(c context.Context, userInfo *auth.UserInfo, hallID uuid.UUID, roomID uuid.UUID) error
	MoveRoom(c context.Context, userInfo *auth.UserInfo, hallID uuid.UUID, roomID uuid.UUID, req *dto.MoveRoomReq) (*dto.RoomRes, error)
	// internal
	GetRoomByID(c context.Context, roomID uuid.UUID) (*models.Room, error)
	IsUserRoomMember(c context.Context, roomID uuid.UUID, userID uuid.UUID) (bool, error)
}

type roomService struct {
	repositories.IHallRepository
	repositories.IFloorRepository
	repositories.IRoomRepository
	repositories.IBanRepsitory

	IPermissionCheckerService

	pool    *pgxpool.Pool
	timeout time.Duration
	mu      sync.RWMutex
}

func NewRoomService(
	hallRepo repositories.IHallRepository,
	floorRepo repositories.IFloorRepository,
	roomRepo repositories.IRoomRepository,
	banRepo repositories.IBanRepsitory,
	permissionChecker IPermissionCheckerService,
	pool *pgxpool.Pool,
) IRoomService {
	return &roomService{
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

func roomToRes(r *models.Room) dto.RoomRes {
	return dto.RoomRes{
		ID:        r.ID,
		HallID:    r.HallID,
		FloorID:   r.FloorID,
		Name:      r.Name,
		RoomType:  r.RoomType,
		Position:  r.Position,
		IsPrivate: r.IsPrivate,
		CreatedAt: r.CreatedAt,
		UpdatedAt: r.UpdatedAt,
	}
}

func (s *roomService) requireManageServers(ctx context.Context, runner database.DBRunner, userID, hallID uuid.UUID) error {
	ok, err := s.IPermissionCheckerService.CanManageServers(ctx, runner, userID, hallID)
	if err != nil {
		return err
	}
	if !ok {
		return utils.ErrorUserCannotManageServer
	}
	return nil
}

// ── CreateRoom ────────────────────────────────────────────────────────────────

func (s *roomService) CreateRoom(c context.Context, userInfo *auth.UserInfo, hallID uuid.UUID, req *dto.CreateRoomReq) (*dto.CreateRoomRes, error) {
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

	hallExists, err := s.IHallRepository.DoesHallExist(ctx, runner, hallID)
	if err != nil || !hallExists {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, utils.ErrorHallNotFound
		}
		if isDeadline(err) {
			return nil, utils.ErrorRequestTimeout
		}
		return nil, utils.ErrorFetchingHall
	}

	// If a floor is specified, verify it belongs to this hall
	if req.FloorID != nil {
		floorExists, err := s.IFloorRepository.DoesFloorExistInHall(ctx, runner, *req.FloorID, hallID)
		if err != nil || !floorExists {
			if isDeadline(err) {
				return nil, utils.ErrorRequestTimeout
			}
			return nil, utils.ErrorFloorNotFound
		}
	}

	canonRoomType, err := utils.ParseRoomType(req.RoomType)
	if err != nil {
		return nil, err
	}

	maxPos, err := s.IRoomRepository.GetMaxPositionInContainer(ctx, runner, hallID, req.FloorID)
	if err != nil {
		if isDeadline(err) {
			return nil, utils.ErrorRequestTimeout
		}
		return nil, utils.ErrorFetchingRoom
	}

	roomID, err := uuid.NewV7()
	if err != nil {
		return nil, utils.ErrorInternal
	}

	isPrivate := false
	if req.IsPrivate != nil {
		isPrivate = *req.IsPrivate
	}

	created, err := s.IRoomRepository.CreateRoom(ctx, runner, &models.Room{
		ID:        roomID,
		HallID:    hallID,
		FloorID:   req.FloorID,
		Name:      canonName,
		RoomType:  string(canonRoomType),
		Position:  maxPos + 1000.0,
		IsPrivate: isPrivate,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	})
	if err != nil {
		if isDeadline(err) {
			return nil, utils.ErrorRequestTimeout
		}
		return nil, utils.ErrorCreatingRoom
	}

	if err := runner.Commit(ctx); err != nil {
		return nil, utils.ErrorInternal
	}

	res := roomToRes(created)
	return &res, nil
}

// ── GetRoom ───────────────────────────────────────────────────────────────────

func (s *roomService) GetRoom(c context.Context, userInfo *auth.UserInfo, hallID uuid.UUID, roomID uuid.UUID) (*dto.GetRoomRes, error) {
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

	room, err := s.IRoomRepository.GetRoomByID(ctx, runner, roomID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, utils.ErrorRoomNotFound
		}
		if isDeadline(err) {
			return nil, utils.ErrorRequestTimeout
		}
		return nil, utils.ErrorFetchingRoom
	}

	// Ensure room belongs to this hall
	if room.HallID != hallID {
		return nil, utils.ErrorRoomNotFound
	}

	res := roomToRes(room)
	return &res, nil
}

// ── GetHallRooms ──────────────────────────────────────────────────────────────

func (s *roomService) GetHallRooms(c context.Context, userInfo *auth.UserInfo, hallID uuid.UUID) (*dto.GetHallRoomsRes, error) {
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

	// Fetch all rooms and floors in parallel
	type roomsResult struct {
		rooms []*models.Room
		err   error
	}
	type floorsResult struct {
		floors []*models.Floor
		err    error
	}

	roomsCh := make(chan roomsResult, 1)
	floorsCh := make(chan floorsResult, 1)

	go func() {
		rooms, err := s.IRoomRepository.GetRoomsByHallID(ctx, s.pool, hallID)
		roomsCh <- roomsResult{rooms, err}
	}()
	go func() {
		floors, err := s.IFloorRepository.GetFloorsByHallID(ctx, s.pool, hallID)
		floorsCh <- floorsResult{floors, err}
	}()

	rr := <-roomsCh
	fr := <-floorsCh

	if rr.err != nil {
		if isDeadline(rr.err) {
			return nil, utils.ErrorRequestTimeout
		}
		return nil, utils.ErrorFetchingRoom
	}
	if fr.err != nil {
		if isDeadline(fr.err) {
			return nil, utils.ErrorRequestTimeout
		}
		return nil, utils.ErrorFetchingFloor
	}

	floorRooms := make(map[uuid.UUID][]dto.RoomRes)
	var topLevel []dto.RoomRes

	for _, rm := range rr.rooms {
		res := roomToRes(rm)
		if rm.FloorID == nil {
			topLevel = append(topLevel, res)
		} else {
			floorRooms[*rm.FloorID] = append(floorRooms[*rm.FloorID], res)
		}
	}

	floors := make([]dto.FloorWithRoomsRes, 0, len(fr.floors))
	for _, f := range fr.floors {
		rooms := floorRooms[f.ID]
		if rooms == nil {
			rooms = []dto.RoomRes{}
		}
		floors = append(floors, dto.FloorWithRoomsRes{
			ID:        f.ID,
			HallID:    f.HallID,
			Name:      f.Name,
			Position:  f.Position,
			IsPrivate: f.IsPrivate,
			Rooms:     rooms,
		})
	}

	if topLevel == nil {
		topLevel = []dto.RoomRes{}
	}

	return &dto.GetHallRoomsRes{
		TopLevel: topLevel,
		Floors:   floors,
	}, nil
}

// ── UpdateRoom ────────────────────────────────────────────────────────────────

func (s *roomService) UpdateRoom(c context.Context, userInfo *auth.UserInfo, hallID uuid.UUID, roomID uuid.UUID, req *dto.UpdateRoomReq) (*dto.UpdateRoomRes, error) {
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

	// Verify room belongs to this hall before updating
	room, err := s.IRoomRepository.GetRoomByID(ctx, runner, roomID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, utils.ErrorRoomNotFound
		}
		if isDeadline(err) {
			return nil, utils.ErrorRequestTimeout
		}
		return nil, utils.ErrorFetchingRoom
	}
	if room.HallID != hallID {
		return nil, utils.ErrorRoomNotFound
	}

	updated, err := s.IRoomRepository.UpdateRoom(ctx, runner, roomID, req.Name, req.IsPrivate)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, utils.ErrorRoomNotFound
		}
		if isDeadline(err) {
			return nil, utils.ErrorRequestTimeout
		}
		return nil, utils.ErrorFetchingRoom
	}

	if err := runner.Commit(ctx); err != nil {
		return nil, utils.ErrorInternal
	}

	res := roomToRes(updated)
	return &res, nil
}

// ── DeleteRoom ────────────────────────────────────────────────────────────────

func (s *roomService) DeleteRoom(c context.Context, userInfo *auth.UserInfo, hallID uuid.UUID, roomID uuid.UUID) error {
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

	room, err := s.IRoomRepository.GetRoomByID(ctx, runner, roomID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return utils.ErrorRoomNotFound
		}
		if isDeadline(err) {
			return utils.ErrorRequestTimeout
		}
		return utils.ErrorFetchingRoom
	}
	if room.HallID != hallID {
		return utils.ErrorRoomNotFound
	}

	if err := s.IRoomRepository.DeleteRoom(ctx, runner, roomID); err != nil {
		if isDeadline(err) {
			return utils.ErrorRequestTimeout
		}
		return utils.ErrorInternal
	}

	return runner.Commit(ctx)
}

// ── MoveRoom ──────────────────────────────────────────────────────────────────

func (s *roomService) MoveRoom(c context.Context, userInfo *auth.UserInfo, hallID uuid.UUID, roomID uuid.UUID, req *dto.MoveRoomReq) (*dto.RoomRes, error) {
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

	// Verify room belongs to this hall
	room, err := s.IRoomRepository.GetRoomByID(ctx, runner, roomID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, utils.ErrorRoomNotFound
		}
		if isDeadline(err) {
			return nil, utils.ErrorRequestTimeout
		}
		return nil, utils.ErrorFetchingRoom
	}
	if room.HallID != hallID {
		return nil, utils.ErrorRoomNotFound
	}

	// Verify target floor belongs to this hall (if targeting a floor)
	if req.NewFloorID != nil {
		floorExists, err := s.IFloorRepository.DoesFloorExistInHall(ctx, runner, *req.NewFloorID, hallID)
		if err != nil || !floorExists {
			if isDeadline(err) {
				return nil, utils.ErrorRequestTimeout
			}
			return nil, utils.ErrorFloorNotFound
		}
	}

	lower, upper, err := s.IRoomRepository.GetRoomPositionBounds(ctx, runner, hallID, req.NewFloorID, req.AfterID)
	if err != nil {
		if isDeadline(err) {
			return nil, utils.ErrorRequestTimeout
		}
		return nil, utils.ErrorFetchingRoom
	}

	newPos := utils.CalcPosition(lower, upper)

	moved, err := s.IRoomRepository.MoveRoom(ctx, runner, roomID, req.NewFloorID, newPos)
	if err != nil {
		if isDeadline(err) {
			return nil, utils.ErrorRequestTimeout
		}
		return nil, utils.ErrorInternal
	}

	if err := runner.Commit(ctx); err != nil {
		return nil, utils.ErrorInternal
	}

	res := roomToRes(moved)
	return &res, nil
}

// ── Internal ──────────────────────────────────────────────────────────────────

func (s *roomService) GetRoomByID(c context.Context, roomID uuid.UUID) (*models.Room, error) {
	ctx, cancel := context.WithTimeout(c, s.timeout)
	defer cancel()

	conn, err := s.pool.Acquire(ctx)
	if err != nil {
		return nil, utils.ErrorInternal
	}
	defer conn.Release()

	room, err := s.IRoomRepository.GetRoomByID(ctx, database.NewConnWrapper(conn), roomID)
	if err != nil {
		return nil, utils.ErrorFetchingRoom
	}
	return room, nil
}

func (s *roomService) IsUserRoomMember(c context.Context, roomID uuid.UUID, userID uuid.UUID) (bool, error) {
	ctx, cancel := context.WithTimeout(c, s.timeout)
	defer cancel()

	conn, err := s.pool.Acquire(ctx)
	if err != nil {
		return false, utils.ErrorInternal
	}
	defer conn.Release()

	return s.IRoomRepository.IsUserRoomMember(ctx, database.NewConnWrapper(conn), roomID, userID)
}
