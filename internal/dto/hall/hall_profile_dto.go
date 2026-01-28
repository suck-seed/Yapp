package dto

import "github.com/google/uuid"

// Contains basic information about hall
// - NAME, ICON, DESCRIPTION
// CAN UPDATE USING PROFILE DTO

type HallProfileUpdateReq struct {
	HallID      uuid.UUID `json:"hall_id" binding:"required"`
	Name        *string   `json:"name" binding:"omitempty"`
	Description *string   `json:"description" binding:"omitempty,max=500"`
}

type HallProfileUpdateRes struct {
}
