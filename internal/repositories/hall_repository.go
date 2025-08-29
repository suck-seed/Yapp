package repositories

type IHallRepository interface {
}

type hallRepository struct {
	db PGXTX
}

func NewHallReposiroty(db PGXTX) IHallRepository {

	return &hallRepository{
		db: db,
	}
}
