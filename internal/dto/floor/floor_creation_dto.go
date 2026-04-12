package dto

import (
	"time"

	"github.com/google/uuid"
)

// ── Create ────────────────────────────────────────────────────────────────────

type CreateFloorReq struct {
	HallID    uuid.UUID `json:"hall_id"    binding:"required"`
	Name      string    `json:"name"       binding:"required"`
	IsPrivate *bool     `json:"is_private" binding:"required"`
}

type CreateFloorRes struct {
	ID        uuid.UUID `json:"id"`
	HallID    uuid.UUID `json:"hall_id"`
	Name      string    `json:"name"`
	Position  float64   `json:"position"`
	IsPrivate bool      `json:"is_private"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// ── Get Single ────────────────────────────────────────────────────────────────

type GetFloorRes struct {
	ID        uuid.UUID `json:"id"`
	HallID    uuid.UUID `json:"hall_id"`
	Name      string    `json:"name"`
	Position  float64   `json:"position"`
	IsPrivate bool      `json:"is_private"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// ── Get All (for a hall) ──────────────────────────────────────────────────────

type GetFloorsReq struct {
	HallID uuid.UUID `form:"hall_id" binding:"required"`
}

type GetFloorsRes struct {
	Floors []GetFloorRes `json:"floors"`
}

// ── Update ────────────────────────────────────────────────────────────────────

type UpdateFloorReq struct {
	// Both fields are optional — at least one must be present (enforced in service)
	Name      *string `json:"name"`
	IsPrivate *bool   `json:"is_private"`
}

type UpdateFloorRes struct {
	ID        uuid.UUID `json:"id"`
	HallID    uuid.UUID `json:"hall_id"`
	Name      string    `json:"name"`
	Position  float64   `json:"position"`
	IsPrivate bool      `json:"is_private"`
	UpdatedAt time.Time `json:"updated_at"`
}

// ── Reorder ───────────────────────────────────────────────────────────────────

// AfterID = nil → place at very top
// AfterID = uuid → place immediately after that floor
type MoveFloorReq struct {
	HallID  uuid.UUID  `json:"hall_id" binding:"required"`
	AfterID *uuid.UUID `json:"after_id" binding:"omitempty"`
}
