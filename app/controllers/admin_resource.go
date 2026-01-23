package controllers

import (
	"fmt"
	"net/http"
	"time"
	"timeLedger/app"
	"timeLedger/app/models"
	"timeLedger/app/repositories"
	"timeLedger/global"
	"timeLedger/global/errInfos"

	"github.com/gin-gonic/gin"
)

type AdminResourceController struct {
	BaseController
	app                *app.App
	centerRepository   *repositories.CenterRepository
	courseRepository   *repositories.CourseRepository
	roomRepository     *repositories.RoomRepository
	offeringRepository *repositories.OfferingRepository
	holidayRepository  *repositories.CenterHolidayRepository
	invitationRepo     *repositories.CenterInvitationRepository
	auditLogRepo       *repositories.AuditLogRepository
}

func NewAdminResourceController(app *app.App) *AdminResourceController {
	return &AdminResourceController{
		app:                app,
		centerRepository:   repositories.NewCenterRepository(app),
		courseRepository:   repositories.NewCourseRepository(app),
		roomRepository:     repositories.NewRoomRepository(app),
		offeringRepository: repositories.NewOfferingRepository(app),
		holidayRepository:  repositories.NewCenterHolidayRepository(app),
		invitationRepo:     repositories.NewCenterInvitationRepository(app),
		auditLogRepo:       repositories.NewAuditLogRepository(app),
	}
}

// GetCenters 取得中心列表
// @Summary 取得中心列表
// @Tags Admin
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} global.ApiResponse{data=[]models.Center}
// @Router /api/v1/admin/centers [get]
func (ctl *AdminResourceController) GetCenters(ctx *gin.Context) {
	centers, err := ctl.centerRepository.List(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, global.ApiResponse{
			Code:    500,
			Message: "Failed to get centers",
		})
		return
	}

	ctx.JSON(http.StatusOK, global.ApiResponse{
		Code:    0,
		Message: "Success",
		Datas:   centers,
	})
}

// CreateCenter 新增中心
// @Summary 新增中心
// @Tags Admin
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body CreateCenterRequest true "中心資訊"
// @Success 200 {object} global.ApiResponse{data=models.Center}
// @Router /api/v1/admin/centers [post]
func (ctl *AdminResourceController) CreateCenter(ctx *gin.Context) {
	var req CreateCenterRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, global.ApiResponse{
			Code:    global.BAD_REQUEST,
			Message: "Invalid request body",
		})
		return
	}

	center := models.Center{
		Name:      req.Name,
		PlanLevel: req.PlanLevel,
		Settings: models.CenterSettings{
			AllowPublicRegister: req.AllowPublicRegister,
			DefaultLanguage:     "zh-TW",
		},
		CreatedAt: time.Now(),
	}

	createdCenter, err := ctl.centerRepository.Create(ctx, center)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, global.ApiResponse{
			Code:    500,
			Message: "Failed to create center",
		})
		return
	}

	ctx.JSON(http.StatusOK, global.ApiResponse{
		Code:    0,
		Message: "Center created",
		Datas:   createdCenter,
	})
}

type CreateCenterRequest struct {
	Name                string `json:"name" binding:"required"`
	PlanLevel           string `json:"plan_level" binding:"required"`
	AllowPublicRegister bool   `json:"allow_public_register"`
}

// GetRooms 取得教室列表
// @Summary 取得教室列表
// @Tags Admin
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} global.ApiResponse{data=[]models.Room}
// @Router /api/v1/admin/rooms [get]
func (ctl *AdminResourceController) GetRooms(ctx *gin.Context) {
	centerID := ctx.GetUint(global.CenterIDKey)

	// If GetUint returns 0, try getting as interface{}
	if centerID == 0 {
		if val, exists := ctx.Get(global.CenterIDKey); exists {
			switch v := val.(type) {
			case uint:
				centerID = v
			case uint64:
				centerID = uint(v)
			case int:
				centerID = uint(v)
			case float64:
				centerID = uint(v)
			}
		}
	}

	if centerID == 0 {
		ctx.JSON(http.StatusBadRequest, global.ApiResponse{
			Code:    global.BAD_REQUEST,
			Message: "Center ID required",
		})
		return
	}

	rooms, err := ctl.roomRepository.ListByCenterID(ctx, centerID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, global.ApiResponse{
			Code:    500,
			Message: "Failed to get rooms",
		})
		return
	}

	ctx.JSON(http.StatusOK, global.ApiResponse{
		Code:    0,
		Message: "Success",
		Datas:   rooms,
	})
}

