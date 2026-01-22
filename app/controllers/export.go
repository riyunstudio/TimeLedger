package controllers

import (
	"fmt"
	"net/http"
	"strings"
	"time"
	"timeLedger/app"
	"timeLedger/app/requests"
	"timeLedger/app/services"
	"timeLedger/global"

	"github.com/gin-gonic/gin"
)

type ExportController struct {
	BaseController
	app       *app.App
	exportSvc services.ExportService
}

func NewExportController(app *app.App) *ExportController {
	return &ExportController{
		app:       app,
		exportSvc: services.NewExportService(app),
	}
}

func (ctl *ExportController) ExportScheduleCSV(ctx *gin.Context) {
	var req requests.ExportScheduleRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, global.ApiResponse{
			Code:    global.BAD_REQUEST,
			Message: "Invalid request body: " + err.Error(),
		})
		return
	}

	if req.CenterID == 0 {
		ctx.JSON(http.StatusBadRequest, global.ApiResponse{
			Code:    global.BAD_REQUEST,
			Message: "Center ID required",
		})
		return
	}

	startDateStr := req.StartDate
	endDateStr := req.EndDate

	if startDateStr == "" || endDateStr == "" {
		ctx.JSON(http.StatusBadRequest, global.ApiResponse{
			Code:    global.BAD_REQUEST,
			Message: "Start date and end date required",
		})
		return
	}

	startDateStr = strings.TrimSuffix(startDateStr, "Z")
	endDateStr = strings.TrimSuffix(endDateStr, "Z")

	var startDate, endDate time.Time
	var err error

	startDate, err = time.Parse("2006-01-02T15:04:05", startDateStr)
	if err != nil {
		startDate, err = time.Parse("2006-01-02", startDateStr)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, global.ApiResponse{
				Code:    global.BAD_REQUEST,
				Message: "Invalid start date format: " + err.Error(),
			})
			return
		}
	}

	endDate, err = time.Parse("2006-01-02T15:04:05", endDateStr)
	if err != nil {
		endDate, err = time.Parse("2006-01-02", endDateStr)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, global.ApiResponse{
				Code:    global.BAD_REQUEST,
				Message: "Invalid end date format: " + err.Error(),
			})
			return
		}
	}

	data, err := ctl.exportSvc.ExportScheduleToCSV(ctx, req.CenterID, startDate, endDate)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, global.ApiResponse{
			Code:    500,
			Message: err.Error(),
		})
		return
	}

	filename := fmt.Sprintf("schedule_%s_%s.csv", startDate.Format("20060102"), endDate.Format("20060102"))
	ctx.Header("Content-Type", "text/csv")
	ctx.Header("Content-Disposition", "attachment; filename="+filename)
	ctx.Data(http.StatusOK, "text/csv", data)
}

func (ctl *ExportController) ExportSchedulePDF(ctx *gin.Context) {
	var req requests.ExportScheduleRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, global.ApiResponse{
			Code:    global.BAD_REQUEST,
			Message: "Invalid request body: " + err.Error(),
		})
		return
	}

	if req.CenterID == 0 {
		ctx.JSON(http.StatusBadRequest, global.ApiResponse{
			Code:    global.BAD_REQUEST,
			Message: "Center ID required",
		})
		return
	}

	startDateStr := req.StartDate
	endDateStr := req.EndDate

	if startDateStr == "" || endDateStr == "" {
		ctx.JSON(http.StatusBadRequest, global.ApiResponse{
			Code:    global.BAD_REQUEST,
			Message: "Start date and end date required",
		})
		return
	}

	startDateStr = strings.TrimSuffix(startDateStr, "Z")
	endDateStr = strings.TrimSuffix(endDateStr, "Z")

	startDate, err := time.Parse("2006-01-02T15:04:05", startDateStr)
	if err != nil {
		startDate, err = time.Parse("2006-01-02", startDateStr)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, global.ApiResponse{
				Code:    global.BAD_REQUEST,
				Message: "Invalid start date format: " + err.Error(),
			})
			return
		}
	}

	endDate, err := time.Parse("2006-01-02T15:04:05", endDateStr)
	if err != nil {
		endDate, err = time.Parse("2006-01-02", endDateStr)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, global.ApiResponse{
				Code:    global.BAD_REQUEST,
				Message: "Invalid end date format: " + err.Error(),
			})
			return
		}
	}

	data, err := ctl.exportSvc.ExportScheduleToPDF(ctx, req.CenterID, startDate, endDate)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, global.ApiResponse{
			Code:    500,
			Message: err.Error(),
		})
		return
	}

	filename := fmt.Sprintf("schedule_%s_%s.txt", startDate.Format("20060102"), endDate.Format("20060102"))
	ctx.Header("Content-Type", "text/plain")
	ctx.Header("Content-Disposition", "attachment; filename="+filename)
	ctx.Data(http.StatusOK, "text/plain", data)
}

func (ctl *ExportController) ExportTeachersCSV(ctx *gin.Context) {
	centerID := ctx.Param("id")
	var centerIDInt uint
	if _, err := fmt.Sscanf(centerID, "%d", &centerIDInt); err != nil {
		ctx.JSON(http.StatusBadRequest, global.ApiResponse{
			Code:    global.BAD_REQUEST,
			Message: "Invalid center ID",
		})
		return
	}

	data, err := ctl.exportSvc.ExportTeachersToCSV(ctx, centerIDInt)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, global.ApiResponse{
			Code:    500,
			Message: err.Error(),
		})
		return
	}

	filename := fmt.Sprintf("teachers_%d.csv", centerIDInt)
	ctx.Header("Content-Type", "text/csv")
	ctx.Header("Content-Disposition", "attachment; filename="+filename)
	ctx.Data(http.StatusOK, "text/csv", data)
}

func (ctl *ExportController) ExportExceptionsCSV(ctx *gin.Context) {
	centerID := ctx.Param("id")
	var centerIDInt uint
	if _, err := fmt.Sscanf(centerID, "%d", &centerIDInt); err != nil {
		ctx.JSON(http.StatusBadRequest, global.ApiResponse{
			Code:    global.BAD_REQUEST,
			Message: "Invalid center ID",
		})
		return
	}

	startDateStr := ctx.Query("start_date")
	endDateStr := ctx.Query("end_date")

	startDate, err := time.Parse("2006-01-02", startDateStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, global.ApiResponse{
			Code:    global.BAD_REQUEST,
			Message: "Invalid start date format",
		})
		return
	}

	endDate, err := time.Parse("2006-01-02", endDateStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, global.ApiResponse{
			Code:    global.BAD_REQUEST,
			Message: "Invalid end date format",
		})
		return
	}

	data, err := ctl.exportSvc.ExportExceptionsToCSV(ctx, centerIDInt, startDate, endDate)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, global.ApiResponse{
			Code:    500,
			Message: err.Error(),
		})
		return
	}

	filename := fmt.Sprintf("exceptions_%s_%s.csv", startDate.Format("20060102"), endDate.Format("20060102"))
	ctx.Header("Content-Type", "text/csv")
	ctx.Header("Content-Disposition", "attachment; filename="+filename)
	ctx.Data(http.StatusOK, "text/csv", data)
}
