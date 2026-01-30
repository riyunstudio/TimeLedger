package services

import (
	"context"
	"fmt"
	"time"
	"timeLedger/app"
	"timeLedger/app/models"
	"timeLedger/app/repositories"
	"timeLedger/global/errInfos"
)

type CourseService struct {
	BaseService
	app          *app.App
	courseRepo   *repositories.CourseRepository
	auditLogRepo *repositories.AuditLogRepository
	cacheService *CacheService
}

func NewCourseService(app *app.App) *CourseService {
	return &CourseService{
		app:          app,
		courseRepo:   repositories.NewCourseRepository(app),
		auditLogRepo: repositories.NewAuditLogRepository(app),
		cacheService: NewCacheService(app),
	}
}

type CourseResponse struct {
	ID               uint      `json:"id"`
	Name             string    `json:"name"`
	DefaultDuration  int       `json:"default_duration"`
	ColorHex         string    `json:"color_hex"`
	RoomBufferMin    int       `json:"room_buffer_min"`
	TeacherBufferMin int       `json:"teacher_buffer_min"`
	IsActive         bool      `json:"is_active"`
	CreatedAt        time.Time `json:"created_at"`
}

type CreateCourseRequest struct {
	Name             string `json:"name" binding:"required"`
	Duration         int    `json:"duration" binding:"required"`
	ColorHex         string `json:"color_hex" binding:"required"`
	RoomBufferMin    int    `json:"room_buffer_min" binding:"gte=0"`
	TeacherBufferMin int    `json:"teacher_buffer_min" binding:"gte=0"`
}

type UpdateCourseRequest struct {
	Name             string `json:"name" binding:"required"`
	Duration         int    `json:"duration" binding:"required"`
	ColorHex         string `json:"color_hex" binding:"required"`
	RoomBufferMin    int    `json:"room_buffer_min" binding:"gte=0"`
	TeacherBufferMin int    `json:"teacher_buffer_min" binding:"gte=0"`
}

func (s *CourseService) GetCourses(ctx context.Context, centerID uint) ([]models.Course, *errInfos.Res, error) {
	// 先從快取取得
	cached, err := s.cacheService.GetCourseList(ctx, centerID)
	if err == nil && cached != nil {
		// 將快取項目轉換為 models.Course
		courses := make([]models.Course, 0, len(cached))
		for _, item := range cached {
			courses = append(courses, models.Course{
				ID:               item.ID,
				CenterID:         centerID,
				Name:             item.Name,
				DefaultDuration:  item.DefaultDuration,
				ColorHex:         item.ColorHex,
				RoomBufferMin:    item.RoomBufferMin,
				TeacherBufferMin: item.TeacherBufferMin,
				IsActive:         item.IsActive,
			})
		}
		return courses, nil, nil
	}

	// 快取未命中，從資料庫取得
	courses, err := s.courseRepo.ListByCenterID(ctx, centerID)
	if err != nil {
		return nil, s.app.Err.New(errInfos.SQL_ERROR), err
	}

	// 存入快取
	cacheItems := make([]CourseCacheItem, 0, len(courses))
	for _, c := range courses {
		cacheItems = append(cacheItems, CourseCacheItem{
			ID:               c.ID,
			Name:             c.Name,
			DefaultDuration:  c.DefaultDuration,
			ColorHex:         c.ColorHex,
			RoomBufferMin:    c.RoomBufferMin,
			TeacherBufferMin: c.TeacherBufferMin,
			IsActive:         c.IsActive,
		})
	}
	_ = s.cacheService.SetCourseList(ctx, centerID, cacheItems)

	return courses, nil, nil
}

func (s *CourseService) GetActiveCourses(ctx context.Context, centerID uint) ([]models.Course, *errInfos.Res, error) {
	courses, err := s.courseRepo.ListActiveByCenterID(ctx, centerID)
	if err != nil {
		return nil, s.app.Err.New(errInfos.SQL_ERROR), err
	}
	return courses, nil, nil
}

