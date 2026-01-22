package controllers

import (
	"fmt"
	"net/http"
	"time"
	"timeLedger/app"
	"timeLedger/app/models"
	"timeLedger/app/repositories"
	"timeLedger/global"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

type AdminUserController struct {
	BaseController
	adminRepository *repositories.AdminUserRepository
	auditLogRepo    *repositories.AuditLogRepository
}

func NewAdminUserController(app *app.App) *AdminUserController {
	return &AdminUserController{
		adminRepository: repositories.NewAdminUserRepository(app),
		auditLogRepo:    repositories.NewAuditLogRepository(app),
	}
}

// GetAdminUsers 取得管理員列表
// @Summary 取得管理員列表
// @Tags Admin
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Center ID"
// @Success 200 {object} global.ApiResponse{data=[]models.AdminUser}
// @Router /api/v1/admin/centers/{id}/users [get]
func (ctl *AdminUserController) GetAdminUsers(ctx *gin.Context) {
	centerID := ctx.Param("id")
	if centerID == "" {
		ctx.JSON(http.StatusBadRequest, global.ApiResponse{
			Code:    global.BAD_REQUEST,
			Message: "Center ID required",
		})
		return
	}

	admins, err := ctl.adminRepository.ListByCenterID(ctl.makeCtx(ctx), 0)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, global.ApiResponse{
			Code:    500,
			Message: "Failed to get admin users",
		})
		return
	}

	ctx.JSON(http.StatusOK, global.ApiResponse{
		Code:    0,
		Message: "Success",
		Datas:   admins,
	})
}

// CreateAdminUser 新增管理員
// @Summary 新增管理員
// @Tags Admin
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Center ID"
// @Param request body CreateAdminUserRequest true "管理員資訊"
// @Success 200 {object} global.ApiResponse{data=models.AdminUser}
// @Router /api/v1/admin/centers/{id}/users [post]
func (ctl *AdminUserController) CreateAdminUser(ctx *gin.Context) {
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

	var req CreateAdminUserRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, global.ApiResponse{
			Code:    global.BAD_REQUEST,
			Message: "Invalid request body",
		})
		return
	}

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)

	admin := models.AdminUser{
		Email:        req.Email,
		PasswordHash: string(hashedPassword),
		Name:         req.Name,
		Role:         req.Role,
		Status:       "ACTIVE",
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	createdAdmin, err := ctl.adminRepository.Create(ctl.makeCtx(ctx), admin)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, global.ApiResponse{
			Code:    500,
			Message: "Failed to create admin user",
		})
		return
	}

	actorID := ctx.GetUint(global.UserIDKey)
	ctl.auditLogRepo.Create(ctx, models.AuditLog{
		CenterID:   centerID,
		ActorType:  "ADMIN",
		ActorID:    actorID,
		Action:     "ADMIN_USER_CREATE",
		TargetType: "AdminUser",
		TargetID:   createdAdmin.ID,
		Payload: models.AuditPayload{
			After: map[string]interface{}{
				"email": req.Email,
				"name":  req.Name,
				"role":  req.Role,
			},
		},
	})

	ctx.JSON(http.StatusOK, global.ApiResponse{
		Code:    0,
		Message: "Admin user created",
		Datas:   createdAdmin,
	})
}

type CreateAdminUserRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
	Name     string `json:"name" binding:"required"`
	Role     string `json:"role" binding:"required,oneof=OWNER ADMIN STAFF"`
}

