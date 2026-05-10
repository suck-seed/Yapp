package models

import (
	"time"

	"github.com/google/uuid"
	dto "github.com/suck-seed/yapp/internal/dto/room"
)

type RoomType string

const (
	AudioRoom RoomType = "audio"
	TextRoom  RoomType = "text"
)

type Room struct {
	ID                   uuid.UUID  `json:"id"                  db:"id"`
	HallID               uuid.UUID  `json:"hall_id"              db:"hall_id"`
	FloorID              *uuid.UUID `json:"floor_id,omitempty"   db:"floor_id"`
	Name                 string     `json:"name"                 db:"name"`
	RoomType             string     `json:"room_type"            db:"room_type"`
	Position             float64    `json:"position"             db:"position"`
	IsPrivate            bool       `json:"is_private"           db:"is_private"`
	SyncWithFloorMembers bool       `json:"sync_with_floor_members"    db:"sync_with_floor_members"`

	CreatedAt time.Time `json:"created_at"           db:"created_at"`
	UpdatedAt time.Time `json:"updated_at"           db:"updated_at"`
}

type RoomMember struct {
	RoomID   uuid.UUID `json:"room_id" db:"room_id"`
	MemberID uuid.UUID `json:"member_id" db:"member_id"`
}

func roomToRes(r *Room) dto.RoomRes {
	return dto.RoomRes{
		ID:                   r.ID,
		HallID:               r.HallID,
		FloorID:              r.FloorID,
		Name:                 r.Name,
		RoomType:             r.RoomType,
		Position:             r.Position,
		IsPrivate:            r.IsPrivate,
		SyncWithFloorMembers: r.SyncWithFloorMembers,
		CreatedAt:            r.CreatedAt,
		UpdatedAt:            r.UpdatedAt,
	}
}
