package api

import (
	"context"
	"net/http"
	"time"

	restful "github.com/emicklei/go-restful/v3"
	"gorm.io/gorm"
)

type HealthComponent struct {
	Status    string `json:"status"`
	LatencyMS int64  `json:"latency_ms"`
}

type HealthResponse struct {
	Status     string                     `json:"status"`
	Components map[string]HealthComponent `json:"components"`
}

type HealthResource struct {
	DB *gorm.DB
}

func (h *HealthResource) Register(ws *restful.WebService) {
	ws.Route(ws.GET("/api/health").To(h.Check).
		Doc("Check database health").
		Writes(HealthResponse{}).
		Returns(http.StatusOK, "Healthy", HealthResponse{}).
		Returns(http.StatusServiceUnavailable, "Database unavailable", HealthResponse{}))
}

func (h *HealthResource) Check(req *restful.Request, resp *restful.Response) {
	start := time.Now()
	status, code := "DOWN", http.StatusServiceUnavailable
	ctx, cancel := context.WithTimeout(req.Request.Context(), 2*time.Second)
	defer cancel()
	if h.DB != nil {
		if db, err := h.DB.DB(); err == nil {
			var value int
			if err := db.QueryRowContext(ctx, "SELECT 1").Scan(&value); err == nil && value == 1 {
				status, code = "UP", http.StatusOK
			}
		}
	}
	resp.AddHeader("Cache-Control", "no-store")
	// Keep connection details and SQL errors out of the public response.
	_ = resp.WriteHeaderAndEntity(code, HealthResponse{
		Status: status,
		Components: map[string]HealthComponent{
			"postgres": {Status: status, LatencyMS: time.Since(start).Milliseconds()},
		},
	})
}
