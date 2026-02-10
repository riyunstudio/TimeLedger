package services

import (
	"context"
	"fmt"
	"time"

	"timeLedger/app"
	"timeLedger/app/models"

	"gorm.io/gorm"
)

// TeacherMergeService 教師合併服務

// TeacherMergeService 教師合併服務
type TeacherMergeService struct {
	BaseService
}

// NewTeacherMergeService 建立教師合併服務實例
func NewTeacherMergeService(app *app.App) *TeacherMergeService {
	baseSvc := NewBaseService(app, "TeacherMergeService")
	return &TeacherMergeService{
		BaseService: *baseSvc,
	}
}

// MergeTeacher 合併兩位教師的資料
// 將 sourceID 教師的所有關聯資料遷移到 targetID，然後軟刪除 sourceID
// 參數：
//   - ctx: 上下文
//   - sourceID: 來源教師 ID（將被合併的教師）
//   - targetID: 目標教師 ID（合併後的主教師）
//   - centerID: 中心 ID，用於篩選會籍和筆記
//
// 返回錯誤：
//   - 若驗證失敗，返回相應錯誤
//   - 若資料庫操作失敗，返回 SQL_ERROR
func (s *TeacherMergeService) MergeTeacher(ctx context.Context, sourceID, targetID, centerID uint) error {
	// 驗證參數
	if sourceID == targetID {
		return fmt.Errorf("來源教師 ID 與目標教師 ID 相同，無法合併")
	}

	if sourceID == 0 || targetID == 0 {
		return fmt.Errorf("教師 ID 不能為零")
	}

	s.Logger.Info("開始合併教師資料",
		"source_teacher_id", sourceID,
		"target_teacher_id", targetID,
		"center_id", centerID)

	// 驗證教師是否存在
	sourceTeacher, err := s.getTeacherByID(sourceID)
	if err != nil {
		s.Logger.Error("取得來源教師失敗", "error", err, "source_teacher_id", sourceID)
		return fmt.Errorf("來源教師不存在: %w", err)
	}

	targetTeacher, err := s.getTeacherByID(targetID)
	if err != nil {
		s.Logger.Error("取得目標教師失敗", "error", err, "target_teacher_id", targetID)
		return fmt.Errorf("目標教師不存在: %w", err)
	}

	s.Logger.Info("教師資料驗證成功",
		"source_name", sourceTeacher.Name,
		"target_name", targetTeacher.Name)

	// 使用交易進行所有操作
	return s.App.MySQL.WDB.Transaction(func(tx *gorm.DB) error {
		// 1. 遷移 schedule_rules（TeacherID）
		if err := s.migrateScheduleRules(tx, sourceID, targetID, centerID); err != nil {
			return err
		}

		// 2. 遷移 schedule_exceptions（NewTeacherID）
		if err := s.migrateScheduleExceptions(tx, sourceID, targetID, centerID); err != nil {
			return err
		}

		// 3. 遷移 session_notes（TeacherID）
		if err := s.migrateSessionNotes(tx, sourceID, targetID); err != nil {
			return err
		}

		// 4. 遷移 personal_events（TeacherID）
		if err := s.migratePersonalEvents(tx, sourceID, targetID); err != nil {
			return err
		}

		// 5. 遷移 offerings（DefaultTeacherID）
		if err := s.migrateOfferings(tx, sourceID, targetID, centerID); err != nil {
			return err
		}

		// 6. 遷移 timetable_cells（TeacherID）
		if err := s.migrateTimetableCells(tx, sourceID, targetID); err != nil {
			return err
		}

		// 7. 遷移 teacher_skills（TeacherID）
		if err := s.migrateTeacherSkills(tx, sourceID, targetID); err != nil {
			return err
		}

		// 8. 遷移 teacher_certificates（TeacherID）
		if err := s.migrateTeacherCertificates(tx, sourceID, targetID); err != nil {
			return err
		}

		// 9. 遷移 teacher_personal_hashtags（TeacherID）
		if err := s.migrateTeacherPersonalHashtags(tx, sourceID, targetID); err != nil {
			return err
		}

		// 10. 處理中心會籍
		if err := s.handleCenterMembership(tx, sourceID, targetID, centerID); err != nil {
			return err
		}

		// 11. 處理教師筆記
		if err := s.handleTeacherNotes(tx, sourceID, targetID, centerID); err != nil {
			return err
		}

		// 12. 軟刪除來源教師
		if err := s.softDeleteTeacher(tx, sourceID); err != nil {
			return err
		}

		s.Logger.Info("教師合併完成",
			"source_teacher_id", sourceID,
			"target_teacher_id", targetID)

		return nil
	})
}

