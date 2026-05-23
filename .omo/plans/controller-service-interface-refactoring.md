# 控制层和服务层接口化重构实施计划

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** 将 CampusLogistics 的 controller 层和 service 层公开方法改为接口暴露，去掉 ServiceProvider 中间载体和所有 New 构造函数

**Architecture:** 每个 service 包定义 `XxxService` 接口 + `XxxServiceProvider` 实现；每个 controller 包定义 `XxxController` 接口 + `XxxControllerProvider` 实现；router 直接创建 service 实例并分发；测试改用结构体字面量初始化

**Tech Stack:** Go 1.26, Gin, GORM

---

## TL;DR

> **Quick Summary**: 对 10 个 service 包、11 个 controller 包进行接口化重构，删除 `internal/service/provider.go`，去掉所有 `NewXxx()` 构造函数，调用方直接使用 `&XxxProvider{Field: value}` 创建实例
>
> **Deliverables**:
> - 10 个 service 包新增接口定义
> - 11 个 controller 包新增接口定义
> - 删除 `internal/service/provider.go`
> - 更新 12 个 router 文件签名
> - 更新 11 个测试文件
>
> **Estimated Effort**: Medium
> **Parallel Execution**: YES - 4 waves
> **Critical Path**: Service 层 → Controller 层 → Router → 测试

---

## Context

### 设计决策摘要

| 维度 | 选择 |
|------|------|
| 接口命名 | 领域命名（`BuildingService`, `UserController`） |
| 实现命名 | Provider 后缀（`BuildingServiceProvider`, `UserControllerProvider`） |
| 构造函数 | 无（直接 `&XxxProvider{Field: val}`） |
| 中间载体 | 删除 ServiceProvider |
| Controller 注入 | 直接依赖各自需要的 Service 接口 |

### 文件改动统计

- **删除**: 1 个（`internal/service/provider.go`）
- **Service 层修改**: 11 个（10 个主要文件 + `operation_log.go`）
- **Controller 层修改**: 11 个
- **Router 层修改**: 12 个（`router.go` + 11 个路由文件）
- **测试层修改**: 11 个
- **总计**: ~46 个文件

---

## Work Objectives

### Core Objective
将 service 和 controller 的公开方法通过 Go 接口暴露，消除对具体类型的直接依赖。

### Must Have
- 每个 service 包定义接口，实现改为 `XxxServiceProvider`，字段导出
- 每个 controller 包定义接口，实现改为 `XxxControllerProvider`，直接依赖 Service 接口
- 删除 `internal/service/provider.go`
- 所有 router 函数签名改为接收 Service 接口参数
- `go build ./...` 编译通过
- `go test ./test/unit_test/ -v` 全部测试通过

### Must NOT Have
- 不修改任何业务逻辑
- 不修改 API 路由路径
- 不修改 HTTP 请求/响应格式
- 不修改 `internal/app/`、`internal/config/`、`internal/model/`、`internal/dao/`、`internal/middleware/`、`api/dto/`、`pkg/`

---

## Verification Strategy

- **编译验证**: `go build ./...` → 无错误
- **测试验证**: `go test ./test/unit_test/ -v` → 全部 PASS
- **LSP 诊断**: 无类型错误

---

## Execution Strategy

### Parallel Execution Waves

```
Wave 1 (Service 层 — 10 个任务可并行):
├── Task 1: building service 接口化
├── Task 2: campus service 接口化
├── Task 3: dorm service 接口化
├── Task 4: menu service 接口化
├── Task 5: notice service 接口化
├── Task 6: repair service 接口化
├── Task 7: role service 接口化
├── Task 8: system service 接口化 (2 个文件)
├── Task 9: utility service 接口化
└── Task 10: user service 接口化 (依赖 menu 接口)

Wave 2 (Controller 层 — 11 个任务，依赖 Wave 1 完成后可并行):
├── Task 11: building controller 接口化
├── Task 12: campus controller 接口化
├── Task 13: common controller 接口化
├── Task 14: dorm controller 接口化
├── Task 15: menu controller 接口化
├── Task 16: notice controller 接口化
├── Task 17: repair controller 接口化
├── Task 18: role controller 接口化
├── Task 19: system controller 接口化
├── Task 20: user controller 接口化
└── Task 21: utility controller 接口化

Wave 3 (Router 层 — 12 个任务，依赖 Wave 2 完成后可并行):
├── Task 22: 删除 provider.go
├── Task 23: 更新 router.go
├── Task 24: 更新 common.go 路由
├── Task 25: 更新 campus.go 路由
├── Task 26: 更新 user.go 路由
├── Task 27: 更新 system.go 路由
├── Task 28: 更新 role.go 路由
├── Task 29: 更新 menu.go 路由
├── Task 30: 更新 repair.go 路由
├── Task 31: 更新 dorm.go 路由
├── Task 32: 更新 utility.go 路由
└── Task 33: 更新 notice.go 路由

Wave 4 (测试层 — 11 个任务，依赖 Wave 2 完成后可并行):
├── Task 34: 更新 user_service_test.go
├── Task 35: 更新 campus_building_service_test.go
├── Task 36: 更新 dorm_service_test.go
├── Task 37: 更新 menu_service_test.go
├── Task 38: 更新 notice_service_test.go
├── Task 39: 更新 repair_controller_test.go
├── Task 40: 更新 repair_service_test.go
├── Task 41: 更新 role_service_test.go
├── Task 42: 更新 system_service_test.go
├── Task 43: 更新 utility_service_test.go

Wave FINAL (验证):
├── Task F1: go build ./... 编译验证
├── Task F2: go test ./test/unit_test/ -v 测试验证
└── Task F3: 代码审查 — 确认接口定义完整、无遗留的 ServiceProvider 引用

Critical Path: Wave 1 → Wave 2 → Wave 3 + Wave 4 → Wave FINAL
```

