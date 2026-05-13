package api

import (
	"net/http"
	"strconv"

	restful "github.com/emicklei/go-restful/v3"
	restfulspec "github.com/emicklei/go-restful-openapi/v2"
	"vblog-core/model"
	"vblog-core/service"
)

// postResp wraps a Post with formatted dates for JSON.
type postResp struct {
	ID        uint         `json:"id"`
	Title     string       `json:"title"`
	Content   string       `json:"content"`
	Excerpt   string       `json:"excerpt"`
	Status    string       `json:"status"`
	Pinned    bool         `json:"pinned"`
	Views     int          `json:"views"`
	ReadTime  int          `json:"read_time"`
	AuthorID  uint         `json:"author_id"`
	Tags      []model.Tag  `json:"tags"`
	CreatedAt string       `json:"created_at"`
	UpdatedAt string       `json:"updated_at"`
}

func newPostResp(p *model.Post) postResp {
	return postResp{
		ID: p.ID, Title: p.Title, Content: p.Content, Excerpt: p.Excerpt,
		Status: p.Status, Pinned: p.Pinned, Views: p.Views, ReadTime: p.ReadTime,
		AuthorID: p.AuthorID, Tags: p.Tags,
		CreatedAt: model.FormatDate(p.CreatedAt),
		UpdatedAt: model.FormatDateTime(p.UpdatedAt),
	}
}

// PostResource handles blog post REST endpoints.
type PostResource struct {
	Service *service.PostService
}

// Register adds post routes to the given WebService.
func (p *PostResource) Register(ws *restful.WebService) {
	ws.Route(ws.GET("/api/posts").To(p.list).
		Doc("List blog posts with pagination and filters").
		Notes("Returns a paginated list of posts. Supports filtering by tag, status, and search by title.").
		Metadata(restfulspec.KeyOpenAPITags, []string{"posts"}).
		Param(ws.QueryParameter("page", "Page number").DataType("integer").DefaultValue("1")).
		Param(ws.QueryParameter("per_page", "Items per page").DataType("integer").DefaultValue("5")).
		Param(ws.QueryParameter("tag", "Filter by tag name").DataType("string")).
		Param(ws.QueryParameter("status", "Filter by status").DataType("string").AllowableValues(map[string]string{"draft": "Draft", "published": "Published"})).
		Param(ws.QueryParameter("search", "Search by title").DataType("string")).
		Writes(PostListResponse{}).
		Returns(200, "OK", PostListResponse{}).
		Returns(500, "Internal Server Error", ErrorResponse{}))

	ws.Route(ws.GET("/api/posts/{id}").To(p.get).
		Doc("Get a single post by ID").
		Notes("Returns the full post content including tags.").
		Metadata(restfulspec.KeyOpenAPITags, []string{"posts"}).
		Param(ws.PathParameter("id", "Post ID").DataType("integer")).
		Writes(postResp{}).
		Returns(200, "OK", postResp{}).
		Returns(404, "Not Found", ErrorResponse{}))

	ws.Route(ws.POST("/api/posts").To(p.create).
		Doc("Create a new blog post").
		Notes("Creates a new blog post. Requires authentication.").
		Metadata(restfulspec.KeyOpenAPITags, []string{"posts"}).
		Reads(model.Post{}).
		Writes(postResp{}).
		Returns(201, "Created", postResp{}).
		Returns(400, "Bad Request", ErrorResponse{}).
		Returns(401, "Unauthorized", ErrorResponse{}))

	ws.Route(ws.PUT("/api/posts/{id}").To(p.update).
		Doc("Update an existing blog post").
		Notes("Updates a blog post by ID. Requires authentication.").
		Metadata(restfulspec.KeyOpenAPITags, []string{"posts"}).
		Param(ws.PathParameter("id", "Post ID").DataType("integer")).
		Reads(model.Post{}).
		Writes(postResp{}).
		Returns(200, "OK", postResp{}).
		Returns(400, "Bad Request", ErrorResponse{}).
		Returns(401, "Unauthorized", ErrorResponse{}).
		Returns(404, "Not Found", ErrorResponse{}))

	ws.Route(ws.DELETE("/api/posts/{id}").To(p.delete).
		Doc("Delete a blog post (soft delete)").
		Notes("Soft deletes a post. Can be restored from trash. Requires authentication.").
		Metadata(restfulspec.KeyOpenAPITags, []string{"posts"}).
		Param(ws.PathParameter("id", "Post ID").DataType("integer")).
		Returns(200, "OK", MessageResponse{}).
		Returns(401, "Unauthorized", ErrorResponse{}).
		Returns(404, "Not Found", ErrorResponse{}))

	ws.Route(ws.GET("/api/posts/trash").To(p.trash).
		Doc("List all soft-deleted posts").
		Notes("Returns all posts in the trash. Requires authentication.").
		Metadata(restfulspec.KeyOpenAPITags, []string{"posts"}).
		Writes(PostListResponse{}).
		Returns(200, "OK", PostListResponse{}).
		Returns(401, "Unauthorized", ErrorResponse{}))

	ws.Route(ws.POST("/api/posts/{id}/restore").To(p.restore).
		Doc("Restore a soft-deleted post").
		Notes("Restores a post from trash. Requires authentication.").
		Metadata(restfulspec.KeyOpenAPITags, []string{"posts"}).
		Param(ws.PathParameter("id", "Post ID").DataType("integer")).
		Returns(200, "OK", MessageResponse{}).
		Returns(401, "Unauthorized", ErrorResponse{}).
		Returns(404, "Not Found", ErrorResponse{}))

	ws.Route(ws.DELETE("/api/posts/{id}/permanent").To(p.permanentDelete).
		Doc("Permanently delete a post").
		Notes("Permanently deletes a post from trash. This action cannot be undone. Requires authentication.").
		Metadata(restfulspec.KeyOpenAPITags, []string{"posts"}).
		Param(ws.PathParameter("id", "Post ID").DataType("integer")).
		Returns(200, "OK", MessageResponse{}).
		Returns(401, "Unauthorized", ErrorResponse{}).
		Returns(404, "Not Found", ErrorResponse{}))
}

