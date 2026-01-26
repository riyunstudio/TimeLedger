package controllers

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"
	"timeLedger/app"
	"timeLedger/app/models"
	"timeLedger/app/repositories"
	"timeLedger/app/services"
	"timeLedger/global"
	"timeLedger/global/errInfos"

	"github.com/gin-gonic/gin"
)

type SchedulingController struct {
	BaseController
	app               *app.App
	validationService services.ScheduleValidationService
	expansionService  services.ScheduleExpansionService
	exceptionService  services.ScheduleExceptionService
	auditLogRepo      *repositories.AuditLogRepository
	centerRepo        *repositories.CenterRepository
}

func NewSchedulingController(app *app.App) *SchedulingController {
	return &SchedulingController{
		app:               app,
		validationService: services.NewScheduleValidationService(app),
		expansionService:  services.NewScheduleExpansionService(app),
		exceptionService:  services.NewScheduleExceptionService(app),
		auditLogRepo:      repositories.NewAuditLogRepository(app),
		centerRepo:        repositories.NewCenterRepository(app),
	}
}

type CheckOverlapRequest struct {
	TeacherID     *uint     `json:"teacher_id"`
	RoomID        uint      `json:"room_id" binding:"required"`
	StartTime     time.Time `json:"start_time" binding:"required"`
	EndTime       time.Time `json:"end_time" binding:"required"`
	ExcludeRuleID *uint     `json:"exclude_rule_id"`
}

type CheckBufferRequest struct {
	TeacherID     uint      `json:"teacher_id" binding:"required"`
	RoomID        uint      `json:"room_id" binding:"required"`
	PrevEndTime   time.Time `json:"prev_end_time" binding:"required"`
	NextStartTime time.Time `json:"next_start_time" binding:"required"`
	CourseID      uint      `json:"course_id" binding:"required"`
}

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

type ReviewExceptionRequest struct {
	Action         string `json:"action" binding:"required"`
	OverrideBuffer bool   `json:"override_buffer"`
	Reason         string `json:"reason"`
}

type ExpandRulesRequest struct {
	RuleIDs   []uint    `json:"rule_ids"` // 可為空，空陣列表示展開所有規則
	StartDate time.Time `json:"start_date" binding:"required"`
	EndDate   time.Time `json:"end_date" binding:"required"`
}

type ValidateFullRequest struct {
	TeacherID           *uint     `json:"teacher_id"`
	RoomID              uint      `json:"room_id" binding:"required"`
	CourseID            uint      `json:"course_id" binding:"required"`
	StartTime           time.Time `json:"start_time" binding:"required"`
	EndTime             time.Time `json:"end_time" binding:"required"`
	ExcludeRuleID       *uint     `json:"exclude_rule_id"`
	AllowBufferOverride bool      `json:"allow_buffer_override"`
}

func (ctl *SchedulingController) CheckOverlap(ctx *gin.Context) {
	var req CheckOverlapRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, global.ApiResponse{
			Code:    global.BAD_REQUEST,
			Message: "Invalid request parameters",
		})
		return
	}

	// 從 JWT token 取得 center_id
	centerID := ctx.GetUint(global.CenterIDKey)
	if centerID == 0 {
		if val, exists := ctx.Get(global.CenterIDKey); exists {
			switch v := val.(type) {
			case uint:
				centerID = v
			case uint64:
				centerID = uint(v)
			case int:
				centerID = uint(v)
			case float64:
				centerID = uint(v)
			}
		}
	}

	if centerID == 0 {
		ctx.JSON(http.StatusUnauthorized, global.ApiResponse{
			Code:    global.UNAUTHORIZED,
			Message: "Center ID not found in token",
		})
		return
	}

	result, err := ctl.validationService.CheckOverlap(ctx, centerID, req.TeacherID, req.RoomID, req.StartTime, req.EndTime, req.ExcludeRuleID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, global.ApiResponse{
			Code:    500,
			Message: err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, global.ApiResponse{
		Code:    0,
		Message: "OK",
		Datas:   result,
	})
}

func (ctl *SchedulingController) CheckTeacherBuffer(ctx *gin.Context) {
	var req CheckBufferRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, global.ApiResponse{
			Code:    global.BAD_REQUEST,
			Message: "Invalid request parameters",
		})
		return
	}

	// 從 JWT token 取得 center_id
	centerID := ctx.GetUint(global.CenterIDKey)
	if centerID == 0 {
		if val, exists := ctx.Get(global.CenterIDKey); exists {
			switch v := val.(type) {
			case uint:
				centerID = v
			case uint64:
				centerID = uint(v)
			case int:
				centerID = uint(v)
			case float64:
				centerID = uint(v)
			}
		}
	}

	if centerID == 0 {
		ctx.JSON(http.StatusUnauthorized, global.ApiResponse{
			Code:    global.UNAUTHORIZED,
			Message: "Center ID not found in token",
		})
		return
	}

	result, err := ctl.validationService.CheckTeacherBuffer(ctx, centerID, req.TeacherID, req.PrevEndTime, req.NextStartTime, req.CourseID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, global.ApiResponse{
			Code:    500,
			Message: err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, global.ApiResponse{
		Code:    0,
		Message: "OK",
		Datas:   result,
	})
}

