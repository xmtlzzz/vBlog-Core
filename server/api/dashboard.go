package api

import (
	restful "github.com/emicklei/go-restful/v3"
	restfulspec "github.com/emicklei/go-restful-openapi/v2"
	"gorm.io/gorm"
	"vblog-core/model"
)

// DashboardResource provides dashboard statistics endpoints.
type DashboardResource struct {
	DB *gorm.DB
}

// Register adds dashboard routes to the given WebService.
func (d *DashboardResource) Register(ws *restful.WebService) {
	ws.Route(ws.GET("/api/dashboard/stats").To(d.Stats).
		Doc("Get dashboard statistics").
		Notes("Returns aggregate statistics including post count, total views, comments, and tags.").
		Metadata(restfulspec.KeyOpenAPITags, []string{"dashboard"}).
		Writes(DashboardStatsResponse{}).
		Returns(200, "OK", DashboardStatsResponse{}))
}

func (d *DashboardResource) Stats(req *restful.Request, resp *restful.Response) {
	var postCount, viewTotal int64
	d.DB.Model(&model.Post{}).
		Select("COUNT(*) FILTER (WHERE status = 'published'), COALESCE(SUM(views), 0)").
		Row().Scan(&postCount, &viewTotal)

	var commentCount, tagCount int64
	d.DB.Raw(`SELECT
		(SELECT COUNT(*) FROM comments),
		(SELECT COUNT(*) FROM tags)`).Row().Scan(&commentCount, &tagCount)

	resp.WriteEntity(map[string]interface{}{
		"total_posts":    postCount,
		"total_views":    viewTotal,
		"total_comments": commentCount,
		"total_tags":     tagCount,
	})
}
