package storage

import (
	"errors"
	"test_pizza/internal/models"
)

var (
	ErrEmptyItems = errors.New("items must not be empty")
	ErrNotFoundOrder = errors.New("order not found")
	ErrOrderIsFinished = errors.New("can't complete the action because the order has already been completed")
)

type Storage interface {
	CreateOrder(items []int) (models.Order, error)
	AddItems(order_id string, items []int) error
	GetOrderById(order_id string) (models.Order, error)
	FinishOrder(order_id string) error
	GetOrdersByStatus(done int)  ([]models.Order, error)
}