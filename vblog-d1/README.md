# vBlog D1

vBlog Core 的「纯 Cloudflare 版」克隆：**同一个 Vue 3 前端 + 零依赖 JS Worker + D1(SQLite) + KV/R2(图片)**，不要 Go 后端、不要 Neon、不要 Fly.io。全部运行在 Cloudflare 免费额度内。

- 技术栈：Vue 3（前端零改动）+ 原生 JS Worker（零 npm 依赖，JWT/PBKDF2 用 WebCrypto）+ D1 + Workers KV（图片，R2 可选）
- 已部署上线，域名见下方「在线预览」

## 这个版本是什么

`vblog-d1` 是 vBlog Core 的**纯 Cloudflare 便捷部署版**：保留同一套 Vue 3 前端与 API 契约，把 Go 后端替换为零依赖 JS Worker，数据库用 D1(SQLite)、图片用 KV/R2，全部运行在 Cloudflare 免费额度内——**不需要公网服务器、不需要 Go/Node 常驻进程**，`deploy.ps1` 一条命令即可构建并上线。

部署细节见「部署」与「在线预览」；相对原版 Go 项目做了什么改动，见「与 vBlog Core 的差异」「性能优化记录」和下方「最近更新」。

## 在线预览（已部署在 Cloudflare）

| 项 | 地址 |
|---|---|
| 主域名 | https://vblog.xmtlz.dev |
| 备用域名（workers.dev） | https://vblog-d1.xmtlzloveasuka.workers.dev（部分网络不可达） |
| 后台入口 | https://vblog.xmtlz.dev/admin |
| 注册管理员（首次） | https://vblog.xmtlz.dev/admin/register |
| 数据库 | D1 `vblog-d1-db`（APAC 区域） |
| 图片存储 | Workers KV（未启用 R2；启用后自动切换） |

> ⚠️ **预览版不支持 gRPC**：Cloudflare Worker 仅提供 HTTP REST，没有 TCP 端口——依赖 gRPC 双向流与变更推送的 **Wails 桌面客户端无法连接此预览版**；Web 博客/后台功能不受影响。

## 架构

```
浏览器 ──► Worker(vblog-d1)
           ├─ /api/*        → JS 路由（复刻 Go 后端 REST 契约，同源免 CORS）
           ├─ 其余路径      → Static Assets 静态资源（web/dist，SPA 兜底）
           ├─ D1 (SQLite)   → 文章/标签/评论/设置/组件/用户/统计（库位于 APAC 区域）
           ├─ KV            → 图片上传（免费 1GB；配了 R2 自动切 R2）
           └─ Cron(每日)    → 访问统计快照（00:05 UTC）
```

## 与 vBlog Core 的差异

| 项 | vBlog Core（原版） | vBlog D1（本克隆） |
|---|---|---|
| 后端 | Go（REST + gRPC） | JS Worker（仅 REST），零依赖 |
| 数据库 | PostgreSQL（Neon） | D1（SQLite，APAC 区域） |
| 部署 | CF Pages + Fly.io + Neon + R2 | CF 一家搞定 |
| 前端 | Vue 3（`web/`） | **同一份 `web/`，零改动** |
| 图片上传 | 本地磁盘/R2 | KV 图床（默认）→ R2（配了 R2_PUBLIC_URL 自动切） |
| 桌面客户端 gRPC | ✅ | ❌（Worker 无 TCP） |
| 认证 | bcrypt + golang-jwt | PBKDF2(WebCrypto) + HS256(WebCrypto) |

> API 响应契约（JSON 字段、状态码、错误文案）与 Go 端一致。
> 数据库 schema 由 `server/model/*.go` 的 GORM 模型移植而来，见 `worker/migrations/0001_init.sql`。

## 部署（一条命令）

```powershell
# 前置：Node 20+；Cloudflare 账号；安装 wrangler 依赖（一次即可）
cd worker
npm install          # 把 wrangler 装到 worker/node_modules

# 登录 Cloudflare（浏览器授权，只需一次；凭据存到 .wrangler/）
cd ..                # 回到 vblog-d1/
$env:WRANGLER_HOME = "$PWD\.wrangler"
cd worker
node node_modules\wrangler\bin\wrangler.js login

# 一键构建 + 部署
cd ..
pwsh ./deploy.ps1
```

### 首次部署的手动步骤（对应的一次性资源）

