package global

import "timeLedger/global/errInfos"

// Gin server 響應回傳格式
type ApiResponse struct {
	Code    errInfos.ErrCode `json:"code"`
	Message string           `json:"message"`
	Datas   any              `json:"datas"`
}

// Gin server 處理結果
type Ret struct {
	Status  int
	Datas   any
	ErrInfo *errInfos.Res
	Err     error
}

// 全專案 Panic 處理後的結構
type Panic struct {
	Panic      string `json:"panic"`
	StackTrace string `json:"stack_trace"`
}

type Pagination struct {
	Page       int64 `json:"page"`
	Limit      int64 `json:"limit"`
	Total      int64 `json:"total"`
	TotalPages int64 `json:"total_pages"`
	HasNext    bool  `json:"has_next"`
	HasPrev    bool  `json:"has_prev"`
}

func NewPagination(page, limit, total int64) Pagination {
	totalPages := (total + limit - 1) / limit
	return Pagination{
		Page:       page,
		Limit:      limit,
		Total:      total,
		TotalPages: totalPages,
		HasNext:    page < totalPages,
		HasPrev:    page > 1,
	}
}

type TimeRange struct {
	StartTime int64
	EndTime   int64
}
