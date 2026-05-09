package api

import (
	"net/http"
	"strconv"

	restful "github.com/emicklei/go-restful/v3"
	"vblog-core/model"
	"vblog-core/service"
)

// CommentResource handles comment REST endpoints.
type CommentResource struct {
	Service *service.CommentService
}

// Register adds comment routes to the given WebService.
func (c *CommentResource) Register(ws *restful.WebService) {
	ws.Route(ws.GET("/api/comments").To(c.list).
		Doc("list comments").
		Param(ws.QueryParameter("page", "page number").DefaultValue("1")).
		Param(ws.QueryParameter("per_page", "items per page").DefaultValue("10")).
		Param(ws.QueryParameter("status", "filter by status")).
		Param(ws.QueryParameter("search", "search in body")))

	ws.Route(ws.POST("/api/comments").To(c.create).
		Doc("create a comment"))

	ws.Route(ws.PATCH("/api/comments/{id}/approve").To(c.approve).
		Doc("approve a comment").
		Param(ws.PathParameter("id", "comment ID")))

	ws.Route(ws.PATCH("/api/comments/{id}/spam").To(c.markSpam).
		Doc("mark comment as spam").
		Param(ws.PathParameter("id", "comment ID")))

	ws.Route(ws.DELETE("/api/comments/{id}").To(c.delete).
		Doc("delete a comment").
		Param(ws.PathParameter("id", "comment ID")))

	// Public endpoints (no auth)
	ws.Route(ws.GET("/api/posts/{postId}/comments").To(c.listByPost).
		Doc("list approved comments for a post").
		Param(ws.PathParameter("postId", "post ID")))
	ws.Route(ws.POST("/api/posts/{postId}/comments").To(c.createPublic).
		Doc("submit a comment on a post").
		Param(ws.PathParameter("postId", "post ID")))
}

func (c *CommentResource) list(req *restful.Request, resp *restful.Response) {
	page, _ := strconv.Atoi(req.QueryParameter("page"))
	if page < 1 {
		page = 1
	}
	perPage, _ := strconv.Atoi(req.QueryParameter("per_page"))
	if perPage < 1 {
		perPage = 10
	}
	status := req.QueryParameter("status")
	search := req.QueryParameter("search")

	comments, total, err := c.Service.List(page, perPage, status, search)
	if err != nil {
		resp.WriteError(http.StatusInternalServerError, err)
		return
	}

	resp.WriteEntity(map[string]interface{}{
		"data":  comments,
		"total": total,
		"page":  page,
	})
}

func (c *CommentResource) create(req *restful.Request, resp *restful.Response) {
	comment := model.Comment{}
	if err := req.ReadEntity(&comment); err != nil {
		resp.WriteError(http.StatusBadRequest, err)
		return
	}
	if err := c.Service.Create(&comment); err != nil {
		resp.WriteError(http.StatusInternalServerError, err)
		return
	}
	resp.WriteHeaderAndEntity(http.StatusCreated, comment)
}

func (c *CommentResource) approve(req *restful.Request, resp *restful.Response) {
	id, err := strconv.Atoi(req.PathParameter("id"))
	if err != nil {
		resp.WriteError(http.StatusBadRequest, err)
		return
	}
	if err := c.Service.Approve(uint(id)); err != nil {
		resp.WriteError(http.StatusInternalServerError, err)
		return
	}
	resp.WriteHeader(http.StatusOK)
}

func (c *CommentResource) markSpam(req *restful.Request, resp *restful.Response) {
	id, err := strconv.Atoi(req.PathParameter("id"))
	if err != nil {
		resp.WriteError(http.StatusBadRequest, err)
		return
	}
	if err := c.Service.MarkSpam(uint(id)); err != nil {
		resp.WriteError(http.StatusInternalServerError, err)
		return
	}
	resp.WriteHeader(http.StatusOK)
}

func (c *CommentResource) delete(req *restful.Request, resp *restful.Response) {
	id, err := strconv.Atoi(req.PathParameter("id"))
	if err != nil {
		resp.WriteError(http.StatusBadRequest, err)
		return
	}
	if err := c.Service.Delete(uint(id)); err != nil {
		resp.WriteError(http.StatusInternalServerError, err)
		return
	}
	resp.WriteHeader(http.StatusOK)
}

func (c *CommentResource) listByPost(req *restful.Request, resp *restful.Response) {
	postId, _ := strconv.ParseUint(req.PathParameter("postId"), 10, 32)
	var comments []model.Comment
	c.Service.DB.Where("post_id = ? AND status = ?", postId, "approved").
		Order("created_at DESC").Find(&comments)
	resp.WriteEntity(map[string]interface{}{"data": comments})
}

func (c *CommentResource) createPublic(req *restful.Request, resp *restful.Response) {
	postId, _ := strconv.ParseUint(req.PathParameter("postId"), 10, 32)
	var comment model.Comment
	if err := req.ReadEntity(&comment); err != nil {
		resp.WriteError(http.StatusBadRequest, err)
		return
	}
	comment.PostID = uint(postId)
	comment.Status = "pending"
	if err := c.Service.Create(&comment); err != nil {
		resp.WriteError(http.StatusInternalServerError, err)
		return
	}
	resp.WriteHeaderAndEntity(http.StatusCreated, map[string]string{"message": "评论已提交，等待审核"})
}
