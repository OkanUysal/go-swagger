package swagger

import (
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// SwagSpec represents the generated swagger spec from swag init
var SwagSpec interface{}

// SetupWithSwag configures Swagger UI using swag-generated docs with runtime host detection
// This function expects docs to be imported: import _ "yourproject/docs"
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

// LoadSwagDocs loads the swagger docs from swag-generated package
// Usage: swagger.LoadSwagDocs(docs.SwaggerInfo.ReadDoc())
func LoadSwagDocs(docJSON string) (interface{}, error) {
	var spec interface{}
	if err := json.Unmarshal([]byte(docJSON), &spec); err != nil {
		return nil, err
	}
	return spec, nil
}
