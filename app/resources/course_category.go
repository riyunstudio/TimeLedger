package resources

import (
	"timeLedger/app"
	"timeLedger/app/models"
)

// CourseCategoryResource 課程類別資源轉換
type CourseCategoryResource struct {
	app *app.App
}

// NewCourseCategoryResource 建立 CourseCategoryResource 實例
func NewCourseCategoryResource(appInstance *app.App) *CourseCategoryResource {
	return &CourseCategoryResource{
		app: appInstance,
	}
}

// ToCategoryResponse 將模型轉換為響應格式
func (r *CourseCategoryResource) ToCategoryResponse(category models.CourseCategory) *models.CourseCategoryResponse {
	return &models.CourseCategoryResponse{
		ID:       category.ID,
		CenterID: category.CenterID,
		Name:     category.Name,
	}
}

// ToCategoryResponses 批量轉換
func (r *CourseCategoryResource) ToCategoryResponses(categories []models.CourseCategory) []models.CourseCategoryResponse {
	if categories == nil {
		return nil
	}

	responses := make([]models.CourseCategoryResponse, len(categories))
	for i, category := range categories {
		responses[i] = *r.ToCategoryResponse(category)
	}
	return responses
}
