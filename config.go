package swagger

// Config holds the Swagger configuration
type Config struct {
	// Title is the API title
	Title string

	// Description is the API description
	Description string

	// Version is the API version
	Version string

	// Host is the API host (e.g., "api.example.com")
	// If empty, will be auto-detected from request
	Host string

	// BasePath is the API base path (e.g., "/api/v1")
	BasePath string

	// Schemes are the API schemes (e.g., ["http", "https"])
	// If empty, will be auto-detected from request
	Schemes []string

	// AutoDetectHost enables automatic host detection from request
	// Default: true
	AutoDetectHost bool

	// Enabled controls whether Swagger UI is enabled
	// Set to false in production
	// Default: true
	Enabled bool

	// UIPath is the path to serve Swagger UI
	// Default: "/swagger"
	UIPath string

	// JSONPath is the path to serve swagger.json
	// Default: "/swagger/doc.json"
	JSONPath string

	// BearerAuth enables JWT Bearer authentication in Swagger
	// Default: false
	BearerAuth bool

	// Contact information
	ContactName  string
	ContactEmail string
	ContactURL   string

	// License information
	LicenseName string
	LicenseURL  string
}

// NewConfig creates a new Config with sensible defaults
func NewConfig() *Config {
	return &Config{
		Title:          "API",
		Description:    "API Documentation",
		Version:        "1.0.0",
		BasePath:       "/",
		AutoDetectHost: true, // Auto-detect by default for Railway
		Enabled:        true,
		UIPath:         "/swagger",
		JSONPath:       "/swagger.json",
		BearerAuth:     false,
	}
}

// DefaultConfig returns a Config with sensible defaults (alias for NewConfig)
func DefaultConfig() *Config {
	return &Config{
		AutoDetectHost: true, // Auto-detect by default for Railway
		Enabled:        true,
		UIPath:         "/swagger",
		JSONPath:       "/swagger.json",
	}
}

// WithTitle sets the API title
func (c *Config) WithTitle(title string) *Config {
	c.Title = title
	return c
}

// WithDescription sets the API description
func (c *Config) WithDescription(description string) *Config {
	c.Description = description
	return c
}

// WithVersion sets the API version
func (c *Config) WithVersion(version string) *Config {
	c.Version = version
	return c
}

// WithHost sets the API host
func (c *Config) WithHost(host string) *Config {
	c.Host = host
	return c
}

// WithBasePath sets the API base path
func (c *Config) WithBasePath(basePath string) *Config {
	c.BasePath = basePath
	return c
}

// WithBearerAuth enables Bearer authentication
func (c *Config) WithBearerAuth(enabled bool) *Config {
	c.BearerAuth = enabled
	return c
}

// WithSchemes sets the API schemes
func (c *Config) WithSchemes(schemes []string) *Config {
	c.Schemes = schemes
	return c
}

// WithContact sets the contact information
func (c *Config) WithContact(name, email, url string) *Config {
	c.ContactName = name
	c.ContactEmail = email
	c.ContactURL = url
	return c
}

// WithLicense sets the license information
func (c *Config) WithLicense(name, url string) *Config {
	c.LicenseName = name
	c.LicenseURL = url
	return c
}

// WithEnabled sets whether Swagger is enabled
func (c *Config) WithEnabled(enabled bool) *Config {
	c.Enabled = enabled
	return c
}

// WithAutoDetectHost enables auto-detection of host
func (c *Config) WithAutoDetectHost(enabled bool) *Config {
	c.AutoDetectHost = enabled
	return c
}

// WithUIPath sets the Swagger UI path
func (c *Config) WithUIPath(path string) *Config {
	c.UIPath = path
	return c
}

// WithJSONPath sets the swagger.json path
func (c *Config) WithJSONPath(path string) *Config {
	c.JSONPath = path
	return c
}
