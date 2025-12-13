package usecase

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/sandroapkhaidze/Golang-Microservices-Ecommerce/user-service/internal/application/dto"
	"github.com/sandroapkhaidze/Golang-Microservices-Ecommerce/user-service/internal/domain/entity"
	"github.com/sandroapkhaidze/Golang-Microservices-Ecommerce/user-service/internal/domain/repository"
)

// RegisterUserUseCase handles user registration business logic
type RegisterUserUseCase struct {
	userRepo repository.UserRepository
}

// NewRegisterUserUseCase creates a new RegisterUserUseCase
func NewRegisterUserUseCase(userRepo repository.UserRepository) *RegisterUserUseCase {
	return &RegisterUserUseCase{
		userRepo: userRepo,
	}
}

// Execute registers a new user
func (uc *RegisterUserUseCase) Execute(ctx context.Context, req dto.RegisterRequest) (*dto.UserResponse, error) {
	// Check if user already exists
	existingUser, _ := uc.userRepo.GetByEmail(ctx, req.Email)
	if existingUser != nil {
		return nil, errors.New("user with this email already exists")
	}

	// Create user entity
	user := &entity.User{
		ID:        uuid.New().String(),
		Email:     req.Email,
		Password:  req.Password,
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Role:      "customer", // Default role
		IsActive:  true,
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	}

	// Hash password (domain logic)
	if err := user.HashPassword(); err != nil {
		return nil, errors.New("failed to hash password")
	}

	// Save to database
	if err := uc.userRepo.Create(ctx, user); err != nil {
		return nil, err
	}

	// Convert to DTO response (don't expose password!)
	response := &dto.UserResponse{
		ID:        user.ID,
		Email:     user.Email,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Role:      user.Role,
		IsActive:  user.IsActive,
		CreatedAt: user.CreatedAt,
	}

	return response, nil
}
