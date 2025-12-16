package persistence

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strconv"

	"github.com/google/uuid"
	"github.com/sandroapkhaidze/Golang-Microservices-Ecommerce/order-service/internal/domain/entity"
	"github.com/sandroapkhaidze/Golang-Microservices-Ecommerce/order-service/internal/domain/repository"
	sqlc "github.com/sandroapkhaidze/Golang-Microservices-Ecommerce/order-service/internal/infrastructure/persistence/sqlc"
)

type PostgresOrderRepository struct {
	queries *sqlc.Queries
	db      *sql.DB
}

func NewPostgresOrderRepository(db *sql.DB) repository.OrderRepository {
	return &PostgresOrderRepository{
		queries: sqlc.New(db),
		db:      db,
	}
}

func (p *PostgresOrderRepository) Create(ctx context.Context, order *entity.Order) error {
	orderUUID, err := uuid.Parse(order.ID)
	if err != nil {
		return errors.New("invalid order ID format")
	}

	correlationUUID, err := uuid.Parse(order.CorrelationID)
	if err != nil {
		return errors.New("invalid correlation ID format")
	}

	tx, err := p.db.Begin()
	if err != nil {
		return fmt.Errorf("could not start transaction: %w", err)
	}

	defer tx.Rollback()

	qtx := p.queries.WithTx(tx)

	err = qtx.CreateOrder(ctx, sqlc.CreateOrderParams{
		ID:            orderUUID,
		UserID:        uuid.MustParse(order.UserID),
		Status:        string(order.Status),
		TotalAmount:   fmt.Sprintf("%.2f", order.TotalAmount),
		CorrelationID: correlationUUID,
		CreatedAt:     order.CreatedAt,
		UpdatedAt:     order.UpdatedAt,
	})

	if err != nil {
		return fmt.Errorf("could not create order: %w", err)
	}

	for _, item := range order.Items {
		itemUUID, err := uuid.Parse(item.ID)
		if err != nil {
			return errors.New("invalid item ID format")
		}

		productUUID, err := uuid.Parse(item.ProductID)
		if err != nil {
			return errors.New("invalid product ID format")
		}

		err = qtx.CreateOrderItem(ctx, sqlc.CreateOrderItemParams{
			ID:        itemUUID,
			OrderID:   orderUUID,
			ProductID: productUUID,
			Quantity:  int32(item.Quantity),
			Price:     fmt.Sprintf("%.2f", item.Price),
		})

		if err != nil {
			return fmt.Errorf("could not create item: %w", err)
		}
	}
	return tx.Commit()
}

func (p *PostgresOrderRepository) GetByID(ctx context.Context, id string) (*entity.Order, error) {
	orderUUID, err := uuid.Parse(id)
	if err != nil {
		return nil, errors.New("invalid order ID format")
	}
	orderRow, err := p.queries.GetOrderByID(ctx, orderUUID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("order not found")
		}
		return nil, fmt.Errorf("could not get order: %w", err)
	}

	itemRows, err := p.queries.GetOrderItemsByOrderID(ctx, orderUUID)
	if err != nil {
		return nil, fmt.Errorf("could not get order items: %w", err)
	}

	order := toOrderEntity(orderRow)
	order.Items = toOrderItemEntities(itemRows)

	return order, nil
}

func (p *PostgresOrderRepository) GetByUserID(ctx context.Context, userID string, limit, offset int) ([]*entity.Order, error) {
	userUUID, err := uuid.Parse(userID)
	if err != nil {
		return nil, errors.New("invalid user ID format")
	}

	orderParams := sqlc.GetOrdersByUserIDParams{
		UserID: userUUID,
		Limit:  int32(limit),
		Offset: int32(offset),
	}

	orderRow, err := p.queries.GetOrdersByUserID(ctx, orderParams)
	if err != nil {
		return nil, fmt.Errorf("could not get order: %w", err)
	}

	orders := make([]*entity.Order, len(orderRow))
	for i, item := range orderRow {
		orderItems, err := p.queries.GetOrderItemsByOrderID(ctx, item.ID)
		if err != nil {
			return nil, fmt.Errorf("could not get order items: %w", err)
		}

		order := toOrderEntity(item)
		order.Items = toOrderItemEntities(orderItems)
		orders[i] = order
	}
	return orders, nil
}

func (p *PostgresOrderRepository) UpdateStatus(ctx context.Context, orderID string, status entity.OrderStatus) error {
	orderUUID, err := uuid.Parse(orderID)
	if err != nil {
		return errors.New("invalid order ID format")
	}
	orderParams := sqlc.UpdateOrderStatusParams{
		ID:     orderUUID,
		Status: string(status),
	}

	err = p.queries.UpdateOrderStatus(ctx, orderParams)
	if err != nil {
		return fmt.Errorf("could not update order status: %w", err)
	}
	return nil
}

func (p *PostgresOrderRepository) GetByCorrelationID(ctx context.Context, correlationID string) (*entity.Order, error) {
	correlationUUID, err := uuid.Parse(correlationID)
	if err != nil {
		return nil, errors.New("invalid correlation ID format")
	}

	orderRow, err := p.queries.GetOrderByCorrelationID(ctx, correlationUUID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("order not found")
		}
		return nil, fmt.Errorf("could not get order: %w", err)
	}

	orderItems, err := p.queries.GetOrderItemsByOrderID(ctx, orderRow.ID)
	if err != nil {
		return nil, fmt.Errorf("could not get order items: %w", err)
	}

	order := toOrderEntity(orderRow)
	order.Items = toOrderItemEntities(orderItems)
	return order, nil
}

func toOrderEntity(row sqlc.Order) *entity.Order {
	return &entity.Order{
		ID:            row.ID.String(),
		UserID:        row.UserID.String(),
		Status:        entity.OrderStatus(row.Status),
		TotalAmount:   parseDecimal(row.TotalAmount), // string → float64
		CorrelationID: row.CorrelationID.String(),
		CreatedAt:     row.CreatedAt,
		UpdatedAt:     row.UpdatedAt,
		Items:         []entity.OrderItem{}, // Will be filled separately
	}
}

func toOrderItemEntities(rows []sqlc.OrderItem) []entity.OrderItem {
	items := make([]entity.OrderItem, len(rows))
	for i, row := range rows {
		items[i] = entity.OrderItem{
			ID:        row.ID.String(),
			OrderID:   row.OrderID.String(),
			ProductID: row.ProductID.String(),
			Quantity:  int(row.Quantity),
			Price:     parseDecimal(row.Price), // string → float64
		}
	}
	return items
}

func parseDecimal(s string) float64 {
	f, _ := strconv.ParseFloat(s, 64)
	return f
}
