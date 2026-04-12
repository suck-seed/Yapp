package services

import (
	"context"
	"errors"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/suck-seed/yapp/internal/auth"
	"github.com/suck-seed/yapp/internal/database"
	dto "github.com/suck-seed/yapp/internal/dto/hall"
	"github.com/suck-seed/yapp/internal/models"
	"github.com/suck-seed/yapp/internal/repositories"
	"github.com/suck-seed/yapp/internal/utils"
)

// consts
const HALLCREATORROLENAME string = "creator"
const HALLDEFAULTROLENAME string = "everyone"

type IHallService interface {

	// -------------- HALLS
	CreateHall(c context.Context, userInfo *auth.UserInfo, req *dto.CreateHallReq) (*dto.CreateHallRes, error)
	JoinHall(c context.Context, userInfo *auth.UserInfo, hallID uuid.UUID) (*dto.JoinHallRes, error)
	GetUserHalls(c context.Context, userInfo *auth.UserInfo) ([]*models.Hall, error)
	GetCurrentHall(c context.Context, userInfo *auth.UserInfo, hallID uuid.UUID) (*dto.GetCurrentHallRes, error)
	DeleteHall(c context.Context, userInfo *auth.UserInfo, hallID uuid.UUID) (*models.Hall, error)
	IsUserHallMember(c context.Context, hallID uuid.UUID, userID uuid.UUID) (bool, error)
	DoesHallExist(c context.Context, hallID uuid.UUID) (bool, error)

	// -------------- HALL PROFILE
	GetHallProfile(c context.Context, userInfo *auth.UserInfo, hallID uuid.UUID) (*dto.GetHallProfileRes, error)
	UpdateHallProfile(c context.Context, userInfo *auth.UserInfo, hallID uuid.UUID, req *dto.HallProfileUpdateReq) (*dto.HallProfileUpdateRes, error)

	// -------------- MEMBERS
	GetHallMembers(c context.Context, userInfo *auth.UserInfo, hallID uuid.UUID) (*dto.GetHallMembersRes, error)
	GetHallMember(c context.Context, userInfo *auth.UserInfo, hallID uuid.UUID, memberID uuid.UUID) (*dto.HallMemberRes, error)
	UpdateHallMemberRole(c context.Context, userInfo *auth.UserInfo, hallID uuid.UUID, memberID uuid.UUID, req *dto.UpdateHallMemberRoleReq) (*dto.UpdateHallMemberRes, error)
	UpdateHallMemberNickname(c context.Context, userInfo *auth.UserInfo, hallID uuid.UUID, memberID uuid.UUID, req *dto.UpdateHallMemberNicknameReq) (*dto.UpdateHallMemberRes, error)
	KickHallMember(c context.Context, userInfo *auth.UserInfo, hallID uuid.UUID, memberID uuid.UUID) (*dto.HallMemberRes, error)

	// -------------- JOIN REQUESTS
	GetCurrentRequests(c context.Context, userInfo *auth.UserInfo, hallID uuid.UUID) (*dto.GetCurrentRequestsRes, error)
	AcceptJoinRequest(c context.Context, userInfo *auth.UserInfo, hallID uuid.UUID, requestID uuid.UUID) (*dto.AcceptJoinRequestRes, error)
	DeclineJoinRequest(c context.Context, userInfo *auth.UserInfo, hallID uuid.UUID, requestID uuid.UUID) (*dto.DeclineJoinRequestRes, error)
}

type hallService struct {
	repositories.IHallRepository
	repositories.IUserRepository
	repositories.IRoleRepository
	repositories.IBanRepsitory

	// Permission checker service
	IPermissionCheckerService

	pool *pgxpool.Pool

	timeout time.Duration
	mu      sync.RWMutex
}

func NewHallService(hallRepo repositories.IHallRepository, userRepo repositories.IUserRepository, roleRepo repositories.IRoleRepository, banRepo repositories.IBanRepsitory, permissionChecker IPermissionCheckerService, pool *pgxpool.Pool) IHallService {
	return &hallService{
		hallRepo,
		userRepo,
		roleRepo,
		banRepo,
		permissionChecker,
		pool,
		time.Duration(2) * time.Second,
		sync.RWMutex{},
	}
}

