package controllers

import (
	"context"
	"time"
	"timeLedger/app"
	"timeLedger/app/models"
	"timeLedger/app/repositories"
	"timeLedger/app/resources"
	"timeLedger/app/services"

	"github.com/gin-gonic/gin"
)

// TeacherInvitationController 教師邀請控制器
// 處理老師邀請、邀請連結管理、公開邀請、管理員邀請管理等功能
type TeacherInvitationController struct {
	BaseController
	app            *app.App
	teacherService *services.TeacherService
	invitationRepo *repositories.CenterInvitationRepository
	centerRepo     *repositories.CenterRepository
	invitationRes  *resources.InvitationResource
}

// NewTeacherInvitationController 建立教師邀請控制器
func NewTeacherInvitationController(app *app.App) *TeacherInvitationController {
	teacherService := services.NewTeacherService(app)

	return &TeacherInvitationController{
		app:            app,
		teacherService: teacherService,
		invitationRepo: repositories.NewCenterInvitationRepository(app),
		centerRepo:     repositories.NewCenterRepository(app),
		invitationRes:  resources.NewInvitationResource(),
	}
}

// ========== 管理員端邀請管理 API ==========

// InvitationStatsResponse 邀請統計回應結構
type InvitationStatsResponse struct {
	Total         int64 `json:"total"`
	Pending       int64 `json:"pending"`
	Accepted      int64 `json:"accepted"`
	Expired       int64 `json:"expired"`
	Rejected      int64 `json:"rejected"`
	RecentPending int64 `json:"recent_pending"`
}

// GetInvitationStats 取得邀請統計資料
// @Summary 取得邀請統計資料
// @Tags Admin - Invitations
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Center ID"
// @Success 200 {object} global.ApiResponse{data=InvitationStatsResponse}
// @Router /api/v1/admin/centers/{id}/invitations/stats [get]
func (ctl *TeacherInvitationController) GetInvitationStats(ctx *gin.Context) {
	helper := NewContextHelper(ctx)

	centerID := helper.MustParamUint("id")
	if centerID == 0 {
		return
	}

	now := time.Now()
	thirtyDaysAgo := now.AddDate(0, 0, -30)

	total, _ := ctl.invitationRepo.CountByCenterID(ctx, centerID)
	pending, _ := ctl.invitationRepo.CountByStatus(ctx, centerID, "PENDING")
	accepted, _ := ctl.invitationRepo.CountByStatus(ctx, centerID, "ACCEPTED")
	expired, _ := ctl.invitationRepo.CountByStatus(ctx, centerID, "EXPIRED")
	rejected, _ := ctl.invitationRepo.CountByStatus(ctx, centerID, "REJECTED")
	recentPending, _ := ctl.invitationRepo.CountByDateRange(ctx, centerID, thirtyDaysAgo, now)

	stats := InvitationStatsResponse{
		Total:         total,
		Pending:       pending,
		Accepted:      accepted,
		Expired:       expired,
		Rejected:      rejected,
		RecentPending: recentPending,
	}

	helper.Success(stats)
}

// PaginationRequest 分頁請求結構
type PaginationRequest struct {
	Page  int `form:"page"`
	Limit int `form:"limit"`
}

// PaginationResponse 分頁回應結構
type PaginationResponse struct {
	Data       interface{} `json:"data"`
	Total      int64       `json:"total"`
	Page       int         `json:"page"`
	Limit      int         `json:"limit"`
	TotalPages int64       `json:"total_pages"`
}

// GetInvitations 取得邀請列表
// @Summary 取得邀請列表
// @Tags Admin - Invitations
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Center ID"
// @Param status query string false "篩選狀態 (PENDING/ACCEPTED/DECLINED/EXPIRED)"
// @Param page query int false "頁碼"
// @Param limit query int false "每頁筆數"
// @Success 200 {object} global.ApiResponse{data=PaginationResponse}
// @Router /api/v1/admin/centers/{id}/invitations [get]
func (ctl *TeacherInvitationController) GetInvitations(ctx *gin.Context) {
	helper := NewContextHelper(ctx)

	centerID := helper.MustParamUint("id")
	if centerID == 0 {
		return
	}

	var req PaginationRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		req.Page = 1
		req.Limit = 20
	}

	status := helper.QueryStringOrDefault("status", "")

	invitations, total, err := ctl.invitationRepo.ListByCenterIDPaginated(ctx, centerID, int(req.Page), int(req.Limit), status)
	if err != nil {
		helper.InternalError("Failed to get invitations")
		return
	}

	helper.Success(PaginationResponse{
		Data:       invitations,
		Total:      total,
		Page:       req.Page,
		Limit:      req.Limit,
		TotalPages: (total + int64(req.Limit) - 1) / int64(req.Limit),
	})
}

