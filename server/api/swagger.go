package api

import (
	"net/http"
	"reflect"
	"strings"

	restful "github.com/emicklei/go-restful/v3"
	restfulspec "github.com/emicklei/go-restful-openapi/v2"
	"github.com/go-openapi/spec"
)

// SwaggerConfig holds the swagger configuration.
type SwaggerConfig struct {
	Title       string
	Description string
	Version     string
	Host        string
	BasePath    string
}

// DefaultSwaggerConfig returns the default swagger configuration.
func DefaultSwaggerConfig() SwaggerConfig {
	return SwaggerConfig{
		Title:       "vBlog Core API",
		Description: "vBlog Core - A customizable, lightweight blog system for geeks and vibe coders",
		Version:     "1.0.0",
		Host:        "localhost:8080",
		BasePath:    "/",
	}
}

// RegisterSwagger registers the swagger spec and UI endpoint.
func RegisterSwagger(container *restful.Container, cfg SwaggerConfig) {
	// Configure OpenAPI spec generation
	config := restfulspec.Config{
		WebServices:                   container.RegisteredWebServices(),
		APIPath:                       "/apidocs.json",
		PostBuildSwaggerObjectHandler: enrichSwaggerObject(cfg),
		// Use TypeOnly naming to strip package prefix from definition names
		// e.g., "model.Component" becomes "Component"
		// Keep time.Time and time.Duration unchanged so they are recognized as primitives
		ModelTypeNameHandler: func(t reflect.Type) (string, bool) {
			name := t.Name()
			if name == "" {
				return "", false
			}
			// Keep standard library time types as-is for proper primitive mapping
			fullName := t.String()
			if fullName == "time.Time" || fullName == "*time.Time" ||
				fullName == "time.Duration" || fullName == "*time.Duration" {
				return "", false
			}
			return name, true
		},
		// Map custom types to swagger formats
		SchemaFormatHandler: func(typeName string) string {
			switch typeName {
			case "gorm.DeletedAt", "*gorm.DeletedAt":
				return "date-time"
			}
			return ""
		},
	}

	// Register the spec endpoint
	container.Add(restfulspec.NewOpenAPIService(config))

	// Register swagger UI at /swagger
	RegisterSwaggerUI(container, cfg)
}

