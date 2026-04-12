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
	"github.com/suck-seed/yapp/internal/constants"
	"github.com/suck-seed/yapp/internal/database"
	dto "github.com/suck-seed/yapp/internal/dto/hall"
	"github.com/suck-seed/yapp/internal/models"
	"github.com/suck-seed/yapp/internal/repositories"
	"github.com/suck-seed/yapp/internal/utils"
)

type IRoleService interface {

	// ------------- HALL ROLES (CRUD)
	ListHallRoles(ctx context.Context, userInfo *auth.UserInfo, hallID uuid.UUID) ([]*dto.HallRoleRes, error)
	GetHallRole(ctx context.Context, userInfo *auth.UserInfo, hallID, roleID uuid.UUID) (*dto.HallRoleRes, error)
	CreateHallRole(ctx context.Context, userInfo *auth.UserInfo, hallID uuid.UUID, req *dto.CreateHallRoleReq) (*dto.HallRoleRes, error)
	UpdateHallRole(ctx context.Context, userInfo *auth.UserInfo, hallID, roleID uuid.UUID, req *dto.UpdateHallRoleReq) (*dto.HallRoleRes, error)
	DeleteHallRole(ctx context.Context, userInfo *auth.UserInfo, hallID, roleID uuid.UUID) (*dto.HallRoleRes, error)

	// ------------- PERMISSIONS
	GetRolePermissions(ctx context.Context, userInfo *auth.UserInfo, hallID, roleID uuid.UUID) (*dto.GetRolePermissionsRes, error)
	GetUserPermissions(ctx context.Context, userInfo *auth.UserInfo, hallID uuid.UUID) (*models.RolePermission, error)

	UpdateRolePermissions(ctx context.Context, userInfo *auth.UserInfo, hallID, roleID uuid.UUID, req *dto.UpdateRolePermissionReq) (*dto.UpdateRolePermissionsRes, error)
}

type roleService struct {
	repositories.IRoleRepository
	repositories.IUserRepository
	repositories.IHallRepository
	repositories.IBanRepsitory

	// Permission checker service
	IPermissionCheckerService

	pool    *pgxpool.Pool
	timeout time.Duration
	mu      sync.RWMutex
}

func NewRoleService(roleRepo repositories.IRoleRepository, userRepo repositories.IUserRepository, hallRepo repositories.IHallRepository, banRepo repositories.IBanRepsitory, permissionChecker IPermissionCheckerService, pool *pgxpool.Pool) IRoleService {
	return &roleService{
		roleRepo,
		userRepo,
		hallRepo,
		banRepo,
		permissionChecker,
		pool,
		time.Duration(2) * time.Second,
		sync.RWMutex{},
	}
}

// ------------- HALL ROLES (CRUD)

func (s *roleService) ListHallRoles(ctx context.Context, userInfo *auth.UserInfo, hallID uuid.UUID) ([]*dto.HallRoleRes, error) {
	ctx, cancel := context.WithTimeout(ctx, s.timeout)
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

	roles, err := s.IRoleRepository.GetAllRole(ctx, runner, hallID)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) || errors.Is(err, context.Canceled) {
			return nil, utils.ErrorRequestTimeout
		}
		return nil, utils.ErrorFetchingRole
	}

	out := make([]*dto.HallRoleRes, 0, len(roles))
	for _, r := range roles {
		out = append(out, hallRoleToDTO(r))
	}
	return out, nil
}

func (s *roleService) GetHallRole(ctx context.Context, userInfo *auth.UserInfo, hallID, roleID uuid.UUID) (*dto.HallRoleRes, error) {
	ctx, cancel := context.WithTimeout(ctx, s.timeout)
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

	role, err := s.IRoleRepository.GetRole(ctx, runner, roleID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, utils.ErrorRoleNotFound
		}
		if errors.Is(err, context.DeadlineExceeded) || errors.Is(err, context.Canceled) {
			return nil, utils.ErrorRequestTimeout
		}
		return nil, utils.ErrorFetchingRole
	}
	if role.HallID != hallID {
		return nil, utils.ErrorRoleDoesntBelongInThisHall
	}

	return hallRoleToDTO(role), nil
}