// UpdateAdminUser 更新管理員
// @Summary 更新管理員
// @Tags Admin
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Center ID"
// @Param admin_id path int true "Admin ID"
// @Param request body UpdateAdminUserRequest true "管理員資訊"
// @Success 200 {object} global.ApiResponse{data=models.AdminUser}
// @Router /api/v1/admin/centers/{id}/users/{admin_id} [put]
func (ctl *AdminUserController) UpdateAdminUser(ctx *gin.Context) {
	centerIDStr := ctx.Param("id")
	adminIDStr := ctx.Param("admin_id")
	if adminIDStr == "" {
		ctx.JSON(http.StatusBadRequest, global.ApiResponse{
			Code:    global.BAD_REQUEST,
			Message: "Admin ID required",
		})
		return
	}

	var centerID, adminID uint
	if _, err := fmt.Sscanf(centerIDStr, "%d", &centerID); err != nil {
		ctx.JSON(http.StatusBadRequest, global.ApiResponse{
			Code:    global.BAD_REQUEST,
			Message: "Invalid center ID",
		})
		return
	}
	if _, err := fmt.Sscanf(adminIDStr, "%d", &adminID); err != nil {
		ctx.JSON(http.StatusBadRequest, global.ApiResponse{
			Code:    global.BAD_REQUEST,
			Message: "Invalid admin ID",
		})
		return
	}

	var req UpdateAdminUserRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, global.ApiResponse{
			Code:    global.BAD_REQUEST,
			Message: "Invalid request body",
		})
		return
	}

	admin := models.AdminUser{
		Name:      req.Name,
		Role:      req.Role,
		Status:    req.Status,
		UpdatedAt: time.Now(),
	}

	if req.Password != "" {
		hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
		admin.PasswordHash = string(hashedPassword)
	}

	if err := ctl.adminRepository.Update(ctl.makeCtx(ctx), admin); err != nil {
		ctx.JSON(http.StatusInternalServerError, global.ApiResponse{
			Code:    500,
			Message: "Failed to update admin user",
		})
		return
	}

	actorID := ctx.GetUint(global.UserIDKey)
	ctl.auditLogRepo.Create(ctx, models.AuditLog{
		CenterID:   centerID,
		ActorType:  "ADMIN",
		ActorID:    actorID,
		Action:     "ADMIN_USER_UPDATE",
		TargetType: "AdminUser",
		TargetID:   adminID,
		Payload: models.AuditPayload{
			After: map[string]interface{}{
				"name":   req.Name,
				"role":   req.Role,
				"status": req.Status,
			},
		},
	})

	ctx.JSON(http.StatusOK, global.ApiResponse{
		Code:    0,
		Message: "Admin user updated",
		Datas:   admin,
	})
}

type UpdateAdminUserRequest struct {
	Name     string `json:"name"`
	Password string `json:"password"`
	Role     string `json:"role" binding:"omitempty,oneof=OWNER ADMIN STAFF"`
	Status   string `json:"status" binding:"omitempty,oneof=ACTIVE INACTIVE"`
}

// DeleteAdminUser 刪除管理員
// @Summary 刪除管理員
// @Tags Admin
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Center ID"
// @Param admin_id path int true "Admin ID"
// @Success 200 {object} global.ApiResponse
// @Router /api/v1/admin/centers/{id}/users/{admin_id} [delete]
func (ctl *AdminUserController) DeleteAdminUser(ctx *gin.Context) {
	centerIDStr := ctx.Param("id")
	adminIDStr := ctx.Param("admin_id")
	if adminIDStr == "" {
		ctx.JSON(http.StatusBadRequest, global.ApiResponse{
			Code:    global.BAD_REQUEST,
			Message: "Admin ID required",
		})
		return
	}

	var centerID, adminID uint
	if _, err := fmt.Sscanf(centerIDStr, "%d", &centerID); err != nil {
		ctx.JSON(http.StatusBadRequest, global.ApiResponse{
			Code:    global.BAD_REQUEST,
			Message: "Invalid center ID",
		})
		return
	}
	if _, err := fmt.Sscanf(adminIDStr, "%d", &adminID); err != nil {
		ctx.JSON(http.StatusBadRequest, global.ApiResponse{
			Code:    global.BAD_REQUEST,
			Message: "Invalid admin ID",
		})
		return
	}

	if err := ctl.adminRepository.Delete(ctl.makeCtx(ctx), adminID); err != nil {
		ctx.JSON(http.StatusInternalServerError, global.ApiResponse{
			Code:    500,
			Message: "Failed to delete admin user",
		})
		return
	}

	actorID := ctx.GetUint(global.UserIDKey)
	ctl.auditLogRepo.Create(ctx, models.AuditLog{
		CenterID:   centerID,
		ActorType:  "ADMIN",
		ActorID:    actorID,
		Action:     "ADMIN_USER_DELETE",
		TargetType: "AdminUser",
		TargetID:   adminID,
		Payload: models.AuditPayload{
			After: map[string]interface{}{
				"status": "DELETED",
			},
		},
	})

	ctx.JSON(http.StatusOK, global.ApiResponse{
		Code:    0,
		Message: "Admin user deleted",
	})
}