func (s *hallService) CreateHall(c context.Context, userInfo *auth.UserInfo, req *dto.CreateHallReq) (*dto.CreateHallRes, error) {
	ctx, cancel := context.WithTimeout(c, s.timeout)
	defer cancel()

	// --------------- TRANSACTION INIT
	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return nil, utils.ErrorInternal
	}
	runner := database.NewTxWrapper(tx)
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
		OwnerID:          userInfo.ID,
	}

	// pass to repo
	hall, err := s.IHallRepository.CreateHall(ctx, runner, newHall)
	if err != nil {
		return nil, utils.ErrorCreatingHall
	}

	// --------------- CREATOR ROLE
	roleId, err := uuid.NewV7()
	if err != nil {
		return nil, utils.ErrorInternal
	}

	// package a role struct
	creatorRole := &models.Role{
		ID:      roleId,
		HallID:  hall.ID,
		Name:    HALLCREATORROLENAME,
		IsAdmin: true,
	}
	creatorRolePermissions := adminRolePermissions(creatorRole.ID)

	// pass to repo
	creatorRoleCRES, err := s.IRoleRepository.CreateRole(ctx, runner, creatorRole)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return nil, utils.ErrorRoleNameAlreadyExists
		}
		if errors.Is(err, context.DeadlineExceeded) || errors.Is(err, context.Canceled) {
			return nil, utils.ErrorRequestTimeout
		}
		return nil, utils.ErrorCreatingHallRole
	}

	// Setting up permission for creator role
	if _, err := s.IRoleRepository.CreateRolePermissions(ctx, runner, creatorRolePermissions); err != nil {
		if errors.Is(err, context.DeadlineExceeded) || errors.Is(err, context.Canceled) {
			return nil, utils.ErrorRequestTimeout
		}
		return nil, utils.ErrorCreatingHallRole
	}

	// -------------- DEFAULT ROLE
	defaultRoleID, err := uuid.NewV7()
	if err != nil {
		return nil, utils.ErrorInternal
	}

	defaultRole := &models.Role{
		ID:        defaultRoleID,
		HallID:    hall.ID,
		Name:      HALLDEFAULTROLENAME,
		IsDefault: true,
		IsAdmin:   false,
	}
	defaultRolePermissions := defaultRolePermissions(defaultRole.ID)

	_, err = s.IRoleRepository.CreateRole(ctx, runner, defaultRole)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return nil, utils.ErrorRoleNameAlreadyExists
		}
		if errors.Is(err, context.DeadlineExceeded) || errors.Is(err, context.Canceled) {
			return nil, utils.ErrorRequestTimeout
		}
		return nil, utils.ErrorCreatingHallRole
	}

	// Setting up permission for default role
	if _, err := s.IRoleRepository.CreateRolePermissions(ctx, runner, defaultRolePermissions); err != nil {
		if errors.Is(err, context.DeadlineExceeded) || errors.Is(err, context.Canceled) {
			return nil, utils.ErrorRequestTimeout
		}
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
		UserID: userInfo.ID,
		RoleID: creatorRoleCRES.ID,
	}

	// pass to repo
	_, err = s.IHallRepository.CreateHallMember(ctx, runner, newHallMember)
	if err != nil {
		return nil, utils.ErrorCreatingHallMember
	}

	// ---------------------- COMMIT
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
func (s *hallService) JoinHall(c context.Context, userInfo *auth.UserInfo, hallID uuid.UUID) (*dto.JoinHallRes, error) {
	ctx, cancel := context.WithTimeout(c, s.timeout)
	defer cancel()

	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return nil, utils.ErrorInternal
	}
	runner := database.NewTxWrapper(tx)
	defer runner.Rollback(ctx)

	// hall must exist
	hall, err := s.IHallRepository.GetHallByID(ctx, runner, hallID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, utils.ErrorHallNotFound
		}
		if errors.Is(err, context.DeadlineExceeded) || errors.Is(err, context.Canceled) {
			return nil, utils.ErrorRequestTimeout
		}
		return nil, utils.ErrorFetchingHall
	}

	// user must not already be a member
	isMember, err := s.IHallRepository.IsUserHallMember(ctx, runner, hallID, userInfo.ID)
	if err != nil {
		return nil, utils.ErrorFetchingUser
	}
	if isMember {
		return nil, utils.ErrorAlreadyHallMember
	}

	// banned users cannot join or request
	isBanned, err := s.IBanRepsitory.IsUserBanned(ctx, runner, hallID, userInfo.ID)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) || errors.Is(err, context.Canceled) {
			return nil, utils.ErrorRequestTimeout
		}
		return nil, utils.ErrorFetchingBan
	}
	if isBanned {
		return nil, utils.ErrorUserAlreadyBanned
	}

	// PUBLIC HALL -> join directly
	if !hall.IsPrivate {
		defaultRole, err := s.IRoleRepository.GetHallDefaultRole(ctx, runner, hallID)
		if err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				return nil, utils.ErrorHallDefaultRoleNotFound
			}
			if errors.Is(err, context.DeadlineExceeded) || errors.Is(err, context.Canceled) {
				return nil, utils.ErrorRequestTimeout
			}
			return nil, utils.ErrorFetchingRole
		}

		memberID, err := uuid.NewV7()
		if err != nil {
			return nil, utils.ErrorInternal
		}

		newMember := &models.HallMember{
			ID:     memberID,
			HallID: hallID,
			UserID: userInfo.ID,
			RoleID: defaultRole.ID,
		}

		createdMember, err := s.IHallRepository.CreateHallMember(ctx, runner, newMember)
		if err != nil {
			if errors.Is(err, context.DeadlineExceeded) || errors.Is(err, context.Canceled) {
				return nil, utils.ErrorRequestTimeout
			}
			return nil, utils.ErrorCreatingHallMember
		}

		if err := runner.Commit(ctx); err != nil {
			return nil, utils.ErrorInternal
		}

		return &dto.JoinHallRes{
			Status:    "joined",
			MemberID:  &createdMember.ID,
			RequestID: nil,
			HallID:    createdMember.HallID,
			UserID:    createdMember.UserID,
			RoleID:    &createdMember.RoleID,
			Nickname:  createdMember.Nickname,
			JoinedAt:  &createdMember.JoinedAt,
			CreatedAt: createdMember.CreatedAt,
			UpdatedAt: createdMember.UpdatedAt,
		}, nil
	}

	// PRIVATE HALL -> create join request
	alreadyRequested, err := s.IHallRepository.DoesPendingJoinRequestExist(ctx, runner, hallID, userInfo.ID)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) || errors.Is(err, context.Canceled) {
			return nil, utils.ErrorRequestTimeout
		}
		return nil, utils.ErrorFetchingJoinRequest
	}
	if alreadyRequested {
		return nil, utils.ErrorJoinRequestAlreadyExists
	}

	requestID, err := uuid.NewV7()
	if err != nil {
		return nil, utils.ErrorInternal
	}

	newRequest := &models.HallRequest{
		ID:     requestID,
		HallID: hallID,
		UserID: userInfo.ID,
	}

	createdRequest, err := s.IHallRepository.CreateJoinRequest(ctx, runner, newRequest)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) || errors.Is(err, context.Canceled) {
			return nil, utils.ErrorRequestTimeout
		}
		return nil, utils.ErrorCreatingJoinRequest
	}

	if err := runner.Commit(ctx); err != nil {
		return nil, utils.ErrorInternal
	}

	return &dto.JoinHallRes{
		Status:    "requested",
		MemberID:  nil,
		RequestID: &createdRequest.ID,
		HallID:    createdRequest.HallID,
		UserID:    createdRequest.UserID,
		RoleID:    nil,
		Nickname:  nil,
		JoinedAt:  nil,
		CreatedAt: createdRequest.CreatedAt,
		UpdatedAt: createdRequest.UpdatedAt,
	}, nil
}

