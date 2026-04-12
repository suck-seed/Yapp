package models

import (
	"time"

	"github.com/google/uuid"
)

type Hall struct {
	ID               uuid.UUID `json:"id" db:"id"`
	Name             string    `json:"name" db:"name"`
	IsPrivate        bool      `json:"is_private" db:"is_private"`
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
	ID        uuid.UUID `db:"id" json:"id"`
	HallID    uuid.UUID `db:"hall_id" json:"hall_id"`
	UserID    uuid.UUID `db:"user_id" json:"user_id"`
	RoleID    uuid.UUID `db:"role_id" json:"role_id"`
	Nickname  *string   `db:"nickname" json:"nickname"`
	JoinedAt  time.Time `db:"joined_at" json:"joined_at"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
	UpdatedAt time.Time `db:"updated_at" json:"updated_at"`
}
type HallRequest struct {
	ID        uuid.UUID `db:"id" json:"id"`
	HallID    uuid.UUID `db:"hall_id" json:"hall_id"`
	UserID    uuid.UUID `db:"user_id" json:"user_id"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
	UpdatedAt time.Time `db:"updated_at" json:"updated_at"`
}
