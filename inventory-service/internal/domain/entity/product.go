package entity

import (
	"errors"
	"time"
)

var (
	ErrInsufficientStock             = errors.New("insufficient stock available")
	ErrProductNotActive              = errors.New("product is not active")
	ErrInvalidQuantity               = errors.New("quantity must be positive")
	ErrCannotReleaseMoreThanReserved = errors.New("cannot release more than reserved")
	ErrCannotConfirmMoreThanReserved = errors.New("cannot confirm more than reserved")
)

type Product struct {
	ID            string    `json:"id"`
	Name          string    `json:"name"`
	Description   string    `json:"description"`
	Price         float64   `json:"price"`
	StockQuantity int32     `json:"stock_quantity"`
	ReservedStock int32     `json:"reserved_stock"`
	IsActive      bool      `json:"is_active"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

func (p *Product) AvailableStock() int32 {
	return p.StockQuantity - p.ReservedStock
}

func (p *Product) CanReserve(quantity int32) error {
	if !p.IsActive {
		return ErrProductNotActive
	}
	if quantity <= 0 {
		return ErrInvalidQuantity
	}
	if p.AvailableStock() < quantity {
		return ErrInsufficientStock
	}
	return nil
}

func (p *Product) ReserveStock(quantity int32) error {
	if err := p.CanReserve(quantity); err != nil {
		return err
	}
	p.ReservedStock += quantity
	p.UpdatedAt = time.Now().UTC()
	return nil
}

func (p *Product) ConfirmSale(quantity int32) error {
	if quantity <= 0 {
		return ErrInvalidQuantity
	}
	if p.ReservedStock < quantity {
		return ErrCannotConfirmMoreThanReserved
	}
	p.StockQuantity -= quantity
	p.ReservedStock -= quantity
	p.UpdatedAt = time.Now().UTC()
	return nil
}

func (p *Product) ReleaseReservation(quantity int32) error {
	if quantity <= 0 {
		return ErrInvalidQuantity
	}
	if p.ReservedStock < quantity {
		return ErrCannotReleaseMoreThanReserved
	}
	p.ReservedStock -= quantity
	p.UpdatedAt = time.Now().UTC()
	return nil
}