// enrichSwaggerObject adds metadata to the swagger spec.
func enrichSwaggerObject(cfg SwaggerConfig) func(*spec.Swagger) {
	return func(s *spec.Swagger) {
		s.Info = &spec.Info{
			InfoProps: spec.InfoProps{
				Title:       cfg.Title,
				Description: cfg.Description,
				Version:     cfg.Version,
				Contact: &spec.ContactInfo{
					ContactInfoProps: spec.ContactInfoProps{
						Name:  "vBlog Core",
						Email: "admin@vblog.com",
					},
				},
			},
		}
		s.Host = cfg.Host
		s.BasePath = cfg.BasePath
		s.Schemes = []string{"http", "https"}
		s.Tags = []spec.Tag{
			{TagProps: spec.TagProps{Name: "posts", Description: "Blog post management"}},
			{TagProps: spec.TagProps{Name: "tags", Description: "Tag management"}},
			{TagProps: spec.TagProps{Name: "comments", Description: "Comment management"}},
			{TagProps: spec.TagProps{Name: "auth", Description: "Authentication"}},
			{TagProps: spec.TagProps{Name: "settings", Description: "Site settings"}},
			{TagProps: spec.TagProps{Name: "components", Description: "Custom components"}},
			{TagProps: spec.TagProps{Name: "dashboard", Description: "Dashboard statistics"}},
			{TagProps: spec.TagProps{Name: "upload", Description: "File upload"}},
			{TagProps: spec.TagProps{Name: "rss", Description: "RSS feed"}},
		}

		// Add security definition for JWT
		s.SecurityDefinitions = map[string]*spec.SecurityScheme{
			"Bearer": {
				SecuritySchemeProps: spec.SecuritySchemeProps{
					Type:        "apiKey",
					Name:        "Authorization",
					In:          "header",
					Description: "JWT token. Format: Bearer {token}",
				},
			},
		}

		// Fix definition names: rename package-prefixed keys to simple names
		// e.g., "model.Component" -> "Component", "api.ErrorResponse" -> "ErrorResponse"
		fixedDefs := spec.Definitions{}
		for key, schema := range s.Definitions {
			// Strip package prefix if present
			if idx := strings.LastIndex(key, "."); idx != -1 {
				key = key[idx+1:]
			}
			fixedDefs[key] = schema
		}

		// Add common response definitions
		fixedDefs["ErrorResponse"] = spec.Schema{
			SchemaProps: spec.SchemaProps{
				Type: spec.StringOrArray{"object"},
				Properties: map[string]spec.Schema{
					"error": {
						SchemaProps: spec.SchemaProps{
							Type:        spec.StringOrArray{"string"},
							Description: "Error message",
						},
					},
				},
			},
		}
		fixedDefs["MessageResponse"] = spec.Schema{
			SchemaProps: spec.SchemaProps{
				Type: spec.StringOrArray{"object"},
				Properties: map[string]spec.Schema{
					"message": {
						SchemaProps: spec.SchemaProps{
							Type:        spec.StringOrArray{"string"},
							Description: "Success message",
						},
					},
				},
			},
		}
		fixedDefs["TokenResponse"] = spec.Schema{
			SchemaProps: spec.SchemaProps{
				Type: spec.StringOrArray{"object"},
				Properties: map[string]spec.Schema{
					"access_token": {
						SchemaProps: spec.SchemaProps{
							Type:        spec.StringOrArray{"string"},
							Description: "JWT access token",
						},
					},
					"refresh_token": {
						SchemaProps: spec.SchemaProps{
							Type:        spec.StringOrArray{"string"},
							Description: "JWT refresh token",
						},
					},
				},
			},
		}
		fixedDefs["DashboardStatsResponse"] = spec.Schema{
			SchemaProps: spec.SchemaProps{
				Type: spec.StringOrArray{"object"},
				Properties: map[string]spec.Schema{
					"total_posts": {
						SchemaProps: spec.SchemaProps{
							Type:        spec.StringOrArray{"integer"},
							Description: "Total number of published posts",
						},
					},
					"total_views": {
						SchemaProps: spec.SchemaProps{
							Type:        spec.StringOrArray{"integer"},
							Description: "Total post views",
						},
					},
					"total_comments": {
						SchemaProps: spec.SchemaProps{
							Type:        spec.StringOrArray{"integer"},
							Description: "Total number of comments",
						},
					},
					"total_tags": {
						SchemaProps: spec.SchemaProps{
							Type:        spec.StringOrArray{"integer"},
							Description: "Total number of tags",
						},
					},
				},
			},
		}
		fixedDefs["UploadResponse"] = spec.Schema{
			SchemaProps: spec.SchemaProps{
				Type: spec.StringOrArray{"object"},
				Properties: map[string]spec.Schema{
					"url": {
						SchemaProps: spec.SchemaProps{
							Type:        spec.StringOrArray{"string"},
							Description: "Uploaded file URL",
						},
					},
				},
			},
		}

		// Fix time.Time types that weren't recognized as primitives by the builder
		fixTimeTypes(fixedDefs)

		s.Definitions = fixedDefs

		// Fix all $ref pointers in paths to use simple names
		for _, path := range s.Paths.Paths {
			fixOperationRefs(path.Get)
			fixOperationRefs(path.Post)
			fixOperationRefs(path.Put)
			fixOperationRefs(path.Patch)
			fixOperationRefs(path.Delete)
		}
	}
}

// fixOperationRefs fixes $ref pointers in an operation to use simple definition names.
func fixOperationRefs(op *spec.Operation) {
	if op == nil {
		return
	}
	// Fix parameters
	for i := range op.Parameters {
		if op.Parameters[i].Schema != nil {
			fixSchemaRef(op.Parameters[i].Schema)
		}
	}
	// Fix responses
	if op.Responses != nil {
		for code := range op.Responses.StatusCodeResponses {
			resp := op.Responses.StatusCodeResponses[code]
			if resp.Schema != nil {
				fixSchemaRef(resp.Schema)
			}
		}
		if op.Responses.Default != nil && op.Responses.Default.Schema != nil {
			fixSchemaRef(op.Responses.Default.Schema)
		}
	}
}

// fixTimeTypes recursively converts time.Time types from "Time" to "string" with date-time format.
func fixTimeTypes(defs spec.Definitions) {
	for name, schema := range defs {
		fixTimeInSchema(&schema)
		defs[name] = schema
	}
}