```powershell
# 1) 建 D1 数据库（区域选 APAC 可显著降低国内访问延迟）并把 database_id 填入 wrangler.toml
node node_modules\wrangler\bin\wrangler.js d1 create vblog-d1-db --location apac
node node_modules\wrangler\bin\wrangler.js d1 migrations apply vblog-d1-db --remote

# 2) 建 KV 命名空间（图片回退图床）并把 id 填入 wrangler.toml
node node_modules\wrangler\bin\wrangler.js kv namespace create IMG

# 3) 注入 JWT 密钥（至少 32 位随机串）
node node_modules\wrangler\bin\wrangler.js secret put JWT_SECRET

# 3b) 注入 Turnstile 人机验证 secret（Cloudflare 控制台 → Turnstile → 你的 widget → Secret Key；勿贴进聊天）
node node_modules\wrangler\bin\wrangler.js secret put TURNSTILE_SECRET

# 4) 可选：R2 图床（需先在 Cloudflare 控制台 → R2 首次访问启用服务）
#    - 控制台建桶 vblog-images，绑定自定义域（如 uploads.blog.xmtlz.dev）
#    - wrangler.toml 取消 [[r2_buckets]] 注释，把公开地址填进 R2_PUBLIC_URL
#    - 之后上传自动走 R2，历史 KV 图片不受影响
```

### 绑定自定义域名

1. Cloudflare 控制台 → **Workers 和 Pages** → `vblog-d1` → **设置 → 域和路由 → 添加** `xxx.xxx.xxx`
2. 域名所在 zone 必须在**同一个 Cloudflare 账号**下（NS 指向 cloudflare）
3. 证书自动签发需要几分钟；期间 `521/522` 或打不开属正常，等 5~10 分钟后刷新

> 为什么强烈建议自定义域名：`*.workers.dev` 在国内网络经常无法访问（DNS 解析正常但 TCP 被拦），自定义域名走 Cloudflare 代理通常可用。

## 本地开发

```powershell
# 终端 1：前端
cd web
npm run dev          # http://localhost:5173（vite 代理 /api → localhost:8787）

# 终端 2：Worker（本地模拟器，不用联网）
cd worker
node node_modules\wrangler\bin\wrangler.js d1 migrations apply vblog-d1-db --local
node node_modules\wrangler\bin\wrangler.js dev --port 8787
# 本地 JWT 密钥写在 worker/.dev.vars（已 gitignore；远程密钥用 wrangler secret put）
```

## 备份与恢复

```powershell
# 一键备份远程库到 backup/vblog-<时间戳>.sql
pwsh ./backup.ps1

# 恢复（整库覆盖导入）
# 注意：对非空库直接 import 会撞已存在的表。恢复到已有库前先删表，
# 或干脆删库重建（wrangler d1 create + 重新绑定 id）后再导入。
node node_modules\wrangler\bin\wrangler.js d1 import vblog-d1-db --remote --file .\backup\vblog-xxx.sql

# 图床图片恢复：backup/images-<时间戳>/ 里的文件名 __ 即原 key 的 /，
# 用 wrangler kv key put 逐个写回 IMG 绑定（key 形如 uploads/<纳秒>.<ext>）

# 时点恢复（Cloudflare 侧自动快照）
node node_modules\wrangler\bin\wrangler.js d1 time-travel vblog-d1-db --remote
```

> 教训记录：删库前必须先备份/导出。D1 删除后无法用 time-travel 恢复（控制台如有「已删除数据库」入口可尝试，30 天内）。

## 最近更新

- **2026-09 页脚 every-qrcode 动态二维码徽章（后台可配置）**
  - 新组件 `web/src/shared/EveryQrBadge.vue`：基于 `@every-qrcode/core` + `renderer-webgpu` 官方底层包（其 Web Component 不暴露 scene 样式参数），作者 GitHub 链接确定性生成 3D 樱花树/地形
  - 交互：点击徽章放大至 104px 并变形为可扫描二维码（内容即 GitHub 主页，已实测解码），再点恢复小徽章；WebGPU 不可用自动降级静态 SVG 二维码
  - 后台「设置 → 页脚二维码徽章」：启用开关 / 模型（tree|terrain）/ 氛围特效（calm|snow|rain|wind）/ 背景色——settings 为自由 KV，后端零改动
  - 性能：徽章含 3D 渲染器，懒加载——页脚进入视口才动态加载（gzip ≈124KB 独立 chunk，首屏不受影响）
  - 顺带：首页 Hero 副标题不再硬编码，读后台「站点描述」；`web/` 与 `vblog-d1/web` 镜像同步修改

