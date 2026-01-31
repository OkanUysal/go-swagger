# Basic Example

A simple API with Swagger documentation and auto-detected host.

## Run

```bash
go run main.go
```

## Access Swagger UI

Open your browser: http://localhost:8080/swagger/index.html

## Generate Swagger Docs (Optional)

If you want to use swaggo/swag for code generation:

```bash
# Install swag
go install github.com/swaggo/swag/cmd/swag@latest

# Generate docs
swag init

# This creates docs/ folder, but go-swagger will override the host at runtime
```

## Endpoints

- `GET /api/v1/ping` - Health check
- `GET /api/v1/users` - Get all users (requires Bearer token)
- `POST /api/v1/users` - Create user (requires Bearer token)
