package dto

type CreateHallReq struct {
	Name             string  `json:"name" binding:"required"`
	IconURL          *string `json:"icon_url" binding:"omitempty,url"`
	IconThumbnailURL *string `json:"icon_thumbnail_url" binding:"omitempty,url"`
	BannerColor      *string `json:"banner_color" binding:"omitempty"`
	Description      *string `json:"description" binding:"omitempty,max=500"`
}
