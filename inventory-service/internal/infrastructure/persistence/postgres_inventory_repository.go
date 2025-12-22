package persistence

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"

	"github.com/google/uuid"
	"github.com/sandroapkhaidze/Golang-Microservices-Ecommerce/inventory-service/internal/domain/entity"
	sqlc "github.com/sandroapkhaidze/Golang-Microservices-Ecommerce/inventory-service/internal/infrastructure/persistence/sqlc"
)

type PostgresInventoryRepository struct {
	db      *sql.DB
	queries *sqlc.Queries
}

func NewPostgresInventoryRepository(db *sql.DB) *PostgresInventoryRepository {
	return &PostgresInventoryRepository{
		db:      db,
		queries: sqlc.New(db),
	}
}

// Create creates a new product
func (r *PostgresInventoryRepository) Create(ctx context.Context, product *entity.Product) error {
	uid, err := parseStringToUUID(product.ID)
	if err != nil {
		return errors.New("invalid user ID format")
	}

	return r.queries.CreateProduct(ctx, sqlc.CreateProductParams{
		ID:            uid,
		Name:          product.Name,
		Description:   sql.NullString{String: product.Description, Valid: product.Description != ""},
		Price:         fmt.Sprintf("%.2f", product.Price),
		StockQuantity: product.StockQuantity,
		ReservedStock: product.ReservedStock,
		IsActive:      product.IsActive,
	})
}

// GetByID retrieves a product by ID
func (r *PostgresInventoryRepository) GetByID(ctx context.Context, id string) (*entity.Product, error) {
	uid, err := parseStringToUUID(id)
	if err != nil {
		return nil, errors.New("invalid user ID format")
	}
	row, err := r.queries.GetProductByID(ctx, uid)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("product not found: %s", id)
		}
		return nil, err
	}

	return r.rowToEntity(row), nil
}

// Update updates an existing product
func (r *PostgresInventoryRepository) Update(ctx context.Context, product *entity.Product) error {
	uid, err := parseStringToUUID(product.ID)
	if err != nil {
		return errors.New("invalid user ID format")
	}
	return r.queries.UpdateProduct(ctx, sqlc.UpdateProductParams{
		ID:            uid,
		Name:          product.Name,
		Description:   sql.NullString{String: product.Description, Valid: product.Description != ""},
		Price:         fmt.Sprintf("%.2f", product.Price),
		StockQuantity: product.StockQuantity,
		ReservedStock: product.ReservedStock,
		IsActive:      product.IsActive,
	})
}

func (r *PostgresInventoryRepository) Delete(ctx context.Context, id string) error {
	uid, err := parseStringToUUID(id)
	if err != nil {
		return errors.New("invalid user ID format")
	}
	return r.queries.DeleteProduct(ctx, uid)
}

func (r *PostgresInventoryRepository) List(ctx context.Context, limit, offset int) ([]*entity.Product, error) {
	rows, err := r.queries.ListProducts(ctx, sqlc.ListProductsParams{
		Limit:  int32(limit),
		Offset: int32(offset),
	})
	if err != nil {
		return nil, err
	}

	products := make([]*entity.Product, len(rows))
	for i, row := range rows {
		products[i] = r.rowToEntity(row)
	}

	return products, nil
}

func (r *PostgresInventoryRepository) GetActiveProducts(ctx context.Context) ([]*entity.Product, error) {
	rows, err := r.queries.GetActiveProducts(ctx)
	if err != nil {
		return nil, err
	}

	products := make([]*entity.Product, len(rows))
	for i, row := range rows {
		products[i] = r.rowToEntity(row)
	}

	return products, nil
}

func (r *PostgresInventoryRepository) GetByIDs(ctx context.Context, ids []string) ([]*entity.Product, error) {
	log.Printf("DEBUG GetByIDs: received %d IDs: %v", len(ids), ids) // ADD THIS

	uids := make([]uuid.UUID, len(ids))
	for i, id := range ids {
		log.Printf("DEBUG: Parsing ID[%d]: %s", i, id) // ADD THIS

		uid, err := parseStringToUUID(id)
		if err != nil {
			log.Printf("DEBUG: Failed to parse ID[%d] '%s': %v", i, id, err) // ADD THIS
			return nil, errors.New("invalid product ID format")
		}
		uids[i] = uid
	}

	log.Printf("DEBUG: Successfully parsed all UUIDs, querying database") // ADD THIS

	rows, err := r.queries.GetProductsByIDs(ctx, uids)
	if err != nil {
		return nil, err
	}

	products := make([]*entity.Product, len(rows))
	for i, row := range rows {
		products[i] = r.rowToEntity(row)
	}

	return products, nil
}

// UpdateMultiple updates multiple products in a single transaction
func (r *PostgresInventoryRepository) UpdateMultiple(ctx context.Context, products []*entity.Product) error {
	// Start transaction
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback() // Rollback if not committed

	// Create queries instance for this transaction
	qtx := r.queries.WithTx(tx)
	// Update each product within the transaction
	for _, product := range products {
		uid, err := parseStringToUUID(product.ID)
		if err != nil {
			return fmt.Errorf("invalid user ID format: %w", err)
		}
		err = qtx.UpdateProduct(ctx, sqlc.UpdateProductParams{
			ID:            uid,
			Name:          product.Name,
			Description:   sql.NullString{String: product.Description, Valid: product.Description != ""},
			Price:         fmt.Sprintf("%.2f", product.Price),
			StockQuantity: product.StockQuantity,
			ReservedStock: product.ReservedStock,
			IsActive:      product.IsActive,
		})
		if err != nil {
			// Transaction will auto-rollback due to defer
			return fmt.Errorf("failed to update product %s: %w", product.ID, err)
		}
	}

	// Commit transaction
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

func (r *PostgresInventoryRepository) rowToEntity(row interface{}) *entity.Product {
	var product entity.Product
	// Type assertion to handle both single row and multiple rows
	switch v := row.(type) {
	case sqlc.Product:
		product.ID = v.ID.String()
		product.Name = v.Name
		product.Description = v.Description.String
		// Parse price from string to float64
		_, _ = fmt.Sscanf(v.Price, "%f", &product.Price)
		product.StockQuantity = v.StockQuantity
		product.ReservedStock = v.ReservedStock
		product.IsActive = v.IsActive
		product.CreatedAt = v.CreatedAt
		product.UpdatedAt = v.UpdatedAt
	}

	return &product
}

// parseStringToUUID converts string to uuid.UUID
func parseStringToUUID(id string) (uuid.UUID, error) {
	return uuid.Parse(id)
}
