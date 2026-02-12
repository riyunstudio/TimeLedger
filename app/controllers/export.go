package controllers

import (
	"fmt"
	"net/http"
	"strings"
	"time"
	"timeLedger/app"
	"timeLedger/app/repositories"
	"timeLedger/app/requests"
	"timeLedger/app/services"
	"timeLedger/global"

	"github.com/gin-gonic/gin"
)

type ExportController struct {
	BaseController
	app       *app.App
	exportSvc services.ExportService
	icsSvc    *services.ICSCalendarService
	imageSvc  *services.ImageService
}

func NewExportController(app *app.App) *ExportController {
	return &ExportController{
		app:       app,
		exportSvc: services.NewExportService(app),
		icsSvc:    services.NewICSCalendarService(app),
		imageSvc:  services.NewImageService(app),
	}
}

func (ctl *ExportController) ExportScheduleCSV(ctx *gin.Context) {
	var req requests.ExportScheduleRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, global.ApiResponse{
			Code:    global.BAD_REQUEST,
			Message: "Invalid request body: " + err.Error(),
		})
		return
	}

	if req.CenterID == 0 {
		ctx.JSON(http.StatusBadRequest, global.ApiResponse{
			Code:    global.BAD_REQUEST,
			Message: "Center ID required",
		})
		return
	}

	startDateStr := req.StartDate
	endDateStr := req.EndDate

	if startDateStr == "" || endDateStr == "" {
		ctx.JSON(http.StatusBadRequest, global.ApiResponse{
			Code:    global.BAD_REQUEST,
			Message: "Start date and end date required",
		})
		return
	}

	startDateStr = strings.TrimSuffix(startDateStr, "Z")
	endDateStr = strings.TrimSuffix(endDateStr, "Z")

	var startDate, endDate time.Time
	var err error

	startDate, err = time.Parse("2006-01-02T15:04:05", startDateStr)
	if err != nil {
		startDate, err = time.Parse("2006-01-02", startDateStr)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, global.ApiResponse{
				Code:    global.BAD_REQUEST,
				Message: "Invalid start date format: " + err.Error(),
			})
			return
		}
	}

	endDate, err = time.Parse("2006-01-02T15:04:05", endDateStr)
	if err != nil {
		endDate, err = time.Parse("2006-01-02", endDateStr)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, global.ApiResponse{
				Code:    global.BAD_REQUEST,
				Message: "Invalid end date format: " + err.Error(),
			})
			return
		}
	}

	data, err := ctl.exportSvc.ExportScheduleToCSV(ctx, req.CenterID, startDate, endDate)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, global.ApiResponse{
			Code:    500,
			Message: err.Error(),
		})
		return
	}

	filename := fmt.Sprintf("schedule_%s_%s.csv", startDate.Format("20060102"), endDate.Format("20060102"))
	ctx.Header("Content-Type", "text/csv")
	ctx.Header("Content-Disposition", "attachment; filename="+filename)
	ctx.Data(http.StatusOK, "text/csv", data)
}

func (ctl *ExportController) ExportSchedulePDF(ctx *gin.Context) {
	var req requests.ExportScheduleRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, global.ApiResponse{
			Code:    global.BAD_REQUEST,
			Message: "Invalid request body: " + err.Error(),
		})
		return
	}

	if req.CenterID == 0 {
		ctx.JSON(http.StatusBadRequest, global.ApiResponse{
			Code:    global.BAD_REQUEST,
			Message: "Center ID required",
		})
		return
	}

	startDateStr := req.StartDate
	endDateStr := req.EndDate

	if startDateStr == "" || endDateStr == "" {
		ctx.JSON(http.StatusBadRequest, global.ApiResponse{
			Code:    global.BAD_REQUEST,
			Message: "Start date and end date required",
		})
		return
	}

	startDateStr = strings.TrimSuffix(startDateStr, "Z")
	endDateStr = strings.TrimSuffix(endDateStr, "Z")

	startDate, err := time.Parse("2006-01-02T15:04:05", startDateStr)
	if err != nil {
		startDate, err = time.Parse("2006-01-02", startDateStr)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, global.ApiResponse{
				Code:    global.BAD_REQUEST,
				Message: "Invalid start date format: " + err.Error(),
			})
			return
		}
	}

	endDate, err := time.Parse("2006-01-02T15:04:05", endDateStr)
	if err != nil {
		endDate, err = time.Parse("2006-01-02", endDateStr)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, global.ApiResponse{
				Code:    global.BAD_REQUEST,
				Message: "Invalid end date format: " + err.Error(),
			})
			return
		}
	}

	data, err := ctl.exportSvc.ExportScheduleToPDF(ctx, req.CenterID, startDate, endDate)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, global.ApiResponse{
			Code:    500,
			Message: err.Error(),
		})
		return
	}

	filename := fmt.Sprintf("schedule_%s_%s.txt", startDate.Format("20060102"), endDate.Format("20060102"))
	ctx.Header("Content-Type", "text/plain")
	ctx.Header("Content-Disposition", "attachment; filename="+filename)
	ctx.Data(http.StatusOK, "text/plain", data)
}

