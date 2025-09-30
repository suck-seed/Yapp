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

type User struct {
	ID                 uuid.UUID    `json:"id" db:"id"`
	Username           string       `json:"username" db:"username"`
	DisplayName        string       `json:"display_name" db:"display_name"`
	Email              string       `json:"email" db:"email"`
	Description        string       `json:"description,omitempty" db:"description"`
	PasswordHash       string       `json:"password_hash" db:"password_hash"`
	PhoneNumber        string       `json:"phone_number,omitempty" db:"phone_number"`
	AvatarURL          string       `json:"avatar_url,omitempty" db:"avatar_url"`
	AvatarThumbnailURL string       `json:"avatar_thumbnail_url,omitempty" db:"avatar_thumbnail_url"`
	FriendPolicy       FriendPolicy `json:"friend_policy" db:"friend_policy"`
	Active             bool         `json:"active" db:"active"`
	// LastSeen     *time.Time   `json:"last_seen,omitempty" db:"last_seen"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}