func (s *hallService) GetUserHalls(c context.Context, userInfo *auth.UserInfo) ([]*models.Hall, error) {
	ctx, cancel := context.WithTimeout(c, s.timeout)
	defer cancel()

	// --------------- CONNECTION INIT
	conn, err := s.pool.Acquire(ctx)
	if err != nil {
		return nil, utils.ErrorInternal
	}
	defer conn.Release()
	runner := database.NewConnWrapper(conn)

	hallIds, err := s.IHallRepository.GetUserHallIDs(ctx, runner, userInfo.ID)
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

func (s *hallService) GetCurrentHall(c context.Context, userInfo *auth.UserInfo, hallID uuid.UUID) (*dto.GetCurrentHallRes, error) {

	ctx, cancel := context.WithTimeout(c, s.timeout)
	defer cancel()

	// --------------- CONNECTION INIT
	conn, err := s.pool.Acquire(ctx)
	if err != nil {
		return nil, utils.ErrorInternal
	}
	defer conn.Release()
	runner := database.NewConnWrapper(conn)

	// fetch the hallInformation
	hall, err := s.IHallRepository.GetHallByID(ctx, runner, hallID)
	if err != nil {
		// check the pgx error type for further validation of error
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, utils.ErrorHallNotFound
		}

		// timeout or cancelled error
		if errors.Is(err, context.DeadlineExceeded) || errors.Is(err, context.Canceled) {
			return nil, utils.ErrorRequestTimeout
		}

		return nil, utils.ErrorFetchingHall
	}

	return &dto.GetCurrentHallRes{
		ID:               hall.ID,
		Name:             hall.Name,
		IsPrivate:        hall.IsPrivate,
		IconURL:          hall.IconURL,
		IconThumbnailURL: hall.IconThumbnailURL,
		BannerColor:      hall.BannerColor,
		Description:      hall.Description,
		CreatedAt:        hall.CreatedAt,
		UpdatedAt:        hall.UpdatedAt,
		OwnerID:          hall.OwnerID,
	}, nil
}

