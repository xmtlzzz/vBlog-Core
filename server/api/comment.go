package api

import (
	"net/http"
	"strconv"

	restful "github.com/emicklei/go-restful/v3"
	restfulspec "github.com/emicklei/go-restful-openapi/v2"
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
		Doc("List all comments with pagination and filters").
		Notes("Returns a paginated list of comments. Admin only.").
		Metadata(restfulspec.KeyOpenAPITags, []string{"comments"}).
		Param(ws.QueryParameter("page", "Page number").DataType("integer").DefaultValue("1")).
		Param(ws.QueryParameter("per_page", "Items per page").DataType("integer").DefaultValue("10")).
		Param(ws.QueryParameter("status", "Filter by status").DataType("string").AllowableValues(map[string]string{"pending": "Pending", "approved": "Approved", "spam": "Spam"})).
		Param(ws.QueryParameter("search", "Search in comment body").DataType("string")).
		Writes(CommentListResponse{}).
		Returns(200, "OK", CommentListResponse{}).
		Returns(401, "Unauthorized", ErrorResponse{}).
		Returns(500, "Internal Server Error", ErrorResponse{}))

	ws.Route(ws.POST("/api/comments").To(c.create).
		Doc("Create a new comment (admin)").
		Notes("Creates a new comment with admin privileges. Requires authentication.").
		Metadata(restfulspec.KeyOpenAPITags, []string{"comments"}).
		Reads(model.Comment{}).
		Writes(model.Comment{}).
		Returns(201, "Created", model.Comment{}).
		Returns(400, "Bad Request", ErrorResponse{}).
		Returns(401, "Unauthorized", ErrorResponse{}))

	ws.Route(ws.PATCH("/api/comments/{id}/approve").To(c.approve).
		Doc("Approve a comment").
		Notes("Approves a pending comment. Requires authentication.").
		Metadata(restfulspec.KeyOpenAPITags, []string{"comments"}).
		Param(ws.PathParameter("id", "Comment ID").DataType("integer")).
		Returns(200, "OK", MessageResponse{}).
		Returns(401, "Unauthorized", ErrorResponse{}).
		Returns(404, "Not Found", ErrorResponse{}))

	ws.Route(ws.PATCH("/api/comments/{id}/spam").To(c.markSpam).
		Doc("Mark comment as spam").
		Notes("Marks a comment as spam. Requires authentication.").
		Metadata(restfulspec.KeyOpenAPITags, []string{"comments"}).
		Param(ws.PathParameter("id", "Comment ID").DataType("integer")).
		Returns(200, "OK", MessageResponse{}).
		Returns(401, "Unauthorized", ErrorResponse{}).
		Returns(404, "Not Found", ErrorResponse{}))

	ws.Route(ws.DELETE("/api/comments/{id}").To(c.delete).
		Doc("Delete a comment").
		Notes("Permanently deletes a comment. Requires authentication.").
		Metadata(restfulspec.KeyOpenAPITags, []string{"comments"}).
		Param(ws.PathParameter("id", "Comment ID").DataType("integer")).
		Returns(200, "OK", MessageResponse{}).
		Returns(401, "Unauthorized", ErrorResponse{}).
		Returns(404, "Not Found", ErrorResponse{}))

	// Public endpoints (no auth)
	ws.Route(ws.GET("/api/posts/{postId}/comments").To(c.listByPost).
		Doc("List approved comments for a specific post").
		Notes("Returns all approved comments for a given post. Public endpoint.").
		Metadata(restfulspec.KeyOpenAPITags, []string{"comments"}).
		Param(ws.PathParameter("postId", "Post ID").DataType("integer")).
		Writes(CommentListResponse{}).
		Returns(200, "OK", CommentListResponse{}))

	ws.Route(ws.POST("/api/posts/{postId}/comments").To(c.createPublic).
		Doc("Submit a comment on a post (public, pending approval)").
		Notes("Submits a new comment on a post. Comments are pending until approved by admin.").
		Metadata(restfulspec.KeyOpenAPITags, []string{"comments"}).
		Param(ws.PathParameter("postId", "Post ID").DataType("integer")).
		Reads(model.Comment{}).
		Writes(MessageResponse{}).
		Returns(201, "Created", MessageResponse{}).
		Returns(400, "Bad Request", ErrorResponse{}))
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
