package dto

type CreateRoom struct {
	FloorID   *string `json:"floor_id" validate:"omitempty,uuid4"`
	Name      string  `json:"name" validate:"required,min=1,max=64"`
	RoomType  string  `json:"room_type" validate:"required,oneof=text voice"`
	IsPrivate bool    `json:"is_private"`
}
