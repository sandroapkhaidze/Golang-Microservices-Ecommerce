package repository

import (
	"context"

	"github.com/sandroapkhaidze/Golang-Microservices-Ecommerce/order-service/internal/domain/entity"
)

type OrderRepository interface {
	Create(ctx context.Context, order *entity.Order) error
	GetByID(ctx context.Context, id string) (*entity.Order, error)
	GetByUserID(ctx context.Context, userID string, limit, offset int) ([]*entity.Order, error)
	UpdateStatus(ctx context.Context, orderID string, status entity.OrderStatus) error
	GetByCorrelationID(ctx context.Context, correlationID string) (*entity.Order, error)
}