// CreateRoom 新增教室
// @Summary 新增教室
// @Tags Admin
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body CreateRoomRequest true "教室資訊"
// @Success 200 {object} global.ApiResponse{data=models.Room}
// @Router /api/v1/admin/rooms [post]
func (ctl *AdminResourceController) CreateRoom(ctx *gin.Context) {
	centerID := ctx.GetUint(global.CenterIDKey)

	if centerID == 0 {
		if val, exists := ctx.Get(global.CenterIDKey); exists {
			switch v := val.(type) {
			case uint:
				centerID = v
			case uint64:
				centerID = uint(v)
			case int:
				centerID = uint(v)
			case float64:
				centerID = uint(v)
			}
		}
	}

	if centerID == 0 {
		ctx.JSON(http.StatusBadRequest, global.ApiResponse{
			Code:    global.BAD_REQUEST,
			Message: "Center ID required",
		})
		return
	}

	var req CreateRoomRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, global.ApiResponse{
			Code:    global.BAD_REQUEST,
			Message: "Invalid request body",
		})
		return
	}

	room := models.Room{
		CenterID:  centerID,
		Name:      req.Name,
		Capacity:  req.Capacity,
		CreatedAt: time.Now(),
	}

	createdRoom, err := ctl.roomRepository.Create(ctx, room)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, global.ApiResponse{
			Code:    500,
			Message: "Failed to create room",
		})
		return
	}

	adminID := ctx.GetUint(global.UserIDKey)
	ctl.auditLogRepo.Create(ctx, models.AuditLog{
		CenterID:   centerID,
		ActorType:  "ADMIN",
		ActorID:    adminID,
		Action:     "CREATE_ROOM",
		TargetType: "Room",
		TargetID:   createdRoom.ID,
		Payload: models.AuditPayload{
			After: room,
		},
	})

	ctx.JSON(http.StatusOK, global.ApiResponse{
		Code:    0,
		Message: "Room created",
		Datas:   createdRoom,
	})
}

type CreateRoomRequest struct {
	Name     string `json:"name" binding:"required"`
	Capacity int    `json:"capacity" binding:"required"`
}

// GetCourses 取得課程列表
// @Summary 取得課程列表
// @Tags Admin
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} global.ApiResponse{data=[]models.Course}
// @Router /api/v1/admin/courses [get]
func (ctl *AdminResourceController) GetCourses(ctx *gin.Context) {
	centerID := ctx.GetUint(global.CenterIDKey)

	if centerID == 0 {
		if val, exists := ctx.Get(global.CenterIDKey); exists {
			switch v := val.(type) {
			case uint:
				centerID = v
			case uint64:
				centerID = uint(v)
			case int:
				centerID = uint(v)
			case float64:
				centerID = uint(v)
			}
		}
	}

	if centerID == 0 {
		ctx.JSON(http.StatusBadRequest, global.ApiResponse{
			Code:    global.BAD_REQUEST,
			Message: "Center ID required",
		})
		return
	}

	courses, err := ctl.courseRepository.ListByCenterID(ctx, centerID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, global.ApiResponse{
			Code:    500,
			Message: "Failed to get courses",
		})
		return
	}

	ctx.JSON(http.StatusOK, global.ApiResponse{
		Code:    0,
		Message: "Success",
		Datas:   courses,
	})
}

