package controllers

import (
	"timeLedger/app"
	"timeLedger/app/requests"
	"timeLedger/app/services"

	"github.com/gin-gonic/gin"
)

// SchedulingController 排課管理控制器（Thin Controller）
type SchedulingController struct {
	BaseController
	app         *app.App
	scheduleSvc services.ScheduleServiceInterface
}

// NewSchedulingController 建立排課控制器
func NewSchedulingController(app *app.App) *SchedulingController {
	return &SchedulingController{
		app:         app,
		scheduleSvc: services.NewScheduleService(app),
	}
}

// requireCenterID 取得並驗證中心 ID（通用模式）
func (ctl *SchedulingController) requireCenterID(helper *ContextHelper) uint {
	centerID := helper.MustCenterID()
	if centerID == 0 {
		return 0
	}
	return centerID
}

// requireRuleID 取得並驗證規則 ID（通用模式）
func (ctl *SchedulingController) requireRuleID(helper *ContextHelper) uint {
	ruleID := helper.MustParamUint("ruleId")
	if ruleID == 0 {
		return 0
	}
	return ruleID
}

// requireExceptionID 取得並驗證例外 ID（通用模式）
func (ctl *SchedulingController) requireExceptionID(helper *ContextHelper) uint {
	exceptionID := helper.MustParamUint("exceptionId")
	if exceptionID == 0 {
		return 0
	}
	return exceptionID
}

// requireAdminID 取得並驗證管理員 ID（通用模式）
func (ctl *SchedulingController) requireAdminID(helper *ContextHelper) uint {
	adminID := helper.MustUserID()
	if adminID == 0 {
		return 0
	}
	return adminID
}

// CheckOverlap 檢查時間衝突
// @Summary 檢查課程時間是否與現有排程衝突
// @Tags Admin - Scheduling
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body requests.CheckOverlapRequest true "衝突檢查請求"
// @Success 200 {object} global.ApiResponse{data=services.OverlapCheckResult}
// @Router /api/v1/admin/scheduling/check-overlap [post]
func (ctl *SchedulingController) CheckOverlap(ctx *gin.Context) {
	helper := NewContextHelper(ctx)

	centerID := ctl.requireCenterID(helper)
	if centerID == 0 {
		return
	}

	var req requests.CheckOverlapRequest
	if !helper.MustBindJSON(&req) {
		return
	}

	// 計算 weekday（如果未提供則從 StartTime 推算）
	checkWeekday := req.Weekday
	if checkWeekday == 0 {
		checkWeekday = int(req.StartTime.Weekday())
		if checkWeekday == 0 {
			checkWeekday = 7
		}
	}

	result, err := ctl.scheduleSvc.CheckOverlap(ctx.Request.Context(), centerID, req.TeacherID, req.RoomID, req.StartTime, req.EndTime, checkWeekday, req.ExcludeRuleID)
	if err != nil {
		helper.InternalError(err.Error())
		return
	}

	helper.Success(result)
}

// CheckTeacherBuffer 檢查老師緩衝時間
// @Summary 檢查老師的緩衝時間是否足夠
// @Tags Admin - Scheduling
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body requests.CheckBufferRequest true "緩衝檢查請求"
// @Success 200 {object} global.ApiResponse{data=services.BufferCheckResult}
// @Router /api/v1/admin/scheduling/check-teacher-buffer [post]
func (ctl *SchedulingController) CheckTeacherBuffer(ctx *gin.Context) {
	helper := NewContextHelper(ctx)

	centerID := ctl.requireCenterID(helper)
	if centerID == 0 {
		return
	}

	var req requests.CheckBufferRequest
	if !helper.MustBindJSON(&req) {
		return
	}

	result, err := ctl.scheduleSvc.CheckTeacherBuffer(ctx.Request.Context(), centerID, req.TeacherID, req.PrevEndTime, req.NextStartTime, req.CourseID)
	if err != nil {
		helper.InternalError(err.Error())
		return
	}

	helper.Success(result)
}

