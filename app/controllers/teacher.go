package controllers

import (
	"crypto/rand"
	"fmt"
	"net/http"
	"time"
	"timeLedger/app"
	"timeLedger/app/models"
	"timeLedger/app/repositories"
	"timeLedger/app/requests"
	"timeLedger/app/resources"
	"timeLedger/app/services"
	"timeLedger/global"
	"timeLedger/global/errInfos"

	"github.com/gin-gonic/gin"
)

type TeacherController struct {
	BaseController
	app               *app.App
	teacherRepository *repositories.TeacherRepository
	membershipRepo    *repositories.CenterMembershipRepository
	scheduleRuleRepo  *repositories.ScheduleRuleRepository
	exceptionRepo     *repositories.ScheduleExceptionRepository
	exceptionService  services.ScheduleExceptionService
	expansionService  services.ScheduleExpansionService
	recurrenceService services.ScheduleRecurrenceService
	auditLogRepo      *repositories.AuditLogRepository
	skillRepo         *repositories.TeacherSkillRepository
	certificateRepo   *repositories.TeacherCertificateRepository
	personalEventRepo *repositories.PersonalEventRepository
	sessionNoteRepo   *repositories.SessionNoteRepository
}

func NewTeacherController(app *app.App) *TeacherController {
	return &TeacherController{
		app:               app,
		teacherRepository: repositories.NewTeacherRepository(app),
		membershipRepo:    repositories.NewCenterMembershipRepository(app),
		scheduleRuleRepo:  repositories.NewScheduleRuleRepository(app),
		exceptionRepo:     repositories.NewScheduleExceptionRepository(app),
		exceptionService:  services.NewScheduleExceptionService(app),
		expansionService:  services.NewScheduleExpansionService(app),
		recurrenceService: services.NewScheduleRecurrenceService(app),
		auditLogRepo:      repositories.NewAuditLogRepository(app),
		skillRepo:         repositories.NewTeacherSkillRepository(app),
		certificateRepo:   repositories.NewTeacherCertificateRepository(app),
		personalEventRepo: repositories.NewPersonalEventRepository(app),
		sessionNoteRepo:   repositories.NewSessionNoteRepository(app),
	}
}

// GetProfile 取得老師個人資料
// @Summary 取得老師個人資料
// @Tags Teacher
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} global.ApiResponse{data=resources.TeacherProfileResource}
// @Router /api/v1/teacher/me/profile [get]
func (ctl *TeacherController) GetProfile(ctx *gin.Context) {
	teacherID := ctx.GetUint(global.UserIDKey)
	if teacherID == 0 {
		ctx.JSON(http.StatusUnauthorized, global.ApiResponse{
			Code:    global.UNAUTHORIZED,
			Message: "Teacher ID required",
		})
		return
	}

	teacher, err := ctl.teacherRepository.GetByID(ctx, teacherID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, global.ApiResponse{
			Code:    500,
			Message: "Failed to get teacher profile",
		})
		return
	}

	response := resources.TeacherProfileResource{
		ID:                teacher.ID,
		LineUserID:        teacher.LineUserID,
		Name:              teacher.Name,
		Email:             teacher.Email,
		Bio:               teacher.Bio,
		City:              teacher.City,
		District:          teacher.District,
		PublicContactInfo: teacher.PublicContactInfo,
		IsOpenToHiring:    teacher.IsOpenToHiring,
	}

	ctx.JSON(http.StatusOK, global.ApiResponse{
		Code:    0,
		Message: "Success",
		Datas:   response,
	})
}

// UpdateProfile 更新老師個人資料
// @Summary 更新老師個人資料
// @Tags Teacher
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body UpdateTeacherProfileRequest true "個人資料"
// @Success 200 {object} global.ApiResponse{data=resources.TeacherProfileResource}
// @Router /api/v1/teacher/me/profile [put]
func (ctl *TeacherController) UpdateProfile(ctx *gin.Context) {
	teacherID := ctx.GetUint(global.UserIDKey)
	if teacherID == 0 {
		ctx.JSON(http.StatusUnauthorized, global.ApiResponse{
			Code:    global.UNAUTHORIZED,
			Message: "Teacher ID required",
		})
		return
	}

	var req UpdateTeacherProfileRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, global.ApiResponse{
			Code:    global.BAD_REQUEST,
			Message: "Invalid request body",
		})
		return
	}

	teacher, err := ctl.teacherRepository.GetByID(ctx, teacherID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, global.ApiResponse{
			Code:    500,
			Message: "Failed to get teacher profile",
		})
		return
	}

	if req.Bio != "" {
		teacher.Bio = req.Bio
	}
	if req.City != "" {
		teacher.City = req.City
	}
	if req.District != "" {
		teacher.District = req.District
	}
	if req.PublicContactInfo != "" {
		teacher.PublicContactInfo = req.PublicContactInfo
	}

	teacher.IsOpenToHiring = req.IsOpenToHiring

	if err := ctl.teacherRepository.Update(ctx, teacher); err != nil {
		ctx.JSON(http.StatusInternalServerError, global.ApiResponse{
			Code:    500,
			Message: "Failed to update teacher profile",
		})
		return
	}

	memberships, _ := ctl.membershipRepo.GetActiveByTeacherID(ctx, teacherID)
	var centerID uint
	if len(memberships) > 0 {
		centerID = memberships[0].CenterID
	}

	ctl.auditLogRepo.Create(ctx, models.AuditLog{
		CenterID:   centerID,
		ActorType:  "TEACHER",
		ActorID:    teacherID,
		Action:     "PROFILE_UPDATE",
		TargetType: "Teacher",
		TargetID:   teacherID,
		Payload: models.AuditPayload{
			After: map[string]interface{}{
				"bio":               req.Bio,
				"city":              req.City,
				"district":          req.District,
				"is_open_to_hiring": req.IsOpenToHiring,
			},
		},
	})

	response := resources.TeacherProfileResource{
		ID:                teacher.ID,
		LineUserID:        teacher.LineUserID,
		Name:              teacher.Name,
		Email:             teacher.Email,
		Bio:               teacher.Bio,
		City:              teacher.City,
		District:          teacher.District,
		PublicContactInfo: teacher.PublicContactInfo,
		IsOpenToHiring:    teacher.IsOpenToHiring,
	}

	ctx.JSON(http.StatusOK, global.ApiResponse{
		Code:    0,
		Message: "Profile updated",
		Datas:   response,
	})
}

type UpdateTeacherProfileRequest struct {
	Bio               string `json:"bio"`
	City              string `json:"city"`
	District          string `json:"district"`
	PublicContactInfo string `json:"public_contact_info"`
	IsOpenToHiring    bool   `json:"is_open_to_hiring"`
}

