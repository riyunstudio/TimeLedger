package controllers

import (
	"timeLedger/app"
	"timeLedger/app/models"
	"timeLedger/app/repositories"
	"timeLedger/app/resources"

	"github.com/gin-gonic/gin"
)

// AdminResourceController 管理員資源控制器
// 處理中心老師列表等資源管理功能
type AdminResourceController struct {
	BaseController
	app                  *app.App
	offeringRepository   *repositories.OfferingRepository
	holidayRepository    *repositories.CenterHolidayRepository
	auditLogRepo         *repositories.AuditLogRepository
	membershipRepo       *repositories.CenterMembershipRepository
	teacherRepository    *repositories.TeacherRepository
	skillRepository      *repositories.TeacherSkillRepository
	certificateRepo      *repositories.TeacherCertificateRepository
	adminTeacherResource *resources.AdminTeacherResource
}

func NewAdminResourceController(app *app.App) *AdminResourceController {
	return &AdminResourceController{
		app:                  app,
		offeringRepository:  repositories.NewOfferingRepository(app),
		holidayRepository:   repositories.NewCenterHolidayRepository(app),
		auditLogRepo:        repositories.NewAuditLogRepository(app),
		membershipRepo:      repositories.NewCenterMembershipRepository(app),
		teacherRepository:   repositories.NewTeacherRepository(app),
		skillRepository:     repositories.NewTeacherSkillRepository(app),
		certificateRepo:     repositories.NewTeacherCertificateRepository(app),
		adminTeacherResource: resources.NewAdminTeacherResource(),
	}
}

// GetTeachers 取得中心的老師列表
// @Summary 取得中心的老師列表
// @Tags Admin
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} global.ApiResponse{data=[]resources.AdminTeacherResponse}
// @Router /api/v1/admin/teachers [get]
func (ctl *AdminResourceController) GetTeachers(ctx *gin.Context) {
	helper := NewContextHelper(ctx)

	centerID := helper.MustCenterID()
	if centerID == 0 {
		helper.Forbidden("Center ID is required")
		return
	}

	// 取得該中心的所有活躍會員的老師 ID
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

	if len(teachersMap) == 0 {
		helper.Success([]resources.AdminTeacherResponse{})
		return
	}

	// 按原始順序重建 slice
	teachers := make([]models.Teacher, 0, len(teacherIDs))
	for _, id := range teacherIDs {
		if teacher, ok := teachersMap[id]; ok {
			teachers = append(teachers, teacher)
		}
	}

	// 批次查詢技能
	skillsMap, _ := ctl.skillRepository.BatchListByTeacherIDs(ctx, teacherIDs)

	// 批次查詢證照
	certificatesMap, _ := ctl.certificateRepo.BatchListByTeacherIDs(ctx, teacherIDs)

	// 使用 Resource 轉換
	responses := ctl.adminTeacherResource.ToAdminTeacherResponses(teachers, skillsMap, certificatesMap)

	helper.Success(responses)
}
