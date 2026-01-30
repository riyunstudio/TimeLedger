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
	BaseService
	app              *app.App
	membershipRepo   *repositories.CenterMembershipRepository
	centerRepo       *repositories.CenterRepository
	scheduleRuleRepo *repositories.ScheduleRuleRepository
	exceptionRepo    *repositories.ScheduleExceptionRepository
	expansionService ScheduleExpansionService
	scheduleService  ScheduleServiceInterface
}

func NewScheduleQueryService(app *app.App) ScheduleQueryService {
	baseSvc := NewBaseService(app, "ScheduleQueryService")
	return &scheduleQueryService{
		BaseService:      *baseSvc,
		app:              app,
		membershipRepo:   repositories.NewCenterMembershipRepository(app),
		centerRepo:       repositories.NewCenterRepository(app),
		scheduleRuleRepo: repositories.NewScheduleRuleRepository(app),
		exceptionRepo:    repositories.NewScheduleExceptionRepository(app),
		expansionService: NewScheduleExpansionService(app),
		scheduleService:  NewScheduleService(app),
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

		// 使用帶快取的老師課表查詢（消除了 N+1 查詢問題）
		expanded, err := s.scheduleService.GetCachedTeacherSchedule(ctx, teacherID, m.CenterID, fromDate, toDate)
		if err != nil {
			s.Logger.Warn("failed to get cached teacher schedule, falling back", "error", err, "center_id", m.CenterID)
			// 如果快取失敗，回退到直接展開
			rules, _ := s.scheduleRuleRepo.ListByTeacherID(ctx, teacherID, m.CenterID)
			expanded = s.expansionService.ExpandRules(ctx, rules, fromDate, toDate, m.CenterID)
		}

		// 建立規則 Map 用於查找課程名稱（使用值類型，因為 GetByID 返回值而非指針）
		ruleMap := make(map[uint]models.ScheduleRule)
		for i := range expanded {
			if expanded[i].OfferingID != 0 {
				// 需要獲取規則來取得課程名稱
				rule, _ := s.scheduleRuleRepo.GetByID(ctx, expanded[i].RuleID)
				if rule.ID != 0 {
					ruleMap[rule.ID] = rule
				}
			}
		}

		for _, item := range expanded {
			// 使用 ExpandedSchedule 中已包含的例外資訊（來自 ExpandRules 批次查詢）
			// 不再進行額外的資料庫查詢
			status := "NORMAL"

			// 從 ExceptionInfo 判斷狀態（來自 ExpandRules 預載入的資料）
			if item.ExceptionInfo != nil {
				if item.ExceptionInfo.Status == "PENDING" {
					status = "PENDING_" + item.ExceptionInfo.Type
				} else if item.ExceptionInfo.Status == "APPROVED" {
					if item.ExceptionInfo.Type == "CANCEL" {
						status = "CANCELLED"
					} else if item.ExceptionInfo.Type == "RESCHEDULE" {
						status = "RESCHEDULED"
					}
				}
			}

			if status != "CANCELLED" {
				// 從 ruleMap 獲取課程名稱
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