---

## TODOs

### Wave 1: Service 层接口化

---

- [x] 1. **building service — 加接口、结构体改名为 Provider**

  **Files:**
  - Modify: `internal/service/building/building_service.go`

  **Changes:**
  - 在文件顶部添加 `BuildingService` 接口，包含所有公开方法签名
  - 将 `type BuildingService struct` 重命名为 `type BuildingServiceProvider struct`
  - 将 `app *app.App` 改为导出字段 `App *app.App`
  - 将 `func NewBuildingService(app *app.App) *BuildingService` 移除
  - 将所有方法的接收者从 `*BuildingService` 改为 `*BuildingServiceProvider`

  在文件顶部 `package building` 后添加：

  ```go
  type BuildingService interface {
      Create(req *dto.BuildingCreateReq) error
      Update(req *dto.BuildingUpdateReq) error
      Delete(ids []string) error
      GetListByPage(req *dto.BuildingListPageReq) (*dto.PageResult, error)
      GetDetail(id string) (*dto.BuildingResult, error)
      GetBuildingsByCampus(campusID string) ([]*dto.BuildingResult, error)
      ImportBuildings(file *multipart.FileHeader) (int, error)
      ExportBuildings(req *dto.BuildingExportReq) (string, error)
  }
  ```

  将结构体定义从：
  ```go
  type BuildingService struct {
      app *app.App
  }
  ```
  改为：
  ```go
  type BuildingServiceProvider struct {
      App *app.App
  }
  ```

  删除 `NewBuildingService` 函数：
  ```go
  // 删除整个函数
  func NewBuildingService(app *app.App) *BuildingService {
      return &BuildingService{app: app}
  }
  ```

  将每个方法的接收者从 `(s *BuildingService)` 改为 `(s *BuildingServiceProvider)`，并将 `s.app` 改为 `s.App`。

  **Verified by:** `go build ./internal/service/building/...` 编译通过

- [x] 2. **campus service — 加接口、结构体改名为 Provider**

  **Files:**
  - Modify: `internal/service/campus/campus_service.go`

  **Changes:** 与 Task 1 相同模式。
  - 接口名: `CampusService`
  - 实现名: `CampusServiceProvider`
  - 字段: `App *app.App`
  - 移除 `NewCampusService`
  - 接收者改为 `*CampusServiceProvider`，`s.app` → `s.App`

  **接口定义：**
  ```go
  type CampusService interface {
      Create(req *dto.CampusCreateReq) error
      Update(req *dto.CampusUpdateReq) error
      Delete(ids []string) error
      GetListByPage(req *dto.CampusListPageReq) (*dto.PageResult, error)
      GetDetail(id string) (*dto.CampusResult, error)
      GetAll() ([]*dto.CampusResult, error)
  }
  ```

- [x] 3. **dorm service — 加接口、结构体改名为 Provider**

  **Files:**
  - Modify: `internal/service/dorm/dorm_service.go`

  **Changes:** 与 Task 1 相同模式。
  - 接口名: `DormService`
  - 实现名: `DormServiceProvider`
  - 字段: `App *app.App`
  - 移除 `NewDormService`
  - 接收者改为 `*DormServiceProvider`，`s.app` → `s.App`

  **接口定义：**
  ```go
  type DormService interface {
      Create(req *dto.DormCreateReq) error
      Update(req *dto.DormUpdateReq) error
      Delete(ids []string) error
      GetListByPage(req *dto.DormListPageReq) (*dto.PageResult, error)
      GetDetail(id string) (*dto.DormResult, error)
      AssignDorm(req *dto.DormAssignReq) error
      TransferDorm(req *dto.DormTransferReq) error
      CheckOut(req *dto.DormCheckOutReq) error
      GetDormUsers(req *dto.DormUserListReq) (*dto.PageResult, error)
      GetCapacityWarning() ([]*dto.DormResult, error)
  }
  ```

- [x] 4. **menu service — 加接口、结构体改名为 Provider**

  **Files:**
  - Modify: `internal/service/menu/menu_service.go`

  **Changes:** 与 Task 1 相同模式。
  - 接口名: `MenuService`
  - 实现名: `MenuServiceProvider`
  - 字段: `App *app.App`
  - 移除 `NewMenuService`
  - 接收者改为 `*MenuServiceProvider`，`s.app` → `s.App`

  **接口定义：**
  ```go
  type MenuService interface {
      CreateMenu(req *dto.CreateMenuReq) error
      UpdateMenu(req *dto.UpdateMenuReq) error
      DelMenu(id []string) error
      GetMenuList(req *dto.MenuListReq) ([]dto.MenuResult, error)
      GetMenuListByPage(req *dto.MenuListByPageReq) (*dto.PageResult, error)
      MenuDetailById(id string) (*dto.MenuResult, error)
      BuildMenuTree(allMenus []model.SysMenu) []dto.MenuResult
  }
  ```

