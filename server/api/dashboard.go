package api

import (
	restful "github.com/emicklei/go-restful/v3"
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
		Doc("get dashboard statistics"))
}

func (d *DashboardResource) Stats(req *restful.Request, resp *restful.Response) {
	var postCount, viewTotal, commentCount, tagCount int64
	d.DB.Model(&model.Post{}).Where("status = ?", "published").Count(&postCount)
	d.DB.Model(&model.Post{}).Select("COALESCE(SUM(views), 0)").Scan(&viewTotal)
	d.DB.Model(&model.Comment{}).Count(&commentCount)
	d.DB.Model(&model.Tag{}).Count(&tagCount)
	resp.WriteEntity(map[string]interface{}{
		"total_posts":    postCount,
		"total_views":    viewTotal,
		"total_comments": commentCount,
		"total_tags":     tagCount,
	})
}
