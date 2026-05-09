# vBlog Core 设计文档

## 1. 项目概述

vBlog Core 是一个可自定义组件的极客博客系统，面向 vibe coder 和 coding learner。核心特点：轻量部署、组件可自定义、前后端分离架构。

### MVP 范围

**包含：**
- 博客前台（首页、文章详情、归档、标签、关于）
- 后台管理（仪表盘、文章管理、标签管理、评论管理、组件自定义、设置）
- Go RESTful API + GORM + PostgreSQL
- JWT 后台认证
- 匿名评论系统
- 自定义组件上传 + iframe 沙箱隔离
- Markdown 编辑器（MdEditor）
- 二进制 + Docker Compose 双部署

**不包含（后续迭代）：**
- gRPC 双向流推送
- Wails 桌面客户端

## 2. 技术栈

| 层级 | 技术 | 版本 |
|------|------|------|
| 前端框架 | Vue 3 | 3.4+ |
| UI 组件库 | Element Plus | 2.x |
| 状态管理 | Pinia | 2.x |
| 路由 | Vue Router | 4.x |
| HTTP 客户端 | Axios | 1.x |
| Markdown 编辑器 | MdEditor | 3.x |
| Markdown 渲染 | markdown-it + highlight.js | - |
| 后端框架 | go-restful | 4.x |
| ORM | GORM | 1.x |
| 数据库 | PostgreSQL | 15+ |
| 认证 | JWT (golang-jwt) | 5.x |
| 部署 | Docker Compose / 二进制 | - |

## 3. 项目结构

```
vBlog Core/
├── server/                    # Go 后端
│   ├── cmd/
│   │   └── main.go            # 入口，启动 HTTP 服务
│   ├── api/
│   │   ├── post.go            # 文章路由和处理器
│   │   ├── tag.go             # 标签路由
│   │   ├── comment.go         # 评论路由
│   │   ├── component.go       # 组件路由
│   │   ├── setting.go         # 设置路由
│   │   ├── auth.go            # 认证路由
│   │   └── dashboard.go       # 仪表盘统计路由
│   ├── model/
│   │   ├── post.go            # 文章模型
│   │   ├── tag.go             # 标签模型
│   │   ├── comment.go         # 评论模型
│   │   ├── component.go       # 组件模型
│   │   ├── setting.go         # 设置模型
│   │   └── user.go            # 用户模型
│   ├── service/
│   │   ├── post.go            # 文章业务逻辑
│   │   ├── tag.go             # 标签业务逻辑
│   │   ├── comment.go         # 评论业务逻辑
│   │   ├── component.go       # 组件业务逻辑
│   │   └── setting.go         # 设置业务逻辑
│   ├── middleware/
│   │   ├── jwt.go             # JWT 认证中间件
│   │   └── cors.go            # CORS 中间件
│   ├── config/
│   │   └── config.go          # 配置加载
│   └── embed.go               # 嵌入 Vue 静态文件
│
├── web/                       # Vue 3 前端
│   ├── src/
│   │   ├── main.js            # Vue 入口
│   │   ├── App.vue            # 根组件
│   │   ├── router/
│   │   │   └── index.js       # 路由定义
│   │   ├── stores/
│   │   │   ├── auth.js        # 认证状态
│   │   │   └── theme.js       # 主题状态
│   │   ├── api/
│   │   │   ├── post.js        # 文章 API
│   │   │   ├── tag.js         # 标签 API
│   │   │   ├── comment.js     # 评论 API
│   │   │   ├── component.js   # 组件 API
│   │   │   ├── setting.js     # 设置 API
│   │   │   └── auth.js        # 认证 API
│   │   ├── blog/              # 博客前台页面
│   │   │   ├── Home.vue       # 首页
│   │   │   ├── Post.vue       # 文章详情
│   │   │   ├── Archives.vue   # 归档
│   │   │   ├── Tags.vue       # 标签
│   │   │   └── About.vue      # 关于
│   │   ├── admin/             # 后台管理页面
│   │   │   ├── Layout.vue     # 后台布局（侧边栏 + 顶栏）
│   │   │   ├── Dashboard.vue  # 仪表盘
│   │   │   ├── Posts.vue      # 文章管理
│   │   │   ├── Tags.vue       # 标签管理
│   │   │   ├── Comments.vue   # 评论管理
│   │   │   ├── Custom.vue     # 组件自定义
│   │   │   ├── Settings.vue   # 设置
│   │   │   └── Login.vue      # 登录页
│   │   ├── shared/            # 共享组件
│   │   │   ├── BlogNav.vue    # 博客导航栏
│   │   │   ├── BlogFooter.vue # 博客页脚
│   │   │   ├── ThemeToggle.vue# 主题切换
│   │   │   ├── PostCard.vue   # 文章卡片
│   │   │   ├── TagBadge.vue   # 标签徽章
│   │   │   └── CustomComponent.vue # 自定义组件 iframe 容器
│   │   └── styles/
│   │       └── variables.css  # CSS 变量（设计令牌）
│   ├── public/
│   ├── index.html
│   ├── vite.config.js
│   └── package.json
│
├── docs/
│   ├── PRD.md
│   └── superpowers/specs/     # 设计文档
│
├── hdx/                       # HTML 原型图（参考用）
│
├── deploy/
│   ├── Dockerfile             # 多阶段构建
│   ├── docker-compose.yml     # 完整部署
│   └── init.sql               # 数据库初始化脚本
│
├── CLAUDE.md
└── README.md
```

