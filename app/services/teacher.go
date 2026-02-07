package services

import (
	"context"
	"crypto/rand"
	"errors"
	"fmt"
	"time"

	"timeLedger/app"
	"timeLedger/app/models"
	"timeLedger/app/repositories"
	"timeLedger/database/redis"
	"timeLedger/global/errInfos"
	jwt "timeLedger/libs/jwt"

	"gorm.io/gorm"
)

// TeacherService 教師相關業務邏輯
// 處理邀請、個人行程等教師專屬業務
type TeacherService struct {
	BaseService
	app            *app.App
	teacherRepo    *repositories.TeacherRepository
	centerRepo     *repositories.CenterRepository
	membershipRepo *repositories.CenterMembershipRepository
	invitationRepo *repositories.CenterInvitationRepository
	adminUserRepo  *repositories.AdminUserRepository
	auditLogRepo   *repositories.AuditLogRepository
	lineBotService LineBotService
	authService    *authService
	redisClient    *redis.Redis
}

// NewTeacherService 建立教師服務
func NewTeacherService(app *app.App) *TeacherService {
	svc := &TeacherService{
		app: app,
	}

	if app.Env != nil {
		svc.redisClient = app.Redis
		svc.authService = NewAuthService(app)
		svc.lineBotService = NewLineBotService(app)
	}

	if app.MySQL != nil {
		svc.teacherRepo = repositories.NewTeacherRepository(app)
		svc.centerRepo = repositories.NewCenterRepository(app)
		svc.membershipRepo = repositories.NewCenterMembershipRepository(app)
		svc.invitationRepo = repositories.NewCenterInvitationRepository(app)
		svc.adminUserRepo = repositories.NewAdminUserRepository(app)
		svc.auditLogRepo = repositories.NewAuditLogRepository(app)
	}

	return svc
}

// ==================== 邀請相關業務邏輯 ====================

// RespondToInvitationRequest 回應邀請請求
type RespondToInvitationRequest struct {
	InvitationID uint
	Response     string // "ACCEPT" 或 "DECLINE"
	TeacherID    uint
}

// RespondToInvitationResult 回應邀請結果
type RespondToInvitationResult struct {
	InvitationID uint
	Status       string
	Message      string
}

// RespondToInvitation 老師回應邀請
func (s *TeacherService) RespondToInvitation(ctx context.Context, req *RespondToInvitationRequest) (*RespondToInvitationResult, *errInfos.Res, error) {
	// 取得邀請記錄
	invitation, err := s.invitationRepo.GetByID(ctx, req.InvitationID)
	if err != nil {
		return nil, s.app.Err.New(errInfos.NOT_FOUND), err
	}

	// 驗證是否為該老師的邀請
	if invitation.TeacherID != 0 && invitation.TeacherID != req.TeacherID {
		teacher, err := s.teacherRepo.GetByID(ctx, req.TeacherID)
		if err != nil || teacher.Email != invitation.Email {
			return nil, s.app.Err.New(errInfos.FORBIDDEN), fmt.Errorf("not authorized")
		}
	}

	// 檢查邀請狀態
	if invitation.Status != models.InvitationStatusPending {
		return nil, s.app.Err.New(errInfos.INVALID_STATUS), fmt.Errorf("invitation already responded")
	}

	// 檢查是否過期（通用邀請無期限，跳過檢查）
	if invitation.ExpiresAt != nil && time.Now().After(*invitation.ExpiresAt) {
		s.invitationRepo.UpdateStatus(ctx, req.InvitationID, models.InvitationStatusExpired)
		return nil, s.app.Err.New(errInfos.INVALID_STATUS), fmt.Errorf("invitation expired")
	}

	// 根據回應更新狀態
	var newStatus models.InvitationStatus
	if req.Response == "ACCEPT" {
		newStatus = models.InvitationStatusAccepted

		// 建立 CenterMembership（如果是老師邀請）
		if invitation.InviteType == models.InvitationTypeTeacher {
			_, err := s.membershipRepo.GetByCenterAndTeacher(ctx, invitation.CenterID, req.TeacherID)
			if err != nil {
				membership := models.CenterMembership{
					CenterID:  invitation.CenterID,
					TeacherID: req.TeacherID,
					Status:    invitation.Role,
				}
				if _, err := s.membershipRepo.Create(ctx, membership); err != nil {
					return nil, s.app.Err.New(errInfos.SQL_ERROR), err
				}

				// 審核日誌
				s.auditLogRepo.Create(ctx, models.AuditLog{
					CenterID:   invitation.CenterID,
					ActorType:  "TEACHER",
					ActorID:    req.TeacherID,
					Action:     "JOIN_CENTER",
					TargetType: "CenterMembership",
					Payload: models.AuditPayload{
						After: map[string]interface{}{
							"teacher_id": req.TeacherID,
							"center_id":  invitation.CenterID,
							"role":       invitation.Role,
						},
					},
				})
			}
		}

		// 更新人才庫狀態（如果是人才庫邀請）
		if invitation.InviteType == models.InvitationTypeTalentPool {
			teacher, _ := s.teacherRepo.GetByID(ctx, req.TeacherID)
			if !teacher.IsOpenToHiring {
				s.teacherRepo.UpdateFields(ctx, req.TeacherID, map[string]interface{}{
					"is_open_to_hiring": true,
				})
			}
		}
	} else {
		newStatus = models.InvitationStatusDeclined
	}

	// 更新邀請狀態
	if err := s.invitationRepo.UpdateStatus(ctx, req.InvitationID, newStatus); err != nil {
		return nil, s.app.Err.New(errInfos.SQL_ERROR), err
	}

	// 記錄審計日誌
	s.auditLogRepo.Create(ctx, models.AuditLog{
		CenterID:   invitation.CenterID,
		ActorType:  "TEACHER",
		ActorID:    req.TeacherID,
		Action:     "INVITATION_" + req.Response,
		TargetType: "CenterInvitation",
		TargetID:   invitation.ID,
		Payload: models.AuditPayload{
			Before: map[string]interface{}{
				"status": string(invitation.Status),
			},
			After: map[string]interface{}{
				"status": string(newStatus),
			},
		},
	})

	message := "已接受邀請"
	if newStatus == models.InvitationStatusDeclined {
		message = "已婉拒邀請"
	}

	return &RespondToInvitationResult{
		InvitationID: req.InvitationID,
		Status:       string(newStatus),
		Message:      message,
	}, nil, nil
}