func (ctl *ExportController) ExportTeachersCSV(ctx *gin.Context) {
	centerID := ctx.Param("id")
	var centerIDInt uint
	if _, err := fmt.Sscanf(centerID, "%d", &centerIDInt); err != nil {
		ctx.JSON(http.StatusBadRequest, global.ApiResponse{
			Code:    global.BAD_REQUEST,
			Message: "Invalid center ID",
		})
		return
	}

	data, err := ctl.exportSvc.ExportTeachersToCSV(ctx, centerIDInt)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, global.ApiResponse{
			Code:    500,
			Message: err.Error(),
		})
		return
	}

	filename := fmt.Sprintf("teachers_%d.csv", centerIDInt)
	ctx.Header("Content-Type", "text/csv")
	ctx.Header("Content-Disposition", "attachment; filename="+filename)
	ctx.Data(http.StatusOK, "text/csv", data)
}

func (ctl *ExportController) ExportExceptionsCSV(ctx *gin.Context) {
	centerID := ctx.Param("id")
	var centerIDInt uint
	if _, err := fmt.Sscanf(centerID, "%d", &centerIDInt); err != nil {
		ctx.JSON(http.StatusBadRequest, global.ApiResponse{
			Code:    global.BAD_REQUEST,
			Message: "Invalid center ID",
		})
		return
	}

	startDateStr := ctx.Query("start_date")
	endDateStr := ctx.Query("end_date")

	startDate, err := time.Parse("2006-01-02", startDateStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, global.ApiResponse{
			Code:    global.BAD_REQUEST,
			Message: "Invalid start date format",
		})
		return
	}

	endDate, err := time.Parse("2006-01-02", endDateStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, global.ApiResponse{
			Code:    global.BAD_REQUEST,
			Message: "Invalid end date format",
		})
		return
	}

	data, err := ctl.exportSvc.ExportExceptionsToCSV(ctx, centerIDInt, startDate, endDate)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, global.ApiResponse{
			Code:    500,
			Message: err.Error(),
		})
		return
	}

	filename := fmt.Sprintf("exceptions_%s_%s.csv", startDate.Format("20060102"), endDate.Format("20060102"))
	ctx.Header("Content-Type", "text/csv")
	ctx.Header("Content-Disposition", "attachment; filename="+filename)
	ctx.Data(http.StatusOK, "text/csv", data)
}

// ==================== ICS Calendar Export ====================

// ExportScheduleToICS 匯出課表為 ICS 格式
// @Summary 匯出課表為 iCalendar 格式
// @Tags Teacher - Export
// @Accept json
// @Produce text/calendar
// @Security BearerAuth
// @Param start_date query string true "開始日期 (YYYY-MM-DD)"
// @Param end_date query string true "結束日期 (YYYY-MM-DD)"
// @Success 200 {file} file "ICS 檔案"
// @Router /api/v1/teacher/me/schedule.ics [get]
func (ctl *ExportController) ExportScheduleToICS(ctx *gin.Context) {
	teacherID := ctx.GetUint(global.UserIDKey)
	if teacherID == 0 {
		ctx.JSON(http.StatusUnauthorized, global.ApiResponse{
			Code:    global.UNAUTHORIZED,
			Message: "Teacher ID required",
		})
		return
	}

	startDateStr := ctx.Query("start_date")
	endDateStr := ctx.Query("end_date")

	if startDateStr == "" || endDateStr == "" {
		ctx.JSON(http.StatusBadRequest, global.ApiResponse{
			Code:    global.BAD_REQUEST,
			Message: "start_date and end_date are required",
		})
		return
	}

	startDate, err := time.Parse("2006-01-02", startDateStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, global.ApiResponse{
			Code:    global.BAD_REQUEST,
			Message: "Invalid start date format",
		})
		return
	}

	endDate, err := time.Parse("2006-01-02", endDateStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, global.ApiResponse{
			Code:    global.BAD_REQUEST,
			Message: "Invalid end date format",
		})
		return
	}

	// 取得老師所屬的中心
	membershipRepo := repositories.NewCenterMembershipRepository(ctl.app)
	memberships, _ := membershipRepo.GetActiveByTeacherID(ctx, teacherID)
	if len(memberships) == 0 {
		ctx.JSON(http.StatusNotFound, global.ApiResponse{
			Code:    404,
			Message: "No center membership found",
		})
		return
	}

	centerID := memberships[0].CenterID

	// 產生 ICS 資料
	data, err := ctl.icsSvc.ExportTeacherScheduleToICS(ctx, teacherID, centerID, startDate, endDate)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, global.ApiResponse{
			Code:    500,
			Message: err.Error(),
		})
		return
	}

	filename := fmt.Sprintf("schedule_%s_%s.ics", startDate.Format("20060102"), endDate.Format("20060102"))
	ctx.Header("Content-Type", "text/calendar; charset=utf-8")
	ctx.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s", filename))
	ctx.Data(http.StatusOK, "text/calendar; charset=utf-8", data)
}