func (p *PostResource) list(req *restful.Request, resp *restful.Response) {
	page, _ := strconv.Atoi(req.QueryParameter("page"))
	if page < 1 {
		page = 1
	}
	perPage, _ := strconv.Atoi(req.QueryParameter("per_page"))
	if perPage < 1 {
		perPage = 5
	}
	tag := req.QueryParameter("tag")
	status := req.QueryParameter("status")
	search := req.QueryParameter("search")

	posts, total, err := p.Service.List(page, perPage, tag, status, search)
	if err != nil {
		resp.WriteError(http.StatusInternalServerError, err)
		return
	}

	list := make([]postResp, len(posts))
	for i := range posts {
		list[i] = newPostResp(&posts[i])
	}

	resp.WriteEntity(map[string]interface{}{
		"data":  list,
		"total": total,
		"page":  page,
	})
}

func (p *PostResource) get(req *restful.Request, resp *restful.Response) {
	id, err := strconv.Atoi(req.PathParameter("id"))
	if err != nil {
		resp.WriteError(http.StatusBadRequest, err)
		return
	}
	post, err := p.Service.GetByID(uint(id))
	if err != nil {
		resp.WriteError(http.StatusNotFound, err)
		return
	}
	resp.WriteEntity(newPostResp(post))
}

func (p *PostResource) create(req *restful.Request, resp *restful.Response) {
	post := model.Post{}
	if err := req.ReadEntity(&post); err != nil {
		resp.WriteError(http.StatusBadRequest, err)
		return
	}
	if err := p.Service.Create(&post); err != nil {
		resp.WriteError(http.StatusInternalServerError, err)
		return
	}
	resp.WriteHeaderAndEntity(http.StatusCreated, post)
}

func (p *PostResource) update(req *restful.Request, resp *restful.Response) {
	id, err := strconv.Atoi(req.PathParameter("id"))
	if err != nil {
		resp.WriteError(http.StatusBadRequest, err)
		return
	}
	post := model.Post{}
	if err := req.ReadEntity(&post); err != nil {
		resp.WriteError(http.StatusBadRequest, err)
		return
	}
	post.ID = uint(id)
	if err := p.Service.Update(&post); err != nil {
		resp.WriteError(http.StatusInternalServerError, err)
		return
	}
	resp.WriteEntity(post)
}

func (p *PostResource) delete(req *restful.Request, resp *restful.Response) {
	id, err := strconv.Atoi(req.PathParameter("id"))
	if err != nil {
		resp.WriteError(http.StatusBadRequest, err)
		return
	}
	if err := p.Service.Delete(uint(id)); err != nil {
		resp.WriteError(http.StatusInternalServerError, err)
		return
	}
	resp.WriteHeader(http.StatusOK)
}

func (p *PostResource) trash(req *restful.Request, resp *restful.Response) {
	posts, err := p.Service.ListTrash()
	if err != nil {
		resp.WriteError(http.StatusInternalServerError, err)
		return
	}
	list := make([]postResp, len(posts))
	for i := range posts {
		list[i] = newPostResp(&posts[i])
	}
	resp.WriteEntity(map[string]interface{}{"data": list})
}

func (p *PostResource) restore(req *restful.Request, resp *restful.Response) {
	id, err := strconv.Atoi(req.PathParameter("id"))
	if err != nil {
		resp.WriteError(http.StatusBadRequest, err)
		return
	}
	if err := p.Service.Restore(uint(id)); err != nil {
		resp.WriteError(http.StatusInternalServerError, err)
		return
	}
	resp.WriteHeader(http.StatusOK)
}

func (p *PostResource) permanentDelete(req *restful.Request, resp *restful.Response) {
	id, err := strconv.Atoi(req.PathParameter("id"))
	if err != nil {
		resp.WriteError(http.StatusBadRequest, err)
		return
	}
	if err := p.Service.PermanentDelete(uint(id)); err != nil {
		resp.WriteError(http.StatusInternalServerError, err)
		return
	}
	resp.WriteHeader(http.StatusOK)
}
