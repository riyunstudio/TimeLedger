package resources

import (
	"time"
	"timeLedger/app/models"
)

type TeacherProfileResource struct {
	ID                uint              `json:"id,omitempty"`
	LineUserID        string            `json:"line_user_id,omitempty"`
	Name              string            `json:"name,omitempty"`
	Email             string            `json:"email,omitempty"`
	Bio               string            `json:"bio,omitempty"`
	City              string            `json:"city,omitempty"`
	District          string            `json:"district,omitempty"`
	PublicContactInfo string            `json:"public_contact_info,omitempty"`
	IsOpenToHiring    bool              `json:"is_open_to_hiring,omitempty"`
	PersonalHashtags  []PersonalHashtag `json:"personal_hashtags,omitempty"`
}

type TeacherPersonalHashtagResource struct {
	ID        uint   `json:"id,omitempty"`
	HashtagID uint   `json:"hashtag_id,omitempty"`
	Name      string `json:"name,omitempty"`
}

type PersonalHashtag struct {
	HashtagID uint   `json:"hashtag_id"`
	Name      string `json:"name"`
	SortOrder int    `json:"sort_order"`
}

type TeacherCertificateResource struct {
	ID        uint      `json:"id,omitempty"`
	Name      string    `json:"name,omitempty"`
	FileURL   string    `json:"file_url,omitempty"`
	IssuedAt  time.Time `json:"issued_at,omitempty"`
	CreatedAt time.Time `json:"created_at,omitempty"`
}

type CenterMembershipResource struct {
	ID         uint      `json:"id,omitempty"`
	CenterID   uint      `json:"center_id,omitempty"`
	CenterName string    `json:"center_name,omitempty"`
	Role       string    `json:"role,omitempty"`
	Status     string    `json:"status,omitempty"`
	CreatedAt  time.Time `json:"created_at,omitempty"`
}

type HashtagResource struct {
	ID         uint   `json:"id,omitempty"`
	Name       string `json:"name,omitempty"`
	UsageCount int    `json:"usage_count,omitempty"`
}

type GeoCityResource struct {
	ID        uint                  `json:"id,omitempty"`
	Name      string                `json:"name,omitempty"`
	Districts []GeoDistrictResource `json:"districts,omitempty"`
}

type GeoDistrictResource struct {
	ID     uint   `json:"id,omitempty"`
	CityID uint   `json:"city_id,omitempty"`
	Name   string `json:"name,omitempty"`
}

// AdminTeacherResource 管理員視角的老師資源（用於中心老師列表）
type AdminTeacherResource struct{}

// NewAdminTeacherResource 建立管理員老師資源轉換實例
func NewAdminTeacherResource() *AdminTeacherResource {
	return &AdminTeacherResource{}
}

// AdminTeacherResponse 管理員取得的老師回應結構
type AdminTeacherResponse struct {
	ID            uint                        `json:"id"`
	Name          string                      `json:"name"`
	Email         string                      `json:"email"`
	Phone         string                      `json:"phone,omitempty"`
	City          string                      `json:"city,omitempty"`
	District      string                      `json:"district,omitempty"`
	Bio           string                      `json:"bio,omitempty"`
	IsActive      bool                        `json:"is_active"`
	IsPlaceholder bool                        `json:"is_placeholder"`
	LineUserID    string                      `json:"line_user_id,omitempty"`
	CreatedAt     time.Time                   `json:"created_at"`
	Skills        []AdminTeacherSkillResponse `json:"skills,omitempty"`
	Certificates  []AdminCertificateResponse  `json:"certificates,omitempty"`
}

// AdminTeacherSkillResponse 管理員視角的老師技能回應結構
type AdminTeacherSkillResponse struct {
	ID        uint   `json:"id"`
	SkillName string `json:"skill_name"`
	Category  string `json:"category,omitempty"`
	Level     string `json:"level,omitempty"`
}

// AdminCertificateResponse 證照回應結構（管理員視角）
type AdminCertificateResponse struct {
	ID        uint      `json:"id"`
	Name      string    `json:"name"`
	FileURL   string    `json:"file_url,omitempty"`
	IssuedAt  time.Time `json:"issued_at"`
	CreatedAt time.Time `json:"created_at"`
}

// ToAdminTeacherResponse 將老師模型轉換為管理員回應
func (r *AdminTeacherResource) ToAdminTeacherResponse(
	teacher *models.Teacher,
	skills []models.TeacherSkill,
	certificates []models.TeacherCertificate,
) *AdminTeacherResponse {
	skillResponses := make([]AdminTeacherSkillResponse, 0, len(skills))
	for _, skill := range skills {
		skillResponses = append(skillResponses, AdminTeacherSkillResponse{
			ID:        skill.ID,
			SkillName: skill.SkillName,
			Category:  skill.Category,
			Level:     skill.Level,
		})
	}

	certResponses := make([]AdminCertificateResponse, 0, len(certificates))
	for _, cert := range certificates {
		certResponses = append(certResponses, AdminCertificateResponse{
			ID:        cert.ID,
			Name:      cert.Name,
			FileURL:   cert.FileURL,
			IssuedAt:  cert.IssuedAt,
			CreatedAt: cert.CreatedAt,
		})
	}

	return &AdminTeacherResponse{
		ID:            teacher.ID,
		Name:          teacher.Name,
		Email:         teacher.Email,
		Phone:         teacher.PublicContactInfo,
		City:          teacher.City,
		District:      teacher.District,
		Bio:           teacher.Bio,
		IsActive:      teacher.IsOpenToHiring,
		IsPlaceholder: teacher.IsPlaceholder,
		LineUserID:    teacher.LineUserID,
		CreatedAt:     teacher.CreatedAt,
		Skills:        skillResponses,
		Certificates:  certResponses,
	}
}

// ToAdminTeacherResponses 將老師列表轉換為管理員回應列表
func (r *AdminTeacherResource) ToAdminTeacherResponses(
	teachers []models.Teacher,
	skillsMap map[uint][]models.TeacherSkill,
	certificatesMap map[uint][]models.TeacherCertificate,
) []AdminTeacherResponse {
	responses := make([]AdminTeacherResponse, 0, len(teachers))
	for _, teacher := range teachers {
		skills := skillsMap[teacher.ID]
		certificates := certificatesMap[teacher.ID]
		responses = append(responses, *r.ToAdminTeacherResponse(&teacher, skills, certificates))
	}
	return responses
}
