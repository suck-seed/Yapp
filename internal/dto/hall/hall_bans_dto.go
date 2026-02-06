package dto

import (
	"time"

	"github.com/google/uuid"
)

// BanUserRequest - POST  request mapping struct for banning a user
type BanUserRequest struct {
	UserID uuid.UUID `json:"user_id" binding:"required"`
	Reason string    `json:"reason" binding:"required,min=1, max-500"`
}

// BanUserResponse - POST response to a ban user request
type BanUserResponse struct {
	ID        uuid.UUID      `json:"id"`
	Reason    string         `json:"reason"`
	UserID    uuid.UUID      `json:"user_id"`
	User      BannedUserInfo `json:"user"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
}

// BannedUserInfo - information about the banned user
type BannedUserInfo struct {
	ID       uuid.UUID `json:"id"`
	Username string    `json:"username"`
	Avatar   *string   `json:"avatar,omitempty"`
}

// AllBannedUserRes - Gets All banned user in a specific hall
type AllBannedUserRes struct {
	Bans []BanSummaryResponse `json:"bans"`
}

// BanSummaryResponse - Gets Response summary for individual ban get query
type BanSummaryResponse struct {
	ID        uuid.UUID `json:"id"`
	UserID    uuid.UUID `json:"user_id"`
	Username  string    `json:"username"`
	AvatarURL *string   `json:"avatar_url,omitempty"`
	Reason    *string   `json:"reason"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
