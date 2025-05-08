package core

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// Mock implementation of OrderRepository
type mockOrderRepo struct {
	saveFunc func(order Order) error
}

func (m *mockOrderRepo) Save(order Order) error {
	return m.saveFunc(order)
}

func TestCreateOrder(t *testing.T) {
	// Success case
	t.Run("success", func(t *testing.T) {
		repo := &mockOrderRepo{
			saveFunc: func(order Order) error {
				// Simulate successful save
				return nil
			},
		}
		service := NewOrderService(repo)

		err := service.CreateOrder(Order{Total: 100})
		assert.NoError(t, err)
	})

	// Failure case: Total less than 0
	t.Run("total less than 0", func(t *testing.T) {
		repo := &mockOrderRepo{
			saveFunc: func(order Order) error {
				// This won't be called due to validation
				return nil
			},
		}
		service := NewOrderService(repo)

		err := service.CreateOrder(Order{Total: -10})
		assert.Error(t, err)
		assert.Equal(t, "total must be positive", err.Error())
	})
}
