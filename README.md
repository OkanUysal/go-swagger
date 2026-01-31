# go-swagger

A lightweight, Railway-friendly Swagger/OpenAPI documentation library for Gin framework with automatic host detection.

## Features

- 🚀 **Railway-Ready**: Automatic host and scheme detection for Railway deployments
- 🔄 **Dynamic Configuration**: Runtime host detection from headers and environment
- 🛡️ **Production Safe**: Enable/disable Swagger UI based on environment
- 🔐 **Bearer Auth**: Built-in JWT/Bearer token authentication support
- ⚡ **Zero Boilerplate**: Simple API with sensible defaults
- 🎯 **Gin Integration**: Native support for Gin web framework

## Problem It Solves

The default `swaggo/swag` library requires hardcoded `@host` annotations (e.g., `@host localhost:8080`), which breaks when deploying to Railway or other cloud platforms. This library solves that by:

1. **Auto-detecting the actual host** from `X-Forwarded-Host`, `Host` header, or `RAILWAY_STATIC_URL`
2. **Auto-detecting the scheme** (http/https) from request headers or TLS
3. **Generating swagger.json dynamically** at runtime instead of static build-time generation

## Installation

```bash
go get github.com/OkanUysal/go-swagger
```

## Dependencies

```bash
go get github.com/gin-gonic/gin@v1.10.0
go get github.com/swaggo/swag@v1.16.3
go get github.com/swaggo/gin-swagger@v1.6.0
go get github.com/swaggo/files@v1.0.1
```

Or add to your `go.mod`:

```go
require (
    github.com/gin-gonic/gin v1.10.0
    github.com/swaggo/swag v1.16.3
    github.com/swaggo/gin-swagger v1.6.0
    github.com/swaggo/files v1.0.1
)
```

## Quick Start

### Option 1: With swag annotations (Recommended)

**1. Add Swagger annotations to your code:**

```go
package main

import (
    "github.com/gin-gonic/gin"
    "github.com/OkanUysal/go-swagger"
    _ "yourproject/docs" // Import generated docs
)

// @title           Your API
// @version         1.0
// @description     API documentation
// @BasePath        /api

func main() {
    r := gin.Default()

    // Setup Swagger with auto host detection
    swaggerConfig := swagger.DefaultConfig()
    swaggerConfig.AutoDetectHost = true // Automatically detect host from request
    swagger.SetupWithSwag(r, docs.SwaggerInfo, swaggerConfig)

    r.Run(":8080")
}

// @Summary      Get user
// @Description  Get user by ID
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "User ID"
// @Success      200  {object}  User
// @Router       /users/{id} [get]
func GetUser(c *gin.Context) {
    // handler code
}
```

**2. Generate docs:**
```bash
swag init
```

**3. Run your app - Host will be auto-detected!**

### Option 2: Programmatic (Without swag)

**Basic Usage**

```go
package main

import (
    "github.com/gin-gonic/gin"
    swagger "github.com/OkanUysal/go-swagger"
)

func main() {
    router := gin.Default()

    // Configure Swagger with auto-detection
    config := swagger.NewConfig().
        WithTitle("My API").
        WithDescription("API Documentation").
        WithVersion("1.0.0").
        WithAutoDetectHost(true).
        WithBearerAuth(true)

    swagger.Setup(router, config)

    router.GET("/ping", func(c *gin.Context) {
        c.JSON(200, gin.H{"message": "pong"})
    })

    router.Run(":8080")
}
```

Now visit `http://localhost:8080/swagger/index.html` to see your Swagger UI.

### Railway Deployment

For Railway deployment, the library automatically detects the host:

```go
config := swagger.NewConfig().
    WithTitle("My Railway API").
    WithAutoDetectHost(true). // Detects from RAILWAY_STATIC_URL or headers
    WithBearerAuth(true)

swagger.Setup(router, config)
```

**Host Detection Priority:**
1. `X-Forwarded-Host` header (from Railway/Nginx reverse proxy)
2. `Host` header (from request)
3. `RAILWAY_STATIC_URL` environment variable
4. `RAILWAY_PUBLIC_DOMAIN` environment variable
5. `API_HOST` environment variable
6. `API_URL` environment variable
7. Default: `localhost:8080`

**Scheme Detection:**
- Checks TLS status
- Checks `X-Forwarded-Proto` header
- Railway domains (`.railway.app`) default to HTTPS
- Production environment defaults to HTTPS

### Environment-Based Configuration

Disable Swagger in production:

```go
import "os"

enabled := os.Getenv("ENV") != "production"

config := swagger.NewConfig().
    WithTitle("My API").
    WithEnabled(enabled). // Disable in production
    WithAutoDetectHost(true)

swagger.Setup(router, config)
```

### Manual Host Configuration

If you prefer manual configuration:

