package controllers

import (
	"net/http"
	"time"
	"timeLedger/app"
	"timeLedger/app/models"
	"timeLedger/app/repositories"
	"timeLedger/app/services"
	"timeLedger/global"

	"github.com/gin-gonic/gin"
)

// TimetableTemplateController 課表模板控制器
type TimetableTemplateController struct {
	app            *app.App
	templateRepo   *repositories.TimetableTemplateRepository
	cellRepo       *repositories.TimetableCellRepository
	auditLogRepo   *repositories.AuditLogRepository
	templateService *services.TimetableTemplateService
}

// NewTimetableTemplateController 建立 TimetableTemplateController 實例
func NewTimetableTemplateController(appInstance *app.App) *TimetableTemplateController {
	return &TimetableTemplateController{
		app:              appInstance,
		templateRepo:     repositories.NewTimetableTemplateRepository(appInstance),
		cellRepo:         repositories.NewTimetableCellRepository(appInstance),
		auditLogRepo:     repositories.NewAuditLogRepository(appInstance),
		templateService:  services.NewTimetableTemplateService(appInstance),
	}
}

// GetTemplates 取得課表模板列表
// @Summary 取得課表模板列表
// @Tags Admin
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} global.ApiResponse{data=[]models.TimetableTemplate}
// @Router /api/v1/admin/templates [get]
func (c *TimetableTemplateController) GetTemplates(ctx *gin.Context) {
	helper := NewContextHelper(ctx)

	centerID := helper.MustCenterID()
	if centerID == 0 {
		return
	}

	templates, err := c.templateRepo.ListByCenterID(ctx.Request.Context(), centerID)
	if err != nil {
		helper.InternalError("Failed to get templates")
		return
	}

	helper.Success(templates)
}

// CreateTemplateRequest 建立模板請求
type CreateTemplateRequest struct {
	Name    string `json:"name" binding:"required"`
	RowType string `json:"row_type"`
}

