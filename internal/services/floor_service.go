package services

import (
	"context"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/suck-seed/yapp/internal/dto"
	"github.com/suck-seed/yapp/internal/models"
	"github.com/suck-seed/yapp/internal/utils"

	"github.com/suck-seed/yapp/internal/repositories"
)

type IFloorService interface {
	CreateFloor(c context.Context, req *dto.CreateFloorReq) (*dto.CreateFloorRes, error)
}

type floorService struct {
	repositories.IHallRepository
	repositories.IFloorRepository
	timeout time.Duration
	mu      sync.RWMutex
}

func NewFloorService(hallRepo repositories.IHallRepository, floorRepo repositories.IFloorRepository) IFloorService {
	return &floorService{
		hallRepo,
		floorRepo,
		time.Duration(2) * time.Second,
		sync.RWMutex{},
	}
}

// ctx, cancel := context.WithTimeout(c, s.timeout)
// 	defer cancel()

func (s *floorService) CreateFloor(c context.Context, req *dto.CreateFloorReq) (*dto.CreateFloorRes, error) {

	ctx, cancel := context.WithTimeout(c, s.timeout)
	defer cancel()

	//	canon name
	canonName, err := utils.SanitizeFloorname(req.Name)
	if err != nil {
		return nil, err
	}

	//	valid hallID
	exists, err := s.IHallRepository.DoesHallExist(ctx, req.HallID)
	if err != nil || !*exists {
		return nil, utils.ErrorHallDoesntExist
	}

	//	create a floorId
	floorID, err := uuid.NewV7()
	if err != nil {
		return nil, utils.ErrorInternal
	}

	floor := &models.Floor{
		ID:        floorID,
		HallID:    req.HallID,
		Name:      canonName,
		IsPrivate: req.IsPrivate,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	floorCRES, err := s.IFloorRepository.CreateFloor(ctx, floor)
	if err != nil {
		return nil, utils.ErrorCreatingFloor
	}

	return &dto.CreateFloorRes{
		ID:        floorCRES.ID,
		HallID:    floorCRES.HallID,
		Name:      floorCRES.Name,
		IsPrivate: floorCRES.IsPrivate,
		CreatedAt: floorCRES.CreatedAt,
		UpdatedAt: floorCRES.UpdatedAt,
	}, nil
}
