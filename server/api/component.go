package api

import (
	"net/http"
	"strconv"

	restful "github.com/emicklei/go-restful/v3"
	restfulspec "github.com/emicklei/go-restful-openapi/v2"
	"vblog-core/model"
	"vblog-core/service"
)

// ComponentResource handles component REST endpoints.
type ComponentResource struct {
	Service *service.ComponentService
}

// Register adds component routes to the given WebService.
func (c *ComponentResource) Register(ws *restful.WebService) {
	ws.Route(ws.GET("/api/components").To(c.List).
		Doc("List all custom components").
		Notes("Returns all custom components. Requires authentication.").
		Metadata(restfulspec.KeyOpenAPITags, []string{"components"}).
		Writes(ComponentListResponse{}).
		Returns(200, "OK", ComponentListResponse{}).
		Returns(401, "Unauthorized", ErrorResponse{}).
		Returns(500, "Internal Server Error", ErrorResponse{}))

	ws.Route(ws.GET("/api/components/active").To(c.ListActive).
		Doc("List active custom components").
		Notes("Returns only active components. Public endpoint.").
		Metadata(restfulspec.KeyOpenAPITags, []string{"components"}).
		Writes([]model.Component{}).
		Returns(200, "OK", []model.Component{}))

	ws.Route(ws.POST("/api/components").To(c.Create).
		Doc("Create a new custom component").
		Notes("Creates a new component. Requires authentication.").
		Metadata(restfulspec.KeyOpenAPITags, []string{"components"}).
		Reads(model.Component{}).
		Writes(model.Component{}).
		Returns(201, "Created", model.Component{}).
		Returns(400, "Bad Request", ErrorResponse{}).
		Returns(401, "Unauthorized", ErrorResponse{}))

	ws.Route(ws.PUT("/api/components/{id}").To(c.Update).
		Doc("Update an existing component").
		Notes("Updates a component by ID. Requires authentication.").
		Metadata(restfulspec.KeyOpenAPITags, []string{"components"}).
		Param(ws.PathParameter("id", "Component ID").DataType("integer")).
		Reads(model.Component{}).
		Writes(model.Component{}).
		Returns(200, "OK", model.Component{}).
		Returns(400, "Bad Request", ErrorResponse{}).
		Returns(401, "Unauthorized", ErrorResponse{}).
		Returns(404, "Not Found", ErrorResponse{}))

	ws.Route(ws.DELETE("/api/components/{id}").To(c.Delete).
		Doc("Delete a component").
		Notes("Deletes a component by ID. Requires authentication.").
		Metadata(restfulspec.KeyOpenAPITags, []string{"components"}).
		Param(ws.PathParameter("id", "Component ID").DataType("integer")).
		Returns(200, "OK", MessageResponse{}).
		Returns(401, "Unauthorized", ErrorResponse{}).
		Returns(404, "Not Found", ErrorResponse{}))

	ws.Route(ws.PATCH("/api/components/{id}/toggle").To(c.Toggle).
		Doc("Toggle component active/inactive status").
		Notes("Toggles a component between active and inactive. Requires authentication.").
		Metadata(restfulspec.KeyOpenAPITags, []string{"components"}).
		Param(ws.PathParameter("id", "Component ID").DataType("integer")).
		Returns(200, "OK", MessageResponse{}).
		Returns(401, "Unauthorized", ErrorResponse{}).
		Returns(404, "Not Found", ErrorResponse{}))
}

func (c *ComponentResource) List(req *restful.Request, resp *restful.Response) {
	components, err := c.Service.List()
	if err != nil {
		resp.WriteError(http.StatusInternalServerError, err)
		return
	}
	resp.WriteEntity(map[string]interface{}{
		"data": components,
	})
}

func (c *ComponentResource) ListActive(req *restful.Request, resp *restful.Response) {
	components, err := c.Service.ListActive()
	if err != nil {
		resp.WriteError(http.StatusInternalServerError, err)
		return
	}
	resp.WriteEntity(components)
}

func (c *ComponentResource) Create(req *restful.Request, resp *restful.Response) {
	comp := model.Component{}
	if err := req.ReadEntity(&comp); err != nil {
		resp.WriteError(http.StatusBadRequest, err)
		return
	}
	if err := c.Service.Create(&comp); err != nil {
		resp.WriteError(http.StatusInternalServerError, err)
		return
	}
	resp.WriteHeaderAndEntity(http.StatusCreated, comp)
}

func (c *ComponentResource) Update(req *restful.Request, resp *restful.Response) {
	id, err := strconv.Atoi(req.PathParameter("id"))
	if err != nil {
		resp.WriteError(http.StatusBadRequest, err)
		return
	}
	comp := model.Component{}
	if err := req.ReadEntity(&comp); err != nil {
		resp.WriteError(http.StatusBadRequest, err)
		return
	}
	comp.ID = uint(id)
	if err := c.Service.Update(&comp); err != nil {
		resp.WriteError(http.StatusInternalServerError, err)
		return
	}
	resp.WriteEntity(comp)
}

func (c *ComponentResource) Delete(req *restful.Request, resp *restful.Response) {
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

func (c *ComponentResource) Toggle(req *restful.Request, resp *restful.Response) {
	id, err := strconv.Atoi(req.PathParameter("id"))
	if err != nil {
		resp.WriteError(http.StatusBadRequest, err)
		return
	}
	if err := c.Service.Toggle(uint(id)); err != nil {
		resp.WriteError(http.StatusInternalServerError, err)
		return
	}
	resp.WriteHeader(http.StatusOK)
}
