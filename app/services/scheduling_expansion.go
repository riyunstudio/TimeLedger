package services

import (
	"context"
	"strconv"
	"strings"
	"time"
	"timeLedger/app"
	"timeLedger/app/models"
	"timeLedger/app/repositories"
)

type ScheduleExpansionServiceImpl struct {
	BaseService
	scheduleRuleRepo *repositories.ScheduleRuleRepository
	exceptionRepo    *repositories.ScheduleExceptionRepository
	auditLogRepo     *repositories.AuditLogRepository
	centerRepo       *repositories.CenterRepository
	holidayRepo      *repositories.CenterHolidayRepository
}

func NewScheduleExpansionService(app *app.App) ScheduleExpansionService {
	baseSvc := NewBaseService(app, "ScheduleExpansionService")
	svc := &ScheduleExpansionServiceImpl{
		BaseService: *baseSvc,
	}

	if app.MySQL != nil {
		svc.scheduleRuleRepo = repositories.NewScheduleRuleRepository(app)
		svc.exceptionRepo = repositories.NewScheduleExceptionRepository(app)
		svc.auditLogRepo = repositories.NewAuditLogRepository(app)
		svc.centerRepo = repositories.NewCenterRepository(app)
		svc.holidayRepo = repositories.NewCenterHolidayRepository(app)
	}

	return svc
}

func (s *ScheduleExpansionServiceImpl) ExpandRules(ctx context.Context, rules []models.ScheduleRule, startDate, endDate time.Time, centerID uint) []ExpandedSchedule {
	var schedules []ExpandedSchedule

	holidays, _ := s.holidayRepo.ListByDateRange(ctx, centerID, startDate, endDate)
	holidayMap := make(map[string]models.CenterHoliday)
	for _, h := range holidays {
		holidayMap[h.Date.Format("2006-01-02")] = h
	}

	ruleIDs := make([]uint, 0, len(rules))
	for _, rule := range rules {
		if rule.Weekday != 0 {
			ruleIDs = append(ruleIDs, rule.ID)
		}
	}

	exceptionsMap := make(map[uint]map[string][]models.ScheduleException)
	if len(ruleIDs) > 0 {
		exceptionsMap, _ = s.exceptionRepo.GetByRuleIDsAndDateRange(ctx, ruleIDs, startDate, endDate)
	}

	for _, rule := range rules {
		if rule.Weekday == 0 {
			continue
		}

		ruleStartDate := rule.EffectiveRange.StartDate
		ruleEndDate := rule.EffectiveRange.EndDate

		date := startDate
		for date.Before(endDate) || date.Equal(endDate) {
			weekday := int(date.Weekday())
			if weekday == 0 {
				weekday = 7
			}

			if weekday == int(rule.Weekday) {
				isWithinEffectiveRange := true
				if !ruleStartDate.IsZero() && date.Before(ruleStartDate) {
					isWithinEffectiveRange = false
				}
				if !ruleEndDate.IsZero() && date.After(ruleEndDate) {
					isWithinEffectiveRange = false
				}

				if isWithinEffectiveRange {
					dateStr := date.Format("2006-01-02")
					isHoliday, exists := holidayMap[dateStr]

					if exists && (isHoliday.ForceCancel || rule.SkipHoliday) {
						date = date.AddDate(0, 0, 1)
						continue
					}

					ruleExceptions := exceptionsMap[rule.ID]
					exceptions := []models.ScheduleException{}
					if ruleExceptions != nil {
						exceptions = ruleExceptions[dateStr]
					}

					skipSession := false
					var pendingException *models.ScheduleException
					var approvedException *models.ScheduleException

					for i := range exceptions {
						exc := &exceptions[i]
						if exc.Status == "CANCELLED" {
							continue
						}

						if exc.Status == "PENDING" {
							if pendingException == nil {
								pendingException = exc
							}
						}

						if exc.Status == "APPROVED" {
							if exc.ExceptionType == "CANCEL" {
								skipSession = true
								approvedException = exc
								break
							}
							if exc.ExceptionType == "REPLACE_TEACHER" && exc.NewTeacherID != nil {
								rule.TeacherID = exc.NewTeacherID
							}
							if exc.ExceptionType == "RESCHEDULE" {
								if approvedException == nil {
									approvedException = exc
								}
							}
						}
					}

					if skipSession {
						date = date.AddDate(0, 0, 1)
						continue
					}

					startParts := strings.Split(rule.StartTime, ":")
					endParts := strings.Split(rule.EndTime, ":")
					startHour, _ := strconv.Atoi(startParts[0])
					endHour, _ := strconv.Atoi(endParts[0])

					isCrossDay := endHour < startHour

					createScheduleEntry := func(scheduleDate time.Time, sTime, eTime string) {
						schedule := ExpandedSchedule{
							RuleID:         rule.ID,
							Date:           scheduleDate,
							StartTime:      sTime,
							EndTime:        eTime,
							RoomID:         rule.RoomID,
							TeacherID:      rule.TeacherID,
							IsHoliday:      exists,
							HasException:   pendingException != nil || approvedException != nil,
							OfferingName:   rule.Offering.Name,
							TeacherName:    rule.Teacher.Name,
							RoomName:       rule.Room.Name,
							OfferingID:     rule.OfferingID,
							EffectiveRange: &rule.EffectiveRange,
							IsCrossDayPart: isCrossDay,
						}

						if pendingException != nil {
							schedule.ExceptionInfo = &ExpandedException{
								ID:           pendingException.ID,
								Type:         pendingException.ExceptionType,
								Status:       pendingException.Status,
								NewTeacherID: pendingException.NewTeacherID,
								NewStartAt:   pendingException.NewStartAt,
								NewEndAt:     pendingException.NewEndAt,
							}
						}

						schedules = append(schedules, schedule)
					}

					if isCrossDay {
						createScheduleEntry(date, rule.StartTime, "24:00")
						nextDay := date.AddDate(0, 0, 1)
						createScheduleEntry(nextDay, "00:00", rule.EndTime)
					} else {
						createScheduleEntry(date, rule.StartTime, rule.EndTime)
					}
				}
			}

			date = date.AddDate(0, 0, 1)
		}
	}

	return schedules
}

