package resources

import (
	"time"
	"timeLedger/app"
	"timeLedger/app/models"
)

// ScheduleRuleResponse 排課規則響應
type ScheduleRuleResponse struct {
	ID            uint      `json:"id"`
	CenterID      uint      `json:"center_id"`
	OfferingID    uint      `json:"offering_id"`
	TeacherID     *uint     `json:"teacher_id,omitempty"`
	RoomID        uint      `json:"room_id"`
	Weekday       int       `json:"weekday"`
	StartTime     string    `json:"start_time"`
	EndTime       string    `json:"end_time"`
	Duration      int       `json:"duration"`
	EffectiveFrom string    `json:"effective_from"`
	EffectiveTo   string    `json:"effective_to"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

// ExpandedScheduleResponse 展開後的課表響應
type ExpandedScheduleResponse struct {
	RuleID         uint           `json:"rule_id"`
	Date           time.Time      `json:"date"`
	StartTime      string         `json:"start_time"`
	EndTime        string         `json:"end_time"`
	RoomID         uint           `json:"room_id"`
	TeacherID      *uint          `json:"teacher_id,omitempty"`
	IsHoliday      bool           `json:"is_holiday"`
	HasException   bool           `json:"has_exception"`
	ExceptionInfo  *ExceptionInfo `json:"exception_info,omitempty"`
	OfferingName   string         `json:"offering_name"`
	TeacherName    string         `json:"teacher_name"`
	RoomName       string         `json:"room_name"`
	OfferingID     uint           `json:"offering_id"`
	EffectiveRange *DateRange     `json:"effective_range,omitempty"`
	IsCrossDayPart bool           `json:"is_cross_day_part"`
}

// ExceptionInfo 例外資訊
type ExceptionInfo struct {
	ID           uint       `json:"id"`
	Type         string     `json:"type"`
	Status       string     `json:"status"`
	NewTeacherID *uint      `json:"new_teacher_id,omitempty"`
	NewStartAt   *time.Time `json:"new_start_at,omitempty"`
	NewEndAt     *time.Time `json:"new_end_at,omitempty"`
}

// DateRange 日期範圍
type DateRange struct {
	StartDate time.Time `json:"start_date"`
	EndDate   time.Time `json:"end_date"`
}

// ScheduleExceptionResponse 例外申請響應
type ScheduleExceptionResponse struct {
	ID            uint       `json:"id"`
	CenterID      uint       `json:"center_id"`
	RuleID        uint       `json:"rule_id"`
	OriginalDate  time.Time  `json:"original_date"`
	ExceptionType string     `json:"exception_type"`
	Status        string     `json:"status"`
	NewStartAt    *time.Time `json:"new_start_at,omitempty"`
	NewEndAt      *time.Time `json:"new_end_at,omitempty"`
	NewTeacherID  *uint      `json:"new_teacher_id,omitempty"`
	NewRoomID     *uint      `json:"new_room_id,omitempty"`
	Reason        string     `json:"reason"`
	ReviewNote    string     `json:"review_note,omitempty"`
	ReviewedBy    *uint      `json:"reviewed_by,omitempty"`
	ReviewedAt    *time.Time `json:"reviewed_at,omitempty"`
	CreatedAt     time.Time  `json:"created_at"`
	UpdatedAt     time.Time  `json:"updated_at"`
	// 關聯資料
	Rule *ScheduleRuleResponse `json:"rule,omitempty"`
}

// OverlapCheckResultResponse 衝突檢查結果響應
type OverlapCheckResultResponse struct {
	HasOverlap bool     `json:"has_overlap"`
	Conflicts  []string `json:"conflicts,omitempty"`
}

// BufferCheckResultResponse 緩衝檢查結果響應
type BufferCheckResultResponse struct {
	Valid     bool                     `json:"valid"`
	Conflicts []BufferConflictResponse `json:"conflicts,omitempty"`
}

// BufferConflictResponse 緩衝衝突響應
type BufferConflictResponse struct {
	Type            string `json:"type"`
	Message         string `json:"message"`
	RequiredMinutes int    `json:"required_minutes"`
	DiffMinutes     int    `json:"diff_minutes"`
	CanOverride     bool   `json:"can_override"`
	PreviousEndTime string `json:"previous_end_time,omitempty"`
	NextStartTime   string `json:"next_start_time,omitempty"`
}

// FullValidationResultResponse 完整驗證結果響應
type FullValidationResultResponse struct {
	Valid     bool                     `json:"valid"`
	Conflicts []BufferConflictResponse `json:"conflicts,omitempty"`
}

// RuleLockStatusResponse 規則鎖定狀態響應
type RuleLockStatusResponse struct {
	IsLocked      bool       `json:"is_locked"`
	LockReason    string     `json:"lock_reason,omitempty"`
	LockAt        *time.Time `json:"lock_at,omitempty"`
	Deadline      time.Time  `json:"deadline"`
	DaysRemaining int        `json:"days_remaining"`
}

// PhaseTransitionResponse 階段轉換響應
type PhaseTransitionResponse struct {
	Date          time.Time `json:"date"`
	PrevRuleID    *uint     `json:"prev_rule_id,omitempty"`
	PrevRoomID    *uint     `json:"prev_room_id,omitempty"`
	PrevTeacherID *uint     `json:"prev_teacher_id,omitempty"`
	PrevStartTime string    `json:"prev_start_time,omitempty"`
	PrevEndTime   string    `json:"prev_end_time,omitempty"`
	NextRuleID    *uint     `json:"next_rule_id,omitempty"`
	NextRoomID    *uint     `json:"next_room_id,omitempty"`
	NextTeacherID *uint     `json:"next_teacher_id,omitempty"`
	NextStartTime string    `json:"next_start_time,omitempty"`
	NextEndTime   string    `json:"next_end_time,omitempty"`
	HasGap        bool      `json:"has_gap"`
}

// TodaySummaryResponse 今日摘要響應
type TodaySummaryResponse struct {
	Sessions               []TodaySessionResponse `json:"sessions"`
	TotalSessions          int                    `json:"total_sessions"`
	CompletedSessions      int                    `json:"completed_sessions"`
	InProgressSessions     int                    `json:"in_progress_sessions"`
	UpcomingSessions       int                    `json:"upcoming_sessions"`
	InProgressTeacherNames []string               `json:"in_progress_teacher_names,omitempty"`
	PendingExceptions      int                    `json:"pending_exceptions"`
	ChangesCount           int                    `json:"changes_count"`
	HasScheduleChanges     bool                   `json:"has_schedule_changes"`
}

// TodaySessionResponse 今日課程響應
type TodaySessionResponse struct {
	ID        uint                  `json:"id"`
	StartTime time.Time             `json:"start_time"`
	EndTime   time.Time             `json:"end_time"`
	Offering  TodayOfferingResponse `json:"offering"`
	Teacher   TodayTeacherResponse  `json:"teacher"`
	Room      TodayRoomResponse     `json:"room"`
	Status    string                `json:"status"`
}

// TodayOfferingResponse 今日課程響應
type TodayOfferingResponse struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
}

// TodayTeacherResponse 今日老師響應
type TodayTeacherResponse struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
}

// TodayRoomResponse 今日教室響應
type TodayRoomResponse struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
}

// ScheduleResource 排課資源轉換
type ScheduleResource struct {
	app *app.App
}

// NewScheduleResource 建立 ScheduleResource
func NewScheduleResource(app *app.App) *ScheduleResource {
	return &ScheduleResource{app: app}
}

// ToRuleResponse 轉換規則為響應格式
func (r *ScheduleResource) ToRuleResponse(rule models.ScheduleRule) *ScheduleRuleResponse {
	return &ScheduleRuleResponse{
		ID:            rule.ID,
		CenterID:      rule.CenterID,
		OfferingID:    rule.OfferingID,
		TeacherID:     rule.TeacherID,
		RoomID:        rule.RoomID,
		Weekday:       rule.Weekday,
		StartTime:     rule.StartTime,
		EndTime:       rule.EndTime,
		Duration:      rule.Duration,
		EffectiveFrom: rule.EffectiveRange.StartDate.Format("2006-01-02"),
		EffectiveTo:   rule.EffectiveRange.EndDate.Format("2006-01-02"),
		CreatedAt:     rule.CreatedAt,
		UpdatedAt:     rule.UpdatedAt,
	}
}

// ToExceptionResponse 轉換例外為響應格式
func (r *ScheduleResource) ToExceptionResponse(exception models.ScheduleException) *ScheduleExceptionResponse {
	return &ScheduleExceptionResponse{
		ID:            exception.ID,
		CenterID:      exception.CenterID,
		RuleID:        exception.RuleID,
		OriginalDate:  exception.OriginalDate,
		ExceptionType: exception.ExceptionType,
		Status:        exception.Status,
		NewStartAt:    exception.NewStartAt,
		NewEndAt:      exception.NewEndAt,
		NewTeacherID:  exception.NewTeacherID,
		NewRoomID:     exception.NewRoomID,
		Reason:        exception.Reason,
		ReviewNote:    exception.ReviewNote,
		ReviewedBy:    exception.ReviewedBy,
		ReviewedAt:    exception.ReviewedAt,
		CreatedAt:     exception.CreatedAt,
		UpdatedAt:     exception.UpdatedAt,
	}
}

// =========================================
// Matrix View Response Types
// =========================================

// MatrixViewResponse 矩陣視圖響應
type MatrixViewResponse struct {
	TimeSlots []int            `json:"time_slots"` // 橫軸時段，如 [9, 10, 11, ...]
	Resources []MatrixResource `json:"resources"`  // 縱軸資源（老師或教室）
	DateRange MatrixDateRange  `json:"date_range"`
}

// MatrixDateRange 日期範圍
type MatrixDateRange struct {
	StartDate string `json:"start_date"` // YYYY-MM-DD
	EndDate   string `json:"end_date"`   // YYYY-MM-DD
}

// MatrixResource 矩陣資源（老師或教室）
type MatrixResource struct {
	ID    uint         `json:"id"`
	Name  string       `json:"name"`
	Type  string       `json:"type"` // "teacher" | "room"
	Items []MatrixItem `json:"items"`
}

// MatrixItem 矩陣項目（課程場次）
type MatrixItem struct {
	ID            uint    `json:"id"`
	RuleID        uint    `json:"rule_id"`
	Title         string  `json:"title"`          // 課程名稱
	Date          string  `json:"date"`           // YYYY-MM-DD
	StartTime     string  `json:"start_time"`     // HH:mm
	EndTime       string  `json:"end_time"`       // HH:mm
	StartHour     int     `json:"start_hour"`     // 開始小時 (用於 CSS 定位)
	StartMinute   int     `json:"start_minute"`   // 開始分鐘
	Duration      int     `json:"duration"`       // 持續分鐘數
	TopOffset     float64 `json:"top_offset"`     // CSS top 百分比 (0-100)
	HeightPercent float64 `json:"height_percent"` // CSS height 百分比 (相對於時段)
	OfferingID    uint    `json:"offering_id"`
	OfferingName  string  `json:"offering_name"`
	TeacherID     *uint   `json:"teacher_id,omitempty"`
	TeacherName   string  `json:"teacher_name"`
	RoomID        uint    `json:"room_id"`
	RoomName      string  `json:"room_name"`
	IsHoliday     bool    `json:"is_holiday"`
	HasException  bool    `json:"has_exception"`
	ExceptionType string  `json:"exception_type,omitempty"`
	IsSuspended   bool    `json:"is_suspended"`    // 是否為停課
	Color         string  `json:"color,omitempty"` // 課程顏色
}

// MatrixViewResource 矩陣視圖資源轉換
type MatrixViewResource struct {
	app *app.App
}

// NewMatrixViewResource 建立 MatrixViewResource
func NewMatrixViewResource(app *app.App) *MatrixViewResource {
	return &MatrixViewResource{app: app}
}
