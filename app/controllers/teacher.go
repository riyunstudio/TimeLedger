package controllers

import (
	"crypto/rand"
	"fmt"
	"net/http"
	"strings"
	"time"
	"timeLedger/app"
	"timeLedger/app/models"
	"timeLedger/app/repositories"
	"timeLedger/app/requests"
	"timeLedger/app/resources"
	"timeLedger/app/services"
	"timeLedger/global"
	"timeLedger/global/errInfos"

	"github.com/gin-gonic/gin"
)

type TeacherController struct {
	BaseController
	app               *app.App
	teacherRepository *repositories.TeacherRepository
	membershipRepo    *repositories.CenterMembershipRepository
	centerRepo        *repositories.CenterRepository
	scheduleRuleRepo  *repositories.ScheduleRuleRepository
	exceptionRepo     *repositories.ScheduleExceptionRepository
	exceptionService  services.ScheduleExceptionService
	expansionService  services.ScheduleExpansionService
	recurrenceService services.ScheduleRecurrenceService
	auditLogRepo      *repositories.AuditLogRepository
	skillRepo         *repositories.TeacherSkillRepository
	certificateRepo   *repositories.TeacherCertificateRepository
	personalEventRepo *repositories.PersonalEventRepository
	sessionNoteRepo   *repositories.SessionNoteRepository
}

func NewTeacherController(app *app.App) *TeacherController {
	return &TeacherController{
		app:               app,
		teacherRepository: repositories.NewTeacherRepository(app),
		membershipRepo:    repositories.NewCenterMembershipRepository(app),
		centerRepo:        repositories.NewCenterRepository(app),
		scheduleRuleRepo:  repositories.NewScheduleRuleRepository(app),
		exceptionRepo:     repositories.NewScheduleExceptionRepository(app),
		exceptionService:  services.NewScheduleExceptionService(app),
		expansionService:  services.NewScheduleExpansionService(app),
		recurrenceService: services.NewScheduleRecurrenceService(app),
		auditLogRepo:      repositories.NewAuditLogRepository(app),
		skillRepo:         repositories.NewTeacherSkillRepository(app),
		certificateRepo:   repositories.NewTeacherCertificateRepository(app),
		personalEventRepo: repositories.NewPersonalEventRepository(app),
		sessionNoteRepo:   repositories.NewSessionNoteRepository(app),
	}
}

// GetProfile 取得老師個人資料
// @Summary 取得老師個人資料
// @Tags Teacher
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} global.ApiResponse{data=resources.TeacherProfileResource}
// @Router /api/v1/teacher/me/profile [get]
func (ctl *TeacherController) GetProfile(ctx *gin.Context) {
	teacherID := ctx.GetUint(global.UserIDKey)
	if teacherID == 0 {
		ctx.JSON(http.StatusUnauthorized, global.ApiResponse{
			Code:    global.UNAUTHORIZED,
			Message: "Teacher ID required",
		})
		return
	}

	teacher, err := ctl.teacherRepository.GetByID(ctx, teacherID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, global.ApiResponse{
			Code:    500,
			Message: "Failed to get teacher profile",
		})
		return
	}

	// 取得個人標籤
	personalHashtags, _ := ctl.teacherRepository.ListPersonalHashtags(ctx, teacherID)
	var hashtagResources []resources.PersonalHashtag
	for _, h := range personalHashtags {
		hashtagResources = append(hashtagResources, resources.PersonalHashtag{
			ID:        h.ID,
			HashtagID: h.HashtagID,
			Name:      h.Name,
		})
	}

	response := resources.TeacherProfileResource{
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
	}

	ctx.JSON(http.StatusOK, global.ApiResponse{
		Code:    0,
		Message: "Success",
		Datas:   response,
	})
}

// SearchHashtags 搜尋標籤
// @Summary 搜尋標籤
// @Tags Hashtag
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param q query string true "搜尋關鍵字"
// @Success 200 {object} global.ApiResponse{data=[]resources.HashtagResource}
// @Router /api/v1/hashtags/search [get]
func (ctl *TeacherController) SearchHashtags(ctx *gin.Context) {
	query := ctx.Query("q")
	if query == "" {
		ctx.JSON(http.StatusBadRequest, global.ApiResponse{
			Code:    global.BAD_REQUEST,
			Message: "Query parameter 'q' is required",
		})
		return
	}

	// 移除 # 符號
	query = strings.TrimPrefix(query, "#")

	hashtagRepo := repositories.NewHashtagRepository(ctl.app)
	hashtags, err := hashtagRepo.Search(ctx, query)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, global.ApiResponse{
			Code:    500,
			Message: "Failed to search hashtags",
		})
		return
	}

	// 如果沒有找到，返回空陣列而不是錯誤
	if hashtags == nil {
		hashtags = []models.Hashtag{}
	}

	// 轉換為 Resource
	var hashtagResources []resources.HashtagResource
	for _, h := range hashtags {
		hashtagResources = append(hashtagResources, resources.HashtagResource{
			ID:         h.ID,
			Name:       h.Name,
			UsageCount: h.UsageCount,
		})
	}

	ctx.JSON(http.StatusOK, global.ApiResponse{
		Code:    0,
		Message: "Success",
		Datas:   hashtagResources,
	})
}

// CreateHashtag 建立新標籤
// @Summary 建立新標籤
// @Tags Hashtag
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body CreateHashtagRequest true "標籤資訊"
// @Success 200 {object} global.ApiResponse{data=models.Hashtag}
// @Router /api/v1/hashtags [post]
func (ctl *TeacherController) CreateHashtag(ctx *gin.Context) {
	var req CreateHashtagRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, global.ApiResponse{
			Code:    global.BAD_REQUEST,
			Message: "Invalid request body",
		})
		return
	}

	// 確保標籤名稱以 # 開頭
	name := req.Name
	if !strings.HasPrefix(name, "#") {
		name = "#" + name
	}

	// 檢查標籤是否已存在
	hashtagRepo := repositories.NewHashtagRepository(ctl.app)
	existing, err := hashtagRepo.GetByName(ctx, name)
	if err == nil && existing != nil {
		// 標籤已存在，直接返回
		ctx.JSON(http.StatusOK, global.ApiResponse{
			Code:    0,
			Message: "Hashtag already exists",
			Datas:   existing,
		})
		return
	}

	// 建立新標籤
	hashtag := &models.Hashtag{
		Name:       name,
		UsageCount: 1,
	}
	if err := hashtagRepo.Create(ctx, hashtag); err != nil {
		ctx.JSON(http.StatusInternalServerError, global.ApiResponse{
			Code:    500,
			Message: "Failed to create hashtag",
		})
		return
	}

	ctx.JSON(http.StatusOK, global.ApiResponse{
		Code:    0,
		Message: "Hashtag created",
		Datas:   hashtag,
	})
}

type CreateHashtagRequest struct {
	Name string `json:"name" binding:"required"`
}

// UpdateProfile 更新老師個人資料
// @Summary 更新老師個人資料
// @Tags Teacher
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body UpdateTeacherProfileRequest true "個人資料"
// @Success 200 {object} global.ApiResponse{data=resources.TeacherProfileResource}
// @Router /api/v1/teacher/me/profile [put]
func (ctl *TeacherController) UpdateProfile(ctx *gin.Context) {
	teacherID := ctx.GetUint(global.UserIDKey)
	if teacherID == 0 {
		ctx.JSON(http.StatusUnauthorized, global.ApiResponse{
			Code:    global.UNAUTHORIZED,
			Message: "Teacher ID required",
		})
		return
	}

	var req UpdateTeacherProfileRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, global.ApiResponse{
			Code:    global.BAD_REQUEST,
			Message: "Invalid request body",
		})
		return
	}

	teacher, err := ctl.teacherRepository.GetByID(ctx, teacherID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, global.ApiResponse{
			Code:    500,
			Message: "Failed to get teacher profile",
		})
		return
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

	// 更新個人標籤
	if len(req.PersonalHashtags) > 0 {
		// 刪除現有標籤
		ctl.app.MySQL.WDB.WithContext(ctx).Where("teacher_id = ?", teacherID).Delete(&models.TeacherPersonalHashtag{})

		// 新增新標籤
		hashtagRepo := repositories.NewHashtagRepository(ctl.app)
		for _, tagName := range req.PersonalHashtags {
			// 確保 # 符號存在
			if !strings.HasPrefix(tagName, "#") {
				tagName = "#" + tagName
			}

			// 查找或創建標籤
			hashtag, err := hashtagRepo.GetByName(ctx, tagName)
			if err != nil {
				// 創建新標籤
				hashtag = &models.Hashtag{Name: tagName, UsageCount: 1}
				hashtagRepo.Create(ctx, hashtag)
			} else {
				// 更新使用次數
				hashtagRepo.IncrementUsage(ctx, tagName)
			}

			// 創建關聯
			ctl.app.MySQL.WDB.WithContext(ctx).Create(&models.TeacherPersonalHashtag{
				TeacherID:  teacherID,
				HashtagID:  hashtag.ID,
				SortOrder:  0,
			})
		}
	}

	if err := ctl.teacherRepository.Update(ctx, teacher); err != nil {
		ctx.JSON(http.StatusInternalServerError, global.ApiResponse{
			Code:    500,
			Message: "Failed to update teacher profile",
		})
		return
	}

	memberships, _ := ctl.membershipRepo.GetActiveByTeacherID(ctx, teacherID)
	var centerID uint
	if len(memberships) > 0 {
		centerID = memberships[0].CenterID
	}

	ctl.auditLogRepo.Create(ctx, models.AuditLog{
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
	})

	response := resources.TeacherProfileResource{
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

	ctx.JSON(http.StatusOK, global.ApiResponse{
		Code:    0,
		Message: "Profile updated",
		Datas:   response,
	})
}

