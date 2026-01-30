package services

import (
	"context"
	"fmt"
	"time"
	"timeLedger/app"
	"timeLedger/app/models"
	"timeLedger/app/repositories"
	"timeLedger/global/errInfos"

	"gorm.io/gorm"
)

// TimetableTemplateService 課表模板相關業務邏輯
type TimetableTemplateService struct {
	BaseService
	app              *app.App
	templateRepo     *repositories.TimetableTemplateRepository
	cellRepo         *repositories.TimetableCellRepository
	scheduleRuleRepo *repositories.ScheduleRuleRepository
	auditLogRepo     *repositories.AuditLogRepository
	ruleValidator    *ScheduleRuleValidator
}

// NewTimetableTemplateService 建立 TimetableTemplateService 實例
func NewTimetableTemplateService(appInstance *app.App) *TimetableTemplateService {
	return &TimetableTemplateService{
		app:              appInstance,
		templateRepo:     repositories.NewTimetableTemplateRepository(appInstance),
		cellRepo:         repositories.NewTimetableCellRepository(appInstance),
		scheduleRuleRepo: repositories.NewScheduleRuleRepository(appInstance),
		auditLogRepo:     repositories.NewAuditLogRepository(appInstance),
		ruleValidator:    NewScheduleRuleValidator(appInstance),
	}
}

// ApplyTemplateInput 套用模板的輸入參數
type ApplyTemplateInput struct {
	TemplateID      uint
	CenterID        uint
	AdminID         uint
	OfferingID      uint
	StartDate       string
	EndDate         string
	Weekdays        []int
	Duration        int
	OverrideBuffer  bool
}

// ApplyTemplateConflictInfo 套用模板衝突資訊
type ApplyTemplateConflictInfo struct {
	Weekday      int    `json:"weekday"`
	StartTime    string `json:"start_time"`
	EndTime      string `json:"end_time"`
	ConflictType string `json:"conflict_type"`
	Message      string `json:"message"`
	RuleID       uint   `json:"rule_id,omitempty"`
	CanOverride  bool   `json:"can_override,omitempty"`
}

// ApplyTemplateOutput 套用模板的輸出資料
type ApplyTemplateOutput struct {
	RulesCreated int    `json:"rules_created"`
	TemplateName string `json:"template_name"`
}

// ApplyTemplateResult 套用模板的結果（包含衝突資訊）
type ApplyTemplateResult struct {
	Output    *ApplyTemplateOutput
	Conflicts []ApplyTemplateConflictInfo
	Valid     bool
}

