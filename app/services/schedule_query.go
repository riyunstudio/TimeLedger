package services

import (
	"context"
	"fmt"
	"time"
	"timeLedger/app"
	"timeLedger/app/models"
	"timeLedger/app/repositories"
)

type scheduleQueryService struct {
	app              *app.App
	membershipRepo   *repositories.CenterMembershipRepository
	centerRepo       *repositories.CenterRepository
	scheduleRuleRepo *repositories.ScheduleRuleRepository
	exceptionRepo    *repositories.ScheduleExceptionRepository
	expansionService ScheduleExpansionService
}

func NewScheduleQueryService(app *app.App) ScheduleQueryService {
	return &scheduleQueryService{
		app:              app,
		membershipRepo:   repositories.NewCenterMembershipRepository(app),
		centerRepo:       repositories.NewCenterRepository(app),
		scheduleRuleRepo: repositories.NewScheduleRuleRepository(app),
		exceptionRepo:    repositories.NewScheduleExceptionRepository(app),
		expansionService: NewScheduleExpansionService(app),
	}
}

func (s *scheduleQueryService) GetTeacherSchedule(ctx context.Context, teacherID uint, fromDate, toDate time.Time) ([]TeacherScheduleItem, error) {
	memberships, err := s.membershipRepo.GetActiveByTeacherID(ctx, teacherID)
	if err != nil {
		return nil, err
	}

	var schedule []TeacherScheduleItem

	for _, m := range memberships {
		center, _ := s.centerRepo.GetByID(ctx, m.CenterID)
		centerName := center.Name

		rules, _ := s.scheduleRuleRepo.ListByTeacherID(ctx, teacherID, m.CenterID)

		// Create a map of rule ID to rule for quick lookup
		ruleMap := make(map[uint]*models.ScheduleRule)
		for i := range rules {
			ruleMap[rules[i].ID] = &rules[i]
		}

		expanded := s.expansionService.ExpandRules(ctx, rules, fromDate, toDate, m.CenterID)

		for _, item := range expanded {
			status := "NORMAL"
			exceptions, _ := s.exceptionRepo.GetByRuleAndDate(ctx, item.RuleID, item.Date)
			for _, exc := range exceptions {
				if exc.Status == "PENDING" {
					status = "PENDING_" + exc.ExceptionType
				} else if exc.Status == "APPROVED" && exc.ExceptionType == "CANCEL" {
					status = "CANCELLED"
				} else if exc.Status == "APPROVED" && exc.ExceptionType == "RESCHEDULE" {
					status = "RESCHEDULED"
				}
			}

			if status != "CANCELLED" {
				// Get offering name from the rule
				offeringName := ""
				if rule, exists := ruleMap[item.RuleID]; exists && rule.OfferingID != 0 {
					offeringName = rule.Offering.Name
				}

				// Create title: "課程名稱 @ 中心名稱"
				title := offeringName
				if centerName != "" {
					if title != "" {
						title = fmt.Sprintf("%s @ %s", offeringName, centerName)
					} else {
						title = centerName
					}
				}
				if title == "" {
					title = "課程"
				}

				schedule = append(schedule, TeacherScheduleItem{
					ID: fmt.Sprintf("center_%d_rule_%d_%s_%s", m.CenterID, item.RuleID, item.Date.Format("20060102"), func() string {
						if item.IsCrossDayPart {
							if item.StartTime == "00:00" {
								return "end"
							}
							return "start"
						}
						return "normal"
					}()),
					Type:           "CENTER_SESSION",
					Title:          title,
					Date:           item.Date.Format("2006-01-02"),
					StartTime:      item.StartTime,
					EndTime:        item.EndTime,
					RoomID:         item.RoomID,
					TeacherID:      item.TeacherID,
					CenterID:       m.CenterID,
					CenterName:     centerName,
					Status:         status,
					RuleID:         item.RuleID,
					IsCrossDayPart: item.IsCrossDayPart,
				})
			}
		}
	}

	return schedule, nil
}
