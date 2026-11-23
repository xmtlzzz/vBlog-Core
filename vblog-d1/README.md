# vBlog D1

vBlog Core 的「纯 Cloudflare 版」：**同一套 Vue 3 前端 + JS Worker + D1(SQLite) + R2(新图片)/KV(历史图片)**，无需 Go 后端、Neon 或 Fly.io，全部运行在 Cloudflare。

- 技术栈：Vue 3 + 原生 JS Worker（Markdown 使用 markdown-it，JWT/PBKDF2 使用 WebCrypto）+ D1 + R2；Workers KV 兼容历史图片
- 已部署上线，域名见下方「在线预览」

## 这个版本是什么

`vblog-d1` 是 vBlog Core 的**纯 Cloudflare 便捷部署版**：保留同一套 Vue 3 前端与 API 契约，把 Go 后端替换为 JS Worker，数据库使用 D1(SQLite)、新图片存储到 R2，历史 KV 图片仍可读取。**不需要公网服务器、不需要 Go/Node 常驻进程**，`deploy.ps1` 一条命令即可构建并上线。

部署细节见「部署」与「在线预览」；相对原版 Go 项目做了什么改动，见「与 vBlog Core 的差异」「性能优化记录」和下方「最近更新」。

## 在线预览（已部署在 Cloudflare）

| 项 | 地址 |
|---|---|
| 主域名 | https://vblog.xmtlz.dev |
| 备用域名（workers.dev） | https://vblog-d1.xmtlzloveasuka.workers.dev（部分网络不可达） |
| 后台入口 | https://vblog.xmtlz.dev/admin |
| 数据库 | D1 `vblog-d1-db`（APAC 区域） |
| 图片存储 | R2 `vblog-images`（新上传）；Workers KV 保留历史图片 |

> ⚠️ **预览版不支持 gRPC**：Cloudflare Worker 仅提供 HTTP REST，没有 TCP 端口——依赖 gRPC 双向流与变更推送的 **Wails 桌面客户端无法连接此预览版**；Web 博客/后台功能不受影响。

## 架构

```
浏览器 ──► Worker(vblog-d1)
           ├─ /api/*        → JS 路由（复刻 Go 后端 REST 契约，同源免 CORS）
           ├─ 其余路径      → Static Assets 静态资源（web/dist，SPA 兜底）
           ├─ D1 (SQLite)   → 文章/标签/评论/设置/组件/用户/统计（库位于 APAC 区域）
           ├─ R2            → 新图片上传；本站 /uploads/ 提供图片访问
           ├─ KV            → 读取历史图片；无 R2 绑定时兼容上传
           └─ Cron(每日)    → 访问统计快照（00:05 UTC）
```

## 与 vBlog Core 的差异

| 项 | vBlog Core（原版） | vBlog D1（本克隆） |
|---|---|---|
| 后端 | Go（REST + gRPC） | JS Worker（仅 REST），零依赖 |
| 数据库 | PostgreSQL（Neon） | D1（SQLite，APAC 区域） |
| 部署 | CF Pages + Fly.io + Neon + R2 | CF 一家搞定 |
| 前端 | Vue 3（`web/`） | **同一份 `web/`，零改动** |
| 图片上传 | 本地磁盘 | R2 优先；公开域名可选，保留历史 KV 图片 |
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

# 4) R2 图床：新账号部署时先建桶，并更新 account_id 和 bucket_name
#    node node_modules\wrangler\bin\wrangler.js r2 bucket create vblog-images --location apac
#    - R2_PUBLIC_URL 留空即可：桶保持私有，通过本站 /uploads/ 访问图片
#    - 可选绑定独立公开域名后填写 R2_PUBLIC_URL；历史 KV 图片仍可访问
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

## Markdown 与本地图片导入

后台「文章 → 导入文章文件夹」选择同时包含 `.md` 和图片的文件夹；多篇 Markdown 时先选择一篇。
图片按 Markdown 所在目录解析相对路径（含 `../`、Windows 反斜杠、中文及 URL 编码文件名），仅上传正文引用到的图片。
缺图时补选图片或图片文件夹；绝对路径、存在多个匹配结果及无法确定的匹配必须人工确认，不能直接继续。
确认后点击「上传图片并进入编辑」，图片压缩上传后自动替换正文路径，再预览、保存草稿或发布。
上传中断可在当前窗口重试，已经成功上传的文件不会再次上传。

「Markdown 上传」仍支持单个文件：纯文字或网络图片直接进入编辑；本地图片则进入补选窗口。补选文档目录、`assets` 或 `images` 文件夹后，会按相对路径自动选择唯一匹配的图片，无需逐张选择。
支持 PNG/JPEG/GIF/WebP，上传文件最大 10MB；普通位图复用编辑器的 WebP 压缩，GIF 保留动图。
网络图片及本站 `/uploads/` 地址保留，代码块、行内代码与普通链接不会被当成图片上传。

## 备份与恢复