type UpdateTeacherProfileRequest struct {
	Bio                 string `json:"bio"`
	City                string `json:"city"`
	District            string `json:"district"`
	PublicContactInfo   string `json:"public_contact_info"`
	IsOpenToHiring      bool   `json:"is_open_to_hiring"`
	PersonalHashtags    []string `json:"personal_hashtags"`
}

// GetCenters 取得老師已加入的中心列表
// @Summary 取得老師已加入的中心列表
// @Tags Teacher
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} global.ApiResponse{data=[]resources.CenterMembershipResource}
// @Router /api/v1/teacher/me/centers [get]
func (ctl *TeacherController) GetCenters(ctx *gin.Context) {
	teacherID := ctx.GetUint(global.UserIDKey)
	if teacherID == 0 {
		ctx.JSON(http.StatusUnauthorized, global.ApiResponse{
			Code:    global.UNAUTHORIZED,
			Message: "Teacher ID required",
		})
		return
	}

	memberships, err := ctl.membershipRepo.GetActiveByTeacherID(ctx, teacherID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, global.ApiResponse{
			Code:    500,
			Message: "Failed to get memberships",
		})
		return
	}

	var centerResources []resources.CenterMembershipResource
	for _, m := range memberships {
		var center models.Center
		ctl.app.MySQL.RDB.WithContext(ctx).First(&center, m.CenterID)
		centerResources = append(centerResources, resources.CenterMembershipResource{
			ID:         m.ID,
			CenterID:   m.CenterID,
			CenterName: center.Name,
			Status:     string(m.Status),
			CreatedAt:  m.CreatedAt,
		})
	}

	ctx.JSON(http.StatusOK, global.ApiResponse{
		Code:    0,
		Message: "Success",
		Datas:   centerResources,
	})
}

// GetSchedule 取得老師的綜合課表（個人行程 + 各中心課程）
// @Summary 取得老師的綜合課表
// @Tags Teacher
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param from query string true "開始日期 (YYYY-MM-DD)"
// @Param to query string true "結束日期 (YYYY-MM-DD)"
// @Success 200 {object} global.ApiResponse{data=[]TeacherScheduleItem}
// @Router /api/v1/teacher/me/schedule [get]
func (ctl *TeacherController) GetSchedule(ctx *gin.Context) {
	teacherID := ctx.GetUint(global.UserIDKey)
	if teacherID == 0 {
		ctx.JSON(http.StatusUnauthorized, global.ApiResponse{
			Code:    global.UNAUTHORIZED,
			Message: "Teacher ID required",
		})
		return
	}

	fromStr := ctx.Query("from")
	toStr := ctx.Query("to")
	if fromStr == "" || toStr == "" {
		ctx.JSON(http.StatusBadRequest, global.ApiResponse{
			Code:    global.BAD_REQUEST,
			Message: "from and to dates are required",
		})
		return
	}

	fromDate, err := time.Parse("2006-01-02", fromStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, global.ApiResponse{
			Code:    global.BAD_REQUEST,
			Message: "Invalid from date format",
		})
		return
	}

	toDate, err := time.Parse("2006-01-02", toStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, global.ApiResponse{
			Code:    global.BAD_REQUEST,
			Message: "Invalid to date format",
		})
		return
	}

	memberships, err := ctl.membershipRepo.GetActiveByTeacherID(ctx, teacherID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, global.ApiResponse{
			Code:    500,
			Message: "Failed to get memberships",
		})
		return
	}

	var schedule []TeacherScheduleItem

	for _, m := range memberships {
		// Get center name
		center, _ := ctl.centerRepo.GetByID(ctx, m.CenterID)
		centerName := center.Name

		rules, _ := ctl.scheduleRuleRepo.ListByCenterID(ctx, m.CenterID)

		// Create a map of rule ID to rule for quick lookup
		ruleMap := make(map[uint]*models.ScheduleRule)
		for i := range rules {
			ruleMap[rules[i].ID] = &rules[i]
		}

		expanded := ctl.expansionService.ExpandRules(ctx, rules, fromDate, toDate, m.CenterID)

		for _, item := range expanded {
			status := "NORMAL"
			exceptions, _ := ctl.exceptionRepo.GetByRuleAndDate(ctx, item.RuleID, item.Date)
			for _, exc := range exceptions {
				if exc.Status == "PENDING" {
					status = "PENDING_" + exc.Type
				} else if exc.Status == "APPROVED" && exc.Type == "CANCEL" {
					status = "CANCELLED"
				} else if exc.Status == "APPROVED" && exc.Type == "RESCHEDULE" {
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
					ID:         fmt.Sprintf("center_%d_rule_%d_%s", m.CenterID, item.RuleID, item.Date.Format("20060102")),
					Type:       "CENTER_SESSION",
					Title:      title,
					Date:       item.Date.Format("2006-01-02"),
					StartTime:  item.StartTime,
					EndTime:    item.EndTime,
					RoomID:     item.RoomID,
					TeacherID:  item.TeacherID,
					CenterID:   m.CenterID,
					CenterName: centerName,
					Status:     status,
				})
			}
		}
	}

	ctx.JSON(http.StatusOK, global.ApiResponse{
		Code:    0,
		Message: "Success",
		Datas:   schedule,
	})
}

// GetCenterScheduleRules 獲取老師在指定中心的排課規則
// @Summary 獲取老師在指定中心的排課規則
// @Tags Teacher
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param center_id path int true "中心 ID"
// @Success 200 {object} global.ApiResponse{data=[]models.ScheduleRule}
// @Router /api/v1/teacher/me/centers/{center_id}/schedule-rules [get]
func (ctl *TeacherController) GetCenterScheduleRules(ctx *gin.Context) {
	teacherID := ctx.GetUint(global.UserIDKey)
	if teacherID == 0 {
		ctx.JSON(http.StatusUnauthorized, global.ApiResponse{
			Code:    global.UNAUTHORIZED,
			Message: "Teacher ID required",
		})
		return
	}

	centerIDStr := ctx.Param("center_id")
	if centerIDStr == "" {
		ctx.JSON(http.StatusBadRequest, global.ApiResponse{
			Code:    global.BAD_REQUEST,
			Message: "Center ID is required",
		})
		return
	}

	// 解析 center_id
	var centerID uint
	if _, err := fmt.Sscanf(centerIDStr, "%d", &centerID); err != nil || centerID == 0 {
		ctx.JSON(http.StatusBadRequest, global.ApiResponse{
			Code:    global.BAD_REQUEST,
			Message: "Invalid center ID",
		})
		return
	}

	// 驗證老師是否屬於該中心
	membership, err := ctl.membershipRepo.GetActiveByTeacherAndCenter(ctx, teacherID, centerIDStr)
	if err != nil || membership == nil {
		ctx.JSON(http.StatusForbidden, global.ApiResponse{
			Code:    global.FORBIDDEN,
			Message: "You are not a member of this center",
		})
		return
	}

	// 獲取該中心的排課規則
	rules, err := ctl.scheduleRuleRepo.ListByCenterID(ctx, centerID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, global.ApiResponse{
			Code:    errInfos.SYSTEM_ERROR,
			Message: "Failed to get schedule rules",
		})
		return
	}

	// 過濾出該老師的課程，並格式化輸出
	type RuleResponse struct {
		ID                  uint   `json:"id"`
		Title               string `json:"title"`
		Weekday             int    `json:"weekday"`
		WeekdayText         string `json:"weekday_text"`
		StartTime           string `json:"start_time"`
		EndTime             string `json:"end_time"`
		EffectiveStartDate  string `json:"effective_start_date"`
		EffectiveEndDate    string `json:"effective_end_date"`
	}

	weekdayTexts := []string{"週日", "週一", "週二", "週三", "週四", "週五", "週六"}

	var teacherRules []RuleResponse
	for _, rule := range rules {
		if rule.TeacherID != nil && *rule.TeacherID == teacherID {
			title := rule.Offering.Name
			if title == "" {
				title = rule.Name
			}

			effectiveStartDate := ""
			effectiveEndDate := ""
			if !rule.EffectiveRange.StartDate.IsZero() {
				effectiveStartDate = rule.EffectiveRange.StartDate.Format("2006-01-02")
			}
			if !rule.EffectiveRange.EndDate.IsZero() {
				effectiveEndDate = rule.EffectiveRange.EndDate.Format("2006-01-02")
			}

			teacherRules = append(teacherRules, RuleResponse{
				ID:                 rule.ID,
				Title:              title,
				Weekday:            rule.Weekday,
				WeekdayText:        weekdayTexts[rule.Weekday],
				StartTime:          rule.StartTime,
				EndTime:            rule.EndTime,
				EffectiveStartDate: effectiveStartDate,
				EffectiveEndDate:   effectiveEndDate,
			})
		}
	}

	ctx.JSON(http.StatusOK, global.ApiResponse{
		Code:    0,
		Message: "Success",
		Datas:   teacherRules,
	})
}

