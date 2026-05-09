package api

import (
	"testing"

	restful "github.com/emicklei/go-restful/v3"
)

func TestSettingResourceRegister(t *testing.T) {
	sr := &SettingResource{}
	ws := new(restful.WebService)
	sr.Register(ws)

	routes := ws.Routes()
	if len(routes) != 3 {
		t.Errorf("expected 3 routes, got %d", len(routes))
	}

	expected := []struct {
		method string
		path   string
	}{
		{"GET", "/api/settings"},
		{"PUT", "/api/settings"},
		{"POST", "/api/settings/reset"},
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
