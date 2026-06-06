# CampusLogistics 项目重构计划

## Context

当前项目存在以下问题：
1. **目录结构混乱**：`internal/user/` 和 `internal/base/` 是重复代码，与 `internal/model/`、`internal/service/user/`、`internal/controller/user/` 功能重叠
2. **滥用 interface**：每个 Service 都定义了对应的 interface（如 `UserService interface`、`MenuService interface`），但每个 interface 只有一个实现，违反了 Go "接口应由使用者定义" 的原则
3. **Controller 命名不规范**：Go/Gin 社区习惯用 `handler` 而非 `controller`
4. **Router 文件过多**：router 拆成了 10 个文件，每个只声明一个 module 的路由，增加维护成本
5. **Provider 命名臃肿**：`XxxServiceProvider`、`XxxControllerProvider` 等命名无实际意义

---

## 重构后目录结构

```
internal/
├── config/           (保留，config.go)
├── model/            (保留，所有数据库 model)
├── dao/              (保留，数据库连接/分页等)
├── handler/          (新建，替换 controller，单个 package，10 个文件)
│   ├── user_handler.go
│   ├── menu_handler.go
│   ├── role_handler.go
│   ├── system_handler.go
│   ├── campus_handler.go
│   ├── building_handler.go
│   ├── dorm_handler.go
│   ├── notice_handler.go
│   ├── repair_handler.go
│   └── utility_handler.go
├── service/          (重构，去掉 interface，单个 package，10 个文件)
│   ├── user_service.go
│   ├── menu_service.go
│   ├── role_service.go
│   ├── system_service.go
│   ├── campus_service.go
│   ├── building_service.go
│   ├── dorm_service.go
│   ├── notice_service.go
│   ├── repair_service.go
│   └── utility_service.go
├── router/
│   └── router.go          (合并所有路由到一个文件)
├── middleware/       (保留不变)
```

根目录保留：`api/dto/`、`pkg/`、`cmd/`、`configs/`、`docs/`、`test/` 不变。

---

## 代码审查发现的问题（重构时一并修复）

### 1. Service 层问题

**service/user/user_service.go**
- `UserService` struct 使用 `DB *gorm.DB` 直接传参，与其他 service 用 `App *app.App` 风格不一致 → 统一改为 `App` 模式
- `NewUserService` 构造函数多余 → 删除，统一在 router 中直接实例化
- `GetUserInfo` 中 `gorm.ErrRecordNotFound` 判断后可直接 return nil, err，无需再判 → 简化
- `UpdateUser` 手机号唯一性检查逻辑有问题：查询到非 `ErrRecordNotFound` 的错误时返回 err，但查询到记录（手机号被占用）时返回业务错误，这个逻辑是对的，只是 `model.SysUser{}` 的 `Updates` 会更新零值字段（如空字符串的 Name）→ 改为 `map[string]interface{}` 只更新非空字段

**service/menu/menu_service.go**
- `BuildMenuTree` 方法暴露为公开方法，但实际上只被 router 调用和测试使用 → 改为私有
- `GetMenuListByPage` 中 `total` 按 `parent_id = '0'` 计数，但列表查询没有这个限制 → 计数与列表不一致

**service/role/role_service.go**
- `CreateRole` 中解码到 `role` 变量但不用，后面直接新建了 `SysRole{}` → 删除无用变量

**service/system/system_service.go**
- `GetOperationLogListByPage` 中时间范围查询逻辑：`!req.OperationTimeStart.IsZero() && req.OperationTimeEnd.IsZero()` 条件反了，应该是两个都非零才做范围查询 → 修复

**service/notice/notice_service.go**
- `GetDetail` 中 `view_count` 自增与返回存在竞态（先 read 再 update 再返回 +1 的值）→ 用 `gorm.Expr("view_count + 1")` 直接 SQL 自增，返回时查询实际值

**service/utility/utility_service.go**
- `getPayStatusName` 方法可以提取到 `pkg/constant/` 中复用
- `GetUserDormUtility` 中 `First(&dormUser)` 没有处理 `ErrRecordNotFound` → 修复

