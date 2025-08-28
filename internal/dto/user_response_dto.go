package dto

import (
	"time"

	"github.com/gofrs/uuid"
	"github.com/suck-seed/yapp/internal/models"
)

//
// ────────────────────────────────
// Responses
// ────────────────────────────────
//

type UserPublic struct {
	ID          string `json:"id"`
	Username    string `json:"username"`
	DisplayName string `json:"display_name"`
	AvatarURL   string `json:"avatar_url,omitempty"`
}

type UserMe struct {
	ID           string              `json:"id"`
	Username     string              `json:"username"`
	DisplayName  *string             `json:"display_name,omitempty"`
	Email        string              `json:"email"`
	PhoneNumber  *string             `json:"phone_number,omitempty"`
	AvatarURL    *string             `json:"avatar_url,omitempty"`
	FriendPolicy models.FriendPolicy `json:"friend_policy"`
	Active       bool                `json:"active"`
	LastSeen     *time.Time          `json:"last_seen,omitempty"`
	CreatedAt    time.Time           `json:"created_at"`
	UpdatedAt    time.Time           `json:"updated_at"`
}

// For app links on response
type AppLink struct {
	Provider Provider `json:"provider"`
	URL      string   `json:"url"`
	Show     bool     `json:"show_on_profile"`
}

// Auth token
type AuthToken struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token,omitempty"`
	ExpiresIn    int64  `json:"expires_in"`
	TokenType    string `json:"token_type,omitempty"` // e.g., "Bearer"
}

type AuthResponse struct {
	User  UserMe    `json:"user"`
	Token AuthToken `json:"token"`
}

type UsernameAvailabilityResponse struct {
	Username  string `json:"username"`
	Available bool   `json:"available"`
}

// Generic ok / error
type OKResponse struct {
	Message string `json:"message"`
}

type ErrorResponse struct {
	Code    string `json:"code"`    // e.g., "conflict_username"
	Message string `json:"message"` // human-readable
}

// Func to convert model.User to public and private
func ToUserPublic(u models.User) UserPublic {
	id := u.ID.String()
	return UserPublic{
		ID:          id,
		Username:    u.Username,
		DisplayName: *u.DisplayName,
		AvatarURL:   *u.AvatarURL,
	}
}

func ToUserMe(u models.User) UserMe {
	return UserMe{
		ID:           uuid.UUID(u.ID).String(),
		Username:     u.Username,
		DisplayName:  u.DisplayName,
		Email:        u.Email,
		PhoneNumber:  u.PhoneNumber,
		AvatarURL:    u.AvatarURL,
		FriendPolicy: u.FriendPolicy,
		Active:       u.Active,
		LastSeen:     u.LastSeen,
		CreatedAt:    u.CreatedAt,
		UpdatedAt:    u.UpdatedAt,
	}
}
