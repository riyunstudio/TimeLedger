package resources

import (
	"time"
	"timeLedger/app/models"
)

// InvitationResource 邀請相關資源轉換
type InvitationResource struct{}

// NewInvitationResource 建立邀請資源轉換實例
func NewInvitationResource() *InvitationResource {
	return &InvitationResource{}
}

// InvitationResponse 邀請回應結構
type InvitationResponse struct {
	ID          uint       `json:"id"`
	CenterID    uint       `json:"center_id"`
	CenterName  string     `json:"center_name"`
	InviteType  string     `json:"invite_type"`
	Status      string     `json:"status"`
	Message     string     `json:"message,omitempty"`
	CreatedAt   time.Time  `json:"created_at"`
	ExpiresAt   *time.Time `json:"expires_at"`
	RespondedAt *time.Time `json:"responded_at,omitempty"`
}

// ToInvitationResponse 將模型轉換為邀請回應
func (r *InvitationResource) ToInvitationResponse(inv *models.CenterInvitation, centerName string) *InvitationResponse {
	return &InvitationResponse{
		ID:          inv.ID,
		CenterID:    inv.CenterID,
		CenterName:  centerName,
		InviteType:  string(inv.InviteType),
		Status:      string(inv.Status),
		Message:     inv.Message,
		CreatedAt:   inv.CreatedAt,
		ExpiresAt:   inv.ExpiresAt,
		RespondedAt: inv.RespondedAt,
	}
}

// ToInvitationResponseList 將模型列表轉換為回應列表
func (r *InvitationResource) ToInvitationResponseList(invitations []models.CenterInvitation, centerNames map[uint]string) []*InvitationResponse {
	if invitations == nil {
		return nil
	}

	result := make([]*InvitationResponse, 0, len(invitations))
	for i := range invitations {
		inv := &invitations[i]
		centerName := ""
		if name, ok := centerNames[inv.CenterID]; ok {
			centerName = name
		}
		result = append(result, r.ToInvitationResponse(inv, centerName))
	}
	return result
}

// PublicInvitationInfo 公開邀請資訊
type PublicInvitationInfo struct {
	ID         uint       `json:"id"`
	CenterID   uint       `json:"center_id"`
	CenterName string     `json:"center_name"`
	Role       string     `json:"role"`
	Status     string     `json:"status"`
	Message    string     `json:"message,omitempty"`
	ExpiresAt  *time.Time `json:"expires_at"`
	InvitedBy  uint       `json:"-"`
}

// ToPublicInvitationInfo 將模型轉換為公開邀請資訊
func (r *InvitationResource) ToPublicInvitationInfo(inv models.CenterInvitation, centerName string) *PublicInvitationInfo {
	return &PublicInvitationInfo{
		ID:         inv.ID,
		CenterID:   inv.CenterID,
		CenterName: centerName,
		Role:       inv.Role,
		Status:     string(inv.Status),
		Message:    inv.Message,
		ExpiresAt:  inv.ExpiresAt,
		InvitedBy:  inv.InvitedBy,
	}
}