// CreateCalendarSubscription 建立課表訂閱
// @Summary 建立課表訂閱連結
// @Tags Teacher - Export
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} global.ApiResponse{data=services.SubscriptionInfo}
// @Router /api/v1/teacher/me/schedule/subscription [post]
func (ctl *ExportController) CreateCalendarSubscription(ctx *gin.Context) {
	teacherID := ctx.GetUint(global.UserIDKey)
	if teacherID == 0 {
		ctx.JSON(http.StatusUnauthorized, global.ApiResponse{
			Code:    global.UNAUTHORIZED,
			Message: "Teacher ID required",
		})
		return
	}

	// 取得老師所屬的中心
	membershipRepo := repositories.NewCenterMembershipRepository(ctl.app)
	memberships, err := membershipRepo.GetActiveByTeacherID(ctx, teacherID)
	if err != nil || len(memberships) == 0 {
		ctx.JSON(http.StatusNotFound, global.ApiResponse{
			Code:    404,
			Message: "No center membership found",
		})
		return
	}

	centerID := memberships[0].CenterID

	// 建立訂閱
	subscription, errInfo, err := ctl.icsSvc.CreateCalendarSubscription(ctx, teacherID, centerID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, global.ApiResponse{
			Code:    errInfo.Code,
			Message: err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, global.ApiResponse{
		Code:    0,
		Message: "Subscription created",
		Datas:   subscription,
	})
}

// UnsubscribeCalendar 取消課表訂閱
// @Summary 取消課表訂閱
// @Tags Teacher - Export
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param token query string true "訂閱 Token"
// @Success 200 {object} global.ApiResponse
// @Router /api/v1/teacher/me/schedule/subscription [delete]
func (ctl *ExportController) UnsubscribeCalendar(ctx *gin.Context) {
	token := ctx.Query("token")
	if token == "" {
		ctx.JSON(http.StatusBadRequest, global.ApiResponse{
			Code:    global.BAD_REQUEST,
			Message: "Token required",
		})
		return
	}

	if err := ctl.icsSvc.Unsubscribe(ctx, token); err != nil {
		ctx.JSON(http.StatusInternalServerError, global.ApiResponse{
			Code:    500,
			Message: err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, global.ApiResponse{
		Code:    0,
		Message: "Unsubscribed successfully",
	})
}

// SubscribeToCalendar 公開訂閱課表（無需認證）
// @Summary 透過 Token 訂閱課表
// @Tags Export
// @Produce text/calendar
// @Param token path string true "訂閱 Token"
// @Success 200 {file} file "ICS 檔案"
// @Router /api/v1/calendar/subscribe/{token}.ics [get]
func (ctl *ExportController) SubscribeToCalendar(ctx *gin.Context) {
	token := ctx.Param("token")
	if token == "" {
		ctx.JSON(http.StatusBadRequest, global.ApiResponse{
			Code:    global.BAD_REQUEST,
			Message: "Token required",
		})
		return
	}

	// 驗證 token 並取得對應的課表資料
	teacherID, centerID, valid, err := ctl.icsSvc.ValidateSubscriptionToken(ctx, token)
	if err != nil || !valid {
		ctx.JSON(http.StatusNotFound, global.ApiResponse{
			Code:    404,
			Message: "Invalid or expired subscription token",
		})
		return
	}

	// 取得訂閱配置（從資料庫或使用預設值）
	startDate := time.Now()
	endDate := startDate.AddDate(0, 1, 0) // 預設訂閱一個月

	// 產生 ICS 資料
	data, err := ctl.icsSvc.ExportTeacherScheduleToICS(ctx, teacherID, centerID, startDate, endDate)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, global.ApiResponse{
			Code:    500,
			Message: err.Error(),
		})
		return
	}

	filename := fmt.Sprintf("schedule_%s.ics", token[:8])
	ctx.Header("Content-Type", "text/calendar; charset=utf-8")
	ctx.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s", filename))
	ctx.Data(http.StatusOK, "text/calendar; charset=utf-8", data)
}

// ==================== Image Export ====================

// ExportScheduleToImage 匯出課表為圖片
// @Summary 匯出課表為圖片
// @Tags Teacher - Export
// @Accept json
// @Produce image/jpeg
// @Security BearerAuth
// @Param start_date query string true "開始日期 (YYYY-MM-DD)"
// @Param end_date query string true "結束日期 (YYYY-MM-DD)"
// @Param background query string false "自訂背景圖路徑"
// @Success 200 {file} jpeg "JPEG 圖片"
// @Router /api/v1/teacher/me/schedule/image [get]
func (ctl *ExportController) ExportScheduleToImage(ctx *gin.Context) {
	teacherID := ctx.GetUint(global.UserIDKey)
	if teacherID == 0 {
		ctx.JSON(http.StatusUnauthorized, global.ApiResponse{
			Code:    global.UNAUTHORIZED,
			Message: "Teacher ID required",
		})
		return
	}

	startDateStr := ctx.Query("start_date")
	endDateStr := ctx.Query("end_date")
	backgroundPath := ctx.Query("background")

	if startDateStr == "" || endDateStr == "" {
		ctx.JSON(http.StatusBadRequest, global.ApiResponse{
			Code:    global.BAD_REQUEST,
			Message: "start_date and end_date are required",
		})
		return
	}

	startDate, err := time.Parse("2006-01-02", startDateStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, global.ApiResponse{
			Code:    global.BAD_REQUEST,
			Message: "Invalid start date format",
		})
		return
	}

	endDate, err := time.Parse("2006-01-02", endDateStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, global.ApiResponse{
			Code:    global.BAD_REQUEST,
			Message: "Invalid end date format",
		})
		return
	}

	// 取得老師所屬的中心
	membershipRepo := repositories.NewCenterMembershipRepository(ctl.app)
	memberships, _ := membershipRepo.GetActiveByTeacherID(ctx, teacherID)
	if len(memberships) == 0 {
		ctx.JSON(http.StatusNotFound, global.ApiResponse{
			Code:    404,
			Message: "No center membership found",
		})
		return
	}

	centerID := memberships[0].CenterID

	// 產生圖片
	data, errInfo, err := ctl.imageSvc.ExportTeacherScheduleToImage(ctx, teacherID, centerID, startDate, endDate, backgroundPath)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, global.ApiResponse{
			Code:    errInfo.Code,
			Message: err.Error(),
		})
		return
	}

	filename := fmt.Sprintf("schedule_%s_%s.jpg", startDate.Format("20060102"), endDate.Format("20060102"))
	ctx.Header("Content-Type", "image/jpeg")
	ctx.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s", filename))
	ctx.Data(http.StatusOK, "image/jpeg", data)
}

