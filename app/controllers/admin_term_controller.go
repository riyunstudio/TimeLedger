package controllers

import (
	"strconv"
	"time"
	"timeLedger/app"
	"timeLedger/app/resources"
	"timeLedger/app/services"

	"github.com/gin-gonic/gin"
)

// AdminTermController 學期管理控制器
type AdminTermController struct {
	BaseController
	app          *app.App
	termService  *services.TermService
	termResource *resources.TermResource
}

// NewAdminTermController 建立 AdminTermController 實例
func NewAdminTermController(app *app.App) *AdminTermController {
	return &AdminTermController{
		app:          app,
		termService:  services.NewTermService(app),
		termResource: resources.NewTermResource(app),
	}
}

// GetTerms 取得學期列表
// @Summary 取得學期列表
// @Tags Admin
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} global.ApiResponse{data=[]resources.TermResponse}
// @Router /api/v1/admin/terms [get]
func (ctl *AdminTermController) GetTerms(ctx *gin.Context) {
	helper := NewContextHelper(ctx)

	// 防止瀏覽器快取
	ctx.Header("Cache-Control", "no-cache, no-store, must-revalidate")
	ctx.Header("Pragma", "no-cache")
	ctx.Header("Expires", "0")

	centerID := helper.MustCenterID()
	if centerID == 0 {
		return
	}

	terms, errInfo, err := ctl.termService.GetTerms(ctx.Request.Context(), centerID)
	if err != nil {
		helper.ErrorWithInfo(errInfo)
		return
	}

	responses := ctl.termResource.ToTermResponses(terms)
	helper.Success(responses)
}

// GetActiveTerms 取得進行中的學期列表
// @Summary 取得進行中的學期列表
// @Tags Admin
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} global.ApiResponse{data=[]resources.TermResponse}
// @Router /api/v1/admin/terms/active [get]
func (ctl *AdminTermController) GetActiveTerms(ctx *gin.Context) {
	helper := NewContextHelper(ctx)

	// 防止瀏覽器快取
	ctx.Header("Cache-Control", "no-cache, no-store, must-revalidate")
	ctx.Header("Pragma", "no-cache")
	ctx.Header("Expires", "0")

	centerID := helper.MustCenterID()
	if centerID == 0 {
		return
	}

	terms, errInfo, err := ctl.termService.GetActiveTerms(ctx.Request.Context(), centerID)
	if err != nil {
		helper.ErrorWithInfo(errInfo)
		return
	}

	responses := ctl.termResource.ToTermResponses(terms)
	helper.Success(responses)
}

// CreateTerm 新增學期
// @Summary 新增學期
// @Tags Admin
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body services.CreateTermRequest true "學期資訊"
// @Success 200 {object} global.ApiResponse{data=resources.TermResponse}
// @Router /api/v1/admin/terms [post]
func (ctl *AdminTermController) CreateTerm(ctx *gin.Context) {
	helper := NewContextHelper(ctx)

	centerID := helper.MustCenterID()
	if centerID == 0 {
		return
	}

	var req services.CreateTermRequest
	if !helper.MustBindJSON(&req) {
		return
	}

	adminID := helper.MustUserID()
	if adminID == 0 {
		return
	}

	term, errInfo, err := ctl.termService.CreateTerm(ctx.Request.Context(), centerID, adminID, &req)
	if err != nil {
		helper.ErrorWithInfo(errInfo)
		return
	}

	response := ctl.termResource.ToTermResponse(*term)
	helper.Success(response)
}

