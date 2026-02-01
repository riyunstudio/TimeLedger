package services

import (
	"timeLedger/app"
	"timeLedger/global/logger"
)

// BaseService 基礎服務結構，提供通用功能
type BaseService struct {
	App        *app.App
	Logger     *ServiceLogger
}

// NewBaseService 建立基礎服務
func NewBaseService(app *app.App, component string) *BaseService {
	return &BaseService{
		App:    app,
		Logger: NewServiceLogger(component),
	}
}

// NewBaseServiceWithLogger 建立基礎服務（自定義日誌器）
func NewBaseServiceWithLogger(app *app.App, log *ServiceLogger) *BaseService {
	return &BaseService{
		App:    app,
		Logger: log,
	}
}

// PaginationParams 分頁參數
type PaginationParams struct {
	Page     int    `json:"page" form:"page"`
	Limit    int    `json:"limit" form:"limit"`
	SortBy   string `json:"sort_by" form:"sort_by"`
	SortOrder string `json:"sort_order" form:"sort_order"`
}

// DefaultPagination 預設分頁設定
func DefaultPagination() *PaginationParams {
	return &PaginationParams{
		Page:     1,
		Limit:    20,
		SortBy:   "id",
		SortOrder: "DESC",
	}
}

// Validate 驗證並修正分頁參數
func (p *PaginationParams) Validate() {
	if p.Page < 1 {
		p.Page = 1
	}
	if p.Limit < 1 {
		p.Limit = 20
	}
	if p.Limit > 100 {
		p.Limit = 100
	}
	if p.SortBy == "" {
		p.SortBy = "id"
	}
	if p.SortOrder != "ASC" && p.SortOrder != "DESC" {
		p.SortOrder = "DESC"
	}
}

// GetOffset 計算分頁偏移量
func (p *PaginationParams) GetOffset() int {
	return (p.Page - 1) * p.Limit
}

// BuildOrderClause 建立排序子句
func (p *PaginationParams) BuildOrderClause() string {
	return p.SortBy + " " + p.SortOrder
}

// PaginationResult 分頁結果
type PaginationResult struct {
	Data       interface{} `json:"data"`
	Total      int64       `json:"total"`
	Page       int         `json:"page"`
	Limit      int         `json:"limit"`
	TotalPages int         `json:"total_pages"`
	HasNext    bool        `json:"has_next"`
	HasPrev    bool        `json:"has_prev"`
}

// NewPaginationResult 建立分頁結果
func NewPaginationResult(data interface{}, total int64, params *PaginationParams) *PaginationResult {
	totalPages := int(total) / params.Limit
	if int(total)%params.Limit > 0 {
		totalPages++
	}

	return &PaginationResult{
		Data:       data,
		Total:      total,
		Page:       params.Page,
		Limit:      params.Limit,
		TotalPages: totalPages,
		HasNext:    params.Page < totalPages,
		HasPrev:    params.Page > 1,
	}
}

// FilterCondition 過濾條件
type FilterCondition struct {
	Field    string
	Operator string // "eq", "ne", "gt", "gte", "lt", "lte", "like", "in", "between"
	Value    interface{}
}

// FilterBuilder 動態過濾建構器
type FilterBuilder struct {
	conditions []FilterCondition
}

// NewFilterBuilder 建立新的過濾建構器
func NewFilterBuilder() *FilterBuilder {
	return &FilterBuilder{
		conditions: make([]FilterCondition, 0),
	}
}

// AddEq 新增等於條件
func (fb *FilterBuilder) AddEq(field string, value interface{}) *FilterBuilder {
	fb.conditions = append(fb.conditions, FilterCondition{
		Field:    field,
		Operator: "eq",
		Value:    value,
	})
	return fb
}

// AddNe 新增不等於條件
func (fb *FilterBuilder) AddNe(field string, value interface{}) *FilterBuilder {
	fb.conditions = append(fb.conditions, FilterCondition{
		Field:    field,
		Operator: "ne",
		Value:    value,
	})
	return fb
}

// AddGt 新增大於條件
func (fb *FilterBuilder) AddGt(field string, value interface{}) *FilterBuilder {
	fb.conditions = append(fb.conditions, FilterCondition{
		Field:    field,
		Operator: "gt",
		Value:    value,
	})
	return fb
}

// AddGte 新增大於等於條件
func (fb *FilterBuilder) AddGte(field string, value interface{}) *FilterBuilder {
	fb.conditions = append(fb.conditions, FilterCondition{
		Field:    field,
		Operator: "gte",
		Value:    value,
	})
	return fb
}

// AddLt 新增小於條件
func (fb *FilterBuilder) AddLt(field string, value interface{}) *FilterBuilder {
	fb.conditions = append(fb.conditions, FilterCondition{
		Field:    field,
		Operator: "lt",
		Value:    value,
	})
	return fb
}