func (s *hallService) DeleteHall(c context.Context, userInfo *auth.UserInfo, hallID uuid.UUID) (*models.Hall, error) {

	ctx, cancel := context.WithTimeout(c, s.timeout)
	defer cancel()

	// --------------- TRANSACTION INIT
	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return nil, utils.ErrorInternal
	}
	runner := database.NewTxWrapper(tx)
	defer runner.Rollback(ctx)

	// Check if userInfo.ID is equal to hall.ownerID
	// - only owner can delete hall (admin cannot)
	ownerID, err := s.IHallRepository.GetHallOwnerID(ctx, runner, hallID)

	// checking for error and condition where uuid.Nil is sent for no ownerID is found
	if err != nil && ownerID == uuid.Nil {
		if errors.Is(err, context.DeadlineExceeded) || errors.Is(err, context.Canceled) {
			return nil, utils.ErrorRequestTimeout
		}

		if errors.Is(err, pgx.ErrNoRows) {
			return nil, utils.ErrorHallNotFound
		}

		return nil, utils.ErrorFetchingHall
	}

	// checking if ownerID is equivalent to userInfo.ID
	if userInfo.ID != ownerID {
		return nil, utils.ErrorCannotDeleteHall
	}

	// everything is valid, go on to DeleteHall

	return nil, nil
}

func (s *hallService) IsUserHallMember(c context.Context, hallID uuid.UUID, userID uuid.UUID) (bool, error) {
	ctx, cancel := context.WithTimeout(c, s.timeout)
	defer cancel()

	// --------------- CONNECTION INIT
	conn, err := s.pool.Acquire(ctx)
	if err != nil {
		return false, utils.ErrorInternal
	}
	defer conn.Release()
	runner := database.NewConnWrapper(conn)

	isMember, err := s.IHallRepository.IsUserHallMember(ctx, runner, hallID, userID)
	if err != nil {
		return false, utils.ErrorInternal
	}

	return isMember, nil
}

func (s *hallService) DoesHallExist(c context.Context, hallID uuid.UUID) (bool, error) {

	ctx, cancel := context.WithTimeout(c, s.timeout)
	defer cancel()

	// --------------- CONNECTION INIT
	conn, err := s.pool.Acquire(ctx)
	if err != nil {
		return false, utils.ErrorInternal
	}
	defer conn.Release()
	runner := database.NewConnWrapper(conn)

	exists, err := s.IHallRepository.DoesHallExist(ctx, runner, hallID)
	if err != nil {
		return false, utils.ErrorInternal
	}

	return exists, nil

}

func (s *hallService) GetHallProfile(c context.Context, userInfo *auth.UserInfo, hallID uuid.UUID) (*dto.GetHallProfileRes, error) {
	ctx, cancel := context.WithTimeout(c, s.timeout)
	defer cancel()

	conn, err := s.pool.Acquire(ctx)
	if err != nil {
		return nil, utils.ErrorInternal
	}
	defer conn.Release()
	runner := database.NewConnWrapper(conn)

	hall, err := s.IHallRepository.GetHallByID(ctx, runner, hallID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, utils.ErrorHallNotFound
		}
		if errors.Is(err, context.DeadlineExceeded) || errors.Is(err, context.Canceled) {
			return nil, utils.ErrorRequestTimeout
		}
		return nil, utils.ErrorFetchingHall
	}

	return &dto.GetHallProfileRes{
		ID:               hall.ID,
		Name:             hall.Name,
		IsPrivate:        hall.IsPrivate,
		IconURL:          hall.IconURL,
		IconThumbnailURL: hall.IconThumbnailURL,
		BannerColor:      hall.BannerColor,
		Description:      hall.Description,
		OwnerID:          hall.OwnerID,
		CreatedAt:        hall.CreatedAt,
		UpdatedAt:        hall.UpdatedAt,
	}, nil
}

