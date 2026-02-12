package services

import (
	"context"
	"encoding/json"
	"time"
	"timeLedger/app/models"
	"timeLedger/global/errInfos"
)

type ValidationResult struct {
	Valid     bool                 `json:"valid"`
	Conflicts []ValidationConflict `json:"conflicts"`
}

type ValidationConflict struct {
	Type             string `json:"type"` // TEACHER_OVERLAP, ROOM_OVERLAP, TEACHER_BUFFER, ROOM_BUFFER
	Message          string `json:"message"`
	CanOverride      bool   `json:"can_override"`
	RequireApproval  bool   `json:"require_approval,omitempty"`
	RequiredMinutes  int    `json:"required_minutes,omitempty"`
	DiffMinutes      int    `json:"diff_minutes,omitempty"`
	ConflictSource   string `json:"conflict_source,omitempty"` // RULE, SESSION, PERSONAL
	ConflictSourceID uint   `json:"conflict_source_id,omitempty"`
	Details          string `json:"details,omitempty"`
}

type ScheduleValidationService interface {
	// CheckOverlap 檢查時段重疊
	// 檢查指定時段是否與現有排課衝突
	CheckOverlap(ctx context.Context, centerID uint, teacherID *uint, roomID uint, startTime, endTime time.Time, weekday int, excludeRuleID *uint) (ValidationResult, error)

	// CheckTeacherBuffer 檢查老師轉場緩衝
	// 檢查老師連續課程是否滿足課程的轉場緩衝時間
	CheckTeacherBuffer(ctx context.Context, centerID uint, teacherID uint, prevEndTime, nextStartTime time.Time, courseID uint) (ValidationResult, error)

	// CheckRoomBuffer 檢查教室清潔緩衝
	// 檢查同一教室連續課程是否滿足課程的清潔緩衝時間
	CheckRoomBuffer(ctx context.Context, centerID uint, roomID uint, prevEndTime, nextStartTime time.Time, courseID uint) (ValidationResult, error)

	// ValidateFull 完整驗證
	// 執行所有檢查（重疊 + 緩衝）
	// 如果 prevEndTime 和 nextStartTime 為 nil，系統會自動計算上一堂課的結束時間
	ValidateFull(ctx context.Context, centerID uint, teacherID *uint, roomID uint, courseID uint, startTime, endTime time.Time, excludeRuleID *uint, allowBufferOverride bool, prevEndTime, nextStartTime *time.Time) (ValidationResult, error)
}

type ScheduleExpansionService interface {
	// ExpandRules 將週規則展開為實際日期
	// 根據規則的週期（weekday）和有效範圍，展開為具體日期
	ExpandRules(ctx context.Context, rules []models.ScheduleRule, startDate, endDate time.Time, centerID uint) []ExpandedSchedule

	// GetEffectiveRuleForDate 取得指定日期的有效規則
	// 用於檢測 phase transition
	GetEffectiveRuleForDate(ctx context.Context, offeringID uint, date time.Time) (*models.ScheduleRule, error)

	// DetectPhaseTransitions 檢測指定範圍內的 phase 變化
	// 返回 phase 變化的日期點
	DetectPhaseTransitions(ctx context.Context, centerID uint, offeringID uint, startDate, endDate time.Time) ([]PhaseTransition, error)

	// GetRulesByEffectiveDateRange 取得指定 effective 日期範圍內的規則
	GetRulesByEffectiveDateRange(ctx context.Context, centerID uint, offeringID uint, startDate, endDate time.Time) ([]models.ScheduleRule, error)
}

