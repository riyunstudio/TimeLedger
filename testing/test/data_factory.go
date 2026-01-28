package test

import (
	"context"
	"fmt"
	"time"

	"timeLedger/app/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// TestDataFactory provides helper functions for creating consistent test data.
// This ensures unique constraints and foreign key relationships are properly satisfied.
type TestDataFactory struct {
	db          *gorm.DB
	uniqueSeq   int
	baseTime    time.Time
}

// NewTestDataFactory creates a new TestDataFactory instance.
func NewTestDataFactory(db *gorm.DB) *TestDataFactory {
	return &TestDataFactory{
		db:        db,
		uniqueSeq: 0,
		baseTime:  time.Now(),
	}
}

// generateUniqueID generates a unique identifier using UUID and sequence.
// This ensures uniqueness across multiple test runs.
func (f *TestDataFactory) generateUniqueID() string {
	f.uniqueSeq++
	uuid := uuid.New().String()
	return fmt.Sprintf("%s-%d-%d", uuid[:8], time.Now().UnixNano()%10000, f.uniqueSeq)
}

// CreateUniqueLineUserID generates a unique LINE user ID for test teachers.
// This prevents unique index conflicts on the teachers.line_user_id column.
func (f *TestDataFactory) CreateUniqueLineUserID() string {
	return fmt.Sprintf("line-%s", f.generateUniqueID())
}

// CreateUniqueEmail generates a unique email address for test entities.
func (f *TestDataFactory) CreateUniqueEmail(prefix string) string {
	return fmt.Sprintf("%s.%s@test.com", prefix, f.generateUniqueID())
}

// CreateUniqueToken generates a unique token for invitations.
func (f *TestDataFactory) CreateUniqueToken() string {
	return fmt.Sprintf("TOKEN-%s", f.generateUniqueID())
}

// CreateTestCenter creates a test center with a unique name.
func (f *TestDataFactory) CreateTestCenter(ctx context.Context, namePrefix string) (*models.Center, error) {
	center := &models.Center{
		Name:      fmt.Sprintf("測試中心 - %s - %s", namePrefix, f.generateUniqueID()[:8]),
		PlanLevel: "STARTER",
		CreatedAt: f.baseTime,
	}
	if err := f.db.WithContext(ctx).Create(center).Error; err != nil {
		return nil, fmt.Errorf("建立測試中心失敗: %w", err)
	}
	return center, nil
}

// CreateTestTeacher creates a test teacher with required fields properly set.
// This ensures LineUserID is unique to avoid index conflicts.
func (f *TestDataFactory) CreateTestTeacher(ctx context.Context, namePrefix string, opts ...TeacherOption) (*models.Teacher, error) {
	teacher := &models.Teacher{
		Name:       fmt.Sprintf("測試老師 - %s - %s", namePrefix, f.generateUniqueID()[:8]),
		Email:      f.CreateUniqueEmail(namePrefix),
		LineUserID: f.CreateUniqueLineUserID(), // Ensure unique LineUserID
		CreatedAt:  f.baseTime,
		UpdatedAt:  f.baseTime,
	}

	// Apply options
	for _, opt := range opts {
		opt(teacher)
	}

	if err := f.db.WithContext(ctx).Create(teacher).Error; err != nil {
		return nil, fmt.Errorf("建立測試老師失敗: %w", err)
	}
	return teacher, nil
}

// TeacherOption is a function type for configuring test teacher options.
type TeacherOption func(*models.Teacher)

// WithTeacherCity sets the city for a test teacher.
func WithTeacherCity(city string) TeacherOption {
	return func(t *models.Teacher) {
		t.City = city
	}
}

// WithTeacherDistrict sets the district for a test teacher.
func WithTeacherDistrict(district string) TeacherOption {
	return func(t *models.Teacher) {
		t.District = district
	}
}

// WithTeacherOpenToHiring sets the IsOpenToHiring flag for a test teacher.
func WithTeacherOpenToHiring(open bool) TeacherOption {
	return func(t *models.Teacher) {
		t.IsOpenToHiring = open
	}
}

// WithTeacherBio sets the bio for a test teacher.
func WithTeacherBio(bio string) TeacherOption {
	return func(t *models.Teacher) {
		t.Bio = bio
	}
}

// CreateTestCourse creates a test course for the center.
// This is required before creating Offerings that reference course_id.
func (f *TestDataFactory) CreateTestCourse(ctx context.Context, centerID uint, namePrefix string) (*models.Course, error) {
	course := &models.Course{
		CenterID:         centerID,
		Name:             fmt.Sprintf("課程 - %s - %s", namePrefix, f.generateUniqueID()[:8]),
		DefaultDuration:  60,
		ColorHex:         "#3498db",
		RoomBufferMin:    10,
		TeacherBufferMin: 10,
		IsActive:         true,
		CreatedAt:        f.baseTime,
		UpdatedAt:        f.baseTime,
	}
	if err := f.db.WithContext(ctx).Create(course).Error; err != nil {
		return nil, fmt.Errorf("建立測試課程失敗: %w", err)
	}
	return course, nil
}

// CreateTestOffering creates a test offering (course plan) for the center.
// This is required before creating ScheduleRules that reference offering_id.
func (f *TestDataFactory) CreateTestOffering(ctx context.Context, centerID uint, courseID uint, namePrefix string) (*models.Offering, error) {
	offering := &models.Offering{
		CenterID:            centerID,
		CourseID:            courseID,
		Name:                fmt.Sprintf("方案 - %s - %s", namePrefix, f.generateUniqueID()[:8]),
		AllowBufferOverride: false,
		IsActive:            true,
		CreatedAt:           f.baseTime,
		UpdatedAt:           f.baseTime,
	}
	if err := f.db.WithContext(ctx).Create(offering).Error; err != nil {
		return nil, fmt.Errorf("建立測試方案失敗: %w", err)
	}
	return offering, nil
}

// CreateTestRoom creates a test room for the center.
func (f *TestDataFactory) CreateTestRoom(ctx context.Context, centerID uint, namePrefix string) (*models.Room, error) {
	room := &models.Room{
		CenterID:  centerID,
		Name:      fmt.Sprintf("教室 - %s - %s", namePrefix, f.generateUniqueID()[:8]),
		Capacity:  10,
		CreatedAt: f.baseTime,
		UpdatedAt: f.baseTime,
	}
	if err := f.db.WithContext(ctx).Create(room).Error; err != nil {
		return nil, fmt.Errorf("建立測試教室失敗: %w", err)
	}
	return room, nil
}

// CreateTestScheduleRule creates a test schedule rule with properly set foreign keys.
// This ensures that offering_id and room_id reference existing records.
func (f *TestDataFactory) CreateTestScheduleRule(ctx context.Context, centerID uint, offeringID, roomID uint, teacherID *uint, namePrefix string) (*models.ScheduleRule, error) {
	rule := &models.ScheduleRule{
		CenterID:       centerID,
		OfferingID:     offeringID,
		TeacherID:      teacherID,
		RoomID:         roomID,
		Name:           fmt.Sprintf("規則 - %s - %s", namePrefix, f.generateUniqueID()[:8]),
		Weekday:        1,
		StartTime:      "09:00",
		EndTime:        "10:00",
		Duration:       60,
		EffectiveRange: models.DateRange{
			StartDate: time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC),
			EndDate:   time.Date(2026, 12, 31, 0, 0, 0, 0, time.UTC),
		},
		CreatedAt:  f.baseTime,
		UpdatedAt:  f.baseTime,
	}
	if err := f.db.WithContext(ctx).Create(rule).Error; err != nil {
		return nil, fmt.Errorf("建立測試規則失敗: %w", err)
	}
	return rule, nil
}

