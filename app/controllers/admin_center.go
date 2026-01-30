package controllers

import (
	"timeLedger/app"
	"timeLedger/app/services"

	"github.com/gin-gonic/gin"
)

// AdminCenterController 中心管理相關 API
type AdminCenterController struct {
	app          *app.App
	centerService *services.CenterService
}

// NewAdminCenterController 建立 AdminCenterController 實例
func NewAdminCenterController(appInstance *app.App) *AdminCenterController {
	return &AdminCenterController{
		app:          appInstance,
		centerService: services.NewCenterService(appInstance),
	}
}

// GetCenters 取得中心列表
// @Summary 取得中心列表
// @Tags Admin
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} global.ApiResponse{data=[]models.Center}
// @Router /api/v1/admin/centers [get]
func (ctl *AdminCenterController) GetCenters(ctx *gin.Context) {
	helper := NewContextHelper(ctx)

	centers, errInfo, err := ctl.centerService.ListCenters(ctx.Request.Context())
	if err != nil {
		helper.ErrorWithInfo(errInfo)
		return
	}

	helper.Success(centers)
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
// @Success 200 {object} global.ApiResponse{data=models.Center}
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

	helper.Created(center)
}
