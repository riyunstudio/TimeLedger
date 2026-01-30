package services

import (
	"context"
	"fmt"
	"time"
	"timeLedger/app"
	"timeLedger/app/models"
	"timeLedger/app/repositories"
	"timeLedger/global"
	"timeLedger/global/errInfos"

	"gorm.io/gorm"
)

// OfferingService 班別相關業務邏輯
type OfferingService struct {
	BaseService
	app          *app.App
	offeringRepo *repositories.OfferingRepository
	courseRepo   *repositories.CourseRepository
	auditLogRepo *repositories.AuditLogRepository
}

// NewOfferingService 建立 OfferingService 實例
func NewOfferingService(appInstance *app.App) *OfferingService {
	return &OfferingService{
		app:          appInstance,
		offeringRepo: repositories.NewOfferingRepository(appInstance),
		courseRepo:   repositories.NewCourseRepository(appInstance),
		auditLogRepo: repositories.NewAuditLogRepository(appInstance),
	}
}

// ListOfferingsInput 查詢班別列表的輸入參數
type ListOfferingsInput struct {
	CenterID uint
	Page     int
	Limit    int
}

// ListOfferingsOutput 查詢班別列表的輸出資料
type ListOfferingsOutput struct {
	Offerings  []models.Offering
	Pagination global.Pagination
}

// ListOfferings 取得班別列表
func (s *OfferingService) ListOfferings(ctx context.Context, input *ListOfferingsInput) (*ListOfferingsOutput, *errInfos.Res, error) {
	offerings, total, err := s.offeringRepo.ListByCenterIDPaginated(
		ctx,
		input.CenterID,
		input.Page,
		input.Limit,
	)
	if err != nil {
		return nil, s.app.Err.New(errInfos.SQL_ERROR), err
	}

	page := int64(input.Page)
	limit := int64(input.Limit)

	return &ListOfferingsOutput{
		Offerings:  offerings,
		Pagination: global.NewPagination(page, limit, total),
	}, nil, nil
}

// CreateOfferingInput 建立班別的輸入參數
type CreateOfferingInput struct {
	CenterID            uint
	AdminID             uint
	CourseID            uint
	DefaultRoomID       *uint
	DefaultTeacherID    *uint
	AllowBufferOverride bool
}

// CreateOffering 建立新班別
func (s *OfferingService) CreateOffering(ctx context.Context, input *CreateOfferingInput) (*models.Offering, *errInfos.Res, error) {
	// 取得課程名稱作為 offering 的名稱
	course, err := s.courseRepo.GetByID(ctx, input.CourseID)
	if err != nil {
		return nil, s.app.Err.New(errInfos.NOT_FOUND), err
	}

	var createdOffering models.Offering

	// 使用交易確保建立 offering 和稽核日誌的原子性
	txErr := s.app.MySQL.WDB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		offering := models.Offering{
			CenterID:            input.CenterID,
			CourseID:            input.CourseID,
			Name:                course.Name,
			DefaultRoomID:       input.DefaultRoomID,
			DefaultTeacherID:    input.DefaultTeacherID,
			AllowBufferOverride: input.AllowBufferOverride,
			CreatedAt:           time.Now(),
			UpdatedAt:           time.Now(),
		}

		if err := tx.Create(&offering).Error; err != nil {
			return fmt.Errorf("failed to create offering: %w", err)
		}
		createdOffering = offering

		// 在交易中記錄稽核日誌
		auditLog := models.AuditLog{
			CenterID:   input.CenterID,
			ActorType:  "ADMIN",
			ActorID:    input.AdminID,
			Action:     "OFFERING_CREATE",
			TargetType: "Offering",
			TargetID:   createdOffering.ID,
			Payload: models.AuditPayload{
				After: map[string]interface{}{
					"course_id":             input.CourseID,
					"course_name":           course.Name,
					"default_room_id":       input.DefaultRoomID,
					"default_teacher_id":    input.DefaultTeacherID,
					"allow_buffer_override": input.AllowBufferOverride,
				},
			},
		}
		if err := tx.Create(&auditLog).Error; err != nil {
			return fmt.Errorf("failed to create audit log: %w", err)
		}

		return nil
	})

	if txErr != nil {
		return nil, s.app.Err.New(errInfos.ERR_TX_FAILED), txErr
	}

	return &createdOffering, nil, nil
}