---

- [x] 5. **notice service — 加接口、结构体改名为 Provider**

  **Files:**
  - Modify: `internal/service/notice/notice_service.go`

  **Changes:** 与 Task 1 相同模式。
  - 接口名: `NoticeService`
  - 实现名: `NoticeServiceProvider`
  - 字段: `App *app.App`
  - 移除 `NewNoticeService`
  - 接收者改为 `*NoticeServiceProvider`，`s.app` → `s.App`

  **接口定义：**
  ```go
  type NoticeService interface {
      Create(creatorID string, req *dto.NoticeCreateReq) error
      Update(req *dto.NoticeUpdateReq) error
      Delete(ids []string) error
      GetListByPage(req *dto.NoticeListPageReq) (*dto.PageResult, error)
      GetPublicList(req *dto.NoticeListPageReq) (*dto.PageResult, error)
      GetDetail(id string) (*dto.NoticeResult, error)
      SetTop(req *dto.NoticeTopReq) error
  }
  ```

- [x] 6. **repair service — 加接口、结构体改名为 Provider**

  **Files:**
  - Modify: `internal/service/repair/repair_service.go`

  **Changes:** 与 Task 1 相同模式。
  - 接口名: `RepairService`
  - 实现名: `RepairServiceProvider`
  - 字段: `App *app.App`
  - 移除 `NewRepairService`
  - 接收者改为 `*RepairServiceProvider`，`s.app` → `s.App`
  - 注意：`GenerateOrderNo` 当前是值接收者 `(RepairService)`，改为指针接收者 `(*RepairServiceProvider)`

  **接口定义：**
  ```go
  type RepairService interface {
      GenerateOrderNo(prefix string) string
      RepairOrderSubmit(userID string, req *dto.RepairOrderSubmitReq) error
      GetListByPage(req *dto.RepairOrderListByPageReq) (*dto.PageResult, error)
      GetDetailById(id string) (*dto.RepairOrderResult, error)
      DelRepairOrderById(id string) error
      UpdateRepairOrder(req dto.UpdateRepairOrderSubmitReq) error
      OrderRecord(req dto.RecordReq) error
  }
  ```

- [x] 7. **role service — 加接口、结构体改名为 Provider**

  **Files:**
  - Modify: `internal/service/role/role_service.go`

  **Changes:** 与 Task 1 相同模式。
  - 接口名: `RoleService`
  - 实现名: `RoleServiceProvider`
  - 字段: `App *app.App`
  - 移除 `NewRoleService`
  - 接收者改为 `*RoleServiceProvider`，`s.app` → `s.App`

  **接口定义：**
  ```go
  type RoleService interface {
      CreateRole(req *dto.CreateRoleReq) error
      UpdateRole(req *dto.UpdateRoleReq) error
      DelRole(id []string) error
      GetRoleList(name string) ([]dto.RoleResult, error)
      GetRoleListByPage(req *dto.RoleListByPageReq) (*dto.PageResult, error)
      RoleDetailById(id string) (*dto.RoleResult, error)
  }
  ```

---

- [x] 8. **system service (2 个源文件) — 加接口、结构体改名为 Provider**

  **Files:**
  - Modify: `internal/service/system/system_service.go`
  - Modify: `internal/service/system/operation_log.go`

  **Changes:** 与 Task 1 相同模式。
  - 接口名: `SystemService`（包含两个文件中的所有公开方法）
  - 实现名: `SystemServiceProvider`
  - 字段: `App *app.App`
  - 移除 `NewSystemService`
  - 两个文件的接收者都改为 `*SystemServiceProvider`，`s.app` → `s.App`

  **接口定义（合并两个文件的方法）：**
  ```go
  type SystemService interface {
      RefreshToken(token string) (*dto.RefreshTokenResult, error)
      GetOperationLogListByPage(req *dto.OperationLogByPageReq) (*dto.PageResult, error)
  }
  ```

  `system_service.go` 中结构体改为：
  ```go
  type SystemServiceProvider struct {
      App *app.App
  }
  ```

  `operation_log.go` 中接收者从 `(s *SystemService)` 改为 `(s *SystemServiceProvider)`。

- [x] 9. **utility service — 加接口、结构体改名为 Provider**

  **Files:**
  - Modify: `internal/service/utility/utility_service.go`

  **Changes:** 与 Task 1 相同模式。
  - 接口名: `UtilityService`
  - 实现名: `UtilityServiceProvider`
  - 字段: `App *app.App`
  - 移除 `NewUtilityService`
  - 接收者改为 `*UtilityServiceProvider`，`s.app` → `s.App`

  **接口定义：**
  ```go
  type UtilityService interface {
      Create(req *dto.UtilityCreateReq) error
      Update(req *dto.UtilityUpdateReq) error
      Delete(ids []string) error
      GetListByPage(req *dto.UtilityListPageReq) (*dto.PageResult, error)
      GetDetail(id string) (*dto.UtilityResult, error)
      Pay(req *dto.UtilityPayReq) error
      BatchPay(req *dto.UtilityBatchPayReq) error
      ImportData(reqs []*dto.UtilityImportReq) (int, error)
      UpdatePrice(req *dto.UtilityPriceReq) error
      GetPrice() (*dto.UtilityPriceResult, error)
      GetStatistics(campusID string, year, month int) (*dto.UtilityStatResult, error)
      GetUnpaidWarning() ([]*dto.UtilityResult, error)
      GetUserDormUtility(userID string, year, month int) (*dto.PageResult, error)
  }
  ```

