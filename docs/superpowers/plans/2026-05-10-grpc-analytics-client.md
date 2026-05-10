# gRPC Analytics + Wails Client Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Add gRPC analytics server with real-time change streaming and build a Wails desktop client for blog monitoring.

**Architecture:** Daily stats snapshots + change log table for streaming. gRPC server on :50051 with unary (GetTrends, GetLatestStats) and server-streaming (WatchChanges) RPCs. Wails client with notification/change center UI.

**Tech Stack:** Go, gRPC (google.golang.org/grpc + protobuf), GORM, PostgreSQL, Wails v2, Vue 3

---

## File Structure

### New Files (Server)

| File | Responsibility |
|------|---------------|
| `server/proto/analytics.proto` | Proto definition |
| `server/proto/analytics.pb.go` | Generated message code |
| `server/proto/analytics_grpc.pb.go` | Generated gRPC stubs |
| `server/model/daily_stats.go` | DailyStats model |
| `server/model/change_log.go` | ChangeLog model |
| `server/model/page_view.go` | PageView model |
| `server/model/daily_stats_test.go` | DailyStats model tests |
| `server/model/change_log_test.go` | ChangeLog model tests |
| `server/model/page_view_test.go` | PageView model tests |
| `server/service/daily_stats.go` | DailyStats CRUD + aggregation |
| `server/service/daily_stats_test.go` | DailyStats service tests |
| `server/service/change_log.go` | ChangeLog CRUD + write helpers |
| `server/service/change_log_test.go` | ChangeLog service tests |
| `server/service/page_view.go` | PV/UV recording + stats |
| `server/service/page_view_test.go` | PageView service tests |
| `server/middleware/pv.go` | PV recording HTTP middleware |
| `server/middleware/pv_test.go` | PV middleware tests |
| `server/grpc/server.go` | gRPC server setup + start |
| `server/grpc/analytics.go` | BlogAnalytics service impl |
| `server/grpc/analytics_test.go` | gRPC analytics tests |
| `server/grpc/auth.go` | API Key auth interceptor |
| `server/grpc/auth_test.go` | Auth interceptor tests |

### Modified Files (Server)

| File | Change |
|------|--------|
| `server/model/migrate.go` | Add DailyStats, ChangeLog, PageView to AutoMigrate |
| `server/cmd/main.go` | Start gRPC server, add PV middleware |
| `server/service/post.go` | Write change_log on Create |
| `server/service/comment.go` | Write change_log on Create |
| `server/service/tag.go` | Write change_log on tag creation |
| `server/go.mod` | Add grpc + protobuf dependencies |

### New Files (Client)

| File | Responsibility |
|------|---------------|
| `client/main.go` | Wails app entry |
| `client/app.go` | App logic: gRPC client, Wails bindings |
| `client/app_test.go` | App logic unit tests |
| `client/proto/` | Copy of generated proto code |
| `client/frontend/index.html` | HTML shell |
| `client/frontend/src/main.js` | Vue 3 entry |
| `client/frontend/src/App.vue` | Main layout |
| `client/frontend/src/components/ChangeCard.vue` | Change notification card |
| `client/frontend/src/components/StatsBar.vue` | Stats summary bar |
| `client/frontend/src/components/TrendPanel.vue` | Trend chart panel |
| `client/frontend/src/components/Settings.vue` | Connection settings |
| `client/wails.json` | Wails project config |

---

## Task 1: Proto Definition + gRPC Skeleton

**Files:**
- Create: `server/proto/analytics.proto`
- Create: `server/proto/analytics.pb.go` (generated)
- Create: `server/proto/analytics_grpc.pb.go` (generated)
- Create: `server/grpc/server.go`
- Create: `server/grpc/server_test.go`
- Modify: `server/go.mod`

- [ ] **Step 1: Install protoc and Go gRPC tools**

```bash
# Check if protoc is installed
protoc --version
# If not installed, install from https://grpc.io/docs/protoc-installation/

# Install Go plugins
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
```

- [ ] **Step 2: Add gRPC dependencies to go.mod**

Run: `cd server && go get google.golang.org/grpc@latest && go get google.golang.org/protobuf@latest`

- [ ] **Step 3: Write the proto file**

Create `server/proto/analytics.proto`:

```protobuf
syntax = "proto3";
package vblog;
option go_package = "vblog-core/proto";

service BlogAnalytics {
  rpc GetTrends(GetTrendsRequest) returns (GetTrendsResponse);
  rpc GetLatestStats(Empty) returns (LatestStats);
  rpc WatchChanges(WatchRequest) returns (stream ChangeEvent);
  rpc Ping(Empty) returns (Empty);
}

message GetTrendsRequest {
  string granularity = 1;
  int32 count = 2;
}

message TrendPoint {
  string label = 1;
  int64 pv = 2;
  int64 uv = 3;
  int64 view_total = 4;
  int64 comment_count = 5;
  int64 post_count = 6;
}

message GetTrendsResponse {
  repeated TrendPoint points = 1;
}

message LatestStats {
  int64 pv_today = 1;
  int64 uv_today = 2;
  int64 total_posts = 3;
  int64 total_views = 4;
  int64 total_comments = 5;
  int64 total_tags = 6;
  int64 pv_yesterday = 7;
  int64 uv_yesterday = 8;
  int64 views_today_delta = 9;
  int64 comments_today_delta = 10;
}

message WatchRequest {
  string api_key = 1;
  int64 since_id = 2;
}

message ChangeEvent {
  int64 id = 1;
  string type = 2;
  string title = 3;
  string detail = 4;
  string timestamp = 5;
}

message Empty {}
```

- [ ] **Step 4: Generate Go code from proto**

```bash
cd server
protoc --go_out=. --go_opt=paths=source_relative \
       --go-grpc_out=. --go-grpc_opt=paths=source_relative \
       proto/analytics.proto
```

Expected: `server/proto/analytics.pb.go` and `server/proto/analytics_grpc.pb.go` generated.

- [ ] **Step 5: Write failing test for gRPC server startup**

Create `server/grpc/server_test.go`:

```go
package grpc

import (
	"net"
	"testing"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	pb "vblog-core/proto"
)

func TestServer_StartAndPing(t *testing.T) {
	lis, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("failed to listen: %v", err)
	}

	srv := NewServer(nil, nil, nil, nil) // nil services for now
	go srv.GrpcServer.Serve(lis)
	defer srv.GrpcServer.Stop()

	conn, err := grpc.NewClient(lis.Addr().String(),
		grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		t.Fatalf("failed to connect: %v", err)
	}
	defer conn.Close()

	client := pb.NewBlogAnalyticsClient(conn)
	_, err = client.Ping(t.Context(), &pb.Empty{})
	if err != nil {
		t.Fatalf("Ping failed: %v", err)
	}
}
```

- [ ] **Step 6: Run test to verify it fails**

Run: `cd server && go test ./grpc/ -v -run TestServer_StartAndPing`
Expected: FAIL — `NewServer` not defined.

- [ ] **Step 7: Implement gRPC server skeleton**

Create `server/grpc/server.go`:

```go
package grpc

import (
	"net"

	"google.golang.org/grpc"
	pb "vblog-core/proto"
)

type Server struct {
	pb.UnimplementedBlogAnalyticsServer
	GrpcServer *grpc.Server
}

func NewServer() *Server {
	s := &Server{
		GrpcServer: grpc.NewServer(),
	}
	pb.RegisterBlogAnalyticsServer(s.GrpcServer, s)
	return s
}

func (s *Server) Start(lis net.Listener) error {
	return s.GrpcServer.Serve(lis)
}

func (s *Server) Ping(ctx context.Context, in *pb.Empty) (*pb.Empty, error) {
	return &pb.Empty{}, nil
}
```

- [ ] **Step 8: Run test to verify it passes**