## 4. 数据模型

### 4.1 数据库表设计

```sql
-- 用户表（后台管理员）
CREATE TABLE users (
    id          SERIAL PRIMARY KEY,
    username    VARCHAR(50) UNIQUE NOT NULL,
    password    VARCHAR(255) NOT NULL,  -- bcrypt 哈希
    email       VARCHAR(100),
    role        VARCHAR(20) DEFAULT 'admin',
    created_at  TIMESTAMP DEFAULT NOW(),
    updated_at  TIMESTAMP DEFAULT NOW()
);

-- 文章表
CREATE TABLE posts (
    id          SERIAL PRIMARY KEY,
    title       VARCHAR(200) NOT NULL,
    content     TEXT NOT NULL,           -- Markdown 原文
    excerpt     VARCHAR(500),            -- 摘要
    status      VARCHAR(20) DEFAULT 'draft',  -- draft/published/archived
    pinned      BOOLEAN DEFAULT FALSE,
    views       INTEGER DEFAULT 0,
    read_time   INTEGER DEFAULT 0,       -- 预估阅读时间（分钟）
    author_id   INTEGER REFERENCES users(id),
    created_at  TIMESTAMP DEFAULT NOW(),
    updated_at  TIMESTAMP DEFAULT NOW()
);

-- 标签表
CREATE TABLE tags (
    id          SERIAL PRIMARY KEY,
    name        VARCHAR(50) UNIQUE NOT NULL,
    description TEXT,
    created_at  TIMESTAMP DEFAULT NOW()
);

-- 文章-标签关联表（多对多）
CREATE TABLE post_tags (
    post_id     INTEGER REFERENCES posts(id) ON DELETE CASCADE,
    tag_id      INTEGER REFERENCES tags(id) ON DELETE CASCADE,
    PRIMARY KEY (post_id, tag_id)
);

-- 评论表
CREATE TABLE comments (
    id          SERIAL PRIMARY KEY,
    post_id     INTEGER REFERENCES posts(id) ON DELETE CASCADE,
    author_name VARCHAR(50) NOT NULL,
    author_email VARCHAR(100),
    body        TEXT NOT NULL,
    status      VARCHAR(20) DEFAULT 'pending',  -- pending/approved/spam
    created_at  TIMESTAMP DEFAULT NOW()
);

-- 组件表
CREATE TABLE components (
    id          SERIAL PRIMARY KEY,
    name        VARCHAR(100) NOT NULL,    -- 组件名，如 PostList
    description TEXT,
    version     VARCHAR(20) DEFAULT '0.1.0',
    code        TEXT,                     -- 组件代码（自定义组件）
    category    VARCHAR(50),              -- layout/content/custom
    origin      VARCHAR(20) DEFAULT 'built-in',  -- built-in/uploaded
    status      VARCHAR(20) DEFAULT 'active',    -- active/inactive
    created_at  TIMESTAMP DEFAULT NOW(),
    updated_at  TIMESTAMP DEFAULT NOW()
);

-- 设置表（键值对）
CREATE TABLE settings (
    key         VARCHAR(100) PRIMARY KEY,
    value       TEXT,
    updated_at  TIMESTAMP DEFAULT NOW()
);
```

