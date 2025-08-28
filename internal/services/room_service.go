package services

type IRoomService interface {
}

type roomService struct {
}

func NewRoomService() IRoomService {
	return &roomService{}
}