// GetCenters 取得老師已加入的中心列表
// @Summary 取得老師已加入的中心列表
// @Tags Teacher
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} global.ApiResponse{data=[]resources.CenterMembershipResource}
// @Router /api/v1/teacher/me/centers [get]
func (ctl *TeacherController) GetCenters(ctx *gin.Context) {
	teacherID := ctx.GetUint(global.UserIDKey)
	if teacherID == 0 {
		ctx.JSON(http.StatusUnauthorized, global.ApiResponse{
			Code:    global.UNAUTHORIZED,
			Message: "Teacher ID required",
		})
		return
	}

	memberships, err := ctl.membershipRepo.GetActiveByTeacherID(ctx, teacherID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, global.ApiResponse{
			Code:    500,
			Message: "Failed to get memberships",
		})
		return
	}

	var centerResources []resources.CenterMembershipResource
	for _, m := range memberships {
		var center models.Center
		ctl.app.MySQL.RDB.WithContext(ctx).First(&center, m.CenterID)
		centerResources = append(centerResources, resources.CenterMembershipResource{
			ID:         m.ID,
			CenterID:   m.CenterID,
			CenterName: center.Name,
			Status:     string(m.Status),
			CreatedAt:  m.CreatedAt,
		})
	}

	ctx.JSON(http.StatusOK, global.ApiResponse{
		Code:    0,
		Message: "Success",
		Datas:   centerResources,
	})
}

// UploadCertificate 上傳證照
// @Summary 上傳證照
// @Tags Teacher
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body UploadCertificateRequest true "證照資訊"
// @Success 200 {object} global.ApiResponse{data=resources.TeacherCertificateResource}
// @Router /api/v1/teacher/me/certificates [post]
func (ctl *TeacherController) UploadCertificate(ctx *gin.Context) {
	teacherID := ctx.GetUint(global.UserIDKey)
	if teacherID == 0 {
		ctx.JSON(http.StatusUnauthorized, global.ApiResponse{
			Code:    global.UNAUTHORIZED,
			Message: "Teacher ID required",
		})
		return
	}

	var req UploadCertificateRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, global.ApiResponse{
			Code:    global.BAD_REQUEST,
			Message: "Invalid request body",
		})
		return
	}

	certificate := models.TeacherCertificate{
		TeacherID: teacherID,
		Name:      req.Name,
		FileURL:   req.FileURL,
	}

	memberships, _ := ctl.membershipRepo.GetActiveByTeacherID(ctx, teacherID)
	var centerID uint
	if len(memberships) > 0 {
		centerID = memberships[0].CenterID
	}

	ctl.auditLogRepo.Create(ctx, models.AuditLog{
		CenterID:   centerID,
		ActorType:  "TEACHER",
		ActorID:    teacherID,
		Action:     "CERTIFICATE_UPLOAD",
		TargetType: "TeacherCertificate",
		TargetID:   certificate.ID,
		Payload: models.AuditPayload{
			After: map[string]interface{}{
				"name":     req.Name,
				"file_url": req.FileURL,
			},
		},
	})

	response := resources.TeacherCertificateResource{
		ID:      certificate.ID,
		Name:    certificate.Name,
		FileURL: certificate.FileURL,
	}

	ctx.JSON(http.StatusOK, global.ApiResponse{
		Code:    0,
		Message: "Certificate uploaded",
		Datas:   response,
	})
}

type UploadCertificateRequest struct {
	Name    string `json:"name" binding:"required"`
	FileURL string `json:"file_url" binding:"required"`
}

