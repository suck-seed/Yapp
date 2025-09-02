package dto

type CreateRoomReq struct {
	// As usual we will fetch the creating user ID from

	HallID string `json:"hall_id" binding:"required"`

	// FloorID can be empty , as room can be in surface level and not in a floor
	FloorID  *string `json:"floor_id,omitempty"`
	Name     string  `json:"name" binding:"required,min=1,max=64"`
	RoomType string  `json:"room_type" binding:"required,oneof=text audio"`

	//  False by default
	IsPrivate bool `json:"is_private"`
}

type CreateRoomRes struct {
	RoomId    string  `json:"room_id"`
	HallID    string  `json:"hall_id"`
	FloorID   *string `json:"floor_id,omitempty"`
	Name      string  `json:"name"`
	RoomType  string  `json:"room_type"`
	IsPrivate bool    `json:"is_private"`
	CreatedAt string  `json:"created_at"`
}

type UpdateRoomReq struct {
	RoomId    string  `json:"room_id"`
	Name      *string `json:"name" binding:"omitempty,min=1,max=64"`
	IsPrivate *bool   `json:"is_private" bindling:"omitempty"`
}

type GetRoomRes struct {
	RoomId    string  `json:"room_id"`
	HallID    string  `json:"hall_id"`
	FloorID   *string `json:"floor_id,omitempty"`
	Name      string  `json:"name"`
	RoomType  string  `json:"room_type"`
	IsPrivate bool    `json:"is_private"`
	CreatedAt string  `json:"created_at"`
	UpdatedAt string  `json:"updated_at"`
}

type GetRoomsByHallReq struct {
	HallID string `json:"hall_id" binding:"required"`
}

type GetRoomsRes struct {
	Rooms []GetRoomRes `json:"rooms"`
	Count int          `json:"count"`
}
