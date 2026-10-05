package api

import (
	"net/http"
	"net/http/httptest"
	"testing"

	restful "github.com/emicklei/go-restful/v3"
)

func TestPostReadWithoutOptionalAuth(t *testing.T) {
	container := restful.NewContainer()
	ws := new(restful.WebService).Path("/").Produces(restful.MIME_JSON)
	(&PostResource{}).Register(ws)
	container.Add(ws)
	response := httptest.NewRecorder()
	container.ServeHTTP(response, httptest.NewRequest(http.MethodGet, "/api/posts/not-a-number", nil))
	if response.Code != http.StatusBadRequest {
		t.Fatalf("invalid post ID returned %d, want 400", response.Code)
	}
}

func TestPostResourceRegister(t *testing.T) {
	p := &PostResource{}
	ws := new(restful.WebService)
	p.Register(ws)

	routes := ws.Routes()
	if len(routes) != 8 {
		t.Errorf("expected 8 routes, got %d", len(routes))
	}

	// Verify route methods and paths
	expected := []struct {
		method string
		path   string
	}{
		{"GET", "/api/posts"},
		{"GET", "/api/posts/{id}"},
		{"POST", "/api/posts"},
		{"PUT", "/api/posts/{id}"},
		{"DELETE", "/api/posts/{id}"},
		{"GET", "/api/posts/trash"},
		{"POST", "/api/posts/{id}/restore"},
		{"DELETE", "/api/posts/{id}/permanent"},
	}

	for i, exp := range expected {
		if routes[i].Method != exp.method {
			t.Errorf("route %d: expected method %s, got %s", i, exp.method, routes[i].Method)
		}
		if routes[i].Path != exp.path {
			t.Errorf("route %d: expected path %s, got %s", i, exp.path, routes[i].Path)
		}
	}
}
