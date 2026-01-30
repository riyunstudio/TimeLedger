package services

import (
	"context"
	"time"
	"timeLedger/app"
	"timeLedger/app/models"
	"timeLedger/app/repositories"
	"timeLedger/global/errInfos"
	"timeLedger/global/logger"
)

// CenterService 中心相關業務邏輯
type CenterService struct {
	BaseService
	app            *app.App
	centerRepo     *repositories.CenterRepository
	holidayRepo    *repositories.CenterHolidayRepository
	invitationRepo *repositories.CenterInvitationRepository
	auditLogRepo   *repositories.AuditLogRepository
	cacheService   *CacheService
}

// NewCenterService 建立 CenterService 實例
func NewCenterService(appInstance *app.App) *CenterService {
	return &CenterService{
		app:            appInstance,
		centerRepo:     repositories.NewCenterRepository(appInstance),
		holidayRepo:    repositories.NewCenterHolidayRepository(appInstance),
		invitationRepo: repositories.NewCenterInvitationRepository(appInstance),
		auditLogRepo:   repositories.NewAuditLogRepository(appInstance),
		cacheService:   NewCacheService(appInstance),
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

// GetCenterBasic 快取取得中心基本資訊（Cache-Aside 模式）
func (s *CenterService) GetCenterBasic(ctx context.Context, centerID uint) (*CenterBasicInfo, *errInfos.Res, error) {
	// 1. 先從快取取得
	cached, err := s.cacheService.GetCenterBasic(ctx, centerID)
	if err == nil && cached != nil {
		return cached, nil, nil
	}

	// 2. 快取未命中，從資料庫取得
	center, err := s.centerRepo.GetByID(ctx, centerID)
	if err != nil {
		return nil, s.app.Err.New(errInfos.NOT_FOUND), err
	}

	// 3. 存入快取
	info := &CenterBasicInfo{
		ID:        center.ID,
		Name:      center.Name,
		PlanLevel: center.PlanLevel,
	}
	_ = s.cacheService.SetCenterBasic(ctx, centerID, info)

	return info, nil, nil
}

// GetCenterSettings 快取取得中心設置（Cache-Aside 模式）
func (s *CenterService) GetCenterSettings(ctx context.Context, centerID uint) (*models.CenterSettings, *errInfos.Res, error) {
	// 1. 先從快取取得
	cached, err := s.cacheService.GetCenterSettings(ctx, centerID)
	if err == nil && cached != nil {
		return cached, nil, nil
	}

	// 2. 快取未命中，從資料庫取得
	center, err := s.centerRepo.GetByID(ctx, centerID)
	if err != nil {
		return nil, s.app.Err.New(errInfos.NOT_FOUND), err
	}

	// 3. 存入快取
	_ = s.cacheService.SetCenterSettings(ctx, centerID, &center.Settings)

	return &center.Settings, nil, nil
}

// UpdateCenterSettings 更新中心設置並清除快取
func (s *CenterService) UpdateCenterSettings(ctx context.Context, centerID uint, adminID uint, settings *models.CenterSettings) (*models.Center, *errInfos.Res, error) {
	// 更新資料庫
	center, err := s.centerRepo.GetByID(ctx, centerID)
	if err != nil {
		return nil, s.app.Err.New(errInfos.NOT_FOUND), err
	}

	center.Settings = *settings

	err = s.centerRepo.Update(ctx, center)
	if err != nil {
		return nil, s.app.Err.New(errInfos.SQL_ERROR), err
	}

	// 清除快取（Cache Invalidation）
	_ = s.cacheService.InvalidateCenterSettings(ctx, centerID)
	_ = s.cacheService.InvalidateCenterBasic(ctx, centerID)

	// 記錄稽核日誌
	s.auditLogRepo.Create(ctx, models.AuditLog{
		CenterID:   centerID,
		ActorType:  "ADMIN",
		ActorID:    adminID,
		Action:     "UPDATE_CENTER_SETTINGS",
		TargetType: "Center",
		TargetID:   centerID,
		Payload: models.AuditPayload{
			Before: center.Settings,
			After:  *settings,
		},
	})

	return &center, nil, nil
}

// ToggleCenterActive 切換中心啟用狀態並清除快取
// 注意：此功能需要 Center 模型包含 IsActive 欄位，目前暫時註解
func (s *CenterService) ToggleCenterActive(ctx context.Context, centerID uint, adminID uint, isActive bool) (*errInfos.Res, error) {
	// TODO: 當 Center 模型新增 IsActive 欄位後，實作此功能
	// 目前暫時跳過，不執行任何操作
	logger.GetLogger().ForComponent("CenterService").Warnw(
		"ToggleCenterActive not implemented",
		"center_id", centerID,
		"admin_id", adminID,
		"is_active", isActive,
		"reason", "Center model missing IsActive field",
	)
	return nil, nil
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