// CreateTemplate 新增課表模板
// @Summary 新增課表模板
// @Tags Admin
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body CreateTemplateRequest true "模板資訊"
// @Success 201 {object} global.ApiResponse{data=models.TimetableTemplate}
// @Router /api/v1/admin/templates [post]
func (c *TimetableTemplateController) CreateTemplate(ctx *gin.Context) {
	helper := NewContextHelper(ctx)

	centerID := helper.MustCenterID()
	if centerID == 0 {
		return
	}

	adminID := helper.MustUserID()
	if adminID == 0 {
		return
	}

	var req CreateTemplateRequest
	if !helper.MustBindJSON(&req) {
		return
	}

	template := models.TimetableTemplate{
		CenterID:  centerID,
		Name:      req.Name,
		RowType:   req.RowType,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	createdTemplate, err := c.templateRepo.Create(ctx.Request.Context(), template)
	if err != nil {
		helper.InternalError("Failed to create template")
		return
	}

	c.auditLogRepo.Create(ctx, models.AuditLog{
		CenterID:   centerID,
		ActorType:  "ADMIN",
		ActorID:    adminID,
		Action:     "TEMPLATE_CREATE",
		TargetType: "TimetableTemplate",
		TargetID:   createdTemplate.ID,
		Payload: models.AuditPayload{
			After: map[string]interface{}{
				"name":     req.Name,
				"row_type": req.RowType,
			},
		},
	})

	helper.Created(createdTemplate)
}

// UpdateTemplateRequest 更新模板請求
type UpdateTemplateRequest struct {
	Name string `json:"name"`
}

// UpdateTemplate 更新課表模板
// @Summary 更新課表模板
// @Tags Admin
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param template_id path int true "Template ID"
// @Param request body UpdateTemplateRequest true "模板資訊"
// @Success 200 {object} global.ApiResponse{data=models.TimetableTemplate}
// @Router /api/v1/admin/templates/{template_id} [put]
func (c *TimetableTemplateController) UpdateTemplate(ctx *gin.Context) {
	helper := NewContextHelper(ctx)

	centerID := helper.MustCenterID()
	if centerID == 0 {
		return
	}

	adminID := helper.MustUserID()
	if adminID == 0 {
		return
	}

	templateID := helper.MustParamUint("templateId")
	if templateID == 0 {
		return
	}

	var req UpdateTemplateRequest
	if !helper.MustBindJSON(&req) {
		return
	}

	template := models.TimetableTemplate{
		ID:        templateID,
		CenterID:  centerID,
		Name:      req.Name,
		UpdatedAt: time.Now(),
	}

	if err := c.templateRepo.Update(ctx.Request.Context(), template); err != nil {
		helper.InternalError("Failed to update template")
		return
	}

	c.auditLogRepo.Create(ctx, models.AuditLog{
		CenterID:   centerID,
		ActorType:  "ADMIN",
		ActorID:    adminID,
		Action:     "TEMPLATE_UPDATE",
		TargetType: "TimetableTemplate",
		TargetID:   templateID,
		Payload: models.AuditPayload{
			After: map[string]interface{}{
				"name": req.Name,
			},
		},
	})

	helper.Success(template)
}

// GetCells 取得模板中的格子
// @Summary 取得模板中的格子
// @Tags Admin
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param template_id path int true "Template ID"
// @Success 200 {object} global.ApiResponse{data=[]models.TimetableCell}
// @Router /api/v1/admin/templates/{template_id}/cells [get]
func (c *TimetableTemplateController) GetCells(ctx *gin.Context) {
	helper := NewContextHelper(ctx)

	templateID := helper.MustParamUint("templateId")
	if templateID == 0 {
		return
	}

	cells, err := c.cellRepo.ListByTemplateID(ctx.Request.Context(), templateID)
	if err != nil {
		helper.InternalError("Failed to get cells")
		return
	}

	helper.Success(cells)
}

// CreateCellRequest 建立格子請求
type CreateCellRequest struct {
	RowNo     int     `json:"row_no"`
	ColNo     int     `json:"col_no"`
	StartTime string  `json:"start_time"`
	EndTime   string  `json:"end_time"`
	RoomID    *uint   `json:"room_id"`
	TeacherID *uint   `json:"teacher_id"`
}

// CreateCells 新增模板中的格子
// @Summary 新增模板中的格子
// @Tags Admin
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param template_id path int true "Template ID"
// @Param request body []CreateCellRequest true "格子資訊"
// @Success 201 {object} global.ApiResponse{data=[]models.TimetableCell}
// @Router /api/v1/admin/templates/{template_id}/cells [post]
func (c *TimetableTemplateController) CreateCells(ctx *gin.Context) {
	helper := NewContextHelper(ctx)

	centerID := helper.MustCenterID()
	if centerID == 0 {
		return
	}

	adminID := helper.MustUserID()
	if adminID == 0 {
		return
	}

	templateID := helper.MustParamUint("templateId")
	if templateID == 0 {
		return
	}

	var req []CreateCellRequest
	if !helper.MustBindJSON(&req) {
		return
	}

	var cells []models.TimetableCell
	for _, cellReq := range req {
		cell := models.TimetableCell{
			TemplateID: templateID,
			RowNo:      cellReq.RowNo,
			ColNo:      cellReq.ColNo,
			StartTime:  cellReq.StartTime,
			EndTime:    cellReq.EndTime,
			RoomID:     cellReq.RoomID,
			TeacherID:  cellReq.TeacherID,
			CreatedAt:  time.Now(),
			UpdatedAt:  time.Now(),
		}
		cells = append(cells, cell)
	}

	// 批次建立所有格子
	var createdCells []models.TimetableCell
	for _, cell := range cells {
		createdCell, err := c.cellRepo.Create(ctx.Request.Context(), cell)
		if err != nil {
			helper.InternalError("Failed to create cell: " + err.Error())
			return
		}
		createdCells = append(createdCells, createdCell)
	}

	c.auditLogRepo.Create(ctx, models.AuditLog{
		CenterID:   centerID,
		ActorType:  "ADMIN",
		ActorID:    adminID,
		Action:     "TIMETABLE_CELLS_CREATE",
		TargetType: "TimetableCell",
		TargetID:   0,
		Payload: models.AuditPayload{
			After: map[string]interface{}{
				"template_id": templateID,
				"cells_count": len(createdCells),
			},
		},
	})

	helper.Created(createdCells)
}

// DeleteCell 刪除格子
// @Summary 刪除格子
// @Tags Admin
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param template_id path int true "Template ID"
// @Param cell_id path int true "Cell ID"
// @Success 200 {object} global.ApiResponse
// @Router /api/v1/admin/templates/{template_id}/cells/{cell_id} [delete]
func (c *TimetableTemplateController) DeleteCell(ctx *gin.Context) {
	helper := NewContextHelper(ctx)

	centerID := helper.MustCenterID()
	if centerID == 0 {
		return
	}

	adminID := helper.MustUserID()
	if adminID == 0 {
		return
	}

	cellID := helper.MustParamUint("cellId")
	if cellID == 0 {
		return
	}

	// 取得 cell 來確認存在
	cell, err := c.cellRepo.GetByID(ctx.Request.Context(), cellID)
	if err != nil {
		helper.NotFound("Cell not found")
		return
	}

	// 刪除格子
	if err := c.cellRepo.Delete(ctx.Request.Context(), cellID); err != nil {
		helper.InternalError("Failed to delete cell: " + err.Error())
		return
	}

	c.auditLogRepo.Create(ctx, models.AuditLog{
		CenterID:   centerID,
		ActorType:  "ADMIN",
		ActorID:    adminID,
		Action:     "TIMETABLE_CELL_DELETE",
		TargetType: "TimetableCell",
		TargetID:   cellID,
		Payload: models.AuditPayload{
			Before: map[string]interface{}{
				"template_id": cell.TemplateID,
				"row_no":      cell.RowNo,
				"col_no":      cell.ColNo,
			},
		},
	})

	helper.Success(nil)
}

// DeleteTemplate 刪除課表模板
// @Summary 刪除課表模板
// @Tags Admin
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param template_id path int true "Template ID"
// @Success 200 {object} global.ApiResponse
// @Router /api/v1/admin/templates/{template_id} [delete]
func (c *TimetableTemplateController) DeleteTemplate(ctx *gin.Context) {
	helper := NewContextHelper(ctx)

	centerID := helper.MustCenterID()
	if centerID == 0 {
		return
	}

	adminID := helper.MustUserID()
	if adminID == 0 {
		return
	}

	templateID := helper.MustParamUint("templateId")
	if templateID == 0 {
		return
	}

	if err := c.templateRepo.DeleteByIDAndCenterID(ctx.Request.Context(), templateID, centerID); err != nil {
		helper.InternalError("Failed to delete template")
		return
	}

	c.auditLogRepo.Create(ctx, models.AuditLog{
		CenterID:   centerID,
		ActorType:  "ADMIN",
		ActorID:    adminID,
		Action:     "TEMPLATE_DELETE",
		TargetType: "TimetableTemplate",
		TargetID:   templateID,
		Payload: models.AuditPayload{
			After: map[string]interface{}{
				"status": "DELETED",
			},
		},
	})

	helper.Success(nil)
}

// ApplyTemplateRequest 套用模板請求
type ApplyTemplateRequest struct {
	OfferingID     uint     `json:"offering_id" binding:"required"`
	StartDate      string   `json:"start_date" binding:"required"`
	EndDate        string   `json:"end_date" binding:"required"`
	Weekdays       []int    `json:"weekdays" binding:"required"`
	Duration       int      `json:"duration"`
	OverrideBuffer bool     `json:"override_buffer"`
}

// ApplyTemplateConflictInfo 套用模板衝突資訊
type ApplyTemplateConflictInfo struct {
	Weekday      int    `json:"weekday"`
	StartTime    string `json:"start_time"`
	EndTime      string `json:"end_time"`
	ConflictType string `json:"conflict_type"`
	Message      string `json:"message"`
	RuleID       uint   `json:"rule_id,omitempty"`
	CanOverride  bool   `json:"can_override,omitempty"`
}

// ApplyTemplate 套用課表模板
// @Summary 套用課表模板
// @Tags Admin
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param template_id path int true "Template ID"
// @Param request body ApplyTemplateRequest true "套用資訊"
// @Success 200 {object} global.ApiResponse
// @Router /api/v1/admin/templates/{template_id}/apply [post]
func (c *TimetableTemplateController) ApplyTemplate(ctx *gin.Context) {
	helper := NewContextHelper(ctx)

	centerID := helper.MustCenterID()
	if centerID == 0 {
		return
	}

	adminID := helper.MustUserID()
	if adminID == 0 {
		return
	}

	templateID := helper.MustParamUint("templateId")
	if templateID == 0 {
		return
	}

	var req ApplyTemplateRequest
	if !helper.MustBindJSON(&req) {
		return
	}

	// 呼叫 Service 層
	result, errInfo, err := c.templateService.ApplyTemplate(ctx.Request.Context(), &services.ApplyTemplateInput{
		TemplateID:     templateID,
		CenterID:       centerID,
		AdminID:        adminID,
		OfferingID:     req.OfferingID,
		StartDate:      req.StartDate,
		EndDate:        req.EndDate,
		Weekdays:       req.Weekdays,
		Duration:       req.Duration,
		OverrideBuffer: req.OverrideBuffer,
	})

	if err != nil {
		helper.ErrorWithInfo(errInfo)
		return
	}

	// 檢查衝突
	if !result.Valid {
		nonOverrideConflicts := 0
		for _, conflict := range result.Conflicts {
			if !conflict.CanOverride {
				nonOverrideConflicts++
			}
		}

		// 如果有不可覆蓋的衝突
		if nonOverrideConflicts > 0 {
			helper.ctx.JSON(http.StatusBadRequest, global.ApiResponse{
				Code:    40002, // OVERLAP error code
				Message: "套用模板會產生時間衝突，請先解決衝突後再嘗試",
				Datas: map[string]interface{}{
					"conflicts":      result.Conflicts,
					"conflict_count": len(result.Conflicts),
				},
			})
			return
		}

		// 只有可覆蓋的衝突（Buffer 衝突）
		helper.ctx.JSON(http.StatusBadRequest, global.ApiResponse{
			Code:    40003, // BUFFER_CONFLICT warning code
			Message: "套用模板會產生緩衝時間衝突，是否繼續？",
			Datas: map[string]interface{}{
				"conflicts":      result.Conflicts,
				"conflict_count": len(result.Conflicts),
				"can_override":   true,
			},
		})
		return
	}

	helper.Success(result.Output)
}

// ValidateApplyTemplateRequest 驗證套用模板請求
type ValidateApplyTemplateRequest struct {
	OfferingID     uint     `json:"offering_id" binding:"required"`
	StartDate      string   `json:"start_date" binding:"required"`
	EndDate        string   `json:"end_date" binding:"required"`
	Weekdays       []int    `json:"weekdays" binding:"required"`
	OverrideBuffer bool     `json:"override_buffer"`
}

// ValidateApplyTemplate 驗證套用模板（不實際產生規則）
// @Summary 驗證套用模板
// @Tags Admin
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param template_id path int true "Template ID"
// @Param request body ValidateApplyTemplateRequest true "驗證資訊"
// @Success 200 {object} global.ApiResponse
// @Router /api/v1/admin/templates/{template_id}/validate-apply [post]
func (c *TimetableTemplateController) ValidateApplyTemplate(ctx *gin.Context) {
	helper := NewContextHelper(ctx)

	centerID := helper.MustCenterID()
	if centerID == 0 {
		return
	}

	templateID := helper.MustParamUint("templateId")
	if templateID == 0 {
		return
	}

	var req ValidateApplyTemplateRequest
	if !helper.MustBindJSON(&req) {
		return
	}

	// 呼叫 Service 層進行驗證
	result, errInfo, err := c.templateService.ValidateApplyTemplate(ctx.Request.Context(), &services.ApplyTemplateValidateInput{
		TemplateID:     templateID,
		CenterID:       centerID,
		OfferingID:     req.OfferingID,
		StartDate:      req.StartDate,
		EndDate:        req.EndDate,
		Weekdays:       req.Weekdays,
		OverrideBuffer: req.OverrideBuffer,
	})

	if err != nil {
		helper.ErrorWithInfo(errInfo)
		return
	}

	helper.Success(map[string]interface{}{
		"valid":     result.Valid,
		"conflicts": result.Conflicts,
	})
}
