package api

import (
	"testing"

	restful "github.com/emicklei/go-restful/v3"
)

func TestAuthResourceRegister(t *testing.T) {
	a := &AuthResource{}
	ws := new(restful.WebService)
	a.Register(ws)

	routes := ws.Routes()
	if len(routes) != 1 {
		t.Errorf("expected 1 route, got %d", len(routes))
	}

	if routes[0].Method != "POST" {
		t.Errorf("expected method POST, got %s", routes[0].Method)
	}
	if routes[0].Path != "/api/auth/login" {
		t.Errorf("expected path /api/auth/login, got %s", routes[0].Path)
	}
}
