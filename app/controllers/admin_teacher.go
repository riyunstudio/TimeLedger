package controllers

import (
	"time"
	"timeLedger/app"
	"timeLedger/app/models"
	"timeLedger/app/repositories"

	"github.com/gin-gonic/gin"
)

// AdminTeacherController 管理員對老師操作相關 API
type AdminTeacherController struct {
	BaseController
	app                   *app.App
	teacherRepository     *repositories.TeacherRepository
	membershipRepo        *repositories.CenterMembershipRepository
	auditLogRepo          *repositories.AuditLogRepository
	centerTeacherNoteRepo *repositories.CenterTeacherNoteRepository
}

func NewAdminTeacherController(app *app.App) *AdminTeacherController {
	return &AdminTeacherController{
		app:                   app,
		teacherRepository:     repositories.NewTeacherRepository(app),
		membershipRepo:        repositories.NewCenterMembershipRepository(app),
		auditLogRepo:          repositories.NewAuditLogRepository(app),
		centerTeacherNoteRepo: repositories.NewCenterTeacherNoteRepository(app),
	}
}

// ListTeachers 取得老師列表（根據當前登入 Admin 的中心過濾）
// @Summary 取得老師列表
// @Tags Admin
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} global.ApiResponse{data=[]TeacherResponse}
// @Router /api/v1/teachers [get]
func (ctl *AdminTeacherController) ListTeachers(ctx *gin.Context) {
	helper := NewContextHelper(ctx)

	centerID := helper.MustCenterID()
	if centerID == 0 {
		return
	}

	// 取得該中心的所有會員老師 ID（包含 ACTIVE 和 INVITED 狀態）
	teacherIDs, err := ctl.membershipRepo.ListTeacherIDsByCenterID(ctx, centerID)
	if err != nil {
		helper.InternalError("Failed to get teacher IDs")
		return
	}

	if len(teacherIDs) == 0 {
		helper.Success([]TeacherResponse{})
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

	helper.Success(teachers)
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
func (ctl *AdminTeacherController) DeleteTeacher(ctx *gin.Context) {
	helper := NewContextHelper(ctx)

	adminID := helper.MustUserID()
	if adminID == 0 {
		return
	}

	teacherID := helper.MustParamUint("id")
	if teacherID == 0 {
		return
	}

	teacher, err := ctl.teacherRepository.GetByID(ctx, teacherID)
	if err != nil {
		helper.NotFound("Teacher not found")
		return
	}

	if err := ctl.teacherRepository.DeleteByID(ctx, teacherID); err != nil {
		helper.InternalError("Failed to delete teacher")
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

	helper.Success(nil)
}

// CenterTeacherNoteResponse 老師評分與備註回應結構
type CenterTeacherNoteResponse struct {
	ID           uint      `json:"id"`
	TeacherID    uint      `json:"teacher_id"`
	Rating       int       `json:"rating"`
	InternalNote string    `json:"internal_note"`
	CreatedAt    time.Time `json:"created_at,omitempty"`
	UpdatedAt    time.Time `json:"updated_at,omitempty"`
}

// UpsertTeacherNoteRequest 新增或更新老師評分與備註請求結構
type UpsertTeacherNoteRequest struct {
	Rating       int    `json:"rating" binding:"required,min=0,max=5"`
	InternalNote string `json:"internal_note"`
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
func (ctl *AdminTeacherController) GetTeacherNote(ctx *gin.Context) {
	helper := NewContextHelper(ctx)

	centerID := helper.MustCenterID()
	if centerID == 0 {
		return
	}

	teacherID := helper.MustParamUint("teacher_id")
	if teacherID == 0 {
		helper.BadRequest("Invalid teacher ID")
		return
	}

	note, err := ctl.centerTeacherNoteRepo.GetByCenterAndTeacher(ctx, centerID, teacherID)
	if err != nil {
		if err.Error() == "record not found" {
			helper.Success(CenterTeacherNoteResponse{
				TeacherID:    teacherID,
				Rating:       0,
				InternalNote: "",
			})
			return
		}
		helper.InternalError("Failed to get teacher note")
		return
	}

	helper.Success(CenterTeacherNoteResponse{
		ID:           note.ID,
		TeacherID:    note.TeacherID,
		Rating:       note.Rating,
		InternalNote: note.InternalNote,
		CreatedAt:    note.CreatedAt,
		UpdatedAt:    note.UpdatedAt,
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
func (ctl *AdminTeacherController) UpsertTeacherNote(ctx *gin.Context) {
	helper := NewContextHelper(ctx)

	centerID := helper.MustCenterID()
	if centerID == 0 {
		return
	}

	teacherID := helper.MustParamUint("teacher_id")
	if teacherID == 0 {
		helper.BadRequest("Invalid teacher ID")
		return
	}

	var req UpsertTeacherNoteRequest
	if !helper.MustBindJSON(&req) {
		return
	}

	// 驗證評分範圍
	if req.Rating < 0 || req.Rating > 5 {
		helper.BadRequest("Rating must be between 0 and 5")
		return
	}

	// 檢查是否已存在評分記錄
	existingNote, err := ctl.centerTeacherNoteRepo.GetByCenterAndTeacher(ctx, centerID, teacherID)
	if err != nil && err.Error() != "record not found" {
		helper.InternalError("Failed to check existing note")
		return
	}

	now := time.Now()
	if existingNote.ID != 0 {
		// 更新現有記錄
		existingNote.Rating = req.Rating
		existingNote.InternalNote = req.InternalNote
		existingNote.UpdatedAt = now

		if err := ctl.centerTeacherNoteRepo.Update(ctx, existingNote); err != nil {
			helper.InternalError("Failed to update teacher note")
			return
		}

		adminID := helper.MustUserID()
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

		helper.Success(CenterTeacherNoteResponse{
			ID:           existingNote.ID,
			TeacherID:    existingNote.TeacherID,
			Rating:       existingNote.Rating,
			InternalNote: existingNote.InternalNote,
			UpdatedAt:    existingNote.UpdatedAt,
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
		helper.InternalError("Failed to create teacher note")
		return
	}

	adminID := helper.MustUserID()
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

	helper.Success(CenterTeacherNoteResponse{
		ID:           newNote.ID,
		TeacherID:    newNote.TeacherID,
		Rating:       newNote.Rating,
		InternalNote: newNote.InternalNote,
		CreatedAt:    newNote.CreatedAt,
		UpdatedAt:    newNote.UpdatedAt,
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
func (ctl *AdminTeacherController) DeleteTeacherNote(ctx *gin.Context) {
	helper := NewContextHelper(ctx)

	centerID := helper.MustCenterID()
	if centerID == 0 {
		return
	}

	teacherID := helper.MustParamUint("teacher_id")
	if teacherID == 0 {
		helper.BadRequest("Invalid teacher ID")
		return
	}

	note, err := ctl.centerTeacherNoteRepo.GetByCenterAndTeacher(ctx, centerID, teacherID)
	if err != nil {
		if err.Error() == "record not found" {
			helper.NotFound("Teacher note not found")
			return
		}
		helper.InternalError("Failed to get teacher note")
		return
	}

	if err := ctl.centerTeacherNoteRepo.DeleteByID(ctx, note.ID); err != nil {
		helper.InternalError("Failed to delete teacher note")
		return
	}

	adminID := helper.MustUserID()
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

	helper.Success(nil)
}
