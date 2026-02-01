package services

import (
	"context"
	"fmt"
	"strings"
	"time"
	"timeLedger/app"
	"timeLedger/app/models"
	"timeLedger/app/repositories"
	"timeLedger/app/resources"
	"timeLedger/global/errInfos"

	"gorm.io/gorm"
)

// TeacherProfileService 教師個人檔案相關業務邏輯
type TeacherProfileService struct {
	BaseService
	app             *app.App
	teacherRepo     *repositories.TeacherRepository
	membershipRepo  *repositories.CenterMembershipRepository
	centerRepo      *repositories.CenterRepository
	skillRepo       *repositories.TeacherSkillRepository
	certificateRepo *repositories.TeacherCertificateRepository
	hashtagRepo     *repositories.HashtagRepository
	auditLogRepo    *repositories.AuditLogRepository
}

// NewTeacherProfileService 建立教師檔案服務
func NewTeacherProfileService(app *app.App) *TeacherProfileService {
	return &TeacherProfileService{
		app:             app,
		teacherRepo:     repositories.NewTeacherRepository(app),
		membershipRepo:  repositories.NewCenterMembershipRepository(app),
		centerRepo:      repositories.NewCenterRepository(app),
		skillRepo:       repositories.NewTeacherSkillRepository(app),
		certificateRepo: repositories.NewTeacherCertificateRepository(app),
		hashtagRepo:     repositories.NewHashtagRepository(app),
		auditLogRepo:    repositories.NewAuditLogRepository(app),
	}
}

// GetProfile 取得老師個人資料
func (s *TeacherProfileService) GetProfile(ctx context.Context, teacherID uint) (*resources.TeacherProfileResource, *errInfos.Res, error) {
	teacher, err := s.teacherRepo.GetByID(ctx, teacherID)
	if err != nil {
		return nil, s.app.Err.New(errInfos.SQL_ERROR), err
	}

	// 取得個人標籤
	personalHashtags, _ := s.teacherRepo.ListPersonalHashtags(ctx, teacherID)
	var hashtagResources []resources.PersonalHashtag
	for _, h := range personalHashtags {
		hashtagResources = append(hashtagResources, resources.PersonalHashtag{
			HashtagID: h.HashtagID,
			Name:      h.Name,
			SortOrder: h.SortOrder,
		})
	}

	return &resources.TeacherProfileResource{
		ID:                teacher.ID,
		LineUserID:        teacher.LineUserID,
		Name:              teacher.Name,
		Email:             teacher.Email,
		Bio:               teacher.Bio,
		City:              teacher.City,
		District:          teacher.District,
		PublicContactInfo: teacher.PublicContactInfo,
		IsOpenToHiring:    teacher.IsOpenToHiring,
		PersonalHashtags:  hashtagResources,
	}, nil, nil
}

// UpdateProfileRequest 更新老師個人資料請求
type UpdateProfileRequest struct {
	Bio               string
	City              string
	District          string
	PublicContactInfo string
	IsOpenToHiring    bool
	PersonalHashtags  []string
}

