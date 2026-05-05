package models

import (
	"time"

	"github.com/google/uuid"
)

type PresenceStatus string

const (
	PresenceStatusOnline  PresenceStatus = "online"
	PresenceStatusOffline PresenceStatus = "offline"
	PresenceStatusAway    PresenceStatus = "away"
	PresenceStatusBusy    PresenceStatus = "busy"
)

type UserPresence struct {
	UserID     uuid.UUID      `json:"user_id"`
	Status     PresenceStatus `json:"status"`
	LastSeenAt *time.Time     `json:"last_seen_at,omitempty"`
	UpdatedAt  time.Time      `json:"updated_at"`
}