- [x] 10. **user service — 加接口、依赖改为 MenuService 接口**

  **Files:**
  - Modify: `internal/service/user/user_service.go`

  **Changes:**
  - 接口名: `UserService`
  - 实现名: `UserServiceProvider`
  - 字段 `App *app.App` → `App *app.App`（导出）
  - 字段 `menu *menu.MenuService` → `Menu menu.MenuService`（接口类型，导出）
  - 移除 `NewUserService`
  - 接收者改为 `*UserServiceProvider`，`s.app` → `s.App`，`s.menu` → `s.Menu`

  **接口定义：**
  ```go
  type UserService interface {
      Register(req *dto.RegisterReq) error
      Login(req *dto.LoginReq) (*dto.LoginResult, error)
      GetUserInfo(c *gin.Context) (*dto.UserInfoResult, error)
      GetListByPage(req *dto.UserListPageReq) (*dto.PageResult, error)
      UpdateUser(req *dto.UserUpdateReq) error
      DelUser(id []string) error
      ResetPassword(req *dto.PasswordReset) error
      GetUserPermission(user_id string) (*dto.UserPermissionResult, error)
  }
  ```

  注意：`menu` 字段类型从 `*menu.MenuService` 改为 `menu.MenuService`（接口），所有使用 `s.menu` 的地方改为 `s.Menu`。

### Wave 2: Controller 层接口化

---

- [x] 11. **building controller — 加接口、直接依赖 BuildingService 接口**

  **Files:**
  - Modify: `internal/controller/building/building_controller.go`

  **Changes:**
  - 移除 `import "github.com/.../internal/service"`（不再需要 ServiceProvider）
  - 添加接口 `BuildingController`
  - 结构体从 `BuildingController` 改为 `BuildingControllerProvider`，字段从 `srv *service.ServiceProvider` 改为 `BuildingSvc building.BuildingService`
  - 移除 `NewBuildingController`
  - 所有方法中 `s.srv.BuildingService.Xxx()` 改为 `s.BuildingSvc.Xxx()`
  - 接收者从 `*BuildingController` 改为 `*BuildingControllerProvider`

  **完整改后代码框架：**
  ```go
  package building

  import (
      "github.com/Struggle-Rabbit/CampusLogistics/api/dto"
      "github.com/Struggle-Rabbit/CampusLogistics/internal/service/building"
      "github.com/Struggle-Rabbit/CampusLogistics/pkg/utils"
      "github.com/gin-gonic/gin"
  )

  type BuildingController interface {
      Create(c *gin.Context)
      Update(c *gin.Context)
      Delete(c *gin.Context)
      GetListByPage(c *gin.Context)
      GetDetail(c *gin.Context)
      GetBuildingsByCampus(c *gin.Context)
      ImportBuildings(c *gin.Context)
      ExportBuildings(c *gin.Context)
  }

  type BuildingControllerProvider struct {
      BuildingSvc building.BuildingService
  }

  // Create 创建楼栋
  func (s *BuildingControllerProvider) Create(c *gin.Context) {
      var req dto.BuildingCreateReq
      if !utils.ShouldBind(c, &req) { return }
      if err := s.BuildingSvc.Create(&req); err != nil {
          utils.Fail(c, err.Error())
          return
      }
      utils.Success(c, "创建成功")
  }
  // ... 其他方法相同模式，s.srv.BuildingService → s.BuildingSvc
  ```

- [x] 12. **campus controller — 加接口、直接依赖 CampusService 接口**

  **Files:**
  - Modify: `internal/controller/campus/campus_controller.go`

  **Changes:** 与 Task 11 相同模式。
  - 接口: `CampusController`
  - 实现: `CampusControllerProvider`
  - 字段: `CampusSvc campus.CampusService`
  - 方法调用: `s.srv.CampusService.Xxx()` → `s.CampusSvc.Xxx()`

  **接口定义：**
  ```go
  type CampusController interface {
      Create(c *gin.Context)
      Update(c *gin.Context)
      Delete(c *gin.Context)
      GetListByPage(c *gin.Context)
      GetDetail(c *gin.Context)
      GetAll(c *gin.Context)
  }
  ```

- [x] 13. **common controller — 加接口、直接依赖 UserService 接口**

  **Files:**
  - Modify: `internal/controller/common/common_controller.go`

  **Changes:** 与 Task 11 相同模式。
  - 接口: `CommonController`
  - 实现: `CommonControllerProvider`
  - 字段: `UserSvc user.UserService`
  - 方法调用: `uCtl.srv.UserService.Xxx()` → `c.UserSvc.Xxx()`

  **接口定义：**
  ```go
  type CommonController interface {
      Login(c *gin.Context)
      Register(c *gin.Context)
  }
  ```