Run: `cd server && go test ./grpc/ -v -run TestServer_StartAndPing`
Expected: PASS

- [ ] **Step 9: Commit**

```bash
git add server/proto/ server/grpc/server.go server/grpc/server_test.go server/go.mod server/go.sum
git commit -m "feat: add proto definition and gRPC server skeleton with Ping RPC"
```

---

## Task 2: Database Models + Migrations

**Files:**
- Create: `server/model/daily_stats.go`
- Create: `server/model/daily_stats_test.go`
- Create: `server/model/change_log.go`
- Create: `server/model/change_log_test.go`
- Create: `server/model/page_view.go`
- Create: `server/model/page_view_test.go`
- Modify: `server/model/migrate.go`

- [ ] **Step 1: Write failing test for DailyStats model**

Create `server/model/daily_stats_test.go`:

```go
package model

import "testing"

func TestDailyStatsTableName(t *testing.T) {
	d := DailyStats{}
	if d.TableName() != "daily_stats" {
		t.Errorf("expected 'daily_stats', got '%s'", d.TableName())
	}
}

func TestDailyStatsDefaultValues(t *testing.T) {
	d := DailyStats{}
	if d.PV != 0 {
		t.Errorf("expected PV 0, got %d", d.PV)
	}
	if d.UV != 0 {
		t.Errorf("expected UV 0, got %d", d.UV)
	}
}
```

- [ ] **Step 2: Run test to verify it fails**

Run: `cd server && go test ./model/ -v -run TestDailyStats`
Expected: FAIL — `DailyStats` not defined.

- [ ] **Step 3: Implement DailyStats model**

Create `server/model/daily_stats.go`:

```go
package model

import "time"

type DailyStats struct {
	ID           uint      `gorm:"primaryKey" json:"id"`
	StatDate     time.Time `gorm:"type:date;uniqueIndex;not null" json:"stat_date"`
	PV           int64     `gorm:"default:0" json:"pv"`
	UV           int64     `gorm:"default:0" json:"uv"`
	PostCount    int       `gorm:"default:0" json:"post_count"`
	ViewTotal    int64     `gorm:"default:0" json:"view_total"`
	CommentCount int       `gorm:"default:0" json:"comment_count"`
	TagCount     int       `gorm:"default:0" json:"tag_count"`
	CreatedAt    time.Time `json:"created_at"`
}

func (DailyStats) TableName() string {
	return "daily_stats"
}
```

- [ ] **Step 4: Run test to verify it passes**

Run: `cd server && go test ./model/ -v -run TestDailyStats`
Expected: PASS

- [ ] **Step 5: Write failing test for ChangeLog model**

Create `server/model/change_log_test.go`:

```go
package model

import "testing"

func TestChangeLogTableName(t *testing.T) {
	c := ChangeLog{}
	if c.TableName() != "change_log" {
		t.Errorf("expected 'change_log', got '%s'", c.TableName())
	}
}

func TestChangeLogDefaultValues(t *testing.T) {
	c := ChangeLog{}
	if c.ChangeType != "" {
		t.Errorf("expected empty change_type, got '%s'", c.ChangeType)
	}
}
```

- [ ] **Step 6: Run test to verify it fails**

Run: `cd server && go test ./model/ -v -run TestChangeLog`
Expected: FAIL — `ChangeLog` not defined.

- [ ] **Step 7: Implement ChangeLog model**

Create `server/model/change_log.go`:

```go
package model

import "time"

type ChangeLog struct {
	ID         uint      `gorm:"primaryKey" json:"id"`
	ChangeType string    `gorm:"size:50;not null" json:"change_type"`
	TargetID   *uint     `json:"target_id"`
	Title      string    `gorm:"size:200" json:"title"`
	Detail     string    `gorm:"type:text" json:"detail"`
	CreatedAt  time.Time `json:"created_at"`
}

func (ChangeLog) TableName() string {
	return "change_log"
}
```

- [ ] **Step 8: Run test to verify it passes**

Run: `cd server && go test ./model/ -v -run TestChangeLog`
Expected: PASS

- [ ] **Step 9: Write failing test for PageView model**

Create `server/model/page_view_test.go`:

```go
package model

import "testing"

func TestPageViewTableName(t *testing.T) {
	p := PageView{}
	if p.TableName() != "page_views" {
		t.Errorf("expected 'page_views', got '%s'", p.TableName())
	}
}
```

- [ ] **Step 10: Run test to verify it fails**

Run: `cd server && go test ./model/ -v -run TestPageView`
Expected: FAIL — `PageView` not defined.

- [ ] **Step 11: Implement PageView model**

Create `server/model/page_view.go`:

```go
package model

import "time"

type PageView struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	IP        string    `gorm:"size:45" json:"ip"`
	Path      string    `gorm:"size:500" json:"path"`
	UserAgent string    `gorm:"type:text" json:"user_agent"`
	CreatedAt time.Time `json:"created_at"`
}

func (PageView) TableName() string {
	return "page_views"
}
```

- [ ] **Step 12: Run test to verify it passes**

Run: `cd server && go test ./model/ -v -run TestPageView`
Expected: PASS

- [ ] **Step 13: Add models to AutoMigrate**

Modify `server/model/migrate.go`:

```go
package model

import "gorm.io/gorm"

func AutoMigrate(db *gorm.DB) error {
	return db.AutoMigrate(
		&User{},
		&Tag{},
		&Post{},
		&Comment{},
		&Component{},
		&Setting{},
		&DailyStats{},
		&ChangeLog{},
		&PageView{},
	)
}
```

- [ ] **Step 14: Run all model tests**

Run: `cd server && go test ./model/ -v`
Expected: All PASS

- [ ] **Step 15: Commit**

```bash
git add server/model/daily_stats.go server/model/daily_stats_test.go \
        server/model/change_log.go server/model/change_log_test.go \
        server/model/page_view.go server/model/page_view_test.go \
        server/model/migrate.go
git commit -m "feat: add DailyStats, ChangeLog, PageView models with tests"
```

---

## Task 3: PageView Service + PV Middleware

**Files:**
- Create: `server/service/page_view.go`
- Create: `server/service/page_view_test.go`
- Create: `server/middleware/pv.go`
- Create: `server/middleware/pv_test.go`

- [ ] **Step 1: Write failing test for PageView service Record**

Create `server/service/page_view_test.go`:

```go
package service

import (
	"testing"
	"vblog-core/model"
	"vblog-core/testutil"
)

func TestPageViewService_Record(t *testing.T) {
	db := testutil.GetTestDB(t)
	svc := NewPageViewService(db)

	err := svc.Record("192.168.1.1", "/posts/1", "Mozilla/5.0")
	if err != nil {
		t.Fatalf("Record failed: %v", err)
	}

	var count int64
	db.Model(&model.PageView{}).Where("ip = ? AND path = ?", "192.168.1.1", "/posts/1").Count(&count)
	if count != 1 {
		t.Errorf("expected 1 page view, got %d", count)
	}

	db.Where("ip = ?", "192.168.1.1").Delete(&model.PageView{})
}

func TestPageViewService_GetPVUV(t *testing.T) {
	db := testutil.GetTestDB(t)
	svc := NewPageViewService(db)

	// Insert test data
	svc.Record("10.0.0.1", "/posts/1", "ua")
	svc.Record("10.0.0.2", "/posts/1", "ua")
	svc.Record("10.0.0.1", "/posts/2", "ua")

	pv, uv, err := svc.GetPVUVToday()
	if err != nil {
		t.Fatalf("GetPVUVToday failed: %v", err)
	}
	if pv < 3 {
		t.Errorf("expected PV >= 3, got %d", pv)
	}
	if uv < 2 {
		t.Errorf("expected UV >= 2, got %d", uv)
	}

	db.Where("ip LIKE ?", "10.0.0.%").Delete(&model.PageView{})
}
```

