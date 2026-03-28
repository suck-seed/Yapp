package repositories

import (
	"context"

	"github.com/google/uuid"
	"github.com/suck-seed/yapp/internal/database"
	"github.com/suck-seed/yapp/internal/models"
)

type IRoleRepository interface {
	// ---------------------------------------- ROLE
	// Role CUD
	CreateRole(ctx context.Context, db database.DBRunner, hallRole *models.Role) (*models.Role, error)
	GetRole(ctx context.Context, db database.DBRunner, roleID uuid.UUID) (*models.Role, error)
	GetAllRole(ctx context.Context, db database.DBRunner, hallID uuid.UUID) ([]*models.Role, error)
	UpdateRole(ctx context.Context, db database.DBRunner, role *models.Role) (*models.Role, error)
	DeleteRole(ctx context.Context, db database.DBRunner, roleID uuid.UUID) (*models.Role, error)

	// --------------------------------------- PERMISSION
	// Permission CUD
	CreateRolePermissions(ctx context.Context, db database.DBRunner, permissions *models.RolePermission) (*models.RolePermission, error)
	GetRolePermissions(ctx context.Context, db database.DBRunner, roleID uuid.UUID) (*models.RolePermission, error)
	UpdateRolePermissions(ctx context.Context, db database.DBRunner, permissions *models.RolePermission) (*models.RolePermission, error)
	DeleteRolePermissions(ctx context.Context, db database.DBRunner, roleID uuid.UUID) (*models.RolePermission, error)

	// User's permissions check
	GetUserPermissionsInHall(ctx context.Context, db database.DBRunner, hallID, userID uuid.UUID) (*models.RolePermission, error)

	// Bulk operations
	GetMultipleRolePermissions(ctx context.Context, db database.DBRunner, roleIDs []uuid.UUID) (map[uuid.UUID]*models.RolePermission, error)

	// Add to IHallRepository interface
	GetHallDefaultRole(ctx context.Context, db database.DBRunner, hallID uuid.UUID) (*models.Role, error)
}

type roleRepository struct {
}

func NewRoleRepository() IRoleRepository {
	return &roleRepository{}
}

