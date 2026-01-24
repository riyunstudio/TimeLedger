package servers

import (
	"log"
	"net/http"
	"timeLedger/app/controllers"
	"timeLedger/app/middleware"
	"timeLedger/app/services"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

type route struct {
	Method      string
	Path        string
	Controller  func(*gin.Context)
	Middlewares []gin.HandlerFunc
}

type actions struct {
	user              *controllers.UserController
	auth              *controllers.AuthController
	teacher           *controllers.TeacherController
	adminResource     *controllers.AdminResourceController
	offering          *controllers.OfferingController
	timetableTemplate *controllers.TimetableTemplateController
	adminUser         *controllers.AdminUserController
	scheduling        *controllers.SchedulingController
	smartMatching     *controllers.SmartMatchingController
	notification      *controllers.NotificationController
	export            *controllers.ExportController
}

// 載入路由
func (s *Server) LoadRoutes() {
	authService := services.NewAuthService(s.app)
	authMiddleware := middleware.NewAuthMiddleware(s.app, authService)

	s.routes = []route{
		// Auth
		{http.MethodPost, "/api/v1/auth/admin/login", s.action.auth.AdminLogin, []gin.HandlerFunc{}},
		{http.MethodPost, "/api/v1/auth/teacher/line/login", s.action.auth.TeacherLineLogin, []gin.HandlerFunc{}},
		{http.MethodPost, "/api/v1/auth/refresh", s.action.auth.RefreshToken, []gin.HandlerFunc{}},
		{http.MethodPost, "/api/v1/auth/logout", s.action.auth.Logout, []gin.HandlerFunc{authMiddleware.Authenticate()}},

		// User
		{http.MethodGet, "/api/v1/user", s.action.user.Get, []gin.HandlerFunc{authMiddleware.Authenticate()}},
		{http.MethodPost, "/api/v1/user", s.action.user.Create, []gin.HandlerFunc{}},
		{http.MethodPut, "/api/v1/user", s.action.user.Update, []gin.HandlerFunc{authMiddleware.Authenticate()}},

		// Teacher - Profile
		{http.MethodGet, "/api/v1/teacher/me/profile", s.action.teacher.GetProfile, []gin.HandlerFunc{authMiddleware.Authenticate()}},
		{http.MethodPut, "/api/v1/teacher/me/profile", s.action.teacher.UpdateProfile, []gin.HandlerFunc{authMiddleware.Authenticate()}},
		{http.MethodGet, "/api/v1/teacher/me/centers", s.action.teacher.GetCenters, []gin.HandlerFunc{authMiddleware.Authenticate()}},
		{http.MethodGet, "/api/v1/teacher/me/skills", s.action.teacher.GetSkills, []gin.HandlerFunc{authMiddleware.Authenticate()}},
		{http.MethodPost, "/api/v1/teacher/me/skills", s.action.teacher.CreateSkill, []gin.HandlerFunc{authMiddleware.Authenticate()}},
		{http.MethodDelete, "/api/v1/teacher/me/skills/:id", s.action.teacher.DeleteSkill, []gin.HandlerFunc{authMiddleware.Authenticate()}},
		{http.MethodGet, "/api/v1/teacher/me/certificates", s.action.teacher.GetCertificates, []gin.HandlerFunc{authMiddleware.Authenticate()}},
		{http.MethodPost, "/api/v1/teacher/me/certificates", s.action.teacher.UploadCertificate, []gin.HandlerFunc{authMiddleware.Authenticate()}},
		{http.MethodDelete, "/api/v1/teacher/me/certificates/:id", s.action.teacher.DeleteCertificate, []gin.HandlerFunc{authMiddleware.Authenticate()}},
		{http.MethodGet, "/api/v1/teacher/me/personal-events", s.action.teacher.GetPersonalEvents, []gin.HandlerFunc{authMiddleware.Authenticate()}},
		{http.MethodPost, "/api/v1/teacher/me/personal-events", s.action.teacher.CreatePersonalEvent, []gin.HandlerFunc{authMiddleware.Authenticate()}},
		{http.MethodPatch, "/api/v1/teacher/me/personal-events/:id", s.action.teacher.UpdatePersonalEvent, []gin.HandlerFunc{authMiddleware.Authenticate()}},
		{http.MethodDelete, "/api/v1/teacher/me/personal-events/:id", s.action.teacher.DeletePersonalEvent, []gin.HandlerFunc{authMiddleware.Authenticate()}},
		{http.MethodGet, "/api/v1/teacher/me/schedule", s.action.teacher.GetSchedule, []gin.HandlerFunc{authMiddleware.Authenticate()}},
		{http.MethodGet, "/api/v1/teacher/sessions/note", s.action.teacher.GetSessionNote, []gin.HandlerFunc{authMiddleware.Authenticate()}},
		{http.MethodPut, "/api/v1/teacher/sessions/note", s.action.teacher.UpsertSessionNote, []gin.HandlerFunc{authMiddleware.Authenticate()}},
		{http.MethodGet, "/api/v1/teacher/exceptions", s.action.teacher.GetExceptions, []gin.HandlerFunc{authMiddleware.Authenticate()}},
		{http.MethodPost, "/api/v1/teacher/exceptions", s.action.teacher.CreateException, []gin.HandlerFunc{authMiddleware.Authenticate()}},
		{http.MethodPost, "/api/v1/teacher/exceptions/:id/revoke", s.action.teacher.RevokeException, []gin.HandlerFunc{authMiddleware.Authenticate()}},
		{http.MethodPost, "/api/v1/teacher/scheduling/check-rule-lock", s.action.teacher.CheckRuleLockStatus, []gin.HandlerFunc{authMiddleware.Authenticate()}},
		{http.MethodPost, "/api/v1/teacher/scheduling/preview-recurrence-edit", s.action.teacher.PreviewRecurrenceEdit, []gin.HandlerFunc{authMiddleware.Authenticate()}},
		{http.MethodPost, "/api/v1/teacher/scheduling/edit-recurring", s.action.teacher.EditRecurringSchedule, []gin.HandlerFunc{authMiddleware.Authenticate()}},
		{http.MethodPost, "/api/v1/teacher/scheduling/delete-recurring", s.action.teacher.DeleteRecurringSchedule, []gin.HandlerFunc{authMiddleware.Authenticate()}},

		// Teacher - Admin
		{http.MethodGet, "/api/v1/teachers", s.action.teacher.ListTeachers, []gin.HandlerFunc{authMiddleware.Authenticate(), authMiddleware.RequireAdmin()}},
		{http.MethodDelete, "/api/v1/teachers/:id", s.action.teacher.DeleteTeacher, []gin.HandlerFunc{authMiddleware.Authenticate(), authMiddleware.RequireAdmin()}},
		{http.MethodPost, "/api/v1/admin/centers/:id/invitations", s.action.teacher.InviteTeacher, []gin.HandlerFunc{authMiddleware.Authenticate(), authMiddleware.RequireCenterAdmin()}},

		// Admin - Offerings
		{http.MethodGet, "/api/v1/admin/offerings", s.action.offering.GetOfferings, []gin.HandlerFunc{authMiddleware.Authenticate(), authMiddleware.RequireCenterAdmin()}},
		{http.MethodPost, "/api/v1/admin/offerings", s.action.offering.CreateOffering, []gin.HandlerFunc{authMiddleware.Authenticate(), authMiddleware.RequireCenterAdmin()}},
		{http.MethodPut, "/api/v1/admin/offerings/:offering_id", s.action.offering.UpdateOffering, []gin.HandlerFunc{authMiddleware.Authenticate(), authMiddleware.RequireCenterAdmin()}},
		{http.MethodDelete, "/api/v1/admin/offerings/:offering_id", s.action.offering.DeleteOffering, []gin.HandlerFunc{authMiddleware.Authenticate(), authMiddleware.RequireCenterAdmin()}},
		{http.MethodPost, "/api/v1/admin/centers/:id/offerings/:offering_id/copy", s.action.offering.CopyOffering, []gin.HandlerFunc{authMiddleware.Authenticate(), authMiddleware.RequireCenterAdmin()}},

		// Admin - Timetable Template
		{http.MethodGet, "/api/v1/admin/centers/:id/templates", s.action.timetableTemplate.GetTemplates, []gin.HandlerFunc{authMiddleware.Authenticate(), authMiddleware.RequireCenterAdmin()}},
		{http.MethodPost, "/api/v1/admin/centers/:id/templates", s.action.timetableTemplate.CreateTemplate, []gin.HandlerFunc{authMiddleware.Authenticate(), authMiddleware.RequireCenterAdmin()}},
		{http.MethodPut, "/api/v1/admin/centers/:id/templates/:templateId", s.action.timetableTemplate.UpdateTemplate, []gin.HandlerFunc{authMiddleware.Authenticate(), authMiddleware.RequireCenterAdmin()}},
		{http.MethodDelete, "/api/v1/admin/centers/:id/templates/:templateId", s.action.timetableTemplate.DeleteTemplate, []gin.HandlerFunc{authMiddleware.Authenticate(), authMiddleware.RequireCenterAdmin()}},
		{http.MethodGet, "/api/v1/admin/centers/:id/templates/:templateId/cells", s.action.timetableTemplate.GetCells, []gin.HandlerFunc{authMiddleware.Authenticate(), authMiddleware.RequireCenterAdmin()}},
		{http.MethodPost, "/api/v1/admin/centers/:id/templates/:templateId/cells", s.action.timetableTemplate.CreateCells, []gin.HandlerFunc{authMiddleware.Authenticate(), authMiddleware.RequireCenterAdmin()}},

		// Admin - Users
		{http.MethodGet, "/api/v1/admin/centers/:id/users", s.action.adminUser.GetAdminUsers, []gin.HandlerFunc{authMiddleware.Authenticate(), authMiddleware.RequireCenterAdmin()}},
		{http.MethodPost, "/api/v1/admin/centers/:id/users", s.action.adminUser.CreateAdminUser, []gin.HandlerFunc{authMiddleware.Authenticate(), authMiddleware.RequireCenterAdmin()}},
		{http.MethodPut, "/api/v1/admin/centers/:id/users/:admin_id", s.action.adminUser.UpdateAdminUser, []gin.HandlerFunc{authMiddleware.Authenticate(), authMiddleware.RequireCenterAdmin()}},
		{http.MethodDelete, "/api/v1/admin/centers/:id/users/:admin_id", s.action.adminUser.DeleteAdminUser, []gin.HandlerFunc{authMiddleware.Authenticate(), authMiddleware.RequireCenterAdmin()}},

		// Admin - Scheduling
		{http.MethodPost, "/api/v1/admin/scheduling/check-overlap", s.action.scheduling.CheckOverlap, []gin.HandlerFunc{authMiddleware.Authenticate(), authMiddleware.RequireCenterAdmin()}},
		{http.MethodPost, "/api/v1/admin/scheduling/check-teacher-buffer", s.action.scheduling.CheckTeacherBuffer, []gin.HandlerFunc{authMiddleware.Authenticate(), authMiddleware.RequireCenterAdmin()}},
		{http.MethodPost, "/api/v1/admin/scheduling/check-room-buffer", s.action.scheduling.CheckRoomBuffer, []gin.HandlerFunc{authMiddleware.Authenticate(), authMiddleware.RequireCenterAdmin()}},
		{http.MethodPost, "/api/v1/admin/scheduling/validate", s.action.scheduling.ValidateFull, []gin.HandlerFunc{authMiddleware.Authenticate(), authMiddleware.RequireCenterAdmin()}},
		{http.MethodGet, "/api/v1/admin/centers/:id/rules", s.action.scheduling.GetRules, []gin.HandlerFunc{authMiddleware.Authenticate(), authMiddleware.RequireCenterAdmin()}},
		{http.MethodPost, "/api/v1/admin/centers/:id/rules", s.action.scheduling.CreateRule, []gin.HandlerFunc{authMiddleware.Authenticate(), authMiddleware.RequireCenterAdmin()}},
		{http.MethodDelete, "/api/v1/admin/centers/:id/rules/:ruleId", s.action.scheduling.DeleteRule, []gin.HandlerFunc{authMiddleware.Authenticate(), authMiddleware.RequireCenterAdmin()}},
		{http.MethodPost, "/api/v1/admin/scheduling/exceptions", s.action.scheduling.CreateException, []gin.HandlerFunc{authMiddleware.Authenticate(), authMiddleware.RequireCenterAdmin()}},
		{http.MethodPost, "/api/v1/admin/scheduling/exceptions/:exceptionId/review", s.action.scheduling.ReviewException, []gin.HandlerFunc{authMiddleware.Authenticate(), authMiddleware.RequireCenterAdmin()}},
		{http.MethodGet, "/api/v1/admin/rules/:ruleId/exceptions", s.action.scheduling.GetExceptionsByRule, []gin.HandlerFunc{authMiddleware.Authenticate(), authMiddleware.RequireCenterAdmin()}},
		{http.MethodGet, "/api/v1/admin/centers/:id/exceptions", s.action.scheduling.GetExceptionsByDateRange, []gin.HandlerFunc{authMiddleware.Authenticate(), authMiddleware.RequireCenterAdmin()}},
		{http.MethodPost, "/api/v1/admin/centers/:id/expand-rules", s.action.scheduling.ExpandRules, []gin.HandlerFunc{authMiddleware.Authenticate(), authMiddleware.RequireCenterAdmin()}},
		{http.MethodPost, "/api/v1/admin/centers/:id/detect-phase-transitions", s.action.scheduling.DetectPhaseTransitions, []gin.HandlerFunc{authMiddleware.Authenticate(), authMiddleware.RequireCenterAdmin()}},
		{http.MethodPost, "/api/v1/admin/scheduling/check-rule-lock", s.action.scheduling.CheckRuleLockStatus, []gin.HandlerFunc{authMiddleware.Authenticate(), authMiddleware.RequireCenterAdmin()}},

		// Admin - Resources
		{http.MethodGet, "/api/v1/admin/rooms", s.action.adminResource.GetRooms, []gin.HandlerFunc{authMiddleware.Authenticate(), authMiddleware.RequireCenterAdmin()}},
		{http.MethodPost, "/api/v1/admin/rooms", s.action.adminResource.CreateRoom, []gin.HandlerFunc{authMiddleware.Authenticate(), authMiddleware.RequireCenterAdmin()}},
		{http.MethodPut, "/api/v1/admin/rooms/:room_id", s.action.adminResource.UpdateRoom, []gin.HandlerFunc{authMiddleware.Authenticate(), authMiddleware.RequireCenterAdmin()}},
		{http.MethodGet, "/api/v1/admin/rooms/active", s.action.adminResource.GetActiveRooms, []gin.HandlerFunc{authMiddleware.Authenticate(), authMiddleware.RequireCenterAdmin()}},
		{http.MethodPatch, "/api/v1/admin/rooms/:room_id/toggle-active", s.action.adminResource.ToggleRoomActive, []gin.HandlerFunc{authMiddleware.Authenticate(), authMiddleware.RequireCenterAdmin()}},
		{http.MethodGet, "/api/v1/admin/courses", s.action.adminResource.GetCourses, []gin.HandlerFunc{authMiddleware.Authenticate(), authMiddleware.RequireCenterAdmin()}},
		{http.MethodPost, "/api/v1/admin/courses", s.action.adminResource.CreateCourse, []gin.HandlerFunc{authMiddleware.Authenticate(), authMiddleware.RequireCenterAdmin()}},
		{http.MethodPut, "/api/v1/admin/courses/:course_id", s.action.adminResource.UpdateCourse, []gin.HandlerFunc{authMiddleware.Authenticate(), authMiddleware.RequireCenterAdmin()}},
		{http.MethodDelete, "/api/v1/admin/courses/:course_id", s.action.adminResource.DeleteCourse, []gin.HandlerFunc{authMiddleware.Authenticate(), authMiddleware.RequireCenterAdmin()}},
		{http.MethodGet, "/api/v1/admin/courses/active", s.action.adminResource.GetActiveCourses, []gin.HandlerFunc{authMiddleware.Authenticate(), authMiddleware.RequireCenterAdmin()}},
		{http.MethodPatch, "/api/v1/admin/courses/:course_id/toggle-active", s.action.adminResource.ToggleCourseActive, []gin.HandlerFunc{authMiddleware.Authenticate(), authMiddleware.RequireCenterAdmin()}},
		{http.MethodGet, "/api/v1/admin/offerings/active", s.action.adminResource.GetActiveOfferings, []gin.HandlerFunc{authMiddleware.Authenticate(), authMiddleware.RequireCenterAdmin()}},
		{http.MethodPatch, "/api/v1/admin/offerings/:offering_id/toggle-active", s.action.adminResource.ToggleOfferingActive, []gin.HandlerFunc{authMiddleware.Authenticate(), authMiddleware.RequireCenterAdmin()}},

		// Admin - Invitations
		{http.MethodGet, "/api/v1/admin/centers/:id/invitations", s.action.adminResource.GetInvitations, []gin.HandlerFunc{authMiddleware.Authenticate(), authMiddleware.RequireCenterAdmin()}},
		{http.MethodGet, "/api/v1/admin/centers/:id/invitations/stats", s.action.adminResource.GetInvitationStats, []gin.HandlerFunc{authMiddleware.Authenticate(), authMiddleware.RequireCenterAdmin()}},

		// Admin - Holidays
		{http.MethodGet, "/api/v1/admin/centers/:id/holidays", s.action.adminResource.GetHolidays, []gin.HandlerFunc{authMiddleware.Authenticate(), authMiddleware.RequireCenterAdmin()}},
		{http.MethodPost, "/api/v1/admin/centers/:id/holidays", s.action.adminResource.CreateHoliday, []gin.HandlerFunc{authMiddleware.Authenticate(), authMiddleware.RequireCenterAdmin()}},
		{http.MethodDelete, "/api/v1/admin/centers/:id/holidays/:holiday_id", s.action.adminResource.DeleteHoliday, []gin.HandlerFunc{authMiddleware.Authenticate(), authMiddleware.RequireCenterAdmin()}},
		{http.MethodPost, "/api/v1/admin/centers/:id/holidays/bulk", s.action.adminResource.BulkCreateHolidays, []gin.HandlerFunc{authMiddleware.Authenticate(), authMiddleware.RequireCenterAdmin()}},

		// Smart Matching
		{http.MethodPost, "/api/v1/admin/smart-matching/matches", s.action.smartMatching.FindMatches, []gin.HandlerFunc{authMiddleware.Authenticate(), authMiddleware.RequireCenterAdmin()}},
		{http.MethodGet, "/api/v1/admin/smart-matching/talent/search", s.action.smartMatching.SearchTalent, []gin.HandlerFunc{authMiddleware.Authenticate(), authMiddleware.RequireCenterAdmin()}},

		// Notification
		{http.MethodGet, "/api/v1/notifications", s.action.notification.ListNotifications, []gin.HandlerFunc{authMiddleware.Authenticate()}},
		{http.MethodGet, "/api/v1/notifications/unread-count", s.action.notification.GetUnreadCount, []gin.HandlerFunc{authMiddleware.Authenticate()}},
		{http.MethodPost, "/api/v1/notifications/:id/read", s.action.notification.MarkAsRead, []gin.HandlerFunc{authMiddleware.Authenticate()}},
		{http.MethodPost, "/api/v1/notifications/read-all", s.action.notification.MarkAllAsRead, []gin.HandlerFunc{authMiddleware.Authenticate()}},
		{http.MethodPost, "/api/v1/notifications/token", s.action.notification.SetNotifyToken, []gin.HandlerFunc{authMiddleware.Authenticate()}},
		{http.MethodPost, "/api/v1/notifications/test", s.action.notification.SendTestNotification, []gin.HandlerFunc{authMiddleware.Authenticate()}},

		// Export
		{http.MethodPost, "/api/v1/admin/export/schedule/csv", s.action.export.ExportScheduleCSV, []gin.HandlerFunc{authMiddleware.Authenticate(), authMiddleware.RequireCenterAdmin()}},
		{http.MethodPost, "/api/v1/admin/export/schedule/pdf", s.action.export.ExportSchedulePDF, []gin.HandlerFunc{authMiddleware.Authenticate(), authMiddleware.RequireCenterAdmin()}},
		{http.MethodGet, "/api/v1/admin/centers/:id/export/teachers/csv", s.action.export.ExportTeachersCSV, []gin.HandlerFunc{authMiddleware.Authenticate(), authMiddleware.RequireCenterAdmin()}},
		{http.MethodGet, "/api/v1/admin/centers/:id/export/exceptions/csv", s.action.export.ExportExceptionsCSV, []gin.HandlerFunc{authMiddleware.Authenticate(), authMiddleware.RequireCenterAdmin()}},
	}
}

// 註冊路由
func (s *Server) registerRoutes(r *gin.Engine) {
	// 健康檢查
	r.GET("/healthy", func(c *gin.Context) {
		c.String(http.StatusOK, "Healthy")
	})

	// Debug - 直接在 middleware 之前
	r.GET("/xyz123", func(c *gin.Context) {
		c.String(http.StatusOK, "xyz123 works!")
	})
	r.GET("/xyz123/", func(c *gin.Context) {
		c.String(http.StatusOK, "xyz123/ works!")
	})

	// Swagger
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// 全局 Middleware - 先註冊 middleware
	r.Use(s.InitMiddleware())
	r.Use(s.RecoverMiddleware())
	r.Use(s.MainMiddleware())

	// 註冊所有路由
	for _, rt := range s.routes {
		switch rt.Method {
		case http.MethodGet:
			r.GET(rt.Path, append(rt.Middlewares, rt.Controller)...)
		case http.MethodPost:
			r.POST(rt.Path, append(rt.Middlewares, rt.Controller)...)
		case http.MethodPut:
			r.PUT(rt.Path, append(rt.Middlewares, rt.Controller)...)
		case http.MethodPatch:
			r.PATCH(rt.Path, append(rt.Middlewares, rt.Controller)...)
		case http.MethodDelete:
			r.DELETE(rt.Path, append(rt.Middlewares, rt.Controller)...)
		}
	}

	log.Printf("Total routes registered: %d", len(s.routes))
}

// 建構控制器
func (s *Server) NewControllers() {
	controllers.NewBaseController(s.app)
	s.action.user = controllers.NewUserController(s.app)

	authService := services.NewAuthService(s.app)
	s.action.auth = controllers.NewAuthController(s.app, authService)
	s.action.teacher = controllers.NewTeacherController(s.app)
	s.action.adminResource = controllers.NewAdminResourceController(s.app)
	s.action.offering = controllers.NewOfferingController(s.app)
	s.action.timetableTemplate = controllers.NewTimetableTemplateController(s.app)
	s.action.adminUser = controllers.NewAdminUserController(s.app)
	s.action.scheduling = controllers.NewSchedulingController(s.app)
	s.action.smartMatching = controllers.NewSmartMatchingController(s.app)
	s.action.notification = controllers.NewNotificationController(s.app)
	s.action.export = controllers.NewExportController(s.app)
}
