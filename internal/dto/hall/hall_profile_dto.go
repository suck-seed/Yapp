package dto

import "github.com/google/uuid"

// Contains basic information about hall
// - NAME, ICON, DESCRIPTION
// CAN UPDATE USING PROFILE DTO

type HallProfileReq struct {
	HallID      uuid.UUID `json:"hall_id"`
	Name        *string   `json:"name"`
	Description *string   `json:"description"`
}

type HallProfileRes struct {
}
