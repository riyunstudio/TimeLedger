package services

import (
	"context"
	"time"
)

type SmartMatchingService interface {
	FindMatches(ctx context.Context, centerID uint, teacherID *uint, roomID uint, startTime, endTime time.Time, requiredSkills []string, excludeTeacherIDs []uint) ([]MatchScore, error)
	SearchTalent(ctx context.Context, searchParams TalentSearchParams) (*TalentSearchResultResponse, error)
	GetTalentStats(ctx context.Context, centerID uint) (*TalentStats, error)
	InviteTalent(ctx context.Context, centerID uint, adminID uint, teacherIDs []uint, message string) (*InviteResult, error)
	GetSearchSuggestions(ctx context.Context, query string) (*SearchSuggestions, error)
	GetAlternativeSlots(ctx context.Context, centerID uint, teacherID uint, originalStart, originalEnd time.Time, duration int) ([]AlternativeSlot, error)
	GetTeacherSessions(ctx context.Context, centerID uint, teacherID uint, startDate, endDate string) (*TeacherSessions, error)
}

// TalentSearchParams - 人才庫搜尋參數
type TalentSearchParams struct {
	CenterID   uint     `json:"center_id"`
	City       string   `json:"city,omitempty"`
	District   string   `json:"district,omitempty"`
	Keyword    string   `json:"keyword,omitempty"`
	Skills     []string `json:"skills,omitempty"`
	Hashtags   []string `json:"hashtags,omitempty"`

	// 分頁參數
	Page  int `json:"page"`
	Limit int `json:"limit"`

	// 排序參數
	SortBy    string `json:"sort_by"`
	SortOrder string `json:"sort_order"`
}

// TalentSearchResult - 人才庫搜尋結果
type TalentSearchResult struct {
	Talents    []TalentResult `json:"talents"`
	Pagination Pagination     `json:"pagination"`
}

// TalentSkillResponse - 人才技能回應結構（符合前端格式）
type TalentSkillResponse struct {
	Category  string `json:"category"`  // 技能分類
	SkillName string `json:"skill_name"` // 技能名稱
}

// TalentCardResponse - 人才卡片回應結構（符合前端 TalentCard 類型）
type TalentCardResponse struct {
	ID               uint                  `json:"id"`                // 教師 ID
	Name             string                `json:"name"`             // 名稱
	AvatarURL        string                `json:"avatar_url,omitempty"` // 頭像 URL
	Bio              string                `json:"bio,omitempty"`    // 簡介
	City             string                `json:"city,omitempty"`   // 縣市
	District         string                `json:"district,omitempty"`   // 區域
	Skills           []TalentSkillResponse `json:"skills"`          // 技能列表
	CertificateCount int                   `json:"certificate_count"` // 證照數量
	Certificates     []CertificateResponse `json:"certificates,omitempty"` // 證照列表
	Rating           int                   `json:"rating"`          // 評分
	ReviewCount      int                   `json:"review_count"`    // 評價數量
	IsOpenToHiring  bool                  `json:"is_open_to_hiring"` // 是否開放應徵
	IsMember        bool                  `json:"is_member"`       // 是否已加入中心
	PersonalHashtags []string              `json:"personal_hashtags,omitempty"` // 個人標籤
	PublicContactInfo string                `json:"public_contact_info,omitempty"` // 聯絡資訊
}

// CertificateResponse - 證照回應結構
type CertificateResponse struct {
	ID         uint   `json:"id"`
	Name       string `json:"name"`
	Issuer     string `json:"issuer,omitempty"`
	IssuedAt   string `json:"issued_at,omitempty"`
	ExpiryDate string `json:"expiry_date,omitempty"`
}

// TalentSearchResultResponse - 人才搜尋回應結構（符合前端格式）
type TalentSearchResultResponse struct {
	Talents    []TalentCardResponse `json:"talents"`
	Pagination Pagination           `json:"pagination"`
}

// Pagination - 分頁資訊
type Pagination struct {
	Page       int  `json:"page"`
	Limit      int  `json:"limit"`
	Total      int  `json:"total"`
	TotalPages int  `json:"total_pages"`
	HasNext    bool `json:"has_next"`
	HasPrev    bool `json:"has_prev"`
}

// TalentStats - 人才庫統計
type TalentStats struct {
	TotalCount       int                       `json:"total_count"`
	OpenHiringCount  int                       `json:"open_hiring_count"`
	MemberCount      int                       `json:"member_count"`
	AverageRating    float64                   `json:"average_rating"`
	MonthlyChange    int                       `json:"monthly_change"`
	MonthlyTrend     []int                     `json:"monthly_trend"`
	PendingInvites   int                       `json:"pending_invites"`
	AcceptedInvites  int                       `json:"accepted_invites"`
	DeclinedInvites  int                       `json:"declined_invites"`
	CityDistribution []CityDistributionItem    `json:"city_distribution"`
	TopSkills        []SkillCountItem          `json:"top_skills"`
}

type CityDistributionItem struct {
	Name  string `json:"name"`
	Count int    `json:"count"`
}

type SkillCountItem struct {
	Name  string `json:"name"`
	Count int    `json:"count"`
}

// InviteResult - 邀請結果
type InviteResult struct {
	InvitedCount  int      `json:"invited_count"`
	FailedCount   int      `json:"failed_count"`
	FailedIDs     []uint   `json:"failed_ids,omitempty"`
	InvitationIDs []uint   `json:"invitation_ids,omitempty"`
}

// SearchSuggestions - 搜尋建議
type SearchSuggestions struct {
	Skills        []SuggestionItem `json:"skills"`
	Tags          []SuggestionItem `json:"tags"`
	Names         []SuggestionItem `json:"names"`
	Trending      []string         `json:"trending"`
}

type SuggestionItem struct {
	Type  string `json:"type"`
	Value string `json:"value"`
}

// AlternativeSlot - 替代時段
type AlternativeSlot struct {
	Date          string                  `json:"date"`
	DateLabel     string                  `json:"date_label"`
	Start         string                  `json:"start"`
	End           string                  `json:"end"`
	Available     bool                    `json:"available"`
	AvailableRooms []RoomInfo             `json:"available_rooms,omitempty"`
	ConflictReason string                 `json:"conflict_reason,omitempty"`
}

type RoomInfo struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
}

// TeacherSessions - 教師課表
type TeacherSessions struct {
	TeacherID   uint         `json:"teacher_id"`
	TeacherName string       `json:"teacher_name"`
	Sessions    []SessionInfo `json:"sessions"`
}

type SessionInfo struct {
	ID         uint   `json:"id"`
	CourseName string `json:"course_name"`
	StartTime  string `json:"start_time"`
	EndTime    string `json:"end_time"`
	RoomName   string `json:"room_name,omitempty"`
	Status     string `json:"status"`
}