```go
config := swagger.NewConfig().
    WithTitle("My API").
    WithHost("api.example.com").
    WithSchemes([]string{"https"}).
    WithBasePath("/v1")

swagger.Setup(router, config)
```

## Configuration Options

```go
config := swagger.NewConfig().
    WithTitle("My API").                                    // API title
    WithDescription("API Documentation").                   // API description
    WithVersion("1.0.0").                                   // API version
    WithAutoDetectHost(true).                               // Enable auto host detection
    WithHost("api.example.com").                            // Manual host (if not auto-detecting)
    WithSchemes([]string{"https"}).                         // Schemes: http, https
    WithBasePath("/v1").                                    // Base path for all routes
    WithUIPath("/swagger").                                 // Swagger UI path
    WithJSONPath("/swagger/doc.json").                      // swagger.json path
    WithBearerAuth(true).                                   // Enable Bearer token auth
    WithEnabled(true).                                      // Enable/disable Swagger UI
    WithContact("API Team", "api@example.com", "https://example.com"). // Contact info
    WithLicense("MIT", "https://opensource.org/licenses/MIT")          // License info
```

## Configuration Struct

```go
type Config struct {
    Title            string   // API title
    Description      string   // API description
    Version          string   // API version
    Host             string   // API host (e.g., "api.example.com")
    BasePath         string   // Base path (e.g., "/v1")
    Schemes          []string // Schemes: ["http", "https"]
    AutoDetectHost   bool     // Auto-detect host from request/env
    Enabled          bool     // Enable/disable Swagger UI
    UIPath           string   // Swagger UI path (default: "/swagger")
    JSONPath         string   // swagger.json path (default: "/swagger/doc.json")
    BearerAuth       bool     // Enable Bearer token authentication
    ContactName      string   // Contact name
    ContactEmail     string   // Contact email
    ContactURL       string   // Contact URL
    LicenseName      string   // License name
    LicenseURL       string   // License URL
}
```

## Advanced Usage

### Custom Swagger Spec

```go
swagger := swagger.New(config)

// Add custom paths
paths := map[string]interface{}{
    "/users": map[string]interface{}{
        "get": map[string]interface{}{
            "summary":     "Get users",
            "description": "Returns all users",
            "responses": map[string]interface{}{
                "200": map[string]interface{}{
                    "description": "Success",
                },
            },
        },
    },
}
swagger.SetPaths(paths)

// Add custom definitions
definitions := map[string]interface{}{
    "User": map[string]interface{}{
        "type": "object",
        "properties": map[string]interface{}{
            "id":   map[string]interface{}{"type": "integer"},
            "name": map[string]interface{}{"type": "string"},
        },
    },
}
swagger.SetDefinitions(definitions)

// Export to JSON
json, err := swagger.ExportJSON()
if err != nil {
    panic(err)
}
fmt.Println(json)
```

### Integration with swaggo/swag

You can still use `swaggo/swag` for code generation and combine it with this library for runtime host detection:

1. Install swag CLI:
```bash
go install github.com/swaggo/swag/cmd/swag@latest
```

2. Add swag comments to your code:
```go
// @title My API
// @version 1.0
// @description API Documentation
// @host localhost:8080  // This will be overridden at runtime
// @BasePath /v1
// @securityDefinitions.apikey Bearer
// @in header
// @name Authorization
func main() {
    // ...
}
```

3. Generate docs:
```bash
swag init
```

4. Use this library for runtime host override:
```go
import _ "myproject/docs" // Import generated docs

config := swagger.NewConfig().
    WithAutoDetectHost(true). // Override @host at runtime
    WithBearerAuth(true)

swagger.Setup(router, config)
```

## Railway Environment Variables

The library automatically reads these Railway variables:

- `RAILWAY_STATIC_URL`: Your Railway deployment URL
- `RAILWAY_PUBLIC_DOMAIN`: Public domain for your app
- `ENV`: Environment (production, staging, development)

Example Railway deployment:
```bash
# Railway automatically sets these
RAILWAY_STATIC_URL=myapi.up.railway.app
ENV=production
```

Your Swagger docs will automatically use `https://myapi.up.railway.app`.

## Examples

See the `examples/` directory for complete examples:

- `examples/basic/`: Simple API with auto-detection
- `examples/railway/`: Railway-optimized setup
- `examples/production/`: Production-ready configuration

## Dependencies

- [gin-gonic/gin](https://github.com/gin-gonic/gin) - Web framework
- [swaggo/gin-swagger](https://github.com/swaggo/gin-swagger) - Gin middleware for Swagger
- [swaggo/files](https://github.com/swaggo/files) - Swagger UI static files

## Testing

```bash
go test -v ./...
```

## License

MIT License - see LICENSE file for details

## Contributing

Contributions are welcome! Please open an issue or submit a pull request.

## Author

Okan Uysal - [GitHub](https://github.com/OkanUysal)
