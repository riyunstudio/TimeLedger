package services

import (
	"context"
	"time"
	"timeLedger/app"
	"timeLedger/app/models"
	"timeLedger/app/repositories"
	"timeLedger/global/errInfos"
)

// CenterService 中心相關業務邏輯
type CenterService struct {
	BaseService
	app            *app.App
	centerRepo     *repositories.CenterRepository
	holidayRepo    *repositories.CenterHolidayRepository
	invitationRepo *repositories.CenterInvitationRepository
	auditLogRepo   *repositories.AuditLogRepository
}

// NewCenterService 建立 CenterService 實例
func NewCenterService(appInstance *app.App) *CenterService {
	return &CenterService{
		app:            appInstance,
		centerRepo:     repositories.NewCenterRepository(appInstance),
		holidayRepo:    repositories.NewCenterHolidayRepository(appInstance),
		invitationRepo: repositories.NewCenterInvitationRepository(appInstance),
		auditLogRepo:   repositories.NewAuditLogRepository(appInstance),
	}
}

// ListCenters 取得所有中心列表
func (s *CenterService) ListCenters(ctx context.Context) ([]models.Center, *errInfos.Res, error) {
	centers, err := s.centerRepo.List(ctx)
	if err != nil {
		return nil, s.app.Err.New(errInfos.SQL_ERROR), err
	}
	return centers, nil, nil
}

// CreateCenterInput 建立中心的輸入資料
type CreateCenterInput struct {
	Name                string
	PlanLevel           string
	AllowPublicRegister bool
}

// CreateCenter 建立新中心
func (s *CenterService) CreateCenter(ctx context.Context, adminID uint, input *CreateCenterInput) (*models.Center, *errInfos.Res, error) {
	center := models.Center{
		Name:      input.Name,
		PlanLevel: input.PlanLevel,
		Settings: models.CenterSettings{
			AllowPublicRegister: input.AllowPublicRegister,
			DefaultLanguage:     "zh-TW",
		},
		CreatedAt: time.Now(),
	}

	createdCenter, err := s.centerRepo.Create(ctx, center)
	if err != nil {
		return nil, s.app.Err.New(errInfos.SQL_ERROR), err
	}

	// 記錄稽核日誌
	s.auditLogRepo.Create(ctx, models.AuditLog{
		CenterID:   createdCenter.ID,
		ActorType:  "ADMIN",
		ActorID:    adminID,
		Action:     "CREATE_CENTER",
		TargetType: "Center",
		TargetID:   createdCenter.ID,
		Payload: models.AuditPayload{
			After: center,
		},
	})

	return &createdCenter, nil, nil
}