// CreateException 老師提出停課/改期申請
// @Summary 老師提出停課/改期申請
// @Tags Teacher
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body TeacherCreateExceptionRequest true "例外申請"
// @Success 200 {object} global.ApiResponse{data=models.ScheduleException}
// @Router /api/v1/teacher/exceptions [post]
func (ctl *TeacherController) CreateException(ctx *gin.Context) {
	teacherID := ctx.GetUint(global.UserIDKey)
	if teacherID == 0 {
		ctx.JSON(http.StatusUnauthorized, global.ApiResponse{
			Code:    global.UNAUTHORIZED,
			Message: "Teacher ID required",
		})
		return
	}

	var req TeacherCreateExceptionRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, global.ApiResponse{
			Code:    global.BAD_REQUEST,
			Message: "Invalid request body",
		})
		return
	}

	exception, err := ctl.exceptionService.CreateException(
		ctx,
		req.CenterID,
		teacherID,
		req.RuleID,
		req.OriginalDate,
		req.Type,
		req.NewStartAt,
		req.NewEndAt,
		req.NewTeacherID,
		req.NewTeacherName,
		req.Reason,
	)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, global.ApiResponse{
			Code:    500,
			Message: err.Error(),
		})
		return
	}

	ctl.auditLogRepo.Create(ctx, models.AuditLog{
		CenterID:   req.CenterID,
		ActorType:  "TEACHER",
		ActorID:    teacherID,
		Action:     "EXCEPTION_CREATE",
		TargetType: "ScheduleException",
		TargetID:   exception.ID,
		Payload: models.AuditPayload{
			After: map[string]interface{}{
				"rule_id":       req.RuleID,
				"original_date": req.OriginalDate,
				"type":          req.Type,
				"reason":        req.Reason,
			},
		},
	})

	ctx.JSON(http.StatusOK, global.ApiResponse{
		Code:    0,
		Message: "Exception submitted",
		Datas:   exception,
	})
}

// RevokeException 老師撤回待審核的例外申請
// @Summary 老師撤回待審核的例外申請
// @Tags Teacher
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Exception ID"
// @Success 200 {object} global.ApiResponse
// @Router /api/v1/teacher/exceptions/{id}/revoke [post]
func (ctl *TeacherController) RevokeException(ctx *gin.Context) {
	teacherID := ctx.GetUint(global.UserIDKey)
	if teacherID == 0 {
		ctx.JSON(http.StatusUnauthorized, global.ApiResponse{
			Code:    global.UNAUTHORIZED,
			Message: "Teacher ID required",
		})
		return
	}

	id := ctx.Param("id")
	var exceptionID uint
	if _, err := fmt.Sscanf(id, "%d", &exceptionID); err != nil {
		ctx.JSON(http.StatusBadRequest, global.ApiResponse{
			Code:    global.BAD_REQUEST,
			Message: "Invalid exception ID",
		})
		return
	}

	err := ctl.exceptionService.RevokeException(ctx, exceptionID, teacherID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, global.ApiResponse{
			Code:    500,
			Message: err.Error(),
		})
		return
	}

	exception, _ := ctl.exceptionRepo.GetByID(ctx, exceptionID)
	if exception.CenterID > 0 {
		ctl.auditLogRepo.Create(ctx, models.AuditLog{
			CenterID:   exception.CenterID,
			ActorType:  "TEACHER",
			ActorID:    teacherID,
			Action:     "EXCEPTION_REVOKE",
			TargetType: "ScheduleException",
			TargetID:   exceptionID,
			Payload: models.AuditPayload{
				After: map[string]interface{}{
					"status": "REVOKED",
				},
			},
		})
	}

	ctx.JSON(http.StatusOK, global.ApiResponse{
		Code:    0,
		Message: "Exception revoked",
	})
}

// GetExceptions 老師查看自己的例外申請列表
// @Summary 老師查看自己的例外申請列表
// @Tags Teacher
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param status query string false "篩選狀態 (PENDING/APPROVED/REJECTED/REVOKED)"
// @Success 200 {object} global.ApiResponse{data=[]models.ScheduleException}
// @Router /api/v1/teacher/exceptions [get]
func (ctl *TeacherController) GetExceptions(ctx *gin.Context) {
	teacherID := ctx.GetUint(global.UserIDKey)
	if teacherID == 0 {
		ctx.JSON(http.StatusUnauthorized, global.ApiResponse{
			Code:    global.UNAUTHORIZED,
			Message: "Teacher ID required",
		})
		return
	}

	memberships, _ := ctl.membershipRepo.GetActiveByTeacherID(ctx, teacherID)
	var centerIDs []uint
	for _, m := range memberships {
		centerIDs = append(centerIDs, m.CenterID)
	}

	var exceptions []models.ScheduleException
	query := ctl.app.MySQL.RDB.WithContext(ctx).
		Table("schedule_exceptions").
		Select("schedule_exceptions.*").
		Joins("JOIN center_memberships ON center_memberships.center_id = schedule_exceptions.center_id").
		Where("center_memberships.teacher_id = ?", teacherID).
		Where("center_memberships.status = ?", "ACTIVE")

	if status := ctx.Query("status"); status != "" {
		// 支援新旧两种状态值（向后兼容）
		if status == "APPROVED" {
			query = query.Where("schedule_exceptions.status IN ('APPROVED', 'APPROVE')")
		} else if status == "REJECTED" {
			query = query.Where("schedule_exceptions.status IN ('REJECTED', 'REJECT')")
		} else {
			query = query.Where("schedule_exceptions.status = ?", status)
		}
	}

	query.Order("schedule_exceptions.created_at DESC").Find(&exceptions)

	ctx.JSON(http.StatusOK, global.ApiResponse{
		Code:    0,
		Message: "Success",
		Datas:   exceptions,
	})
}

type TeacherScheduleItem struct {
	ID         string `json:"id"`
	Type       string `json:"type"`
	Title      string `json:"title"`
	Date       string `json:"date"`
	StartTime  string `json:"start_time"`
	EndTime    string `json:"end_time"`
	RoomID     uint   `json:"room_id"`
	TeacherID  *uint  `json:"teacher_id"`
	CenterID   uint   `json:"center_id"`
	CenterName string `json:"center_name"`
	Status     string `json:"status"`
}

type TeacherCreateExceptionRequest struct {
	CenterID       uint       `json:"center_id" binding:"required"`
	RuleID         uint       `json:"rule_id" binding:"required"`
	OriginalDate   time.Time  `json:"original_date" binding:"required"`
	Type           string     `json:"type" binding:"required,oneof=CANCEL RESCHEDULE REPLACE_TEACHER"`
	NewStartAt     *time.Time `json:"new_start_at"`
	NewEndAt       *time.Time `json:"new_end_at"`
	NewTeacherID   *uint      `json:"new_teacher_id"`
	NewTeacherName string     `json:"new_teacher_name"`
	Reason         string     `json:"reason" binding:"required"`
}

// GetSessionNote 取得課堂筆記
// @Summary 取得課堂筆記
// @Tags Teacher
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param rule_id query uint true "課程規則ID"
// @Param session_date query string true "課程日期 (YYYY-MM-DD)"
// @Success 200 {object} global.ApiResponse{data=resources.SessionNoteResource}
// @Router /api/v1/teacher/sessions/note [get]
func (ctl *TeacherController) GetSessionNote(ctx *gin.Context) {
	teacherID := ctx.GetUint(global.UserIDKey)
	if teacherID == 0 {
		ctx.JSON(http.StatusUnauthorized, global.ApiResponse{
			Code:    global.UNAUTHORIZED,
			Message: "Teacher ID required",
		})
		return
	}

	ruleID, err := ctx.GetQuery("rule_id")
	if !err {
		ctx.JSON(http.StatusBadRequest, global.ApiResponse{
			Code:    global.BAD_REQUEST,
			Message: "rule_id required",
		})
		return
	}

	sessionDateStr, err := ctx.GetQuery("session_date")
	if !err {
		ctx.JSON(http.StatusBadRequest, global.ApiResponse{
			Code:    global.BAD_REQUEST,
			Message: "session_date required",
		})
		return
	}

	var ruleIDUint uint
	fmt.Sscanf(ruleID, "%d", &ruleIDUint)

	sessionDate, parseErr := time.Parse("2006-01-02", sessionDateStr)
	if parseErr != nil {
		ctx.JSON(http.StatusBadRequest, global.ApiResponse{
			Code:    global.BAD_REQUEST,
			Message: "Invalid session_date format, use YYYY-MM-DD",
		})
		return
	}

	noteRepo := repositories.NewSessionNoteRepository(ctl.app)
	note, isNew, dbErr := noteRepo.GetOrCreate(ctx, teacherID, ruleIDUint, sessionDate)
	if dbErr != nil {
		ctx.JSON(http.StatusInternalServerError, global.ApiResponse{
			Code:    500,
			Message: dbErr.Error(),
		})
		return
	}

	response := resources.SessionNoteResource{
		ID:          note.ID,
		RuleID:      note.RuleID,
		SessionDate: note.SessionDate.Format("2006-01-02"),
		Content:     note.Content,
		PrepNote:    note.PrepNote,
		UpdatedAt:   note.UpdatedAt,
	}

	ctx.JSON(http.StatusOK, global.ApiResponse{
		Code:    0,
		Message: "Success",
		Datas: map[string]interface{}{
			"note":   response,
			"is_new": isNew,
		},
	})
}

