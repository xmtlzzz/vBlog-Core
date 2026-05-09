package api

import (
	"testing"

	restful "github.com/emicklei/go-restful/v3"
)

func TestDashboardResourceRegister(t *testing.T) {
	d := &DashboardResource{}
	ws := new(restful.WebService)
	d.Register(ws)

	routes := ws.Routes()
	if len(routes) != 1 {
		t.Errorf("expected 1 route, got %d", len(routes))
	}

	if routes[0].Method != "GET" {
		t.Errorf("expected method GET, got %s", routes[0].Method)
	}
	if routes[0].Path != "/api/dashboard/stats" {
		t.Errorf("expected path /api/dashboard/stats, got %s", routes[0].Path)
	}
}
