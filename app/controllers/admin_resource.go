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
	app               *app.App
	centerRepository  *repositories.CenterRepository
	courseRepository  *repositories.CourseRepository
	roomRepository    *repositories.RoomRepository
	holidayRepository *repositories.CenterHolidayRepository
	auditLogRepo      *repositories.AuditLogRepository
}

func NewAdminResourceController(app *app.App) *AdminResourceController {
	return &AdminResourceController{
		app:               app,
		centerRepository:  repositories.NewCenterRepository(app),
		courseRepository:  repositories.NewCourseRepository(app),
		roomRepository:    repositories.NewRoomRepository(app),
		holidayRepository: repositories.NewCenterHolidayRepository(app),
		auditLogRepo:      repositories.NewAuditLogRepository(app),
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
