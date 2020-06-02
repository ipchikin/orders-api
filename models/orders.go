package models

type OrdersModel struct {
	BaseModel
}

type Coordinates struct {
	Latitude  string
	Longitude string
}

type Orders struct {
	ID          int
	Origin      Coordinates
	Destination Coordinates
	Distance    int
	Status      string
}