// CreateCourse 新增課程
// @Summary 新增課程
// @Tags Admin
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body CreateCourseRequest true "課程資訊"
// @Success 200 {object} global.ApiResponse{data=models.Course}
// @Router /api/v1/admin/courses [post]
func (ctl *AdminResourceController) CreateCourse(ctx *gin.Context) {
	centerID := ctx.GetUint(global.CenterIDKey)

	if centerID == 0 {
		if val, exists := ctx.Get(global.CenterIDKey); exists {
			switch v := val.(type) {
			case uint:
				centerID = v
			case uint64:
				centerID = uint(v)
			case int:
				centerID = uint(v)
			case float64:
				centerID = uint(v)
			}
		}
	}

	if centerID == 0 {
		ctx.JSON(http.StatusBadRequest, global.ApiResponse{
			Code:    global.BAD_REQUEST,
			Message: "Center ID required",
		})
		return
	}

	var req CreateCourseRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, global.ApiResponse{
			Code:    global.BAD_REQUEST,
			Message: "Invalid request body",
		})
		return
	}

	course := models.Course{
		CenterID:         centerID,
		Name:             req.Name,
		DefaultDuration:  req.Duration,
		ColorHex:         req.ColorHex,
		RoomBufferMin:    req.RoomBufferMin,
		TeacherBufferMin: req.TeacherBufferMin,
		CreatedAt:        time.Now(),
	}

	createdCourse, err := ctl.courseRepository.Create(ctx, course)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, global.ApiResponse{
			Code:    500,
			Message: "Failed to create course",
		})
		return
	}

	adminID := ctx.GetUint(global.UserIDKey)
	ctl.auditLogRepo.Create(ctx, models.AuditLog{
		CenterID:   centerID,
		ActorType:  "ADMIN",
		ActorID:    adminID,
		Action:     "CREATE_COURSE",
		TargetType: "Course",
		TargetID:   createdCourse.ID,
		Payload: models.AuditPayload{
			After: course,
		},
	})

	ctx.JSON(http.StatusOK, global.ApiResponse{
		Code:    0,
		Message: "Course created",
		Datas:   createdCourse,
	})
}

type CreateCourseRequest struct {
	Name             string `json:"name" binding:"required"`
	Duration         int    `json:"duration" binding:"required"`
	ColorHex         string `json:"color_hex" binding:"required"`
	RoomBufferMin    int    `json:"room_buffer_min" binding:"required"`
	TeacherBufferMin int    `json:"teacher_buffer_min" binding:"required"`
}

// DeleteCourse 刪除課程
// @Summary 刪除課程
// @Tags Admin
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param course_id path int true "Course ID"
// @Success 200 {object} global.ApiResponse
// @Router /api/v1/admin/courses/{course_id} [delete]
func (ctl *AdminResourceController) DeleteCourse(ctx *gin.Context) {
	centerID := ctx.GetUint(global.CenterIDKey)

	if centerID == 0 {
		if val, exists := ctx.Get(global.CenterIDKey); exists {
			switch v := val.(type) {
			case uint:
				centerID = v
			case uint64:
				centerID = uint(v)
			case int:
				centerID = uint(v)
			case float64:
				centerID = uint(v)
			}
		}
	}

	if centerID == 0 {
		ctx.JSON(http.StatusBadRequest, global.ApiResponse{
			Code:    global.BAD_REQUEST,
			Message: "Center ID required",
		})
		return
	}

	courseIDStr := ctx.Param("course_id")
	if courseIDStr == "" {
		ctx.JSON(http.StatusBadRequest, global.ApiResponse{
			Code:    global.BAD_REQUEST,
			Message: "Course ID required",
		})
		return
	}

	var courseID uint
	if _, err := fmt.Sscanf(courseIDStr, "%d", &courseID); err != nil {
		ctx.JSON(http.StatusBadRequest, global.ApiResponse{
			Code:    global.BAD_REQUEST,
			Message: "Invalid course ID",
		})
		return
	}

	if err := ctl.courseRepository.DeleteByID(ctl.makeCtx(ctx), courseID, centerID); err != nil {
		ctx.JSON(http.StatusInternalServerError, global.ApiResponse{
			Code:    500,
			Message: "Failed to delete course",
		})
		return
	}

	adminID := ctx.GetUint(global.UserIDKey)
	ctl.auditLogRepo.Create(ctx, models.AuditLog{
		CenterID:   centerID,
		ActorType:  "ADMIN",
		ActorID:    adminID,
		Action:     "DELETE_COURSE",
		TargetType: "Course",
		TargetID:   courseID,
		Payload: models.AuditPayload{
			After: map[string]interface{}{
				"status": "DELETED",
			},
		},
	})

	ctx.JSON(http.StatusOK, global.ApiResponse{
		Code:    0,
		Message: "Course deleted",
	})
}

