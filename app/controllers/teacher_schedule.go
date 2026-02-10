package controllers

import (
	"fmt"
	"time"
	"timeLedger/app"
	"timeLedger/app/models"
	"timeLedger/app/repositories"
	"timeLedger/app/services"

	"github.com/gin-gonic/gin"
)

// TeacherScheduleController 老師課表相關 API
type TeacherScheduleController struct {
	BaseController
	app               *app.App
	scheduleRuleRepo  *repositories.ScheduleRuleRepository
	scheduleQuery     services.ScheduleQueryService
	membershipRepo    *repositories.CenterMembershipRepository
	exceptionService  services.ScheduleExceptionService
	recurrenceService services.ScheduleRecurrenceService
	auditLogRepo      *repositories.AuditLogRepository
}

func NewTeacherScheduleController(app *app.App) *TeacherScheduleController {
	return &TeacherScheduleController{
		app:               app,
		scheduleRuleRepo:  repositories.NewScheduleRuleRepository(app),
		scheduleQuery:     services.NewScheduleQueryService(app),
		membershipRepo:    repositories.NewCenterMembershipRepository(app),
		exceptionService:  services.NewScheduleExceptionService(app),
		recurrenceService: services.NewScheduleRecurrenceService(app),
		auditLogRepo:      repositories.NewAuditLogRepository(app),
	}
}

// GetSchedule 取得老師的綜合課表（個人行程 + 各中心課程）
// @Summary 取得老師的綜合課表
// @Tags Teacher
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param from query string true "開始日期 (YYYY-MM-DD)"
// @Param to query string true "結束日期 (YYYY-MM-DD)"
// @Success 200 {object} global.ApiResponse{data=[]services.TeacherScheduleItem}
// @Router /api/v1/teacher/me/schedule [get]
func (ctl *TeacherScheduleController) GetSchedule(ctx *gin.Context) {
	helper := NewContextHelper(ctx)

	teacherID := helper.MustUserID()
	if teacherID == 0 {
		return
	}

	from, to := helper.MustQueryDateRange("from", "to")

	schedule, err := ctl.scheduleQuery.GetTeacherSchedule(ctx, teacherID, from, to)
	if err != nil {
		helper.InternalError("Failed to get schedule")
		return
	}

	helper.Success(schedule)
}

// GetCenterScheduleRules 獲取老師在指定中心的排課規則
// @Summary 獲取老師在指定中心的排課規則
// @Tags Teacher
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param center_id path int true "中心 ID"
// @Success 200 {object} global.ApiResponse{data=[]models.ScheduleRule}
// @Router /api/v1/teacher/me/centers/{center_id}/schedule-rules [get]
func (ctl *TeacherScheduleController) GetCenterScheduleRules(ctx *gin.Context) {
	helper := NewContextHelper(ctx)

	teacherID := helper.MustUserID()
	if teacherID == 0 {
		return
	}

	centerID := helper.MustParamUint("center_id")
	if centerID == 0 {
		return
	}

	// 驗證老師是否屬於該中心
	membership, err := ctl.membershipRepo.GetActiveByTeacherAndCenter(ctx, teacherID, fmt.Sprintf("%d", centerID))
	if err != nil || membership == nil {
		helper.Forbidden("You are not a member of this center")
		return
	}

	// 獲取該中心的排課規則
	rules, err := ctl.scheduleRuleRepo.ListByCenterID(ctx, centerID)
	if err != nil {
		helper.InternalError("Failed to get schedule rules")
		return
	}

	// 過濾出該老師的課程，並格式化輸出
	type RuleResponse struct {
		ID                 uint   `json:"id"`
		Title              string `json:"title"`
		Weekday            int    `json:"weekday"`
		WeekdayText        string `json:"weekday_text"`
		StartTime          string `json:"start_time"`
		EndTime            string `json:"end_time"`
		EffectiveStartDate string `json:"effective_start_date"`
		EffectiveEndDate   string `json:"effective_end_date"`
	}

	weekdayTexts := []string{"週日", "週一", "週二", "週三", "週四", "週五", "週六"}

	var teacherRules []RuleResponse
	for _, rule := range rules {
		if rule.TeacherID != nil && *rule.TeacherID == teacherID {
			title := rule.Offering.Name
			if title == "" {
				title = rule.Name
			}

			effectiveStartDate := ""
			effectiveEndDate := ""
			if !rule.EffectiveRange.StartDate.IsZero() {
				effectiveStartDate = rule.EffectiveRange.StartDate.Format("2006-01-02")
			}
			if !rule.EffectiveRange.EndDate.IsZero() {
				effectiveEndDate = rule.EffectiveRange.EndDate.Format("2006-01-02")
			}

			// 確保 weekday 在有效範圍內 (0-6)
			weekdayText := ""
			weekday := rule.Weekday
			if weekday == 7 {
				weekday = 0 // 將 7 視為週日
			}
			if weekday >= 0 && weekday < len(weekdayTexts) {
				weekdayText = weekdayTexts[weekday]
			}

			teacherRules = append(teacherRules, RuleResponse{
				ID:                 rule.ID,
				Title:              title,
				Weekday:            rule.Weekday,
				WeekdayText:        weekdayText,
				StartTime:          rule.StartTime,
				EndTime:            rule.EndTime,
				EffectiveStartDate: effectiveStartDate,
				EffectiveEndDate:   effectiveEndDate,
			})
		}
	}

	helper.Success(teacherRules)
}

