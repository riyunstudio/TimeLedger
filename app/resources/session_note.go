package resources

import (
	"time"
)

type SessionNoteResource struct {
	ID          uint      `json:"id"`
	RuleID      *uint     `json:"rule_id"`
	SessionDate string    `json:"session_date"`
	Content     string    `json:"content"`
	PrepNote    string    `json:"prep_note"`
	UpdatedAt   time.Time `json:"updated_at"`
}