// CleanupTestData removes test data created by the factory.
// This ensures proper cleanup order to avoid foreign key constraint violations.
func (f *TestDataFactory) CleanupTestData(ctx context.Context) error {
	// Clean up in reverse order of dependencies
	// Note: Must delete child tables before parent tables to avoid foreign key constraint violations
	tables := []string{
		"schedule_exceptions",
		"schedule_rules",
		"center_invitations",
		"center_memberships",
		"center_teacher_notes",
		"personal_events",
		"teachers",
		"centers",
		"offerings",
		"courses",
		"rooms",
	}

	for _, table := range tables {
		if err := f.db.WithContext(ctx).Exec(fmt.Sprintf("DELETE FROM %s WHERE id > 0", table)).Error; err != nil {
			// Log but continue - some tables may not exist in test environment
			fmt.Printf("[WARN] Cleanup table %s failed: %v\n", table, err)
		}
	}
	return nil
}

// SetupCompleteTestData creates a complete set of related test data.
// This is useful for tests that need multiple related entities.
func (f *TestDataFactory) SetupCompleteTestData(ctx context.Context, testName string) (*CompleteTestData, error) {
	// Create center
	center, err := f.CreateTestCenter(ctx, testName)
	if err != nil {
		return nil, err
	}

	// Create course
	course, err := f.CreateTestCourse(ctx, center.ID, testName)
	if err != nil {
		// Cleanup on failure
		f.db.WithContext(ctx).Where("id = ?", center.ID).Delete(&models.Center{})
		return nil, err
	}

	// Create offering
	offering, err := f.CreateTestOffering(ctx, center.ID, course.ID, testName)
	if err != nil {
		f.db.WithContext(ctx).Where("id = ?", center.ID).Delete(&models.Center{})
		f.db.WithContext(ctx).Where("center_id = ?", center.ID).Delete(&models.Course{})
		return nil, err
	}

	// Create room
	room, err := f.CreateTestRoom(ctx, center.ID, testName)
	if err != nil {
		f.db.WithContext(ctx).Where("id = ?", center.ID).Delete(&models.Center{})
		f.db.WithContext(ctx).Where("center_id = ?", center.ID).Delete(&models.Course{})
		f.db.WithContext(ctx).Where("center_id = ?", center.ID).Delete(&models.Offering{})
		return nil, err
	}

	// Create teacher (open to hiring)
	teacher, err := f.CreateTestTeacher(ctx, testName, WithTeacherOpenToHiring(true))
	if err != nil {
		f.db.WithContext(ctx).Where("id = ?", center.ID).Delete(&models.Center{})
		f.db.WithContext(ctx).Where("center_id = ?", center.ID).Delete(&models.Course{})
		f.db.WithContext(ctx).Where("center_id = ?", center.ID).Delete(&models.Offering{})
		f.db.WithContext(ctx).Where("center_id = ?", center.ID).Delete(&models.Room{})
		return nil, err
	}

	// Create schedule rule
	rule, err := f.CreateTestScheduleRule(ctx, center.ID, offering.ID, room.ID, &teacher.ID, testName)
	if err != nil {
		f.db.WithContext(ctx).Where("id = ?", center.ID).Delete(&models.Center{})
		f.db.WithContext(ctx).Where("center_id = ?", center.ID).Delete(&models.Course{})
		f.db.WithContext(ctx).Where("center_id = ?", center.ID).Delete(&models.Offering{})
		f.db.WithContext(ctx).Where("center_id = ?", center.ID).Delete(&models.Room{})
		f.db.WithContext(ctx).Where("id = ?", teacher.ID).Delete(&models.Teacher{})
		return nil, err
	}

	return &CompleteTestData{
		Center:   center,
		Course:   course,
		Offering: offering,
		Room:     room,
		Teacher:  teacher,
		Rule:     rule,
	}, nil
}

