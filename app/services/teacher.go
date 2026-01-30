package services

import (
	"context"
	"crypto/rand"
	"fmt"
	"time"
	"timeLedger/app"
	"timeLedger/app/models"
	"timeLedger/app/repositories"
	"timeLedger/global/errInfos"
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
}

// NewTeacherService 建立教師服務
func NewTeacherService(app *app.App) *TeacherService {
	lineBotService := NewLineBotService(app)

	return &TeacherService{
		app:            app,
		teacherRepo:    repositories.NewTeacherRepository(app),
		centerRepo:     repositories.NewCenterRepository(app),
		membershipRepo: repositories.NewCenterMembershipRepository(app),
		invitationRepo: repositories.NewCenterInvitationRepository(app),
		adminUserRepo:  repositories.NewAdminUserRepository(app),
		auditLogRepo:   repositories.NewAuditLogRepository(app),
		lineBotService: lineBotService,
	}
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

	// 檢查是否過期
	if time.Now().After(invitation.ExpiresAt) {
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

	invitation := models.CenterInvitation{
		CenterID:   req.CenterID,
		Email:      req.Email,
		Token:      token,
		InviteType: models.InvitationTypeTeacher,
		Role:       req.Role,
		Status:     models.InvitationStatusPending,
		Message:    req.Message,
		CreatedAt:  time.Now(),
		ExpiresAt:  expiresAt,
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

	invitation := models.CenterInvitation{
		CenterID:   req.CenterID,
		Email:      req.Email,
		Token:      token,
		InviteType: models.InvitationTypeTeacher,
		Role:       req.Role,
		Status:     models.InvitationStatusPending,
		Message:    req.Message,
		CreatedAt:  time.Now(),
		ExpiresAt:  expiresAt,
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

// InvitationLinkResponse 邀請連結回應結構
type InvitationLinkResponse struct {
	ID         uint      `json:"id"`
	CenterID   uint      `json:"center_id"`
	CenterName string    `json:"center_name"`
	Email      string    `json:"email"`
	Role       string    `json:"role"`
	Token      string    `json:"token"`
	InviteLink string    `json:"invite_link"`
	Status     string    `json:"status"`
	Message    string    `json:"message,omitempty"`
	CreatedAt  time.Time `json:"created_at"`
	ExpiresAt  time.Time `json:"expires_at"`
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
	Token      string
	LineUserID string
}

// AcceptInvitationByLinkResult 透過連結接受邀請結果
type AcceptInvitationByLinkResult struct {
	InvitationID uint
	Status       string
	CenterID     uint
	CenterName   string
	Role         string
}

// AcceptInvitationByLink 透過連結接受邀請
func (s *TeacherService) AcceptInvitationByLink(ctx context.Context, req *AcceptInvitationByLinkRequest) (*AcceptInvitationByLinkResult, *errInfos.Res, error) {
	// 透過 token 取得邀請記錄
	invitation, err := s.invitationRepo.GetByToken(ctx, req.Token)
	if err != nil {
		return nil, s.app.Err.New(errInfos.NOT_FOUND), err
	}

	// 檢查狀態
	if invitation.Status != models.InvitationStatusPending {
		if time.Now().After(invitation.ExpiresAt) {
			s.invitationRepo.UpdateStatus(ctx, invitation.ID, models.InvitationStatusExpired)
		}
		return nil, s.app.Err.New(errInfos.INVALID_STATUS), fmt.Errorf("invitation already responded or expired")
	}

	// 檢查是否過期
	if time.Now().After(invitation.ExpiresAt) {
		s.invitationRepo.UpdateStatus(ctx, invitation.ID, models.InvitationStatusExpired)
		return nil, s.app.Err.New(errInfos.INVALID_STATUS), fmt.Errorf("invitation expired")
	}

	// 取得老師資料
	teacher, err := s.teacherRepo.GetByLineUserID(ctx, req.LineUserID)
	if err != nil {
		return nil, s.app.Err.New(errInfos.NOT_FOUND), fmt.Errorf("teacher not found")
	}

	// 驗證 Email
	if invitation.Email != "" && invitation.Email != teacher.Email {
		return nil, s.app.Err.New(errInfos.FORBIDDEN), fmt.Errorf("email mismatch")
	}

	// 檢查是否已經是中心成員
	_, err = s.membershipRepo.GetByCenterAndTeacher(ctx, invitation.CenterID, teacher.ID)
	if err == nil {
		// 已經是成員，更新狀態
		s.invitationRepo.UpdateStatus(ctx, invitation.ID, models.InvitationStatusAccepted)
		return &AcceptInvitationByLinkResult{
			InvitationID: invitation.ID,
			Status:       "ALREADY_MEMBER",
			CenterID:     invitation.CenterID,
			Role:         invitation.Role,
		}, nil, nil
	}

	// 建立 CenterMembership
	membership := models.CenterMembership{
		CenterID:  invitation.CenterID,
		TeacherID: teacher.ID,
		Status:    invitation.Role,
	}
	if _, err := s.membershipRepo.Create(ctx, membership); err != nil {
		return nil, s.app.Err.New(errInfos.SQL_ERROR), err
	}

	// 更新邀請狀態
	now := time.Now()
	if err := s.invitationRepo.UpdateFields(ctx, invitation.ID, map[string]interface{}{
		"status":       models.InvitationStatusAccepted,
		"responded_at": &now,
	}); err != nil {
		return nil, s.app.Err.New(errInfos.SQL_ERROR), err
	}

	// 審核日誌
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
				"invite_type": "LINK",
			},
		},
	})

	// 取得中心資訊
	center, err := s.centerRepo.GetByID(ctx, invitation.CenterID)
	centerName := ""
	if err == nil {
		centerName = center.Name
	}

	// 發送 LINE 通知（異步）
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
		_ = s.lineBotService.SendInvitationAcceptedNotification(notifyCtx, adminPtrs, &teacher, centerName, invitation.Role)
	}()

	return &AcceptInvitationByLinkResult{
		InvitationID: invitation.ID,
		Status:       "ACCEPTED",
		CenterID:     invitation.CenterID,
		CenterName:   centerName,
		Role:         invitation.Role,
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
