package dto

type PageResult struct {
	List        interface{} `json:"list"`        // 数据列表
	Total       int64       `json:"total"`       // 总条数
	CurrentPage int         `json:"currentPage"` // 当前页
	PageSize    int         `json:"pageSize"`    // 每页条数
}

type PageReq struct {
	CurrentPage int `json:"currentPage"` // 当前页
	PageSize    int `json:"pageSize"`    // 每页条数
}
