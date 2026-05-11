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
	"github.com/suck-seed/yapp/internal/realtime"
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
	GetAccessibleRoomsForUser(c context.Context, userInfo *auth.UserInfo) (map[uuid.UUID]uuid.UUID, error)

	// Room member management
	GetRoomMembers(c context.Context, userInfo *auth.UserInfo, hallID uuid.UUID, roomID uuid.UUID) (*dto.GetRoomMembersRes, error)
	GetRoomMember(c context.Context, userInfo *auth.UserInfo, hallID uuid.UUID, roomID uuid.UUID, memberID uuid.UUID) (*dto.GetRoomMemberRes, error)
	AddRoomMember(c context.Context, userInfo *auth.UserInfo, hallID uuid.UUID, roomID uuid.UUID, memberID uuid.UUID) (*dto.RoomAccessMemberRes, error)
	RemoveRoomMember(c context.Context, userInfo *auth.UserInfo, hallID uuid.UUID, roomID uuid.UUID, memberID uuid.UUID) (*dto.RoomAccessMemberRes, error)
	SyncRoomMembersToFloor(c context.Context, userInfo *auth.UserInfo, hallID uuid.UUID, roomID uuid.UUID) (*dto.RoomRes, error)
}

type roomService struct {
	repositories.IHallRepository
	repositories.IFloorRepository
	repositories.IRoomRepository
	repositories.IBanRepsitory

	IPermissionCheckerService

	EventPublisher realtime.Publisher

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
	eventPublisher realtime.Publisher,
	pool *pgxpool.Pool,
) IRoomService {
	return &roomService{
		hallRepo,
		floorRepo,
		roomRepo,
		banRepo,
		permissionChecker,
		eventPublisher,
		pool,
		time.Duration(2) * time.Second,
		sync.RWMutex{},
	}
}

// ── helpers ───────────────────────────────────────────────────────────────────

func roomToRes(r *models.Room) dto.RoomRes {
	return dto.RoomRes{
		ID:                   r.ID,
		HallID:               r.HallID,
		FloorID:              r.FloorID,
		Name:                 r.Name,
		RoomType:             r.RoomType,
		Position:             r.Position,
		IsPrivate:            r.IsPrivate,
		SyncWithFloorMembers: r.SyncWithFloorMembers,
		CreatedAt:            r.CreatedAt,
		UpdatedAt:            r.UpdatedAt,
	}
}

