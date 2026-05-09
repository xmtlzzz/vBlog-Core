# 集成测试实现计划

**Goal:** 给后端 service 层补充数据库级别的集成测试，覆盖完整 CRUD 流程。

**Tech Stack:** Go testing + GORM + PostgreSQL (192.168.81.101)

**TDD 规则：**
- 先写测试 → 运行确认失败 → 写实现（如有） → 运行确认通过 → 提交
- 测试文件与实现文件同目录

**数据库连接：** 从 `.env` 读取，使用 `config.Load()` 获取配置

---

## Task 1: 测试基础设施

**Files:**
- Create: `server/testutil/db.go`

创建测试数据库连接工具函数：
- `SetupTestDB()` — 连接数据库，返回 *gorm.DB
- `CleanupDB(db, models...)` — 清理测试数据

---

## Task 2: PostService 集成测试

**Files:**
- Create: `server/service/post_integration_test.go`

测试用例：
- TestPostService_Create — 创建文章，验证字段
- TestPostService_List — 创建多篇，验证分页和筛选
- TestPostService_GetByID — 获取单篇，验证 Tags 预加载
- TestPostService_Update — 更新文章，验证 read_time 重算
- TestPostService_Delete — 软删除，验证不存在

---

## Task 3: TagService 集成测试

**Files:**
- Create: `server/service/tag_integration_test.go`

测试用例：
- TestTagService_Create — 创建标签
- TestTagService_List — 列表含 post_count
- TestTagService_Update — 更新标签
- TestTagService_Delete — 删除标签

---

## Task 4: CommentService 集成测试

**Files:**
- Create: `server/service/comment_integration_test.go`

测试用例：
- TestCommentService_Create — 创建评论，默认 pending
- TestCommentService_List — 按状态筛选
- TestCommentService_Approve — 批准评论
- TestCommentService_MarkSpam — 标记垃圾
- TestCommentService_Delete — 删除评论

---

## Task 5: ComponentService + SettingService 集成测试

**Files:**
- Create: `server/service/component_integration_test.go`
- Create: `server/service/setting_integration_test.go`

---

## Task 6: AuthService 集成测试

**Files:**
- Create: `server/service/auth_integration_test.go`

测试用例：
- TestAuthService_Register — 注册用户，验证密码已哈希
- TestAuthService_Login — 正确密码登录
- TestAuthService_LoginWrongPassword — 错误密码失败
- TestAuthService_LoginNonexistent — 不存在用户失败
