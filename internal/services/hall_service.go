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

	GetUserHalls(c context.Context) ([]*models.Hall, error)
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

func (s *hallService) CreateHall(c context.Context, req *dto.CreateHallReq) (*dto.CreateHallRes, error) {
	ctx, cancel := context.WithTimeout(c, s.timeout)
	defer cancel()

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
	userIdString, _, err := auth.CurrentUserFromContext(c)
	if err != nil {
		return nil, err
	}

	userId, err := uuid.Parse(userIdString)
	if err != nil {
		return nil, utils.ErrorInvalidUserIdInContext
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
		Owner:            userId,
	}

	// pass to repo
	hall, err := s.IHallRepository.CreateHall(ctx, newHall)
	if err != nil {
		return nil, utils.ErrorCreatingHall
	}

	//
	// Role Creation
	//

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
	role, err := s.IHallRepository.CreateHallRole(ctx, newRole)
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
		UserID: userId,
		RoleID: role.ID,
	}

	// pass to repo
	err = s.IHallRepository.CreateHallMember(ctx, newHallMember)
	if err != nil {
		return nil, utils.ErrorCreatingHallMember
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
		CreatedBy:        hall.Owner,
	}, nil
}

func (s *hallService) GetUserHalls(c context.Context) ([]*models.Hall, error) {
	ctx, cancel := context.WithTimeout(c, s.timeout)
	defer cancel()

	// get userId from context.Context()
	userIdString, _, err := auth.CurrentUserFromContext(c)
	if err != nil {
		return nil, err
	}

	userId, err := uuid.Parse(userIdString)
	if err != nil {
		return nil, utils.ErrorInvalidUserIdInContext
	}

	hallIds, err := s.IHallRepository.GetUserHallIDs(ctx, userId)
	if err != nil {
		return nil, utils.ErrorInternal
	}

	var halls []*models.Hall
	for _, hallId := range hallIds {
		hall, err := s.IHallRepository.GetHallByID(ctx, hallId)
		if err != nil {
			return nil, utils.ErrorInternal
		}
		halls = append(halls, hall)
	}

	return halls, nil
}

func (s *hallService) IsMember(c context.Context, hallID *uuid.UUID, userId *uuid.UUID) (bool, error) {
	ctx, cancel := context.WithTimeout(c, s.timeout)
	defer cancel()

	print(ctx)

	return true, nil
}
