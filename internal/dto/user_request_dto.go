package dto

import "github.com/suck-seed/yapp/internal/models"

// REQUESTS

type CreateUserReq struct {
	Username    string `json:"username" validate:"required,min=3,max=32"`
	Password    string `json:"password" validate:"required,min=8"`
	Email       string `json:"email" validate:"required,email"`
	PhoneNumber string `json:"phone_number" validate:"omitempty"`
	DisplayName string `json:"display_name" validate:"omitempty,max=64"`
	// if empty, server may default display_name = username
}

type UserLoginRew struct {
	UsernameOrEmail string `json:"username_or_email" validate:"required"`
	Password        string `json:"password" validate:"required"`
	// Optional device info (if you track sessions per device)
	DeviceID   string `json:"device_id" validate:"omitempty,uuid4"`
	DeviceName string `json:"device_name" validate:"omitempty,max=64"`
}

type RefreshTokenReq struct {
	RefreshToken string `json:"refresh_token" validate:"required"`
}

type RevokeSessionReq struct {
	DeviceID string `json:"device_id" validate:"omitempty,uuid4"`
	All      bool   `json:"all"` // if true, revoke all sessions
}

type ForgotPasswordReq struct {
	Email string `json:"email" validate:"required,email"`
}

type ResetPasswordConfirmReq struct {
	Token       string `json:"token" validate:"required"`
	NewPassword string `json:"new_password" validate:"required,min=8"`
}

type UpdateProfileReq struct {
	DisplayName *string `json:"display_name" validate:"omitempty,max=64"`
	PhoneNumber *string `json:"phone_number" validate:"omitempty"`
	AvatarURL   *string `json:"avatar_url" validate:"omitempty,url"`
}

type ChangePasswordReq struct {
	OldPassword string `json:"old_password" validate:"required"`
	NewPassword string `json:"new_password" validate:"required,min=8"`
}

type UpdateUsernameReq struct {
	NewUsername string `json:"new_username" validate:"required,min=3,max=32"`
}

// Email change (typically triggers re-verify)
type UpdateEmailReq struct {
	NewEmail string `json:"new_email" validate:"required,email"`
}

// Friend policy update
type UpdateFriendPolicyReq struct {
	FriendPolicy models.FriendPolicy `json:"friend_policy" validate:"required,oneof=everyone friends no_one"`
}

type SetPresenceReq struct {
	Active bool `json:"active"`
}

// App links (profile)
type Provider string

const (
	ProviderSpotify   Provider = "spotify"
	ProviderReddit    Provider = "reddit"
	ProviderTwitter   Provider = "twitter"
	ProviderSteam     Provider = "steam"
	ProviderInstagram Provider = "instagram"
)

type UpsertAppLinkReq struct {
	Provider Provider `json:"provider" validate:"required,oneof=spotify reddit twitter steam instagram"`
	URL      string   `json:"url" validate:"required,url"`
	Show     bool     `json:"show_on_profile" validate:"omitempty"`
}

type DeleteAppLinkReq struct {
	Provider Provider `json:"provider" validate:"required,oneof=spotify reddit twitter steam instagram"`
}

// Search & pagination
type UserSearchReq struct {
	Query  string `json:"query" validate:"omitempty,max=64"` // username/display/email substring
	Limit  int    `json:"limit" validate:"omitempty,min=1,max=100"`
	Offset int    `json:"offset" validate:"omitempty,min=0"`
}
