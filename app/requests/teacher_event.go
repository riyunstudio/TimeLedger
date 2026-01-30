package requests

import (
	"time"

	"timeLedger/app/models"
)

// CreatePersonalEventRequest 新增個人行程請求
type CreatePersonalEventRequest struct {
	Title          string                 `json:"title" binding:"required"`
	StartAt        time.Time              `json:"start_at" binding:"required"`
	EndAt          time.Time              `json:"end_at" binding:"required"`
	IsAllDay       bool                   `json:"is_all_day"`
	ColorHex       string                 `json:"color_hex"`
	RecurrenceRule *models.RecurrenceRule `json:"recurrence_rule"`
}

// UpdatePersonalEventRequest 更新個人行程請求
type UpdatePersonalEventRequest struct {
	Title          *string                `json:"title"`
	StartAt        *time.Time             `json:"start_at"`
	EndAt          *time.Time             `json:"end_at"`
	IsAllDay       *bool                  `json:"is_all_day"`
	ColorHex       *string                `json:"color_hex"`
	RecurrenceRule *models.RecurrenceRule `json:"recurrence_rule"`
	UpdateMode     string                 `json:"update_mode" binding:"required,oneof=SINGLE FUTURE ALL"`
}

// UpdatePersonalEventResponse 更新回應結構
type UpdatePersonalEventResponse struct {
	UpdatedCount int64  `json:"updated_count"`
	Message      string `json:"message"`
}

// UpdatePersonalEventNoteRequest 更新備註請求
type UpdatePersonalEventNoteRequest struct {
	Content string `json:"content" binding:"required"`
}

// PersonalEventNoteResponse 個人行程備註回應
type PersonalEventNoteResponse struct {
	Content string `json:"content"`
}
