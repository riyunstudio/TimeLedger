package controllers

import (
	"timeLedger/app"
	"timeLedger/app/models"
	"timeLedger/app/repositories"

	"github.com/gin-gonic/gin"
)

// AdminTeacherController 管理員對老師操作相關 API
type AdminTeacherController struct {
	BaseController
	app               *app.App
	teacherRepository *repositories.TeacherRepository
	membershipRepo    *repositories.CenterMembershipRepository
	auditLogRepo      *repositories.AuditLogRepository
}

func NewAdminTeacherController(app *app.App) *AdminTeacherController {
	return &AdminTeacherController{
		app:               app,
		teacherRepository: repositories.NewTeacherRepository(app),
		membershipRepo:    repositories.NewCenterMembershipRepository(app),
		auditLogRepo:      repositories.NewAuditLogRepository(app),
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
