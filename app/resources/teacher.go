package resources

import (
	"time"
)

type TeacherProfileResource struct {
	ID                uint   `json:"id,omitempty"`
	LineUserID        string `json:"line_user_id,omitempty"`
	Name              string `json:"name,omitempty"`
	Email             string `json:"email,omitempty"`
	Bio               string `json:"bio,omitempty"`
	City              string `json:"city,omitempty"`
	District          string `json:"district,omitempty"`
	PublicContactInfo string `json:"public_contact_info,omitempty"`
	IsOpenToHiring    bool   `json:"is_open_to_hiring,omitempty"`
}

type TeacherCertificateResource struct {
	ID        uint      `json:"id,omitempty"`
	Name      string    `json:"name,omitempty"`
	FileURL   string    `json:"file_url,omitempty"`
	IssuedAt  time.Time `json:"issued_at,omitempty"`
	CreatedAt time.Time `json:"created_at,omitempty"`
}

type CenterMembershipResource struct {
	ID         uint      `json:"id,omitempty"`
	CenterID   uint      `json:"center_id,omitempty"`
	CenterName string    `json:"center_name,omitempty"`
	Status     string    `json:"status,omitempty"`
	CreatedAt  time.Time `json:"created_at,omitempty"`
}
