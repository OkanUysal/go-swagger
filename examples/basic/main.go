package main

import (
	"log"

	swagger "github.com/OkanUysal/go-swagger"
	"github.com/gin-gonic/gin"
)

// @title Basic API
// @version 1.0
// @description A simple API with Swagger documentation
// @BasePath /api/v1
// @securityDefinitions.apikey Bearer
// @in header
// @name Authorization
func main() {
	// Create Gin router
	router := gin.Default()

	// Configure Swagger with auto-detection
	swaggerConfig := swagger.NewConfig().
		WithTitle("Basic API").
		WithDescription("A simple API with auto-detected host").
		WithVersion("1.0.0").
		WithBasePath("/api/v1").
		WithAutoDetectHost(true).
		WithBearerAuth(true).
		WithContact("API Team", "api@example.com", "https://example.com")

	swagger.Setup(router, swaggerConfig)

	// API routes
	v1 := router.Group("/api/v1")
	{
		v1.GET("/ping", pingHandler)
		v1.GET("/users", getUsersHandler)
		v1.POST("/users", createUserHandler)
	}

	log.Println("Server starting on :8080")
	log.Println("Swagger UI: http://localhost:8080/swagger/index.html")

	if err := router.Run(":8080"); err != nil {
		log.Fatal(err)
	}
}

// User represents a user model
type User struct {
	ID    int    `json:"id" example:"1"`
	Name  string `json:"name" example:"John Doe"`
	Email string `json:"email" example:"john@example.com"`
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
		"status":  "healthy",
	})
}

// getUsersHandler godoc
// @Summary Get all users
// @Description Get a list of all users
// @Tags users
// @Produce json
// @Security Bearer
// @Success 200 {array} User
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
// @Description Create a new user with the provided data
// @Tags users
// @Accept json
// @Produce json
// @Security Bearer
// @Param user body CreateUserRequest true "User data"
// @Success 201 {object} User
// @Failure 400 {object} map[string]string
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
