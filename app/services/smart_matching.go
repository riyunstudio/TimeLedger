package services

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"
	"timeLedger/app"
	"timeLedger/app/models"
	"timeLedger/app/repositories"
)

// 繁簡轉換對照表（常見地名用字）
var traditionalToSimplified = map[string]string{
	"臺": "台",
	"縣": "县",
	"市": "市", // 相同
}

// normalizeCityName 將城市名稱正規化為資料庫格式（繁體中文）
func normalizeCityName(name string) string {
	if name == "" {
		return name
	}

	// 嘗試將簡體轉換為繁體（因為資料庫使用繁體）
	result := name
	for traditional, simplified := range traditionalToSimplified {
		// 檢查是否包含簡體版本
		if strings.Contains(result, simplified) && !strings.Contains(result, traditional) {
			// 簡體版本存在但繁體不存在，進行替換
			result = strings.ReplaceAll(result, simplified, traditional)
		}
	}

	// 特殊地名對照
	specialMappings := map[string]string{
		"台北市": "臺北市",
		"台中市": "臺中市",
		"台南市": "臺南市",
		"台東縣": "臺東縣",
		"台西鄉": "臺西鄉",
	}
	if mapped, ok := specialMappings[result]; ok {
		return mapped
	}

	return result
}

// normalizeDistrictName 將區域名稱正規化
func normalizeDistrictName(name string) string {
	if name == "" {
		return name
	}

	// 應用相同的繁簡轉換
	result := name
	for traditional, simplified := range traditionalToSimplified {
		if strings.Contains(result, simplified) && !strings.Contains(result, traditional) {
			result = strings.ReplaceAll(result, simplified, traditional)
		}
	}
	return result
}

type MatchAvailability string

const (
	MatchAvailable      MatchAvailability = "AVAILABLE"
	MatchBufferConflict MatchAvailability = "BUFFER_CONFLICT"
	MatchOverlap        MatchAvailability = "OVERLAP"
)

type SmartMatchingServiceImpl struct {
	BaseService
	app                    *app.App
	teacherRepository      *repositories.TeacherRepository
	scheduleRuleRepo       *repositories.ScheduleRuleRepository
	scheduleExceptionRepo  *repositories.ScheduleExceptionRepository
	teacherSkillRepo       *repositories.TeacherSkillRepository
	hashtagRepo            *repositories.HashtagRepository
	teacherCertificateRepo *repositories.TeacherCertificateRepository
	centerTeacherNoteRepo  *repositories.CenterTeacherNoteRepository
	centerInvitationRepo   *repositories.CenterInvitationRepository
	membershipRepo         *repositories.CenterMembershipRepository
	validationSvc          ScheduleValidationService
	notificationService    NotificationService
}

func NewSmartMatchingService(app *app.App) SmartMatchingService {
	return &SmartMatchingServiceImpl{
		app:                    app,
		teacherRepository:      repositories.NewTeacherRepository(app),
		scheduleRuleRepo:       repositories.NewScheduleRuleRepository(app),
		scheduleExceptionRepo:  repositories.NewScheduleExceptionRepository(app),
		teacherSkillRepo:       repositories.NewTeacherSkillRepository(app),
		hashtagRepo:            repositories.NewHashtagRepository(app),
		teacherCertificateRepo: repositories.NewTeacherCertificateRepository(app),
		centerTeacherNoteRepo:  repositories.NewCenterTeacherNoteRepository(app),
		centerInvitationRepo:   repositories.NewCenterInvitationRepository(app),
		membershipRepo:         repositories.NewCenterMembershipRepository(app),
		validationSvc:          NewScheduleValidationService(app),
		notificationService:    NewNotificationService(app),
	}
}

// FindMatches searches for available teachers matching the criteria
// 搜尋範圍：
// 1. 該中心已有排課的老師
// 2. 開放對外徵才的老師（可邀請加入）
func (s *SmartMatchingServiceImpl) FindMatches(ctx context.Context, centerID uint, teacherID *uint, roomIDs []uint, startTime, endTime time.Time, requiredSkills []string, excludeTeacherIDs []uint) ([]MatchScore, error) {
	// 建立排除教師 ID 的 Map（O(1) 查找）
	excludeMap := make(map[uint]bool)
	for _, id := range excludeTeacherIDs {
		excludeMap[id] = true
	}

	// 取得候選老師列表
	candidateTeacherIDs, err := s.getCandidateTeachers(ctx, centerID)
	if err != nil {
		return nil, err
	}

	// 過濾排除的老師
	filteredTeacherIDs := make([]uint, 0, len(candidateTeacherIDs))
	for _, id := range candidateTeacherIDs {
		if !excludeMap[id] {
			filteredTeacherIDs = append(filteredTeacherIDs, id)
		}
	}

	if len(filteredTeacherIDs) == 0 {
		return []MatchScore{}, nil
	}

	// 批量查詢所有教師數據
	teachersMap, err := s.teacherRepository.BatchGetByIDs(ctx, filteredTeacherIDs)
	if err != nil {
		return nil, err
	}

	// 批量查詢所有技能數據
	skillsMap, err := s.teacherSkillRepo.BatchListByTeacherIDs(ctx, filteredTeacherIDs)
	if err != nil {
		return nil, err
	}

	// 批量查詢所有中心教師備註
	centerNotesMap, err := s.centerTeacherNoteRepo.BatchGetByCenterAndTeachers(ctx, centerID, filteredTeacherIDs)
	if err != nil {
		return nil, err
	}

	var matches []MatchScore

	// 對每個候選老師進行評分
	for _, currentTeacherID := range filteredTeacherIDs {
		teacher, ok := teachersMap[currentTeacherID]
		if !ok {
			continue
		}

		skills := skillsMap[currentTeacherID]
		centerNote := centerNotesMap[currentTeacherID]

		// 檢查老師是否為中心成員
		isMember := centerNote.ID != 0

		// 計算各項分數
		skillScore := calculateSkillMatchScore(skills, requiredSkills)
		regionScore := calculateRegionMatchScore(teacher.City, teacher.District, "", "")
		hashtagScore := calculateHashtagScore(skills, []string{})
		internalScore := calculateInternalScore(centerNote.Rating, centerNote.InternalNote)

		// 可用性檢查（只有中心成員才需要檢查排課衝突）
		availability := MatchAvailable
		availabilityScore := 0

		if isMember {
			// 中心成員：檢查時段是否有空
			availability, availabilityScore = s.checkTeacherAvailability(ctx, centerID, currentTeacherID, roomIDs, startTime, endTime)
		} else if teacher.IsOpenToHiring {
			// 外部老師（開放徵才）：預設可用，但標記為可邀請
			availability = MatchAvailable
			availabilityScore = 40
		} else {
			// 不符合條件的老師（已排除）
			continue
		}

		totalScore := availabilityScore + internalScore + skillScore + regionScore + hashtagScore

		match := MatchScore{
			TeacherID:         currentTeacherID,
			Name:              teacher.Name,
			Score:             totalScore,
			Availability:      availability,
			AvailabilityScore: availabilityScore,
			InternalScore:     internalScore,
			SkillScore:        skillScore,
			RegionScore:       regionScore,
			Skills:            extractSkills(skills),
			Hashtags:          extractHashtags(skills),
			IsMember:          isMember,
			InternalNote:      centerNote.InternalNote,
			InternalRating:    centerNote.Rating,
			PublicContactInfo: teacher.PublicContactInfo,
		}

		matches = append(matches, match)
	}

	sortMatchesByScore(matches)

	return matches, nil
}