// GetSchedules 取得老師的課表（支援 start_date/end_date 參數）
// @Summary 取得老師的課表
// @Tags Teacher
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param start_date query string true "開始日期 (YYYY-MM-DD)"
// @Param end_date query string true "結束日期 (YYYY-MM-DD)"
// @Success 200 {object} global.ApiResponse{data=[]services.TeacherScheduleItem}
// @Router /api/v1/teacher/schedules [get]
func (ctl *TeacherScheduleController) GetSchedules(ctx *gin.Context) {
	helper := NewContextHelper(ctx)

	teacherID := helper.MustUserID()
	if teacherID == 0 {
		return
	}

	// 支援兩種參數名稱
	fromStr := helper.QueryStringOrDefault("start_date", helper.QueryStringOrDefault("from", ""))
	toStr := helper.QueryStringOrDefault("end_date", helper.QueryStringOrDefault("to", ""))

	if fromStr == "" || toStr == "" {
		helper.BadRequest("start_date and end_date are required")
		return
	}

	fromDate, err := time.Parse("2006-01-02", fromStr)
	if err != nil {
		helper.BadRequest("Invalid start_date format")
		return
	}

	toDate, err := time.Parse("2006-01-02", toStr)
	if err != nil {
		helper.BadRequest("Invalid end_date format")
		return
	}

	// 使用 ScheduleQueryService 取得課表
	schedule, err := ctl.scheduleQuery.GetTeacherSchedule(ctx, teacherID, fromDate, toDate)
	if err != nil {
		helper.InternalError("Failed to get schedule")
		return
	}

	helper.Success(schedule)
}

// CheckTeacherRuleLockRequest 檢查老師規則鎖定狀態請求
type CheckTeacherRuleLockRequest struct {
	RuleID        uint   `json:"rule_id" binding:"required"`
	ExceptionDate string `json:"exception_date" binding:"required"`
}

// CheckTeacherRuleLockResponse 檢查老師規則鎖定狀態回應
type CheckTeacherRuleLockResponse struct {
	IsLocked      bool       `json:"is_locked"`
	LockReason    string     `json:"lock_reason,omitempty"`
	Deadline      *time.Time `json:"deadline,omitempty"`
	DaysRemaining int        `json:"days_remaining"`
}

// CheckRuleLockStatus 檢查老師規則鎖定狀態
// @Summary 檢查老師規則鎖定狀態
// @Tags Teacher
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body CheckTeacherRuleLockRequest true "檢查規則鎖定請求"
// @Success 200 {object} global.ApiResponse{data=CheckTeacherRuleLockResponse}
// @Router /api/v1/teacher/scheduling/check-rule-lock [post]
func (ctl *TeacherScheduleController) CheckRuleLockStatus(ctx *gin.Context) {
	helper := NewContextHelper(ctx)

	teacherID := helper.MustUserID()
	if teacherID == 0 {
		return
	}

	var req CheckTeacherRuleLockRequest
	if !helper.MustBindJSON(&req) {
		return
	}

	rule, err := ctl.scheduleRuleRepo.GetByID(ctx, req.RuleID)
	if err != nil {
		helper.NotFound("Rule not found")
		return
	}

	exceptionDate, err := time.Parse("2006-01-02", req.ExceptionDate)
	if err != nil {
		helper.BadRequest("Invalid date format, expected YYYY-MM-DD")
		return
	}

	allowed, errInfo, _ := ctl.exceptionService.CheckExceptionDeadline(ctx, rule.CenterID, req.RuleID, exceptionDate)

	response := CheckTeacherRuleLockResponse{
		IsLocked:   !allowed,
		LockReason: errInfo.Msg,
	}

	if !allowed {
		helper.Success(response)
		return
	}

	helper.Success(response)
}

