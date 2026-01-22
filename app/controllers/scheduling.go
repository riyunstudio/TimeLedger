package controllers

import (
	"fmt"
	"net/http"
	"time"
	"timeLedger/app"
	"timeLedger/app/models"
	"timeLedger/app/repositories"
	"timeLedger/app/services"
	"timeLedger/global"

	"github.com/gin-gonic/gin"
)

type SchedulingController struct {
	BaseController
	app               *app.App
	validationService services.ScheduleValidationService
	expansionService  services.ScheduleExpansionService
	exceptionService  services.ScheduleExceptionService
	auditLogRepo      *repositories.AuditLogRepository
}

func NewSchedulingController(app *app.App) *SchedulingController {
	return &SchedulingController{
		app:               app,
		validationService: services.NewScheduleValidationService(app),
		expansionService:  services.NewScheduleExpansionService(app),
		exceptionService:  services.NewScheduleExceptionService(app),
		auditLogRepo:      repositories.NewAuditLogRepository(app),
	}
}

type CheckOverlapRequest struct {
	CenterID      uint      `json:"center_id" binding:"required"`
	TeacherID     *uint     `json:"teacher_id"`
	RoomID        uint      `json:"room_id" binding:"required"`
	StartTime     time.Time `json:"start_time" binding:"required"`
	EndTime       time.Time `json:"end_time" binding:"required"`
	ExcludeRuleID *uint     `json:"exclude_rule_id"`
}

type CheckBufferRequest struct {
	CenterID      uint      `json:"center_id" binding:"required"`
	TeacherID     uint      `json:"teacher_id" binding:"required"`
	RoomID        uint      `json:"room_id" binding:"required"`
	PrevEndTime   time.Time `json:"prev_end_time" binding:"required"`
	NextStartTime time.Time `json:"next_start_time" binding:"required"`
	CourseID      uint      `json:"course_id" binding:"required"`
}

type CreateExceptionRequest struct {
	CenterID     uint       `json:"center_id" binding:"required"`
	RuleID       uint       `json:"rule_id" binding:"required"`
	OriginalDate time.Time  `json:"original_date" binding:"required"`
	Type         string     `json:"type" binding:"required"`
	NewStartAt   *time.Time `json:"new_start_at"`
	NewEndAt     *time.Time `json:"new_end_at"`
	NewTeacherID *uint      `json:"new_teacher_id"`
	NewRoomID    *uint      `json:"new_room_id"`
	Reason       string     `json:"reason" binding:"required"`
}

type ReviewExceptionRequest struct {
	Action         string `json:"action" binding:"required"`
	OverrideBuffer bool   `json:"override_buffer"`
	Reason         string `json:"reason"`
}

type ExpandRulesRequest struct {
	RuleIDs   []uint    `json:"rule_ids" binding:"required"`
	StartDate time.Time `json:"start_date" binding:"required"`
	EndDate   time.Time `json:"end_date" binding:"required"`
}

type ValidateFullRequest struct {
	CenterID            uint      `json:"center_id" binding:"required"`
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

	result, err := ctl.validationService.CheckOverlap(ctx, req.CenterID, req.TeacherID, req.RoomID, req.StartTime, req.EndTime, req.ExcludeRuleID)
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

	result, err := ctl.validationService.CheckTeacherBuffer(ctx, req.CenterID, req.TeacherID, req.PrevEndTime, req.NextStartTime, req.CourseID)
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

	result, err := ctl.validationService.CheckRoomBuffer(ctx, req.CenterID, req.RoomID, req.PrevEndTime, req.NextStartTime, req.CourseID)
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

	adminID := ctx.GetUint(global.UserIDKey)
	result, err := ctl.validationService.ValidateFull(ctx, req.CenterID, req.TeacherID, req.RoomID, req.CourseID, req.StartTime, req.EndTime, req.ExcludeRuleID, req.AllowBufferOverride)
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
		CenterID:   req.CenterID,
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

	adminID := ctx.GetUint(global.UserIDKey)
	exception, err := ctl.exceptionService.CreateException(ctx, req.CenterID, adminID, req.RuleID, req.OriginalDate, req.Type, req.NewStartAt, req.NewEndAt, req.NewTeacherID, req.Reason)
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
	centerID := ctx.Param("id")
	var centerIDInt uint
	if _, err := fmt.Sscanf(centerID, "%d", &centerIDInt); err != nil {
		ctx.JSON(http.StatusBadRequest, global.ApiResponse{
			Code:    global.BAD_REQUEST,
			Message: "Invalid center ID",
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

	exceptions, err := ctl.exceptionService.GetExceptionsByDateRange(ctx, centerIDInt, startDate, endDate)
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
	centerID := ctx.Param("id")
	var centerIDInt uint
	if _, err := fmt.Sscanf(centerID, "%d", &centerIDInt); err != nil {
		ctx.JSON(http.StatusBadRequest, global.ApiResponse{
			Code:    global.BAD_REQUEST,
			Message: "Invalid center ID",
		})
		return
	}

	var req ExpandRulesRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, global.ApiResponse{
			Code:    global.BAD_REQUEST,
			Message: "Invalid request parameters",
		})
		return
	}

	scheduleRuleRepo := repositories.NewScheduleRuleRepository(ctl.app)
	rules, err := scheduleRuleRepo.ListByCenterID(ctx, centerIDInt)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, global.ApiResponse{
			Code:    500,
			Message: err.Error(),
		})
		return
	}

	var filteredRules []models.ScheduleRule
	if len(req.RuleIDs) > 0 {
		for _, rule := range rules {
			for _, ruleID := range req.RuleIDs {
				if rule.ID == ruleID {
					filteredRules = append(filteredRules, rule)
					break
				}
			}
		}
	} else {
		filteredRules = rules
	}

	expandedSchedules := ctl.expansionService.ExpandRules(ctx, filteredRules, req.StartDate, req.EndDate, centerIDInt)

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