// getTeacherByID 根據 ID 取得教師
func (s *TeacherMergeService) getTeacherByID(id uint) (*models.Teacher, error) {
	var teacher models.Teacher
	result := s.App.MySQL.RDB.WithContext(context.Background()).
		Where("id = ?", id).
		First(&teacher)

	if result.Error != nil {
		return nil, result.Error
	}

	return &teacher, nil
}

// migrateScheduleRules 遷移排課規則
func (s *TeacherMergeService) migrateScheduleRules(tx *gorm.DB, sourceID, targetID, centerID uint) error {
	result := tx.Model(&models.ScheduleRule{}).
		Where("teacher_id = ? AND center_id = ?", sourceID, centerID).
		Update("teacher_id", targetID)

	if result.Error != nil {
		s.Logger.Error("遷移排課規則失敗", "error", result.Error)
		return fmt.Errorf("遷移排課規則失敗: %w", result.Error)
	}

	if result.RowsAffected > 0 {
		s.Logger.Info("已遷移排課規則",
			"count", result.RowsAffected,
			"source_teacher_id", sourceID,
			"target_teacher_id", targetID)
	}

	return nil
}

// migrateScheduleExceptions 遷移例外記錄（包含 NewTeacherID）
func (s *TeacherMergeService) migrateScheduleExceptions(tx *gorm.DB, sourceID, targetID, centerID uint) error {
	// 更新原始教師（Rule 的 teacher_id）
	// 使用子查詢來過濾屬於該教師的規則
	result1 := tx.Model(&models.ScheduleException{}).
		Where("center_id = ?", centerID).
		Where("rule_id IN (SELECT id FROM schedule_rules WHERE teacher_id = ?)", sourceID).
		Where("(new_teacher_id IS NULL OR new_teacher_id = ?)", sourceID).
		Update("new_teacher_id", targetID)

	if result1.Error != nil {
		s.Logger.Error("遷移例外記錄（原始教師）失敗", "error", result1.Error)
		return fmt.Errorf("遷移例外記錄失敗: %w", result1.Error)
	}

	// 更新新教師（NewTeacherID）
	result2 := tx.Model(&models.ScheduleException{}).
		Where("new_teacher_id = ? AND center_id = ?", sourceID, centerID).
		Update("new_teacher_id", targetID)

	if result2.Error != nil {
		s.Logger.Error("遷移例外記錄（新教師）失敗", "error", result2.Error)
		return fmt.Errorf("遷移例外記錄失敗: %w", result2.Error)
	}

	totalAffected := result1.RowsAffected + result2.RowsAffected
	if totalAffected > 0 {
		s.Logger.Info("已遷移例外記錄",
			"count", totalAffected,
			"source_teacher_id", sourceID,
			"target_teacher_id", targetID)
	}

	return nil
}

// migrateSessionNotes 遷移課程筆記
func (s *TeacherMergeService) migrateSessionNotes(tx *gorm.DB, sourceID, targetID uint) error {
	result := tx.Model(&models.SessionNote{}).
		Where("teacher_id = ?", sourceID).
		Update("teacher_id", targetID)

	if result.Error != nil {
		s.Logger.Error("遷移課程筆記失敗", "error", result.Error)
		return fmt.Errorf("遷移課程筆記失敗: %w", result.Error)
	}

	if result.RowsAffected > 0 {
		s.Logger.Info("已遷移課程筆記",
			"count", result.RowsAffected,
			"source_teacher_id", sourceID,
			"target_teacher_id", targetID)
	}

	return nil
}