// ========== 通用輔助方法 ==========

// requireTeacherID 取得並驗證老師 ID（老師端 API）
func (ctl *TeacherInvitationController) requireTeacherID(helper *ContextHelper) uint {
	teacherID := helper.MustUserID()
	if teacherID == 0 {
		return 0
	}
	return teacherID
}

// requireAdminID 取得並驗證管理員 ID（管理員端 API）
func (ctl *TeacherInvitationController) requireAdminID(helper *ContextHelper) uint {
	adminID := helper.MustUserID()
	if adminID == 0 {
		return 0
	}
	return adminID
}

// requireCenterID 從 URL 參數取得中心 ID
func (ctl *TeacherInvitationController) requireCenterID(helper *ContextHelper) uint {
	centerID := helper.MustParamUint("id")
	if centerID == 0 {
		return 0
	}
	return centerID
}

// requireAdminAndCenterID 取得管理員 ID 和中心 ID
func (ctl *TeacherInvitationController) requireAdminAndCenterID(helper *ContextHelper) (uint, uint) {
	adminID := ctl.requireAdminID(helper)
	if adminID == 0 {
		return 0, 0
	}
	centerID := ctl.requireCenterID(helper)
	if centerID == 0 {
		return adminID, 0
	}
	return adminID, centerID
}

// getCenterName 取得中心名稱
func (ctl *TeacherInvitationController) getCenterName(ctx context.Context, centerID uint) string {
	center, err := ctl.centerRepo.GetByID(ctx, centerID)
	if err == nil {
		return center.Name
	}
	return ""
}

// buildInvitationLinks 批次產生邀請連結
func (ctl *TeacherInvitationController) buildInvitationLinks(ctx context.Context, invitations []models.CenterInvitation, centerID uint) []services.InvitationLinkResponse {
	baseURL := ctl.app.Env.FrontendBaseURL
	if baseURL == "" {
		baseURL = "https://timeledger.app"
	}

	centerName := ctl.getCenterName(ctx, centerID)
	response := make([]services.InvitationLinkResponse, 0, len(invitations))

	for _, inv := range invitations {
		response = append(response, services.InvitationLinkResponse{
			ID:         inv.ID,
			CenterID:   centerID,
			CenterName: centerName,
			Email:      inv.Email,
			Role:       inv.Role,
			Token:      inv.Token,
			InviteLink: baseURL + "/invite/" + inv.Token,
			Status:     string(inv.Status),
			Message:    inv.Message,
			CreatedAt:  inv.CreatedAt,
			ExpiresAt:  inv.ExpiresAt,
		})
	}

	return response
}

// validateInvitationToken 驗證邀請 token
func (ctl *TeacherInvitationController) validateInvitationToken(helper *ContextHelper) (models.CenterInvitation, bool) {
	token := helper.GinContext().Param("token")
	if token == "" {
		helper.BadRequest("Token required")
		return models.CenterInvitation{}, false
	}

	invitation, err := ctl.invitationRepo.GetByToken(helper.GinContext().Request.Context(), token)
	if err != nil {
		helper.NotFound("Invitation not found")
		return models.CenterInvitation{}, false
	}

	// 檢查是否過期（通用邀請無期限，跳過檢查）
	if invitation.InviteType != models.InvitationTypeGeneral && invitation.ExpiresAt != nil && time.Now().After(*invitation.ExpiresAt) {
		if invitation.Status == models.InvitationStatusPending {
			ctl.invitationRepo.UpdateStatus(helper.GinContext().Request.Context(), invitation.ID, models.InvitationStatusExpired)
		}
		helper.Conflict("Invitation has expired")
		return models.CenterInvitation{}, false
	}

	return invitation, true
}

// ========== 老師端邀請 API ==========

// GetTeacherInvitations 取得老師的邀請列表
// @Summary 取得老師的邀請列表
// @Tags Teacher - Invitations
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param status query string false "篩選狀態 (PENDING/ACCEPTED/DECLINED/EXPIRED)"
// @Success 200 {object} global.ApiResponse{data=[]resources.InvitationResponse}
// @Router /api/v1/teacher/me/invitations [get]
func (ctl *TeacherInvitationController) GetTeacherInvitations(ctx *gin.Context) {
	helper := NewContextHelper(ctx)
	teacherID := ctl.requireTeacherID(helper)
	if teacherID == 0 {
		return
	}

	status := helper.QueryStringOrDefault("status", "")
	invitations, err := ctl.invitationRepo.ListByTeacherWithStatus(ctx.Request.Context(), teacherID, status)
	if err != nil {
		helper.InternalError("Failed to get invitations")
		return
	}

	// 批量查詢中心名稱
	centerNames := make(map[uint]string)
	for _, inv := range invitations {
		if _, exists := centerNames[inv.CenterID]; !exists {
			centerNames[inv.CenterID] = ctl.getCenterName(ctx.Request.Context(), inv.CenterID)
		}
	}

	response := ctl.invitationRes.ToInvitationResponseList(invitations, centerNames)
	helper.Success(response)
}

