package controllers

import (
	"context"
	"fmt"
	"time"
	"timeLedger/app"
	"timeLedger/app/models"
	"timeLedger/app/repositories"
	"timeLedger/app/resources"
	"timeLedger/app/services"
	"timeLedger/global/logger"

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
	teacherMergeService    *services.TeacherMergeService
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
		teacherMergeService:   services.NewTeacherMergeService(app),
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

// CreatePlaceholderTeacher 建立佔位老師（用於中心暫時無 LINE 帳號的老師）
// @Summary 建立佔位老師
// @Tags Admin
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body CreatePlaceholderTeacherRequest true "佔位老師資料"
// @Success 200 {object} global.ApiResponse{data=resources.AdminTeacherResponse}
// @Router /api/v1/admin/teachers/placeholder [post]
func (ctl *AdminTeacherController) CreatePlaceholderTeacher(ctx *gin.Context) {
	helper := NewContextHelper(ctx)

	centerID := helper.MustCenterID()
	if centerID == 0 {
		return
	}

	adminID := helper.MustUserID()
	if adminID == 0 {
		return
	}

	var req CreatePlaceholderTeacherRequest
	if !helper.MustBindJSON(&req) {
		return
	}

	var teacherID uint

	// 使用交易確保資料一致性
	err := ctl.teacherRepository.Transaction(ctx.Request.Context(), func(txRepo *repositories.TeacherRepository) error {
		// 建立佔位老師
		now := time.Now()
		teacher := models.Teacher{
			Name:          req.Name,
			Email:         req.Email,
			LineUserID:    "", // 佔位老師無 LINE ID
			IsPlaceholder: true,
			CreatedAt:     now,
			UpdatedAt:     now,
		}

		createdTeacher, err := txRepo.Create(ctx.Request.Context(), teacher)
		if err != nil {
			return fmt.Errorf("failed to create placeholder teacher: %w", err)
		}

		teacherID = createdTeacher.ID

		// 建立中心會籍（狀態為 ACTIVE）
		membership := models.CenterMembership{
			CenterID:  centerID,
			TeacherID: createdTeacher.ID,
			Role:      "TEACHER",
			Status:    "ACTIVE",
			CreatedAt: now,
			UpdatedAt: now,
		}

		_, err = ctl.membershipRepo.Create(ctx.Request.Context(), membership)
		if err != nil {
			return fmt.Errorf("failed to create center membership: %w", err)
		}

		// 審計日誌
		ctl.auditLogRepo.Create(ctx.Request.Context(), models.AuditLog{
			CenterID:   centerID,
			ActorType:  "ADMIN",
			ActorID:    adminID,
			Action:     "CREATE_PLACEHOLDER_TEACHER",
			TargetType: "Teacher",
			TargetID:   createdTeacher.ID,
			Payload: models.AuditPayload{
				After: map[string]interface{}{
					"name":            createdTeacher.Name,
					"email":           createdTeacher.Email,
					"is_placeholder":   true,
					"membership_status": "ACTIVE",
				},
			},
		})

		return nil
	})

	if err != nil {
		helper.InternalError("Failed to create placeholder teacher: " + err.Error())
		return
	}

	// 重新查詢建立的老師（確保取得完整資料）
	teacher, err := ctl.teacherRepository.GetByID(ctx.Request.Context(), teacherID)
	if err != nil {
		helper.InternalError("Failed to fetch created teacher")
		return
	}

	// 轉換回應（無技能和證照）
	response := ctl.adminTeacherResource.ToAdminTeacherResponse(&teacher, nil, nil)
	helper.Success(response)
}

// DeleteTeacher 刪除老師（全局刪除，僅限系統管理員）
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

	// 同步刪除該老師的所有會籍（軟刪除）
	if _, err := ctl.membershipRepo.DeleteWhere(ctx, "teacher_id = ?", teacherID); err != nil {
		// 這裡失敗了也沒關係，Teacher 已刪除，但為了資料一致性我們嘗試刪除
		logger.GetLogger().Errorw("Failed to cleanup memberships after deleting teacher", "teacher_id", teacherID, "error", err)
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

// RemoveFromCenter 將老師從中心移除（僅移除會籍）
// @Summary 將老師從中心移除
// @Tags Admin
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Center ID"
// @Param teacher_id path int true "Teacher ID"
// @Success 200 {object} global.ApiResponse
// @Router /api/v1/admin/centers/{id}/teachers/{teacher_id} [delete]
func (ctl *AdminTeacherController) RemoveFromCenter(ctx *gin.Context) {
	helper := NewContextHelper(ctx)

	adminID := helper.MustUserID()
	if adminID == 0 {
		return
	}

	centerID := helper.MustCenterID()
	if centerID == 0 {
		return
	}

	teacherID := helper.MustParamUint("teacher_id")
	if teacherID == 0 {
		return
	}

	// 檢查會籍是否存在
	membership, err := ctl.membershipRepo.GetByCenterAndTeacher(ctx, centerID, teacherID)
	if err != nil {
		helper.NotFound("Teacher is not a member of this center")
		return
	}

	// 刪除會籍（軟刪除）
	if err := ctl.membershipRepo.DeleteByID(ctx, membership.ID); err != nil {
		helper.InternalError("Failed to remove teacher from center")
		return
	}

	// 審計日誌
	ctl.auditLogRepo.Create(ctx, models.AuditLog{
		CenterID:   centerID,
		ActorType:  "ADMIN",
		ActorID:    adminID,
		Action:     "REMOVE_TEACHER_FROM_CENTER",
		TargetType: "CenterMembership",
		TargetID:   membership.ID,
		Payload: models.AuditPayload{
			Before: map[string]interface{}{
				"teacher_id": teacherID,
				"center_id":  centerID,
				"role":       membership.Role,
				"status":     membership.Status,
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

// CreatePlaceholderTeacherRequest 建立佔位老師請求結構
type CreatePlaceholderTeacherRequest struct {
	Name  string `json:"name" binding:"required,min=1,max=255"`
	Email string `json:"email" binding:"omitempty,email"`
}

// MergeTeachersRequest 合併教師請求結構
type MergeTeachersRequest struct {
	SourceTeacherID uint `json:"source_teacher_id" binding:"required"`
	TargetTeacherID uint `json:"target_teacher_id" binding:"required"`
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

// MergeTeachers 合併兩位教師的資料
// @Summary 合併教師資料
// @Description 將來源教師的所有關聯資料遷移到目標教師，然後軟刪除來源教師
// @Tags Admin
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body MergeTeachersRequest true "合併請求資料"
// @Success 200 {object} global.ApiResponse
// @Router /api/v1/admin/teachers/merge [post]
func (ctl *AdminTeacherController) MergeTeachers(ctx *gin.Context) {
	helper := NewContextHelper(ctx)

	// 取得管理員 ID
	adminID := helper.MustUserID()
	if adminID == 0 {
		return
	}

	// 取得中心 ID（從 JWT Token）
	centerID := helper.MustCenterID()
	if centerID == 0 {
		return
	}

	// 綁定請求資料
	var req MergeTeachersRequest
	if !helper.MustBindJSON(&req) {
		return
	}

	// 驗證參數
	if req.SourceTeacherID == req.TargetTeacherID {
		helper.BadRequest("來源教師與目標教師相同，無法合併")
		return
	}

	if req.SourceTeacherID == 0 || req.TargetTeacherID == 0 {
		helper.BadRequest("教師 ID 不能為零")
		return
	}

	// 驗證兩位教師是否都屬於該中心
	if err := ctl.validateTeachersBelongToCenter(centerID, req.SourceTeacherID, req.TargetTeacherID); err != nil {
		helper.BadRequest(err.Error())
		return
	}

	// 執行合併操作
	if err := ctl.teacherMergeService.MergeTeacher(ctx.Request.Context(), req.SourceTeacherID, req.TargetTeacherID, centerID); err != nil {
		helper.InternalError("合併教師失敗: " + err.Error())
		return
	}

	// 記錄審計日誌
	ctl.auditLogRepo.Create(ctx, models.AuditLog{
		CenterID:   centerID,
		ActorType:  "ADMIN",
		ActorID:    adminID,
		Action:     "MERGE_TEACHERS",
		TargetType: "Teacher",
		TargetID:   req.TargetTeacherID,
		Payload: models.AuditPayload{
			After: map[string]interface{}{
				"source_teacher_id": req.SourceTeacherID,
				"target_teacher_id": req.TargetTeacherID,
			},
		},
	})

	helper.Success(map[string]string{
		"message":            "教師合併成功",
		"source_teacher_id":  fmt.Sprintf("%d", req.SourceTeacherID),
		"target_teacher_id":  fmt.Sprintf("%d", req.TargetTeacherID),
	})
}

// validateTeachersBelongToCenter 驗證兩位教師是否都屬於該中心
func (ctl *AdminTeacherController) validateTeachersBelongToCenter(centerID, sourceID, targetID uint) error {
	// 檢查來源教師
	sourceMembership, err := ctl.membershipRepo.GetByCenterAndTeacher(
		context.Background(), centerID, sourceID)
	if err != nil {
		return fmt.Errorf("來源教師不屬於該中心")
	}
	if sourceMembership.Status != "ACTIVE" {
		return fmt.Errorf("來源教師的會籍狀態不是 ACTIVE")
	}

	// 檢查目標教師
	targetMembership, err := ctl.membershipRepo.GetByCenterAndTeacher(
		context.Background(), centerID, targetID)
	if err != nil {
		return fmt.Errorf("目標教師不屬於該中心")
	}
	if targetMembership.Status != "ACTIVE" {
		return fmt.Errorf("目標教師的會籍狀態不是 ACTIVE")
	}

	return nil
}
