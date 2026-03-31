package dto

// RegisterReq 注册请求参数
type RegisterReq struct {
	UserCode string `json:"user_code" binding:"required"` // 学号/工号
	Name     string `json:"name" binding:"required"`
	Phone    string `json:"phone" binding:"required,phone"`
	Password string `json:"password" binding:"required,min=8,max=20"`
	UserType int    `json:"user_type" binding:"required,oneof=1 2"` // 1-学生 2-教职工
}