### 2. Controller（Handler）层问题

**controller/user/user_controller.go**
- `QueryDetail` 从 context 取 user_id 查自己信息，但路由上没有参数可以查他人信息 → 符合当前需求，保留
- `DelUser` 中 `data["id"].([]string)` 类型断言失败会直接报"请选择要删除的数据"，但前端传的可能是一个 string → 改进类型解析

**controller/menu/menu_controller.go**
- `CreateMenu`、`UpdateMenu`、`GetListByPage` 都用了 `_ = utils.ShouldBind(...)` 忽略返回值 → 应该检查绑定结果

**controller/role/role_controller.go**
- `QueryDetail` 用 `ShouldBindJSON` 解析 query 参数不合理（这是 GET 请求）→ 改为 `c.GetQuery("id")`

### 3. Router 层问题

**router/router.go**
- `initRouter` 中 service 全部在函数内创建，但每个 service 依赖 `*app.App` → 去掉 app，直接用 config + db
- 路由文件分散，每个文件一个 `LoadXxxRouter` 函数 → 合并

### 4. Model 层问题

**model/base.go**
- `BaseModelIntId` 未被任何 model 使用 → 删除

**model/sys_user.go**
- `UserID` 在 `CreatorID`、`RepairOrder.UserID` 中用 `string` 存"学号/工号"，但 JWT 里存的是 `sysUser.ID`（雪花 ID）→ 这是一个设计不一致，本次重构备注但不改动（需产品确认）

### 5. go.mod 问题
- 依赖写了两段 `require (...)` → 合并为一段

---

## 具体操作步骤

### Step 1：删除冗余目录
- 删除 `internal/user/`（handler.go/service.go/model.go/route.go 全是重复代码）
- 删除 `internal/base/`（model.go 与 `model/base.go` 完全重复）
- 删除 `internal/app/`（仅包装一个 struct，可直接内联到 router）
- 删除 `internal/controller/`（整体替换为 handler）

### Step 2：Service 层重构（去 interface + 修 bug）
每个 service 文件执行：
- **删除** `XxxService interface{}` 定义
- **重命名** struct：`XxxServiceProvider` → `svc`（小写私有，包内唯一实现无需暴露类型名）
  - 实际 struct 命名为小写（如 `type svc struct` 有语法冲突，所以用 `type service struct`）
  - 对外导出的构造函数 `NewXxxService()` 保留，返回 `*service`（小写类型）
- **统一** user_service.go 风格：去掉 `DB *gorm.DB` + `Menu menu.MenuService`，改为 `App *app.App` + `MenuSvc *menu.Service`
- 修复上述代码审查发现的所有 bug

### Step 3：Controller → Handler（去 interface + 修 bug）
- 创建 `internal/handler/` 目录（单个 package）
- 每个 handler 文件：
  - 删除 `XxxController interface{}`
  - struct 命名为小写：`type handler struct{ svc *service }`
  - 导出构造函数 `NewXxxHandler(svc *service) *handler`
- 修复上述代码审查发现的绑定/类型断言问题

### Step 4：Router 合并（单文件）
- 删除所有 router 子文件
- 在 `router.go` 中 `initRouter` 函数内：创建 service → 创建 handler → 注册路由
- 去掉 service 接口参数传递，直接使用 struct 指针

### Step 5：更新 main.go
- 去掉 `internal/app` import
- `Run()` 函数参数从 `*app.App` 改为 `(cfg *config.Config, db *gorm.DB)`

### Step 6：更新测试文件
- 所有 `XxxServiceProvider` → 对应的新 `NewXxxService()` 返回值
- 所有 interface 类型断言 → struct 指针类型
- 测试逻辑根据重构后的 service 字段调整

### Step 7：清理与合并 go.mod
- 合并两段 `require` 为一段
- 运行 `go mod tidy`

### Step 8：验证
1. `go build ./...` 编译通过
2. `go test ./...` 测试通过
3. 确认 Swagger 文档可用
