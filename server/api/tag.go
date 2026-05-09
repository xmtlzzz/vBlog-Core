package api

import (
	"net/http"
	"strconv"

	restful "github.com/emicklei/go-restful/v3"
	"vblog-core/model"
	"vblog-core/service"
)

// TagResource handles tag REST endpoints.
type TagResource struct {
	Service *service.TagService
}

// Register adds tag routes to the given WebService.
func (t *TagResource) Register(ws *restful.WebService) {
	ws.Route(ws.GET("/api/tags").To(t.list).
		Doc("list tags"))

	ws.Route(ws.POST("/api/tags").To(t.create).
		Doc("create a tag"))

	ws.Route(ws.PUT("/api/tags/{id}").To(t.update).
		Doc("update a tag").
		Param(ws.PathParameter("id", "tag ID")))

	ws.Route(ws.DELETE("/api/tags/{id}").To(t.delete).
		Doc("delete a tag").
		Param(ws.PathParameter("id", "tag ID")))
}

func (t *TagResource) list(req *restful.Request, resp *restful.Response) {
	tags, err := t.Service.List()
	if err != nil {
		resp.WriteError(http.StatusInternalServerError, err)
		return
	}
	resp.WriteEntity(map[string]interface{}{
		"data": tags,
	})
}

func (t *TagResource) create(req *restful.Request, resp *restful.Response) {
	tag := model.Tag{}
	if err := req.ReadEntity(&tag); err != nil {
		resp.WriteError(http.StatusBadRequest, err)
		return
	}
	if err := t.Service.Create(&tag); err != nil {
		resp.WriteError(http.StatusInternalServerError, err)
		return
	}
	resp.WriteHeaderAndEntity(http.StatusCreated, tag)
}

func (t *TagResource) update(req *restful.Request, resp *restful.Response) {
	id, err := strconv.Atoi(req.PathParameter("id"))
	if err != nil {
		resp.WriteError(http.StatusBadRequest, err)
		return
	}
	tag := model.Tag{}
	if err := req.ReadEntity(&tag); err != nil {
		resp.WriteError(http.StatusBadRequest, err)
		return
	}
	tag.ID = uint(id)
	if err := t.Service.Update(&tag); err != nil {
		resp.WriteError(http.StatusInternalServerError, err)
		return
	}
	resp.WriteEntity(tag)
}

func (t *TagResource) delete(req *restful.Request, resp *restful.Response) {
	id, err := strconv.Atoi(req.PathParameter("id"))
	if err != nil {
		resp.WriteError(http.StatusBadRequest, err)
		return
	}
	if err := t.Service.Delete(uint(id)); err != nil {
		resp.WriteError(http.StatusInternalServerError, err)
		return
	}
	resp.WriteHeader(http.StatusOK)
}