### 4.2 GORM 模型

```go
// model/post.go
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
```

## 5. API 设计

### 5.1 认证

```
POST /api/auth/login       # 登录，返回 JWT token
POST /api/auth/refresh     # 刷新 token
```

### 5.2 文章

```
GET    /api/posts           # 列表（分页、筛选）
GET    /api/posts/:id       # 详情
POST   /api/posts           # 创建（需认证）
PUT    /api/posts/:id       # 更新（需认证）
DELETE /api/posts/:id       # 删除（需认证）
GET    /api/posts/export    # 导出 CSV（需认证）
```

查询参数：`?page=1&per_page=5&tag=Go&status=published&search=关键词`

### 5.3 标签

```
GET    /api/tags            # 列表（含文章数）
POST   /api/tags            # 创建（需认证）
PUT    /api/tags/:id        # 更新（需认证）
DELETE /api/tags/:id        # 删除（需认证）
```

### 5.4 评论

```
GET    /api/comments        # 列表（分页、筛选）（需认证）
POST   /api/comments        # 提交评论（匿名）
PATCH  /api/comments/:id/approve  # 批准（需认证）
PATCH  /api/comments/:id/spam     # 标记垃圾（需认证）
DELETE /api/comments/:id          # 删除（需认证）
```

### 5.5 组件

```
GET    /api/components      # 列表
POST   /api/components      # 上传自定义组件（需认证）
PUT    /api/components/:id  # 更新（需认证）
DELETE /api/components/:id  # 删除（需认证）
PATCH  /api/components/:id/toggle  # 启用/禁用（需认证）
```

### 5.6 设置

```
GET    /api/settings        # 获取所有设置
PUT    /api/settings        # 保存设置（需认证）
POST   /api/settings/reset  # 重置默认（需认证）
```

### 5.7 仪表盘

```
GET    /api/dashboard/stats # 统计数据（需认证）
```

## 6. 前端路由

```javascript
// 博客前台
{ path: '/', component: 'blog/Home.vue' }
{ path: '/post/:id', component: 'blog/Post.vue' }
{ path: '/archives', component: 'blog/Archives.vue' }
{ path: '/tags', component: 'blog/Tags.vue' }
{ path: '/about', component: 'blog/About.vue' }

// 后台管理
{ path: '/admin/login', component: 'admin/Login.vue' }
{ path: '/admin', component: 'admin/Layout.vue', children: [
    { path: '', component: 'admin/Dashboard.vue' },
    { path: 'posts', component: 'admin/Posts.vue' },
    { path: 'tags', component: 'admin/Tags.vue' },
    { path: 'comments', component: 'admin/Comments.vue' },
    { path: 'custom', component: 'admin/Custom.vue' },
    { path: 'settings', component: 'admin/Settings.vue' },
]}
```

## 7. 核心功能设计

### 7.1 主题系统

- CSS 变量定义在 `:root` 和 `[data-theme="dark"]`
- Pinia store 管理主题状态
- `localStorage` 持久化，首次加载尊重 `prefers-color-scheme`
- 所有组件通过 CSS 变量自动适配

### 7.2 JWT 认证

- 登录返回 `access_token`（2小时过期）和 `refresh_token`（7天过期）
- Axios 拦截器自动附加 token
- 401 时自动刷新 token，刷新失败跳转登录页
- 后台路由守卫检查登录状态

### 7.3 自定义组件系统

**架构：**
- 自定义组件代码存储在数据库 `components` 表
- 前端通过 iframe 沙箱渲染自定义组件
- 组件间通过 `postMessage` 通信

**组件上传流程：**
1. 管理员在后台输入组件名、描述、版本、Vue SFC 代码
2. 代码保存到数据库
3. 渲染时创建 iframe，注入 Vue 运行时 + 组件代码
4. iframe 通过 postMessage 与主页面通信

**内置组件（可配置，不可删除）：**
- `<PostList />` — 文章列表（卡片/列表/时间线三种模式）
- `<Sidebar />` — 侧边栏
- `<Header />` — 页头
- `<MarkdownRenderer />` — Markdown 渲染器
- `<CodeBlock />` — 代码块
- `<CommentBox />` — 评论框

### 7.4 评论系统