// UpsertSessionNote 新增或更新課堂筆記
// @Summary 新增或更新課堂筆記
// @Tags Teacher
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body requests.UpsertSessionNoteRequest true "筆記內容"
// @Success 200 {object} global.ApiResponse{data=resources.SessionNoteResource}
// @Router /api/v1/teacher/sessions/note [put]
func (ctl *TeacherController) UpsertSessionNote(ctx *gin.Context) {
	teacherID := ctx.GetUint(global.UserIDKey)
	if teacherID == 0 {
		ctx.JSON(http.StatusUnauthorized, global.ApiResponse{
			Code:    global.UNAUTHORIZED,
			Message: "Teacher ID required",
		})
		return
	}

	var req requests.UpsertSessionNoteRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, global.ApiResponse{
			Code:    global.BAD_REQUEST,
			Message: "Invalid request body: " + err.Error(),
		})
		return
	}

	sessionDate, parseErr := time.Parse("2006-01-02", req.SessionDate)
	if parseErr != nil {
		ctx.JSON(http.StatusBadRequest, global.ApiResponse{
			Code:    global.BAD_REQUEST,
			Message: "Invalid session_date format, use YYYY-MM-DD",
		})
		return
	}

	noteRepo := repositories.NewSessionNoteRepository(ctl.app)
	note, _, err := noteRepo.GetOrCreate(ctx, teacherID, req.RuleID, sessionDate)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, global.ApiResponse{
			Code:    500,
			Message: err.Error(),
		})
		return
	}

	note.Content = req.Content
	note.PrepNote = req.PrepNote
	note.UpdatedAt = time.Now()

	if err := noteRepo.Update(ctx, note); err != nil {
		ctx.JSON(http.StatusInternalServerError, global.ApiResponse{
			Code:    500,
			Message: err.Error(),
		})
		return
	}

	response := resources.SessionNoteResource{
		ID:          note.ID,
		RuleID:      note.RuleID,
		SessionDate: note.SessionDate.Format("2006-01-02"),
		Content:     note.Content,
		PrepNote:    note.PrepNote,
		UpdatedAt:   note.UpdatedAt,
	}

	ctx.JSON(http.StatusOK, global.ApiResponse{
		Code:    0,
		Message: "Note saved",
		Datas:   response,
	})
}

// GetSkills 取得老師技能列表
// @Summary 取得老師技能列表
// @Tags Teacher
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} global.ApiResponse{data=[]models.TeacherSkill}
// @Router /api/v1/teacher/me/skills [get]
func (ctl *TeacherController) GetSkills(ctx *gin.Context) {
	teacherID := ctx.GetUint(global.UserIDKey)
	if teacherID == 0 {
		ctx.JSON(http.StatusUnauthorized, global.ApiResponse{
			Code:    global.UNAUTHORIZED,
			Message: "Teacher ID required",
		})
		return
	}

	skills, err := ctl.skillRepo.ListByTeacherID(ctx, teacherID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, global.ApiResponse{
			Code:    500,
			Message: err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, global.ApiResponse{
		Code:    0,
		Message: "Success",
		Datas:   skills,
	})
}

// CreateSkill 新增老師技能
// @Summary 新增老師技能
// @Tags Teacher
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body CreateSkillRequest true "技能資訊"
// @Success 200 {object} global.ApiResponse{data=models.TeacherSkill}
// @Router /api/v1/teacher/me/skills [post]
func (ctl *TeacherController) CreateSkill(ctx *gin.Context) {
	teacherID := ctx.GetUint(global.UserIDKey)
	if teacherID == 0 {
		ctx.JSON(http.StatusUnauthorized, global.ApiResponse{
			Code:    global.UNAUTHORIZED,
			Message: "Teacher ID required",
		})
		return
	}

	var req CreateSkillRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, global.ApiResponse{
			Code:    global.BAD_REQUEST,
			Message: "Invalid request body: " + err.Error(),
		})
		return
	}

	skill := models.TeacherSkill{
		TeacherID: teacherID,
		Category:  req.Category,
		SkillName: req.SkillName,
		Level:     req.Level,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := ctl.skillRepo.Create(ctx, &skill); err != nil {
		ctx.JSON(http.StatusInternalServerError, global.ApiResponse{
			Code:    500,
			Message: err.Error(),
		})
		return
	}

	// 建立技能標籤關聯
	if len(req.HashtagIDs) > 0 {
		for _, hashtagID := range req.HashtagIDs {
			skillHashtag := models.TeacherSkillHashtag{
				TeacherSkillID: skill.ID,
				HashtagID:      hashtagID,
			}
			if err := ctl.app.MySQL.WDB.WithContext(ctx).Create(&skillHashtag).Error; err != nil {
				// 記錄錯誤但不影響主要流程
				fmt.Printf("Failed to create skill hashtag: %v\n", err)
			}
		}
	}

	// 重新載入技能（含標籤）
	skillHashtags := []models.TeacherSkillHashtag{}
	ctl.app.MySQL.RDB.WithContext(ctx).
		Where("teacher_skill_id = ?", skill.ID).
		Find(&skillHashtags)

	ctx.JSON(http.StatusOK, global.ApiResponse{
		Code:    0,
		Message: "Skill created",
		Datas:   skill,
	})
}

// DeleteSkill 刪除老師技能
// @Summary 刪除老師技能
// @Tags Teacher
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path uint true "技能ID"
// @Success 200 {object} global.ApiResponse
// @Router /api/v1/teacher/me/skills/{id} [delete]
func (ctl *TeacherController) DeleteSkill(ctx *gin.Context) {
	teacherID := ctx.GetUint(global.UserIDKey)
	if teacherID == 0 {
		ctx.JSON(http.StatusUnauthorized, global.ApiResponse{
			Code:    global.UNAUTHORIZED,
			Message: "Teacher ID required",
		})
		return
	}

	var id uint
	if _, err := fmt.Sscanf(ctx.Param("id"), "%d", &id); err != nil {
		ctx.JSON(http.StatusBadRequest, global.ApiResponse{
			Code:    global.BAD_REQUEST,
			Message: "Invalid skill ID",
		})
		return
	}

	skill, err := ctl.skillRepo.GetByID(ctx, id)
	if err != nil {
		ctx.JSON(http.StatusNotFound, global.ApiResponse{
			Code:    404,
			Message: "Skill not found",
		})
		return
	}

	if skill.TeacherID != teacherID {
		ctx.JSON(http.StatusForbidden, global.ApiResponse{
			Code:    403,
			Message: "Not authorized to delete this skill",
		})
		return
	}

	if err := ctl.skillRepo.Delete(ctx, id); err != nil {
		ctx.JSON(http.StatusInternalServerError, global.ApiResponse{
			Code:    500,
			Message: err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, global.ApiResponse{
		Code:    0,
		Message: "Skill deleted",
	})
}

// UpdateSkill 更新老師技能
// @Summary 更新老師技能
// @Tags Teacher
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path uint true "技能ID"
// @Param request body UpdateSkillRequest true "技能資訊"
// @Success 200 {object} global.ApiResponse{data=models.TeacherSkill}
// @Router /api/v1/teacher/me/skills/{id} [put]
func (ctl *TeacherController) UpdateSkill(ctx *gin.Context) {
	teacherID := ctx.GetUint(global.UserIDKey)
	if teacherID == 0 {
		ctx.JSON(http.StatusUnauthorized, global.ApiResponse{
			Code:    global.UNAUTHORIZED,
			Message: "Teacher ID required",
		})
		return
	}

	var id uint
	if _, err := fmt.Sscanf(ctx.Param("id"), "%d", &id); err != nil {
		ctx.JSON(http.StatusBadRequest, global.ApiResponse{
			Code:    global.BAD_REQUEST,
			Message: "Invalid skill ID",
		})
		return
	}

	var req UpdateSkillRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, global.ApiResponse{
			Code:    global.BAD_REQUEST,
			Message: "Invalid request body: " + err.Error(),
		})
		return
	}

	skill, err := ctl.skillRepo.GetByID(ctx, id)
	if err != nil {
		ctx.JSON(http.StatusNotFound, global.ApiResponse{
			Code:    404,
			Message: "Skill not found",
		})
		return
	}

	if skill.TeacherID != teacherID {
		ctx.JSON(http.StatusForbidden, global.ApiResponse{
			Code:    403,
			Message: "Not authorized to update this skill",
		})
		return
	}

	skill.Category = req.Category
	skill.SkillName = req.SkillName
	skill.UpdatedAt = time.Now()

	if err := ctl.skillRepo.Update(ctx, skill); err != nil {
		ctx.JSON(http.StatusInternalServerError, global.ApiResponse{
			Code:    500,
			Message: err.Error(),
		})
		return
	}

	// 更新技能標籤
	if len(req.Hashtags) > 0 {
		// 刪除現有標籤
		ctl.app.MySQL.WDB.WithContext(ctx).Where("teacher_skill_id = ?", skill.ID).Delete(&models.TeacherSkillHashtag{})

		hashtagRepo := repositories.NewHashtagRepository(ctl.app)
		for _, tagName := range req.Hashtags {
			// 確保 # 符號存在
			if !strings.HasPrefix(tagName, "#") {
				tagName = "#" + tagName
			}

			// 查找或創建標籤
			hashtag, err := hashtagRepo.GetByName(ctx, tagName)
			if err != nil {
				hashtag = &models.Hashtag{Name: tagName, UsageCount: 1}
				hashtagRepo.Create(ctx, hashtag)
			} else {
				hashtagRepo.IncrementUsage(ctx, tagName)
			}

			// 創建關聯
			ctl.app.MySQL.WDB.WithContext(ctx).Create(&models.TeacherSkillHashtag{
				TeacherSkillID: skill.ID,
				HashtagID:      hashtag.ID,
			})
		}
	}

	ctx.JSON(http.StatusOK, global.ApiResponse{
		Code:    0,
		Message: "Skill updated",
		Datas:   skill,
	})
}