func (ctl *SchedulingController) CheckRoomBuffer(ctx *gin.Context) {
	var req CheckBufferRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, global.ApiResponse{
			Code:    global.BAD_REQUEST,
			Message: "Invalid request parameters",
		})
		return
	}

	// 從 JWT token 取得 center_id
	centerID := ctx.GetUint(global.CenterIDKey)
	if centerID == 0 {
		if val, exists := ctx.Get(global.CenterIDKey); exists {
			switch v := val.(type) {
			case uint:
				centerID = v
			case uint64:
				centerID = uint(v)
			case int:
				centerID = uint(v)
			case float64:
				centerID = uint(v)
			}
		}
	}

	if centerID == 0 {
		ctx.JSON(http.StatusUnauthorized, global.ApiResponse{
			Code:    global.UNAUTHORIZED,
			Message: "Center ID not found in token",
		})
		return
	}

	result, err := ctl.validationService.CheckRoomBuffer(ctx, centerID, req.RoomID, req.PrevEndTime, req.NextStartTime, req.CourseID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, global.ApiResponse{
			Code:    500,
			Message: err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, global.ApiResponse{
		Code:    0,
		Message: "OK",
		Datas:   result,
	})
}

func (ctl *SchedulingController) ValidateFull(ctx *gin.Context) {
	var req ValidateFullRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, global.ApiResponse{
			Code:    global.BAD_REQUEST,
			Message: "Invalid request parameters",
		})
		return
	}

	// 從 JWT token 取得 center_id
	centerID := ctx.GetUint(global.CenterIDKey)
	if centerID == 0 {
		if val, exists := ctx.Get(global.CenterIDKey); exists {
			switch v := val.(type) {
			case uint:
				centerID = v
			case uint64:
				centerID = uint(v)
			case int:
				centerID = uint(v)
			case float64:
				centerID = uint(v)
			}
		}
	}

	if centerID == 0 {
		ctx.JSON(http.StatusUnauthorized, global.ApiResponse{
			Code:    global.UNAUTHORIZED,
			Message: "Center ID not found in token",
		})
		return
	}

	adminID := ctx.GetUint(global.UserIDKey)
	result, err := ctl.validationService.ValidateFull(ctx, centerID, req.TeacherID, req.RoomID, req.CourseID, req.StartTime, req.EndTime, req.ExcludeRuleID, req.AllowBufferOverride)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, global.ApiResponse{
			Code:    500,
			Message: err.Error(),
		})
		return
	}

	conflictCount := 0
	if !result.Valid {
		conflictCount = len(result.Conflicts)
	}

	ctl.auditLogRepo.Create(ctx, models.AuditLog{
		CenterID:   centerID,
		ActorType:  "ADMIN",
		ActorID:    adminID,
		Action:     "VALIDATE_SCHEDULE",
		TargetType: "ScheduleValidation",
		TargetID:   0,
		Payload: models.AuditPayload{
			After: map[string]interface{}{
				"teacher_id":     req.TeacherID,
				"room_id":        req.RoomID,
				"course_id":      req.CourseID,
				"start_time":     req.StartTime,
				"end_time":       req.EndTime,
				"valid":          result.Valid,
				"conflict_count": conflictCount,
			},
		},
	})

	ctx.JSON(http.StatusOK, global.ApiResponse{
		Code:    0,
		Message: "OK",
		Datas:   result,
	})
}

func (ctl *SchedulingController) CreateException(ctx *gin.Context) {
	var req CreateExceptionRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, global.ApiResponse{
			Code:    global.BAD_REQUEST,
			Message: "Invalid request parameters",
		})
		return
	}

	// 從 JWT token 取得 center_id
	centerID := ctx.GetUint(global.CenterIDKey)
	if centerID == 0 {
		if val, exists := ctx.Get(global.CenterIDKey); exists {
			switch v := val.(type) {
			case uint:
				centerID = v
			case uint64:
				centerID = uint(v)
			case int:
				centerID = uint(v)
			case float64:
				centerID = uint(v)
			}
		}
	}

	if centerID == 0 {
		ctx.JSON(http.StatusUnauthorized, global.ApiResponse{
			Code:    global.UNAUTHORIZED,
			Message: "Center ID not found in token",
		})
		return
	}

	adminID := ctx.GetUint(global.UserIDKey)
	exception, err := ctl.exceptionService.CreateException(ctx, centerID, adminID, req.RuleID, req.OriginalDate, req.Type, req.NewStartAt, req.NewEndAt, req.NewTeacherID, req.NewTeacherName, req.Reason)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, global.ApiResponse{
			Code:    500,
			Message: err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, global.ApiResponse{
		Code:    0,
		Message: "Exception created successfully",
		Datas:   exception,
	})
}

func (ctl *SchedulingController) ReviewException(ctx *gin.Context) {
	id := ctx.Param("exceptionId")
	var req ReviewExceptionRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, global.ApiResponse{
			Code:    global.BAD_REQUEST,
			Message: "Invalid request parameters",
		})
		return
	}

	var exceptionID uint
	if _, err := fmt.Sscanf(id, "%d", &exceptionID); err != nil {
		ctx.JSON(http.StatusBadRequest, global.ApiResponse{
			Code:    global.BAD_REQUEST,
			Message: "Invalid exception ID",
		})
		return
	}

	adminID := ctx.GetUint(global.UserIDKey)
	err := ctl.exceptionService.ReviewException(ctx, exceptionID, adminID, req.Action, req.OverrideBuffer, req.Reason)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, global.ApiResponse{
			Code:    500,
			Message: err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, global.ApiResponse{
		Code:    0,
		Message: "Exception reviewed successfully",
	})
}

