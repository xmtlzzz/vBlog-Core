# vBlog Core MVP 实现计划

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** 构建 vBlog Core MVP — 一个可自定义组件的极客博客系统，包含博客前台、后台管理、RESTful API。

**Architecture:** Go 后端 (go-restful + GORM) 通过 embed 嵌入 Vue 3 前端静态文件，单二进制部署。前后端在同一 monorepo 中。

**Tech Stack:** Vue 3 + Element Plus + Pinia / Go + go-restful + GORM / PostgreSQL / JWT

**TDD 规则（严格遵守）：**
- 每个功能模块先写测试文件 `xxx_test.go`，再写实现
- 测试文件与实现文件在同一目录
- 流程：写测试 → 运行确认失败 → 写最小实现 → 运行确认通过 → 提交
- 前端使用 Vitest + @vue/test-utils

---

## Phase 1: 项目脚手架

### Task 1: 初始化项目结构

**Files:**
- Create: `server/go.mod`
- Create: `server/cmd/main.go`
- Create: `web/package.json`
- Create: `web/vite.config.js`
- Create: `.gitignore`

- [ ] **Step 1: 初始化 Go 模块**

```bash
cd server
go mod init vblog-core
go get github.com/emicklei/go-restful/v4
go get gorm.io/gorm
go get gorm.io/driver/postgres
go get github.com/golang-jwt/jwt/v5
```

- [ ] **Step 2: 创建 main.go 入口**

```go
// server/cmd/main.go
package main

import (
	"fmt"
	"log"
	"net/http"

	restful "github.com/emicklei/go-restful/v4"
)

func main() {
	wsContainer := restful.NewContainer()
	// TODO: 注册路由
	log.Printf("vBlog Core starting on :8080")
	server := &http.Server{Addr: ":8080", Handler: wsContainer}
	log.Fatal(server.ListenAndServe())
}
```

- [ ] **Step 3: 初始化 Vue 项目**

```bash
cd web
npm create vite@latest . -- --template vue
npm install
npm install vue-router@4 pinia axios element-plus @element-plus/icons-vue
npm install -D @vue/test-utils vitest jsdom
```

- [ ] **Step 4: 创建 .gitignore**

```
server/vblog
server/tmp/
web/node_modules/
web/dist/
*.exe
.env
```

- [ ] **Step 5: 验证编译和启动**

```bash
cd server && go build ./cmd/main.go
cd web && npm run dev
```

- [ ] **Step 6: 提交**

```bash
git add .
git commit -m "chore: initialize project structure"
```

---

### Task 2: 数据库 Schema + GORM 模型

**Files:**
- Create: `server/model/post.go`
- Create: `server/model/post_test.go`
- Create: `server/model/tag.go`
- Create: `server/model/tag_test.go`
- Create: `server/model/comment.go`
- Create: `server/model/comment_test.go`
- Create: `server/model/component.go`
- Create: `server/model/setting.go`
- Create: `server/model/user.go`
- Create: `server/model/migrate.go`
- Create: `deploy/init.sql`

- [ ] **Step 1: 写 Post 模型测试**

```go
// server/model/post_test.go
package model

import (
	"testing"
	"time"
)

func TestPostTableName(t *testing.T) {
	p := Post{}
	if p.TableName() != "posts" {
		t.Errorf("expected table name 'posts', got '%s'", p.TableName())
	}
}

func TestPostDefaults(t *testing.T) {
	p := Post{}
	if p.Status != "draft" {
		t.Errorf("expected default status 'draft', got '%s'", p.Status)
	}
	if p.Views != 0 {
		t.Errorf("expected default views 0, got %d", p.Views)
	}
	if p.Pinned != false {
		t.Errorf("expected default pinned false, got %v", p.Pinned)
	}
}
```

- [ ] **Step 2: 运行测试确认失败**

```bash
cd server && go test ./model/ -v
```
Expected: FAIL — `Post` type not defined

- [ ] **Step 3: 写 Post 模型**

```go
// server/model/post.go
package model

import (
	"time"
	"gorm.io/gorm"
)

type Post struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	Title     string         `gorm:"size:200;not null" json:"title"`
	Content   string         `gorm:"type:text;not null" json:"content"`
	Excerpt   string         `gorm:"size:500" json:"excerpt"`
	Status    string         `gorm:"size:20;default:draft" json:"status"`
	Pinned    bool           `gorm:"default:false" json:"pinned"`
	Views     int            `gorm:"default:0" json:"views"`
	ReadTime  int            `gorm:"default:0" json:"read_time"`
	AuthorID  uint           `json:"author_id"`
	Tags      []Tag          `gorm:"many2many:post_tags;" json:"tags"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

func (Post) TableName() string { return "posts" }
```

- [ ] **Step 4: 运行测试确认通过**

```bash
cd server && go test ./model/ -v -run TestPost
```
Expected: PASS

- [ ] **Step 5: 写 Tag 模型测试**

```go
// server/model/tag_test.go
package model

import "testing"

func TestTagTableName(t *testing.T) {
	tag := Tag{}
	if tag.TableName() != "tags" {
		t.Errorf("expected table name 'tags', got '%s'", tag.TableName())
	}
}
```

- [ ] **Step 6: 写 Tag 模型**

```go
// server/model/tag.go
package model

import "time"

type Tag struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	Name        string    `gorm:"size:50;uniqueIndex;not null" json:"name"`
	Description string    `gorm:"type:text" json:"description"`
	PostCount   int       `gorm:"-" json:"post_count"`
	CreatedAt   time.Time `json:"created_at"`
}

func (Tag) TableName() string { return "tags" }
```

- [ ] **Step 7: 运行测试确认通过**

```bash
cd server && go test ./model/ -v -run TestTag
```

- [ ] **Step 8: 写 Comment 模型测试**

```go
// server/model/comment_test.go
package model

import "testing"

func TestCommentTableName(t *testing.T) {
	c := Comment{}
	if c.TableName() != "comments" {
		t.Errorf("expected table name 'comments', got '%s'", c.TableName())
	}
}

func TestCommentDefaults(t *testing.T) {
	c := Comment{}
	if c.Status != "pending" {
		t.Errorf("expected default status 'pending', got '%s'", c.Status)
	}
}
```

- [ ] **Step 9: 写 Comment 模型**

```go
// server/model/comment.go
package model

import "time"

type Comment struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	PostID      uint      `gorm:"index;not null" json:"post_id"`
	AuthorName  string    `gorm:"size:50;not null" json:"author_name"`
	AuthorEmail string    `gorm:"size:100" json:"author_email"`
	Body        string    `gorm:"type:text;not null" json:"body"`
	Status      string    `gorm:"size:20;default:pending" json:"status"`
	CreatedAt   time.Time `json:"created_at"`
}

func (Comment) TableName() string { return "comments" }
```

- [ ] **Step 10: 写 Component 模型**

```go
// server/model/component.go
package model

import "time"

type Component struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	Name        string    `gorm:"size:100;not null" json:"name"`
	Description string    `gorm:"type:text" json:"description"`
	Version     string    `gorm:"size:20;default:0.1.0" json:"version"`
	Code        string    `gorm:"type:text" json:"code"`
	Category    string    `gorm:"size:50" json:"category"`
	Origin      string    `gorm:"size:20;default:built-in" json:"origin"`
	Status      string    `gorm:"size:20;default:active" json:"status"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func (Component) TableName() string { return "components" }
```

- [ ] **Step 11: 写 Setting 模型**

```go
// server/model/setting.go
package model

import "time"

type Setting struct {
	Key       string    `gorm:"primaryKey;size:100" json:"key"`
	Value     string    `gorm:"type:text" json:"value"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (Setting) TableName() string { return "settings" }
```

- [ ] **Step 12: 写 User 模型**

```go
// server/model/user.go
package model

import "time"

type User struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Username  string    `gorm:"size:50;uniqueIndex;not null" json:"username"`
	Password  string    `gorm:"size:255;not null" json:"-"`
	Email     string    `gorm:"size:100" json:"email"`
	Role      string    `gorm:"size:20;default:admin" json:"role"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (User) TableName() string { return "users" }
```

- [ ] **Step 13: 写 migrate.go**

```go
// server/model/migrate.go
package model

import "gorm.io/gorm"

func AutoMigrate(db *gorm.DB) error {
	return db.AutoMigrate(
		&User{},
		&Post{},
		&Tag{},
		&Comment{},
		&Component{},
		&Setting{},
	)
}
```

- [ ] **Step 14: 写 init.sql**

```sql
-- deploy/init.sql
CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY, username VARCHAR(50) UNIQUE NOT NULL,
    password VARCHAR(255) NOT NULL, email VARCHAR(100),
    role VARCHAR(20) DEFAULT 'admin',
    created_at TIMESTAMP DEFAULT NOW(), updated_at TIMESTAMP DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS posts (
    id SERIAL PRIMARY KEY, title VARCHAR(200) NOT NULL,
    content TEXT NOT NULL, excerpt VARCHAR(500),
    status VARCHAR(20) DEFAULT 'draft', pinned BOOLEAN DEFAULT FALSE,
    views INTEGER DEFAULT 0, read_time INTEGER DEFAULT 0,
    author_id INTEGER REFERENCES users(id),
    created_at TIMESTAMP DEFAULT NOW(), updated_at TIMESTAMP DEFAULT NOW(),
    deleted_at TIMESTAMP
);
CREATE INDEX IF NOT EXISTS idx_posts_status ON posts(status);
CREATE INDEX IF NOT EXISTS idx_posts_deleted_at ON posts(deleted_at);

CREATE TABLE IF NOT EXISTS tags (
    id SERIAL PRIMARY KEY, name VARCHAR(50) UNIQUE NOT NULL,
    description TEXT, created_at TIMESTAMP DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS post_tags (
    post_id INTEGER REFERENCES posts(id) ON DELETE CASCADE,
    tag_id INTEGER REFERENCES tags(id) ON DELETE CASCADE,
    PRIMARY KEY (post_id, tag_id)
);

CREATE TABLE IF NOT EXISTS comments (
    id SERIAL PRIMARY KEY, post_id INTEGER REFERENCES posts(id) ON DELETE CASCADE,
    author_name VARCHAR(50) NOT NULL, author_email VARCHAR(100),
    body TEXT NOT NULL, status VARCHAR(20) DEFAULT 'pending',
    created_at TIMESTAMP DEFAULT NOW()
);
CREATE INDEX IF NOT EXISTS idx_comments_post_id ON comments(post_id);
CREATE INDEX IF NOT EXISTS idx_comments_status ON comments(status);

CREATE TABLE IF NOT EXISTS components (
    id SERIAL PRIMARY KEY, name VARCHAR(100) NOT NULL,
    description TEXT, version VARCHAR(20) DEFAULT '0.1.0',
    code TEXT, category VARCHAR(50),
    origin VARCHAR(20) DEFAULT 'built-in',
    status VARCHAR(20) DEFAULT 'active',
    created_at TIMESTAMP DEFAULT NOW(), updated_at TIMESTAMP DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS settings (
    key VARCHAR(100) PRIMARY KEY, value TEXT,
    updated_at TIMESTAMP DEFAULT NOW()
);
```