// GetCertificates 取得老師證照列表
// @Summary 取得老師證照列表
// @Tags Teacher
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} global.ApiResponse{data=[]models.TeacherCertificate}
// @Router /api/v1/teacher/me/certificates [get]
func (ctl *TeacherController) GetCertificates(ctx *gin.Context) {
	teacherID := ctx.GetUint(global.UserIDKey)
	if teacherID == 0 {
		ctx.JSON(http.StatusUnauthorized, global.ApiResponse{
			Code:    global.UNAUTHORIZED,
			Message: "Teacher ID required",
		})
		return
	}

	certificates, err := ctl.certificateRepo.ListByTeacherID(ctx, teacherID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, global.ApiResponse{
			Code:    500,
			Message: err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, global.ApiResponse{
		Code:    0,
		Message: "Success",
		Datas:   certificates,
	})
}

// CreateCertificate 新增老師證照
// @Summary 新增老師證照
// @Tags Teacher
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body CreateCertificateRequest true "證照資訊"
// @Success 200 {object} global.ApiResponse{data=models.TeacherCertificate}
// @Router /api/v1/teacher/me/certificates [post]
func (ctl *TeacherController) CreateCertificate(ctx *gin.Context) {
	teacherID := ctx.GetUint(global.UserIDKey)
	if teacherID == 0 {
		ctx.JSON(http.StatusUnauthorized, global.ApiResponse{
			Code:    global.UNAUTHORIZED,
			Message: "Teacher ID required",
		})
		return
	}

	var req CreateCertificateRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, global.ApiResponse{
			Code:    global.BAD_REQUEST,
			Message: "Invalid request body: " + err.Error(),
		})
		return
	}

	certificate := models.TeacherCertificate{
		TeacherID: teacherID,
		Name:      req.Name,
		FileURL:   req.FileURL,
		IssuedAt:  req.IssuedAt,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := ctl.certificateRepo.Create(ctx, &certificate); err != nil {
		ctx.JSON(http.StatusInternalServerError, global.ApiResponse{
			Code:    500,
			Message: err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, global.ApiResponse{
		Code:    0,
		Message: "Certificate created",
		Datas:   certificate,
	})
}

// DeleteCertificate 刪除老師證照
// @Summary 刪除老師證照
// @Tags Teacher
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path uint true "證照ID"
// @Success 200 {object} global.ApiResponse
// @Router /api/v1/teacher/me/certificates/{id} [delete]
func (ctl *TeacherController) DeleteCertificate(ctx *gin.Context) {
	teacherID := ctx.GetUint(global.UserIDKey)
	if teacherID == 0 {
		ctx.JSON(http.StatusUnauthorized, global.ApiResponse{
			Code:    global.UNAUTHORIZED,
			Message: "Teacher ID required",
		})
		return
	}

	var id uint
	if _, err := fmt.Sscanf(ctx.Param("id"), "%d", &id); err != nil {
		ctx.JSON(http.StatusBadRequest, global.ApiResponse{
			Code:    global.BAD_REQUEST,
			Message: "Invalid certificate ID",
		})
		return
	}

	certificate, err := ctl.certificateRepo.GetByID(ctx, id)
	if err != nil {
		ctx.JSON(http.StatusNotFound, global.ApiResponse{
			Code:    404,
			Message: "Certificate not found",
		})
		return
	}

	if certificate.TeacherID != teacherID {
		ctx.JSON(http.StatusForbidden, global.ApiResponse{
			Code:    403,
			Message: "Not authorized to delete this certificate",
		})
		return
	}

	if err := ctl.certificateRepo.Delete(ctx, id); err != nil {
		ctx.JSON(http.StatusInternalServerError, global.ApiResponse{
			Code:    500,
			Message: err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, global.ApiResponse{
		Code:    0,
		Message: "Certificate deleted",
	})
}

// GetPersonalEvents 取得老師個人行程列表
// @Summary 取得老師個人行程列表
// @Tags Teacher
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} global.ApiResponse{data=[]models.PersonalEvent}
// @Router /api/v1/teacher/me/personal-events [get]
func (ctl *TeacherController) GetPersonalEvents(ctx *gin.Context) {
	teacherID := ctx.GetUint(global.UserIDKey)
	if teacherID == 0 {
		ctx.JSON(http.StatusUnauthorized, global.ApiResponse{
			Code:    global.UNAUTHORIZED,
			Message: "Teacher ID required",
		})
		return
	}

	fromStr := ctx.Query("from")
	toStr := ctx.Query("to")

	var events []models.PersonalEvent
	var err error

	if fromStr != "" && toStr != "" {
		from, err := time.Parse("2006-01-02", fromStr)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, global.ApiResponse{
				Code:    400,
				Message: "Invalid from date format",
			})
			return
		}
		to, err := time.Parse("2006-01-02", toStr)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, global.ApiResponse{
				Code:    400,
				Message: "Invalid to date format",
			})
			return
		}
		// Add one day to to date to make it inclusive
		to = to.AddDate(0, 0, 1)
		events, err = ctl.personalEventRepo.GetByTeacherAndDateRange(ctx, teacherID, from, to)
	} else {
		events, err = ctl.personalEventRepo.ListByTeacherID(ctx, teacherID)
	}

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, global.ApiResponse{
			Code:    500,
			Message: err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, global.ApiResponse{
		Code:    0,
		Message: "Success",
		Datas:   events,
	})
}

