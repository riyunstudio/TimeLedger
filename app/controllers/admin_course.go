package controllers

import (
	"timeLedger/app"
	"timeLedger/app/resources"
	"timeLedger/app/services"

	"github.com/gin-gonic/gin"
)

type AdminCourseController struct {
	BaseController
	app           *app.App
	courseService *services.CourseService
	courseResource *resources.CourseResource
}

func NewAdminCourseController(app *app.App) *AdminCourseController {
	return &AdminCourseController{
		app:            app,
		courseService:  services.NewCourseService(app),
		courseResource: resources.NewCourseResource(app),
	}
}

// GetCourses 取得課程列表
// @Summary 取得課程列表
// @Tags Admin
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} global.ApiResponse{data=[]resources.CourseResponse}
// @Router /api/v1/admin/courses [get]
func (ctl *AdminCourseController) GetCourses(ctx *gin.Context) {
	helper := NewContextHelper(ctx)

	centerID := helper.MustCenterID()
	if centerID == 0 {
		return
	}

	courses, errInfo, err := ctl.courseService.GetCourses(ctx.Request.Context(), centerID)
	if err != nil {
		helper.ErrorWithInfo(errInfo)
		return
	}

	responses := ctl.courseResource.ToCourseResponses(courses)
	helper.Success(responses)
}

// CreateCourse 新增課程
// @Summary 新增課程
// @Tags Admin
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body services.CreateCourseRequest true "課程資訊"
// @Success 200 {object} global.ApiResponse{data=resources.CourseResponse}
// @Router /api/v1/admin/courses [post]
func (ctl *AdminCourseController) CreateCourse(ctx *gin.Context) {
	helper := NewContextHelper(ctx)

	centerID := helper.MustCenterID()
	if centerID == 0 {
		return
	}

	var req services.CreateCourseRequest
	if !helper.MustBindJSON(&req) {
		return
	}

	adminID := helper.MustUserID()
	if adminID == 0 {
		return
	}

	course, errInfo, err := ctl.courseService.CreateCourse(ctx.Request.Context(), centerID, adminID, &req)
	if err != nil {
		helper.ErrorWithInfo(errInfo)
		return
	}

	response := ctl.courseResource.ToCourseResponse(*course)
	helper.Success(response)
}

// UpdateCourse 更新課程
// @Summary 更新課程
// @Tags Admin
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param course_id path int true "Course ID"
// @Param request body services.UpdateCourseRequest true "課程資訊"
// @Success 200 {object} global.ApiResponse{data=resources.CourseResponse}
// @Router /api/v1/admin/courses/{course_id} [put]
func (ctl *AdminCourseController) UpdateCourse(ctx *gin.Context) {
	helper := NewContextHelper(ctx)

	centerID := helper.MustCenterID()
	if centerID == 0 {
		return
	}

	courseID := helper.MustParamUint("course_id")
	if courseID == 0 {
		return
	}

	var req services.UpdateCourseRequest
	if !helper.MustBindJSON(&req) {
		return
	}

	adminID := helper.MustUserID()
	if adminID == 0 {
		return
	}

	course, errInfo, err := ctl.courseService.UpdateCourse(ctx.Request.Context(), centerID, adminID, courseID, &req)
	if err != nil {
		helper.ErrorWithInfo(errInfo)
		return
	}

	response := ctl.courseResource.ToCourseResponse(*course)
	helper.Success(response)
}

// DeleteCourse 刪除課程
// @Summary 刪除課程
// @Tags Admin
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param course_id path int true "Course ID"
// @Success 200 {object} global.ApiResponse
// @Router /api/v1/admin/courses/{course_id} [delete]
func (ctl *AdminCourseController) DeleteCourse(ctx *gin.Context) {
	helper := NewContextHelper(ctx)

	centerID := helper.MustCenterID()
	if centerID == 0 {
		return
	}

	courseID := helper.MustParamUint("course_id")
	if courseID == 0 {
		return
	}

	adminID := helper.MustUserID()
	if adminID == 0 {
		return
	}

	errInfo, err := ctl.courseService.DeleteCourse(ctx.Request.Context(), centerID, adminID, courseID)
	if err != nil {
		helper.ErrorWithInfo(errInfo)
		return
	}

	helper.Success(nil)
}

// GetActiveCourses 取得已啟用的課程列表
// @Summary 取得已啟用的課程列表
// @Tags Admin
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} global.ApiResponse{data=[]resources.CourseResponse}
// @Router /api/v1/admin/courses/active [get]
func (ctl *AdminCourseController) GetActiveCourses(ctx *gin.Context) {
	helper := NewContextHelper(ctx)

	centerID := helper.MustCenterID()
	if centerID == 0 {
		return
	}

	courses, errInfo, err := ctl.courseService.GetActiveCourses(ctx.Request.Context(), centerID)
	if err != nil {
		helper.ErrorWithInfo(errInfo)
		return
	}

	responses := ctl.courseResource.ToCourseResponses(courses)
	helper.Success(responses)
}

// ToggleCourseActive 切換課程啟用狀態
// @Summary 切換課程啟用狀態
// @Tags Admin
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param course_id path int true "Course ID"
// @Param request body services.ToggleActiveRequest true "啟用狀態"
// @Success 200 {object} global.ApiResponse
// @Router /api/v1/admin/courses/{course_id}/toggle-active [patch]
func (ctl *AdminCourseController) ToggleCourseActive(ctx *gin.Context) {
	helper := NewContextHelper(ctx)

	centerID := helper.MustCenterID()
	if centerID == 0 {
		return
	}

	courseID := helper.MustParamUint("course_id")
	if courseID == 0 {
		return
	}

	var req services.ToggleActiveRequest
	if !helper.MustBindJSON(&req) {
		return
	}

	adminID := helper.MustUserID()
	if adminID == 0 {
		return
	}

	errInfo, err := ctl.courseService.ToggleCourseActive(ctx.Request.Context(), centerID, adminID, courseID, req.IsActive)
	if err != nil {
		helper.ErrorWithInfo(errInfo)
		return
	}

	helper.Success(nil)
}
