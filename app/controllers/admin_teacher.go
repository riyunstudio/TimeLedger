package controllers

import (
	"time"
	"timeLedger/app"
	"timeLedger/app/models"
	"timeLedger/app/repositories"
	"timeLedger/app/resources"

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
	adminTeacherResource  *resources.AdminTeacherResource
	teacherNoteResource   *resources.TeacherNoteResource
}

func NewAdminTeacherController(app *app.App) *AdminTeacherController {
	return &AdminTeacherController{
		app:                   app,
		teacherRepository:     repositories.NewTeacherRepository(app),
		membershipRepo:        repositories.NewCenterMembershipRepository(app),
		auditLogRepo:          repositories.NewAuditLogRepository(app),
		centerTeacherNoteRepo: repositories.NewCenterTeacherNoteRepository(app),
		adminTeacherResource:  resources.NewAdminTeacherResource(),
		teacherNoteResource:   resources.NewTeacherNoteResource(),
	}
}

// ListTeachers 取得老師列表（根據當前登入 Admin 的中心過濾）
// @Summary 取得老師列表
// @Tags Admin
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} global.ApiResponse{data=[]resources.AdminTeacherResponse}
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
		helper.Success([]resources.AdminTeacherResponse{})
		return
	}

	// 批次查詢老師資料
	teachersMap, err := ctl.teacherRepository.BatchGetByIDs(ctx, teacherIDs)
	if err != nil {
		helper.InternalError("Failed to get teachers")
		return
	}

	// 按原始順序重建 slice
	teachers := make([]models.Teacher, 0, len(teacherIDs))
	for _, id := range teacherIDs {
		if teacher, ok := teachersMap[id]; ok {
			teachers = append(teachers, teacher)
		}
	}

	// 使用 Resource 轉換（無技能和證照）
	responses := ctl.adminTeacherResource.ToAdminTeacherResponses(teachers, make(map[uint][]models.TeacherSkill), make(map[uint][]models.TeacherCertificate))

	helper.Success(responses)
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
// @Success 200 {object} global.ApiResponse{data=resources.TeacherNoteResponse}
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
			helper.Success(ctl.teacherNoteResource.ToEmptyTeacherNoteResponse(teacherID))
			return
		}
		helper.InternalError("Failed to get teacher note")
		return
	}

	helper.Success(ctl.teacherNoteResource.ToTeacherNoteResponse(&note))
}

// UpsertTeacherNote 新增或更新老師評分與備註
// @Summary 新增或更新老師評分與備註
// @Tags Admin
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param teacher_id path int true "Teacher ID"
// @Param request body UpsertTeacherNoteRequest true "評分與備註"
// @Success 200 {object} global.ApiResponse{data=resources.TeacherNoteResponse}
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

		helper.Success(ctl.teacherNoteResource.ToTeacherNoteResponse(&existingNote))
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

	helper.Success(ctl.teacherNoteResource.ToTeacherNoteResponse(&newNote))
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
