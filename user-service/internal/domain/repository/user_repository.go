package repository

import (
	"github.com/sandroapkhaidze/Golang-Microservices-Ecommerce/user-service/internal/domain/entity"

	"context"
)

// UserRepository defines the interface for user data operations
// This is an INTERFACE - implementations will be in infrastructure layer
type UserRepository interface {
	// Create creates a new user
	Create(ctx context.Context, user *entity.User) error

	// GetByID retrieves a user by ID
	GetByID(ctx context.Context, id string) (*entity.User, error)

	// GetByEmail retrieves a user by email
	GetByEmail(ctx context.Context, email string) (*entity.User, error)

	// Update updates an existing user
	Update(ctx context.Context, user *entity.User) error

	// Delete soft deletes a user (sets IsActive to false)
	Delete(ctx context.Context, id string) error

	// List retrieves all users with pagination
	List(ctx context.Context, limit, offset int) ([]*entity.User, error)
}
