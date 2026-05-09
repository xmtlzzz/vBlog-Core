# vBlog Core

可自定义组件的极客博客系统。

## 技术栈

- **前端**: Vue 3 + Element Plus + Pinia
- **后端**: Go + go-restful + GORM
- **数据库**: PostgreSQL
- **认证**: JWT

## 快速开始

### Docker Compose（推荐）

```bash
cd deploy
docker compose up -d
```

访问 http://localhost:8080

### 本地开发

**后端:**
```bash
cd server
go run ./cmd/main.go
```

**前端:**
```bash
cd web
npm install
npm run dev
```

前端开发服务器运行在 http://localhost:5173，API 请求代理到 http://localhost:8080。

## 环境变量

在项目根目录创建 `.env` 文件：

```env
DB_HOST=192.168.81.101
DB_PORT=5432
DB_NAME=vblog
DB_USER=vblog
DB_PASSWORD=your_password
JWT_SECRET=your_secret
SERVER_PORT=8080
```

## 项目结构

```
vBlog Core/
├── server/          # Go 后端
│   ├── cmd/         # 入口
│   ├── api/         # REST API handlers
│   ├── service/     # 业务逻辑
│   ├── model/       # 数据模型
│   ├── middleware/   # JWT 中间件
│   └── config/      # 配置
├── web/             # Vue 3 前端
│   └── src/
│       ├── blog/    # 博客前台页面
│       ├── admin/   # 后台管理页面
│       ├── shared/  # 共享组件
│       ├── stores/  # Pinia 状态
│       └── api/     # API 请求封装
├── deploy/          # 部署配置
├── docs/            # 文档
└── hdx/             # HTML 原型（参考）
```

## 运行测试

```bash
cd server && go test ./... -v
```