func (s *hallService) UpdateHallProfile(c context.Context, userInfo *auth.UserInfo, hallID uuid.UUID, req *dto.HallProfileUpdateReq) (*dto.HallProfileUpdateRes, error) {
	ctx, cancel := context.WithTimeout(c, s.timeout)
	defer cancel()

	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return nil, utils.ErrorInternal
	}
	runner := database.NewTxWrapper(tx)
	defer runner.Rollback(ctx)

	// Only the owner can edit profile
	ownerID, err := s.IHallRepository.GetHallOwnerID(ctx, runner, hallID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, utils.ErrorHallNotFound
		}
		if errors.Is(err, context.DeadlineExceeded) || errors.Is(err, context.Canceled) {
			return nil, utils.ErrorRequestTimeout
		}
		return nil, utils.ErrorFetchingHall
	}

	if userInfo.ID != ownerID {
		return nil, utils.ErrorUnauthorizedToUpdateHall
	}

	// Build the fields map — only include what was actually sent
	fields := make(map[string]any)

	if req.Name != nil {
		canon, err := utils.SanitizeHallname(*req.Name)
		if err != nil {
			return nil, err
		}
		fields["name"] = canon
	}

	if req.Description != nil {
		canon, err := utils.SanitizeText(req.Description)
		if err != nil {
			return nil, err
		}
		fields["description"] = canon
	}

	if req.BannerColor != nil {
		canon, err := utils.SanitizeColorFormat(req.BannerColor)
		if err != nil {
			return nil, err
		}
		fields["banner_color"] = canon
	}

	if req.IsPrivate != nil {
		fields["is_private"] = *req.IsPrivate
	}

	if len(fields) == 0 {
		return nil, utils.ErrorNoFieldsToUpdate // add this sentinel if not present
	}

	hall, err := s.IHallRepository.UpdateHallProfile(ctx, runner, hallID, fields)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) || errors.Is(err, context.Canceled) {
			return nil, utils.ErrorRequestTimeout
		}
		return nil, utils.ErrorInternal
	}

	if err := runner.Commit(ctx); err != nil {
		return nil, utils.ErrorInternal
	}

	return &dto.HallProfileUpdateRes{
		ID:               hall.ID,
		Name:             hall.Name,
		IsPrivate:        hall.IsPrivate,
		IconURL:          hall.IconURL,
		IconThumbnailURL: hall.IconThumbnailURL,
		BannerColor:      *hall.BannerColor,
		Description:      hall.Description,
		OwnerID:          hall.OwnerID,
		UpdatedAt:        hall.UpdatedAt,
	}, nil
}

func (s *hallService) GetHallMembers(c context.Context, userInfo *auth.UserInfo, hallID uuid.UUID) (*dto.GetHallMembersRes, error) {
	ctx, cancel := context.WithTimeout(c, s.timeout)
	defer cancel()

	conn, err := s.pool.Acquire(ctx)
	if err != nil {
		return nil, utils.ErrorInternal
	}
	defer conn.Release()
	runner := database.NewConnWrapper(conn)

	isMember, err := s.IHallRepository.IsUserHallMember(ctx, runner, hallID, userInfo.ID)
	if err != nil {
		return nil, utils.ErrorInternal
	}
	if !isMember {
		return nil, utils.ErrorUserDoesntBelongHall
	}

	rows, err := s.IHallRepository.ListHallMembers(ctx, runner, hallID)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) || errors.Is(err, context.Canceled) {
			return nil, utils.ErrorRequestTimeout
		}
		return nil, utils.ErrorInternal
	}

	out := make([]*dto.HallMemberRes, 0, len(rows))
	for _, m := range rows {
		out = append(out, &dto.HallMemberRes{
			ID:        m.ID,
			HallID:    m.HallID,
			UserID:    m.UserID,
			RoleID:    m.RoleID,
			Nickname:  m.Nickname,
			JoinedAt:  m.JoinedAt,
			UpdatedAt: m.UpdatedAt,
		})
	}

	return &dto.GetHallMembersRes{
		Members: out,
		Total:   len(out),
	}, nil
}

func (s *hallService) GetHallMember(c context.Context, userInfo *auth.UserInfo, hallID uuid.UUID, memberID uuid.UUID) (*dto.HallMemberRes, error) {
	ctx, cancel := context.WithTimeout(c, s.timeout)
	defer cancel()

	conn, err := s.pool.Acquire(ctx)
	if err != nil {
		return nil, utils.ErrorInternal
	}
	defer conn.Release()
	runner := database.NewConnWrapper(conn)

	isMember, err := s.IHallRepository.IsUserHallMember(ctx, runner, hallID, userInfo.ID)
	if err != nil {
		return nil, utils.ErrorInternal
	}
	if !isMember {
		return nil, utils.ErrorUserDoesntBelongHall
	}

	member, err := s.IHallRepository.GetHallMemberByID(ctx, runner, hallID, memberID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, utils.ErrorMemberNotFound
		}
		if errors.Is(err, context.DeadlineExceeded) || errors.Is(err, context.Canceled) {
			return nil, utils.ErrorRequestTimeout
		}
		return nil, utils.ErrorInternal
	}

	return &dto.HallMemberRes{
		ID:        member.ID,
		HallID:    member.HallID,
		UserID:    member.UserID,
		RoleID:    member.RoleID,
		Nickname:  member.Nickname,
		JoinedAt:  member.JoinedAt,
		UpdatedAt: member.UpdatedAt,
	}, nil
}