// getCandidateTeachers 取得候選老師列表（中心成員 + 開放徵才）
func (s *SmartMatchingServiceImpl) getCandidateTeachers(ctx context.Context, centerID uint) ([]uint, error) {
	teacherIDs := make([]uint, 0)

	// 1. 取得該中心已有排課規則的老師
	rules, err := s.scheduleRuleRepo.ListByCenterID(ctx, centerID)
	if err != nil {
		return nil, err
	}

	ruleTeacherMap := make(map[uint]bool) // 記錄哪些老師有排課規則
	for _, rule := range rules {
		if rule.TeacherID != nil {
			teacherIDs = append(teacherIDs, *rule.TeacherID)
			ruleTeacherMap[*rule.TeacherID] = true
		}
	}

	// 2. 取得中心成員（已加入該中心的老師）
	memberTeacherIDs, err := s.membershipRepo.ListTeacherIDsByCenterID(ctx, centerID)
	if err != nil {
		return nil, err
	}

	for _, tid := range memberTeacherIDs {
		// 如果成員還沒在列表中，加入
		if !ruleTeacherMap[tid] {
			teacherIDs = append(teacherIDs, tid)
			ruleTeacherMap[tid] = true
		}
	}

	// 3. 取得開放對外徵才的老師
	allTeachers, err := s.teacherRepository.List(ctx)
	if err != nil {
		return nil, err
	}

	for _, t := range allTeachers {
		// 只加入開放徵才且不在列表中的老師
		if t.IsOpenToHiring && !ruleTeacherMap[t.ID] {
			teacherIDs = append(teacherIDs, t.ID)
		}
	}

	return teacherIDs, nil
}

// checkTeacherAvailability 檢查老師在某時段是否可用
// 修正：應該檢查所有教室，回傳最好的結果
func (s *SmartMatchingServiceImpl) checkTeacherAvailability(ctx context.Context, centerID uint, teacherID uint, roomIDs []uint, startTime, endTime time.Time) (MatchAvailability, int) {
	teacherIDPtr := &teacherID

	// 優先級：OVERLAP (0) > BUFFER_CONFLICT (15) > AVAILABLE (40)
	// 只要有一個教室重疊就是重疊，全部緩衝衝突才是緩衝衝突
	hasOverlap := false
	hasBufferConflict := false

	// 檢查每個教室
	for _, roomID := range roomIDs {
		result, err := s.validationSvc.CheckOverlap(ctx, centerID, teacherIDPtr, roomID, startTime, endTime, 0, nil, "")
		if err != nil {
			// 發生錯誤時，視為可用（保守策略）
			return MatchAvailable, 40
		}

		// 檢查是否有重疊
		if !result.Valid {
			hasOverlap = true
		} else if len(result.Conflicts) > 0 {
			// 有緩衝衝突但無重疊
			hasBufferConflict = true
		}
	}

	// 回傳最好的結果（優先級：可用 > 緩衝衝突 > 重疊）
	if hasOverlap {
		return MatchOverlap, 0
	}
	if hasBufferConflict {
		return MatchBufferConflict, 15
	}

	return MatchAvailable, 40
}

// NewPagination - 建立分頁資訊
func NewPagination(page, limit, total int) Pagination {
	totalPages := 0
	if limit > 0 {
		totalPages = (total + limit - 1) / limit
	}
	return Pagination{
		Page:       page,
		Limit:      limit,
		Total:      total,
		TotalPages: totalPages,
		HasNext:    page < totalPages,
		HasPrev:    page > 1,
	}
}

