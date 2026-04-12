package dto

import (
	"time"

	"github.com/google/uuid"
)

// Shared response shape
type RoomRes struct {
	ID        uuid.UUID  `json:"id"`
	HallID    uuid.UUID  `json:"hall_id"`
	FloorID   *uuid.UUID `json:"floor_id"`
	Name      string     `json:"name"`
	RoomType  string     `json:"room_type"`
	Position  float64    `json:"position"`
	IsPrivate bool       `json:"is_private"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
}

type CreateRoomReq struct {
	HallID    uuid.UUID  `json:"hall_id"    binding:"required"`
	FloorID   *uuid.UUID `json:"floor_id"   binding:"omitempty"`
	Name      string     `json:"name"       binding:"required,min=1,max=64"`
	RoomType  string     `json:"room_type"  binding:"required,oneof=text audio"`
	IsPrivate *bool      `json:"is_private" binding:"omitempty"`
}

type CreateRoomRes = RoomRes
type GetRoomRes = RoomRes
type UpdateRoomRes = RoomRes

type GetHallRoomsReq struct {
	HallID uuid.UUID `form:"hall_id" binding:"required"`
}

type FloorWithRoomsRes struct {
	ID        uuid.UUID `json:"id"`
	HallID    uuid.UUID `json:"hall_id"`
	Name      string    `json:"name"`
	Position  float64   `json:"position"`
	IsPrivate bool      `json:"is_private"`
	Rooms     []RoomRes `json:"rooms"`
}

type GetHallRoomsRes struct {
	TopLevel []RoomRes           `json:"top_level"`
	Floors   []FloorWithRoomsRes `json:"floors"`
}

type UpdateRoomReq struct {
	Name      *string `json:"name"       binding:"omitempty,min=1,max=64"`
	IsPrivate *bool   `json:"is_private" binding:"omitempty"`
}

// ── The only drag-drop endpoint you need ──────────────────────────────────────
//
// Covers ALL cases in one call:
//   - reorder within same container  → same new_floor_id, different after_id
//   - move into a floor + place      → new_floor_id = floorUUID, after_id = where to land
//   - move out of floor + place      → new_floor_id = nil, after_id = where to land
//
// after_id = nil  →  insert at the very top of the target container
// after_id = uuid →  insert immediately after that room

type MoveRoomReq struct {
	HallID     uuid.UUID  `json:"hall_id"     binding:"required"`
	NewFloorID *uuid.UUID `json:"new_floor_id" binding:"omitempty"` // nil = top-level
	AfterID    *uuid.UUID `json:"after_id"    binding:"omitempty"`  // nil = place at top
}
