package main

import (
	"log"
	"os"

	swagger "github.com/OkanUysal/go-swagger"
	"github.com/gin-gonic/gin"
)

// @title Railway API
// @version 1.0
// @description API deployed on Railway with auto-detected host
// @BasePath /api/v1
// @securityDefinitions.apikey Bearer
// @in header
// @name Authorization
func main() {
	// Set Gin mode based on environment
	env := os.Getenv("ENV")
	if env == "production" || env == "staging" {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.Default()

	// Railway-optimized Swagger configuration
	swaggerEnabled := env != "production" // Disable in production

	swaggerConfig := swagger.NewConfig().
		WithTitle("Railway API").
		WithDescription("API deployed on Railway with automatic host detection").
		WithVersion("1.0.0").
		WithBasePath("/api/v1").
		WithAutoDetectHost(true). // Automatically detects Railway URL
		WithEnabled(swaggerEnabled).
		WithBearerAuth(true).
		WithContact("API Team", "api@example.com", "https://example.com")

	swagger.Setup(router, swaggerConfig)

	// Health check (no auth required)
	router.GET("/health", healthHandler)

	// API routes
	v1 := router.Group("/api/v1")
	{
		v1.GET("/ping", pingHandler)

		// Protected routes
		protected := v1.Group("")
		// protected.Use(authMiddleware()) // Add your auth middleware here
		{
			protected.GET("/users", getUsersHandler)
			protected.POST("/users", createUserHandler)
			protected.GET("/profile", getProfileHandler)
		}
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Starting server on port %s", port)
	log.Printf("Environment: %s", env)
	log.Printf("Swagger enabled: %v", swaggerEnabled)

	if swaggerEnabled {
		railwayURL := os.Getenv("RAILWAY_STATIC_URL")
		if railwayURL != "" {
			log.Printf("Swagger UI: https://%s/swagger/index.html", railwayURL)
		} else {
			log.Printf("Swagger UI: http://localhost:%s/swagger/index.html", port)
		}
	}

	if err := router.Run(":" + port); err != nil {
		log.Fatal(err)
	}
}

// healthHandler godoc
// @Summary Health check
// @Description Check if the service is healthy
// @Tags health
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Router /health [get]
func healthHandler(c *gin.Context) {
	c.JSON(200, gin.H{
		"status": "healthy",
		"env":    os.Getenv("ENV"),
	})
}

// pingHandler godoc
// @Summary Ping endpoint
// @Description Check if the API is alive
// @Tags health
// @Produce json
// @Success 200 {object} map[string]string
// @Router /ping [get]
func pingHandler(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "pong",
		"host":    c.Request.Host,
		"scheme":  getScheme(c),
	})
}

// User represents a user model
type User struct {
	ID    int    `json:"id" example:"1"`
	Name  string `json:"name" example:"John Doe"`
	Email string `json:"email" example:"john@example.com"`
}

// getUsersHandler godoc
// @Summary Get all users
// @Description Get a list of all users (protected endpoint)
// @Tags users
// @Produce json
// @Security Bearer
// @Success 200 {array} User
// @Failure 401 {object} map[string]string
// @Router /users [get]
func getUsersHandler(c *gin.Context) {
	users := []User{
		{ID: 1, Name: "John Doe", Email: "john@example.com"},
		{ID: 2, Name: "Jane Smith", Email: "jane@example.com"},
	}
	c.JSON(200, users)
}

// CreateUserRequest represents the request body for creating a user
type CreateUserRequest struct {
	Name  string `json:"name" binding:"required" example:"John Doe"`
	Email string `json:"email" binding:"required" example:"john@example.com"`
}

// createUserHandler godoc
// @Summary Create a new user
// @Description Create a new user with the provided data (protected endpoint)
// @Tags users
// @Accept json
// @Produce json
// @Security Bearer
// @Param user body CreateUserRequest true "User data"
// @Success 201 {object} User
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Router /users [post]
func createUserHandler(c *gin.Context) {
	var req CreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	user := User{
		ID:    3,
		Name:  req.Name,
		Email: req.Email,
	}
	c.JSON(201, user)
}

// ProfileResponse represents the user profile response
type ProfileResponse struct {
	ID       int      `json:"id" example:"1"`
	Name     string   `json:"name" example:"John Doe"`
	Email    string   `json:"email" example:"john@example.com"`
	Roles    []string `json:"roles" example:"admin,user"`
	Verified bool     `json:"verified" example:"true"`
}

// getProfileHandler godoc
// @Summary Get user profile
// @Description Get the authenticated user's profile (protected endpoint)
// @Tags users
// @Produce json
// @Security Bearer
// @Success 200 {object} ProfileResponse
// @Failure 401 {object} map[string]string
// @Router /profile [get]
func getProfileHandler(c *gin.Context) {
	profile := ProfileResponse{
		ID:       1,
		Name:     "John Doe",
		Email:    "john@example.com",
		Roles:    []string{"admin", "user"},
		Verified: true,
	}
	c.JSON(200, profile)
}

func getScheme(c *gin.Context) string {
	if c.Request.TLS != nil {
		return "https"
	}
	if proto := c.Request.Header.Get("X-Forwarded-Proto"); proto != "" {
		return proto
	}
	return "http"
}
