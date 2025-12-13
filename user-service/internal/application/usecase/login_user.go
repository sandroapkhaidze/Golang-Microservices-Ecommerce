package usecase

import (
	"context"
	"errors"

	"github.com/sandroapkhaidze/Golang-Microservices-Ecommerce/user-service/internal/application/dto"
	"github.com/sandroapkhaidze/Golang-Microservices-Ecommerce/user-service/internal/domain/repository"
)

// LoginUserUseCase handles user authentication
type LoginUserUseCase struct {
	userRepo repository.UserRepository
}

// NewLoginUserUseCase creates a new LoginUserUseCase
func NewLoginUserUseCase(userRepo repository.UserRepository) *LoginUserUseCase {
	return &LoginUserUseCase{
		userRepo: userRepo,
	}
}

// Execute authenticates a user and returns user data
func (uc *LoginUserUseCase) Execute(ctx context.Context, req dto.LoginRequest) (*dto.UserResponse, error) {
	// Get user by email
	user, err := uc.userRepo.GetByEmail(ctx, req.Email)
	if err != nil {
		return nil, errors.New("invalid email or password")
	}

	// Check if user is active
	if !user.IsActive {
		return nil, errors.New("user account is disabled")
	}

	// Verify password (domain logic)
	if !user.CheckPassword(req.Password) {
		return nil, errors.New("invalid email or password")
	}

	// Convert to DTO (hide sensitive data)
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
