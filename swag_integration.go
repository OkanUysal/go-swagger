package swagger

import (
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// SwagSpec holds the global swagger spec parsed from swag init
var SwagSpec interface{}

// SetupWithSwag configures Swagger UI using swag-generated documentation with runtime host detection.
//
// This is the recommended approach for Railway deployments. It allows you to:
// - Use standard swag annotations (@Summary, @Router, etc.)
// - Auto-detect host from Railway headers/environment
// - Override @host annotation at runtime
//
// Example:
//
//	import _ "myapp/docs" // Import swag-generated docs
//
//	swaggerConfig := swagger.DefaultConfig()
//	swaggerConfig.AutoDetectHost = true
//
//	swagSpec, _ := swagger.LoadSwagDocs(docs.SwaggerInfo.ReadDoc())
//	swagger.SetupWithSwag(router, swagSpec, swaggerConfig)
//
// Parameters:
//   - router: Gin engine instance
//   - swagSpec: Parsed swagger spec from LoadSwagDocs()
//   - config: Configuration options (use DefaultConfig() for defaults)
func SetupWithSwag(router *gin.Engine, swagSpec interface{}, config *Config) {
	if config == nil {
		config = DefaultConfig()
	}

	// Skip if disabled
	if !config.Enabled {
		return
	}

	SwagSpec = swagSpec

	// Serve dynamic swagger.json with auto-detected host
	router.GET(config.JSONPath, func(c *gin.Context) {
		// Parse the swag-generated spec
		specMap, ok := swagSpec.(map[string]interface{})
		if !ok {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid swagger spec"})
			return
		}

		// Clone the spec to avoid modifying the original
		dynamicSpec := make(map[string]interface{})
		for k, v := range specMap {
			dynamicSpec[k] = v
		}

		// Override host and schemes with auto-detection if enabled
		if config.AutoDetectHost {
			dynamicSpec["host"] = detectHost(c)
			dynamicSpec["schemes"] = []string{detectScheme(c)}
		} else if config.Host != "" {
			dynamicSpec["host"] = config.Host
			if len(config.Schemes) > 0 {
				dynamicSpec["schemes"] = config.Schemes
			} else {
				dynamicSpec["schemes"] = []string{detectScheme(c)}
			}
		}

		c.JSON(http.StatusOK, dynamicSpec)
	})

	// Serve Swagger UI
	url := ginSwagger.URL(config.JSONPath)
	router.GET(config.UIPath+"/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, url))
}

// LoadSwagDocs parses swag-generated documentation into a swagger spec.
//
// Example:
//
//	import _ "myapp/docs"
//
//	swagSpec, err := swagger.LoadSwagDocs(docs.SwaggerInfo.ReadDoc())
//	if err != nil {
//	    log.Fatal(err)
//	}
//
// Parameters:
//   - docJSON: JSON string from docs.SwaggerInfo.ReadDoc()
//
// Returns:
//   - interface{}: Parsed swagger spec (map[string]interface{})
//   - error: Parse error if JSON is invalid
func LoadSwagDocs(docJSON string) (interface{}, error) {
	var spec interface{}
	if err := json.Unmarshal([]byte(docJSON), &spec); err != nil {
		return nil, err
	}
	return spec, nil
}

// SetupFromDocs is a convenience function that loads swag-generated docs and sets up Swagger UI.
// This expects that docs are imported with: import _ "yourapp/docs"
//
// Example:
//
//	import (
//		_ "myapp/docs"
//		docs "myapp/docs"
//	)
//
//	swagger.SetupFromDocs(router, docs.SwaggerInfo, nil)
//
// Parameters:
//   - router: Gin engine instance
//   - swaggerInfo: docs.SwaggerInfo from your generated docs package
//   - config: Optional configuration (pass nil for defaults)
func SetupFromDocs(router *gin.Engine, swaggerInfo interface{}, config *Config) error {
	if config == nil {
		config = DefaultConfig()
	}

	// Extract ReadDoc method from SwaggerInfo
	type swaggerInfoInterface interface {
		ReadDoc() string
	}

	info, ok := swaggerInfo.(swaggerInfoInterface)
	if !ok {
		return http.ErrNotSupported
	}

	// Load swagger docs
	swagSpec, err := LoadSwagDocs(info.ReadDoc())
	if err != nil {
		return err
	}

	// Setup with loaded spec
	SetupWithSwag(router, swagSpec, config)
	return nil
}