// BulkCreateHolidays 批次建立假日
// @Summary 批次建立假日
// @Tags Admin
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Center ID"
// @Param request body BulkCreateHolidaysRequest true "假日列表"
// @Success 200 {object} global.ApiResponse{data=BulkCreateHolidaysResponse}
// @Router /api/v1/admin/centers/{id}/holidays/bulk [post]
func (ctl *AdminResourceController) BulkCreateHolidays(ctx *gin.Context) {
	centerIDStr := ctx.Param("id")
	if centerIDStr == "" {
		ctx.JSON(http.StatusBadRequest, global.ApiResponse{
			Code:    global.BAD_REQUEST,
			Message: "Center ID required",
		})
		return
	}

	var centerID uint
	if _, err := fmt.Sscanf(centerIDStr, "%d", &centerID); err != nil {
		ctx.JSON(http.StatusBadRequest, global.ApiResponse{
			Code:    global.BAD_REQUEST,
			Message: "Invalid center ID",
		})
		return
	}

	var req BulkCreateHolidaysRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, global.ApiResponse{
			Code:    global.BAD_REQUEST,
			Message: "Invalid request body",
		})
		return
	}

	if len(req.Holidays) == 0 {
		ctx.JSON(http.StatusBadRequest, global.ApiResponse{
			Code:    errInfos.PARAMS_VALIDATE_ERROR,
			Message: ctl.app.Err.New(errInfos.PARAMS_VALIDATE_ERROR).Msg,
		})
		return
	}

	if len(req.Holidays) > 100 {
		ctx.JSON(http.StatusBadRequest, global.ApiResponse{
			Code:    errInfos.LIMIT_EXCEEDED,
			Message: ctl.app.Err.New(errInfos.LIMIT_EXCEEDED).Msg,
		})
		return
	}

	holidays := make([]models.CenterHoliday, 0, len(req.Holidays))
	now := time.Now()

	for _, h := range req.Holidays {
		parsedDate, err := time.Parse("2006-01-02", h.Date)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, global.ApiResponse{
				Code:    errInfos.PARAMS_VALIDATE_ERROR,
				Message: "Invalid date format, expected YYYY-MM-DD: " + h.Date,
			})
			return
		}

		holidays = append(holidays, models.CenterHoliday{
			CenterID:  centerID,
			Date:      parsedDate,
			Name:      h.Name,
			CreatedAt: now,
			UpdatedAt: now,
		})
	}

	createdHolidays, createdCount, err := ctl.holidayRepository.BulkCreateWithSkipDuplicate(ctl.makeCtx(ctx), holidays)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, global.ApiResponse{
			Code:    errInfos.SQL_ERROR,
			Message: ctl.app.Err.New(errInfos.SQL_ERROR).Msg,
		})
		return
	}

	adminID := ctx.GetUint(global.UserIDKey)
	ctl.auditLogRepo.Create(ctx, models.AuditLog{
		CenterID:   centerID,
		ActorType:  "ADMIN",
		ActorID:    adminID,
		Action:     "BULK_CREATE_HOLIDAYS",
		TargetType: "CenterHoliday",
		TargetID:   0,
		Payload: models.AuditPayload{
			After: map[string]interface{}{
				"requested_count": len(req.Holidays),
				"created_count":   createdCount,
				"skipped_count":   len(req.Holidays) - int(createdCount),
			},
		},
	})

	ctx.JSON(http.StatusOK, global.ApiResponse{
		Code:    0,
		Message: "Holidays imported",
		Datas: BulkCreateHolidaysResponse{
			TotalRequested: len(req.Holidays),
			TotalCreated:   int(createdCount),
			TotalSkipped:   len(req.Holidays) - int(createdCount),
			Holidays:       createdHolidays,
		},
	})
}

type BulkCreateHolidaysRequest struct {
	Holidays []HolidayItem `json:"holidays" binding:"required,dive"`
}

type HolidayItem struct {
	Date string `json:"date" binding:"required"`
	Name string `json:"name" binding:"required"`
}

type BulkCreateHolidaysResponse struct {
	TotalRequested int                    `json:"total_requested"`
	TotalCreated   int                    `json:"total_created"`
	TotalSkipped   int                    `json:"total_skipped"`
	Holidays       []models.CenterHoliday `json:"holidays"`
}

