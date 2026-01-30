package controllers

import (
	"strings"
	"time"
	"timeLedger/app"
	"timeLedger/app/models"
	"timeLedger/app/repositories"
	"timeLedger/app/requests"
	"timeLedger/app/services"

	"github.com/gin-gonic/gin"
)

type SmartMatchingController struct {
	BaseController
	app              *app.App
	smartMatchingSvc services.SmartMatchingService
	auditLogRepo     *repositories.AuditLogRepository
}

type FindMatchesRequest struct {
	TeacherID         *uint     `json:"teacher_id"`
	RoomID            uint      `json:"room_id" binding:"required"`
	StartTime         time.Time `json:"start_time" binding:"required"`
	EndTime           time.Time `json:"end_time" binding:"required"`
	RequiredSkills    []string  `json:"required_skills"`
	ExcludeTeacherIDs []uint    `json:"exclude_teacher_ids"`
}

type TalentSearchParams struct {
	City     string   `form:"city"`
	District string   `form:"district"`
	Keyword  string   `form:"keyword"`
	Skills   []string `form:"skills"`
	Hashtags []string `form:"hashtags"`
}

func NewSmartMatchingController(app *app.App) *SmartMatchingController {
	return &SmartMatchingController{
		app:              app,
		smartMatchingSvc: services.NewSmartMatchingService(app),
		auditLogRepo:     repositories.NewAuditLogRepository(app),
	}
}

// requireCenterID 取得並驗證中心 ID（通用模式）
func (ctl *SmartMatchingController) requireCenterID(helper *ContextHelper) uint {
	centerID := helper.MustCenterID()
	if centerID == 0 {
		return 0
	}
	return centerID
}

// FindMatches 智慧媒合搜尋
// @Summary 智慧媒合搜尋
// @Description 根據條件搜尋可配合的老師
// @Tags Smart Matching
// @Accept json
// @Produce json
// @Param request body FindMatchesRequest true "搜尋條件"
// @Success 200 {object} global.ApiResponse
// @Router /admin/smart-matching/matches [post]
func (ctl *SmartMatchingController) FindMatches(ctx *gin.Context) {
	helper := NewContextHelper(ctx)

	var req FindMatchesRequest
	if !helper.MustBindJSON(&req) {
		return
	}

	centerID := ctl.requireCenterID(helper)
	if centerID == 0 {
		return
	}

	// 直接同步處理（Go 處理速度很快，不需要 WebSocket）
	matches, err := ctl.smartMatchingSvc.FindMatches(ctx.Request.Context(), centerID, req.TeacherID, req.RoomID, req.StartTime, req.EndTime, req.RequiredSkills, req.ExcludeTeacherIDs)
	if err != nil {
		helper.InternalError(err.Error())
		return
	}

	actorID := helper.MustUserID()
	ctl.auditLogRepo.Create(ctx, models.AuditLog{
		CenterID:   centerID,
		ActorType:  "ADMIN",
		ActorID:    actorID,
		Action:     "SMART_MATCH_FIND",
		TargetType: "SmartMatching",
		TargetID:   0,
		Payload: models.AuditPayload{
			After: map[string]interface{}{
				"room_id":         req.RoomID,
				"start_time":      req.StartTime,
				"end_time":        req.EndTime,
				"required_skills": req.RequiredSkills,
				"matches_count":   len(matches),
			},
		},
	})

	helper.Success(matches)
}