// UpdateOfferingInput 更新班別的輸入參數
type UpdateOfferingInput struct {
	CenterID            uint
	AdminID             uint
	OfferingID          uint
	Name                *string
	DefaultRoomID       *uint
	DefaultTeacherID    *uint
	AllowBufferOverride bool
}

// UpdateOffering 更新班別
func (s *OfferingService) UpdateOffering(ctx context.Context, input *UpdateOfferingInput) (*models.Offering, *errInfos.Res, error) {
	// 查詢現有 offering
	existingOffering, err := s.offeringRepo.GetByID(ctx, input.OfferingID)
	if err != nil {
		return nil, s.app.Err.New(errInfos.NOT_FOUND), err
	}

	// 驗證權限
	if existingOffering.CenterID != input.CenterID {
		return nil, s.app.Err.New(errInfos.FORBIDDEN), nil
	}

	// 使用交易確保更新 offering 和稽核日誌的原子性
	var updatedOffering models.Offering

	txErr := s.app.MySQL.WDB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// 更新字段
		existingOffering.DefaultRoomID = input.DefaultRoomID
		existingOffering.DefaultTeacherID = input.DefaultTeacherID
		existingOffering.AllowBufferOverride = input.AllowBufferOverride
		existingOffering.UpdatedAt = time.Now()

		// 如果有提供 name，則更新
		if input.Name != nil && *input.Name != "" {
			existingOffering.Name = *input.Name
		}

		if err := tx.Save(&existingOffering).Error; err != nil {
			return fmt.Errorf("failed to update offering: %w", err)
		}
		updatedOffering = existingOffering

		// 在交易中記錄稽核日誌
		auditLog := models.AuditLog{
			CenterID:   input.CenterID,
			ActorType:  "ADMIN",
			ActorID:    input.AdminID,
			Action:     "OFFERING_UPDATE",
			TargetType: "Offering",
			TargetID:   input.OfferingID,
			Payload: models.AuditPayload{
				After: map[string]interface{}{
					"default_room_id":       input.DefaultRoomID,
					"default_teacher_id":    input.DefaultTeacherID,
					"allow_buffer_override": input.AllowBufferOverride,
				},
			},
		}
		if err := tx.Create(&auditLog).Error; err != nil {
			return fmt.Errorf("failed to create audit log: %w", err)
		}

		return nil
	})

	if txErr != nil {
		return nil, s.app.Err.New(errInfos.ERR_TX_FAILED), txErr
	}

	return &updatedOffering, nil, nil
}

// DeleteOfferingInput 刪除班別的輸入參數
type DeleteOfferingInput struct {
	CenterID uint
	AdminID  uint
	ID       uint
}

// DeleteOffering 刪除班別
func (s *OfferingService) DeleteOffering(ctx context.Context, input *DeleteOfferingInput) *errInfos.Res {
	// 使用交易確保刪除 offering 和稽核日誌的原子性
	txErr := s.app.MySQL.WDB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// 先刪除 offering
		if err := tx.Table("offerings").
			Where("id = ? AND center_id = ?", input.ID, input.CenterID).
			Delete(&models.Offering{}).Error; err != nil {
			return fmt.Errorf("failed to delete offering: %w", err)
		}

		// 在交易中記錄稽核日誌
		auditLog := models.AuditLog{
			CenterID:   input.CenterID,
			ActorType:  "ADMIN",
			ActorID:    input.AdminID,
			Action:     "OFFERING_DELETE",
			TargetType: "Offering",
			TargetID:   input.ID,
			Payload: models.AuditPayload{
				After: map[string]interface{}{
					"status": "DELETED",
				},
			},
		}
		if err := tx.Create(&auditLog).Error; err != nil {
			return fmt.Errorf("failed to create audit log: %w", err)
		}

		return nil
	})

	if txErr != nil {
		return s.app.Err.New(errInfos.ERR_TX_FAILED)
	}

	return nil
}

// CopyOfferingInput 複製班別的輸入參數
type CopyOfferingInput struct {
	CenterID    uint
	AdminID     uint
	OfferingID  uint
	NewName     string
	CopyTeacher bool
}

// CopyOfferingOutput 複製班別的輸出資料
type CopyOfferingOutput struct {
	ID          uint   `json:"id"`
	Name        string `json:"name"`
	CourseID    uint   `json:"course_id"`
	RulesCopied int    `json:"rules_copied"`
}