func roomMemberToRes(m *models.HallMember) *dto.RoomMemberRes {
	return &dto.RoomMemberRes{
		ID:        m.ID,
		HallID:    m.HallID,
		UserID:    m.UserID,
		RoleID:    m.RoleID,
		Nickname:  m.Nickname,
		JoinedAt:  m.JoinedAt,
		CreatedAt: m.CreatedAt,
		UpdatedAt: m.UpdatedAt,
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
		if utils.IsDeadline(err) {
			return nil, utils.ErrorRequestTimeout
		}
		return nil, utils.ErrorFetchingHall
	}

	// If a floor is specified, verify it belongs to this hall
	var isFloorPrivate bool
	if req.FloorID != nil {
		floorExists, err := s.IFloorRepository.DoesFloorExistInHall(ctx, runner, *req.FloorID, hallID)
		if err != nil || !floorExists {
			if utils.IsDeadline(err) {
				return nil, utils.ErrorRequestTimeout
			}
			return nil, utils.ErrorFloorNotFound
		}

		isFloorPrivate, err = s.IFloorRepository.IsFloorPrivate(ctx, runner, *req.FloorID)
		if err != nil {
			if utils.IsDeadline(err) {
				return nil, utils.ErrorRequestTimeout
			}
			return nil, utils.ErrorFetchingFloor
		}
	}

	canonRoomType, err := utils.ParseRoomType(req.RoomType)
	if err != nil {
		return nil, err
	}

	maxPos, err := s.IRoomRepository.GetMaxPositionInContainer(ctx, runner, hallID, req.FloorID)
	if err != nil {
		if utils.IsDeadline(err) {
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

	// Sync to Floor is floor is private
	syncWithFloorMembers := false

	// Rule:
	// Any room inside a private floor must also be private
	// and should sync with floor members by default.
	if req.FloorID != nil && isFloorPrivate {
		isPrivate = true
		syncWithFloorMembers = true
	}

	now := time.Now()
	created, err := s.IRoomRepository.CreateRoom(ctx, runner, &models.Room{
		ID:                   roomID,
		HallID:               hallID,
		FloorID:              req.FloorID,
		Name:                 canonName,
		RoomType:             string(canonRoomType),
		Position:             maxPos + 1000.0,
		IsPrivate:            isPrivate,
		SyncWithFloorMembers: syncWithFloorMembers,
		CreatedAt:            now,
		UpdatedAt:            now,
	})
	if err != nil {
		if utils.IsDeadline(err) {
			return nil, utils.ErrorRequestTimeout
		}
		return nil, utils.ErrorCreatingRoom
	}

	// Important fix:
	// If room is inside a private floor, make sure the creator is part of floor_members
	// before syncing room_members from floor_members.
	if created.FloorID != nil && isFloorPrivate {
		creatorMemberID, err := s.IHallRepository.GetHallMemberID(ctx, runner, hallID, userInfo.ID)
		if err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				return nil, utils.ErrorUserDoesntBelongHall
			}
			if utils.IsDeadline(err) {
				return nil, utils.ErrorRequestTimeout
			}
			return nil, utils.ErrorFetchingHallMembers
		}

		// Idempotent because repository uses ON CONFLICT DO NOTHING.
		if err := s.IFloorRepository.AddFloorMember(ctx, runner, *created.FloorID, creatorMemberID); err != nil {
			if utils.IsDeadline(err) {
				return nil, utils.ErrorRequestTimeout
			}
			return nil, utils.ErrorCreatingFloorMember
		}

		// Sync all synced private rooms in this floor, including the newly created room.
		if err := s.IRoomRepository.SyncRoomsInFloorFromFloorMembers(ctx, runner, *created.FloorID); err != nil {
			if utils.IsDeadline(err) {
				return nil, utils.ErrorRequestTimeout
			}
			return nil, utils.ErrorCreatingRoomMember
		}

	} else if created.IsPrivate {
		// Standalone private room OR private room inside a public floor.
		// At minimum, room creator should have access.

		creatorMemberID, err := s.IHallRepository.GetHallMemberID(ctx, runner, hallID, userInfo.ID)
		if err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				return nil, utils.ErrorUserDoesntBelongHall
			}
			if utils.IsDeadline(err) {
				return nil, utils.ErrorRequestTimeout
			}
			return nil, utils.ErrorInternal
		}

		if err := s.IRoomRepository.AddRoomMember(ctx, runner, created.ID, creatorMemberID); err != nil {
			if utils.IsDeadline(err) {
				return nil, utils.ErrorRequestTimeout
			}
			return nil, utils.ErrorCreatingRoomMember
		}
	}

	if err := runner.Commit(ctx); err != nil {
		return nil, utils.ErrorInternal
	}

	// PUBLISH EVENT
	// For private rooms, hub will resync this user.
	// For public rooms, hub will subscribe hall clients directly.
	publishHubEvent(s.EventPublisher, realtime.HubEvent{
		Type:      realtime.HubEventRoomCreated,
		HallID:    created.HallID,
		RoomID:    created.ID,
		UserID:    userInfo.ID,
		IsPrivate: created.IsPrivate,
	})

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
		if utils.IsDeadline(err) {
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
		if utils.IsDeadline(rr.err) {
			return nil, utils.ErrorRequestTimeout
		}
		return nil, utils.ErrorFetchingRoom
	}
	if fr.err != nil {
		if utils.IsDeadline(fr.err) {
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

	fields := make(map[string]any)

	if req.Name != nil {
		canon, err := utils.SanitizeFloorname(*req.Name)
		if err != nil {
			return nil, err
		}
		fields["name"] = canon
	}

	if req.IsPrivate != nil {
		fields["is_private"] = *req.IsPrivate
	}

	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return nil, utils.ErrorInternal
	}
	runner := database.NewTxWrapper(tx)
	defer runner.Rollback(ctx)

	// check if appropriate permissions to manage room and things
	if err := s.requireManageServers(ctx, runner, userInfo.ID, hallID); err != nil {
		return nil, err
	}

	// Verify room belongs to this hall before updating
	room, err := s.IRoomRepository.GetRoomByID(ctx, runner, roomID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, utils.ErrorRoomNotFound
		}
		if utils.IsDeadline(err) {
			return nil, utils.ErrorRequestTimeout
		}
		return nil, utils.ErrorFetchingRoom
	}
	if room.HallID != hallID {
		return nil, utils.ErrorRoomNotFound
	}

	// Check if the room is inside a private floor
	// in that case, cannot update it to public sorry mate
	if room.FloorID != nil && req.IsPrivate != nil && !*req.IsPrivate {
		isFloorPrivate, err := s.IFloorRepository.IsFloorPrivate(ctx, runner, *room.FloorID)
		if err != nil {
			if utils.IsDeadline(err) {
				return nil, utils.ErrorRequestTimeout
			}
			return nil, utils.ErrorFetchingFloor
		}

		if isFloorPrivate {
			return nil, utils.ErrorCannotMakeRoomInsidePrivateFloorPublic
		}
	}

	updated, err := s.IRoomRepository.UpdateRoom(ctx, runner, roomID, fields)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, utils.ErrorRoomNotFound
		}
		if utils.IsDeadline(err) {
			return nil, utils.ErrorRequestTimeout
		}
		return nil, utils.ErrorFetchingRoom
	}

	// Privacy transition handling
	if req.IsPrivate != nil {
		wasPrivate := room.IsPrivate
		nowPrivate := *req.IsPrivate

		// public -> private
		creatorMemberID, err := s.IHallRepository.GetHallMemberID(ctx, runner, hallID, userInfo.ID)
		if err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				return nil, utils.ErrorUserDoesntBelongHall
			}
			if utils.IsDeadline(err) {
				return nil, utils.ErrorRequestTimeout
			}
			return nil, utils.ErrorFetchingHallMembers
		}

		if err := s.IRoomRepository.AddRoomMember(ctx, runner, roomID, creatorMemberID); err != nil {
			if utils.IsDeadline(err) {
				return nil, utils.ErrorRequestTimeout
			}
			return nil, utils.ErrorCreatingRoomMember
		}

		// private -> public
		if wasPrivate && !nowPrivate {
			if err := s.IRoomRepository.ClearRoomMembers(ctx, runner, roomID); err != nil {
				if utils.IsDeadline(err) {
					return nil, utils.ErrorRequestTimeout
				}
				return nil, utils.ErrorFetchingRoom
			}
		}
	}

	if err := runner.Commit(ctx); err != nil {
		return nil, utils.ErrorInternal
	}

	// PUBLISH EVENT
	// room's privacy status change checker
	// if changed publish event HubEventRoomPrivacyChanged
	if room.IsPrivate != updated.IsPrivate {
		publishHubEvent(s.EventPublisher, realtime.HubEvent{
			Type:      realtime.HubEventRoomPrivacyChanged,
			HallID:    updated.HallID,
			RoomID:    updated.ID,
			IsPrivate: updated.IsPrivate,
		})
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
		if utils.IsDeadline(err) {
			return utils.ErrorRequestTimeout
		}
		return utils.ErrorFetchingRoom
	}
	if room.HallID != hallID {
		return utils.ErrorRoomNotFound
	}

	if err := s.IRoomRepository.DeleteRoom(ctx, runner, roomID); err != nil {
		if utils.IsDeadline(err) {
			return utils.ErrorRequestTimeout
		}
		return utils.ErrorInternal
	}

	if err := runner.Commit(ctx); err != nil {
		return utils.ErrorInternal
	}

	// PUBLISH EVENT
	publishHubEvent(s.EventPublisher, realtime.HubEvent{
		Type:   realtime.HubEventRoomDeleted,
		HallID: room.HallID,
		RoomID: room.ID,
	})

	return nil

}

