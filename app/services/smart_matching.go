package services

import (
	"context"
	"strings"
	"time"
	"timeLedger/app"
	"timeLedger/app/models"
	"timeLedger/app/repositories"
)

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
	}
}

// FindMatches searches for available teachers matching the criteria
func (s *SmartMatchingServiceImpl) FindMatches(ctx context.Context, centerID uint, teacherID *uint, roomID uint, startTime, endTime time.Time, requiredSkills []string, excludeTeacherIDs []uint) ([]MatchScore, error) {
	var matches []MatchScore

	rules, err := s.scheduleRuleRepo.ListByCenterID(ctx, centerID)
	if err != nil {
		return nil, err
	}

	for _, rule := range rules {
		if rule.TeacherID == nil {
			continue
		}

		currentTeacherID := *rule.TeacherID

		isExcluded := false
		for _, excludedID := range excludeTeacherIDs {
			if currentTeacherID == excludedID {
				isExcluded = true
				break
			}
		}

		if isExcluded {
			continue
		}

		exceptions, _ := s.scheduleExceptionRepo.GetByRuleAndDate(ctx, rule.ID, startTime)
		hasCancel := false
		for _, exc := range exceptions {
			if exc.Type == "CANCEL" && exc.Status == "APPROVED" {
				hasCancel = true
				break
			}
		}

		if hasCancel {
			continue
		}

		teacher, err := s.teacherRepository.GetByID(ctx, currentTeacherID)
		if err != nil {
			continue
		}

		if !teacher.IsOpenToHiring {
			continue
		}

		skills, err := s.teacherSkillRepo.ListByTeacherID(ctx, currentTeacherID)
		if err != nil {
			continue
		}

		centerNote, _ := s.centerTeacherNoteRepo.GetByCenterAndTeacher(ctx, centerID, currentTeacherID)

		availability := MatchAvailable
		if teacherID != nil && currentTeacherID == *teacherID && roomID == rule.RoomID {
			availability = MatchOverlap
		}

		skillScore := calculateSkillMatchScore(skills, requiredSkills)
		regionScore := calculateRegionMatchScore(teacher.City, teacher.District, "", "")
		hashtagScore := calculateHashtagScore(skills, []string{})
		internalScore := calculateInternalScore(centerNote.Rating, centerNote.InternalNote)

		availabilityScore := 0
		switch availability {
		case MatchAvailable:
			availabilityScore = 40
		case MatchBufferConflict:
			availabilityScore = 15
		case MatchOverlap:
			availabilityScore = 0
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
			IsMember:          centerNote.ID != 0,
			InternalNote:      centerNote.InternalNote,
			InternalRating:    centerNote.Rating,
			PublicContactInfo: teacher.PublicContactInfo,
		}

		matches = append(matches, match)
	}

	sortMatchesByScore(matches)

	return matches, nil
}

// SearchTalent searches for teachers based on criteria
func (s *SmartMatchingServiceImpl) SearchTalent(ctx context.Context, searchParams TalentSearchParams) ([]TalentResult, error) {
	var results []TalentResult

	teachers, err := s.teacherRepository.List(ctx)
	if err != nil {
		return nil, err
	}

	for _, teacher := range teachers {
		if !teacher.IsOpenToHiring {
			continue
		}

		if searchParams.CenterID > 0 {
			// Filter by center membership if CenterID is specified
			// This is a simplified check - in production you'd query membership
		}

		if searchParams.City != "" && teacher.City != searchParams.City {
			continue
		}

		if searchParams.District != "" && teacher.District != searchParams.District {
			continue
		}

		if searchParams.Keyword != "" {
			if !strings.Contains(teacher.Name, searchParams.Keyword) && !strings.Contains(teacher.Bio, searchParams.Keyword) {
				continue
			}
		}

		skillsList, err := s.teacherSkillRepo.ListByTeacherID(ctx, teacher.ID)
		if err != nil {
			continue
		}

		if len(searchParams.Skills) > 0 {
			hasSkill := false
			for _, requiredSkill := range searchParams.Skills {
				for _, teacherSkill := range skillsList {
					if teacherSkill.SkillName == requiredSkill || teacherSkill.Category == requiredSkill {
						hasSkill = true
						break
					}
				}
			}
			if !hasSkill {
				continue
			}
		}

		// 標籤篩選
		if len(searchParams.Hashtags) > 0 {
			personalHashtags := s.extractPersonalHashtags(ctx, teacher.ID)
			hasMatchingHashtag := false
			
			for _, requiredTag := range searchParams.Hashtags {
				// 移除 # 符號進行比對
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

		certificates, _ := s.teacherCertificateRepo.ListByTeacherID(ctx, teacher.ID)

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
			PersonalHashtags:  s.extractPersonalHashtags(ctx, teacher.ID),
			IsOpenToHiring:    teacher.IsOpenToHiring,
			Certificates:      certificatesList,
			PublicContactInfo: teacher.PublicContactInfo,
		}

		results = append(results, result)
	}

	return results, nil
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
		hashtags := make([]string, 0, 0)
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
	hashtags, _ := s.teacherRepository.ListPersonalHashtags(ctx, teacherID)

	result := make([]string, 0, len(hashtags))
	for _, tag := range hashtags {
		result = append(result, "#"+tag.Name)
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
	FindMatches(ctx context.Context, centerID uint, teacherID *uint, roomID uint, startTime, endTime time.Time, requiredSkills []string, excludeTeacherIDs []uint) ([]MatchScore, error)
	SearchTalent(ctx context.Context, searchParams TalentSearchParams) ([]TalentResult, error)
}

var _ MatchService = (*SmartMatchingServiceImpl)(nil)
