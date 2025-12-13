package persistence

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/sandroapkhaidze/Golang-Microservices-Ecommerce/user-service/internal/domain/entity"
	"github.com/sandroapkhaidze/Golang-Microservices-Ecommerce/user-service/internal/domain/repository"
	sqlc "github.com/sandroapkhaidze/Golang-Microservices-Ecommerce/user-service/internal/infrastructure/persistence/sqlc"
)

// PostgresUserRepository implements repository.UserRepository using sqlc
type PostgresUserRepository struct {
	queries *sqlc.Queries
}

// NewPostgresUserRepository creates a new PostgreSQL user repository
func NewPostgresUserRepository(db *sql.DB) repository.UserRepository {
	return &PostgresUserRepository{
		queries: sqlc.New(db),
	}
}

// parseStringToUUID converts string to uuid.UUID
func parseStringToUUID(id string) (uuid.UUID, error) {
	return uuid.Parse(id)
}

// Create creates a new user in the database
func (r *PostgresUserRepository) Create(ctx context.Context, user *entity.User) error {
	uid, err := parseStringToUUID(user.ID)
	if err != nil {
		return errors.New("invalid user ID format")
	}

	err = r.queries.CreateUser(ctx, sqlc.CreateUserParams{
		ID:        uid,
		Email:     user.Email,
		Password:  user.Password,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Role:      user.Role,
		IsActive:  user.IsActive,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	})

	if err != nil {
		return errors.New("failed to create user")
	}

	return nil
}

// GetByID retrieves a user by ID
func (r *PostgresUserRepository) GetByID(ctx context.Context, id string) (*entity.User, error) {
	uid, err := parseStringToUUID(id)
	if err != nil {
		return nil, errors.New("invalid user ID format")
	}

	user, err := r.queries.GetUserByID(ctx, uid)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("user not found")
		}
		return nil, err
	}

	return toEntity(&user), nil
}

// GetByEmail retrieves a user by email
func (r *PostgresUserRepository) GetByEmail(ctx context.Context, email string) (*entity.User, error) {
	user, err := r.queries.GetUserByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("user not found")
		}
		return nil, err
	}

	return toEntity(&user), nil
}

// Update updates an existing user
func (r *PostgresUserRepository) Update(ctx context.Context, user *entity.User) error {
	uid, err := parseStringToUUID(user.ID)
	if err != nil {
		return errors.New("invalid user ID format")
	}

	err = r.queries.UpdateUser(ctx, sqlc.UpdateUserParams{
		ID:        uid,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Email:     user.Email,
		UpdatedAt: time.Now().UTC(),
	})

	if err != nil {
		return errors.New("failed to update user")
	}

	return nil
}

// Delete soft deletes a user (sets IsActive to false)
func (r *PostgresUserRepository) Delete(ctx context.Context, id string) error {
	uid, err := parseStringToUUID(id)
	if err != nil {
		return errors.New("invalid user ID format")
	}

	err = r.queries.DeleteUser(ctx, sqlc.DeleteUserParams{
		ID:        uid,
		UpdatedAt: time.Now().UTC(),
	})

	if err != nil {
		return errors.New("failed to delete user")
	}

	return nil
}

// List retrieves all users with pagination
func (r *PostgresUserRepository) List(ctx context.Context, limit, offset int) ([]*entity.User, error) {
	users, err := r.queries.ListUsers(ctx, sqlc.ListUsersParams{
		Limit:  int32(limit),
		Offset: int32(offset),
	})

	if err != nil {
		return nil, errors.New("failed to list users")
	}

	result := make([]*entity.User, len(users))
	for i := range users {
		result[i] = toEntity(&users[i])
	}

	return result, nil
}

// toEntity converts sqlc.User to domain entity
func toEntity(user *sqlc.User) *entity.User {
	if user == nil {
		return nil
	}

	return &entity.User{
		ID:        user.ID.String(),
		Email:     user.Email,
		Password:  user.Password,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Role:      user.Role,
		IsActive:  user.IsActive,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}
}
