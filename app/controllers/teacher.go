package controllers

import (
	"timeLedger/app"
	"timeLedger/app/repositories"
	"timeLedger/app/services"

	"github.com/gin-gonic/gin"
)

// TeacherController 教師相關 API（保留核心功能）
type TeacherController struct {
	BaseController
	app                  *app.App
	teacherRepository    *repositories.TeacherRepository
	centerRepo           *repositories.CenterRepository
	scheduleRuleRepo     *repositories.ScheduleRuleRepository
	exceptionRepo        *repositories.ScheduleExceptionRepository
	expansionService     services.ScheduleExpansionService
	scheduleQueryService services.ScheduleQueryService
	sessionNoteRepo      *repositories.SessionNoteRepository
	adminUserRepo        *repositories.AdminUserRepository
	lineBotService       services.LineBotService
	teacherService       *services.TeacherService
}

func NewTeacherController(app *app.App) *TeacherController {
	lineBotService := services.NewLineBotService(app)
	teacherService := services.NewTeacherService(app)

	return &TeacherController{
		app:                  app,
		teacherRepository:    repositories.NewTeacherRepository(app),
		centerRepo:           repositories.NewCenterRepository(app),
		scheduleRuleRepo:     repositories.NewScheduleRuleRepository(app),
		exceptionRepo:        repositories.NewScheduleExceptionRepository(app),
		expansionService:     services.NewScheduleExpansionService(app),
		scheduleQueryService: services.NewScheduleQueryService(app),
		sessionNoteRepo:      repositories.NewSessionNoteRepository(app),
		adminUserRepo:        repositories.NewAdminUserRepository(app),
		lineBotService:       lineBotService,
		teacherService:       teacherService,
	}
}

// PublicRegister 公開註冊老師（LINE Bot 自主註冊）
// POST /api/v1/teacher/public/register
func (c *TeacherController) PublicRegister(ctx *gin.Context) {
	var req services.PublicRegisterRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(400, gin.H{
			"code":    400,
			"message": "無效的請求參數：" + err.Error(),
		})
		return
	}

	result, errInfo, err := c.teacherService.RegisterPublic(ctx.Request.Context(), &req)
	if err != nil {
		status := 500
		errorCode := 500
		if errInfo != nil {
			errorCode = int(errInfo.Code)
			switch errInfo.Code {
			case 120409: // DUPLICATE
				status = 409
			case 120404: // NOT_FOUND
				status = 404
			case 120403: // FORBIDDEN
				status = 403
			}
		}
		ctx.JSON(status, gin.H{
			"code":    errorCode,
			"message": err.Error(),
		})
		return
	}

	ctx.JSON(200, gin.H{
		"code":    0,
		"message": "註冊成功",
		"data":    result,
	})
}
