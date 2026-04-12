package models

import (
	"github.com/google/uuid"
	"time"
)



type HallInvite struct{
	ID        uuid.UUID  `db:"id"`
    HallID    uuid.UUID  `db:"hall_id"`
    CreatedBy uuid.UUID  `db:"created_by"`
    Code      string     `db:"code"`
    RoleID    *uuid.UUID `db:"role_id"`
    // nil  →if no role assigned on join

    MaxUses   *int       `db:"max_uses"`
    // nil  → unlimited

    UsedCount int        `db:"used_count"` // keeps track of used
    ExpiresAt *time.Time `db:"expires_at"`
    // nil → never expires

    CreatedAt time.Time  `db:"created_at"`
}


// Helper functions
func (i *HallInvite) IsExpired() bool{
	if i.ExpiresAt == nil {
		return false
	}

	return time.Now().After(*i.ExpiresAt)
}

func (i *HallInvite) IsExhausted() bool {
    if i.MaxUses == nil {
        return false
    }
    return i.UsedCount >= *i.MaxUses
}

// IsValid is a convenience check for handlers/services.
func (i *HallInvite) IsValid() bool {
    return !i.IsExpired() && !i.IsExhausted()
}