// RespondToInvitation 老師回應邀請
// @Summary 老師回應邀請
// @Tags Teacher - Invitations
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body services.RespondToInvitationRequest true "回應請求"
// @Success 200 {object} global.ApiResponse
// @Router /api/v1/teacher/me/invitations/respond [post]
func (ctl *TeacherInvitationController) RespondToInvitation(ctx *gin.Context) {
	helper := NewContextHelper(ctx)
	teacherID := ctl.requireTeacherID(helper)
	if teacherID == 0 {
		return
	}

	var req services.RespondToInvitationRequest
	if !helper.MustBindJSON(&req) {
		return
	}
	req.TeacherID = teacherID

	result, errInfo, err := ctl.teacherService.RespondToInvitation(ctx.Request.Context(), &req)
	if err != nil {
		helper.ErrorWithInfo(errInfo)
		return
	}

	helper.Success(gin.H{
		"invitation_id": result.InvitationID,
		"status":        result.Status,
		"message":       result.Message,
	})
}

// GetPendingInvitationsCount 取得待處理邀請數量
// @Summary 取得待處理邀請數量
// @Tags Teacher - Invitations
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} global.ApiResponse{data=map[string]int}
// @Router /api/v1/teacher/me/invitations/pending-count [get]
func (ctl *TeacherInvitationController) GetPendingInvitationsCount(ctx *gin.Context) {
	helper := NewContextHelper(ctx)
	teacherID := ctl.requireTeacherID(helper)
	if teacherID == 0 {
		return
	}

	invitations, err := ctl.invitationRepo.GetPendingByTeacher(ctx.Request.Context(), teacherID)
	if err != nil {
		helper.InternalError("Failed to get pending invitations")
		return
	}

	helper.Success(gin.H{"count": len(invitations)})
}

// ========== 管理員端邀請 API ==========

// InviteTeacher 邀請老師加入中心
// @Summary 邀請老師加入中心
// @Tags Admin
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Center ID"
// @Param request body services.InviteTeacherRequest true "邀請資訊"
// @Success 200 {object} global.ApiResponse{data=models.CenterInvitation}
// @Router /api/v1/admin/centers/{id}/invitations [post]
func (ctl *TeacherInvitationController) InviteTeacher(ctx *gin.Context) {
	helper := NewContextHelper(ctx)
	_, centerID := ctl.requireAdminAndCenterID(helper)
	if centerID == 0 {
		return
	}

	var req services.InviteTeacherRequest
	if !helper.MustBindJSON(&req) {
		return
	}
	req.CenterID = centerID
	req.AdminID = ctl.requireAdminID(helper)

	result, errInfo, err := ctl.teacherService.InviteTeacher(ctx.Request.Context(), &req)
	if err != nil {
		helper.ErrorWithInfo(errInfo)
		return
	}

	helper.Success(result.Invitation)
}

// ========== 邀請連結管理 API ==========

// GenerateInvitationLink 產生邀請連結
// @Summary 產生邀請連結
// @Tags Admin - Invitations
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Center ID"
// @Param request body services.GenerateInvitationLinkRequest true "邀請資訊"
// @Success 200 {object} global.ApiResponse{data=services.InvitationLinkResponse}
// @Router /api/v1/admin/centers/{id}/invitations/generate-link [post]
func (ctl *TeacherInvitationController) GenerateInvitationLink(ctx *gin.Context) {
	helper := NewContextHelper(ctx)
	_, centerID := ctl.requireAdminAndCenterID(helper)
	if centerID == 0 {
		return
	}

	var req services.GenerateInvitationLinkRequest
	if !helper.MustBindJSON(&req) {
		return
	}
	req.CenterID = centerID
	req.AdminID = ctl.requireAdminID(helper)

	result, errInfo, err := ctl.teacherService.GenerateInvitationLink(ctx.Request.Context(), &req)
	if err != nil {
		helper.ErrorWithInfo(errInfo)
		return
	}

	helper.Success(result.InvitationLinkResponse)
}

