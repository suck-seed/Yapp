package dto

import (
	"time"

	"github.com/google/uuid"
)

// -------------------- GET MEMBERS

type HallMemberRes struct {
	ID        uuid.UUID  `json:"id"`
	HallID    uuid.UUID  `json:"hall_id"`
	UserID    uuid.UUID  `json:"user_id"`
	RoleID    *uuid.UUID `json:"role_id"`
	Nickname  *string    `json:"nickname"`
	JoinedAt  time.Time  `json:"joined_at"`
	UpdatedAt time.Time  `json:"updated_at"`
}

type GetHallMembersRes struct {
	Members []*HallMemberRes `json:"members"`
	Total   int              `json:"total"`
}
