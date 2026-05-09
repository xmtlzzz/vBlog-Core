package main

import (
	"log"
	"net/http"
	"os"
	"path/filepath"

	restful "github.com/emicklei/go-restful/v3"
	"vblog-core/api"
	"vblog-core/config"
	"vblog-core/middleware"
	"vblog-core/model"
	"vblog-core/service"
)

func main() {
	cfg := config.Load()

	db, err := cfg.DB.Connect()
	if err != nil {
		log.Fatalf("database connection failed: %v", err)
	}
	if err := model.AutoMigrate(db); err != nil {
		log.Fatalf("database migration failed: %v", err)
	}

	// Services
	postSvc := service.NewPostService(db)
	tagSvc := service.NewTagService(db)
	commentSvc := service.NewCommentService(db)
	componentSvc := service.NewComponentService(db)
	settingSvc := service.NewSettingService(db)
	authSvc := service.NewAuthService(db)

	wsContainer := restful.NewContainer()
	wsContainer.EnableContentEncoding(true)

	jwtFilter := middleware.JWTFilter(cfg.JWT.Secret)

	// All API routes in one WebService
	ws := new(restful.WebService).Path("/").Produces(restful.MIME_JSON)
	// Public routes
	(&api.PostResource{Service: postSvc}).Register(ws)
	(&api.TagResource{Service: tagSvc}).Register(ws)
	(&api.CommentResource{Service: commentSvc}).Register(ws)
	(&api.SettingResource{Service: settingSvc}).Register(ws)
	(&api.AuthResource{Service: authSvc, Secret: cfg.JWT.Secret}).Register(ws)
	(&api.RSSResource{DB: db}).Register(ws)
	// Admin routes (JWT protected)
	ws.Route(ws.GET("/api/dashboard/stats").Filter(jwtFilter).To((&api.DashboardResource{DB: db}).Stats))
	ws.Route(ws.GET("/api/components").Filter(jwtFilter).To((&api.ComponentResource{Service: componentSvc}).List))
	ws.Route(ws.POST("/api/components").Filter(jwtFilter).To((&api.ComponentResource{Service: componentSvc}).Create))
	ws.Route(ws.PUT("/api/components/{id}").Filter(jwtFilter).To((&api.ComponentResource{Service: componentSvc}).Update))
	ws.Route(ws.DELETE("/api/components/{id}").Filter(jwtFilter).To((&api.ComponentResource{Service: componentSvc}).Delete))
	ws.Route(ws.PATCH("/api/components/{id}/toggle").Filter(jwtFilter).To((&api.ComponentResource{Service: componentSvc}).Toggle))
	wsContainer.Add(ws)

	// Serve frontend static files if ./static directory exists (Docker deployment)
	if _, err := os.Stat("static"); err == nil {
		staticDir, _ := filepath.Abs("static")
		wsContainer.ServeMux.Handle("/", http.FileServer(http.Dir(staticDir)))
		log.Printf("serving static files from %s", staticDir)
	}

	log.Printf("vBlog Core starting on :%s", cfg.Server.Port)
	server := &http.Server{Addr: ":" + cfg.Server.Port, Handler: wsContainer}
	log.Fatal(server.ListenAndServe())
}
