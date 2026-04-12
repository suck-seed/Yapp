// internal/api/rest/dto/invite.go
package dto

import (
    "time"
    "github.com/google/uuid"
)

// ---------- enums / constants ----------

type ExpireAfterOption string

const (
    Expire30Min ExpireAfterOption = "30min"
    Expire1Hr   ExpireAfterOption = "1hr"
    Expire6Hr   ExpireAfterOption = "6hr"
    Expire12Hr  ExpireAfterOption = "12hr"
    Expire1Day  ExpireAfterOption = "1day"
    Expire7Days ExpireAfterOption = "7days"
    ExpireNever ExpireAfterOption = "never"
)

var ValidExpireOptions = map[ExpireAfterOption]bool{
    Expire30Min: true, Expire1Hr: true, Expire6Hr: true,
    Expire12Hr: true, Expire1Day: true, Expire7Days: true,
    ExpireNever: true,
}

// nil MaxUses in request = no limit
var ValidMaxUses = map[int]bool{
    1: true, 5: true, 10: true, 25: true, 50: true, 100: true,
}

// ---------- requests ----------

type CreateInviteLinkReq struct {
    ExpireAfter ExpireAfterOption `json:"expire_after"` // required
    MaxUses     *int              `json:"max_uses"`     // nil = no limit; else must be in ValidMaxUses
    RoleID      *uuid.UUID        `json:"role_id"`      // nil = no role assigned
}

// ---------- responses ----------

// InviteLinkRes is used for management endpoints (list, create, revoke).
type InviteLinkRes struct {
    ID        uuid.UUID  `json:"id"`
    HallID    uuid.UUID  `json:"hall_id"`
    CreatedBy uuid.UUID  `json:"created_by"`
    Code      string     `json:"code"`
    URL       string     `json:"url"`       // full join URL, assembled by service
    RoleID    *uuid.UUID `json:"role_id"`
    MaxUses   *int       `json:"max_uses"`
    UsedCount int        `json:"used_count"`
    ExpiresAt *time.Time `json:"expires_at"`
    CreatedAt time.Time  `json:"created_at"`
    IsValid   bool       `json:"is_valid"`
}

// InviteInfoRes is the public-facing response when a user opens an invite link.
// No omitempty — every field is always present.
type InviteInfoRes struct {
    Code      string     `json:"code"`
    HallID    uuid.UUID  `json:"hall_id"`
    HallName  string     `json:"hall_name"`
    HallImage string     `json:"hall_image"`
    RoleID    *uuid.UUID `json:"role_id"`
    RoleName  string     `json:"role_name"`  // empty string when no role, not omitted
    MaxUses   *int       `json:"max_uses"`
    UsedCount int        `json:"used_count"`
    ExpiresAt *time.Time `json:"expires_at"`
    IsValid   bool       `json:"is_valid"`
}

// AcceptInviteLinkRes is returned after a successful join via invite.
type AcceptInviteLinkRes struct {
    HallID   uuid.UUID `json:"hall_id"`
    MemberID uuid.UUID `json:"member_id"`
    RoleID   *uuid.UUID `json:"role_id"`
    JoinedAt time.Time `json:"joined_at"`
}
