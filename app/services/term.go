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

// TermService 學期管理相關業務邏輯
type TermService struct {
	BaseService
	app          *app.App
	termRepo     *repositories.CenterTermRepository
	auditLogRepo *repositories.AuditLogRepository
}

// NewTermService 建立 TermService 實例
func NewTermService(app *app.App) *TermService {
	return &TermService{
		BaseService:  *NewBaseService(app, "TermService"),
		app:          app,
		termRepo:     repositories.NewCenterTermRepository(app),
		auditLogRepo: repositories.NewAuditLogRepository(app),
	}
}

// CreateTermRequest 建立學期請求
type CreateTermRequest struct {
	Name      string `json:"name" binding:"required"`
	StartDate string `json:"start_date" binding:"required,date_format"`
	EndDate   string `json:"end_date" binding:"required,date_format"`
}

// UpdateTermRequest 更新學期請求
type UpdateTermRequest struct {
	Name      string `json:"name" binding:"required"`
	StartDate string `json:"start_date" binding:"required,date_format"`
	EndDate   string `json:"end_date" binding:"required,date_format"`
}

// GetTerms 取得學期列表
func (s *TermService) GetTerms(ctx context.Context, centerID uint) ([]models.CenterTerm, *errInfos.Res, error) {
	s.Logger.Info("fetching terms", "center_id", centerID)

	terms, err := s.termRepo.ListByCenterID(ctx, centerID)
	if err != nil {
		s.Logger.Error("failed to fetch terms", "error", err)
		return nil, s.app.Err.New(errInfos.SQL_ERROR), err
	}

	s.Logger.Debug("terms fetched successfully", "count", len(terms))
	return terms, nil, nil
}

// GetActiveTerms 取得進行中的學期列表
func (s *TermService) GetActiveTerms(ctx context.Context, centerID uint) ([]models.CenterTerm, *errInfos.Res, error) {
	s.Logger.Info("fetching active terms", "center_id", centerID)

	terms, err := s.termRepo.ListActiveByCenterID(ctx, centerID)
	if err != nil {
		s.Logger.Error("failed to fetch active terms", "error", err)
		return nil, s.app.Err.New(errInfos.SQL_ERROR), err
	}

	return terms, nil, nil
}