- [ ] **Step 2: Run test to verify it fails**

Run: `cd server && go test ./service/ -v -run TestPageViewService`
Expected: FAIL — `NewPageViewService` not defined.

- [ ] **Step 3: Implement PageView service**

Create `server/service/page_view.go`:

```go
package service

import (
	"time"

	"gorm.io/gorm"
	"vblog-core/model"
)

type PageViewService struct {
	DB *gorm.DB
}

func NewPageViewService(db *gorm.DB) *PageViewService {
	return &PageViewService{DB: db}
}

func (s *PageViewService) Record(ip, path, userAgent string) error {
	return s.DB.Create(&model.PageView{
		IP:        ip,
		Path:      path,
		UserAgent: userAgent,
	}).Error
}

func (s *PageViewService) GetPVUVToday() (pv int64, uv int64, err error) {
	today := time.Now().Format("2006-01-02")
	err = s.DB.Model(&model.PageView{}).
		Where("DATE(created_at) = ?", today).
		Count(&pv).Error
	if err != nil {
		return
	}
	err = s.DB.Model(&model.PageView{}).
		Where("DATE(created_at) = ?", today).
		Distinct("ip").
		Count(&uv).Error
	return
}

func (s *PageViewService) GetPVUVByDate(date string) (pv int64, uv int64, err error) {
	err = s.DB.Model(&model.PageView{}).
		Where("DATE(created_at) = ?", date).
		Count(&pv).Error
	if err != nil {
		return
	}
	err = s.DB.Model(&model.PageView{}).
		Where("DATE(created_at) = ?", date).
		Distinct("ip").
		Count(&uv).Error
	return
}
```

- [ ] **Step 4: Run test to verify it passes**

Run: `cd server && go test ./service/ -v -run TestPageViewService`
Expected: PASS

- [ ] **Step 5: Write failing test for PV middleware**

Create `server/middleware/pv_test.go`:

```go
package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestPVMiddleware_RecordsView(t *testing.T) {
	called := false
	inner := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		called = true
		w.WriteHeader(http.StatusOK)
	})

	handler := PVMiddleware(nil)(inner) // nil service for now, just test wrapping
	req := httptest.NewRequest("GET", "/posts/1", nil)
	req.RemoteAddr = "127.0.0.1:12345"
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	if !called {
		t.Error("expected inner handler to be called")
	}
	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", w.Code)
	}
}
```

- [ ] **Step 6: Run test to verify it fails**

Run: `cd server && go test ./middleware/ -v -run TestPVMiddleware`
Expected: FAIL — `PVMiddleware` not defined.

- [ ] **Step 7: Implement PV middleware**

Create `server/middleware/pv.go`:

```go
package middleware

import (
	"net/http"
	"strings"

	"vblog-core/service"
)

func PVMiddleware(pvSvc *service.PageViewService) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if pvSvc != nil && !strings.HasPrefix(r.URL.Path, "/api") {
				ip := r.RemoteAddr
				if forwarded := r.Header.Get("X-Forwarded-For"); forwarded != "" {
					ip = strings.Split(forwarded, ",")[0]
				}
				go pvSvc.Record(strings.TrimSpace(ip), r.URL.Path, r.UserAgent())
			}
			next.ServeHTTP(w, r)
		})
	}
}
```

- [ ] **Step 8: Run test to verify it passes**

Run: `cd server && go test ./middleware/ -v -run TestPVMiddleware`
Expected: PASS

- [ ] **Step 9: Commit**

```bash
git add server/service/page_view.go server/service/page_view_test.go \
        server/middleware/pv.go server/middleware/pv_test.go
git commit -m "feat: add PageView service and PV recording middleware"
```

---

## Task 4: ChangeLog Service

**Files:**
- Create: `server/service/change_log.go`
- Create: `server/service/change_log_test.go`

- [ ] **Step 1: Write failing test for ChangeLog service Write**

Create `server/service/change_log_test.go`:

```go
package service

import (
	"testing"
	"vblog-core/model"
	"vblog-core/testutil"
)

func TestChangeLogService_Write(t *testing.T) {
	db := testutil.GetTestDB(t)
	svc := NewChangeLogService(db)

	err := svc.Write("new_post", nil, "Test Post", `{"id":1}`)
	if err != nil {
		t.Fatalf("Write failed: %v", err)
	}

	var log model.ChangeLog
	err = db.Order("id DESC").First(&log).Error
	if err != nil {
		t.Fatalf("query failed: %v", err)
	}
	if log.ChangeType != "new_post" {
		t.Errorf("expected 'new_post', got '%s'", log.ChangeType)
	}
	if log.Title != "Test Post" {
		t.Errorf("expected 'Test Post', got '%s'", log.Title)
	}

	db.Where("change_type = ?", "new_post").Delete(&model.ChangeLog{})
}

func TestChangeLogService_GetAfterID(t *testing.T) {
	db := testutil.GetTestDB(t)
	svc := NewChangeLogService(db)

	svc.Write("new_post", nil, "Post 1", "")
	svc.Write("new_comment", nil, "Comment 1", "")

	var last model.ChangeLog
	db.Order("id DESC").First(&last)

	logs, err := svc.GetAfterID(last.ID - 1)
	if err != nil {
		t.Fatalf("GetAfterID failed: %v", err)
	}
	if len(logs) < 1 {
		t.Errorf("expected >= 1 log, got %d", len(logs))
	}

	db.Where("1 = 1").Delete(&model.ChangeLog{})
}

func TestChangeLogService_GetLatestID(t *testing.T) {
	db := testutil.GetTestDB(t)
	svc := NewChangeLogService(db)

	svc.Write("test_type", nil, "test", "")

	id, err := svc.GetLatestID()
	if err != nil {
		t.Fatalf("GetLatestID failed: %v", err)
	}
	if id == 0 {
		t.Error("expected non-zero ID")
	}

	db.Where("change_type = ?", "test_type").Delete(&model.ChangeLog{})
}
```

- [ ] **Step 2: Run test to verify it fails**

Run: `cd server && go test ./service/ -v -run TestChangeLogService`
Expected: FAIL — `NewChangeLogService` not defined.

- [ ] **Step 3: Implement ChangeLog service**

Create `server/service/change_log.go`:

```go
package service

import (
	"gorm.io/gorm"
	"vblog-core/model"
)

type ChangeLogService struct {
	DB *gorm.DB
}

func NewChangeLogService(db *gorm.DB) *ChangeLogService {
	return &ChangeLogService{DB: db}
}

func (s *ChangeLogService) Write(changeType string, targetID *uint, title, detail string) error {
	return s.DB.Create(&model.ChangeLog{
		ChangeType: changeType,
		TargetID:   targetID,
		Title:      title,
		Detail:     detail,
	}).Error
}

func (s *ChangeLogService) GetAfterID(afterID int64) ([]model.ChangeLog, error) {
	var logs []model.ChangeLog
	err := s.DB.Where("id > ?", afterID).Order("id ASC").Find(&logs).Error
	return logs, err
}

func (s *ChangeLogService) GetLatestID() (int64, error) {
	var log model.ChangeLog
	err := s.DB.Order("id DESC").First(&log).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return 0, nil
		}
		return 0, err
	}
	return int64(log.ID), nil
}
```

- [ ] **Step 4: Run test to verify it passes**

Run: `cd server && go test ./service/ -v -run TestChangeLogService`
Expected: PASS

- [ ] **Step 5: Commit**

```bash
git add server/service/change_log.go server/service/change_log_test.go
git commit -m "feat: add ChangeLog service with Write, GetAfterID, GetLatestID"
```

---

## Task 5: Integrate Change Log into Post/Comment/Tag Services

**Files:**
- Modify: `server/service/post.go`
- Modify: `server/service/comment.go`
- Modify: `server/service/tag.go`

- [ ] **Step 1: Write failing test for Post create with change log**

