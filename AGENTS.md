# CampusLogistics

智慧校园后勤管理系统 — Go + Gin + GORM。

## 技术栈

- **Go 1.26.1**, **Gin** Web 框架, **GORM** ORM, **SQLite** (开发) / **MySQL** (生产)
- **JWT** 认证（access + refresh token）, **Snowflake** 字符串 ID, `go-playground/validator`, `zap` 日志, `viper` 配置
- 密码: **bcrypt** (`golang.org/x/crypto/bcrypt`)
- 模块路径: `github.com/Struggle-Rabbit/CampusLogistics`

## 架构

```
cmd/server/          # 入口（初始化顺序: Snowflake → config → logger → validator → DB → router.Run）
internal/app/        # App 结构体, 持有 config + *gorm.DB（依赖注入根）
internal/config/     # Viper 加载 configs/config-{ENV}.yaml（ENV 环境变量, 默认 "dev"）
internal/router/     # 11 个路由文件, 将 building 路由定义在 campus.go 中（非独立文件）
internal/controller/ # 11 个领域包（含 common, 用于 register/login，无对应 service）
internal/service/    # 10 个领域包 + ServiceProvider 注入层（无 common service）
internal/model/      # 13 个 GORM 模型（Campus、Building、DormRoom 定义在同一文件）
internal/dao/        # DB 初始化, Redis 空实现, 分页辅助, GORM 查询辅助
api/dto/             # 11 个请求/响应 DTO 文件
pkg/                 # 日志、工具（Snowflake、JWT、校验器、响应辅助）、常量、错误
deploy/              # Dockerfile + docker-compose.yml 均为空模板; k8s 为空占位
scripts/sql/dev.db   # SQLite 开发数据库（启动时自动迁移）
docs/                # Swagger 生成文档（swaggo/swag）
test/unit_test/      # 11 个测试文件, 包名 unittest, 使用 glebarez/sqlite 内存库
test/api_test/       # 空占位（尚无集成测试）
```

## 命令

```bash
# 启动开发服务器（SQLite, 启动时自动迁移）
go run ./cmd/server/

# 运行全部单元测试
go test ./test/unit_test/ -v

# 运行单个测试
go test ./test/unit_test/ -v -Run TestUserService_Register

# 编译二进制（输出到当前目录）
go build -o server ./cmd/server/

# 重新生成 Swagger 文档（需安装 swag CLI）
swag init -g cmd/server/main.go -o docs/
```

## 关键约定

- **环境切换**: `ENV` 环境变量控制加载 `configs/config-dev.yaml`（默认）或 `configs/config-prod.yaml`。也可通过 `APP_` 前缀环境变量覆盖（`.` 转 `_`，如 `APP_JWT_SECRET`）。
- **数据库**: 开发模式使用 SQLite + 启动时自动迁移 13 个模型；生产模式使用 MySQL（手动迁移）。
- **ID 生成**: Snowflake 字符串 ID，在 GORM `BeforeCreate` hook 中自动生成（`utils.GenStringID()`），不要手动设置。
- **JWT 认证**: access token + refresh token。refresh token 直接存储在 `sys_user.refresh_token` 字段（未使用 Redis）。
- **中间件执行顺序**: Recovery → RequestID → Logger → CORS（全局）→ OperationLog（全局 `/api/v1`）→ JWTAuth + PermissionValidator（按路由）。
- **OperationLog**: 自动脱敏 `password`/`token` 等敏感字段，异步写入 `sys_operation_log` 表。
- **权限校验**: `PermissionValidator(perms)` 通过 `sys_role_menu` + `sys_user_role` 多对多关联查询。
- **`common` 控制器**: 没有对应的 service 包，register/login 直接调用 `UserService`。
- **Service 依赖**: `UserService` 显式依赖 `MenuService`（在 `NewServiceProvider` 中注入）。
- **自定义校验器**: `utils.InitValidator()` 注册了 `mobile` tag（`^1[3-9]\d{9}$`）。
- **分页**: `dao.Paginate(page, size)`，最大限制 100 条。
- **CORS**: 动态允许来源（取请求 `Origin` 头）。
- **空占位包**: `pkg/errx/` 仅为 `.gitkeep`；部署文件均为空模板。
- **无 Makefile / 无 CI/CD** — 所有操作通过原生 `go` 命令完成。

