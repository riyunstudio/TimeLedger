package requests

import "github.com/gin-gonic/gin"

type ExportScheduleRequest struct {
	CenterID  uint   `json:"center_id" binding:"required"`
	StartDate string `json:"start_date"`
	EndDate   string `json:"end_date"`
}

func ValidateExportSchedule(ctx *gin.Context) (*ExportScheduleRequest, error) {
	return Validate[ExportScheduleRequest](ctx)
}
