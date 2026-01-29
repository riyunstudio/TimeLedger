package controllers

import (
	"fmt"
	"net/http"
	"strings"
	"time"
	"timeLedger/app"
	"timeLedger/app/models"
	"timeLedger/app/repositories"
	"timeLedger/app/requests"
	"timeLedger/app/services"
	"timeLedger/global"

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

func (ctl *SmartMatchingController) FindMatches(ctx *gin.Context) {
	var req FindMatchesRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, global.ApiResponse{
			Code:    global.BAD_REQUEST,
			Message: "Invalid request parameters",
		})
		return
	}

	// 從 JWT token 取得 center_id
	centerID := ctx.GetUint(global.CenterIDKey)
	if centerID == 0 {
		if val, exists := ctx.Get(global.CenterIDKey); exists {
			switch v := val.(type) {
			case uint:
				centerID = v
			case uint64:
				centerID = uint(v)
			case int:
				centerID = uint(v)
			case float64:
				centerID = uint(v)
			}
		}
	}

	if centerID == 0 {
		ctx.JSON(http.StatusUnauthorized, global.ApiResponse{
			Code:    global.UNAUTHORIZED,
			Message: "Center ID not found in token",
		})
		return
	}

	// 直接同步處理（Go 處理速度很快，不需要 WebSocket）
	matches, err := ctl.smartMatchingSvc.FindMatches(ctx, centerID, req.TeacherID, req.RoomID, req.StartTime, req.EndTime, req.RequiredSkills, req.ExcludeTeacherIDs)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, global.ApiResponse{
			Code:    500,
			Message: err.Error(),
		})
		return
	}

	actorID := ctx.GetUint(global.UserIDKey)
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

	ctx.JSON(http.StatusOK, global.ApiResponse{
		Code:    0,
		Message: "OK",
		Datas:   matches,
	})
}

func (ctl *SmartMatchingController) SearchTalent(ctx *gin.Context) {
	uriReq, err := requests.ValidateTalentSearch(ctx)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, global.ApiResponse{
			Code:    global.BAD_REQUEST,
			Message: err.Error(),
		})
		return
	}

	centerIDVal, exists := ctx.Get(global.CenterIDKey)
	if !exists {
		ctx.JSON(http.StatusBadRequest, global.ApiResponse{
			Code:    global.BAD_REQUEST,
			Message: "Center ID required",
		})
		return
	}

	centerID, ok := centerIDVal.(uint)
	if !ok {
		ctx.JSON(http.StatusBadRequest, global.ApiResponse{
			Code:    global.BAD_REQUEST,
			Message: "Invalid center ID format",
		})
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

	result, err := ctl.smartMatchingSvc.SearchTalent(ctx, searchParams)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, global.ApiResponse{
			Code:    500,
			Message: err.Error(),
		})
		return
	}

	actorID := ctx.GetUint(global.UserIDKey)
	ctl.auditLogRepo.Create(ctx, models.AuditLog{
		CenterID:   uriReq.CenterID,
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

	ctx.JSON(http.StatusOK, global.ApiResponse{
		Code:    0,
		Message: "OK",
		Datas:   result,
	})
}

// TalentStatsRequest - 人才庫統計請求
type TalentStatsRequest struct {
	City     string `form:"city"`
	District string `form:"district"`
}

// TalentStatsResponse - 人才庫統計回應
type TalentStatsResponse struct {
	TotalCount       int                      `json:"total_count"`
	OpenHiringCount  int                      `json:"open_hiring_count"`
	MemberCount      int                      `json:"member_count"`
	AverageRating    float64                  `json:"average_rating"`
	MonthlyChange    int                      `json:"monthly_change"`
	MonthlyTrend     []int                    `json:"monthly趋势"`
	PendingInvites   int                      `json:"pending_invites"`
	AcceptedInvites  int                      `json:"accepted_invites"`
	DeclinedInvites  int                      `json:"declined_invites"`
	CityDistribution []CityDistributionItem   `json:"city_distribution"`
	TopSkills        []SkillCountItem         `json:"top_skills"`
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
func (ctl *SmartMatchingController) GetTalentStats(ctx *gin.Context) {
	centerIDVal, exists := ctx.Get(global.CenterIDKey)
	if !exists {
		ctx.JSON(http.StatusBadRequest, global.ApiResponse{
			Code:    global.BAD_REQUEST,
			Message: "Center ID required",
		})
		return
	}

	centerID, ok := centerIDVal.(uint)
	if !ok {
		ctx.JSON(http.StatusBadRequest, global.ApiResponse{
			Code:    global.BAD_REQUEST,
			Message: "Invalid center ID format",
		})
		return
	}

	stats, err := ctl.smartMatchingSvc.GetTalentStats(ctx, centerID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, global.ApiResponse{
			Code:    500,
			Message: err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, global.ApiResponse{
		Code:    0,
		Message: "OK",
		Datas:   stats,
	})
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
func (ctl *SmartMatchingController) InviteTalent(ctx *gin.Context) {
	var req InviteTalentRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, global.ApiResponse{
			Code:    global.BAD_REQUEST,
			Message: "Invalid request: " + err.Error(),
		})
		return
	}

	centerIDVal, exists := ctx.Get(global.CenterIDKey)
	if !exists {
		ctx.JSON(http.StatusBadRequest, global.ApiResponse{
			Code:    global.BAD_REQUEST,
			Message: "Center ID required",
		})
		return
	}

	centerID, ok := centerIDVal.(uint)
	if !ok {
		ctx.JSON(http.StatusBadRequest, global.ApiResponse{
			Code:    global.BAD_REQUEST,
			Message: "Invalid center ID format",
		})
		return
	}

	actorID := ctx.GetUint(global.UserIDKey)

	result, err := ctl.smartMatchingSvc.InviteTalent(ctx, centerID, actorID, req.TeacherIDs, req.Message)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, global.ApiResponse{
			Code:    500,
			Message: err.Error(),
		})
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

	ctx.JSON(http.StatusOK, global.ApiResponse{
		Code:    0,
		Message: "OK",
		Datas:   result,
	})
}

// GetSearchSuggestions - 取得搜尋建議
func (ctl *SmartMatchingController) GetSearchSuggestions(ctx *gin.Context) {
	query := ctx.Query("q")
	if query == "" {
		query = ctx.Query("keyword")
	}

	suggestions, err := ctl.smartMatchingSvc.GetSearchSuggestions(ctx, query)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, global.ApiResponse{
			Code:    500,
			Message: err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, global.ApiResponse{
		Code:    0,
		Message: "OK",
		Datas:   suggestions,
	})
}

// GetAlternativeSlotsRequest - 替代時段請求
type GetAlternativeSlotsRequest struct {
	TeacherID  uint   `json:"teacher_id" binding:"required"`
	OriginalStart time.Time `json:"original_start" binding:"required"`
	OriginalEnd   time.Time `json:"original_end" binding:"required"`
	Duration     int    `json:"duration"`
}

// GetAlternativeSlots - 取得替代時段建議
func (ctl *SmartMatchingController) GetAlternativeSlots(ctx *gin.Context) {
	var req GetAlternativeSlotsRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, global.ApiResponse{
			Code:    global.BAD_REQUEST,
			Message: "Invalid request: " + err.Error(),
		})
		return
	}

	centerIDVal, exists := ctx.Get(global.CenterIDKey)
	if !exists {
		ctx.JSON(http.StatusBadRequest, global.ApiResponse{
			Code:    global.BAD_REQUEST,
			Message: "Center ID required",
		})
		return
	}

	centerID, ok := centerIDVal.(uint)
	if !ok {
		ctx.JSON(http.StatusBadRequest, global.ApiResponse{
			Code:    global.BAD_REQUEST,
			Message: "Invalid center ID format",
		})
		return
	}

	alternatives, err := ctl.smartMatchingSvc.GetAlternativeSlots(ctx, centerID, req.TeacherID, req.OriginalStart, req.OriginalEnd, req.Duration)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, global.ApiResponse{
			Code:    500,
			Message: err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, global.ApiResponse{
		Code:    0,
		Message: "OK",
		Datas:   alternatives,
	})
}

