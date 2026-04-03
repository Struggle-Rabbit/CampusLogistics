package dto

import "github.com/Struggle-Rabbit/CampusLogistics/internal/model"

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
	Name     string `json:"name" binding:"required"`
	Mobile   string `json:"mobile" binding:"required,mobile"`
	Password string `json:"password" binding:"required,min=8,max=20"`
	UserType string `json:"userType" binding:"required,oneof=1 2"`
}

// LoginReq 登录请求参数
type LoginReq struct {
	Account  string `json:"account" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// LoginRes 登录响应参数
type LoginResult struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
}

// UserInfo 用户信息响应
type UserInfoResult struct {
	model.SysUser
}
