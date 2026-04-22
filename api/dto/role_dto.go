package dto

import "time"

type CreateRoleReq struct {
	RoleName    string `json:"role_name" binding:"required"`
	RoleCode    string `json:"role_code" binding:"required"`
	Status      string `json:"status" binding:"required"`
	IsBuiltIn   int    `json:"is_built_in"`
	Description string `json:"description"`
}

type UpdateRoleReq struct {
	ID          string `json:"id" binding:"required"`
	RoleName    string `json:"role_name"`
	RoleCode    string `json:"role_code"`
	Status      string `json:"status"`
	IsBuiltIn   int    `json:"is_built_in"`
	Description string `json:"description"`
}

type RoleListReq struct {
	RoleName string `json:"role_name"`
	Status   string `json:"status"`
}

type RoleListByPageReq struct {
	PageReq
	RoleListReq
}

type RoleResult struct {
	ID          string    `json:"id"`
	RoleName    string    `json:"role_name"`
	RoleCode    string    `json:"role_code"`
	Status      string    `json:"status"`
	IsBuiltIn   int       `json:"is_built_in"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
