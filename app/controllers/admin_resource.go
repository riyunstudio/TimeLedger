package controllers

import (
	"fmt"
	"net/http"
	"time"
	"timeLedger/app"
	"timeLedger/app/models"
	"timeLedger/app/repositories"
	"timeLedger/global"
	"timeLedger/global/errInfos"

	"github.com/gin-gonic/gin"
)

type AdminResourceController struct {
	BaseController
	app                   *app.App
	offeringRepository    *repositories.OfferingRepository
	holidayRepository     *repositories.CenterHolidayRepository
	invitationRepo        *repositories.CenterInvitationRepository
	auditLogRepo          *repositories.AuditLogRepository
	centerTeacherNoteRepo *repositories.CenterTeacherNoteRepository
	membershipRepo        *repositories.CenterMembershipRepository
	teacherRepository     *repositories.TeacherRepository
	skillRepository       *repositories.TeacherSkillRepository
	certificateRepo       *repositories.TeacherCertificateRepository
}

func NewAdminResourceController(app *app.App) *AdminResourceController {
	return &AdminResourceController{
		app:                   app,
		offeringRepository:    repositories.NewOfferingRepository(app),
		holidayRepository:     repositories.NewCenterHolidayRepository(app),
		invitationRepo:        repositories.NewCenterInvitationRepository(app),
		auditLogRepo:          repositories.NewAuditLogRepository(app),
		centerTeacherNoteRepo: repositories.NewCenterTeacherNoteRepository(app),
		membershipRepo:        repositories.NewCenterMembershipRepository(app),
		teacherRepository:     repositories.NewTeacherRepository(app),
		skillRepository:       repositories.NewTeacherSkillRepository(app),
		certificateRepo:       repositories.NewTeacherCertificateRepository(app),
	}
}

type InvitationStatsResponse struct {
	Total         int64 `json:"total"`
	Pending       int64 `json:"pending"`
	Accepted      int64 `json:"accepted"`
	Expired       int64 `json:"expired"`
	Rejected      int64 `json:"rejected"`
	RecentPending int64 `json:"recent_pending"`
}

func (ctl *AdminResourceController) GetInvitationStats(ctx *gin.Context) {
	centerID := ctl.getCenterID(ctx)
	if centerID == 0 {
		ctl.respondError(ctx, global.BAD_REQUEST, "Center ID required")
		return
	}

	now := time.Now()
	thirtyDaysAgo := now.AddDate(0, 0, -30)

	total, _ := ctl.invitationRepo.CountByCenterID(ctx, centerID)
	pending, _ := ctl.invitationRepo.CountByStatus(ctx, centerID, "PENDING")
	accepted, _ := ctl.invitationRepo.CountByStatus(ctx, centerID, "ACCEPTED")
	expired, _ := ctl.invitationRepo.CountByStatus(ctx, centerID, "EXPIRED")
	rejected, _ := ctl.invitationRepo.CountByStatus(ctx, centerID, "REJECTED")
	recentPending, _ := ctl.invitationRepo.CountByDateRange(ctx, centerID, thirtyDaysAgo, now)

	stats := InvitationStatsResponse{
		Total:         total,
		Pending:       pending,
		Accepted:      accepted,
		Expired:       expired,
		Rejected:      rejected,
		RecentPending: recentPending,
	}

	ctx.JSON(http.StatusOK, global.ApiResponse{
		Code:    0,
		Message: "Success",
		Datas:   stats,
	})
}

func (ctl *AdminResourceController) GetInvitations(ctx *gin.Context) {
	centerID := ctl.getCenterID(ctx)
	if centerID == 0 {
		ctl.respondError(ctx, global.BAD_REQUEST, "Center ID required")
		return
	}

	var req PaginationRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		req.Page = 1
		req.Limit = 20
	}

	status := ctx.Query("status")

	invitations, total, err := ctl.invitationRepo.ListByCenterIDPaginated(ctx, centerID, int(req.Page), int(req.Limit), status)
	if err != nil {
		ctl.respondError(ctx, errInfos.SQL_ERROR, "Failed to get invitations")
		return
	}

	ctx.JSON(http.StatusOK, global.ApiResponse{
		Code:    0,
		Message: "Success",
		Datas: PaginationResponse{
			Data:       invitations,
			Total:      total,
			Page:       req.Page,
			Limit:      req.Limit,
			TotalPages: (total + int64(req.Limit) - 1) / int64(req.Limit),
		},
	})
}

