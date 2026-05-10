package api

import (
	"net/http"
	"strconv"

	restful "github.com/emicklei/go-restful/v3"
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
		Doc("list posts").
		Param(ws.QueryParameter("page", "page number").DefaultValue("1")).
		Param(ws.QueryParameter("per_page", "items per page").DefaultValue("5")).
		Param(ws.QueryParameter("tag", "filter by tag name")).
		Param(ws.QueryParameter("status", "filter by status")).
		Param(ws.QueryParameter("search", "search by title")))

	ws.Route(ws.GET("/api/posts/{id}").To(p.get).
		Doc("get a post").
		Param(ws.PathParameter("id", "post ID")))

	ws.Route(ws.POST("/api/posts").To(p.create).
		Doc("create a post"))

	ws.Route(ws.PUT("/api/posts/{id}").To(p.update).
		Doc("update a post").
		Param(ws.PathParameter("id", "post ID")))

	ws.Route(ws.DELETE("/api/posts/{id}").To(p.delete).
		Doc("delete a post").
		Param(ws.PathParameter("id", "post ID")))

	ws.Route(ws.GET("/api/posts/trash").To(p.trash).
		Doc("list deleted posts"))

	ws.Route(ws.POST("/api/posts/{id}/restore").To(p.restore).
		Doc("restore a deleted post").
		Param(ws.PathParameter("id", "post ID")))

	ws.Route(ws.DELETE("/api/posts/{id}/permanent").To(p.permanentDelete).
		Doc("permanently delete a post").
		Param(ws.PathParameter("id", "post ID")))
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