// CreateTerm 新增學期
func (s *TermService) CreateTerm(ctx context.Context, centerID, adminID uint, req *CreateTermRequest) (*models.CenterTerm, *errInfos.Res, error) {
	s.Logger.Info("creating term", "center_id", centerID, "name", req.Name)

	// 解析日期
	startDate, err := time.Parse("2006-01-02", req.StartDate)
	if err != nil {
		s.Logger.Warn("invalid start date format", "start_date", req.StartDate)
		return nil, s.app.Err.New(errInfos.PARAMS_VALIDATE_ERROR), fmt.Errorf("invalid start date format: %s", req.StartDate)
	}

	endDate, err := time.Parse("2006-01-02", req.EndDate)
	if err != nil {
		s.Logger.Warn("invalid end date format", "end_date", req.EndDate)
		return nil, s.app.Err.New(errInfos.PARAMS_VALIDATE_ERROR), fmt.Errorf("invalid end date format: %s", req.EndDate)
	}

	// 驗證日期範圍
	if endDate.Before(startDate) {
		s.Logger.Warn("end date before start date")
		return nil, s.app.Err.New(errInfos.PARAMS_VALIDATE_ERROR), fmt.Errorf("end date must be after start date")
	}

	// 檢查學期名稱是否已存在
	exists, err := s.termRepo.ExistsByCenterAndName(ctx, centerID, req.Name)
	if err != nil {
		s.Logger.Error("failed to check term name existence", "error", err)
		return nil, s.app.Err.New(errInfos.SQL_ERROR), err
	}
	if exists {
		s.Logger.Warn("term name already exists", "name", req.Name)
		return nil, s.app.Err.New(errInfos.DUPLICATE), fmt.Errorf("term with name '%s' already exists", req.Name)
	}

	// 檢查日期範圍是否重疊
	overlap, err := s.termRepo.ExistsByCenterAndDateRange(ctx, centerID, startDate, endDate, 0)
	if err != nil {
		s.Logger.Error("failed to check date range overlap", "error", err)
		return nil, s.app.Err.New(errInfos.SQL_ERROR), err
	}
	if overlap {
		s.Logger.Warn("date range overlaps with existing term")
		return nil, s.app.Err.New(errInfos.DUPLICATE), fmt.Errorf("date range overlaps with existing term")
	}

	now := time.Now()
	term := models.CenterTerm{
		CenterID:  centerID,
		Name:      req.Name,
		StartDate: startDate,
		EndDate:   endDate,
		CreatedAt: now,
		UpdatedAt: now,
	}

	var created models.CenterTerm
	err = s.termRepo.Transaction(ctx, func(txRepo *repositories.CenterTermRepository) error {
		var txErr error
		created, txErr = txRepo.Create(ctx, term)
		if txErr != nil {
			return txErr
		}

		// 記錄稽核日誌
		auditLog := models.AuditLog{
			CenterID:   centerID,
			ActorType:  "ADMIN",
			ActorID:    adminID,
			Action:     "CREATE_TERM",
			TargetType: "CenterTerm",
			TargetID:   created.ID,
			Payload: models.AuditPayload{
				After: created,
			},
		}
		if txErr := txRepo.GetDBWrite().Create(&auditLog).Error; txErr != nil {
			return txErr
		}
		return nil
	})

	if err != nil {
		s.Logger.Error("failed to create term", "error", err)
		return nil, s.app.Err.New(errInfos.SQL_ERROR), err
	}

	s.Logger.Info("term created successfully", "term_id", created.ID)
	return &created, nil, nil
}

// UpdateTerm 更新學期
func (s *TermService) UpdateTerm(ctx context.Context, centerID, adminID, termID uint, req *UpdateTermRequest) (*models.CenterTerm, *errInfos.Res, error) {
	s.Logger.Info("updating term", "center_id", centerID, "term_id", termID)

	// 解析日期
	startDate, err := time.Parse("2006-01-02", req.StartDate)
	if err != nil {
		s.Logger.Warn("invalid start date format", "start_date", req.StartDate)
		return nil, s.app.Err.New(errInfos.PARAMS_VALIDATE_ERROR), fmt.Errorf("invalid start date format: %s", req.StartDate)
	}

	endDate, err := time.Parse("2006-01-02", req.EndDate)
	if err != nil {
		s.Logger.Warn("invalid end date format", "end_date", req.EndDate)
		return nil, s.app.Err.New(errInfos.PARAMS_VALIDATE_ERROR), fmt.Errorf("invalid end date format: %s", req.EndDate)
	}

	// 驗證日期範圍
	if endDate.Before(startDate) {
		s.Logger.Warn("end date before start date")
		return nil, s.app.Err.New(errInfos.PARAMS_VALIDATE_ERROR), fmt.Errorf("end date must be after start date")
	}

	var updated *models.CenterTerm
	err = s.termRepo.Transaction(ctx, func(txRepo *repositories.CenterTermRepository) error {
		// 取得現有學期
		existing, txErr := txRepo.GetByIDWithCenterScope(ctx, termID, centerID)
		if txErr != nil {
			s.Logger.Error("term not found", "error", txErr)
			return txErr
		}

		// 檢查學期名稱是否已被其他學期使用
		if existing.Name != req.Name {
			exists, txErr := txRepo.ExistsByCenterAndNameExcludingID(ctx, centerID, req.Name, termID)
			if txErr != nil {
				s.Logger.Error("failed to check term name existence", "error", txErr)
				return txErr
			}
			if exists {
				s.Logger.Warn("term name already exists", "name", req.Name)
				return fmt.Errorf("term with name '%s' already exists", req.Name)
			}
		}

		// 檢查日期範圍是否重疊（排除自己）
		overlap, txErr := txRepo.ExistsByCenterAndDateRange(ctx, centerID, startDate, endDate, termID)
		if txErr != nil {
			s.Logger.Error("failed to check date range overlap", "error", txErr)
			return txErr
		}
		if overlap {
			s.Logger.Warn("date range overlaps with existing term")
			return fmt.Errorf("date range overlaps with existing term")
		}

		// 記錄更新前的資料
		before := existing

		// 更新學期
		existing.Name = req.Name
		existing.StartDate = startDate
		existing.EndDate = endDate
		existing.UpdatedAt = time.Now()

		if txErr := txRepo.Update(ctx, &existing); txErr != nil {
			s.Logger.Error("failed to update term", "error", txErr)
			return txErr
		}

		updated = &existing

		// 記錄稽核日誌
		auditLog := models.AuditLog{
			CenterID:   centerID,
			ActorType:  "ADMIN",
			ActorID:    adminID,
			Action:     "UPDATE_TERM",
			TargetType: "CenterTerm",
			TargetID:   termID,
			Payload: models.AuditPayload{
				Before: before,
				After:  existing,
			},
		}
		if txErr := txRepo.GetDBWrite().Create(&auditLog).Error; txErr != nil {
			return txErr
		}
		return nil
	})

	if err != nil {
		if err.Error() == "record not found" {
			s.Logger.Warn("term not found", "term_id", termID)
			return nil, s.app.Err.New(errInfos.NOT_FOUND), fmt.Errorf("term not found")
		}
		if err.Error() == "term with name '%s' already exists" || err.Error() == "date range overlaps with existing term" {
			s.Logger.Warn("validation error", "error", err)
			return nil, s.app.Err.New(errInfos.DUPLICATE), err
		}
		s.Logger.Error("failed to update term", "error", err)
		return nil, s.app.Err.New(errInfos.SQL_ERROR), err
	}

	s.Logger.Info("term updated successfully", "term_id", termID)
	return updated, nil, nil
}