- [x] 14. **dorm controller — 加接口、直接依赖 DormService 接口**

  **Files:**
  - Modify: `internal/controller/dorm/dorm_controller.go`

  **Changes:** 与 Task 11 相同模式。
  - 接口: `DormController`
  - 实现: `DormControllerProvider`
  - 字段: `DormSvc dorm.DormService`
  - 方法调用: `s.srv.DormService.Xxx()` → `s.DormSvc.Xxx()`

  **接口定义：**
  ```go
  type DormController interface {
      Create(c *gin.Context)
      Update(c *gin.Context)
      Delete(c *gin.Context)
      GetListByPage(c *gin.Context)
      GetDetail(c *gin.Context)
      AssignDorm(c *gin.Context)
      TransferDorm(c *gin.Context)
      CheckOut(c *gin.Context)
      GetDormUsers(c *gin.Context)
      GetCapacityWarning(c *gin.Context)
  }
  ```

---

- [x] 15. **menu controller — 加接口、直接依赖 MenuService 接口**

  **Files:**
  - Modify: `internal/controller/menu/menu_controller.go`

  **Changes:** 与 Task 11 相同模式。
  - 接口: `MenuController`
  - 实现: `MenuControllerProvider`
  - 字段: `MenuSvc menu.MenuService`
  - 方法调用: `mc.srv.MenuService.Xxx()` → `c.MenuSvc.Xxx()`

  **接口定义：**
  ```go
  type MenuController interface {
      CreateMenu(c *gin.Context)
      DelMenu(c *gin.Context)
      UpdateMenu(c *gin.Context)
      GetListByPage(c *gin.Context)
      GetList(c *gin.Context)
      QueryDetail(c *gin.Context)
  }
  ```

- [x] 16. **notice controller — 加接口、直接依赖 NoticeService 接口**

  **Files:**
  - Modify: `internal/controller/notice/notice_controller.go`

  **Changes:** 与 Task 11 相同模式。
  - 接口: `NoticeController`
  - 实现: `NoticeControllerProvider`
  - 字段: `NoticeSvc notice.NoticeService`
  - 方法调用: `ctl.srv.NoticeService.Xxx()` → `c.NoticeSvc.Xxx()`

  **接口定义：**
  ```go
  type NoticeController interface {
      Create(c *gin.Context)
      Update(c *gin.Context)
      Delete(c *gin.Context)
      GetListByPage(c *gin.Context)
      GetPublicList(c *gin.Context)
      GetDetail(c *gin.Context)
      SetTop(c *gin.Context)
  }
  ```

- [x] 17. **repair controller — 加接口、直接依赖 RepairService 接口**

  **Files:**
  - Modify: `internal/controller/repair/repair_controller.go`

  **Changes:** 与 Task 11 相同模式。
  - 接口: `RepairController`
  - 实现: `RepairControllerProvider`
  - 字段: `RepairSvc repair.RepairService`
  - 方法调用: `ctl.srv.RepairService.Xxx()` → `c.RepairSvc.Xxx()`

  **接口定义：**
  ```go
  type RepairController interface {
      RepairOrderSubmit(c *gin.Context)
      GetListByPage(c *gin.Context)
      GetDetailById(c *gin.Context)
      UpdateRepairOrder(c *gin.Context)
      OrderRecord(c *gin.Context)
      DelRepairOrder(c *gin.Context)
  }
  ```

- [x] 18. **role controller — 加接口、直接依赖 RoleService 接口**

  **Files:**
  - Modify: `internal/controller/role/role_controller.go`

  **Changes:** 与 Task 11 相同模式。
  - 接口: `RoleController`
  - 实现: `RoleControllerProvider`
  - 字段: `RoleSvc role.RoleService`
  - 方法调用: `s.srv.RoleService.Xxx()` → `s.RoleSvc.Xxx()`

  **接口定义：**
  ```go
  type RoleController interface {
      CreateRole(c *gin.Context)
      DelRole(c *gin.Context)
      UpdateRole(c *gin.Context)
      GetListByPage(c *gin.Context)
      GetList(c *gin.Context)
      QueryDetail(c *gin.Context)
  }
  ```

---

- [x] 19. **system controller — 加接口、直接依赖 SystemService 接口**

  **Files:**
  - Modify: `internal/controller/system/system_controller.go`

  **Changes:** 与 Task 11 相同模式。
  - 接口: `SystemController`
  - 实现: `SystemControllerProvider`
  - 字段: `SystemSvc system.SystemService`
  - 方法调用: `s.srv.SystemService.Xxx()` → `s.SystemSvc.Xxx()`

  **接口定义：**
  ```go
  type SystemController interface {
      GetOperationLogListByPage(c *gin.Context)
      RefreshToken(c *gin.Context)
  }
  ```