// InviteTeacherRequest 邀請老師請求
type InviteTeacherRequest struct {
	CenterID uint
	AdminID  uint
	Email    string
	Role     string
	Message  string
}

// InviteTeacherResult 邀請老師結果
type InviteTeacherResult struct {
	Invitation models.CenterInvitation
}

// InviteTeacher 邀請老師加入中心
func (s *TeacherService) InviteTeacher(ctx context.Context, req *InviteTeacherRequest) (*InviteTeacherResult, *errInfos.Res, error) {
	token := s.generateInviteToken()
	expiresAt := time.Now().Add(72 * time.Hour)
	expiresAtPtr := &expiresAt

	invitation := models.CenterInvitation{
		CenterID:   req.CenterID,
		Email:      req.Email,
		Token:      token,
		InviteType: models.InvitationTypeTeacher,
		Role:       req.Role,
		Status:     models.InvitationStatusPending,
		Message:    req.Message,
		CreatedAt:  time.Now(),
		ExpiresAt:  expiresAtPtr,
	}

	result, err := s.invitationRepo.Create(ctx, invitation)
	if err != nil {
		return nil, s.app.Err.New(errInfos.SQL_ERROR), err
	}
	invitation.ID = result.ID

	// 審核日誌
	s.auditLogRepo.Create(ctx, models.AuditLog{
		CenterID:   req.CenterID,
		ActorType:  "ADMIN",
		ActorID:    req.AdminID,
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

	return &InviteTeacherResult{Invitation: invitation}, nil, nil
}

// GenerateInvitationLinkRequest 產生邀請連結請求
type GenerateInvitationLinkRequest struct {
	CenterID uint
	AdminID  uint
	Email    string
	Role     string
	Message  string
}

// GenerateInvitationLinkResult 產生邀請連結結果
type GenerateInvitationLinkResult struct {
	InvitationLinkResponse
}

// GenerateInvitationLink 產生邀請連結
func (s *TeacherService) GenerateInvitationLink(ctx context.Context, req *GenerateInvitationLinkRequest) (*GenerateInvitationLinkResult, *errInfos.Res, error) {
	token := s.generateInviteToken()
	expiresAt := time.Now().Add(72 * time.Hour)
	expiresAtPtr := &expiresAt

	invitation := models.CenterInvitation{
		CenterID:   req.CenterID,
		Email:      req.Email,
		Token:      token,
		InviteType: models.InvitationTypeTeacher,
		Role:       req.Role,
		Status:     models.InvitationStatusPending,
		Message:    req.Message,
		CreatedAt:  time.Now(),
		ExpiresAt:  expiresAtPtr,
	}

	if _, err := s.invitationRepo.Create(ctx, invitation); err != nil {
		return nil, s.app.Err.New(errInfos.SQL_ERROR), err
	}

	// 取得中心資訊
	center, err := s.centerRepo.GetByID(ctx, req.CenterID)
	centerName := ""
	if err == nil {
		centerName = center.Name
	}

	// 產生邀請連結
	baseURL := s.app.Env.FrontendBaseURL
	if baseURL == "" {
		baseURL = "https://timeledger.app"
	}
	inviteLink := fmt.Sprintf("%s/invite/%s", baseURL, token)

	// 審核日誌
	s.auditLogRepo.Create(ctx, models.AuditLog{
		CenterID:   req.CenterID,
		ActorType:  "ADMIN",
		ActorID:    req.AdminID,
		Action:     "GENERATE_INVITATION_LINK",
		TargetType: "CenterInvitation",
		TargetID:   invitation.ID,
		Payload: models.AuditPayload{
			After: map[string]interface{}{
				"email":       req.Email,
				"role":        req.Role,
				"status":      "PENDING",
				"expires_at":  expiresAt,
				"invite_link": inviteLink,
			},
		},
	})

	result := GenerateInvitationLinkResult{
		InvitationLinkResponse: InvitationLinkResponse{
			ID:         invitation.ID,
			CenterID:   req.CenterID,
			CenterName: centerName,
			Email:      req.Email,
			Role:       req.Role,
			Token:      token,
			InviteLink: inviteLink,
			Status:     string(invitation.Status),
			Message:    req.Message,
			CreatedAt:  invitation.CreatedAt,
			ExpiresAt:  invitation.ExpiresAt,
		},
	}

	return &result, nil, nil
}

// GenerateGeneralInvitationLinkRequest 產生通用邀請連結請求
type GenerateGeneralInvitationLinkRequest struct {
	CenterID uint
	AdminID  uint
	Role     string // 角色：TEACHER 或 SUBSTITUTE
	Message  string // 邀請訊息
}

// GenerateGeneralInvitationLinkResult 產生通用邀請連結結果
type GenerateGeneralInvitationLinkResult struct {
	InvitationLinkResponse
}

// GenerateGeneralInvitationLink 產生通用邀請連結（不綁定 Email，無期限）
func (s *TeacherService) GenerateGeneralInvitationLink(ctx context.Context, req *GenerateGeneralInvitationLinkRequest) (*GenerateGeneralInvitationLinkResult, *errInfos.Res, error) {
	// 產生新的通用邀請連結（無期限）
	token := s.generateInviteToken()

	invitation := models.CenterInvitation{
		CenterID:   req.CenterID,
		Email:      "", // 通用邀請不綁定 Email
		Token:      token,
		InviteType: models.InvitationTypeGeneral,
		Role:       req.Role,
		Status:     models.InvitationStatusPending,
		Message:    req.Message,
		CreatedAt:  time.Now(),
		ExpiresAt:  nil, // 通用邀請無期限
	}

	if _, err := s.invitationRepo.Create(ctx, invitation); err != nil {
		return nil, s.app.Err.New(errInfos.SQL_ERROR), err
	}

	// 取得中心資訊
	center, err := s.centerRepo.GetByID(ctx, req.CenterID)
	centerName := ""
	if err == nil {
		centerName = center.Name
	}

	// 產生邀請連結
	baseURL := s.app.Env.FrontendBaseURL
	if baseURL == "" {
		baseURL = "https://timeledger.app"
	}
	inviteLink := fmt.Sprintf("%s/invite/%s", baseURL, token)

	// 審核日誌
	s.auditLogRepo.Create(ctx, models.AuditLog{
		CenterID:   req.CenterID,
		ActorType:  "ADMIN",
		ActorID:    req.AdminID,
		Action:     "GENERATE_GENERAL_INVITATION_LINK",
		TargetType: "CenterInvitation",
		TargetID:   invitation.ID,
		Payload: models.AuditPayload{
			After: map[string]interface{}{
				"role":        req.Role,
				"status":      "PENDING",
				"expires_at":  nil,
				"invite_link": inviteLink,
				"invite_type": "GENERAL",
			},
		},
	})

	result := GenerateGeneralInvitationLinkResult{
		InvitationLinkResponse: InvitationLinkResponse{
			ID:         invitation.ID,
			CenterID:   req.CenterID,
			CenterName: centerName,
			Email:      "",
			Role:       invitation.Role,
			Token:      token,
			InviteLink: inviteLink,
			Status:     string(invitation.Status),
			Message:    invitation.Message,
			CreatedAt:  invitation.CreatedAt,
			ExpiresAt:  nil, // 通用邀請無期限
		},
	}

	return &result, nil, nil
}

// ToggleGeneralInvitationStatusRequest 切換通用邀請狀態請求
type ToggleGeneralInvitationStatusRequest struct {
	CenterID uint
	AdminID  uint
}

// ToggleGeneralInvitationStatus 啟用或停用通用邀請連結
func (s *TeacherService) ToggleGeneralInvitationStatus(ctx context.Context, centerID, adminID uint) (*errInfos.Res, error) {
	// 取得現有的通用邀請
	existingGeneral, err := s.invitationRepo.GetGeneralByCenterID(ctx, centerID)
	if err != nil || existingGeneral.ID == 0 {
		return s.app.Err.New(errInfos.NOT_FOUND), fmt.Errorf("no general invitation found")
	}

	// 切換狀態
	var newStatus models.InvitationStatus
	if existingGeneral.Status == models.InvitationStatusPending {
		newStatus = models.InvitationStatusExpired
	} else if existingGeneral.Status == models.InvitationStatusExpired {
		// 如果是停用狀態，產生新的連結
		return s.generateNewGeneralLink(ctx, centerID, adminID, existingGeneral.Role, existingGeneral.Message)
	}

	if err := s.invitationRepo.UpdateStatus(ctx, existingGeneral.ID, newStatus); err != nil {
		return s.app.Err.New(errInfos.SQL_ERROR), err
	}

	// 審核日誌
	s.auditLogRepo.Create(ctx, models.AuditLog{
		CenterID:   centerID,
		ActorType:  "ADMIN",
		ActorID:    adminID,
		Action:     "TOGGLE_GENERAL_INVITATION_STATUS",
		TargetType: "CenterInvitation",
		TargetID:   existingGeneral.ID,
		Payload: models.AuditPayload{
			Before: map[string]interface{}{
				"status": string(existingGeneral.Status),
			},
			After: map[string]interface{}{
				"status": string(newStatus),
			},
		},
	})

	return nil, nil
}

// generateNewGeneralLink 產生新的通用邀請連結（內部方法）
func (s *TeacherService) generateNewGeneralLink(ctx context.Context, centerID, adminID uint, role, message string) (*errInfos.Res, error) {
	// 停用舊的通用邀請
	existingGeneral, _ := s.invitationRepo.GetGeneralByCenterID(ctx, centerID)
	if existingGeneral.ID > 0 {
		s.invitationRepo.UpdateStatus(ctx, existingGeneral.ID, models.InvitationStatusExpired)
	}

	// 產生新的通用邀請連結
	token := s.generateInviteToken()

	invitation := models.CenterInvitation{
		CenterID:   centerID,
		Email:      "",
		Token:      token,
		InviteType: models.InvitationTypeGeneral,
		Role:       role,
		Status:     models.InvitationStatusPending,
		Message:    message,
		CreatedAt:  time.Now(),
		ExpiresAt:  nil, // 通用邀請無期限
	}

	if _, err := s.invitationRepo.Create(ctx, invitation); err != nil {
		return s.app.Err.New(errInfos.SQL_ERROR), err
	}

	// 審核日誌
	s.auditLogRepo.Create(ctx, models.AuditLog{
		CenterID:   centerID,
		ActorType:  "ADMIN",
		ActorID:    adminID,
		Action:     "REGENERATE_GENERAL_INVITATION_LINK",
		TargetType: "CenterInvitation",
		TargetID:   invitation.ID,
		Payload: models.AuditPayload{
			After: map[string]interface{}{
				"role":        role,
				"status":      "PENDING",
				"expires_at":  nil,
				"invite_link": fmt.Sprintf("%s/invite/%s", s.app.Env.FrontendBaseURL, token),
				"invite_type": "GENERAL",
			},
		},
	})

	return nil, nil
}

// InvitationLinkResponse 邀請連結回應結構
type InvitationLinkResponse struct {
	ID         uint       `json:"id"`
	CenterID   uint       `json:"center_id"`
	CenterName string     `json:"center_name"`
	Email      string     `json:"email"`
	Role       string     `json:"role"`
	Token      string     `json:"token"`
	InviteLink string     `json:"invite_link"`
	Status     string     `json:"status"`
	Message    string     `json:"message,omitempty"`
	CreatedAt  time.Time  `json:"created_at"`
	ExpiresAt  *time.Time `json:"expires_at"`
}

// RevokeInvitationLink 撤回邀請連結
func (s *TeacherService) RevokeInvitationLink(ctx context.Context, invitationID, adminID uint) (*errInfos.Res, error) {
	// 取得邀請記錄
	invitation, err := s.invitationRepo.GetByID(ctx, invitationID)
	if err != nil {
		return s.app.Err.New(errInfos.NOT_FOUND), err
	}

	// 檢查狀態
	if invitation.Status != models.InvitationStatusPending {
		return s.app.Err.New(errInfos.INVALID_STATUS), fmt.Errorf("only pending invitations can be revoked")
	}

	// 更新狀態為過期
	if err := s.invitationRepo.UpdateStatus(ctx, invitationID, models.InvitationStatusExpired); err != nil {
		return s.app.Err.New(errInfos.SQL_ERROR), err
	}

	// 審核日誌
	s.auditLogRepo.Create(ctx, models.AuditLog{
		CenterID:   invitation.CenterID,
		ActorType:  "ADMIN",
		ActorID:    adminID,
		Action:     "REVOKE_INVITATION_LINK",
		TargetType: "CenterInvitation",
		TargetID:   invitationID,
		Payload: models.AuditPayload{
			Before: map[string]interface{}{
				"status": string(invitation.Status),
			},
			After: map[string]interface{}{
				"status": "EXPIRED",
			},
		},
	})

	return nil, nil
}

// AcceptInvitationByLinkRequest 透過連結接受邀請請求
type AcceptInvitationByLinkRequest struct {
	Token       string `json:"token"`
	LineUserID  string `json:"line_user_id"`
	AccessToken string `json:"access_token"`
	Email       string `json:"email"` // 從 LINE ID Token 獲取的 email（可選）
}

// AcceptInvitationByLinkResult 透過連結接受邀請結果
type AcceptInvitationByLinkResult struct {
	InvitationID uint         `json:"invitation_id"`
	Status       string       `json:"status"`
	CenterID     uint         `json:"center_id"`
	CenterName   string       `json:"center_name"`
	Role         string       `json:"role"`
	Token        string       `json:"token"`   // JWT Token for auto-login
	Teacher      *TeacherInfo `json:"teacher"` // Teacher info for frontend
}

// TeacherInfo 老師資訊（用於 API 回傳）
type TeacherInfo struct {
	ID         uint   `json:"id"`
	Name       string `json:"name"`
	Email      string `json:"email"`
	LineUserID string `json:"line_user_id"`
	AvatarURL  string `json:"avatar_url,omitempty"`
}

// AcceptInvitationByLink 透過連結接受邀請
func (s *TeacherService) AcceptInvitationByLink(ctx context.Context, req *AcceptInvitationByLinkRequest) (*AcceptInvitationByLinkResult, *errInfos.Res, error) {
	// 取得 LINE Profile（驗證 token 並取得用戶資料）
	lineProfile, err := s.authService.getLineProfile(req.AccessToken, req.LineUserID)
	if err != nil {
		return nil, s.app.Err.New(errInfos.UNAUTHORIZED), err
	}

	// 透過 token 取得邀請記錄
	invitation, err := s.invitationRepo.GetByToken(ctx, req.Token)
	if err != nil {
		return nil, s.app.Err.New(errInfos.NOT_FOUND), err
	}

	// 檢查狀態
	if invitation.Status != models.InvitationStatusPending {
		// 通用邀請無期限，不需要檢查過期
		if invitation.InviteType != models.InvitationTypeGeneral && invitation.ExpiresAt != nil && time.Now().After(*invitation.ExpiresAt) {
			s.invitationRepo.UpdateStatus(ctx, invitation.ID, models.InvitationStatusExpired)
		}
		return nil, s.app.Err.New(errInfos.INVALID_STATUS), fmt.Errorf("invitation already responded or expired")
	}

	// 檢查是否過期（通用邀請無期限，跳過檢查）
	if invitation.InviteType != models.InvitationTypeGeneral && invitation.ExpiresAt != nil && time.Now().After(*invitation.ExpiresAt) {
		s.invitationRepo.UpdateStatus(ctx, invitation.ID, models.InvitationStatusExpired)
		return nil, s.app.Err.New(errInfos.INVALID_STATUS), fmt.Errorf("invitation expired")
	}

	// 使用分布式鎖防止並發創建老師帳號（Race Condition）
	var teacher models.Teacher
	err = s.tryLockTeacherCreation(ctx, req.LineUserID, func() error {
		var innerErr error
		// 取得老師資料
		teacher, innerErr = s.teacherRepo.GetByLineUserID(ctx, req.LineUserID)
		if innerErr != nil {
			// 新老師：自動建立老師帳號，使用 LINE 的 displayName 作為名稱
			// Email 優先級：邀請指定 > LINE ID Token > 預設
			var teacherEmail string

			if invitation.Email != "" {
				// 指定邀請使用邀請中的 Email
				teacherEmail = invitation.Email
			} else if req.Email != "" {
				// 使用 LINE ID Token 獲取的 Email
				teacherEmail = req.Email
			} else {
				// 都沒有 email，使用固定格式
				teacherEmail = "teacher@timeledger.tw"
			}

			newTeacher := models.Teacher{
				LineUserID: req.LineUserID,
				Email:      teacherEmail,
				Name:       lineProfile.DisplayName, // 使用 LINE 的顯示名稱
				AvatarURL:  lineProfile.PictureURL,  // 使用 LINE 的頭像
			}

			createdTeacher, err := s.teacherRepo.Create(ctx, newTeacher)
			if err != nil {
				return fmt.Errorf("failed to create teacher: %w", err)
			}
			teacher = createdTeacher
		}
		return nil
	})

	if err != nil {
		return nil, s.app.Err.New(errInfos.SQL_ERROR), err
	}

	// 通用邀請跳過 Email 驗證（已移除 Email 檢查，允許 Email 不符時也能加入）

	// 檢查是否已經是中心成員
	existingMembership, err := s.membershipRepo.GetByCenterAndTeacher(ctx, invitation.CenterID, teacher.ID)
	if err == nil {
		// 通用邀請不更新邀請狀態（保持 PENDING，可重複使用）
		if invitation.InviteType != models.InvitationTypeGeneral {
			s.invitationRepo.UpdateStatus(ctx, invitation.ID, models.InvitationStatusAccepted)
		}

		// 同步成員 Role（如果不同）
		if existingMembership.Role != invitation.Role {
			if err := s.membershipRepo.UpdateFields(ctx, existingMembership.ID, map[string]interface{}{
				"role":   invitation.Role,
				"status": "ACTIVE",
			}); err != nil {
				s.Logger.Warn("failed to update membership role", "membership_id", existingMembership.ID, "error", err)
			}
		}

		// 通用邀請不回填 TeacherID（允許多人使用同一個邀請連結）
		// 只有非通用邀請且 TeacherID 為 0 時才回填
		if invitation.InviteType != models.InvitationTypeGeneral && invitation.TeacherID == 0 {
			if err := s.invitationRepo.UpdateFields(ctx, invitation.ID, map[string]interface{}{
				"teacher_id": teacher.ID,
			}); err != nil {
				s.Logger.Warn("failed to update invitation teacher_id", "invitation_id", invitation.ID, "error", err)
			}
		}

		// 取得中心名稱
		centerName := ""
		center, err := s.centerRepo.GetByID(ctx, invitation.CenterID)
		if err == nil {
			centerName = center.Name
		}

		// 生成 JWT Token
		claims := jwt.Claims{
			UserType:   "TEACHER",
			UserID:     teacher.ID,
			CenterID:   invitation.CenterID,
			LineUserID: teacher.LineUserID,
		}
		token, _ := s.authService.GenerateToken(claims)

		return &AcceptInvitationByLinkResult{
			InvitationID: invitation.ID,
			Status:       "ALREADY_MEMBER",
			CenterID:     invitation.CenterID,
			CenterName:   centerName,
			Role:         invitation.Role,
			Token:        token,
			Teacher: &TeacherInfo{
				ID:         teacher.ID,
				Name:       teacher.Name,
				Email:      teacher.Email,
				LineUserID: teacher.LineUserID,
				AvatarURL:  teacher.AvatarURL,
			},
		}, nil, nil
	}

	// 建立 CenterMembership（包含 Role 欄位）
	// 使用資料庫事務確保原子性
	membership := models.CenterMembership{
		CenterID:  invitation.CenterID,
		TeacherID: teacher.ID,
		Role:      invitation.Role,
		Status:    "ACTIVE", // 設置為 ACTIVE 狀態，供 GetActiveByTeacherID 和 ListTeacherIDsByCenterID 查詢
	}

	// 執行事務操作
	err = s.app.MySQL.WDB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// 在事務中創建 Membership
		if _, err := s.membershipRepo.Create(ctx, membership); err != nil {
			return fmt.Errorf("failed to create membership: %w", err)
		}

		// 通用邀請不回填 TeacherID（允許多人使用同一個邀請連結）
		// 只有非通用邀請且 TeacherID 為 0 時才回填
		if invitation.InviteType != models.InvitationTypeGeneral && invitation.TeacherID == 0 {
			if err := s.invitationRepo.UpdateFields(ctx, invitation.ID, map[string]interface{}{
				"teacher_id": teacher.ID,
			}); err != nil {
				return fmt.Errorf("failed to update invitation teacher_id: %w", err)
			}
		}

		// 通用邀請不更新邀請狀態（保持 PENDING，可重複使用）
		if invitation.InviteType != models.InvitationTypeGeneral {
			// 更新邀請狀態（非通用邀請）
			now := time.Now()
			if err := s.invitationRepo.UpdateFields(ctx, invitation.ID, map[string]interface{}{
				"status":       models.InvitationStatusAccepted,
				"responded_at": &now,
			}); err != nil {
				return fmt.Errorf("failed to update invitation status: %w", err)
			}
		}

		// 審核日誌（在事務中創建）
		s.auditLogRepo.Create(ctx, models.AuditLog{
			CenterID:   invitation.CenterID,
			ActorType:  "TEACHER",
			ActorID:    teacher.ID,
			Action:     "JOIN_CENTER_VIA_LINK",
			TargetType: "CenterInvitation",
			TargetID:   invitation.ID,
			Payload: models.AuditPayload{
				After: map[string]interface{}{
					"teacher_id":  teacher.ID,
					"center_id":   invitation.CenterID,
					"role":        invitation.Role,
					"status":      "ACCEPTED",
					"invite_type": string(invitation.InviteType),
				},
			},
		})

		return nil
	})

	if err != nil {
		return nil, s.app.Err.New(errInfos.SQL_ERROR), err
	}

	// 取得中心資訊
	center, err := s.centerRepo.GetByID(ctx, invitation.CenterID)
	centerName := ""
	if err == nil {
		centerName = center.Name
	}

	// 發送 LINE 通知（異步）
	// 建立 teacher 副本以確保 goroutine 執行時資料仍然有效
	teacherForNotification := teacher
	go func() {
		notifyCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		admins, err := s.adminUserRepo.GetByCenterID(notifyCtx, invitation.CenterID)
		if err != nil || len(admins) == 0 {
			return
		}
		// 轉換為 []*models.AdminUser
		adminPtrs := make([]*models.AdminUser, len(admins))
		for i := range admins {
			adminPtrs[i] = &admins[i]
		}
		_ = s.lineBotService.SendInvitationAcceptedNotification(notifyCtx, adminPtrs, &teacherForNotification, centerName, invitation.Role)
	}()

	// 生成 JWT Token
	claims := jwt.Claims{
		UserType:   "TEACHER",
		UserID:     teacher.ID,
		CenterID:   invitation.CenterID,
		LineUserID: teacher.LineUserID,
	}
	token, err := s.authService.GenerateToken(claims)
	if err != nil {
		// Token 生成失敗不應該阻止成功響應，記錄日誌即可
		s.Logger.Error("failed to generate token after accepting invitation", "error", err)
	}

	return &AcceptInvitationByLinkResult{
		InvitationID: invitation.ID,
		Status:       "ACCEPTED",
		CenterID:     invitation.CenterID,
		CenterName:   centerName,
		Role:         invitation.Role,
		Token:        token,
		Teacher: &TeacherInfo{
			ID:         teacher.ID,
			Name:       teacher.Name,
			Email:      teacher.Email,
			LineUserID: teacher.LineUserID,
			AvatarURL:  teacher.AvatarURL,
		},
	}, nil, nil
}

