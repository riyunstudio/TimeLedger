package services

import (
	"bytes"
	"context"
	"encoding/csv"
	"fmt"
	"time"
	"timeLedger/app"
	"timeLedger/app/repositories"
)

type ExportService interface {
	ExportScheduleToCSV(ctx context.Context, centerID uint, startDate, endDate time.Time) ([]byte, error)
	ExportScheduleToPDF(ctx context.Context, centerID uint, startDate, endDate time.Time) ([]byte, error)
	ExportTeachersToCSV(ctx context.Context, centerID uint) ([]byte, error)
	ExportExceptionsToCSV(ctx context.Context, centerID uint, startDate, endDate time.Time) ([]byte, error)
	GenerateScheduleCSV(ctx context.Context, centerID uint, startDate, endDate time.Time) ([][]string, error)
	GenerateTeacherCSV(ctx context.Context, centerID uint) ([][]string, error)
	GenerateExceptionCSV(ctx context.Context, centerID uint) ([][]string, error)
	ValidatePDFParams(centerID uint, startDate, endDate time.Time) error
}

type ExportServiceImpl struct {
	BaseService
	app                   *app.App
	scheduleRuleRepo      *repositories.ScheduleRuleRepository
	scheduleExceptionRepo *repositories.ScheduleExceptionRepository
	teacherRepo           *repositories.TeacherRepository
}

func NewExportService(app *app.App) ExportService {
	return &ExportServiceImpl{
		app:                   app,
		scheduleRuleRepo:      repositories.NewScheduleRuleRepository(app),
		scheduleExceptionRepo: repositories.NewScheduleExceptionRepository(app),
		teacherRepo:           repositories.NewTeacherRepository(app),
	}
}

func (s *ExportServiceImpl) ExportScheduleToCSV(ctx context.Context, centerID uint, startDate, endDate time.Time) ([]byte, error) {
	rules, err := s.scheduleRuleRepo.ListByCenterID(ctx, centerID)
	if err != nil {
		return nil, err
	}

	var records [][]string
	records = append(records, []string{"Date", "Weekday", "Start Time", "End Time", "Teacher", "Room", "Course"})

	for _, rule := range rules {
		teacherName := "未指定"
		if rule.TeacherID != nil {
			teacher, err := s.teacherRepo.GetByID(ctx, *rule.TeacherID)
			if err == nil {
				teacherName = teacher.Name
			}
		}

		records = append(records, []string{
			startDate.Format("2006-01-02"),
			fmt.Sprintf("%d", rule.Weekday),
			rule.StartTime,
			rule.EndTime,
			teacherName,
			"",
			"",
		})
	}

	var output bytes.Buffer
	writer := csv.NewWriter(&output)
	for _, record := range records {
		if err := writer.Write(record); err != nil {
			return nil, err
		}
	}
	writer.Flush()

	return output.Bytes(), nil
}

func (s *ExportServiceImpl) ExportScheduleToPDF(ctx context.Context, centerID uint, startDate, endDate time.Time) ([]byte, error) {
	rules, err := s.scheduleRuleRepo.ListByCenterID(ctx, centerID)
	if err != nil {
		return nil, err
	}

	content := "排課表\n\n"
	content += fmt.Sprintf("開始日期: %s\n", startDate.Format("2006-01-02"))
	content += fmt.Sprintf("結束日期: %s\n\n", endDate.Format("2006-01-02"))
	content += fmt.Sprintf("%-12s %-10s %-10s %-10s %-20s %-15s\n", "日期", "星期", "開始", "結束", "老師", "教室")
	content += "--------------------------------------------------------------------------------\n"

	for _, rule := range rules {
		teacherName := "未指定"
		if rule.TeacherID != nil {
			teacher, err := s.teacherRepo.GetByID(ctx, *rule.TeacherID)
			if err == nil {
				teacherName = teacher.Name
			}
		}

		weekday := s.getWeekdayName(rule.Weekday)
		content += fmt.Sprintf("%-12s %-10s %-10s %-10s %-20s %-15s\n",
			startDate.Format("2006-01-02"),
			weekday,
			rule.StartTime,
			rule.EndTime,
			teacherName,
			"",
		)
	}

	return []byte(content), nil
}

func (s *ExportServiceImpl) ExportTeachersToCSV(ctx context.Context, centerID uint) ([]byte, error) {
	teachers, err := s.teacherRepo.List(ctx)
	if err != nil {
		return nil, err
	}

	var records [][]string
	records = append(records, []string{"ID", "姓名", "Email", "城市", "區域", "Bio", "開放求職"})

	for _, teacher := range teachers {
		records = append(records, []string{
			fmt.Sprintf("%d", teacher.ID),
			teacher.Name,
			teacher.Email,
			teacher.City,
			teacher.District,
			teacher.Bio,
			fmt.Sprintf("%v", teacher.IsOpenToHiring),
		})
	}

	var output bytes.Buffer
	writer := csv.NewWriter(&output)
	for _, record := range records {
		if err := writer.Write(record); err != nil {
			return nil, err
		}
	}
	writer.Flush()

	return output.Bytes(), nil
}

