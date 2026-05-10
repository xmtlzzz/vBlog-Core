package main

import (
	"log"
	"net"
	"net/http"
	"os"
	"path/filepath"

	restful "github.com/emicklei/go-restful/v3"
	"github.com/robfig/cron/v3"
	"vblog-core/api"
	"vblog-core/config"
	grpcpkg "vblog-core/grpc"
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
	dailyStatsSvc := service.NewDailyStatsService(db)
	changeLogSvc := service.NewChangeLogService(db)
	pageViewSvc := service.NewPageViewService(db)

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
	// Public stats
	ws.Route(ws.GET("/api/dashboard/stats").To((&api.DashboardResource{DB: db}).Stats))
	// Public active components
	ws.Route(ws.GET("/api/components/active").To((&api.ComponentResource{Service: componentSvc}).ListActive))
	// Admin routes (JWT protected)
	ws.Route(ws.GET("/api/components").Filter(jwtFilter).To((&api.ComponentResource{Service: componentSvc}).List))
	ws.Route(ws.POST("/api/components").Filter(jwtFilter).To((&api.ComponentResource{Service: componentSvc}).Create))
	ws.Route(ws.PUT("/api/components/{id}").Filter(jwtFilter).To((&api.ComponentResource{Service: componentSvc}).Update))
	ws.Route(ws.DELETE("/api/components/{id}").Filter(jwtFilter).To((&api.ComponentResource{Service: componentSvc}).Delete))
	ws.Route(ws.PATCH("/api/components/{id}/toggle").Filter(jwtFilter).To((&api.ComponentResource{Service: componentSvc}).Toggle))
	ws.Route(ws.POST("/api/upload").Filter(jwtFilter).To((&api.UploadResource{Dir: "static/uploads"}).Upload))
	wsContainer.Add(ws)

	// Daily snapshot cron
	c := cron.New()
	c.AddFunc("5 0 * * *", func() {
		if err := dailyStatsSvc.Snapshot(); err != nil {
			log.Printf("daily snapshot failed: %v", err)
		}
	})
	c.Start()

	// Build final handler: static files → go-restful API → SPA fallback
	handler := http.Handler(wsContainer)
	if _, err := os.Stat("static"); err == nil {
		staticDir, _ := filepath.Abs("static")
		fs := http.FileServer(http.Dir(staticDir))
		indexFile := filepath.Join(staticDir, "index.html")
		handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Let API routes go through go-restful
			if len(r.URL.Path) >= 4 && r.URL.Path[:4] == "/api" {
				wsContainer.ServeHTTP(w, r)
				return
			}
			// Try to serve static file directly
			path := filepath.Join(staticDir, r.URL.Path)
			if _, err := os.Stat(path); err == nil {
				fs.ServeHTTP(w, r)
				return
			}
			// SPA fallback
			http.ServeFile(w, r, indexFile)
		})
		log.Printf("serving static files from %s", staticDir)
	}

	// Wrap handler with PV middleware
	pvMiddleware := middleware.PVMiddleware(pageViewSvc)
	handler = pvMiddleware(handler)

	// Start gRPC server
	grpcSrv := grpcpkg.NewServer(dailyStatsSvc, changeLogSvc, pageViewSvc, settingSvc)
	grpcLis, err := net.Listen("tcp", ":"+cfg.Server.GrpcPort)
	if err != nil {
		log.Fatalf("gRPC listen failed: %v", err)
	}
	go func() {
		log.Printf("gRPC server starting on :%s", cfg.Server.GrpcPort)
		if err := grpcSrv.Start(grpcLis); err != nil {
			log.Fatalf("gRPC server failed: %v", err)
		}
	}()

	listenAddr := cfg.Server.Addr + ":" + cfg.Server.Port
	log.Printf("vBlog Core starting on %s", listenAddr)
	server := &http.Server{Addr: listenAddr, Handler: handler}
	log.Fatal(server.ListenAndServe())
}
