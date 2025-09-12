package repositories

type IRoomRepository interface {
}

type roomRepository struct {
	db PGXTX
}

func NewRoomRepository(db PGXTX) IRoomRepository {

	return &roomRepository{
		db: db,
	}
}
