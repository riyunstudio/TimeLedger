package controllers

import (
	"akali/app"
	"akali/global"
	"akali/global/errInfos"
	"context"

	"github.com/gin-gonic/gin"
)

type BaseController struct {
	app *app.App
}

// 建構函式
func NewBaseController(app *app.App) *BaseController {
	return &BaseController{
		app: app,
	}
}

// 把 gin ctx 取指定內容寫入 context ctx
func (b *BaseController) makeCtx(ctx *gin.Context) (c context.Context) {
	if val, exists := ctx.Get(global.TraceIDKey); exists {
		if tid, ok := val.(string); ok {
			c = context.WithValue(context.Background(), global.TraceIDKey, tid)
			return
		}
	}
	c = context.WithValue(context.Background(), global.TraceIDKey, "")
	return
}

// JSON 輔助回傳
func (b *BaseController) JSON(ctx *gin.Context, response global.Ret) {
	// 檢查 errInfo, 如果是nil就初始化
	if response.ErrInfo == nil {
		response.ErrInfo = &errInfos.Res{}
	}

	// 存到 gin.Context
	ctx.Set(global.RetKey, response)

	// 客戶端回傳
	if response.ErrInfo.Code == 0 {
		ctx.JSON(response.Status, global.ApiResponse{
			Code:    0,
			Message: "OK",
			Datas:   response.Datas,
		})
	} else {
		ctx.JSON(response.Status, global.ApiResponse{
			Code:    response.ErrInfo.Code,
			Message: response.ErrInfo.Msg,
			Datas:   response.Datas,
		})
	}
}