// generateInviteToken 產生邀請 token
func (s *TeacherService) generateInviteToken() string {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return ""
	}
	return fmt.Sprintf("%x", b)
}

// ==================== 公開註冊相關業務邏輯 ====================

// PublicRegisterRequest 公開註冊請求
type PublicRegisterRequest struct {
	LineUserID  string `json:"line_user_id" binding:"required"`
	AccessToken string `json:"access_token" binding:"required"`
	Name        string `json:"name" binding:"required"`
	Email       string `json:"email" binding:"required,email"`
}

// PublicRegisterResult 公開註冊結果
type PublicRegisterResult struct {
	Token   string `json:"token"`
	Teacher any    `json:"teacher"`
}

// RegisterPublic 公開註冊老師（LINE Bot 自主註冊）
func (s *TeacherService) RegisterPublic(ctx context.Context, req *PublicRegisterRequest) (*PublicRegisterResult, *errInfos.Res, error) {
	// 取得 LINE Profile（驗證 token 並取得用戶資料）
	lineProfile, err := s.authService.getLineProfile(req.AccessToken, req.LineUserID)
	if err != nil {
		return nil, s.app.Err.New(errInfos.UNAUTHORIZED), err
	}

	// 檢查是否已存在相同 LineUserID 的老師
	existingTeacher, err := s.teacherRepo.GetByLineUserID(ctx, req.LineUserID)
	if err == nil {
		// 已經註冊過，直接登入並產生新 token
		centerID, _ := s.teacherRepo.GetCenterID(ctx, existingTeacher.ID)

		claims := jwt.Claims{
			UserType:   "TEACHER",
			UserID:     existingTeacher.ID,
			CenterID:   centerID,
			LineUserID: existingTeacher.LineUserID,
		}

		token, err := s.authService.GenerateToken(claims)
		if err != nil {
			return nil, s.app.Err.New(errInfos.SYSTEM_ERROR), err
		}

		return &PublicRegisterResult{
			Token: token,
			Teacher: map[string]interface{}{
				"id":           existingTeacher.ID,
				"name":         existingTeacher.Name,
				"email":        existingTeacher.Email,
				"line_user_id": existingTeacher.LineUserID,
			},
		}, nil, nil
	}

	// 檢查 Email 是否已被使用
	_, err = s.teacherRepo.GetByEmail(ctx, req.Email)
	if err == nil {
		return nil, s.app.Err.New(errInfos.DUPLICATE), fmt.Errorf("email already registered")
	}

	// 建立新老師帳號（預設加入人才庫），使用 LINE 的頭像
	newTeacher := models.Teacher{
		LineUserID:     req.LineUserID,
		Name:           req.Name, // 使用用戶提供的名稱
		Email:          req.Email,
		AvatarURL:      lineProfile.PictureURL, // 使用 LINE 的頭像
		IsOpenToHiring: true,                   // 預設加入人才庫
	}

	createdTeacher, err := s.teacherRepo.Create(ctx, newTeacher)
	if err != nil {
		return nil, s.app.Err.New(errInfos.SQL_ERROR), fmt.Errorf("failed to create teacher: %w", err)
	}

	// 產生 JWT token
	claims := jwt.Claims{
		UserType:   "TEACHER",
		UserID:     createdTeacher.ID,
		CenterID:   0, // 新註冊老師尚未加入任何中心
		LineUserID: createdTeacher.LineUserID,
	}

	token, err := s.authService.GenerateToken(claims)
	if err != nil {
		return nil, s.app.Err.New(errInfos.SYSTEM_ERROR), err
	}

	return &PublicRegisterResult{
		Token: token,
		Teacher: map[string]interface{}{
			"id":                createdTeacher.ID,
			"name":              createdTeacher.Name,
			"email":             createdTeacher.Email,
			"line_user_id":      createdTeacher.LineUserID,
			"is_open_to_hiring": createdTeacher.IsOpenToHiring,
		},
	}, nil, nil
}