type GetHolidaysRequest struct {
	StartDate string `form:"start_date"`
	EndDate   string `form:"end_date"`
}

func (ctl *AdminResourceController) GetHolidays(ctx *gin.Context) {
	centerIDStr := ctx.Param("id")
	var centerID uint
	if _, err := fmt.Sscanf(centerIDStr, "%d", &centerID); err != nil {
		ctx.JSON(http.StatusBadRequest, global.ApiResponse{
			Code:    global.BAD_REQUEST,
			Message: "Invalid center ID",
		})
		return
	}

	var req GetHolidaysRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, global.ApiResponse{
			Code:    global.BAD_REQUEST,
			Message: "Invalid query parameters",
		})
		return
	}

	var holidays []models.CenterHoliday
	var err error

	if req.StartDate != "" && req.EndDate != "" {
		startDate, _ := time.Parse("2006-01-02", req.StartDate)
		endDate, _ := time.Parse("2006-01-02", req.EndDate)
		holidays, err = ctl.holidayRepository.ListByDateRange(ctx, centerID, startDate, endDate)
	} else {
		holidays, err = ctl.holidayRepository.ListByCenterID(ctx, centerID)
	}

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, global.ApiResponse{
			Code:    errInfos.SQL_ERROR,
			Message: ctl.app.Err.New(errInfos.SQL_ERROR).Msg,
		})
		return
	}

	ctx.JSON(http.StatusOK, global.ApiResponse{
		Code:    0,
		Message: "OK",
		Datas:   holidays,
	})
}

type CreateHolidayRequest struct {
	Date string `json:"date" binding:"required"`
	Name string `json:"name" binding:"required"`
}

func (ctl *AdminResourceController) CreateHoliday(ctx *gin.Context) {
	centerIDStr := ctx.Param("id")
	var centerID uint
	if _, err := fmt.Sscanf(centerIDStr, "%d", &centerID); err != nil {
		ctx.JSON(http.StatusBadRequest, global.ApiResponse{
			Code:    global.BAD_REQUEST,
			Message: "Invalid center ID",
		})
		return
	}

	var req CreateHolidayRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, global.ApiResponse{
			Code:    errInfos.PARAMS_VALIDATE_ERROR,
			Message: "Invalid request body",
		})
		return
	}

	parsedDate, err := time.Parse("2006-01-02", req.Date)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, global.ApiResponse{
			Code:    errInfos.PARAMS_VALIDATE_ERROR,
			Message: "Invalid date format, expected YYYY-MM-DD",
		})
		return
	}

	exists, err := ctl.holidayRepository.ExistsByCenterAndDate(ctx, centerID, parsedDate)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, global.ApiResponse{
			Code:    errInfos.SQL_ERROR,
			Message: ctl.app.Err.New(errInfos.SQL_ERROR).Msg,
		})
		return
	}

	if exists {
		ctx.JSON(http.StatusConflict, global.ApiResponse{
			Code:    errInfos.DUPLICATE,
			Message: "Holiday already exists for this date",
		})
		return
	}

	holiday := models.CenterHoliday{
		CenterID: centerID,
		Date:     parsedDate,
		Name:     req.Name,
	}

	created, err := ctl.holidayRepository.Create(ctx, holiday)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, global.ApiResponse{
			Code:    errInfos.SQL_ERROR,
			Message: ctl.app.Err.New(errInfos.SQL_ERROR).Msg,
		})
		return
	}

	adminID := ctx.GetUint(global.UserIDKey)
	ctl.auditLogRepo.Create(ctx, models.AuditLog{
		CenterID:   centerID,
		ActorType:  "ADMIN",
		ActorID:    adminID,
		Action:     "CREATE_HOLIDAY",
		TargetType: "CenterHoliday",
		TargetID:   created.ID,
		Payload: models.AuditPayload{
			After: holiday,
		},
	})

	ctx.JSON(http.StatusOK, global.ApiResponse{
		Code:    0,
		Message: "Holiday created",
		Datas:   created,
	})
}

