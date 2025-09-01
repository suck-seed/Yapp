package repositories

type IFloorRepository interface {
}

type floorRepository struct {
	db PGXTX
}

func NewFloorRepository(db PGXTX) IFloorRepository {

	return &floorRepository{
		db: db,
	}

}