## 测试

- 测试位于 `test/unit_test/`，**包名 `unittest`**。
- `SetupTestDB()` 创建临时 SQLite 内存库，自动迁移全部 13 个模型，返回 `(*gorm.DB, *app.App)`。
- 测试需手动设置 `config.GlobalConfig`（含 JWT Secret）并调用 `utils.InitSnowflake()`（已封装在 `SetupTestDB()` 中）。
- 使用 `stretchr/testify` 做断言。
- 手动构造 service 实例进行依赖注入（如 `user.NewUserService(appInstance, menuSvc)`）。

## 模块功能说明

### 路由层 (`internal/router/`)
11 个路由文件，定义 API 端点与权限绑定：

| 路由文件 | 职责 | 公共/需认证 |
|----------|------|------------|
| `router.go` | 根路由，初始化 Gin 引擎+中间件栈+ServiceProvider | - |
| `common.go` | `/api/v1/register`（注册）, `/login`（登录）, `/notice/public`（公共公告列表） | 公共 |
| `user.go` | 用户 CRUD、密码重置、获取权限 | 需 JWT |
| `system.go` | 操作日志列表查询 | 需 JWT |
| `role.go` | 角色 CRUD | 需 JWT |
| `menu.go` | 菜单 CRUD | 需 JWT |
| `repair.go` | 报修单提交/列表/详情/更新/流转/删除 | 需 JWT |
| `campus.go` | 校区 CRUD + 楼栋 CRUD + 楼栋导入/导出 | 需 JWT |
| `dorm.go` | 宿舍 CRUD + 分配/调宿/迁出/容量预警 | 需 JWT |
| `utility.go` | 水电费 CRUD + 缴费/批量缴费/单价设置/统计/欠费预警 | 需 JWT |
| `notice.go` | 公告 CRUD + 置顶 | 需 JWT |

### 控制器层 (`internal/controller/`)
11 个包，处理 HTTP 请求/响应（数据校验 + 调用 service）：

| 控制器 | HTTP 端点前缀 | 核心接口 | 对应服务 |
|--------|-------------|---------|---------|
| `common` | (无前缀) | POST `/login`, `/register` | UserService |
| `user` | `/api/v1/user` | 列表/详情/新增/编辑/删除/重置密码/权限查询 | UserService |
| `system` | (无前缀) | POST `/OperationLogList` | SystemService |
| `role` | `/api/v1/role` | 列表/详情/新增/编辑/删除 | RoleService |
| `menu` | `/api/v1/menu` | 列表/详情/新增/编辑/删除 | MenuService |
| `repair` | `/api/v1/repair` | 提交/列表/详情/更新/流转/删除 | RepairService |
| `campus` | `/api/v1/campus` | 新增/编辑/删除/列表/详情/所有校区 | CampusService |
| `building` | `/api/v1/building` | 新增/编辑/删除/列表/详情/按校区/导入/导出 | BuildingService |
| `dorm` | `/api/v1/dorm` | 新增/编辑/删除/列表/详情/分配/调宿/迁出/成员/预警 | DormService |
| `utility` | `/api/v1/utility` | 新增/编辑/删除/列表/详情/缴费/批量缴费/单价/统计/预警/个人查询 | UtilityService |
| `notice` | `/api/v1/notice` | 新增/编辑/删除/列表/详情/置顶 | NoticeService |

### 服务层 (`internal/service/`)
10 个领域包 + 1 个注入器，封装核心业务逻辑：

