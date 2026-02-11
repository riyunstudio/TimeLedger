package controllers

import (
	"fmt"
	"timeLedger/app"
	"timeLedger/app/resources"
	"timeLedger/app/services"

	"github.com/gin-gonic/gin"
)

// OfferingController 班別管理控制器
type OfferingController struct {
	app              *app.App
	offeringService  *services.OfferingService
	offeringResource *resources.OfferingResource
}

// NewOfferingController 建立 OfferingController 實例
func NewOfferingController(appInstance *app.App) *OfferingController {
	return &OfferingController{
		app:              appInstance,
		offeringService:  services.NewOfferingService(appInstance),
		offeringResource: resources.NewOfferingResource(appInstance),
	}
}

// GetOfferings 取得班別列表
// @Summary 取得班別列表
// @Tags Admin
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Items per page" default(20)
// @Success 200 {object} global.ApiResponse{data=services.ListOfferingsOutput}
// @Router /api/v1/admin/offerings [get]
func (c *OfferingController) GetOfferings(ctx *gin.Context) {
	helper := NewContextHelper(ctx)

	centerID := helper.MustCenterID()
	if centerID == 0 {
		return
	}

	page := helper.QueryStringOrDefault("page", "1")
	limit := helper.QueryStringOrDefault("limit", "20")

	var pageInt, limitInt int
	if _, err := fmt.Sscanf(page, "%d", &pageInt); err != nil || pageInt < 1 {
		pageInt = 1
	}
	if _, err := fmt.Sscanf(limit, "%d", &limitInt); err != nil || limitInt < 1 || limitInt > 100 {
		limitInt = 20
	}

	result, errInfo, err := c.offeringService.ListOfferings(ctx.Request.Context(), &services.ListOfferingsInput{
		CenterID: centerID,
		Page:     pageInt,
		Limit:    limitInt,
	})
	if err != nil {
		helper.ErrorWithInfo(errInfo)
		return
	}

	helper.Success(result)
}

// CreateOfferingRequest 新增班別請求
type CreateOfferingRequest struct {
	CourseID            uint  `json:"course_id" binding:"required"`
	DefaultRoomID       *uint `json:"default_room_id"`
	DefaultTeacherID    *uint `json:"default_teacher_id"`
	AllowBufferOverride bool  `json:"allow_buffer_override"`
}

// CreateOffering 新增班別
// @Summary 新增班別
// @Tags Admin
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body CreateOfferingRequest true "班別資訊"
// @Success 201 {object} global.ApiResponse{data=models.Offering}
// @Router /api/v1/admin/offerings [post]
func (c *OfferingController) CreateOffering(ctx *gin.Context) {
	helper := NewContextHelper(ctx)

	centerID := helper.MustCenterID()
	if centerID == 0 {
		return
	}

	adminID := helper.MustUserID()
	if adminID == 0 {
		return
	}

	var req CreateOfferingRequest
	if !helper.MustBindJSON(&req) {
		return
	}

	result, errInfo, err := c.offeringService.CreateOffering(ctx.Request.Context(), &services.CreateOfferingInput{
		CenterID:            centerID,
		AdminID:             adminID,
		CourseID:            req.CourseID,
		DefaultRoomID:       req.DefaultRoomID,
		DefaultTeacherID:    req.DefaultTeacherID,
		AllowBufferOverride: req.AllowBufferOverride,
	})
	if err != nil {
		helper.ErrorWithInfo(errInfo)
		return
	}

	helper.Created(result)
}

// UpdateOfferingRequest 更新班別請求
type UpdateOfferingRequest struct {
	Name                *string `json:"name"`
	DefaultRoomID       *uint   `json:"default_room_id"`
	DefaultTeacherID    *uint   `json:"default_teacher_id"`
	AllowBufferOverride bool    `json:"allow_buffer_override"`
}

// UpdateOffering 更新班別
// @Summary 更新班別
// @Tags Admin
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param offering_id path int true "Offering ID"
// @Param request body UpdateOfferingRequest true "班別資訊"
// @Success 200 {object} global.ApiResponse{data=models.Offering}
// @Router /api/v1/admin/offerings/{offering_id} [put]
func (c *OfferingController) UpdateOffering(ctx *gin.Context) {
	helper := NewContextHelper(ctx)

	centerID := helper.MustCenterID()
	if centerID == 0 {
		return
	}

	adminID := helper.MustUserID()
	if adminID == 0 {
		return
	}

	offeringID := helper.MustParamUint("offering_id")
	if offeringID == 0 {
		return
	}

	var req UpdateOfferingRequest
	if !helper.MustBindJSON(&req) {
		return
	}

	result, errInfo, err := c.offeringService.UpdateOffering(ctx.Request.Context(), &services.UpdateOfferingInput{
		CenterID:            centerID,
		AdminID:             adminID,
		OfferingID:          offeringID,
		Name:                req.Name,
		DefaultRoomID:       req.DefaultRoomID,
		DefaultTeacherID:    req.DefaultTeacherID,
		AllowBufferOverride: req.AllowBufferOverride,
	})
	if err != nil {
		helper.ErrorWithInfo(errInfo)
		return
	}

	helper.Success(result)
}