func (s *roleService) CreateHallRole(ctx context.Context, userInfo *auth.UserInfo, hallID uuid.UUID, req *dto.CreateHallRoleReq) (*dto.HallRoleRes, error) {
	ctx, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()

	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return nil, utils.ErrorInternal
	}
	runner := database.NewTxWrapper(tx)
	defer runner.Rollback(ctx)

	can, err := s.canCreateHallRole(ctx, runner, userInfo.ID, hallID)
	if err != nil {
		return nil, err
	}
	if !can {
		return nil, utils.ErrorUserCannotCreateHallRoles
	}

	nameIn := req.Name
	sanName, err := utils.SanitizeText(&nameIn)
	if err != nil {
		return nil, err
	}
	if sanName == nil || *sanName == "" {
		return nil, utils.ErrorInvalidInput
	}

	var colorPtr *string
	if req.Color != nil {
		c, err := utils.SanitizeColorFormat(req.Color)
		if err != nil {
			return nil, err
		}
		colorPtr = c
	}

	var iconPtr *string
	if req.IconURL != nil {
		i, err := utils.SanitizeText(req.IconURL)
		if err != nil {
			return nil, err
		}
		iconPtr = i
	}

	roleID, err := uuid.NewV7()
	if err != nil {
		return nil, utils.ErrorInternal
	}

	newRole := &models.Role{
		ID:        roleID,
		HallID:    hallID,
		Name:      *sanName,
		Color:     colorPtr,
		IconURL:   iconPtr,
		IsDefault: false,
		IsAdmin:   false,
	}

	saved, err := s.IRoleRepository.CreateRole(ctx, runner, newRole)
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

	perms := defaultRolePermissions(saved.ID)
	if _, err := s.IRoleRepository.CreateRolePermissions(ctx, runner, perms); err != nil {
		if errors.Is(err, context.DeadlineExceeded) || errors.Is(err, context.Canceled) {
			return nil, utils.ErrorRequestTimeout
		}
		return nil, utils.ErrorCreatingHallRole
	}

	if err := runner.Commit(ctx); err != nil {
		return nil, utils.ErrorInternal
	}

	return hallRoleToDTO(saved), nil
}

func (s *roleService) UpdateHallRole(ctx context.Context, userInfo *auth.UserInfo, hallID, roleID uuid.UUID, req *dto.UpdateHallRoleReq) (*dto.HallRoleRes, error) {
	ctx, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()

	if req.Name == nil && req.Color == nil && req.IconURL == nil {
		return nil, utils.ErrorNoFieldsToUpdate
	}

	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return nil, utils.ErrorInternal
	}
	runner := database.NewTxWrapper(tx)
	defer runner.Rollback(ctx)

	canManage, err := s.CanManageRoles(ctx, runner, userInfo.ID, hallID)
	if err != nil {
		return nil, err
	}
	if !canManage {
		return nil, utils.ErrorUserCannotManageRoles
	}

	role, err := s.IRoleRepository.GetRole(ctx, runner, roleID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, utils.ErrorRoleNotFound
		}
		if errors.Is(err, context.DeadlineExceeded) || errors.Is(err, context.Canceled) {
			return nil, utils.ErrorRequestTimeout
		}
		return nil, utils.ErrorFetchingRole
	}
	if role.HallID != hallID {
		return nil, utils.ErrorRoleDoesntBelongInThisHall
	}
	if role.IsDefault || role.IsAdmin {
		return nil, utils.ErrorCannotModifyProtectedHallRole
	}

	updated := *role
	if req.Name != nil {
		san, err := utils.SanitizeText(req.Name)
		if err != nil {
			return nil, err
		}
		if san == nil || *san == "" {
			return nil, utils.ErrorInvalidInput
		}
		updated.Name = *san
	}
	if req.Color != nil {
		c, err := utils.SanitizeColorFormat(req.Color)
		if err != nil {
			return nil, err
		}
		updated.Color = c
	}
	if req.IconURL != nil {
		i, err := utils.SanitizeText(req.IconURL)
		if err != nil {
			return nil, err
		}
		updated.IconURL = i
	}

	saved, err := s.IRoleRepository.UpdateRole(ctx, runner, &updated)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, utils.ErrorRoleNotFound
		}
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return nil, utils.ErrorRoleNameAlreadyExists
		}
		if errors.Is(err, context.DeadlineExceeded) || errors.Is(err, context.Canceled) {
			return nil, utils.ErrorRequestTimeout
		}
		return nil, utils.ErrorInternal
	}

	if err := runner.Commit(ctx); err != nil {
		return nil, utils.ErrorInternal
	}

	return hallRoleToDTO(saved), nil
}