// CreatePersonalEvent 新增老師個人行程
// @Summary 新增老師個人行程
// @Tags Teacher
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body CreatePersonalEventRequest true "行程資訊"
// @Success 200 {object} global.ApiResponse{data=models.PersonalEvent}
// @Failure 400 {object} global.ApiResponse
// @Failure 409 {object} global.ApiResponse
// @Router /api/v1/teacher/me/personal-events [post]
func (ctl *TeacherController) CreatePersonalEvent(ctx *gin.Context) {
	teacherID := ctx.GetUint(global.UserIDKey)
	if teacherID == 0 {
		ctx.JSON(http.StatusUnauthorized, global.ApiResponse{
			Code:    global.UNAUTHORIZED,
			Message: "Teacher ID required",
		})
		return
	}

	var req CreatePersonalEventRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, global.ApiResponse{
			Code:    global.BAD_REQUEST,
			Message: "Invalid request body: " + err.Error(),
		})
		return
	}

	// 檢查個人行程是否與中心課程衝突
	// 取得老師所屬的所有中心
	memberships, err := ctl.membershipRepo.GetActiveByTeacherID(ctx, teacherID)
	if err != nil {
		// 靜默處理錯誤，繼續執行
	}
	if len(memberships) == 0 {
		// 老師沒有加入任何中心，創建個人行程
	} else {
		var centerIDs []uint
		for _, m := range memberships {
			centerIDs = append(centerIDs, m.CenterID)
		}

		// 檢查每個中心的課程衝突
		for _, centerID := range centerIDs {
			conflicts, err := ctl.scheduleRuleRepo.CheckPersonalEventConflict(ctx, teacherID, centerID, req.StartAt, req.EndAt)
			if err != nil {
				continue
			}

			if len(conflicts) > 0 {
				// 發現衝突，阻擋操作並返回錯誤
				conflictMessages := []string{}
				for _, rule := range conflicts {
					conflictMessages = append(conflictMessages, fmt.Sprintf(
						"您於 %s %s-%s 在中心 %d 有課程「%s」的安排，時間衝突",
						req.StartAt.Format("2006-01-02"),
						rule.StartTime,
						rule.EndTime,
						centerID,
						rule.Offering.Name,
					))
				}
				ctx.JSON(http.StatusConflict, global.ApiResponse{
					Code:    409,
					Message: strings.Join(conflictMessages, "; "),
				})
				return
			}
		}
	}

	event := models.PersonalEvent{
		TeacherID: teacherID,
		Title:     req.Title,
		StartAt:   req.StartAt,
		EndAt:     req.EndAt,
		IsAllDay:  req.IsAllDay,
		ColorHex:  req.ColorHex,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if req.RecurrenceRule != nil {
		event.RecurrenceRule = *req.RecurrenceRule
	}

	if err := ctl.personalEventRepo.Create(ctx, &event); err != nil {
		ctx.JSON(http.StatusInternalServerError, global.ApiResponse{
			Code:    500,
			Message: err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, global.ApiResponse{
		Code:    0,
		Message: "Personal event created",
		Datas:   event,
	})
}

// DeletePersonalEvent 刪除老師個人行程
// @Summary 刪除老師個人行程
// @Tags Teacher
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path uint true "行程ID"
// @Success 200 {object} global.ApiResponse
// @Router /api/v1/teacher/me/personal-events/{id} [delete]
func (ctl *TeacherController) DeletePersonalEvent(ctx *gin.Context) {
	teacherID := ctx.GetUint(global.UserIDKey)
	if teacherID == 0 {
		ctx.JSON(http.StatusUnauthorized, global.ApiResponse{
			Code:    global.UNAUTHORIZED,
			Message: "Teacher ID required",
		})
		return
	}

	var id uint
	if _, err := fmt.Sscanf(ctx.Param("id"), "%d", &id); err != nil {
		ctx.JSON(http.StatusBadRequest, global.ApiResponse{
			Code:    global.BAD_REQUEST,
			Message: "Invalid personal event ID",
		})
		return
	}

	event, err := ctl.personalEventRepo.GetByID(ctx, id)
	if err != nil {
		ctx.JSON(http.StatusNotFound, global.ApiResponse{
			Code:    404,
			Message: "Personal event not found",
		})
		return
	}

	if event.TeacherID != teacherID {
		ctx.JSON(http.StatusForbidden, global.ApiResponse{
			Code:    403,
			Message: "Not authorized to delete this personal event",
		})
		return
	}

	if err := ctl.personalEventRepo.Delete(ctx, id); err != nil {
		ctx.JSON(http.StatusInternalServerError, global.ApiResponse{
			Code:    500,
			Message: err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, global.ApiResponse{
		Code:    0,
		Message: "Personal event deleted",
	})
}

type CreateSkillRequest struct {
	Category   string `json:"category" binding:"required"`
	SkillName  string `json:"skill_name" binding:"required"`
	Level      string `json:"level"`
	HashtagIDs []uint `json:"hashtag_ids"`
}

type UpdateSkillRequest struct {
	Category   string `json:"category" binding:"required"`
	SkillName  string `json:"skill_name" binding:"required"`
	Hashtags   []string `json:"hashtags"`
}

type CreateCertificateRequest struct {
	Name     string    `json:"name" binding:"required"`
	FileURL  string    `json:"file_url" binding:"required"`
	IssuedAt time.Time `json:"issued_at" binding:"required"`
}

type CreatePersonalEventRequest struct {
	Title          string                 `json:"title" binding:"required"`
	StartAt        time.Time              `json:"start_at" binding:"required"`
	EndAt          time.Time              `json:"end_at" binding:"required"`
	IsAllDay       bool                   `json:"is_all_day"`
	ColorHex       string                 `json:"color_hex"`
	RecurrenceRule *models.RecurrenceRule `json:"recurrence_rule"`
}

type UpdatePersonalEventRequest struct {
	Title          *string                `json:"title"`
	StartAt        *time.Time             `json:"start_at"`
	EndAt          *time.Time             `json:"end_at"`
	IsAllDay       *bool                  `json:"is_all_day"`
	ColorHex       *string                `json:"color_hex"`
	RecurrenceRule *models.RecurrenceRule `json:"recurrence_rule"`
	UpdateMode     string     `json:"update_mode" binding:"required,oneof=SINGLE FUTURE ALL"`
}

type UpdatePersonalEventResponse struct {
	UpdatedCount int64  `json:"updated_count"`
	Message      string `json:"message"`
}

// UpdatePersonalEvent 更新老師個人行程
// @Summary 更新老師個人行程
// @Tags Teacher
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path uint true "行程ID"
// @Param request body UpdatePersonalEventRequest true "行程資訊"
// @Success 200 {object} global.ApiResponse{data=UpdatePersonalEventResponse}
// @Router /api/v1/teacher/me/personal-events/{id} [patch]
func (ctl *TeacherController) UpdatePersonalEvent(ctx *gin.Context) {
	teacherID := ctx.GetUint(global.UserIDKey)
	if teacherID == 0 {
		ctx.JSON(http.StatusUnauthorized, global.ApiResponse{
			Code:    global.UNAUTHORIZED,
			Message: "Teacher ID required",
		})
		return
	}

	var id uint
	if _, err := fmt.Sscanf(ctx.Param("id"), "%d", &id); err != nil {
		ctx.JSON(http.StatusBadRequest, global.ApiResponse{
			Code:    global.BAD_REQUEST,
			Message: "Invalid personal event ID",
		})
		return
	}

	event, err := ctl.personalEventRepo.GetByID(ctx, id)
	if err != nil {
		ctx.JSON(http.StatusNotFound, global.ApiResponse{
			Code:    errInfos.NOT_FOUND,
			Message: ctl.app.Err.New(errInfos.NOT_FOUND).Msg,
		})
		return
	}

	if event.TeacherID != teacherID {
		ctx.JSON(http.StatusForbidden, global.ApiResponse{
			Code:    global.FORBIDDEN,
			Message: "Not authorized to update this personal event",
		})
		return
	}

	var req UpdatePersonalEventRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, global.ApiResponse{
			Code:    global.BAD_REQUEST,
			Message: "Invalid request body: " + err.Error(),
		})
		return
	}

	// 檢查個人行程是否與中心課程衝突（如果時間有變更）
	if req.StartAt != nil && req.EndAt != nil {
		// 取得老師所屬的所有中心
		memberships, err := ctl.membershipRepo.GetActiveByTeacherID(ctx, teacherID)
		if err != nil {
			// 靜默處理錯誤，繼續執行
		}
		if len(memberships) == 0 {
			// 老師沒有加入任何中心
		} else {
			var centerIDs []uint
			for _, m := range memberships {
				centerIDs = append(centerIDs, m.CenterID)
			}

			// 檢查每個中心的課程衝突
			for _, centerID := range centerIDs {
				conflicts, err := ctl.scheduleRuleRepo.CheckPersonalEventConflict(ctx, teacherID, centerID, *req.StartAt, *req.EndAt)
				if err != nil {
					continue
				}

				if len(conflicts) > 0 {
					// 發現衝突，阻擋操作並返回錯誤
					conflictMessages := []string{}
					for _, rule := range conflicts {
						conflictMessages = append(conflictMessages, fmt.Sprintf(
							"您於 %s %s-%s 在中心 %d 有課程「%s」的安排，時間衝突",
							req.StartAt.Format("2006-01-02"),
							rule.StartTime,
							rule.EndTime,
							centerID,
							rule.Offering.Name,
						))
					}
					ctx.JSON(http.StatusConflict, global.ApiResponse{
						Code:    409,
						Message: strings.Join(conflictMessages, "; "),
					})
					return
				}
			}
		}
	}

	now := time.Now()
	var updatedCount int64 = 1

	switch req.UpdateMode {
	case "SINGLE":
		if req.Title != nil {
			event.Title = *req.Title
		}
		if req.StartAt != nil {
			event.StartAt = *req.StartAt
		}
		if req.EndAt != nil {
			event.EndAt = *req.EndAt
		}
		if req.IsAllDay != nil {
			event.IsAllDay = *req.IsAllDay
		}
		if req.ColorHex != nil {
			event.ColorHex = *req.ColorHex
		}
		if req.RecurrenceRule != nil {
			event.RecurrenceRule = *req.RecurrenceRule
		}
		event.UpdatedAt = now
		if err := ctl.personalEventRepo.Update(ctx, event); err != nil {
			ctx.JSON(http.StatusInternalServerError, global.ApiResponse{
				Code:    errInfos.SQL_ERROR,
				Message: ctl.app.Err.New(errInfos.SQL_ERROR).Msg,
			})
			return
		}

	case "FUTURE":
		if event.RecurrenceRule.Type == "" {
			ctx.JSON(http.StatusBadRequest, global.ApiResponse{
				Code:    errInfos.PARAMS_VALIDATE_ERROR,
				Message: "update_mode FUTURE requires a recurring event",
			})
			return
		}

		repoReq := repositories.UpdateEventRequest{
			Title:    req.Title,
			StartAt:  req.StartAt,
			EndAt:    req.EndAt,
			IsAllDay: req.IsAllDay,
			ColorHex: req.ColorHex,
		}
		updatedCount, err = ctl.personalEventRepo.UpdateFutureOccurrences(ctx, id, teacherID, repoReq, now)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, global.ApiResponse{
				Code:    errInfos.SQL_ERROR,
				Message: ctl.app.Err.New(errInfos.SQL_ERROR).Msg,
			})
			return
		}

	case "ALL":
		if event.RecurrenceRule.Type == "" {
			ctx.JSON(http.StatusBadRequest, global.ApiResponse{
				Code:    errInfos.PARAMS_VALIDATE_ERROR,
				Message: "update_mode ALL requires a recurring event",
			})
			return
		}

		repoReq := repositories.UpdateEventRequest{
			Title:    req.Title,
			StartAt:  req.StartAt,
			EndAt:    req.EndAt,
			IsAllDay: req.IsAllDay,
			ColorHex: req.ColorHex,
		}
		updatedCount, err = ctl.personalEventRepo.UpdateAllOccurrences(ctx, id, teacherID, repoReq, now)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, global.ApiResponse{
				Code:    errInfos.SQL_ERROR,
				Message: ctl.app.Err.New(errInfos.SQL_ERROR).Msg,
			})
			return
		}
	}

	ctx.JSON(http.StatusOK, global.ApiResponse{
		Code:    0,
		Message: "Personal event updated",
		Datas: UpdatePersonalEventResponse{
			UpdatedCount: updatedCount,
			Message:      "Updated " + string(req.UpdateMode) + " occurrence(s)",
		},
	})
}

