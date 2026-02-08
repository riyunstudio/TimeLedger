package controllers

import (
	"strconv"

	"timeLedger/app"
	"timeLedger/app/resources"
	"timeLedger/app/services"

	"github.com/gin-gonic/gin"
)

// AdminCenterController 中心管理相關 API
type AdminCenterController struct {
	app            *app.App
	centerService  *services.CenterService
	centerResource *resources.CenterResource
}

// NewAdminCenterController 建立 AdminCenterController 實例
func NewAdminCenterController(appInstance *app.App) *AdminCenterController {
	return &AdminCenterController{
		app:            appInstance,
		centerService:  services.NewCenterService(appInstance),
		centerResource: resources.NewCenterResource(appInstance),
	}
}

// UpdateSettingsRequest 更新中心設定的請求結構
type UpdateSettingsRequest struct {
	DefaultCourseDuration *int64 `json:"default_course_duration"`
}

// UpdateSettings 更新中心設定
// @Summary 更新中心設定
// @Tags Admin
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param center_id path int true "中心 ID"
// @Param request body UpdateSettingsRequest true "設定資訊"
// @Success 200 {object} global.ApiResponse{data=resources.CenterResponse}
// @Router /api/v1/admin/centers/{center_id}/settings [patch]
func (ctl *AdminCenterController) UpdateSettings(ctx *gin.Context) {
	helper := NewContextHelper(ctx)

	centerID, err := strconv.ParseUint(ctx.Param("center_id"), 10, 32)
	if err != nil || centerID == 0 {
		helper.BadRequest("無效的中心 ID")
		return
	}

	var req UpdateSettingsRequest
	if !helper.MustBindJSON(&req) {
		return
	}

	// 取得現有設定
	settings, errInfo, err := ctl.centerService.GetCenterSettings(ctx.Request.Context(), uint(centerID))
	if err != nil {
		helper.ErrorWithInfo(errInfo)
		return
	}

	// 更新設定
	if req.DefaultCourseDuration != nil {
		settings.DefaultCourseDuration = int(*req.DefaultCourseDuration)
	}

	// 取得管理員 ID
	adminID := helper.MustUserID()
	if adminID == 0 {
		return
	}

	// 儲存變更
	updatedCenter, errInfo, err := ctl.centerService.UpdateCenterSettings(ctx.Request.Context(), uint(centerID), adminID, settings)
	if err != nil {
		helper.ErrorWithInfo(errInfo)
		return
	}

	response := ctl.centerResource.ToCenterResponse(*updatedCenter)
	helper.Success(response)
}

// GetSettings 取得中心設定
// @Summary 取得中心設定
// @Tags Admin
// @Produce json
// @Security BearerAuth
// @Param center_id path int true "中心 ID"
// @Success 200 {object} global.ApiResponse{data=resources.CenterSettingsResponse}
// @Router /api/v1/admin/centers/{center_id}/settings [get]
func (ctl *AdminCenterController) GetSettings(ctx *gin.Context) {
	helper := NewContextHelper(ctx)

	centerID, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil || centerID == 0 {
		helper.BadRequest("無效的中心 ID")
		return
	}

	settings, errInfo, err := ctl.centerService.GetCenterSettings(ctx.Request.Context(), uint(centerID))
	if err != nil {
		helper.ErrorWithInfo(errInfo)
		return
	}

	response := ctl.centerResource.ToSettingsResponse(*settings)
	helper.Success(response)
}

// GetCenters 取得中心列表
// @Summary 取得中心列表
// @Tags Admin
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} global.ApiResponse{data=[]resources.CenterResponse}
// @Router /api/v1/admin/centers [get]
func (ctl *AdminCenterController) GetCenters(ctx *gin.Context) {
	helper := NewContextHelper(ctx)

	centers, errInfo, err := ctl.centerService.ListCenters(ctx.Request.Context())
	if err != nil {
		helper.ErrorWithInfo(errInfo)
		return
	}

	responses := ctl.centerResource.ToCenterResponses(centers)
	helper.Success(responses)
}

// CreateCenterRequest 新增中心的請求結構
type CreateCenterRequest struct {
	Name                string `json:"name" binding:"required"`
	PlanLevel           string `json:"plan_level" binding:"required"`
	AllowPublicRegister bool   `json:"allow_public_register"`
}

// CreateCenter 新增中心
// @Summary 新增中心
// @Tags Admin
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body CreateCenterRequest true "中心資訊"
// @Success 200 {object} global.ApiResponse{data=resources.CenterResponse}
// @Router /api/v1/admin/centers [post]
func (ctl *AdminCenterController) CreateCenter(ctx *gin.Context) {
	helper := NewContextHelper(ctx)

	var req CreateCenterRequest
	if !helper.MustBindJSON(&req) {
		return
	}

	adminID := helper.MustUserID()
	if adminID == 0 {
		return
	}

	input := &services.CreateCenterInput{
		Name:                req.Name,
		PlanLevel:           req.PlanLevel,
		AllowPublicRegister: req.AllowPublicRegister,
	}

	center, errInfo, err := ctl.centerService.CreateCenter(ctx.Request.Context(), adminID, input)
	if err != nil {
		helper.ErrorWithInfo(errInfo)
		return
	}

	response := ctl.centerResource.ToCenterResponse(*center)
	helper.Created(response)
}