// GetSchedule 取得老師的綜合課表（個人行程 + 各中心課程）
// @Summary 取得老師的綜合課表
// @Tags Teacher
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param from query string true "開始日期 (YYYY-MM-DD)"
// @Param to query string true "結束日期 (YYYY-MM-DD)"
// @Success 200 {object} global.ApiResponse{data=[]TeacherScheduleItem}
// @Router /api/v1/teacher/me/schedule [get]
func (ctl *TeacherController) GetSchedule(ctx *gin.Context) {
	teacherID := ctx.GetUint(global.UserIDKey)
	if teacherID == 0 {
		ctx.JSON(http.StatusUnauthorized, global.ApiResponse{
			Code:    global.UNAUTHORIZED,
			Message: "Teacher ID required",
		})
		return
	}

	fromStr := ctx.Query("from")
	toStr := ctx.Query("to")
	if fromStr == "" || toStr == "" {
		ctx.JSON(http.StatusBadRequest, global.ApiResponse{
			Code:    global.BAD_REQUEST,
			Message: "from and to dates are required",
		})
		return
	}

	fromDate, err := time.Parse("2006-01-02", fromStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, global.ApiResponse{
			Code:    global.BAD_REQUEST,
			Message: "Invalid from date format",
		})
		return
	}

	toDate, err := time.Parse("2006-01-02", toStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, global.ApiResponse{
			Code:    global.BAD_REQUEST,
			Message: "Invalid to date format",
		})
		return
	}

	memberships, err := ctl.membershipRepo.GetActiveByTeacherID(ctx, teacherID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, global.ApiResponse{
			Code:    500,
			Message: "Failed to get memberships",
		})
		return
	}

	var schedule []TeacherScheduleItem

	for _, m := range memberships {
		rules, _ := ctl.scheduleRuleRepo.ListByCenterID(ctx, m.CenterID)
		expanded := ctl.expansionService.ExpandRules(ctx, rules, fromDate, toDate, m.CenterID)

		for _, item := range expanded {
			status := "NORMAL"
			exceptions, _ := ctl.exceptionRepo.GetByRuleAndDate(ctx, item.RuleID, item.Date)
			for _, exc := range exceptions {
				if exc.Status == "PENDING" {
					status = "PENDING_" + exc.Type
				} else if exc.Status == "APPROVED" && exc.Type == "CANCEL" {
					status = "CANCELLED"
				} else if exc.Status == "APPROVED" && exc.Type == "RESCHEDULE" {
					status = "RESCHEDULED"
				}
			}

			if status != "CANCELLED" {
				schedule = append(schedule, TeacherScheduleItem{
					ID:         fmt.Sprintf("center_%d_rule_%d_%s", m.CenterID, item.RuleID, item.Date.Format("20060102")),
					Type:       "CENTER_SESSION",
					Title:      "Center Session",
					Date:       item.Date.Format("2006-01-02"),
					StartTime:  item.StartTime,
					EndTime:    item.EndTime,
					RoomID:     item.RoomID,
					TeacherID:  item.TeacherID,
					CenterID:   m.CenterID,
					CenterName: "",
					Status:     status,
				})
			}
		}
	}

	ctx.JSON(http.StatusOK, global.ApiResponse{
		Code:    0,
		Message: "Success",
		Datas:   schedule,
	})
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
func (ctl *TeacherController) CreateException(ctx *gin.Context) {
	teacherID := ctx.GetUint(global.UserIDKey)
	if teacherID == 0 {
		ctx.JSON(http.StatusUnauthorized, global.ApiResponse{
			Code:    global.UNAUTHORIZED,
			Message: "Teacher ID required",
		})
		return
	}

	var req TeacherCreateExceptionRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, global.ApiResponse{
			Code:    global.BAD_REQUEST,
			Message: "Invalid request body",
		})
		return
	}

	exception, err := ctl.exceptionService.CreateException(
		ctx,
		req.CenterID,
		teacherID,
		req.RuleID,
		req.OriginalDate,
		req.Type,
		req.NewStartAt,
		req.NewEndAt,
		req.NewTeacherID,
		req.Reason,
	)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, global.ApiResponse{
			Code:    500,
			Message: err.Error(),
		})
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

	ctx.JSON(http.StatusOK, global.ApiResponse{
		Code:    0,
		Message: "Exception submitted",
		Datas:   exception,
	})
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
func (ctl *TeacherController) RevokeException(ctx *gin.Context) {
	teacherID := ctx.GetUint(global.UserIDKey)
	if teacherID == 0 {
		ctx.JSON(http.StatusUnauthorized, global.ApiResponse{
			Code:    global.UNAUTHORIZED,
			Message: "Teacher ID required",
		})
		return
	}

	id := ctx.Param("id")
	var exceptionID uint
	if _, err := fmt.Sscanf(id, "%d", &exceptionID); err != nil {
		ctx.JSON(http.StatusBadRequest, global.ApiResponse{
			Code:    global.BAD_REQUEST,
			Message: "Invalid exception ID",
		})
		return
	}

	err := ctl.exceptionService.RevokeException(ctx, exceptionID, teacherID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, global.ApiResponse{
			Code:    500,
			Message: err.Error(),
		})
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

	ctx.JSON(http.StatusOK, global.ApiResponse{
		Code:    0,
		Message: "Exception revoked",
	})
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
func (ctl *TeacherController) GetExceptions(ctx *gin.Context) {
	teacherID := ctx.GetUint(global.UserIDKey)
	if teacherID == 0 {
		ctx.JSON(http.StatusUnauthorized, global.ApiResponse{
			Code:    global.UNAUTHORIZED,
			Message: "Teacher ID required",
		})
		return
	}

	memberships, _ := ctl.membershipRepo.GetActiveByTeacherID(ctx, teacherID)
	var centerIDs []uint
	for _, m := range memberships {
		centerIDs = append(centerIDs, m.CenterID)
	}

	var exceptions []models.ScheduleException
	query := ctl.app.MySQL.RDB.WithContext(ctx).
		Table("schedule_exceptions").
		Select("schedule_exceptions.*").
		Joins("JOIN center_memberships ON center_memberships.center_id = schedule_exceptions.center_id").
		Where("center_memberships.teacher_id = ?", teacherID).
		Where("center_memberships.status = ?", "ACTIVE")

	if status := ctx.Query("status"); status != "" {
		query = query.Where("schedule_exceptions.status = ?", status)
	}

	query.Order("schedule_exceptions.created_at DESC").Find(&exceptions)

	ctx.JSON(http.StatusOK, global.ApiResponse{
		Code:    0,
		Message: "Success",
		Datas:   exceptions,
	})
}

type TeacherScheduleItem struct {
	ID         string `json:"id"`
	Type       string `json:"type"`
	Title      string `json:"title"`
	Date       string `json:"date"`
	StartTime  string `json:"start_time"`
	EndTime    string `json:"end_time"`
	RoomID     uint   `json:"room_id"`
	TeacherID  *uint  `json:"teacher_id"`
	CenterID   uint   `json:"center_id"`
	CenterName string `json:"center_name"`
	Status     string `json:"status"`
}

type TeacherCreateExceptionRequest struct {
	CenterID     uint       `json:"center_id" binding:"required"`
	RuleID       uint       `json:"rule_id" binding:"required"`
	OriginalDate time.Time  `json:"original_date" binding:"required"`
	Type         string     `json:"type" binding:"required,oneof=CANCEL RESCHEDULE"`
	NewStartAt   *time.Time `json:"new_start_at"`
	NewEndAt     *time.Time `json:"new_end_at"`
	NewTeacherID *uint      `json:"new_teacher_id"`
	Reason       string     `json:"reason" binding:"required"`
}

// GetSessionNote 取得課堂筆記
// @Summary 取得課堂筆記
// @Tags Teacher
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param rule_id query uint true "課程規則ID"
// @Param session_date query string true "課程日期 (YYYY-MM-DD)"
// @Success 200 {object} global.ApiResponse{data=resources.SessionNoteResource}
// @Router /api/v1/teacher/sessions/note [get]
func (ctl *TeacherController) GetSessionNote(ctx *gin.Context) {
	teacherID := ctx.GetUint(global.UserIDKey)
	if teacherID == 0 {
		ctx.JSON(http.StatusUnauthorized, global.ApiResponse{
			Code:    global.UNAUTHORIZED,
			Message: "Teacher ID required",
		})
		return
	}

	ruleID, err := ctx.GetQuery("rule_id")
	if !err {
		ctx.JSON(http.StatusBadRequest, global.ApiResponse{
			Code:    global.BAD_REQUEST,
			Message: "rule_id required",
		})
		return
	}

	sessionDateStr, err := ctx.GetQuery("session_date")
	if !err {
		ctx.JSON(http.StatusBadRequest, global.ApiResponse{
			Code:    global.BAD_REQUEST,
			Message: "session_date required",
		})
		return
	}

	var ruleIDUint uint
	fmt.Sscanf(ruleID, "%d", &ruleIDUint)

	sessionDate, parseErr := time.Parse("2006-01-02", sessionDateStr)
	if parseErr != nil {
		ctx.JSON(http.StatusBadRequest, global.ApiResponse{
			Code:    global.BAD_REQUEST,
			Message: "Invalid session_date format, use YYYY-MM-DD",
		})
		return
	}

	noteRepo := repositories.NewSessionNoteRepository(ctl.app)
	note, isNew, dbErr := noteRepo.GetOrCreate(ctx, teacherID, ruleIDUint, sessionDate)
	if dbErr != nil {
		ctx.JSON(http.StatusInternalServerError, global.ApiResponse{
			Code:    500,
			Message: dbErr.Error(),
		})
		return
	}

	response := resources.SessionNoteResource{
		ID:          note.ID,
		RuleID:      note.RuleID,
		SessionDate: note.SessionDate.Format("2006-01-02"),
		Content:     note.Content,
		PrepNote:    note.PrepNote,
		UpdatedAt:   note.UpdatedAt,
	}

	ctx.JSON(http.StatusOK, global.ApiResponse{
		Code:    0,
		Message: "Success",
		Datas: map[string]interface{}{
			"note":   response,
			"is_new": isNew,
		},
	})
}

// UpsertSessionNote 新增或更新課堂筆記
// @Summary 新增或更新課堂筆記
// @Tags Teacher
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body requests.UpsertSessionNoteRequest true "筆記內容"
// @Success 200 {object} global.ApiResponse{data=resources.SessionNoteResource}
// @Router /api/v1/teacher/sessions/note [put]
func (ctl *TeacherController) UpsertSessionNote(ctx *gin.Context) {
	teacherID := ctx.GetUint(global.UserIDKey)
	if teacherID == 0 {
		ctx.JSON(http.StatusUnauthorized, global.ApiResponse{
			Code:    global.UNAUTHORIZED,
			Message: "Teacher ID required",
		})
		return
	}

	var req requests.UpsertSessionNoteRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, global.ApiResponse{
			Code:    global.BAD_REQUEST,
			Message: "Invalid request body: " + err.Error(),
		})
		return
	}

	sessionDate, parseErr := time.Parse("2006-01-02", req.SessionDate)
	if parseErr != nil {
		ctx.JSON(http.StatusBadRequest, global.ApiResponse{
			Code:    global.BAD_REQUEST,
			Message: "Invalid session_date format, use YYYY-MM-DD",
		})
		return
	}

	noteRepo := repositories.NewSessionNoteRepository(ctl.app)
	note, _, err := noteRepo.GetOrCreate(ctx, teacherID, req.RuleID, sessionDate)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, global.ApiResponse{
			Code:    500,
			Message: err.Error(),
		})
		return
	}

	note.Content = req.Content
	note.PrepNote = req.PrepNote
	note.UpdatedAt = time.Now()

	if err := noteRepo.Update(ctx, note); err != nil {
		ctx.JSON(http.StatusInternalServerError, global.ApiResponse{
			Code:    500,
			Message: err.Error(),
		})
		return
	}

	response := resources.SessionNoteResource{
		ID:          note.ID,
		RuleID:      note.RuleID,
		SessionDate: note.SessionDate.Format("2006-01-02"),
		Content:     note.Content,
		PrepNote:    note.PrepNote,
		UpdatedAt:   note.UpdatedAt,
	}

	ctx.JSON(http.StatusOK, global.ApiResponse{
		Code:    0,
		Message: "Note saved",
		Datas:   response,
	})
}

