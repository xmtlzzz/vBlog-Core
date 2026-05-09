package api

import (
	"net/http"
	"strconv"

	restful "github.com/emicklei/go-restful/v3"
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
		Doc("list components"))

	ws.Route(ws.POST("/api/components").To(c.Create).
		Doc("create a component"))

	ws.Route(ws.PUT("/api/components/{id}").To(c.Update).
		Doc("update a component").
		Param(ws.PathParameter("id", "component ID")))

	ws.Route(ws.DELETE("/api/components/{id}").To(c.Delete).
		Doc("delete a component").
		Param(ws.PathParameter("id", "component ID")))

	ws.Route(ws.PATCH("/api/components/{id}/toggle").To(c.Toggle).
		Doc("toggle component active/inactive").
		Param(ws.PathParameter("id", "component ID")))
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
