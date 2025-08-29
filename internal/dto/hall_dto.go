package dto

type CreateHall struct {
	Name        string `json:"name" validate:"required"`
	IconURL     string `json:"icon_url" validate:"omitempty,url"`
	BannerColor string `json:"banner_color" validate:"omitempty"`
	Description string `json:"description" validate:"omitempty,max=500"`
}
