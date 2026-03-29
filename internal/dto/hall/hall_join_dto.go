package dto

import (
	"time"

	"github.com/google/uuid"
)

type JoinHallRes struct {
	MemberID uuid.UUID `json:"member_id"`
	HallID   uuid.UUID `json:"hall_id"`
	UserID   uuid.UUID `json:"user_id"`
	RoleID   uuid.UUID `json:"role_id"`
	Nickname *string   `json:"nickname"`
	JoinedAt time.Time `json:"joined_at"`
}