// CompleteTestData holds a complete set of related test data entities.
type CompleteTestData struct {
	Center   *models.Center
	Course   *models.Course
	Offering *models.Offering
	Room     *models.Room
	Teacher  *models.Teacher
	Rule     *models.ScheduleRule
}

// Cleanup removes all test data created in this set.
func (c *CompleteTestData) Cleanup(db *gorm.DB, ctx context.Context) {
	// Clean up in reverse order of dependencies
	if c.Rule != nil {
		db.WithContext(ctx).Where("id = ?", c.Rule.ID).Delete(&models.ScheduleRule{})
	}
	if c.Teacher != nil {
		// First clean up related data
		db.WithContext(ctx).Where("teacher_id = ?", c.Teacher.ID).Delete(&models.PersonalEvent{})
		db.WithContext(ctx).Where("teacher_id = ?", c.Teacher.ID).Delete(&models.CenterMembership{})
		db.WithContext(ctx).Where("teacher_id = ?", c.Teacher.ID).Delete(&models.CenterTeacherNote{})
		db.WithContext(ctx).Where("id = ?", c.Teacher.ID).Delete(&models.Teacher{})
	}
	if c.Center != nil {
		db.WithContext(ctx).Where("center_id = ?", c.Center.ID).Delete(&models.CenterInvitation{})
		db.WithContext(ctx).Where("center_id = ?", c.Center.ID).Delete(&models.CenterHoliday{})
		db.WithContext(ctx).Where("center_id = ?", c.Center.ID).Delete(&models.ScheduleException{})
		db.WithContext(ctx).Where("center_id = ?", c.Center.ID).Delete(&models.Offering{})
		db.WithContext(ctx).Where("center_id = ?", c.Center.ID).Delete(&models.Course{})
		db.WithContext(ctx).Where("center_id = ?", c.Center.ID).Delete(&models.Room{})
		db.WithContext(ctx).Where("center_id = ?", c.Center.ID).Delete(&models.CenterTeacherNote{})
		db.WithContext(ctx).Where("center_id = ?", c.Center.ID).Delete(&models.CenterMembership{})
		db.WithContext(ctx).Where("id = ?", c.Center.ID).Delete(&models.Center{})
	}
}