// PreviewRecurrenceEditRequest 預覽循環編輯請求
type PreviewRecurrenceEditRequest struct {
	RuleID   uint   `json:"rule_id" binding:"required"`
	EditDate string `json:"edit_date" binding:"required"`
	Mode     string `json:"mode" binding:"required,oneof=SINGLE FUTURE ALL"`
}

// PreviewRecurrenceEdit 預覽循環編輯影響範圍
// @Summary 預覽循環編輯影響範圍
// @Tags Teacher
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body PreviewRecurrenceEditRequest true "預覽循環編輯請求"
// @Success 200 {object} global.ApiResponse
// @Router /api/v1/teacher/scheduling/preview-recurrence-edit [post]
func (ctl *TeacherScheduleController) PreviewRecurrenceEdit(ctx *gin.Context) {
	helper := NewContextHelper(ctx)

	teacherID := helper.MustUserID()
	if teacherID == 0 {
		return
	}

	var req PreviewRecurrenceEditRequest
	if !helper.MustBindJSON(&req) {
		return
	}

	editDate, err := time.Parse("2006-01-02", req.EditDate)
	if err != nil {
		helper.BadRequest("Invalid date format, expected YYYY-MM-DD")
		return
	}

	preview, err := ctl.recurrenceService.PreviewAffectedSessions(ctx, req.RuleID, editDate, services.RecurrenceEditMode(req.Mode))
	if err != nil {
		helper.InternalError("Failed to preview affected sessions")
		return
	}

	helper.Success(preview)
}

// EditRecurringSchedule 編輯循環排課
// @Summary 編輯循環排課
// @Tags Teacher
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body services.RecurrenceEditRequest true "編輯循環排課請求"
// @Success 200 {object} global.ApiResponse
// @Router /api/v1/teacher/scheduling/edit-recurring [post]
func (ctl *TeacherScheduleController) EditRecurringSchedule(ctx *gin.Context) {
	helper := NewContextHelper(ctx)

	teacherID := helper.MustUserID()
	if teacherID == 0 {
		return
	}

	var req services.RecurrenceEditRequest
	if !helper.MustBindJSON(&req) {
		return
	}

	rule, err := ctl.scheduleRuleRepo.GetByID(ctx, req.RuleID)
	if err != nil {
		helper.NotFound("Rule not found")
		return
	}

	result, err := ctl.recurrenceService.EditRecurringSchedule(ctx, rule.CenterID, teacherID, req)
	if err != nil {
		helper.InternalError(err.Error())
		return
	}

	helper.Success(result)
}

// DeleteRecurringScheduleRequest 刪除循環排課請求
type DeleteRecurringScheduleRequest struct {
	RuleID   uint   `json:"rule_id" binding:"required"`
	EditDate string `json:"edit_date" binding:"required"`
	Mode     string `json:"mode" binding:"required,oneof=SINGLE FUTURE ALL"`
	Reason   string `json:"reason"`
}

// DeleteRecurringSchedule 刪除循環排課
// @Summary 刪除循環排課
// @Tags Teacher
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body DeleteRecurringScheduleRequest true "刪除循環排課請求"
// @Success 200 {object} global.ApiResponse
// @Router /api/v1/teacher/scheduling/delete-recurring [post]
func (ctl *TeacherScheduleController) DeleteRecurringSchedule(ctx *gin.Context) {
	helper := NewContextHelper(ctx)

	teacherID := helper.MustUserID()
	if teacherID == 0 {
		return
	}

	var req DeleteRecurringScheduleRequest
	if !helper.MustBindJSON(&req) {
		return
	}

	editDate, err := time.Parse("2006-01-02", req.EditDate)
	if err != nil {
		helper.BadRequest("Invalid date format, expected YYYY-MM-DD")
		return
	}

	rule, err := ctl.scheduleRuleRepo.GetByID(ctx, req.RuleID)
	if err != nil {
		helper.NotFound("Rule not found")
		return
	}

	result, err := ctl.recurrenceService.DeleteRecurringSchedule(ctx, rule.CenterID, teacherID, req.RuleID, editDate, services.RecurrenceEditMode(req.Mode), req.Reason)
	if err != nil {
		helper.InternalError(err.Error())
		return
	}

	// 記錄審核日誌
	ctl.auditLogRepo.Create(ctx, models.AuditLog{
		CenterID:   rule.CenterID,
		ActorType:  "TEACHER",
		ActorID:    teacherID,
		Action:     "DELETE_RECURRING_SCHEDULE",
		TargetType: "ScheduleRule",
		TargetID:   req.RuleID,
		Payload: models.AuditPayload{
			After: map[string]interface{}{
				"mode":   req.Mode,
				"reason": req.Reason,
			},
		},
	})

	helper.Success(result)
}
