package dto

import "time"

type OperationLogReq struct {
	UserID      string    `json:"user_id"`
	UserName    string    `json:"user_name"`
	IP          string    `json:"ip"`
	OperationAt time.Time `json:"operation_at"`
}

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
	OperationLogReq
}
