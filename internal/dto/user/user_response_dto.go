package dto

import (
	"time"

	"github.com/google/uuid"
	"github.com/suck-seed/yapp/internal/models"
)

// Responses

type SignupUserRes struct {
	ID       uuid.UUID `json:"id"`
	Username string    `json:"username"`
	Email    string    `json:"email"`
}

type SigninUserRes struct {
	AccessToken string `json:"-"`
	UserMe
	Success bool `json:"success"`
}

type UserPublic struct {
	ID                 uuid.UUID `json:"id"`
	Username           string    `json:"username"`
	DisplayName        string    `json:"display_name"`
	AvatarURL          *string   `json:"avatar_url"`
	AvatarThumbnailURL *string   `json:"avatar_thumbnail_url"`
	Description        *string   `json:"description"`
	AppLinks           []AppLink `json:"app_links"`
	FriendCount        int       `json:"friend_count"`
	MutualFriendCount  int       `json:"mutual_friend_count,omitempty"`
	IsFriend           bool      `json:"is_friend"`
}

type UserMe struct {
	ID                 uuid.UUID `json:"id"`
	Username           string    `json:"username"`
	DisplayName        string    `json:"display_name"`
	Email              string    `json:"email"`
	PhoneNumber        *string   `json:"phone_number"`
	AvatarURL          *string   `json:"avatar_url"`
	AvatarThumbnailURL *string   `json:"avatar_thumbnail_url"`
	Description        *string   `json:"description"`
	FriendPolicy       string    `json:"friend_policy"`
	AppLinks           []AppLink `json:"app_links"`
	CreatedAt          string    `json:"created_at"`
	UpdatedAt          string    `json:"updated_at"`
}

type AppLink struct {
	Provider models.AppProvider `json:"provider"`
	URL      string             `json:"url"`
	Show     bool               `json:"show_on_profile"`
}

type FriendRequestRes struct {
	ID        uuid.UUID  `json:"id"`
	Sender    UserPublic `json:"sender"`
	Receiver  UserPublic `json:"receiver"`
	CreatedAt time.Time  `json:"created_at"`
}

type MutualFriendRes struct {
	Users []*UserPublic `json:"users"`
	Total int           `json:"total"`
}

type FriendListRes struct {
	Users []*UserPublic `json:"users"`
	Total int           `json:"total"`
}

type UpsertAppLinkRes struct {
	ID            uuid.UUID          `json:"id"`
	UserID        uuid.UUID          `json:"user_id"`
	Provider      models.AppProvider `json:"provider"`
	URL           string             `json:"url"`
	ShowOnProfile bool               `json:"show_on_profile"`
	CreatedAt     time.Time          `json:"created_at"`
	UpdatedAt     time.Time          `json:"updated_at"`
}

func ToUserPublic(u models.User) UserPublic {
	return UserPublic{
		ID:                 u.ID,
		Username:           u.Username,
		DisplayName:        u.DisplayName,
		AvatarURL:          u.AvatarURL,
		AvatarThumbnailURL: u.AvatarThumbnailURL,
		Description:        u.Description,
		AppLinks:           []AppLink{},
		FriendCount:        0,
		IsFriend:           false,
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
		AppLinks:           []AppLink{},
		CreatedAt:          u.CreatedAt.Format(time.RFC3339),
		UpdatedAt:          u.UpdatedAt.Format(time.RFC3339),
	}
}