// SearchTalent searches for teachers based on criteria with pagination and sorting
// 優化版本：減少 N+1 查詢問題
func (s *SmartMatchingServiceImpl) SearchTalent(ctx context.Context, searchParams TalentSearchParams) (*TalentSearchResultResponse, error) {
	// 取得該中心的成員 ID 列表（如果指定了中心）
	memberIDMap := make(map[uint]bool)
	if searchParams.CenterID > 0 {
		memberTeacherIDs, err := s.membershipRepo.ListTeacherIDsByCenterID(ctx, searchParams.CenterID)
		if err == nil {
			for _, id := range memberTeacherIDs {
				memberIDMap[id] = true
			}
		}
	}

	// 第一階段：基本過濾（減少載入的資料量）
	teachers, err := s.teacherRepository.List(ctx)
	if err != nil {
		return nil, err
	}

	// 收集需要批量查詢的教師 ID
	candidateTeacherIDs := make([]uint, 0)
	filteredTeachers := make(map[uint]*models.Teacher)

	for i := range teachers {
		teacher := &teachers[i]

		// 只搜尋開放徵才的老師
		if !teacher.IsOpenToHiring {
			continue
		}

		// 如果指定了中心，排除已加入該中心的老師
		if searchParams.City != "" {
			normalizedSearchCity := normalizeCityName(searchParams.City)
			normalizedTeacherCity := normalizeCityName(teacher.City)
			if normalizedTeacherCity != normalizedSearchCity {
				continue
			}
		}

		if searchParams.District != "" {
			normalizedSearchDistrict := normalizeDistrictName(searchParams.District)
			normalizedTeacherDistrict := normalizeDistrictName(teacher.District)
			if normalizedTeacherDistrict != normalizedSearchDistrict {
				continue
			}
		}

		if searchParams.Keyword != "" {
			keywordLower := strings.ToLower(searchParams.Keyword)
			if !strings.Contains(strings.ToLower(teacher.Name), keywordLower) &&
				!strings.Contains(strings.ToLower(teacher.Bio), keywordLower) {
				continue
			}
		}

		// 如果指定了中心，排除已加入該中心的老師
		if searchParams.CenterID > 0 && memberIDMap[teacher.ID] {
			continue
		}

		candidateTeacherIDs = append(candidateTeacherIDs, teacher.ID)
		filteredTeachers[teacher.ID] = teacher
	}

	if len(candidateTeacherIDs) == 0 {
		// 沒有符合條件的老師，直接返回空結果
		return &TalentSearchResultResponse{
			Talents:   []TalentCardResponse{},
			Pagination: Pagination{Page: 1, Limit: 20, Total: 0},
		}, nil
	}

	// 第二階段：批量查詢相關資料（解決 N+1 問題）
	skillsMap, err := s.teacherSkillRepo.BatchListByTeacherIDs(ctx, candidateTeacherIDs)
	if err != nil {
		s.Logger.Warn("failed to batch query teacher skills", "error", err.Error())
	}
	certificatesMap, err := s.teacherCertificateRepo.BatchListByTeacherIDs(ctx, candidateTeacherIDs)
	if err != nil {
		s.Logger.Warn("failed to batch query teacher certificates", "error", err.Error())
	}
	notesMap, err := s.centerTeacherNoteRepo.BatchGetByCenterAndTeachers(ctx, searchParams.CenterID, candidateTeacherIDs)
	if err != nil {
		s.Logger.Warn("failed to batch query center teacher notes", "error", err.Error())
	}

	// 第三階段：根據技能和標籤過濾，並建立結果
	var allResults []TalentResult

	for _, teacherID := range candidateTeacherIDs {
		teacher, ok := filteredTeachers[teacherID]
		if !ok {
			continue
		}

		skillsList := skillsMap[teacherID]
		certificates := certificatesMap[teacherID]
		centerNote := notesMap[teacherID]

		// 技能篩選
		if len(searchParams.Skills) > 0 {
			hasSkill := false
			for _, requiredSkill := range searchParams.Skills {
				for _, teacherSkill := range skillsList {
					if teacherSkill.SkillName == requiredSkill || teacherSkill.Category == requiredSkill {
						hasSkill = true
						break
					}
				}
				if hasSkill {
					break
				}
			}
			if !hasSkill {
				continue
			}
		}

		// 標籤篩選
		if len(searchParams.Hashtags) > 0 {
			personalHashtags := s.extractPersonalHashtagsBatch(ctx, teacherID, skillsList)
			hasMatchingHashtag := false

			for _, requiredTag := range searchParams.Hashtags {
				normalizedTag := strings.TrimPrefix(requiredTag, "#")
				for _, personalTag := range personalHashtags {
					personalTagNormalized := strings.TrimPrefix(personalTag, "#")
					if strings.Contains(strings.ToLower(personalTagNormalized), strings.ToLower(normalizedTag)) ||
						strings.Contains(strings.ToLower(normalizedTag), strings.ToLower(personalTagNormalized)) {
						hasMatchingHashtag = true
						break
					}
				}
				if hasMatchingHashtag {
					break
				}
			}

			// 也檢查技能標籤
			if !hasMatchingHashtag {
				for _, skill := range skillsList {
					for _, tag := range skill.Hashtags {
						if tag.Hashtag.Name == "" {
							continue
						}
						tagNormalized := strings.TrimPrefix(tag.Hashtag.Name, "#")
						for _, reqTag := range searchParams.Hashtags {
							normalizedReqTag := strings.TrimPrefix(reqTag, "#")
							if strings.Contains(strings.ToLower(tagNormalized), strings.ToLower(normalizedReqTag)) ||
								strings.Contains(strings.ToLower(normalizedReqTag), strings.ToLower(tagNormalized)) {
								hasMatchingHashtag = true
								break
							}
						}
						if hasMatchingHashtag {
							break
						}
					}
					if hasMatchingHashtag {
						break
					}
				}
			}

			if !hasMatchingHashtag {
				continue
			}
		}

		certificatesList := make([]Certificate, 0, len(certificates))
		for _, cert := range certificates {
			certificatesList = append(certificatesList, Certificate{
				Name:     cert.Name,
				FileURL:  cert.FileURL,
				IssuedAt: cert.IssuedAt,
			})
		}

		result := TalentResult{
			TeacherID:         teacher.ID,
			Name:              teacher.Name,
			Bio:               teacher.Bio,
			City:              teacher.City,
			District:          teacher.District,
			Skills:            extractSkills(skillsList),
			PersonalHashtags:  s.extractPersonalHashtagsBatch(ctx, teacherID, skillsList),
			IsOpenToHiring:    teacher.IsOpenToHiring,
			Certificates:      certificatesList,
			PublicContactInfo: teacher.PublicContactInfo,
			IsMember:          centerNote.ID != 0,
			InternalRating:    centerNote.Rating,
		}

		allResults = append(allResults, result)
	}

	// 排序
	sortedResults := s.sortTalentResults(allResults, searchParams.SortBy, searchParams.SortOrder)

	// 分頁
	limit := searchParams.Limit
	if limit <= 0 {
		limit = 20 // 預設每頁 20 筆
	}
	page := searchParams.Page
	if page <= 0 {
		page = 1 // 預設第 1 頁
	}
	total := len(sortedResults)
	offset := (page - 1) * limit
	paginatedResults := sortedResults
	if offset < total {
		end := offset + limit
		if end > total {
			end = total
		}
		paginatedResults = sortedResults[offset:end]
	}

	pagination := NewPagination(searchParams.Page, searchParams.Limit, total)

	// 轉換為符合前端格式的結果
	talentCards := make([]TalentCardResponse, 0, len(paginatedResults))
	for _, result := range paginatedResults {
		// 轉換技能格式：將 Skill 改為 TalentSkillResponse（使用 skill_name）
		skills := make([]TalentSkillResponse, 0, len(result.Skills))
		for _, skill := range result.Skills {
			skills = append(skills, TalentSkillResponse{
				Category:  skill.Category,
				SkillName: skill.Name,
			})
		}

		// 轉換證照格式
		certificates := make([]CertificateResponse, 0, len(result.Certificates))
		for _, cert := range result.Certificates {
			issuedAt := ""
			if !cert.IssuedAt.IsZero() {
				issuedAt = cert.IssuedAt.Format("2006-01-02")
			}
			certificates = append(certificates, CertificateResponse{
				ID:         0, // 簡化結構，ID 非必要不返回
				Name:       cert.Name,
				Issuer:     "", // Certificate 結構沒有 issuer
				IssuedAt:   issuedAt,
				ExpiryDate: "", // Certificate 結構沒有 expiry_date
			})
		}

		talentCard := TalentCardResponse{
			ID:               result.TeacherID,
			Name:             result.Name,
			Bio:              result.Bio,
			City:             result.City,
			District:         result.District,
			Skills:           skills,
			CertificateCount:  len(result.Certificates),
			Certificates:     certificates,
			Rating:           result.InternalRating,
			ReviewCount:      0, // 暫無評價數量，設為 0
			IsOpenToHiring:   result.IsOpenToHiring,
			IsMember:         result.IsMember,
			PersonalHashtags: result.PersonalHashtags,
			PublicContactInfo: result.PublicContactInfo,
		}
		talentCards = append(talentCards, talentCard)
	}

	return &TalentSearchResultResponse{
		Talents:    talentCards,
		Pagination: pagination,
	}, nil
}

