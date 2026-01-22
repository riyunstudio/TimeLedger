package requests

import (
	"timeLedger/app"
	"timeLedger/global/errInfos"

	"github.com/gin-gonic/gin"
)

type SessionNoteRequest struct {
	BaseRequest
	app *app.App
}

func NewSessionNoteRequest(app *app.App) *SessionNoteRequest {
	return &SessionNoteRequest{
		app: app,
	}
}

type GetSessionNoteRequest struct {
	RuleID      uint   `json:"rule_id" binding:"required"`
	SessionDate string `json:"session_date" binding:"required"`
}

func (r *SessionNoteRequest) Get(ctx *gin.Context) (req *GetSessionNoteRequest, errInfo *errInfos.Res, err error) {
	req, err = Validate[GetSessionNoteRequest](ctx)
	if err != nil {
		return nil, r.app.Err.New(errInfos.PARAMS_VALIDATE_ERROR), err
	}
	return
}

type UpsertSessionNoteRequest struct {
	RuleID      uint   `json:"rule_id" binding:"required"`
	SessionDate string `json:"session_date" binding:"required"`
	Content     string `json:"content"`
	PrepNote    string `json:"prep_note"`
}

func (r *SessionNoteRequest) Upsert(ctx *gin.Context) (req *UpsertSessionNoteRequest, errInfo *errInfos.Res, err error) {
	req, err = Validate[UpsertSessionNoteRequest](ctx)
	if err != nil {
		return nil, r.app.Err.New(errInfos.PARAMS_VALIDATE_ERROR), err
	}
	return
}
