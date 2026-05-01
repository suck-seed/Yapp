package dto

import "github.com/google/uuid"

type GetUserPresenceRes struct {
	UserID   uuid.UUID `json:"user_id"`
	IsOnline bool      `json:"is_online"`
}

type BatchUserPresenceReq struct {
	UserIDs []uuid.UUID `json:"user_ids" binding:"required"`
}

type BatchUserPresenceItem struct {
	UserID   uuid.UUID `json:"user_id"`
	IsOnline bool      `json:"is_online"`
}

type BatchUserPresenceRes struct {
	Users []*BatchUserPresenceItem `json:"users"`
	Total int                      `json:"total"`
}
