package internal_test

import (
	"FunPDF/internal"
	"FunPDF/internal/handler"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestRouterRegistersTranslatorRoutes(t *testing.T) {
	gin.SetMode(gin.TestMode)
	engine := gin.New()
	internal.NewRouter(
		handler.NewFileHandler(),
		handler.NewAlbumHandler(),
		handler.NewTranslatorHandler(),
		handler.NewProviderHandler(),
		handler.NewModelHandler(),
		handler.NewChatSessionHandler(),
	).Setup(engine)

	routes := map[string]bool{}
	for _, route := range engine.Routes() {
		routes[route.Method+" "+route.Path] = true
	}

	expected := []string{
		"GET /api/translators",
		"POST /api/translators",
		"POST /api/translators/:translator_name",
	}
	for _, route := range expected {
		if !routes[route] {
			t.Fatalf("missing route %s", route)
		}
	}
	for route := range routes {
		if strings.Contains(route, " /api/translators") {
			found := false
			for _, expectedRoute := range expected {
				if route == expectedRoute {
					found = true
					break
				}
			}
			if !found {
				t.Fatalf("unexpected translator route %s", route)
			}
		}
	}
}
