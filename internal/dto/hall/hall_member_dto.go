package dto

import (
	"time"

	"github.com/google/uuid"
	dto "github.com/suck-seed/yapp/internal/dto/user"
)

// -------------------- GET MEMBERS

type HallMemberRes struct {
	ID        uuid.UUID `json:"id"`
	HallID    uuid.UUID `json:"hall_id"`
	UserID    uuid.UUID `json:"user_id"`
	RoleID    uuid.UUID `json:"role_id"`
	Nickname  *string   `json:"nickname"`
	JoinedAt  time.Time `json:"joined_at"`
	UpdatedAt time.Time `json:"updated_at"`

	// Presence of Specific Hall HallMember
	Presence *dto.UserPresenceRes `json:"presence"`
}

type GetHallMembersRes struct {
	Members []*HallMemberRes `json:"members"`
	Total   int              `json:"total"`
}

// -------------------- UPDATE MEMBER ROLE
type UpdateHallMemberRoleReq struct {
	RoleID uuid.UUID `json:"role_id" binding:"required"`
}

// -------------------- UPDATE MEMBER NICKNAME
type UpdateHallMemberNicknameReq struct {
	// Nickname is required in the JSON body; send "" to clear. Omitted or null returns "no field to update".
	Nickname *string `json:"nickname" binding:"omitempty,max=64"`
}

// shared response — same for both
type UpdateHallMemberRes struct {
	ID        uuid.UUID `json:"id"`
	HallID    uuid.UUID `json:"hall_id"`
	UserID    uuid.UUID `json:"user_id"`
	RoleID    uuid.UUID `json:"role_id"`
	Nickname  *string   `json:"nickname"`
	JoinedAt  time.Time `json:"joined_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
