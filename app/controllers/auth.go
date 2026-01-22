package controllers

import (
	"net/http"
	"timeLedger/app"
	"timeLedger/app/services"
	"timeLedger/global"

	"github.com/gin-gonic/gin"
)

type AuthController struct {
	BaseController
	authService services.AuthService
}

func NewAuthController(app *app.App, authService services.AuthService) *AuthController {
	return &AuthController{
		authService: authService,
	}
}

// AdminLogin 管理員登入
// @Summary 管理員登入
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body AdminLoginRequest true "登入資訊"
// @Success 200 {object} global.ApiResponse{data=services.LoginResponse}
// @Router /api/v1/auth/admin/login [post]
func (ctl *AuthController) AdminLogin(ctx *gin.Context) {
	var req AdminLoginRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, global.ApiResponse{
			Code:    global.BAD_REQUEST,
			Message: "Invalid request body",
		})
		return
	}

	response, err := ctl.authService.AdminLogin(ctx, req.Email, req.Password)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, global.ApiResponse{
			Code:    global.UNAUTHORIZED,
			Message: err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, global.ApiResponse{
		Code:    0,
		Message: "Login successful",
		Datas:   response,
	})
}

type AdminLoginRequest struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// TeacherLineLogin 老師 LINE 登入
// @Summary 老師 LINE 登入
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body TeacherLineLoginRequest true "登入資訊"
// @Success 200 {object} global.ApiResponse{data=services.LoginResponse}
// @Router /api/v1/auth/teacher/line/login [post]
func (ctl *AuthController) TeacherLineLogin(ctx *gin.Context) {
	var req TeacherLineLoginRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, global.ApiResponse{
			Code:    global.BAD_REQUEST,
			Message: "Invalid request body",
		})
		return
	}

	response, err := ctl.authService.TeacherLineLogin(ctx, req.LineUserID, req.AccessToken)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, global.ApiResponse{
			Code:    global.UNAUTHORIZED,
			Message: err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, global.ApiResponse{
		Code:    0,
		Message: "Login successful",
		Datas:   response,
	})
}

type TeacherLineLoginRequest struct {
	LineUserID  string `json:"line_user_id" binding:"required"`
	AccessToken string `json:"access_token" binding:"required"`
}

// RefreshToken 刷新 Token
// @Summary 刷新 Token
// @Tags Auth
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} global.ApiResponse{data=services.LoginResponse}
// @Router /api/v1/auth/refresh [post]
func (ctl *AuthController) RefreshToken(ctx *gin.Context) {
	token := ctx.GetHeader("Authorization")
	if token == "" {
		ctx.JSON(http.StatusUnauthorized, global.ApiResponse{
			Code:    global.UNAUTHORIZED,
			Message: "Authorization header required",
		})
		return
	}

	response, err := ctl.authService.RefreshToken(ctx, token)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, global.ApiResponse{
			Code:    global.UNAUTHORIZED,
			Message: err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, global.ApiResponse{
		Code:    0,
		Message: "Token refreshed",
		Datas:   response,
	})
}

// Logout 登出
// @Summary 登出
// @Tags Auth
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} global.ApiResponse
// @Router /api/v1/auth/logout [post]
func (ctl *AuthController) Logout(ctx *gin.Context) {
	token := ctx.GetHeader("Authorization")
	if token == "" {
		ctx.JSON(http.StatusUnauthorized, global.ApiResponse{
			Code:    global.UNAUTHORIZED,
			Message: "Authorization header required",
		})
		return
	}

	err := ctl.authService.Logout(ctx, token)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, global.ApiResponse{
			Code:    500,
			Message: err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, global.ApiResponse{
		Code:    0,
		Message: "Logout successful",
	})
}