// sortTalentResults - 排序人才結果
func (s *SmartMatchingServiceImpl) sortTalentResults(results []TalentResult, sortBy, sortOrder string) []TalentResult {
	sorted := make([]TalentResult, len(results))
	copy(sorted, results)

	// 複製 slice 以避免修改原始資料
	for i := range sorted {
		sorted[i] = results[i]
	}

	sortFunc := func(i, j int) bool {
		a, b := sorted[i], sorted[j]
		var comparison int

		switch sortBy {
		case "name":
			comparison = strings.Compare(a.Name, b.Name)
		case "skills":
			comparison = len(a.Skills) - len(b.Skills)
		case "rating":
			comparison = a.InternalRating - b.InternalRating
		case "city":
			comparison = strings.Compare(a.City, b.City)
		default:
			comparison = strings.Compare(a.Name, b.Name)
		}

		if sortOrder == "desc" {
			return comparison > 0
		}
		return comparison < 0
	}

	// 使用氣泡排序（資料量小時簡單有效）
	for i := 0; i < len(sorted); i++ {
		for j := 0; j < len(sorted)-i-1; j++ {
			if sortFunc(j, j+1) {
				sorted[j], sorted[j+1] = sorted[j+1], sorted[j]
			}
		}
	}

	return sorted
}

func calculateSkillMatchScore(teacherSkills []models.TeacherSkill, requiredSkills []string) int {
	if len(requiredSkills) == 0 {
		return 10
	}

	matchCount := 0
	for _, requiredSkill := range requiredSkills {
		for _, teacherSkill := range teacherSkills {
			if teacherSkill.SkillName == requiredSkill {
				matchCount++
				break
			}
		}
	}

	return int(float64(matchCount) / float64(len(requiredSkills)) * 10)
}

func calculateHashtagScore(skills []models.TeacherSkill, requiredHashtags []string) int {
	if len(requiredHashtags) == 0 {
		return 0
	}

	hashtags := extractHashtags(skills)
	matchCount := 0
	for _, required := range requiredHashtags {
		for _, tag := range hashtags {
			if tag == required {
				matchCount++
				break
			}
		}
	}

	if matchCount > len(requiredHashtags) {
		matchCount = len(requiredHashtags)
	}

	return matchCount * 8
}

func calculateRegionMatchScore(teacherCity, teacherDistrict, searchCity, searchDistrict string) int {
	if searchCity == "" && searchDistrict == "" {
		return 10
	}
	if searchDistrict != "" && teacherDistrict == searchDistrict {
		return 10
	}
	if searchCity != "" && teacherCity == searchCity {
		return 10
	}
	return 0
}

