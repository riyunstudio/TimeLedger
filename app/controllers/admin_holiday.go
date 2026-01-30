package controllers

import (
	"timeLedger/app"
	"timeLedger/app/services"

	"github.com/gin-gonic/gin"
)

type AdminHolidayController struct {
	BaseController
	app            *app.App
	holidayService *services.HolidayService
}

func NewAdminHolidayController(app *app.App) *AdminHolidayController {
	return &AdminHolidayController{
		app:            app,
		holidayService: services.NewHolidayService(app),
	}
}

// GetHolidays 取得假日列表
// @Summary 取得假日列表
// @Tags Admin
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Center ID"
// @Param start_date query string false "開始日期"
// @Param end_date query string false "結束日期"
// @Success 200 {object} global.ApiResponse{data=[]models.CenterHoliday}
// @Router /api/v1/admin/centers/{id}/holidays [get]
func (ctl *AdminHolidayController) GetHolidays(ctx *gin.Context) {
	helper := NewContextHelper(ctx)

	centerID := helper.MustParamUint("id")
	if centerID == 0 {
		return
	}

	var req services.GetHolidaysRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		helper.BadRequest("Invalid query parameters")
		return
	}

	holidays, errInfo, err := ctl.holidayService.GetHolidays(ctx.Request.Context(), centerID, &req)
	if err != nil {
		helper.ErrorWithInfo(errInfo)
		return
	}

	helper.Success(holidays)
}

// CreateHoliday 新增假日
// @Summary 新增假日
// @Tags Admin
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Center ID"
// @Param request body services.CreateHolidayRequest true "假日資訊"
// @Success 200 {object} global.ApiResponse{data=models.CenterHoliday}
// @Router /api/v1/admin/centers/{id}/holidays [post]
func (ctl *AdminHolidayController) CreateHoliday(ctx *gin.Context) {
	helper := NewContextHelper(ctx)

	centerID := helper.MustParamUint("id")
	if centerID == 0 {
		return
	}

	var req services.CreateHolidayRequest
	if !helper.MustBindJSON(&req) {
		return
	}

	adminID := helper.MustUserID()
	if adminID == 0 {
		return
	}

	holiday, errInfo, err := ctl.holidayService.CreateHoliday(ctx.Request.Context(), centerID, adminID, &req)
	if err != nil {
		helper.ErrorWithInfo(errInfo)
		return
	}

	helper.Success(holiday)
}

// BulkCreateHolidays 批次建立假日
// @Summary 批次建立假日
// @Tags Admin
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Center ID"
// @Param request body services.BulkCreateHolidaysRequest true "假日列表"
// @Success 200 {object} global.ApiResponse{data=services.BulkCreateHolidaysResponse}
// @Router /api/v1/admin/centers/{id}/holidays/bulk [post]
func (ctl *AdminHolidayController) BulkCreateHolidays(ctx *gin.Context) {
	helper := NewContextHelper(ctx)

	centerID := helper.MustParamUint("id")
	if centerID == 0 {
		return
	}

	var req services.BulkCreateHolidaysRequest
	if !helper.MustBindJSON(&req) {
		return
	}

	adminID := helper.MustUserID()
	if adminID == 0 {
		return
	}

	response, errInfo, err := ctl.holidayService.BulkCreateHolidays(ctx.Request.Context(), centerID, adminID, &req)
	if err != nil {
		helper.ErrorWithInfo(errInfo)
		return
	}

	helper.Success(response)
}

// DeleteHoliday 刪除假日
// @Summary 刪除假日
// @Tags Admin
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Center ID"
// @Param holiday_id path int true "Holiday ID"
// @Success 200 {object} global.ApiResponse
// @Router /api/v1/admin/centers/{id}/holidays/{holiday_id} [delete]
func (ctl *AdminHolidayController) DeleteHoliday(ctx *gin.Context) {
	helper := NewContextHelper(ctx)

	centerID := helper.MustParamUint("id")
	if centerID == 0 {
		return
	}

	holidayID := helper.MustParamUint("holiday_id")
	if holidayID == 0 {
		return
	}

	adminID := helper.MustUserID()
	if adminID == 0 {
		return
	}

	errInfo, err := ctl.holidayService.DeleteHoliday(ctx.Request.Context(), centerID, adminID, holidayID)
	if err != nil {
		helper.ErrorWithInfo(errInfo)
		return
	}

	helper.Success(nil)
}