// GetSkills 取得老師技能列表
// @Summary 取得老師技能列表
// @Tags Teacher
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} global.ApiResponse{data=[]models.TeacherSkill}
// @Router /api/v1/teacher/me/skills [get]
func (ctl *TeacherController) GetSkills(ctx *gin.Context) {
	teacherID := ctx.GetUint(global.UserIDKey)
	if teacherID == 0 {
		ctx.JSON(http.StatusUnauthorized, global.ApiResponse{
			Code:    global.UNAUTHORIZED,
			Message: "Teacher ID required",
		})
		return
	}

	skills, err := ctl.skillRepo.ListByTeacherID(ctx, teacherID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, global.ApiResponse{
			Code:    500,
			Message: err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, global.ApiResponse{
		Code:    0,
		Message: "Success",
		Datas:   skills,
	})
}

// CreateSkill 新增老師技能
// @Summary 新增老師技能
// @Tags Teacher
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body CreateSkillRequest true "技能資訊"
// @Success 200 {object} global.ApiResponse{data=models.TeacherSkill}
// @Router /api/v1/teacher/me/skills [post]
func (ctl *TeacherController) CreateSkill(ctx *gin.Context) {
	teacherID := ctx.GetUint(global.UserIDKey)
	if teacherID == 0 {
		ctx.JSON(http.StatusUnauthorized, global.ApiResponse{
			Code:    global.UNAUTHORIZED,
			Message: "Teacher ID required",
		})
		return
	}

	var req CreateSkillRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, global.ApiResponse{
			Code:    global.BAD_REQUEST,
			Message: "Invalid request body: " + err.Error(),
		})
		return
	}

	skill := models.TeacherSkill{
		TeacherID: teacherID,
		Category:  req.Category,
		SkillName: req.SkillName,
		Level:     req.Level,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := ctl.skillRepo.Create(ctx, &skill); err != nil {
		ctx.JSON(http.StatusInternalServerError, global.ApiResponse{
			Code:    500,
			Message: err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, global.ApiResponse{
		Code:    0,
		Message: "Skill created",
		Datas:   skill,
	})
}

// DeleteSkill 刪除老師技能
// @Summary 刪除老師技能
// @Tags Teacher
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path uint true "技能ID"
// @Success 200 {object} global.ApiResponse
// @Router /api/v1/teacher/me/skills/{id} [delete]
func (ctl *TeacherController) DeleteSkill(ctx *gin.Context) {
	teacherID := ctx.GetUint(global.UserIDKey)
	if teacherID == 0 {
		ctx.JSON(http.StatusUnauthorized, global.ApiResponse{
			Code:    global.UNAUTHORIZED,
			Message: "Teacher ID required",
		})
		return
	}

	var id uint
	if _, err := fmt.Sscanf(ctx.Param("id"), "%d", &id); err != nil {
		ctx.JSON(http.StatusBadRequest, global.ApiResponse{
			Code:    global.BAD_REQUEST,
			Message: "Invalid skill ID",
		})
		return
	}

	skill, err := ctl.skillRepo.GetByID(ctx, id)
	if err != nil {
		ctx.JSON(http.StatusNotFound, global.ApiResponse{
			Code:    404,
			Message: "Skill not found",
		})
		return
	}

	if skill.TeacherID != teacherID {
		ctx.JSON(http.StatusForbidden, global.ApiResponse{
			Code:    403,
			Message: "Not authorized to delete this skill",
		})
		return
	}

	if err := ctl.skillRepo.Delete(ctx, id); err != nil {
		ctx.JSON(http.StatusInternalServerError, global.ApiResponse{
			Code:    500,
			Message: err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, global.ApiResponse{
		Code:    0,
		Message: "Skill deleted",
	})
}

// GetCertificates 取得老師證照列表
// @Summary 取得老師證照列表
// @Tags Teacher
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} global.ApiResponse{data=[]models.TeacherCertificate}
// @Router /api/v1/teacher/me/certificates [get]
func (ctl *TeacherController) GetCertificates(ctx *gin.Context) {
	teacherID := ctx.GetUint(global.UserIDKey)
	if teacherID == 0 {
		ctx.JSON(http.StatusUnauthorized, global.ApiResponse{
			Code:    global.UNAUTHORIZED,
			Message: "Teacher ID required",
		})
		return
	}

	certificates, err := ctl.certificateRepo.ListByTeacherID(ctx, teacherID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, global.ApiResponse{
			Code:    500,
			Message: err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, global.ApiResponse{
		Code:    0,
		Message: "Success",
		Datas:   certificates,
	})
}

// CreateCertificate 新增老師證照
// @Summary 新增老師證照
// @Tags Teacher
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body CreateCertificateRequest true "證照資訊"
// @Success 200 {object} global.ApiResponse{data=models.TeacherCertificate}
// @Router /api/v1/teacher/me/certificates [post]
func (ctl *TeacherController) CreateCertificate(ctx *gin.Context) {
	teacherID := ctx.GetUint(global.UserIDKey)
	if teacherID == 0 {
		ctx.JSON(http.StatusUnauthorized, global.ApiResponse{
			Code:    global.UNAUTHORIZED,
			Message: "Teacher ID required",
		})
		return
	}

	var req CreateCertificateRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, global.ApiResponse{
			Code:    global.BAD_REQUEST,
			Message: "Invalid request body: " + err.Error(),
		})
		return
	}

	certificate := models.TeacherCertificate{
		TeacherID: teacherID,
		Name:      req.Name,
		FileURL:   req.FileURL,
		IssuedAt:  req.IssuedAt,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := ctl.certificateRepo.Create(ctx, &certificate); err != nil {
		ctx.JSON(http.StatusInternalServerError, global.ApiResponse{
			Code:    500,
			Message: err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, global.ApiResponse{
		Code:    0,
		Message: "Certificate created",
		Datas:   certificate,
	})
}

// DeleteCertificate 刪除老師證照
// @Summary 刪除老師證照
// @Tags Teacher
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path uint true "證照ID"
// @Success 200 {object} global.ApiResponse
// @Router /api/v1/teacher/me/certificates/{id} [delete]
func (ctl *TeacherController) DeleteCertificate(ctx *gin.Context) {
	teacherID := ctx.GetUint(global.UserIDKey)
	if teacherID == 0 {
		ctx.JSON(http.StatusUnauthorized, global.ApiResponse{
			Code:    global.UNAUTHORIZED,
			Message: "Teacher ID required",
		})
		return
	}

	var id uint
	if _, err := fmt.Sscanf(ctx.Param("id"), "%d", &id); err != nil {
		ctx.JSON(http.StatusBadRequest, global.ApiResponse{
			Code:    global.BAD_REQUEST,
			Message: "Invalid certificate ID",
		})
		return
	}

	certificate, err := ctl.certificateRepo.GetByID(ctx, id)
	if err != nil {
		ctx.JSON(http.StatusNotFound, global.ApiResponse{
			Code:    404,
			Message: "Certificate not found",
		})
		return
	}

	if certificate.TeacherID != teacherID {
		ctx.JSON(http.StatusForbidden, global.ApiResponse{
			Code:    403,
			Message: "Not authorized to delete this certificate",
		})
		return
	}

	if err := ctl.certificateRepo.Delete(ctx, id); err != nil {
		ctx.JSON(http.StatusInternalServerError, global.ApiResponse{
			Code:    500,
			Message: err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, global.ApiResponse{
		Code:    0,
		Message: "Certificate deleted",
	})
}

// GetPersonalEvents 取得老師個人行程列表
// @Summary 取得老師個人行程列表
// @Tags Teacher
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} global.ApiResponse{data=[]models.PersonalEvent}
// @Router /api/v1/teacher/me/personal-events [get]
func (ctl *TeacherController) GetPersonalEvents(ctx *gin.Context) {
	teacherID := ctx.GetUint(global.UserIDKey)
	if teacherID == 0 {
		ctx.JSON(http.StatusUnauthorized, global.ApiResponse{
			Code:    global.UNAUTHORIZED,
			Message: "Teacher ID required",
		})
		return
	}

	fromStr := ctx.Query("from")
	toStr := ctx.Query("to")

	var events []models.PersonalEvent
	var err error

	if fromStr != "" && toStr != "" {
		from, err := time.Parse("2006-01-02", fromStr)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, global.ApiResponse{
				Code:    400,
				Message: "Invalid from date format",
			})
			return
		}
		to, err := time.Parse("2006-01-02", toStr)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, global.ApiResponse{
				Code:    400,
				Message: "Invalid to date format",
			})
			return
		}
		// Add one day to to date to make it inclusive
		to = to.AddDate(0, 0, 1)
		events, err = ctl.personalEventRepo.GetByTeacherAndDateRange(ctx, teacherID, from, to)
	} else {
		events, err = ctl.personalEventRepo.ListByTeacherID(ctx, teacherID)
	}

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, global.ApiResponse{
			Code:    500,
			Message: err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, global.ApiResponse{
		Code:    0,
		Message: "Success",
		Datas:   events,
	})
}

// CreatePersonalEvent 新增老師個人行程
// @Summary 新增老師個人行程
// @Tags Teacher
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body CreatePersonalEventRequest true "行程資訊"
// @Success 200 {object} global.ApiResponse{data=models.PersonalEvent}
// @Router /api/v1/teacher/me/personal-events [post]
func (ctl *TeacherController) CreatePersonalEvent(ctx *gin.Context) {
	teacherID := ctx.GetUint(global.UserIDKey)
	if teacherID == 0 {
		ctx.JSON(http.StatusUnauthorized, global.ApiResponse{
			Code:    global.UNAUTHORIZED,
			Message: "Teacher ID required",
		})
		return
	}

	var req CreatePersonalEventRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, global.ApiResponse{
			Code:    global.BAD_REQUEST,
			Message: "Invalid request body: " + err.Error(),
		})
		return
	}

	event := models.PersonalEvent{
		TeacherID: teacherID,
		Title:     req.Title,
		StartAt:   req.StartAt,
		EndAt:     req.EndAt,
		IsAllDay:  req.IsAllDay,
		ColorHex:  req.ColorHex,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if req.RecurrenceRule != nil {
		event.RecurrenceRule = *req.RecurrenceRule
	}

	if err := ctl.personalEventRepo.Create(ctx, &event); err != nil {
		ctx.JSON(http.StatusInternalServerError, global.ApiResponse{
			Code:    500,
			Message: err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, global.ApiResponse{
		Code:    0,
		Message: "Personal event created",
		Datas:   event,
	})
}

// DeletePersonalEvent 刪除老師個人行程
// @Summary 刪除老師個人行程
// @Tags Teacher
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path uint true "行程ID"
// @Success 200 {object} global.ApiResponse
// @Router /api/v1/teacher/me/personal-events/{id} [delete]
func (ctl *TeacherController) DeletePersonalEvent(ctx *gin.Context) {
	teacherID := ctx.GetUint(global.UserIDKey)
	if teacherID == 0 {
		ctx.JSON(http.StatusUnauthorized, global.ApiResponse{
			Code:    global.UNAUTHORIZED,
			Message: "Teacher ID required",
		})
		return
	}

	var id uint
	if _, err := fmt.Sscanf(ctx.Param("id"), "%d", &id); err != nil {
		ctx.JSON(http.StatusBadRequest, global.ApiResponse{
			Code:    global.BAD_REQUEST,
			Message: "Invalid personal event ID",
		})
		return
	}

	event, err := ctl.personalEventRepo.GetByID(ctx, id)
	if err != nil {
		ctx.JSON(http.StatusNotFound, global.ApiResponse{
			Code:    404,
			Message: "Personal event not found",
		})
		return
	}

	if event.TeacherID != teacherID {
		ctx.JSON(http.StatusForbidden, global.ApiResponse{
			Code:    403,
			Message: "Not authorized to delete this personal event",
		})
		return
	}

	if err := ctl.personalEventRepo.Delete(ctx, id); err != nil {
		ctx.JSON(http.StatusInternalServerError, global.ApiResponse{
			Code:    500,
			Message: err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, global.ApiResponse{
		Code:    0,
		Message: "Personal event deleted",
	})
}

type CreateSkillRequest struct {
	Category  string `json:"category" binding:"required"`
	SkillName string `json:"skill_name" binding:"required"`
	Level     string `json:"level" binding:"required,oneof=BASIC INTERMEDIATE ADVANCED"`
}

type CreateCertificateRequest struct {
	Name     string    `json:"name" binding:"required"`
	FileURL  string    `json:"file_url" binding:"required"`
	IssuedAt time.Time `json:"issued_at" binding:"required"`
}

type CreatePersonalEventRequest struct {
	Title          string                 `json:"title" binding:"required"`
	StartAt        time.Time              `json:"start_at" binding:"required"`
	EndAt          time.Time              `json:"end_at" binding:"required"`
	IsAllDay       bool                   `json:"is_all_day"`
	ColorHex       string                 `json:"color_hex"`
	RecurrenceRule *models.RecurrenceRule `json:"recurrence_rule"`
}

type UpdatePersonalEventRequest struct {
	Title          *string                `json:"title"`
	StartAt        *time.Time             `json:"start_at"`
	EndAt          *time.Time             `json:"end_at"`
	IsAllDay       *bool                  `json:"is_all_day"`
	ColorHex       *string                `json:"color_hex"`
	RecurrenceRule *models.RecurrenceRule `json:"recurrence_rule"`
	UpdateMode     UpdateMode             `json:"update_mode" binding:"required,oneof=SINGLE FUTURE ALL"`
}

type UpdateMode string

const (
	UpdateModeSingle UpdateMode = "SINGLE"
	UpdateModeFuture UpdateMode = "FUTURE"
	UpdateModeAll    UpdateMode = "ALL"
)

type UpdatePersonalEventResponse struct {
	UpdatedCount int64  `json:"updated_count"`
	Message      string `json:"message"`
}

// UpdatePersonalEvent 更新老師個人行程
// @Summary 更新老師個人行程
// @Tags Teacher
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path uint true "行程ID"
// @Param request body UpdatePersonalEventRequest true "行程資訊"
// @Success 200 {object} global.ApiResponse{data=UpdatePersonalEventResponse}
// @Router /api/v1/teacher/me/personal-events/{id} [patch]
func (ctl *TeacherController) UpdatePersonalEvent(ctx *gin.Context) {
	teacherID := ctx.GetUint(global.UserIDKey)
	if teacherID == 0 {
		ctx.JSON(http.StatusUnauthorized, global.ApiResponse{
			Code:    global.UNAUTHORIZED,
			Message: "Teacher ID required",
		})
		return
	}

	var id uint
	if _, err := fmt.Sscanf(ctx.Param("id"), "%d", &id); err != nil {
		ctx.JSON(http.StatusBadRequest, global.ApiResponse{
			Code:    global.BAD_REQUEST,
			Message: "Invalid personal event ID",
		})
		return
	}

	event, err := ctl.personalEventRepo.GetByID(ctx, id)
	if err != nil {
		ctx.JSON(http.StatusNotFound, global.ApiResponse{
			Code:    errInfos.NOT_FOUND,
			Message: ctl.app.Err.New(errInfos.NOT_FOUND).Msg,
		})
		return
	}

	if event.TeacherID != teacherID {
		ctx.JSON(http.StatusForbidden, global.ApiResponse{
			Code:    global.FORBIDDEN,
			Message: "Not authorized to update this personal event",
		})
		return
	}

	var req UpdatePersonalEventRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, global.ApiResponse{
			Code:    global.BAD_REQUEST,
			Message: "Invalid request body: " + err.Error(),
		})
		return
	}

	now := time.Now()
	var updatedCount int64 = 1

	switch req.UpdateMode {
	case UpdateModeSingle:
		if req.Title != nil {
			event.Title = *req.Title
		}
		if req.StartAt != nil {
			event.StartAt = *req.StartAt
		}
		if req.EndAt != nil {
			event.EndAt = *req.EndAt
		}
		if req.IsAllDay != nil {
			event.IsAllDay = *req.IsAllDay
		}
		if req.ColorHex != nil {
			event.ColorHex = *req.ColorHex
		}
		if req.RecurrenceRule != nil {
			event.RecurrenceRule = *req.RecurrenceRule
		}
		event.UpdatedAt = now
		if err := ctl.personalEventRepo.Update(ctx, event); err != nil {
			ctx.JSON(http.StatusInternalServerError, global.ApiResponse{
				Code:    errInfos.SQL_ERROR,
				Message: ctl.app.Err.New(errInfos.SQL_ERROR).Msg,
			})
			return
		}

	case UpdateModeFuture:
		if event.RecurrenceRule.Type == "" {
			ctx.JSON(http.StatusBadRequest, global.ApiResponse{
				Code:    errInfos.PARAMS_VALIDATE_ERROR,
				Message: "update_mode FUTURE requires a recurring event",
			})
			return
		}

		repoReq := repositories.UpdateEventRequest{
			Title:    req.Title,
			StartAt:  req.StartAt,
			EndAt:    req.EndAt,
			IsAllDay: req.IsAllDay,
			ColorHex: req.ColorHex,
		}
		updatedCount, err = ctl.personalEventRepo.UpdateFutureOccurrences(ctx, id, teacherID, repoReq, now)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, global.ApiResponse{
				Code:    errInfos.SQL_ERROR,
				Message: ctl.app.Err.New(errInfos.SQL_ERROR).Msg,
			})
			return
		}

	case UpdateModeAll:
		if event.RecurrenceRule.Type == "" {
			ctx.JSON(http.StatusBadRequest, global.ApiResponse{
				Code:    errInfos.PARAMS_VALIDATE_ERROR,
				Message: "update_mode ALL requires a recurring event",
			})
			return
		}

		repoReq := repositories.UpdateEventRequest{
			Title:    req.Title,
			StartAt:  req.StartAt,
			EndAt:    req.EndAt,
			IsAllDay: req.IsAllDay,
			ColorHex: req.ColorHex,
		}
		updatedCount, err = ctl.personalEventRepo.UpdateAllOccurrences(ctx, id, teacherID, repoReq, now)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, global.ApiResponse{
				Code:    errInfos.SQL_ERROR,
				Message: ctl.app.Err.New(errInfos.SQL_ERROR).Msg,
			})
			return
		}
	}

	ctx.JSON(http.StatusOK, global.ApiResponse{
		Code:    0,
		Message: "Personal event updated",
		Datas: UpdatePersonalEventResponse{
			UpdatedCount: updatedCount,
			Message:      "Updated " + string(req.UpdateMode) + " occurrence(s)",
		},
	})
}

