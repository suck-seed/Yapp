package dto

type CreateRoom struct {
	FloorID   *string `json:"floor_id" binding:"omitempty,uuid4"`
	Name      string  `json:"name" binding:"required,min=1,max=64"`
	RoomType  string  `json:"room_type" binding:"required,oneof=text voice"`
	IsPrivate bool    `json:"is_private"`
}
