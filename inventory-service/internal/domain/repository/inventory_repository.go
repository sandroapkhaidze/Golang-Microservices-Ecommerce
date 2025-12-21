package repository

import (
	"context"

	"github.com/sandroapkhaidze/Golang-Microservices-Ecommerce/inventory-service/internal/domain/entity"
)

type InventoryRepository interface {
	Create(ctx context.Context, product *entity.Product) error
	GetByID(ctx context.Context, id string) (*entity.Product, error)
	Update(ctx context.Context, product *entity.Product) error
	Delete(ctx context.Context, id string) error
	List(ctx context.Context, limit, offset int) ([]*entity.Product, error)
	GetActiveProducts(ctx context.Context) ([]*entity.Product, error)
	GetByIDs(ctx context.Context, ids []string) ([]*entity.Product, error)
	UpdateMultiple(ctx context.Context, products []*entity.Product) error
}