// CheckRoomBuffer 檢查教室緩衝時間
// @Summary 檢查教室的緩衝時間是否足夠
// @Tags Admin - Scheduling
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body requests.CheckBufferRequest true "緩衝檢查請求"
// @Success 200 {object} global.ApiResponse{data=services.BufferCheckResult}
// @Router /api/v1/admin/scheduling/check-room-buffer [post]
func (ctl *SchedulingController) CheckRoomBuffer(ctx *gin.Context) {
	helper := NewContextHelper(ctx)

	centerID := ctl.requireCenterID(helper)
	if centerID == 0 {
		return
	}

	var req requests.CheckBufferRequest
	if !helper.MustBindJSON(&req) {
		return
	}

	result, err := ctl.scheduleSvc.CheckRoomBuffer(ctx.Request.Context(), centerID, req.RoomID, req.PrevEndTime, req.NextStartTime, req.CourseID)
	if err != nil {
		helper.InternalError(err.Error())
		return
	}

	helper.Success(result)
}

// ValidateFull 完整驗證排課
// @Summary 完整驗證排課（硬衝突 + 緩衝檢查）
// @Tags Admin - Scheduling
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body requests.ValidateFullRequest true "完整驗證請求"
// @Success 200 {object} global.ApiResponse{data=services.FullValidationResult}
// @Router /api/v1/admin/scheduling/validate [post]
func (ctl *SchedulingController) ValidateFull(ctx *gin.Context) {
	helper := NewContextHelper(ctx)

	centerID := ctl.requireCenterID(helper)
	if centerID == 0 {
		return
	}

	adminID := ctl.requireAdminID(helper)
	if adminID == 0 {
		return
	}

	var req requests.ValidateFullRequest
	if !helper.MustBindJSON(&req) {
		return
	}

	result, err := ctl.scheduleSvc.ValidateFull(ctx.Request.Context(), centerID, req.TeacherID, req.RoomID, req.CourseID, req.StartTime, req.EndTime, req.ExcludeRuleID, req.AllowBufferOverride)
	if err != nil {
		helper.InternalError(err.Error())
		return
	}

	helper.Success(result)
}

// GetRules 取得排課規則列表
// @Summary 取得中心的所有排課規則
// @Tags Admin - Scheduling
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} global.ApiResponse{data=[]models.ScheduleRule}
// @Router /api/v1/admin/scheduling/rules [get]
func (ctl *SchedulingController) GetRules(ctx *gin.Context) {
	helper := NewContextHelper(ctx)

	centerID := ctl.requireCenterID(helper)
	if centerID == 0 {
		return
	}

	rules, err := ctl.scheduleSvc.GetRules(ctx.Request.Context(), centerID)
	if err != nil {
		helper.InternalError(err.Error())
		return
	}

	helper.Success(rules)
}

// CreateRule 建立排課規則
// @Summary 建立新的排課規則
// @Tags Admin - Scheduling
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body requests.CreateRuleRequest true "規則資訊"
// @Success 200 {object} global.ApiResponse{data=[]models.ScheduleRule}
// @Router /api/v1/admin/scheduling/rules [post]
func (ctl *SchedulingController) CreateRule(ctx *gin.Context) {
	helper := NewContextHelper(ctx)

	centerID := ctl.requireCenterID(helper)
	if centerID == 0 {
		return
	}

	adminID := ctl.requireAdminID(helper)
	if adminID == 0 {
		return
	}

	var req requests.CreateRuleRequest
	if !helper.MustBindJSON(&req) {
		return
	}

	svcReq := &services.CreateScheduleRuleRequest{
		Name:           req.Name,
		OfferingID:     req.OfferingID,
		TeacherID:      req.TeacherID,
		RoomID:         req.RoomID,
		StartTime:      req.StartTime,
		EndTime:        req.EndTime,
		Duration:       req.Duration,
		Weekdays:       req.Weekdays,
		StartDate:      req.StartDate,
		EndDate:        req.EndDate,
		OverrideBuffer: req.OverrideBuffer,
	}

	rules, errInfo, err := ctl.scheduleSvc.CreateRule(ctx.Request.Context(), centerID, adminID, svcReq)
	if err != nil {
		if errInfo != nil {
			helper.ErrorWithInfo(errInfo)
		} else {
			helper.InternalError(err.Error())
		}
		return
	}

	helper.Success(rules)
}