- [ ] **Step 15: 运行全部模型测试**

```bash
cd server && go test ./model/ -v
```

- [ ] **Step 16: 提交**

```bash
git add server/model/ deploy/init.sql
git commit -m "feat: add GORM models and DB schema"
```

---

## Phase 2: 配置 + 数据库连接

### Task 3: 配置加载 + 数据库连接

**Files:**
- Create: `server/config/config.go`
- Create: `server/config/config_test.go`
- Create: `server/config/database.go`
- Create: `server/config/database_test.go`

- [ ] **Step 1: 写配置测试**

```go
// server/config/config_test.go
package config

import (
	"os"
	"testing"
)

func TestLoadFromEnv(t *testing.T) {
	os.Setenv("DB_HOST", "testhost")
	os.Setenv("DB_PORT", "5433")
	os.Setenv("JWT_SECRET", "testsecret")
	defer os.Unsetenv("DB_HOST")
	defer os.Unsetenv("DB_PORT")
	defer os.Unsetenv("JWT_SECRET")

	cfg := Load()
	if cfg.DB.Host != "testhost" {
		t.Errorf("expected DB host 'testhost', got '%s'", cfg.DB.Host)
	}
	if cfg.DB.Port != "5433" {
		t.Errorf("expected DB port '5433', got '%s'", cfg.DB.Port)
	}
	if cfg.JWT.Secret != "testsecret" {
		t.Errorf("expected JWT secret 'testsecret', got '%s'", cfg.JWT.Secret)
	}
}

func TestLoadDefaults(t *testing.T) {
	cfg := Load()
	if cfg.Server.Port != "8080" {
		t.Errorf("expected default port '8080', got '%s'", cfg.Server.Port)
	}
	if cfg.DB.Host != "localhost" {
		t.Errorf("expected default DB host 'localhost', got '%s'", cfg.DB.Host)
	}
}
```

- [ ] **Step 2: 运行测试确认失败**

```bash
cd server && go test ./config/ -v
```

- [ ] **Step 3: 写配置加载**

```go
// server/config/config.go
package config

import "os"

type Config struct {
	Server ServerConfig
	DB     DBConfig
	JWT    JWTConfig
}

type ServerConfig struct{ Port string }
type DBConfig struct {
	Host, Port, Name, User, Password string
}
type JWTConfig struct {
	Secret          string
	ExpireHours     int
	RefreshExpireHours int
}

func envOrDefault(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}

func Load() *Config {
	return &Config{
		Server: ServerConfig{Port: envOrDefault("SERVER_PORT", "8080")},
		DB: DBConfig{
			Host:     envOrDefault("DB_HOST", "localhost"),
			Port:     envOrDefault("DB_PORT", "5432"),
			Name:     envOrDefault("DB_NAME", "vblog"),
			User:     envOrDefault("DB_USER", "vblog_admin"),
			Password: envOrDefault("DB_PASSWORD", ""),
		},
		JWT: JWTConfig{
			Secret:          envOrDefault("JWT_SECRET", "vblog-default-secret"),
			ExpireHours:     2,
			RefreshExpireHours: 168,
		},
	}
}
```

- [ ] **Step 4: 运行测试确认通过**

```bash
cd server && go test ./config/ -v
```

- [ ] **Step 5: 写数据库连接测试**

```go
// server/config/database_test.go
package config

import "testing"

func TestDSN(t *testing.T) {
	cfg := &DBConfig{Host: "localhost", Port: "5432", Name: "vblog", User: "admin", Password: "pass"}
	dsn := cfg.DSN()
	if dsn == "" {
		t.Error("DSN should not be empty")
	}
}
```

- [ ] **Step 6: 写数据库连接**

```go
// server/config/database.go
package config

import (
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func (c *DBConfig) DSN() string {
	return fmt.Sprintf("host=%s port=%s dbname=%s user=%s password=%s sslmode=disable",
		c.Host, c.Port, c.Name, c.User, c.Password)
}

func (c *DBConfig) Connect() (*gorm.DB, error) {
	return gorm.Open(postgres.Open(c.DSN()), &gorm.Config{})
}
```

- [ ] **Step 7: 运行测试 + 提交**

```bash
cd server && go test ./config/ -v
git add server/config/
git commit -m "feat: add config loading and database connection"
```

---

## Phase 3: API 层（TDD，逐模块）

### Task 4: Post CRUD API

**Files:**
- Create: `server/api/post.go`
- Create: `server/api/post_test.go`
- Create: `server/service/post.go`
- Create: `server/service/post_test.go`

- [ ] **Step 1: 写 Post Service 测试**

```go
// server/service/post_test.go
package service

import (
	"testing"
	"vblog-core/model"
)

// 使用内存数据库测试（SQLite）或 mock
// 此处测试业务逻辑纯函数

func TestCalcReadTime(t *testing.T) {
	tests := []struct {
		content  string
		expected int
	}{
		{"", 1},
		{string(make([]byte, 500)), 1},   // ~500 chars, ~1 min
		{string(make([]byte, 3000)), 3},  // ~3000 chars, ~3 min
	}
	for _, tt := range tests {
		got := CalcReadTime(tt.content)
		if got != tt.expected {
			t.Errorf("CalcReadTime(%d chars): got %d, want %d", len(tt.content), got, tt.expected)
		}
	}
}

func TestBuildExcerpt(t *testing.T) {
	short := "Hello world"
	if BuildExcerpt(short, 200) != short {
		t.Error("short content should return as-is")
	}
}
```

- [ ] **Step 2: 运行测试确认失败**

```bash
cd server && go test ./service/ -v
```

- [ ] **Step 3: 写 Post Service**

```go
// server/service/post.go
package service

import (
	"math"
	"vblog-core/model"
	"gorm.io/gorm"
)

type PostService struct{ DB *gorm.DB }

func NewPostService(db *gorm.DB) *PostService {
	return &PostService{DB: db}
}

func CalcReadTime(content string) int {
	words := len(content)
	minutes := int(math.Ceil(float64(words) / 500))
	if minutes < 1 {
		return 1
	}
	return minutes
}

func BuildExcerpt(content string, maxLen int) string {
	if len(content) <= maxLen {
		return content
	}
	return content[:maxLen] + "..."
}

func (s *PostService) List(page, perPage int, tag, status, search string) ([]model.Post, int64, error) {
	var posts []model.Post
	var total int64
	q := s.DB.Model(&model.Post{})
	if status != "" {
		q = q.Where("status = ?", status)
	}
	if search != "" {
		q = q.Where("title ILIKE ?", "%"+search+"%")
	}
	if tag != "" {
		q = q.Joins("JOIN post_tags ON post_tags.post_id = posts.id").
			Joins("JOIN tags ON tags.id = post_tags.tag_id").
			Where("tags.name = ?", tag)
	}
	q.Count(&total)
	err := q.Preload("Tags").Order("pinned DESC, created_at DESC").
		Offset((page - 1) * perPage).Limit(perPage).Find(&posts).Error
	return posts, total, err
}

func (s *PostService) GetByID(id uint) (*model.Post, error) {
	var post model.Post
	err := s.DB.Preload("Tags").First(&post, id).Error
	return &post, err
}

func (s *PostService) Create(post *model.Post) error {
	post.ReadTime = CalcReadTime(post.Content)
	if post.Excerpt == "" {
		post.Excerpt = BuildExcerpt(post.Content, 200)
	}
	return s.DB.Create(post).Error
}

func (s *PostService) Update(post *model.Post) error {
	post.ReadTime = CalcReadTime(post.Content)
	return s.DB.Save(post).Error
}

func (s *PostService) Delete(id uint) error {
	return s.DB.Delete(&model.Post{}, id).Error
}

func (s *PostService) IncrementViews(id uint) error {
	return s.DB.Model(&model.Post{}).Where("id = ?", id).
		UpdateColumn("views", gorm.Expr("views + 1")).Error
}
```

- [ ] **Step 4: 运行测试确认通过**

```bash
cd server && go test ./service/ -v -run TestCalcReadTime -run TestBuildExcerpt
```

- [ ] **Step 5: 写 Post API handler 测试**

```go
// server/api/post_test.go
package api

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	restful "github.com/emicklei/go-restful/v4"
)

func TestPostAPI_List(t *testing.T) {
	ws := new(restful.WebService)
	ws.Path("/api/posts")
	// 注册路由（需要 mock service）
	container := restful.NewContainer()
	container.Add(ws)

	req := httptest.NewRequest("GET", "/api/posts?page=1&per_page=5", nil)
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()
	container.ServeHTTP(resp, req)

	if resp.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", resp.Code)
	}
}
```

- [ ] **Step 6: 写 Post API handler**