// ==================== 分布式鎖相關方法 ====================

const (
	// 分布式鎖前綴
	lockPrefix = "teacher:lock:"
	// 鎖的 TTL（秒）
	lockTTL = 10
)

// acquireTeacherCreationLock 嘗試獲取教師創建分布式鎖
// 使用 Redis SETNX 實現，防止並發創建重複帳號
func (s *TeacherService) acquireTeacherCreationLock(ctx context.Context, lineUserID string) (string, error) {
	if s.redisClient == nil || s.redisClient.DB0 == nil {
		// Redis 未配置，跳過鎖（單機部署環境）
		return "", nil
	}

	lockKey := lockPrefix + lineUserID

	// 使用 SET NX EX 原子操作獲取鎖
	// 如果鍵不存在，則設置值並過期時間
	// 返回 OK 表示成功獲取鎖
	result, err := s.redisClient.DB0.SetNX(ctx, lockKey, "locked", time.Duration(lockTTL)*time.Second).Result()
	if err != nil {
		s.Logger.Error("failed to acquire distributed lock", "key", lockKey, "error", err)
		return "", err
	}

	if result {
		// 成功獲取鎖
		s.Logger.Debug("acquired teacher creation lock", "line_user_id", lineUserID)
		return lockKey, nil
	}

	// 獲取鎖失敗
	s.Logger.Debug("failed to acquire teacher creation lock (already locked)", "line_user_id", lineUserID)
	return "", errors.New("system is busy, please try again")
}