Add to `server/service/post_test.go` (create if not exists):

```go
package service

import (
	"testing"
	"vblog-core/model"
	"vblog-core/testutil"
)

func TestPostService_Create_WritesChangeLog(t *testing.T) {
	db := testutil.GetTestDB(t)
	postSvc := NewPostService(db)
	logSvc := NewChangeLogService(db)

	beforeID, _ := logSvc.GetLatestID()

	post := &model.Post{
		Title:   "Test Change Log Post",
		Content: "Some content",
		Status:  "published",
	}
	err := postSvc.Create(post)
	if err != nil {
		t.Fatalf("Create failed: %v", err)
	}

	logs, _ := logSvc.GetAfterID(beforeID)
	found := false
	for _, l := range logs {
		if l.ChangeType == "new_post" && l.Title == "Test Change Log Post" {
			found = true
			break
		}
	}
	if !found {
		t.Error("expected change_log entry for new_post")
	}

	db.Unscoped().Delete(post)
	db.Where("change_type = ? AND title = ?", "new_post", "Test Change Log Post").Delete(&model.ChangeLog{})
}
```

- [ ] **Step 2: Run test to verify it fails**

Run: `cd server && go test ./service/ -v -run TestPostService_Create_WritesChangeLog`
Expected: FAIL — post Create doesn't write change_log.

- [ ] **Step 3: Add ChangeLogService dependency to PostService**

Modify `server/service/post.go` — change PostService struct and NewPostService:

```go
type PostService struct {
	DB      *gorm.DB
	LogSvc  *ChangeLogService
}

func NewPostService(db *gorm.DB) *PostService {
	return &PostService{DB: db, LogSvc: NewChangeLogService(db)}
}
```

Add change log write at end of Create method, after `s.DB.Create(post).Error`:

```go
func (s *PostService) Create(post *model.Post) error {
	post.ReadTime = CalcReadTime(post.Content)
	if post.Excerpt == "" {
		post.Excerpt = BuildExcerpt(post.Content, 200)
	}
	if len(post.Tags) > 0 {
		resolved, err := s.resolveTags(post.Tags)
		if err != nil {
			return err
		}
		post.Tags = resolved
	}
	if err := s.DB.Create(post).Error; err != nil {
		return err
	}
	if post.Status == "published" {
		s.LogSvc.Write("new_post", &post.ID, post.Title, "")
	}
	return nil
}
```

- [ ] **Step 4: Run test to verify it passes**

Run: `cd server && go test ./service/ -v -run TestPostService_Create_WritesChangeLog`
Expected: PASS

- [ ] **Step 5: Write failing test for Comment create with change log**

Add to `server/service/comment_test.go`:

```go
func TestCommentService_Create_WritesChangeLog(t *testing.T) {
	db := testutil.GetTestDB(t)
	commentSvc := NewCommentService(db)
	logSvc := NewChangeLogService(db)

	beforeID, _ := logSvc.GetLatestID()

	comment := &model.Comment{
		PostID: 1,
		Author: "tester",
		Body:   "Test comment",
	}
	err := commentSvc.Create(comment)
	if err != nil {
		t.Fatalf("Create failed: %v", err)
	}

	logs, _ := logSvc.GetAfterID(beforeID)
	found := false
	for _, l := range logs {
		if l.ChangeType == "new_comment" {
			found = true
			break
		}
	}
	if !found {
		t.Error("expected change_log entry for new_comment")
	}

	db.Unscoped().Delete(comment)
	db.Where("change_type = ?", "new_comment").Delete(&model.ChangeLog{})
}
```

- [ ] **Step 6: Run test to verify it fails**

Run: `cd server && go test ./service/ -v -run TestCommentService_Create_WritesChangeLog`
Expected: FAIL

- [ ] **Step 7: Add ChangeLogService to CommentService**

Modify `server/service/comment.go`:

```go
type CommentService struct {
	DB     *gorm.DB
	LogSvc *ChangeLogService
}

func NewCommentService(db *gorm.DB) *CommentService {
	return &CommentService{DB: db, LogSvc: NewChangeLogService(db)}
}

func (s *CommentService) Create(c *model.Comment) error {
	c.Status = "pending"
	if err := s.DB.Create(c).Error; err != nil {
		return err
	}
	s.LogSvc.Write("new_comment", &c.ID, c.Body, "")
	return nil
}
```

- [ ] **Step 8: Run test to verify it passes**

Run: `cd server && go test ./service/ -v -run TestCommentService_Create_WritesChangeLog`
Expected: PASS

- [ ] **Step 9: Commit**

```bash
git add server/service/post.go server/service/comment.go
git commit -m "feat: integrate change log into post and comment creation"
```

---

## Task 6: DailyStats Service + Snapshot Cron

**Files:**
- Create: `server/service/daily_stats.go`
- Create: `server/service/daily_stats_test.go`
- Modify: `server/cmd/main.go` (add cron)

- [ ] **Step 1: Write failing test for DailyStats Snapshot**

Create `server/service/daily_stats_test.go`:

```go
package service

import (
	"testing"
	"time"
	"vblog-core/model"
	"vblog-core/testutil"
)

func TestDailyStatsService_Snapshot(t *testing.T) {
	db := testutil.GetTestDB(t)
	svc := NewDailyStatsService(db)

	err := svc.Snapshot()
	if err != nil {
		t.Fatalf("Snapshot failed: %v", err)
	}

	today := time.Now().Format("2006-01-02")
	var stats model.DailyStats
	err = db.Where("stat_date = ?", today).First(&stats).Error
	if err != nil {
		t.Fatalf("query failed: %v", err)
	}

	// Cleanup
	db.Where("stat_date = ?", today).Delete(&model.DailyStats{})
}

func TestDailyStatsService_GetTrends(t *testing.T) {
	db := testutil.GetTestDB(t)
	svc := NewDailyStatsService(db)

	// Insert test data
	now := time.Now()
	for i := 0; i < 5; i++ {
		date := now.AddDate(0, 0, -i)
		db.Create(&model.DailyStats{
			StatDate:  date,
			PV:        int64(100 + i),
			UV:        int64(50 + i),
			ViewTotal: int64(1000 + i*10),
		})
	}

	points, err := svc.GetTrends("day", 5)
	if err != nil {
		t.Fatalf("GetTrends failed: %v", err)
	}
	if len(points) < 1 {
		t.Errorf("expected >= 1 point, got %d", len(points))
	}

	// Cleanup
	db.Where("pv >= ?", 100).Delete(&model.DailyStats{})
}
```

- [ ] **Step 2: Run test to verify it fails**

Run: `cd server && go test ./service/ -v -run TestDailyStatsService`
Expected: FAIL — `NewDailyStatsService` not defined.

- [ ] **Step 3: Implement DailyStats service**

Create `server/service/daily_stats.go`:

```go
package service

import (
	"time"

	"gorm.io/gorm"
	"vblog-core/model"
)

type DailyStatsService struct {
	DB *gorm.DB
}

func NewDailyStatsService(db *gorm.DB) *DailyStatsService {
	return &DailyStatsService{DB: db}
}

func (s *DailyStatsService) Snapshot() error {
	today := time.Now().Format("2006-01-02")

	var postCount int64
	s.DB.Model(&model.Post{}).Where("status = ?", "published").Count(&postCount)

	var viewTotal int64
	s.DB.Model(&model.Post{}).Select("COALESCE(SUM(views), 0)").Scan(&viewTotal)

	var commentCount int64
	s.DB.Model(&model.Comment{}).Count(&commentCount)

	var tagCount int64
	s.DB.Model(&model.Tag{}).Count(&tagCount)

	var pvToday, uvToday int64
	s.DB.Model(&model.PageView{}).Where("DATE(created_at) = ?", today).Count(&pvToday)
	s.DB.Model(&model.PageView{}).Where("DATE(created_at) = ?", today).Distinct("ip").Count(&uvToday)

	return s.DB.Where("stat_date = ?", today).
		Assign(model.DailyStats{
			PV:           pvToday,
			UV:           uvToday,
			PostCount:    int(postCount),
			ViewTotal:    viewTotal,
			CommentCount: int(commentCount),
			TagCount:     int(tagCount),
		}).FirstOrCreate(&model.DailyStats{StatDate: time.Now()}).Error
}

func (s *DailyStatsService) GetTrends(granularity string, count int) ([]model.DailyStats, error) {
	var stats []model.DailyStats
	query := s.DB.Order("stat_date DESC").Limit(count)

	switch granularity {
	case "week":
		// Group by week — return raw daily, client groups
		err := query.Find(&stats).Error
		return stats, err
	case "month":
		err := query.Find(&stats).Error
		return stats, err
	default: // "day"
		err := query.Find(&stats).Error
		return stats, err
	}
}
```