// ListTeachers 取得老師列表（根據當前登入 Admin 的中心過濾）
// @Summary 取得老師列表
// @Tags Admin
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} global.ApiResponse{data=[]TeacherResponse}
// @Router /api/v1/teachers [get]
func (ctl *TeacherController) ListTeachers(ctx *gin.Context) {
	// 從 JWT Token 取得 center_id
	centerID := ctx.GetUint(global.CenterIDKey)
	if centerID == 0 {
		ctx.JSON(http.StatusUnauthorized, global.ApiResponse{
			Code:    global.UNAUTHORIZED,
			Message: "Center ID required",
		})
		return
	}

	// 取得該中心的所有會員老師 ID（包含 ACTIVE 和 INVITED 狀態）
	teacherIDs, err := ctl.membershipRepo.ListTeacherIDsByCenterID(ctx, centerID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, global.ApiResponse{
			Code:    500,
			Message: "Failed to get teacher IDs",
		})
		return
	}

	if len(teacherIDs) == 0 {
		ctx.JSON(http.StatusOK, global.ApiResponse{
			Code:    0,
			Message: "Success",
			Datas:   []TeacherResponse{},
		})
		return
	}

	// 取得老師詳細資料
	teachers := make([]TeacherResponse, 0, len(teacherIDs))
	for _, teacherID := range teacherIDs {
		teacher, err := ctl.teacherRepository.GetByID(ctx, teacherID)
		if err != nil {
			continue
		}

		teachers = append(teachers, TeacherResponse{
			ID:        teacher.ID,
			Name:      teacher.Name,
			Email:     teacher.Email,
			City:      teacher.City,
			District:  teacher.District,
			Bio:       teacher.Bio,
			CreatedAt: teacher.CreatedAt,
		})
	}

	ctx.JSON(http.StatusOK, global.ApiResponse{
		Code:    0,
		Message: "Success",
		Datas:   teachers,
	})
}