type PaginationRequest struct {
	Page  int `form:"page"`
	Limit int `form:"limit"`
}

type PaginationResponse struct {
	Data       interface{} `json:"data"`
	Total      int64       `json:"total"`
	Page       int         `json:"page"`
	Limit      int         `json:"limit"`
	TotalPages int64       `json:"total_pages"`
}

func (ctl *AdminResourceController) getCenterID(ctx *gin.Context) uint {
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
	return centerID
}

func (ctl *AdminResourceController) getUintParam(ctx *gin.Context, paramName string) (uint, error) {
	paramStr := ctx.Param(paramName)
	if paramStr == "" {
		return 0, fmt.Errorf("parameter not found")
	}
	var result uint
	_, err := fmt.Sscanf(paramStr, "%d", &result)
	return result, err
}

func (ctl *AdminResourceController) respondError(ctx *gin.Context, code errInfos.ErrCode, message string) {
	ctx.JSON(http.StatusBadRequest, global.ApiResponse{
		Code:    code,
		Message: message,
	})
}

// GetTeacherNote 取得老師評分與備註
// @Summary 取得老師評分與備註
// @Tags Admin
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param teacher_id path int true "Teacher ID"
// @Success 200 {object} global.ApiResponse{data=CenterTeacherNoteResponse}
// @Router /api/v1/admin/teachers/{teacher_id}/note [get]
func (ctl *AdminResourceController) GetTeacherNote(ctx *gin.Context) {
	centerID := ctl.getCenterID(ctx)
	teacherID, err := ctl.getUintParam(ctx, "teacher_id")
	if err != nil {
		ctl.respondError(ctx, global.BAD_REQUEST, "Invalid teacher ID")
		return
	}

	note, err := ctl.centerTeacherNoteRepo.GetByCenterAndTeacher(ctx, centerID, teacherID)
	if err != nil {
		// 如果沒有找到，回傳空的評分記錄
		if err.Error() == "record not found" {
			ctx.JSON(http.StatusOK, global.ApiResponse{
				Code:    0,
				Message: "Success",
				Datas: CenterTeacherNoteResponse{
					TeacherID:    teacherID,
					Rating:       0,
					InternalNote: "",
				},
			})
			return
		}
		ctl.respondError(ctx, errInfos.SQL_ERROR, "Failed to get teacher note")
		return
	}

	ctx.JSON(http.StatusOK, global.ApiResponse{
		Code:    0,
		Message: "Success",
		Datas: CenterTeacherNoteResponse{
			ID:           note.ID,
			TeacherID:    note.TeacherID,
			Rating:       note.Rating,
			InternalNote: note.InternalNote,
			CreatedAt:    note.CreatedAt,
			UpdatedAt:    note.UpdatedAt,
		},
	})
}

