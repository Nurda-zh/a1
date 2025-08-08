package repository

import (
	"github.com/Nurda-zh/a1/order-service/internal/domain"
)

type OrderRepo interface {
	Create(order *domain.Order) (string, error)
	GetByID(id string) (*domain.Order, error)
	UpdateStatus(id string, status domain.OrderStatus) error
	ListByUser(userID string, page, pageSize int64) ([]*domain.Order, int64, error)
}
