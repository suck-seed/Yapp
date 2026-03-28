package models

import (
	"time"

	"github.com/google/uuid"
)

type Role struct {
	ID        uuid.UUID `db:"id" json:"id"`
	HallID    uuid.UUID `db:"hall_id" json:"hall_id"`
	Name      string    `db:"name" json:"name"`
	Color     *string   `db:"color" json:"color,omitempty"`
	IconURL   *string   `db:"icon_url" json:"icon_url,omitempty"`
	IsDefault bool      `db:"is_default" json:"is_default"`
	IsAdmin   bool      `db:"is_admin" json:"is_admin"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
	UpdatedAt time.Time `db:"updated_at" json:"updated_at"`
}

// RolePermission represents all permissions for a role
type RolePermission struct {
	RoleID uuid.UUID `db:"role_id" json:"role_id"`

	// General Permissions
	ViewChannels    bool `db:"view_channels" json:"view_channels"`
	ManageChannels  bool `db:"manage_channels" json:"manage_channels"`
	ManageRoles     bool `db:"manage_roles" json:"manage_roles"`
	ManageServers   bool `db:"manage_servers" json:"manage_servers"`
	ChangeNickname  bool `db:"change_nickname" json:"change_nickname"`
	ManageNicknames bool `db:"manage_nicknames" json:"manage_nicknames"`
	KickMembers     bool `db:"kick_members" json:"kick_members"`
	BanMembers      bool `db:"ban_members" json:"ban_members"`

	// Text Permissions
	TextSendMessages   bool `db:"text_send_messages" json:"text_send_messages"`
	TextAttachFiles    bool `db:"text_attach_files" json:"text_attach_files"`
	TextMentionRoles   bool `db:"text_mention_roles" json:"text_mention_roles"`
	TextManageMessages bool `db:"text_manage_messages" json:"text_manage_messages"`
	TextReadHistory    bool `db:"text_read_history" json:"text_read_history"`
	TextSendVoice      bool `db:"text_send_voice" json:"text_send_voice"`

	// Voice Permissions
	VoiceConnect     bool `db:"voice_connect" json:"voice_connect"`
	VoiceSpeak       bool `db:"voice_speak" json:"voice_speak"`
	VoiceVideo       bool `db:"voice_video" json:"voice_video"`
	VoiceMuteMembers bool `db:"voice_mute_members" json:"voice_mute_members"`
}