// UpdateTerm 更新學期
// @Summary 更新學期
// @Tags Admin
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param term_id path int true "Term ID"
// @Param request body services.UpdateTermRequest true "學期資訊"
// @Success 200 {object} global.ApiResponse{data=resources.TermResponse}
// @Router /api/v1/admin/terms/{term_id} [put]
func (ctl *AdminTermController) UpdateTerm(ctx *gin.Context) {
	helper := NewContextHelper(ctx)

	centerID := helper.MustCenterID()
	if centerID == 0 {
		return
	}

	termID := helper.MustParamUint("term_id")
	if termID == 0 {
		return
	}

	var req services.UpdateTermRequest
	if !helper.MustBindJSON(&req) {
		return
	}

	adminID := helper.MustUserID()
	if adminID == 0 {
		return
	}

	term, errInfo, err := ctl.termService.UpdateTerm(ctx.Request.Context(), centerID, adminID, termID, &req)
	if err != nil {
		helper.ErrorWithInfo(errInfo)
		return
	}

	response := ctl.termResource.ToTermResponse(*term)
	helper.Success(response)
}

// DeleteTerm 刪除學期
// @Summary 刪除學期
// @Tags Admin
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param term_id path int true "Term ID"
// @Success 200 {object} global.ApiResponse
// @Router /api/v1/admin/terms/{term_id} [delete]
func (ctl *AdminTermController) DeleteTerm(ctx *gin.Context) {
	helper := NewContextHelper(ctx)

	centerID := helper.MustCenterID()
	if centerID == 0 {
		return
	}

	termID := helper.MustParamUint("term_id")
	if termID == 0 {
		return
	}

	adminID := helper.MustUserID()
	if adminID == 0 {
		return
	}

	errInfo, err := ctl.termService.DeleteTerm(ctx.Request.Context(), centerID, adminID, termID)
	if err != nil {
		helper.ErrorWithInfo(errInfo)
		return
	}

	helper.Success(nil)
}

// GetOccupancyRules 取得佔用規則（按老師或教室分組）
// @Summary 取得佔用規則
// @Tags Admin
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param teacher_id query int false "Teacher ID"
// @Param room_id query int false "Room ID"
// @Param start_date query string true "開始日期 (YYYY-MM-DD)"
// @Param end_date query string true "結束日期 (YYYY-MM-DD)"
// @Success 200 {object} global.ApiResponse{data=[]resources.OccupancyRulesByDayOfWeek}
// @Router /api/v1/admin/occupancy/rules [get]
func (ctl *AdminTermController) GetOccupancyRules(ctx *gin.Context) {
	helper := NewContextHelper(ctx)

	centerID := helper.MustCenterID()
	if centerID == 0 {
		return
	}

	// 解析查詢參數
	var teacherID, roomID *uint

	if teacherIDStr := ctx.Query("teacher_id"); teacherIDStr != "" {
		id, err := strconv.ParseUint(teacherIDStr, 10, 64)
		if err != nil {
			helper.BadRequest("Invalid teacher_id format")
			return
		}
		uid := uint(id)
		teacherID = &uid
	}

	if roomIDStr := ctx.Query("room_id"); roomIDStr != "" {
		id, err := strconv.ParseUint(roomIDStr, 10, 64)
		if err != nil {
			helper.BadRequest("Invalid room_id format")
			return
		}
		uid := uint(id)
		roomID = &uid
	}

	// 必須提供 teacher_id 或 room_id
	if teacherID == nil && roomID == nil {
		helper.BadRequest("Either teacher_id or room_id is required")
		return
	}

	// 解析日期範圍
	startDateStr := ctx.Query("start_date")
	endDateStr := ctx.Query("end_date")

	if startDateStr == "" || endDateStr == "" {
		helper.BadRequest("start_date and end_date are required")
		return
	}

	startDate, err := time.Parse("2006-01-02", startDateStr)
	if err != nil {
		helper.BadRequest("Invalid start_date format, expected YYYY-MM-DD")
		return
	}

	endDate, err := time.Parse("2006-01-02", endDateStr)
	if err != nil {
		helper.BadRequest("Invalid end_date format, expected YYYY-MM-DD")
		return
	}

	// 驗證日期範圍
	if endDate.Before(startDate) {
		helper.BadRequest("end_date must be after or equal to start_date")
		return
	}

	// 防止查詢過大範圍（最多 1 年）
	maxDuration := time.Hour * 24 * 365
	if endDate.Sub(startDate) > maxDuration {
		helper.BadRequest("Date range cannot exceed 1 year")
		return
	}

	groups, errInfo, err := ctl.termService.GetOccupancyRules(ctx.Request.Context(), centerID, teacherID, roomID, startDate, endDate)
	if err != nil {
		helper.ErrorWithInfo(errInfo)
		return
	}

	// 直接構造響應類型
	dayNames := []string{"", "週一", "週二", "週三", "週四", "週五", "週六", "週日"}
	responses := make([]resources.OccupancyRulesByDayOfWeek, len(groups))

	for i, group := range groups {
		rules := make([]resources.OccupancyRuleInfo, len(group.Rules))
		for j, rule := range group.Rules {
			rules[j] = resources.OccupancyRuleInfo{
				RuleID:       rule.RuleID,
				OfferingID:   rule.OfferingID,
				OfferingName: rule.OfferingName,
				Weekday:      rule.Weekday,
				StartTime:    rule.StartTime,
				EndTime:      rule.EndTime,
				Duration:     rule.Duration,
				TeacherID:    rule.TeacherID,
				TeacherName:  rule.TeacherName,
				RoomID:       rule.RoomID,
				RoomName:     rule.RoomName,
			}
		}

		dayName := ""
		if group.DayOfWeek >= 1 && group.DayOfWeek <= 7 {
			dayName = dayNames[group.DayOfWeek]
		}

		responses[i] = resources.OccupancyRulesByDayOfWeek{
			DayOfWeek: group.DayOfWeek,
			DayName:   dayName,
			Rules:     rules,
		}
	}

	helper.Success(responses)
}