- [x] 20. **user controller — 加接口、直接依赖 UserService 接口**

  **Files:**
  - Modify: `internal/controller/user/user_controller.go`

  **Changes:** 与 Task 11 相同模式。
  - 接口: `UserController`
  - 实现: `UserControllerProvider`
  - 字段: `UserSvc user.UserService`
  - 方法调用: `s.srv.UserService.Xxx()` → `s.UserSvc.Xxx()`

  **接口定义：**
  ```go
  type UserController interface {
      DelUser(c *gin.Context)
      UpdateUser(c *gin.Context)
      ResetPassword(c *gin.Context)
      GetListByPage(c *gin.Context)
      QueryDetail(c *gin.Context)
      GetUserPermission(c *gin.Context)
  }
  ```

- [x] 21. **utility controller — 加接口、直接依赖 UtilityService 接口**

  **Files:**
  - Modify: `internal/controller/utility/utility_controller.go`

  **Changes:** 与 Task 11 相同模式。
  - 接口: `UtilityController`
  - 实现: `UtilityControllerProvider`
  - 字段: `UtilitySvc utility.UtilityService`
  - 方法调用: `s.srv.UtilityService.Xxx()` → `s.UtilitySvc.Xxx()`

  **接口定义：**
  ```go
  type UtilityController interface {
      Create(c *gin.Context)
      Update(c *gin.Context)
      Delete(c *gin.Context)
      GetListByPage(c *gin.Context)
      GetDetail(c *gin.Context)
      Pay(c *gin.Context)
      BatchPay(c *gin.Context)
      UpdatePrice(c *gin.Context)
      GetPrice(c *gin.Context)
      GetStatistics(c *gin.Context)
      GetUnpaidWarning(c *gin.Context)
      GetUserDormUtility(c *gin.Context)
  }
  ```

---

### Wave 3: Router 层 — 删除 provider.go + 更新路由签名 (12 个文件)

---

- [x] 22. **删除 provider.go**

  **Files:**
  - Delete: `internal/service/provider.go`

  **Changes:** 直接删除整个文件。该文件中的 `ServiceProvider` 结构体和 `NewServiceProvider` 函数不再需要。

- [x] 23. **更新 router.go — 直接创建所有 service 实例并分发**

  **Files:**
  - Modify: `internal/router/router.go`

  **Changes:**
  - 移除 `import "github.com/.../internal/service"`
  - 添加所有 10 个 service 包的 import
  - 在 `initRouter` 中，去掉 `srv := service.NewServiceProvider(app)`，改为直接创建每个 service
  - 将 `srv` 传给各路由函数改为传具体的 service 实例

  ```go
  package router

  import (
      "fmt"
      "github.com/Struggle-Rabbit/CampusLogistics/internal/app"
      "github.com/Struggle-Rabbit/CampusLogistics/internal/config"
      "github.com/Struggle-Rabbit/CampusLogistics/internal/middleware"
      "github.com/Struggle-Rabbit/CampusLogistics/internal/service/building"
      "github.com/Struggle-Rabbit/CampusLogistics/internal/service/campus"
      "github.com/Struggle-Rabbit/CampusLogistics/internal/service/dorm"
      "github.com/Struggle-Rabbit/CampusLogistics/internal/service/menu"
      "github.com/Struggle-Rabbit/CampusLogistics/internal/service/notice"
      "github.com/Struggle-Rabbit/CampusLogistics/internal/service/repair"
      "github.com/Struggle-Rabbit/CampusLogistics/internal/service/role"
      "github.com/Struggle-Rabbit/CampusLogistics/internal/service/system"
      "github.com/Struggle-Rabbit/CampusLogistics/internal/service/user"
      "github.com/Struggle-Rabbit/CampusLogistics/internal/service/utility"
      "github.com/Struggle-Rabbit/CampusLogistics/pkg/logger"
      "github.com/gin-gonic/gin"
      "go.uber.org/zap"
      _ "github.com/Struggle-Rabbit/CampusLogistics/docs"
      swaggerFiles "github.com/swaggo/files"
      ginSwagger "github.com/swaggo/gin-swagger"
  )

  func initRouter(app *app.App) *gin.Engine {
      r := gin.Default()
      r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
      r.Use(middleware.Recovery(), middleware.RequestID(), middleware.Logger(), middleware.CORS())

      menuSvc := &menu.MenuServiceProvider{App: app}
      roleSvc := &role.RoleServiceProvider{App: app}
      userSvc := &user.UserServiceProvider{App: app, Menu: menuSvc}
      systemSvc := &system.SystemServiceProvider{App: app}
      repairSvc := &repair.RepairServiceProvider{App: app}
      campusSvc := &campus.CampusServiceProvider{App: app}
      buildingSvc := &building.BuildingServiceProvider{App: app}
      dormSvc := &dorm.DormServiceProvider{App: app}
      utilitySvc := &utility.UtilityServiceProvider{App: app}
      noticeSvc := &notice.NoticeServiceProvider{App: app}

      api := r.Group("/api/v1")
      api.Use(middleware.OperationLogMiddleware())
      {
          LoadCommonRouter(api, userSvc, noticeSvc)
          api.Use(middleware.JWTAuth())
          {
              LoadUserRouter(api, userSvc)
              LoadSystemRouter(api, systemSvc)
              LoadRoleRouter(api, roleSvc)
              LoadMenuRouter(api, menuSvc)
              LoadRepairRouter(api, repairSvc)
              LoadCampusRouter(api, campusSvc)
              LoadBuildingRouter(api, buildingSvc)
              LoadDormRouter(api, dormSvc)
              LoadUtilityRouter(api, utilitySvc)
              LoadNoticeRouter(api, noticeSvc)
          }
      }
      return r
  }
  // Run 函数不变
  ```