// GetTeacherSessionsRequest - 教師課表請求
type GetTeacherSessionsRequest struct {
	StartDate string `form:"start_date" binding:"required"`
	EndDate   string `form:"end_date" binding:"required"`
}

// TeacherSessionResponse - 教師課表回應
type TeacherSessionResponse struct {
	TeacherID uint              `json:"teacher_id"`
	TeacherName string          `json:"teacher_name"`
	Sessions   []SessionItem    `json:"sessions"`
}

type SessionItem struct {
	ID          uint   `json:"id"`
	CourseName  string `json:"course_name"`
	StartTime   string `json:"start_time"`
	EndTime     string `json:"end_time"`
	RoomName    string `json:"room_name,omitempty"`
	Status      string `json:"status"`
}

// GetTeacherSessions - 取得教師課表
func (ctl *SmartMatchingController) GetTeacherSessions(ctx *gin.Context) {
	teacherID := ctx.Param("teacher_id")
	if teacherID == "" {
		ctx.JSON(http.StatusBadRequest, global.ApiResponse{
			Code:    global.BAD_REQUEST,
			Message: "Teacher ID is required",
		})
		return
	}

	var req GetTeacherSessionsRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, global.ApiResponse{
			Code:    global.BAD_REQUEST,
			Message: "Invalid request: " + err.Error(),
		})
		return
	}

	centerIDVal, exists := ctx.Get(global.CenterIDKey)
	if !exists {
		ctx.JSON(http.StatusBadRequest, global.ApiResponse{
			Code:    global.BAD_REQUEST,
			Message: "Center ID required",
		})
		return
	}

	centerID, ok := centerIDVal.(uint)
	if !ok {
		ctx.JSON(http.StatusBadRequest, global.ApiResponse{
			Code:    global.BAD_REQUEST,
			Message: "Invalid center ID format",
		})
		return
	}

	var tID uint
	if _, err := fmt.Sscanf(teacherID, "%d", &tID); err != nil {
		ctx.JSON(http.StatusBadRequest, global.ApiResponse{
			Code:    global.BAD_REQUEST,
			Message: "Invalid teacher ID format",
		})
		return
	}

	sessions, err := ctl.smartMatchingSvc.GetTeacherSessions(ctx, centerID, tID, req.StartDate, req.EndDate)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, global.ApiResponse{
			Code:    500,
			Message: err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, global.ApiResponse{
		Code:    0,
		Message: "OK",
		Datas:   sessions,
	})
}
