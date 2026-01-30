package resources

import (
	"time"
	"timeLedger/app"
	"timeLedger/app/models"
)

// CourseResource 課程資源轉換
type CourseResource struct {
	app *app.App
}

// NewCourseResource 建立 CourseResource 實例
func NewCourseResource(appInstance *app.App) *CourseResource {
	return &CourseResource{
		app: appInstance,
	}
}

// CourseResponse 課程響應結構
type CourseResponse struct {
	ID               uint      `json:"id"`
	CenterID         uint      `json:"center_id"`
	Name             string    `json:"name"`
	DefaultDuration  int       `json:"default_duration"`
	ColorHex         string    `json:"color_hex"`
	RoomBufferMin    int       `json:"room_buffer_min"`
	TeacherBufferMin int       `json:"teacher_buffer_min"`
	IsActive         bool      `json:"is_active"`
	CreatedAt        time.Time `json:"created_at"`
}

// ToCourseResponse 將課程模型轉換為響應格式
func (r *CourseResource) ToCourseResponse(course models.Course) *CourseResponse {
	return &CourseResponse{
		ID:               course.ID,
		CenterID:         course.CenterID,
		Name:             course.Name,
		DefaultDuration:  course.DefaultDuration,
		ColorHex:         course.ColorHex,
		RoomBufferMin:    course.RoomBufferMin,
		TeacherBufferMin: course.TeacherBufferMin,
		IsActive:         course.IsActive,
		CreatedAt:        course.CreatedAt,
	}
}

// ToCourseResponses 批量將課程模型轉換為響應格式
func (r *CourseResource) ToCourseResponses(courses []models.Course) []CourseResponse {
	if courses == nil {
		return nil
	}

	responses := make([]CourseResponse, len(courses))
	for i, course := range courses {
		responses[i] = *r.ToCourseResponse(course)
	}
	return responses
}
