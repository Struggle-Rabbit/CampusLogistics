package dto

import "time"

type CreateMenuReq struct {
	ParentID    string `json:"parent_id"`
	Name        string `json:"name" binding:"required"`  // 菜单名称
	Path        string `json:"path"`                     // 路由
	Component   string `json:"component"`                // 前端组件
	Type        int    `json:"type" binding:"required"`  // 1目录2菜单3按钮
	Perms       string `json:"perms" binding:"required"` // 权限标识 sys:user:list
	Icon        string `json:"icon"`
	Sort        int    `json:"sort"`
	Status      int    `json:"status" binding:"required"` // 状态
	Description string `json:"description"`               // 描述
}

type UpdateMenuReq struct {
	ID          string `json:"id" binding:"required"`
	ParentID    string `json:"parent_id"`
	Name        string `json:"name"`      // 菜单名称
	Path        string `json:"path"`      // 路由
	Component   string `json:"component"` // 前端组件
	Type        int    `json:"type"`      // 1目录2菜单3按钮
	Perms       string `json:"perms"`     // 权限标识 sys:user:list
	Icon        string `json:"icon"`
	Sort        int    `json:"sort"`
	Status      int    `json:"status"`      // 状态
	Description string `json:"description"` // 描述
}

type MenuListReq struct {
	ParentID string `json:"parent_id"`
	Name     string `json:"name"`   // 菜单名称
	Type     int    `json:"type"`   // 1目录2菜单3按钮
	Perms    string `json:"perms"`  // 权限标识 sys:user:list
	Status   int    `json:"status"` // 状态
}

type MenuListByPageReq struct {
	PageReq
	MenuListReq
}

type MenuResult struct {
	ID          string        `json:"id"`
	ParentID    string        `json:"parent_id"`
	Name        string        `json:"name"`      // 菜单名称
	Path        string        `json:"path"`      // 路由
	Component   string        `json:"component"` // 前端组件
	Type        int           `json:"type"`      // 1目录2菜单3按钮
	Perms       string        `json:"perms"`     // 权限标识 sys:user:list
	Icon        string        `json:"icon"`
	Sort        int           `json:"sort"`
	Status      int           `json:"status"`      // 状态
	Description string        `json:"description"` // 描述
	CreatedAt   time.Time     `json:"created_at"`
	UpdatedAt   time.Time     `json:"updated_at"`
	Children    []*MenuResult `json:"childen"`
}
