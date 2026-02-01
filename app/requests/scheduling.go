package requests

import "time"

// CheckOverlapRequest 衝突檢查請求
type CheckOverlapRequest struct {
	TeacherID     *uint   `json:"teacher_id"`
	RoomID        uint    `json:"room_id" binding:"required"`
	StartTime     string  `json:"start_time" binding:"required"`
	EndTime       string  `json:"end_time" binding:"required"`
	Weekday       int     `json:"weekday"` // 可選，如果未提供則從 StartTime 推算
	ExcludeRuleID *uint   `json:"exclude_rule_id"`
}

// CheckBufferRequest 緩衝時間檢查請求
type CheckBufferRequest struct {
	TeacherID     uint   `json:"teacher_id" binding:"required"`
	RoomID        uint   `json:"room_id" binding:"required"`
	PrevEndTime   string `json:"prev_end_time" binding:"required"`
	NextStartTime string `json:"next_start_time" binding:"required"`
	CourseID      uint   `json:"course_id" binding:"required"`
}

// CreateExceptionRequest 建立例外請求
type CreateExceptionRequest struct {
	RuleID         uint       `json:"rule_id" binding:"required"`
	OriginalDate   time.Time  `json:"original_date" binding:"required"`
	Type           string     `json:"type" binding:"required"`
	NewStartAt     *time.Time `json:"new_start_at"`
	NewEndAt       *time.Time `json:"new_end_at"`
	NewTeacherID   *uint      `json:"new_teacher_id"`
	NewTeacherName string     `json:"new_teacher_name"`
	NewRoomID      *uint      `json:"new_room_id"`
	Reason         string     `json:"reason" binding:"required"`
}

// ReviewExceptionRequest 審核例外請求
type ReviewExceptionRequest struct {
	Action         string `json:"action" binding:"required"`
	OverrideBuffer bool   `json:"override_buffer"`
	Reason         string `json:"reason"`
}

// ExpandRulesRequest 展開規則請求
type ExpandRulesRequest struct {
	RuleIDs   []uint   `json:"rule_ids"`
	StartDate string   `json:"start_date" binding:"required,date_format"`
	EndDate   string   `json:"end_date" binding:"required,date_format"`
}

// ValidateFullRequest 完整驗證請求
type ValidateFullRequest struct {
	TeacherID           *uint    `json:"teacher_id"`
	RoomID              uint     `json:"room_id" binding:"required"`
	CourseID            uint     `json:"course_id" binding:"required"`
	StartTime           string   `json:"start_time" binding:"required,time_format"`
	EndTime             string   `json:"end_time" binding:"required,time_format"`
	ExcludeRuleID       *uint    `json:"exclude_rule_id"`
	AllowBufferOverride bool     `json:"allow_buffer_override"`
	// 以下欄位可選，如果未提供，系統會自動計算上一堂課的結束時間
	PrevEndTime   *string `json:"prev_end_time"`
	NextStartTime *string `json:"next_start_time"`
}

// CreateRuleRequest 建立排課規則請求
type CreateRuleRequest struct {
	Name           string  `json:"name" binding:"required"`
	OfferingID     uint    `json:"offering_id" binding:"required"`
	TeacherID      uint    `json:"teacher_id"`
	RoomID         uint    `json:"room_id" binding:"required"`
	StartTime      string  `json:"start_time" binding:"required,time_format"`
	EndTime        string  `json:"end_time" binding:"required,time_format"`
	Duration       int     `json:"duration" binding:"required"`
	Weekdays       []int   `json:"weekdays" binding:"required,min=1"`
	StartDate      string  `json:"start_date" binding:"required,date_format"`
	EndDate        *string `json:"end_date"`
	OverrideBuffer bool    `json:"override_buffer"`
}

// UpdateRuleRequest 更新排課規則請求
type UpdateRuleRequest struct {
	Name       string  `json:"name"`
	OfferingID uint    `json:"offering_id"`
	TeacherID  *uint   `json:"teacher_id"`
	RoomID     uint    `json:"room_id"`
	StartTime  string  `json:"start_time"`
	EndTime    string  `json:"end_time"`
	Duration   int     `json:"duration"`
	Weekdays   []int   `json:"weekdays"`
	StartDate  string  `json:"start_date"`
	EndDate    *string `json:"end_date"`
	// 更新模式：SINGLE - 只修改這一天，FUTURE - 修改這天及之後，ALL - 修改所有
	UpdateMode string `json:"update_mode"`
}

// DetectPhaseTransitionsRequest 偵測階段轉換請求
type DetectPhaseTransitionsRequest struct {
	OfferingID uint      `json:"offering_id" binding:"required"`
	StartDate  string    `json:"start_date" binding:"required,date_format"`
	EndDate    string    `json:"end_date" binding:"required,date_format"`
}

// CheckRuleLockStatusRequest 檢查規則鎖定狀態請求
type CheckRuleLockStatusRequest struct {
	RuleID        uint      `json:"rule_id" binding:"required"`
	ExceptionDate time.Time `json:"exception_date" binding:"required"`
}

// CreateScheduleRuleFromOfferingRequest 從開課建立規則請求
type CreateScheduleRuleFromOfferingRequest struct {
	OfferingID  uint     `json:"offering_id" binding:"required"`
	RoomID      uint     `json:"room_id" binding:"required"`
	TeacherID   *uint    `json:"teacher_id"`
	StartTime   string   `json:"start_time" binding:"required,time_format"`
	EndTime     string   `json:"end_time" binding:"required,time_format"`
	Weekdays    []int    `json:"weekdays" binding:"required,min=1"`
	StartDate   string   `json:"start_date" binding:"required,date_format"`
	EndDate     *string  `json:"end_date"`
	Duration    int      `json:"duration" binding:"required"`
	OverrideBuffer bool `json:"override_buffer"`
}

// BatchCreateRulesRequest 批量建立規則請求
type BatchCreateRulesRequest struct {
	Rules []CreateRuleRequest `json:"rules" binding:"required,min=1,dive"`
}

// UpdateRulesBatchRequest 批量更新規則請求
type UpdateRulesBatchRequest struct {
	Updates []UpdateRuleRequest `json:"updates" binding:"required,min=1,dive"`
	DeleteIDs []uint           `json:"delete_ids"`
}