- [x] 24. **更新 common.go 路由**

  **Files:**
  - Modify: `internal/router/common.go`

  **Changes:**
  - 签名: `(api *gin.RouterGroup, userSvc user.UserService, noticeSvc notice.NoticeService)`
  - `common.NewCommonController(srv)` → `&common.CommonControllerProvider{UserSvc: userSvc}`
  - `notice.NewNoticeController(srv)` → `&notice.NoticeControllerProvider{NoticeSvc: noticeSvc}`

  ```go
  package router

  import (
      "github.com/Struggle-Rabbit/CampusLogistics/internal/controller/common"
      "github.com/Struggle-Rabbit/CampusLogistics/internal/controller/notice"
      "github.com/Struggle-Rabbit/CampusLogistics/internal/service/notice"
      "github.com/Struggle-Rabbit/CampusLogistics/internal/service/user"
      "github.com/gin-gonic/gin"
  )

  func LoadCommonRouter(api *gin.RouterGroup, userSvc user.UserService, noticeSvc notice.NoticeService) {
      commonCtl := &common.CommonControllerProvider{UserSvc: userSvc}
      noticeCtl := &notice.NoticeControllerProvider{NoticeSvc: noticeSvc}
      api.POST("/register", commonCtl.Register)
      api.POST("/login", commonCtl.Login)
      noticeGroup := api.Group("/notice")
      {
          noticeGroup.GET("/public", noticeCtl.GetPublicList)
      }
  }
  ```

- [x] 25. **更新 campus.go 路由（含 LoadCampusRouter + LoadBuildingRouter）**

  **Files:**
  - Modify: `internal/router/campus.go`

  **Changes:**
  - `LoadCampusRouter(api, campusSvc campus.CampusService)`，`campus.NewCampusController(srv)` → `&campus.CampusControllerProvider{CampusSvc: campusSvc}`
  - `LoadBuildingRouter(api, buildingSvc building.BuildingService)`，`building.NewBuildingController(srv)` → `&building.BuildingControllerProvider{BuildingSvc: buildingSvc}`

- [x] 26. **更新 user.go 路由**

  **Files:**
  - Modify: `internal/router/user.go`
  - 签名: `(api *gin.RouterGroup, userSvc user.UserService)`
  - `user.NewUserController(srv)` → `&user.UserControllerProvider{UserSvc: userSvc}`

- [x] 27. **更新 system.go 路由**

  **Files:**
  - Modify: `internal/router/system.go`
  - 签名: `(api *gin.RouterGroup, systemSvc system.SystemService)`
  - `system.NewSystemController(srv)` → `&system.SystemControllerProvider{SystemSvc: systemSvc}`

- [x] 28. **更新 role.go 路由**

  **Files:**
  - Modify: `internal/router/role.go`
  - 签名: `(api *gin.RouterGroup, roleSvc role.RoleService)`
  - `role.NewRoleController(srv)` → `&role.RoleControllerProvider{RoleSvc: roleSvc}`

- [x] 29. **更新 menu.go 路由**

  **Files:**
  - Modify: `internal/router/menu.go`
  - 签名: `(api *gin.RouterGroup, menuSvc menu.MenuService)`
  - `menu.NewMenuController(srv)` → `&menu.MenuControllerProvider{MenuSvc: menuSvc}`

- [x] 30. **更新 repair.go 路由**

  **Files:**
  - Modify: `internal/router/repair.go`
  - 签名: `(api *gin.RouterGroup, repairSvc repair.RepairService)`
  - `repair.NewRepairController(srv)` → `&repair.RepairControllerProvider{RepairSvc: repairSvc}`

- [x] 31. **更新 dorm.go 路由**

  **Files:**
  - Modify: `internal/router/dorm.go`
  - 签名: `(api *gin.RouterGroup, dormSvc dorm.DormService)`
  - `dorm.NewDormController(srv)` → `&dorm.DormControllerProvider{DormSvc: dormSvc}`

- [x] 32. **更新 utility.go 路由**

  **Files:**
  - Modify: `internal/router/utility.go`
  - 签名: `(api *gin.RouterGroup, utilitySvc utility.UtilityService)`
  - `utility.NewUtilityController(srv)` → `&utility.UtilityControllerProvider{UtilitySvc: utilitySvc}`

- [x] 33. **更新 notice.go 路由**

  **Files:**
  - Modify: `internal/router/notice.go`
  - 签名: `(api *gin.RouterGroup, noticeSvc notice.NoticeService)`
  - `notice.NewNoticeController(srv)` → `&notice.NoticeControllerProvider{NoticeSvc: noticeSvc}`

---

### Wave 4: 测试层更新 (11 个文件)

