package controllers

import (
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
		CenterID: centerID,
		City:     uriReq.City,
		District: uriReq.District,
		Keyword:  uriReq.Keyword,
		Skills:   skills,
		Hashtags: hashtags,
	}

	results, err := ctl.smartMatchingSvc.SearchTalent(ctx, searchParams)
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
				"results_count": len(results),
			},
		},
	})

	ctx.JSON(http.StatusOK, global.ApiResponse{
		Code:    0,
		Message: "OK",
		Datas:   results,
	})
}
