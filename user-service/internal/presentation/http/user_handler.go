package http

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sandroapkhaidze/Golang-Microservices-Ecommerce/user-service/internal/application/dto"
	"github.com/sandroapkhaidze/Golang-Microservices-Ecommerce/user-service/internal/application/usecase"
)

// UserHandler handles HTTP requests for user operations
type UserHandler struct {
	registerUseCase *usecase.RegisterUserUseCase
	loginUseCase    *usecase.LoginUserUseCase
}

// NewUserHandler creates a new user handler
func NewUserHandler(
	registerUseCase *usecase.RegisterUserUseCase,
	loginUseCase *usecase.LoginUserUseCase,
) *UserHandler {
	return &UserHandler{
		registerUseCase: registerUseCase,
		loginUseCase:    loginUseCase,
	}
}

// Register handles user registration
// @Summary Register a new user
// @Description Create a new user account
// @Tags users
// @Accept json
// @Produce json
// @Param request body dto.RegisterRequest true "Registration details"
// @Success 201 {object} dto.UserResponse
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/v1/users/register [post]
func (h *UserHandler) Register(c *gin.Context) {
	var req dto.RegisterRequest

	// Bind and validate request
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Execute use case
	user, err := h.registerUseCase.Execute(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, user)
}

// Login handles user authentication
// @Summary User login
// @Description Authenticate user and return user data
// @Tags users
// @Accept json
// @Produce json
// @Param request body dto.LoginRequest true "Login credentials"
// @Success 200 {object} dto.UserResponse
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Router /api/v1/users/login [post]
func (h *UserHandler) Login(c *gin.Context) {
	var req dto.LoginRequest

	// Bind and validate request
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Execute use case
	user, err := h.loginUseCase.Execute(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, user)
}

// Health check endpoint
// @Summary Health check
// @Description Check if the service is running
// @Tags health
// @Produce json
// @Success 200 {object} map[string]string
// @Router /health [get]
func (h *UserHandler) Health(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":  "healthy",
		"service": "user-service",
	})
}
