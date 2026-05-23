# 数据模型层 — internal/model/

13 个 GORM 模型 + 3 个基础嵌入模型。所有模型通过 GORM AutoMigrate 自动建表（开发环境）。

## 模型关系图

```
Campus (1) ──→ Building (N) ──→ DormRoom (N) ──→ DormUser (N) ──→ SysUser (N)
                                              ├── DormUtility (N)
                                              └── UtilityPrice (1)

SysRole (N) ←── sys_user_role ──→ SysUser (N)
SysRole (N) ←── sys_role_menu ──→ SysMenu (N)

SysUser (1) ──→ RepairOrder (N) ──→ RepairRecord (N)
SysUser (1) ──→ Notice (N)
SysUser (N) ──→ SysOperationLog (N)
```

## 基础嵌入模型

| 模型 | 字段 | 用途 |
|------|------|------|
| `BaseModel` | ID(string,Snowflake), CreatedAt, UpdatedAt | 普通表 |
| `BaseModelWithDelete` | 同 BaseModel + DeletedAt(gorm.DeletedAt) | 需软删除的表 |
| `BaseModelIntId` | ID(uint,自增), CreatedAt, UpdatedAt | 备用 |

ID 通过 `BeforeCreate` GORM hook 自动调用 `utils.GenStringID()` 生成。

## 领域模型清单

### 校区宿舍域
| 模型 | 表名 | 关键字段 | 关联 |
|------|------|---------|------|
| `Campus` | `campus` | campus_name, address, contact, phone | 1:N → Building |
| `Building` | `building` | campus_id, building_no, floor_count, room_count | N:1 ← Campus; 1:N → DormRoom |
| `DormRoom` | `dorm_room` | building_id, room_no, floor, room_type, max_count, current_count | N:1 ← Building; 1:N → DormUser |
| `DormUser` | `dorm_user` | room_id, user_id, check_in_time, check_out_time, status | N:1 ← DormRoom; N:1 → SysUser |
| `DormUtility` | `dorm_utility` | room_id, year, month, water_usage, electric_usage, amount, pay_status | N:1 ← DormRoom |
| `UtilityPrice` | `utility_price` | water_price, electric_price | 全局单例配置 |

### 用户权限域
| 模型 | 表名 | 关键字段 | 关联 |
|------|------|---------|------|
| `SysUser` | `sys_user` | user_code(唯一), name, mobile(唯一), password(bcrypt), status, user_type, refresh_token | N:M → SysRole (sys_user_role) |
| `SysRole` | `sys_role` | role_name, role_code(唯一), status, is_built_in | N:M → SysUser (sys_user_role); N:M → SysMenu (sys_role_menu) |
| `SysMenu` | `sys_menu` | parent_id, name, path, component, type, perms(唯一), icon, sort, status | 自引用(树形); N:M → SysRole (sys_role_menu) |

### 后勤服务域
| 模型 | 表名 | 关键字段 | 关联 |
|------|------|---------|------|
| `Notice` | `notice` | title, content, notice_type, is_top, publish_time, view_count, creator_id, attachments(JSON) | 软删除 |
| `RepairOrder` | `repair_order` | order_no(唯一), user_id, repair_type, address, description, images(JSON), status(6状态), handler_id | 1:N → RepairRecord |
| `RepairRecord` | `repair_record` | order_id, operator_id, old_status, new_status, remark | N:1 ← RepairOrder |
| `SysOperationLog` | `operation_log` | user_id, method, path, params, status_code, ip, user_agent, operation_at | - |

## CONVENTIONS

- 字符串 ID（Snowflake），非自增整数
- 软删除模型继承 `BaseModelWithDelete`
- JSON 字段使用 GORM JSON 序列化器（`serializer:json`）
- 密码字段在 JSON 序列化时隐藏（`json:"-"`）
- 所有表名通过 `TableName()` 方法显式指定

## ANTI-PATTERNS

- `Building`、`DormRoom`、`Campus` 定义在 `campus.go` 同一文件中，文件名与实际模型名不一致
- 表名命名不一致：`sys_user`/`sys_role` 使用 `sys_` 前缀，但 `operation_log` 无前缀