func (s *ScheduleExpansionServiceImpl) GetEffectiveRuleForDate(ctx context.Context, offeringID uint, date time.Time) (*models.ScheduleRule, error) {
	var rules []models.ScheduleRule
	// 處理零值時間
	startDateStr := "0001-01-01"
	if !date.IsZero() {
		startDateStr = date.Format("2006-01-02")
	}

	err := s.App.MySQL.RDB.WithContext(ctx).
		Where("offering_id = ?", offeringID).
		Where("weekday = ?", func() int {
			weekday := int(date.Weekday())
			if weekday == 0 {
				return 7
			}
			return weekday
		}()).
		Where("COALESCE(NULLIF(NULLIF(JSON_UNQUOTE(JSON_EXTRACT(effective_range, '$.start_date')), ''), 'null'), '0001-01-01') <= ?", startDateStr).
		Where("COALESCE(NULLIF(NULLIF(NULLIF(JSON_UNQUOTE(JSON_EXTRACT(effective_range, '$.end_date')), ''), 'null'), '0001-01-01 00:00:00'), '9999-12-31') >= ?", startDateStr).
		Find(&rules).Error

	if err != nil {
		return nil, err
	}

	if len(rules) == 0 {
		return nil, nil
	}

	return &rules[0], nil
}

func (s *ScheduleExpansionServiceImpl) DetectPhaseTransitions(ctx context.Context, centerID uint, offeringID uint, startDate, endDate time.Time) ([]PhaseTransition, error) {
	rules, err := s.scheduleRuleRepo.ListByOfferingID(ctx, offeringID)
	if err != nil {
		return nil, err
	}

	if len(rules) <= 1 {
		return []PhaseTransition{}, nil
	}

	ruleByDate := make(map[string]*models.ScheduleRule)
	for i := range rules {
		rule := &rules[i]
		ruleStart := rule.EffectiveRange.StartDate
		ruleEnd := rule.EffectiveRange.EndDate

		current := ruleStart
		for !current.After(ruleEnd) {
			if !current.Before(startDate) && !current.After(endDate) {
				weekday := int(current.Weekday())
				if weekday == 0 {
					weekday = 7
				}

				if weekday == int(rule.Weekday) {
					dateStr := current.Format("2006-01-02")
					ruleByDate[dateStr] = rule
				}
			}
			current = current.AddDate(0, 0, 1)
		}
	}

	var transitions []PhaseTransition
	date := startDate
	prevRule := (*models.ScheduleRule)(nil)

	for date.Before(endDate) || date.Equal(endDate) {
		dateStr := date.Format("2006-01-02")
		currentRule := ruleByDate[dateStr]

		if prevRule != nil && currentRule != nil {
			prevRuleID := prevRule.ID
			currRuleID := currentRule.ID

			if prevRuleID != currRuleID ||
				(prevRule.RoomID != currentRule.RoomID) ||
				!ptrEqual(prevRule.TeacherID, currentRule.TeacherID) ||
				(prevRule.StartTime != currentRule.StartTime) ||
				(prevRule.EndTime != currentRule.EndTime) {

				transition := PhaseTransition{
					Date:          date,
					PrevRuleID:    &prevRuleID,
					PrevRoomID:    &prevRule.RoomID,
					PrevTeacherID: prevRule.TeacherID,
					PrevStartTime: prevRule.StartTime,
					PrevEndTime:   prevRule.EndTime,
					NextRuleID:    &currRuleID,
					NextRoomID:    &currentRule.RoomID,
					NextTeacherID: currentRule.TeacherID,
					NextStartTime: currentRule.StartTime,
					NextEndTime:   currentRule.EndTime,
					HasGap:        false,
				}
				transitions = append(transitions, transition)
			}
		} else if prevRule != nil && currentRule == nil {
			prevRuleID := prevRule.ID
			transition := PhaseTransition{
				Date:       date,
				PrevRuleID: &prevRuleID,
				PrevRoomID: &prevRule.RoomID,
				HasGap:     true,
			}
			transitions = append(transitions, transition)
		} else if prevRule == nil && currentRule != nil {
			currRuleID := currentRule.ID
			transition := PhaseTransition{
				Date:       date,
				NextRuleID: &currRuleID,
				NextRoomID: &currentRule.RoomID,
			}
			transitions = append(transitions, transition)
		}

		prevRule = currentRule
		date = date.AddDate(0, 0, 1)
	}

	return transitions, nil
}

