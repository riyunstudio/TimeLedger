package controllers

import (
	"fmt"
	"net/http"
	"time"
	"timeLedger/app"
	"timeLedger/app/models"
	"timeLedger/app/repositories"
	"timeLedger/global"
	"timeLedger/global/errInfos"

	"github.com/gin-gonic/gin"
)

type OfferingController struct {
	BaseController
	offeringRepository *repositories.OfferingRepository
	courseRepository   *repositories.CourseRepository
	roomRepository     *repositories.RoomRepository
	teacherRepository  *repositories.TeacherRepository
	auditLogRepo       *repositories.AuditLogRepository
}

func NewOfferingController(app *app.App) *OfferingController {
	return &OfferingController{
		offeringRepository: repositories.NewOfferingRepository(app),
		courseRepository:   repositories.NewCourseRepository(app),
		roomRepository:     repositories.NewRoomRepository(app),
		teacherRepository:  repositories.NewTeacherRepository(app),
		auditLogRepo:       repositories.NewAuditLogRepository(app),
	}
}

type GetOfferingsResponse struct {
	Offerings  []models.Offering `json:"offerings"`
	Pagination global.Pagination `json:"pagination"`
}

// GetOfferings 取得班別列表
// @Summary 取得班別列表
// @Tags Admin
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Items per page" default(20)
// @Success 200 {object} global.ApiResponse{data=GetOfferingsResponse}
// @Router /api/v1/admin/offerings [get]
func (ctl *OfferingController) GetOfferings(ctx *gin.Context) {
	centerID := ctx.GetUint(global.CenterIDKey)
	if centerID == 0 {
		ctx.JSON(http.StatusBadRequest, global.ApiResponse{
			Code:    global.BAD_REQUEST,
			Message: "Center ID required",
		})
		return
	}

	page := int64(1)
	limit := int64(20)
	if p := ctx.Query("page"); p != "" {
		fmt.Sscanf(p, "%d", &page)
	}
	if l := ctx.Query("limit"); l != "" {
		fmt.Sscanf(l, "%d", &limit)
	}

	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 20
	}

	offerings, total, err := ctl.offeringRepository.ListByCenterIDPaginated(ctl.makeCtx(ctx), centerID, page, limit)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, global.ApiResponse{
			Code:    errInfos.SQL_ERROR,
			Message: ctl.app.Err.New(errInfos.SQL_ERROR).Msg,
		})
		return
	}

	pagination := global.NewPagination(page, limit, total)

	ctx.JSON(http.StatusOK, global.ApiResponse{
		Code:    0,
		Message: "Success",
		Datas: GetOfferingsResponse{
			Offerings:  offerings,
			Pagination: pagination,
		},
	})
}

// CreateOffering 新增班別
// @Summary 新增班別
// @Tags Admin
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body CreateOfferingRequest true "班別資訊"
// @Success 200 {object} global.ApiResponse{data=models.Offering}
// @Router /api/v1/admin/offerings [post]
func (ctl *OfferingController) CreateOffering(ctx *gin.Context) {
	centerID := ctx.GetUint(global.CenterIDKey)
	if centerID == 0 {
		ctx.JSON(http.StatusBadRequest, global.ApiResponse{
			Code:    global.BAD_REQUEST,
			Message: "Center ID required",
		})
		return
	}

	var req CreateOfferingRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, global.ApiResponse{
			Code:    global.BAD_REQUEST,
			Message: "Invalid request body",
		})
		return
	}

	offering := models.Offering{
		CenterID:            centerID,
		CourseID:            req.CourseID,
		DefaultRoomID:       req.DefaultRoomID,
		DefaultTeacherID:    req.DefaultTeacherID,
		AllowBufferOverride: req.AllowBufferOverride,
		CreatedAt:           time.Now(),
		UpdatedAt:           time.Now(),
	}

	createdOffering, err := ctl.offeringRepository.Create(ctl.makeCtx(ctx), offering)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, global.ApiResponse{
			Code:    500,
			Message: "Failed to create offering",
		})
		return
	}

	actorID := ctx.GetUint(global.UserIDKey)
	ctl.auditLogRepo.Create(ctx, models.AuditLog{
		CenterID:   centerID,
		ActorType:  "ADMIN",
		ActorID:    actorID,
		Action:     "OFFERING_CREATE",
		TargetType: "Offering",
		TargetID:   createdOffering.ID,
		Payload: models.AuditPayload{
			After: map[string]interface{}{
				"course_id":             req.CourseID,
				"default_room_id":       req.DefaultRoomID,
				"default_teacher_id":    req.DefaultTeacherID,
				"allow_buffer_override": req.AllowBufferOverride,
			},
		},
	})

	ctx.JSON(http.StatusOK, global.ApiResponse{
		Code:    0,
		Message: "Offering created",
		Datas:   createdOffering,
	})
}