// migratePersonalEvents 遷移私人行程
func (s *TeacherMergeService) migratePersonalEvents(tx *gorm.DB, sourceID, targetID uint) error {
	result := tx.Model(&models.PersonalEvent{}).
		Where("teacher_id = ?", sourceID).
		Update("teacher_id", targetID)

	if result.Error != nil {
		s.Logger.Error("遷移私人行程失敗", "error", result.Error)
		return fmt.Errorf("遷移私人行程失敗: %w", result.Error)
	}

	if result.RowsAffected > 0 {
		s.Logger.Info("已遷移私人行程",
			"count", result.RowsAffected,
			"source_teacher_id", sourceID,
			"target_teacher_id", targetID)
	}

	return nil
}

// migrateOfferings 遷移課程項目（DefaultTeacherID）
func (s *TeacherMergeService) migrateOfferings(tx *gorm.DB, sourceID, targetID, centerID uint) error {
	result := tx.Model(&models.Offering{}).
		Where("default_teacher_id = ? AND center_id = ?", sourceID, centerID).
		Update("default_teacher_id", targetID)

	if result.Error != nil {
		s.Logger.Error("遷移課程項目失敗", "error", result.Error)
		return fmt.Errorf("遷移課程項目失敗: %w", result.Error)
	}

	if result.RowsAffected > 0 {
		s.Logger.Info("已遷移課程項目",
			"count", result.RowsAffected,
			"source_teacher_id", sourceID,
			"target_teacher_id", targetID)
	}

	return nil
}

// migrateTimetableCells 遷移課表格子
func (s *TeacherMergeService) migrateTimetableCells(tx *gorm.DB, sourceID, targetID uint) error {
	result := tx.Model(&models.TimetableCell{}).
		Where("teacher_id = ?", sourceID).
		Update("teacher_id", targetID)

	if result.Error != nil {
		s.Logger.Error("遷移課表格子失敗", "error", result.Error)
		return fmt.Errorf("遷移課表格子失敗: %w", result.Error)
	}

	if result.RowsAffected > 0 {
		s.Logger.Info("已遷移課表格子",
			"count", result.RowsAffected,
			"source_teacher_id", sourceID,
			"target_teacher_id", targetID)
	}

	return nil
}

// migrateTeacherSkills 遷移教師技能
func (s *TeacherMergeService) migrateTeacherSkills(tx *gorm.DB, sourceID, targetID uint) error {
	// 先查詢來源教師的所有技能
	var sourceSkills []models.TeacherSkill
	if err := tx.Where("teacher_id = ?", sourceID).Find(&sourceSkills).Error; err != nil {
		s.Logger.Error("查詢來源教師技能失敗", "error", err)
		return fmt.Errorf("查詢來源教師技能失敗: %w", err)
	}

	if len(sourceSkills) == 0 {
		s.Logger.Info("來源教師無技能需要遷移")
		return nil
	}

	// 取得目標教師的現有技能
	var targetSkills []models.TeacherSkill
	if err := tx.Where("teacher_id = ?", targetID).Find(&targetSkills).Error; err != nil {
		s.Logger.Error("查詢目標教師技能失敗", "error", err)
		return fmt.Errorf("查詢目標教師技能失敗: %w", err)
	}

	// 建立目標技能的地圖（用於快速查找重複）
	targetSkillMap := make(map[string]bool) // key: category+skill_name+level
	for _, ts := range targetSkills {
		key := ts.Category + "|" + ts.SkillName + "|" + ts.Level
		targetSkillMap[key] = true
	}

	// 找出需要遷移的技能（不重複的）
	var skillsToMigrate []uint
	for _, skill := range sourceSkills {
		key := skill.Category + "|" + skill.SkillName + "|" + skill.Level
		if !targetSkillMap[key] {
			skillsToMigrate = append(skillsToMigrate, skill.ID)
		}
	}

	// 批量更新需要遷移的技能
	if len(skillsToMigrate) > 0 {
		result := tx.Model(&models.TeacherSkill{}).
			Where("id IN ?", skillsToMigrate).
			Update("teacher_id", targetID)

		if result.Error != nil {
			s.Logger.Error("遷移教師技能失敗", "error", result.Error)
			return fmt.Errorf("遷移教師技能失敗: %w", result.Error)
		}

		if result.RowsAffected > 0 {
			s.Logger.Info("已遷移教師技能",
				"count", result.RowsAffected,
				"source_teacher_id", sourceID,
				"target_teacher_id", targetID)
		}
	}

	// 清理來源教師剩餘的重複技能（那些因目標教師已有而未被更新的技能）
	if err := tx.Where("teacher_id = ?", sourceID).Delete(&models.TeacherSkill{}).Error; err != nil {
		s.Logger.Error("清理來源教師重複技能失敗", "error", err)
		return fmt.Errorf("清理來源教師重複技能失敗: %w", err)
	}

	return nil
}

