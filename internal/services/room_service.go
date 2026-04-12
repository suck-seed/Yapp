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
	CreateRoom(c context.Context, userInfo *auth.UserInfo, req *dto.CreateRoomReq) (*dto.CreateRoomRes, error)
	GetRoom(c context.Context, roomID uuid.UUID) (*dto.GetRoomRes, error)
	GetHallRooms(c context.Context, hallID uuid.UUID) (*dto.GetHallRoomsRes, error)
	UpdateRoom(c context.Context, roomID uuid.UUID, req *dto.UpdateRoomReq) (*dto.UpdateRoomRes, error)
	DeleteRoom(c context.Context, roomID uuid.UUID) error

	MoveRoom(c context.Context, roomID uuid.UUID, req *dto.MoveRoomReq) (*dto.RoomRes, error)
	// internal
	GetRoomByID(c context.Context, roomID uuid.UUID) (*models.Room, error)
	IsUserRoomMember(c context.Context, roomID uuid.UUID, userID uuid.UUID) (bool, error)
}

type roomService struct {
	repositories.IHallRepository
	repositories.IFloorRepository
	repositories.IRoomRepository
	pool    *pgxpool.Pool
	timeout time.Duration
	mu      sync.RWMutex
}

func NewRoomService(
	hallRepo repositories.IHallRepository,
	floorRepo repositories.IFloorRepository,
	roomRepo repositories.IRoomRepository,
	pool *pgxpool.Pool,
) IRoomService {
	return &roomService{
		IHallRepository:  hallRepo,
		IFloorRepository: floorRepo,
		IRoomRepository:  roomRepo,
		pool:             pool,
		timeout:          2 * time.Second,
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

// ── CreateRoom ────────────────────────────────────────────────────────────────

func (s *roomService) CreateRoom(c context.Context, userInfo *auth.UserInfo, req *dto.CreateRoomReq) (*dto.CreateRoomRes, error) {
	ctx, cancel := context.WithTimeout(c, s.timeout)
	defer cancel()

	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return nil, utils.ErrorInternal
	}
	runner := database.NewTxWrapper(tx)
	defer runner.Rollback(ctx)

	canonName, err := utils.SanitizeFloorname(req.Name)
	if err != nil {
		return nil, err
	}

	hallExists, err := s.IHallRepository.DoesHallExist(ctx, runner, req.HallID)
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
		floorExists, err := s.IFloorRepository.DoesFloorExistInHall(ctx, runner, *req.FloorID, req.HallID)
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

	// Position is scoped to the container (floor or top-level)
	maxPos, err := s.IRoomRepository.GetMaxPositionInContainer(ctx, runner, req.HallID, req.FloorID)
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
		HallID:    req.HallID,
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

func (s *roomService) GetRoom(c context.Context, roomID uuid.UUID) (*dto.GetRoomRes, error) {
	ctx, cancel := context.WithTimeout(c, s.timeout)
	defer cancel()

	room, err := s.IRoomRepository.GetRoomByID(ctx, s.pool, roomID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, utils.ErrorRoomNotFound
		}
		if isDeadline(err) {
			return nil, utils.ErrorRequestTimeout
		}
		return nil, utils.ErrorFetchingRoom
	}

	res := roomToRes(room)
	return &res, nil
}

// ── GetHallRooms ──────────────────────────────────────────────────────────────

func (s *roomService) GetHallRooms(c context.Context, hallID uuid.UUID) (*dto.GetHallRoomsRes, error) {
	ctx, cancel := context.WithTimeout(c, s.timeout)
	defer cancel()

	// Fetch all rooms and floors in parallel using goroutines
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

	// Group rooms into their containers in Go — no complex SQL needed
	// Index floor rooms by floor_id for O(1) lookup
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
			rooms = []dto.RoomRes{} // never return null for empty floors
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

func (s *roomService) UpdateRoom(c context.Context, roomID uuid.UUID, req *dto.UpdateRoomReq) (*dto.UpdateRoomRes, error) {
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

func (s *roomService) DeleteRoom(c context.Context, roomID uuid.UUID) error {
	ctx, cancel := context.WithTimeout(c, s.timeout)
	defer cancel()

	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return utils.ErrorInternal
	}
	runner := database.NewTxWrapper(tx)
	defer runner.Rollback(ctx)

	_, err = s.IRoomRepository.GetRoomByID(ctx, runner, roomID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return utils.ErrorRoomNotFound
		}
		if isDeadline(err) {
			return utils.ErrorRequestTimeout
		}
		return utils.ErrorFetchingRoom
	}

	if err := s.IRoomRepository.DeleteRoom(ctx, runner, roomID); err != nil {
		if isDeadline(err) {
			return utils.ErrorRequestTimeout
		}
		return utils.ErrorInternal
	}

	return runner.Commit(ctx)
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

// MoveRoom handles every drag-drop case atomically:
//   - reorder within same floor
//   - drag into a floor and place at exact spot
//   - drag out of floor and place at exact spot among top-level rooms
func (s *roomService) MoveRoom(c context.Context, roomID uuid.UUID, req *dto.MoveRoomReq) (*dto.RoomRes, error) {
	ctx, cancel := context.WithTimeout(c, s.timeout)
	defer cancel()

	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return nil, utils.ErrorInternal
	}
	runner := database.NewTxWrapper(tx)
	defer runner.Rollback(ctx)

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
	if room.HallID != req.HallID {
		return nil, utils.ErrorUserDoesntBelongRoom
	}

	// Verify target floor belongs to this hall (if targeting a floor)
	if req.NewFloorID != nil {
		floorExists, err := s.IFloorRepository.DoesFloorExistInHall(ctx, runner, *req.NewFloorID, req.HallID)
		if err != nil || !floorExists {
			if isDeadline(err) {
				return nil, utils.ErrorRequestTimeout
			}
			return nil, utils.ErrorFloorNotFound
		}
	}

	lower, upper, err := s.IRoomRepository.GetRoomPositionBounds(ctx, runner, req.HallID, req.NewFloorID, req.AfterID)
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