func (s *roleService) DeleteHallRole(ctx context.Context, userInfo *auth.UserInfo, hallID, roleID uuid.UUID) (*dto.HallRoleRes, error) {
	ctx, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()

	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return nil, utils.ErrorInternal
	}
	runner := database.NewTxWrapper(tx)
	defer runner.Rollback(ctx)

	canManage, err := s.CanManageRoles(ctx, runner, userInfo.ID, hallID)
	if err != nil {
		return nil, err
	}
	if !canManage {
		return nil, utils.ErrorUserCannotManageRoles
	}

	role, err := s.IRoleRepository.GetRole(ctx, runner, roleID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, utils.ErrorRoleNotFound
		}
		if errors.Is(err, context.DeadlineExceeded) || errors.Is(err, context.Canceled) {
			return nil, utils.ErrorRequestTimeout
		}
		return nil, utils.ErrorFetchingRole
	}
	if role.HallID != hallID {
		return nil, utils.ErrorRoleDoesntBelongInThisHall
	}
	if role.IsDefault || role.IsAdmin {
		return nil, utils.ErrorCannotModifyProtectedHallRole
	}

	deleted, err := s.IRoleRepository.DeleteRole(ctx, runner, roleID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, utils.ErrorRoleNotFound
		}
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23503" {
			return nil, utils.ErrorCannotDeleteRoleInUse
		}
		if errors.Is(err, context.DeadlineExceeded) || errors.Is(err, context.Canceled) {
			return nil, utils.ErrorRequestTimeout
		}
		return nil, utils.ErrorInternal
	}

	if err := runner.Commit(ctx); err != nil {
		return nil, utils.ErrorInternal
	}

	return hallRoleToDTO(deleted), nil
}

func (s *roleService) canCreateHallRole(ctx context.Context, runner database.DBRunner, userID, hallID uuid.UUID) (bool, error) {
	ownerID, err := s.IHallRepository.GetHallOwnerID(ctx, runner, hallID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return false, utils.ErrorHallNotFound
		}
		if errors.Is(err, context.DeadlineExceeded) || errors.Is(err, context.Canceled) {
			return false, utils.ErrorRequestTimeout
		}
		return false, utils.ErrorFetchingHall
	}
	if ownerID == userID {
		return true, nil
	}

	member, err := s.IHallRepository.GetHallMemberByUserID(ctx, runner, hallID, userID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return false, utils.ErrorUserDoesntBelongHall
		}
		if errors.Is(err, context.DeadlineExceeded) || errors.Is(err, context.Canceled) {
			return false, utils.ErrorRequestTimeout
		}
		return false, utils.ErrorInternal
	}

	memberRole, err := s.IRoleRepository.GetRole(ctx, runner, member.RoleID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return false, utils.ErrorRoleNotFound
		}
		if errors.Is(err, context.DeadlineExceeded) || errors.Is(err, context.Canceled) {
			return false, utils.ErrorRequestTimeout
		}
		return false, utils.ErrorFetchingRole
	}
	if memberRole.IsAdmin {
		return true, nil
	}

	return s.CanManageRoles(ctx, runner, userID, hallID)
}

func hallRoleToDTO(r *models.Role) *dto.HallRoleRes {
	return &dto.HallRoleRes{
		ID:        r.ID,
		HallID:    r.HallID,
		Name:      r.Name,
		Color:     r.Color,
		IconURL:   r.IconURL,
		IsDefault: r.IsDefault,
		IsAdmin:   r.IsAdmin,
		CreatedAt: r.CreatedAt,
		UpdatedAt: r.UpdatedAt,
	}
}

func defaultRolePermissions(roleID uuid.UUID) *models.RolePermission {
	return &models.RolePermission{
		RoleID:             roleID,
		ViewChannels:       true,
		ManageChannels:     false,
		ManageRoles:        false,
		ManageServers:      false,
		ManageInvites:      false,
		ManageRequests:     false,
		ChangeNickname:     true,
		ManageNicknames:    false,
		KickMembers:        false,
		BanMembers:         false,
		TextSendMessages:   true,
		TextAttachFiles:    true,
		TextMentionRoles:   true,
		TextManageMessages: false,
		TextReadHistory:    true,
		TextSendVoice:      true,
		VoiceConnect:       true,
		VoiceSpeak:         true,
		VoiceVideo:         false,
		VoiceMuteMembers:   false,
	}
}

