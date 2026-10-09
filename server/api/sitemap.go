package api

import (
	"encoding/xml"
	"fmt"
	"time"

	restfulspec "github.com/emicklei/go-restful-openapi/v2"
	restful "github.com/emicklei/go-restful/v3"
	"gorm.io/gorm"
	"vblog-core/model"
)

type SitemapResource struct{ DB *gorm.DB }

type SitemapURLSet struct {
	XMLName xml.Name     `xml:"urlset"`
	Xmlns   string       `xml:"xmlns,attr"`
	URLs    []SitemapURL `xml:"url"`
}

type SitemapURL struct {
	Loc        string  `xml:"loc"`
	LastMod    string  `xml:"lastmod,omitempty"`
	ChangeFreq string  `xml:"changefreq,omitempty"`
	Priority   float32 `xml:"priority,omitempty"`
}

func (s *SitemapResource) Register(ws *restful.WebService) {
	ws.Route(ws.GET("/sitemap.xml").To(s.sitemap).
		Doc("Get Sitemap XML of the site").
		Notes("Returns sitemap.xml for search engine indexing.").
		Metadata(restfulspec.KeyOpenAPITags, []string{"sitemap"}).
		Produces("application/xml"))

	ws.Route(ws.GET("/api/sitemap.xml").To(s.sitemap).
		Doc("Get Sitemap XML of the site (API path)").
		Metadata(restfulspec.KeyOpenAPITags, []string{"sitemap"}).
		Produces("application/xml"))
}

func (s *SitemapResource) sitemap(req *restful.Request, resp *restful.Response) {
	scheme := "http"
	if req.Request.TLS != nil || req.Request.Header.Get("X-Forwarded-Proto") == "https" {
		scheme = "https"
	}
	host := req.Request.Header.Get("X-Forwarded-Host")
	if host == "" {
		host = req.Request.Host
	}
	baseURL := fmt.Sprintf("%s://%s", scheme, host)

	today := time.Now().Format("2006-01-02")

	staticRoutes := []struct {
		path     string
		freq     string
		priority float32
	}{
		{"/", "daily", 1.0},
		{"/archives", "weekly", 0.8},
		{"/modules", "weekly", 0.7},
		{"/tags", "weekly", 0.7},
		{"/friends", "weekly", 0.6},
		{"/about", "monthly", 0.6},
	}

	var urls []SitemapURL
	for _, r := range staticRoutes {
		urls = append(urls, SitemapURL{
			Loc:        baseURL + r.path,
			LastMod:    today,
			ChangeFreq: r.freq,
			Priority:   r.priority,
		})
	}

	var posts []model.Post
	if err := s.DB.Where("status = ?", "published").Order("created_at DESC").Find(&posts).Error; err == nil {
		for _, p := range posts {
			lastMod := p.UpdatedAt.Format("2006-01-02")
			if p.UpdatedAt.IsZero() {
				lastMod = p.CreatedAt.Format("2006-01-02")
			}
			urls = append(urls, SitemapURL{
				Loc:        fmt.Sprintf("%s/post/%d", baseURL, p.ID),
				LastMod:    lastMod,
				ChangeFreq: "weekly",
				Priority:   0.8,
			})
		}
	}

	urlSet := SitemapURLSet{
		Xmlns: "http://www.sitemaps.org/schemas/sitemap/0.9",
		URLs:  urls,
	}

	resp.Header().Set("Content-Type", "application/xml; charset=utf-8")
	resp.Header().Set("Cache-Control", "public, max-age=3600")
	resp.Write([]byte(xml.Header))
	encoder := xml.NewEncoder(resp)
	encoder.Indent("", "  ")
	_ = encoder.Encode(urlSet)
}
