package dto

import (
	"time"

	"github.com/google/uuid"
)

// All Banned User in Hall
type AllBannedUserRes struct {
	HallID uuid.UUID    `json:"hall_id"`
	Bans   []BannedUser `json:"bans"`
}

type BannedUser struct {
	BanID     uuid.UUID `json:"ban_id"`
	UserID    uuid.UUID `json:"user_id"`
	Username  string    `json:"username"`
	AvatarURL *string   `json:"avatar_url,omitempty"`
	Reason    *string   `json:"reason"`
}

// Banning an user
type BanUserReq struct {
	UserID uuid.UUID `json:"user_id"`
	Reason *string   `json:"reason,omitempty"`
}

type BanUserResponse struct {
	BanID     uuid.UUID `json:"ban_id"`
	Reason    *string   `json:"reason"`
	BannedBy  uuid.UUID `json:"banned_by"`
	UserID    uuid.UUID `json:"user_id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
