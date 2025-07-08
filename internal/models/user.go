package models

import (
	"github.com/google/uuid"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// TODO DO proper implementation here
type User struct {
	UserId    uuid.UUID             `json:"user_id"`
	Username  string                `json:"username"`
	Avatar    string                `json:"avatar"`
	Email     string                `json:"email"`
	Password  string                `json:"password"`
	Locale    string                `json:"locale"`
	Number    string                `json:"number"`
	CreatedAt timestamppb.Timestamp `json:"created_at"`
}