// ── MoveRoom ──────────────────────────────────────────────────────────────────

func (s *roomService) MoveRoom(c context.Context, userInfo *auth.UserInfo, hallID uuid.UUID, roomID uuid.UUID, req *dto.MoveRoomReq) (*dto.RoomRes, error) {
	ctx, cancel := context.WithTimeout(c, s.timeout)
	defer cancel()

	// Invalid conflicting move instructions
	if req.MoveToTopLevel != nil &&
		*req.MoveToTopLevel &&
		req.NewFloorID != nil {

		return nil, utils.ErrorInvalidMoveRoomPayload
	}

	// nothing to move
	if req.MoveToTopLevel == nil &&
		req.NewFloorID == nil &&
		req.AfterID == nil {
		return nil, utils.ErrorNoFieldsToUpdate
	}

	// sync_private only makes sense when moving into a floor
	if req.SyncPrivate != nil && *req.SyncPrivate && req.NewFloorID == nil {
		return nil, utils.ErrorInvalidMoveRoomPayload
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

	// Verify room belongs to this hall
	room, err := s.IRoomRepository.GetRoomByID(ctx, runner, roomID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, utils.ErrorRoomNotFound
		}
		if utils.IsDeadline(err) {
			return nil, utils.ErrorRequestTimeout
		}
		return nil, utils.ErrorFetchingRoom
	}

	if room.HallID != hallID {
		return nil, utils.ErrorRoomNotFound
	}

	newFloorID := room.FloorID

	// Verify target floor belongs to this hall (if targeting a floor)
	if req.MoveToTopLevel != nil && *req.MoveToTopLevel {
		newFloorID = nil
	}

	if req.NewFloorID != nil {
		floorExists, err := s.IFloorRepository.DoesFloorExistInHall(ctx, runner, *req.NewFloorID, hallID)
		if err != nil || !floorExists {
			if utils.IsDeadline(err) {
				return nil, utils.ErrorRequestTimeout
			}
			return nil, utils.ErrorFloorNotFound
		}

		newFloorID = req.NewFloorID
	}

	lower, upper, err := s.IRoomRepository.GetRoomPositionBounds(ctx, runner, hallID, newFloorID, req.AfterID)
	if err != nil {
		if utils.IsDeadline(err) {
			return nil, utils.ErrorRequestTimeout
		}
		return nil, utils.ErrorFetchingRoom
	}

	newPos := utils.CalcPosition(lower, upper)

	moved, err := s.IRoomRepository.MoveRoom(ctx, runner, roomID, newFloorID, newPos)
	if err != nil {
		if utils.IsDeadline(err) {
			return nil, utils.ErrorRequestTimeout
		}
		return nil, utils.ErrorMovingRoom
	}

	if req.NewFloorID != nil && req.SyncPrivate != nil && *req.SyncPrivate {
		isFloorPrivate, err := s.IFloorRepository.IsFloorPrivate(ctx, runner, *req.NewFloorID)
		if err != nil {
			if utils.IsDeadline(err) {
				return nil, utils.ErrorRequestTimeout
			}
			return nil, utils.ErrorFetchingFloor
		}

		if isFloorPrivate {
			if err := s.IRoomRepository.ClearRoomMembers(ctx, runner, roomID); err != nil {
				if utils.IsDeadline(err) {
					return nil, utils.ErrorRequestTimeout
				}
				return nil, utils.ErrorFetchingRoom
			}

			if err := s.IRoomRepository.SyncRoomMembersFromFloor(ctx, runner, roomID, *req.NewFloorID); err != nil {
				if utils.IsDeadline(err) {
					return nil, utils.ErrorRequestTimeout
				}
				return nil, utils.ErrorCreatingRoomMember
			}

			private := true
			moved, err = s.IRoomRepository.UpdateRoom(ctx, runner, roomID, map[string]any{
				"is_private": private,
			})
			if err != nil {
				if utils.IsDeadline(err) {
					return nil, utils.ErrorRequestTimeout
				}
				return nil, utils.ErrorFetchingRoom
			}
		}
	}

	// If moved out of floor, do nothing to room_members.
	// Existing room_members are retained.

	if err := runner.Commit(ctx); err != nil {
		return nil, utils.ErrorInternal
	}

	// PUBLISH EVENT
	publishHubEvent(s.EventPublisher, realtime.HubEvent{
		Type:   realtime.HubEventRoomMoved,
		HallID: moved.HallID,
		RoomID: moved.ID,
	})

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
	runner := database.NewConnWrapper(conn)

	room, err := s.IRoomRepository.GetRoomByID(ctx, runner, roomID)
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
	runner := database.NewConnWrapper(conn)

	return s.IRoomRepository.IsUserRoomMember(ctx, runner, roomID, userID)
}