// CopyOffering 複製班別
func (s *OfferingService) CopyOffering(ctx context.Context, input *CopyOfferingInput) (*CopyOfferingOutput, *errInfos.Res, error) {
	// 取得原始班別
	original, err := s.offeringRepo.GetByIDAndCenterID(ctx, input.OfferingID, input.CenterID)
	if err != nil {
		return nil, s.app.Err.New(errInfos.NOT_FOUND), err
	}

	// 決定新班別的老師
	newTeacherID := original.DefaultTeacherID
	if !input.CopyTeacher {
		newTeacherID = nil
	}

	var createdOffering models.Offering

	// 使用交易確保建立新 offering 和稽核日誌的原子性
	txErr := s.app.MySQL.WDB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		newOffering := models.Offering{
			CenterID:            input.CenterID,
			CourseID:            original.CourseID,
			Name:                input.NewName,
			DefaultRoomID:       original.DefaultRoomID,
			DefaultTeacherID:    newTeacherID,
			AllowBufferOverride: original.AllowBufferOverride,
			IsActive:            true,
			CreatedAt:           time.Now(),
			UpdatedAt:           time.Now(),
		}

		if err := tx.Create(&newOffering).Error; err != nil {
			return fmt.Errorf("failed to create offering: %w", err)
		}
		createdOffering = newOffering

		// 在交易中記錄稽核日誌
		auditLog := models.AuditLog{
			CenterID:   input.CenterID,
			ActorType:  "ADMIN",
			ActorID:    input.AdminID,
			Action:     "OFFERING_COPY",
			TargetType: "Offering",
			TargetID:   createdOffering.ID,
			Payload: models.AuditPayload{
				After: map[string]interface{}{
					"original_offering_id": original.ID,
					"new_offering_id":      createdOffering.ID,
					"name":                 input.NewName,
					"copy_teacher":         input.CopyTeacher,
				},
			},
		}
		if err := tx.Create(&auditLog).Error; err != nil {
			return fmt.Errorf("failed to create audit log: %w", err)
		}

		return nil
	})

	if txErr != nil {
		return nil, s.app.Err.New(errInfos.ERR_TX_FAILED), txErr
	}

	return &CopyOfferingOutput{
		ID:          createdOffering.ID,
		Name:        createdOffering.Name,
		CourseID:    createdOffering.CourseID,
		RulesCopied: 0,
	}, nil, nil
}

// ToggleOfferingActiveInput 切換班別啟用狀態的輸入參數
type ToggleOfferingActiveInput struct {
	CenterID uint
	AdminID  uint
	ID       uint
	IsActive bool
}

// ToggleOfferingActive 切換班別啟用狀態
func (s *OfferingService) ToggleOfferingActive(ctx context.Context, input *ToggleOfferingActiveInput) *errInfos.Res {
	// 使用交易確保更新狀態和稽核日誌的原子性
	txErr := s.app.MySQL.WDB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// 更新啟用狀態
		if err := tx.Table("offerings").
			Where("id = ? AND center_id = ?", input.ID, input.CenterID).
			Update("is_active", input.IsActive).Error; err != nil {
			return fmt.Errorf("failed to toggle offering active status: %w", err)
		}

		// 在交易中記錄稽核日誌
		auditLog := models.AuditLog{
			CenterID:   input.CenterID,
			ActorType:  "ADMIN",
			ActorID:    input.AdminID,
			Action:     "OFFERING_TOGGLE_ACTIVE",
			TargetType: "Offering",
			TargetID:   input.ID,
			Payload: models.AuditPayload{
				After: map[string]interface{}{
					"is_active": input.IsActive,
				},
			},
		}
		if err := tx.Create(&auditLog).Error; err != nil {
			return fmt.Errorf("failed to create audit log: %w", err)
		}

		return nil
	})

	if txErr != nil {
		return s.app.Err.New(errInfos.ERR_TX_FAILED)
	}

	return nil
}

// GetActiveOfferingsInput 取得啟用班別列表的輸入參數
type GetActiveOfferingsInput struct {
	CenterID uint
}

// GetActiveOfferings 取得啟用的班別列表
func (s *OfferingService) GetActiveOfferings(ctx context.Context, input *GetActiveOfferingsInput) ([]models.Offering, error) {
	return s.offeringRepo.ListActiveByCenterID(ctx, input.CenterID)
}