func (ctl *SchedulingController) GetExceptionsByRule(ctx *gin.Context) {
	ruleID := ctx.Param("ruleId")
	var ruleIDInt uint
	if _, err := fmt.Sscanf(ruleID, "%d", &ruleIDInt); err != nil {
		ctx.JSON(http.StatusBadRequest, global.ApiResponse{
			Code:    global.BAD_REQUEST,
			Message: "Invalid rule ID",
		})
		return
	}

	exceptions, err := ctl.exceptionService.GetExceptionsByRule(ctx, ruleIDInt)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, global.ApiResponse{
			Code:    500,
			Message: err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, global.ApiResponse{
		Code:    0,
		Message: "OK",
		Datas:   exceptions,
	})
}

func (ctl *SchedulingController) GetExceptionsByDateRange(ctx *gin.Context) {
	centerID := ctx.GetUint(global.CenterIDKey)
	if centerID == 0 {
		ctx.JSON(http.StatusBadRequest, global.ApiResponse{
			Code:    global.BAD_REQUEST,
			Message: "Center ID required",
		})
		return
	}

	startDateStr := ctx.Query("start_date")
	endDateStr := ctx.Query("end_date")

	startDate, err := time.Parse("2006-01-02", startDateStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, global.ApiResponse{
			Code:    global.BAD_REQUEST,
			Message: "Invalid start date format",
		})
		return
	}

	endDate, err := time.Parse("2006-01-02", endDateStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, global.ApiResponse{
			Code:    global.BAD_REQUEST,
			Message: "Invalid end date format",
		})
		return
	}

	exceptions, err := ctl.exceptionService.GetExceptionsByDateRange(ctx, centerID, startDate, endDate)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, global.ApiResponse{
			Code:    500,
			Message: err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, global.ApiResponse{
		Code:    0,
		Message: "OK",
		Datas:   exceptions,
	})
}

// GetPendingExceptions 獲取中心所有待審核的例外申請
// @Summary 獲取中心所有待審核的例外申請
// @Tags Admin - Scheduling
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} global.ApiResponse{data=[]models.ScheduleException}
// @Router /api/v1/admin/exceptions/pending [get]
func (ctl *SchedulingController) GetPendingExceptions(ctx *gin.Context) {
	centerID := ctx.GetUint(global.CenterIDKey)
	if centerID == 0 {
		ctx.JSON(http.StatusBadRequest, global.ApiResponse{
			Code:    global.BAD_REQUEST,
			Message: "Center ID required",
		})
		return
	}

	exceptions, err := ctl.exceptionService.GetPendingExceptions(ctx, centerID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, global.ApiResponse{
			Code:    500,
			Message: err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, global.ApiResponse{
		Code:    0,
		Message: "OK",
		Datas:   exceptions,
	})
}

// GetAllExceptions 取得所有例外申請（可依狀態篩選）
// @Summary 取得所有例外申請
// @Tags Admin - Scheduling
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param status query string false "狀態篩選：PENDING, APPROVED, REJECTED, REVOKED"
// @Success 200 {object} global.ApiResponse{data=[]models.ScheduleException}
// @Router /api/v1/admin/exceptions/all [get]
func (ctl *SchedulingController) GetAllExceptions(ctx *gin.Context) {
	centerID := ctx.GetUint(global.CenterIDKey)
	if centerID == 0 {
		ctx.JSON(http.StatusBadRequest, global.ApiResponse{
			Code:    global.BAD_REQUEST,
			Message: "Center ID required",
		})
		return
	}

	status := ctx.Query("status")

	exceptions, err := ctl.exceptionService.GetAllExceptions(ctx, centerID, status)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, global.ApiResponse{
			Code:    500,
			Message: err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, global.ApiResponse{
		Code:    0,
		Message: "OK",
		Datas:   exceptions,
	})
}

func (ctl *SchedulingController) ExpandRules(ctx *gin.Context) {
	centerID := ctx.GetUint(global.CenterIDKey)
	if centerID == 0 {
		ctx.JSON(http.StatusBadRequest, global.ApiResponse{
			Code:    global.BAD_REQUEST,
			Message: "Center ID required",
		})
		return
	}

	// 手動解析請求參數，支援 YYYY-MM-DD 格式
	var reqBody struct {
		RuleIDs   []uint  `json:"rule_ids"`
		StartDate string  `json:"start_date"`
		EndDate   string  `json:"end_date"`
	}
	if err := ctx.ShouldBindJSON(&reqBody); err != nil {
		ctx.JSON(http.StatusBadRequest, global.ApiResponse{
			Code:    global.BAD_REQUEST,
			Message: "Invalid request parameters",
		})
		return
	}

	// 解析日期（支援 YYYY-MM-DD 格式）
	startDate, err := time.Parse("2006-01-02", reqBody.StartDate)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, global.ApiResponse{
			Code:    global.BAD_REQUEST,
			Message: "Invalid start_date format, expected YYYY-MM-DD",
		})
		return
	}

	endDate, err := time.Parse("2006-01-02", reqBody.EndDate)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, global.ApiResponse{
			Code:    global.BAD_REQUEST,
			Message: "Invalid end_date format, expected YYYY-MM-DD",
		})
		return
	}

	scheduleRuleRepo := repositories.NewScheduleRuleRepository(ctl.app)
	rules, err := scheduleRuleRepo.ListByCenterID(ctx, centerID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, global.ApiResponse{
			Code:    500,
			Message: err.Error(),
		})
		return
	}

	var filteredRules []models.ScheduleRule
	if len(reqBody.RuleIDs) > 0 {
		for _, rule := range rules {
			for _, ruleID := range reqBody.RuleIDs {
				if rule.ID == ruleID {
					filteredRules = append(filteredRules, rule)
					break
				}
			}
		}
	} else {
		filteredRules = rules
	}

	expandedSchedules := ctl.expansionService.ExpandRules(ctx, filteredRules, startDate, endDate, centerID)

	ctx.JSON(http.StatusOK, global.ApiResponse{
		Code:    0,
		Message: "OK",
		Datas:   expandedSchedules,
	})
}