func (ctl *AdminResourceController) DeleteHoliday(ctx *gin.Context) {
	centerIDStr := ctx.Param("id")
	var centerID uint
	if _, err := fmt.Sscanf(centerIDStr, "%d", &centerID); err != nil {
		ctx.JSON(http.StatusBadRequest, global.ApiResponse{
			Code:    global.BAD_REQUEST,
			Message: "Invalid center ID",
		})
		return
	}

	holidayIDStr := ctx.Param("holiday_id")
	var holidayID uint
	if _, err := fmt.Sscanf(holidayIDStr, "%d", &holidayID); err != nil {
		ctx.JSON(http.StatusBadRequest, global.ApiResponse{
			Code:    global.BAD_REQUEST,
			Message: "Invalid holiday ID",
		})
		return
	}

	holiday, err := ctl.holidayRepository.GetByID(ctx, holidayID)
	if err != nil {
		ctx.JSON(http.StatusNotFound, global.ApiResponse{
			Code:    errInfos.NOT_FOUND,
			Message: "Holiday not found",
		})
		return
	}

	if holiday.CenterID != centerID {
		ctx.JSON(http.StatusForbidden, global.ApiResponse{
			Code:    errInfos.FORBIDDEN,
			Message: "Holiday does not belong to this center",
		})
		return
	}

	if err := ctl.holidayRepository.Delete(ctx, holidayID); err != nil {
		ctx.JSON(http.StatusInternalServerError, global.ApiResponse{
			Code:    errInfos.SQL_ERROR,
			Message: ctl.app.Err.New(errInfos.SQL_ERROR).Msg,
		})
		return
	}

	adminID := ctx.GetUint(global.UserIDKey)
	ctl.auditLogRepo.Create(ctx, models.AuditLog{
		CenterID:   centerID,
		ActorType:  "ADMIN",
		ActorID:    adminID,
		Action:     "DELETE_HOLIDAY",
		TargetType: "CenterHoliday",
		TargetID:   holidayID,
		Payload: models.AuditPayload{
			Before: holiday,
		},
	})

	ctx.JSON(http.StatusOK, global.ApiResponse{
		Code:    0,
		Message: "Holiday deleted",
	})
}

func (ctl *AdminResourceController) GetActiveRooms(ctx *gin.Context) {
	centerID := ctl.getCenterID(ctx)
	if centerID == 0 {
		ctl.respondError(ctx, global.BAD_REQUEST, "Center ID required")
		return
	}

	rooms, err := ctl.roomRepository.ListActiveByCenterID(ctx, centerID)
	if err != nil {
		ctl.respondError(ctx, errInfos.SQL_ERROR, "Failed to get active rooms")
		return
	}

	ctx.JSON(http.StatusOK, global.ApiResponse{
		Code:    0,
		Message: "Success",
		Datas:   rooms,
	})
}

func (ctl *AdminResourceController) GetActiveCourses(ctx *gin.Context) {
	centerID := ctl.getCenterID(ctx)
	if centerID == 0 {
		ctl.respondError(ctx, global.BAD_REQUEST, "Center ID required")
		return
	}

	courses, err := ctl.courseRepository.ListActiveByCenterID(ctx, centerID)
	if err != nil {
		ctl.respondError(ctx, errInfos.SQL_ERROR, "Failed to get active courses")
		return
	}

	ctx.JSON(http.StatusOK, global.ApiResponse{
		Code:    0,
		Message: "Success",
		Datas:   courses,
	})
}

func (ctl *AdminResourceController) GetActiveOfferings(ctx *gin.Context) {
	centerID := ctl.getCenterID(ctx)
	if centerID == 0 {
		ctl.respondError(ctx, global.BAD_REQUEST, "Center ID required")
		return
	}

	offerings, err := ctl.offeringRepository.ListActiveByCenterID(ctx, centerID)
	if err != nil {
		ctl.respondError(ctx, errInfos.SQL_ERROR, "Failed to get active offerings")
		return
	}

	ctx.JSON(http.StatusOK, global.ApiResponse{
		Code:    0,
		Message: "Success",
		Datas:   offerings,
	})
}

