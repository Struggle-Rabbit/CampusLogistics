package dto

import "time"

// UserListPage 分页列表请求参数
type UserListPageReq struct {
	UserListReq
	PageReq
}

// UserList 列表请求参数
type UserListReq struct {
	Name     string `json:"name"`
	Mobile   string `json:"mobile"`
	UserType string `json:"userType"`
	Status   string `json:"status"`
}

// RegisterReq 注册请求参数
type RegisterReq struct {
	Name     string `json:"name" binding:"required" example:"张三"`
	Mobile   string `json:"mobile" binding:"required,mobile" example:"13811111111"`
	Password string `json:"password" binding:"required,min=8,max=20" example:"12345678"`
	UserType string `json:"userType" binding:"required,oneof=00 01 02" example:"02"`
}

// LoginReq 登录请求参数
type LoginReq struct {
	Account  string `json:"account" binding:"required" example:"13811111111"`
	Password string `json:"password" binding:"required" example:"12345678"`
}

// LoginRes 登录响应参数
type LoginResult struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
}

// UserInfo 用户信息响应
type UserInfoResult struct {
	ID        string        `json:"id"`
	UserCode  string        `json:"user_code"`
	Name      string        `json:"name"`
	Mobile    string        `json:"mobile"`
	RoleIDs   []string      `json:"role_ids"`
	Status    int           `json:"status"`
	Avatar    string        `json:"avatar"`
	UserType  string        `json:"user_type"`
	CreatedAt time.Time     `json:"created_at"`
	UpdatedAt time.Time     `json:"updated_at"`
	Roles     []*RoleResult `json:"roles"`
}

type UserUpdateReq struct {
	ID       string `json:"id"`
	Name     string `json:"name"`      // 姓名
	Mobile   string `json:"mobile"`    // 手机号（唯一）
	Status   int    `json:"status"`    // 状态：1-启用 2-禁用
	Avatar   string `json:"avatar"`    // 头像URL
	UserType string `json:"user_type"` // 用户类型: 00-管理员  01-职工  02-学生
}

type PasswordReset struct {
	OldPassword string `json:"old_password"`
	NewPassword string `json:"new_password"`
	Mobile      string `json:"mobile"`
}

type UserPermissionResult struct {
	UserId  string       `json:"user_id"`
	RoleIDs []string     `json:"role_ids"`
	MenuIDs []string     `json:"menu_ids"`
	Roles   []RoleResult `json:"roles"`
	Menus   []MenuResult `json:"menus"`
}