// UploadBackgroundImage 上傳自訂背景圖
// @Summary 上傳自訂背景圖
// @Tags Teacher - Export
// @Accept multipart/form-data
// @Produce json
// @Security BearerAuth
// @Param file formData file true "背景圖片"
// @Success 200 {object} global.ApiResponse{data=TeacherBackgroundResponse}
// @Router /api/v1/teacher/me/backgrounds [post]
func (ctl *ExportController) UploadBackgroundImage(ctx *gin.Context) {
	teacherID := ctx.GetUint(global.UserIDKey)
	if teacherID == 0 {
		ctx.JSON(http.StatusUnauthorized, global.ApiResponse{
			Code:    global.UNAUTHORIZED,
			Message: "Teacher ID required",
		})
		return
	}

	file, err := ctx.FormFile("file")
	if err != nil {
		ctx.JSON(http.StatusBadRequest, global.ApiResponse{
			Code:    global.BAD_REQUEST,
			Message: "No file uploaded: " + err.Error(),
		})
		return
	}

	// 檢查檔案類型
	allowedTypes := map[string]bool{"image/jpeg": true, "image/png": true}
	if !allowedTypes[file.Header.Get("Content-Type")] {
		ctx.JSON(http.StatusBadRequest, global.ApiResponse{
			Code:    global.BAD_REQUEST,
			Message: "Invalid file type. Allowed: jpeg, png",
		})
		return
	}

	// 開啟檔案
	src, err := file.Open()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, global.ApiResponse{
			Code:    500,
			Message: "Failed to open file: " + err.Error(),
		})
		return
	}
	defer src.Close()

	// 上傳背景圖
	background, errInfo, err := ctl.imageSvc.UploadBackgroundImage(ctx, teacherID, src, file.Filename)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, global.ApiResponse{
			Code:    errInfo.Code,
			Message: err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, global.ApiResponse{
		Code:    0,
		Message: "Background image uploaded",
		Datas: TeacherBackgroundResponse{
			ID:        background.ID,
			TeacherID: background.TeacherID,
			FileURL:   background.FileURL,
			CreatedAt: background.CreatedAt,
		},
	})
}

