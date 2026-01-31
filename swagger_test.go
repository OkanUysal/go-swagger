package swagger

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	t.Run("with custom config", func(t *testing.T) {
		config := NewConfig().
			WithTitle("Test API").
			WithDescription("Test Description").
			WithVersion("1.0.0").
			WithBearerAuth(true)

		swagger := New(config)

		assert.NotNil(t, swagger)
		assert.Equal(t, "Test API", swagger.spec.Info.Title)
		assert.Equal(t, "Test Description", swagger.spec.Info.Description)
		assert.Equal(t, "1.0.0", swagger.spec.Info.Version)
		assert.True(t, swagger.config.BearerAuth)
		assert.Contains(t, swagger.spec.SecurityDefinitions, "Bearer")
	})

	t.Run("with nil config uses defaults", func(t *testing.T) {
		swagger := New(nil)

		assert.NotNil(t, swagger)
		assert.Equal(t, "API", swagger.spec.Info.Title)
		assert.Equal(t, "1.0.0", swagger.spec.Info.Version)
	})

	t.Run("with contact info", func(t *testing.T) {
		config := NewConfig().
			WithContact("John Doe", "john@example.com", "https://example.com")

		swagger := New(config)

		assert.NotNil(t, swagger.spec.Info.Contact)
		assert.Equal(t, "John Doe", swagger.spec.Info.Contact.Name)
		assert.Equal(t, "john@example.com", swagger.spec.Info.Contact.Email)
	})

	t.Run("with license info", func(t *testing.T) {
		config := NewConfig().
			WithLicense("MIT", "https://opensource.org/licenses/MIT")

		swagger := New(config)

		assert.NotNil(t, swagger.spec.Info.License)
		assert.Equal(t, "MIT", swagger.spec.Info.License.Name)
		assert.Equal(t, "https://opensource.org/licenses/MIT", swagger.spec.Info.License.URL)
	})
}

func TestSetPathsAndDefinitions(t *testing.T) {
	swagger := New(DefaultConfig())

	paths := map[string]interface{}{
		"/users": map[string]interface{}{
			"get": map[string]interface{}{
				"summary": "Get users",
			},
		},
	}

	definitions := map[string]interface{}{
		"User": map[string]interface{}{
			"type": "object",
			"properties": map[string]interface{}{
				"id": map[string]interface{}{
					"type": "integer",
				},
			},
		},
	}

	swagger.SetPaths(paths)
	swagger.SetDefinitions(definitions)

	assert.Equal(t, paths, swagger.spec.Paths)
	assert.Equal(t, definitions, swagger.spec.Definitions)
}

func TestExportJSON(t *testing.T) {
	config := NewConfig().
		WithTitle("Test API").
		WithVersion("1.0.0").
		WithHost("api.example.com")

	swagger := New(config)

	json, err := swagger.ExportJSON()

	assert.NoError(t, err)
	assert.Contains(t, json, "Test API")
	assert.Contains(t, json, "1.0.0")
	assert.Contains(t, json, "api.example.com")
}

func TestDetectHost(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("from X-Forwarded-Host", func(t *testing.T) {
		c, _ := gin.CreateTestContext(nil)
		c.Request = &http.Request{
			Header: http.Header{
				"X-Forwarded-Host": {"api.railway.app"},
			},
		}

		host := detectHost(c)
		assert.Equal(t, "api.railway.app", host)
	})

	t.Run("from Host header", func(t *testing.T) {
		c, _ := gin.CreateTestContext(nil)
		c.Request = &http.Request{
			Host:   "api.example.com",
			Header: http.Header{},
		}

		host := detectHost(c)
		assert.Equal(t, "api.example.com", host)
	})

	t.Run("strips standard ports", func(t *testing.T) {
		c, _ := gin.CreateTestContext(nil)
		c.Request = &http.Request{
			Host:   "api.example.com:443",
			Header: http.Header{},
		}

		host := detectHost(c)
		assert.Equal(t, "api.example.com", host)
	})

	t.Run("default fallback", func(t *testing.T) {
		c, _ := gin.CreateTestContext(nil)
		c.Request = &http.Request{
			Header: http.Header{},
		}

		host := detectHost(c)
		assert.Equal(t, "localhost:8080", host)
	})
}

func TestDetectScheme(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("from X-Forwarded-Proto", func(t *testing.T) {
		c, _ := gin.CreateTestContext(nil)
		c.Request = &http.Request{
			Header: http.Header{
				"X-Forwarded-Proto": {"https"},
			},
		}

		scheme := detectScheme(c)
		assert.Equal(t, "https", scheme)
	})

	t.Run("railway domain defaults to https", func(t *testing.T) {
		c, _ := gin.CreateTestContext(nil)
		c.Request = &http.Request{
			Host:   "api.railway.app",
			Header: http.Header{},
		}

		scheme := detectScheme(c)
		assert.Equal(t, "https", scheme)
	})

	t.Run("default to http", func(t *testing.T) {
		c, _ := gin.CreateTestContext(nil)
		c.Request = &http.Request{
			Host:   "localhost:8080",
			Header: http.Header{},
		}

		scheme := detectScheme(c)
		assert.Equal(t, "http", scheme)
	})
}

func TestIsLocalhost(t *testing.T) {
	assert.True(t, isLocalhost("localhost:8080"))
	assert.True(t, isLocalhost("127.0.0.1:8080"))
	assert.False(t, isLocalhost("api.example.com"))
	assert.False(t, isLocalhost("api.railway.app"))
}

func ExampleSetup() {
	router := gin.New()

	config := NewConfig().
		WithTitle("My API").
		WithDescription("API Documentation").
		WithVersion("1.0.0").
		WithAutoDetectHost(true).
		WithBearerAuth(true)

	Setup(router, config)

	fmt.Println("Swagger UI available at /swagger/index.html")
	// Output: Swagger UI available at /swagger/index.html
}

func ExampleNew() {
	config := NewConfig().
		WithTitle("E-commerce API").
		WithDescription("RESTful API for e-commerce platform").
		WithVersion("2.0.0").
		WithHost("api.example.com").
		WithSchemes([]string{"https"}).
		WithBasePath("/v2").
		WithBearerAuth(true).
		WithContact("API Team", "api@example.com", "https://example.com/contact")

	swagger := New(config)

	json, _ := swagger.ExportJSON()
	fmt.Println("Swagger spec generated")
	_ = json
	// Output: Swagger spec generated
}