// UpsertTeacherNote 新增或更新老師評分與備註
// @Summary 新增或更新老師評分與備註
// @Tags Admin
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param teacher_id path int true "Teacher ID"
// @Param request body UpsertTeacherNoteRequest true "評分與備註"
// @Success 200 {object} global.ApiResponse{data=CenterTeacherNoteResponse}
// @Router /api/v1/admin/teachers/{teacher_id}/note [put]
func (ctl *AdminResourceController) UpsertTeacherNote(ctx *gin.Context) {
	centerID := ctl.getCenterID(ctx)
	teacherID, err := ctl.getUintParam(ctx, "teacher_id")
	if err != nil {
		ctl.respondError(ctx, global.BAD_REQUEST, "Invalid teacher ID")
		return
	}

	var req UpsertTeacherNoteRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctl.respondError(ctx, errInfos.PARAMS_VALIDATE_ERROR, "Invalid request body")
		return
	}

	// 驗證評分範圍
	if req.Rating < 0 || req.Rating > 5 {
		ctl.respondError(ctx, errInfos.PARAMS_VALIDATE_ERROR, "Rating must be between 0 and 5")
		return
	}

	// 檢查是否已存在評分記錄
	existingNote, err := ctl.centerTeacherNoteRepo.GetByCenterAndTeacher(ctx, centerID, teacherID)
	if err != nil && err.Error() != "record not found" {
		ctl.respondError(ctx, errInfos.SQL_ERROR, "Failed to check existing note")
		return
	}

	now := time.Now()
	if existingNote.ID != 0 {
		// 更新現有記錄
		existingNote.Rating = req.Rating
		existingNote.InternalNote = req.InternalNote
		existingNote.UpdatedAt = now

		if err := ctl.centerTeacherNoteRepo.Update(ctx, existingNote); err != nil {
			ctl.respondError(ctx, errInfos.SQL_ERROR, "Failed to update teacher note")
			return
		}

		adminID := ctx.GetUint(global.UserIDKey)
		ctl.auditLogRepo.Create(ctx, models.AuditLog{
			CenterID:   centerID,
			ActorType:  "ADMIN",
			ActorID:    adminID,
			Action:     "UPDATE_TEACHER_NOTE",
			TargetType: "CenterTeacherNote",
			TargetID:   existingNote.ID,
			Payload: models.AuditPayload{
				Before: map[string]interface{}{
					"rating": existingNote.Rating,
					"note":   existingNote.InternalNote,
				},
				After: map[string]interface{}{
					"rating": req.Rating,
					"note":   req.InternalNote,
				},
			},
		})

		ctx.JSON(http.StatusOK, global.ApiResponse{
			Code:    0,
			Message: "Teacher note updated",
			Datas: CenterTeacherNoteResponse{
				ID:           existingNote.ID,
				TeacherID:    existingNote.TeacherID,
				Rating:       existingNote.Rating,
				InternalNote: existingNote.InternalNote,
				UpdatedAt:    existingNote.UpdatedAt,
			},
		})
		return
	}

	// 建立新記錄
	newNote := models.CenterTeacherNote{
		CenterID:     centerID,
		TeacherID:    teacherID,
		Rating:       req.Rating,
		InternalNote: req.InternalNote,
		CreatedAt:    now,
		UpdatedAt:    now,
	}

	_, createErr := ctl.centerTeacherNoteRepo.Create(ctx, newNote)
	if createErr != nil {
		ctl.respondError(ctx, errInfos.SQL_ERROR, "Failed to create teacher note")
		return
	}

	adminID := ctx.GetUint(global.UserIDKey)
	ctl.auditLogRepo.Create(ctx, models.AuditLog{
		CenterID:   centerID,
		ActorType:  "ADMIN",
		ActorID:    adminID,
		Action:     "CREATE_TEACHER_NOTE",
		TargetType: "CenterTeacherNote",
		TargetID:   newNote.ID,
		Payload: models.AuditPayload{
			After: map[string]interface{}{
				"teacher_id": teacherID,
				"rating":     req.Rating,
				"note":       req.InternalNote,
			},
		},
	})

	ctx.JSON(http.StatusOK, global.ApiResponse{
		Code:    0,
		Message: "Teacher note created",
		Datas: CenterTeacherNoteResponse{
			ID:           newNote.ID,
			TeacherID:    newNote.TeacherID,
			Rating:       newNote.Rating,
			InternalNote: newNote.InternalNote,
			CreatedAt:    newNote.CreatedAt,
			UpdatedAt:    newNote.UpdatedAt,
		},
	})
}

// DeleteTeacherNote 刪除老師評分與備註
// @Summary 刪除老師評分與備註
// @Tags Admin
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param teacher_id path int true "Teacher ID"
// @Success 200 {object} global.ApiResponse
// @Router /api/v1/admin/teachers/{teacher_id}/note [delete]
func (ctl *AdminResourceController) DeleteTeacherNote(ctx *gin.Context) {
	centerID := ctl.getCenterID(ctx)
	teacherID, err := ctl.getUintParam(ctx, "teacher_id")
	if err != nil {
		ctl.respondError(ctx, global.BAD_REQUEST, "Invalid teacher ID")
		return
	}

	note, err := ctl.centerTeacherNoteRepo.GetByCenterAndTeacher(ctx, centerID, teacherID)
	if err != nil {
		if err.Error() == "record not found" {
			ctl.respondError(ctx, errInfos.NOT_FOUND, "Teacher note not found")
			return
		}
		ctl.respondError(ctx, errInfos.SQL_ERROR, "Failed to get teacher note")
		return
	}

	if err := ctl.centerTeacherNoteRepo.DeleteByID(ctx, note.ID); err != nil {
		ctl.respondError(ctx, errInfos.SQL_ERROR, "Failed to delete teacher note")
		return
	}

	adminID := ctx.GetUint(global.UserIDKey)
	ctl.auditLogRepo.Create(ctx, models.AuditLog{
		CenterID:   centerID,
		ActorType:  "ADMIN",
		ActorID:    adminID,
		Action:     "DELETE_TEACHER_NOTE",
		TargetType: "CenterTeacherNote",
		TargetID:   note.ID,
		Payload: models.AuditPayload{
			Before: map[string]interface{}{
				"teacher_id": teacherID,
				"rating":     note.Rating,
				"note":       note.InternalNote,
			},
		},
	})

	ctx.JSON(http.StatusOK, global.ApiResponse{
		Code:    0,
		Message: "Teacher note deleted",
	})
}

