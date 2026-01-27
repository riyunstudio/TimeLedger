package controllers

import (
	"fmt"
	"net/http"
	"time"
	"timeLedger/app"
	"timeLedger/app/models"
	"timeLedger/app/repositories"
	"timeLedger/app/requests"
	"timeLedger/app/services"
	"timeLedger/global"

	"github.com/gin-gonic/gin"
)

type TimetableTemplateController struct {
	BaseController
	templateRepository    *repositories.TimetableTemplateRepository
	cellRepository        *repositories.TimetableCellRepository
	scheduleRuleRepo      *repositories.ScheduleRuleRepository
	personalEventRepo     *repositories.PersonalEventRepository
	auditLogRepo          *repositories.AuditLogRepository
	ruleValidator         *services.ScheduleRuleValidator
}

func NewTimetableTemplateController(app *app.App) *TimetableTemplateController {
	return &TimetableTemplateController{
		templateRepository: repositories.NewTimetableTemplateRepository(app),
		cellRepository:     repositories.NewTimetableCellRepository(app),
		scheduleRuleRepo:   repositories.NewScheduleRuleRepository(app),
		personalEventRepo:  repositories.NewPersonalEventRepository(app),
		auditLogRepo:       repositories.NewAuditLogRepository(app),
		ruleValidator:      services.NewScheduleRuleValidator(app),
	}
}

// getCenterID 從 JWT Token 取得 center_id
// 如果找不到或為 0，回傳錯誤
func (ctl *TimetableTemplateController) getCenterID(ctx *gin.Context) (uint, bool) {
	centerID, exists := ctx.Get(global.CenterIDKey)
	if !exists {
		ctx.JSON(http.StatusUnauthorized, global.ApiResponse{
			Code:    global.UNAUTHORIZED,
			Message: "Center ID not found in token",
		})
		return 0, false
	}

	centerIDUint := centerID.(uint)
	if centerIDUint == 0 {
		ctx.JSON(http.StatusForbidden, global.ApiResponse{
			Code:    global.FORBIDDEN,
			Message: "Center ID is required",
		})
		return 0, false
	}

	return centerIDUint, true
}

func (ctl *TimetableTemplateController) GetTemplates(ctx *gin.Context) {
	centerID, ok := ctl.getCenterID(ctx)
	if !ok {
		return
	}

	templates, err := ctl.templateRepository.ListByCenterID(ctl.makeCtx(ctx), centerID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, global.ApiResponse{
			Code:    500,
			Message: "Failed to get templates",
		})
		return
	}

	ctx.JSON(http.StatusOK, global.ApiResponse{
		Code:    0,
		Message: "Success",
		Datas:   templates,
	})
}

