package entity

type OrderItem struct {
	Base
	OrderID   uint32
	ProductID string
	Quantity  int
	Price     float64
	Subtotal  float64

	Order *Order
}
