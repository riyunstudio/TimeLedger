package resources

import (
	"time"
	"timeLedger/app"
	"timeLedger/app/models"
)

// TermResource 學期資源轉換
type TermResource struct {
	app *app.App
}

// NewTermResource 建立 TermResource 實例
func NewTermResource(appInstance *app.App) *TermResource {
	return &TermResource{
		app: appInstance,
	}
}

// TermResponse 學期響應結構
type TermResponse struct {
	ID        uint      `json:"id"`
	CenterID  uint      `json:"center_id"`
	Name      string    `json:"name"`
	StartDate string    `json:"start_date"`
	EndDate   string    `json:"end_date"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// ToTermResponse 將學期模型轉換為響應格式
func (r *TermResource) ToTermResponse(term models.CenterTerm) *TermResponse {
	return &TermResponse{
		ID:        term.ID,
		CenterID:  term.CenterID,
		Name:      term.Name,
		StartDate: term.StartDate.Format("2006-01-02"),
		EndDate:   term.EndDate.Format("2006-01-02"),
		CreatedAt: term.CreatedAt,
		UpdatedAt: term.UpdatedAt,
	}
}

// ToTermResponses 批量將學期模型轉換為響應格式
func (r *TermResource) ToTermResponses(terms []models.CenterTerm) []TermResponse {
	if terms == nil {
		return nil
	}

	responses := make([]TermResponse, len(terms))
	for i, term := range terms {
		responses[i] = *r.ToTermResponse(term)
	}
	return responses
}

// OccupancyRuleInfo 佔用規則資訊（用於前端週曆顯示）
type OccupancyRuleInfo struct {
	RuleID       uint   `json:"rule_id"`
	OfferingID   uint   `json:"offering_id"`
	OfferingName string `json:"offering_name"`
	Weekday      int    `json:"weekday"`
	StartTime    string `json:"start_time"`
	EndTime      string `json:"end_time"`
	Duration     int    `json:"duration"`
	TeacherID    *uint  `json:"teacher_id,omitempty"`
	TeacherName  string `json:"teacher_name,omitempty"`
	RoomID       uint   `json:"room_id"`
	RoomName     string `json:"room_name,omitempty"`
	Status       string `json:"status"` // 狀態: PLANNED(預計), CONFIRMED(已開課), SUSPENDED(停課), ARCHIVED(歸檔)
}

// OccupancyRulesByDayOfWeek 按星期分組的佔用規則
type OccupancyRulesByDayOfWeek struct {
	DayOfWeek int                  `json:"day_of_week"` // 1-7 (週一到週日)
	DayName   string               `json:"day_name"`     // "週一", "週二", etc.
	Rules     []OccupancyRuleInfo  `json:"rules"`
}

// CopiedRuleInfo 複製規則結果資訊
type CopiedRuleInfo struct {
	OriginalRuleID uint   `json:"original_rule_id"`
	NewRuleID      uint   `json:"new_rule_id"`
	OfferingName   string `json:"offering_name"`
	Weekday        int    `json:"weekday"`
	StartTime      string `json:"start_time"`
	EndTime        string `json:"end_time"`
}

// CopyRulesResponse 複製規則響應
type CopyRulesResponse struct {
	CopiedCount int             `json:"copied_count"`
	Rules       []CopiedRuleInfo `json:"rules"`
}
