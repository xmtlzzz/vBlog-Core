package api

import (
	"testing"

	restful "github.com/emicklei/go-restful/v3"
)

func TestPostResourceRegister(t *testing.T) {
	p := &PostResource{}
	ws := new(restful.WebService)
	p.Register(ws)

	routes := ws.Routes()
	if len(routes) != 5 {
		t.Errorf("expected 5 routes, got %d", len(routes))
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