func (s *ExportServiceImpl) ExportExceptionsToCSV(ctx context.Context, centerID uint, startDate, endDate time.Time) ([]byte, error) {
	exceptions, err := s.scheduleExceptionRepo.ListByCenterID(ctx, centerID)
	if err != nil {
		return nil, err
	}

	var records [][]string
	records = append(records, []string{"ID", "原始日期", "類型", "狀態", "新開始時間", "新結束時間", "原因"})

	for _, exc := range exceptions {
		var newStartAt, newEndAt string
		if exc.NewStartAt != nil {
			newStartAt = exc.NewStartAt.Format("2006-01-02 15:04:05")
		}
		if exc.NewEndAt != nil {
			newEndAt = exc.NewEndAt.Format("2006-01-02 15:04:05")
		}

		records = append(records, []string{
			fmt.Sprintf("%d", exc.ID),
			exc.OriginalDate.Format("2006-01-02"),
			exc.ExceptionType,
			exc.Status,
			newStartAt,
			newEndAt,
			exc.Reason,
		})
	}

	var output bytes.Buffer
	writer := csv.NewWriter(&output)
	for _, record := range records {
		if err := writer.Write(record); err != nil {
			return nil, err
		}
	}
	writer.Flush()

	return output.Bytes(), nil
}

func (s *ExportServiceImpl) getWeekdayName(weekday int) string {
	names := []string{"", "星期一", "星期二", "星期三", "星期四", "星期五", "星期六", "星期日"}
	if weekday >= 1 && weekday <= 7 {
		return names[weekday]
	}
	return ""
}

// GenerateScheduleCSV generates CSV data as 2D slice for export
func (s *ExportServiceImpl) GenerateScheduleCSV(ctx context.Context, centerID uint, startDate, endDate time.Time) ([][]string, error) {
	rules, err := s.scheduleRuleRepo.ListByCenterID(ctx, centerID)
	if err != nil {
		return nil, err
	}

	var records [][]string
	records = append(records, []string{"Date", "Weekday", "Start Time", "End Time", "Teacher", "Room", "Course"})

	for _, rule := range rules {
		teacherName := "未指定"
		if rule.TeacherID != nil {
			teacher, err := s.teacherRepo.GetByID(ctx, *rule.TeacherID)
			if err == nil {
				teacherName = teacher.Name
			}
		}

		records = append(records, []string{
			startDate.Format("2006-01-02"),
			fmt.Sprintf("%d", rule.Weekday),
			rule.StartTime,
			rule.EndTime,
			teacherName,
			"",
			"",
		})
	}

	return records, nil
}

// GenerateTeacherCSV generates teacher CSV data as 2D slice for export
func (s *ExportServiceImpl) GenerateTeacherCSV(ctx context.Context, centerID uint) ([][]string, error) {
	teachers, err := s.teacherRepo.List(ctx)
	if err != nil {
		return nil, err
	}

	var records [][]string
	records = append(records, []string{"ID", "姓名", "Email", "城市", "區域", "Bio", "開放求職"})

	for _, teacher := range teachers {
		records = append(records, []string{
			fmt.Sprintf("%d", teacher.ID),
			teacher.Name,
			teacher.Email,
			teacher.City,
			teacher.District,
			teacher.Bio,
			fmt.Sprintf("%v", teacher.IsOpenToHiring),
		})
	}

	return records, nil
}

// GenerateExceptionCSV generates exception CSV data as 2D slice for export
func (s *ExportServiceImpl) GenerateExceptionCSV(ctx context.Context, centerID uint) ([][]string, error) {
	exceptions, err := s.scheduleExceptionRepo.ListByCenterID(ctx, centerID)
	if err != nil {
		return nil, err
	}

	var records [][]string
	records = append(records, []string{"ID", "原始日期", "類型", "狀態", "新開始時間", "新結束時間", "原因"})

	for _, exc := range exceptions {
		var newStartAt, newEndAt string
		if exc.NewStartAt != nil {
			newStartAt = exc.NewStartAt.Format("2006-01-02 15:04:05")
		}
		if exc.NewEndAt != nil {
			newEndAt = exc.NewEndAt.Format("2006-01-02 15:04:05")
		}

		records = append(records, []string{
			fmt.Sprintf("%d", exc.ID),
			exc.OriginalDate.Format("2006-01-02"),
			exc.ExceptionType,
			exc.Status,
			newStartAt,
			newEndAt,
			exc.Reason,
		})
	}

	return records, nil
}

// ValidatePDFParams validates PDF export parameters
func (s *ExportServiceImpl) ValidatePDFParams(centerID uint, startDate, endDate time.Time) error {
	if startDate.After(endDate) {
		return fmt.Errorf("start date must be before or equal to end date")
	}
	return nil
}
