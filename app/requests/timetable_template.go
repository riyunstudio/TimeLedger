package requests

import "github.com/gin-gonic/gin"

type TemplateIDRequest struct {
	CenterID   uint `uri:"id" binding:"required"`
	TemplateID uint `uri:"template_id" binding:"required"`
}

type TemplateCellsRequest struct {
	CenterID   uint `uri:"id" binding:"required"`
	TemplateID uint `uri:"template_id" binding:"required"`
}

type CreateTemplateRequest struct {
	CenterID uint   `json:"-"`
	Name     string `json:"name" binding:"required"`
	RowType  string `json:"row_type" binding:"required,oneof=ROOM TEACHER"`
}

type UpdateTemplateRequest struct {
	Name string `json:"name"`
}

type CreateCellRequest struct {
	RowNo     int    `json:"row_no" binding:"required"`
	ColNo     int    `json:"col_no" binding:"required"`
	StartTime string `json:"start_time" binding:"required"`
	EndTime   string `json:"end_time" binding:"required"`
	RoomID    *uint  `json:"room_id"`
	TeacherID *uint  `json:"teacher_id"`
}

func ValidateTemplateID(ctx *gin.Context) (*TemplateIDRequest, error) {
	return ValidateURI[TemplateIDRequest](ctx)
}

func ValidateTemplateCells(ctx *gin.Context) (*TemplateCellsRequest, error) {
	return ValidateURI[TemplateCellsRequest](ctx)
}