// CopyRulesRequest 複製規則請求結構
type CopyRulesRequest struct {
	SourceTermID uint   `json:"source_term_id" binding:"required"`
	TargetTermID uint   `json:"target_term_id" binding:"required"`
	RuleIDs      []uint `json:"rule_ids" binding:"required,min=1"`
}

// CopyRules 批量複製規則到目標學期
// @Summary 批量複製規則
// @Tags Admin
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body CopyRulesRequest true "複製規則請求"
// @Success 200 {object} global.ApiResponse{data=resources.CopyRulesResponse}
// @Router /api/v1/admin/terms/copy-rules [post]
func (ctl *AdminTermController) CopyRules(ctx *gin.Context) {
	helper := NewContextHelper(ctx)

	centerID := helper.MustCenterID()
	if centerID == 0 {
		return
	}

	adminID := helper.MustUserID()
	if adminID == 0 {
		return
	}

	var req CopyRulesRequest
	if !helper.MustBindJSON(&req) {
		return
	}

	// 驗證規則數量（最多 100 條）
	if len(req.RuleIDs) > 100 {
		helper.BadRequest("Maximum 100 rules can be copied at once")
		return
	}

	// 轉換為服務層請求格式
	serviceReq := &services.CopyRulesRequest{
		SourceTermID: req.SourceTermID,
		TargetTermID: req.TargetTermID,
		RuleIDs:      req.RuleIDs,
	}

	result, errInfo, err := ctl.termService.CopyRules(ctx.Request.Context(), centerID, adminID, serviceReq)
	if err != nil {
		helper.ErrorWithInfo(errInfo)
		return
	}

	// 直接構造響應類型
	rules := make([]resources.CopiedRuleInfo, len(result.Rules))
	for i, rule := range result.Rules {
		rules[i] = resources.CopiedRuleInfo{
			OriginalRuleID: rule.OriginalRuleID,
			NewRuleID:     rule.NewRuleID,
			OfferingName:   rule.OfferingName,
			Weekday:       rule.Weekday,
			StartTime:     rule.StartTime,
			EndTime:       rule.EndTime,
		}
	}

	response := resources.CopyRulesResponse{
		CopiedCount: result.CopiedCount,
		Rules:       rules,
	}

	helper.Success(response)
}
