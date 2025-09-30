package dto

import (
	"github.com/google/uuid"
	"time"

	"github.com/suck-seed/yapp/internal/models"
)

// Responses

type SignupUserRes struct {
	ID       uuid.UUID `json:"id"`
	Username string    `json:"username" binding:"required,min=3,max=32"`
}

type SigninUserRes struct {
	AccessToken string
	ID          uuid.UUID `json:"id" db:"id"`
	Username    string    `json:"username" db:"username"`
}

type UserPublic struct {
	ID                 uuid.UUID `json:"id"`
	Username           string    `json:"username"`
	DisplayName        string    `json:"display_name"`
	AvatarURL          *string   `json:"avatar_url,omitempty"`
	AvatarThumbnailURL *string   `json:"avatar_thumbnail_url,omitempty"`
}

type UserMe struct {
	ID                 uuid.UUID `json:"id"`
	Username           string    `json:"username"`
	DisplayName        string    `json:"display_name"`
	Email              string    `json:"email"`
	PhoneNumber        *string   `json:"phone_number,omitempty"`
	AvatarURL          *string   `json:"avatar_url,omitempty"`
	AvatarThumbnailURL *string   `json:"avatar_thumbnail_url,omitempty"`
	Description        *string   `json:"description,omitempty"`
	FriendPolicy       string    `json:"friend_policy"`
	Active             bool      `json:"active"`
	//LastSeen     *time.Time          `json:"last_seen,omitempty"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
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

// Func to convert model.User to public and private
func ToUserPublic(u models.User) UserPublic {
	return UserPublic{
		ID:                 u.ID,
		Username:           u.Username,
		DisplayName:        u.DisplayName,
		AvatarURL:          u.AvatarURL,
		AvatarThumbnailURL: u.AvatarThumbnailURL,
	}
}

func ToUserMe(u models.User) UserMe {

	return UserMe{
		ID:                 u.ID,
		Username:           u.Username,
		DisplayName:        u.DisplayName,
		Email:              u.Email,
		PhoneNumber:        u.PhoneNumber,
		AvatarURL:          u.AvatarURL,
		AvatarThumbnailURL: u.AvatarThumbnailURL,
		Description:        u.Description,
		FriendPolicy:       string(u.FriendPolicy),
		Active:             u.Active,
		// LastSeen:     u.LastSeen,
		CreatedAt: u.CreatedAt.Format(time.RFC3339),
		UpdatedAt: u.UpdatedAt.Format(time.RFC3339),
	}
}