func (s *roomService) GetAccessibleRoomsForUser(c context.Context, userInfo *auth.UserInfo) (map[uuid.UUID]uuid.UUID, error) {

	ctx, cancel := context.WithTimeout(c, s.timeout)
	defer cancel()

	conn, err := s.pool.Acquire(ctx)
	if err != nil {
		return nil, utils.ErrorInternal
	}
	defer conn.Release()
	runner := database.NewConnWrapper(conn)

	// fetch all the hall the user is associated with
	hallIDs, err := s.IHallRepository.GetUserHallIDs(ctx, runner, userInfo.ID)
	if err != nil {
		return nil, utils.ErrorInternal
	}

	accessibleRooms := make(map[uuid.UUID]uuid.UUID)

	// range through the hall
	// fetch all the room associated with the hallID
	for _, hallID := range hallIDs {

		// fetch all roomID of current hall

		// full room struct brings too much information
		// useing []*RoomIDandPrivate which only contains RoomID and IsPrivate
		roomInfo, err := s.IRoomRepository.GetRoomsIDandPrivateInfoByHallID(ctx, runner, hallID)
		if err != nil {
			if utils.IsDeadline(err) {
				return nil, utils.ErrorRequestTimeout
			}
			return nil, utils.ErrorFetchingRoom
		}

		// check if private
		// if private, check if accessible
		for _, rm := range roomInfo {
			if rm.IsPrivate {

				isRoomMember, err := s.IRoomRepository.IsUserRoomMember(ctx, runner, rm.RoomID, userInfo.ID)
				if err != nil {
					if utils.IsDeadline(err) {
						return nil, utils.ErrorRequestTimeout
					}
					return nil, utils.ErrorFetchingRoom

				}
				if !isRoomMember {
					// skip this room and do not make it accessible
					continue
				}

			}

			// map current RoomId -> hallID
			accessibleRooms[rm.RoomID] = hallID

		}

	}

	return accessibleRooms, nil
}

