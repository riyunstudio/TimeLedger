package controllers

import (
	"time"
	"timeLedger/app"
	"timeLedger/app/repositories"

	"github.com/gin-gonic/gin"
)

// AdminResourceController 管理員資源控制器
// 處理中心老師列表等資源管理功能
type AdminResourceController struct {
	BaseController
	app                *app.App
	offeringRepository *repositories.OfferingRepository
	holidayRepository  *repositories.CenterHolidayRepository
	auditLogRepo       *repositories.AuditLogRepository
	membershipRepo     *repositories.CenterMembershipRepository
	teacherRepository  *repositories.TeacherRepository
	skillRepository    *repositories.TeacherSkillRepository
	certificateRepo    *repositories.TeacherCertificateRepository
}

func NewAdminResourceController(app *app.App) *AdminResourceController {
	return &AdminResourceController{
		app:                app,
		offeringRepository: repositories.NewOfferingRepository(app),
		holidayRepository:  repositories.NewCenterHolidayRepository(app),
		auditLogRepo:       repositories.NewAuditLogRepository(app),
		membershipRepo:     repositories.NewCenterMembershipRepository(app),
		teacherRepository:  repositories.NewTeacherRepository(app),
		skillRepository:    repositories.NewTeacherSkillRepository(app),
		certificateRepo:    repositories.NewTeacherCertificateRepository(app),
	}
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
			IsActive:     teacher.IsOpenToHiring,
			CreatedAt:    teacher.CreatedAt,
			Skills:       skillResponses,
			Certificates: certResponses,
		})
	}

	helper.Success(teachers)
}

// TeacherResponse 老師回應結構
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

// TeacherSkillResponse 老師技能回應結構
type TeacherSkillResponse struct {
	ID        uint   `json:"id"`
	SkillName string `json:"skill_name"`
	Category  string `json:"category,omitempty"`
	Level     string `json:"level,omitempty"`
}

// CertificateResponse 證照回應結構
type CertificateResponse struct {
	ID        uint      `json:"id"`
	Name      string    `json:"name"`
	FileURL   string    `json:"file_url,omitempty"`
	IssuedAt  time.Time `json:"issued_at"`
	CreatedAt time.Time `json:"created_at"`
}