func (r *roleRepository) CreateRole(ctx context.Context, db database.DBRunner, hallRole *models.Role) (*models.Role, error) {

	query := `
    INSERT INTO roles (id, hall_id, name, color, icon_url, is_default, is_admin)
    VALUES ($1, $2, $3, $4, $5, $6, $7)
    RETURNING id, hall_id, name, color, icon_url, is_default, is_admin, created_at, updated_at
    `

	row := db.QueryRow(ctx, query,
		hallRole.ID,
		hallRole.HallID,
		hallRole.Name,
		hallRole.Color,
		hallRole.IconURL,
		hallRole.IsDefault,
		hallRole.IsAdmin,
	)

	saved := &models.Role{}
	err := row.Scan(
		&saved.ID,
		&saved.HallID,
		&saved.Name,
		&saved.Color,
		&saved.IconURL,
		&saved.IsDefault,
		&saved.IsAdmin,
		&saved.CreatedAt,
		&saved.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return saved, nil
}
func (r *roleRepository) GetRole(ctx context.Context, db database.DBRunner, roleID uuid.UUID) (*models.Role, error) {

	query := `
    SELECT
    	id, hall_id, name, color, icon_url, is_default, is_admin, created_at, updated_at
    FROM roles
    WHERE id = $1
    `

	saved := &models.Role{}

	err := db.QueryRow(ctx, query, roleID).Scan(
		&saved.ID,
		&saved.HallID,
		&saved.Name,
		&saved.Color,
		&saved.IconURL,
		&saved.IsDefault,
		&saved.IsAdmin,
		&saved.CreatedAt,
		&saved.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return saved, nil

}

func (r *roleRepository) GetAllRole(ctx context.Context, db database.DBRunner, hallID uuid.UUID) ([]*models.Role, error) {

	query := `
    SELECT
    	id, hall_id, name, color, icon_url, is_default, is_admin, created_at, updated_at
    FROM roles
    WHERE hall_id = $1
    `

	rows, err := db.Query(ctx, query, hallID)
	if err != nil {
		return nil, nil
	}
	defer rows.Close()

	roles := []*models.Role{}
	for rows.Next() {
		currentRole := &models.Role{}
		err := rows.Scan(
			&currentRole.ID,
			&currentRole.HallID,
			&currentRole.Name,
			&currentRole.Color,
			&currentRole.IconURL,
			&currentRole.IsDefault,
			&currentRole.IsAdmin,
			&currentRole.CreatedAt,
			&currentRole.UpdatedAt,
		)

		// Scan error
		if err != nil {
			return nil, err
		}

		roles = append(roles, currentRole)
	}

	// Error iterating rows
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return roles, nil

}

func (r *roleRepository) UpdateRole(ctx context.Context, db database.DBRunner, role *models.Role) (*models.Role, error) {

	updatedRole := &models.Role{}

	query := `
        UPDATE users
        SET name = $1, color = $2, icon_url = $3, is_default = $4, is_admin = $5
        WHERE id = $6
        RETURNING id, hall_id, name, color, icon_url, is_default, is_admin, created_at, updated_at
    `
	err := db.QueryRow(ctx, query, role.Name, role.Color, role.IconURL, role.IsDefault, role.IsAdmin).Scan(
		&updatedRole.ID,
		&updatedRole.HallID,
		&updatedRole.Name,
		&updatedRole.Color,
		&updatedRole.IconURL,
		&updatedRole.IsDefault,
		&updatedRole.IsAdmin,
		&updatedRole.CreatedAt,
		&updatedRole.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}
	return updatedRole, nil

}

func (r *roleRepository) DeleteRole(ctx context.Context, db database.DBRunner, roleID uuid.UUID) (*models.Role, error) {

	query := `
			DELETE FROM roles WHERE id = $1
			RETURNING id, hall_id, name, color, icon_url, is_default, is_admin, created_at, updated_at
		`

	deleted := &models.Role{}

	row := db.QueryRow(ctx, query, roleID)

	err := row.Scan(
		&deleted.ID,
		&deleted.HallID,
		&deleted.Name,
		&deleted.Color,
		&deleted.IconURL,
		&deleted.IsDefault,
		&deleted.IsAdmin,
		&deleted.CreatedAt,
		&deleted.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return deleted, nil
}

// --------------------------------------- PERMISSION
// Permission CUD
func (r *roleRepository) CreateRolePermissions(ctx context.Context, db database.DBRunner, permissions *models.RolePermission) (*models.RolePermission, error) {

	query := `
    INSERT INTO role_permissions (
        role_id,
        view_channels, manage_channels, manage_roles, manage_servers,
        change_nickname, manage_nicknames, kick_members, ban_members,
        text_send_messages, text_attach_files, text_mention_roles,
        text_manage_messages, text_read_history, text_send_voice,
        voice_connect, voice_speak, voice_video, voice_mute_members
    )
    VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19)
    RETURNING 	role_id,
    			view_channels, manage_channels, manage_roles, manage_servers,
       			change_nickname, manage_nicknames, kick_members, ban_members,
          		text_send_messages, text_attach_files, text_mention_roles,
            	text_manage_messages, text_read_history, text_send_voice,
             	voice_connect, voice_speak, voice_video, voice_mute_members

    `

	saved := &models.RolePermission{}

	row := db.QueryRow(ctx, query,
		permissions.RoleID,
		permissions.ViewChannels, permissions.ManageChannels,
		permissions.ManageRoles,
		permissions.ManageServers,
		permissions.ChangeNickname, permissions.ManageNicknames, permissions.KickMembers,
		permissions.BanMembers,
		permissions.TextSendMessages, permissions.TextAttachFiles, permissions.TextMentionRoles,
		permissions.TextManageMessages, permissions.TextReadHistory, permissions.TextSendVoice,
		permissions.VoiceConnect,
		permissions.VoiceSpeak,
		permissions.VoiceVideo, permissions.VoiceMuteMembers,
	)

	err := row.Scan(
		&saved.RoleID,
		&saved.ViewChannels, &saved.ManageChannels,
		&saved.ManageRoles,
		&saved.ManageServers,
		&saved.ChangeNickname, &saved.ManageNicknames, &saved.KickMembers,
		&saved.BanMembers,
		&saved.TextSendMessages, &saved.TextAttachFiles, &saved.TextMentionRoles,
		&saved.TextManageMessages, &saved.TextReadHistory, &saved.TextSendVoice,
		&saved.VoiceConnect,
		&saved.VoiceSpeak,
		&saved.VoiceVideo, &saved.VoiceMuteMembers,
	)
	if err != nil {
		return nil, err
	}

	return saved, nil

}
func (r *roleRepository) GetRolePermissions(ctx context.Context, db database.DBRunner, roleID uuid.UUID) (*models.RolePermission, error) {

	query := `
    SELECT
        role_id,
        view_channels, manage_channels, manage_roles, manage_servers,
        change_nickname, manage_nicknames, kick_members, ban_members,
        text_send_messages, text_attach_files, text_mention_roles,
        text_manage_messages, text_read_history, text_send_voice,
        voice_connect, voice_speak, voice_video, voice_mute_members
    FROM role_permissions
    WHERE role_id = $1
    `

	permissions := &models.RolePermission{}

	err := db.QueryRow(ctx, query, roleID).Scan(
		&permissions.RoleID,
		&permissions.ViewChannels, &permissions.ManageChannels, &permissions.ManageRoles, &permissions.ManageServers,
		&permissions.ChangeNickname, &permissions.ManageNicknames, &permissions.KickMembers, &permissions.BanMembers,
		&permissions.TextSendMessages, &permissions.TextAttachFiles, &permissions.TextMentionRoles,
		&permissions.TextManageMessages, &permissions.TextReadHistory, &permissions.TextSendVoice,
		&permissions.VoiceConnect, &permissions.VoiceSpeak, &permissions.VoiceVideo, &permissions.VoiceMuteMembers,
	)

	if err != nil {
		return nil, err
	}

	return permissions, nil
}
func (r *roleRepository) UpdateRolePermissions(ctx context.Context, db database.DBRunner, permissions *models.RolePermission) (*models.RolePermission, error) {

	query := `
    UPDATE role_permissions SET
        view_channels = $2, manage_channels = $3, manage_roles = $4, manage_servers = $5,
        change_nickname = $6, manage_nicknames = $7, kick_members = $8, ban_members = $9,
        text_send_messages = $10, text_attach_files = $11, text_mention_roles = $12,
        text_manage_messages = $13, text_read_history = $14, text_send_voice = $15,
        voice_connect = $16, voice_speak = $17, voice_video = $18, voice_mute_members = $19
    WHERE role_id = $1

    RETURNING 	role_id,
    			view_channels, manage_channels, manage_roles, manage_servers,
       			change_nickname, manage_nicknames, kick_members, ban_members,
          		text_send_messages, text_attach_files, text_mention_roles,
            	text_manage_messages, text_read_history, text_send_voice,
             	voice_connect, voice_speak, voice_video, voice_mute_members
    `

	saved := &models.RolePermission{}

	row := db.QueryRow(ctx, query,
		permissions.RoleID,
		permissions.ViewChannels, permissions.ManageChannels,
		permissions.ManageRoles,
		permissions.ManageServers,
		permissions.ChangeNickname, permissions.ManageNicknames, permissions.KickMembers,
		permissions.BanMembers,
		permissions.TextSendMessages, permissions.TextAttachFiles, permissions.TextMentionRoles,
		permissions.TextManageMessages, permissions.TextReadHistory, permissions.TextSendVoice,
		permissions.VoiceConnect,
		permissions.VoiceSpeak,
		permissions.VoiceVideo, permissions.VoiceMuteMembers,
	)

	err := row.Scan(
		&saved.RoleID,
		&saved.ViewChannels, &saved.ManageChannels,
		&saved.ManageRoles,
		&saved.ManageServers,
		&saved.ChangeNickname, &saved.ManageNicknames, &saved.KickMembers,
		&saved.BanMembers,
		&saved.TextSendMessages, &saved.TextAttachFiles, &saved.TextMentionRoles,
		&saved.TextManageMessages, &saved.TextReadHistory, &saved.TextSendVoice,
		&saved.VoiceConnect,
		&saved.VoiceSpeak,
		&saved.VoiceVideo, &saved.VoiceMuteMembers,
	)
	if err != nil {
		return nil, err
	}

	return saved, nil
}
func (r *roleRepository) DeleteRolePermissions(ctx context.Context, db database.DBRunner, roleID uuid.UUID) (*models.RolePermission, error) {
	query := `
		DELETE FROM role_permissions WHERE role_id = $1
	 	RETURNING 	role_id,
    			view_channels, manage_channels, manage_roles, manage_servers,
       			change_nickname, manage_nicknames, kick_members, ban_members,
          		text_send_messages, text_attach_files, text_mention_roles,
            	text_manage_messages, text_read_history, text_send_voice,
             	voice_connect, voice_speak, voice_video, voice_mute_members
	`

	saved := &models.RolePermission{}

	row := db.QueryRow(ctx, query, roleID)

	err := row.Scan(
		&saved.RoleID,
		&saved.ViewChannels, &saved.ManageChannels,
		&saved.ManageRoles,
		&saved.ManageServers,
		&saved.ChangeNickname, &saved.ManageNicknames, &saved.KickMembers,
		&saved.BanMembers,
		&saved.TextSendMessages, &saved.TextAttachFiles, &saved.TextMentionRoles,
		&saved.TextManageMessages, &saved.TextReadHistory, &saved.TextSendVoice,
		&saved.VoiceConnect,
		&saved.VoiceSpeak,
		&saved.VoiceVideo, &saved.VoiceMuteMembers,
	)
	if err != nil {
		return nil, err
	}

	return saved, nil
}

// User's permissions check
func (r *roleRepository) GetUserPermissionsInHall(ctx context.Context, db database.DBRunner, hallID, userID uuid.UUID) (*models.RolePermission, error) {

	query := `
    SELECT
        bool_or(rp.view_channels) as view_channels,
        bool_or(rp.manage_channels) as manage_channels,
        bool_or(rp.manage_roles) as manage_roles,
        bool_or(rp.manage_servers) as manage_servers,
        bool_or(rp.change_nickname) as change_nickname,
        bool_or(rp.manage_nicknames) as manage_nicknames,
        bool_or(rp.kick_members) as kick_members,
        bool_or(rp.ban_members) as ban_members,
        bool_or(rp.text_send_messages) as text_send_messages,
        bool_or(rp.text_attach_files) as text_attach_files,
        bool_or(rp.text_mention_roles) as text_mention_roles,
        bool_or(rp.text_manage_messages) as text_manage_messages,
        bool_or(rp.text_read_history) as text_read_history,
        bool_or(rp.text_send_voice) as text_send_voice,
        bool_or(rp.voice_connect) as voice_connect,
        bool_or(rp.voice_speak) as voice_speak,
        bool_or(rp.voice_video) as voice_video,
        bool_or(rp.voice_mute_members) as voice_mute_members
    FROM hall_members hm
    JOIN roles r ON hm.role_id = r.id
    JOIN role_permissions rp on r.id = rp.role_id
    WHERE hm.hall_id = $1 AND hm.user_id = $2
    GROUP BY hm.user_id
    `

	permissions := &models.RolePermission{
		RoleID: uuid.Nil,
	}

	err := db.QueryRow(ctx, query, hallID, userID).Scan(&permissions.ViewChannels, &permissions.ManageChannels, &permissions.ManageRoles, &permissions.ManageServers,
		&permissions.ChangeNickname, &permissions.ManageNicknames, &permissions.KickMembers, &permissions.BanMembers,
		&permissions.TextSendMessages, &permissions.TextAttachFiles, &permissions.TextMentionRoles,
		&permissions.TextManageMessages, &permissions.TextReadHistory, &permissions.TextSendVoice,
		&permissions.VoiceConnect, &permissions.VoiceSpeak, &permissions.VoiceVideo, &permissions.VoiceMuteMembers)

	if err != nil {
		return nil, err
	}

	return permissions, nil
}

// Bulk operations
func (r *roleRepository) GetMultipleRolePermissions(ctx context.Context, db database.DBRunner, roleIDs []uuid.UUID) (map[uuid.UUID]*models.RolePermission, error) {

	if len(roleIDs) == 0 {
		return make(map[uuid.UUID]*models.RolePermission), nil
	}

	query := `

	SELECT
        role_id,
        view_channels, manage_channels, manage_roles, manage_servers,
        change_nickname, manage_nicknames, kick_members, ban_members,
        text_send_messages, text_attach_files, text_mention_roles,
        text_manage_messages, text_read_history, text_send_voice,
        voice_connect, voice_speak, voice_video, voice_mute_members
    FROM role_permissions
    WHERE role_id = ANY($1)

    `
	rows, err := db.Query(ctx, query, roleIDs)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	permissionMap := make(map[uuid.UUID]*models.RolePermission)

	for rows.Next() {
		currentPermissions := &models.RolePermission{}

		err := rows.Scan(&currentPermissions.RoleID,
			&currentPermissions.ViewChannels, &currentPermissions.ManageChannels, &currentPermissions.ManageRoles, &currentPermissions.ManageServers,
			&currentPermissions.ChangeNickname, &currentPermissions.ManageNicknames, &currentPermissions.KickMembers, &currentPermissions.BanMembers,
			&currentPermissions.TextSendMessages, &currentPermissions.TextAttachFiles, &currentPermissions.TextMentionRoles,
			&currentPermissions.TextManageMessages, &currentPermissions.TextReadHistory, &currentPermissions.TextSendVoice,
			&currentPermissions.VoiceConnect, &currentPermissions.VoiceSpeak, &currentPermissions.VoiceVideo, &currentPermissions.VoiceMuteMembers)

		if err != nil {
			return nil, err
		}

		permissionMap[currentPermissions.RoleID] = currentPermissions

	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return permissionMap, nil
}

func (r *roleRepository) GetHallDefaultRole(ctx context.Context, db database.DBRunner, hallID uuid.UUID) (*models.Role, error) {
	query := `
		SELECT id, hall_id, name, color, icon_url, is_default, is_admin, created_at, updated_at
		FROM roles
		WHERE hall_id = $1 AND is_default = true
		LIMIT 1
	`

	role := &models.Role{}
	err := db.QueryRow(ctx, query, hallID).Scan(
		&role.ID,
		&role.HallID,
		&role.Name,
		&role.Color,
		&role.IconURL,
		&role.IsDefault,
		&role.IsAdmin,
		&role.CreatedAt,
		&role.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return role, nil
}
