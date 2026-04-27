package dto

import "time"

type OperationLogResult struct {
	ID          string    `json:"id"`
	UserID      string    `json:"user_id"`
	UserName    string    `json:"user_name"`
	Method      string    `json:"method"`
	Path        string    `json:"path"`
	Params      string    `json:"params"`
	StatusCode  int       `json:"status_code"`
	IP          string    `json:"ip"`
	UserAgent   string    `json:"user_agent"`
	OperationAt time.Time `json:"operation_at"`
}

type OperationLogByPageReq struct {
	PageReq
	UserID             string    `json:"user_id"`
	OperationTimeStart time.Time `json:"operation_time_start"`
	OperationTimeEnd   time.Time `json:"operation_time_end"`
	IP                 string    `json:"ip"`
}

type RefreshTokenResult struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}
