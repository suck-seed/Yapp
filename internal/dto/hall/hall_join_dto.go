package dto

import (
	"time"

	"github.com/google/uuid"
)

// POST /halls/:hallID/join
type JoinHallRes struct {
	Status    string     `json:"status"` // "joined" or "requested"
	MemberID  *uuid.UUID `json:"member_id"`
	RequestID *uuid.UUID `json:"request_id"`
	HallID    uuid.UUID  `json:"hall_id"`
	UserID    uuid.UUID  `json:"user_id"`
	RoleID    *uuid.UUID `json:"role_id"`
	Nickname  *string    `json:"nickname"`
	JoinedAt  *time.Time `json:"joined_at"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
}

// GET /halls/:hallID/settings/requests
type HallRequestRes struct {
	ID        uuid.UUID `json:"id"`
	HallID    uuid.UUID `json:"hall_id"`
	UserID    uuid.UUID `json:"user_id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type GetCurrentRequestsRes struct {
	Requests []*HallRequestRes `json:"requests"`
	Total    int               `json:"total"`
}

// PATCH /halls/:hallID/settings/requests/:requestID/accept
type AcceptJoinRequestRes struct {
	RequestID uuid.UUID `json:"request_id"`
	MemberID  uuid.UUID `json:"member_id"`
	HallID    uuid.UUID `json:"hall_id"`
	UserID    uuid.UUID `json:"user_id"`
	RoleID    uuid.UUID `json:"role_id"`
	JoinedAt  time.Time `json:"joined_at"`
}

// DELETE /halls/:hallID/settings/requests/:requestID
type DeclineJoinRequestRes struct {
	ID        uuid.UUID `json:"id"`
	HallID    uuid.UUID `json:"hall_id"`
	UserID    uuid.UUID `json:"user_id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
