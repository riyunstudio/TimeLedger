package resources

import (
	"time"
	"timeLedger/app/models"
)

type SessionNoteResource struct {
	ID          uint      `json:"id"`
	RuleID      *uint     `json:"rule_id"`
	SessionDate string    `json:"session_date"`
	Content     string    `json:"content"`
	PrepNote    string    `json:"prep_note"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// TeacherNoteResource 老師評分與備註資源轉換
type TeacherNoteResource struct{}

// NewTeacherNoteResource 建立老師評分與備註資源轉換實例
func NewTeacherNoteResource() *TeacherNoteResource {
	return &TeacherNoteResource{}
}

// TeacherNoteResponse 老師評分與備註回應結構
type TeacherNoteResponse struct {
	ID           uint      `json:"id"`
	TeacherID    uint      `json:"teacher_id"`
	Rating       int       `json:"rating"`
	InternalNote string    `json:"internal_note"`
	CreatedAt    time.Time `json:"created_at,omitempty"`
	UpdatedAt    time.Time `json:"updated_at,omitempty"`
}

// ToTeacherNoteResponse 將模型轉換為回應
func (r *TeacherNoteResource) ToTeacherNoteResponse(note *models.CenterTeacherNote) *TeacherNoteResponse {
	return &TeacherNoteResponse{
		ID:           note.ID,
		TeacherID:    note.TeacherID,
		Rating:       note.Rating,
		InternalNote: note.InternalNote,
		CreatedAt:    note.CreatedAt,
		UpdatedAt:    note.UpdatedAt,
	}
}

// ToEmptyTeacherNoteResponse 建立空的回應（當沒有評分記錄時）
func (r *TeacherNoteResource) ToEmptyTeacherNoteResponse(teacherID uint) *TeacherNoteResponse {
	return &TeacherNoteResponse{
		TeacherID:    teacherID,
		Rating:       0,
		InternalNote: "",
	}
}