// migrateTeacherCertificates 遷移教師證照
func (s *TeacherMergeService) migrateTeacherCertificates(tx *gorm.DB, sourceID, targetID uint) error {
	// 先查詢來源教師的所有證照
	var sourceCerts []models.TeacherCertificate
	if err := tx.Where("teacher_id = ?", sourceID).Find(&sourceCerts).Error; err != nil {
		s.Logger.Error("查詢來源教師證照失敗", "error", err)
		return fmt.Errorf("查詢來源教師證照失敗: %w", err)
	}

	if len(sourceCerts) == 0 {
		s.Logger.Info("來源教師無證照需要遷移")
		return nil
	}

	// 取得目標教師的現有證照
	var targetCerts []models.TeacherCertificate
	if err := tx.Where("teacher_id = ?", targetID).Find(&targetCerts).Error; err != nil {
		s.Logger.Error("查詢目標教師證照失敗", "error", err)
		return fmt.Errorf("查詢目標教師證照失敗: %w", err)
	}

	// 建立目標證照的地圖（用於快速查找重複）
	targetCertMap := make(map[string]bool) // key: name
	for _, tc := range targetCerts {
		targetCertMap[tc.Name] = true
	}

	// 找出需要遷移的證照（不重複的）
	var certsToMigrate []uint
	for _, cert := range sourceCerts {
		if !targetCertMap[cert.Name] {
			certsToMigrate = append(certsToMigrate, cert.ID)
		}
	}

	// 批量更新需要遷移的證照
	if len(certsToMigrate) > 0 {
		result := tx.Model(&models.TeacherCertificate{}).
			Where("id IN ?", certsToMigrate).
			Update("teacher_id", targetID)

		if result.Error != nil {
			s.Logger.Error("遷移教師證照失敗", "error", result.Error)
			return fmt.Errorf("遷移教師證照失敗: %w", result.Error)
		}

		if result.RowsAffected > 0 {
			s.Logger.Info("已遷移教師證照",
				"count", result.RowsAffected,
				"source_teacher_id", sourceID,
				"target_teacher_id", targetID)
		}
	}

	return nil
}