- [ ] **Step 4: Run test to verify it passes**

Run: `cd server && go test ./service/ -v -run TestDailyStatsService`
Expected: PASS

- [ ] **Step 5: Add daily snapshot cron to main.go**

Modify `server/cmd/main.go` — add import and cron setup after services:

```go
import (
	// ... existing imports ...
	"github.com/robfig/cron/v3"
)

// After services are created:
dailyStatsSvc := service.NewDailyStatsService(db)

// Cron: daily snapshot at 00:05
c := cron.New()
c.AddFunc("5 0 * * *", func() {
	if err := dailyStatsSvc.Snapshot(); err != nil {
		log.Printf("daily snapshot failed: %v", err)
	}
})
c.Start()
```

Also add `go get github.com/robfig/cron/v3` dependency.

- [ ] **Step 6: Run all service tests**

Run: `cd server && go test ./service/ -v`
Expected: All PASS

- [ ] **Step 7: Commit**

```bash
git add server/service/daily_stats.go server/service/daily_stats_test.go server/cmd/main.go server/go.mod server/go.sum
git commit -m "feat: add DailyStats service with snapshot cron"
```

---

## Task 7: gRPC Analytics Service Implementation

**Files:**
- Modify: `server/grpc/server.go`
- Modify: `server/grpc/analytics.go` (rename from server.go or split)
- Create: `server/grpc/analytics_test.go`
- Modify: `server/grpc/server_test.go`

- [ ] **Step 1: Update Server struct to hold service dependencies**

Modify `server/grpc/server.go`:

```go
package grpc

import (
	"context"
	"net"

	"google.golang.org/grpc"
	"vblog-core/service"
	pb "vblog-core/proto"
)

type Server struct {
	pb.UnimplementedBlogAnalyticsServer
	GrpcServer   *grpc.Server
	DailyStatsSvc *service.DailyStatsService
	ChangeLogSvc  *service.ChangeLogService
	PageViewSvc   *service.PageViewService
	SettingSvc    *service.SettingService
}

func NewServer(ds *service.DailyStatsService, cl *service.ChangeLogService, pv *service.PageViewService, st *service.SettingService) *Server {
	s := &Server{
		GrpcServer:    grpc.NewServer(),
		DailyStatsSvc: ds,
		ChangeLogSvc:  cl,
		PageViewSvc:   pv,
		SettingSvc:    st,
	}
	pb.RegisterBlogAnalyticsServer(s.GrpcServer, s)
	return s
}

func (s *Server) Start(lis net.Listener) error {
	return s.GrpcServer.Serve(lis)
}

func (s *Server) Ping(ctx context.Context, in *pb.Empty) (*pb.Empty, error) {
	return &pb.Empty{}, nil
}
```

- [ ] **Step 2: Write failing test for GetLatestStats**

Add to `server/grpc/analytics_test.go`:

```go
package grpc

import (
	"net"
	"testing"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"vblog-core/testutil"
	"vblog-core/service"
	pb "vblog-core/proto"
)

func TestGetLatestStats(t *testing.T) {
	db := testutil.GetTestDB(t)
	ds := service.NewDailyStatsService(db)
	cl := service.NewChangeLogService(db)
	pv := service.NewPageViewService(db)
	st := service.NewSettingService(db)

	lis, _ := net.Listen("tcp", "127.0.0.1:0")
	srv := NewServer(ds, cl, pv, st)
	go srv.GrpcServer.Serve(lis)
	defer srv.GrpcServer.Stop()

	conn, _ := grpc.NewClient(lis.Addr().String(),
		grpc.WithTransportCredentials(insecure.NewCredentials()))
	defer conn.Close()

	client := pb.NewBlogAnalyticsClient(conn)
	stats, err := client.GetLatestStats(t.Context(), &pb.Empty{})
	if err != nil {
		t.Fatalf("GetLatestStats failed: %v", err)
	}
	if stats == nil {
		t.Fatal("expected non-nil stats")
	}
}
```

- [ ] **Step 3: Run test to verify it fails**

Run: `cd server && go test ./grpc/ -v -run TestGetLatestStats`
Expected: FAIL — `GetLatestStats` not implemented.

- [ ] **Step 4: Implement GetLatestStats**

Add to `server/grpc/server.go` (or create `server/grpc/analytics.go`):

```go
func (s *Server) GetLatestStats(ctx context.Context, in *pb.Empty) (*pb.LatestStats, error) {
	today := time.Now().Format("2006-01-02")
	yesterday := time.Now().AddDate(0, 0, -1).Format("2006-01-02")

	pvToday, uvToday, _ := s.PageViewSvc.GetPVUVByDate(today)
	pvYesterday, uvYesterday, _ := s.PageViewSvc.GetPVUVByDate(yesterday)

	var postCount, viewTotal, commentCount, tagCount int64
	s.DailyStatsSvc.DB.Model(&model.Post{}).Where("status = ?", "published").Count(&postCount)
	s.DailyStatsSvc.DB.Model(&model.Post{}).Select("COALESCE(SUM(views), 0)").Scan(&viewTotal)
	s.DailyStatsSvc.DB.Model(&model.Comment{}).Count(&commentCount)
	s.DailyStatsSvc.DB.Model(&model.Tag{}).Count(&tagCount)

	return &pb.LatestStats{
		PvToday:    pvToday,
		UvToday:    uvToday,
		TotalPosts: postCount,
		TotalViews: viewTotal,
		TotalComments: commentCount,
		TotalTags:  tagCount,
		PvYesterday: pvYesterday,
		UvYesterday: uvYesterday,
	}, nil
}
```

- [ ] **Step 5: Run test to verify it passes**

Run: `cd server && go test ./grpc/ -v -run TestGetLatestStats`
Expected: PASS

- [ ] **Step 6: Write failing test for GetTrends**

```go
func TestGetTrends(t *testing.T) {
	db := testutil.GetTestDB(t)
	ds := service.NewDailyStatsService(db)
	cl := service.NewChangeLogService(db)
	pv := service.NewPageViewService(db)
	st := service.NewSettingService(db)

	lis, _ := net.Listen("tcp", "127.0.0.1:0")
	srv := NewServer(ds, cl, pv, st)
	go srv.GrpcServer.Serve(lis)
	defer srv.GrpcServer.Stop()

	conn, _ := grpc.NewClient(lis.Addr().String(),
		grpc.WithTransportCredentials(insecure.NewCredentials()))
	defer conn.Close()

	client := pb.NewBlogAnalyticsClient(conn)
	resp, err := client.GetTrends(t.Context(), &pb.GetTrendsRequest{
		Granularity: "day",
		Count:       7,
	})
	if err != nil {
		t.Fatalf("GetTrends failed: %v", err)
	}
	if resp == nil {
		t.Fatal("expected non-nil response")
	}
}
```

- [ ] **Step 7: Run test to verify it fails**

