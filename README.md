# go-swagger

🚀 Railway-friendly Swagger/OpenAPI documentation for Gin with automatic host detection.

## Why This Library?

The standard `swaggo/swag` requires hardcoded `@host` annotations (e.g., `@host localhost:8080`), which breaks when deploying to Railway or cloud platforms. 

**This library solves that** by auto-detecting the host at runtime from headers and environment variables.

## Features

- ✅ **Auto Host Detection** - Works on Railway, Heroku, any cloud platform
- ✅ **Standard swag annotations** - Use familiar `@Summary`, `@Router`, etc.
- ✅ **Zero configuration** - Just `swag init` and go
- ✅ **Production ready** - Enable/disable based on environment

## Installation

```bash
go get github.com/OkanUysal/go-swagger
go get github.com/swaggo/swag/cmd/swag@latest
```

## Quick Start

### 1. Add annotations to your code

```go
package main

import (
    "github.com/gin-gonic/gin"
    "github.com/OkanUysal/go-swagger"
    _ "yourproject/docs"
)

// @title           My API
// @version         1.0
// @description     API documentation
// @host            localhost:8080
// @BasePath        /api
func main() {
    r := gin.Default()

    // Setup Swagger with auto host detection
    swaggerConfig := swagger.DefaultConfig()
    swaggerConfig.AutoDetectHost = true
    
    swagSpec, _ := swagger.LoadSwagDocs(docs.SwaggerInfo.ReadDoc())
    swagger.SetupWithSwag(r, swagSpec, swaggerConfig)

    r.Run(":8080")
}

// @Summary      Get users
// @Description  Returns list of users
// @Tags         users
// @Accept       json
// @Produce      json
// @Success      200  {array}  User
// @Router       /users [get]
func GetUsers(c *gin.Context) {
    c.JSON(200, []User{})
}
```

### 2. Generate docs

```bash
swag init
```

### 3. Run your app

```bash
go run main.go
```

Visit `http://localhost:8080/swagger/index.html`

## Railway Deployment

**No changes needed!** The `@host localhost:8080` annotation will be automatically updated to your Railway domain at runtime.

```bash
railway up
```

Your Swagger will work at `https://your-app.railway.app/swagger/index.html` with the correct host!

## How It Works

**Host Detection Priority:**
1. `X-Forwarded-Host` header (Railway, Nginx)
2. `RAILWAY_STATIC_URL` environment variable
3. `Host` header from request
4. Fallback to annotation host

**Scheme Detection:**
- Railway domains → `https`
- TLS connection → `https`
- `X-Forwarded-Proto` header
- Default → annotation scheme

## Configuration

### Basic config (recommended)

```go
swaggerConfig := swagger.DefaultConfig()
swaggerConfig.AutoDetectHost = true // Enable auto-detection
```

### Disable in production

```go
import "os"

swaggerConfig := swagger.DefaultConfig()
swaggerConfig.AutoDetectHost = true
swaggerConfig.Enabled = os.Getenv("ENV") != "production"
```

### Custom paths

```go
swaggerConfig := swagger.DefaultConfig()
swaggerConfig.AutoDetectHost = true
swaggerConfig.UIPath = "/docs"         // Swagger UI at /docs/index.html
swaggerConfig.JSONPath = "/docs.json"  // swagger.json at /docs.json
```

## Config Options

```go
type Config struct {
    AutoDetectHost  bool     // Enable auto host detection
    Enabled         bool     // Enable/disable Swagger UI
    UIPath          string   // Swagger UI path (default: "/swagger")
    JSONPath        string   // swagger.json path (default: "/swagger.json")
    Host            string   // Manual host override
    Schemes         []string // Manual schemes override
}
```

## Default Config

```go
swagger.DefaultConfig() returns:
{
    AutoDetectHost: true,
    Enabled:        true,
    UIPath:         "/swagger",
    JSONPath:       "/swagger.json",
}
```

## Complete Example

```go
package main

import (
    "os"
    "github.com/gin-gonic/gin"
    "github.com/OkanUysal/go-swagger"
    _ "github.com/myapp/docs"
)

// @title           My API
// @version         1.0.0
// @description     REST API documentation
// @host            localhost:8080
// @BasePath        /api
// @schemes         http https
func main() {
    r := gin.Default()

    // Swagger setup
    swaggerConfig := swagger.DefaultConfig()
    swaggerConfig.AutoDetectHost = true
    swaggerConfig.Enabled = os.Getenv("ENV") != "production"
    
    swagSpec, err := swagger.LoadSwagDocs(docs.SwaggerInfo.ReadDoc())
    if err == nil {
        swagger.SetupWithSwag(r, swagSpec, swaggerConfig)
    }

    // Routes
    api := r.Group("/api")
    {
        api.GET("/users", GetUsers)
        api.POST("/users", CreateUser)
    }

    r.Run(":8080")
}

// @Summary      Get users
// @Tags         users
// @Success      200  {array}  User
// @Router       /users [get]
func GetUsers(c *gin.Context) {
    c.JSON(200, []User{})
}

// @Summary      Create user
// @Tags         users
// @Param        user  body  User  true  "User object"
// @Success      201   {object}  User
// @Router       /users [post]
func CreateUser(c *gin.Context) {
    var user User
    c.BindJSON(&user)
    c.JSON(201, user)
}
```

## License

MIT
