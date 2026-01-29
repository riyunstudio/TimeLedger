package models

import (
	"time"

	"gorm.io/gorm"
)

// InvitationType 邀請類型
type InvitationType string

const (
	InvitationTypeTalentPool InvitationType = "TALENT_POOL" // 人才庫邀請
	InvitationTypeTeacher    InvitationType = "TEACHER"     // 老師邀請
	InvitationTypeMember     InvitationType = "MEMBER"      // 會員邀請
)

// InvitationStatus 邀請狀態
type InvitationStatus string

const (
	InvitationStatusPending  InvitationStatus = "PENDING"  // 待處理
	InvitationStatusAccepted InvitationStatus = "ACCEPTED" // 已接受
	InvitationStatusDeclined InvitationStatus = "DECLINED" // 已拒絕
	InvitationStatusExpired  InvitationStatus = "EXPIRED"  // 已過期
)

type CenterInvitation struct {
	ID          uint              `gorm:"primaryKey" json:"id"`
	CenterID    uint              `gorm:"type:bigint unsigned;not null;index" json:"center_id"`
	TeacherID   uint              `gorm:"type:bigint unsigned;index" json:"teacher_id"` // 被邀請的老師ID
	InvitedBy   uint              `gorm:"type:bigint unsigned;not null" json:"invited_by"` // 邀請人（管理員）
	Email       string            `gorm:"type:varchar(255)" json:"email"`
	Token       string            `gorm:"type:varchar(255);uniqueIndex;not null" json:"token"`
	InviteType  InvitationType    `gorm:"type:varchar(20);default:'TALENT_POOL';not null" json:"invite_type"`
	Role        string            `gorm:"type:varchar(20);default:'TEACHER'" json:"role"` // 角色：TEACHER 或 SUBSTITUTE
	Status      InvitationStatus  `gorm:"type:varchar(20);default:'PENDING';not null;index" json:"status"`
	Message     string            `gorm:"type:text" json:"message"` // 邀請訊息
	RespondedAt *time.Time        `gorm:"type:datetime" json:"responded_at"`
	CreatedAt   time.Time         `gorm:"type:datetime;not null" json:"created_at"`
	ExpiresAt   time.Time         `gorm:"type:datetime;not null;index" json:"expires_at"`
	DeletedAt   gorm.DeletedAt    `gorm:"index" json:"-"`
}

func (CenterInvitation) TableName() string {
	return "center_invitations"
}

// AutoMigrate 自动迁移表结构
func (CenterInvitation) AutoMigrate(db *gorm.DB) error {
	return db.AutoMigrate(&CenterInvitation{})
}
