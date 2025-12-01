package services

import (
	"context"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/suck-seed/yapp/internal/auth"
	"github.com/suck-seed/yapp/internal/database"
	"github.com/suck-seed/yapp/internal/dto"
	"github.com/suck-seed/yapp/internal/models"
	"github.com/suck-seed/yapp/internal/repositories"
	"github.com/suck-seed/yapp/internal/utils"
)

type IHallService interface {
	CreateHall(c context.Context, req *dto.CreateHallReq) (*dto.CreateHallRes, error)
	IsUserHallMember(c context.Context, hallID *uuid.UUID, userId *uuid.UUID) (*bool, error)
	DoesHallExist(c context.Context, HallId *uuid.UUID) (*bool, error)

	GetUserHalls(c context.Context) ([]*models.Hall, error)
}

type hallService struct {
	repositories.IHallRepository
	pool *pgxpool.Pool

	timeout time.Duration
	mu      sync.RWMutex
}

func NewHallService(hallRepo repositories.IHallRepository, pool *pgxpool.Pool) IHallService {
	return &hallService{
		hallRepo,
		pool,
		time.Duration(2) * time.Second,
		sync.RWMutex{},
	}
}

func (s *hallService) CreateHall(c context.Context, req *dto.CreateHallReq) (*dto.CreateHallRes, error) {
	ctx, cancel := context.WithTimeout(c, s.timeout)
	defer cancel()

	// Begin a Transaction
	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return nil, utils.ErrorInternal
	}

	// Runner from aquired transaction
	runner := database.NewTxWrapper(tx)

	// Rollback if anamoly occurs
	defer runner.Rollback(ctx)

	// sanatize req
	canonHallname, err := utils.SanitizeHallname(req.Name)
	if err != nil {
		return nil, err
	}
	canonBannerColor, err := utils.SanitizeColorFormat(req.BannerColor)
	if err != nil {
		return nil, err
	}
	canonDescription, err := utils.SanitizeText(req.Description)
	if err != nil {
		return nil, err
	}

	// get userId from context.Context()
	userId, _, err := auth.CurrentUserFromContext(c)
	if err != nil {
		return nil, err
	}

	//
	// Hall Creation
	//

	// generate hall id
	hallId, err := uuid.NewV7()
	if err != nil {
		return nil, utils.ErrorInternal
	}

	// package a hall struct

	newHall := &models.Hall{
		ID:               hallId,
		Name:             canonHallname,
		IconURL:          req.IconURL,
		IconThumbnailURL: req.IconThumbnailURL,
		BannerColor:      canonBannerColor,
		Description:      canonDescription,
		OwnerID:          *userId,
	}

	// pass to repo
	hall, err := s.IHallRepository.CreateHall(ctx, runner, newHall)
	if err != nil {
		return nil, utils.ErrorCreatingHall
	}

	// generate role id
	roleId, err := uuid.NewV7()
	if err != nil {
		return nil, utils.ErrorInternal
	}

	// package a role struct
	newRole := &models.Role{
		ID:      roleId,
		HallID:  hall.ID,
		Name:    "creator",
		IsAdmin: true,
	}

	// pass to repo
	role, err := s.IHallRepository.CreateHallRole(ctx, runner, newRole)
	if err != nil {
		return nil, utils.ErrorCreatingHallRole
	}

	//
	// Hall Member Creation
	//

	// generate hall member id
	hallMemberID, err := uuid.NewV7()
	if err != nil {
		return nil, utils.ErrorInternal
	}

	// package a hall-member struct
	newHallMember := &models.HallMember{
		ID:     hallMemberID,
		HallID: hall.ID,
		UserID: *userId,
		RoleID: role.ID,
	}

	// pass to repo
	err = s.IHallRepository.CreateHallMember(ctx, runner, newHallMember)
	if err != nil {
		return nil, utils.ErrorCreatingHallMember
	}

	// Commit before returning data to handler
	if err := runner.Commit(ctx); err != nil {
		return nil, utils.ErrorInternal
	}

	return &dto.CreateHallRes{
		ID:               hall.ID,
		Name:             hall.Name,
		IconURL:          hall.IconURL,
		IconThumbnailURL: hall.IconThumbnailURL,
		BannerColor:      hall.BannerColor,
		Description:      hall.Description,
		CreatedAt:        hall.CreatedAt,
		UpdatedAt:        hall.UpdatedAt,
		OwnerID:          hall.OwnerID,
	}, nil
}

func (s *hallService) GetUserHalls(c context.Context) ([]*models.Hall, error) {
	ctx, cancel := context.WithTimeout(c, s.timeout)
	defer cancel()

	// Acquire a conn from pool
	conn, err := s.pool.Acquire(ctx)
	if err != nil {
		return nil, utils.ErrorInternal
	}

	// Release after everything is done
	defer conn.Release()

	// Runner from the acquired conn
	runner := database.NewConnWrapper(conn)

	// get userId from context.Context()
	userId, _, err := auth.CurrentUserFromContext(c)
	if err != nil {
		return nil, err
	}

	hallIds, err := s.IHallRepository.GetUserHallIDs(ctx, runner, *userId)
	if err != nil {
		return nil, utils.ErrorInternal
	}

	var halls []*models.Hall
	for _, hallId := range hallIds {
		hall, err := s.IHallRepository.GetHallByID(ctx, runner, hallId)
		if err != nil {
			return nil, utils.ErrorInternal
		}
		halls = append(halls, hall)
	}

	return halls, nil
}

func (s *hallService) IsUserHallMember(c context.Context, hallID *uuid.UUID, userID *uuid.UUID) (*bool, error) {
	ctx, cancel := context.WithTimeout(c, s.timeout)
	defer cancel()

	// Acquire a conn from pool
	conn, err := s.pool.Acquire(ctx)
	if err != nil {
		return nil, utils.ErrorInternal
	}

	// Release after everything is done
	defer conn.Release()

	// Runner from the acquired conn
	runner := database.NewConnWrapper(conn)

	isMember, err := s.IHallRepository.IsUserHallMember(ctx, runner, *hallID, *userID)
	if err != nil {
		return nil, utils.ErrorInternal
	}

	return isMember, nil
}

func (s *hallService) DoesHallExist(c context.Context, HallId *uuid.UUID) (*bool, error) {

	ctx, cancel := context.WithTimeout(c, s.timeout)
	defer cancel()

	// Acquire a conn from pool
	conn, err := s.pool.Acquire(ctx)
	if err != nil {
		return nil, utils.ErrorInternal
	}

	// Release after everything is done
	defer conn.Release()

	// Runner from the acquired conn
	runner := database.NewConnWrapper(conn)

	exists, err := s.IHallRepository.DoesHallExist(ctx, runner, *HallId)
	if err != nil {
		return nil, utils.ErrorInternal
	}

	return exists, nil

}
