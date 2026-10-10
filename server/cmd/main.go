package main

import (
	"log"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	restfulspec "github.com/emicklei/go-restful-openapi/v2"
	restful "github.com/emicklei/go-restful/v3"
	"github.com/robfig/cron/v3"
	"vblog-core/api"
	"vblog-core/config"
	grpcpkg "vblog-core/grpc"
	"vblog-core/middleware"
	"vblog-core/model"
	"vblog-core/service"
	"vblog-core/webmeta"
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

	// Backfill change_log with existing content (one-time, skips if already done)
	if err := changeLogSvc.Backfill(); err != nil {
		log.Printf("change_log backfill failed: %v", err)
	}

	wsContainer := restful.NewContainer()
	wsContainer.EnableContentEncoding(true)

	jwtFilter := middleware.JWTFilter(cfg.JWT.Secret)

	// All API routes in one WebService
	ws := new(restful.WebService).Path("/").Produces(restful.MIME_JSON)
	// Public routes (admin write routes get jwtFilter via each resource's Auth field)
	(&api.PostResource{Service: postSvc, Auth: jwtFilter, ReadAuth: middleware.OptionalJWTFilter(cfg.JWT.Secret)}).Register(ws)
	(&api.TagResource{Service: tagSvc, Auth: jwtFilter}).Register(ws)
	(&api.CommentResource{Service: commentSvc, Auth: jwtFilter, TurnstileSecret: cfg.Turnstile.Secret, TurnstileHostnames: cfg.Turnstile.Hostnames}).Register(ws)
	(&api.SettingResource{Service: settingSvc, Auth: jwtFilter}).Register(ws)
	(&api.AuthResource{Service: authSvc, Secret: cfg.JWT.Secret, Auth: jwtFilter}).Register(ws)
	(&api.RSSResource{DB: db}).Register(ws)
	(&api.SitemapResource{DB: db}).Register(ws)
	// Public stats
	ws.Route(ws.GET("/api/dashboard/stats").To((&api.DashboardResource{DB: db}).Stats).
		Doc("Get dashboard statistics").
		Notes("Returns aggregate statistics including post count, total views, comments, and tags.").
		Metadata(restfulspec.KeyOpenAPITags, []string{"dashboard"}).
		Writes(api.DashboardStatsResponse{}).
		Returns(200, "OK", api.DashboardStatsResponse{}))
	// Public active components
	ws.Route(ws.GET("/api/components/active").To((&api.ComponentResource{Service: componentSvc}).ListActive).
		Doc("List active custom components").
		Notes("Returns only active components. Public endpoint.").
		Metadata(restfulspec.KeyOpenAPITags, []string{"components"}).
		Writes([]model.Component{}).
		Returns(200, "OK", []model.Component{}))
	// Admin routes (JWT protected)
	ws.Route(ws.GET("/api/components").Filter(jwtFilter).To((&api.ComponentResource{Service: componentSvc}).List).
		Doc("List all components").
		Notes("Returns all custom components. Requires authentication.").
		Metadata(restfulspec.KeyOpenAPITags, []string{"components"}).
		Writes(api.ComponentListResponse{}).
		Returns(200, "OK", api.ComponentListResponse{}).
		Returns(401, "Unauthorized", api.ErrorResponse{}))
	ws.Route(ws.POST("/api/components").Filter(jwtFilter).To((&api.ComponentResource{Service: componentSvc}).Create).
		Doc("Create a component").
		Notes("Creates a new component. Requires authentication.").
		Metadata(restfulspec.KeyOpenAPITags, []string{"components"}).
		Reads(model.Component{}).
		Writes(model.Component{}).
		Returns(201, "Created", model.Component{}).
		Returns(400, "Bad Request", api.ErrorResponse{}).
		Returns(401, "Unauthorized", api.ErrorResponse{}))
	ws.Route(ws.PUT("/api/components/{id}").Filter(jwtFilter).To((&api.ComponentResource{Service: componentSvc}).Update).
		Doc("Update a component").
		Notes("Updates a component by ID. Requires authentication.").
		Metadata(restfulspec.KeyOpenAPITags, []string{"components"}).
		Param(ws.PathParameter("id", "Component ID").DataType("integer")).
		Reads(model.Component{}).
		Writes(model.Component{}).
		Returns(200, "OK", model.Component{}).
		Returns(400, "Bad Request", api.ErrorResponse{}).
		Returns(401, "Unauthorized", api.ErrorResponse{}).
		Returns(404, "Not Found", api.ErrorResponse{}))
	ws.Route(ws.DELETE("/api/components/{id}").Filter(jwtFilter).To((&api.ComponentResource{Service: componentSvc}).Delete).
		Doc("Delete a component").
		Notes("Deletes a component by ID. Requires authentication.").
		Metadata(restfulspec.KeyOpenAPITags, []string{"components"}).
		Param(ws.PathParameter("id", "Component ID").DataType("integer")).
		Returns(200, "OK", api.MessageResponse{}).
		Returns(401, "Unauthorized", api.ErrorResponse{}).
		Returns(404, "Not Found", api.ErrorResponse{}))
	ws.Route(ws.PATCH("/api/components/{id}/toggle").Filter(jwtFilter).To((&api.ComponentResource{Service: componentSvc}).Toggle).
		Doc("Toggle component status").
		Notes("Toggles a component between active and inactive. Requires authentication.").
		Metadata(restfulspec.KeyOpenAPITags, []string{"components"}).
		Param(ws.PathParameter("id", "Component ID").DataType("integer")).
		Returns(200, "OK", api.MessageResponse{}).
		Returns(401, "Unauthorized", api.ErrorResponse{}).
		Returns(404, "Not Found", api.ErrorResponse{}))
	ws.Route(ws.POST("/api/upload").Filter(jwtFilter).To((&api.UploadResource{Dir: "static/uploads"}).Upload).
		Doc("Upload an image").
		Notes("Uploads an image file. Only image files are allowed. Requires authentication.").
		Metadata(restfulspec.KeyOpenAPITags, []string{"upload"}).
		Reads(api.File{}).
		Writes(api.UploadResponse{}).
		Returns(200, "OK", api.UploadResponse{}).
		Returns(400, "Bad Request", api.ErrorResponse{}).
		Returns(401, "Unauthorized", api.ErrorResponse{}).
		Returns(500, "Internal Server Error", api.ErrorResponse{}))
	wsContainer.Add(ws)

	// Daily snapshot cron
	c := cron.New()
	c.AddFunc("5 0 * * *", func() {
		if err := dailyStatsSvc.Snapshot(); err != nil {
			log.Printf("daily snapshot failed: %v", err)
		}
	})
	c.Start()

	// Register Swagger UI
	swaggerCfg := api.DefaultSwaggerConfig()
	swaggerCfg.Host = cfg.Server.Addr + ":" + cfg.Server.Port
	api.RegisterSwagger(wsContainer, swaggerCfg)

	// Build final handler: static files → go-restful API → SPA fallback
	handler := http.Handler(wsContainer)
	if _, err := os.Stat("static"); err == nil {
		staticDir, _ := filepath.Abs("static")
		fs := http.FileServer(http.Dir(staticDir))
		indexFile := filepath.Join(staticDir, "index.html")
		handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Let API routes and swagger go through go-restful
			if len(r.URL.Path) >= 4 && r.URL.Path[:4] == "/api" {
				wsContainer.ServeHTTP(w, r)
				return
			}
			if len(r.URL.Path) >= 8 && r.URL.Path[:8] == "/swagger" {
				wsContainer.ServeHTTP(w, r)
				return
			}
			if r.URL.Path == "/apidocs.json" {
				wsContainer.ServeHTTP(w, r)
				return
			}
			if r.URL.Path == "/sitemap.xml" || r.URL.Path == "/sitemap" ||
				r.URL.Path == "/feed.xml" || r.URL.Path == "/rss.xml" ||
				r.URL.Path == "/feed" || r.URL.Path == "/rss" {
				wsContainer.ServeHTTP(w, r)
				return
			}
			// Try to serve static file directly
			path := filepath.Join(staticDir, r.URL.Path)
			if info, err := os.Stat(path); err == nil && !info.IsDir() && r.URL.Path != "/index.html" {
				fs.ServeHTTP(w, r)
				return
			}
			// SPA fallback
			index, err := os.ReadFile(indexFile)
			if err != nil {
				http.Error(w, "页面暂不可用", http.StatusServiceUnavailable)
				return
			}
			settings, _ := settingSvc.GetAll()
			scheme := "http"
			if r.TLS != nil {
				scheme = "https"
			}
			page := webmeta.NewPage(settings, scheme+"://"+r.Host, r.URL.Path)
			if strings.HasPrefix(r.URL.Path, "/post/") {
				id, parseErr := strconv.ParseUint(strings.TrimPrefix(r.URL.Path, "/post/"), 10, 64)
				var post model.Post
				if parseErr == nil && db.Where("id = ? AND status = ?", id, "published").First(&post).Error == nil {
					page.Title, page.Article = post.Title, true
					content := post.Excerpt
					if content == "" {
						content = post.Content
					}
					page.Description = service.BuildExcerpt(content, 160)
				} else {
					page.Title, page.Missing = "文章不存在", true
				}
			}
			w.Header().Set("Content-Type", "text/html; charset=utf-8")
			w.Header().Set("Cache-Control", "no-cache")
			if page.Missing {
				w.WriteHeader(http.StatusNotFound)
			}
			if r.Method != http.MethodHead {
				_, _ = w.Write(webmeta.Render(index, page))
			}
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