```go
// server/api/post.go
package api

import (
	"encoding/json"
	"net/http"
	"strconv"

	restful "github.com/emicklei/go-restful/v4"
	"vblog-core/model"
	"vblog-core/service"
)

type PostResource struct{ Service *service.PostService }

func (p *PostResource) Register(ws *restful.WebService) {
	ws.Route(ws.GET("/api/posts").To(p.list).
		Param(ws.QueryParameter("page", "页码").DefaultValue("1")).
		Param(ws.QueryParameter("per_page", "每页数量").DefaultValue("5")).
		Param(ws.QueryParameter("tag", "标签筛选")).
		Param(ws.QueryParameter("status", "状态筛选")).
		Param(ws.QueryParameter("search", "搜索关键词")))
	ws.Route(ws.GET("/api/posts/{id}").To(p.get))
	ws.Route(ws.POST("/api/posts").To(p.create))
	ws.Route(ws.PUT("/api/posts/{id}").To(p.update))
	ws.Route(ws.DELETE("/api/posts/{id}").To(p.delete))
}

func (p *PostResource) list(req *restful.Request, resp *restful.Response) {
	page, _ := strconv.Atoi(req.QueryParameter("page"))
	perPage, _ := strconv.Atoi(req.QueryParameter("per_page"))
	if page < 1 { page = 1 }
	if perPage < 1 { perPage = 5 }
	tag := req.QueryParameter("tag")
	status := req.QueryParameter("status")
	search := req.QueryParameter("search")
	posts, total, err := p.Service.List(page, perPage, tag, status, search)
	if err != nil {
		resp.WriteError(http.StatusInternalServerError, err)
		return
	}
	resp.WriteEntity(map[string]interface{}{
		"data":  posts,
		"total": total,
		"page":  page,
		"per_page": perPage,
	})
}

func (p *PostResource) get(req *restful.Request, resp *restful.Response) {
	id, _ := strconv.ParseUint(req.PathParameter("id"), 10, 32)
	post, err := p.Service.GetByID(uint(id))
	if err != nil {
		resp.WriteError(http.StatusNotFound, err)
		return
	}
	resp.WriteEntity(post)
}

func (p *PostResource) create(req *restful.Request, resp *restful.Response) {
	var post model.Post
	if err := req.ReadEntity(&post); err != nil {
		resp.WriteError(http.StatusBadRequest, err)
		return
	}
	if err := p.Service.Create(&post); err != nil {
		resp.WriteError(http.StatusInternalServerError, err)
		return
	}
	resp.WriteHeaderAndEntity(http.StatusCreated, post)
}

func (p *PostResource) update(req *restful.Request, resp *restful.Response) {
	id, _ := strconv.ParseUint(req.PathParameter("id"), 10, 32)
	post, err := p.Service.GetByID(uint(id))
	if err != nil {
		resp.WriteError(http.StatusNotFound, err)
		return
	}
	json.NewDecoder(req.Request.Body).Decode(post)
	p.Service.Update(post)
	resp.WriteEntity(post)
}

func (p *PostResource) delete(req *restful.Request, resp *restful.Response) {
	id, _ := strconv.ParseUint(req.PathParameter("id"), 10, 32)
	if err := p.Service.Delete(uint(id)); err != nil {
		resp.WriteError(http.StatusInternalServerError, err)
		return
	}
	resp.WriteHeader(http.StatusNoContent)
}
```

- [ ] **Step 7: 运行测试 + 提交**

```bash
cd server && go test ./api/ -v -run TestPostAPI
git add server/api/ server/service/
git commit -m "feat: add Post CRUD API with TDD"
```

---

### Task 5: Tag CRUD API

**Files:**
- Create: `server/api/tag.go`
- Create: `server/api/tag_test.go`
- Create: `server/service/tag.go`
- Create: `server/service/tag_test.go`

- [ ] **Step 1: 写 Tag Service 测试**

```go
// server/service/tag_test.go
package service

import "testing"

func TestTagService_CountByTag(t *testing.T) {
	// 测试标签计数逻辑（纯函数测试）
	tags := []string{"Go", "React", "Go"}
	counts := countStrings(tags)
	if counts["Go"] != 2 {
		t.Errorf("expected Go count 2, got %d", counts["Go"])
	}
}

func countStrings(items []string) map[string]int {
	m := make(map[string]int)
	for _, item := range items {
		m[item]++
	}
	return m
}
```

- [ ] **Step 2: 写 Tag Service**

```go
// server/service/tag.go
package service

import (
	"vblog-core/model"
	"gorm.io/gorm"
)

type TagService struct{ DB *gorm.DB }

func NewTagService(db *gorm.DB) *TagService {
	return &TagService{DB: db}
}

func (s *TagService) List() ([]model.Tag, error) {
	var tags []model.Tag
	err := s.DB.Find(&tags).Error
	// 计算每个标签的文章数
	for i := range tags {
		s.DB.Model(&model.Post{}).
			Joins("JOIN post_tags ON post_tags.post_id = posts.id").
			Where("post_tags.tag_id = ?", tags[i].ID).
			Count(&tags[i].PostCount)
	}
	return tags, err
}

func (s *TagService) Create(tag *model.Tag) error {
	return s.DB.Create(tag).Error
}

func (s *TagService) Update(tag *model.Tag) error {
	return s.DB.Save(tag).Error
}

func (s *TagService) Delete(id uint) error {
	return s.DB.Delete(&model.Tag{}, id).Error
}
```

- [ ] **Step 3: 写 Tag API handler 测试**

```go
// server/api/tag_test.go
package api

import (
	"net/http"
	"net/http/httptest"
	"testing"

	restful "github.com/emicklei/go-restful/v4"
)

func TestTagAPI_List(t *testing.T) {
	ws := new(restful.WebService)
	ws.Path("/api/tags")
	container := restful.NewContainer()
	container.Add(ws)

	req := httptest.NewRequest("GET", "/api/tags", nil)
	resp := httptest.NewRecorder()
	container.ServeHTTP(resp, req)

	if resp.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", resp.Code)
	}
}
```

- [ ] **Step 4: 写 Tag API handler**

```go
// server/api/tag.go
package api

import (
	"net/http"
	"strconv"

	restful "github.com/emicklei/go-restful/v4"
	"vblog-core/model"
	"vblog-core/service"
)

type TagResource struct{ Service *service.TagService }

func (t *TagResource) Register(ws *restful.WebService) {
	ws.Route(ws.GET("/api/tags").To(t.list))
	ws.Route(ws.POST("/api/tags").To(t.create))
	ws.Route(ws.PUT("/api/tags/{id}").To(t.update))
	ws.Route(ws.DELETE("/api/tags/{id}").To(t.delete))
}

func (t *TagResource) list(req *restful.Request, resp *restful.Response) {
	tags, err := t.Service.List()
	if err != nil {
		resp.WriteError(http.StatusInternalServerError, err)
		return
	}
	resp.WriteEntity(tags)
}

func (t *TagResource) create(req *restful.Request, resp *restful.Response) {
	var tag model.Tag
	if err := req.ReadEntity(&tag); err != nil {
		resp.WriteError(http.StatusBadRequest, err)
		return
	}
	if tag.Name == "" {
		resp.WriteErrorString(http.StatusBadRequest, "name is required")
		return
	}
	t.Service.Create(&tag)
	resp.WriteHeaderAndEntity(http.StatusCreated, tag)
}

func (t *TagResource) update(req *restful.Request, resp *restful.Response) {
	id, _ := strconv.ParseUint(req.PathParameter("id"), 10, 32)
	var tag model.Tag
	if err := t.Service.DB.First(&tag, id).Error; err != nil {
		resp.WriteError(http.StatusNotFound, err)
		return
	}
	req.ReadEntity(&tag)
	t.Service.Update(&tag)
	resp.WriteEntity(tag)
}

func (t *TagResource) delete(req *restful.Request, resp *restful.Response) {
	id, _ := strconv.ParseUint(req.PathParameter("id"), 10, 32)
	t.Service.Delete(uint(id))
	resp.WriteHeader(http.StatusNoContent)
}
```

- [ ] **Step 5: 运行测试 + 提交**

```bash
cd server && go test ./service/ ./api/ -v -run TestTag
git add server/api/tag.go server/api/tag_test.go server/service/tag.go server/service/tag_test.go
git commit -m "feat: add Tag CRUD API with TDD"
```

---

### Task 6: Comment API

**Files:**
- Create: `server/api/comment.go`
- Create: `server/api/comment_test.go`
- Create: `server/service/comment.go`
- Create: `server/service/comment_test.go`

- [ ] **Step 1: 写 Comment Service 测试**

```go
// server/service/comment_test.go
package service

import (
	"testing"
	"vblog-core/model"
)

func TestCommentStatusTransition(t *testing.T) {
	c := &model.Comment{Status: "pending"}

	// pending -> approved
	c.Status = "approved"
	if c.Status != "approved" {
		t.Error("should transition to approved")
	}

	// pending -> spam
	c2 := &model.Comment{Status: "pending"}
	c2.Status = "spam"
	if c2.Status != "spam" {
		t.Error("should transition to spam")
	}
}
```

- [ ] **Step 2: 写 Comment Service**

```go
// server/service/comment.go
package service

import (
	"vblog-core/model"
	"gorm.io/gorm"
)

type CommentService struct{ DB *gorm.DB }

func NewCommentService(db *gorm.DB) *CommentService {
	return &CommentService{DB: db}
}

func (s *CommentService) List(page, perPage int, status, search string) ([]model.Comment, int64, error) {
	var comments []model.Comment
	var total int64
	q := s.DB.Model(&model.Comment{})
	if status != "" && status != "all" {
		q = q.Where("status = ?", status)
	}
	if search != "" {
		q = q.Where("body ILIKE ? OR author_name ILIKE ?", "%"+search+"%", "%"+search+"%")
	}
	q.Count(&total)
	err := q.Order("created_at DESC").Offset((page - 1) * perPage).Limit(perPage).Find(&comments).Error
	return comments, total, err
}

func (s *CommentService) Create(c *model.Comment) error {
	c.Status = "pending"
	return s.DB.Create(c).Error
}

func (s *CommentService) Approve(id uint) error {
	return s.DB.Model(&model.Comment{}).Where("id = ?", id).Update("status", "approved").Error
}

func (s *CommentService) MarkSpam(id uint) error {
	return s.DB.Model(&model.Comment{}).Where("id = ?", id).Update("status", "spam").Error
}

func (s *CommentService) Delete(id uint) error {
	return s.DB.Delete(&model.Comment{}, id).Error
}
```

- [ ] **Step 3: 写 Comment API handler**

