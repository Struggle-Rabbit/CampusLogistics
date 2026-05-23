# 服务层 — internal/service/

10 个领域包 + 1 个注入器，封装核心业务逻辑。所有 service 在 `ServiceProvider` 中统一初始化。

## 服务列表

| 服务 | 包路径 | 核心函数 | 特殊说明 |
|------|--------|---------|---------|
| `UserService` | `user/` | `Register`, `Login`, `DelUser`, `UpdateUser`, `ResetPassword`, `GetListByPage`, `GetUserInfo`, `GetUserPermission` | 依赖 `MenuService`；注册使用 bcrypt；登录签发 JWT（access+refresh）；密码重置校验旧密码；`GetUserInfo` 直接从 `gin.Context` 取 userID（违反分层） |
| `RoleService` | `role/` | 角色 CRUD，角色-菜单关联 | 内置角色（`is_built_in=1`）不可删除 |
| `MenuService` | `menu/` | 菜单 CRUD（返回树形结构），菜单-角色关联 | 菜单类型：1目录 2菜单 3按钮 |
| `SystemService` | `system/` | `GetOperationLogListByPage` | 操作日志纯查询 |
| `RepairService` | `repair/` | 提交/列表/详情/更新/流转/删除 | 6 状态流转（待分配→待处理→处理中→已完成/已驳回/已撤销）；流转写入 `RepairRecord` |
| `CampusService` | `campus/` | 校区 CRUD + 全量列表 | 基础 CRUD |
| `BuildingService` | `building/` | 楼栋 CRUD + 按校区查询 + 批量导入/导出 | 支持 Excel 导入/导出 |
| `DormService` | `dorm/` | 宿舍 CRUD + 分配/调宿/迁出/容量预警 | 入住更新 `dorm_user` 关联表；迁出标记状态；预警检测超员宿舍 |
| `UtilityService` | `utility/` | 水电费 CRUD + 缴费/批量缴费/单价管理/统计/预警 | 录入时自动计算金额；同宿舍同月份不可重复；预警检测长期未缴费 |
| `NoticeService` | `notice/` | 公告 CRUD + 置顶控制 + 公开列表 | 置顶最多 3 条；公开列表无需认证 |

## 依赖注入

`ServiceProvider`（`provider.go`）在 `router.go` 的 `initRouter()` 中创建：

```
menuSvc → roleSvc → userSvc(依赖menuSvc) → systemSvc → ... → ServiceProvider
```

所有 controller 共享同一个 ServiceProvider 实例。

## CONVENTIONS

- Service 只通过 `app.App` 获取 `*gorm.DB`，不直接引用 DAO
- 事务在 service 层控制（使用 `tx := app.DB.Begin()`）
- 所有 service 方法返回 `(T, error)` 或 `error`

## ANTI-PATTERNS

- `UserService.GetUserInfo()` 直接接收 `*gin.Context` 而非 userID 字符串，破坏分层隔离
- 部分 service 方法未使用事务（如批量操作可能需要回滚）
