package webmeta

import (
	"fmt"
	"html"
	"net/url"
	"regexp"
	"strings"
	"vblog-core/service"
)

type Page struct {
	Site, Title, Description, URL string
	Article, Missing              bool
}

func NewPage(settings map[string]string, origin, path string) Page {
	site := settings["site_title"]
	if site == "" {
		site = "vBlog"
	}
	description := settings["description"]
	if description == "" {
		description = settings["site_description"]
	}
	if description == "" {
		description = "记录开发、技术与日常。"
	}
	base, err := url.Parse(settings["site_url"])
	if err != nil || base == nil || (base.Scheme != "https" && base.Scheme != "http") || base.Host == "" {
		base, _ = url.Parse(origin)
	}
	canonical := base.ResolveReference(&url.URL{Path: path}).String()
	titles := map[string]string{"/archives": "归档", "/modules": "模块", "/tags": "标签", "/about": "关于"}
	return Page{Site: site, Title: titles[path], Description: service.BuildExcerpt(description, 160), URL: canonical}
}

var oldMetadata = regexp.MustCompile(`(?is)<title\b[^>]*>.*?</title>|<meta\b[^>]*(?:name=["']description["']|property=["']og:[^"']*["'])[^>]*>|<link\b[^>]*rel=["']canonical["'][^>]*>`)

func Render(index []byte, page Page) []byte {
	title := page.Site
	if page.Title != "" {
		title = page.Title + " · " + page.Site
	}
	ogTitle := page.Title
	if ogTitle == "" {
		ogTitle = page.Site
	}
	kind := "website"
	if page.Article {
		kind = "article"
	}
	head := fmt.Sprintf(`<title>%s</title><meta name="description" content="%s"><meta property="og:title" content="%s"><meta property="og:description" content="%s"><meta property="og:site_name" content="%s"><meta property="og:type" content="%s"><meta property="og:url" content="%s"><link rel="canonical" href="%s">`, html.EscapeString(title), html.EscapeString(page.Description), html.EscapeString(ogTitle), html.EscapeString(page.Description), html.EscapeString(page.Site), kind, html.EscapeString(page.URL), html.EscapeString(page.URL))
	if page.Missing {
		head += `<meta name="robots" content="noindex">`
	}
	clean := oldMetadata.ReplaceAllString(string(index), "")
	return []byte(strings.Replace(clean, "</head>", head+"</head>", 1))
}
