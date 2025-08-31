package models

type Hub struct {

	// maps roomID to roomStruct
	Rooms map[string]*Room
}
