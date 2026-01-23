package services

import (
	"context"
	"time"
	"timeLedger/app/models"
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
	CheckOverlap(ctx context.Context, centerID uint, teacherID *uint, roomID uint, startTime, endTime time.Time, excludeRuleID *uint) (ValidationResult, error)

	// CheckTeacherBuffer 檢查老師轉場緩衝
	// 檢查老師連續課程是否滿足課程的轉場緩衝時間
	CheckTeacherBuffer(ctx context.Context, centerID uint, teacherID uint, prevEndTime, nextStartTime time.Time, courseID uint) (ValidationResult, error)

	// CheckRoomBuffer 檢查教室清潔緩衝
	// 檢查同一教室連續課程是否滿足課程的清潔緩衝時間
	CheckRoomBuffer(ctx context.Context, centerID uint, roomID uint, prevEndTime, nextStartTime time.Time, courseID uint) (ValidationResult, error)

	// ValidateFull 完整驗證
	// 執行所有檢查（重疊 + 緩衝）
	ValidateFull(ctx context.Context, centerID uint, teacherID *uint, roomID uint, courseID uint, startTime, endTime time.Time, excludeRuleID *uint, allowBufferOverride bool) (ValidationResult, error)
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
	// 支援停課(CANCEL)、改期(RESCHEDULE)、新增(ADD)
	CreateException(ctx context.Context, centerID uint, teacherID uint, ruleID uint, originalDate time.Time, exceptionType string, newStartAt, newEndAt *time.Time, newTeacherID *uint, reason string) (models.ScheduleException, error)

	// CheckExceptionDeadline 檢查是否超過異動截止日
	// 檢查規則的 lock_at 或中心的 exception_lead_days
	CheckExceptionDeadline(ctx context.Context, centerID uint, ruleID uint, exceptionDate time.Time) (allowed bool, reason string, err error)

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
}

type ExpandedSchedule struct {
	RuleID       uint      `json:"rule_id"`
	Date         time.Time `json:"date"`
	StartTime    string    `json:"start_time"`
	EndTime      string    `json:"end_time"`
	RoomID       uint      `json:"room_id"`
	TeacherID    *uint     `json:"teacher_id"`
	IsHoliday    bool      `json:"is_holiday"`
	HasException bool      `json:"has_exception"`
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