type CreateRuleRequest struct {
	Name       string  `json:"name" binding:"required"`
	OfferingID uint    `json:"offering_id" binding:"required"`
	TeacherID  uint    `json:"teacher_id"`
	RoomID     uint    `json:"room_id" binding:"required"`
	StartTime  string  `json:"start_time" binding:"required"`
	EndTime    string  `json:"end_time" binding:"required"`
	Duration   int     `json:"duration" binding:"required"`
	Weekdays   []int   `json:"weekdays" binding:"required,min=1"`
	StartDate  string  `json:"start_date" binding:"required"`
	EndDate    *string `json:"end_date"`
}

func (ctl *SchedulingController) CreateRule(ctx *gin.Context) {
	centerID := ctx.GetUint(global.CenterIDKey)
	if centerID == 0 {
		ctx.JSON(http.StatusBadRequest, global.ApiResponse{
			Code:    global.BAD_REQUEST,
			Message: "Center ID required",
		})
		return
	}

	var req CreateRuleRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, global.ApiResponse{
			Code:    global.BAD_REQUEST,
			Message: "Invalid request parameters: " + err.Error(),
		})
		return
	}

	startDate, err := time.Parse("2006-01-02", req.StartDate)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, global.ApiResponse{
			Code:    global.BAD_REQUEST,
			Message: "Invalid start_date format",
		})
		return
	}

	var endDate time.Time
	if req.EndDate != nil && *req.EndDate != "" {
		endDate, err = time.Parse("2006-01-02", *req.EndDate)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, global.ApiResponse{
				Code:    global.BAD_REQUEST,
				Message: "Invalid end_date format",
			})
			return
		}
	} else {
		endDate = time.Date(2099, 12, 31, 0, 0, 0, 0, time.UTC)
	}

	scheduleRuleRepo := repositories.NewScheduleRuleRepository(ctl.app)

	// 計算檢查日期（從開始日期起的未來一週）
	checkDate := startDate

	// 檢查每個 weekday 是否有衝突
	var overlappingRules []models.ScheduleRule
	var personalEventConflicts []models.PersonalEvent
	for _, weekday := range req.Weekdays {
		// 找到該 weekday 的日期
		current := checkDate
		weekdayDiff := weekday - int(current.Weekday())
		if weekdayDiff <= 0 {
			weekdayDiff += 7
		}
		targetDate := current.AddDate(0, 0, weekdayDiff)

		overlaps, eventConflicts, err := scheduleRuleRepo.CheckOverlap(ctx, centerID, req.RoomID, &req.TeacherID, weekday, req.StartTime, req.EndTime, nil, targetDate)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, global.ApiResponse{
				Code:    500,
				Message: "Failed to check overlap: " + err.Error(),
			})
			return
		}
		overlappingRules = append(overlappingRules, overlaps...)
		personalEventConflicts = append(personalEventConflicts, eventConflicts...)
	}

	if len(overlappingRules) > 0 || len(personalEventConflicts) > 0 {
		// 產生衝突訊息 - 去重，只顯示衝突的星期
		conflictSet := make(map[string]bool) // 用於去重
		conflictMessages := []string{}
		
		// 排課規則衝突
		for _, rule := range overlappingRules {
			dayNames := []string{"", "週一", "週二", "週三", "週四", "週五", "週六", "週日"}
			key := fmt.Sprintf("rule-%d-%s-%s", rule.Weekday, rule.StartTime, rule.EndTime)
			
			if !conflictSet[key] {
				conflictSet[key] = true
				msg := fmt.Sprintf("%s %s-%s", dayNames[rule.Weekday], rule.StartTime, rule.EndTime)
				if rule.TeacherID != nil {
					msg += " 老師已有排課"
				} else {
					msg += " 教室已有排課"
				}
				conflictMessages = append(conflictMessages, msg)
			}
		}

		// 個人行程衝突
		for _, event := range personalEventConflicts {
			key := fmt.Sprintf("event-%s", event.Title)
			
			if !conflictSet[key] {
				conflictSet[key] = true
				eventDate := event.StartAt
				dayNames := []string{"", "週日", "週一", "週二", "週三", "週四", "週五", "週六"}
				msg := fmt.Sprintf("%s %s-%s 「%s」",
					dayNames[eventDate.Weekday()],
					event.StartAt.Format("15:04"),
					event.EndAt.Format("15:04"),
					event.Title)
				msg += " 老師個人行程"
				conflictMessages = append(conflictMessages, msg)
			}
		}

		ctx.JSON(http.StatusBadRequest, global.ApiResponse{
			Code:    40002, // OVERLAP error code
			Message: "排課時間與現有規則或個人行程衝突",
			Datas: map[string]interface{}{
				"conflicts": conflictMessages,
			},
		})
		return
	}

	var createdRules []models.ScheduleRule

	for _, weekday := range req.Weekdays {
		rule := models.ScheduleRule{
			CenterID:   centerID,
			OfferingID: req.OfferingID,
			TeacherID:  &req.TeacherID,
			RoomID:     req.RoomID,
			Weekday:    weekday,
			StartTime:  req.StartTime,
			EndTime:    req.EndTime,
			EffectiveRange: models.DateRange{
				StartDate: startDate,
				EndDate:   endDate,
			},
		}

		createdRule, err := scheduleRuleRepo.Create(ctx, rule)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, global.ApiResponse{
				Code:    500,
				Message: "Failed to create schedule rule: " + err.Error(),
			})
			return
		}
		createdRules = append(createdRules, createdRule)
	}

	ctx.JSON(http.StatusOK, global.ApiResponse{
		Code:    0,
		Message: "Schedule rules created successfully",
		Datas:   createdRules,
	})
}

