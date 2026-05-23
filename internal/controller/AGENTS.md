# 控制器层 — internal/controller/

11 个领域包，处理 HTTP 请求解析、参数校验、调用 service 并返回统一 JSON 响应。

## 控制器列表

| 包 | 文件 | 核心接口 | 权限标识 |
|----|------|---------|---------|
| `common` | `common_controller.go` | POST `/login`, `/register` | 无需认证（JWTAuth 之前） |
| `user` | `user_controller.go` | 列表/详情/编辑/删除/重置密码/权限查询 | `sys:user:*` |
| `system` | `system_controller.go` | 操作日志分页查询 | `sys:optLog` |
| `role` | `role_controller.go` | 角色 CRUD | `sys:role:*` |
| `menu` | `menu_controller.go` | 菜单 CRUD（树形） | `sys:menu:*` |
| `repair` | `repair_controller.go` | 提交/列表/详情/更新/流转/删除 | `repair:*` |
| `campus` | `campus_controller.go` | 校区 CRUD | `campus:*` |
| `building` | `building_controller.go` | 楼栋 CRUD + 导入/导出 | `building:*` |
| `dorm` | `dorm_controller.go` | 宿舍 CRUD + 分配/调宿/迁出/预警 | `dorm:*` |
| `utility` | `utility_controller.go` | 水电费 CRUD + 缴费/单价/统计/预警 | `utility:*` |
| `notice` | `notice_controller.go` | 公告 CRUD + 置顶 | `notice:*` |

## 通用模式

- 所有控制器通过 `ServiceProvider` 注入，不直接创建 service
- 请求绑定 + 校验使用 `utils.ShouldBind(c, &req)`
- 成功返回 `utils.Success(c, data, msg)`，失败返回 `utils.Fail(c, errMsg)`
- 路由注册在 `internal/router/` 中完成，控制器不涉及 HTTP 方法/路径配置

## WHERE TO LOOK

| 需求 | 操作 |
|------|------|
| 新增接口 | 在对应包添加 handler → 在 router 注册路由+权限 |
| 修改参数校验 | 编辑对应 DTO（`api/dto/`）+ controller 中的 binding tag |
| 调整权限粒度 | 修改 router 中的 `PermissionValidator(perms)` 参数 |
