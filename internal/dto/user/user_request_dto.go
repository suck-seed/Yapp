package dto

import (
	"github.com/google/uuid"
	"github.com/suck-seed/yapp/internal/models"
)

// REQUESTS

type SignupUserReq struct {
	Username    string `json:"username" binding:"required,min=3,max=32"`
	Password    string `json:"password" binding:"required,min=8"`
	Email       string `json:"email" binding:"required,email"`
	DisplayName string `json:"display_name" binding:"required,min=1,max=64"`
}

type SigninUserReq struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type UpdateUserMeReq struct {
	DisplayName        *string              `json:"display_name" binding:"omitempty,min=1,max=64"`
	Description        *string              `json:"description" binding:"omitempty,max=200"`
	AvatarURL          *string              `json:"avatar_url" binding:"omitempty,url"`
	AvatarThumbnailURL *string              `json:"avatar_thumbnail_url" binding:"omitempty,url"`
	FriendPolicy       *models.FriendPolicy `json:"friend_policy" binding:"omitempty,oneof=everyone friends no_one"`
}

type UpdateUsernameReq struct {
	NewUsername string `json:"new_username" binding:"required,min=3,max=32"`
}

type UpdateEmailReq struct {
	NewEmail string `json:"new_email" binding:"required,email"`
}

type SendFriendRequestReq struct {
	ReceiverID uuid.UUID `json:"receiver_id" binding:"required"`
}

type RespondFriendRequestReq struct {
	Action string `json:"action" binding:"required,oneof=accept decline"`
}

type UpsertAppLinkReq struct {
	Provider models.AppProvider `json:"provider" binding:"required,oneof=spotify reddit twitter steam"`
	URL      string             `json:"url" binding:"required,url"`
	Show     bool               `json:"show_on_profile"`
}