// ── Room Members ──────────────────────────────────────────────────────────────

func (s *roomService) AddRoomMember(c context.Context, userInfo *auth.UserInfo, hallID uuid.UUID, roomID uuid.UUID, memberID uuid.UUID) (*dto.RoomAccessMemberRes, error) {
	ctx, cancel := context.WithTimeout(c, s.timeout)
	defer cancel()

	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return nil, utils.ErrorInternal
	}
	runner := database.NewTxWrapper(tx)
	defer runner.Rollback(ctx)

	// check if appropriate permissions to manage room and things
	if err := s.requireManageServers(ctx, runner, userInfo.ID, hallID); err != nil {
		return nil, err
	}

	// Verify room belongs to this hall before adding member
	room, err := s.IRoomRepository.GetRoomByID(ctx, runner, roomID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, utils.ErrorRoomNotFound
		}
		if utils.IsDeadline(err) {
			return nil, utils.ErrorRequestTimeout
		}
		return nil, utils.ErrorFetchingRoom
	}

	if room.HallID != hallID {
		return nil, utils.ErrorRoomNotFound
	}

	// Got to be private for adding members
	if !room.IsPrivate {
		return nil, utils.ErrorRoomIsNotPrivate
	}

	// Verify member belongs to this hall
	member, err := s.IHallRepository.GetHallMemberByID(ctx, runner, hallID, memberID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, utils.ErrorMemberNotFound
		}
		if utils.IsDeadline(err) {
			return nil, utils.ErrorRequestTimeout
		}
		return nil, utils.ErrorInternal
	}

	// If a room inside a private floor has its member updated separately,
	// it is no longer in sync with the floor members.
	if room.FloorID != nil {
		if err := s.IRoomRepository.SetRoomFloorMemberSync(ctx, runner, roomID, false); err != nil {
			if utils.IsDeadline(err) {
				return nil, utils.ErrorRequestTimeout
			}
			return nil, utils.ErrorFetchingRoom
		}
	}

	if err := s.IRoomRepository.AddRoomMember(ctx, runner, roomID, memberID); err != nil {
		if utils.IsDeadline(err) {
			return nil, utils.ErrorRequestTimeout
		}
		return nil, utils.ErrorCreatingRoomMember
	}

	if err := runner.Commit(ctx); err != nil {
		return nil, utils.ErrorInternal
	}

	// PUBLISH EVENT
	publishHubEvent(s.EventPublisher, realtime.HubEvent{
		Type:     realtime.HubEventRoomMemberAdded,
		HallID:   hallID,
		RoomID:   roomID,
		MemberID: memberID,
		UserID:   member.UserID,
	})

	return &dto.RoomAccessMemberRes{
		RoomID:   roomID,
		MemberID: memberID,
	}, nil
}