type CenterTeacherNoteResponse struct {
	ID           uint      `json:"id"`
	TeacherID    uint      `json:"teacher_id"`
	Rating       int       `json:"rating"`
	InternalNote string    `json:"internal_note"`
	CreatedAt    time.Time `json:"created_at,omitempty"`
	UpdatedAt    time.Time `json:"updated_at,omitempty"`
}

type UpsertTeacherNoteRequest struct {
	Rating       int    `json:"rating" binding:"required,min=0,max=5"`
	InternalNote string `json:"internal_note"`
}

// GetTeachers 取得中心的老師列表
// @Summary 取得中心的老師列表
// @Tags Admin
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} global.ApiResponse{data=[]TeacherResponse}
// @Router /api/v1/admin/teachers [get]
func (ctl *AdminResourceController) GetTeachers(ctx *gin.Context) {
	// 從 JWT Token 取得 center_id
	centerID := ctl.getCenterID(ctx)
	if centerID == 0 {
		ctl.respondError(ctx, global.FORBIDDEN, "Center ID is required")
		return
	}

	// 取得該中心的所有活躍會員的老師 ID
	teacherIDs, err := ctl.membershipRepo.ListTeacherIDsByCenterID(ctx, centerID)
	if err != nil {
		ctl.respondError(ctx, errInfos.SQL_ERROR, "Failed to get teacher IDs")
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

		// 取得技能
		skills, _ := ctl.skillRepository.ListByTeacherID(ctx, teacherID)
		var skillResponses []TeacherSkillResponse
		for _, skill := range skills {
			skillResponses = append(skillResponses, TeacherSkillResponse{
				ID:        skill.ID,
				SkillName: skill.SkillName,
				Category:  skill.Category,
				Level:     skill.Level,
			})
		}

		// 取得證照
		certificates, _ := ctl.certificateRepo.ListByTeacherID(ctx, teacherID)
		var certResponses []CertificateResponse
		for _, cert := range certificates {
			certResponses = append(certResponses, CertificateResponse{
				ID:        cert.ID,
				Name:      cert.Name,
				FileURL:   cert.FileURL,
				IssuedAt:  cert.IssuedAt,
				CreatedAt: cert.CreatedAt,
			})
		}

		teachers = append(teachers, TeacherResponse{
			ID:           teacher.ID,
			Name:         teacher.Name,
			Email:        teacher.Email,
			City:         teacher.City,
			District:     teacher.District,
			Bio:          teacher.Bio,
			IsActive:     teacher.IsOpenToHiring, // 使用 IsOpenToHiring 作為活躍狀態
			CreatedAt:    teacher.CreatedAt,
			Skills:       skillResponses,
			Certificates: certResponses,
		})
	}

	ctx.JSON(http.StatusOK, global.ApiResponse{
		Code:    0,
		Message: "Success",
		Datas:   teachers,
	})
}

type TeacherResponse struct {
	ID           uint                   `json:"id"`
	Name         string                 `json:"name"`
	Email        string                 `json:"email"`
	Phone        string                 `json:"phone,omitempty"`
	City         string                 `json:"city,omitempty"`
	District     string                 `json:"district,omitempty"`
	Bio          string                 `json:"bio,omitempty"`
	IsActive     bool                   `json:"is_active"`
	CreatedAt    time.Time              `json:"created_at"`
	Skills       []TeacherSkillResponse `json:"skills,omitempty"`
	Certificates []CertificateResponse  `json:"certificates,omitempty"`
}

type TeacherSkillResponse struct {
	ID        uint   `json:"id"`
	SkillName string `json:"skill_name"`
	Category  string `json:"category,omitempty"`
	Level     string `json:"level,omitempty"`
}

type CertificateResponse struct {
	ID        uint      `json:"id"`
	Name      string    `json:"name"`
	FileURL   string    `json:"file_url,omitempty"`
	IssuedAt  time.Time `json:"issued_at"`
	CreatedAt time.Time `json:"created_at"`
}