```powershell
# 一键备份远程库到 backup/vblog-<时间戳>.sql
pwsh ./backup.ps1

# 恢复（整库覆盖导入）
# 注意：对非空库直接 import 会撞已存在的表。恢复到已有库前先删表，
# 或干脆删库重建（wrangler d1 create + 重新绑定 id）后再导入。
node node_modules\wrangler\bin\wrangler.js d1 import vblog-d1-db --remote --file .\backup\vblog-xxx.sql

# 图片备份：backup/images-kv-<时间戳>/ 与 images-r2-<时间戳>/ 各含 manifest.json。
# manifest 中 objects[].key 是原始 key，objects[].file 是对应二进制文件名。
# 恢复 KV：wrangler kv key put <key> --binding IMG --remote --path <file>
# 恢复 R2：wrangler r2 object put vblog-images/<key> --remote --file <file> --content-type <原类型>
# 旧版 images-<时间戳>/ 备份仍按文件名中的 __ 还原为 /，再写回 KV。

# 时点恢复（Cloudflare 侧自动快照）
node node_modules\wrangler\bin\wrangler.js d1 time-travel vblog-d1-db --remote
```

> 教训记录：删库前必须先备份/导出。D1 删除后无法用 time-travel 恢复（控制台如有「已删除数据库」入口可尝试，30 天内）。

## 最近更新

- **2026-09 体验补全（后台可配置化 + 移动端）**
  - 首页标语 `hero_title`：打字机大标题后台可配（支持换行），设置就绪后才开始打字，并做 HTML 转义
  - 关于页技术栈 `about_tech`：每行「名称|角色|图标（emoji 或 vue）」，整项留空使用内置默认卡片
  - QR 徽章配色预设 `qr_palette`：樱花（默认）/森林/海洋/落日（仅树模型）；二维码放大视图左侧显示「手机扫码访问 GitHub」提示
  - 浏览器标签页标题跟随后台「站点标题」（App.vue 统一设置，index.html 兜底 `vBlog`）
  - 移动端导航（≤640px）：两行布局——品牌+操作一行、链接一行横向滑动，消除 375px 下折行与页面横向滚动
  - 页脚版权年份动态化（不再写死）

- **2026-09 页脚 every-qrcode 动态二维码徽章（后台可配置）**
  - 新组件 `web/src/shared/EveryQrBadge.vue`：基于 `@every-qrcode/core` + `renderer-webgpu` 官方底层包（其 Web Component 不暴露 scene 样式参数），作者 GitHub 链接确定性生成 3D 樱花树/地形
  - 交互：点击徽章放大至约 2.6 倍并变形为可扫描二维码（内容即 GitHub 主页，已实测解码），再点恢复小徽章；WebGPU 不可用自动降级静态 SVG 二维码
  - 后台「设置 → 页脚二维码徽章」：启用开关 / 模型（tree|terrain）/ 氛围特效（calm|snow|rain|wind）/ 配色预设（樱花|森林|海洋|落日）/ 尺寸（28–120px，默认 40）/ 背景色——settings 为自由 KV，后端零改动
  - 树模型视图 `setZoom(1.2)` 更饱满且不裁树冠（>1.25 会裁顶）；二维码视图保持 1.0 不裁静区
  - 项目二维码导航：`qr_links` 每行「名称|URL」、只填 URL（名称自动取域名）或「名称|URL|配色」手动指定；「模块」页顶部生成卡片墙（`QrLinksSection.vue`），配色未指定时按行自动轮换预设（樱花/森林/海洋/落雪，相邻卡片不同色），每个地址确定性生成专属 3D 徽章，点击原地变形为可扫描二维码（payload 已实测解码正确）；有项目徽章时模块页不再叠加「暂无功能模块」空状态
  - 模块页二维码卡片带氛围特效切换按钮（默认/平静/刮风/下雨/落雪，仅本地视图实时切换，不写回后台配置）
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
│   ├── package.json         # Wrangler 与 Markdown 摘要依赖
│   └── .dev.vars            # 本地开发密钥（gitignore）
├── deploy.ps1               # 一键构建 + 部署
├── backup.ps1               # 一键备份远程 D1
└── README.md
```

## 本地回归检查

两套前端分别在 `web/` 和 `../web/` 运行 `npm test`、`npm run build`。Worker 在 `worker/` 运行 `npm test`。
前端构建后，在 `worker/` 运行 `npm run test:storage:runtime`，验证真实 workerd 中的 R2 上传、图片读取、HEAD/ETag、历史 KV 兼容及健康检查；测试资源独立于生产绑定。

运行完整 D1 回归时，在 `worker/` 依次执行下面命令；测试配置只使用本地数据库，不连接生产绑定。第二个终端执行最后一行，检查公开文章、登录后草稿访问、保留正文与标签的草稿转换，以及服务端分享元信息。

```powershell
npx wrangler d1 execute DB --local --config tests/wrangler.toml --persist-to .wrangler/ui-tests --file migrations/0001_init.sql
npx wrangler d1 execute DB --local --config tests/wrangler.toml --persist-to .wrangler/ui-tests --file tests/fixtures.sql
npx wrangler dev --local --config tests/wrangler.toml --persist-to .wrangler/ui-tests --port 8787
# 另一个终端（仍在 worker/）
npm run test:runtime
```

预览使用 `tests/fixtures.sql` 中的演示文章。线上测试文章可在后台「全部文章」转为草稿，并可通过编辑页重新发布；程序不按文件名或标题自动过滤文章。
