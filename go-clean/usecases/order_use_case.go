package usecases

import "potae/entities"

type OrderUseCase interface {
	CreateOrder(order entities.Order) error
}

type OrderServiec struct {
	repo OrderRepository
}

func NewOrderService(repo OrderRepository) OrderUseCase {
	return &OrderServiec{
		repo: repo,
	}
}

func (s *OrderServiec) CreateOrder(order entities.Order) error {
	if err := s.repo.Save(order); err != nil {
		return err
	}
	return nil
}