func adminRolePermissions(roleID uuid.UUID) *models.RolePermission {
	return &models.RolePermission{
		RoleID:             roleID,
		ViewChannels:       true,
		ManageChannels:     true,
		ManageRoles:        true,
		ManageServers:      true,
		ManageInvites:      true,
		ManageRequests:     true,
		ChangeNickname:     true,
		ManageNicknames:    true,
		KickMembers:        true,
		BanMembers:         true,
		TextSendMessages:   true,
		TextAttachFiles:    true,
		TextMentionRoles:   true,
		TextManageMessages: true,
		TextReadHistory:    true,
		TextSendVoice:      true,
		VoiceConnect:       true,
		VoiceSpeak:         true,
		VoiceVideo:         true,
		VoiceMuteMembers:   true,
	}
}

// ------------- PERMISSIONS
func (s *roleService) GetRolePermissions(ctx context.Context, userInfo *auth.UserInfo, hallID, roleID uuid.UUID) (*dto.GetRolePermissionsRes, error) {

	ctx, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()

	// ------------ CONNECTION INIT
	conn, err := s.pool.Acquire(ctx)
	if err != nil {
		return nil, utils.ErrorInternal
	}
	defer conn.Release()
	runner := database.NewConnWrapper(conn)

	// Does user belong to hall
	ok, err := s.IHallRepository.IsUserHallMember(ctx, runner, hallID, userInfo.ID)
	if err != nil {
		return nil, utils.ErrorInternal
	}
	if !ok {
		return nil, utils.ErrorUserDoesntBelongHall
	}

	// Getting the role
	role, err := s.IRoleRepository.GetRole(ctx, runner, roleID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, utils.ErrorRoleNotFound
		}

		if errors.Is(err, context.DeadlineExceeded) || errors.Is(err, context.Canceled) {
			return nil, utils.ErrorRequestTimeout
		}

		return nil, utils.ErrorFetchingRole
	}

	// Verifying role belongs to the hall
	if role.HallID != hallID {
		return nil, utils.ErrorRoleDoesntBelongInThisHall
	}

	// GetPermissions
	permissions, err := s.IRoleRepository.GetRolePermissions(ctx, runner, roleID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, utils.ErrorPermissionsNotFound
		}

		if errors.Is(err, context.DeadlineExceeded) || errors.Is(err, context.Canceled) {
			return nil, utils.ErrorRequestTimeout
		}

		return nil, utils.ErrorFetchingPermission
	}

	// build the permission structure
	permissionResponse := s.buildPermissionResponse(role, permissions)

	return permissionResponse, nil
}

func (s *roleService) GetUserPermissions(ctx context.Context, userInfo *auth.UserInfo, hallID uuid.UUID) (*models.RolePermission, error) {

	ctx, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()

	// ------------ CONNECTION INIT
	conn, err := s.pool.Acquire(ctx)
	if err != nil {
		return nil, utils.ErrorInternal
	}
	defer conn.Release()
	runner := database.NewConnWrapper(conn)

	// Check if user is hall.OwnerID
	// If user is owner, we kinda skip the check and return a permission struct with all dial true
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

	if hall.OwnerID == userInfo.ID {
		// Owner has all permissions
		return s.getAllPermissionsEnabled(), nil
	}

	// For all other scope of users, fetch from repo and return the bare struct
	// No formatting and binding into categories required
	rolePermissions, err := s.IRoleRepository.GetUserPermissionsInHall(ctx, runner, hallID, userInfo.ID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			// If repo sends ErrNoRows
			// Highly condition that the userID doesnt belong in the hall or prolly doesnt exist
			return nil, utils.ErrorUserDoesntBelongHall
		}

		if errors.Is(err, context.DeadlineExceeded) || errors.Is(err, context.Canceled) {
			return nil, utils.ErrorRequestTimeout
		}

		return nil, utils.ErrorFetchingPermission
	}

	return rolePermissions, nil
}

