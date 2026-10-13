package api

import (
	"encoding/xml"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	restful "github.com/emicklei/go-restful/v3"
)

func TestSitemapResourceRegister(t *testing.T) {
	s := &SitemapResource{}
	ws := new(restful.WebService)
	s.Register(ws)

	routes := ws.Routes()
	if len(routes) != 8 { // 4 paths * 2 methods (GET and HEAD)
		t.Fatalf("expected 8 routes, got %d", len(routes))
	}
}

func TestSitemapOutput(t *testing.T) {
	container := restful.NewContainer()
	ws := new(restful.WebService).Path("/").Produces(restful.MIME_XML)
	(&SitemapResource{}).Register(ws)
	container.Add(ws)

	// Test GET /sitemap.xml
	req := httptest.NewRequest(http.MethodGet, "/sitemap.xml", nil)
	req.Host = "example.com"
	rec := httptest.NewRecorder()
	container.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}

	body := rec.Body.String()
	if !strings.Contains(body, `<urlset xmlns="http://www.sitemaps.org/schemas/sitemap/0.9">`) {
		t.Fatalf("missing urlset tag: %s", body)
	}
	if !strings.Contains(body, `http://example.com/`) {
		t.Fatalf("missing base URL in sitemap: %s", body)
	}

	var parsed SitemapURLSet
	if err := xml.Unmarshal(rec.Body.Bytes(), &parsed); err != nil {
		t.Fatalf("invalid XML output: %v", err)
	}

	// Test HEAD /sitemap.xml
	headReq := httptest.NewRequest(http.MethodHead, "/sitemap.xml", nil)
	headRec := httptest.NewRecorder()
	container.ServeHTTP(headRec, headReq)

	if headRec.Code != http.StatusOK {
		t.Fatalf("HEAD expected 200, got %d", headRec.Code)
	}
	if headRec.Body.Len() != 0 {
		t.Fatalf("HEAD should not write body, got %d bytes", headRec.Body.Len())
	}
}
