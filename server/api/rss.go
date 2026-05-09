package api

import (
	"encoding/xml"
	"fmt"
	"time"

	restful "github.com/emicklei/go-restful/v3"
	"gorm.io/gorm"
	"vblog-core/model"
)

type RSSResource struct{ DB *gorm.DB }

type RSSFeed struct {
	XMLName xml.Name   `xml:"rss"`
	Version string     `xml:"version,attr"`
	Channel RSSChannel `xml:"channel"`
}

type RSSChannel struct {
	Title       string    `xml:"title"`
	Link        string    `xml:"link"`
	Description string    `xml:"description"`
	Language    string    `xml:"language"`
	LastBuild   string    `xml:"lastBuildDate"`
	Items       []RSSItem `xml:"item"`
}

type RSSItem struct {
	Title       string `xml:"title"`
	Link        string `xml:"link"`
	Description string `xml:"description"`
	PubDate     string `xml:"pubDate"`
	GUID        string `xml:"guid"`
}

func (r *RSSResource) Register(ws *restful.WebService) {
	ws.Route(ws.GET("/api/rss").To(r.feed))
}

func (r *RSSResource) feed(req *restful.Request, resp *restful.Response) {
	var posts []model.Post
	r.DB.Preload("Tags").Where("status = ?", "published").Order("created_at DESC").Limit(20).Find(&posts)

	var items []RSSItem
	for _, p := range posts {
		items = append(items, RSSItem{
			Title:       p.Title,
			Link:        fmt.Sprintf("/post/%d", p.ID),
			Description: p.Excerpt,
			PubDate:     p.CreatedAt.Format(time.RFC1123Z),
			GUID:        fmt.Sprintf("post-%d", p.ID),
		})
	}

	// Get site title from settings
	siteTitle := "vBlog"
	var setting model.Setting
	if err := r.DB.Where("key = ?", "site_title").First(&setting).Error; err == nil && setting.Value != "" {
		siteTitle = setting.Value
	}

	feed := RSSFeed{
		Version: "2.0",
		Channel: RSSChannel{
			Title:       siteTitle,
			Link:        "/",
			Description: "RSS Feed for " + siteTitle,
			Language:    "zh-CN",
			LastBuild:   time.Now().Format(time.RFC1123Z),
			Items:       items,
		},
	}

	resp.Header().Set("Content-Type", "application/xml; charset=utf-8")
	resp.Write([]byte(xml.Header))
	encoder := xml.NewEncoder(resp)
	encoder.Indent("", "  ")
	encoder.Encode(feed)
}
