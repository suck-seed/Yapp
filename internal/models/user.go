package models

import (
	"time"

	"github.com/google/uuid"
)

// FriendPolicy type to match Postgres ENUM
type FriendPolicy string

const (
	FriendPolicyEveryone FriendPolicy = "everyone"
	FriendPolicyFriends  FriendPolicy = "friends"
	FriendPolicyNoOne    FriendPolicy = "no_one"
)

type AppProvider string

const (
	AppProviderSpotify AppProvider = "spotify"
	AppProviderReddit  AppProvider = "reddit"
	AppProviderTwitter AppProvider = "twitter"
	AppProviderSteam   AppProvider = "steam"
)

type User struct {
	ID                 uuid.UUID    `json:"id" db:"id"`
	Username           string       `json:"username" db:"username"`
	DisplayName        string       `json:"display_name" db:"display_name"`
	Email              string       `json:"email" db:"email"`
	PasswordHash       string       `json:"password_hash" db:"password_hash"`
	Description        *string      `json:"description,omitempty" db:"description"`
	PhoneNumber        *string      `json:"phone_number,omitempty" db:"phone_number"`
	AvatarURL          *string      `json:"avatar_url,omitempty" db:"avatar_url"`
	AvatarThumbnailURL *string      `json:"avatar_thumbnail_url,omitempty" db:"avatar_thumbnail_url"`
	FriendPolicy       FriendPolicy `json:"friend_policy" db:"friend_policy"`
	CreatedAt          time.Time    `json:"created_at" db:"created_at"`
	UpdatedAt          time.Time    `json:"updated_at" db:"updated_at"`
}

type Friend struct {
	UserID1   uuid.UUID `json:"user_id_1" db:"user_id_1"`
	UserID2   uuid.UUID `json:"user_id_2" db:"user_id_2"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}

type FriendRequest struct {
	ID         uuid.UUID `json:"id" db:"id"`
	SenderID   uuid.UUID `json:"sender_id" db:"sender_id"`
	ReceiverID uuid.UUID `json:"receiver_id" db:"receiver_id"`
	CreatedAt  time.Time `json:"created_at" db:"created_at"`
}

type UserAppLink struct {
	ID            uuid.UUID   `json:"id" db:"id"`
	UserID        uuid.UUID   `json:"user_id" db:"user_id"`
	Provider      AppProvider `json:"provider" db:"provider"`
	URL           string      `json:"url" db:"url"`
	ShowOnProfile bool        `json:"show_on_profile" db:"show_on_profile"`
	CreatedAt     time.Time   `json:"created_at" db:"created_at"`
	UpdatedAt     time.Time   `json:"updated_at" db:"updated_at"`
}