func calculateInternalScore(rating int, internalNote string) int {
	score := 0

	if rating > 0 {
		score += rating * 6
	}

	if internalNote != "" {
		positiveKeywords := []string{"推薦", "優秀", "穩定", "配合度高", "專業"}
		for _, keyword := range positiveKeywords {
			if strings.Contains(internalNote, keyword) {
				score += 10
				break
			}
		}
	}

	return score
}

func extractSkills(skills []models.TeacherSkill) []Skill {
	result := make([]Skill, 0, len(skills))
	for _, skill := range skills {
		hashtags := make([]string, 0)
		for _, tag := range skill.Hashtags {
			if tag.Hashtag.Name != "" {
				hashtags = append(hashtags, tag.Hashtag.Name)
			}
		}
		result = append(result, Skill{
			Category: skill.Category,
			Name:     skill.SkillName,
			Hashtags: hashtags,
		})
	}
	return result
}

func extractHashtags(skills []models.TeacherSkill) []string {
	hashtagsMap := make(map[string]bool)
	for _, skill := range skills {
		for _, tag := range skill.Hashtags {
			if tag.Hashtag.Name != "" {
				hashtagsMap[tag.Hashtag.Name] = true
			}
		}
	}

	hashtags := make([]string, 0, len(hashtagsMap))
	for hashtag := range hashtagsMap {
		hashtags = append(hashtags, "#"+hashtag)
	}
	return hashtags
}

func (s *SmartMatchingServiceImpl) extractPersonalHashtags(ctx context.Context, teacherID uint) []string {
	hashtags, err := s.teacherRepository.ListPersonalHashtags(ctx, teacherID)
	if err != nil {
		s.Logger.Warn("failed to list personal hashtags", "teacher_id", teacherID, "error", err.Error())
		return []string{}
	}

	result := make([]string, 0, len(hashtags))
	for _, tag := range hashtags {
		result = append(result, "#"+tag.Name)
	}
	return result
}

// extractPersonalHashtagsBatch 從已載入的技能中提取標籤（避免 N+1）
func (s *SmartMatchingServiceImpl) extractPersonalHashtagsBatch(ctx context.Context, teacherID uint, skills []models.TeacherSkill) []string {
	hashtagMap := make(map[string]bool)

	// 從技能中提取標籤
	for _, skill := range skills {
		for _, tag := range skill.Hashtags {
			if tag.Hashtag.Name != "" {
				hashtagMap[tag.Hashtag.Name] = true
			}
		}
	}

	// 轉換為字串切片
	result := make([]string, 0, len(hashtagMap))
	for tag := range hashtagMap {
		result = append(result, "#"+tag)
	}
	return result
}

func sortMatchesByScore(matches []MatchScore) {
	for i := 0; i < len(matches); i++ {
		for j := i + 1; j < len(matches); j++ {
			if matches[i].Score < matches[j].Score {
				matches[i], matches[j] = matches[j], matches[i]
			}
		}
	}
}

