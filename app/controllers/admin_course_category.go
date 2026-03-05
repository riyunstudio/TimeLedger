package controllers

import (
	"strconv"

	"timeLedger/app"
	"timeLedger/app/resources"
	"timeLedger/app/services"

	"github.com/gin-gonic/gin"
)

// AdminCourseCategoryController 課程類別管理 API
type AdminCourseCategoryController struct {
	app              *app.App
	categoryService  *services.CourseCategoryService
	categoryResource *resources.CourseCategoryResource
}

// NewAdminCourseCategoryController 建立 AdminCourseCategoryController 實例
func NewAdminCourseCategoryController(appInstance *app.App) *AdminCourseCategoryController {
	return &AdminCourseCategoryController{
		app:              appInstance,
		categoryService:  services.NewCourseCategoryService(appInstance),
		categoryResource: resources.NewCourseCategoryResource(appInstance),
	}
}

// requireCenterID 取得並驗證中心 ID
func (ctl *AdminCourseCategoryController) requireCenterID(helper *ContextHelper) uint {
	centerID := helper.MustCenterID()
	if centerID == 0 {
		helper.BadRequest("無效的中心 ID")
		return 0
	}
	return centerID
}

// GetCategories 取得課程類別列表
// @Summary 取得課程類別列表
// @Tags Admin - Course Category
// @Produce json
// @Security BearerAuth
// @Success 200 {object} global.ApiResponse{data=[]models.CourseCategoryResponse}
// @Router /api/v1/admin/course-categories [get]
func (ctl *AdminCourseCategoryController) GetCategories(ctx *gin.Context) {
	helper := NewContextHelper(ctx)

	centerID := ctl.requireCenterID(helper)
	if centerID == 0 {
		return
	}

	categories, errInfo, err := ctl.categoryService.GetCategories(ctx.Request.Context(), centerID)
	if err != nil {
		helper.ErrorWithInfo(errInfo)
		return
	}

	response := ctl.categoryResource.ToCategoryResponses(categories)
	helper.Success(response)
}

// CreateCategoryRequest 建立課程類別請求
type CreateCategoryRequest struct {
	Name string `json:"name" binding:"required"`
}

// CreateCategory 建立課程類別
// @Summary 建立課程類別
// @Tags Admin - Course Category
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body CreateCategoryRequest true "類別資訊"
// @Success 200 {object} global.ApiResponse{data=models.CourseCategoryResponse}
// @Router /api/v1/admin/course-categories [post]
func (ctl *AdminCourseCategoryController) CreateCategory(ctx *gin.Context) {
	helper := NewContextHelper(ctx)

	centerID := ctl.requireCenterID(helper)
	if centerID == 0 {
		return
	}

	var req CreateCategoryRequest
	if !helper.MustBindJSON(&req) {
		return
	}

	category, errInfo, err := ctl.categoryService.CreateCategory(ctx.Request.Context(), centerID, req.Name)
	if err != nil {
		helper.ErrorWithInfo(errInfo)
		return
	}

	response := ctl.categoryResource.ToCategoryResponse(*category)
	helper.Created(response)
}

// UpdateCategoryRequest 更新課程類別請求
type UpdateCategoryRequest struct {
	Name string `json:"name" binding:"required"`
}

// UpdateCategory 更新課程類別
// @Summary 更新課程類別
// @Tags Admin - Course Category
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "類別 ID"
// @Param request body UpdateCategoryRequest true "類別資訊"
// @Success 200 {object} global.ApiResponse{data=models.CourseCategoryResponse}
// @Router /api/v1/admin/course-categories/{id} [put]
func (ctl *AdminCourseCategoryController) UpdateCategory(ctx *gin.Context) {
	helper := NewContextHelper(ctx)

	centerID := ctl.requireCenterID(helper)
	if centerID == 0 {
		return
	}

	categoryID, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil || categoryID == 0 {
		helper.BadRequest("無效的類別 ID")
		return
	}

	var req UpdateCategoryRequest
	if !helper.MustBindJSON(&req) {
		return
	}

	category, errInfo, err := ctl.categoryService.UpdateCategory(ctx.Request.Context(), centerID, uint(categoryID), req.Name)
	if err != nil {
		helper.ErrorWithInfo(errInfo)
		return
	}

	response := ctl.categoryResource.ToCategoryResponse(*category)
	helper.Success(response)
}

// DeleteCategory 刪除課程類別
// @Summary 刪除課程類別
// @Tags Admin - Course Category
// @Produce json
// @Security BearerAuth
// @Param id path int true "類別 ID"
// @Success 200 {object} global.ApiResponse
// @Router /api/v1/admin/course-categories/{id} [delete]
func (ctl *AdminCourseCategoryController) DeleteCategory(ctx *gin.Context) {
	helper := NewContextHelper(ctx)

	centerID := ctl.requireCenterID(helper)
	if centerID == 0 {
		return
	}

	categoryID, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil || categoryID == 0 {
		helper.BadRequest("無效的類別 ID")
		return
	}

	errInfo, err := ctl.categoryService.DeleteCategory(ctx.Request.Context(), centerID, uint(categoryID))
	if err != nil {
		helper.ErrorWithInfo(errInfo)
		return
	}

	helper.Success(nil)
}
