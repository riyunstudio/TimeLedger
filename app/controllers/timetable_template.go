package controllers

import (
	"fmt"
	"net/http"
	"time"
	"timeLedger/app"
	"timeLedger/app/models"
	"timeLedger/app/repositories"
	"timeLedger/app/requests"
	"timeLedger/global"

	"github.com/gin-gonic/gin"
)

type TimetableTemplateController struct {
	BaseController
	templateRepository *repositories.TimetableTemplateRepository
	cellRepository     *repositories.TimetableCellRepository
	auditLogRepo       *repositories.AuditLogRepository
}

func NewTimetableTemplateController(app *app.App) *TimetableTemplateController {
	return &TimetableTemplateController{
		templateRepository: repositories.NewTimetableTemplateRepository(app),
		cellRepository:     repositories.NewTimetableCellRepository(app),
		auditLogRepo:       repositories.NewAuditLogRepository(app),
	}
}

func (ctl *TimetableTemplateController) GetTemplates(ctx *gin.Context) {
	templates, err := ctl.templateRepository.ListByCenterID(ctl.makeCtx(ctx), 0)
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
	centerIDStr := ctx.Param("id")
	var centerID uint
	if _, err := fmt.Sscanf(centerIDStr, "%d", &centerID); err != nil || centerID == 0 {
		ctx.JSON(http.StatusBadRequest, global.ApiResponse{
			Code:    global.BAD_REQUEST,
			Message: "Invalid center ID",
		})
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
	centerID, _ := ctx.Get("center_id")

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

	centerIDUint, _ := centerID.(uint)

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
		CenterID:  centerIDUint,
		Name:      req.Name,
		UpdatedAt: time.Now(),
	}

	if err := ctl.templateRepository.UpdateByIDAndCenterID(ctl.makeCtx(ctx), templateID, centerIDUint, template); err != nil {
		ctx.JSON(http.StatusInternalServerError, global.ApiResponse{
			Code:    500,
			Message: "Failed to update template",
		})
		return
	}

	actorID := ctx.GetUint(global.UserIDKey)
	ctl.auditLogRepo.Create(ctx, models.AuditLog{
		CenterID:   centerIDUint,
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
	centerID, _ := ctx.Get("center_id")
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

	centerIDUint, _ := centerID.(uint)

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

	actorID := ctx.GetUint(global.UserIDKey)
	ctl.auditLogRepo.Create(ctx, models.AuditLog{
		CenterID:   centerIDUint,
		ActorType:  "ADMIN",
		ActorID:    actorID,
		Action:     "TIMETABLE_CELLS_CREATE",
		TargetType: "TimetableCell",
		TargetID:   0,
		Payload: models.AuditPayload{
			After: map[string]interface{}{
				"template_id": templateID,
				"cells_count": len(cells),
			},
		},
	})

	ctx.JSON(http.StatusOK, global.ApiResponse{
		Code:    0,
		Message: "Cells created",
		Datas:   cells,
	})
}

func (ctl *TimetableTemplateController) DeleteTemplate(ctx *gin.Context) {
	centerID, _ := ctx.Get("center_id")
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

	centerIDUint, _ := centerID.(uint)

	if err := ctl.templateRepository.DeleteByIDAndCenterID(ctl.makeCtx(ctx), templateID, centerIDUint); err != nil {
		ctx.JSON(http.StatusInternalServerError, global.ApiResponse{
			Code:    500,
			Message: "Failed to delete template",
		})
		return
	}

	actorID := ctx.GetUint(global.UserIDKey)
	ctl.auditLogRepo.Create(ctx, models.AuditLog{
		CenterID:   centerIDUint,
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