// GetInvitationLinks 取得邀請連結列表
// @Summary 取得邀請連結列表
// @Tags Admin - Invitations
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Center ID"
// @Success 200 {object} global.ApiResponse{data=[]services.InvitationLinkResponse}
// @Router /api/v1/admin/centers/{id}/invitations/links [get]
func (ctl *TeacherInvitationController) GetInvitationLinks(ctx *gin.Context) {
	helper := NewContextHelper(ctx)
	centerID := ctl.requireCenterID(helper)
	if centerID == 0 {
		return
	}

	invitations, err := ctl.invitationRepo.GetPendingByCenter(ctx.Request.Context(), centerID)
	if err != nil {
		helper.InternalError("Failed to get invitation links")
		return
	}

	response := ctl.buildInvitationLinks(ctx.Request.Context(), invitations, centerID)
	helper.Success(response)
}

// RevokeInvitationLink 撤回邀請連結
// @Summary 撤回邀請連結
// @Tags Admin - Invitations
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Invitation ID"
// @Success 200 {object} global.ApiResponse
// @Router /api/v1/admin/invitations/links/{id} [delete]
func (ctl *TeacherInvitationController) RevokeInvitationLink(ctx *gin.Context) {
	helper := NewContextHelper(ctx)
	adminID := ctl.requireAdminID(helper)
	if adminID == 0 {
		return
	}

	invitationID := helper.MustParamUint("id")
	if invitationID == 0 {
		return
	}

	errInfo, err := ctl.teacherService.RevokeInvitationLink(ctx.Request.Context(), invitationID, adminID)
	if err != nil {
		helper.ErrorWithInfo(errInfo)
		return
	}

	helper.Success(gin.H{"message": "Invitation link revoked"})
}

// ========== 通用邀請連結 API ==========

// GenerateGeneralInvitationLink 產生通用邀請連結
// @Summary 產生通用邀請連結
// @Tags Admin - Invitations
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Center ID"
// @Param request body services.GenerateGeneralInvitationLinkRequest true "邀請資訊"
// @Success 200 {object} global.ApiResponse{data=services.InvitationLinkResponse}
// @Router /api/v1/admin/centers/{id}/invitations/general-link [post]
func (ctl *TeacherInvitationController) GenerateGeneralInvitationLink(ctx *gin.Context) {
	helper := NewContextHelper(ctx)
	_, centerID := ctl.requireAdminAndCenterID(helper)
	if centerID == 0 {
		return
	}

	var req services.GenerateGeneralInvitationLinkRequest
	if !helper.MustBindJSON(&req) {
		return
	}
	req.CenterID = centerID
	req.AdminID = ctl.requireAdminID(helper)

	result, errInfo, err := ctl.teacherService.GenerateGeneralInvitationLink(ctx.Request.Context(), &req)
	if err != nil {
		helper.ErrorWithInfo(errInfo)
		return
	}

	helper.Success(result.InvitationLinkResponse)
}

// GetGeneralInvitationLink 取得通用邀請連結
// @Summary 取得通用邀請連結
// @Tags Admin - Invitations
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Center ID"
// @Success 200 {object} global.ApiResponse{data=services.InvitationLinkResponse}
// @Router /api/v1/admin/centers/{id}/invitations/general-link [get]
func (ctl *TeacherInvitationController) GetGeneralInvitationLink(ctx *gin.Context) {
	helper := NewContextHelper(ctx)
	centerID := ctl.requireCenterID(helper)
	if centerID == 0 {
		return
	}

	invitations, err := ctl.invitationRepo.GetPendingByCenter(ctx.Request.Context(), centerID)
	if err != nil {
		helper.InternalError("Failed to get invitation links")
		return
	}

	// 找出通用邀請
	var generalInvitation *models.CenterInvitation
	for _, inv := range invitations {
		if inv.InviteType == models.InvitationTypeGeneral {
			generalInvitation = &inv
			break
		}
	}

	if generalInvitation == nil {
		helper.Success(nil)
		return
	}

	response := ctl.buildInvitationLinks(ctx.Request.Context(), []models.CenterInvitation{*generalInvitation}, centerID)
	if len(response) > 0 {
		helper.Success(response[0])
		return
	}

	helper.Success(nil)
}

// ToggleGeneralInvitationStatus 啟用或停用通用邀請連結
// @Summary 啟用或停用通用邀請連結
// @Tags Admin - Invitations
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Center ID"
// @Success 200 {object} global.ApiResponse
// @Router /api/v1/admin/centers/{id}/invitations/general-link/toggle [post]
func (ctl *TeacherInvitationController) ToggleGeneralInvitationStatus(ctx *gin.Context) {
	helper := NewContextHelper(ctx)
	_, centerID := ctl.requireAdminAndCenterID(helper)
	if centerID == 0 {
		return
	}

	adminID := ctl.requireAdminID(helper)

	errInfo, err := ctl.teacherService.ToggleGeneralInvitationStatus(ctx.Request.Context(), centerID, adminID)
	if err != nil {
		helper.ErrorWithInfo(errInfo)
		return
	}

	helper.Success(gin.H{"message": "General invitation link status toggled"})
}