// GetTalentStats - 取得人才庫統計資料
func (s *SmartMatchingServiceImpl) GetTalentStats(ctx context.Context, centerID uint) (*TalentStats, error) {
	// 取得該中心的成員數量
	memberTeacherIDs, err := s.membershipRepo.ListTeacherIDsByCenterID(ctx, centerID)
	if err != nil {
		return nil, err
	}
	memberCount := len(memberTeacherIDs)

	// 建立成員 ID Map，方便快速查找
	memberIDMap := make(map[uint]bool)
	for _, id := range memberTeacherIDs {
		memberIDMap[id] = true
	}

	// 取得所有老師（只取有必要的欄位）
	teachers, err := s.teacherRepository.List(ctx)
	if err != nil {
		return nil, err
	}

	// 過濾：只統計開放徵才且非該中心成員的老師
	var openHiringCount int
	var totalRating float64
	ratingCount := 0
	cityMap := make(map[string]int)
	skillMap := make(map[string]int)
	openHiringTeacherIDs := make([]uint, 0)

	for _, teacher := range teachers {
		// 只統計開放徵才且不是該中心成員的老師
		if !teacher.IsOpenToHiring || memberIDMap[teacher.ID] {
			continue
		}
		openHiringCount++
		openHiringTeacherIDs = append(openHiringTeacherIDs, teacher.ID)

		// 統計城市分布
		if teacher.City != "" {
			cityMap[teacher.City]++
		}

		// 取得技能
		skills, err := s.teacherSkillRepo.ListByTeacherID(ctx, teacher.ID)
		if err != nil {
			s.Logger.Warn("failed to list teacher skills", "teacher_id", teacher.ID, "error", err.Error())
		}
		for _, skill := range skills {
			skillMap[skill.SkillName]++
		}

		// 取得中心備註（用於評分）- 這裡雖然老師不是成員，但可能有之前的評分
		centerNote, err := s.centerTeacherNoteRepo.GetByCenterAndTeacher(ctx, centerID, teacher.ID)
		if err != nil {
			s.Logger.Warn("failed to get center teacher note", "teacher_id", teacher.ID, "error", err.Error())
		}
		if centerNote.Rating > 0 {
			totalRating += float64(centerNote.Rating)
			ratingCount++
		}
	}

	// 批量查詢成員的技能和評分（提升效能）
	if len(memberTeacherIDs) > 0 {
		memberSkillsMap, err := s.teacherSkillRepo.BatchListByTeacherIDs(ctx, memberTeacherIDs)
		if err != nil {
			s.Logger.Warn("failed to batch query member skills", "error", err.Error())
		}
		memberNotesMap, err := s.centerTeacherNoteRepo.BatchGetByCenterAndTeachers(ctx, centerID, memberTeacherIDs)
		if err != nil {
			s.Logger.Warn("failed to batch query member notes", "error", err.Error())
		}

		for _, teacherID := range memberTeacherIDs {
			// 成員也納入技能統計
			if skills, ok := memberSkillsMap[teacherID]; ok {
				for _, skill := range skills {
					skillMap[skill.SkillName]++
				}
			}
			// 成員評分
			if note, ok := memberNotesMap[teacherID]; ok && note.Rating > 0 {
				totalRating += float64(note.Rating)
				ratingCount++
			}
		}
	}

	// 計算城市分布
	cityDistribution := make([]CityDistributionItem, 0, len(cityMap))
	for city, count := range cityMap {
		cityDistribution = append(cityDistribution, CityDistributionItem{
			Name:  city,
			Count: count,
		})
	}
	// 按數量排序
	for i := 0; i < len(cityDistribution); i++ {
		for j := i + 1; j < len(cityDistribution); j++ {
			if cityDistribution[i].Count < cityDistribution[j].Count {
				cityDistribution[i], cityDistribution[j] = cityDistribution[j], cityDistribution[i]
			}
		}
	}
	if len(cityDistribution) > 5 {
		cityDistribution = cityDistribution[:5]
	}

	// 計算熱門技能
	topSkills := make([]SkillCountItem, 0, len(skillMap))
	for skill, count := range skillMap {
		topSkills = append(topSkills, SkillCountItem{
			Name:  skill,
			Count: count,
		})
	}
	// 按數量排序
	for i := 0; i < len(topSkills); i++ {
		for j := i + 1; j < len(topSkills); j++ {
			if topSkills[i].Count < topSkills[j].Count {
				topSkills[i], topSkills[j] = topSkills[j], topSkills[i]
			}
		}
	}
	if len(topSkills) > 5 {
		topSkills = topSkills[:5]
	}

	// 計算平均評分
	var averageRating float64
	if ratingCount > 0 {
		averageRating = totalRating / float64(ratingCount)
	}

	// 從資料庫取得真實的邀請統計
	pendingInvites, acceptedInvites, declinedInvites, err := s.centerInvitationRepo.CountByCenter(ctx, centerID)
	if err != nil {
		// 如果取得失敗，使用預設值
		pendingInvites, acceptedInvites, declinedInvites = 0, 0, 0
	}

	// 計算月趨勢（簡化版：最近6個月的變化）
	monthlyTrend := []int{
		openHiringCount * 80 / 100,
		openHiringCount * 85 / 100,
		openHiringCount * 88 / 100,
		openHiringCount * 90 / 100,
		openHiringCount * 95 / 100,
		openHiringCount,
	}

	// 計算月變化
	monthlyChange := 0
	if len(monthlyTrend) > 1 {
		prevMonth := monthlyTrend[len(monthlyTrend)-2]
		if prevMonth > 0 {
			monthlyChange = int(float64(openHiringCount-prevMonth) / float64(prevMonth) * 100)
		}
	}

	return &TalentStats{
		TotalCount:       openHiringCount + memberCount, // 可搜尋的總人數（開放徵才 + 中心成員）
		OpenHiringCount:  openHiringCount,
		MemberCount:      memberCount,
		AverageRating:    averageRating,
		MonthlyChange:    monthlyChange,
		MonthlyTrend:     monthlyTrend,
		PendingInvites:   int(pendingInvites),
		AcceptedInvites:  int(acceptedInvites),
		DeclinedInvites:  int(declinedInvites),
		CityDistribution: cityDistribution,
		TopSkills:        topSkills,
	}, nil
}

// InviteTalent - 邀請人才合作
func (s *SmartMatchingServiceImpl) InviteTalent(ctx context.Context, centerID uint, adminID uint, teacherIDs []uint, message string) (*InviteResult, error) {
	result := &InviteResult{
		InvitedCount:  0,
		FailedCount:   0,
		FailedIDs:     make([]uint, 0),
		InvitationIDs: make([]uint, 0),
	}

	// 生成邀請碼
	generateToken := func() string {
		return fmt.Sprintf("INV-%d-%s", centerID, strconv.FormatInt(time.Now().UnixNano(), 36))
	}

	// 邀請過期時間（7天後）
	expiresAt := time.Now().Add(7 * 24 * time.Hour)
	expiresAtPtr := &expiresAt

	for _, teacherID := range teacherIDs {
		// 檢查老師是否存在且開放徵才
		teacher, err := s.teacherRepository.GetByID(ctx, teacherID)
		if err != nil {
			result.FailedCount++
			result.FailedIDs = append(result.FailedIDs, teacherID)
			continue
		}

		if !teacher.IsOpenToHiring {
			result.FailedCount++
			result.FailedIDs = append(result.FailedIDs, teacherID)
			continue
		}

		// 檢查是否已有待處理的邀請
		hasPending, err := s.centerInvitationRepo.HasPendingInvitation(ctx, teacherID, centerID)
		if err != nil {
			result.FailedCount++
			result.FailedIDs = append(result.FailedIDs, teacherID)
			continue
		}
		if hasPending {
			// 已有待處理邀請，標記為失敗
			result.FailedCount++
			result.FailedIDs = append(result.FailedIDs, teacherID)
			continue
		}

		// 建立邀請記錄
		invitation := models.CenterInvitation{
			CenterID:   centerID,
			TeacherID:  teacherID,
			InvitedBy:  adminID,
			Email:      teacher.Email,
			Token:      generateToken(),
			Status:     models.InvitationStatusPending,
			InviteType: models.InvitationTypeTalentPool,
			Message:    message,
			ExpiresAt:  expiresAtPtr,
			CreatedAt:  time.Now(),
		}

		// 儲存到資料庫
		created, err := s.centerInvitationRepo.Create(ctx, invitation)
		if err != nil {
			result.FailedCount++
			result.FailedIDs = append(result.FailedIDs, teacherID)
			continue
		}

		result.InvitedCount++
		result.InvitationIDs = append(result.InvitationIDs, created.ID)

		// 發送 LINE 通知給老師（如果有 LineUserID）
		// 使用 goroutine 非同步發送，不影響主流程
		if teacher.LineUserID != "" {
			go func(tID uint, tName string, token string) {
				// 取得中心名稱（簡化處理，實際應該查詢中心資料）
				centerName := "TimeLedger 中心"

				// 發送人才庫邀請通知
				if err := s.notificationService.SendTalentInvitationNotification(context.Background(), tID, centerName, token); err != nil {
					fmt.Printf("[ERROR] 發送人才庫邀請通知失敗: %v (teacher_id=%d)\n", err, tID)
				} else {
					fmt.Printf("[INFO] 人才庫邀請通知已發送給 %s\n", tName)
				}
			}(teacherID, teacher.Name, created.Token)
		}
	}

	return result, nil
}