func (ctl *SchedulingController) GetRules(ctx *gin.Context) {
	centerID := ctx.GetUint(global.CenterIDKey)
	if centerID == 0 {
		ctx.JSON(http.StatusBadRequest, global.ApiResponse{
			Code:    global.BAD_REQUEST,
			Message: "Center ID required",
		})
		return
	}

	scheduleRuleRepo := repositories.NewScheduleRuleRepository(ctl.app)
	rules, err := scheduleRuleRepo.ListByCenterID(ctx, centerID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, global.ApiResponse{
			Code:    500,
			Message: err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, global.ApiResponse{
		Code:    0,
		Message: "OK",
		Datas:   rules,
	})
}

func (ctl *SchedulingController) DeleteRule(ctx *gin.Context) {
	centerID := ctx.GetUint(global.CenterIDKey)
	if centerID == 0 {
		ctx.JSON(http.StatusBadRequest, global.ApiResponse{
			Code:    global.BAD_REQUEST,
			Message: "Center ID required",
		})
		return
	}

	ruleIDStr := ctx.Param("ruleId")
	var ruleID uint
	if _, err := fmt.Sscanf(ruleIDStr, "%d", &ruleID); err != nil {
		ctx.JSON(http.StatusBadRequest, global.ApiResponse{
			Code:    global.BAD_REQUEST,
			Message: "Invalid rule ID",
		})
		return
	}

	scheduleRuleRepo := repositories.NewScheduleRuleRepository(ctl.app)
	if err := scheduleRuleRepo.DeleteByIDAndCenterID(ctx, ruleID, centerID); err != nil {
		ctx.JSON(http.StatusInternalServerError, global.ApiResponse{
			Code:    500,
			Message: "Failed to delete rule: " + err.Error(),
		})
		return
	}

	actorID := ctx.GetUint(global.UserIDKey)
	ctl.auditLogRepo.Create(ctx, models.AuditLog{
		CenterID:   centerID,
		ActorType:  "ADMIN",
		ActorID:    actorID,
		Action:     "DELETE_SCHEDULE_RULE",
		TargetType: "ScheduleRule",
		TargetID:   ruleID,
		Payload: models.AuditPayload{
			After: map[string]interface{}{
				"status": "DELETED",
			},
		},
	})

	ctx.JSON(http.StatusOK, global.ApiResponse{
		Code:    0,
		Message: "Schedule rule deleted successfully",
	})
}

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
	UpdateMode string  `json:"update_mode"`
}

const (
	UpdateModeSingle  = "SINGLE"
	UpdateModeFuture  = "FUTURE"
	UpdateModeAll     = "ALL"
)

func (ctl *SchedulingController) UpdateRule(ctx *gin.Context) {
	centerID := ctx.GetUint(global.CenterIDKey)
	if centerID == 0 {
		ctx.JSON(http.StatusBadRequest, global.ApiResponse{
			Code:    global.BAD_REQUEST,
			Message: "Center ID required",
		})
		return
	}

	ruleIDStr := ctx.Param("ruleId")
	var ruleID uint
	if _, err := fmt.Sscanf(ruleIDStr, "%d", &ruleID); err != nil {
		ctx.JSON(http.StatusBadRequest, global.ApiResponse{
			Code:    global.BAD_REQUEST,
			Message: "Invalid rule ID",
		})
		return
	}

	var req UpdateRuleRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, global.ApiResponse{
			Code:    global.BAD_REQUEST,
			Message: "Invalid request parameters: " + err.Error(),
		})
		return
	}

	// 解析日期
	var startDate, endDate time.Time
	var err error

	// 如果提供了 start_date，則解析；否則保持為零值（表示不更新）
	if req.StartDate != "" {
		startDate, err = time.Parse("2006-01-02", req.StartDate)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, global.ApiResponse{
				Code:    global.BAD_REQUEST,
				Message: "Invalid start_date format",
			})
			return
		}
	}

	// 如果提供了 end_date，則解析；否則保持為零值（表示不更新）
	if req.EndDate != nil && *req.EndDate != "" {
		endDate, err = time.Parse("2006-01-02", *req.EndDate)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, global.ApiResponse{
				Code:    global.BAD_REQUEST,
				Message: "Invalid end_date format",
			})
			return
		}
	}

	scheduleRuleRepo := repositories.NewScheduleRuleRepository(ctl.app)

	// 檢查規則是否存在且屬於該中心
	existingRule, err := scheduleRuleRepo.GetByIDAndCenterID(ctx, ruleID, centerID)
	if err != nil {
		ctx.JSON(http.StatusNotFound, global.ApiResponse{
			Code:    errInfos.NOT_FOUND,
			Message: "Rule not found",
		})
		return
	}

	// 找出同一組的所有規則（相同的 name 或 offering_id 和時間）
	allRules, err := scheduleRuleRepo.ListByCenterID(ctx, centerID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, global.ApiResponse{
			Code:    500,
			Message: "Failed to fetch rules: " + err.Error(),
		})
		return
	}

	// 找出與當前規則相關聯的規則（相同的 offering_id、start_time、end_time）
	var relatedRules []models.ScheduleRule
	for _, rule := range allRules {
		if rule.OfferingID == existingRule.OfferingID &&
			rule.StartTime == existingRule.StartTime &&
			rule.EndTime == existingRule.EndTime &&
			rule.ID != existingRule.ID {
			relatedRules = append(relatedRules, rule)
		}
	}

	// 根據 update_mode 處理不同的更新策略
	var resultRules []models.ScheduleRule

	switch req.UpdateMode {
	case UpdateModeFuture:
		// FUTURE: 截斷現有規則，創建新規則段
		resultRules, err = handleFutureUpdate(ctx, scheduleRuleRepo, centerID, existingRule, relatedRules, req, startDate, endDate)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, global.ApiResponse{
				Code:    500,
				Message: "Failed to update rule (FUTURE): " + err.Error(),
			})
			return
		}
	case UpdateModeSingle:
		// SINGLE: 只更新這個規則（不改變 weekday 結構）
		exceptionRepo := repositories.NewScheduleExceptionRepository(ctl.app)
		resultRules, err = handleSingleUpdate(ctx, scheduleRuleRepo, exceptionRepo, centerID, existingRule, req, startDate)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, global.ApiResponse{
				Code:    500,
				Message: "Failed to update rule (SINGLE): " + err.Error(),
			})
			return
		}
	default:
		// ALL 或空：更新所有相關規則
		resultRules, err = handleAllUpdate(ctx, scheduleRuleRepo, centerID, existingRule, relatedRules, req, startDate, endDate)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, global.ApiResponse{
				Code:    500,
				Message: "Failed to update rule (ALL): " + err.Error(),
			})
			return
		}
	}

	ctx.JSON(http.StatusOK, global.ApiResponse{
		Code:    0,
		Message: "Schedule rules updated successfully",
		Datas:   resultRules,
	})
}