// UpdateProfile 更新老師個人資料（使用交易確保原子性）
func (s *TeacherProfileService) UpdateProfile(ctx context.Context, teacherID uint, req *UpdateProfileRequest) (*resources.TeacherProfileResource, *errInfos.Res, error) {
	teacher, err := s.teacherRepo.GetByID(ctx, teacherID)
	if err != nil {
		return nil, s.app.Err.New(errInfos.SQL_ERROR), err
	}

	if req.Bio != "" {
		teacher.Bio = req.Bio
	}
	if req.City != "" {
		teacher.City = req.City
	}
	if req.District != "" {
		teacher.District = req.District
	}
	if req.PublicContactInfo != "" {
		teacher.PublicContactInfo = req.PublicContactInfo
	}

	teacher.IsOpenToHiring = req.IsOpenToHiring

	// 使用交易更新資料和標籤
	var result *resources.TeacherProfileResource

	txErr := s.app.MySQL.WDB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// 更新老師資料
		if err := tx.Save(&teacher).Error; err != nil {
			return fmt.Errorf("failed to update teacher: %w", err)
		}

		// 更新個人標籤
		if len(req.PersonalHashtags) > 0 {
			if err := s.updatePersonalHashtagsWithTx(tx, ctx, teacherID, req.PersonalHashtags); err != nil {
				return err
			}
		}

		// 記錄審核日誌
		memberships, _ := s.membershipRepo.GetActiveByTeacherID(ctx, teacherID)
		var centerID uint
		if len(memberships) > 0 {
			centerID = memberships[0].CenterID
		}

		auditLog := models.AuditLog{
			CenterID:   centerID,
			ActorType:  "TEACHER",
			ActorID:    teacherID,
			Action:     "PROFILE_UPDATE",
			TargetType: "Teacher",
			TargetID:   teacherID,
			Payload: models.AuditPayload{
				After: map[string]interface{}{
					"bio":               req.Bio,
					"city":              req.City,
					"district":          req.District,
					"is_open_to_hiring": req.IsOpenToHiring,
				},
			},
		}
		if err := tx.Create(&auditLog).Error; err != nil {
			return fmt.Errorf("failed to create audit log: %w", err)
		}

		result = &resources.TeacherProfileResource{
			ID:                teacher.ID,
			LineUserID:        teacher.LineUserID,
			Name:              teacher.Name,
			Email:             teacher.Email,
			Bio:               teacher.Bio,
			City:              teacher.City,
			District:          teacher.District,
			PublicContactInfo: teacher.PublicContactInfo,
			IsOpenToHiring:    teacher.IsOpenToHiring,
		}

		return nil
	})

	if txErr != nil {
		return nil, s.app.Err.New(errInfos.ERR_TX_FAILED), txErr
	}

	return result, nil, nil
}

// updatePersonalHashtagsWithTx 更新個人標籤（交易版本）
func (s *TeacherProfileService) updatePersonalHashtagsWithTx(tx *gorm.DB, ctx context.Context, teacherID uint, tags []string) error {
	// 刪除現有標籤（使用交易連接）
	if err := tx.Where("teacher_id = ?", teacherID).Delete(&models.TeacherPersonalHashtag{}).Error; err != nil {
		return fmt.Errorf("failed to delete personal hashtags: %w", err)
	}

	// 創建交易感知的 HashtagRepository，使用交易 DB 連線
	hashtagTxRepo := repositories.NewTransactionRepo[models.Hashtag](ctx, tx, "hashtags")

	sortOrder := 0
	for _, tagName := range tags {
		// 確保 # 符號存在
		name := tagName
		if !strings.HasPrefix(name, "#") {
			name = "#" + name
		}

		// 查找標籤
		hashtag, err := s.hashtagRepo.GetByName(ctx, name)
		if err != nil {
			// 標籤不存在，創建新標籤
			hashtagModel := models.Hashtag{Name: name, UsageCount: 1}
			createdHashtag, err := hashtagTxRepo.Create(ctx, hashtagModel)
			if err != nil {
				return fmt.Errorf("failed to create hashtag '%s': %w", name, err)
			}
			hashtag = &createdHashtag
		} else {
			// 標籤存在，更新使用次數（使用交易連線）
			if err := tx.WithContext(ctx).Model(&models.Hashtag{}).Where("id = ?", hashtag.ID).UpdateColumn("usage_count", "usage_count + 1").Error; err != nil {
				return fmt.Errorf("failed to increment hashtag usage: %w", err)
			}
		}

		// 創建關聯
		personalHashtag := models.TeacherPersonalHashtag{
			TeacherID: teacherID,
			HashtagID: hashtag.ID,
			SortOrder: sortOrder,
		}
		if err := tx.Create(&personalHashtag).Error; err != nil {
			return fmt.Errorf("failed to create personal hashtag: %w", err)
		}
		sortOrder++
	}

	return nil
}

// updatePersonalHashtags 更新個人標籤
func (s *TeacherProfileService) updatePersonalHashtags(ctx context.Context, teacherID uint, tags []string) error {
	// 刪除現有標籤
	if err := s.teacherRepo.DeleteAllPersonalHashtags(ctx, teacherID); err != nil {
		return fmt.Errorf("failed to delete personal hashtags: %w", err)
	}

	sortOrder := 0
	for _, tagName := range tags {
		// 確保 # 符號存在
		name := tagName
		if !strings.HasPrefix(name, "#") {
			name = "#" + name
		}

		// 查找標籤
		hashtag, err := s.hashtagRepo.GetByName(ctx, name)
		if err != nil {
			// 標籤不存在，創建新標籤
			hashtagModel := models.Hashtag{Name: name, UsageCount: 1}
			createdHashtag, err := s.hashtagRepo.Create(ctx, hashtagModel)
			if err != nil {
				return fmt.Errorf("failed to create hashtag '%s': %w", name, err)
			}
			hashtag = &createdHashtag
		} else {
			// 標籤存在，更新使用次數
			if err := s.hashtagRepo.IncrementUsage(ctx, name); err != nil {
				return fmt.Errorf("failed to increment hashtag usage: %w", err)
			}
		}

		// 創建關聯
		if err := s.teacherRepo.CreatePersonalHashtag(ctx, teacherID, hashtag.ID, sortOrder); err != nil {
			return fmt.Errorf("failed to create personal hashtag: %w", err)
		}
		sortOrder++
	}

	return nil
}

