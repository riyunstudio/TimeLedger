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
	auth              *controllers.AuthController
	teacher           *controllers.TeacherController
	adminTeacher      *controllers.AdminTeacherController
	adminCenter       *controllers.AdminCenterController
	adminRoom         *controllers.AdminRoomController
	adminCourse       *controllers.AdminCourseController
	adminHoliday      *controllers.AdminHolidayController
	teacherProfile    *controllers.TeacherProfileController
	teacherSchedule   *controllers.TeacherScheduleController
	teacherSession    *controllers.TeacherSessionController
	teacherEvent      *controllers.TeacherEventController
	teacherException  *controllers.TeacherExceptionController
	teacherInvitation *controllers.TeacherInvitationController
	geo               *controllers.GeoController
	adminResource     *controllers.AdminResourceController
	offering          *controllers.OfferingController
	timetableTemplate *controllers.TimetableTemplateController
	adminUser         *controllers.AdminUserController
	scheduling        *controllers.SchedulingController
	smartMatching     *controllers.SmartMatchingController
	notification      *controllers.NotificationController
	adminNotification *controllers.AdminNotificationController
	export            *controllers.ExportController
	lineBot           *controllers.LineBotController
	r2Test            *controllers.R2TestController
}

// 載入路由
func (s *Server) LoadRoutes() {
	authService := services.NewAuthService(s.app)
	authMiddleware := middleware.NewAuthMiddleware(s.app, authService)

	s.routes = []route{
		// R2 Storage Test - 無需認證
		{http.MethodGet, "/api/test/r2-status", s.action.r2Test.StatusTest, []gin.HandlerFunc{}},
		{http.MethodPost, "/api/test/upload", s.action.r2Test.UploadTest, []gin.HandlerFunc{}},
		{http.MethodPost, "/api/test/upload-batch", s.action.r2Test.BatchUploadTest, []gin.HandlerFunc{}},

		// Auth
		{http.MethodPost, "/api/v1/auth/admin/login", s.action.auth.AdminLogin, []gin.HandlerFunc{}},
		{http.MethodPost, "/api/v1/auth/teacher/line/login", s.action.auth.TeacherLineLogin, []gin.HandlerFunc{}},
		{http.MethodPost, "/api/v1/auth/refresh", s.action.auth.RefreshToken, []gin.HandlerFunc{}},
		{http.MethodPost, "/api/v1/auth/logout", s.action.auth.Logout, []gin.HandlerFunc{authMiddleware.Authenticate()}},

		// Geo - Cities & Districts
		{http.MethodGet, "/api/v1/geo/cities", s.action.geo.ListCities, []gin.HandlerFunc{authMiddleware.Authenticate()}},

		// Teacher - Profile
		{http.MethodGet, "/api/v1/teacher/me/profile", s.action.teacherProfile.GetProfile, []gin.HandlerFunc{authMiddleware.Authenticate()}},
		{http.MethodPut, "/api/v1/teacher/me/profile", s.action.teacherProfile.UpdateProfile, []gin.HandlerFunc{authMiddleware.Authenticate()}},
		// Teacher - Centers
		{http.MethodGet, "/api/v1/teacher/me/centers", s.action.teacherProfile.GetCenters, []gin.HandlerFunc{authMiddleware.Authenticate()}},
		// Teacher - Skills
		{http.MethodGet, "/api/v1/teacher/me/skills", s.action.teacherProfile.GetSkills, []gin.HandlerFunc{authMiddleware.Authenticate()}},
		{http.MethodPost, "/api/v1/teacher/me/skills", s.action.teacherProfile.CreateSkill, []gin.HandlerFunc{authMiddleware.Authenticate()}},
		{http.MethodPut, "/api/v1/teacher/me/skills/:id", s.action.teacherProfile.UpdateSkill, []gin.HandlerFunc{authMiddleware.Authenticate()}},
		{http.MethodDelete, "/api/v1/teacher/me/skills/:id", s.action.teacherProfile.DeleteSkill, []gin.HandlerFunc{authMiddleware.Authenticate()}},
		// Teacher - Certificates
		{http.MethodGet, "/api/v1/teacher/me/certificates", s.action.teacherProfile.GetCertificates, []gin.HandlerFunc{authMiddleware.Authenticate()}},
		{http.MethodPost, "/api/v1/teacher/me/certificates", s.action.teacherProfile.CreateCertificate, []gin.HandlerFunc{authMiddleware.Authenticate()}},
		{http.MethodPost, "/api/v1/teacher/me/certificates/upload", s.action.teacherProfile.UploadCertificateFile, []gin.HandlerFunc{authMiddleware.Authenticate()}},
		{http.MethodDelete, "/api/v1/teacher/me/certificates/:id", s.action.teacherProfile.DeleteCertificate, []gin.HandlerFunc{authMiddleware.Authenticate()}},
		// Hashtags
		{http.MethodGet, "/api/v1/hashtags/search", s.action.teacherProfile.SearchHashtags, []gin.HandlerFunc{authMiddleware.Authenticate()}},
		{http.MethodPost, "/api/v1/hashtags", s.action.teacherProfile.CreateHashtag, []gin.HandlerFunc{authMiddleware.Authenticate()}},
		// Teacher - Personal Events
		{http.MethodGet, "/api/v1/teacher/me/personal-events", s.action.teacherEvent.GetPersonalEvents, []gin.HandlerFunc{authMiddleware.Authenticate()}},
		{http.MethodPost, "/api/v1/teacher/me/personal-events", s.action.teacherEvent.CreatePersonalEvent, []gin.HandlerFunc{authMiddleware.Authenticate()}},
		{http.MethodPatch, "/api/v1/teacher/me/personal-events/:id", s.action.teacherEvent.UpdatePersonalEvent, []gin.HandlerFunc{authMiddleware.Authenticate()}},
		{http.MethodDelete, "/api/v1/teacher/me/personal-events/:id", s.action.teacherEvent.DeletePersonalEvent, []gin.HandlerFunc{authMiddleware.Authenticate()}},
		{http.MethodGet, "/api/v1/teacher/me/personal-events/:id/note", s.action.teacherEvent.GetPersonalEventNote, []gin.HandlerFunc{authMiddleware.Authenticate()}},
		{http.MethodPut, "/api/v1/teacher/me/personal-events/:id/note", s.action.teacherEvent.UpdatePersonalEventNote, []gin.HandlerFunc{authMiddleware.Authenticate()}},
		// Teacher - Schedule
		{http.MethodGet, "/api/v1/teacher/me/schedule", s.action.teacherSchedule.GetSchedule, []gin.HandlerFunc{authMiddleware.Authenticate()}},
		{http.MethodGet, "/api/v1/teacher/schedules", s.action.teacherSchedule.GetSchedules, []gin.HandlerFunc{authMiddleware.Authenticate()}},
		{http.MethodGet, "/api/v1/teacher/me/centers/:center_id/schedule-rules", s.action.teacherSchedule.GetCenterScheduleRules, []gin.HandlerFunc{authMiddleware.Authenticate()}},
		// Teacher - Sessions
		{http.MethodGet, "/api/v1/teacher/sessions/note", s.action.teacherSession.GetSessionNote, []gin.HandlerFunc{authMiddleware.Authenticate()}},
		{http.MethodPut, "/api/v1/teacher/sessions/note", s.action.teacherSession.UpsertSessionNote, []gin.HandlerFunc{authMiddleware.Authenticate()}},
		// Teacher - Exceptions
		{http.MethodGet, "/api/v1/teacher/exceptions", s.action.teacherException.GetExceptions, []gin.HandlerFunc{authMiddleware.Authenticate()}},
		{http.MethodPost, "/api/v1/teacher/exceptions", s.action.teacherException.CreateException, []gin.HandlerFunc{authMiddleware.Authenticate()}},
		{http.MethodPost, "/api/v1/teacher/exceptions/:id/revoke", s.action.teacherException.RevokeException, []gin.HandlerFunc{authMiddleware.Authenticate()}},
		// Teacher - Scheduling
		{http.MethodPost, "/api/v1/teacher/scheduling/check-rule-lock", s.action.teacherSchedule.CheckRuleLockStatus, []gin.HandlerFunc{authMiddleware.Authenticate()}},
		{http.MethodPost, "/api/v1/teacher/scheduling/preview-recurrence-edit", s.action.teacherSchedule.PreviewRecurrenceEdit, []gin.HandlerFunc{authMiddleware.Authenticate()}},
		{http.MethodPost, "/api/v1/teacher/scheduling/edit-recurring", s.action.teacherSchedule.EditRecurringSchedule, []gin.HandlerFunc{authMiddleware.Authenticate()}},
		{http.MethodPost, "/api/v1/teacher/scheduling/delete-recurring", s.action.teacherSchedule.DeleteRecurringSchedule, []gin.HandlerFunc{authMiddleware.Authenticate()}},

		// Teacher - Invitations
		{http.MethodGet, "/api/v1/teacher/me/invitations", s.action.teacherInvitation.GetTeacherInvitations, []gin.HandlerFunc{authMiddleware.Authenticate()}},
		{http.MethodPost, "/api/v1/teacher/me/invitations/respond", s.action.teacherInvitation.RespondToInvitation, []gin.HandlerFunc{authMiddleware.Authenticate()}},
		{http.MethodGet, "/api/v1/teacher/me/invitations/pending-count", s.action.teacherInvitation.GetPendingInvitationsCount, []gin.HandlerFunc{authMiddleware.Authenticate()}},
		// Teacher - Public Registration (LINE Bot 自主註冊)
		{http.MethodPost, "/api/v1/teacher/public/register", s.action.teacher.PublicRegister, []gin.HandlerFunc{}},

		// Admin - Teacher Management
		{http.MethodGet, "/api/v1/teachers", s.action.adminTeacher.ListTeachers, []gin.HandlerFunc{authMiddleware.Authenticate(), authMiddleware.RequireAdmin()}},
		{http.MethodPost, "/api/v1/admin/teachers/placeholder", s.action.adminTeacher.CreatePlaceholderTeacher, []gin.HandlerFunc{authMiddleware.Authenticate(), authMiddleware.RequireCenterAdmin()}},
		{http.MethodPost, "/api/v1/admin/teachers/merge", s.action.adminTeacher.MergeTeachers, []gin.HandlerFunc{authMiddleware.Authenticate(), authMiddleware.RequireCenterAdmin()}},
		{http.MethodDelete, "/api/v1/teachers/:id", s.action.adminTeacher.DeleteTeacher, []gin.HandlerFunc{authMiddleware.Authenticate(), authMiddleware.RequireAdmin()}},
		{http.MethodDelete, "/api/v1/admin/centers/:id/teachers/:teacher_id", s.action.adminTeacher.RemoveFromCenter, []gin.HandlerFunc{authMiddleware.Authenticate(), authMiddleware.RequireCenterAdmin()}},
		{http.MethodPost, "/api/v1/admin/centers/:id/invitations", s.action.teacherInvitation.InviteTeacher, []gin.HandlerFunc{authMiddleware.Authenticate(), authMiddleware.RequireCenterAdmin()}},

		// Admin - Center Management
		{http.MethodGet, "/api/v1/admin/centers", s.action.adminCenter.GetCenters, []gin.HandlerFunc{authMiddleware.Authenticate(), authMiddleware.RequireCenterAdmin()}},
		{http.MethodPost, "/api/v1/admin/centers", s.action.adminCenter.CreateCenter, []gin.HandlerFunc{authMiddleware.Authenticate(), authMiddleware.RequireCenterAdmin()}},

		// Admin - Teacher Resources
		{http.MethodGet, "/api/v1/admin/teachers", s.action.adminResource.GetTeachers, []gin.HandlerFunc{authMiddleware.Authenticate(), authMiddleware.RequireCenterAdmin()}},

		// Admin - Offerings
		{http.MethodGet, "/api/v1/admin/offerings", s.action.offering.GetOfferings, []gin.HandlerFunc{authMiddleware.Authenticate(), authMiddleware.RequireCenterAdmin()}},
		{http.MethodPost, "/api/v1/admin/offerings", s.action.offering.CreateOffering, []gin.HandlerFunc{authMiddleware.Authenticate(), authMiddleware.RequireCenterAdmin()}},
		{http.MethodPut, "/api/v1/admin/offerings/:offering_id", s.action.offering.UpdateOffering, []gin.HandlerFunc{authMiddleware.Authenticate(), authMiddleware.RequireCenterAdmin()}},
		{http.MethodDelete, "/api/v1/admin/offerings/:offering_id", s.action.offering.DeleteOffering, []gin.HandlerFunc{authMiddleware.Authenticate(), authMiddleware.RequireCenterAdmin()}},
		{http.MethodPost, "/api/v1/admin/centers/:id/offerings/:offering_id/copy", s.action.offering.CopyOffering, []gin.HandlerFunc{authMiddleware.Authenticate(), authMiddleware.RequireCenterAdmin()}},

		// Admin - Timetable Template
		{http.MethodGet, "/api/v1/admin/templates", s.action.timetableTemplate.GetTemplates, []gin.HandlerFunc{authMiddleware.Authenticate(), authMiddleware.RequireCenterAdmin()}},
		{http.MethodPost, "/api/v1/admin/templates", s.action.timetableTemplate.CreateTemplate, []gin.HandlerFunc{authMiddleware.Authenticate(), authMiddleware.RequireCenterAdmin()}},
		{http.MethodPut, "/api/v1/admin/templates/:templateId", s.action.timetableTemplate.UpdateTemplate, []gin.HandlerFunc{authMiddleware.Authenticate(), authMiddleware.RequireCenterAdmin()}},
		{http.MethodDelete, "/api/v1/admin/templates/:templateId", s.action.timetableTemplate.DeleteTemplate, []gin.HandlerFunc{authMiddleware.Authenticate(), authMiddleware.RequireCenterAdmin()}},
		{http.MethodGet, "/api/v1/admin/templates/:templateId/cells", s.action.timetableTemplate.GetCells, []gin.HandlerFunc{authMiddleware.Authenticate(), authMiddleware.RequireCenterAdmin()}},
		{http.MethodPost, "/api/v1/admin/templates/:templateId/cells", s.action.timetableTemplate.CreateCells, []gin.HandlerFunc{authMiddleware.Authenticate(), authMiddleware.RequireCenterAdmin()}},
		{http.MethodPut, "/api/v1/admin/templates/:templateId/cells/reorder", s.action.timetableTemplate.ReorderCells, []gin.HandlerFunc{authMiddleware.Authenticate(), authMiddleware.RequireCenterAdmin()}},
		{http.MethodDelete, "/api/v1/admin/templates/cells/:cellId", s.action.timetableTemplate.DeleteCell, []gin.HandlerFunc{authMiddleware.Authenticate(), authMiddleware.RequireCenterAdmin()}},
		{http.MethodPost, "/api/v1/admin/templates/:templateId/apply", s.action.timetableTemplate.ApplyTemplate, []gin.HandlerFunc{authMiddleware.Authenticate(), authMiddleware.RequireCenterAdmin()}},
		{http.MethodPost, "/api/v1/admin/templates/:templateId/validate-apply", s.action.timetableTemplate.ValidateApplyTemplate, []gin.HandlerFunc{authMiddleware.Authenticate(), authMiddleware.RequireCenterAdmin()}},

		// Admin - LINE 綁定
		{http.MethodGet, "/api/v1/admin/me/line-binding", s.action.adminUser.GetLINEBindingStatus, []gin.HandlerFunc{authMiddleware.Authenticate(), authMiddleware.RequireCenterAdmin()}},
		{http.MethodPost, "/api/v1/admin/me/line/bind", s.action.adminUser.InitLINEBinding, []gin.HandlerFunc{authMiddleware.Authenticate(), authMiddleware.RequireCenterAdmin()}},
		{http.MethodDelete, "/api/v1/admin/me/line/unbind", s.action.adminUser.UnbindLINE, []gin.HandlerFunc{authMiddleware.Authenticate(), authMiddleware.RequireCenterAdmin()}},
		{http.MethodGet, "/api/v1/admin/me/line/notify-settings", s.action.adminUser.GetLINENotifySettings, []gin.HandlerFunc{authMiddleware.Authenticate(), authMiddleware.RequireCenterAdmin()}},
		{http.MethodPatch, "/api/v1/admin/me/line/notify-settings", s.action.adminUser.UpdateLINENotifySettings, []gin.HandlerFunc{authMiddleware.Authenticate(), authMiddleware.RequireCenterAdmin()}},
		// Admin - LINE QR Code
		{http.MethodGet, "/api/v1/admin/me/line/qrcode", s.action.lineBot.GenerateLINEBindingQR, []gin.HandlerFunc{authMiddleware.Authenticate(), authMiddleware.RequireCenterAdmin()}},
		{http.MethodGet, "/api/v1/admin/me/line/qrcode-with-code", s.action.lineBot.GenerateVerificationCodeQR, []gin.HandlerFunc{authMiddleware.Authenticate(), authMiddleware.RequireCenterAdmin()}},

		// Admin - Profile
		{http.MethodGet, "/api/v1/admin/me/profile", s.action.adminUser.GetAdminProfile, []gin.HandlerFunc{authMiddleware.Authenticate(), authMiddleware.RequireCenterAdmin()}},
		{http.MethodPost, "/api/v1/admin/me/change-password", s.action.adminUser.ChangePassword, []gin.HandlerFunc{authMiddleware.Authenticate(), authMiddleware.RequireCenterAdmin()}},

		// Admin - Management (僅 OWNER 可執行)
		{http.MethodGet, "/api/v1/admin/admins", s.action.adminUser.ListAdmins, []gin.HandlerFunc{authMiddleware.Authenticate(), authMiddleware.RequireCenterAdmin()}},
		{http.MethodPost, "/api/v1/admin/admins", s.action.adminUser.CreateAdmin, []gin.HandlerFunc{authMiddleware.Authenticate(), authMiddleware.RequireCenterAdmin()}},
		{http.MethodPost, "/api/v1/admin/admins/toggle-status", s.action.adminUser.ToggleAdminStatus, []gin.HandlerFunc{authMiddleware.Authenticate(), authMiddleware.RequireCenterAdmin()}},
		{http.MethodPost, "/api/v1/admin/admins/reset-password", s.action.adminUser.ResetAdminPassword, []gin.HandlerFunc{authMiddleware.Authenticate(), authMiddleware.RequireCenterAdmin()}},
		{http.MethodPost, "/api/v1/admin/admins/change-role", s.action.adminUser.ChangeAdminRole, []gin.HandlerFunc{authMiddleware.Authenticate(), authMiddleware.RequireCenterAdmin()}},

		// Admin - Scheduling
		{http.MethodPost, "/api/v1/admin/scheduling/check-overlap", s.action.scheduling.CheckOverlap, []gin.HandlerFunc{authMiddleware.Authenticate(), authMiddleware.RequireCenterAdmin()}},
		{http.MethodPost, "/api/v1/admin/scheduling/check-teacher-buffer", s.action.scheduling.CheckTeacherBuffer, []gin.HandlerFunc{authMiddleware.Authenticate(), authMiddleware.RequireCenterAdmin()}},
		{http.MethodPost, "/api/v1/admin/scheduling/check-room-buffer", s.action.scheduling.CheckRoomBuffer, []gin.HandlerFunc{authMiddleware.Authenticate(), authMiddleware.RequireCenterAdmin()}},
		{http.MethodPost, "/api/v1/admin/scheduling/validate", s.action.scheduling.ValidateFull, []gin.HandlerFunc{authMiddleware.Authenticate(), authMiddleware.RequireCenterAdmin()}},
		// Dashboard
		{http.MethodGet, "/api/v1/admin/dashboard/today-summary", s.action.scheduling.GetTodaySummary, []gin.HandlerFunc{authMiddleware.Authenticate(), authMiddleware.RequireCenterAdmin()}},
		{http.MethodGet, "/api/v1/admin/rules", s.action.scheduling.GetRules, []gin.HandlerFunc{authMiddleware.Authenticate(), authMiddleware.RequireCenterAdmin()}},
		{http.MethodPost, "/api/v1/admin/rules", s.action.scheduling.CreateRule, []gin.HandlerFunc{authMiddleware.Authenticate(), authMiddleware.RequireCenterAdmin()}},
		{http.MethodPut, "/api/v1/admin/rules/:ruleId", s.action.scheduling.UpdateRule, []gin.HandlerFunc{authMiddleware.Authenticate(), authMiddleware.RequireCenterAdmin()}},
		{http.MethodDelete, "/api/v1/admin/rules/:ruleId", s.action.scheduling.DeleteRule, []gin.HandlerFunc{authMiddleware.Authenticate(), authMiddleware.RequireCenterAdmin()}},
		{http.MethodPost, "/api/v1/admin/scheduling/exceptions", s.action.scheduling.CreateException, []gin.HandlerFunc{authMiddleware.Authenticate(), authMiddleware.RequireCenterAdmin()}},
		{http.MethodPost, "/api/v1/admin/scheduling/exceptions/:exceptionId/review", s.action.scheduling.ReviewException, []gin.HandlerFunc{authMiddleware.Authenticate(), authMiddleware.RequireCenterAdmin()}},
		{http.MethodGet, "/api/v1/admin/rules/:ruleId/exceptions", s.action.scheduling.GetExceptionsByRule, []gin.HandlerFunc{authMiddleware.Authenticate(), authMiddleware.RequireCenterAdmin()}},
		{http.MethodGet, "/api/v1/admin/exceptions", s.action.scheduling.GetExceptionsByDateRange, []gin.HandlerFunc{authMiddleware.Authenticate(), authMiddleware.RequireCenterAdmin()}},
		{http.MethodGet, "/api/v1/admin/exceptions/pending", s.action.scheduling.GetPendingExceptions, []gin.HandlerFunc{authMiddleware.Authenticate(), authMiddleware.RequireCenterAdmin()}},
		{http.MethodGet, "/api/v1/admin/exceptions/all", s.action.scheduling.GetAllExceptions, []gin.HandlerFunc{authMiddleware.Authenticate(), authMiddleware.RequireCenterAdmin()}},
		{http.MethodPost, "/api/v1/admin/expand-rules", s.action.scheduling.ExpandRules, []gin.HandlerFunc{authMiddleware.Authenticate(), authMiddleware.RequireCenterAdmin()}},
		{http.MethodPost, "/api/v1/admin/detect-phase-transitions", s.action.scheduling.DetectPhaseTransitions, []gin.HandlerFunc{authMiddleware.Authenticate(), authMiddleware.RequireCenterAdmin()}},
		{http.MethodPost, "/api/v1/admin/scheduling/check-rule-lock", s.action.scheduling.CheckRuleLockStatus, []gin.HandlerFunc{authMiddleware.Authenticate(), authMiddleware.RequireCenterAdmin()}},

		// Admin - Rooms
		{http.MethodGet, "/api/v1/admin/rooms", s.action.adminRoom.GetRooms, []gin.HandlerFunc{authMiddleware.Authenticate(), authMiddleware.RequireCenterAdmin()}},
		{http.MethodPost, "/api/v1/admin/rooms", s.action.adminRoom.CreateRoom, []gin.HandlerFunc{authMiddleware.Authenticate(), authMiddleware.RequireCenterAdmin()}},
		{http.MethodPut, "/api/v1/admin/rooms/:room_id", s.action.adminRoom.UpdateRoom, []gin.HandlerFunc{authMiddleware.Authenticate(), authMiddleware.RequireCenterAdmin()}},
		{http.MethodGet, "/api/v1/admin/rooms/active", s.action.adminRoom.GetActiveRooms, []gin.HandlerFunc{authMiddleware.Authenticate(), authMiddleware.RequireCenterAdmin()}},
		{http.MethodPatch, "/api/v1/admin/rooms/:room_id/toggle-active", s.action.adminRoom.ToggleRoomActive, []gin.HandlerFunc{authMiddleware.Authenticate(), authMiddleware.RequireCenterAdmin()}},
		// Admin - Resources (非 Room/Course 路由)
		{http.MethodGet, "/api/v1/admin/courses", s.action.adminCourse.GetCourses, []gin.HandlerFunc{authMiddleware.Authenticate(), authMiddleware.RequireCenterAdmin()}},
		{http.MethodPost, "/api/v1/admin/courses", s.action.adminCourse.CreateCourse, []gin.HandlerFunc{authMiddleware.Authenticate(), authMiddleware.RequireCenterAdmin()}},
		{http.MethodPut, "/api/v1/admin/courses/:course_id", s.action.adminCourse.UpdateCourse, []gin.HandlerFunc{authMiddleware.Authenticate(), authMiddleware.RequireCenterAdmin()}},
		{http.MethodDelete, "/api/v1/admin/courses/:course_id", s.action.adminCourse.DeleteCourse, []gin.HandlerFunc{authMiddleware.Authenticate(), authMiddleware.RequireCenterAdmin()}},
		{http.MethodGet, "/api/v1/admin/courses/active", s.action.adminCourse.GetActiveCourses, []gin.HandlerFunc{authMiddleware.Authenticate(), authMiddleware.RequireCenterAdmin()}},
		{http.MethodPatch, "/api/v1/admin/courses/:course_id/toggle-active", s.action.adminCourse.ToggleCourseActive, []gin.HandlerFunc{authMiddleware.Authenticate(), authMiddleware.RequireCenterAdmin()}},
		{http.MethodGet, "/api/v1/admin/offerings/active", s.action.offering.GetActiveOfferings, []gin.HandlerFunc{authMiddleware.Authenticate(), authMiddleware.RequireCenterAdmin()}},
		{http.MethodPatch, "/api/v1/admin/offerings/:offering_id/toggle-active", s.action.offering.ToggleOfferingActive, []gin.HandlerFunc{authMiddleware.Authenticate(), authMiddleware.RequireCenterAdmin()}},

		// Admin - Invitations
		{http.MethodGet, "/api/v1/admin/centers/:id/invitations", s.action.teacherInvitation.GetInvitations, []gin.HandlerFunc{authMiddleware.Authenticate(), authMiddleware.RequireCenterAdmin()}},
		{http.MethodGet, "/api/v1/admin/centers/:id/invitations/stats", s.action.teacherInvitation.GetInvitationStats, []gin.HandlerFunc{authMiddleware.Authenticate(), authMiddleware.RequireCenterAdmin()}},
		// Admin - Invitation Links
		{http.MethodPost, "/api/v1/admin/centers/:id/invitations/generate-link", s.action.teacherInvitation.GenerateInvitationLink, []gin.HandlerFunc{authMiddleware.Authenticate(), authMiddleware.RequireCenterAdmin()}},
		{http.MethodGet, "/api/v1/admin/centers/:id/invitations/links", s.action.teacherInvitation.GetInvitationLinks, []gin.HandlerFunc{authMiddleware.Authenticate(), authMiddleware.RequireCenterAdmin()}},
		{http.MethodDelete, "/api/v1/admin/invitations/links/:id", s.action.teacherInvitation.RevokeInvitationLink, []gin.HandlerFunc{authMiddleware.Authenticate(), authMiddleware.RequireCenterAdmin()}},
		// Admin - General Invitation Links
		{http.MethodGet, "/api/v1/admin/centers/:id/invitations/general-link", s.action.teacherInvitation.GetGeneralInvitationLink, []gin.HandlerFunc{authMiddleware.Authenticate(), authMiddleware.RequireCenterAdmin()}},
		{http.MethodPost, "/api/v1/admin/centers/:id/invitations/general-link", s.action.teacherInvitation.GenerateGeneralInvitationLink, []gin.HandlerFunc{authMiddleware.Authenticate(), authMiddleware.RequireCenterAdmin()}},
		{http.MethodPost, "/api/v1/admin/centers/:id/invitations/general-link/toggle", s.action.teacherInvitation.ToggleGeneralInvitationStatus, []gin.HandlerFunc{authMiddleware.Authenticate(), authMiddleware.RequireCenterAdmin()}},
		{http.MethodPost, "/api/v1/admin/centers/:id/invitations/general-link/regenerate", s.action.teacherInvitation.RegenerateGeneralInvitationLink, []gin.HandlerFunc{authMiddleware.Authenticate(), authMiddleware.RequireCenterAdmin()}},
		// Public Invitation APIs (no auth required)
		{http.MethodGet, "/api/v1/invitations/:token", s.action.teacherInvitation.GetPublicInvitation, []gin.HandlerFunc{}},
		{http.MethodPost, "/api/v1/invitations/:token/accept", s.action.teacherInvitation.AcceptInvitationByLink, []gin.HandlerFunc{}},

		// Admin - Holidays
		{http.MethodGet, "/api/v1/admin/centers/:id/holidays", s.action.adminHoliday.GetHolidays, []gin.HandlerFunc{authMiddleware.Authenticate(), authMiddleware.RequireCenterAdmin()}},
		{http.MethodPost, "/api/v1/admin/centers/:id/holidays", s.action.adminHoliday.CreateHoliday, []gin.HandlerFunc{authMiddleware.Authenticate(), authMiddleware.RequireCenterAdmin()}},
		{http.MethodDelete, "/api/v1/admin/centers/:id/holidays/:holiday_id", s.action.adminHoliday.DeleteHoliday, []gin.HandlerFunc{authMiddleware.Authenticate(), authMiddleware.RequireCenterAdmin()}},
		{http.MethodPost, "/api/v1/admin/centers/:id/holidays/bulk", s.action.adminHoliday.BulkCreateHolidays, []gin.HandlerFunc{authMiddleware.Authenticate(), authMiddleware.RequireCenterAdmin()}},

		// Admin - Teacher Notes (評分與備註)
		{http.MethodGet, "/api/v1/admin/teachers/:teacher_id/note", s.action.adminTeacher.GetTeacherNote, []gin.HandlerFunc{authMiddleware.Authenticate(), authMiddleware.RequireCenterAdmin()}},
		{http.MethodPut, "/api/v1/admin/teachers/:teacher_id/note", s.action.adminTeacher.UpsertTeacherNote, []gin.HandlerFunc{authMiddleware.Authenticate(), authMiddleware.RequireCenterAdmin()}},
		{http.MethodDelete, "/api/v1/admin/teachers/:teacher_id/note", s.action.adminTeacher.DeleteTeacherNote, []gin.HandlerFunc{authMiddleware.Authenticate(), authMiddleware.RequireCenterAdmin()}},

		// Smart Matching
		{http.MethodPost, "/api/v1/admin/smart-matching/matches", s.action.smartMatching.FindMatches, []gin.HandlerFunc{authMiddleware.Authenticate(), authMiddleware.RequireCenterAdmin()}},
		{http.MethodGet, "/api/v1/admin/smart-matching/talent/search", s.action.smartMatching.SearchTalent, []gin.HandlerFunc{authMiddleware.Authenticate(), authMiddleware.RequireCenterAdmin()}},
		{http.MethodGet, "/api/v1/admin/smart-matching/talent/stats", s.action.smartMatching.GetTalentStats, []gin.HandlerFunc{authMiddleware.Authenticate(), authMiddleware.RequireCenterAdmin()}},
		{http.MethodPost, "/api/v1/admin/smart-matching/talent/invite", s.action.smartMatching.InviteTalent, []gin.HandlerFunc{authMiddleware.Authenticate(), authMiddleware.RequireCenterAdmin()}},
		{http.MethodGet, "/api/v1/admin/smart-matching/suggestions", s.action.smartMatching.GetSearchSuggestions, []gin.HandlerFunc{authMiddleware.Authenticate(), authMiddleware.RequireCenterAdmin()}},
		{http.MethodPost, "/api/v1/admin/smart-matching/alternatives", s.action.smartMatching.GetAlternativeSlots, []gin.HandlerFunc{authMiddleware.Authenticate(), authMiddleware.RequireCenterAdmin()}},
		{http.MethodGet, "/api/v1/admin/teachers/:teacher_id/sessions", s.action.smartMatching.GetTeacherSessions, []gin.HandlerFunc{authMiddleware.Authenticate(), authMiddleware.RequireCenterAdmin()}},

		// Notification
		{http.MethodGet, "/api/v1/notifications", s.action.notification.ListNotifications, []gin.HandlerFunc{authMiddleware.Authenticate()}},
		{http.MethodGet, "/api/v1/notifications/unread-count", s.action.notification.GetUnreadCount, []gin.HandlerFunc{authMiddleware.Authenticate()}},
		{http.MethodPost, "/api/v1/notifications/:id/read", s.action.notification.MarkAsRead, []gin.HandlerFunc{authMiddleware.Authenticate()}},
		{http.MethodPost, "/api/v1/notifications/read-all", s.action.notification.MarkAllAsRead, []gin.HandlerFunc{authMiddleware.Authenticate()}},
		{http.MethodPost, "/api/v1/notifications/token", s.action.notification.SetNotifyToken, []gin.HandlerFunc{authMiddleware.Authenticate()}},
		{http.MethodPost, "/api/v1/notifications/test", s.action.notification.SendTestNotification, []gin.HandlerFunc{authMiddleware.Authenticate()}},
		// Notification Queue Stats (Admin only)
		{http.MethodGet, "/api/v1/admin/notifications/queue-stats", s.action.notification.GetQueueStats, []gin.HandlerFunc{authMiddleware.Authenticate(), authMiddleware.RequireCenterAdmin()}},
		// Admin Notification - Broadcast (Admin only)
		{http.MethodPost, "/api/v1/admin/notifications/broadcast", s.action.adminNotification.Broadcast, []gin.HandlerFunc{authMiddleware.Authenticate(), authMiddleware.RequireCenterAdmin()}},
		// Health Check
		{http.MethodGet, "/api/v1/admin/health/redis", s.action.notification.CheckRedisHealth, []gin.HandlerFunc{authMiddleware.Authenticate(), authMiddleware.RequireCenterAdmin()}},

		// Export
		{http.MethodPost, "/api/v1/admin/export/schedule/csv", s.action.export.ExportScheduleCSV, []gin.HandlerFunc{authMiddleware.Authenticate(), authMiddleware.RequireCenterAdmin()}},
		{http.MethodPost, "/api/v1/admin/export/schedule/pdf", s.action.export.ExportSchedulePDF, []gin.HandlerFunc{authMiddleware.Authenticate(), authMiddleware.RequireCenterAdmin()}},
		{http.MethodGet, "/api/v1/admin/centers/:id/export/teachers/csv", s.action.export.ExportTeachersCSV, []gin.HandlerFunc{authMiddleware.Authenticate(), authMiddleware.RequireCenterAdmin()}},
		{http.MethodGet, "/api/v1/admin/centers/:id/export/exceptions/csv", s.action.export.ExportExceptionsCSV, []gin.HandlerFunc{authMiddleware.Authenticate(), authMiddleware.RequireCenterAdmin()}},

		// Teacher - ICS Calendar Export
		{http.MethodGet, "/api/v1/teacher/me/schedule.ics", s.action.export.ExportScheduleToICS, []gin.HandlerFunc{authMiddleware.Authenticate()}},
		{http.MethodPost, "/api/v1/teacher/me/schedule/subscription", s.action.export.CreateCalendarSubscription, []gin.HandlerFunc{authMiddleware.Authenticate()}},
		{http.MethodDelete, "/api/v1/teacher/me/schedule/subscription", s.action.export.UnsubscribeCalendar, []gin.HandlerFunc{authMiddleware.Authenticate()}},

		// Teacher - Image Export
		{http.MethodGet, "/api/v1/teacher/me/schedule/image", s.action.export.ExportScheduleToImage, []gin.HandlerFunc{authMiddleware.Authenticate()}},
		{http.MethodPost, "/api/v1/teacher/me/backgrounds", s.action.export.UploadBackgroundImage, []gin.HandlerFunc{authMiddleware.Authenticate()}},
		{http.MethodGet, "/api/v1/teacher/me/backgrounds", s.action.export.GetBackgroundImages, []gin.HandlerFunc{authMiddleware.Authenticate()}},
		{http.MethodDelete, "/api/v1/teacher/me/backgrounds", s.action.export.DeleteBackgroundImage, []gin.HandlerFunc{authMiddleware.Authenticate()}},

		// Public - Calendar Subscription (no auth required)
		{http.MethodGet, "/api/v1/calendar/subscribe/:token.ics", s.action.export.SubscribeToCalendar, []gin.HandlerFunc{}},

		// LINE Bot Webhook (不需要認證)
		{http.MethodPost, "/api/v1/line/webhook", s.action.lineBot.HandleWebhook, []gin.HandlerFunc{}},
		{http.MethodGet, "/api/v1/line/health", s.action.lineBot.HealthCheck, []gin.HandlerFunc{}},
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

	// Response Sanitizer - 確保 API 回應是乾淨的 JSON
	responseSanitizer := middleware.NewResponseSanitizer()
	r.Use(responseSanitizer.Sanitize())

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

	authService := services.NewAuthService(s.app)
	s.action.auth = controllers.NewAuthController(s.app, authService)
	s.action.teacher = controllers.NewTeacherController(s.app)
	s.action.adminTeacher = controllers.NewAdminTeacherController(s.app)
	s.action.adminCenter = controllers.NewAdminCenterController(s.app)
	s.action.adminRoom = controllers.NewAdminRoomController(s.app)
	s.action.adminCourse = controllers.NewAdminCourseController(s.app)
	s.action.adminHoliday = controllers.NewAdminHolidayController(s.app)
	s.action.teacherProfile = controllers.NewTeacherProfileController(s.app)
	s.action.teacherSchedule = controllers.NewTeacherScheduleController(s.app)
	s.action.teacherSession = controllers.NewTeacherSessionController(s.app)
	s.action.teacherEvent = controllers.NewTeacherEventController(s.app)
	s.action.teacherException = controllers.NewTeacherExceptionController(s.app)
	s.action.teacherInvitation = controllers.NewTeacherInvitationController(s.app)
	s.action.geo = controllers.NewGeoController(s.app)
	s.action.adminResource = controllers.NewAdminResourceController(s.app)
	s.action.offering = controllers.NewOfferingController(s.app)
	s.action.timetableTemplate = controllers.NewTimetableTemplateController(s.app)
	s.action.adminUser = controllers.NewAdminUserController(s.app)
	s.action.scheduling = controllers.NewSchedulingController(s.app)
	s.action.smartMatching = controllers.NewSmartMatchingController(s.app)
	s.action.notification = controllers.NewNotificationController(s.app)
	s.action.adminNotification = controllers.NewAdminNotificationController(s.app)
	s.action.export = controllers.NewExportController(s.app)
	s.action.lineBot = controllers.NewLineBotController(s.app)
	s.action.r2Test = controllers.NewR2TestController(s.app)
}
