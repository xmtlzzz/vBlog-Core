package api

import (
	"encoding/xml"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	restful "github.com/emicklei/go-restful/v3"
)

func TestRSSResourceRegister(t *testing.T) {
	r := &RSSResource{}
	ws := new(restful.WebService)
	r.Register(ws)

	routes := ws.Routes()
	if len(routes) != 16 { // 8 paths * 2 methods (GET and HEAD)
		t.Fatalf("expected 16 routes, got %d", len(routes))
	}
}

func TestRSSFeedOutput(t *testing.T) {
	container := restful.NewContainer()
	ws := new(restful.WebService).Path("/").Produces(restful.MIME_XML)
	(&RSSResource{}).Register(ws)
	container.Add(ws)

	// Test GET /feed.xml
	req := httptest.NewRequest(http.MethodGet, "/feed.xml", nil)
	req.Host = "example.com"
	rec := httptest.NewRecorder()
	container.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}

	body := rec.Body.String()
	if !strings.Contains(body, `<rss version="2.0" xmlns:atom="http://www.w3.org/2005/Atom">`) {
		t.Fatalf("missing rss tag with atom namespace: %s", body)
	}
	if !strings.Contains(body, `<atom:link href="http://example.com/feed.xml" rel="self" type="application/rss+xml"`) {
		t.Fatalf("missing atom:link: %s", body)
	}

	var parsed RSSFeed
	if err := xml.Unmarshal(rec.Body.Bytes(), &parsed); err != nil {
		t.Fatalf("invalid XML output: %v", err)
	}

	// Test HEAD /feed.xml
	headReq := httptest.NewRequest(http.MethodHead, "/feed.xml", nil)
	headRec := httptest.NewRecorder()
	container.ServeHTTP(headRec, headReq)

	if headRec.Code != http.StatusOK {
		t.Fatalf("HEAD expected 200, got %d", headRec.Code)
	}
	if headRec.Body.Len() != 0 {
		t.Fatalf("HEAD should not write body, got %d bytes", headRec.Body.Len())
	}
}