func (s *roleService) UpdateRolePermissions(ctx context.Context, userInfo *auth.UserInfo, hallID, roleID uuid.UUID, req *dto.UpdateRolePermissionReq) (*dto.UpdateRolePermissionsRes, error) {

	ctx, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()

	// ------------ TRANSACTION INIT
	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return nil, utils.ErrorInternal
	}

	runner := database.NewTxWrapper(tx)
	defer runner.Rollback(ctx)

	// Verify if the user have right to manage roles
	canManageRoles, err := s.CanManageRoles(ctx, runner, userInfo.ID, hallID)
	if err != nil {
		return nil, err
	}
	if !canManageRoles {
		return nil, utils.ErrorUserCannotManageRoles
	}

	// fetching the corresponding role
	role, err := s.IRoleRepository.GetRole(ctx, runner, roleID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, utils.ErrorRoleNotFound
		}
		if errors.Is(err, context.DeadlineExceeded) || errors.Is(err, context.Canceled) {
			return nil, utils.ErrorRequestTimeout
		}

		return nil, utils.ErrorFetchingRole
	}

	// Verifying if the role belong to this hall
	if role.HallID != hallID {
		return nil, utils.ErrorRoleDoesntBelongInThisHall
	}

	// check if role if default role ( i.e the admin role or as per the user has changed it to )
	// default role has all access, and cannot have its permission updated
	if role.IsDefault {
		return nil, utils.ErrorCannotUpdateDefaultRolePermission
	}

	if role.IsAdmin {
		return nil, utils.ErrorCannotUpdateAdminRolePermission
	}

	currentPermission, err := s.IRoleRepository.GetRolePermissions(ctx, runner, roleID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, utils.ErrorPermissionsNotFound
		}
		if errors.Is(err, context.DeadlineExceeded) || errors.Is(err, context.Canceled) {
			return nil, utils.ErrorRequestTimeout
		}
		return nil, utils.ErrorFetchingPermission
	}

	// apply permission update
	updatedPermissions := s.applyPermissionUpdates(currentPermission, req)
	permissions, err := s.IRoleRepository.UpdateRolePermissions(ctx, runner, updatedPermissions)
	if err != nil {

		if errors.Is(err, pgx.ErrNoRows) {
			return nil, utils.ErrorPermissionsNotFound
		}

		if errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded) {
			return nil, utils.ErrorRequestTimeout
		}

		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			return nil, utils.ErrorUpdatingPermissions
		}

		return nil, utils.ErrorUpdatingPermissions
	}

	// ---------------------- COMMIT
	if err := runner.Commit(ctx); err != nil {
		return nil, utils.ErrorInternal
	}

	// building response
	permissionRes := s.buildPermissionResponse(role, permissions)

	res := &dto.UpdateRolePermissionsRes{
		Success:    true,
		Message:    "Permission Updated",
		RoleID:     roleID,
		Categories: permissionRes.Categories,
	}

	return res, nil

}

// HELPER FUNCTIONS

// buildPermissionResponse builds categorized permission response
func (s *roleService) buildPermissionResponse(role *models.Role, permissions *models.RolePermission) *dto.GetRolePermissionsRes {
	response := &dto.GetRolePermissionsRes{
		RoleID:     role.ID,
		RoleName:   role.Name,
		IsAdmin:    role.IsAdmin,
		IsDefault:  role.IsDefault,
		Categories: []dto.PermissionCategory{},
	}

	// If admin role, all permissions are enabled
	if role.IsAdmin {
		categories := s.buildCategoriesFromMetadata(permissions, true)
		response.Categories = categories
		return response
	}

	categories := s.buildCategoriesFromMetadata(permissions, false)
	response.Categories = categories

	return response
}

// buildCategoriesFromMetadata organizes permissions into categories
func (s *roleService) buildCategoriesFromMetadata(permissions *models.RolePermission, isAdmin bool) []dto.PermissionCategory {
	permsByCategory := constants.GetPermissionByCategory()

	categories := []dto.PermissionCategory{}

	// Build in order: general, text, voice
	for _, categoryKey := range []string{"general", "text", "voice"} {
		perms := permsByCategory[categoryKey]
		categoryMeta := constants.CategoryMetadata[categoryKey]

		category := dto.PermissionCategory{
			Name:        categoryMeta.Name,
			Description: categoryMeta.Description,
			Permissions: []dto.PermissionDetail{},
		}

		for _, perm := range perms {
			detail := dto.PermissionDetail{
				Key:         perm.Key,
				Name:        perm.Name,
				Description: perm.Description,
				IsEnabled:   s.getPermissionValue(permissions, perm.Key, isAdmin),
			}
			category.Permissions = append(category.Permissions, detail)
		}

		categories = append(categories, category)
	}

	return categories
}