type ScheduleExceptionService interface {
	// CreateException 創建例外單
	// 支援停課(CANCEL)、改期(RESCHEDULE)、代課(REPLACE_TEACHER)
	CreateException(ctx context.Context, centerID uint, teacherID uint, ruleID uint, req *CreateExceptionRequest) (models.ScheduleException, *errInfos.Res, error)

	// CheckExceptionDeadline 檢查是否超過異動截止日
	// 檢查規則的 lock_at 或中心的 exception_lead_days
	CheckExceptionDeadline(ctx context.Context, centerID uint, ruleID uint, exceptionDate time.Time) (allowed bool, errInfo *errInfos.Res, err error)

	// RevokeException 老師撤回待審核的例外單
	RevokeException(ctx context.Context, exceptionID uint, teacherID uint) error

	// ReviewException 審核例外單
	// 只有 ADMIN 角色可以審核，核准時會執行 Re-validation
	ReviewException(ctx context.Context, exceptionID uint, adminID uint, action string, overrideBuffer bool, reason string) error

	// GetExceptionsByRule 取得某規則的所有例外
	GetExceptionsByRule(ctx context.Context, ruleID uint) ([]models.ScheduleException, error)

	// GetExceptionsByDateRange 取得日期範圍內的例外
	GetExceptionsByDateRange(ctx context.Context, centerID uint, startDate, endDate time.Time) ([]models.ScheduleException, error)

	// GetPendingExceptions 取得待審核的例外單
	GetPendingExceptions(ctx context.Context, centerID uint) ([]models.ScheduleException, error)

	// GetAllExceptions 取得所有例外單（可依狀態篩選）
	// status 為空時返回所有狀態的例外
	GetAllExceptions(ctx context.Context, centerID uint, status string) ([]models.ScheduleException, error)
}

type ExpandedSchedule struct {
	RuleID        uint               `json:"rule_id"`
	Date          time.Time          `json:"date"`
	StartTime     string             `json:"start_time"`
	EndTime       string             `json:"end_time"`
	RoomID        uint               `json:"room_id"`
	TeacherID     *uint              `json:"teacher_id"`
	IsHoliday     bool               `json:"is_holiday"`
	HasException  bool               `json:"has_exception"`
	Status        string             `json:"status"` // 狀態: PLANNED(預計), CONFIRMED(已開課), SUSPENDED(停課), ARCHIVED(歸檔)
	ExceptionInfo *ExpandedException `json:"exception_info,omitempty"`
	// 關聯資料
	OfferingName   string            `json:"offering_name,omitempty"`
	TeacherName    string            `json:"teacher_name,omitempty"`
	RoomName       string            `json:"room_name,omitempty"`
	OfferingID     uint              `json:"offering_id,omitempty"`
	EffectiveRange *models.DateRange `json:"effective_range,omitempty"`
	IsCrossDayPart bool              `json:"is_cross_day_part,omitempty"` // 跨日課程的一部分
}

// MarshalJSON 確保 Date 欄位以 ISO 8601 日期格式序列化
func (es ExpandedSchedule) MarshalJSON() ([]byte, error) {
	type Alias ExpandedSchedule
	return json.Marshal(struct {
		Alias
		Date string `json:"date"`
	}{
		Alias: Alias(es),
		Date:  es.Date.Format("2006-01-02"),
	})
}

type ExpandedException struct {
	ID           uint       `json:"id"`
	Type         string     `json:"type"`   // CANCEL, RESCHEDULE, REPLACE_TEACHER
	Status       string     `json:"status"` // PENDING, APPROVED, REJECTED, REVOKED
	NewTeacherID *uint      `json:"new_teacher_id,omitempty"`
	NewStartAt   *time.Time `json:"new_start_at,omitempty"`
	NewEndAt     *time.Time `json:"new_end_at,omitempty"`
}