// DeleteTeacher 刪除老師
// @Summary 刪除老師
// @Tags Admin
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Teacher ID"
// @Success 200 {object} global.ApiResponse
// @Router /api/v1/teachers/{id} [delete]
func (ctl *TeacherController) DeleteTeacher(ctx *gin.Context) {
	adminID := ctx.GetUint(global.UserIDKey)
	if adminID == 0 {
		ctx.JSON(http.StatusUnauthorized, global.ApiResponse{
			Code:    global.UNAUTHORIZED,
			Message: "Admin ID required",
		})
		return
	}

	var teacherID uint
	if _, err := fmt.Sscanf(ctx.Param("id"), "%d", &teacherID); err != nil {
		ctx.JSON(http.StatusBadRequest, global.ApiResponse{
			Code:    global.BAD_REQUEST,
			Message: "Invalid teacher ID",
		})
		return
	}

	teacher, err := ctl.teacherRepository.GetByID(ctx, teacherID)
	if err != nil {
		ctx.JSON(http.StatusNotFound, global.ApiResponse{
			Code:    404,
			Message: "Teacher not found",
		})
		return
	}

	if err := ctl.teacherRepository.DeleteByID(ctx, teacherID); err != nil {
		ctx.JSON(http.StatusInternalServerError, global.ApiResponse{
			Code:    500,
			Message: "Failed to delete teacher",
		})
		return
	}

	memberships, _ := ctl.membershipRepo.GetActiveByTeacherID(ctx, teacherID)
	var centerID uint
	if len(memberships) > 0 {
		centerID = memberships[0].CenterID
	}

	ctl.auditLogRepo.Create(ctx, models.AuditLog{
		CenterID:   centerID,
		ActorType:  "ADMIN",
		ActorID:    adminID,
		Action:     "DELETE_TEACHER",
		TargetType: "Teacher",
		TargetID:   teacherID,
		Payload: models.AuditPayload{
			After: map[string]interface{}{
				"status": "DELETED",
				"name":   teacher.Name,
				"email":  teacher.Email,
			},
		},
	})

	ctx.JSON(http.StatusOK, global.ApiResponse{
		Code:    0,
		Message: "Teacher deleted",
	})
}

