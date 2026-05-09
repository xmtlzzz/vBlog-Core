package api

import (
	"testing"

	restful "github.com/emicklei/go-restful/v3"
)

func TestTagResourceRegister(t *testing.T) {
	tr := &TagResource{}
	ws := new(restful.WebService)
	tr.Register(ws)

	routes := ws.Routes()
	if len(routes) != 4 {
		t.Errorf("expected 4 routes, got %d", len(routes))
	}

	expected := []struct {
		method string
		path   string
	}{
		{"GET", "/api/tags"},
		{"POST", "/api/tags"},
		{"PUT", "/api/tags/{id}"},
		{"DELETE", "/api/tags/{id}"},
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
