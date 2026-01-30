package controllers

import (
	"timeLedger/app"
	"timeLedger/app/requests"
	"timeLedger/app/services"

	"github.com/gin-gonic/gin"
)

// TeacherEventController 老師個人行程相關 API
type TeacherEventController struct {
	BaseController
	app              *app.App
	personalEventSvc *services.PersonalEventService
}

func NewTeacherEventController(app *app.App) *TeacherEventController {
	return &TeacherEventController{
		app:              app,
		personalEventSvc: services.NewPersonalEventService(app),
	}
}

// requireTeacherID 取得並驗證老師 ID（通用模式）
func (ctl *TeacherEventController) requireTeacherID(helper *ContextHelper) uint {
	teacherID := helper.MustUserID()
	if teacherID == 0 {
		return 0
	}
	return teacherID
}

// requireEventID 取得並驗證行程 ID（通用模式）
func (ctl *TeacherEventController) requireEventID(helper *ContextHelper) uint {
	eventID := helper.MustParamUint("id")
	if eventID == 0 {
		return 0
	}
	return eventID
}

// requireTeacherAndEventID 取得老師 ID 和行程 ID（通用模式）
func (ctl *TeacherEventController) requireTeacherAndEventID(helper *ContextHelper) (uint, uint) {
	teacherID := ctl.requireTeacherID(helper)
	if teacherID == 0 {
		return 0, 0
	}
	eventID := ctl.requireEventID(helper)
	if eventID == 0 {
		return teacherID, 0
	}
	return teacherID, eventID
}

// GetPersonalEvents 取得老師個人行程列表
// @Summary 取得老師個人行程列表
// @Tags Teacher - Events
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param from query string false "開始日期 (YYYY-MM-DD)"
// @Param to query string false "結束日期 (YYYY-MM-DD)"
// @Success 200 {object} global.ApiResponse{data=[]models.PersonalEvent}
// @Router /api/v1/teacher/me/personal-events [get]
func (ctl *TeacherEventController) GetPersonalEvents(ctx *gin.Context) {
	helper := NewContextHelper(ctx)
	teacherID := ctl.requireTeacherID(helper)
	if teacherID == 0 {
		return
	}

	from, to, _ := helper.QueryDateRange("from", "to")
	events, errInfo, err := ctl.personalEventSvc.GetPersonalEvents(ctx.Request.Context(), teacherID, from, to)
	if err != nil {
		helper.ErrorWithInfo(errInfo)
		return
	}

	helper.Success(events)
}

// CreatePersonalEvent 新增老師個人行程
// @Summary 新增老師個人行程
// @Tags Teacher - Events
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body requests.CreatePersonalEventRequest true "行程資訊"
// @Success 200 {object} global.ApiResponse{data=models.PersonalEvent}
// @Failure 400 {object} global.ApiResponse
// @Failure 409 {object} global.ApiResponse
// @Router /api/v1/teacher/me/personal-events [post]
func (ctl *TeacherEventController) CreatePersonalEvent(ctx *gin.Context) {
	helper := NewContextHelper(ctx)
	teacherID := ctl.requireTeacherID(helper)
	if teacherID == 0 {
		return
	}

	var req requests.CreatePersonalEventRequest
	if !helper.MustBindJSON(&req) {
		return
	}

	svcReq := &services.CreatePersonalEventRequest{
		Title:          req.Title,
		StartAt:        req.StartAt,
		EndAt:          req.EndAt,
		IsAllDay:       req.IsAllDay,
		ColorHex:       req.ColorHex,
		RecurrenceRule: req.RecurrenceRule,
	}

	event, errInfo, err := ctl.personalEventSvc.CreatePersonalEvent(ctx.Request.Context(), teacherID, svcReq)
	if err != nil {
		helper.ErrorWithInfo(errInfo)
		return
	}

	helper.Success(event)
}