Run: `cd server && go test ./grpc/ -v -run TestGetTrends`
Expected: FAIL — `GetTrends` not implemented.

- [ ] **Step 8: Implement GetTrends**

```go
func (s *Server) GetTrends(ctx context.Context, in *pb.GetTrendsRequest) (*pb.GetTrendsResponse, error) {
	count := int(in.Count)
	if count <= 0 {
		count = 7
	}

	stats, err := s.DailyStatsSvc.GetTrends(in.Granularity, count)
	if err != nil {
		return nil, err
	}

	points := make([]*pb.TrendPoint, len(stats))
	for i, st := range stats {
		points[i] = &pb.TrendPoint{
			Label:        st.StatDate.Format("2006-01-02"),
			Pv:           st.PV,
			Uv:           st.UV,
			ViewTotal:    st.ViewTotal,
			CommentCount: int64(st.CommentCount),
			PostCount:    int64(st.PostCount),
		}
	}

	return &pb.GetTrendsResponse{Points: points}, nil
}
```

- [ ] **Step 9: Run test to verify it passes**

Run: `cd server && go test ./grpc/ -v -run TestGetTrends`
Expected: PASS

- [ ] **Step 10: Write failing test for WatchChanges**

```go
func TestWatchChanges(t *testing.T) {
	db := testutil.GetTestDB(t)
	ds := service.NewDailyStatsService(db)
	cl := service.NewChangeLogService(db)
	pv := service.NewPageViewService(db)
	st := service.NewSettingService(db)

	// Set API key for auth
	st.Set("grpc_api_key", "test-key-123")

	lis, _ := net.Listen("tcp", "127.0.0.1:0")
	srv := NewServer(ds, cl, pv, st)
	go srv.GrpcServer.Serve(lis)
	defer srv.GrpcServer.Stop()

	conn, _ := grpc.NewClient(lis.Addr().String(),
		grpc.WithTransportCredentials(insecure.NewCredentials()))
	defer conn.Close()

	client := pb.NewBlogAnalyticsClient(conn)
	stream, err := client.WatchChanges(t.Context(), &pb.WatchRequest{
		ApiKey: "test-key-123",
	})
	if err != nil {
		t.Fatalf("WatchChanges failed: %v", err)
	}

	// Write a change while watching
	cl.Write("new_post", nil, "Stream Test Post", "")

	// Receive event
	event, err := stream.Recv()
	if err != nil {
		t.Fatalf("stream.Recv failed: %v", err)
	}
	if event.Type != "new_post" {
		t.Errorf("expected 'new_post', got '%s'", event.Type)
	}

	// Cleanup
	db.Where("title = ?", "Stream Test Post").Delete(&model.ChangeLog{})
	db.Where("key = ?", "grpc_api_key").Delete(&model.Setting{})
}
```

- [ ] **Step 11: Run test to verify it fails**

Run: `cd server && go test ./grpc/ -v -run TestWatchChanges`
Expected: FAIL — `WatchChanges` not implemented.

- [ ] **Step 12: Implement WatchChanges**

```go
func (s *Server) WatchChanges(in *pb.WatchRequest, stream pb.BlogAnalytics_WatchChangesServer) error {
	// Validate API key
	apiKey, _ := s.SettingSvc.Get("grpc_api_key")
	if apiKey == "" || apiKey != in.ApiKey {
		return status.Error(codes.Unauthenticated, "invalid api_key")
	}

	// Catch-up: send missed changes
	sinceID := in.SinceId
	logs, _ := s.ChangeLogSvc.GetAfterID(sinceID)
	for _, l := range logs {
		event := &pb.ChangeEvent{
			Id:        int64(l.ID),
			Type:      l.ChangeType,
			Title:     l.Title,
			Detail:    l.Detail,
			Timestamp: l.CreatedAt.Format(time.RFC3339),
		}
		if err := stream.Send(event); err != nil {
			return err
		}
		sinceID = int64(l.ID)
	}

	// Watch loop
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-stream.Context().Done():
			return nil
		case <-ticker.C:
			newLogs, err := s.ChangeLogSvc.GetAfterID(sinceID)
			if err != nil {
				continue
			}
			for _, l := range newLogs {
				event := &pb.ChangeEvent{
					Id:        int64(l.ID),
					Type:      l.ChangeType,
					Title:     l.Title,
					Detail:    l.Detail,
					Timestamp: l.CreatedAt.Format(time.RFC3339),
				}
				if err := stream.Send(event); err != nil {
					return err
				}
				sinceID = int64(l.ID)
			}
		}
	}
}
```

- [ ] **Step 13: Run test to verify it passes**

Run: `cd server && go test ./grpc/ -v -run TestWatchChanges`
Expected: PASS

- [ ] **Step 14: Commit**

```bash
git add server/grpc/
git commit -m "feat: implement gRPC GetLatestStats, GetTrends, WatchChanges"
```

---

## Task 8: API Key Auth + gRPC Interceptor

**Files:**
- Create: `server/grpc/auth.go`
- Create: `server/grpc/auth_test.go`
- Modify: `server/grpc/server.go` (add interceptor)

- [ ] **Step 1: Write failing test for auth interceptor**

Create `server/grpc/auth_test.go`:

```go
package grpc

import (
	"testing"
	"vblog-core/service"
	"vblog-core/testutil"
)

func TestAuthInterceptor_ValidKey(t *testing.T) {
	db := testutil.GetTestDB(t)
	st := service.NewSettingService(db)
	st.Set("grpc_api_key", "valid-key")

	interceptor := NewAuthInterceptor(st)
	// Test that valid key passes
	err := interceptor.ValidateKey("valid-key")
	if err != nil {
		t.Fatalf("expected valid key to pass: %v", err)
	}

	db.Where("key = ?", "grpc_api_key").Delete(&model.Setting{})
}

func TestAuthInterceptor_InvalidKey(t *testing.T) {
	db := testutil.GetTestDB(t)
	st := service.NewSettingService(db)
	st.Set("grpc_api_key", "valid-key")

	interceptor := NewAuthInterceptor(st)
	err := interceptor.ValidateKey("wrong-key")
	if err == nil {
		t.Fatal("expected invalid key to fail")
	}

	db.Where("key = ?", "grpc_api_key").Delete(&model.Setting{})
}

func TestAuthInterceptor_NoKey(t *testing.T) {
	db := testutil.GetTestDB(t)
	st := service.NewSettingService(db)
	// No key set

	interceptor := NewAuthInterceptor(st)
	err := interceptor.ValidateKey("anything")
	if err == nil {
		t.Fatal("expected no key configured to fail")
	}
}
```

- [ ] **Step 2: Run test to verify it fails**

Run: `cd server && go test ./grpc/ -v -run TestAuthInterceptor`
Expected: FAIL — `NewAuthInterceptor` not defined.

- [ ] **Step 3: Implement auth interceptor**

Create `server/grpc/auth.go`:

```go
package grpc

import (
	"errors"
	"vblog-core/service"
)

type AuthInterceptor struct {
	SettingSvc *service.SettingService
}

func NewAuthInterceptor(st *service.SettingService) *AuthInterceptor {
	return &AuthInterceptor{SettingSvc: st}
}

func (a *AuthInterceptor) ValidateKey(apiKey string) error {
	stored, err := a.SettingSvc.Get("grpc_api_key")
	if err != nil || stored == "" {
		return errors.New("grpc api key not configured")
	}
	if stored != apiKey {
		return errors.New("invalid api key")
	}
	return nil
}
```

- [ ] **Step 4: Run test to verify it passes**

Run: `cd server && go test ./grpc/ -v -run TestAuthInterceptor`
Expected: PASS

- [ ] **Step 5: Integrate auth into WatchChanges**

Update `WatchChanges` in `server/grpc/server.go` to use `AuthInterceptor`:

