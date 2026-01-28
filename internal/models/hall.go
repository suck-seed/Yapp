package models

import (
	"time"

	"github.com/google/uuid"
)

type Hall struct {
	ID               uuid.UUID `json:"id" db:"id"`
	Name             string    `json:"name" db:"name"`
	IconURL          *string   `json:"icon_url,omitempty" db:"icon_url"`
	IconThumbnailURL *string   `json:"icon_thumbnail_url,omitempty" db:"icon_thumbnail_url"`
	BannerColor      *string   `json:"banner_color,omitempty" db:"banner_color"`
	Description      *string   `json:"description,omitempty" db:"description"`
	CreatedAt        time.Time `json:"created_at" db:"created_at"`
	UpdatedAt        time.Time `json:"updated_at" db:"updated_at"`

	// createdby
	OwnerID uuid.UUID `json:"owner_id" db:"owner_id"`
}

type HallMember struct {
	ID        uuid.UUID `db:"id"`
	HallID    uuid.UUID `db:"hall_id"`
	UserID    uuid.UUID `db:"user_id"`
	RoleID    uuid.UUID `db:"role_id"`
	JoinedAt  time.Time `db:"joined_at"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

type Role struct {
	ID        uuid.UUID `db:"id"`
	HallID    uuid.UUID `db:"hall_id"`
	Name      string    `db:"name"`
	Color     *string   `db:"color"`
	IconURL   *string   `db:"icon_url"`
	IsDefault bool      `db:"is_default"`
	IsAdmin   bool      `db:"is_admin"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

type HallBans struct {
	ID       uuid.UUID `db:"id"`
	Reason   string    `db:"reason"`
	BannedBy uuid.UUID `db:"banned_by"`

	UserID uuid.UUID `db:"user_id"`
	HallID uuid.UUID `db:"hall_id"`

	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}