// CreateTestPersonalEvent creates a test personal event for a teacher.
func (f *TestDataFactory) CreateTestPersonalEvent(ctx context.Context, teacherID uint, namePrefix string) (*models.PersonalEvent, error) {
	event := &models.PersonalEvent{
		TeacherID: teacherID,
		Title:     fmt.Sprintf("私人行程 - %s - %s", namePrefix, f.generateUniqueID()[:8]),
		StartAt:   f.baseTime.Add(24 * time.Hour),
		EndAt:     f.baseTime.Add(26 * time.Hour),
		CreatedAt: f.baseTime,
		UpdatedAt: f.baseTime,
	}
	if err := f.db.WithContext(ctx).Create(event).Error; err != nil {
		return nil, fmt.Errorf("建立測試私人行程失敗: %w", err)
	}
	return event, nil
}

// CreateTestCenterHoliday creates a test holiday for the center.
func (f *TestDataFactory) CreateTestCenterHoliday(ctx context.Context, centerID uint, namePrefix string) (*models.CenterHoliday, error) {
	holiday := &models.CenterHoliday{
		CenterID: centerID,
		Name:     fmt.Sprintf("假日 - %s - %s", namePrefix, f.generateUniqueID()[:8]),
		Date:     f.baseTime.AddDate(0, 0, 7),
		CreatedAt: f.baseTime,
		UpdatedAt: f.baseTime,
	}
	if err := f.db.WithContext(ctx).Create(holiday).Error; err != nil {
		return nil, fmt.Errorf("建立測試假日失敗: %w", err)
	}
	return holiday, nil
}

// CreateTestScheduleException creates a test schedule exception request.
func (f *TestDataFactory) CreateTestScheduleException(ctx context.Context, centerID, ruleID uint, exceptionType string) (*models.ScheduleException, error) {
	exception := &models.ScheduleException{
		CenterID:       centerID,
		RuleID:         ruleID,
		ExceptionType:  exceptionType,
		OriginalDate:   f.baseTime.AddDate(0, 0, 3),
		Reason:         fmt.Sprintf("測試原因 - %s", f.generateUniqueID()[:8]),
		Status:         "PENDING",
		CreatedAt:      f.baseTime,
		UpdatedAt:      f.baseTime,
	}
	if err := f.db.WithContext(ctx).Create(exception).Error; err != nil {
		return nil, fmt.Errorf("建立測試例外失敗: %w", err)
	}
	return exception, nil
}

// CreateTestCenterMembership creates a membership relationship between teacher and center.
func (f *TestDataFactory) CreateTestCenterMembership(ctx context.Context, teacherID, centerID uint, status string) (*models.CenterMembership, error) {
	membership := &models.CenterMembership{
		TeacherID: teacherID,
		CenterID:  centerID,
		Status:    status,
		CreatedAt: f.baseTime,
		UpdatedAt: f.baseTime,
	}
	if err := f.db.WithContext(ctx).Create(membership).Error; err != nil {
		return nil, fmt.Errorf("建立測試中心會員關係失敗: %w", err)
	}
	return membership, nil
}