func (s *ScheduleExpansionServiceImpl) GetRulesByEffectiveDateRange(ctx context.Context, centerID uint, offeringID uint, startDate, endDate time.Time) ([]models.ScheduleRule, error) {
	var rules []models.ScheduleRule
	// 處理零值時間
	startDateStr := "0001-01-01"
	if !startDate.IsZero() {
		startDateStr = startDate.Format("2006-01-02")
	}
	endDateStr := "9999-12-31"
	if !endDate.IsZero() {
		endDateStr = endDate.Format("2006-01-02")
	}

	err := s.App.MySQL.RDB.WithContext(ctx).
		Where("center_id = ?", centerID).
		Where("offering_id = ?", offeringID).
		Where("COALESCE(NULLIF(NULLIF(JSON_UNQUOTE(JSON_EXTRACT(effective_range, '$.start_date')), ''), 'null'), '0001-01-01') <= ?", endDateStr).
		Where("COALESCE(NULLIF(NULLIF(NULLIF(JSON_UNQUOTE(JSON_EXTRACT(effective_range, '$.end_date')), ''), 'null'), '0001-01-01 00:00:00'), '9999-12-31') >= ?", startDateStr).
		Order("JSON_UNQUOTE(JSON_EXTRACT(effective_range, '$.start_date')) ASC").
		Find(&rules).Error

	return rules, err
}

func ptrEqual(a, b *uint) bool {
	if a == nil && b == nil {
		return true
	}
	if a == nil || b == nil {
		return false
	}
	return *a == *b
}