func (ctl *AdminResourceController) ToggleCourseActive(ctx *gin.Context) {
	centerID := ctl.getCenterID(ctx)
	if centerID == 0 {
		ctl.respondError(ctx, global.BAD_REQUEST, "Center ID required")
		return
	}

	courseID, err := ctl.getUintParam(ctx, "course_id")
	if err != nil {
		ctl.respondError(ctx, global.BAD_REQUEST, "Invalid course ID")
		return
	}

	var req ToggleActiveRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctl.respondError(ctx, errInfos.PARAMS_VALIDATE_ERROR, "Invalid request body")
		return
	}

	if err := ctl.courseRepository.ToggleActive(ctx, courseID, centerID, req.IsActive); err != nil {
		ctl.respondError(ctx, errInfos.SQL_ERROR, "Failed to toggle course active status: "+err.Error())
		return
	}

	adminID := ctx.GetUint(global.UserIDKey)
	ctl.auditLogRepo.Create(ctx, models.AuditLog{
		CenterID:   centerID,
		ActorType:  "ADMIN",
		ActorID:    adminID,
		Action:     "TOGGLE_COURSE_ACTIVE",
		TargetType: "Course",
		TargetID:   courseID,
		Payload: models.AuditPayload{
			After: map[string]interface{}{
				"is_active": req.IsActive,
			},
		},
	})

	ctx.JSON(http.StatusOK, global.ApiResponse{
		Code:    0,
		Message: "Course active status updated",
	})
}

func (ctl *AdminResourceController) ToggleRoomActive(ctx *gin.Context) {
	centerID := ctl.getCenterID(ctx)
	if centerID == 0 {
		ctl.respondError(ctx, global.BAD_REQUEST, "Center ID required")
		return
	}

	roomID, err := ctl.getUintParam(ctx, "room_id")
	if err != nil {
		ctl.respondError(ctx, global.BAD_REQUEST, "Invalid room ID")
		return
	}

	var req ToggleActiveRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctl.respondError(ctx, errInfos.PARAMS_VALIDATE_ERROR, "Invalid request body")
		return
	}

	if err := ctl.roomRepository.ToggleActive(ctx, roomID, centerID, req.IsActive); err != nil {
		ctl.respondError(ctx, errInfos.SQL_ERROR, "Failed to toggle room active status")
		return
	}

	adminID := ctx.GetUint(global.UserIDKey)
	ctl.auditLogRepo.Create(ctx, models.AuditLog{
		CenterID:   centerID,
		ActorType:  "ADMIN",
		ActorID:    adminID,
		Action:     "TOGGLE_ROOM_ACTIVE",
		TargetType: "Room",
		TargetID:   roomID,
		Payload: models.AuditPayload{
			After: map[string]interface{}{
				"is_active": req.IsActive,
			},
		},
	})

	ctx.JSON(http.StatusOK, global.ApiResponse{
		Code:    0,
		Message: "Room active status updated",
	})
}

func (ctl *AdminResourceController) ToggleOfferingActive(ctx *gin.Context) {
	centerID := ctl.getCenterID(ctx)
	if centerID == 0 {
		ctl.respondError(ctx, global.BAD_REQUEST, "Center ID required")
		return
	}

	offeringID, err := ctl.getUintParam(ctx, "offering_id")
	if err != nil {
		ctl.respondError(ctx, global.BAD_REQUEST, "Invalid offering ID")
		return
	}

	var req ToggleActiveRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctl.respondError(ctx, errInfos.PARAMS_VALIDATE_ERROR, "Invalid request body")
		return
	}

	if err := ctl.offeringRepository.ToggleActive(ctx, offeringID, centerID, req.IsActive); err != nil {
		ctl.respondError(ctx, errInfos.SQL_ERROR, "Failed to toggle offering active status")
		return
	}

	adminID := ctx.GetUint(global.UserIDKey)
	ctl.auditLogRepo.Create(ctx, models.AuditLog{
		CenterID:   centerID,
		ActorType:  "ADMIN",
		ActorID:    adminID,
		Action:     "TOGGLE_OFFERING_ACTIVE",
		TargetType: "Offering",
		TargetID:   offeringID,
		Payload: models.AuditPayload{
			After: map[string]interface{}{
				"is_active": req.IsActive,
			},
		},
	})

	ctx.JSON(http.StatusOK, global.ApiResponse{
		Code:    0,
		Message: "Offering active status updated",
	})
}

