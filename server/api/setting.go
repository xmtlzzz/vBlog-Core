package api

import (
	"net/http"

	restful "github.com/emicklei/go-restful/v3"
	"vblog-core/service"
)

// SettingResource handles site settings REST endpoints.
type SettingResource struct {
	Service *service.SettingService
}

// Register adds settings routes to the given WebService.
func (s *SettingResource) Register(ws *restful.WebService) {
	ws.Route(ws.GET("/api/settings").To(s.get).
		Doc("get all settings"))

	ws.Route(ws.PUT("/api/settings").To(s.save).
		Doc("save settings"))

	ws.Route(ws.POST("/api/settings/reset").To(s.reset).
		Doc("reset settings to defaults"))
}

func (s *SettingResource) get(req *restful.Request, resp *restful.Response) {
	settings, err := s.Service.GetAll()
	if err != nil {
		resp.WriteError(http.StatusInternalServerError, err)
		return
	}
	resp.WriteEntity(settings)
}

func (s *SettingResource) save(req *restful.Request, resp *restful.Response) {
	settings := map[string]string{}
	if err := req.ReadEntity(&settings); err != nil {
		resp.WriteError(http.StatusBadRequest, err)
		return
	}
	if err := s.Service.Save(settings); err != nil {
		resp.WriteError(http.StatusInternalServerError, err)
		return
	}
	resp.WriteHeader(http.StatusOK)
}

func (s *SettingResource) reset(req *restful.Request, resp *restful.Response) {
	if err := s.Service.Reset(); err != nil {
		resp.WriteError(http.StatusInternalServerError, err)
		return
	}
	resp.WriteHeader(http.StatusOK)
}