// GetCenters 取得老師已加入的中心列表
func (s *TeacherProfileService) GetCenters(ctx context.Context, teacherID uint) ([]resources.CenterMembershipResource, *errInfos.Res, error) {
	memberships, err := s.membershipRepo.GetActiveByTeacherID(ctx, teacherID)
	if err != nil {
		return nil, s.app.Err.New(errInfos.SQL_ERROR), err
	}

	var centerResources []resources.CenterMembershipResource
	for _, m := range memberships {
		center, _ := s.centerRepo.GetByID(ctx, m.CenterID)
		centerResources = append(centerResources, resources.CenterMembershipResource{
			ID:         m.ID,
			CenterID:   m.CenterID,
			CenterName: center.Name,
			Status:     string(m.Status),
			CreatedAt:  m.CreatedAt,
		})
	}

	return centerResources, nil, nil
}

// GetSkills 取得老師技能列表
func (s *TeacherProfileService) GetSkills(ctx context.Context, teacherID uint) ([]models.TeacherSkill, *errInfos.Res, error) {
	skills, err := s.skillRepo.ListByTeacherID(ctx, teacherID)
	if err != nil {
		return nil, s.app.Err.New(errInfos.SQL_ERROR), err
	}

	return skills, nil, nil
}

// CreateSkillRequest 新增技能請求
type CreateSkillRequest struct {
	Category   string
	SkillName  string
	Level      string
	HashtagIDs []uint
}

// CreateSkill 新增老師技能（使用交易確保原子性）
func (s *TeacherProfileService) CreateSkill(ctx context.Context, teacherID uint, req *CreateSkillRequest) (*models.TeacherSkill, *errInfos.Res, error) {
	skill := models.TeacherSkill{
		TeacherID: teacherID,
		Category:  req.Category,
		SkillName: req.SkillName,
		Level:     req.Level,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// 使用交易建立技能和標籤關聯
	var createdSkill *models.TeacherSkill

	txErr := s.app.MySQL.WDB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// 建立技能
		if err := tx.Create(&skill).Error; err != nil {
			return fmt.Errorf("failed to create skill: %w", err)
		}
		createdSkill = &skill

		// 建立技能標籤關聯
		if len(req.HashtagIDs) > 0 {
			for _, hashtagID := range req.HashtagIDs {
				skillHashtag := models.TeacherSkillHashtag{
					TeacherSkillID: skill.ID,
					HashtagID:      hashtagID,
				}
				if err := tx.Create(&skillHashtag).Error; err != nil {
					return fmt.Errorf("failed to create skill hashtag: %w", err)
				}
			}
		}

		return nil
	})

	if txErr != nil {
		return nil, s.app.Err.New(errInfos.ERR_TX_FAILED), txErr
	}

	return createdSkill, nil, nil
}

// UpdateSkillRequest 更新技能請求
type UpdateSkillRequest struct {
	Category  string
	SkillName string
	Hashtags  []string
}

// UpdateSkill 更新老師技能
func (s *TeacherProfileService) UpdateSkill(ctx context.Context, skillID, teacherID uint, req *UpdateSkillRequest) (*models.TeacherSkill, *errInfos.Res, error) {
	skill, err := s.skillRepo.GetByID(ctx, skillID)
	if err != nil {
		return nil, s.app.Err.New(errInfos.NOT_FOUND), err
	}

	if skill.TeacherID != teacherID {
		return nil, s.app.Err.New(errInfos.FORBIDDEN), nil
	}

	skill.Category = req.Category
	skill.SkillName = req.SkillName
	skill.UpdatedAt = time.Now()

	if err := s.skillRepo.Update(ctx, skill); err != nil {
		return nil, s.app.Err.New(errInfos.SQL_ERROR), err
	}

	// 更新技能標籤
	if len(req.Hashtags) > 0 {
		if err := s.updateSkillHashtags(ctx, skill.ID, req.Hashtags); err != nil {
			return nil, s.app.Err.New(errInfos.SQL_ERROR), err
		}
	}

	return &skill, nil, nil
}