func (ctl *TimetableTemplateController) CreateTemplate(ctx *gin.Context) {
	centerID, ok := ctl.getCenterID(ctx)
	if !ok {
		return
	}

	var req requests.CreateTemplateRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, global.ApiResponse{
			Code:    global.BAD_REQUEST,
			Message: "Invalid request body: " + err.Error(),
		})
		return
	}

	req.CenterID = centerID

	template := models.TimetableTemplate{
		CenterID:  req.CenterID,
		Name:      req.Name,
		RowType:   req.RowType,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	createdTemplate, err := ctl.templateRepository.Create(ctl.makeCtx(ctx), template)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, global.ApiResponse{
			Code:    500,
			Message: "Failed to create template",
		})
		return
	}

	actorID := ctx.GetUint(global.UserIDKey)
	ctl.auditLogRepo.Create(ctx, models.AuditLog{
		CenterID:   req.CenterID,
		ActorType:  "ADMIN",
		ActorID:    actorID,
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

	ctx.JSON(http.StatusOK, global.ApiResponse{
		Code:    0,
		Message: "Template created",
		Datas:   createdTemplate,
	})
}

func (ctl *TimetableTemplateController) UpdateTemplate(ctx *gin.Context) {
	centerID, ok := ctl.getCenterID(ctx)
	if !ok {
		return
	}

	allParams := make(map[string]string)
	for _, param := range ctx.Params {
		allParams[param.Key] = param.Value
	}

	templateIDStr := ctx.Param("templateId")

	if templateIDStr == "" {
		ctx.JSON(http.StatusBadRequest, global.ApiResponse{
			Code:    global.BAD_REQUEST,
			Message: "Template ID required",
			Datas:   allParams,
		})
		return
	}

	var templateID uint
	if _, err := fmt.Sscanf(templateIDStr, "%d", &templateID); err != nil || templateID == 0 {
		ctx.JSON(http.StatusBadRequest, global.ApiResponse{
			Code:    global.BAD_REQUEST,
			Message: "Invalid template ID format",
		})
		return
	}

	var req requests.UpdateTemplateRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, global.ApiResponse{
			Code:    global.BAD_REQUEST,
			Message: "Invalid request body",
		})
		return
	}

	template := models.TimetableTemplate{
		ID:        templateID,
		CenterID:  centerID,
		Name:      req.Name,
		UpdatedAt: time.Now(),
	}

	if err := ctl.templateRepository.UpdateByIDAndCenterID(ctl.makeCtx(ctx), templateID, centerID, template); err != nil {
		ctx.JSON(http.StatusInternalServerError, global.ApiResponse{
			Code:    500,
			Message: "Failed to update template",
		})
		return
	}

	actorID := ctx.GetUint(global.UserIDKey)
	ctl.auditLogRepo.Create(ctx, models.AuditLog{
		CenterID:   centerID,
		ActorType:  "ADMIN",
		ActorID:    actorID,
		Action:     "TEMPLATE_UPDATE",
		TargetType: "TimetableTemplate",
		TargetID:   templateID,
		Payload: models.AuditPayload{
			After: map[string]interface{}{
				"name": req.Name,
			},
		},
	})

	ctx.JSON(http.StatusOK, global.ApiResponse{
		Code:    0,
		Message: "Template updated",
		Datas:   template,
	})
}

func (ctl *TimetableTemplateController) GetCells(ctx *gin.Context) {
	templateIDStr := ctx.Param("templateId")

	if templateIDStr == "" {
		ctx.JSON(http.StatusBadRequest, global.ApiResponse{
			Code:    global.BAD_REQUEST,
			Message: "Template ID required",
		})
		return
	}

	var templateID uint
	if _, err := fmt.Sscanf(templateIDStr, "%d", &templateID); err != nil || templateID == 0 {
		ctx.JSON(http.StatusBadRequest, global.ApiResponse{
			Code:    global.BAD_REQUEST,
			Message: "Invalid template ID format",
		})
		return
	}

	cells, err := ctl.cellRepository.ListByTemplateID(ctl.makeCtx(ctx), templateID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, global.ApiResponse{
			Code:    500,
			Message: "Failed to get cells",
		})
		return
	}

	ctx.JSON(http.StatusOK, global.ApiResponse{
		Code:    0,
		Message: "Success",
		Datas:   cells,
	})
}