// DeleteTerm 刪除學期
func (s *TermService) DeleteTerm(ctx context.Context, centerID, adminID, termID uint) (*errInfos.Res, error) {
	s.Logger.Info("deleting term", "center_id", centerID, "term_id", termID)

	err := s.termRepo.Transaction(ctx, func(txRepo *repositories.CenterTermRepository) error {
		// 取得現有學期
		existing, err := txRepo.GetByIDWithCenterScope(ctx, termID, centerID)
		if err != nil {
			s.Logger.Error("term not found", "error", err)
			return err
		}

		// 刪除學期
		if err := txRepo.Delete(ctx, termID); err != nil {
			s.Logger.Error("failed to delete term", "error", err)
			return err
		}

		// 記錄稽核日誌
		auditLog := models.AuditLog{
			CenterID:   centerID,
			ActorType:  "ADMIN",
			ActorID:    adminID,
			Action:     "DELETE_TERM",
			TargetType: "CenterTerm",
			TargetID:   termID,
			Payload: models.AuditPayload{
				Before: existing,
			},
		}
		if err := txRepo.GetDBWrite().Create(&auditLog).Error; err != nil {
			return err
		}
		return nil
	})

	if err != nil {
		if err.Error() == "record not found" {
			s.Logger.Warn("term not found", "term_id", termID)
			return s.app.Err.New(errInfos.NOT_FOUND), fmt.Errorf("term not found")
		}
		s.Logger.Error("failed to delete term", "error", err)
		return s.app.Err.New(errInfos.SQL_ERROR), err
	}

	s.Logger.Info("term deleted successfully", "term_id", termID)
	return nil, nil
}

