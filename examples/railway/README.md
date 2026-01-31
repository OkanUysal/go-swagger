# Railway Example

API optimized for Railway deployment with automatic host detection.

## Features

- Auto-detects Railway URL from `RAILWAY_STATIC_URL`
- Disables Swagger UI in production
- Environment-based configuration
- Bearer token authentication

## Local Development

```bash
go run main.go
```

Access: http://localhost:8080/swagger/index.html

## Railway Deployment

1. Create a new project on Railway
2. Connect your GitHub repository
3. Set environment variables:
   - `ENV=production` or `ENV=staging`
   - Railway automatically provides `RAILWAY_STATIC_URL` and `PORT`

4. Deploy!

Your API will be available at: `https://your-app.up.railway.app`

## Environment Variables

- `ENV`: Environment (development, staging, production)
- `PORT`: Port to listen on (Railway sets this automatically)
- `RAILWAY_STATIC_URL`: Your Railway domain (set automatically)
- `RAILWAY_PUBLIC_DOMAIN`: Alternative Railway domain

## Endpoints

- `GET /health` - Health check (public)
- `GET /api/v1/ping` - Ping endpoint (public)
- `GET /api/v1/users` - Get users (requires Bearer token)
- `POST /api/v1/users` - Create user (requires Bearer token)
- `GET /api/v1/profile` - Get profile (requires Bearer token)

## Swagger UI Access

- **Development**: http://localhost:8080/swagger/index.html
- **Staging**: Enabled at https://your-staging-url.up.railway.app/swagger/index.html
- **Production**: Disabled (set `ENV=production`)

## Testing Bearer Auth

In Swagger UI, click "Authorize" and enter:

```
Bearer your-jwt-token-here
```