// SearchTalent 人才庫搜尋
// @Summary 人才庫搜尋
// @Description 搜尋符合條件的人才
// @Tags Smart Matching
// @Accept json
// @Produce json
// @Param city query string false "縣市"
// @Param district query string false "區域"
// @Param keyword query string false "關鍵字"
// @Param skills query string false "技能（逗號分隔）"
// @Param hashtags query string false "標籤（逗號分隔）"
// @Param page query int false "頁碼"
// @Param limit query int false "每頁筆數"
// @Success 200 {object} global.ApiResponse
// @Router /admin/smart-matching/talent/search [get]
func (ctl *SmartMatchingController) SearchTalent(ctx *gin.Context) {
	helper := NewContextHelper(ctx)
	centerID := ctl.requireCenterID(helper)
	if centerID == 0 {
		return
	}

	uriReq, err := requests.ValidateTalentSearch(ctx)
	if err != nil {
		helper.BadRequest(err.Error())
		return
	}

	skills := []string{}
	if uriReq.Skills != "" {
		skills = strings.Split(uriReq.Skills, ",")
		for i := range skills {
			skills[i] = strings.TrimSpace(skills[i])
		}
	}

	hashtags := []string{}
	if uriReq.Hashtags != "" {
		hashtags = strings.Split(uriReq.Hashtags, ",")
		for i := range hashtags {
			hashtags[i] = strings.TrimSpace(hashtags[i])
		}
	}

	searchParams := services.TalentSearchParams{
		CenterID:   centerID,
		City:       uriReq.City,
		District:   uriReq.District,
		Keyword:    uriReq.Keyword,
		Skills:     skills,
		Hashtags:   hashtags,
		Page:       uriReq.Page,
		Limit:      uriReq.Limit,
		SortBy:     uriReq.SortBy,
		SortOrder:  uriReq.SortOrder,
	}

	result, err := ctl.smartMatchingSvc.SearchTalent(ctx.Request.Context(), searchParams)
	if err != nil {
		helper.InternalError(err.Error())
		return
	}

	actorID := helper.MustUserID()
	ctl.auditLogRepo.Create(ctx, models.AuditLog{
		CenterID:   centerID,
		ActorType:  "ADMIN",
		ActorID:    actorID,
		Action:     "TALENT_SEARCH",
		TargetType: "Teacher",
		TargetID:   0,
		Payload: models.AuditPayload{
			After: map[string]interface{}{
				"city":          uriReq.City,
				"district":      uriReq.District,
				"keyword":       uriReq.Keyword,
				"skills":        uriReq.Skills,
				"hashtags":      uriReq.Hashtags,
				"page":          uriReq.Page,
				"limit":         uriReq.Limit,
				"sort_by":       uriReq.SortBy,
				"sort_order":    uriReq.SortOrder,
				"results_count": len(result.Talents),
				"total_count":   result.Pagination.Total,
			},
		},
	})

	helper.Success(result)
}

// TalentStatsRequest - 人才庫統計請求
type TalentStatsRequest struct {
	City     string `form:"city"`
	District string `form:"district"`
}

// TalentStatsResponse - 人才庫統計回應
type TalentStatsResponse struct {
	TotalCount       int                    `json:"total_count"`
	OpenHiringCount  int                    `json:"open_hiring_count"`
	MemberCount      int                    `json:"member_count"`
	AverageRating    float64                `json:"average_rating"`
	MonthlyChange    int                    `json:"monthly_change"`
	MonthlyTrend     []int                  `json:"monthly_trend"`
	PendingInvites   int                    `json:"pending_invites"`
	AcceptedInvites  int                    `json:"accepted_invites"`
	DeclinedInvites  int                    `json:"declined_invites"`
	CityDistribution []CityDistributionItem `json:"city_distribution"`
	TopSkills        []SkillCountItem       `json:"top_skills"`
}

type CityDistributionItem struct {
	Name  string `json:"name"`
	Count int    `json:"count"`
}

type SkillCountItem struct {
	Name  string `json:"name"`
	Count int    `json:"count"`
}

// GetTalentStats - 取得人才庫統計資料
// @Summary 取得人才庫統計資料
// @Description 取得人才庫的統計資訊
// @Tags Smart Matching
// @Accept json
// @Produce json
// @Param city query string false "縣市"
// @Param district query string false "區域"
// @Success 200 {object} TalentStatsResponse
// @Router /admin/smart-matching/talent/stats [get]
func (ctl *SmartMatchingController) GetTalentStats(ctx *gin.Context) {
	helper := NewContextHelper(ctx)
	centerID := ctl.requireCenterID(helper)
	if centerID == 0 {
		return
	}

	stats, err := ctl.smartMatchingSvc.GetTalentStats(ctx.Request.Context(), centerID)
	if err != nil {
		helper.InternalError(err.Error())
		return
	}

	helper.Success(stats)
}

// InviteTalentRequest - 邀請人才請求
type InviteTalentRequest struct {
	TeacherIDs []uint `json:"teacher_ids" binding:"required,min=1,max=20"`
	Message    string `json:"message"`
}

// InviteTalentResponse - 邀請人才回應
type InviteTalentResponse struct {
	InvitedCount  int      `json:"invited_count"`
	FailedCount   int      `json:"failed_count"`
	FailedIDs     []uint   `json:"failed_ids,omitempty"`
	InvitationIDs []uint   `json:"invitation_ids,omitempty"`
}

// InviteTalent - 邀請人才合作
// @Summary 邀請人才合作
// @Description 邀請老師加入人才庫
// @Tags Smart Matching
// @Accept json
// @Produce json
// @Param request body InviteTalentRequest true "邀請資訊"
// @Success 200 {object} InviteTalentResponse
// @Router /admin/smart-matching/talent/invite [post]
func (ctl *SmartMatchingController) InviteTalent(ctx *gin.Context) {
	helper := NewContextHelper(ctx)
	centerID := ctl.requireCenterID(helper)
	if centerID == 0 {
		return
	}

	var req InviteTalentRequest
	if !helper.MustBindJSON(&req) {
		return
	}

	actorID := helper.MustUserID()

	result, err := ctl.smartMatchingSvc.InviteTalent(ctx.Request.Context(), centerID, actorID, req.TeacherIDs, req.Message)
	if err != nil {
		helper.InternalError(err.Error())
		return
	}

	// 記錄審核日誌
	ctl.auditLogRepo.Create(ctx, models.AuditLog{
		CenterID:   centerID,
		ActorType:  "ADMIN",
		ActorID:    actorID,
		Action:     "TALENT_INVITE",
		TargetType: "Teacher",
		TargetID:   0,
		Payload: models.AuditPayload{
			After: map[string]interface{}{
				"teacher_ids":    req.TeacherIDs,
				"invited_count":  result.InvitedCount,
				"failed_count":   result.FailedCount,
				"invitation_ids": result.InvitationIDs,
			},
		},
	})

	helper.Success(result)
}