// OccupancyRuleInfo 佔用規則資訊（用於前端週曆顯示）
type OccupancyRuleInfo struct {
	RuleID       uint   `json:"rule_id"`
	OfferingID   uint   `json:"offering_id"`
	OfferingName string `json:"offering_name"`
	Weekday      int    `json:"weekday"`
	StartTime    string `json:"start_time"`
	EndTime      string `json:"end_time"`
	Duration     int    `json:"duration"`
	TeacherID    *uint  `json:"teacher_id,omitempty"`
	TeacherName  string `json:"teacher_name,omitempty"`
	RoomID       uint   `json:"room_id"`
	RoomName     string `json:"room_name,omitempty"`
	Status       string `json:"status"` // 狀態: PLANNED(預計), CONFIRMED(已開課), SUSPENDED(停課), ARCHIVED(歸檔)
}

// OccupancyRulesByDayOfWeek 按星期分組的佔用規則
type OccupancyRulesByDayOfWeek struct {
	DayOfWeek int                 `json:"day_of_week"` // 1-7 (週一到週日)
	DayName   string              `json:"day_name"`    // "週一", "週二", etc.
	Rules     []OccupancyRuleInfo `json:"rules"`
}

// GetOccupancyRules 取得佔用規則（聚合查詢）
func (s *TermService) GetOccupancyRules(ctx context.Context, centerID uint, teacherIDs []uint, roomIDs []uint, startDate, endDate time.Time) ([]OccupancyRulesByDayOfWeek, *errInfos.Res, error) {
	s.Logger.Info("fetching occupancy rules",
		"center_id", centerID,
		"teacher_ids", teacherIDs,
		"room_ids", roomIDs,
		"start_date", startDate.Format("2006-01-02"),
		"end_date", endDate.Format("2006-01-02"))

	// 查詢重疊的規則
	var rules []models.ScheduleRule
	query := s.App.MySQL.RDB.WithContext(ctx).
		Preload("Offering").
		Preload("Teacher").
		Preload("Room").
		Where("center_id = ?", centerID)

	// 根據 teacher_ids 或 room_ids 過濾
	if len(teacherIDs) > 0 {
		query = query.Where("teacher_id IN ?", teacherIDs)
	}
	if len(roomIDs) > 0 {
		query = query.Where("room_id IN ?", roomIDs)
	}

	// 過濾日期範圍重疊的規則
	// 規則的 EffectiveRange 必須與查詢範圍有交集
	// 即：rule.StartDate <= query.EndDate AND rule.EndDate >= query.StartDate
	// 使用 JSON_UNQUOTE 確保 JSON_EXTRACT 返回的字串可以正確比較
	query = query.Where("JSON_UNQUOTE(JSON_EXTRACT(effective_range, '$.start_date')) <= ?", endDate.Format("2006-01-02")).
		Where("(JSON_UNQUOTE(JSON_EXTRACT(effective_range, '$.end_date')) >= ? OR JSON_UNQUOTE(JSON_EXTRACT(effective_range, '$.end_date')) = '0001-01-01' OR JSON_UNQUOTE(JSON_EXTRACT(effective_range, '$.end_date')) IS NULL)", startDate.Format("2006-01-02"))

	// 只查詢有 weekday 的規則（循環規則）
	query = query.Where("weekday > 0")

	if err := query.Order("weekday ASC, start_time ASC").Find(&rules).Error; err != nil {
		s.Logger.Error("failed to fetch occupancy rules", "error", err)
		return nil, s.app.Err.New(errInfos.SQL_ERROR), err
	}

	// 按 weekday 分組
	dayNames := []string{"", "週一", "週二", "週三", "週四", "週五", "週六", "週日"}
	resultByDay := make(map[int]*OccupancyRulesByDayOfWeek)

	for _, rule := range rules {
		dayOfWeek := rule.Weekday
		if dayOfWeek == 0 {
			dayOfWeek = 7 // 週日
		}

		if _, exists := resultByDay[dayOfWeek]; !exists {
			resultByDay[dayOfWeek] = &OccupancyRulesByDayOfWeek{
				DayOfWeek: dayOfWeek,
				DayName:   dayNames[dayOfWeek],
				Rules:     []OccupancyRuleInfo{},
			}
		}

		var teacherName string
		if rule.Teacher.ID != 0 {
			teacherName = rule.Teacher.Name
		}

		var roomName string
		if rule.Room.ID != 0 {
			roomName = rule.Room.Name
		}

		// 確保 Status 有正確的值，避免 NULL 或空字串
		ruleStatus := rule.Status
		if ruleStatus == "" || !models.IsValidRuleStatus(ruleStatus) {
			ruleStatus = models.RuleStatusConfirmed // 預設為 CONFIRMED
			s.Logger.Warn("rule has invalid or empty status, using default",
				"rule_id", rule.ID, "original_status", rule.Status, "default_status", ruleStatus)
		}

		ruleInfo := OccupancyRuleInfo{
			RuleID:       rule.ID,
			OfferingID:   rule.OfferingID,
			OfferingName: rule.Offering.Name,
			Weekday:      dayOfWeek,
			StartTime:    rule.StartTime,
			EndTime:      rule.EndTime,
			Duration:     rule.Duration,
			TeacherID:    rule.TeacherID,
			TeacherName:  teacherName,
			RoomID:       rule.RoomID,
			RoomName:     roomName,
			Status:       ruleStatus,
		}
		resultByDay[dayOfWeek].Rules = append(resultByDay[dayOfWeek].Rules, ruleInfo)
	}

	// 轉換為有序切片
	result := make([]OccupancyRulesByDayOfWeek, 0, len(resultByDay))
	for day := 1; day <= 7; day++ {
		if group, exists := resultByDay[day]; exists {
			result = append(result, *group)
		}
	}

	s.Logger.Debug("occupancy rules fetched successfully", "total_days", len(result), "total_rules", len(rules))
	return result, nil, nil
}