type CreateOfferingRequest struct {
	CourseID            uint  `json:"course_id" binding:"required"`
	DefaultRoomID       *uint `json:"default_room_id"`
	DefaultTeacherID    *uint `json:"default_teacher_id"`
	AllowBufferOverride bool  `json:"allow_buffer_override"`
}

// UpdateOffering 更新班別
// @Summary 更新班別
// @Tags Admin
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param offering_id path int true "Offering ID"
// @Param request body UpdateOfferingRequest true "班別資訊"
// @Success 200 {object} global.ApiResponse{data=models.Offering}
// @Router /api/v1/admin/offerings/{offering_id} [put]
func (ctl *OfferingController) UpdateOffering(ctx *gin.Context) {
	centerID := ctx.GetUint(global.CenterIDKey)
	offeringIDStr := ctx.Param("offering_id")
	if offeringIDStr == "" {
		ctx.JSON(http.StatusBadRequest, global.ApiResponse{
			Code:    global.BAD_REQUEST,
			Message: "Offering ID required",
		})
		return
	}

	var offeringID uint
	if _, err := fmt.Sscanf(offeringIDStr, "%d", &offeringID); err != nil {
		ctx.JSON(http.StatusBadRequest, global.ApiResponse{
			Code:    global.BAD_REQUEST,
			Message: "Invalid offering ID",
		})
		return
	}

	var req UpdateOfferingRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, global.ApiResponse{
			Code:    global.BAD_REQUEST,
			Message: "Invalid request body",
		})
		return
	}

	offering := models.Offering{
		CenterID:            centerID,
		DefaultRoomID:       req.DefaultRoomID,
		DefaultTeacherID:    req.DefaultTeacherID,
		AllowBufferOverride: req.AllowBufferOverride,
		UpdatedAt:           time.Now(),
	}

	if err := ctl.offeringRepository.Update(ctl.makeCtx(ctx), offering); err != nil {
		ctx.JSON(http.StatusInternalServerError, global.ApiResponse{
			Code:    500,
			Message: "Failed to update offering",
		})
		return
	}

	actorID := ctx.GetUint(global.UserIDKey)
	ctl.auditLogRepo.Create(ctx, models.AuditLog{
		CenterID:   centerID,
		ActorType:  "ADMIN",
		ActorID:    actorID,
		Action:     "OFFERING_UPDATE",
		TargetType: "Offering",
		TargetID:   offeringID,
		Payload: models.AuditPayload{
			After: map[string]interface{}{
				"default_room_id":       req.DefaultRoomID,
				"default_teacher_id":    req.DefaultTeacherID,
				"allow_buffer_override": req.AllowBufferOverride,
			},
		},
	})

	ctx.JSON(http.StatusOK, global.ApiResponse{
		Code:    0,
		Message: "Offering updated",
		Datas:   offering,
	})
}

type UpdateOfferingRequest struct {
	DefaultRoomID       *uint `json:"default_room_id"`
	DefaultTeacherID    *uint `json:"default_teacher_id"`
	AllowBufferOverride bool  `json:"allow_buffer_override"`
}

type CopyOfferingRequest struct {
	NewName     string `json:"new_name" binding:"required"`
	CopyTeacher bool   `json:"copy_teacher"`
}

type CopyOfferingResponse struct {
	ID          uint   `json:"id"`
	Name        string `json:"name"`
	CourseID    uint   `json:"course_id"`
	RulesCopied int    `json:"rules_copied"`
}

// CopyOffering 複製班別