// InviteTeacher 邀請老師加入中心
// @Summary 邀請老師加入中心
// @Tags Admin
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Center ID"
// @Param request body InviteTeacherRequest true "邀請資訊"
// @Success 200 {object} global.ApiResponse{data=models.CenterInvitation}
// @Router /api/v1/admin/centers/{id}/invitations [post]
func (ctl *TeacherController) InviteTeacher(ctx *gin.Context) {
	adminID := ctx.GetUint(global.UserIDKey)
	if adminID == 0 {
		ctx.JSON(http.StatusUnauthorized, global.ApiResponse{
			Code:    global.UNAUTHORIZED,
			Message: "Admin ID required",
		})
		return
	}

	centerIDStr := ctx.Param("id")
	if centerIDStr == "" {
		ctx.JSON(http.StatusBadRequest, global.ApiResponse{
			Code:    global.BAD_REQUEST,
			Message: "Center ID required",
		})
		return
	}

	var centerID uint
	if _, err := fmt.Sscanf(centerIDStr, "%d", &centerID); err != nil {
		ctx.JSON(http.StatusBadRequest, global.ApiResponse{
			Code:    global.BAD_REQUEST,
			Message: "Invalid center ID",
		})
		return
	}

	var req InviteTeacherRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, global.ApiResponse{
			Code:    global.BAD_REQUEST,
			Message: "Invalid request body: " + err.Error(),
		})
		return
	}

	token := generateInviteToken()
	expiresAt := time.Now().Add(72 * time.Hour)

	invitation := models.CenterInvitation{
		CenterID:  centerID,
		Email:     req.Email,
		Token:     token,
		Status:    "PENDING",
		CreatedAt: time.Now(),
		ExpiresAt: expiresAt,
	}

	if err := ctl.app.MySQL.WDB.WithContext(ctx).Create(&invitation).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, global.ApiResponse{
			Code:    500,
			Message: "Failed to create invitation",
		})
		return
	}

	ctl.auditLogRepo.Create(ctx, models.AuditLog{
		CenterID:   centerID,
		ActorType:  "ADMIN",
		ActorID:    adminID,
		Action:     "INVITE_TEACHER",
		TargetType: "CenterInvitation",
		TargetID:   invitation.ID,
		Payload: models.AuditPayload{
			After: map[string]interface{}{
				"email":      req.Email,
				"status":     "PENDING",
				"expires_at": expiresAt,
			},
		},
	})

	ctx.JSON(http.StatusOK, global.ApiResponse{
		Code:    0,
		Message: "Invitation sent",
		Datas:   invitation,
	})
}

type InviteTeacherRequest struct {
	Email   string `json:"email" binding:"required,email"`
	Role    string `json:"role" binding:"required,oneof=TEACHER SUBSTITUTE"`
	Message string `json:"message"`
}

func generateInviteToken() string {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return ""
	}
	return fmt.Sprintf("%x", b)
}

