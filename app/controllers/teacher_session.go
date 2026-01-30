package controllers

import (
	"time"
	"timeLedger/app"
	"timeLedger/app/repositories"
	"timeLedger/app/requests"
	"timeLedger/app/resources"

	"github.com/gin-gonic/gin"
)

// TeacherSessionController 老師課堂筆記相關 API
type TeacherSessionController struct {
	BaseController
	app         *app.App
	sessionNote *repositories.SessionNoteRepository
}

func NewTeacherSessionController(app *app.App) *TeacherSessionController {
	return &TeacherSessionController{
		app:         app,
		sessionNote: repositories.NewSessionNoteRepository(app),
	}
}

// GetSessionNote 取得課堂筆記
// @Summary 取得課堂筆記
// @Tags Teacher - Sessions
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param rule_id query uint true "課程規則ID"
// @Param session_date query string true "課程日期 (YYYY-MM-DD)"
// @Success 200 {object} global.ApiResponse{data=resources.SessionNoteResource}
// @Router /api/v1/teacher/sessions/note [get]
func (ctl *TeacherSessionController) GetSessionNote(ctx *gin.Context) {
	helper := NewContextHelper(ctx)

	teacherID := helper.MustUserID()
	if teacherID == 0 {
		return
	}

	ruleID := helper.MustQueryUint("rule_id")
	if ruleID == 0 {
		return
	}

	sessionDateStr, _ := helper.QueryString("session_date")
	if sessionDateStr == "" {
		helper.BadRequest("session_date required")
		return
	}

	sessionDate, err := time.Parse("2006-01-02", sessionDateStr)
	if err != nil {
		helper.BadRequest("Invalid session_date format, use YYYY-MM-DD")
		return
	}

	note, isNew, dbErr := ctl.sessionNote.GetOrCreate(ctx, teacherID, ruleID, sessionDate)
	if dbErr != nil {
		helper.InternalError(dbErr.Error())
		return
	}

	response := resources.SessionNoteResource{
		ID:          note.ID,
		RuleID:      note.RuleID,
		SessionDate: note.SessionDate.Format("2006-01-02"),
		Content:     note.Content,
		PrepNote:    note.PrepNote,
		UpdatedAt:   note.UpdatedAt,
	}

	helper.Success(map[string]interface{}{
		"note":   response,
		"is_new": isNew,
	})
}

// UpsertSessionNote 新增或更新課堂筆記
// @Summary 新增或更新課堂筆記
// @Tags Teacher - Sessions
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body requests.UpsertSessionNoteRequest true "筆記內容"
// @Success 200 {object} global.ApiResponse{data=resources.SessionNoteResource}
// @Router /api/v1/teacher/sessions/note [put]
func (ctl *TeacherSessionController) UpsertSessionNote(ctx *gin.Context) {
	helper := NewContextHelper(ctx)

	teacherID := helper.MustUserID()
	if teacherID == 0 {
		return
	}

	var req requests.UpsertSessionNoteRequest
	if !helper.MustBindJSON(&req) {
		return
	}

	sessionDate, err := time.Parse("2006-01-02", req.SessionDate)
	if err != nil {
		helper.BadRequest("Invalid session_date format, use YYYY-MM-DD")
		return
	}

	note, _, err := ctl.sessionNote.GetOrCreate(ctx, teacherID, req.RuleID, sessionDate)
	if err != nil {
		helper.InternalError(err.Error())
		return
	}

	note.Content = req.Content
	note.PrepNote = req.PrepNote
	note.UpdatedAt = time.Now()

	if err := ctl.sessionNote.Update(ctx, note); err != nil {
		helper.InternalError(err.Error())
		return
	}

	response := resources.SessionNoteResource{
		ID:          note.ID,
		RuleID:      note.RuleID,
		SessionDate: note.SessionDate.Format("2006-01-02"),
		Content:     note.Content,
		PrepNote:    note.PrepNote,
		UpdatedAt:   note.UpdatedAt,
	}

	helper.Success(response)
}