func (ctl *TimetableTemplateController) CreateCells(ctx *gin.Context) {
	centerID, ok := ctl.getCenterID(ctx)
	if !ok {
		return
	}

	templateIDStr := ctx.Param("templateId")

	if templateIDStr == "" {
		ctx.JSON(http.StatusBadRequest, global.ApiResponse{
			Code:    global.BAD_REQUEST,
			Message: "Template ID required",
		})
		return
	}

	var templateID uint
	if _, err := fmt.Sscanf(templateIDStr, "%d", &templateID); err != nil || templateID == 0 {
		ctx.JSON(http.StatusBadRequest, global.ApiResponse{
			Code:    global.BAD_REQUEST,
			Message: "Invalid template ID format",
		})
		return
	}

	var req []requests.CreateCellRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, global.ApiResponse{
			Code:    global.BAD_REQUEST,
			Message: "Invalid request body",
		})
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
		createdCell, err := ctl.cellRepository.Create(ctl.makeCtx(ctx), cell)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, global.ApiResponse{
				Code:    500,
				Message: "Failed to create cell: " + err.Error(),
			})
			return
		}
		createdCells = append(createdCells, createdCell)
	}

	actorID := ctx.GetUint(global.UserIDKey)
	ctl.auditLogRepo.Create(ctx, models.AuditLog{
		CenterID:   centerID,
		ActorType:  "ADMIN",
		ActorID:    actorID,
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

	ctx.JSON(http.StatusOK, global.ApiResponse{
		Code:    0,
		Message: "Cells created",
		Datas:   createdCells,
	})
}

func (ctl *TimetableTemplateController) DeleteCell(ctx *gin.Context) {
	centerID, ok := ctl.getCenterID(ctx)
	if !ok {
		return
	}

	cellIDStr := ctx.Param("cellId")

	if cellIDStr == "" {
		ctx.JSON(http.StatusBadRequest, global.ApiResponse{
			Code:    global.BAD_REQUEST,
			Message: "Cell ID required",
		})
		return
	}

	var cellID uint
	if _, err := fmt.Sscanf(cellIDStr, "%d", &cellID); err != nil || cellID == 0 {
		ctx.JSON(http.StatusBadRequest, global.ApiResponse{
			Code:    global.BAD_REQUEST,
			Message: "Invalid cell ID format",
		})
		return
	}

	// 取得 cell 來確認存在
	cell, err := ctl.cellRepository.GetByID(ctl.makeCtx(ctx), cellID)
	if err != nil {
		ctx.JSON(http.StatusNotFound, global.ApiResponse{
			Code:    http.StatusNotFound,
			Message: "Cell not found",
		})
		return
	}

	// 刪除格子
	if err := ctl.cellRepository.Delete(ctl.makeCtx(ctx), cellID); err != nil {
		ctx.JSON(http.StatusInternalServerError, global.ApiResponse{
			Code:    500,
			Message: "Failed to delete cell: " + err.Error(),
		})
		return
	}

	actorID := ctx.GetUint(global.UserIDKey)
	ctl.auditLogRepo.Create(ctx, models.AuditLog{
		CenterID:   centerID,
		ActorType:  "ADMIN",
		ActorID:    actorID,
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

	ctx.JSON(http.StatusOK, global.ApiResponse{
		Code:    0,
		Message: "Cell deleted",
	})
}

func (ctl *TimetableTemplateController) DeleteTemplate(ctx *gin.Context) {
	centerID, ok := ctl.getCenterID(ctx)
	if !ok {
		return
	}

	templateIDStr := ctx.Param("templateId")

	if templateIDStr == "" {
		ctx.JSON(http.StatusBadRequest, global.ApiResponse{
			Code:    global.BAD_REQUEST,
			Message: "Template ID required",
		})
		return
	}

	var templateID uint
	if _, err := fmt.Sscanf(templateIDStr, "%d", &templateID); err != nil || templateID == 0 {
		ctx.JSON(http.StatusBadRequest, global.ApiResponse{
			Code:    global.BAD_REQUEST,
			Message: "Invalid template ID format",
		})
		return
	}

	if err := ctl.templateRepository.DeleteByIDAndCenterID(ctl.makeCtx(ctx), templateID, centerID); err != nil {
		ctx.JSON(http.StatusInternalServerError, global.ApiResponse{
			Code:    500,
			Message: "Failed to delete template",
		})
		return
	}

	actorID := ctx.GetUint(global.UserIDKey)
	ctl.auditLogRepo.Create(ctx, models.AuditLog{
		CenterID:   centerID,
		ActorType:  "ADMIN",
		ActorID:    actorID,
		Action:     "TEMPLATE_DELETE",
		TargetType: "TimetableTemplate",
		TargetID:   templateID,
		Payload: models.AuditPayload{
			After: map[string]interface{}{
				"status": "DELETED",
			},
		},
	})

	ctx.JSON(http.StatusOK, global.ApiResponse{
		Code:    0,
		Message: "Template deleted",
	})
}

type ApplyTemplateRequest struct {
	OfferingID     uint     `json:"offering_id" binding:"required"`
	StartDate      string   `json:"start_date" binding:"required"`
	EndDate        string   `json:"end_date" binding:"required"`
	Weekdays       []int    `json:"weekdays" binding:"required"`
	Duration       int      `json:"duration"`
	OverrideBuffer bool     `json:"override_buffer"` // 允許覆蓋 Buffer 衝突
}

// ApplyTemplateConflictInfo 套用模板衝突資訊（向前相容）
type ApplyTemplateConflictInfo struct {
	Weekday      int    `json:"weekday"`
	StartTime    string `json:"start_time"`
	EndTime      string `json:"end_time"`
	ConflictType string `json:"conflict_type"` // "ROOM_OVERLAP", "TEACHER_OVERLAP", "PERSONAL_EVENT", "TEACHER_BUFFER", "ROOM_BUFFER"
	Message      string `json:"message"`
	RuleID       uint   `json:"rule_id,omitempty"`
	CanOverride  bool   `json:"can_override,omitempty"`
}

func (ctl *TimetableTemplateController) ApplyTemplate(ctx *gin.Context) {
	centerID, ok := ctl.getCenterID(ctx)
	if !ok {
		return
	}

	templateIDStr := ctx.Param("templateId")

	if templateIDStr == "" {
		ctx.JSON(http.StatusBadRequest, global.ApiResponse{
			Code:    global.BAD_REQUEST,
			Message: "Template ID required",
		})
		return
	}

	var templateID uint
	if _, err := fmt.Sscanf(templateIDStr, "%d", &templateID); err != nil || templateID == 0 {
		ctx.JSON(http.StatusBadRequest, global.ApiResponse{
			Code:    global.BAD_REQUEST,
			Message: "Invalid template ID format",
		})
		return
	}

	// 解析請求
	var req ApplyTemplateRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, global.ApiResponse{
			Code:    global.BAD_REQUEST,
			Message: "Invalid request body: " + err.Error(),
		})
		return
	}

	// 解析日期
	startDate, err := time.Parse("2006-01-02", req.StartDate)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, global.ApiResponse{
			Code:    global.BAD_REQUEST,
			Message: "Invalid start_date format",
		})
		return
	}
	endDate, err := time.Parse("2006-01-02", req.EndDate)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, global.ApiResponse{
			Code:    global.BAD_REQUEST,
			Message: "Invalid end_date format",
		})
		return
	}

	// 取得模板
	template, err := ctl.templateRepository.GetByID(ctl.makeCtx(ctx), templateID)
	if err != nil || template.CenterID != centerID {
		ctx.JSON(http.StatusNotFound, global.ApiResponse{
			Code:    http.StatusNotFound,
			Message: "Template not found",
		})
		return
	}

	// 取得模板中的 cells
	cells, err := ctl.cellRepository.ListByTemplateID(ctl.makeCtx(ctx), templateID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, global.ApiResponse{
			Code:    500,
			Message: "Failed to get template cells",
		})
		return
	}

	// 使用 ScheduleRuleValidator 進行完整驗證（Overlap + Buffer）
	// 如果請求中指定了 override_buffer，則允許覆蓋 Buffer 衝突
	validationSummary, err := ctl.ruleValidator.ValidateForApplyTemplate(
		ctl.makeCtx(ctx),
		centerID,
		req.OfferingID,
		req.Weekdays,
		cells,
		req.StartDate,
		req.EndDate,
		req.OverrideBuffer,
	)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, global.ApiResponse{
			Code:    500,
			Message: "Failed to validate template: " + err.Error(),
		})
		return
	}

	// 如果有不可覆蓋的衝突，回傳衝突資訊
	if !validationSummary.Valid {
		// 計算不可覆蓋的衝突數量
		nonOverrideConflicts := 0
		for _, conflict := range validationSummary.AllConflicts {
			if !conflict.CanOverride {
				nonOverrideConflicts++
			}
		}

		// 如果還有不可覆蓋的衝突，才回傳錯誤
		if nonOverrideConflicts > 0 {
			// 轉換衝突格式以維持向前相容
			var allConflicts []ApplyTemplateConflictInfo
			for _, conflict := range validationSummary.AllConflicts {
				allConflicts = append(allConflicts, ApplyTemplateConflictInfo{
					Weekday:      conflict.Weekday,
					StartTime:    conflict.StartTime,
					EndTime:      conflict.EndTime,
					ConflictType: conflict.ConflictType,
					Message:      conflict.Message,
					RuleID:       conflict.RuleID,
					CanOverride:  conflict.CanOverride,
				})
			}

			ctx.JSON(http.StatusBadRequest, global.ApiResponse{
				Code:    40002, // OVERLAP error code
				Message: "套用模板會產生時間衝突，請先解決衝突後再嘗試",
				Datas: map[string]interface{}{
					"conflicts":      allConflicts,
					"conflict_count": len(allConflicts),
				},
			})
			return
		}

		// 如果只有可覆蓋的衝突，但管理員選擇覆蓋，則繼續執行
		// 否則回傳警告資訊
		if !req.OverrideBuffer {
			var allConflicts []ApplyTemplateConflictInfo
			for _, conflict := range validationSummary.AllConflicts {
				allConflicts = append(allConflicts, ApplyTemplateConflictInfo{
					Weekday:      conflict.Weekday,
					StartTime:    conflict.StartTime,
					EndTime:      conflict.EndTime,
					ConflictType: conflict.ConflictType,
					Message:      conflict.Message,
					RuleID:       conflict.RuleID,
					CanOverride:  conflict.CanOverride,
				})
			}

			ctx.JSON(http.StatusBadRequest, global.ApiResponse{
				Code:    40003, // BUFFER_CONFLICT warning code
				Message: "套用模板會產生緩衝時間衝突，是否繼續？",
				Datas: map[string]interface{}{
					"conflicts":      allConflicts,
					"conflict_count": len(allConflicts),
					"can_override":   true,
				},
			})
			return
		}
	}

	// 為每個 weekday 和每個 cell 建立 schedule rule
	var rules []models.ScheduleRule
	for _, weekday := range req.Weekdays {
		for _, cell := range cells {
			rule := models.ScheduleRule{
				CenterID:    centerID,
				OfferingID:  req.OfferingID,
				TeacherID:   cell.TeacherID,
				RoomID:      *cell.RoomID,
				Weekday:     weekday,
				StartTime:   cell.StartTime,
				EndTime:     cell.EndTime,
				EffectiveRange: models.DateRange{
					StartDate: startDate,
					EndDate:   endDate,
				},
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			}
			rules = append(rules, rule)
		}
	}

	// 建立 schedule rules
	createdRules, err := ctl.scheduleRuleRepo.BulkCreate(ctl.makeCtx(ctx), rules)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, global.ApiResponse{
			Code:    500,
			Message: "Failed to create schedule rules: " + err.Error(),
		})
		return
	}

	actorID := ctx.GetUint(global.UserIDKey)
	ctl.auditLogRepo.Create(ctx, models.AuditLog{
		CenterID:   centerID,
		ActorType:  "ADMIN",
		ActorID:    actorID,
		Action:     "TEMPLATE_APPLY",
		TargetType: "ScheduleRule",
		TargetID:   0,
		Payload: models.AuditPayload{
			After: map[string]interface{}{
				"template_id":   templateID,
				"offering_id":   req.OfferingID,
				"start_date":    req.StartDate,
				"end_date":      req.EndDate,
				"weekdays":      req.Weekdays,
				"rules_created": len(createdRules),
			},
		},
	})

	ctx.JSON(http.StatusOK, global.ApiResponse{
		Code:    0,
		Message: "Template applied successfully",
		Datas: map[string]interface{}{
			"rules_created": len(createdRules),
			"template_name": template.Name,
		},
	})
}
