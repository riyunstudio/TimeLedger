package services

import (
	"context"

	jwt "timeLedger/libs/jwt"
)

type IdentInfo struct {
	ID         uint   `json:"id"`
	UserType   string `json:"user_type"` // ADMIN or TEACHER
	Name       string `json:"name"`
	Email      string `json:"email,omitempty"`
	LineUserID string `json:"line_user_id,omitempty"`
	CenterID   uint   `json:"center_id,omitempty"`
}

type LoginResponse struct {
	Token string    `json:"token"`
	User  IdentInfo `json:"user"`
}

type AuthService interface {
	// Admin Login
	AdminLogin(ctx context.Context, email, password string) (LoginResponse, error)

	// Teacher Login
	TeacherLineLogin(ctx context.Context, lineUserID, accessToken string) (LoginResponse, error)

	// Generate Token
	GenerateToken(claims jwt.Claims) (string, error)

	// Validate Token
	ValidateToken(token string) (*jwt.Claims, error)

	// Refresh Token (optional)
	RefreshToken(ctx context.Context, token string) (LoginResponse, error)

	// Logout
	Logout(ctx context.Context, token string) error
}