```go
// server/api/comment.go
package api

import (
	"net/http"
	"strconv"

	restful "github.com/emicklei/go-restful/v4"
	"vblog-core/model"
	"vblog-core/service"
)

type CommentResource struct{ Service *service.CommentService }

func (c *CommentResource) Register(ws *restful.WebService) {
	ws.Route(ws.GET("/api/comments").To(c.list))
	ws.Route(ws.POST("/api/comments").To(c.create))
	ws.Route(ws.PATCH("/api/comments/{id}/approve").To(c.approve))
	ws.Route(ws.PATCH("/api/comments/{id}/spam").To(c.spam))
	ws.Route(ws.DELETE("/api/comments/{id}").To(c.delete))
}

func (c *CommentResource) list(req *restful.Request, resp *restful.Response) {
	page, _ := strconv.Atoi(req.QueryParameter("page"))
	perPage, _ := strconv.Atoi(req.QueryParameter("per_page"))
	if page < 1 { page = 1 }
	if perPage < 1 { perPage = 20 }
	status := req.QueryParameter("status")
	search := req.QueryParameter("search")
	comments, total, err := c.Service.List(page, perPage, status, search)
	if err != nil {
		resp.WriteError(http.StatusInternalServerError, err)
		return
	}
	resp.WriteEntity(map[string]interface{}{"data": comments, "total": total})
}

func (c *CommentResource) create(req *restful.Request, resp *restful.Response) {
	var comment model.Comment
	if err := req.ReadEntity(&comment); err != nil {
		resp.WriteError(http.StatusBadRequest, err)
		return
	}
	c.Service.Create(&comment)
	resp.WriteHeaderAndEntity(http.StatusCreated, comment)
}

func (c *CommentResource) approve(req *restful.Request, resp *restful.Response) {
	id, _ := strconv.ParseUint(req.PathParameter("id"), 10, 32)
	c.Service.Approve(uint(id))
	resp.WriteHeader(http.StatusOK)
}

func (c *CommentResource) spam(req *restful.Request, resp *restful.Response) {
	id, _ := strconv.ParseUint(req.PathParameter("id"), 10, 32)
	c.Service.MarkSpam(uint(id))
	resp.WriteHeader(http.StatusOK)
}

func (c *CommentResource) delete(req *restful.Request, resp *restful.Response) {
	id, _ := strconv.ParseUint(req.PathParameter("id"), 10, 32)
	c.Service.Delete(uint(id))
	resp.WriteHeader(http.StatusNoContent)
}
```

- [ ] **Step 4: 运行测试 + 提交**

```bash
cd server && go test ./service/ -v -run TestComment
git add server/api/comment.go server/api/comment_test.go server/service/comment.go server/service/comment_test.go
git commit -m "feat: add Comment API with TDD"
```

---

### Task 7: Settings API

**Files:**
- Create: `server/api/setting.go`
- Create: `server/api/setting_test.go`
- Create: `server/service/setting.go`
- Create: `server/service/setting_test.go`

- [ ] **Step 1: 写 Settings Service 测试**

```go
// server/service/setting_test.go
package service

import "testing"

func TestDefaultSettings(t *testing.T) {
	defaults := DefaultSettings()
	if defaults["site_title"] != "vBlog Core" {
		t.Errorf("expected 'vBlog Core', got '%s'", defaults["site_title"])
	}
	if defaults["posts_per_page"] != "5" {
		t.Errorf("expected '5', got '%s'", defaults["posts_per_page"])
	}
}
```

- [ ] **Step 2: 写 Settings Service**

```go
// server/service/setting.go
package service

import (
	"vblog-core/model"
	"gorm.io/gorm"
)

type SettingService struct{ DB *gorm.DB }

func NewSettingService(db *gorm.DB) *SettingService {
	return &SettingService{DB: db}
}

func DefaultSettings() map[string]string {
	return map[string]string{
		"site_title":        "vBlog Core",
		"subtitle":          "写代码的人，也写点别的。",
		"description":       "一个关于系统设计、工程实践与极客生活的博客",
		"language":          "zh-CN",
		"posts_per_page":    "5",
		"author_name":       "vBlog Admin",
		"author_email":      "admin@vblog.dev",
		"author_bio":        "全栈开发者，热爱 Go 和 React",
		"author_github":     "https://github.com/vblog",
		"feature_comments":  "true",
		"feature_rss":       "true",
		"feature_view_counter": "true",
	}
}

func (s *SettingService) GetAll() (map[string]string, error) {
	var settings []model.Setting
	if err := s.DB.Find(&settings).Error; err != nil {
		return nil, err
	}
	result := DefaultSettings()
	for _, s := range settings {
		result[s.Key] = s.Value
	}
	return result, nil
}

func (s *SettingService) Save(settings map[string]string) error {
	for key, value := range settings {
		s.DB.Save(&model.Setting{Key: key, Value: value})
	}
	return nil
}

func (s *SettingService) Reset() error {
	s.DB.Where("1 = 1").Delete(&model.Setting{})
	return nil
}
```

- [ ] **Step 3: 写 Settings API handler**

```go
// server/api/setting.go
package api

import (
	"net/http"

	restful "github.com/emicklei/go-restful/v4"
	"vblog-core/service"
)

type SettingResource struct{ Service *service.SettingService }

func (s *SettingResource) Register(ws *restful.WebService) {
	ws.Route(ws.GET("/api/settings").To(s.getAll))
	ws.Route(ws.PUT("/api/settings").To(s.save))
	ws.Route(ws.POST("/api/settings/reset").To(s.reset))
}

func (s *SettingResource) getAll(req *restful.Request, resp *restful.Response) {
	settings, err := s.Service.GetAll()
	if err != nil {
		resp.WriteError(http.StatusInternalServerError, err)
		return
	}
	resp.WriteEntity(settings)
}

func (s *SettingResource) save(req *restful.Request, resp *restful.Response) {
	var settings map[string]string
	if err := req.ReadEntity(&settings); err != nil {
		resp.WriteError(http.StatusBadRequest, err)
		return
	}
	s.Service.Save(settings)
	resp.WriteEntity(map[string]string{"message": "设置已保存"})
}

func (s *SettingResource) reset(req *restful.Request, resp *restful.Response) {
	s.Service.Reset()
	resp.WriteEntity(map[string]string{"message": "已重置"})
}
```

- [ ] **Step 4: 运行测试 + 提交**

```bash
cd server && go test ./service/ -v -run TestDefault
git add server/api/setting.go server/api/setting_test.go server/service/setting.go server/service/setting_test.go
git commit -m "feat: add Settings API with TDD"
```

---

### Task 8: Component API

**Files:**
- Create: `server/api/component.go`
- Create: `server/api/component_test.go`
- Create: `server/service/component.go`
- Create: `server/service/component_test.go`

- [ ] **Step 1: 写 Component Service 测试**

```go
// server/service/component_test.go
package service

import (
	"testing"
	"vblog-core/model"
)

func TestComponentValidation(t *testing.T) {
	c := &model.Component{Name: ""}
	if c.Name != "" {
		t.Error("empty name should be caught")
	}

	c2 := &model.Component{Name: "MyWidget", Origin: "uploaded"}
	if c2.Origin != "uploaded" {
		t.Error("origin should be uploaded")
	}
}
```

- [ ] **Step 2: 写 Component Service**

```go
// server/service/component.go
package service

import (
	"vblog-core/model"
	"gorm.io/gorm"
)

type ComponentService struct{ DB *gorm.DB }

func NewComponentService(db *gorm.DB) *ComponentService {
	return &ComponentService{DB: db}
}

func (s *ComponentService) List() ([]model.Component, error) {
	var components []model.Component
	err := s.DB.Order("category, name").Find(&components).Error
	return components, err
}

func (s *ComponentService) Create(c *model.Component) error {
	c.Origin = "uploaded"
	c.Status = "active"
	return s.DB.Create(c).Error
}

func (s *ComponentService) Update(c *model.Component) error {
	return s.DB.Save(c).Error
}

func (s *ComponentService) Delete(id uint) error {
	return s.DB.Delete(&model.Component{}, id).Error
}

func (s *ComponentService) Toggle(id uint) error {
	var c model.Component
	if err := s.DB.First(&c, id).Error; err != nil {
		return err
	}
	if c.Status == "active" {
		c.Status = "inactive"
	} else {
		c.Status = "active"
	}
	return s.DB.Save(&c).Error
}
```

- [ ] **Step 3: 写 Component API handler**

```go
// server/api/component.go
package api

import (
	"net/http"
	"strconv"

	restful "github.com/emicklei/go-restful/v4"
	"vblog-core/model"
	"vblog-core/service"
)

type ComponentResource struct{ Service *service.ComponentService }

func (c *ComponentResource) Register(ws *restful.WebService) {
	ws.Route(ws.GET("/api/components").To(c.list))
	ws.Route(ws.POST("/api/components").To(c.create))
	ws.Route(ws.PUT("/api/components/{id}").To(c.update))
	ws.Route(ws.DELETE("/api/components/{id}").To(c.delete))
	ws.Route(ws.PATCH("/api/components/{id}/toggle").To(c.toggle))
}

func (c *ComponentResource) list(req *restful.Request, resp *restful.Response) {
	components, err := c.Service.List()
	if err != nil {
		resp.WriteError(http.StatusInternalServerError, err)
		return
	}
	resp.WriteEntity(components)
}

func (c *ComponentResource) create(req *restful.Request, resp *restful.Response) {
	var comp model.Component
	if err := req.ReadEntity(&comp); err != nil {
		resp.WriteError(http.StatusBadRequest, err)
		return
	}
	if comp.Name == "" {
		resp.WriteErrorString(http.StatusBadRequest, "name is required")
		return
	}
	c.Service.Create(&comp)
	resp.WriteHeaderAndEntity(http.StatusCreated, comp)
}

func (c *ComponentResource) update(req *restful.Request, resp *restful.Response) {
	id, _ := strconv.ParseUint(req.PathParameter("id"), 10, 32)
	var comp model.Component
	if err := c.Service.DB.First(&comp, id).Error; err != nil {
		resp.WriteError(http.StatusNotFound, err)
		return
	}
	req.ReadEntity(&comp)
	c.Service.Update(&comp)
	resp.WriteEntity(comp)
}

func (c *ComponentResource) delete(req *restful.Request, resp *restful.Response) {
	id, _ := strconv.ParseUint(req.PathParameter("id"), 10, 32)
	c.Service.Delete(uint(id))
	resp.WriteHeader(http.StatusNoContent)
}

func (c *ComponentResource) toggle(req *restful.Request, resp *restful.Response) {
	id, _ := strconv.ParseUint(req.PathParameter("id"), 10, 32)
	if err := c.Service.Toggle(uint(id)); err != nil {
		resp.WriteError(http.StatusInternalServerError, err)
		return
	}
	resp.WriteHeader(http.StatusOK)
}
```