// GetSearchSuggestions - 取得搜尋建議（優化版本：加入 centerID 並修復 N+1）
func (s *SmartMatchingServiceImpl) GetSearchSuggestions(ctx context.Context, centerID uint, query string) (*SearchSuggestions, error) {
	suggestions := &SearchSuggestions{
		Skills:   make([]SuggestionItem, 0),
		Tags:     make([]SuggestionItem, 0),
		Names:    make([]SuggestionItem, 0),
		Trending: []string{"瑜珈", "鋼琴", "舞蹈", "美術", "英語"},
	}

	// 取得該中心的成員 ID 列表
	memberIDMap := make(map[uint]bool)
	memberTeacherIDs, err := s.membershipRepo.ListTeacherIDsByCenterID(ctx, centerID)
	if err == nil {
		for _, id := range memberTeacherIDs {
			memberIDMap[id] = true
		}
	}

	// 取得所有開放徵才的老師
	teachers, err := s.teacherRepository.List(ctx)
	if err != nil {
		return nil, err
	}

	// 收集候選老師 ID（開放徵才且不是該中心成員）
	candidateIDs := make([]uint, 0)
	for _, teacher := range teachers {
		if teacher.IsOpenToHiring && !memberIDMap[teacher.ID] {
			candidateIDs = append(candidateIDs, teacher.ID)
		}
	}

	// 批量查詢技能（解決 N+1 問題）
	skillsMap, err := s.teacherSkillRepo.BatchListByTeacherIDs(ctx, candidateIDs)
	if err != nil {
		s.Logger.Warn("failed to batch query skills for suggestions", "error", err.Error())
	}

	skillSet := make(map[string]bool)
	tagSet := make(map[string]bool)
	nameSet := make(map[string]bool)

	queryLower := strings.ToLower(query)

	for _, teacher := range teachers {
		// 排除該中心成員
		if memberIDMap[teacher.ID] {
			continue
		}

		// 只處理開放徵才的老師
		if !teacher.IsOpenToHiring {
			continue
		}

		// 姓名匹配
		if query != "" && strings.Contains(strings.ToLower(teacher.Name), queryLower) && !nameSet[teacher.Name] {
			nameSet[teacher.Name] = true
			suggestions.Names = append(suggestions.Names, SuggestionItem{
				Type:  "name",
				Value: teacher.Name,
			})
		}

		// 技能匹配（使用批量查詢結果）
		if query != "" {
			skills := skillsMap[teacher.ID]
			for _, skill := range skills {
				if strings.Contains(strings.ToLower(skill.SkillName), queryLower) && !skillSet[skill.SkillName] {
					skillSet[skill.SkillName] = true
					suggestions.Skills = append(suggestions.Skills, SuggestionItem{
						Type:  "skill",
						Value: skill.SkillName,
					})
				}
			}

			// 標籤匹配（從技能中提取）
			for _, skill := range skills {
				for _, tag := range skill.Hashtags {
					tagName := tag.Hashtag.Name
					if tagName == "" {
						continue
					}
					tagClean := strings.TrimPrefix(tagName, "#")
					if strings.Contains(strings.ToLower(tagClean), queryLower) && !tagSet[tagClean] {
						tagSet[tagClean] = true
						suggestions.Tags = append(suggestions.Tags, SuggestionItem{
							Type:  "tag",
							Value: tagClean,
						})
					}
				}
			}
		}
	}

	// 限制數量
	if len(suggestions.Skills) > 5 {
		suggestions.Skills = suggestions.Skills[:5]
	}
	if len(suggestions.Tags) > 5 {
		suggestions.Tags = suggestions.Tags[:5]
	}
	if len(suggestions.Names) > 5 {
		suggestions.Names = suggestions.Names[:5]
	}

	return suggestions, nil
}

// GetAlternativeSlots - 取得替代時段建議
func (s *SmartMatchingServiceImpl) GetAlternativeSlots(ctx context.Context, centerID uint, teacherID uint, originalStart, originalEnd time.Time, duration int) ([]AlternativeSlot, error) {
	if duration == 0 {
		duration = 90 // 預設 90 分鐘
	}

	slots := make([]AlternativeSlot, 0)

	// 取得該時段區間的規則
	rules, err := s.scheduleRuleRepo.ListByCenterID(ctx, centerID)
	if err != nil {
		return nil, err
	}

	// 產生未來 7 天的替代時段
	for day := 1; day <= 7; day++ {
		date := originalStart.AddDate(0, 0, day)
		dateStr := date.Format("2006-01-02")

		// 產生三個時段：上午、下午、晚上
		timeSlots := []string{"09:00", "14:00", "16:00"}

		for _, startTime := range timeSlots {
			hourParts := strings.Split(startTime, ":")
			hour, err := strconv.Atoi(hourParts[0])
			if err != nil {
				s.Logger.Warn("failed to parse hour", "start_time", startTime, "error", err.Error())
				continue
			}
			endHour := hour + duration/60
			endTime := fmt.Sprintf("%02d:30", endHour)

			// 檢查是否與現有課表衝突
			hasConflict := false
			for _, rule := range rules {
				if rule.TeacherID == nil || *rule.TeacherID != teacherID {
					continue
				}

				// 檢查日期和時間是否衝突
				if int(date.Weekday()) == rule.Weekday {
					ruleStart := rule.StartTime
					ruleEnd := rule.EndTime

					// 簡單的時間衝突檢查
					if startTime < ruleEnd && endTime > ruleStart {
						hasConflict = true
						break
					}
				}
			}

			slot := AlternativeSlot{
				Date:      dateStr,
				DateLabel: fmt.Sprintf("%d/%d", date.Month(), date.Day()),
				Start:     startTime,
				End:       endTime,
				Available: !hasConflict,
			}

			if hasConflict {
				slot.ConflictReason = "與現有課程衝突"
			}

			slots = append(slots, slot)
		}
	}

	return slots, nil
}

