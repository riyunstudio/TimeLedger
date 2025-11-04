package controllers

import (
	"akali/global"
	"akali/global/errInfos"

	"github.com/gin-gonic/gin"
)

type BaseController struct{}

// JSON 輔助回傳
func (b *BaseController) JSON(ctx *gin.Context, response global.Ret) {
	// 檢查 errInfo, 如果是nil就初始化
	if response.ErrInfo == nil {
		response.ErrInfo = &errInfos.Err{}
	}

	// 存到 gin.Context
	ctx.Set(global.RetKey, response)

	// 客戶端回傳
	if response.ErrInfo.Code == 0 {
		ctx.JSON(response.Status, global.ApiResponse{
			Code:    response.ErrInfo.Code,
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