type ToggleActiveRequest struct {
	IsActive bool `json:"is_active"`
}

type InvitationStatsResponse struct {
	Total         int64 `json:"total"`
	Pending       int64 `json:"pending"`
	Accepted      int64 `json:"accepted"`
	Expired       int64 `json:"expired"`
	Rejected      int64 `json:"rejected"`
	RecentPending int64 `json:"recent_pending"`
}

func (ctl *AdminResourceController) GetInvitationStats(ctx *gin.Context) {
	centerID := ctl.getCenterID(ctx)
	if centerID == 0 {
		ctl.respondError(ctx, global.BAD_REQUEST, "Center ID required")
		return
	}

	now := time.Now()
	thirtyDaysAgo := now.AddDate(0, 0, -30)

	total, _ := ctl.invitationRepo.CountByCenterID(ctx, centerID)
	pending, _ := ctl.invitationRepo.CountByStatus(ctx, centerID, "PENDING")
	accepted, _ := ctl.invitationRepo.CountByStatus(ctx, centerID, "ACCEPTED")
	expired, _ := ctl.invitationRepo.CountByStatus(ctx, centerID, "EXPIRED")
	rejected, _ := ctl.invitationRepo.CountByStatus(ctx, centerID, "REJECTED")
	recentPending, _ := ctl.invitationRepo.CountByDateRange(ctx, centerID, thirtyDaysAgo, now, "PENDING")

	stats := InvitationStatsResponse{
		Total:         total,
		Pending:       pending,
		Accepted:      accepted,
		Expired:       expired,
		Rejected:      rejected,
		RecentPending: recentPending,
	}

	ctx.JSON(http.StatusOK, global.ApiResponse{
		Code:    0,
		Message: "Success",
		Datas:   stats,
	})
}

func (ctl *AdminResourceController) GetInvitations(ctx *gin.Context) {
	centerID := ctl.getCenterID(ctx)
	if centerID == 0 {
		ctl.respondError(ctx, global.BAD_REQUEST, "Center ID required")
		return
	}

	var req PaginationRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		req.Page = 1
		req.Limit = 20
	}

	status := ctx.Query("status")

	invitations, total, err := ctl.invitationRepo.ListByCenterIDPaginated(ctx, centerID, int64(req.Page), int64(req.Limit), status)
	if err != nil {
		ctl.respondError(ctx, errInfos.SQL_ERROR, "Failed to get invitations")
		return
	}

	ctx.JSON(http.StatusOK, global.ApiResponse{
		Code:    0,
		Message: "Success",
		Datas: PaginationResponse{
			Data:       invitations,
			Total:      total,
			Page:       req.Page,
			Limit:      req.Limit,
			TotalPages: (total + int64(req.Limit) - 1) / int64(req.Limit),
		},
	})
}

type PaginationRequest struct {
	Page  int `form:"page"`
	Limit int `form:"limit"`
}

type PaginationResponse struct {
	Data       interface{} `json:"data"`
	Total      int64       `json:"total"`
	Page       int         `json:"page"`
	Limit      int         `json:"limit"`
	TotalPages int64       `json:"total_pages"`
}

func (ctl *AdminResourceController) getCenterID(ctx *gin.Context) uint {
	centerID := ctx.GetUint(global.CenterIDKey)
	if centerID == 0 {
		if val, exists := ctx.Get(global.CenterIDKey); exists {
			switch v := val.(type) {
			case uint:
				centerID = v
			case uint64:
				centerID = uint(v)
			case int:
				centerID = uint(v)
			case float64:
				centerID = uint(v)
			}
		}
	}
	return centerID
}

func (ctl *AdminResourceController) getUintParam(ctx *gin.Context, paramName string) (uint, error) {
	paramStr := ctx.Param(paramName)
	if paramStr == "" {
		return 0, fmt.Errorf("parameter not found")
	}
	var result uint
	_, err := fmt.Sscanf(paramStr, "%d", &result)
	return result, err
}

func (ctl *AdminResourceController) respondError(ctx *gin.Context, code errInfos.ErrCode, message string) {
	ctx.JSON(http.StatusBadRequest, global.ApiResponse{
		Code:    code,
		Message: message,
	})
}