// ApplyTemplate 套用課表模板
func (s *TimetableTemplateService) ApplyTemplate(ctx context.Context, input *ApplyTemplateInput) (*ApplyTemplateResult, *errInfos.Res, error) {
	// 解析日期
	startDate, err := time.Parse("2006-01-02", input.StartDate)
	if err != nil {
		return nil, s.app.Err.New(errInfos.PARAMS_VALIDATE_ERROR), err
	}
	endDate, err := time.Parse("2006-01-02", input.EndDate)
	if err != nil {
		return nil, s.app.Err.New(errInfos.PARAMS_VALIDATE_ERROR), err
	}

	// 取得模板
	template, err := s.templateRepo.GetByID(ctx, input.TemplateID)
	if err != nil {
		return nil, s.app.Err.New(errInfos.NOT_FOUND), err
	}

	// 驗證權限
	if template.CenterID != input.CenterID {
		return nil, s.app.Err.New(errInfos.FORBIDDEN), nil
	}

	// 取得模板中的 cells
	cells, err := s.cellRepo.ListByTemplateID(ctx, input.TemplateID)
	if err != nil {
		return nil, s.app.Err.New(errInfos.SQL_ERROR), err
	}

	// 使用 ScheduleRuleValidator 進行完整驗證
	validationSummary, err := s.ruleValidator.ValidateForApplyTemplate(
		ctx,
		input.CenterID,
		input.OfferingID,
		input.Weekdays,
		cells,
		input.StartDate,
		input.EndDate,
		input.OverrideBuffer,
	)
	if err != nil {
		return nil, s.app.Err.New(errInfos.SCHED_OVERLAP), err
	}

	// 分析衝突
	nonOverrideConflicts := 0
	var allConflicts []ApplyTemplateConflictInfo
	for _, conflict := range validationSummary.AllConflicts {
		allConflicts = append(allConflicts, ApplyTemplateConflictInfo{
			Weekday:      conflict.Weekday,
			StartTime:    conflict.StartTime,
			EndTime:      conflict.EndTime,
			ConflictType: conflict.ConflictType,
			Message:      conflict.Message,
			RuleID:       conflict.RuleID,
			CanOverride:  conflict.CanOverride,
		})
		if !conflict.CanOverride {
			nonOverrideConflicts++
		}
	}

	// 如果有不可覆蓋的衝突，回傳衝突資訊
	if !validationSummary.Valid && nonOverrideConflicts > 0 {
		return &ApplyTemplateResult{
			Valid:     false,
			Conflicts: allConflicts,
		}, nil, nil
	}

	// 如果只有可覆蓋的衝突，但未指定覆蓋，回傳警告
	if !validationSummary.Valid && !input.OverrideBuffer && len(allConflicts) > 0 {
		return &ApplyTemplateResult{
			Valid:     false,
			Conflicts: allConflicts,
		}, nil, nil
	}

	// 產生 Schedule Rules
	rules := s.generateScheduleRules(ctx, input.CenterID, input.OfferingID, cells, input.Weekdays, startDate, endDate)

	// 使用交易確保建立規則和稽核日誌的原子性
	var createdRules []models.ScheduleRule

	txErr := s.app.MySQL.WDB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// 在交易中建立 Schedule Rules
		if len(rules) > 0 {
			if err := tx.Table("schedule_rules").Create(&rules).Error; err != nil {
				return fmt.Errorf("failed to create schedule rules: %w", err)
			}
			createdRules = rules
		}

		// 在交易中記錄稽核日誌
		auditLog := models.AuditLog{
			CenterID:   input.CenterID,
			ActorType:  "ADMIN",
			ActorID:    input.AdminID,
			Action:     "TEMPLATE_APPLY",
			TargetType: "ScheduleRule",
			TargetID:   0,
			Payload: models.AuditPayload{
				After: map[string]interface{}{
					"template_id":   input.TemplateID,
					"offering_id":   input.OfferingID,
					"start_date":    input.StartDate,
					"end_date":      input.EndDate,
					"weekdays":      input.Weekdays,
					"rules_created": len(createdRules),
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

	return &ApplyTemplateResult{
		Valid: true,
		Output: &ApplyTemplateOutput{
			RulesCreated: len(createdRules),
			TemplateName: template.Name,
		},
		Conflicts: nil,
	}, nil, nil
}

// generateScheduleRules 產生排課規則
func (s *TimetableTemplateService) generateScheduleRules(
	ctx context.Context,
	centerID, offeringID uint,
	cells []models.TimetableCell,
	weekdays []int,
	startDate, endDate time.Time,
) []models.ScheduleRule {
	var rules []models.ScheduleRule
	for _, weekday := range weekdays {
		for _, cell := range cells {
			rule := models.ScheduleRule{
				CenterID:    centerID,
				OfferingID:  offeringID,
				TeacherID:   cell.TeacherID,
				RoomID:      *cell.RoomID,
				Weekday:     weekday,
				StartTime:   cell.StartTime,
				EndTime:     cell.EndTime,
				EffectiveRange: models.DateRange{
					StartDate: startDate,
					EndDate:   endDate,
				},
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			}
			rules = append(rules, rule)
		}
	}
	return rules
}

// ApplyTemplateValidateInput 驗證套用模板的輸入參數
type ApplyTemplateValidateInput struct {
	TemplateID     uint
	CenterID       uint
	OfferingID     uint
	Weekdays       []int
	StartDate      string
	EndDate        string
	OverrideBuffer bool
}

// ApplyTemplateValidateResult 驗證結果
type ApplyTemplateValidateResult struct {
	Valid     bool
	Conflicts []ApplyTemplateConflictInfo
}

// ValidateApplyTemplate 驗證套用模板（不實際產生規則）
func (s *TimetableTemplateService) ValidateApplyTemplate(ctx context.Context, input *ApplyTemplateValidateInput) (*ApplyTemplateValidateResult, *errInfos.Res, error) {
	// 取得模板中的 cells
	cells, err := s.cellRepo.ListByTemplateID(ctx, input.TemplateID)
	if err != nil {
		return nil, s.app.Err.New(errInfos.SQL_ERROR), err
	}

	// 驗證模板所屬中心
	template, err := s.templateRepo.GetByID(ctx, input.TemplateID)
	if err != nil {
		return nil, s.app.Err.New(errInfos.NOT_FOUND), err
	}
	if template.CenterID != input.CenterID {
		return nil, s.app.Err.New(errInfos.FORBIDDEN), nil
	}

	// 使用 ScheduleRuleValidator 進行驗證
	validationSummary, err := s.ruleValidator.ValidateForApplyTemplate(
		ctx,
		input.CenterID,
		input.OfferingID,
		input.Weekdays,
		cells,
		input.StartDate,
		input.EndDate,
		input.OverrideBuffer,
	)
	if err != nil {
		return nil, s.app.Err.New(errInfos.SCHED_OVERLAP), err
	}

	// 轉換衝突格式
	var conflicts []ApplyTemplateConflictInfo
	nonOverrideConflicts := 0
	for _, conflict := range validationSummary.AllConflicts {
		conflicts = append(conflicts, ApplyTemplateConflictInfo{
			Weekday:      conflict.Weekday,
			StartTime:    conflict.StartTime,
			EndTime:      conflict.EndTime,
			ConflictType: conflict.ConflictType,
			Message:      conflict.Message,
			RuleID:       conflict.RuleID,
			CanOverride:  conflict.CanOverride,
		})
		if !conflict.CanOverride {
			nonOverrideConflicts++
		}
	}

	valid := validationSummary.Valid || (len(conflicts) > 0 && nonOverrideConflicts == 0 && input.OverrideBuffer)

	return &ApplyTemplateValidateResult{
		Valid:     valid,
		Conflicts: conflicts,
	}, nil, nil
}