// GetBackgroundImages 取得自訂背景圖列表
// @Summary 取得自訂背景圖列表
// @Tags Teacher - Export
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} global.ApiResponse{data=[]TeacherBackgroundResponse}
// @Router /api/v1/teacher/me/backgrounds [get]
func (ctl *ExportController) GetBackgroundImages(ctx *gin.Context) {
	teacherID := ctx.GetUint(global.UserIDKey)
	if teacherID == 0 {
		ctx.JSON(http.StatusUnauthorized, global.ApiResponse{
			Code:    global.UNAUTHORIZED,
			Message: "Teacher ID required",
		})
		return
	}

	backgrounds, err := ctl.imageSvc.GetTeacherBackgroundImages(ctx, teacherID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, global.ApiResponse{
			Code:    500,
			Message: err.Error(),
		})
		return
	}

	// 轉換為 response 格式
	responses := make([]TeacherBackgroundResponse, len(backgrounds))
	for i, bg := range backgrounds {
		responses[i] = TeacherBackgroundResponse{
			ID:        bg.ID,
			TeacherID: bg.TeacherID,
			FileURL:   bg.FileURL,
			CreatedAt: bg.CreatedAt,
		}
	}

	ctx.JSON(http.StatusOK, global.ApiResponse{
		Code:    0,
		Message: "Success",
		Datas:   responses,
	})
}

// DeleteBackgroundImage 刪除自訂背景圖
// @Summary 刪除自訂背景圖
// @Tags Teacher - Export
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "背景圖 ID"
// @Success 200 {object} global.ApiResponse
// @Router /api/v1/teacher/me/backgrounds/{id} [delete]
func (ctl *ExportController) DeleteBackgroundImage(ctx *gin.Context) {
	teacherID := ctx.GetUint(global.UserIDKey)
	if teacherID == 0 {
		ctx.JSON(http.StatusUnauthorized, global.ApiResponse{
			Code:    global.UNAUTHORIZED,
			Message: "Teacher ID required",
		})
		return
	}

	// 解析背景圖 ID
	var backgroundID uint
	if _, err := fmt.Sscanf(ctx.Param("id"), "%d", &backgroundID); err != nil || backgroundID == 0 {
		ctx.JSON(http.StatusBadRequest, global.ApiResponse{
			Code:    global.BAD_REQUEST,
			Message: "Invalid background ID",
		})
		return
	}

	if err := ctl.imageSvc.DeleteBackgroundImage(ctx, teacherID, backgroundID); err != nil {
		ctx.JSON(http.StatusInternalServerError, global.ApiResponse{
			Code:    500,
			Message: err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, global.ApiResponse{
		Code:    0,
		Message: "Background image deleted",
	})
}

// TeacherBackgroundResponse 背景圖回應結構
type TeacherBackgroundResponse struct {
	ID        uint      `json:"id"`
	TeacherID uint      `json:"teacher_id"`
	FileURL   string    `json:"file_url"`
	CreatedAt time.Time `json:"created_at"`
}

// UploadBackgroundResponse 上傳背景圖回應（已廢用，請使用 TeacherBackgroundResponse）
type UploadBackgroundResponse struct {
	ImagePath string `json:"image_path"`
	FileName  string `json:"file_name"`
}
