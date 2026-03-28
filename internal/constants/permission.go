package constants

type PermissionMetadata struct {
	Key          string
	Name         string
	Description  string
	Category     string
	DefaultValue bool
}

// static defination of all the permissions

var AllPermissions = []PermissionMetadata{
	// ============= GENERAL PERMISSIONS =============
	{
		Key:          "view_channels",
		Name:         "View Channels",
		Description:  "Allows members to view channels",
		Category:     "general",
		DefaultValue: true,
	},
	{
		Key:          "manage_channels",
		Name:         "Manage Channels",
		Description:  "Create, edit, and delete channels",
		Category:     "general",
		DefaultValue: false,
	},
	{
		Key:          "manage_roles",
		Name:         "Manage Roles",
		Description:  "Create and edit roles below their highest role",
		Category:     "general",
		DefaultValue: false,
	},
	{
		Key:          "manage_servers",
		Name:         "Manage Server",
		Description:  "Change server name, description, and settings",
		Category:     "general",
		DefaultValue: false,
	},
	{
		Key:          "change_nickname",
		Name:         "Change Nickname",
		Description:  "Change their own nickname in the server",
		Category:     "general",
		DefaultValue: true,
	},
	{
		Key:          "manage_nicknames",
		Name:         "Manage Nicknames",
		Description:  "Change nicknames of other members",
		Category:     "general",
		DefaultValue: false,
	},
	{
		Key:          "kick_members",
		Name:         "Kick Members",
		Description:  "Remove members from the server (they can rejoin)",
		Category:     "general",
		DefaultValue: false,
	},
	{
		Key:          "ban_members",
		Name:         "Ban Members",
		Description:  "Permanently ban members from the server",
		Category:     "general",
		DefaultValue: false,
	},

	// ============= TEXT PERMISSIONS =============
	{
		Key:          "text_send_messages",
		Name:         "Send Messages",
		Description:  "Send messages in text channels",
		Category:     "text",
		DefaultValue: true,
	},
	{
		Key:          "text_attach_files",
		Name:         "Attach Files",
		Description:  "Upload files and media to messages",
		Category:     "text",
		DefaultValue: true,
	},
	{
		Key:          "text_mention_roles",
		Name:         "Mention @everyone, @here, and All Roles",
		Description:  "Use @everyone, @here, and mention all roles",
		Category:     "text",
		DefaultValue: true,
	},
	{
		Key:          "text_manage_messages",
		Name:         "Manage Messages",
		Description:  "Delete and pin messages from other members",
		Category:     "text",
		DefaultValue: false,
	},
	{
		Key:          "text_read_history",
		Name:         "Read Message History",
		Description:  "View previous messages in channels",
		Category:     "text",
		DefaultValue: true,
	},
	{
		Key:          "text_send_voice",
		Name:         "Send Voice Messages",
		Description:  "Send voice messages in text channels",
		Category:     "text",
		DefaultValue: true,
	},

	// ============= VOICE PERMISSIONS =============
	{
		Key:          "voice_connect",
		Name:         "Connect",
		Description:  "Join voice channels",
		Category:     "voice",
		DefaultValue: true,
	},
	{
		Key:          "voice_speak",
		Name:         "Speak",
		Description:  "Speak in voice channels",
		Category:     "voice",
		DefaultValue: true,
	},
	{
		Key:          "voice_video",
		Name:         "Video",
		Description:  "Share video in voice channels",
		Category:     "voice",
		DefaultValue: false,
	},
	{
		Key:          "voice_mute_members",
		Name:         "Mute Members",
		Description:  "Server mute other members in voice channels",
		Category:     "voice",
		DefaultValue: false,
	},
}

type PermissionCategory struct {
	Name        string
	Description string
	Order       int
}

var CategoryMetadata = map[string]PermissionCategory{
	"general": {
		Name:        "General Permissions",
		Description: "Fundamental Server Permissions",
		Order:       1,
	},
	"text": {
		Name:        "Text Channel Permissions",
		Description: "Permissions for text channels and messaging",
		Order:       2,
	},
	"voice": {
		Name:        "Voice Channel Permissions",
		Description: "Permissions for voice channels",
		Order:       3,
	},
}

func GetPermissionByCategory() map[string][]PermissionMetadata {
	categories := make(map[string][]PermissionMetadata)

	for _, permission := range AllPermissions {
		categories[permission.Category] = append(categories[permission.Category], permission)
	}

	return categories
}

func GetDefaultPermissions() map[string]bool {
	defaults := make(map[string]bool)

	for _, permission := range AllPermissions {
		defaults[permission.Key] = permission.DefaultValue
	}

	return defaults
}