func (s *roomService) RemoveRoomMember(c context.Context, userInfo *auth.UserInfo, hallID uuid.UUID, roomID uuid.UUID, memberID uuid.UUID) (*dto.RoomAccessMemberRes, error) {
	ctx, cancel := context.WithTimeout(c, s.timeout)
	defer cancel()

	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return nil, utils.ErrorInternal
	}
	runner := database.NewTxWrapper(tx)
	defer runner.Rollback(ctx)

	// check if appropriate permissions to manage room and things
	if err := s.requireManageServers(ctx, runner, userInfo.ID, hallID); err != nil {
		return nil, err
	}

	// Verify room belongs to this hall before removing member
	room, err := s.IRoomRepository.GetRoomByID(ctx, runner, roomID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, utils.ErrorRoomNotFound
		}
		if utils.IsDeadline(err) {
			return nil, utils.ErrorRequestTimeout
		}
		return nil, utils.ErrorFetchingRoom
	}

	if room.HallID != hallID {
		return nil, utils.ErrorRoomNotFound
	}

	// Got to be private for removing members
	if !room.IsPrivate {
		return nil, utils.ErrorRoomIsNotPrivate
	}

	// Verify member belongs to this hall
	member, err := s.IHallRepository.GetHallMemberByID(ctx, runner, hallID, memberID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, utils.ErrorMemberNotFound
		}
		if utils.IsDeadline(err) {
			return nil, utils.ErrorRequestTimeout
		}
		return nil, utils.ErrorInternal
	}

	// If a room inside a private floor has its member updated separately,
	// it is no longer in sync with the floor members.
	if room.FloorID != nil {
		if err := s.IRoomRepository.SetRoomFloorMemberSync(ctx, runner, roomID, false); err != nil {
			if utils.IsDeadline(err) {
				return nil, utils.ErrorRequestTimeout
			}
			return nil, utils.ErrorFetchingRoom
		}
	}

	if err := s.IRoomRepository.RemoveRoomMember(ctx, runner, roomID, memberID); err != nil {
		if utils.IsDeadline(err) {
			return nil, utils.ErrorRequestTimeout
		}
		return nil, utils.ErrorDeletingRoomMember
	}

	if err := runner.Commit(ctx); err != nil {
		return nil, utils.ErrorInternal
	}

	// PUBLISH EVENT
	publishHubEvent(s.EventPublisher, realtime.HubEvent{
		Type:     realtime.HubEventRoomMemberRemoved,
		HallID:   hallID,
		RoomID:   roomID,
		MemberID: memberID,
		UserID:   member.UserID,
	})

	return &dto.RoomAccessMemberRes{
		RoomID:   roomID,
		MemberID: memberID,
	}, nil
}

