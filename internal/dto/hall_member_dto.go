package dto

import (
	"time"

	"github.com/google/uuid"
)

type AddUserToHallReq struct {
	UserID uuid.UUID `json:"user_id" binding:"required"`
	Role   string    `json:"role" binding:"required,oneof=member admin owner"`
}

type HallUserRes struct {
	HallID   uuid.UUID `json:"hall_id"`
	UserID   uuid.UUID `json:"user_id"`
	Role     string    `json:"role"`
	JoinedAt time.Time `json:"joined_at"`
}

type HallMemberRes struct {
	UserID    uuid.UUID `json:"user_id"`
	Username  string    `json:"username"`
	AvatarURL *string   `json:"avatar_url"`
	Role      string    `json:"role"`
	JoinedAt  time.Time `json:"joined_at"`
}

type GetHallMembersRes struct {
	Members []HallMemberRes `json:"members"`
	Total   int             `json:"total"`
}