// ListTeachers 取得老師列表（根據當前登入 Admin 的中心過濾）
// @Summary 取得老師列表
// @Tags Admin
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} global.ApiResponse{data=[]TeacherResponse}
// @Router /api/v1/teachers [get]
func (ctl *TeacherController) ListTeachers(ctx *gin.Context) {
	// 從 JWT Token 取得 center_id
	centerID := ctx.GetUint(global.CenterIDKey)
	if centerID == 0 {
		ctx.JSON(http.StatusUnauthorized, global.ApiResponse{
			Code:    global.UNAUTHORIZED,
			Message: "Center ID required",
		})
		return
	}

	// 取得該中心的所有會員老師 ID（包含 ACTIVE 和 INVITED 狀態）
	teacherIDs, err := ctl.membershipRepo.ListTeacherIDsByCenterID(ctx, centerID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, global.ApiResponse{
			Code:    500,
			Message: "Failed to get teacher IDs",
		})
		return
	}

	if len(teacherIDs) == 0 {
		ctx.JSON(http.StatusOK, global.ApiResponse{
			Code:    0,
			Message: "Success",
			Datas:   []TeacherResponse{},
		})
		return
	}

	// 取得老師詳細資料
	teachers := make([]TeacherResponse, 0, len(teacherIDs))
	for _, teacherID := range teacherIDs {
		teacher, err := ctl.teacherRepository.GetByID(ctx, teacherID)
		if err != nil {
			continue
		}

		teachers = append(teachers, TeacherResponse{
			ID:        teacher.ID,
			Name:      teacher.Name,
			Email:     teacher.Email,
			City:      teacher.City,
			District:  teacher.District,
			Bio:       teacher.Bio,
			CreatedAt: teacher.CreatedAt,
		})
	}

	ctx.JSON(http.StatusOK, global.ApiResponse{
		Code:    0,
		Message: "Success",
		Datas:   teachers,
	})
}

// GetCenterTeachers 取得指定中心的老師列表（供老師申請代課時選擇）
// @Summary 取得指定中心的老師列表
// @Tags Teacher
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param center_id path int true "Center ID"
// @Success 200 {object} global.ApiResponse{data=[]TeacherResponse}
// @Router /api/v1/teacher/me/centers/{center_id}/teachers [get]
func (ctl *TeacherController) GetCenterTeachers(ctx *gin.Context) {
	// 從 URL 參數取得 center_id
	centerIDStr := ctx.Param("center_id")
	var centerID uint
	if _, err := fmt.Sscanf(centerIDStr, "%d", &centerID); err != nil {
		ctx.JSON(http.StatusBadRequest, global.ApiResponse{
			Code:    global.BAD_REQUEST,
			Message: "Invalid center ID",
		})
		return
	}

	// 驗證老師是否屬於該中心
	teacherID := ctx.GetUint(global.UserIDKey)
	if teacherID == 0 {
		ctx.JSON(http.StatusUnauthorized, global.ApiResponse{
			Code:    global.UNAUTHORIZED,
			Message: "Teacher ID required",
		})
		return
	}

	// 檢查老師是否為該中心的會員
	memberships, err := ctl.membershipRepo.GetActiveByTeacherID(ctx, teacherID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, global.ApiResponse{
			Code:    500,
			Message: "Failed to get memberships",
		})
		return
	}

	isMember := false
	for _, m := range memberships {
		if m.CenterID == centerID && (m.Status == "ACTIVE" || m.Status == "INVITED") {
			isMember = true
			break
		}
	}

	if !isMember {
		ctx.JSON(http.StatusForbidden, global.ApiResponse{
			Code:    global.FORBIDDEN,
			Message: "You are not a member of this center",
		})
		return
	}

	// 取得該中心的所有會員老師 ID（排除自己）
	teacherIDs, err := ctl.membershipRepo.ListTeacherIDsByCenterID(ctx, centerID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, global.ApiResponse{
			Code:    500,
			Message: "Failed to get teacher IDs",
		})
		return
	}

	if len(teacherIDs) == 0 {
		ctx.JSON(http.StatusOK, global.ApiResponse{
			Code:    0,
			Message: "Success",
			Datas:   []TeacherResponse{},
		})
		return
	}

	// 取得老師詳細資料（排除自己）
	teachers := make([]TeacherResponse, 0, len(teacherIDs))
	for _, tID := range teacherIDs {
		// 排除自己
		if tID == teacherID {
			continue
		}

		teacher, err := ctl.teacherRepository.GetByID(ctx, tID)
		if err != nil {
			continue
		}

		teachers = append(teachers, TeacherResponse{
			ID:        teacher.ID,
			Name:      teacher.Name,
			Email:     teacher.Email,
			City:      teacher.City,
			District:  teacher.District,
			Bio:       teacher.Bio,
			CreatedAt: teacher.CreatedAt,
		})
	}

	ctx.JSON(http.StatusOK, global.ApiResponse{
		Code:    0,
		Message: "Success",
		Datas:   teachers,
	})
}

// DeleteTeacher 刪除老師
// @Summary 刪除老師
// @Tags Admin
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Teacher ID"
// @Success 200 {object} global.ApiResponse
// @Router /api/v1/teachers/{id} [delete]
func (ctl *TeacherController) DeleteTeacher(ctx *gin.Context) {
	adminID := ctx.GetUint(global.UserIDKey)
	if adminID == 0 {
		ctx.JSON(http.StatusUnauthorized, global.ApiResponse{
			Code:    global.UNAUTHORIZED,
			Message: "Admin ID required",
		})
		return
	}

	var teacherID uint
	if _, err := fmt.Sscanf(ctx.Param("id"), "%d", &teacherID); err != nil {
		ctx.JSON(http.StatusBadRequest, global.ApiResponse{
			Code:    global.BAD_REQUEST,
			Message: "Invalid teacher ID",
		})
		return
	}

	teacher, err := ctl.teacherRepository.GetByID(ctx, teacherID)
	if err != nil {
		ctx.JSON(http.StatusNotFound, global.ApiResponse{
			Code:    404,
			Message: "Teacher not found",
		})
		return
	}

	if err := ctl.teacherRepository.DeleteByID(ctx, teacherID); err != nil {
		ctx.JSON(http.StatusInternalServerError, global.ApiResponse{
			Code:    500,
			Message: "Failed to delete teacher",
		})
		return
	}

	memberships, _ := ctl.membershipRepo.GetActiveByTeacherID(ctx, teacherID)
	var centerID uint
	if len(memberships) > 0 {
		centerID = memberships[0].CenterID
	}

	ctl.auditLogRepo.Create(ctx, models.AuditLog{
		CenterID:   centerID,
		ActorType:  "ADMIN",
		ActorID:    adminID,
		Action:     "DELETE_TEACHER",
		TargetType: "Teacher",
		TargetID:   teacherID,
		Payload: models.AuditPayload{
			After: map[string]interface{}{
				"status": "DELETED",
				"name":   teacher.Name,
				"email":  teacher.Email,
			},
		},
	})

	ctx.JSON(http.StatusOK, global.ApiResponse{
		Code:    0,
		Message: "Teacher deleted",
	})
}

