package services

import (
	"context"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/suck-seed/yapp/internal/auth"
	"github.com/suck-seed/yapp/internal/dto"
	"github.com/suck-seed/yapp/internal/models"
	"github.com/suck-seed/yapp/internal/repositories"
	"github.com/suck-seed/yapp/internal/utils"
)

type IHallService interface {
	CreateHall(c context.Context, req *dto.CreateHallReq) (*dto.CreateHallRes, error)
	IsMember(c context.Context, hallID *uuid.UUID, userId *uuid.UUID) (bool, error)
}

type hallService struct {
	repositories.IHallRepository
	timeout time.Duration
	mu      sync.RWMutex
}

func NewHallService(hallRepo repositories.IHallRepository) IHallService {
	return &hallService{
		hallRepo,
		time.Duration(2) * time.Second,
		sync.RWMutex{},
	}
}

// ctx, cancel := context.WithTimeout(c, s.timeout)
// 	defer cancel()

func (s *hallService) CreateHall(c context.Context, req *dto.CreateHallReq) (*dto.CreateHallRes, error) {

	ctx, cancel := context.WithTimeout(c, s.timeout)
	defer cancel()

	// sanatize req
	canonHallname, err := utils.SanatizeHallname(req.Name)
	if err != nil {
		return nil, err
	}
	canonBannerColor, err := utils.SanatizeColorFormat(req.BannerColor)
	if err != nil {
		return nil, err
	}
	canonDescription, err := utils.SanatizeText(req.Description)
	if err != nil {
		return nil, err
	}

	// get userId from context.Context()
	userIdString, _, err := auth.CurrentUserFromContext(c)
	if err != nil {
		return nil, err
	}

	userId, err := uuid.Parse(userIdString)
	if err != nil {
		return nil, utils.ErrorInvalidUserIdInContext
	}

	// check if hallName already existing hallname
	hallByName, _ := s.IHallRepository.GetHallByName(ctx, canonHallname)

	if hallByName != nil {
		return nil, utils.ErrorHallAlreadyExist
	}

	// generate id
	id, err := uuid.NewV7()
	if err != nil {
		return nil, utils.ErrorInternal
	}

	// package a hall struct
	hall := &models.Hall{
		ID:          id,
		Name:        canonHallname,
		IconURL:     req.IconURL,
		BannerColor: canonBannerColor,
		Description: canonDescription,
		Owner:       userId,
	}

	// passing to repo
	hallCRES, err := s.IHallRepository.CreateHall(ctx, hall)
	if err != nil {
		return nil, utils.ErrorCreatingHall
	}

	// additional setup

	return &dto.CreateHallRes{
		ID:          hallCRES.ID,
		Name:        hallCRES.Name,
		IconURL:     hallCRES.IconURL,
		BannerColor: hallCRES.BannerColor,
		Description: hallCRES.Description,
		CreatedAt:   hall.CreatedAt,
		UpdatedAt:   hall.UpdatedAt,
		Owner:       hallCRES.Owner,
	}, nil
}

func (s *hallService) IsMember(c context.Context, hallID *uuid.UUID, userId *uuid.UUID) (bool, error) {

	ctx, cancel := context.WithTimeout(c, s.timeout)
	defer cancel()

	print(ctx)

	return true, nil
}