- [ ] 34. **更新 user_service_test.go**

  **Files:**
  - Modify: `test/unit_test/user_service_test.go`

  **Changes:**
  - `user.NewUserService(appInstance, menuSvc)` → `&user.UserServiceProvider{App: appInstance, Menu: menuSvc}`

  ```go
  // 在 TestUserService_Register, TestUserService_Login, TestUserService_ResetPassword 中
  // 重构前:
  // menuSvc := menu.NewMenuService(appInstance)
  // svc := user.NewUserService(appInstance, menuSvc)
  //
  // 重构后:
  menuSvc := &menu.MenuServiceProvider{App: appInstance}
  svc := &user.UserServiceProvider{App: appInstance, Menu: menuSvc}
  ```

- [ ] 35. **更新 campus_building_service_test.go**

  **Files:**
  - Modify: `test/unit_test/campus_building_service_test.go`

  **Changes:**
  - `campus.NewCampusService(appInstance)` → `&campus.CampusServiceProvider{App: appInstance}`
  - `building.NewBuildingService(appInstance)` → `&building.BuildingServiceProvider{App: appInstance}`

- [ ] 36. **更新 dorm_service_test.go**

  **Files:**
  - Modify: `test/unit_test/dorm_service_test.go`
  - `dorm.NewDormService(appInstance)` → `&dorm.DormServiceProvider{App: appInstance}`

- [ ] 37. **更新 menu_service_test.go**

  **Files:**
  - Modify: `test/unit_test/menu_service_test.go`
  - `menu.NewMenuService(appInstance)` → `&menu.MenuServiceProvider{App: appInstance}`

- [ ] 38. **更新 notice_service_test.go**

  **Files:**
  - Modify: `test/unit_test/notice_service_test.go`
  - `notice.NewNoticeService(appInstance)` → `&notice.NoticeServiceProvider{App: appInstance}`

- [ ] 39. **更新 repair_controller_test.go**

  **Files:**
  - Modify: `test/unit_test/repair_controller_test.go`

  **Changes:**
  - 移除 `import "github.com/.../internal/service"`
  - `srvProvider := service.NewServiceProvider(appInstance)` → 直接创建 repairSvc
  - `repair.NewRepairController(srvProvider)` → `&repair.RepairControllerProvider{RepairSvc: repairSvc}`

  ```go
  // 重构前:
  // srvProvider := service.NewServiceProvider(appInstance)
  // ctl := repair.NewRepairController(srvProvider)
  //
  // 重构后:
  repairSvc := &repair.RepairServiceProvider{App: appInstance}
  ctl := &repair.RepairControllerProvider{RepairSvc: repairSvc}
  ```

- [ ] 40. **更新 repair_service_test.go**

  **Files:**
  - Modify: `test/unit_test/repair_service_test.go`
  - `repair.NewRepairService(appInstance)` → `&repair.RepairServiceProvider{App: appInstance}`

- [ ] 41. **更新 role_service_test.go**

  **Files:**
  - Modify: `test/unit_test/role_service_test.go`
  - `role.NewRoleService(appInstance)` → `&role.RoleServiceProvider{App: appInstance}`

- [ ] 42. **更新 system_service_test.go**

  **Files:**
  - Modify: `test/unit_test/system_service_test.go`
  - `system.NewSystemService(appInstance)` → `&system.SystemServiceProvider{App: appInstance}`

- [ ] 43. **更新 utility_service_test.go**

  **Files:**
  - Modify: `test/unit_test/utility_service_test.go`
  - `utility.NewUtilityService(appInstance)` → `&utility.UtilityServiceProvider{App: appInstance}`

---

## Final Verification Wave

- [ ] F1. **编译验证**
  Run: `go build ./...`
  Expected: 编译成功，无错误
  Evidence: `.omo/evidence/build-output.txt`

- [ ] F2. **测试验证**
  Run: `go test ./test/unit_test/ -v`
  Expected: 全部 PASS
  Evidence: `.omo/evidence/test-output.txt`

- [ ] F3. **代码审查**
  - 确认所有 service 包都有接口定义
  - 确认所有 controller 包都有接口定义
  - 确认 `internal/service/provider.go` 已删除
  - 确认没有剩余的 `NewXxxService` / `NewXxxController` 调用
  - 确认没有 `*service.ServiceProvider` 类型的引用
  Evidence: `.omo/evidence/review-log.md`

---

## Commit Strategy

分组提交，每 Wave 一个 commit：

1. `refactor(service): add interfaces and rename implementations to Provider pattern`
   - All Wave 1 service changes
2. `refactor(controller): add interfaces and rename implementations to Provider pattern`
   - All Wave 2 controller changes
3. `refactor(router): remove ServiceProvider, wire services directly`
   - All Wave 3 router changes
4. `refactor(test): update tests to use direct Provider initialization`
   - All Wave 4 test changes

---

## Success Criteria

### Verification Commands
```bash
go build ./...          # 编译成功
go test ./test/unit_test/ -v   # 全部测试通过
```

### Final Checklist
- [ ] `go build ./...` 编译通过
- [ ] `go test ./test/unit_test/ -v` 全部 PASS
- [ ] `internal/service/provider.go` 已删除
- [ ] 所有 `NewXxxService` / `NewXxxController` 函数已移除
- [ ] 所有 controller 直接依赖 Service 接口，不再经过 ServiceProvider
- [ ] 所有 router 函数签名为 Service 接口参数
