package resources

import (
	"time"
)

type SessionNoteResource struct {
	ID          uint      `json:"id,omitempty"`
	RuleID      *uint     `json:"rule_id,omitempty"`
	SessionDate string    `json:"session_date,omitempty"`
	Content     string    `json:"content,omitempty"`
	PrepNote    string    `json:"prep_note,omitempty"`
	UpdatedAt   time.Time `json:"updated_at,omitempty"`
}
