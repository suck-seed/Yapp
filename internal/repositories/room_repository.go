package repositories

type IRoomRepository interface {
}

type roomRepository struct {
	db PGXTX
}

func NewRoomReposiroty(db PGXTX) IHallRepository {

	return &hallRepository{
		db: db,
	}
}