func (s *hallService) UpdateHallMemberRole(c context.Context, userInfo *auth.UserInfo, hallID uuid.UUID, memberID uuid.UUID, req *dto.UpdateHallMemberRoleReq) (*dto.UpdateHallMemberRes, error) {
	ctx, cancel := context.WithTimeout(c, s.timeout)
	defer cancel()

	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return nil, utils.ErrorInternal
	}
	runner := database.NewTxWrapper(tx)
	defer runner.Rollback(ctx)

	target, err := s.IHallRepository.GetHallMemberByID(ctx, runner, hallID, memberID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, utils.ErrorMemberNotFound
		}
		return nil, utils.ErrorInternal
	}

	canManage, err := s.CanManageRoles(ctx, runner, userInfo.ID, hallID)
	if err != nil {
		return nil, err
	}
	if !canManage {
		return nil, utils.ErrorUserCannotManageRoles
	}

	role, err := s.IRoleRepository.GetRole(ctx, runner, req.RoleID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, utils.ErrorRoleNotFound
		}
		return nil, utils.ErrorFetchingRole
	}
	if role.HallID != hallID {
		return nil, utils.ErrorRoleDoesntBelongInThisHall
	}

	updated, err := s.IHallRepository.UpdateHallMember(ctx, runner, hallID, target.UserID, map[string]any{
		"role_id": req.RoleID,
	})
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) || errors.Is(err, context.Canceled) {
			return nil, utils.ErrorRequestTimeout
		}
		return nil, utils.ErrorInternal
	}

	if err := runner.Commit(ctx); err != nil {
		return nil, utils.ErrorInternal
	}

	return &dto.UpdateHallMemberRes{
		ID:        updated.ID,
		HallID:    updated.HallID,
		UserID:    updated.UserID,
		RoleID:    updated.RoleID,
		Nickname:  updated.Nickname,
		JoinedAt:  updated.JoinedAt,
		UpdatedAt: updated.UpdatedAt,
	}, nil
}

func (s *hallService) UpdateHallMemberNickname(c context.Context, userInfo *auth.UserInfo, hallID uuid.UUID, memberID uuid.UUID, req *dto.UpdateHallMemberNicknameReq) (*dto.UpdateHallMemberRes, error) {
	ctx, cancel := context.WithTimeout(c, s.timeout)
	defer cancel()

	if req.Nickname == nil {
		return nil, utils.ErrorNoFieldsToUpdate
	}

	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return nil, utils.ErrorInternal
	}
	runner := database.NewTxWrapper(tx)
	defer runner.Rollback(ctx)

	target, err := s.IHallRepository.GetHallMemberByID(ctx, runner, hallID, memberID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, utils.ErrorMemberNotFound
		}
		return nil, utils.ErrorInternal
	}

	var allowed bool
	if userInfo.ID == target.UserID {
		allowed, err = s.CanChangeNickname(ctx, runner, userInfo.ID, hallID)
	} else {
		allowed, err = s.CanManageNicknames(ctx, runner, userInfo.ID, hallID)
	}
	if err != nil {
		return nil, err
	}
	if !allowed {
		if userInfo.ID == target.UserID {
			return nil, utils.ErrorUserCannotChangeNickname
		}
		return nil, utils.ErrorUserCannotManageNicknames
	}

	var nicknameVal any
	if *req.Nickname == "" {
		nicknameVal = nil
	} else {
		san, err := utils.SanitizeText(req.Nickname)
		if err != nil {
			return nil, err
		}
		if san == nil || *san == "" {
			nicknameVal = nil
		} else {
			nicknameVal = *san
		}
	}

	updated, err := s.IHallRepository.UpdateHallMember(ctx, runner, hallID, target.UserID, map[string]any{
		"nickname": nicknameVal,
	})
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) || errors.Is(err, context.Canceled) {
			return nil, utils.ErrorRequestTimeout
		}
		return nil, utils.ErrorInternal
	}

	if err := runner.Commit(ctx); err != nil {
		return nil, utils.ErrorInternal
	}

	return &dto.UpdateHallMemberRes{
		ID:        updated.ID,
		HallID:    updated.HallID,
		UserID:    updated.UserID,
		RoleID:    updated.RoleID,
		Nickname:  updated.Nickname,
		JoinedAt:  updated.JoinedAt,
		UpdatedAt: updated.UpdatedAt,
	}, nil
}