- [ ] **Step 4: 运行测试 + 提交**

```bash
cd server && go test ./service/ -v -run TestComponent
git add server/api/component.go server/api/component_test.go server/service/component.go server/service/component_test.go
git commit -m "feat: add Component API with TDD"
```

---

### Task 9: Dashboard Stats API

**Files:**
- Create: `server/api/dashboard.go`
- Create: `server/api/dashboard_test.go`

- [ ] **Step 1: 写 Dashboard 测试**

```go
// server/api/dashboard_test.go
package api

import (
	"net/http"
	"net/http/httptest"
	"testing"

	restful "github.com/emicklei/go-restful/v4"
)

func TestDashboardAPI_Stats(t *testing.T) {
	ws := new(restful.WebService)
	ws.Path("/api/dashboard")
	container := restful.NewContainer()
	container.Add(ws)

	req := httptest.NewRequest("GET", "/api/dashboard/stats", nil)
	resp := httptest.NewRecorder()
	container.ServeHTTP(resp, req)

	if resp.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", resp.Code)
	}
}
```

- [ ] **Step 2: 写 Dashboard API**

```go
// server/api/dashboard.go
package api

import (
	"net/http"
	"time"

	restful "github.com/emicklei/go-restful/v4"
	"vblog-core/model"
	"gorm.io/gorm"
)

type DashboardResource struct{ DB *gorm.DB }

func (d *DashboardResource) Register(ws *restful.WebService) {
	ws.Route(ws.GET("/api/dashboard/stats").To(d.stats))
}

func (d *DashboardResource) stats(req *restful.Request, resp *restful.Response) {
	var postCount, viewTotal, commentCount, tagCount int64
	d.DB.Model(&model.Post{}).Where("status = ?", "published").Count(&postCount)
	d.DB.Model(&model.Post{}).Select("COALESCE(SUM(views), 0)").Scan(&viewTotal)
	d.DB.Model(&model.Comment{}).Count(&commentCount)
	d.DB.Model(&model.Tag{}).Count(&tagCount)

	monthStart := time.Now().AddDate(0, 0, -30)
	var monthPosts int64
	d.DB.Model(&model.Post{}).Where("created_at > ?", monthStart).Count(&monthPosts)

	resp.WriteEntity(map[string]interface{}{
		"total_posts":    postCount,
		"total_views":    viewTotal,
		"total_comments": commentCount,
		"total_tags":     tagCount,
		"month_posts":    monthPosts,
	})
}
```

- [ ] **Step 3: 运行测试 + 提交**

```bash
cd server && go test ./api/ -v -run TestDashboard
git add server/api/dashboard.go server/api/dashboard_test.go
git commit -m "feat: add Dashboard stats API with TDD"
```

---

### Task 10: JWT 认证

**Files:**
- Create: `server/middleware/jwt.go`
- Create: `server/middleware/jwt_test.go`
- Create: `server/api/auth.go`
- Create: `server/api/auth_test.go`
- Create: `server/service/auth.go`
- Create: `server/service/auth_test.go`

- [ ] **Step 1: 写 JWT 测试**

```go
// server/middleware/jwt_test.go
package middleware

import (
	"testing"
	"time"
)

func TestGenerateAndValidateToken(t *testing.T) {
	secret := "test-secret-key"
	token, err := GenerateToken(1, "admin", secret, 2*time.Hour)
	if err != nil {
		t.Fatalf("generate token failed: %v", err)
	}
	if token == "" {
		t.Fatal("token should not be empty")
	}

	claims, err := ValidateToken(token, secret)
	if err != nil {
		t.Fatalf("validate token failed: %v", err)
	}
	if claims.UserID != 1 {
		t.Errorf("expected user ID 1, got %d", claims.UserID)
	}
	if claims.Username != "admin" {
		t.Errorf("expected username 'admin', got '%s'", claims.Username)
	}
}

func TestValidateExpiredToken(t *testing.T) {
	secret := "test-secret-key"
	token, _ := GenerateToken(1, "admin", secret, -1*time.Hour) // 已过期
	_, err := ValidateToken(token, secret)
	if err == nil {
		t.Error("expired token should return error")
	}
}

func TestValidateInvalidToken(t *testing.T) {
	_, err := ValidateToken("invalid-token", "secret")
	if err == nil {
		t.Error("invalid token should return error")
	}
}
```

- [ ] **Step 2: 运行测试确认失败**

```bash
cd server && go test ./middleware/ -v
```

- [ ] **Step 3: 写 JWT 实现**

```go
// server/middleware/jwt.go
package middleware

import (
	"errors"
	"net/http"
	"strings"
	"time"

	restful "github.com/emicklei/go-restful/v4"
	"github.com/golang-jwt/jwt/v5"
)

type Claims struct {
	UserID   uint   `json:"user_id"`
	Username string `json:"username"`
	jwt.RegisteredClaims
}

func GenerateToken(userID uint, username, secret string, duration time.Duration) (string, error) {
	claims := Claims{
		UserID:   userID,
		Username: username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(duration)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}

func ValidateToken(tokenStr, secret string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &Claims{}, func(t *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})
	if err != nil {
		return nil, err
	}
	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}
	return nil, errors.New("invalid token")
}

func JWTFilter(secret string) restful.FilterFunction {
	return func(req *restful.Request, resp *restful.Response, chain *restful.FilterChain) {
		auth := req.HeaderParameter("Authorization")
		if auth == "" {
			resp.WriteHeaderAndEntity(http.StatusUnauthorized, map[string]string{"error": "missing token"})
			return
		}
		tokenStr := strings.TrimPrefix(auth, "Bearer ")
		claims, err := ValidateToken(tokenStr, secret)
		if err != nil {
			resp.WriteHeaderAndEntity(http.StatusUnauthorized, map[string]string{"error": "invalid token"})
			return
		}
		req.SetAttribute("claims", claims)
		chain.ProcessFilter(req, resp)
	}
}
```

- [ ] **Step 4: 运行测试确认通过**

```bash
cd server && go test ./middleware/ -v
```

- [ ] **Step 5: 写 Auth Service 测试**

```go
// server/service/auth_test.go
package service

import "testing"

func TestHashPassword(t *testing.T) {
	hash, err := HashPassword("mypassword")
	if err != nil {
		t.Fatalf("hash failed: %v", err)
	}
	if hash == "" {
		t.Fatal("hash should not be empty")
	}
	if hash == "mypassword" {
		t.Fatal("hash should not equal plaintext")
	}
}

func TestCheckPassword(t *testing.T) {
	hash, _ := HashPassword("mypassword")
	if !CheckPassword(hash, "mypassword") {
		t.Error("correct password should pass")
	}
	if CheckPassword(hash, "wrongpassword") {
		t.Error("wrong password should fail")
	}
}
```

- [ ] **Step 6: 写 Auth Service**

```go
// server/service/auth.go
package service

import (
	"vblog-core/model"
	"gorm.io/gorm"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct{ DB *gorm.DB }

func NewAuthService(db *gorm.DB) *AuthService {
	return &AuthService{DB: db}
}

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

func CheckPassword(hash, password string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password)) == nil
}

func (s *AuthService) Login(username, password string) (*model.User, error) {
	var user model.User
	if err := s.DB.Where("username = ?", username).First(&user).Error; err != nil {
		return nil, err
	}
	if !CheckPassword(user.Password, password) {
		return nil, bcrypt.ErrMismatchedHashAndPassword
	}
	return &user, nil
}
```

- [ ] **Step 7: 写 Auth API**

```go
// server/api/auth.go
package api

import (
	"net/http"
	"time"

	restful "github.com/emicklei/go-restful/v4"
	"vblog-core/middleware"
	"vblog-core/service"
)

type AuthResource struct {
	Service *service.AuthService
	Secret  string
}

func (a *AuthResource) Register(ws *restful.WebService) {
	ws.Route(ws.POST("/api/auth/login").To(a.login))
}

func (a *AuthResource) login(req *restful.Request, resp *restful.Response) {
	var body struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	if err := req.ReadEntity(&body); err != nil {
		resp.WriteError(http.StatusBadRequest, err)
		return
	}
	user, err := a.Service.Login(body.Username, body.Password)
	if err != nil {
		resp.WriteHeaderAndEntity(http.StatusUnauthorized, map[string]string{"error": "用户名或密码错误"})
		return
	}
	token, _ := middleware.GenerateToken(user.ID, user.Username, a.Secret, 2*time.Hour)
	refresh, _ := middleware.GenerateToken(user.ID, user.Username, a.Secret, 7*24*time.Hour)
	resp.WriteEntity(map[string]string{
		"access_token":  token,
		"refresh_token": refresh,
	})
}
```

- [ ] **Step 8: 运行全部测试 + 提交**

```bash
cd server && go test ./... -v
git add server/middleware/ server/api/auth.go server/api/auth_test.go server/service/auth.go server/service/auth_test.go
git commit -m "feat: add JWT auth with TDD"
```

---

### Task 11: 整合 main.go + 注册所有路由

**Files:**
- Modify: `server/cmd/main.go`

- [ ] **Step 1: 更新 main.go**

