# vBlog Core

一个面向极客和 Vibe Coder 的可自定义轻量博客系统。使用 Markdown 写作，支持组件自定义、标签分类、评论系统、RSS 订阅等功能。

## 功能特性

### 博客前台

- **首页**：文章列表、统计概览、Ctrl+F 搜索、分页浏览（支持页码跳转）
- **文章详情**：Markdown 渲染、代码高亮、目录导航（TOC）、阅读量统计、上/下篇导航
- **归档页面**：按时间线浏览所有文章
- **标签页面**：按标签分类浏览，显示每个标签对应的文章数量
- **关于页面**：博主信息展示
- **评论系统**：访客评论，管理员可开关
- **RSS 订阅**：自动生成 RSS 2.0 Feed
- **主题切换**：亮色/暗色主题，自动持久化

### 后台管理

- **仪表盘**：文章总数、总阅读量、评论数、标签数统计
- **文章管理**：创建、编辑、删除文章；Markdown 编辑器（md-editor-v3）支持图片上传；Markdown 文件导入
- **标签管理**：标签 CRUD，显示文章计数
- **评论管理**：查看、删除评论
- **组件自定义**：iframe 沙盒自定义组件，支持启用/禁用切换
- **系统设置**：站点信息、作者信息、功能开关
- **图片上传**：编辑器内粘贴/拖拽图片自动上传

## 技术栈

| 层级 | 技术 |
|------|------|
| 前端 | Vue 3 + Element Plus + Pinia + Vue Router 4 + md-editor-v3 |
| 后端 | Go + go-restful/v3 + GORM |
| 数据库 | PostgreSQL |
| 认证 | JWT（golang-jwt/v5） |
| 配置 | Viper + TOML |
| 前端构建 | Vite |

## 项目结构

```
vBlog Core/
├── server/                  # Go 后端
│   ├── cmd/
│   │   ├── main.go          # 服务入口
│   │   └── seed/            # 测试数据种子脚本
│   ├── api/                 # REST API 处理器
│   │   ├── post.go          # 文章 API
│   │   ├── tag.go           # 标签 API
│   │   ├── comment.go       # 评论 API
│   │   ├── setting.go       # 设置 API
│   │   ├── auth.go          # 登录/注册 API
│   │   ├── component.go     # 组件 API
│   │   ├── dashboard.go     # 仪表盘统计 API
│   │   ├── rss.go           # RSS Feed
│   │   └── upload.go        # 图片上传 API
│   ├── service/             # 业务逻辑层
│   ├── model/               # 数据模型（GORM）
│   ├── middleware/           # JWT 中间件
│   └── config/              # 配置（Viper + TOML）
│       ├── config.go        # 配置加载
│       ├── config.toml      # 配置文件
│       └── database.go      # 数据库连接
├── web/                     # Vue 3 前端
│   └── src/
│       ├── blog/            # 博客前台页面
│       ├── admin/           # 后台管理页面
│       ├── shared/          # 共享组件（导航、页脚、评论）
│       ├── stores/          # Pinia 状态管理
│       ├── api/             # Axios 请求封装
│       ├── styles/          # 全局样式与设计变量
│       └── utils/           # 工具函数
├── hdx/                     # HTML 原型（UI 参考）
├── docs/                    # 文档
└── deploy/                  # 部署配置
```

## 快速开始

### 环境要求

- Go 1.22+
- Node.js 18+
- PostgreSQL 14+

### 1. 配置数据库

创建 PostgreSQL 数据库：

```sql
CREATE DATABASE vblog;
CREATE USER vblog WITH PASSWORD 'your_password';
GRANT ALL PRIVILEGES ON DATABASE vblog TO vblog;
```

### 2. 修改配置文件

编辑 `server/config/config.toml`：

```toml
[http]
addr = "0.0.0.0"
port = 8080

[postgres]
host = "127.0.0.1"
port = 5432
name = "vblog"
user = "vblog"
password = "your_password"

[jwt]
secret = "your-jwt-secret"
```

### 3. 启动后端

```bash
cd server
go run ./cmd/main.go
```

后端启动后自动创建数据库表结构，监听 `0.0.0.0:8080`。

### 4. 启动前端

```bash
cd web
npm install
npm run dev
```

前端开发服务器运行在 `http://localhost:5173`，API 请求自动代理到 `http://localhost:8080`。

### 5. 构建生产版本

```bash
cd web
npm run build
```

构建产物输出到 `web/dist/`。将 `dist/` 目录内容复制到 `server/static/`，Go 服务器会自动提供静态文件服务并处理 SPA 路由。

### 6. 注册管理员

首次使用需注册管理员账号：

1. 访问 `/admin/register` 注册账号
2. 登录后即可进入后台管理

### 7. 填充测试数据（可选）

```bash
cd server
go run ./cmd/seed/
```

自动创建 50 篇测试文章和 20 个标签。

## 运行测试

```bash
cd server
go test ./... -v
```

项目采用 TDD 开发模式，每个模块的测试文件与实现在同一目录下。

## API 接口

### 公开接口

| 方法 | 路径 | 说明 |
|------|------|------|
| GET | `/api/posts` | 文章列表（支持分页、搜索、标签筛选） |
| GET | `/api/posts/{id}` | 文章详情（自动递增阅读量） |
| GET | `/api/tags` | 标签列表（含文章计数） |
| GET | `/api/comments?post_id={id}` | 文章评论 |
| POST | `/api/comments` | 提交评论 |
| GET | `/api/settings` | 站点设置 |
| GET | `/api/dashboard/stats` | 统计数据 |
| POST | `/api/auth/login` | 用户登录 |
| POST | `/api/auth/register` | 用户注册 |
| GET | `/api/rss` | RSS Feed |

### 管理接口（需 JWT）

| 方法 | 路径 | 说明 |
|------|------|------|
| POST | `/api/posts` | 创建文章 |
| PUT | `/api/posts/{id}` | 更新文章 |
| DELETE | `/api/posts/{id}` | 删除文章 |
| POST/PUT/DELETE | `/api/tags` | 标签 CRUD |
| PUT | `/api/settings` | 保存设置 |
| POST | `/api/upload` | 上传图片 |
| CRUD | `/api/components` | 组件管理 |

## 设计系统

项目使用 CSS 变量实现主题系统，定义在 `web/src/styles/variables.css`：

```css
/* 亮色主题 */
--bg: #fafafa; --surface: #ffffff; --fg: #171717; --muted: #737373;
--border: #e5e5e5; --accent: #2563eb; --accent-soft: #eff6ff;

/* 暗色主题 */
--bg: #0a0a0a; --surface: #141414; --fg: #ededed; --muted: #a3a3a3;
--border: #262626; --accent: #3b82f6; --accent-soft: #172554;
```

主题切换通过 `data-theme="dark"` 属性实现，自动持久化到 localStorage。

## 许可证

MIT