// handleFutureUpdate 處理 FUTURE 模式：截斷現有規則，創建新規則段
func handleFutureUpdate(ctx context.Context, ruleRepo *repositories.ScheduleRuleRepository, centerID uint, existingRule models.ScheduleRule, relatedRules []models.ScheduleRule, req UpdateRuleRequest, startDate, endDate time.Time) ([]models.ScheduleRule, error) {
	var result []models.ScheduleRule

	// 如果沒有提供 startDate，無法執行 FUTURE 模式
	if startDate.IsZero() {
		return nil, errors.New("start_date is required for FUTURE update mode")
	}

	cutoffDate := startDate.AddDate(0, 0, -1) // 前一天

	// 1. 更新現有規則的 end_date 到 cutoffDate
	for _, rule := range append(relatedRules, existingRule) {
		rule.EffectiveRange.EndDate = cutoffDate
		if err := ruleRepo.Update(ctx, rule); err != nil {
			return nil, err
		}
		result = append(result, rule)
	}

	// 2. 創建新的規則段
	newRules := createNewRuleSegment(centerID, existingRule, relatedRules, req, startDate, endDate)
	for _, newRule := range newRules {
		created, err := ruleRepo.Create(ctx, newRule)
		if err != nil {
			return nil, err
		}
		result = append(result, created)
	}

	return result, nil
}

// handleSingleUpdate 處理 SINGLE 模式：針對單一日期建立例外單
func handleSingleUpdate(ctx context.Context, ruleRepo *repositories.ScheduleRuleRepository, exceptionRepo *repositories.ScheduleExceptionRepository, centerID uint, existingRule models.ScheduleRule, req UpdateRuleRequest, targetDate time.Time) ([]models.ScheduleRule, error) {
	// SINGLE 模式：針對 targetDate 這一天進行變更
	// 1. 建立 CANCEL 例外單（取消這一天的原場次）

	now := time.Now()
	cancelException := models.ScheduleException{
		CenterID:     centerID,
		RuleID:       existingRule.ID,
		OriginalDate: targetDate,
		Type:         "CANCEL",
		Status:       "APPROVED", // SINGLE 模式直接核准
		Reason:       req.Name,   // 使用更新原因作為取消原因
		ReviewedAt:   &now,
	}

	if _, err := exceptionRepo.Create(ctx, cancelException); err != nil {
		return nil, fmt.Errorf("failed to create cancel exception: %w", err)
	}

	// 2. 建立新時間的規則（視為新規則，只針對這一天）
	weekday := int(targetDate.Weekday())
	if weekday == 0 {
		weekday = 7 // 週日轉換為 7
	}

	newRule := models.ScheduleRule{
		CenterID:       centerID,
		OfferingID:     req.OfferingID,
		TeacherID:      req.TeacherID,
		RoomID:         req.RoomID,
		Name:           req.Name,
		Weekday:        weekday,
		StartTime:      req.StartTime,
		EndTime:        req.EndTime,
		Duration:       req.Duration,
		EffectiveRange: models.DateRange{
			StartDate: targetDate,
			EndDate:   targetDate, // 只有這一天
		},
	}

	created, err := ruleRepo.Create(ctx, newRule)
	if err != nil {
		return nil, fmt.Errorf("failed to create new rule for reschedule: %w", err)
	}

	// 3. 建立 RESCHEDULE 例外單關聯新舊規則（如果新時間與原時間不同）
	startHour := parseHour(req.StartTime)
	endHour := parseHour(req.EndTime)

	newStartAt := time.Date(targetDate.Year(), targetDate.Month(), targetDate.Day(), startHour, 0, 0, 0, time.UTC)
	newEndAt := time.Date(targetDate.Year(), targetDate.Month(), targetDate.Day(), endHour, 0, 0, 0, time.UTC)

	rescheduleException := models.ScheduleException{
		CenterID:     centerID,
		RuleID:       existingRule.ID, // 關聯原規則
		OriginalDate: targetDate,
		Type:         "RESCHEDULE",
		Status:       "APPROVED",
		Reason:       req.Name,
		NewStartAt:   &newStartAt,
		NewEndAt:     &newEndAt,
		ReviewedAt:   &now,
	}

	if _, err := exceptionRepo.Create(ctx, rescheduleException); err != nil {
		// 例外單建立失敗不影響主要流程，記錄錯誤但繼續
		fmt.Printf("Warning: failed to create reschedule exception: %v\n", err)
	}

	return []models.ScheduleRule{existingRule, created}, nil
}