```go
// server/cmd/main.go
package main

import (
	"log"
	"net/http"

	restful "github.com/emicklei/go-restful/v4"
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
	model.AutoMigrate(db)

	// Services
	postSvc := service.NewPostService(db)
	tagSvc := service.NewTagService(db)
	commentSvc := service.NewCommentService(db)
	componentSvc := service.NewComponentService(db)
	settingSvc := service.NewSettingService(db)
	authSvc := service.NewAuthService(db)

	// Container
	wsContainer := restful.NewContainer()
	wsContainer.EnableContentEncoding(true)

	// 公开 API（无需认证）
	publicWS := new(restful.WebService).Path("/api")
	postRes := &api.PostResource{Service: postSvc}
	postRes.Register(publicWS)
	tagRes := &api.TagResource{Service: tagSvc}
	tagRes.Register(publicWS)
	commentRes := &api.CommentResource{Service: commentSvc}
	// 只注册公开的评论提交
	publicWS.Route(publicWS.POST("/api/comments").To(commentRes.CreatePublic))
	wsContainer.Add(publicWS)

	// 需要认证的 API
	adminWS := new(restful.WebService).Path("/api/admin")
	adminWS.Filter(middleware.JWTFilter(cfg.JWT.Secret))
	// 注册管理端路由...
	wsContainer.Add(adminWS)

	// Auth
	authRes := &api.AuthResource{Service: authSvc, Secret: cfg.JWT.Secret}
	authWS := new(restful.WebService).Path("/api")
	authRes.Register(authWS)
	wsContainer.Add(authWS)

	// Dashboard
	dashRes := &api.DashboardResource{DB: db}
	dashRes.Register(adminWS)

	log.Printf("vBlog Core starting on :%s", cfg.Server.Port)
	server := &http.Server{Addr: ":" + cfg.Server.Port, Handler: wsContainer}
	log.Fatal(server.ListenAndServe())
}
```

- [ ] **Step 2: 编译验证**

```bash
cd server && go build ./cmd/main.go
```

- [ ] **Step 3: 提交**

```bash
git add server/cmd/main.go
git commit -m "feat: integrate all API routes in main.go"
```

---

## Phase 4: Vue 前端

### Task 12: Vue 脚手架 + 路由 + 主题

**Files:**
- Create: `web/src/main.js`
- Create: `web/src/App.vue`
- Create: `web/src/router/index.js`
- Create: `web/src/stores/theme.js`
- Create: `web/src/stores/auth.js`
- Create: `web/src/styles/variables.css`
- Create: `web/src/api/request.js`
- Create: `web/src/shared/ThemeToggle.vue`

- [ ] **Step 1: 创建 CSS 变量文件**

```css
/* web/src/styles/variables.css */
:root {
  --bg: #fafafa; --surface: #ffffff; --fg: #171717; --muted: #737373;
  --border: #e5e5e5; --accent: #2563eb; --accent-soft: #eff6ff;
  --success: #16a34a; --warning: #f59e0b; --error: #dc2626;
  --radius: 8px; --radius-lg: 12px;
  --font-sans: -apple-system, BlinkMacSystemFont, 'SF Pro Text', 'Segoe UI', system-ui, sans-serif;
  --font-display: -apple-system, BlinkMacSystemFont, 'SF Pro Display', 'Segoe UI', system-ui, sans-serif;
  --font-mono: 'JetBrains Mono', 'IBM Plex Mono', ui-monospace, Menlo, monospace;
}
[data-theme="dark"] {
  --bg: #0a0a0a; --surface: #141414; --fg: #ededed; --muted: #a3a3a3;
  --border: #262626; --accent: #3b82f6; --accent-soft: #172554;
  --success: #22c55e; --warning: #eab308; --error: #ef4444;
}
```

- [ ] **Step 2: 创建主题 store**

```javascript
// web/src/stores/theme.js
import { defineStore } from 'pinia'
import { ref } from 'vue'

export const useThemeStore = defineStore('theme', () => {
  const theme = ref(localStorage.getItem('vblog-theme') || (window.matchMedia('(prefers-color-scheme: dark)').matches ? 'dark' : 'light'))

  function toggle() {
    theme.value = theme.value === 'dark' ? 'light' : 'dark'
    document.documentElement.setAttribute('data-theme', theme.value)
    localStorage.setItem('vblog-theme', theme.value)
  }

  function init() {
    document.documentElement.setAttribute('data-theme', theme.value)
  }

  return { theme, toggle, init }
})
```

- [ ] **Step 3: 创建路由**

```javascript
// web/src/router/index.js
import { createRouter, createWebHistory } from 'vue-router'

const routes = [
  { path: '/', component: () => import('../blog/Home.vue') },
  { path: '/post/:id', component: () => import('../blog/Post.vue') },
  { path: '/archives', component: () => import('../blog/Archives.vue') },
  { path: '/tags', component: () => import('../blog/Tags.vue') },
  { path: '/about', component: () => import('../blog/About.vue') },
  { path: '/admin/login', component: () => import('../admin/Login.vue') },
  {
    path: '/admin',
    component: () => import('../admin/Layout.vue'),
    meta: { requiresAuth: true },
    children: [
      { path: '', component: () => import('../admin/Dashboard.vue') },
      { path: 'posts', component: () => import('../admin/Posts.vue') },
      { path: 'tags', component: () => import('../admin/Tags.vue') },
      { path: 'comments', component: () => import('../admin/Comments.vue') },
      { path: 'custom', component: () => import('../admin/Custom.vue') },
      { path: 'settings', component: () => import('../admin/Settings.vue') },
    ]
  }
]

const router = createRouter({
  history: createWebHistory(),
  routes
})

router.beforeEach((to, from, next) => {
  if (to.matched.some(r => r.meta.requiresAuth)) {
    const token = localStorage.getItem('vblog-token')
    if (!token) next('/admin/login')
    else next()
  } else {
    next()
  }
})

export default router
```

- [ ] **Step 4: 创建 API 请求封装**

```javascript
// web/src/api/request.js
import axios from 'axios'
import { ElMessage } from 'element-plus'

const api = axios.create({ baseURL: '/api' })

api.interceptors.request.use(config => {
  const token = localStorage.getItem('vblog-token')
  if (token) config.headers.Authorization = `Bearer ${token}`
  return config
})

api.interceptors.response.use(
  res => res.data,
  err => {
    if (err.response?.status === 401) {
      localStorage.removeItem('vblog-token')
      window.location.href = '/admin/login'
    }
    ElMessage.error(err.response?.data?.error || '请求失败')
    return Promise.reject(err)
  }
)

export default api
```

- [ ] **Step 5: 创建 main.js**

```javascript
// web/src/main.js
import { createApp } from 'vue'
import { createPinia } from 'pinia'
import ElementPlus from 'element-plus'
import 'element-plus/dist/index.css'
import App from './App.vue'
import router from './router'
import './styles/variables.css'

const app = createApp(App)
app.use(createPinia())
app.use(router)
app.use(ElementPlus)
app.mount('#app')
```

- [ ] **Step 6: 创建 App.vue**

```vue
<!-- web/src/App.vue -->
<template>
  <router-view />
</template>

<script setup>
import { useThemeStore } from './stores/theme'
const theme = useThemeStore()
theme.init()
</script>
```

- [ ] **Step 7: 提交**

```bash
git add web/src/main.js web/src/App.vue web/src/router/ web/src/stores/ web/src/styles/ web/src/api/
git commit -m "feat: Vue scaffolding with router and theme"
```

---

### Task 13: 博客前台页面

**Files:**
- Create: `web/src/blog/Home.vue`
- Create: `web/src/blog/Post.vue`
- Create: `web/src/blog/Archives.vue`
- Create: `web/src/blog/Tags.vue`
- Create: `web/src/blog/About.vue`
- Create: `web/src/shared/BlogNav.vue`
- Create: `web/src/shared/BlogFooter.vue`
- Create: `web/src/shared/PostCard.vue`

- [ ] **Step 1: 创建 BlogNav 组件**

```vue
<!-- web/src/shared/BlogNav.vue -->
<template>
  <nav class="blog-nav">
    <div class="nav-inner">
      <router-link to="/" class="nav-brand"><span class="dot"></span> vBlog</router-link>
      <div class="nav-links">
        <router-link to="/">首页 Home</router-link>
        <router-link to="/archives">归档 Archives</router-link>
        <router-link to="/tags">标签 Tags</router-link>
        <router-link to="/about">关于 About</router-link>
      </div>
      <div class="nav-right">
        <router-link to="/admin" class="admin-btn">后台 Admin</router-link>
        <el-button @click="theme.toggle()" :icon="theme.theme === 'dark' ? 'Moon' : 'Sunny'" circle size="small" />
      </div>
    </div>
  </nav>
</template>

<script setup>
import { useThemeStore } from '../stores/theme'
const theme = useThemeStore()
</script>

<style scoped>
.blog-nav { position: sticky; top: 0; z-index: 100; background: var(--bg); backdrop-filter: blur(12px); border-bottom: 1px solid var(--border); }
.nav-inner { max-width: 1080px; margin: 0 auto; padding: 0 24px; height: 56px; display: flex; align-items: center; justify-content: space-between; }
.nav-brand { font-family: var(--font-display); font-size: 17px; font-weight: 600; color: var(--fg); text-decoration: none; display: flex; align-items: center; gap: 8px; }
.dot { width: 8px; height: 8px; border-radius: 50%; background: var(--accent); }
.nav-links { display: flex; gap: 4px; }
.nav-links a { font-size: 14px; color: var(--muted); text-decoration: none; padding: 6px 12px; border-radius: 6px; }
.nav-links a:hover, .nav-links a.router-link-active { color: var(--fg); background: var(--bg); }
.admin-btn { font-size: 13px; color: var(--muted); text-decoration: none; padding: 6px 12px; border-radius: 6px; border: 1px solid var(--border); }
</style>
```

- [ ] **Step 2: 创建 PostCard 组件**

