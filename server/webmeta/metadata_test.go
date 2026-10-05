package webmeta

import (
	"strings"
	"testing"
)

func TestArticleSharingMetadata(t *testing.T) {
	page := NewPage(map[string]string{"site_title": "vBlog", "site_url": "https://vblog.xmtlz.dev"}, "http://localhost:8080", "/post/3")
	page.Title = `项目 "导航" <script>`
	page.Description = "项目介绍"
	page.Article = true
	rendered := string(Render([]byte(`<html><head><title>vBlog</title><meta name="description" content="old"><meta property="og:title" content="old"></head><body><div id="app"></div></body></html>`), page))
	for _, expected := range []string{`项目 &#34;导航&#34; &lt;script&gt; · vBlog`, `content="article"`, `href="https://vblog.xmtlz.dev/post/3"`, `<div id="app">`} {
		if !strings.Contains(rendered, expected) {
			t.Errorf("missing %s in %s", expected, rendered)
		}
	}
	if strings.Contains(rendered, `content="old"`) || strings.Count(rendered, `<title>`) != 1 {
		t.Fatal("stale or duplicate metadata")
	}
}