| 服务 | 核心职责 | 依赖 |
|------|---------|------|
| `UserService` | 注册（bcrypt 加密）、登录（JWT 签发）、用户 CRUD、密码重置、菜单权限查询（通过 MenuService） | MenuService |
| `RoleService` | 角色 CRUD、角色-菜单权限关联 | - |
| `MenuService` | 菜单 CRUD（树形结构）、菜单-角色关联 | - |
| `SystemService` | 操作日志查询 | - |
| `RepairService` | 报修单提交/状态流转（6 种状态）、报修记录追踪 | - |
| `CampusService` | 校区 CRUD | - |
| `BuildingService` | 楼栋 CRUD、按校区查询、批量导入/导出（Excel） | - |
| `DormService` | 宿舍 CRUD、学生入住/调宿/迁出、容量预警 | - |
| `UtilityService` | 水电费录入（自动算费）、缴费/批量缴费、单价配置、月度统计、欠费预警 | - |
| `NoticeService` | 公告 CRUD、置顶控制（最多 3 条）、公开列表查询 | - |

`ServiceProvider` 是依赖注入根，在 `router.go` 中统一创建所有 service 并注入 controller。

### 数据模型层 (`internal/model/`)
13 个 GORM 模型 + 2 个基础嵌入模型：

| 模型 | 表名 | 业务含义 |
|------|------|---------|
| `BaseModel` | (嵌入) | 字符串 ID（Snowflake）+ 创建/更新时间 |
| `BaseModelWithDelete` | (嵌入) | 同上 + 软删除字段 |
| `BaseModelIntId` | (嵌入) | 自增 uint ID |
| `Campus` | `campus` | 校区信息（名称、地址、联系方式） |
| `Building` | `building` | 楼栋（关联校区、楼层数、房间数） |
| `DormRoom` | `dorm_room` | 宿舍（关联楼栋、楼层、类型、最大人数、当前人数） |
| `DormUser` | `dorm_user` | 宿舍人员关联（入住/迁出时间、状态） |
| `DormUtility` | `dorm_utility` | 水电费记录（年份/月份、用量、金额、缴费状态） |
| `UtilityPrice` | `utility_price` | 水电费单价配置（水价/电价） |
| `SysUser` | `sys_user` | 用户（学号/工号、姓名、手机号、密码 bcrypt、角色多对多） |
| `SysRole` | `sys_role` | 角色（名称、编码、状态、内置标记、用户/菜单多对多） |
| `SysMenu` | `sys_menu` | 菜单/权限树（父级ID、路由、组件、权限标识、类型1目录2菜单3按钮） |
| `SysOperationLog` | `operation_log` | 操作日志（用户、请求方法/路径/参数、IP、浏览器、状态码） |
| `Notice` | `notice` | 公告（标题、内容、类型、置顶、发布时间、浏览量、附件JSON） |
| `RepairOrder` | `repair_order` | 报修单（单号、类型、地点、描述、图片JSON、6 状态流转） |
| `RepairRecord` | `repair_record` | 报修流转记录（操作人、旧/新状态、备注） |

### 中间件层 (`internal/middleware/`)
6 个 Gin 中间件，按顺序注册：

| 中间件 | 注册位置 | 功能 |
|--------|---------|------|
| `Recovery` | 全局 | Panic 恢复，防止进程崩溃 |
| `RequestID` | 全局 | 为每个请求注入唯一追踪 ID |
| `Logger` | 全局 | 请求日志记录（zap） |
| `CORS` | 全局 | 动态允许来源（取请求 Origin 头） |
| `OperationLog` | `/api/v1` 组 | 自动脱敏 password/token 等敏感字段，异步写入 operation_log 表 |
| `JWTAuth` | `/api/v1` 子组 | Bearer Token 解析、用户信息注入 gin.Context |
| `PermissionValidator(perms)` | 按路由 | 通过 sys_role_menu + sys_user_role 多对多校验权限标识 |

### 数据访问层 (`internal/dao/`)
| 文件 | 职责 |
|------|------|
| `sql.go` | DB 初始化（SQLite 开发/MySQL 生产）、自动迁移 |
| `redis.go` | Redis 空实现（占位，后续迭代） |
| `page.go` | `Paginate(page, size)` 分页辅助（最大 100 条） |
| `grom_func.go` | GORM 查询辅助函数 |

### DTO 层 (`api/dto/`)
11 个 DTO 文件，定义请求体和响应体的 Go 结构体，与 `go-playground/validator` 配合做输入校验。