// parseHour 從時間字串（如 "14:00"）解析小時
func parseHour(timeStr string) int {
	if timeStr == "" {
		return 0
	}
	parts := strings.Split(timeStr, ":")
	if len(parts) < 1 {
		return 0
	}
	hour, _ := strconv.Atoi(parts[0])
	return hour
}

// handleAllUpdate 處理 ALL 模式：更新所有相關規則
func handleAllUpdate(ctx context.Context, ruleRepo *repositories.ScheduleRuleRepository, centerID uint, existingRule models.ScheduleRule, relatedRules []models.ScheduleRule, req UpdateRuleRequest, startDate, endDate time.Time) ([]models.ScheduleRule, error) {
	// 收集所有相關的 weekday
	existingWeekdays := []int{existingRule.Weekday}
	for _, rule := range relatedRules {
		existingWeekdays = append(existingWeekdays, rule.Weekday)
	}

	// 轉換 req.Weekdays 為 map 便於查找
	newWeekdayMap := make(map[int]bool)
	for _, w := range req.Weekdays {
		newWeekdayMap[w] = true
	}

	// 追蹤需要更新和創建的規則
	var updatedRules []models.ScheduleRule
	var createdRules []models.ScheduleRule
	var deletedRuleIDs []uint

	// 1. 處理現有規則：更新或標記刪除
	for _, rule := range append(relatedRules, existingRule) {
		if newWeekdayMap[rule.Weekday] {
			// 這個 weekday 仍然需要保留，更新它
			updatedRule := rule
			if req.Name != "" {
				updatedRule.Name = req.Name
			}
			if req.OfferingID != 0 {
				updatedRule.OfferingID = req.OfferingID
			}
			if req.TeacherID != nil {
				updatedRule.TeacherID = req.TeacherID
			}
			if req.RoomID != 0 {
				updatedRule.RoomID = req.RoomID
			}
			if req.StartTime != "" {
				updatedRule.StartTime = req.StartTime
			}
			if req.EndTime != "" {
				updatedRule.EndTime = req.EndTime
			}
			if req.Duration != 0 {
				updatedRule.Duration = req.Duration
			}
			// 只有當新日期與現有日期不同時才更新（避免不必要的更新）
			if !startDate.IsZero() && startDate != rule.EffectiveRange.StartDate {
				updatedRule.EffectiveRange.StartDate = startDate
			}
			// endDate 為零值時保持現有值
			if !endDate.IsZero() {
				updatedRule.EffectiveRange.EndDate = endDate
			}

			updatedRules = append(updatedRules, updatedRule)
			delete(newWeekdayMap, rule.Weekday)
		} else {
			// 這個 weekday 不再需要，標記刪除
			deletedRuleIDs = append(deletedRuleIDs, rule.ID)
		}
	}

	// 2. 創建新的規則
	for weekday := range newWeekdayMap {
		newRule := createSingleRule(centerID, existingRule, req, weekday, startDate, endDate)
		createdRules = append(createdRules, *newRule)
	}

	// 3. 執行資料庫操作
	// 先刪除
	for _, id := range deletedRuleIDs {
		if err := ruleRepo.Delete(ctx, id); err != nil {
			return nil, err
		}
	}

	// 更新
	for _, rule := range updatedRules {
		if err := ruleRepo.Update(ctx, rule); err != nil {
			return nil, err
		}
	}

	// 創建
	for _, rule := range createdRules {
		created, err := ruleRepo.Create(ctx, rule)
		if err != nil {
			return nil, err
		}
		updatedRules = append(updatedRules, created)
	}

	return updatedRules, nil
}

// createNewRuleSegment 為 FUTURE 模式創建新的規則段
func createNewRuleSegment(centerID uint, existingRule models.ScheduleRule, relatedRules []models.ScheduleRule, req UpdateRuleRequest, startDate, endDate time.Time) []models.ScheduleRule {
	// 收集所有需要創建的 weekday
	weekdaysToCreate := []int{existingRule.Weekday}
	for _, rule := range relatedRules {
		weekdaysToCreate = append(weekdaysToCreate, rule.Weekday)
	}

	var newRules []models.ScheduleRule
	for _, weekday := range weekdaysToCreate {
		newRule := createSingleRule(centerID, existingRule, req, weekday, startDate, endDate)
		newRules = append(newRules, *newRule)
	}

	return newRules
}