func (s *hallService) KickHallMember(c context.Context, userInfo *auth.UserInfo, hallID uuid.UUID, memberID uuid.UUID) (*dto.HallMemberRes, error) {
	ctx, cancel := context.WithTimeout(c, s.timeout)
	defer cancel()

	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return nil, utils.ErrorInternal
	}
	runner := database.NewTxWrapper(tx)
	defer runner.Rollback(ctx)

	target, err := s.IHallRepository.GetHallMemberByID(ctx, runner, hallID, memberID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, utils.ErrorMemberNotFound
		}
		return nil, utils.ErrorInternal
	}

	ownerID, err := s.IHallRepository.GetHallOwnerID(ctx, runner, hallID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, utils.ErrorHallNotFound
		}
		if errors.Is(err, context.DeadlineExceeded) || errors.Is(err, context.Canceled) {
			return nil, utils.ErrorRequestTimeout
		}
		return nil, utils.ErrorFetchingHall
	}

	if target.UserID == ownerID {
		return nil, utils.ErrorCannotKickHallOwner
	}
	if target.UserID == userInfo.ID {
		return nil, utils.ErrorCannotKickYourself
	}

	canKick, err := s.CanKickMembers(ctx, runner, userInfo.ID, hallID)
	if err != nil {
		return nil, err
	}
	if !canKick {
		return nil, utils.ErrorUserCannotKickMembers
	}

	res := &dto.HallMemberRes{
		ID:        target.ID,
		HallID:    target.HallID,
		UserID:    target.UserID,
		RoleID:    target.RoleID,
		Nickname:  target.Nickname,
		JoinedAt:  target.JoinedAt,
		UpdatedAt: target.UpdatedAt,
	}

	if err := s.IHallRepository.KickHallMember(ctx, runner, hallID, target.UserID); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, utils.ErrorMemberNotFound
		}
		if errors.Is(err, context.DeadlineExceeded) || errors.Is(err, context.Canceled) {
			return nil, utils.ErrorRequestTimeout
		}
		return nil, utils.ErrorInternal
	}

	if err := runner.Commit(ctx); err != nil {
		return nil, utils.ErrorInternal
	}

	return res, nil
}

// ------------------------ JOIN
func (s *hallService) GetCurrentRequests(c context.Context, userInfo *auth.UserInfo, hallID uuid.UUID) (*dto.GetCurrentRequestsRes, error) {
	ctx, cancel := context.WithTimeout(c, s.timeout)
	defer cancel()

	conn, err := s.pool.Acquire(ctx)
	if err != nil {
		return nil, utils.ErrorInternal
	}
	defer conn.Release()
	runner := database.NewConnWrapper(conn)

	canManage, err := s.CanManageRequests(ctx, runner, userInfo.ID, hallID)
	if err != nil {
		return nil, err
	}
	if !canManage {
		return nil, utils.ErrorUserCannotManageRequests
	}

	requests, err := s.IHallRepository.GetAllHallRequests(ctx, runner, hallID)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) || errors.Is(err, context.Canceled) {
			return nil, utils.ErrorRequestTimeout
		}
		return nil, utils.ErrorFetchingJoinRequest
	}

	out := make([]*dto.HallRequestRes, 0, len(requests))
	for _, req := range requests {
		out = append(out, &dto.HallRequestRes{
			ID:        req.ID,
			HallID:    req.HallID,
			UserID:    req.UserID,
			CreatedAt: req.CreatedAt,
			UpdatedAt: req.UpdatedAt,
		})
	}

	return &dto.GetCurrentRequestsRes{
		Requests: out,
		Total:    len(out),
	}, nil
}