// GetTeacherSessions - 取得教師課表
func (s *SmartMatchingServiceImpl) GetTeacherSessions(ctx context.Context, centerID uint, teacherID uint, startDate, endDate string) (*TeacherSessions, error) {
	// 取得老師資料
	teacher, err := s.teacherRepository.GetByID(ctx, teacherID)
	if err != nil {
		return nil, err
	}

	// 解析日期
	start, err := time.Parse("2006-01-02", startDate)
	if err != nil {
		return nil, fmt.Errorf("無效的開始日期格式: %w", err)
	}
	end, err := time.Parse("2006-01-02", endDate)
	if err != nil {
		return nil, fmt.Errorf("無效的結束日期格式: %w", err)
	}

	// 取得該時段區間的規則
	rules, err := s.scheduleRuleRepo.ListByCenterID(ctx, centerID)
	if err != nil {
		return nil, err
	}

	sessions := make([]SessionInfo, 0)

	for _, rule := range rules {
		if rule.TeacherID == nil || *rule.TeacherID != teacherID {
			continue
		}

		// 檢查規則是否在日期範圍內
		if rule.EffectiveRange.StartDate.After(end) || rule.EffectiveRange.EndDate.Before(start) {
			continue
		}

		// 展開循環規則產生具體日期的課表
		currentDate := start
		for !currentDate.After(end) {
			// 檢查是否為規則的星期
			if int(currentDate.Weekday()) == rule.Weekday {
				// 檢查是否有取消例外
				exceptions, err := s.scheduleExceptionRepo.GetByRuleAndDate(ctx, rule.ID, currentDate)
				if err != nil {
					// 查詢例外失敗時，保守起見假設沒有取消
					s.Logger.Warn("failed to query schedule exceptions", "rule_id", rule.ID, "date", currentDate.Format("2006-01-02"), "error", err.Error())
				}
				hasCancel := false
				for _, exc := range exceptions {
					if exc.ExceptionType == "CANCEL" && exc.Status == "APPROVED" {
						hasCancel = true
						break
					}
				}

				if !hasCancel {
					sessionDate := currentDate.Format("2006-01-02")
					startDateTime := fmt.Sprintf("%sT%s:00", sessionDate, rule.StartTime)
					endDateTime := fmt.Sprintf("%sT%s:00", sessionDate, rule.EndTime)

					sessions = append(sessions, SessionInfo{
						ID:         rule.ID,
						CourseName: rule.Name,
						StartTime:  startDateTime,
						EndTime:    endDateTime,
						RoomName:   rule.Room.Name,
						Status:     "SCHEDULED",
					})
				}
			}

			currentDate = currentDate.AddDate(0, 0, 1)
		}
	}

	return &TeacherSessions{
		TeacherID:   teacherID,
		TeacherName: teacher.Name,
		Sessions:    sessions,
	}, nil
}

type MatchScore struct {
	TeacherID         uint              `json:"teacher_id"`
	Name              string            `json:"name"`
	Score             int               `json:"score"`
	Availability      MatchAvailability `json:"availability"`
	AvailabilityScore int               `json:"availability_score"`
	InternalScore     int               `json:"internal_score"`
	SkillScore        int               `json:"skill_score"`
	RegionScore       int               `json:"region_score"`
	Skills            []Skill           `json:"skills,omitempty"`
	Hashtags          []string          `json:"hashtags,omitempty"`
	IsMember          bool              `json:"is_member"`
	InternalNote      string            `json:"internal_note,omitempty"`
	InternalRating    int               `json:"internal_rating,omitempty"`
	PublicContactInfo string            `json:"public_contact_info,omitempty"`
}

type TalentResult struct {
	TeacherID         uint          `json:"teacher_id"`
	Name              string        `json:"name"`
	Bio               string        `json:"bio,omitempty"`
	City              string        `json:"city,omitempty"`
	District          string        `json:"district,omitempty"`
	Skills            []Skill       `json:"skills,omitempty"`
	PersonalHashtags  []string      `json:"personal_hashtags,omitempty"`
	IsOpenToHiring    bool          `json:"is_open_to_hiring"`
	Certificates      []Certificate `json:"certificates,omitempty"`
	PublicContactInfo string        `json:"public_contact_info,omitempty"`
	IsMember          bool          `json:"is_member"`       // 是否為中心成員
	InternalRating    int           `json:"internal_rating"` // 中心內部評分
}

type Skill struct {
	Category string   `json:"category"`
	Name     string   `json:"name"`
	Hashtags []string `json:"hashtags,omitempty"`
}

type Certificate struct {
	Name     string    `json:"name"`
	FileURL  string    `json:"file_url"`
	IssuedAt time.Time `json:"issued_at"`
}

type MatchService interface {
	FindMatches(ctx context.Context, centerID uint, teacherID *uint, roomIDs []uint, startTime, endTime time.Time, requiredSkills []string, excludeTeacherIDs []uint) ([]MatchScore, error)
	SearchTalent(ctx context.Context, searchParams TalentSearchParams) (*TalentSearchResultResponse, error)
}

var _ MatchService = (*SmartMatchingServiceImpl)(nil)