### 共享包 (`pkg/`)
| 包 | 文件 | 职责 |
|----|------|------|
| `utils` | `jwt.go` | JWT 生成/解析（access + refresh token） |
| `utils` | `validator.go` | 自定义校验器注册（如 `mobile` tag: `^1[3-9]\d{9}$`） |
| `utils` | `id_generator.go` | Snowflake 字符串 ID 生成器 |
| `utils` | `response.go` | 统一 JSON 响应格式（Success/Fail/Unauth/NoPermission） |
| `utils` | `utils.go` | 通用工具函数 |
| `constant` | `status_code.go` | HTTP 状态码常量 |
| `constant` | `res_msg.go` | 响应消息常量 |
| `constant` | `user_const.go` | 用户类型/状态常量 |
| `constant` | `menu_const.go` | 菜单类型常量 |
| `logger` | `logger.go` | zap 日志初始化/全局实例 |
| `errx` | (空占位) | 预留自定义错误包 |

### 配置层 (`internal/config/`)
Viper 加载 `configs/config-{ENV}.yaml`，支持 `APP_` 前缀环境变量覆盖（`.` 转 `_`，如 `APP_JWT_SECRET`）。

### 应用启动 (`internal/app/`)
`App` 结构体持有 `Config` + `*gorm.DB`，作为依赖注入的根对象。

### 部署 (`deploy/`)
- `docker/Dockerfile` + `docker-compose.yml`：空模板，未配置
- `k8s/`：.gitkeep 占位，未配置

### 测试 (`test/unit_test/`)
11 个测试文件（包名 `unittest`），使用 `glebarez/sqlite` 内存库 + `stretchr/testify` 断言。

## WHERE TO LOOK

| 任务 | 位置 | 备注 |
|------|------|------|
| 新增 API 端点 | `internal/router/` 加路由 → `internal/controller/` 加处理器 → `internal/service/` 加业务 → `api/dto/` 加 DTO | 遵循现有分层 |
| 新增数据模型 | `internal/model/` 加结构体 → 在 `dao/sql.go` 的 `AutoMigrate` 列表中注册 | 继承 BaseModel/BaseModelWithDelete |
| 修改权限系统 | `internal/middleware/permission.go`（校验逻辑）+ `internal/model/sys_menu.go`（权限标识）+ `internal/model/sys_role.go`（角色-菜单关联） | - |
| 修改认证逻辑 | `internal/middleware/jwt.go` + `pkg/utils/jwt.go` + `internal/service/user/user_service.go` | - |
| 调整数据库连接 | `internal/dao/sql.go` | 切换 SQLite/MySQL |
| 修改 API 响应格式 | `pkg/utils/response.go` | - |
| 添加单元测试 | `test/unit_test/` | 使用 `SetupTestDB()` 创建内存库 |
| Swagger 文档 | `docs/`（生成）+ 控制器中的 swagger 注解 | 运行 `swag init` 重新生成 |
| 配置项 | `configs/config-dev.yaml` / `configs/config-prod.yaml` | 开发/生产分离 |
| 常量定义 | `pkg/constant/` | - |

## ANTI-PATTERNS (本项目)

- **勿试图直接运行 `go test ./test/unit_test/` 前不检查包名**：测试包名为 `unittest`（非 `unit_test`），编译时会先进入 `dao/` 等依赖包的 init/全局状态初始化流程
- **勿在 service 中直接操作 `gin.Context`**：目前 `UserService.GetUserInfo()` 直接从 `c.Get("userID")` 读取，应改为入参透传 userID
- **勿在非 `common` 控制器外直接访问 `UserService`**：controller→service 调用关系应保持单向，controller 只调用自己依赖的 service
- **勿手动设置 ID**：所有模型通过 `BeforeCreate` hook + `utils.GenStringID()` 自动生成，DB 层不应干预
- **勿在生产环境使用 SQLite + 自动迁移**：开发环境使用 SQLite，生产环境使用 MySQL 并手动管理迁移
- **勿依赖 `pkg/errx/` 中的内容**：该包目前仅为 .gitkeep 占位，无实际实现
- **勿跳过 `OperationLogMiddleware` 的敏感字段脱敏**：所有新增接口需确保 password/token 等敏感参数在日志中被过滤
