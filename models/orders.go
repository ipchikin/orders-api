package models

import (
	"github.com/google/uuid"
)

type OrdersModel struct {
	BaseModel
}

// type Coordinates struct {
// 	Latitude  string
// 	Longitude string
// }

// type Orders struct {
// 	ID          int64
// 	Origin      Coordinates
// 	Destination Coordinates
// 	Distance    uint32
// 	Status      string
// }

// Place
func (om *OrdersModel) Place(origin, destination [2]string, distance uint32, status string) (id string, err error) {
	id = uuid.New().String()
	_, err = om.DB.Exec(`INSERT INTO orders (id, origin_lat, origin_long, destination_lat, destination_long, distance, status) VALUES (?, ?, ?, ?, ?, ?, ?)`,
		id, origin[0], origin[1], destination[0], destination[1], distance, status)
	if err != nil {
		return
	}

	return
}
