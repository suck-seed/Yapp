package services

type IHallService interface {
}

type hallService struct {
}

func NewHallService() IMessageService {
	return &messageService{}
}
