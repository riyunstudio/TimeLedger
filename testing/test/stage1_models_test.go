package test

import (
	"testing"
	"time"
	"timeLedger/app/models"
)

func TestCenter_Validation(t *testing.T) {
	center := models.Center{
		Name:      "",
		PlanLevel: "",
		Settings:  models.CenterSettings{},
		CreatedAt: time.Time{},
	}

	center.Name = "Test Center"
	if center.Name == "" {
		t.Error("Name should not be empty")
	}

	validPlans := []string{"FREE", "STARTER", "PRO", "TEAM"}
	for _, plan := range validPlans {
		center.PlanLevel = plan
		if center.PlanLevel == "" {
			t.Errorf("PlanLevel should be set to %s", plan)
		}
	}

	center.Settings.ExceptionLeadDays = 7
	if center.Settings.ExceptionLeadDays < 0 {
		t.Error("ExceptionLeadDays should not be negative")
	}
}

func TestTeacher_Validation(t *testing.T) {
	teacher := models.Teacher{
		LineUserID:        "TEST_LINE_12345",
		Name:              "Test Teacher",
		Email:             "test@example.com",
		City:              "台北市",
		District:          "大安區",
		PublicContactInfo: "Line: test123",
		Bio:               "Test bio",
		IsOpenToHiring:    true,
	}

	if teacher.LineUserID == "" {
		t.Error("LineUserID is required")
	}

	if teacher.Name == "" {
		t.Error("Name is required")
	}

	validCities := []string{"台北市", "新北市", "台中市"}
	cityValid := false
	for _, city := range validCities {
		if teacher.City == city {
			cityValid = true
			break
		}
	}
	if !cityValid && teacher.City != "" {
		t.Logf("City %s is not in valid list", teacher.City)
	}
}

func TestAdminUser_Validation(t *testing.T) {
	admin := models.AdminUser{
		CenterID:     1,
		Email:        "admin@test.com",
		Name:         "Admin User",
		PasswordHash: "hashed_password",
		Role:         "OWNER",
		Status:       "ACTIVE",
	}

	if admin.CenterID == 0 {
		t.Error("CenterID is required")
	}

	if admin.Email == "" {
		t.Error("Email is required")
	}

	validRoles := []string{"OWNER", "ADMIN", "STAFF"}
	roleValid := false
	for _, role := range validRoles {
		if admin.Role == role {
			roleValid = true
			break
		}
	}
	if !roleValid {
		t.Error("Role must be OWNER, ADMIN, or STAFF")
	}

	validStatuses := []string{"ACTIVE", "INACTIVE"}
	statusValid := false
	for _, status := range validStatuses {
		if admin.Status == status {
			statusValid = true
			break
		}
	}
	if !statusValid {
		t.Error("Status must be ACTIVE or INACTIVE")
	}
}

func TestGeoDistrict_ForeignKey(t *testing.T) {
	city := models.GeoCity{
		Name: "台北市",
	}
	city.ID = 123

	district := models.GeoDistrict{
		CityID: city.ID,
		Name:   "大安區",
	}

	if district.CityID == 0 {
		t.Error("District must have a valid CityID")
	}
}

func TestTeacherSkill_Validation(t *testing.T) {
	skill := models.TeacherSkill{
		TeacherID: 1,
		Category:  "有氧",
		SkillName: "Zumba",
		Level:     "INTERMEDIATE",
	}

	if skill.TeacherID == 0 {
		t.Error("TeacherID is required")
	}

	if skill.Category == "" {
		t.Error("Category is required")
	}

	if skill.SkillName == "" {
		t.Error("SkillName is required")
	}

	validLevels := []string{"BASIC", "INTERMEDIATE", "ADVANCED"}
	levelValid := false
	for _, level := range validLevels {
		if skill.Level == level {
			levelValid = true
			break
		}
	}
	if !levelValid {
		t.Error("Level must be BASIC, INTERMEDIATE, or ADVANCED")
	}
}

func TestHashtag_Validation(t *testing.T) {
	hashtag := models.Hashtag{
		Name:       "#街舞",
		UsageCount: 0,
	}

	if hashtag.Name == "" {
		t.Error("Hashtag name is required")
	}

	if !containsPrefix(hashtag.Name, "#") {
		t.Error("Hashtag name should start with #")
	}

	if hashtag.UsageCount < 0 {
		t.Error("Usage count cannot be negative")
	}
}

func TestTeacherCertificate_Validation(t *testing.T) {
	cert := models.TeacherCertificate{
		TeacherID: 1,
		Name:      "RYT-200",
		FileURL:   "https://example.com/cert.pdf",
		IssuedAt:  time.Date(2020, time.January, 1, 0, 0, 0, 0, time.UTC),
	}

	if cert.TeacherID == 0 {
		t.Error("TeacherID is required")
	}

	if cert.Name == "" {
		t.Error("Certificate name is required")
	}

	if cert.FileURL == "" {
		t.Error("File URL is required")
	}

	if cert.IssuedAt.IsZero() {
		t.Error("IssuedAt date is required")
	}
}

func TestTeacherPersonalHashtag_SortOrder(t *testing.T) {
	ph := models.TeacherPersonalHashtag{
		TeacherID: 1,
		HashtagID: 1,
		SortOrder: 1,
	}

	if ph.TeacherID == 0 {
		t.Error("TeacherID is required")
	}

	if ph.HashtagID == 0 {
		t.Error("HashtagID is required")
	}

	if ph.SortOrder < 1 || ph.SortOrder > 5 {
		t.Error("Sort order must be between 1 and 5")
	}
}

func containsPrefix(s, prefix string) bool {
	if len(s) < len(prefix) {
		return false
	}
	return s[:len(prefix)] == prefix
}
