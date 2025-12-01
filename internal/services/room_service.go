package services

import (
	"context"
	"sync"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/suck-seed/yapp/internal/database"
	"github.com/suck-seed/yapp/internal/dto"
	"github.com/suck-seed/yapp/internal/utils"

	"github.com/google/uuid"
	"github.com/suck-seed/yapp/internal/models"
	"github.com/suck-seed/yapp/internal/repositories"
)

type IRoomService interface {
	CreateRoom(c context.Context, req *dto.CreateRoomReq) (*dto.CreateRoomRes, error)
	GetRoomByID(c context.Context, roomId *uuid.UUID) (*models.Room, error)

	IsUserRoomMember(c context.Context, roomId *uuid.UUID, userId *uuid.UUID) (*bool, error)
}

type roomService struct {
	repositories.IHallRepository
	repositories.IFloorRepository
	repositories.IRoomRepository
	pool *pgxpool.Pool

	timeout time.Duration
	mu      sync.RWMutex
}

func NewRoomService(hallRepo repositories.IHallRepository, floorRepo repositories.IFloorRepository, roomRepo repositories.IRoomRepository, pool *pgxpool.Pool) IRoomService {
	return &roomService{
		hallRepo,
		floorRepo,
		roomRepo,
		pool,
		time.Duration(2) * time.Second,
		sync.RWMutex{},
	}
}

// ctx, cancel := context.WithTimeout(c, s.timeout)
// 	defer cancel()

func (s *roomService) CreateRoom(c context.Context, req *dto.CreateRoomReq) (*dto.CreateRoomRes, error) {

	ctx, cancel := context.WithTimeout(c, s.timeout)
	defer cancel()

	// ---------- TRANSACTION INIT
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
	if err != nil {
		return nil, utils.ErrorHallDoesntExist
	}
	if !*hallExists {
		return nil, utils.ErrorHallDoesntExist
	}

	if req.FloorID != nil {
		floorExists, err := s.IFloorRepository.DoesFloorExistsInRoom(ctx, runner, req.FloorID, &req.HallID)
		if err != nil {
			return nil, utils.ErrorFloorDoesntExistInHall
		}

		if !floorExists {
			return nil, utils.ErrorFloorDoesntExistInHall
		}
	}

	//	this can be done better, prolly enums
	canonRoomType, err := utils.ParseRoomType(req.RoomType)
	if err != nil {
		return nil, err
	}

	roomID, err := uuid.NewV7()
	if err != nil {
		return nil, utils.ErrorInternal
	}

	roomCRES, err := s.IRoomRepository.CreateRoom(ctx, runner, &models.Room{
		ID:        roomID,
		HallId:    req.HallID,
		FloorId:   req.FloorID,
		Name:      canonName,
		RoomType:  string(canonRoomType),
		IsPrivate: req.IsPrivate,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	})
	if err != nil {
		return nil, utils.ErrorCreatingRoom
	}

	// Commit before returning data to handler
	if err := runner.Commit(ctx); err != nil {
		return nil, utils.ErrorInternal
	}

	return &dto.CreateRoomRes{
		ID:        roomCRES.ID,
		HallID:    roomCRES.HallId,
		FloorID:   roomCRES.FloorId,
		Name:      roomCRES.Name,
		RoomType:  roomCRES.RoomType,
		IsPrivate: roomCRES.IsPrivate,
		CreatedAt: roomCRES.CreatedAt,
		UpdatedAt: roomCRES.UpdatedAt,
	}, nil
}

func (s *roomService) GetRoomByID(c context.Context, rooomId *uuid.UUID) (*models.Room, error) {

	ctx, cancel := context.WithTimeout(c, s.timeout)
	defer cancel()

	conn, err := s.pool.Acquire(ctx)
	if err != nil {
		return nil, utils.ErrorInternal
	}

	defer conn.Release()
	runner := database.NewConnWrapper(conn)

	roomCRES, err := s.IRoomRepository.GetRoomByID(ctx, runner, rooomId)
	if err != nil {
		return nil, utils.ErrorFetchingRoom
	}

	return roomCRES, nil

}

func (s *roomService) IsUserRoomMember(c context.Context, roomId *uuid.UUID, userId *uuid.UUID) (*bool, error) {

	ctx, cancel := context.WithTimeout(c, s.timeout)
	defer cancel()

	conn, err := s.pool.Acquire(ctx)
	if err != nil {
		return nil, utils.ErrorInternal
	}

	defer conn.Release()
	runner := database.NewConnWrapper(conn)

	isMember, err := s.IRoomRepository.IsUserRoomMember(ctx, runner, roomId, userId)
	if err != nil {
		return nil, err
	}

	return isMember, nil
}
