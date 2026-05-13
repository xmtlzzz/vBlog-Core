package api

import (
	"net/http"
	"strconv"

	restful "github.com/emicklei/go-restful/v3"
	restfulspec "github.com/emicklei/go-restful-openapi/v2"
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
		Doc("List all tags").
		Notes("Returns all tags with their post counts.").
		Metadata(restfulspec.KeyOpenAPITags, []string{"tags"}).
		Writes(TagListResponse{}).
		Returns(200, "OK", TagListResponse{}).
		Returns(500, "Internal Server Error", ErrorResponse{}))

	ws.Route(ws.POST("/api/tags").To(t.create).
		Doc("Create a new tag").
		Notes("Creates a new tag. Tag name must be unique. Requires authentication.").
		Metadata(restfulspec.KeyOpenAPITags, []string{"tags"}).
		Reads(model.Tag{}).
		Writes(model.Tag{}).
		Returns(201, "Created", model.Tag{}).
		Returns(400, "Bad Request", ErrorResponse{}).
		Returns(401, "Unauthorized", ErrorResponse{}))

	ws.Route(ws.PUT("/api/tags/{id}").To(t.update).
		Doc("Update an existing tag").
		Notes("Updates a tag by ID. Requires authentication.").
		Metadata(restfulspec.KeyOpenAPITags, []string{"tags"}).
		Param(ws.PathParameter("id", "Tag ID").DataType("integer")).
		Reads(model.Tag{}).
		Writes(model.Tag{}).
		Returns(200, "OK", model.Tag{}).
		Returns(400, "Bad Request", ErrorResponse{}).
		Returns(401, "Unauthorized", ErrorResponse{}).
		Returns(404, "Not Found", ErrorResponse{}))

	ws.Route(ws.DELETE("/api/tags/{id}").To(t.delete).
		Doc("Delete a tag").
		Notes("Deletes a tag by ID. Requires authentication.").
		Metadata(restfulspec.KeyOpenAPITags, []string{"tags"}).
		Param(ws.PathParameter("id", "Tag ID").DataType("integer")).
		Returns(200, "OK", MessageResponse{}).
		Returns(401, "Unauthorized", ErrorResponse{}).
		Returns(404, "Not Found", ErrorResponse{}))
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
