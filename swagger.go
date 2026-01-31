package swagger

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// Swagger manages the Swagger documentation
type Swagger struct {
	config *Config
	spec   *SwaggerSpec
}

// SwaggerSpec represents the OpenAPI/Swagger specification
type SwaggerSpec struct {
	Swagger             string                        `json:"swagger"`
	Info                Info                          `json:"info"`
	Host                string                        `json:"host"`
	BasePath            string                        `json:"basePath"`
	Schemes             []string                      `json:"schemes"`
	Paths               map[string]interface{}        `json:"paths,omitempty"`
	Definitions         map[string]interface{}        `json:"definitions,omitempty"`
	SecurityDefinitions map[string]SecurityDefinition `json:"securityDefinitions,omitempty"`
}

// Info represents the API information
type Info struct {
	Title       string   `json:"title"`
	Description string   `json:"description"`
	Version     string   `json:"version"`
	Contact     *Contact `json:"contact,omitempty"`
	License     *License `json:"license,omitempty"`
}

// Contact represents the contact information
type Contact struct {
	Name  string `json:"name,omitempty"`
	Email string `json:"email,omitempty"`
	URL   string `json:"url,omitempty"`
}

// License represents the license information
type License struct {
	Name string `json:"name"`
	URL  string `json:"url,omitempty"`
}

// SecurityDefinition represents a security scheme
type SecurityDefinition struct {
	Type string `json:"type"`
	In   string `json:"in,omitempty"`
	Name string `json:"name,omitempty"`
}

// New creates a new Swagger instance
func New(config *Config) *Swagger {
	if config == nil {
		config = DefaultConfig()
	}

	spec := &SwaggerSpec{
		Swagger:  "2.0",
		Host:     config.Host,
		BasePath: config.BasePath,
		Schemes:  config.Schemes,
		Info: Info{
			Title:       config.Title,
			Description: config.Description,
			Version:     config.Version,
		},
		Paths:       make(map[string]interface{}),
		Definitions: make(map[string]interface{}),
	}

	// Add contact if provided
	if config.ContactName != "" || config.ContactEmail != "" || config.ContactURL != "" {
		spec.Info.Contact = &Contact{
			Name:  config.ContactName,
			Email: config.ContactEmail,
			URL:   config.ContactURL,
		}
	}

	// Add license if provided
	if config.LicenseName != "" {
		spec.Info.License = &License{
			Name: config.LicenseName,
			URL:  config.LicenseURL,
		}
	}

	// Add Bearer auth if enabled
	if config.BearerAuth {
		spec.SecurityDefinitions = map[string]SecurityDefinition{
			"Bearer": {
				Type: "apiKey",
				In:   "header",
				Name: "Authorization",
			},
		}
	}

	return &Swagger{
		config: config,
		spec:   spec,
	}
}

// Setup configures Swagger UI routes on a Gin router
func Setup(router *gin.Engine, config *Config) {
	swagger := New(config)

	// Skip if disabled
	if !config.Enabled {
		return
	}

	// Serve dynamic swagger.json
	router.GET(config.JSONPath, swagger.docHandler)

	// Serve Swagger UI
	url := ginSwagger.URL(config.JSONPath)
	router.GET(config.UIPath+"/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, url))
}

// docHandler serves the swagger.json with dynamic host and scheme
func (s *Swagger) docHandler(c *gin.Context) {
	// Update host if auto-detection is enabled
	if s.config.AutoDetectHost {
		s.spec.Host = detectHost(c)
		s.spec.Schemes = []string{detectScheme(c)}
	} else if s.config.Host != "" {
		s.spec.Host = s.config.Host
		if len(s.config.Schemes) > 0 {
			s.spec.Schemes = s.config.Schemes
		} else {
			s.spec.Schemes = []string{detectScheme(c)}
		}
	}

	c.JSON(http.StatusOK, s.spec)
}

// SetPaths sets the API paths (from swag generated docs)
func (s *Swagger) SetPaths(paths map[string]interface{}) {
	s.spec.Paths = paths
}

// SetDefinitions sets the API definitions (from swag generated docs)
func (s *Swagger) SetDefinitions(definitions map[string]interface{}) {
	s.spec.Definitions = definitions
}

// GetSpec returns the current Swagger specification
func (s *Swagger) GetSpec() *SwaggerSpec {
	return s.spec
}

// ExportJSON exports the Swagger spec as JSON string
func (s *Swagger) ExportJSON() (string, error) {
	data, err := json.MarshalIndent(s.spec, "", "  ")
	if err != nil {
		return "", fmt.Errorf("failed to marshal swagger spec: %w", err)
	}
	return string(data), nil
}