```go
func (s *Server) WatchChanges(in *pb.WatchRequest, stream pb.BlogAnalytics_WatchChangesServer) error {
	auth := NewAuthInterceptor(s.SettingSvc)
	if err := auth.ValidateKey(in.ApiKey); err != nil {
		return status.Error(codes.Unauthenticated, err.Error())
	}
	// ... rest of implementation
}
```

- [ ] **Step 6: Run all gRPC tests**

Run: `cd server && go test ./grpc/ -v`
Expected: All PASS

- [ ] **Step 7: Commit**

```bash
git add server/grpc/auth.go server/grpc/auth_test.go server/grpc/server.go
git commit -m "feat: add API Key auth interceptor for gRPC"
```

---

## Task 9: Start gRPC Server in main.go

**Files:**
- Modify: `server/cmd/main.go`
- Modify: `server/config/config.go`

- [ ] **Step 1: Add gRPC port to config**

Modify `server/config/config.go` — add GrpcPort to ServerConfig:

```go
type ServerConfig struct {
	Addr      string
	Port      string
	GrpcPort  string
}
```

In Load(), add:
```go
cfg.Server.GrpcPort = v.GetString("http.grpc_port")
if cfg.Server.GrpcPort == "" {
    cfg.Server.GrpcPort = "50051"
}
```

- [ ] **Step 2: Add gRPC server startup to main.go**

Modify `server/cmd/main.go` — add imports and gRPC startup:

```go
import (
	// ... existing imports ...
	"net"
	"google.golang.org/grpc"
	grpcpkg "vblog-core/grpc"
)

// After HTTP server setup, before ListenAndServe:

// gRPC server
dailyStatsSvc := service.NewDailyStatsService(db)
changeLogSvc := service.NewChangeLogService(db)
pageViewSvc := service.NewPageViewService(db)

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
```

- [ ] **Step 3: Add PV middleware to HTTP handler**

Wrap the static file handler with PV middleware in main.go.

- [ ] **Step 4: Build and verify server starts**

Run: `cd server && go build ./cmd/`
Expected: Build succeeds.

- [ ] **Step 5: Run all server tests**

Run: `cd server && go test ./... -v`
Expected: All PASS

- [ ] **Step 6: Commit**

```bash
git add server/cmd/main.go server/config/config.go
git commit -m "feat: start gRPC server alongside HTTP, add PV middleware"
```

---

## Task 10: Wails Client Scaffolding

**Files:**
- Create: `client/main.go`
- Create: `client/app.go`
- Create: `client/wails.json`
- Create: `client/frontend/index.html`
- Create: `client/frontend/src/main.js`
- Create: `client/frontend/src/App.vue`
- Create: `client/go.mod`

- [ ] **Step 1: Initialize Wails project**

```bash
cd "D:/Desktop/code/vibe/vBlog Core"
wails init -n client -t vue
```

Or manually create the structure if wails CLI not available.

- [ ] **Step 2: Write failing test for App.Connect**

Create `client/app_test.go`:

```go
package main

import "testing"

func TestApp_Connect_InvalidAddress(t *testing.T) {
	app := NewApp()
	err := app.Connect("invalid:99999", "test-key")
	if err == nil {
		t.Fatal("expected error for invalid address")
	}
}

func TestApp_Disconnect(t *testing.T) {
	app := NewApp()
	err := app.Disconnect()
	if err != nil {
		t.Fatalf("Disconnect on unconnected app should not error: %v", err)
	}
}
```

- [ ] **Step 3: Run test to verify it fails**

Run: `cd client && go test -v -run TestApp`
Expected: FAIL — `NewApp` not defined.

- [ ] **Step 4: Implement App struct**

Create `client/app.go`:

```go
package main

import (
	"context"
	"errors"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	pb "vblog-core/proto"
)

type App struct {
	conn   *grpc.ClientConn
	client pb.BlogAnalyticsClient
	ctx    context.Context
	cancel context.CancelFunc
}

func NewApp() *App {
	ctx, cancel := context.WithCancel(context.Background())
	return &App{ctx: ctx, cancel: cancel}
}

func (a *App) Connect(addr, apiKey string) error {
	if addr == "" {
		return errors.New("address is required")
	}
	conn, err := grpc.NewClient(addr,
		grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return err
	}
	a.conn = conn
	a.client = pb.NewBlogAnalyticsClient(conn)
	return nil
}

func (a *App) Disconnect() error {
	a.cancel()
	if a.conn != nil {
		return a.conn.Close()
	}
	return nil
}
```

- [ ] **Step 5: Run test to verify it passes**

Run: `cd client && go test -v -run TestApp`
Expected: PASS

- [ ] **Step 6: Implement GetLatestStats and GetTrends bindings**

Add to `client/app.go`:

```go
func (a *App) GetLatestStats() (*pb.LatestStats, error) {
	if a.client == nil {
		return nil, errors.New("not connected")
	}
	return a.client.GetLatestStats(a.ctx, &pb.Empty{})
}

func (a *App) GetTrends(granularity string, count int32) (*pb.GetTrendsResponse, error) {
	if a.client == nil {
		return nil, errors.New("not connected")
	}
	return a.client.GetTrends(a.ctx, &pb.GetTrendsRequest{
		Granularity: granularity,
		Count:       count,
	})
}
```

- [ ] **Step 7: Implement WatchChanges with Wails Events**

Add to `client/app.go`:

```go
import "github.com/wailsapp/wails/v2/pkg/runtime"

func (a *App) WatchChanges(apiKey string, sinceID int64) error {
	if a.client == nil {
		return errors.New("not connected")
	}

	stream, err := a.client.WatchChanges(a.ctx, &pb.WatchRequest{
		ApiKey:  apiKey,
		SinceId: sinceID,
	})
	if err != nil {
		return err
	}

	go func() {
		for {
			event, err := stream.Recv()
			if err != nil {
				return
			}
			runtime.EventsEmit(a.ctx, "change", event)
		}
	}()
	return nil
}
```

- [ ] **Step 8: Write main.go**

Create `client/main.go`:

```go
package main

import (
	"embed"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
)

//go:embed all:frontend/dist
var assets embed.FS

func main() {
	app := NewApp()

	err := wails.Run(&options.App{
		Title:     "vBlog Monitor",
		Width:     420,
		Height:    640,
		Assets:    assets,
		OnStartup: app.startup,
		Bind: []interface{}{
			app,
		},
	})
	if err != nil {
		println("Error:", err.Error())
	}
}
```

- [ ] **Step 9: Create frontend skeleton**

Create `client/frontend/src/App.vue`:

```vue
<template>
  <div class="app">
    <header class="app-header">
      <h1>vBlog Monitor</h1>
      <button @click="showSettings = true">Settings</button>
    </header>
    <StatsBar :stats="stats" />
    <div class="change-feed">
      <ChangeCard v-for="c in changes" :key="c.id" :change="c" />
    </div>
    <TrendPanel :points="trends" />
    <Settings v-if="showSettings" @close="showSettings = false" @connect="handleConnect" />
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import StatsBar from './components/StatsBar.vue'
import ChangeCard from './components/ChangeCard.vue'
import TrendPanel from './components/TrendPanel.vue'
import Settings from './components/Settings.vue'

const stats = ref({})
const changes = ref([])
const trends = ref([])
const showSettings = ref(false)

// Wails runtime events will be set up here
</script>
```

- [ ] **Step 10: Build client**

```bash
cd client
go mod tidy
go build .
```

Expected: Build succeeds.

- [ ] **Step 11: Commit**

```bash
git add client/
git commit -m "feat: Wails client scaffolding with Connect, GetLatestStats, GetTrends, WatchChanges"
```

---

## Task 11: Client Frontend Components

**Files:**
- Create: `client/frontend/src/components/ChangeCard.vue`
- Create: `client/frontend/src/components/StatsBar.vue`
- Create: `client/frontend/src/components/TrendPanel.vue`
- Create: `client/frontend/src/components/Settings.vue`
- Modify: `client/frontend/src/App.vue`