// updateSkillHashtags 更新技能標籤
func (s *TeacherProfileService) updateSkillHashtags(ctx context.Context, skillID uint, tags []string) error {
	// 刪除現有標籤
	if err := s.skillRepo.DeleteAllHashtags(ctx, skillID); err != nil {
		return fmt.Errorf("failed to delete skill hashtags: %w", err)
	}

	for _, tagName := range tags {
		// 確保 # 符號存在
		name := tagName
		if !strings.HasPrefix(name, "#") {
			name = "#" + name
		}

		// 查找標籤
		hashtag, err := s.hashtagRepo.GetByName(ctx, name)
		if err != nil {
			// 標籤不存在，創建新標籤
			hashtagModel := models.Hashtag{Name: name, UsageCount: 1}
			createdHashtag, err := s.hashtagRepo.Create(ctx, hashtagModel)
			if err != nil {
				return fmt.Errorf("failed to create hashtag '%s': %w", name, err)
			}
			hashtag = &createdHashtag
		} else {
			// 標籤存在，更新使用次數
			if err := s.hashtagRepo.IncrementUsage(ctx, name); err != nil {
				return fmt.Errorf("failed to increment hashtag usage: %w", err)
			}
		}

		// 創建關聯
		if err := s.skillRepo.CreateHashtag(ctx, skillID, hashtag.ID); err != nil {
			return fmt.Errorf("failed to create skill hashtag: %w", err)
		}
	}

	return nil
}

// DeleteSkill 刪除老師技能
func (s *TeacherProfileService) DeleteSkill(ctx context.Context, skillID, teacherID uint) *errInfos.Res {
	skill, err := s.skillRepo.GetByID(ctx, skillID)
	if err != nil {
		return s.app.Err.New(errInfos.NOT_FOUND)
	}

	if skill.TeacherID != teacherID {
		return s.app.Err.New(errInfos.FORBIDDEN)
	}

	if err := s.skillRepo.DeleteByID(ctx, skillID); err != nil {
		return s.app.Err.New(errInfos.SQL_ERROR)
	}

	return nil
}

// GetCertificates 取得老師證照列表
func (s *TeacherProfileService) GetCertificates(ctx context.Context, teacherID uint) ([]models.TeacherCertificate, *errInfos.Res, error) {
	certificates, err := s.certificateRepo.ListByTeacherID(ctx, teacherID)
	if err != nil {
		return nil, s.app.Err.New(errInfos.SQL_ERROR), err
	}

	return certificates, nil, nil
}

// CreateCertificateRequest 新增證照請求
type CreateCertificateRequest struct {
	Name     string
	FileURL  string
	IssuedAt time.Time
}

// CreateCertificate 新增老師證照
func (s *TeacherProfileService) CreateCertificate(ctx context.Context, teacherID uint, req *CreateCertificateRequest) (*models.TeacherCertificate, *errInfos.Res, error) {
	certificate := models.TeacherCertificate{
		TeacherID: teacherID,
		Name:      req.Name,
		FileURL:   req.FileURL,
		IssuedAt:  req.IssuedAt,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	createdCert, err := s.certificateRepo.Create(ctx, certificate)
	if err != nil {
		return nil, s.app.Err.New(errInfos.SQL_ERROR), err
	}

	return &createdCert, nil, nil
}

// DeleteCertificate 刪除老師證照
func (s *TeacherProfileService) DeleteCertificate(ctx context.Context, certificateID, teacherID uint) *errInfos.Res {
	certificate, err := s.certificateRepo.GetByID(ctx, certificateID)
	if err != nil {
		return s.app.Err.New(errInfos.NOT_FOUND)
	}

	if certificate.TeacherID != teacherID {
		return s.app.Err.New(errInfos.FORBIDDEN)
	}

	if err := s.certificateRepo.DeleteByID(ctx, certificateID); err != nil {
		return s.app.Err.New(errInfos.SQL_ERROR)
	}

	return nil
}