- **2026-08 Cloudflare Turnstile 人机验证（评论）**
  - 前端 `web/src/shared/CommentSection.vue` 显式渲染 Turnstile（action=`comment`，主题随站点，提交后 reset 支持重试）
  - 后端 Worker `verifyTurnstile()`：`siteverify` 校验 `success / action / hostname`，missing 或伪造 token 一律 `403`（在线已实测拦截）
  - 配置：`TURNSTILE_SECRET`（secret）+ `TURNSTILE_HOSTNAMES=vblog.xmtlz.dev`（`[vars]`，生产不含 localhost；本地开发在 `.dev.vars` 用 `localhost,127.0.0.1`）

- **2026-08-22 熊猫吉祥物 & 性能与部署收尾**
  - 品牌图：浏览器图标用 A1、首页顶栏用 A2、关于页头像用 B1（`web/public/favicon.png` / `nav-mascot.png` / `avatar.png`，AI 生成资源不入 git）
  - 在线预览域名 `vblog.xmtlz.dev` 绑定并验证可用（workers.dev 备用，注意部分网络不可达）
  - 性能优化：D1 重建至 APAC 区域（国内访问延迟大幅下降）、文章标签单查询消除 N+1、公开只读接口边缘缓存、page_view 降载
  - 图片存储默认回退 Workers KV（未启用 R2 时可用）；每日 00:05(UTC) cron 统计快照
  - 新增 `deploy.ps1`（一键构建 + 部署）与 `backup.ps1`（远程 D1 导出备份，破坏性操作前必跑）

## 性能优化记录（已内置）

| 优化 | 说明 |
|---|---|
| D1 区域 = APAC | 库建在亚太，避免美国往返 ~200ms 延迟；`wrangler d1 create --location apac` |
| 标签查询零 N+1 | 一条 `json_group_array` 子查询带出文章标签（列表/详情/回收站共用） |
| 边缘缓存 | `/api/posts` 30s、settings/tags/stats/components 60s、RSS 300s（`cf.cacheTtl`）；详情页不缓存以保持阅读量准确 |
| page_view 降载 | 仅真实页面浏览写一行（API 请求不再写），减轻 D1 写入配额压力 |

首页加载路径：`posts`(2 次 D1 查询) + 3 个缓存命中接口 → 实测本地 64ms。

## 免费额度速查（2025 现行）

| 资源 | 免费额度 | 说明 |
|---|---|---|
| Worker 请求 | 10 万次/天 | 个人博客绰绰有余 |
| D1 | 5GB 存储；每天 500 万次读、10 万次写 | 写配额注意别被刷评论打爆 |
| KV | 1GB 存储、每天 1000 次写、10 万次读 | 无 R2 时图片回退到这里；读有边缘缓存 |
| R2 | 10GB 存储、出站不收费 | 需先在控制台启用 R2 服务再建桶 |
| 静态资源 | 计入 Worker 请求 | SPA 由边缘分发 |

## 已知取舍与排障

- **桌面客户端（Wails + gRPC）不可用**：Web 前后台功能不受影响；后台「gRPC 状态」页无接口，属预期。
- **workers.dev 打不开**：国内网络常见，换自定义域名。
- **自定义域名 521/522**：证书签发中，等 5~10 分钟；zone 必须在同一 CF 账号。
- **上传报「upload storage not configured」**：wrangler.toml 里 KV id 没填（或两台都没配置）。
- **改了代码部署后没生效**：`node node_modules\wrangler\bin\wrangler.js tail`（或控制台日志）看实时日志；确认 `deploy.ps1` 输出了新 Version ID。
- **中文乱码排查**：Worker 全程 UTF-8；PowerShell 5.1 控制台显示 `?` 是编码假象，浏览器正常。

## 目录结构

```
vblog-d1/
├── web/                     # Vue 3 前端（与 vBlog Core 相同，未改动）
├── worker/
│   ├── src/index.js         # Worker 入口：路由 + 全部 API + cron
│   ├── src/util.js          # 零依赖工具：JWT / PBKDF2 / 日期
│   ├── migrations/0001_init.sql   # D1 schema
│   ├── wrangler.toml        # Worker/Assets/D1/KV/R2 配置
│   ├── package.json         # 本地 wrangler 依赖（devDependency）
│   └── .dev.vars            # 本地开发密钥（gitignore）
├── deploy.ps1               # 一键构建 + 部署
├── backup.ps1               # 一键备份远程 D1
└── README.md
```