func (s *hallService) AcceptJoinRequest(c context.Context, userInfo *auth.UserInfo, hallID uuid.UUID, requestID uuid.UUID) (*dto.AcceptJoinRequestRes, error) {
	ctx, cancel := context.WithTimeout(c, s.timeout)
	defer cancel()

	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return nil, utils.ErrorInternal
	}
	runner := database.NewTxWrapper(tx)
	defer runner.Rollback(ctx)

	canManage, err := s.CanManageRequests(ctx, runner, userInfo.ID, hallID)
	if err != nil {
		return nil, err
	}
	if !canManage {
		return nil, utils.ErrorUserCannotManageRequests
	}

	request, err := s.IHallRepository.GetJoinRequestByID(ctx, runner, requestID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, utils.ErrorJoinRequestNotFound
		}
		if errors.Is(err, context.DeadlineExceeded) || errors.Is(err, context.Canceled) {
			return nil, utils.ErrorRequestTimeout
		}
		return nil, utils.ErrorFetchingJoinRequest
	}

	if request.HallID != hallID {
		return nil, utils.ErrorJoinRequestDoesntBelongToHall
	}

	isMember, err := s.IHallRepository.IsUserHallMember(ctx, runner, hallID, request.UserID)
	if err != nil {
		return nil, utils.ErrorFetchingUser
	}
	if isMember {
		return nil, utils.ErrorAlreadyHallMember
	}

	isBanned, err := s.IBanRepsitory.IsUserBanned(ctx, runner, hallID, request.UserID)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) || errors.Is(err, context.Canceled) {
			return nil, utils.ErrorRequestTimeout
		}
		return nil, utils.ErrorFetchingBan
	}
	if isBanned {
		return nil, utils.ErrorUserAlreadyBanned
	}

	defaultRole, err := s.IRoleRepository.GetHallDefaultRole(ctx, runner, hallID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, utils.ErrorHallDefaultRoleNotFound
		}
		if errors.Is(err, context.DeadlineExceeded) || errors.Is(err, context.Canceled) {
			return nil, utils.ErrorRequestTimeout
		}
		return nil, utils.ErrorFetchingRole
	}

	memberID, err := uuid.NewV7()
	if err != nil {
		return nil, utils.ErrorInternal
	}

	newMember := &models.HallMember{
		ID:     memberID,
		HallID: hallID,
		UserID: request.UserID,
		RoleID: defaultRole.ID,
	}

	member, err := s.IHallRepository.CreateHallMember(ctx, runner, newMember)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) || errors.Is(err, context.Canceled) {
			return nil, utils.ErrorRequestTimeout
		}
		return nil, utils.ErrorCreatingHallMember
	}

	_, err = s.IHallRepository.DeleteJoinRequest(ctx, runner, requestID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, utils.ErrorJoinRequestNotFound
		}
		if errors.Is(err, context.DeadlineExceeded) || errors.Is(err, context.Canceled) {
			return nil, utils.ErrorRequestTimeout
		}
		return nil, utils.ErrorDeletingJoinRequest
	}

	if err := runner.Commit(ctx); err != nil {
		return nil, utils.ErrorInternal
	}

	return &dto.AcceptJoinRequestRes{
		RequestID: request.ID,
		MemberID:  member.ID,
		HallID:    member.HallID,
		UserID:    member.UserID,
		RoleID:    member.RoleID,
		JoinedAt:  member.JoinedAt,
	}, nil
}

func (s *hallService) DeclineJoinRequest(c context.Context, userInfo *auth.UserInfo, hallID uuid.UUID, requestID uuid.UUID) (*dto.DeclineJoinRequestRes, error) {
	ctx, cancel := context.WithTimeout(c, s.timeout)
	defer cancel()

	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return nil, utils.ErrorInternal
	}
	runner := database.NewTxWrapper(tx)
	defer runner.Rollback(ctx)

	canManage, err := s.CanManageRequests(ctx, runner, userInfo.ID, hallID)
	if err != nil {
		return nil, err
	}
	if !canManage {
		return nil, utils.ErrorUserCannotManageRequests
	}

	request, err := s.IHallRepository.GetJoinRequestByID(ctx, runner, requestID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, utils.ErrorJoinRequestNotFound
		}
		if errors.Is(err, context.DeadlineExceeded) || errors.Is(err, context.Canceled) {
			return nil, utils.ErrorRequestTimeout
		}
		return nil, utils.ErrorFetchingJoinRequest
	}

	if request.HallID != hallID {
		return nil, utils.ErrorJoinRequestDoesntBelongToHall
	}

	deleted, err := s.IHallRepository.DeleteJoinRequest(ctx, runner, requestID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, utils.ErrorJoinRequestNotFound
		}
		if errors.Is(err, context.DeadlineExceeded) || errors.Is(err, context.Canceled) {
			return nil, utils.ErrorRequestTimeout
		}
		return nil, utils.ErrorDeletingJoinRequest
	}

	if err := runner.Commit(ctx); err != nil {
		return nil, utils.ErrorInternal
	}

	return &dto.DeclineJoinRequestRes{
		ID:        deleted.ID,
		HallID:    deleted.HallID,
		UserID:    deleted.UserID,
		CreatedAt: deleted.CreatedAt,
		UpdatedAt: deleted.UpdatedAt,
	}, nil
}
