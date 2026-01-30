package controllers

import (
	"time"
	"timeLedger/app"
	"timeLedger/app/models"
	"timeLedger/app/repositories"
	"timeLedger/app/services"

	"github.com/gin-gonic/gin"
)

// TeacherExceptionController 老師例外申請相關 API
type TeacherExceptionController struct {
	BaseController
	app              *app.App
	exceptionRepo    *repositories.ScheduleExceptionRepository
	exceptionService services.ScheduleExceptionService
	membershipRepo   *repositories.CenterMembershipRepository
	auditLogRepo     *repositories.AuditLogRepository
}

func NewTeacherExceptionController(app *app.App) *TeacherExceptionController {
	return &TeacherExceptionController{
		app:              app,
		exceptionRepo:    repositories.NewScheduleExceptionRepository(app),
		exceptionService: services.NewScheduleExceptionService(app),
		membershipRepo:   repositories.NewCenterMembershipRepository(app),
		auditLogRepo:     repositories.NewAuditLogRepository(app),
	}
}

// TeacherCreateExceptionRequest 老師建立例外申請請求結構
type TeacherCreateExceptionRequest struct {
	CenterID     uint   `json:"center_id" binding:"required"`
	RuleID       uint   `json:"rule_id" binding:"required"`
	OriginalDate string `json:"original_date" binding:"required"`
	Type         string `json:"type" binding:"required,oneof=CANCEL RESCHEDULE REPLACE_TEACHER"`
	Reason       string `json:"reason" binding:"required"`
	NewStartAt   string `json:"new_start_at"`
	NewEndAt     string `json:"new_end_at"`
	NewTeacherID *uint  `json:"new_teacher_id"`
	NewTeacherName string `json:"new_teacher_name"`
}

// CreateException 老師提出停課/改期申請
// @Summary 老師提出停課/改期申請
// @Tags Teacher
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body TeacherCreateExceptionRequest true "例外申請"
// @Success 200 {object} global.ApiResponse{data=models.ScheduleException}
// @Router /api/v1/teacher/exceptions [post]
func (ctl *TeacherExceptionController) CreateException(ctx *gin.Context) {
	helper := NewContextHelper(ctx)

	teacherID := helper.MustUserID()
	if teacherID == 0 {
		return
	}

	var req TeacherCreateExceptionRequest
	if !helper.MustBindJSON(&req) {
		return
	}

	// 解析時間欄位
	originalDate, err := time.Parse("2006-01-02", req.OriginalDate)
	if err != nil {
		helper.BadRequest("Invalid original_date format, expected YYYY-MM-DD")
		return
	}

	var newStartAt, newEndAt *time.Time
	if req.NewStartAt != "" {
		t, err := time.Parse(time.RFC3339, req.NewStartAt)
		if err != nil {
			helper.BadRequest("Invalid new_start_at format, expected RFC3339")
			return
		}
		newStartAt = &t
	}
	if req.NewEndAt != "" {
		t, err := time.Parse(time.RFC3339, req.NewEndAt)
		if err != nil {
			helper.BadRequest("Invalid new_end_at format, expected RFC3339")
			return
		}
		newEndAt = &t
	}

	exception, err := ctl.exceptionService.CreateException(
		ctx,
		req.CenterID,
		teacherID,
		req.RuleID,
		originalDate,
		req.Type,
		newStartAt,
		newEndAt,
		req.NewTeacherID,
		req.NewTeacherName,
		req.Reason,
	)
	if err != nil {
		helper.InternalError(err.Error())
		return
	}

	ctl.auditLogRepo.Create(ctx, models.AuditLog{
		CenterID:   req.CenterID,
		ActorType:  "TEACHER",
		ActorID:    teacherID,
		Action:     "EXCEPTION_CREATE",
		TargetType: "ScheduleException",
		TargetID:   exception.ID,
		Payload: models.AuditPayload{
			After: map[string]interface{}{
				"rule_id":       req.RuleID,
				"original_date": req.OriginalDate,
				"type":          req.Type,
				"reason":        req.Reason,
			},
		},
	})

	helper.Success(exception)
}

// RevokeException 老師撤回待審核的例外申請
// @Summary 老師撤回待審核的例外申請
// @Tags Teacher
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Exception ID"
// @Success 200 {object} global.ApiResponse
// @Router /api/v1/teacher/exceptions/{id}/revoke [post]
func (ctl *TeacherExceptionController) RevokeException(ctx *gin.Context) {
	helper := NewContextHelper(ctx)

	teacherID := helper.MustUserID()
	if teacherID == 0 {
		return
	}

	exceptionID := helper.MustParamUint("id")
	if exceptionID == 0 {
		return
	}

	err := ctl.exceptionService.RevokeException(ctx, exceptionID, teacherID)
	if err != nil {
		helper.InternalError(err.Error())
		return
	}

	exception, _ := ctl.exceptionRepo.GetByID(ctx, exceptionID)
	if exception.CenterID > 0 {
		ctl.auditLogRepo.Create(ctx, models.AuditLog{
			CenterID:   exception.CenterID,
			ActorType:  "TEACHER",
			ActorID:    teacherID,
			Action:     "EXCEPTION_REVOKE",
			TargetType: "ScheduleException",
			TargetID:   exceptionID,
			Payload: models.AuditPayload{
				After: map[string]interface{}{
					"status": "REVOKED",
				},
			},
		})
	}

	helper.Success(nil)
}

// GetExceptions 老師查看自己的例外申請列表
// @Summary 老師查看自己的例外申請列表
// @Tags Teacher
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param status query string false "篩選狀態 (PENDING/APPROVED/REJECTED/REVOKED)"
// @Success 200 {object} global.ApiResponse{data=[]models.ScheduleException}
// @Router /api/v1/teacher/exceptions [get]
func (ctl *TeacherExceptionController) GetExceptions(ctx *gin.Context) {
	helper := NewContextHelper(ctx)

	teacherID := helper.MustUserID()
	if teacherID == 0 {
		return
	}

	// 從 Query 取得狀態篩選
	status := helper.QueryStringOrDefault("status", "")

	exceptions, err := ctl.exceptionRepo.GetByTeacherID(ctx, teacherID, status)
	if err != nil {
		helper.InternalError("Failed to get exceptions")
		return
	}

	helper.Success(exceptions)
}
