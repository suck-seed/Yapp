package dto

import "github.com/suck-seed/yapp/internal/models"

// REQUESTS

type SignupUserReq struct {
	Username string `json:"username" binding:"required,min=3,max=32"`
	Password string `json:"password" binding:"required,min=8"`
	Email    string `json:"email" binding:"required,email"`
	// PhoneNumber *string `json:"phone_number" binding:"omitempty"`
	DisplayName string `json:"display_name" binding:"required,min=1,max=64"`
}

type SigninUserReq struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`

	// // Optional device info (if you track sessions per device)
	// DeviceID   string `json:"device_id" binding:"omitempty,uuid4"`
	// DeviceName string `json:"device_name" binding:"omitempty,max=64"`
}

type RefreshTokenReq struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

type RevokeSessionReq struct {
	DeviceID string `json:"device_id" binding:"omitempty,uuid4"`
	All      bool   `json:"all"` // if true, revoke all sessions
}

type ForgotPasswordReq struct {
	Email string `json:"email" binding:"required,email"`
}

type ResetPasswordConfirmReq struct {
	Token       string `json:"token" binding:"required"`
	NewPassword string `json:"new_password" binding:"required,min=8"`
}

type UpdateProfileReq struct {
	DisplayName string `json:"display_name" binding:"required,min=1,max=64"`
	// PhoneNumber *string `json:"phone_number" binding:"omitempty"`
	AvatarURL *string `json:"avatar_url" binding:"omitempty,url"`
}

type ChangePasswordReq struct {
	OldPassword string `json:"old_password" binding:"required"`
	NewPassword string `json:"new_password" binding:"required,min=8"`
}

type UpdateUsernameReq struct {
	NewUsername string `json:"new_username" binding:"required,min=3,max=32"`
}

// Email change (typically triggers re-verify)
type UpdateEmailReq struct {
	NewEmail string `json:"new_email" binding:"required,email"`
}

// Friend policy update
type UpdateFriendPolicyReq struct {
	FriendPolicy models.FriendPolicy `json:"friend_policy" binding:"required,oneof=everyone friends no_one"`
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
	Provider Provider `json:"provider" binding:"required,oneof=spotify reddit twitter steam instagram"`
	URL      string   `json:"url" binding:"required,url"`
	Show     bool     `json:"show_on_profile" binding:"omitempty"`
}

type DeleteAppLinkReq struct {
	Provider Provider `json:"provider" binding:"required,oneof=spotify reddit twitter steam instagram"`
}

// Search & pagination
type UserSearchReq struct {
	Query  string `json:"query" binding:"omitempty,max=64"` // username/display/email substring
	Limit  int    `json:"limit" binding:"omitempty,min=1,max=100"`
	Offset int    `json:"offset" binding:"omitempty,min=0"`
}
