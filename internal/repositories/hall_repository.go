package repositories

type IHallRepository interface {
}

type hallRepository struct {
	db PGXTX
}

func NewHallRepository(db PGXTX) IHallRepository {

	return &hallRepository{
		db: db,
	}
}
