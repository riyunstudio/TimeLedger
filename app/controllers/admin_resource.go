package controllers

import (
	"timeLedger/app"
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
	teacherRepository     *repositories.TeacherRepository
	skillRepository      *repositories.TeacherSkillRepository
	certificateRepo      *repositories.TeacherCertificateRepository
	noteRepo             *repositories.CenterTeacherNoteRepository
	adminTeacherResource *resources.AdminTeacherResource
}

func NewAdminResourceController(app *app.App) *AdminResourceController {
	return &AdminResourceController{
		app:                  app,
		offeringRepository:   repositories.NewOfferingRepository(app),
		holidayRepository:   repositories.NewCenterHolidayRepository(app),
		auditLogRepo:        repositories.NewAuditLogRepository(app),
		membershipRepo:      repositories.NewCenterMembershipRepository(app),
		teacherRepository:   repositories.NewTeacherRepository(app),
		skillRepository:     repositories.NewTeacherSkillRepository(app),
		certificateRepo:      repositories.NewTeacherCertificateRepository(app),
		noteRepo:            repositories.NewCenterTeacherNoteRepository(app),
		adminTeacherResource: resources.NewAdminTeacherResource(),
	}
}

// GetTeachers 取得中心的老師列表
// @Summary 取得中心的老師列表
// @Tags Admin
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param query query string false "關鍵字搜尋"
// @Param page query int false "頁碼，預設 1"
// @Param limit query int false "每頁筆數，預設 20"
// @Success 200 {object} global.ApiResponse{data=resources.PaginationResponse}
// @Router /api/v1/admin/teachers [get]
func (ctl *AdminResourceController) GetTeachers(ctx *gin.Context) {
	helper := NewContextHelper(ctx)

	centerID := helper.MustCenterID()
	if centerID == 0 {
		helper.Forbidden("Center ID is required")
		return
	}

	// 取得查詢參數
	query := helper.QueryStringOrDefault("query", "")
	page := helper.QueryIntOrDefault("page", 1)
	limit := helper.QueryIntOrDefault("limit", 20)

	// 使用分頁查詢取得老師列表
	teachers, total, err := ctl.membershipRepo.ListTeachersByCenterPaginated(ctx.Request.Context(), centerID, query, page, limit)
	if err != nil {
		helper.InternalError("Failed to get teachers")
		return
	}

	if len(teachers) == 0 {
		helper.Success(resources.NewPaginationResponse([]resources.AdminTeacherResponse{}, total, page, limit))
		return
	}

	// 提取老師 IDs
	teacherIDs := make([]uint, 0, len(teachers))
	for _, t := range teachers {
		teacherIDs = append(teacherIDs, t.ID)
	}

	// 批次查詢技能
	skillsMap, _ := ctl.skillRepository.BatchListByTeacherIDs(ctx.Request.Context(), teacherIDs)

	// 批次查詢證照
	certificatesMap, _ := ctl.certificateRepo.BatchListByTeacherIDs(ctx.Request.Context(), teacherIDs)

	// 批次查詢老師備註
	notesMap, _ := ctl.noteRepo.BatchGetByCenterAndTeachers(ctx.Request.Context(), centerID, teacherIDs)

	// 使用 Resource 轉換
	responses := ctl.adminTeacherResource.ToAdminTeacherResponses(teachers, skillsMap, certificatesMap, notesMap)

	helper.Success(resources.NewPaginationResponse(responses, total, page, limit))
}