// migrateTeacherPersonalHashtags 遷移教師個人標籤
func (s *TeacherMergeService) migrateTeacherPersonalHashtags(tx *gorm.DB, sourceID, targetID uint) error {
	// TeacherPersonalHashtag 是複合主鍵表（TeacherID + HashtagID），
	// 沒有 ID 欄位，需要用「刪除後新增」的模式

	// 1. 先查詢來源教師的所有標籤
	var sourceTags []models.TeacherPersonalHashtag
	if err := tx.Where("teacher_id = ?", sourceID).Find(&sourceTags).Error; err != nil {
		s.Logger.Error("查詢來源教師標籤失敗", "error", err)
		return fmt.Errorf("查詢來源教師標籤失敗: %w", err)
	}

	if len(sourceTags) == 0 {
		s.Logger.Info("來源教師無標籤需要遷移")
		return nil
	}

	// 2. 取得目標教師的現有標籤
	var targetTags []models.TeacherPersonalHashtag
	if err := tx.Where("teacher_id = ?", targetID).Find(&targetTags).Error; err != nil {
		s.Logger.Error("查詢目標教師標籤失敗", "error", err)
		return fmt.Errorf("查詢目標教師標籤失敗: %w", err)
	}

	// 3. 建立目標標籤的地圖（用於快速查找重複）
	targetTagMap := make(map[uint]bool) // key: hashtag_id
	for _, tt := range targetTags {
		targetTagMap[tt.HashtagID] = true
	}

	// 4. 找出需要遷移的標籤（不重複的）
	var tagsToMigrate []models.TeacherPersonalHashtag
	for _, tag := range sourceTags {
		if !targetTagMap[tag.HashtagID] {
			tagsToMigrate = append(tagsToMigrate, tag)
		}
	}

	if len(tagsToMigrate) == 0 {
		s.Logger.Info("來源教師所有標籤都與目標教師重複，無需遷移")
		// 清理來源教師的標籤
		if err := tx.Where("teacher_id = ?", sourceID).Delete(&models.TeacherPersonalHashtag{}).Error; err != nil {
			s.Logger.Error("清理來源教師標籤失敗", "error", err)
			return fmt.Errorf("清理來源教師標籤失敗: %w", err)
		}
		return nil
	}

	// 5. 刪除來源教師的標籤
	if err := tx.Where("teacher_id = ?", sourceID).Delete(&models.TeacherPersonalHashtag{}).Error; err != nil {
		s.Logger.Error("刪除來源教師標籤失敗", "error", err)
		return fmt.Errorf("刪除來源教師標籤失敗: %w", err)
	}

	// 6. 新增不重複的標籤到目標教師
	var newTags []models.TeacherPersonalHashtag
	for _, tag := range tagsToMigrate {
		newTags = append(newTags, models.TeacherPersonalHashtag{
			TeacherID: targetID,
			HashtagID: tag.HashtagID,
			SortOrder: tag.SortOrder,
		})
	}

	if len(newTags) > 0 {
		if err := tx.Create(&newTags).Error; err != nil {
			s.Logger.Error("新增教師標籤失敗", "error", err)
			return fmt.Errorf("新增教師標籤失敗: %w", err)
		}

		s.Logger.Info("已遷移教師個人標籤",
			"count", len(newTags),
			"source_teacher_id", sourceID,
			"target_teacher_id", targetID)
	}

	return nil
}

// handleCenterMembership 處理中心會籍
// 如果目標教師已有會籍，則刪除來源教師的會籍
// 如果目標教師沒有會籍，則將來源教師的會籍遷移給目標教師
func (s *TeacherMergeService) handleCenterMembership(tx *gorm.DB, sourceID, targetID, centerID uint) error {
	// 檢查目標教師在該中心是否有會籍
	var targetMembership models.CenterMembership
	targetResult := tx.Where("teacher_id = ? AND center_id = ?", targetID, centerID).First(&targetMembership)

	if targetResult.Error == nil {
		// 目標教師已有會籍，刪除來源教師的會籍
		deleteResult := tx.Where("teacher_id = ? AND center_id = ?", sourceID, centerID).Delete(&models.CenterMembership{})
		if deleteResult.Error != nil {
			s.Logger.Error("刪除來源教師會籍失敗", "error", deleteResult.Error)
			return fmt.Errorf("刪除來源教師會籍失敗: %w", deleteResult.Error)
		}

		s.Logger.Info("目標教師已有會籍，已刪除來源教師會籍",
			"source_teacher_id", sourceID,
			"target_teacher_id", targetID,
			"center_id", centerID)
	} else if targetResult.Error == gorm.ErrRecordNotFound {
		// 目標教師沒有會籍，將來源教師的會籍遷移
		updateResult := tx.Model(&models.CenterMembership{}).
			Where("teacher_id = ? AND center_id = ?", sourceID, centerID).
			Update("teacher_id", targetID)

		if updateResult.Error != nil {
			s.Logger.Error("遷移會籍失敗", "error", updateResult.Error)
			return fmt.Errorf("遷移會籍失敗: %w", updateResult.Error)
		}

		if updateResult.RowsAffected > 0 {
			s.Logger.Info("已遷移會籍至目標教師",
				"source_teacher_id", sourceID,
				"target_teacher_id", targetID,
				"center_id", centerID)
		}
	} else {
		// 其他錯誤
		s.Logger.Error("查詢目標教師會籍失敗", "error", targetResult.Error)
		return fmt.Errorf("查詢目標教師會籍失敗: %w", targetResult.Error)
	}

	return nil
}

