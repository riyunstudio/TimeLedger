package services

import (
	"context"
	"time"
)

type SmartMatchingService interface {
	FindMatches(ctx context.Context, centerID uint, teacherID *uint, roomID uint, startTime, endTime time.Time, requiredSkills []string, excludeTeacherIDs []uint) ([]MatchScore, error)
	SearchTalent(ctx context.Context, searchParams TalentSearchParams) ([]TalentResult, error)
}

type TalentSearchParams struct {
	City     string   `json:"city,omitempty"`
	District string   `json:"district,omitempty"`
	Keyword  string   `json:"keyword,omitempty"`
	Skills   []string `json:"skills,omitempty"`
	Hashtags []string `json:"hashtags,omitempty"`
}
