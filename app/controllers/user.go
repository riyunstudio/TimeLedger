package controllers

import (
	"akali/app"
	"akali/app/requests"
	"akali/app/services"
	"akali/global"
	"net/http"

	"github.com/gin-gonic/gin"
)

type UserController struct {
	BaseController
	UserRequest *requests.UserRequest
	UserService *services.UserService
}

// 建構函式
func NewUserController(app *app.App) *UserController {
	return &UserController{
		UserRequest: requests.NewUserRequest(app),
		UserService: services.NewUserService(app),
	}
}

// @Summary 查詢使用者
// @description
// @Tags User
// @Param Content-Type header string true "Content-Type" default(application/json)
// @Param Tid header string false "TraceID"
// @Param ID query int true "會員ID"
// @Success 200 {object} apiResponse{datas=resources.UserGetResource} "回傳"
// @Router /user [get]
func (ctl *UserController) Get(ctx *gin.Context) {
	req, eInfo, err := ctl.UserRequest.Get(ctx)
	if err != nil {
		ctl.JSON(ctx, global.Ret{Status: http.StatusBadRequest, ErrInfo: eInfo, Err: err})
		return
	}

	datas, eInfo, err := ctl.UserService.Get(req)
	if err != nil {
		ctl.JSON(ctx, global.Ret{Status: http.StatusInternalServerError, ErrInfo: eInfo, Err: err})
		return
	}
	ctl.JSON(ctx, global.Ret{Status: http.StatusOK, Datas: datas})
}

// @Summary 新增使用者
// @description
// @Tags User
// @Param Content-Type header string true "Content-Type" default(application/json)
// @Param Tid header string false "TraceID"
// @param Params body requests.UserCreateRequest true "參數"
// @Success 200 {object} apiResponse{datas=resources.UserCreateResource} "回傳"
// @Router /user [post]
func (ctl *UserController) Create(ctx *gin.Context) {
	req, eInfo, err := ctl.UserRequest.Create(ctx)
	if err != nil {
		ctl.JSON(ctx, global.Ret{Status: http.StatusBadRequest, ErrInfo: eInfo, Err: err})
		return
	}

	datas, eInfo, err := ctl.UserService.Create(req)
	if err != nil {
		ctl.JSON(ctx, global.Ret{Status: http.StatusInternalServerError, ErrInfo: eInfo, Err: err})
		return
	}
	ctl.JSON(ctx, global.Ret{Status: http.StatusOK, Datas: datas})
}

// @Summary 修改使用者
// @description
// @Tags User
// @Param Content-Type header string true "Content-Type" default(application/json)
// @Param Tid header string false "TraceID"
// @param Params body requests.UserUpdateRequest true "參數"
// @Success 200 {object} apiResponse{datas=resources.UserUpdateResource} "回傳"
// @Router /user [put]
func (ctl *UserController) Update(ctx *gin.Context) {
	req, eInfo, err := ctl.UserRequest.Update(ctx)
	if err != nil {
		ctl.JSON(ctx, global.Ret{Status: http.StatusBadRequest, ErrInfo: eInfo, Err: err})
		return
	}

	datas, eInfo, err := ctl.UserService.Update(req)
	if err != nil {
		ctl.JSON(ctx, global.Ret{Status: http.StatusInternalServerError, ErrInfo: eInfo, Err: err})
		return
	}
	ctl.JSON(ctx, global.Ret{Status: http.StatusOK, Datas: datas})
}
