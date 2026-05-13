package api

import (
	"testing"

	restful "github.com/emicklei/go-restful/v3"
)

func TestComponentResourceRegister(t *testing.T) {
	cr := &ComponentResource{}
	ws := new(restful.WebService)
	cr.Register(ws)

	routes := ws.Routes()
	if len(routes) != 6 {
		t.Errorf("expected 6 routes, got %d", len(routes))
	}

	expected := []struct {
		method string
		path   string
	}{
		{"GET", "/api/components"},
		{"GET", "/api/components/active"},
		{"POST", "/api/components"},
		{"PUT", "/api/components/{id}"},
		{"DELETE", "/api/components/{id}"},
		{"PATCH", "/api/components/{id}/toggle"},
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
