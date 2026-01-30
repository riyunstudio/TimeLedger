package controllers

import (
	"timeLedger/app"
	"timeLedger/app/repositories"
	"timeLedger/app/services"
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
}

func NewTeacherController(app *app.App) *TeacherController {
	lineBotService := services.NewLineBotService(app)

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
	}
}