- 匿名评论：读者填写昵称、邮箱（可选）、评论内容
- 评论默认状态为 `pending`，管理员审核后显示
- 支持标记为垃圾评论
- 评论展示在文章详情页底部

### 7.5 Markdown 渲染

- 后端存储 Markdown 原文
- 前端使用 `markdown-it` 渲染为 HTML
- `highlight.js` 代码高亮
- 支持：标题、列表、代码块、引用、表格、图片

### 7.6 设置系统

设置以键值对存储，分为 4 组：

**通用设置：**
- `site_title` — 站点标题
- `subtitle` — 副标题
- `description` — 站点描述
- `language` — 语言（zh-CN/en/bilingual）
- `posts_per_page` — 每页文章数

**作者信息：**
- `author_name` — 作者名
- `author_email` — 邮箱
- `author_bio` — 简介
- `author_github` — GitHub 链接

**功能开关：**
- `feature_comments` — 评论系统
- `feature_rss` — RSS 输出
- `feature_view_counter` — 阅读统计

**关于页面：**
- `about_content` — JSON 格式存储技术栈、联系方式等

## 8. UI 设计规范

### 8.1 设计令牌（CSS 变量）

```css
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

### 8.2 博客前台布局

- 导航栏：56px 高，固定顶部，毛玻璃背景
- 内容区：最大宽度 1080px，居中
- 页脚：版权信息 + 链接
- 响应式：640px 断点

### 8.3 后台管理布局

- 侧边栏：220px 宽，3 个分组（概览/内容/系统）
- 顶栏：56px 高，面包屑 + 操作按钮
- 内容区：弹性填充
- 响应式：900px 隐藏侧边栏

## 9. 部署方案

### 9.1 二进制部署

```bash
# 构建
cd web && npm run build        # 构建前端
cd ../server && go build -o vblog ./cmd/main.go  # 构建后端

# 运行
./vblog --config config.yaml
```

### 9.2 Docker Compose

```yaml
# docker-compose.yml
services:
  postgres:
    image: postgres:15
    environment:
      POSTGRES_DB: vblog
      POSTGRES_USER: vblog_admin
      POSTGRES_PASSWORD: ${DB_PASSWORD}
    volumes:
      - pgdata:/var/lib/postgresql/data
      - ./deploy/init.sql:/docker-entrypoint-initdb.d/init.sql

  vblog:
    build: .
    ports:
      - "8080:8080"
    depends_on:
      - postgres
    environment:
      DB_HOST: postgres
      DB_PORT: 5432
      DB_NAME: vblog
      DB_USER: vblog_admin
      DB_PASSWORD: ${DB_PASSWORD}
      JWT_SECRET: ${JWT_SECRET}

volumes:
  pgdata:
```

### 9.3 配置文件

```yaml
# config.yaml
server:
  port: 8080
  mode: release

database:
  host: localhost
  port: 5432
  name: vblog
  user: vblog_admin
  password: ""

jwt:
  secret: ""
  expire_hours: 2
  refresh_expire_hours: 168
```

## 10. 开发流程

### 10.1 开发命令

```bash
# 后端
cd server
go run ./cmd/main.go          # 启动后端（热重载用 air）
go test ./...                  # 运行测试

# 前端
cd web
npm install                    # 安装依赖
npm run dev                    # 开发服务器（Vite）
npm run build                  # 构建生产版本
npm run lint                   # ESLint 检查
```

### 10.2 数据库初始化

```bash
# 使用 init.sql
psql -U vblog_admin -d vblog -f deploy/init.sql

# 或使用 GORM AutoMigrate（开发环境）
```

## 11. 实现阶段

### Phase 1：基础框架
1. 项目结构搭建（monorepo）
2. Go 后端 + go-restful 基础路由
3. GORM 模型 + PostgreSQL 连接
4. Vue 3 + Element Plus 前端脚手架
5. 设计令牌 + 主题系统

### Phase 2：核心功能
1. 文章 CRUD（后台管理 + API）
2. 标签 CRUD
3. 博客首页（文章列表、分页、标签筛选）
4. 文章详情页（Markdown 渲染）

### Phase 3：扩展功能
1. JWT 认证系统
2. 评论系统
3. 归档页、标签页、关于页
4. 仪表盘统计

### Phase 4：高级功能
1. 自定义组件系统（iframe 沙箱）
2. 设置管理
3. 部署脚本（Docker + 二进制）
4. README 和文档
