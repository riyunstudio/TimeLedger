package resources

import (
	"time"
	"timeLedger/app"
	"timeLedger/app/models"
)

// OfferingResource 開課資源轉換
type OfferingResource struct {
	app *app.App
}

// NewOfferingResource 建立 OfferingResource 實例
func NewOfferingResource(appInstance *app.App) *OfferingResource {
	return &OfferingResource{
		app: appInstance,
	}
}

// OfferingResponse 開課響應結構（含課程時長）
type OfferingResponse struct {
	ID                   uint           `json:"id"`
	CenterID             uint           `json:"center_id"`
	CourseID             uint           `json:"course_id"`
	Name                 string         `json:"name"`
	CourseName           string         `json:"course_name,omitempty"`
	CourseDuration       int            `json:"course_duration"` // 課程時長（分鐘），優先使用課程設定，若無則使用中心全局設定
	DefaultRoomID        *uint          `json:"default_room_id,omitempty"`
	DefaultTeacherID     *uint          `json:"default_teacher_id,omitempty"`
	AllowBufferOverride  bool           `json:"allow_buffer_override"`
	IsActive             bool           `json:"is_active"`
	CreatedAt            time.Time      `json:"created_at"`
	UpdatedAt            time.Time      `json:"updated_at"`
}

// ToOfferingResponse 將開課模型轉換為響應格式
func (r *OfferingResource) ToOfferingResponse(offering models.Offering) *OfferingResponse {
	response := &OfferingResponse{
		ID:                  offering.ID,
		CenterID:            offering.CenterID,
		CourseID:            offering.CourseID,
		Name:                offering.Name,
		AllowBufferOverride: offering.AllowBufferOverride,
		IsActive:            offering.IsActive,
		CreatedAt:           offering.CreatedAt,
		UpdatedAt:           offering.UpdatedAt,
	}

	// 如果有關聯課程，帶入課程名稱
	if offering.Course.ID != 0 {
		response.CourseName = offering.Course.Name
	}

	// 課程時長處理策略：
	// 1. 優先使用課程的 DefaultDuration（如果設定且 > 0）
	// 2. 否則使用中心的全局設定 DefaultCourseDuration
	response.CourseDuration = r.getCourseDuration(offering)

	// 處理預設教室
	if offering.DefaultRoomID != nil {
		response.DefaultRoomID = offering.DefaultRoomID
	}

	// 處理預設教師
	if offering.DefaultTeacherID != nil {
		response.DefaultTeacherID = offering.DefaultTeacherID
	}

	return response
}

// getCourseDuration 取得課程時長
// 優先順序：課程 DefaultDuration > 中心全局設定
func (r *OfferingResource) getCourseDuration(offering models.Offering) int {
	// 優先使用課程的 DefaultDuration
	if offering.Course.ID != 0 && offering.Course.DefaultDuration > 0 {
		return offering.Course.DefaultDuration
	}

	// 否則從中心獲取全局設定
	var center models.Center
	err := r.app.MySQL.RDB.WithContext(nil).
		Select("settings").
		Where("id = ?", offering.CenterID).
		First(&center).Error

	if err != nil {
		// 如果獲取失敗，返回預設值 60
		return 60
	}

	// 如果中心設定中有 DefaultCourseDuration 且 > 0，使用該值
	if center.Settings.DefaultCourseDuration > 0 {
		return center.Settings.DefaultCourseDuration
	}

	// 最終預設值
	return 60
}

// ToOfferingResponses 批量將開課模型轉換為響應格式
func (r *OfferingResource) ToOfferingResponses(offerings []models.Offering) []OfferingResponse {
	if offerings == nil {
		return nil
	}

	responses := make([]OfferingResponse, len(offerings))
	for i, offering := range offerings {
		responses[i] = *r.ToOfferingResponse(offering)
	}
	return responses
}