// InviteTeacher 邀請老師加入中心
// @Summary 邀請老師加入中心
// @Tags Admin
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Center ID"
// @Param request body InviteTeacherRequest true "邀請資訊"
// @Success 200 {object} global.ApiResponse{data=models.CenterInvitation}
// @Router /api/v1/admin/centers/{id}/invitations [post]
func (ctl *TeacherController) InviteTeacher(ctx *gin.Context) {
	adminID := ctx.GetUint(global.UserIDKey)
	if adminID == 0 {
		ctx.JSON(http.StatusUnauthorized, global.ApiResponse{
			Code:    global.UNAUTHORIZED,
			Message: "Admin ID required",
		})
		return
	}

	centerIDStr := ctx.Param("id")
	if centerIDStr == "" {
		ctx.JSON(http.StatusBadRequest, global.ApiResponse{
			Code:    global.BAD_REQUEST,
			Message: "Center ID required",
		})
		return
	}

	var centerID uint
	if _, err := fmt.Sscanf(centerIDStr, "%d", &centerID); err != nil {
		ctx.JSON(http.StatusBadRequest, global.ApiResponse{
			Code:    global.BAD_REQUEST,
			Message: "Invalid center ID",
		})
		return
	}

	var req InviteTeacherRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, global.ApiResponse{
			Code:    global.BAD_REQUEST,
			Message: "Invalid request body: " + err.Error(),
		})
		return
	}

	token := generateInviteToken()
	expiresAt := time.Now().Add(72 * time.Hour)

	invitation := models.CenterInvitation{
		CenterID:  centerID,
		Email:     req.Email,
		Token:     token,
		Status:    "PENDING",
		CreatedAt: time.Now(),
		ExpiresAt: expiresAt,
	}

	if err := ctl.app.MySQL.WDB.WithContext(ctx).Create(&invitation).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, global.ApiResponse{
			Code:    500,
			Message: "Failed to create invitation",
		})
		return
	}

	ctl.auditLogRepo.Create(ctx, models.AuditLog{
		CenterID:   centerID,
		ActorType:  "ADMIN",
		ActorID:    adminID,
		Action:     "INVITE_TEACHER",
		TargetType: "CenterInvitation",
		TargetID:   invitation.ID,
		Payload: models.AuditPayload{
			After: map[string]interface{}{
				"email":      req.Email,
				"status":     "PENDING",
				"expires_at": expiresAt,
			},
		},
	})

	ctx.JSON(http.StatusOK, global.ApiResponse{
		Code:    0,
		Message: "Invitation sent",
		Datas:   invitation,
	})
}

type InviteTeacherRequest struct {
	Email   string `json:"email" binding:"required,email"`
	Role    string `json:"role" binding:"required,oneof=TEACHER SUBSTITUTE"`
	Message string `json:"message"`
}

func generateInviteToken() string {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return ""
	}
	return fmt.Sprintf("%x", b)
}

type CheckTeacherRuleLockRequest struct {
	RuleID        uint   `json:"rule_id" binding:"required"`
	ExceptionDate string `json:"exception_date" binding:"required"`
}

type CheckTeacherRuleLockResponse struct {
	IsLocked      bool       `json:"is_locked"`
	LockReason    string     `json:"lock_reason,omitempty"`
	Deadline      *time.Time `json:"deadline,omitempty"`
	DaysRemaining int        `json:"days_remaining"`
}

func (ctl *TeacherController) CheckRuleLockStatus(ctx *gin.Context) {
	teacherID := ctx.GetUint(global.UserIDKey)
	if teacherID == 0 {
		ctx.JSON(http.StatusUnauthorized, global.ApiResponse{
			Code:    errInfos.UNAUTHORIZED,
			Message: "Teacher ID not found",
		})
		return
	}

	var req CheckTeacherRuleLockRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, global.ApiResponse{
			Code:    errInfos.PARAMS_VALIDATE_ERROR,
			Message: "Invalid request parameters",
		})
		return
	}

	rule, err := ctl.scheduleRuleRepo.GetByID(ctx, req.RuleID)
	if err != nil {
		ctx.JSON(http.StatusNotFound, global.ApiResponse{
			Code:    errInfos.NOT_FOUND,
			Message: "Rule not found",
		})
		return
	}

	exceptionDate, err := time.Parse("2006-01-02", req.ExceptionDate)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, global.ApiResponse{
			Code:    errInfos.PARAMS_VALIDATE_ERROR,
			Message: "Invalid date format, expected YYYY-MM-DD",
		})
		return
	}

	allowed, reasonStr, _ := ctl.exceptionService.CheckExceptionDeadline(ctx, rule.CenterID, req.RuleID, exceptionDate)

	response := CheckTeacherRuleLockResponse{
		IsLocked:   !allowed,
		LockReason: reasonStr,
	}

	if !allowed {
		ctx.JSON(http.StatusOK, global.ApiResponse{
			Code:    0,
			Message: "Rule is locked",
			Datas:   response,
		})
		return
	}

	ctx.JSON(http.StatusOK, global.ApiResponse{
		Code:    0,
		Message: "Rule is available for exception",
		Datas:   response,
	})
}

type PreviewRecurrenceEditRequest struct {
	RuleID   uint   `json:"rule_id" binding:"required"`
	EditDate string `json:"edit_date" binding:"required"`
	Mode     string `json:"mode" binding:"required,oneof=SINGLE FUTURE ALL"`
}

func (ctl *TeacherController) PreviewRecurrenceEdit(ctx *gin.Context) {
	teacherID := ctx.GetUint(global.UserIDKey)
	if teacherID == 0 {
		ctx.JSON(http.StatusUnauthorized, global.ApiResponse{
			Code:    errInfos.UNAUTHORIZED,
			Message: "Teacher ID not found",
		})
		return
	}

	var req PreviewRecurrenceEditRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, global.ApiResponse{
			Code:    errInfos.PARAMS_VALIDATE_ERROR,
			Message: "Invalid request parameters",
		})
		return
	}

	editDate, err := time.Parse("2006-01-02", req.EditDate)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, global.ApiResponse{
			Code:    errInfos.PARAMS_VALIDATE_ERROR,
			Message: "Invalid date format, expected YYYY-MM-DD",
		})
		return
	}

	preview, err := ctl.recurrenceService.PreviewAffectedSessions(ctx, req.RuleID, editDate, services.RecurrenceEditMode(req.Mode))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, global.ApiResponse{
			Code:    errInfos.SYSTEM_ERROR,
			Message: "Failed to preview affected sessions",
		})
		return
	}

	ctx.JSON(http.StatusOK, global.ApiResponse{
		Code:    0,
		Message: "Preview generated",
		Datas:   preview,
	})
}

func (ctl *TeacherController) EditRecurringSchedule(ctx *gin.Context) {
	teacherID := ctx.GetUint(global.UserIDKey)
	if teacherID == 0 {
		ctx.JSON(http.StatusUnauthorized, global.ApiResponse{
			Code:    errInfos.UNAUTHORIZED,
			Message: "Teacher ID not found",
		})
		return
	}

	var req services.RecurrenceEditRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, global.ApiResponse{
			Code:    errInfos.PARAMS_VALIDATE_ERROR,
			Message: "Invalid request parameters",
		})
		return
	}

	rule, err := ctl.scheduleRuleRepo.GetByID(ctx, req.RuleID)
	if err != nil {
		ctx.JSON(http.StatusNotFound, global.ApiResponse{
			Code:    errInfos.NOT_FOUND,
			Message: "Rule not found",
		})
		return
	}

	result, err := ctl.recurrenceService.EditRecurringSchedule(ctx, rule.CenterID, teacherID, req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, global.ApiResponse{
			Code:    errInfos.SYSTEM_ERROR,
			Message: err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, global.ApiResponse{
		Code:    0,
		Message: "Schedule edited successfully",
		Datas:   result,
	})
}

type DeleteRecurringScheduleRequest struct {
	RuleID   uint   `json:"rule_id" binding:"required"`
	EditDate string `json:"edit_date" binding:"required"`
	Mode     string `json:"mode" binding:"required,oneof=SINGLE FUTURE ALL"`
	Reason   string `json:"reason"`
}

func (ctl *TeacherController) DeleteRecurringSchedule(ctx *gin.Context) {
	teacherID := ctx.GetUint(global.UserIDKey)
	if teacherID == 0 {
		ctx.JSON(http.StatusUnauthorized, global.ApiResponse{
			Code:    errInfos.UNAUTHORIZED,
			Message: "Teacher ID not found",
		})
		return
	}

	var req DeleteRecurringScheduleRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, global.ApiResponse{
			Code:    errInfos.PARAMS_VALIDATE_ERROR,
			Message: "Invalid request parameters",
		})
		return
	}

	editDate, err := time.Parse("2006-01-02", req.EditDate)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, global.ApiResponse{
			Code:    errInfos.PARAMS_VALIDATE_ERROR,
			Message: "Invalid date format, expected YYYY-MM-DD",
		})
		return
	}

	rule, err := ctl.scheduleRuleRepo.GetByID(ctx, req.RuleID)
	if err != nil {
		ctx.JSON(http.StatusNotFound, global.ApiResponse{
			Code:    errInfos.NOT_FOUND,
			Message: "Rule not found",
		})
		return
	}

	result, err := ctl.recurrenceService.DeleteRecurringSchedule(ctx, rule.CenterID, teacherID, req.RuleID, editDate, services.RecurrenceEditMode(req.Mode), req.Reason)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, global.ApiResponse{
			Code:    errInfos.SYSTEM_ERROR,
			Message: err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, global.ApiResponse{
		Code:    0,
		Message: "Schedule deleted successfully",
		Datas:   result,
	})
}
