package global

import "akali/global/errInfos"

// Gin server 響應回傳格式
type ApiResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Datas   any    `json:"datas"`
}

// Gin server 處理結果
type Ret struct {
	Status  int
	Datas   any
	ErrInfo *errInfos.Err
	Err     error
}

// 全專案 Panic 處理後的結構
type Panic struct {
	Panic      string `json:"panic"`
	StackTrace string `json:"stack_trace"`
}

type Pagination struct {
	Total int64
	Page  int64
	Limit int64
}

type TimeRange struct {
	StartTime int64
	EndTime   int64
}