type CheckTeacherRuleLockRequest struct {
	RuleID        uint   `json:"rule_id" binding:"required"`
	ExceptionDate string `json:"exception_date" binding:"required"`
}

type CheckTeacherRuleLockResponse struct {
	IsLocked      bool       `json:"is_locked"`
	LockReason    string     `json:"lock_reason,omitempty"`
	Deadline      *time.Time `json:"deadline,omitempty"`
	DaysRemaining int        `json:"days_remaining"`
}

func (ctl *TeacherController) CheckRuleLockStatus(ctx *gin.Context) {
	teacherID := ctx.GetUint(global.UserIDKey)
	if teacherID == 0 {
		ctx.JSON(http.StatusUnauthorized, global.ApiResponse{
			Code:    errInfos.UNAUTHORIZED,
			Message: "Teacher ID not found",
		})
		return
	}

	var req CheckTeacherRuleLockRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, global.ApiResponse{
			Code:    errInfos.PARAMS_VALIDATE_ERROR,
			Message: "Invalid request parameters",
		})
		return
	}

	rule, err := ctl.scheduleRuleRepo.GetByID(ctx, req.RuleID)
	if err != nil {
		ctx.JSON(http.StatusNotFound, global.ApiResponse{
			Code:    errInfos.NOT_FOUND,
			Message: "Rule not found",
		})
		return
	}

	exceptionDate, err := time.Parse("2006-01-02", req.ExceptionDate)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, global.ApiResponse{
			Code:    errInfos.PARAMS_VALIDATE_ERROR,
			Message: "Invalid date format, expected YYYY-MM-DD",
		})
		return
	}

	allowed, reasonStr, _ := ctl.exceptionService.CheckExceptionDeadline(ctx, rule.CenterID, req.RuleID, exceptionDate)

	response := CheckTeacherRuleLockResponse{
		IsLocked:   !allowed,
		LockReason: reasonStr,
	}

	if !allowed {
		ctx.JSON(http.StatusOK, global.ApiResponse{
			Code:    0,
			Message: "Rule is locked",
			Datas:   response,
		})
		return
	}

	ctx.JSON(http.StatusOK, global.ApiResponse{
		Code:    0,
		Message: "Rule is available for exception",
		Datas:   response,
	})
}

type PreviewRecurrenceEditRequest struct {
	RuleID   uint   `json:"rule_id" binding:"required"`
	EditDate string `json:"edit_date" binding:"required"`
	Mode     string `json:"mode" binding:"required,oneof=SINGLE FUTURE ALL"`
}

func (ctl *TeacherController) PreviewRecurrenceEdit(ctx *gin.Context) {
	teacherID := ctx.GetUint(global.UserIDKey)
	if teacherID == 0 {
		ctx.JSON(http.StatusUnauthorized, global.ApiResponse{
			Code:    errInfos.UNAUTHORIZED,
			Message: "Teacher ID not found",
		})
		return
	}

	var req PreviewRecurrenceEditRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, global.ApiResponse{
			Code:    errInfos.PARAMS_VALIDATE_ERROR,
			Message: "Invalid request parameters",
		})
		return
	}

	editDate, err := time.Parse("2006-01-02", req.EditDate)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, global.ApiResponse{
			Code:    errInfos.PARAMS_VALIDATE_ERROR,
			Message: "Invalid date format, expected YYYY-MM-DD",
		})
		return
	}

	preview, err := ctl.recurrenceService.PreviewAffectedSessions(ctx, req.RuleID, editDate, services.RecurrenceEditMode(req.Mode))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, global.ApiResponse{
			Code:    errInfos.SYSTEM_ERROR,
			Message: "Failed to preview affected sessions",
		})
		return
	}

	ctx.JSON(http.StatusOK, global.ApiResponse{
		Code:    0,
		Message: "Preview generated",
		Datas:   preview,
	})
}

func (ctl *TeacherController) EditRecurringSchedule(ctx *gin.Context) {
	teacherID := ctx.GetUint(global.UserIDKey)
	if teacherID == 0 {
		ctx.JSON(http.StatusUnauthorized, global.ApiResponse{
			Code:    errInfos.UNAUTHORIZED,
			Message: "Teacher ID not found",
		})
		return
	}

	var req services.RecurrenceEditRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, global.ApiResponse{
			Code:    errInfos.PARAMS_VALIDATE_ERROR,
			Message: "Invalid request parameters",
		})
		return
	}

	rule, err := ctl.scheduleRuleRepo.GetByID(ctx, req.RuleID)
	if err != nil {
		ctx.JSON(http.StatusNotFound, global.ApiResponse{
			Code:    errInfos.NOT_FOUND,
			Message: "Rule not found",
		})
		return
	}

	result, err := ctl.recurrenceService.EditRecurringSchedule(ctx, rule.CenterID, teacherID, req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, global.ApiResponse{
			Code:    errInfos.SYSTEM_ERROR,
			Message: err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, global.ApiResponse{
		Code:    0,
		Message: "Schedule edited successfully",
		Datas:   result,
	})
}

type DeleteRecurringScheduleRequest struct {
	RuleID   uint   `json:"rule_id" binding:"required"`
	EditDate string `json:"edit_date" binding:"required"`
	Mode     string `json:"mode" binding:"required,oneof=SINGLE FUTURE ALL"`
	Reason   string `json:"reason"`
}

func (ctl *TeacherController) DeleteRecurringSchedule(ctx *gin.Context) {
	teacherID := ctx.GetUint(global.UserIDKey)
	if teacherID == 0 {
		ctx.JSON(http.StatusUnauthorized, global.ApiResponse{
			Code:    errInfos.UNAUTHORIZED,
			Message: "Teacher ID not found",
		})
		return
	}

	var req DeleteRecurringScheduleRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, global.ApiResponse{
			Code:    errInfos.PARAMS_VALIDATE_ERROR,
			Message: "Invalid request parameters",
		})
		return
	}

	editDate, err := time.Parse("2006-01-02", req.EditDate)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, global.ApiResponse{
			Code:    errInfos.PARAMS_VALIDATE_ERROR,
			Message: "Invalid date format, expected YYYY-MM-DD",
		})
		return
	}

	rule, err := ctl.scheduleRuleRepo.GetByID(ctx, req.RuleID)
	if err != nil {
		ctx.JSON(http.StatusNotFound, global.ApiResponse{
			Code:    errInfos.NOT_FOUND,
			Message: "Rule not found",
		})
		return
	}

	result, err := ctl.recurrenceService.DeleteRecurringSchedule(ctx, rule.CenterID, teacherID, req.RuleID, editDate, services.RecurrenceEditMode(req.Mode), req.Reason)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, global.ApiResponse{
			Code:    errInfos.SYSTEM_ERROR,
			Message: err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, global.ApiResponse{
		Code:    0,
		Message: "Schedule deleted successfully",
		Datas:   result,
	})
}