// RegenerateGeneralInvitationLink 重新產生通用邀請連結
// @Summary 重新產生通用邀請連結
// @Tags Admin - Invitations
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Center ID"
// @Param request body services.GenerateGeneralInvitationLinkRequest true "邀請資訊"
// @Success 200 {object} global.ApiResponse{data=services.InvitationLinkResponse}
// @Router /api/v1/admin/centers/{id}/invitations/general-link/regenerate [post]
func (ctl *TeacherInvitationController) RegenerateGeneralInvitationLink(ctx *gin.Context) {
	helper := NewContextHelper(ctx)
	_, centerID := ctl.requireAdminAndCenterID(helper)
	if centerID == 0 {
		return
	}

	var req services.GenerateGeneralInvitationLinkRequest
	if !helper.MustBindJSON(&req) {
		return
	}
	req.CenterID = centerID
	req.AdminID = ctl.requireAdminID(helper)

	result, errInfo, err := ctl.teacherService.GenerateGeneralInvitationLink(ctx.Request.Context(), &req)
	if err != nil {
		helper.ErrorWithInfo(errInfo)
		return
	}

	helper.Success(result.InvitationLinkResponse)
}

// ========== 公開邀請 API（無需認證）==========

// AcceptInvitationByLinkRequest 透過連結接受邀請請求
type AcceptInvitationByLinkRequest struct {
	LineUserID  string `json:"line_user_id" binding:"required"`
	AccessToken string `json:"access_token" binding:"required"`
	Email       string `json:"email"`
}

// GetPublicInvitation 取得公開邀請資訊
// @Summary 取得公開邀請資訊
// @Tags Invitations
// @Accept json
// @Produce json
// @Param token path string true "Invitation Token"
// @Success 200 {object} global.ApiResponse{data=resources.PublicInvitationInfo}
// @Router /api/v1/invitations/{token} [get]
func (ctl *TeacherInvitationController) GetPublicInvitation(ctx *gin.Context) {
	helper := NewContextHelper(ctx)
	invitation, valid := ctl.validateInvitationToken(helper)
	if !valid {
		return
	}

	centerName := ctl.getCenterName(ctx.Request.Context(), invitation.CenterID)
	response := ctl.invitationRes.ToPublicInvitationInfo(invitation, centerName)
	helper.Success(response)
}

// Accept透過連結接受邀請
// @InvitationByLink Summary 透過連結接受邀請
// @Tags Invitations
// @Accept json
// @Produce json
// @Param token path string true "Invitation Token"
// @Param request body AcceptInvitationByLinkRequest true "LINE ID Token"
// @Success 200 {object} global.ApiResponse
// @Router /api/v1/invitations/{token}/accept [post]
func (ctl *TeacherInvitationController) AcceptInvitationByLink(ctx *gin.Context) {
	helper := NewContextHelper(ctx)
	invitation, valid := ctl.validateInvitationToken(helper)
	if !valid {
		return
	}

	var req AcceptInvitationByLinkRequest
	if !helper.MustBindJSON(&req) {
		return
	}

	serviceReq := &services.AcceptInvitationByLinkRequest{
		Token:       invitation.Token,
		LineUserID:  req.LineUserID,
		AccessToken: req.AccessToken,
		Email:       req.Email,
	}

	result, errInfo, err := ctl.teacherService.AcceptInvitationByLink(ctx.Request.Context(), serviceReq)
	if err != nil {
		helper.ErrorWithInfo(errInfo)
		return
	}

	// 處理已成為成員的情況
	if result.Status == "ALREADY_MEMBER" {
		helper.Success(gin.H{
			"invitation_id": result.InvitationID,
			"status":        result.Status,
			"center_id":     result.CenterID,
			"center_name":   result.CenterName,
			"role":          result.Role,
			"token":         result.Token,
			"teacher":       result.Teacher,
			"message":       "You are already a member of this center",
		})
		return
	}

	helper.Success(gin.H{
		"invitation_id": result.InvitationID,
		"status":        result.Status,
		"center_id":     result.CenterID,
		"center_name":   result.CenterName,
		"role":          result.Role,
		"token":         result.Token,
		"teacher":       result.Teacher,
		"message":       "Successfully joined the center",
	})
}