func (s *CourseService) CreateCourse(ctx context.Context, centerID, adminID uint, req *CreateCourseRequest) (*models.Course, *errInfos.Res, error) {
	course := models.Course{
		CenterID:         centerID,
		Name:             req.Name,
		DefaultDuration:  req.Duration,
		ColorHex:         req.ColorHex,
		RoomBufferMin:    req.RoomBufferMin,
		TeacherBufferMin: req.TeacherBufferMin,
		IsActive:         true,
		CreatedAt:        time.Now(),
	}

	var createdCourse models.Course
	err := s.courseRepo.Transaction(ctx, func(txRepo *repositories.CourseRepository) error {
		var txErr error
		createdCourse, txErr = txRepo.Create(ctx, course)
		if txErr != nil {
			return txErr
		}

		// 記錄稽核日誌
		auditLog := models.AuditLog{
			CenterID:   centerID,
			ActorType:  "ADMIN",
			ActorID:    adminID,
			Action:     "CREATE_COURSE",
			TargetType: "Course",
			TargetID:   createdCourse.ID,
			Payload: models.AuditPayload{
				After: createdCourse,
			},
		}

		if err := txRepo.GetDBWrite().Create(&auditLog).Error; err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return nil, s.app.Err.New(errInfos.SQL_ERROR), err
	}

	// 清除課程列表快取
	_ = s.cacheService.InvalidateCourseList(ctx, centerID)

	return &createdCourse, nil, nil
}

func (s *CourseService) UpdateCourse(ctx context.Context, centerID, adminID, courseID uint, req *UpdateCourseRequest) (*models.Course, *errInfos.Res, error) {
	var updatedCourse models.Course
	err := s.courseRepo.Transaction(ctx, func(txRepo *repositories.CourseRepository) error {
		// 查詢現有課程
		course, err := txRepo.GetByID(ctx, courseID)
		if err != nil {
			return err
		}

		// 驗證權限
		if course.CenterID != centerID {
			return fmt.Errorf("permission denied")
		}

		before := course

		// 更新
		course.Name = req.Name
		course.DefaultDuration = req.Duration
		course.ColorHex = req.ColorHex
		course.RoomBufferMin = req.RoomBufferMin
		course.TeacherBufferMin = req.TeacherBufferMin
		course.UpdatedAt = time.Now()

		if err := txRepo.Update(ctx, course); err != nil {
			return err
		}

		// 記錄稽核日誌
		auditLog := models.AuditLog{
			CenterID:   centerID,
			ActorType:  "ADMIN",
			ActorID:    adminID,
			Action:     "UPDATE_COURSE",
			TargetType: "Course",
			TargetID:   course.ID,
			Payload: models.AuditPayload{
				Before: before,
				After:  course,
			},
		}

		if err := txRepo.GetDBWrite().Create(&auditLog).Error; err != nil {
			return err
		}

		updatedCourse = course
		return nil
	})

	if err != nil {
		if err.Error() == "permission denied" {
			return nil, s.app.Err.New(errInfos.FORBIDDEN), err
		}
		return nil, s.app.Err.New(errInfos.SQL_ERROR), err
	}

	// 清除課程列表快取
	_ = s.cacheService.InvalidateCourseList(ctx, centerID)

	return &updatedCourse, nil, nil
}

func (s *CourseService) DeleteCourse(ctx context.Context, centerID, adminID, courseID uint) (*errInfos.Res, error) {
	err := s.courseRepo.Transaction(ctx, func(txRepo *repositories.CourseRepository) error {
		// 驗證是否存在且屬於該中心
		course, err := txRepo.GetByID(ctx, courseID)
		if err != nil {
			return err
		}
		if course.CenterID != centerID {
			return fmt.Errorf("permission denied")
		}

		if err := txRepo.DeleteByIDWithCenterScope(ctx, courseID, centerID); err != nil {
			return err
		}

		// 記錄稽核日誌
		auditLog := models.AuditLog{
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
		}

		if err := txRepo.GetDBWrite().Create(&auditLog).Error; err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		if err.Error() == "permission denied" {
			return s.app.Err.New(errInfos.FORBIDDEN), err
		}
		return s.app.Err.New(errInfos.SQL_ERROR), err
	}

	// 清除課程列表快取
	_ = s.cacheService.InvalidateCourseList(ctx, centerID)

	return nil, nil
}

func (s *CourseService) ToggleCourseActive(ctx context.Context, centerID, adminID, courseID uint, isActive bool) (*errInfos.Res, error) {
	err := s.courseRepo.Transaction(ctx, func(txRepo *repositories.CourseRepository) error {
		// 驗證
		course, err := txRepo.GetByID(ctx, courseID)
		if err != nil {
			return err
		}
		if course.CenterID != centerID {
			return fmt.Errorf("permission denied")
		}

		if err := txRepo.ToggleActive(ctx, courseID, centerID, isActive); err != nil {
			return err
		}

		// 記錄稽核日誌
		auditLog := models.AuditLog{
			CenterID:   centerID,
			ActorType:  "ADMIN",
			ActorID:    adminID,
			Action:     "TOGGLE_COURSE_ACTIVE",
			TargetType: "Course",
			TargetID:   courseID,
			Payload: models.AuditPayload{
				After: map[string]interface{}{
					"is_active": isActive,
				},
			},
		}

		if err := txRepo.GetDBWrite().Create(&auditLog).Error; err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		if err.Error() == "permission denied" {
			return s.app.Err.New(errInfos.FORBIDDEN), err
		}
		return s.app.Err.New(errInfos.SQL_ERROR), err
	}

	// 清除課程列表快取
	_ = s.cacheService.InvalidateCourseList(ctx, centerID)

	return nil, nil
}

func (s *CourseService) ToCourseResponse(course models.Course) CourseResponse {
	return CourseResponse{
		ID:               course.ID,
		Name:             course.Name,
		DefaultDuration:  course.DefaultDuration,
		ColorHex:         course.ColorHex,
		RoomBufferMin:    course.RoomBufferMin,
		TeacherBufferMin: course.TeacherBufferMin,
		IsActive:         course.IsActive,
		CreatedAt:        course.CreatedAt,
	}
}

func (s *CourseService) ToCourseResponses(courses []models.Course) []CourseResponse {
	responses := make([]CourseResponse, 0, len(courses))
	for _, c := range courses {
		responses = append(responses, s.ToCourseResponse(c))
	}
	return responses
}
