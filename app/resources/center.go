package resources

import (
	"time"
	"timeLedger/app"
	"timeLedger/app/models"
)

// CenterResource 中心資源轉換
type CenterResource struct {
	app *app.App
}

// NewCenterResource 建立 CenterResource 實例
func NewCenterResource(appInstance *app.App) *CenterResource {
	return &CenterResource{
		app: appInstance,
	}
}

// CenterResponse 中心響應結構
type CenterResponse struct {
	ID        uint           `json:"id"`
	Name      string         `json:"name"`
	PlanLevel string         `json:"plan_level"`
	Settings  CenterSettings `json:"settings"`
	CreatedAt time.Time      `json:"created_at"`
}

// CenterSettings 中心設置響應
type CenterSettings struct {
	AllowPublicRegister   bool   `json:"allow_public_register"`
	DefaultLanguage       string `json:"default_language"`
	ExceptionLeadDays     int    `json:"exception_lead_days"`
	DefaultCourseDuration int    `json:"default_course_duration"`
}

// CenterSettingsResponse 僅包含設置的響應
type CenterSettingsResponse struct {
	AllowPublicRegister   bool   `json:"allow_public_register"`
	DefaultLanguage       string `json:"default_language"`
	ExceptionLeadDays     int    `json:"exception_lead_days"`
	DefaultCourseDuration int    `json:"default_course_duration"`
}

// ToCenterResponse 將中心模型轉換為響應格式
func (r *CenterResource) ToCenterResponse(center models.Center) *CenterResponse {
	return &CenterResponse{
		ID:        center.ID,
		Name:      center.Name,
		PlanLevel: center.PlanLevel,
		Settings: CenterSettings{
			AllowPublicRegister:   center.Settings.AllowPublicRegister,
			DefaultLanguage:       center.Settings.DefaultLanguage,
			ExceptionLeadDays:     center.Settings.ExceptionLeadDays,
			DefaultCourseDuration: center.Settings.DefaultCourseDuration,
		},
		CreatedAt: center.CreatedAt,
	}
}

// ToCenterResponses 批量將中心模型轉換為響應格式
func (r *CenterResource) ToCenterResponses(centers []models.Center) []CenterResponse {
	if centers == nil {
		return nil
	}

	responses := make([]CenterResponse, len(centers))
	for i, center := range centers {
		responses[i] = *r.ToCenterResponse(center)
	}
	return responses
}

// ToSettingsResponse 將中心設置模型轉換為響應格式
func (r *CenterResource) ToSettingsResponse(settings models.CenterSettings) *CenterSettingsResponse {
	return &CenterSettingsResponse{
		AllowPublicRegister:   settings.AllowPublicRegister,
		DefaultLanguage:       settings.DefaultLanguage,
		ExceptionLeadDays:     settings.ExceptionLeadDays,
		DefaultCourseDuration: settings.DefaultCourseDuration,
	}
}
