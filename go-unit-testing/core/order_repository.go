package core

// Secondary port
type OrderRepository interface {
	Save(order Order) error
}
