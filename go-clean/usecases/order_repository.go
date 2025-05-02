package usecases

import "potae/entities"

type OrderRepository interface {
	Save(order entities.Order) error
}
