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

type HolidayService struct {
	BaseService
	app          *app.App
	holidayRepo  *repositories.CenterHolidayRepository
	auditLogRepo *repositories.AuditLogRepository
}

func NewHolidayService(app *app.App) *HolidayService {
	return &HolidayService{
		app:          app,
		holidayRepo:  repositories.NewCenterHolidayRepository(app),
		auditLogRepo: repositories.NewAuditLogRepository(app),
	}
}

type HolidayItem struct {
	Date        string `json:"date" binding:"required"`
	Name        string `json:"name" binding:"required"`
	ForceCancel bool   `json:"force_cancel"`
}

type BulkCreateHolidaysRequest struct {
	Holidays []HolidayItem `json:"holidays" binding:"required,dive"`
}

type BulkCreateHolidaysResponse struct {
	TotalRequested int                    `json:"total_requested"`
	TotalCreated   int                    `json:"total_created"`
	TotalSkipped   int                    `json:"total_skipped"`
	Holidays       []models.CenterHoliday `json:"holidays"`
}

type CreateHolidayRequest struct {
	Date        string `json:"date" binding:"required"`
	Name        string `json:"name" binding:"required"`
	ForceCancel bool   `json:"force_cancel"`
}

type GetHolidaysRequest struct {
	StartDate string `form:"start_date"`
	EndDate   string `form:"end_date"`
}

func (s *HolidayService) GetHolidays(ctx context.Context, centerID uint, req *GetHolidaysRequest) ([]models.CenterHoliday, *errInfos.Res, error) {
	var holidays []models.CenterHoliday
	var err error

	if req.StartDate != "" && req.EndDate != "" {
		startDate, _ := time.Parse("2006-01-02", req.StartDate)
		endDate, _ := time.Parse("2006-01-02", req.EndDate)
		holidays, err = s.holidayRepo.ListByDateRange(ctx, centerID, startDate, endDate)
	} else {
		holidays, err = s.holidayRepo.ListByCenterID(ctx, centerID)
	}

	if err != nil {
		return nil, s.app.Err.New(errInfos.SQL_ERROR), err
	}
	return holidays, nil, nil
}

func (s *HolidayService) CreateHoliday(ctx context.Context, centerID, adminID uint, req *CreateHolidayRequest) (*models.CenterHoliday, *errInfos.Res, error) {
	parsedDate, err := time.Parse("2006-01-02", req.Date)
	if err != nil {
		return nil, s.app.Err.New(errInfos.PARAMS_VALIDATE_ERROR), fmt.Errorf("invalid date format")
	}

	exists, err := s.holidayRepo.ExistsByCenterAndDate(ctx, centerID, parsedDate)
	if err != nil {
		return nil, s.app.Err.New(errInfos.SQL_ERROR), err
	}
	if exists {
		return nil, s.app.Err.New(errInfos.DUPLICATE), fmt.Errorf("holiday already exists")
	}

	holiday := models.CenterHoliday{
		CenterID:     centerID,
		Date:         parsedDate,
		Name:         req.Name,
		ForceCancel:  req.ForceCancel,
	}

	var created models.CenterHoliday
	err = s.holidayRepo.Transaction(ctx, func(txRepo *repositories.CenterHolidayRepository) error {
		var txErr error
		created, txErr = txRepo.Create(ctx, holiday)
		if txErr != nil {
			return txErr
		}

		// 記錄稽核日誌
		auditLog := models.AuditLog{
			CenterID:   centerID,
			ActorType:  "ADMIN",
			ActorID:    adminID,
			Action:     "CREATE_HOLIDAY",
			TargetType: "CenterHoliday",
			TargetID:   created.ID,
			Payload: models.AuditPayload{
				After: created,
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

	return &created, nil, nil
}

func (s *HolidayService) BulkCreateHolidays(ctx context.Context, centerID, adminID uint, req *BulkCreateHolidaysRequest) (*BulkCreateHolidaysResponse, *errInfos.Res, error) {
	if len(req.Holidays) == 0 {
		return nil, s.app.Err.New(errInfos.PARAMS_VALIDATE_ERROR), nil
	}

	holidays := make([]models.CenterHoliday, 0, len(req.Holidays))
	now := time.Now()

	for _, h := range req.Holidays {
		parsedDate, err := time.Parse("2006-01-02", h.Date)
		if err != nil {
			return nil, s.app.Err.New(errInfos.PARAMS_VALIDATE_ERROR), fmt.Errorf("invalid date: %s", h.Date)
		}
		holidays = append(holidays, models.CenterHoliday{
			CenterID:     centerID,
			Date:         parsedDate,
			Name:         h.Name,
			ForceCancel:  h.ForceCancel,
			CreatedAt:    now,
			UpdatedAt:    now,
		})
	}

	createdHolidays, createdCount, err := s.holidayRepo.BulkCreateWithSkipDuplicate(ctx, holidays)
	if err != nil {
		return nil, s.app.Err.New(errInfos.SQL_ERROR), err
	}

	// 記錄稽核日誌
	s.auditLogRepo.Create(ctx, models.AuditLog{
		CenterID:   centerID,
		ActorType:  "ADMIN",
		ActorID:    adminID,
		Action:     "BULK_CREATE_HOLIDAYS",
		TargetType: "CenterHoliday",
		Payload: models.AuditPayload{
			After: map[string]interface{}{
				"requested_count": len(req.Holidays),
				"created_count":   createdCount,
				"skipped_count":   len(req.Holidays) - int(createdCount),
			},
		},
	})

	return &BulkCreateHolidaysResponse{
		TotalRequested: len(req.Holidays),
		TotalCreated:   int(createdCount),
		TotalSkipped:   len(req.Holidays) - int(createdCount),
		Holidays:       createdHolidays,
	}, nil, nil
}

func (s *HolidayService) DeleteHoliday(ctx context.Context, centerID, adminID, holidayID uint) (*errInfos.Res, error) {
	err := s.holidayRepo.Transaction(ctx, func(txRepo *repositories.CenterHolidayRepository) error {
		holiday, err := txRepo.GetByID(ctx, holidayID)
		if err != nil {
			return err
		}
		if holiday.CenterID != centerID {
			return fmt.Errorf("permission denied")
		}

		if err := txRepo.Delete(ctx, holidayID); err != nil {
			return err
		}

		// 記錄稽核日誌
		auditLog := models.AuditLog{
			CenterID:   centerID,
			ActorType:  "ADMIN",
			ActorID:    adminID,
			Action:     "DELETE_HOLIDAY",
			TargetType: "CenterHoliday",
			TargetID:   holidayID,
			Payload: models.AuditPayload{
				Before: holiday,
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

	return nil, nil
}