type PhaseTransition struct {
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

// MarshalJSON 確保 Date 欄位以 ISO 8601 日期格式序列化
func (pt PhaseTransition) MarshalJSON() ([]byte, error) {
	type Alias PhaseTransition
	return json.Marshal(struct {
		Alias
		Date string `json:"date"`
	}{
		Alias: Alias(pt),
		Date:  pt.Date.Format("2006-01-02"),
	})
}

type RecurrenceEditMode string

const (
	RecurrenceEditSingle RecurrenceEditMode = "SINGLE"
	RecurrenceEditFuture RecurrenceEditMode = "FUTURE"
	RecurrenceEditAll    RecurrenceEditMode = "ALL"
)

type RecurrenceEditRequest struct {
	RuleID       uint               `json:"rule_id" binding:"required"`
	EditDate     time.Time          `json:"edit_date" binding:"required"`
	Mode         RecurrenceEditMode `json:"mode" binding:"required,oneof=SINGLE FUTURE ALL"`
	NewStartTime string             `json:"new_start_time,omitempty"`
	NewEndTime   string             `json:"new_end_time,omitempty"`
	NewRoomID    *uint              `json:"new_room_id,omitempty"`
	NewTeacherID *uint              `json:"new_teacher_id,omitempty"`
	Reason       string             `json:"reason"`
}

type RecurrenceEditPreview struct {
	Mode           RecurrenceEditMode `json:"mode"`
	AffectedCount  int                `json:"affected_count"`
	AffectedDates  []time.Time        `json:"affected_dates"`
	WillCreateRule bool               `json:"will_create_rule,omitempty"`
	NewRuleID      *uint              `json:"new_rule_id,omitempty"`
}

type RecurrenceEditResult struct {
	Mode             RecurrenceEditMode         `json:"mode"`
	CancelExceptions []models.ScheduleException `json:"cancel_exceptions"`
	AddExceptions    []models.ScheduleException `json:"add_exceptions,omitempty"`
	UpdatedRule      *models.ScheduleRule       `json:"updated_rule,omitempty"`
	NewRule          *models.ScheduleRule       `json:"new_rule,omitempty"`
	AffectedCount    int                        `json:"affected_count"`
}

type ScheduleRecurrenceService interface {
	// PreviewAffectedSessions 預覽將被影響的場次
	// 根據編輯模式和日期，返回受影響的場次列表
	PreviewAffectedSessions(ctx context.Context, ruleID uint, editDate time.Time, mode RecurrenceEditMode) (RecurrenceEditPreview, error)

	// EditRecurringSchedule 編輯循環排課
	// 支援 SINGLE/FUTURE/ALL 三種模式
	// SINGLE: 產生 CANCEL + ADD 例外單
	// FUTURE: 產生 CANCEL + ADD 例外單，並創建新規則
	// ALL: 修改規則的有效範圍或基本設定
	EditRecurringSchedule(ctx context.Context, centerID uint, teacherID uint, req RecurrenceEditRequest) (RecurrenceEditResult, error)

	// DeleteRecurringSchedule 刪除循環排課
	// SINGLE: 產生 CANCEL 例外單
	// FUTURE: 產生 CANCEL 例外單，並創建新規則（空規則或調整有效範圍）
	// ALL: 軟刪除規則
	DeleteRecurringSchedule(ctx context.Context, centerID uint, teacherID uint, ruleID uint, editDate time.Time, mode RecurrenceEditMode, reason string) (RecurrenceEditResult, error)
}

// TeacherScheduleItem 老師課表項目
type TeacherScheduleItem struct {
	ID             string `json:"id"`
	Type           string `json:"type"`
	Title          string `json:"title"`
	Date           string `json:"date"`
	StartTime      string `json:"start_time"`
	EndTime        string `json:"end_time"`
	RoomID         uint   `json:"room_id"`
	TeacherID      *uint  `json:"teacher_id"`
	CenterID       uint   `json:"center_id"`
	CenterName     string `json:"center_name"`
	Status         string `json:"status"` // 狀態: PLANNED(預計), CONFIRMED(已開課), SUSPENDED(停課), ARCHIVED(歸檔)
	RuleID         uint   `json:"rule_id"`
	IsCrossDayPart bool   `json:"is_cross_day_part,omitempty"`
}

type ScheduleQueryService interface {
	// GetTeacherSchedule 取得老師的綜合課表
	GetTeacherSchedule(ctx context.Context, teacherID uint, fromDate, toDate time.Time) ([]TeacherScheduleItem, error)
}
