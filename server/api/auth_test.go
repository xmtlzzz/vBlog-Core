package api

import (
	"testing"

	restful "github.com/emicklei/go-restful/v3"
)

func TestAuthResourceRegister(t *testing.T) {
	dummyAuth := func(req *restful.Request, resp *restful.Response, chain *restful.FilterChain) {
		chain.ProcessFilter(req, resp)
	}
	a := &AuthResource{Auth: dummyAuth}
	ws := new(restful.WebService)
	a.Register(ws)

	routes := ws.Routes()
	if len(routes) != 2 {
		t.Fatalf("expected 2 routes, got %d", len(routes))
	}

	if routes[0].Method != "POST" || routes[0].Path != "/api/auth/login" {
		t.Errorf("expected POST /api/auth/login, got %s %s", routes[0].Method, routes[0].Path)
	}
	if routes[1].Method != "GET" || routes[1].Path != "/api/auth/me" {
		t.Errorf("expected GET /api/auth/me, got %s %s", routes[1].Method, routes[1].Path)
	}
}
