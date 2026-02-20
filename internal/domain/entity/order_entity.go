package entity

type Order struct {
	Base
	UserID     uint32
	Status     string
	TotalPrice float64

	Items []*OrderItem
}