// getPermissionValue gets the value of a permission by key
func (s *roleService) getPermissionValue(permissions *models.RolePermission, key string, isAdmin bool) bool {
	if isAdmin {
		return true
	}

	switch key {
	case "view_channels":
		return permissions.ViewChannels
	case "manage_channels":
		return permissions.ManageChannels
	case "manage_roles":
		return permissions.ManageRoles
	case "manage_servers":
		return permissions.ManageServers
	case "manage_invites":
		return permissions.ManageInvites
	case "manage_requests":
		return permissions.ManageRequests
	case "change_nickname":
		return permissions.ChangeNickname
	case "manage_nicknames":
		return permissions.ManageNicknames
	case "kick_members":
		return permissions.KickMembers
	case "ban_members":
		return permissions.BanMembers
	case "text_send_messages":
		return permissions.TextSendMessages
	case "text_attach_files":
		return permissions.TextAttachFiles
	case "text_mention_roles":
		return permissions.TextMentionRoles
	case "text_manage_messages":
		return permissions.TextManageMessages
	case "text_read_history":
		return permissions.TextReadHistory
	case "text_send_voice":
		return permissions.TextSendVoice
	case "voice_connect":
		return permissions.VoiceConnect
	case "voice_speak":
		return permissions.VoiceSpeak
	case "voice_video":
		return permissions.VoiceVideo
	case "voice_mute_members":
		return permissions.VoiceMuteMembers
	default:
		return false
	}
}

// applyPermissionUpdates applies partial updates to permissions
func (s *roleService) applyPermissionUpdates(current *models.RolePermission, updates *dto.UpdateRolePermissionReq) *models.RolePermission {
	updated := *current // Copy current permissions

	// General
	if updates.ViewChannels != nil {
		updated.ViewChannels = *updates.ViewChannels
	}
	if updates.ManageChannels != nil {
		updated.ManageChannels = *updates.ManageChannels
	}
	if updates.ManageRoles != nil {
		updated.ManageRoles = *updates.ManageRoles
	}
	if updates.ManageServers != nil {
		updated.ManageServers = *updates.ManageServers
	}
	if updates.ManageInvites != nil {
		updated.ManageInvites = *updates.ManageInvites
	}
	if updates.ManageInvites != nil {
		updated.ManageRequests = *updates.ManageRequests
	}
	if updates.ChangeNickname != nil {
		updated.ChangeNickname = *updates.ChangeNickname
	}
	if updates.ManageNicknames != nil {
		updated.ManageNicknames = *updates.ManageNicknames
	}
	if updates.KickMembers != nil {
		updated.KickMembers = *updates.KickMembers
	}
	if updates.BanMembers != nil {
		updated.BanMembers = *updates.BanMembers
	}

	// Text
	if updates.TextSendMessages != nil {
		updated.TextSendMessages = *updates.TextSendMessages
	}
	if updates.TextAttachFiles != nil {
		updated.TextAttachFiles = *updates.TextAttachFiles
	}
	if updates.TextMentionRoles != nil {
		updated.TextMentionRoles = *updates.TextMentionRoles
	}
	if updates.TextManageMessages != nil {
		updated.TextManageMessages = *updates.TextManageMessages
	}
	if updates.TextReadHistory != nil {
		updated.TextReadHistory = *updates.TextReadHistory
	}
	if updates.TextSendVoice != nil {
		updated.TextSendVoice = *updates.TextSendVoice
	}

	// Voice
	if updates.VoiceConnect != nil {
		updated.VoiceConnect = *updates.VoiceConnect
	}
	if updates.VoiceSpeak != nil {
		updated.VoiceSpeak = *updates.VoiceSpeak
	}
	if updates.VoiceVideo != nil {
		updated.VoiceVideo = *updates.VoiceVideo
	}
	if updates.VoiceMuteMembers != nil {
		updated.VoiceMuteMembers = *updates.VoiceMuteMembers
	}

	return &updated
}

// getAllPermissionsEnabled returns a permission object with all permissions enabled
func (s *roleService) getAllPermissionsEnabled() *models.RolePermission {
	return &models.RolePermission{
		RoleID:             uuid.Nil,
		ViewChannels:       true,
		ManageChannels:     true,
		ManageRoles:        true,
		ManageServers:      true,
		ManageInvites:      true,
		ManageRequests:     true,
		ChangeNickname:     true,
		ManageNicknames:    true,
		KickMembers:        true,
		BanMembers:         true,
		TextSendMessages:   true,
		TextAttachFiles:    true,
		TextMentionRoles:   true,
		TextManageMessages: true,
		TextReadHistory:    true,
		TextSendVoice:      true,
		VoiceConnect:       true,
		VoiceSpeak:         true,
		VoiceVideo:         true,
		VoiceMuteMembers:   true,
	}
}