// UpdateRule 更新排課規則
// @Summary 更新排課規則
// @Tags Admin - Scheduling
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path uint true "規則ID"
// @Param request body requests.UpdateRuleRequest true "規則資訊"
// @Success 200 {object} global.ApiResponse{data=[]models.ScheduleRule}
// @Router /api/v1/admin/scheduling/rules/{id} [put]
func (ctl *SchedulingController) UpdateRule(ctx *gin.Context) {
	helper := NewContextHelper(ctx)

	centerID := ctl.requireCenterID(helper)
	if centerID == 0 {
		return
	}

	adminID := ctl.requireAdminID(helper)
	if adminID == 0 {
		return
	}

	ruleID := ctl.requireRuleID(helper)
	if ruleID == 0 {
		return
	}

	var req requests.UpdateRuleRequest
	if !helper.MustBindJSON(&req) {
		return
	}

	svcReq := &services.UpdateScheduleRuleRequest{
		Name:       req.Name,
		OfferingID: req.OfferingID,
		TeacherID:  req.TeacherID,
		RoomID:     req.RoomID,
		StartTime:  req.StartTime,
		EndTime:    req.EndTime,
		Duration:   req.Duration,
		Weekdays:   req.Weekdays,
		StartDate:  req.StartDate,
		EndDate:    req.EndDate,
		UpdateMode: req.UpdateMode,
	}

	rules, errInfo, err := ctl.scheduleSvc.UpdateRule(ctx.Request.Context(), centerID, adminID, ruleID, svcReq)
	if err != nil {
		if errInfo != nil {
			helper.ErrorWithInfo(errInfo)
		} else {
			helper.InternalError(err.Error())
		}
		return
	}

	helper.Success(rules)
}

// DeleteRule 刪除排課規則
// @Summary 刪除排課規則
// @Tags Admin - Scheduling
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path uint true "規則ID"
// @Success 200 {object} global.ApiResponse
// @Router /api/v1/admin/scheduling/rules/{id} [delete]
func (ctl *SchedulingController) DeleteRule(ctx *gin.Context) {
	helper := NewContextHelper(ctx)

	centerID := ctl.requireCenterID(helper)
	if centerID == 0 {
		return
	}

	adminID := ctl.requireAdminID(helper)
	if adminID == 0 {
		return
	}

	ruleID := ctl.requireRuleID(helper)
	if ruleID == 0 {
		return
	}

	err := ctl.scheduleSvc.DeleteRule(ctx.Request.Context(), centerID, adminID, ruleID)
	if err != nil {
		helper.InternalError(err.Error())
		return
	}

	helper.Success(gin.H{"message": "Schedule rule deleted successfully"})
}

// ExpandRules 展開排課規則
// @Summary 展開排課規則為具體課程場次
// @Tags Admin - Scheduling
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body requests.ExpandRulesRequest true "展開請求"
// @Success 200 {object} global.ApiResponse{data=[]services.ExpandedSchedule}
// @Router /api/v1/admin/scheduling/expand [post]
func (ctl *SchedulingController) ExpandRules(ctx *gin.Context) {
	helper := NewContextHelper(ctx)

	centerID := ctl.requireCenterID(helper)
	if centerID == 0 {
		return
	}

	var req requests.ExpandRulesRequest
	if !helper.MustBindJSON(&req) {
		return
	}

	svcReq := &services.ExpandRulesRequest{
		RuleIDs:   req.RuleIDs,
		StartDate: req.StartDate,
		EndDate:   req.EndDate,
	}

	schedules, err := ctl.scheduleSvc.ExpandRules(ctx.Request.Context(), centerID, svcReq)
	if err != nil {
		helper.InternalError(err.Error())
		return
	}

	helper.Success(schedules)
}

// GetExceptionsByRule 取得規則的所有例外申請
// @Summary 取得指定規則的所有例外申請
// @Tags Admin - Scheduling
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path uint true "規則ID"
// @Success 200 {object} global.ApiResponse{data=[]models.ScheduleException}
// @Router /api/v1/admin/scheduling/rules/{id}/exceptions [get]
func (ctl *SchedulingController) GetExceptionsByRule(ctx *gin.Context) {
	helper := NewContextHelper(ctx)

	ruleID := ctl.requireRuleID(helper)
	if ruleID == 0 {
		return
	}

	exceptions, err := ctl.scheduleSvc.GetExceptionsByRule(ctx.Request.Context(), ruleID)
	if err != nil {
		helper.InternalError(err.Error())
		return
	}

	helper.Success(exceptions)
}

