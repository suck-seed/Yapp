package dto

// REQUESTS

type UserSignup struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
	Email    string `json:"email" validate:"required"`
	Number   string `json:"number" validate:"omitempty"`
	Display  string `json:"display_name" validate:"omitempty"`
	// Disolay Name if empty is same as the username
}

type UpsertAppLink struct {
	Provider string `json:"provider" validate:"required, oneof= spotify reddit twitter steam instagram"`
	// Valid Providers : spority reddit twitter steam instagram
	URL  string `json:"url" validate:"required"`
	Show bool   `json:"show_on_profile" validate:"omitempty"`
	// by default, Do not show, so can omit empty
}

type UserLogin struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}

// RESPONSE

type UserPublic struct {
	ID          string `json:"id"`
	Username    string `json:"username"`
	DisplayName string `json:"display_name"`
	Email       string `json:"email"`
	AvatarURL   string `json:"avatar_url,omitempty"`
}

//  to send back after logging in to a device
type AuthToken struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   int64  `json:"expires_in"`
}