// AddLte 新增小於等於條件
func (fb *FilterBuilder) AddLte(field string, value interface{}) *FilterBuilder {
	fb.conditions = append(fb.conditions, FilterCondition{
		Field:    field,
		Operator: "lte",
		Value:    value,
	})
	return fb
}

// AddLike 新增 LIKE 條件
func (fb *FilterBuilder) AddLike(field string, value string) *FilterBuilder {
	fb.conditions = append(fb.conditions, FilterCondition{
		Field:    field,
		Operator: "like",
		Value:    "%" + value + "%",
	})
	return fb
}

// AddIn 新增 IN 條件
func (fb *FilterBuilder) AddIn(field string, values []interface{}) *FilterBuilder {
	fb.conditions = append(fb.conditions, FilterCondition{
		Field:    field,
		Operator: "in",
		Value:    values,
	})
	return fb
}

// AddNotIn 新增 NOT IN 條件
func (fb *FilterBuilder) AddNotIn(field string, values []interface{}) *FilterBuilder {
	fb.conditions = append(fb.conditions, FilterCondition{
		Field:    field,
		Operator: "not_in",
		Value:    values,
	})
	return fb
}

// AddBetween 新增 BETWEEN 條件
func (fb *FilterBuilder) AddBetween(field string, start, end interface{}) *FilterBuilder {
	fb.conditions = append(fb.conditions, FilterCondition{
		Field:    field,
		Operator: "between",
		Value:    []interface{}{start, end},
	})
	return fb
}

// AddDateBetween 新增日期範圍條件（自動處理時間）
func (fb *FilterBuilder) AddDateBetween(field string, startDate, endDate string) *FilterBuilder {
	fb.conditions = append(fb.conditions, FilterCondition{
		Field:    field,
		Operator: "date_between",
		Value:    []string{startDate, endDate},
	})
	return fb
}

// AddCenterScope 新增中心範圍條件
func (fb *FilterBuilder) AddCenterScope(centerID uint) *FilterBuilder {
	fb.conditions = append(fb.conditions, FilterCondition{
		Field:    "center_id",
		Operator: "eq",
		Value:    centerID,
	})
	return fb
}

// Build 建構條件切片
func (fb *FilterBuilder) Build() []FilterCondition {
	return fb.conditions
}

// Len 返回條件數量
func (fb *FilterBuilder) Len() int {
	return len(fb.conditions)
}

// IsEmpty 檢查是否為空
func (fb *FilterBuilder) IsEmpty() bool {
	return len(fb.conditions) == 0
}

// ServiceLogger 服務層日誌工具
type ServiceLogger struct {
	logger   *logger.Logger
	component string
	enabled   bool
}

// NewServiceLogger 建立服務層日誌工具
func NewServiceLogger(component string) *ServiceLogger {
	var log *logger.Logger
	enabled := false

	// 嘗試取得日誌器，如果未初始化則使用無操作日誌
	defer func() {
		if r := recover(); r != nil {
			// 日誌未初始化，使用無操作模式
			enabled = false
		}
	}()

	log = logger.GetLogger()
	enabled = log != nil

	return &ServiceLogger{
		logger:   log,
		component: component,
		enabled:   enabled,
	}
}

// Debug 記錄除錯訊息
func (sl *ServiceLogger) Debug(message string, keysAndValues ...interface{}) {
	if !sl.enabled || sl.logger == nil {
		return
	}
	sl.logger.ForComponent(sl.component).Debugw(message, keysAndValues...)
}

// Info 記錄資訊訊息
func (sl *ServiceLogger) Info(message string, keysAndValues ...interface{}) {
	if !sl.enabled || sl.logger == nil {
		return
	}
	sl.logger.ForComponent(sl.component).Infow(message, keysAndValues...)
}

// Warn 記錄警告訊息
func (sl *ServiceLogger) Warn(message string, keysAndValues ...interface{}) {
	if !sl.enabled || sl.logger == nil {
		return
	}
	sl.logger.ForComponent(sl.component).Warnw(message, keysAndValues...)
}

// Error 記錄錯誤訊息
func (sl *ServiceLogger) Error(message string, keysAndValues ...interface{}) {
	if !sl.enabled || sl.logger == nil {
		return
	}
	sl.logger.ForComponent(sl.component).Errorw(message, keysAndValues...)
}

// ErrorWithErr 記錄錯誤訊息（含錯誤物件）
func (sl *ServiceLogger) ErrorWithErr(message string, err error, keysAndValues ...interface{}) {
	if !sl.enabled || sl.logger == nil {
		return
	}
	sl.logger.ForComponent(sl.component).Errorw(message, append(keysAndValues, "error", err.Error())...)
}