// CreateException 建立例外申請
// @Summary 建立新的例外申請
// @Tags Teacher - Scheduling
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body requests.CreateExceptionRequest true "例外資訊"
// @Success 200 {object} global.ApiResponse{data=models.ScheduleException}
// @Router /api/v1/teacher/me/exceptions [post]
func (ctl *SchedulingController) CreateException(ctx *gin.Context) {
	helper := NewContextHelper(ctx)

	centerID := ctl.requireCenterID(helper)
	if centerID == 0 {
		return
	}

	teacherID := ctl.requireAdminID(helper) // 老師端使用 MustUserID
	if teacherID == 0 {
		return
	}

	var req requests.CreateExceptionRequest
	if !helper.MustBindJSON(&req) {
		return
	}

	svcReq := &services.CreateExceptionRequest{
		RuleID:         req.RuleID,
		OriginalDate:   req.OriginalDate,
		Type:           req.Type,
		NewStartAt:     req.NewStartAt,
		NewEndAt:       req.NewEndAt,
		NewTeacherID:   req.NewTeacherID,
		NewTeacherName: req.NewTeacherName,
		NewRoomID:      req.NewRoomID,
		Reason:         req.Reason,
	}

	exception, err := ctl.scheduleSvc.CreateException(ctx.Request.Context(), centerID, teacherID, req.RuleID, svcReq)
	if err != nil {
		helper.InternalError(err.Error())
		return
	}

	helper.Success(exception)
}

// ReviewException 審核例外申請
// @Summary 審核例外申請（核准/拒絕）
// @Tags Admin - Scheduling
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path uint true "例外ID"
// @Param request body requests.ReviewExceptionRequest true "審核資訊"
// @Success 200 {object} global.ApiResponse
// @Router /api/v1/admin/scheduling/exceptions/{id}/review [post]
func (ctl *SchedulingController) ReviewException(ctx *gin.Context) {
	helper := NewContextHelper(ctx)

	adminID := ctl.requireAdminID(helper)
	if adminID == 0 {
		return
	}

	exceptionID := ctl.requireExceptionID(helper)
	if exceptionID == 0 {
		return
	}

	var req requests.ReviewExceptionRequest
	if !helper.MustBindJSON(&req) {
		return
	}

	svcReq := &services.ReviewExceptionRequest{
		Action:         req.Action,
		OverrideBuffer: req.OverrideBuffer,
		Reason:         req.Reason,
	}

	err := ctl.scheduleSvc.ReviewException(ctx.Request.Context(), exceptionID, adminID, svcReq)
	if err != nil {
		helper.InternalError(err.Error())
		return
	}

	helper.Success(gin.H{"message": "Exception reviewed successfully"})
}

// GetExceptionsByDateRange 取得日期範圍內的例外申請
// @Summary 取得指定日期範圍內的所有例外申請
// @Tags Admin - Scheduling
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param start_date query string false "開始日期 (YYYY-MM-DD)"
// @Param end_date query string false "結束日期 (YYYY-MM-DD)"
// @Success 200 {object} global.ApiResponse{data=[]models.ScheduleException}
// @Router /api/v1/admin/scheduling/exceptions [get]
func (ctl *SchedulingController) GetExceptionsByDateRange(ctx *gin.Context) {
	helper := NewContextHelper(ctx)

	centerID := ctl.requireCenterID(helper)
	if centerID == 0 {
		return
	}

	startDate, endDate := helper.MustQueryDateRange("start_date", "end_date")
	if startDate.IsZero() && endDate.IsZero() {
		return
	}

	exceptions, err := ctl.scheduleSvc.GetExceptionsByDateRange(ctx.Request.Context(), centerID, startDate, endDate)
	if err != nil {
		helper.InternalError(err.Error())
		return
	}

	helper.Success(exceptions)
}

