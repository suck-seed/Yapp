package dto

import "github.com/google/uuid"

// Each role when created, consists of the default permission values and can only be changed after the creation

// UpdateRolePermissionReq - used to update the permission accessible to each role,
type UpdateRolePermissionReq struct {

	// Using pointer to differentiate between not-send through the request and false values
	// General Permissions
	ViewChannels    *bool `json:"view_channels"`
	ManageChannels  *bool `json:"manage_channels"`
	ManageRoles     *bool `json:"manage_roles"`
	ManageServers   *bool `json:"manage_servers"`
	ChangeNickname  *bool `json:"change_nickname"`
	ManageNicknames *bool `json:"manage_nicknames"`
	KickMembers     *bool `json:"kick_members"`
	BanMembers      *bool `json:"ban_members"`

	// Text
	TextSendMessages   *bool `json:"text_send_messages"`
	TextAttachFiles    *bool `json:"text_attach_files"`
	TextMentionRoles   *bool `json:"text_mention_roles"`
	TextManageMessages *bool `json:"text_manage_messages"`
	TextReadHistory    *bool `json:"text_read_history"`
	TextSendVoice      *bool `json:"text_send_voice"`

	// Voice
	VoiceConnect     *bool `json:"voice_connect"`
	VoiceSpeak       *bool `json:"voice_speak"`
	VoiceVideo       *bool `json:"voice_video"`
	VoiceMuteMembers *bool `json:"voice_mute_members"`
}

// GetRolePermissionsRes - GET /halls/:hallID/settings/roles/:roleID/permissions
type GetRolePermissionsRes struct {
	RoleID     uuid.UUID            `json:"role_id"`
	RoleName   string               `json:"role_name"`
	IsAdmin    bool                 `json:"is_admin"`   // If true, all permissions are enabled
	IsDefault  bool                 `json:"is_default"` // If true, this is the @everyone role
	Categories []PermissionCategory `json:"categories"`
}

// PermissionCategory - grouped permissions for UI
type PermissionCategory struct {
	Name        string             `json:"name"`        // e.g., "General Permissions"
	Description string             `json:"description"` // e.g., "Fundamental server permissions"
	Permissions []PermissionDetail `json:"permissions"`
}

// PermissionDetail - individual permission with metadata about the specific permission
type PermissionDetail struct {
	Key         string `json:"key"`         // e.g., "manage_channels"
	Name        string `json:"name"`        // e.g., "Manage Channels"
	Description string `json:"description"` // e.g., "Create, edit, and delete channels"
	IsEnabled   bool   `json:"is_enabled"`  // current value for this role
}

// UpdateRolePermissionsRes - PATCH response
type UpdateRolePermissionsRes struct {
	Success    bool                 `json:"success"`
	Message    string               `json:"message"`
	RoleID     uuid.UUID            `json:"role_id"`
	Categories []PermissionCategory `json:"categories"` // Return updated permissions
}