// UpdatePersonalEvent 更新老師個人行程
// @Summary 更新老師個人行程
// @Tags Teacher - Events
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path uint true "行程ID"
// @Param request body requests.UpdatePersonalEventRequest true "行程資訊"
// @Success 200 {object} global.ApiResponse{data=requests.UpdatePersonalEventResponse}
// @Router /api/v1/teacher/me/personal-events/{id} [patch]
func (ctl *TeacherEventController) UpdatePersonalEvent(ctx *gin.Context) {
	helper := NewContextHelper(ctx)
	teacherID, eventID := ctl.requireTeacherAndEventID(helper)
	if teacherID == 0 || eventID == 0 {
		return
	}

	var req requests.UpdatePersonalEventRequest
	if !helper.MustBindJSON(&req) {
		return
	}

	svcReq := &services.UpdatePersonalEventRequest{
		Title:          req.Title,
		StartAt:        req.StartAt,
		EndAt:          req.EndAt,
		IsAllDay:       req.IsAllDay,
		ColorHex:       req.ColorHex,
		RecurrenceRule: req.RecurrenceRule,
		UpdateMode:     req.UpdateMode,
	}

	result, errInfo, err := ctl.personalEventSvc.UpdatePersonalEvent(ctx.Request.Context(), eventID, teacherID, svcReq)
	if err != nil {
		helper.ErrorWithInfo(errInfo)
		return
	}

	helper.Success(requests.UpdatePersonalEventResponse{
		UpdatedCount: result.UpdatedCount,
		Message:      result.Message,
	})
}

// DeletePersonalEvent 刪除老師個人行程
// @Summary 刪除老師個人行程
// @Tags Teacher - Events
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path uint true "行程ID"
// @Success 200 {object} global.ApiResponse
// @Router /api/v1/teacher/me/personal-events/{id} [delete]
func (ctl *TeacherEventController) DeletePersonalEvent(ctx *gin.Context) {
	helper := NewContextHelper(ctx)
	_, eventID := ctl.requireTeacherAndEventID(helper)
	if eventID == 0 {
		return
	}

	errInfo := ctl.personalEventSvc.DeletePersonalEvent(ctx.Request.Context(), eventID, ctl.requireTeacherID(helper))
	if errInfo != nil {
		helper.ErrorWithInfo(errInfo)
		return
	}

	helper.Success(gin.H{"message": "Personal event deleted"})
}

// GetPersonalEventNote 取得個人行程備註
// @Summary 取得個人行程備註
// @Tags Teacher - Events
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path uint true "行程ID"
// @Success 200 {object} global.ApiResponse{data=requests.PersonalEventNoteResponse}
// @Router /api/v1/teacher/me/personal-events/{id}/note [get]
func (ctl *TeacherEventController) GetPersonalEventNote(ctx *gin.Context) {
	helper := NewContextHelper(ctx)
	teacherID, eventID := ctl.requireTeacherAndEventID(helper)
	if teacherID == 0 || eventID == 0 {
		return
	}

	note, errInfo, err := ctl.personalEventSvc.GetPersonalEventNote(ctx.Request.Context(), eventID, teacherID)
	if err != nil {
		helper.ErrorWithInfo(errInfo)
		return
	}

	helper.Success(requests.PersonalEventNoteResponse{Content: note})
}

// UpdatePersonalEventNote 更新個人行程備註
// @Summary 更新個人行程備註
// @Tags Teacher - Events
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path uint true "行程ID"
// @Param request body requests.UpdatePersonalEventNoteRequest true "備註內容"
// @Success 200 {object} global.ApiResponse
// @Router /api/v1/teacher/me/personal-events/{id}/note [put]
func (ctl *TeacherEventController) UpdatePersonalEventNote(ctx *gin.Context) {
	helper := NewContextHelper(ctx)
	teacherID, eventID := ctl.requireTeacherAndEventID(helper)
	if teacherID == 0 || eventID == 0 {
		return
	}

	var req requests.UpdatePersonalEventNoteRequest
	if !helper.MustBindJSON(&req) {
		return
	}

	svcReq := &services.UpdatePersonalEventNoteRequest{Content: req.Content}
	errInfo := ctl.personalEventSvc.UpdatePersonalEventNote(ctx.Request.Context(), eventID, teacherID, svcReq)
	if errInfo != nil {
		helper.ErrorWithInfo(errInfo)
		return
	}

	helper.Success(gin.H{"message": "Note updated successfully"})
}