// releaseTeacherCreationLock 釋放教師創建分布式鎖
func (s *TeacherService) releaseTeacherCreationLock(ctx context.Context, lockKey string) {
	if s.redisClient == nil || s.redisClient.DB0 == nil || lockKey == "" {
		return
	}

	// 刪除鎖鍵
	err := s.redisClient.DB0.Del(ctx, lockKey).Err()
	if err != nil {
		s.Logger.Warn("failed to release distributed lock", "key", lockKey, "error", err)
		return
	}

	s.Logger.Debug("released teacher creation lock", "key", lockKey)
}

// tryLockTeacherCreation 嘗試獲取鎖並執行操作
// 如果獲取鎖失敗，會重試指定的次數
func (s *TeacherService) tryLockTeacherCreation(ctx context.Context, lineUserID string, fn func() error) error {
	lockKey, err := s.acquireTeacherCreationLock(ctx, lineUserID)
	if err != nil {
		// 獲取鎖失敗，進行重試
		maxRetries := 3
		for i := 0; i < maxRetries; i++ {
			time.Sleep(time.Duration(i+1) * 100 * time.Millisecond) // 指數退避
			lockKey, err = s.acquireTeacherCreationLock(ctx, lineUserID)
			if err == nil {
				break
			}
		}

		if err != nil {
			return fmt.Errorf("failed to acquire lock after %d retries: %w", maxRetries, err)
		}
	}

	// 確保最終釋放鎖
	defer func() {
		s.releaseTeacherCreationLock(ctx, lockKey)
	}()

	// 執行被保護的操作
	return fn()
}