// DeleteOffering 刪除班別
// @Summary 刪除班別
// @Tags Admin
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param offering_id path int true "Offering ID"
// @Success 204 {object} global.ApiResponse
// @Router /api/v1/admin/offerings/{offering_id} [delete]
func (c *OfferingController) DeleteOffering(ctx *gin.Context) {
	helper := NewContextHelper(ctx)

	centerID := helper.MustCenterID()
	if centerID == 0 {
		return
	}

	adminID := helper.MustUserID()
	if adminID == 0 {
		return
	}

	offeringID := helper.MustParamUint("offering_id")
	if offeringID == 0 {
		return
	}

	errInfo := c.offeringService.DeleteOffering(ctx.Request.Context(), &services.DeleteOfferingInput{
		CenterID: centerID,
		AdminID:  adminID,
		ID:       offeringID,
	})
	if errInfo != nil {
		helper.ErrorWithInfo(errInfo)
		return
	}

	helper.NoContent()
}

// CopyOfferingRequest 複製班別請求
type CopyOfferingRequest struct {
	NewName     string `json:"new_name" binding:"required"`
	CopyTeacher bool   `json:"copy_teacher"`
}

// CopyOfferingResponse 複製班別響應
type CopyOfferingResponse struct {
	ID          uint   `json:"id"`
	Name        string `json:"name"`
	CourseID    uint   `json:"course_id"`
	RulesCopied int    `json:"rules_copied"`
}

// CopyOffering 複製班別
// @Summary 複製班別
// @Tags Admin
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Center ID"
// @Param offering_id path int true "Offering ID"
// @Param request body CopyOfferingRequest true "複製資訊"
// @Success 201 {object} global.ApiResponse{data=CopyOfferingResponse}
// @Router /api/v1/admin/centers/{id}/offerings/{offering_id}/copy [post]
func (c *OfferingController) CopyOffering(ctx *gin.Context) {
	helper := NewContextHelper(ctx)

	centerID := helper.MustParamUint("id")
	if centerID == 0 {
		return
	}

	adminID := helper.MustUserID()
	if adminID == 0 {
		return
	}

	offeringID := helper.MustParamUint("offering_id")
	if offeringID == 0 {
		return
	}

	var req CopyOfferingRequest
	if !helper.MustBindJSON(&req) {
		return
	}

	result, errInfo, err := c.offeringService.CopyOffering(ctx.Request.Context(), &services.CopyOfferingInput{
		CenterID:    centerID,
		AdminID:     adminID,
		OfferingID:  offeringID,
		NewName:     req.NewName,
		CopyTeacher: req.CopyTeacher,
	})
	if err != nil {
		helper.ErrorWithInfo(errInfo)
		return
	}

	helper.Created(CopyOfferingResponse{
		ID:          result.ID,
		Name:        result.Name,
		CourseID:    result.CourseID,
		RulesCopied: result.RulesCopied,
	})
}

// GetActiveOfferings 取得啟用的班別列表
// @Summary 取得啟用的班別列表
// @Tags Admin
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} global.ApiResponse{data=[]resources.OfferingResponse}
// @Router /api/v1/admin/offerings/active [get]
func (c *OfferingController) GetActiveOfferings(ctx *gin.Context) {
	helper := NewContextHelper(ctx)

	centerID := helper.MustCenterID()
	if centerID == 0 {
		return
	}

	offerings, err := c.offeringService.GetActiveOfferings(ctx.Request.Context(), &services.GetActiveOfferingsInput{
		CenterID: centerID,
	})
	if err != nil {
		helper.InternalError("Failed to get active offerings")
		return
	}

	// 轉換為含課程時長的響應格式
	response := c.offeringResource.ToOfferingResponses(offerings)
	helper.Success(response)
}

// ToggleActiveRequest 切換狀態請求
type ToggleActiveRequest struct {
	IsActive bool `json:"is_active"`
}

// ToggleOfferingActive 切換班別啟用狀態
// @Summary 切換班別啟用狀態
// @Tags Admin
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param offering_id path int true "Offering ID"
// @Param request body ToggleActiveRequest true "狀態資訊"
// @Success 200 {object} global.ApiResponse
// @Router /api/v1/admin/offerings/{offering_id}/toggle-active [patch]
func (c *OfferingController) ToggleOfferingActive(ctx *gin.Context) {
	helper := NewContextHelper(ctx)

	centerID := helper.MustCenterID()
	if centerID == 0 {
		return
	}

	adminID := helper.MustUserID()
	if adminID == 0 {
		return
	}

	offeringID := helper.MustParamUint("offering_id")
	if offeringID == 0 {
		return
	}

	var req ToggleActiveRequest
	if !helper.MustBindJSON(&req) {
		return
	}

	errInfo := c.offeringService.ToggleOfferingActive(ctx.Request.Context(), &services.ToggleOfferingActiveInput{
		CenterID: centerID,
		AdminID:  adminID,
		ID:       offeringID,
		IsActive: req.IsActive,
	})
	if errInfo != nil {
		helper.ErrorWithInfo(errInfo)
		return
	}

	helper.Success(nil)
}