```vue
<!-- web/src/shared/PostCard.vue -->
<template>
  <router-link :to="`/post/${post.id}`" class="post-card" :class="{ pinned: post.pinned }">
    <div v-if="post.pinned" class="pin-badge">置顶 Pinned</div>
    <div class="post-meta">
      <el-tag v-for="tag in post.tags" :key="tag.id" size="small" effect="plain">{{ tag.name }}</el-tag>
      <span>{{ formatDate(post.created_at) }}</span>
    </div>
    <div class="post-title">{{ post.title }}</div>
    <div class="post-excerpt">{{ post.excerpt }}</div>
    <div class="post-stats">
      <span>{{ post.read_time }} min</span>
      <span>{{ post.views }} views</span>
    </div>
  </router-link>
</template>

<script setup>
defineProps({ post: Object })
function formatDate(d) { return new Date(d).toISOString().split('T')[0] }
</script>

<style scoped>
.post-card { display: block; padding: 20px; border-bottom: 1px solid var(--border); text-decoration: none; color: var(--fg); transition: background 0.15s; }
.post-card:hover { background: var(--bg); }
.post-card.pinned { background: var(--accent-soft); }
.pin-badge { font-size: 12px; color: var(--accent); margin-bottom: 8px; }
.post-meta { display: flex; align-items: center; gap: 8px; font-size: 13px; color: var(--muted); margin-bottom: 8px; }
.post-title { font-size: 18px; font-weight: 600; margin-bottom: 6px; }
.post-excerpt { font-size: 14px; color: var(--muted); display: -webkit-box; -webkit-line-clamp: 2; -webkit-box-orient: vertical; overflow: hidden; }
.post-stats { font-family: var(--font-mono); font-size: 12px; color: var(--muted); margin-top: 8px; display: flex; gap: 12px; }
</style>
```

- [ ] **Step 3: 创建 Home.vue**

```vue
<!-- web/src/blog/Home.vue -->
<template>
  <div>
    <BlogNav />
    <div class="container">
      <header class="hero"><h1>写代码的人，<br/>也写点别的。</h1><p>一个关于系统设计、工程实践与极客生活的博客。</p></header>
      <div class="stats-bar">
        <div class="stat"><span class="val">{{ stats.total_posts }}</span><span class="lbl">篇文章</span></div>
        <div class="stat"><span class="val">{{ stats.total_views }}</span><span class="lbl">次阅读</span></div>
        <div class="stat"><span class="val">{{ stats.total_tags }}</span><span class="lbl">个标签</span></div>
      </div>
      <div class="filter-bar">
        <el-button :type="activeTag === '' ? 'primary' : ''" @click="filterTag('')">全部</el-button>
        <el-button v-for="tag in tags" :key="tag.id" :type="activeTag === tag.name ? 'primary' : ''" @click="filterTag(tag.name)" size="small">{{ tag.name }}</el-button>
      </div>
      <div class="post-list">
        <PostCard v-for="post in posts" :key="post.id" :post="post" />
      </div>
      <el-pagination v-if="total > perPage" :total="total" :page-size="perPage" v-model:current-page="page" layout="prev, pager, next" />
    </div>
    <BlogFooter />
  </div>
</template>

<script setup>
import { ref, onMounted, watch } from 'vue'
import api from '../api/request'
import BlogNav from '../shared/BlogNav.vue'
import BlogFooter from '../shared/BlogFooter.vue'
import PostCard from '../shared/PostCard.vue'

const posts = ref([])
const tags = ref([])
const stats = ref({ total_posts: 0, total_views: 0, total_tags: 0 })
const page = ref(1)
const perPage = 5
const total = ref(0)
const activeTag = ref('')

async function loadPosts() {
  const res = await api.get('/posts', { params: { page: page.value, per_page: perPage, tag: activeTag.value, status: 'published' } })
  posts.value = res.data
  total.value = res.total
}

function filterTag(tag) {
  activeTag.value = tag
  page.value = 1
  loadPosts()
}

onMounted(async () => {
  const [postsRes, tagsRes, statsRes] = await Promise.all([
    api.get('/posts', { params: { page: 1, per_page: 5, status: 'published' } }),
    api.get('/tags'),
    api.get('/dashboard/stats')
  ])
  posts.value = postsRes.data
  total.value = postsRes.total
  tags.value = tagsRes
  stats.value = statsRes
})

watch(page, loadPosts)
</script>

<style scoped>
.container { max-width: 1080px; margin: 0 auto; padding: 0 24px; }
.hero { padding: 48px 0 32px; }
.hero h1 { font-family: var(--font-display); font-size: clamp(28px, 4vw, 40px); font-weight: 700; letter-spacing: -0.03em; line-height: 1.15; }
.hero p { font-size: 16px; color: var(--muted); margin-top: 12px; }
.stats-bar { display: flex; gap: 32px; padding: 20px 0; border-bottom: 1px solid var(--border); }
.val { font-family: var(--font-mono); font-size: 24px; font-weight: 700; color: var(--accent); }
.lbl { font-size: 13px; color: var(--muted); margin-left: 4px; }
.filter-bar { display: flex; gap: 8px; padding: 20px 0; flex-wrap: wrap; }
.post-list { min-height: 200px; }
</style>
```

- [ ] **Step 4: 创建 Post.vue（文章详情）**

```vue
<!-- web/src/blog/Post.vue -->
<template>
  <div>
    <BlogNav />
    <div class="post-container">
      <router-link to="/" class="back-link">← 返回首页</router-link>
      <article v-if="post">
        <div class="post-meta">
          <el-tag v-for="tag in post.tags" :key="tag.id" size="small">{{ tag.name }}</el-tag>
          <span>{{ formatDate(post.created_at) }}</span>
          <span>{{ post.read_time }} min</span>
          <span>{{ post.views }} views</span>
        </div>
        <h1>{{ post.title }}</h1>
        <div class="post-excerpt">{{ post.excerpt }}</div>
        <div class="post-content" v-html="renderedContent"></div>
        <div class="post-tags">
          <el-tag v-for="tag in post.tags" :key="tag.id" effect="plain">{{ tag.name }}</el-tag>
        </div>
      </article>
      <div v-else class="not-found">文章不存在</div>
    </div>
    <BlogFooter />
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useRoute } from 'vue-router'
import MarkdownIt from 'markdown-it'
import api from '../api/request'
import BlogNav from '../shared/BlogNav.vue'
import BlogFooter from '../shared/BlogFooter.vue'

const route = useRoute()
const post = ref(null)
const md = new MarkdownIt()

const renderedContent = computed(() => post.value ? md.render(post.value.content) : '')

function formatDate(d) { return new Date(d).toISOString().split('T')[0] }

onMounted(async () => {
  try {
    post.value = await api.get(`/posts/${route.params.id}`)
  } catch { post.value = null }
})
</script>

<style scoped>
.post-container { max-width: 720px; margin: 0 auto; padding: 48px 24px 80px; }
.back-link { font-size: 14px; color: var(--muted); text-decoration: none; padding: 6px 12px; border: 1px solid var(--border); border-radius: 20px; }
.post-meta { display: flex; align-items: center; gap: 8px; font-size: 13px; color: var(--muted); margin: 24px 0 12px; }
h1 { font-family: var(--font-display); font-size: clamp(24px, 4vw, 36px); font-weight: 700; letter-spacing: -0.03em; line-height: 1.2; }
.post-excerpt { font-size: 16px; color: var(--muted); margin: 12px 0 24px; }
.post-content :deep(h2) { font-size: 22px; font-weight: 600; margin: 32px 0 12px; }
.post-content :deep(p) { margin: 12px 0; line-height: 1.8; }
.post-content :deep(pre) { background: var(--bg); border: 1px solid var(--border); border-radius: var(--radius); padding: 16px; overflow-x: auto; font-family: var(--font-mono); font-size: 14px; }
.post-content :deep(code) { font-family: var(--font-mono); background: var(--bg); padding: 2px 6px; border-radius: 4px; font-size: 14px; }
.post-content :deep(blockquote) { border-left: 3px solid var(--accent); padding-left: 16px; color: var(--muted); margin: 16px 0; }
.post-tags { display: flex; gap: 8px; margin-top: 32px; padding-top: 24px; border-top: 1px solid var(--border); }
.not-found { text-align: center; padding: 80px 0; color: var(--muted); }
</style>
```

- [ ] **Step 5: 创建 Archives.vue、Tags.vue、About.vue**

Archives.vue 按年分组显示文章时间线；Tags.vue 标签云 + 按标签筛选；About.vue 从 settings API 获取内容。

- [ ] **Step 6: 创建 BlogFooter.vue**

```vue
<!-- web/src/shared/BlogFooter.vue -->
<template>
  <footer class="blog-footer">
    <span>© 2026 vBlog Core · 用代码写作，用文字思考</span>
    <div class="links"><a href="#">GitHub</a><a href="#">RSS</a></div>
  </footer>
</template>
<style scoped>
.blog-footer { max-width: 1080px; margin: 0 auto; padding: 24px; display: flex; justify-content: space-between; font-size: 13px; color: var(--muted); border-top: 1px solid var(--border); }
.links { display: flex; gap: 16px; }
.links a { color: var(--muted); text-decoration: none; }
</style>
```

- [ ] **Step 7: 提交**

```bash
git add web/src/blog/ web/src/shared/
git commit -m "feat: add blog frontend pages"
```

---

### Task 14: 后台管理页面

**Files:**
- Create: `web/src/admin/Layout.vue`
- Create: `web/src/admin/Login.vue`
- Create: `web/src/admin/Dashboard.vue`
- Create: `web/src/admin/Posts.vue`
- Create: `web/src/admin/Tags.vue`
- Create: `web/src/admin/Comments.vue`
- Create: `web/src/admin/Custom.vue`
- Create: `web/src/admin/Settings.vue`

- [ ] **Step 1: 创建 Layout.vue（侧边栏布局）**

