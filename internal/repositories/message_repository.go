package repositories

type IMessageRepository interface {
}

type messageRepository struct {
	db PGXTX
}

func NewMessageReposiroty(db PGXTX) IMessageRepository {

	return &messageRepository{
		db: db,
	}
}