- [ ] **Step 1: Create StatsBar component**

Create `client/frontend/src/components/StatsBar.vue`:

```vue
<template>
  <div class="stats-bar" v-if="stats">
    <div class="stat">
      <span class="stat-val">{{ stats.pvToday || 0 }}</span>
      <span class="stat-label">PV Today</span>
      <span class="stat-delta" v-if="pvDelta !== null">{{ pvDelta > 0 ? '+' : '' }}{{ pvDelta }}%</span>
    </div>
    <div class="stat">
      <span class="stat-val">{{ stats.uvToday || 0 }}</span>
      <span class="stat-label">UV Today</span>
      <span class="stat-delta" v-if="uvDelta !== null">{{ uvDelta > 0 ? '+' : '' }}{{ uvDelta }}%</span>
    </div>
    <div class="stat">
      <span class="stat-val">{{ stats.totalViews || 0 }}</span>
      <span class="stat-label">Total Views</span>
    </div>
    <div class="stat">
      <span class="stat-val">{{ stats.totalPosts || 0 }}</span>
      <span class="stat-label">Posts</span>
    </div>
  </div>
</template>

<script setup>
import { computed } from 'vue'
const props = defineProps({ stats: Object })

const pvDelta = computed(() => {
  if (!props.stats?.pvYesterday) return null
  return Math.round(((props.stats.pvToday - props.stats.pvYesterday) / props.stats.pvYesterday) * 100)
})
const uvDelta = computed(() => {
  if (!props.stats?.uvYesterday) return null
  return Math.round(((props.stats.uvToday - props.stats.uvYesterday) / props.stats.uvYesterday) * 100)
})
</script>
```

- [ ] **Step 2: Create ChangeCard component**

Create `client/frontend/src/components/ChangeCard.vue`:

```vue
<template>
  <div class="change-card">
    <span class="change-icon">{{ icon }}</span>
    <div class="change-body">
      <div class="change-title">{{ change.title }}</div>
      <div class="change-time">{{ change.timestamp }}</div>
    </div>
  </div>
</template>

<script setup>
import { computed } from 'vue'
const props = defineProps({ change: Object })

const icons = {
  new_post: 'New post',
  new_comment: 'Comment',
  view_milestone: 'Milestone',
  pv_milestone: 'Target',
  tag_added: 'Tag',
}
const icon = computed(() => icons[props.change?.type] || 'Update')
</script>
```

- [ ] **Step 3: Create TrendPanel component**

Create `client/frontend/src/components/TrendPanel.vue`:

```vue
<template>
  <div class="trend-panel">
    <div class="trend-header">
      <h3>Trends</h3>
      <div class="trend-tabs">
        <button v-for="g in ['day','week','month']" :key="g"
          :class="{ active: granularity === g }"
          @click="$emit('change', g)">{{ g }}</button>
      </div>
    </div>
    <div class="trend-chart">
      <div v-for="(p, i) in points" :key="i" class="trend-bar-group">
        <div class="trend-bar" :style="{ height: barHeight(p.pv) + '%' }"></div>
        <span class="trend-label">{{ p.label.slice(5) }}</span>
      </div>
    </div>
  </div>
</template>

<script setup>
import { computed } from 'vue'
const props = defineProps({ points: Array, granularity: String })

const maxPv = computed(() => Math.max(...(props.points?.map(p => p.pv) || [1])))
const barHeight = (pv) => Math.round((pv / maxPv.value) * 100)
</script>
```

- [ ] **Step 4: Create Settings component**

Create `client/frontend/src/components/Settings.vue`:

```vue
<template>
  <div class="settings-overlay" @click.self="$emit('close')">
    <div class="settings-panel">
      <h2>Connection</h2>
      <label>
        Server Address
        <input v-model="addr" placeholder="localhost:50051" />
      </label>
      <label>
        API Key
        <input v-model="apiKey" type="password" placeholder="Enter API key" />
      </label>
      <div class="settings-actions">
        <button @click="$emit('close')">Cancel</button>
        <button class="primary" @click="connect">Connect</button>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref } from 'vue'
const emit = defineEmits(['close', 'connect'])
const addr = ref(localStorage.getItem('vblog_addr') || 'localhost:50051')
const apiKey = ref(localStorage.getItem('vblog_key') || '')

function connect() {
  localStorage.setItem('vblog_addr', addr.value)
  localStorage.setItem('vblog_key', apiKey.value)
  emit('connect', { addr: addr.value, apiKey: apiKey.value })
}
</script>
```

- [ ] **Step 5: Update App.vue to wire everything**

Update `client/frontend/src/App.vue`:

```vue
<template>
  <div class="app">
    <header class="app-header">
      <h1>vBlog Monitor</h1>
      <button @click="showSettings = true">Settings</button>
    </header>
    <StatsBar :stats="stats" />
    <div class="change-feed">
      <ChangeCard v-for="c in changes" :key="c.id" :change="c" />
      <div v-if="!changes.length" class="empty">No changes yet</div>
    </div>
    <TrendPanel :points="trends" :granularity="granularity" @change="loadTrends" />
    <Settings v-if="showSettings" @close="showSettings = false" @connect="handleConnect" />
  </div>
</template>

<script setup>
import { ref, onMounted, onUnmounted } from 'vue'
import StatsBar from './components/StatsBar.vue'
import ChangeCard from './components/ChangeCard.vue'
import TrendPanel from './components/TrendPanel.vue'
import Settings from './components/Settings.vue'
import { Connect, GetLatestStats, GetTrends, WatchChanges } from '../wailsjs/main/App'

const stats = ref({})
const changes = ref([])
const trends = ref([])
const granularity = ref('day')
const showSettings = ref(false)
let refreshInterval = null

async function handleConnect({ addr, apiKey }) {
  try {
    await Connect(addr, apiKey)
    showSettings.value = false
    await refresh()
    await WatchChanges(apiKey, 0)
    refreshInterval = setInterval(refresh, 30000)
  } catch (e) {
    alert('Connection failed: ' + e)
  }
}

async function refresh() {
  try {
    stats.value = await GetLatestStats()
    await loadTrends(granularity.value)
  } catch {}
}

async function loadTrends(g) {
  granularity.value = g
  const resp = await GetTrends(g, 14)
  trends.value = resp?.points || []
}

onMounted(() => {
  // Listen for Wails events
  if (window.runtime) {
    window.runtime.EventsOn('change', (event) => {
      changes.value.unshift(event)
      if (changes.value.length > 50) changes.value.pop()
    })
  }
})

onUnmounted(() => {
  if (refreshInterval) clearInterval(refreshInterval)
})
</script>
```

- [ ] **Step 6: Build frontend**

```bash
cd client/frontend
npm install
npm run build
```

Expected: Build succeeds.

- [ ] **Step 7: Commit**

```bash
git add client/frontend/
git commit -m "feat: client frontend with StatsBar, ChangeCard, TrendPanel, Settings"
```

---

## Task 12: Final Integration + Cleanup

- [ ] **Step 1: Run all server tests**

Run: `cd server && go test ./... -v`
Expected: All PASS

- [ ] **Step 2: Build server**

Run: `cd server && go build ./cmd/`
Expected: Build succeeds.

- [ ] **Step 3: Build client**

Run: `cd client && go build .`
Expected: Build succeeds.

- [ ] **Step 4: Add grpc_api_key to seed data**

Update `server/cmd/seed/main.go` to insert a default API key:

```go
db.FirstOrCreate(&model.Setting{Key: "grpc_api_key", Value: "change-me-in-production"})
```

- [ ] **Step 5: Commit**

```bash
git add .
git commit -m "feat: final integration, seed grpc_api_key, all tests passing"
```