// createSingleRule 創建單個規則
func createSingleRule(centerID uint, existingRule models.ScheduleRule, req UpdateRuleRequest, weekday int, startDate, endDate time.Time) *models.ScheduleRule {
	// 如果日期為零值，使用現有規則的日期
	effectiveStartDate := startDate
	if startDate.IsZero() {
		effectiveStartDate = existingRule.EffectiveRange.StartDate
	}
	effectiveEndDate := endDate
	if endDate.IsZero() {
		effectiveEndDate = existingRule.EffectiveRange.EndDate
	}

	newRule := &models.ScheduleRule{
		CenterID:   centerID,
		OfferingID: existingRule.OfferingID,
		TeacherID:  existingRule.TeacherID,
		RoomID:     existingRule.RoomID,
		Name:       existingRule.Name,
		Weekday:    weekday,
		StartTime:  existingRule.StartTime,
		EndTime:    existingRule.EndTime,
		Duration:   existingRule.Duration,
		EffectiveRange: models.DateRange{
			StartDate: effectiveStartDate,
			EndDate:   effectiveEndDate,
		},
	}
	if req.Name != "" {
		newRule.Name = req.Name
	}
	if req.OfferingID != 0 {
		newRule.OfferingID = req.OfferingID
	}
	if req.TeacherID != nil {
		newRule.TeacherID = req.TeacherID
	}
	if req.RoomID != 0 {
		newRule.RoomID = req.RoomID
	}
	if req.StartTime != "" {
		newRule.StartTime = req.StartTime
	}
	if req.EndTime != "" {
		newRule.EndTime = req.EndTime
	}
	if req.Duration != 0 {
		newRule.Duration = req.Duration
	}
	return newRule
}

type DetectPhaseTransitionsRequest struct {
	OfferingID uint      `json:"offering_id" binding:"required"`
	StartDate  time.Time `json:"start_date" binding:"required"`
	EndDate    time.Time `json:"end_date" binding:"required"`
}

func (ctl *SchedulingController) DetectPhaseTransitions(ctx *gin.Context) {
	centerID := ctx.GetUint(global.CenterIDKey)
	if centerID == 0 {
		ctx.JSON(http.StatusBadRequest, global.ApiResponse{
			Code:    global.BAD_REQUEST,
			Message: "Center ID required",
		})
		return
	}

	var req DetectPhaseTransitionsRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, global.ApiResponse{
			Code:    global.BAD_REQUEST,
			Message: "Invalid request parameters: " + err.Error(),
		})
		return
	}

	transitions, err := ctl.expansionService.DetectPhaseTransitions(ctx, centerID, req.OfferingID, req.StartDate, req.EndDate)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, global.ApiResponse{
			Code:    500,
			Message: err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, global.ApiResponse{
		Code:    0,
		Message: "OK",
		Datas:   transitions,
	})
}

type CheckRuleLockStatusRequest struct {
	RuleID        uint      `json:"rule_id" binding:"required"`
	ExceptionDate time.Time `json:"exception_date" binding:"required"`
}

type CheckRuleLockStatusResponse struct {
	IsLocked      bool       `json:"is_locked"`
	LockReason    string     `json:"lock_reason,omitempty"`
	LockAt        *time.Time `json:"lock_at,omitempty"`
	Deadline      time.Time  `json:"deadline"`
	DaysRemaining int        `json:"days_remaining"`
}

func (ctl *SchedulingController) CheckRuleLockStatus(ctx *gin.Context) {
	centerID := ctx.GetUint(global.CenterIDKey)
	if centerID == 0 {
		ctx.JSON(http.StatusBadRequest, global.ApiResponse{
			Code:    global.BAD_REQUEST,
			Message: "Center ID required",
		})
		return
	}

	var req CheckRuleLockStatusRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, global.ApiResponse{
			Code:    global.BAD_REQUEST,
			Message: "Invalid request parameters: " + err.Error(),
		})
		return
	}

	scheduleRuleRepo := repositories.NewScheduleRuleRepository(ctl.app)
	rule, err := scheduleRuleRepo.GetByID(ctx, req.RuleID)
	if err != nil {
		ctx.JSON(http.StatusNotFound, global.ApiResponse{
			Code:    errInfos.NOT_FOUND,
			Message: "Rule not found",
		})
		return
	}

	if rule.CenterID != centerID {
		ctx.JSON(http.StatusForbidden, global.ApiResponse{
			Code:    errInfos.FORBIDDEN,
			Message: "Rule does not belong to this center",
		})
		return
	}

	now := time.Now()
	response := CheckRuleLockStatusResponse{
		DaysRemaining: -1,
	}

	if rule.LockAt != nil && now.After(*rule.LockAt) {
		response.IsLocked = true
		response.LockReason = "已超過異動截止日"
		response.LockAt = rule.LockAt
	}

	center, _ := ctl.centerRepo.GetByID(ctx, centerID)
	leadDays := center.Settings.ExceptionLeadDays
	if leadDays <= 0 {
		leadDays = 14
	}

	deadline := req.ExceptionDate.AddDate(0, 0, -leadDays)
	response.Deadline = deadline
	daysRemaining := int(deadline.Sub(now).Hours() / 24)
	response.DaysRemaining = daysRemaining

	if daysRemaining < 0 && !response.IsLocked {
		response.IsLocked = true
		response.LockReason = fmt.Sprintf("已超過異動截止日（需提前 %d 天申請）", leadDays)
	}

	ctx.JSON(http.StatusOK, global.ApiResponse{
		Code:    0,
		Message: "OK",
		Datas:   response,
	})
}