func (s *roomService) SyncRoomMembersToFloor(c context.Context, userInfo *auth.UserInfo, hallID uuid.UUID, roomID uuid.UUID) (*dto.RoomRes, error) {
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
		if utils.IsDeadline(err) {
			return nil, utils.ErrorRequestTimeout
		}
		return nil, utils.ErrorFetchingRoom
	}

	if room.HallID != hallID {
		return nil, utils.ErrorRoomNotFound
	}

	// Room must be inside a floor to sync with floor members
	if room.FloorID == nil {
		return nil, utils.ErrorRoomNotInFloor
	}

	// Room must be private to have room_members
	if !room.IsPrivate {
		return nil, utils.ErrorRoomIsNotPrivate
	}

	// Floor must be private, otherwise floor_members should not control access
	isFloorPrivate, err := s.IFloorRepository.IsFloorPrivate(ctx, runner, *room.FloorID)
	if err != nil {
		if utils.IsDeadline(err) {
			return nil, utils.ErrorRequestTimeout
		}
		return nil, utils.ErrorFetchingFloor
	}

	if !isFloorPrivate {
		return nil, utils.ErrorFloorIsNotPrivate
	}

	// Replace this room's member list with the parent floor's member list
	if err := s.IRoomRepository.ReplaceRoomMembersFromFloor(ctx, runner, roomID, *room.FloorID); err != nil {
		if utils.IsDeadline(err) {
			return nil, utils.ErrorRequestTimeout
		}
		return nil, utils.ErrorCreatingRoomMember
	}

	updated, err := s.IRoomRepository.UpdateRoom(ctx, runner, roomID, map[string]any{
		"sync_with_floor_members": true,
	})
	if err != nil {
		if utils.IsDeadline(err) {
			return nil, utils.ErrorRequestTimeout
		}
		return nil, utils.ErrorFetchingRoom
	}

	if err := runner.Commit(ctx); err != nil {
		return nil, utils.ErrorInternal
	}

	// PUBLISH EVENT
	publishHubEvent(s.EventPublisher, realtime.HubEvent{
		Type:   realtime.HubEventHallAccessResync,
		HallID: hallID,
		RoomID: roomID,
	})

	res := roomToRes(updated)
	return &res, nil
}

func (s *roomService) GetRoomMembers(c context.Context, userInfo *auth.UserInfo, hallID uuid.UUID, roomID uuid.UUID) (*dto.GetRoomMembersRes, error) {
	ctx, cancel := context.WithTimeout(c, s.timeout)
	defer cancel()

	conn, err := s.pool.Acquire(ctx)
	if err != nil {
		return nil, utils.ErrorInternal
	}
	defer conn.Release()
	runner := database.NewConnWrapper(conn)

	if err := s.requireManageServers(ctx, runner, userInfo.ID, hallID); err != nil {
		return nil, err
	}

	room, err := s.IRoomRepository.GetRoomByID(ctx, runner, roomID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, utils.ErrorRoomNotFound
		}
		if utils.IsDeadline(err) {
			return nil, utils.ErrorRequestTimeout
		}
		return nil, utils.ErrorFetchingRoom
	}

	if room.HallID != hallID {
		return nil, utils.ErrorRoomNotFound
	}

	if !room.IsPrivate {
		return nil, utils.ErrorRoomIsNotPrivate
	}

	members, err := s.IRoomRepository.ListRoomMembers(ctx, runner, hallID, roomID)
	if err != nil {
		if utils.IsDeadline(err) {
			return nil, utils.ErrorRequestTimeout
		}
		return nil, utils.ErrorFetchingHallMembers
	}

	out := make([]*dto.RoomMemberRes, 0, len(members))
	for _, m := range members {
		out = append(out, roomMemberToRes(m))
	}

	return &dto.GetRoomMembersRes{
		Members: out,
		Total:   len(out),
	}, nil
}

func (s *roomService) GetRoomMember(c context.Context, userInfo *auth.UserInfo, hallID uuid.UUID, roomID uuid.UUID, memberID uuid.UUID) (*dto.GetRoomMemberRes, error) {
	ctx, cancel := context.WithTimeout(c, s.timeout)
	defer cancel()

	conn, err := s.pool.Acquire(ctx)
	if err != nil {
		return nil, utils.ErrorInternal
	}
	defer conn.Release()
	runner := database.NewConnWrapper(conn)

	if err := s.requireManageServers(ctx, runner, userInfo.ID, hallID); err != nil {
		return nil, err
	}

	room, err := s.IRoomRepository.GetRoomByID(ctx, runner, roomID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, utils.ErrorRoomNotFound
		}
		if utils.IsDeadline(err) {
			return nil, utils.ErrorRequestTimeout
		}
		return nil, utils.ErrorFetchingRoom
	}

	if room.HallID != hallID {
		return nil, utils.ErrorRoomNotFound
	}

	if !room.IsPrivate {
		return nil, utils.ErrorRoomIsNotPrivate
	}

	member, err := s.IRoomRepository.GetRoomMember(ctx, runner, hallID, roomID, memberID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, utils.ErrorMemberNotFound
		}
		if utils.IsDeadline(err) {
			return nil, utils.ErrorRequestTimeout
		}
		return nil, utils.ErrorFetchingHallMembers
	}

	return roomMemberToRes(member), nil
}