// GetPendingExceptions 取得待審核的例外申請
// @Summary 取得所有待審核的例外申請
// @Tags Admin - Scheduling
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} global.ApiResponse{data=[]models.ScheduleException}
// @Router /api/v1/admin/scheduling/exceptions/pending [get]
func (ctl *SchedulingController) GetPendingExceptions(ctx *gin.Context) {
	helper := NewContextHelper(ctx)

	centerID := ctl.requireCenterID(helper)
	if centerID == 0 {
		return
	}

	exceptions, err := ctl.scheduleSvc.GetPendingExceptions(ctx.Request.Context(), centerID)
	if err != nil {
		helper.InternalError(err.Error())
		return
	}

	helper.Success(exceptions)
}

// GetAllExceptions 取得所有例外申請
// @Summary 取得所有例外申請（可依狀態篩選）
// @Tags Admin - Scheduling
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param status query string false "狀態篩選：PENDING, APPROVED, REJECTED, REVOKED"
// @Success 200 {object} global.ApiResponse{data=[]models.ScheduleException}
// @Router /api/v1/admin/scheduling/exceptions/all [get]
func (ctl *SchedulingController) GetAllExceptions(ctx *gin.Context) {
	helper := NewContextHelper(ctx)

	centerID := ctl.requireCenterID(helper)
	if centerID == 0 {
		return
	}

	status := helper.QueryStringOrDefault("status", "")

	exceptions, err := ctl.scheduleSvc.GetAllExceptions(ctx.Request.Context(), centerID, status)
	if err != nil {
		helper.InternalError(err.Error())
		return
	}

	helper.Success(exceptions)
}

// DetectPhaseTransitions 偵測階段轉換
// @Summary 偵測課程序列中的階段轉換點
// @Tags Admin - Scheduling
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body requests.DetectPhaseTransitionsRequest true "偵測請求"
// @Success 200 {object} global.ApiResponse{data=[]services.PhaseTransition}
// @Router /api/v1/admin/scheduling/phase-transitions [post]
func (ctl *SchedulingController) DetectPhaseTransitions(ctx *gin.Context) {
	helper := NewContextHelper(ctx)

	centerID := ctl.requireCenterID(helper)
	if centerID == 0 {
		return
	}

	var req requests.DetectPhaseTransitionsRequest
	if !helper.MustBindJSON(&req) {
		return
	}

	transitions, err := ctl.scheduleSvc.DetectPhaseTransitions(ctx.Request.Context(), centerID, req.OfferingID, req.StartDate, req.EndDate)
	if err != nil {
		helper.InternalError(err.Error())
		return
	}

	helper.Success(transitions)
}

// CheckRuleLockStatus 檢查規則鎖定狀態
// @Summary 檢查規則是否已超過異動截止日
// @Tags Admin - Scheduling
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body requests.CheckRuleLockStatusRequest true "檢查請求"
// @Success 200 {object} global.ApiResponse{data=services.RuleLockStatus}
// @Router /api/v1/admin/scheduling/rules/check-lock [post]
func (ctl *SchedulingController) CheckRuleLockStatus(ctx *gin.Context) {
	helper := NewContextHelper(ctx)

	centerID := ctl.requireCenterID(helper)
	if centerID == 0 {
		return
	}

	var req requests.CheckRuleLockStatusRequest
	if !helper.MustBindJSON(&req) {
		return
	}

	status, err := ctl.scheduleSvc.CheckRuleLockStatus(ctx.Request.Context(), centerID, req.RuleID, req.ExceptionDate)
	if err != nil {
		helper.InternalError(err.Error())
		return
	}

	helper.Success(status)
}

// GetTodaySummary 取得今日課表摘要
// @Summary 取得管理員後台首頁的今日課表摘要
// @Tags Admin - Dashboard
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} global.ApiResponse{data=services.TodaySummary}
// @Router /api/v1/admin/dashboard/today-summary [get]
func (ctl *SchedulingController) GetTodaySummary(ctx *gin.Context) {
	helper := NewContextHelper(ctx)

	centerID := ctl.requireCenterID(helper)
	if centerID == 0 {
		return
	}

	summary, err := ctl.scheduleSvc.GetTodaySummary(ctx.Request.Context(), centerID)
	if err != nil {
		helper.InternalError(err.Error())
		return
	}

	helper.Success(summary)
}
