package dto

import (
	"time"

	"github.com/google/uuid"
)

// BanUserReq - POST  request mapping struct for banning a user
type BanUserReq struct {
	UserID uuid.UUID `json:"user_id" binding:"required"`
	Reason string    `json:"reason" binding:"required,min=1,max=500"`
}

// BanUserRes - POST response to a ban user request
type BanUserRes struct {
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
	Bans []BanSummaryRes `json:"bans"`
}

// BanSummaryRes - Gets Response summary for individual ban get query
type BanSummaryRes struct {
	ID        uuid.UUID `json:"id"`
	UserID    uuid.UUID `json:"user_id"`
	Username  string    `json:"username"`
	AvatarURL *string   `json:"avatar_url"`
	Reason    *string   `json:"reason"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// UnbanRes - response after unbanning
type UnbanRes struct {
	UserID   uuid.UUID `json:"user_id"`
	Username string    `json:"username"`
	Message  string    `json:"message"`
}
