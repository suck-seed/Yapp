package dto

import (
	"time"

	"github.com/google/uuid"
	"github.com/suck-seed/yapp/internal/models"
)

type UpdatePresenceReq struct {
	Status models.PresenceStatus `json:"status" binding:"required"`
}

type UserPresenceRes struct {
	UserID     uuid.UUID             `json:"user_id"`
	Status     models.PresenceStatus `json:"status"`
	LastSeenAt *time.Time            `json:"last_seen_at,omitempty"`
	UpdatedAt  time.Time             `json:"updated_at"`
}

type GetManyPresenceReq struct {
	UserIDs []uuid.UUID `json:"user_ids" binding:"required"`
}

type GetManyPresenceRes struct {
	Presences []*UserPresenceRes `json:"presences"`
}