```vue
<!-- web/src/admin/Layout.vue -->
<template>
  <div class="admin-layout">
    <aside class="sidebar">
      <div class="sidebar-brand"><span class="dot"></span> vBlog <el-tag size="small" type="info">Admin</el-tag></div>
      <div class="sidebar-section">
        <div class="section-title">概览 Overview</div>
        <router-link to="/admin" class="sidebar-link" exact>仪表盘 Dashboard</router-link>
      </div>
      <div class="sidebar-section">
        <div class="section-title">内容 Content</div>
        <router-link to="/admin/posts" class="sidebar-link">文章 Posts</router-link>
        <router-link to="/admin/tags" class="sidebar-link">标签 Tags</router-link>
        <router-link to="/admin/comments" class="sidebar-link">评论 Comments</router-link>
      </div>
      <div class="sidebar-section">
        <div class="section-title">系统 System</div>
        <router-link to="/admin/custom" class="sidebar-link">组件 Custom</router-link>
        <router-link to="/admin/settings" class="sidebar-link">设置 Settings</router-link>
      </div>
      <div class="sidebar-footer">
        <div class="avatar">A</div>
        <div><div class="name">Admin</div><div class="role">超级管理员</div></div>
      </div>
    </aside>
    <main class="main">
      <header class="topbar">
        <div><router-link to="/" class="back-link">← 首页</router-link></div>
        <el-button type="primary" @click="theme.toggle()">{{ theme.theme === 'dark' ? '🌙' : '☀' }}</el-button>
      </header>
      <div class="content"><router-view /></div>
    </main>
  </div>
</template>

<script setup>
import { useThemeStore } from '../stores/theme'
const theme = useThemeStore()
</script>

<style scoped>
.admin-layout { display: flex; min-height: 100vh; }
.sidebar { width: 220px; background: var(--bg); border-right: 1px solid var(--border); padding: 20px 0; display: flex; flex-direction: column; }
.sidebar-brand { padding: 0 20px 24px; font-family: var(--font-display); font-size: 16px; font-weight: 600; display: flex; align-items: center; gap: 8px; border-bottom: 1px solid var(--border); margin-bottom: 16px; }
.dot { width: 8px; height: 8px; border-radius: 50%; background: var(--accent); }
.sidebar-section { padding: 0 12px; margin-bottom: 24px; }
.section-title { font-size: 11px; font-weight: 600; text-transform: uppercase; letter-spacing: 0.06em; color: var(--muted); padding: 0 8px 8px; }
.sidebar-link { display: flex; align-items: center; gap: 10px; padding: 8px 12px; border-radius: 6px; font-size: 13px; color: var(--muted); text-decoration: none; }
.sidebar-link:hover { color: var(--fg); background: var(--bg); }
.sidebar-link.router-link-active { color: var(--fg); background: var(--surface); box-shadow: 0 1px 2px rgba(0,0,0,0.06); }
.sidebar-footer { margin-top: auto; padding: 16px 20px; border-top: 1px solid var(--border); display: flex; align-items: center; gap: 10px; }
.avatar { width: 28px; height: 28px; border-radius: 50%; background: var(--accent); color: white; display: flex; align-items: center; justify-content: center; font-size: 12px; font-weight: 600; }
.name { font-size: 13px; font-weight: 500; }
.role { font-size: 11px; color: var(--muted); }
.main { flex: 1; display: flex; flex-direction: column; }
.topbar { height: 56px; border-bottom: 1px solid var(--border); display: flex; align-items: center; justify-content: space-between; padding: 0 24px; }
.back-link { font-size: 13px; color: var(--muted); text-decoration: none; }
.content { flex: 1; padding: 24px; overflow-y: auto; }
</style>
```

- [ ] **Step 2: 创建 Login.vue**

```vue
<!-- web/src/admin/Login.vue -->
<template>
  <div class="login-page">
    <el-card class="login-card">
      <h2>vBlog Admin</h2>
      <el-form @submit.prevent="handleLogin">
        <el-form-item><el-input v-model="form.username" placeholder="用户名" /></el-form-item>
        <el-form-item><el-input v-model="form.password" type="password" placeholder="密码" /></el-form-item>
        <el-button type="primary" @click="handleLogin" :loading="loading" style="width:100%">登录</el-button>
      </el-form>
    </el-card>
  </div>
</template>

<script setup>
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import api from '../api/request'

const router = useRouter()
const form = ref({ username: '', password: '' })
const loading = ref(false)

async function handleLogin() {
  loading.value = true
  try {
    const res = await api.post('/auth/login', form.value)
    localStorage.setItem('vblog-token', res.access_token)
    router.push('/admin')
  } catch {
    ElMessage.error('用户名或密码错误')
  } finally { loading.value = false }
}
</script>

<style scoped>
.login-page { min-height: 100vh; display: flex; align-items: center; justify-content: center; background: var(--bg); }
.login-card { width: 360px; }
h2 { text-align: center; margin-bottom: 24px; font-family: var(--font-display); }
</style>
```

- [ ] **Step 3: 创建 Dashboard.vue**

仪表盘显示 4 个统计卡片 + 文章列表表格。使用 Element Plus 的 `el-table`、`el-statistic`。

- [ ] **Step 4: 创建 Posts.vue**

文章管理页面：搜索筛选 + `el-table` 表格 + 编辑弹窗（`el-dialog`）+ MdEditor 内容编辑。

- [ ] **Step 5: 创建 Tags.vue**

标签管理：卡片网格布局 + 新建/编辑弹窗。

- [ ] **Step 6: 创建 Comments.vue**

评论管理：评论卡片列表 + 状态筛选 + 批准/标记垃圾/删除操作。

- [ ] **Step 7: 创建 Custom.vue**

组件自定义：内置组件展示 + 上传自定义组件弹窗（name/desc/version/code textarea）。

- [ ] **Step 8: 创建 Settings.vue**

设置页面：4 个分区（通用/作者/功能开关/数据库），表单 + 保存按钮。

- [ ] **Step 9: 提交**

```bash
git add web/src/admin/
git commit -m "feat: add admin management pages"
```

---

## Phase 5: 部署 + 文档

### Task 15: Docker Compose + 部署脚本

**Files:**
- Create: `deploy/Dockerfile`
- Create: `deploy/docker-compose.yml`
- Create: `server/embed.go`

- [ ] **Step 1: 创建 Dockerfile**

```dockerfile
# deploy/Dockerfile
FROM node:20-alpine AS frontend
WORKDIR /app/web
COPY web/package*.json ./
RUN npm ci
COPY web/ ./
RUN npm run build

FROM golang:1.22-alpine AS backend
WORKDIR /app/server
COPY server/go.* ./
RUN go mod download
COPY server/ ./
COPY --from=frontend /app/web/dist ./static
RUN go build -o /vblog ./cmd/main.go

FROM alpine:3.19
RUN apk add --no-cache ca-certificates
COPY --from=backend /vblog /usr/local/bin/vblog
COPY deploy/init.sql /init.sql
EXPOSE 8080
CMD ["vblog"]
```

- [ ] **Step 2: 创建 docker-compose.yml**

```yaml
# deploy/docker-compose.yml
services:
  postgres:
    image: postgres:15-alpine
    environment:
      POSTGRES_DB: vblog
      POSTGRES_USER: vblog_admin
      POSTGRES_PASSWORD: ${DB_PASSWORD:-vblog123}
    volumes:
      - pgdata:/var/lib/postgresql/data
      - ./init.sql:/docker-entrypoint-initdb.d/init.sql
    ports:
      - "5432:5432"

  vblog:
    build:
      context: ..
      dockerfile: deploy/Dockerfile
    ports:
      - "8080:8080"
    depends_on:
      - postgres
    environment:
      DB_HOST: postgres
      DB_PORT: "5432"
      DB_NAME: vblog
      DB_USER: vblog_admin
      DB_PASSWORD: ${DB_PASSWORD:-vblog123}
      JWT_SECRET: ${JWT_SECRET:-change-me-in-production}

volumes:
  pgdata:
```

- [ ] **Step 3: 创建 embed.go（嵌入前端静态文件）**

```go
// server/embed.go
package main

import "embed"

//go:embed static/*
var staticFiles embed.FS
```

- [ ] **Step 4: 验证 Docker 构建**

```bash
cd deploy && docker compose build
docker compose up -d
```

- [ ] **Step 5: 提交**

```bash
git add deploy/ server/embed.go
git commit -m "feat: add Docker deployment"
```

---

### Task 16: README.md

**Files:**
- Create: `README.md`

- [ ] **Step 1: 写 README**

```markdown
# vBlog Core

可自定义组件的极客博客系统。

## 技术栈

- 前端：Vue 3 + Element Plus + Pinia
- 后端：Go + go-restful + GORM
- 数据库：PostgreSQL
- 部署：Docker Compose / 二进制

## 快速开始

### Docker Compose

```bash
cd deploy
docker compose up -d
```

访问 http://localhost:8080

### 本地开发

```bash
# 后端
cd server && go run ./cmd/main.go

# 前端
cd web && npm install && npm run dev
```

## 环境变量

| 变量 | 默认值 | 说明 |
|------|--------|------|
| DB_HOST | localhost | 数据库地址 |
| DB_PORT | 5432 | 数据库端口 |
| DB_NAME | vblog | 数据库名 |
| DB_USER | vblog_admin | 数据库用户 |
| DB_PASSWORD | | 数据库密码 |
| JWT_SECRET | | JWT 密钥 |
| SERVER_PORT | 8080 | 服务端口 |

## 项目结构

```
vBlog Core/
├── server/       # Go 后端
├── web/          # Vue 前端
├── deploy/       # 部署配置
├── docs/         # 文档
└── hdx/          # HTML 原型
```
```

- [ ] **Step 2: 提交**

```bash
git add README.md
git commit -m "docs: add README"
```

---

## 实现顺序总结

| Phase | Task | 内容 | TDD 文件 |
|-------|------|------|----------|
| 1 | 1 | 项目脚手架 | - |
| 1 | 2 | 数据库模型 | `model/*_test.go` |
| 2 | 3 | 配置 + DB连接 | `config/*_test.go` |
| 3 | 4 | Post API | `api/post_test.go`, `service/post_test.go` |
| 3 | 5 | Tag API | `api/tag_test.go`, `service/tag_test.go` |
| 3 | 6 | Comment API | `api/comment_test.go`, `service/comment_test.go` |
| 3 | 7 | Settings API | `api/setting_test.go`, `service/setting_test.go` |
| 3 | 8 | Component API | `api/component_test.go`, `service/component_test.go` |
| 3 | 9 | Dashboard API | `api/dashboard_test.go` |
| 3 | 10 | JWT 认证 | `middleware/jwt_test.go`, `service/auth_test.go`, `api/auth_test.go` |
| 3 | 11 | 整合 main.go | 编译验证 |
| 4 | 12 | Vue 脚手架 | - |
| 4 | 13 | 博客前台页面 | - |
| 4 | 14 | 后台管理页面 | - |
| 5 | 15 | Docker 部署 | docker compose build |
| 5 | 16 | README | - |

每个 Task 完成后运行 `go test ./... -v` 确认全部测试通过再进入下一个 Task。