// handleTeacherNotes 處理教師筆記
// 將來源教師的筆記內容附加到目標教師筆記後方
func (s *TeacherMergeService) handleTeacherNotes(tx *gorm.DB, sourceID, targetID, centerID uint) error {
	var sourceNote models.CenterTeacherNote
	var targetNote models.CenterTeacherNote

	// 取得來源教師筆記
	sourceResult := tx.Where("teacher_id = ? AND center_id = ?", sourceID, centerID).First(&sourceNote)
	if sourceResult.Error != nil && sourceResult.Error != gorm.ErrRecordNotFound {
		s.Logger.Error("取得來源教師筆記失敗", "error", sourceResult.Error)
		return fmt.Errorf("取得來源教師筆記失敗: %w", sourceResult.Error)
	}

	// 取得目標教師筆記
	targetResult := tx.Where("teacher_id = ? AND center_id = ?", targetID, centerID).First(&targetNote)
	if targetResult.Error != nil && targetResult.Error != gorm.ErrRecordNotFound {
		s.Logger.Error("取得目標教師筆記失敗", "error", targetResult.Error)
		return fmt.Errorf("取得目標教師筆記失敗: %w", targetResult.Error)
	}

	// 情況1：來源和目標都有筆記，附加內容
	if sourceResult.Error == nil && targetResult.Error == nil {
		now := time.Now()
		separator := "\n\n--- 合併自 [" + fmt.Sprintf("%d", sourceNote.TeacherID) + "] ---\n\n"

		targetNote.InternalNote += separator + sourceNote.InternalNote
		targetNote.UpdatedAt = now

		updateResult := tx.Model(&targetNote).
			Where("id = ?", targetNote.ID).
			Updates(map[string]interface{}{
				"internal_note": targetNote.InternalNote,
				"updated_at":    now,
			})

		if updateResult.Error != nil {
			s.Logger.Error("更新目標教師筆記失敗", "error", updateResult.Error)
			return fmt.Errorf("更新目標教師筆記失敗: %w", updateResult.Error)
		}

		// 刪除來源教師筆記
		deleteResult := tx.Where("id = ?", sourceNote.ID).Delete(&models.CenterTeacherNote{})
		if deleteResult.Error != nil {
			s.Logger.Error("刪除來源教師筆記失敗", "error", deleteResult.Error)
			return fmt.Errorf("刪除來源教師筆記失敗: %w", deleteResult.Error)
		}

		s.Logger.Info("已合併教師筆記",
			"source_teacher_id", sourceID,
			"target_teacher_id", targetID)
	} else if sourceResult.Error == nil && targetResult.Error == gorm.ErrRecordNotFound {
		// 情況2：只有來源有筆記，遷移給目標
		moveResult := tx.Model(&models.CenterTeacherNote{}).
			Where("id = ?", sourceNote.ID).
			Updates(map[string]interface{}{
				"teacher_id": targetID,
				"updated_at": time.Now(),
			})

		if moveResult.Error != nil {
			s.Logger.Error("遷移來源教師筆記失敗", "error", moveResult.Error)
			return fmt.Errorf("遷移來源教師筆記失敗: %w", moveResult.Error)
		}

		s.Logger.Info("已遷移教師筆記至目標教師",
			"source_teacher_id", sourceID,
			"target_teacher_id", targetID)
	} else if sourceResult.Error == gorm.ErrRecordNotFound {
		// 情況3：只有目標有筆記，來源沒有筆記，無需操作
		s.Logger.Info("來源教師無筆記，無需合併",
			"source_teacher_id", sourceID,
			"target_teacher_id", targetID)
	}

	return nil
}

// softDeleteTeacher 軟刪除來源教師
func (s *TeacherMergeService) softDeleteTeacher(tx *gorm.DB, sourceID uint) error {
	result := tx.Model(&models.Teacher{}).
		Where("id = ?", sourceID).
		Update("deleted_at", time.Now())

	if result.Error != nil {
		s.Logger.Error("軟刪除來源教師失敗", "error", result.Error)
		return fmt.Errorf("軟刪除來源教師失敗: %w", result.Error)
	}

	s.Logger.Info("已軟刪除來源教師",
		"source_teacher_id", sourceID)

	return nil
}