// GetSearchSuggestions - 取得搜尋建議
// @Summary 取得搜尋建議
// @Description 取得搜尋建議關鍵字
// @Tags Smart Matching
// @Accept json
// @Produce json
// @Param q query string false "搜尋關鍵字"
// @Success 200 {object} global.ApiResponse
// @Router /admin/smart-matching/suggestions [get]
func (ctl *SmartMatchingController) GetSearchSuggestions(ctx *gin.Context) {
	helper := NewContextHelper(ctx)

	query := helper.QueryStringOrDefault("q", "")
	if query == "" {
		query = helper.QueryStringOrDefault("keyword", "")
	}

	suggestions, err := ctl.smartMatchingSvc.GetSearchSuggestions(ctx.Request.Context(), query)
	if err != nil {
		helper.InternalError(err.Error())
		return
	}

	helper.Success(suggestions)
}

// GetAlternativeSlotsRequest - 替代時段請求
type GetAlternativeSlotsRequest struct {
	TeacherID    uint      `json:"teacher_id" binding:"required"`
	OriginalStart time.Time `json:"original_start" binding:"required"`
	OriginalEnd   time.Time `json:"original_end" binding:"required"`
	Duration     int       `json:"duration"`
}

// GetAlternativeSlots - 取得替代時段建議
// @Summary 取得替代時段建議
// @Description 取得替代時段建議
// @Tags Smart Matching
// @Accept json
// @Produce json
// @Param request body GetAlternativeSlotsRequest true "時段資訊"
// @Success 200 {object} global.ApiResponse
// @Router /admin/smart-matching/alternatives [post]
func (ctl *SmartMatchingController) GetAlternativeSlots(ctx *gin.Context) {
	helper := NewContextHelper(ctx)
	centerID := ctl.requireCenterID(helper)
	if centerID == 0 {
		return
	}

	var req GetAlternativeSlotsRequest
	if !helper.MustBindJSON(&req) {
		return
	}

	alternatives, err := ctl.smartMatchingSvc.GetAlternativeSlots(ctx.Request.Context(), centerID, req.TeacherID, req.OriginalStart, req.OriginalEnd, req.Duration)
	if err != nil {
		helper.InternalError(err.Error())
		return
	}

	helper.Success(alternatives)
}

// GetTeacherSessionsRequest - 教師課表請求
type GetTeacherSessionsRequest struct {
	StartDate string `form:"start_date" binding:"required"`
	EndDate   string `form:"end_date" binding:"required"`
}

// TeacherSessionResponse - 教師課表回應
type TeacherSessionResponse struct {
	TeacherID   uint         `json:"teacher_id"`
	TeacherName string       `json:"teacher_name"`
	Sessions    []SessionItem `json:"sessions"`
}

type SessionItem struct {
	ID         uint   `json:"id"`
	CourseName string `json:"course_name"`
	StartTime  string `json:"start_time"`
	EndTime    string `json:"end_time"`
	RoomName   string `json:"room_name,omitempty"`
	Status     string `json:"status"`
}

// GetTeacherSessions - 取得教師課表
// @Summary 取得教師課表
// @Description 取得指定教師的課表
// @Tags Smart Matching
// @Accept json
// @Produce json
// @Param teacher_id path uint true "教師ID"
// @Param start_date query string true "開始日期"
// @Param end_date query string true "結束日期"
// @Success 200 {object} TeacherSessionResponse
// @Router /admin/smart-matching/teachers/{teacher_id}/sessions [get]
func (ctl *SmartMatchingController) GetTeacherSessions(ctx *gin.Context) {
	helper := NewContextHelper(ctx)
	centerID := ctl.requireCenterID(helper)
	if centerID == 0 {
		return
	}

	teacherID := helper.MustParamUint("teacher_id")
	if teacherID == 0 {
		return
	}

	var req GetTeacherSessionsRequest
	if !helper.MustBindJSON(&req) {
		return
	}

	sessions, err := ctl.smartMatchingSvc.GetTeacherSessions(ctx.Request.Context(), centerID, teacherID, req.StartDate, req.EndDate)
	if err != nil {
		helper.InternalError(err.Error())
		return
	}

	helper.Success(sessions)
}
