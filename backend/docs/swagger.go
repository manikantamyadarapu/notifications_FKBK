package docs

import (
	_ "embed"
	"net/http"
)

//go:embed swagger.yaml
var swaggerYAML []byte

func RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/swagger.yaml", serveSwaggerSpec)
	mux.HandleFunc("/docs", serveSwaggerUI)
}

func serveSwaggerSpec(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "application/yaml")
	_, _ = w.Write(swaggerYAML)
}

func serveSwaggerUI(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	_, _ = w.Write([]byte(`<!doctype html>
<html>
  <head>
    <meta charset="utf-8" />
    <title>Notification Service API Docs</title>
    <link rel="stylesheet" href="https://unpkg.com/swagger-ui-dist@5/swagger-ui.css" />
  </head>
  <body>
    <div id="swagger-ui"></div>
    <script src="https://unpkg.com/swagger-ui-dist@5/swagger-ui-bundle.js"></script>
    <script>
      window.ui = SwaggerUIBundle({ url: '/swagger.yaml', dom_id: '#swagger-ui' });
    </script>
  </body>
</html>`))
}