func fixTimeInSchema(schema *spec.Schema) {
	if schema == nil {
		return
	}
	for k, v := range schema.Properties {
		if v.Type != nil && len(v.Type) > 0 && v.Type[0] == "Time" {
			v.Type = spec.StringOrArray{"string"}
			v.Format = "date-time"
			schema.Properties[k] = v
		}
		fixTimeInSchema(&v)
		schema.Properties[k] = v
	}
	if schema.Items != nil && schema.Items.Schema != nil {
		fixTimeInSchema(schema.Items.Schema)
	}
	if schema.AdditionalProperties != nil && schema.AdditionalProperties.Schema != nil {
		fixTimeInSchema(schema.AdditionalProperties.Schema)
	}
	for i := range schema.AllOf {
		fixTimeInSchema(&schema.AllOf[i])
	}
}

// fixSchemaRef fixes a schema's $ref to use simple definition names.
func fixSchemaRef(schema *spec.Schema) {
	if schema == nil {
		return
	}
	if schema.Ref.String() != "" {
		ref := schema.Ref.String()
		// Strip package prefix from #/definitions/model.Component -> #/definitions/Component
		if idx := strings.LastIndex(ref, "."); idx != -1 {
			// Find the last / before the .
			slashIdx := strings.LastIndex(ref[:idx], "/")
			if slashIdx != -1 {
				newRef := ref[:slashIdx+1] + ref[idx+1:]
				schema.Ref = spec.MustCreateRef(newRef)
			}
		}
	}
	// Fix items (for arrays)
	if schema.Items != nil && schema.Items.Schema != nil {
		fixSchemaRef(schema.Items.Schema)
	}
	// Fix additional properties
	if schema.AdditionalProperties != nil && schema.AdditionalProperties.Schema != nil {
		fixSchemaRef(schema.AdditionalProperties.Schema)
	}
	// Fix allOf
	for i := range schema.AllOf {
		fixSchemaRef(&schema.AllOf[i])
	}
	// Fix properties
	for k, v := range schema.Properties {
		fixSchemaRef(&v)
		schema.Properties[k] = v
	}
}

// RegisterSwaggerUI serves the swagger-ui static files.
func RegisterSwaggerUI(container *restful.Container, cfg SwaggerConfig) {
	ws := new(restful.WebService).
		Path("/swagger").
		Produces("text/html")

	ws.Route(ws.GET("").To(func(req *restful.Request, resp *restful.Response) {
		resp.Header().Set("Content-Type", "text/html; charset=utf-8")
		resp.Write([]byte(swaggerUIHTML(cfg)))
	}))

	ws.Route(ws.GET("/").To(func(req *restful.Request, resp *restful.Response) {
		http.Redirect(resp.ResponseWriter, req.Request, "/swagger", http.StatusMovedPermanently)
	}))

	container.Add(ws)
}

func swaggerUIHTML(cfg SwaggerConfig) string {
	return `<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <title>` + cfg.Title + ` - Swagger UI</title>
    <link rel="stylesheet" type="text/css" href="https://unpkg.com/swagger-ui-dist@5.11.0/swagger-ui.css">
    <link rel="icon" type="image/png" href="https://unpkg.com/swagger-ui-dist@5.11.0/favicon-32x32.png" sizes="32x32">
    <style>
        html { box-sizing: border-box; overflow-y: scroll; }
        *, *:before, *:after { box-sizing: inherit; }
        body { margin: 0; background: #fafafa; }
        .topbar { display: none !important; }
    </style>
</head>
<body>
    <div id="swagger-ui"></div>
    <script src="https://unpkg.com/swagger-ui-dist@5.11.0/swagger-ui-bundle.js" charset="UTF-8"></script>
    <script src="https://unpkg.com/swagger-ui-dist@5.11.0/swagger-ui-standalone-preset.js" charset="UTF-8"></script>
    <script>
        window.onload = function() {
            window.ui = SwaggerUIBundle({
                url: "` + cfg.BasePath + `apidocs.json",
                dom_id: '#swagger-ui',
                deepLinking: true,
                presets: [
                    SwaggerUIBundle.presets.apis,
                    SwaggerUIStandalonePreset
                ],
                plugins: [
                    SwaggerUIBundle.plugins.DownloadUrl
                ],
                layout: "StandaloneLayout"
            });
        };
    </script>
</body>
</html>`
}
