package services

type IMessageService interface {
}

type messageService struct {
}

func NewMessageService() IHallService {
	return &hallService{}
}
