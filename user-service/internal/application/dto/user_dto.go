package dto

import "time"

// RegisterRequest represents user registration request
type RegisterRequest struct {
	Email     string `json:"email" binding:"required,email"`
	Password  string `json:"password" binding:"required,min=8"`
	FirstName string `json:"first_name" binding:"required"`
	LastName  string `json:"last_name" binding:"required"`
}

// LoginRequest represents user login request
type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

// LoginResponse represents login response with JWT token
type LoginResponse struct {
	Token     string       `json:"token"`
	User      UserResponse `json:"user"`
	ExpiresAt time.Time    `json:"expires_at"`
}

// UserResponse represents user data in API responses
type UserResponse struct {
	ID        string    `json:"id"`
	Email     string    `json:"email"`
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
	Role      string    `json:"role"`
	IsActive  bool      `json:"is_active"`
	CreatedAt time.Time `json:"created_at"`
}

// UpdateUserRequest represents user update request
type UpdateUserRequest struct {
	FirstName *string `json:"first_name,omitempty"`
	LastName  *string `json:"last_name,omitempty"`
	Email     *string `json:"email,omitempty" binding:"omitempty,email"`
}

// ListUsersResponse represents paginated user list
type ListUsersResponse struct {
	Users  []UserResponse `json:"users"`
	Total  int            `json:"total"`
	Limit  int            `json:"limit"`
	Offset int            `json:"offset"`
}
