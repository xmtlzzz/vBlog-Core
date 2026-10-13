package api

import (
	"encoding/xml"
	"fmt"
	"net/http"
	"time"

	restfulspec "github.com/emicklei/go-restful-openapi/v2"
	restful "github.com/emicklei/go-restful/v3"
	"gorm.io/gorm"
	"vblog-core/model"
)

type RSSResource struct{ DB *gorm.DB }

type RSSFeed struct {
	XMLName   xml.Name   `xml:"rss"`
	Version   string     `xml:"version,attr"`
	XmlnsAtom string     `xml:"xmlns:atom,attr"`
	Channel   RSSChannel `xml:"channel"`
}

type RSSAtomLink struct {
	XMLName xml.Name `xml:"atom:link"`
	Href    string   `xml:"href,attr"`
	Rel     string   `xml:"rel,attr"`
	Type    string   `xml:"type,attr"`
}

type RSSChannel struct {
	Title       string      `xml:"title"`
	Link        string      `xml:"link"`
	AtomLink    RSSAtomLink `xml:"atom:link"`
	Description string      `xml:"description"`
	Language    string      `xml:"language"`
	LastBuild   string      `xml:"lastBuildDate"`
	Items       []RSSItem   `xml:"item"`
}

type RSSGUID struct {
	IsPermaLink string `xml:"isPermaLink,attr"`
	Value       string `xml:",chardata"`
}

type RSSItem struct {
	Title       string  `xml:"title"`
	Link        string  `xml:"link"`
	Description string  `xml:"description"`
	PubDate     string  `xml:"pubDate"`
	GUID        RSSGUID `xml:"guid"`
}

func (r *RSSResource) Register(ws *restful.WebService) {
	paths := []string{
		"/feed.xml", "/rss.xml", "/feed", "/rss",
		"/api/feed.xml", "/api/rss.xml", "/api/feed", "/api/rss",
	}
	for _, p := range paths {
		ws.Route(ws.GET(p).To(r.feed).
			Doc("Get RSS feed of published posts").
			Notes("Returns an RSS 2.0 feed of the 20 most recent published posts.").
			Metadata(restfulspec.KeyOpenAPITags, []string{"rss"}).
			Writes(RSSFeed{}).
			Produces("application/xml").
			Returns(200, "OK", RSSFeed{}))
		ws.Route(ws.HEAD(p).To(r.feed).
			Doc("Get RSS feed headers").
			Metadata(restfulspec.KeyOpenAPITags, []string{"rss"}).
			Produces("application/xml"))
	}
}

func (r *RSSResource) feed(req *restful.Request, resp *restful.Response) {
	scheme := "http"
	if req.Request.TLS != nil || req.Request.Header.Get("X-Forwarded-Proto") == "https" {
		scheme = "https"
	}
	host := req.Request.Header.Get("X-Forwarded-Host")
	if host == "" {
		host = req.Request.Host
	}
	baseURL := fmt.Sprintf("%s://%s", scheme, host)

	var posts []model.Post
	if r.DB != nil {
		r.DB.Preload("Tags").Where("status = ?", "published").Order("created_at DESC").Limit(20).Find(&posts)
	}

	var items []RSSItem
	for _, p := range posts {
		items = append(items, RSSItem{
			Title:       p.Title,
			Link:        fmt.Sprintf("%s/post/%d", baseURL, p.ID),
			Description: p.Excerpt,
			PubDate:     p.CreatedAt.Format(time.RFC1123Z),
			GUID: RSSGUID{
				IsPermaLink: "true",
				Value:       fmt.Sprintf("%s/post/%d", baseURL, p.ID),
			},
		})
	}

	// Get site title from settings
	siteTitle := "vBlog"
	if r.DB != nil {
		var setting model.Setting
		if err := r.DB.Where("key = ?", "site_title").First(&setting).Error; err == nil && setting.Value != "" {
			siteTitle = setting.Value
		}
	}

	feed := RSSFeed{
		Version:   "2.0",
		XmlnsAtom: "http://www.w3.org/2005/Atom",
		Channel: RSSChannel{
			Title: siteTitle,
			Link:  baseURL + "/",
			AtomLink: RSSAtomLink{
				Href: fmt.Sprintf("%s/feed.xml", baseURL),
				Rel:  "self",
				Type: "application/rss+xml",
			},
			Description: "RSS Feed for " + siteTitle,
			Language:    "zh-CN",
			LastBuild:   time.Now().Format(time.RFC1123Z),
			Items:       items,
		},
	}

	resp.Header().Set("Content-Type", "application/xml; charset=utf-8")
	resp.Header().Set("Cache-Control", "public, max-age=600")
	if req.Request.Method == http.MethodHead {
		resp.WriteHeader(http.StatusOK)
		return
	}

	resp.Write([]byte(xml.Header))
	encoder := xml.NewEncoder(resp)
	encoder.Indent("", "  ")
	_ = encoder.Encode(feed)
}