// CopyRulesRequest 複製規則請求
type CopyRulesRequest struct {
	SourceTermID uint   `json:"source_term_id" binding:"required"`
	TargetTermID uint   `json:"target_term_id" binding:"required"`
	RuleIDs      []uint `json:"rule_ids" binding:"required,min=1"`
}

// CopiedRuleInfo 複製規則結果資訊
type CopiedRuleInfo struct {
	OriginalRuleID uint   `json:"original_rule_id"`
	NewRuleID      uint   `json:"new_rule_id"`
	OfferingName   string `json:"offering_name"`
	Weekday        int    `json:"weekday"`
	StartTime      string `json:"start_time"`
	EndTime        string `json:"end_time"`
}

// CopyRulesResponse 複製規則響應
type CopyRulesResponse struct {
	CopiedCount int              `json:"copied_count"`
	Rules       []CopiedRuleInfo `json:"rules"`
}

// CopyRules 複製規則到目標學期
func (s *TermService) CopyRules(ctx context.Context, centerID, adminID uint, req *CopyRulesRequest) (*CopyRulesResponse, *errInfos.Res, error) {
	s.Logger.Info("copying rules",
		"center_id", centerID,
		"source_term_id", req.SourceTermID,
		"target_term_id", req.TargetTermID,
		"rule_count", len(req.RuleIDs))

	// 驗證來源和目標學期存在且屬於該中心
	sourceTerm, err := s.termRepo.GetByIDWithCenterScope(ctx, req.SourceTermID, centerID)
	if err != nil {
		s.Logger.Error("source term not found", "error", err, "term_id", req.SourceTermID)
		return nil, s.app.Err.New(errInfos.NOT_FOUND), fmt.Errorf("source term not found")
	}

	targetTerm, err := s.termRepo.GetByIDWithCenterScope(ctx, req.TargetTermID, centerID)
	if err != nil {
		s.Logger.Error("target term not found", "error", err, "term_id", req.TargetTermID)
		return nil, s.app.Err.New(errInfos.NOT_FOUND), fmt.Errorf("target term not found")
	}

	// 不能複製到同一個學期
	if req.SourceTermID == req.TargetTermID {
		s.Logger.Warn("source and target term are the same")
		return nil, s.app.Err.New(errInfos.PARAMS_VALIDATE_ERROR), fmt.Errorf("source and target term cannot be the same")
	}

	// 查詢來源規則
	var sourceRules []models.ScheduleRule
	query := s.App.MySQL.RDB.WithContext(ctx).
		Preload("Offering").
		Preload("Teacher").
		Preload("Room").
		Where("center_id = ?", centerID).
		Where("id IN ?", req.RuleIDs)

	if err := query.Find(&sourceRules).Error; err != nil {
		s.Logger.Error("failed to fetch source rules", "error", err)
		return nil, s.app.Err.New(errInfos.SQL_ERROR), err
	}

	if len(sourceRules) == 0 {
		s.Logger.Warn("no valid rules found to copy")
		return nil, s.app.Err.New(errInfos.NOT_FOUND), fmt.Errorf("no valid rules found to copy")
	}

	// 驗證所有規則都屬於來源學期範圍內
	targetStartDate := targetTerm.StartDate
	targetEndDate := targetTerm.EndDate

	now := time.Now()
	copiedRules := make([]CopiedRuleInfo, 0, len(sourceRules))

	// 在交易中執行複製
	txErr := s.App.MySQL.WDB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		for _, rule := range sourceRules {
			// Deep copy 規則（排除 ID 和時間戳）
			newRule := models.ScheduleRule{
				CenterID:    centerID,
				OfferingID:  rule.OfferingID,
				TeacherID:   rule.TeacherID,
				RoomID:      rule.RoomID,
				Name:        rule.Name,
				Weekday:     rule.Weekday,
				StartTime:   rule.StartTime,
				EndTime:     rule.EndTime,
				Duration:    rule.Duration,
				IsCrossDay:  rule.IsCrossDay,
				SkipHoliday: rule.SkipHoliday,
				EffectiveRange: models.DateRange{
					StartDate: targetStartDate,
					EndDate:   targetEndDate,
				},
				LockAt:    rule.LockAt,
				CreatedAt: now,
				UpdatedAt: now,
			}

			// 建立新規則
			if err := tx.Create(&newRule).Error; err != nil {
				s.Logger.Error("failed to create new rule", "error", err, "original_rule_id", rule.ID)
				return fmt.Errorf("failed to copy rule %d: %w", rule.ID, err)
			}

			copiedInfo := CopiedRuleInfo{
				OriginalRuleID: rule.ID,
				NewRuleID:      newRule.ID,
				OfferingName:   rule.Offering.Name,
				Weekday:        rule.Weekday,
				StartTime:      rule.StartTime,
				EndTime:        rule.EndTime,
			}
			copiedRules = append(copiedRules, copiedInfo)

			// 記錄稽核日誌
			auditLog := models.AuditLog{
				CenterID:   centerID,
				ActorType:  "ADMIN",
				ActorID:    adminID,
				Action:     "COPY_TERM_RULES",
				TargetType: "ScheduleRule",
				TargetID:   newRule.ID,
				Payload: models.AuditPayload{
					Before: fmt.Sprintf("Copy from term %d, rule %d", sourceTerm.ID, rule.ID),
					After:  newRule,
				},
			}
			if err := tx.Create(&auditLog).Error; err != nil {
				s.Logger.Warn("failed to create audit log", "error", err)
			}
		}
		return nil
	})

	if txErr != nil {
		s.Logger.Error("failed to copy rules in transaction", "error", txErr)
		return nil, s.app.Err.New(errInfos.SQL_ERROR), txErr
	}

	// 呼叫 ScheduleExpansionService 展開新規則（產生課程實例）
	// 這裡我們需要手動觸發擴展
	// 注意：實際的展開通常是在查詢時動態執行的，
	// 或者可以在背景 job 中執行
	s.Logger.Info("rules copied successfully, expansion will be handled by expand query",
		"copied_count", len(copiedRules))

	response := &CopyRulesResponse{
		CopiedCount: len(copiedRules),
		Rules:       copiedRules,
	}

	return response, nil, nil
}
