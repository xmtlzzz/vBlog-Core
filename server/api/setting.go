package api

import (
	"net/http"

	restful "github.com/emicklei/go-restful/v3"
	restfulspec "github.com/emicklei/go-restful-openapi/v2"
	"vblog-core/model"
	"vblog-core/service"
)

// SettingResource handles site settings REST endpoints.
type SettingResource struct {
	Service *service.SettingService
}

// Register adds settings routes to the given WebService.
func (s *SettingResource) Register(ws *restful.WebService) {
	ws.Route(ws.GET("/api/settings").To(s.get).
		Doc("Get all site settings").
		Notes("Returns all site settings as key-value pairs.").
		Metadata(restfulspec.KeyOpenAPITags, []string{"settings"}).
		Writes([]model.Setting{}).
		Returns(200, "OK", []model.Setting{}).
		Returns(500, "Internal Server Error", ErrorResponse{}))

	ws.Route(ws.PUT("/api/settings").To(s.save).
		Doc("Save site settings").
		Notes("Saves multiple site settings. Requires authentication.").
		Metadata(restfulspec.KeyOpenAPITags, []string{"settings"}).
		Reads(map[string]string{}).
		Returns(200, "OK", MessageResponse{}).
		Returns(400, "Bad Request", ErrorResponse{}).
		Returns(401, "Unauthorized", ErrorResponse{}))

	ws.Route(ws.POST("/api/settings/reset").To(s.reset).
		Doc("Reset settings to defaults").
		Notes("Resets all settings to their default values. Requires authentication.").
		Metadata(restfulspec.KeyOpenAPITags, []string{"settings"}).
		Returns(200, "OK", MessageResponse{}).
		Returns(401, "Unauthorized", ErrorResponse{}))
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
