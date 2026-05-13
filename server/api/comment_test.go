package api

import (
	"testing"

	restful "github.com/emicklei/go-restful/v3"
)

func TestCommentResourceRegister(t *testing.T) {
	cr := &CommentResource{}
	ws := new(restful.WebService)
	cr.Register(ws)

	routes := ws.Routes()
	if len(routes) != 7 {
		t.Errorf("expected 7 routes, got %d", len(routes))
	}

	expected := []struct {
		method string
		path   string
	}{
		{"GET", "/api/comments"},
		{"POST", "/api/comments"},
		{"PATCH", "/api/comments/{id}/approve"},
		{"PATCH", "/api/comments/{id}/spam"},
		{"DELETE", "/api/comments/{id}"},
		{"GET", "/api/posts/{postId}/comments"},
		{"POST", "/api/posts/{postId}/comments"},
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
