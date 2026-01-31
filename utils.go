package swagger

import (
	"os"
	"strings"

	"github.com/gin-gonic/gin"
)

// detectHost automatically detects the host from the request or environment
func detectHost(c *gin.Context) string {
	// 1. Check X-Forwarded-Host header (reverse proxy, Railway, Nginx)
	if host := c.Request.Header.Get("X-Forwarded-Host"); host != "" {
		return host
	}

	// 2. Check Host header
	if host := c.Request.Host; host != "" {
		// Remove port for cleaner display if it's standard port
		if strings.HasSuffix(host, ":80") || strings.HasSuffix(host, ":443") {
			return strings.Split(host, ":")[0]
		}
		return host
	}

	// 3. Check environment variables
	// Railway specific
	if railwayURL := os.Getenv("RAILWAY_STATIC_URL"); railwayURL != "" {
		return railwayURL
	}

	if railwayURL := os.Getenv("RAILWAY_PUBLIC_DOMAIN"); railwayURL != "" {
		return railwayURL
	}

	// Generic API host
	if apiHost := os.Getenv("API_HOST"); apiHost != "" {
		return apiHost
	}

	if apiURL := os.Getenv("API_URL"); apiURL != "" {
		return apiURL
	}

	// 4. Default fallback
	return "localhost:8080"
}

// detectScheme automatically detects the scheme (http/https) from the request
func detectScheme(c *gin.Context) string {
	// 1. Check if TLS is enabled
	if c.Request.TLS != nil {
		return "https"
	}

	// 2. Check X-Forwarded-Proto header (reverse proxy)
	if proto := c.Request.Header.Get("X-Forwarded-Proto"); proto != "" {
		return proto
	}

	// 3. Check if host contains Railway domain (usually https)
	host := c.Request.Host
	if strings.Contains(host, ".railway.app") || strings.Contains(host, ".up.railway.app") {
		return "https"
	}

	// 4. Check environment
	if os.Getenv("ENV") == "production" || os.Getenv("ENV") == "staging" {
		return "https"
	}

	// 5. Default to http
	return "http"
}

// getEnvWithDefault gets an environment variable with a default value
func getEnvWithDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// isProduction checks if the environment is production
func isProduction() bool {
	env := strings.ToLower(os.Getenv("ENV"))
	return env == "production" || env == "prod"
}

// isLocalhost checks if the host is localhost
func isLocalhost(host string) bool {
	return strings.Contains(host, "localhost") || strings.Contains(host, "127.0.0.1")
}