// DeleteOffering 刪除班別
// @Summary 刪除班別
// @Tags Admin
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param offering_id path int true "Offering ID"
// @Success 200 {object} global.ApiResponse
// @Router /api/v1/admin/offerings/{offering_id} [delete]
func (ctl *OfferingController) DeleteOffering(ctx *gin.Context) {
	centerID := ctx.GetUint(global.CenterIDKey)
	offeringIDStr := ctx.Param("offering_id")
	if offeringIDStr == "" {
		ctx.JSON(http.StatusBadRequest, global.ApiResponse{
			Code:    global.BAD_REQUEST,
			Message: "Offering ID required",
		})
		return
	}

	var offeringID uint
	if _, err := fmt.Sscanf(offeringIDStr, "%d", &offeringID); err != nil {
		ctx.JSON(http.StatusBadRequest, global.ApiResponse{
			Code:    global.BAD_REQUEST,
			Message: "Invalid offering ID",
		})
		return
	}

	if err := ctl.offeringRepository.DeleteByID(ctl.makeCtx(ctx), offeringID, centerID); err != nil {
		ctx.JSON(http.StatusInternalServerError, global.ApiResponse{
			Code:    500,
			Message: "Failed to delete offering",
		})
		return
	}

	actorID := ctx.GetUint(global.UserIDKey)
	ctl.auditLogRepo.Create(ctx, models.AuditLog{
		CenterID:   centerID,
		ActorType:  "ADMIN",
		ActorID:    actorID,
		Action:     "OFFERING_DELETE",
		TargetType: "Offering",
		TargetID:   offeringID,
		Payload: models.AuditPayload{
			After: map[string]interface{}{
				"status": "DELETED",
			},
		},
	})

	ctx.JSON(http.StatusOK, global.ApiResponse{
		Code:    0,
		Message: "Offering deleted",
	})
}

// CopyOffering 複製班別
// @Summary 複製班別
// @Tags Admin
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Center ID"
// @Param offering_id path int true "Offering ID"
// @Param request body CopyOfferingRequest true "複製資訊"
// @Success 200 {object} global.ApiResponse{data=CopyOfferingResponse}
// @Router /api/v1/admin/centers/{id}/offerings/{offering_id}/copy [post]
func (ctl *OfferingController) CopyOffering(ctx *gin.Context) {
	centerIDStr := ctx.Param("id")
	offeringIDStr := ctx.Param("offeringId")

	var centerID, offeringID uint
	if _, err := fmt.Sscanf(centerIDStr, "%d", &centerID); err != nil {
		ctx.JSON(http.StatusBadRequest, global.ApiResponse{
			Code:    global.BAD_REQUEST,
			Message: "Invalid center ID",
		})
		return
	}
	if _, err := fmt.Sscanf(offeringIDStr, "%d", &offeringID); err != nil {
		ctx.JSON(http.StatusBadRequest, global.ApiResponse{
			Code:    global.BAD_REQUEST,
			Message: "Invalid offering ID",
		})
		return
	}

	var req CopyOfferingRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, global.ApiResponse{
			Code:    global.BAD_REQUEST,
			Message: "Invalid request body",
		})
		return
	}

	original, err := ctl.offeringRepository.GetByIDAndCenterID(ctl.makeCtx(ctx), offeringID, centerID)
	if err != nil {
		ctx.JSON(http.StatusNotFound, global.ApiResponse{
			Code:    errInfos.NOT_FOUND,
			Message: ctl.app.Err.New(errInfos.NOT_FOUND).Msg,
		})
		return
	}

	newTeacherID := original.DefaultTeacherID
	if !req.CopyTeacher {
		newTeacherID = nil
	}

	newOffering := models.Offering{
		CenterID:            centerID,
		CourseID:            original.CourseID,
		Name:                req.NewName,
		DefaultRoomID:       original.DefaultRoomID,
		DefaultTeacherID:    newTeacherID,
		AllowBufferOverride: original.AllowBufferOverride,
		IsActive:            true,
		CreatedAt:           time.Now(),
		UpdatedAt:           time.Now(),
	}

	createdOffering, err := ctl.offeringRepository.Create(ctl.makeCtx(ctx), newOffering)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, global.ApiResponse{
			Code:    errInfos.SQL_ERROR,
			Message: ctl.app.Err.New(errInfos.SQL_ERROR).Msg,
		})
		return
	}

	actorID := ctx.GetUint(global.UserIDKey)
	ctl.auditLogRepo.Create(ctx, models.AuditLog{
		CenterID:   centerID,
		ActorType:  "ADMIN",
		ActorID:    actorID,
		Action:     "OFFERING_COPY",
		TargetType: "Offering",
		TargetID:   createdOffering.ID,
		Payload: models.AuditPayload{
			After: map[string]interface{}{
				"original_offering_id": original.ID,
				"new_offering_id":      createdOffering.ID,
				"name":                 req.NewName,
				"copy_teacher":         req.CopyTeacher,
			},
		},
	})

	ctx.JSON(http.StatusOK, global.ApiResponse{
		Code:    0,
		Message: "Offering copied successfully",
		Datas: CopyOfferingResponse{
			ID:          createdOffering.ID,
			Name:        createdOffering.Name,
			CourseID:    createdOffering.CourseID,
			RulesCopied: 0,
		},
	})
}
