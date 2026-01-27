package services

import (
	"context"
	"errors"

	"golang.org/x/crypto/bcrypt"
	"timeLedger/app"
	"timeLedger/app/repositories"
	jwt "timeLedger/libs/jwt"
)

type authService struct {
	BaseService
	app               *app.App
	adminRepository   *repositories.AdminUserRepository
	teacherRepository *repositories.TeacherRepository
	jwt               *jwt.JWT
}

func NewAuthService(app *app.App) *authService {
	return &authService{
		app:               app,
		adminRepository:   repositories.NewAdminUserRepository(app),
		teacherRepository: repositories.NewTeacherRepository(app),
		jwt:               jwt.NewJWT(app.Env.JWTSecret),
	}
}

func (s *authService) AdminLogin(ctx context.Context, email, password string) (LoginResponse, error) {
	admin, err := s.adminRepository.GetByEmail(ctx, email)
	if err != nil {
		return LoginResponse{}, errors.New("admin not found")
	}

	if admin.Status != "ACTIVE" {
		return LoginResponse{}, errors.New("admin account is inactive")
	}

	// 驗證 center_id 必須有效
	if admin.CenterID == 0 {
		return LoginResponse{}, errors.New("admin account is not associated with any center")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(admin.PasswordHash), []byte(password)); err != nil {
		return LoginResponse{}, errors.New("invalid password")
	}

	claims := jwt.Claims{
		UserType: "ADMIN",
		UserID:   admin.ID,
		CenterID: admin.CenterID,
	}

	token, err := s.GenerateToken(claims)
	if err != nil {
		return LoginResponse{}, err
	}

	return LoginResponse{
		Token: token,
		User: IdentInfo{
			ID:       admin.ID,
			UserType: "ADMIN",
			Name:     admin.Name,
			Email:    admin.Email,
			CenterID: admin.CenterID,
		},
	}, nil
}

func (s *authService) TeacherLineLogin(ctx context.Context, lineUserID, accessToken string) (LoginResponse, error) {
	teacher, err := s.teacherRepository.GetByLineUserID(ctx, lineUserID)
	if err != nil {
		return LoginResponse{}, errors.New("teacher not found")
	}

	// 取得老師所屬的中心 ID
	centerID, err := s.teacherRepository.GetCenterID(ctx, teacher.ID)
	if err != nil {
		centerID = 0 // 如果沒有會籍，設為 0
	}

	claims := jwt.Claims{
		UserType:   "TEACHER",
		UserID:     teacher.ID,
		CenterID:   centerID,
		LineUserID: teacher.LineUserID,
	}

	token, err := s.GenerateToken(claims)
	if err != nil {
		return LoginResponse{}, err
	}

	return LoginResponse{
		Token: token,
		User: IdentInfo{
			ID:         teacher.ID,
			UserType:   "TEACHER",
			Name:       teacher.Name,
			Email:      teacher.Email,
			LineUserID: teacher.LineUserID,
			CenterID:   centerID,
		},
	}, nil
}

func (s *authService) GenerateToken(claims jwt.Claims) (string, error) {
	return s.jwt.GenerateToken(claims)
}

func (s *authService) ValidateToken(tokenString string) (*jwt.Claims, error) {
	return s.jwt.ValidateToken(tokenString)
}

func (s *authService) RefreshToken(ctx context.Context, token string) (LoginResponse, error) {
	claims, err := s.ValidateToken(token)
	if err != nil {
		return LoginResponse{}, err
	}

	if claims.UserType == "ADMIN" {
		admin, err := s.adminRepository.GetByID(ctx, claims.UserID)
		if err != nil {
			return LoginResponse{}, err
		}
		return LoginResponse{
			Token: token,
			User: IdentInfo{
				ID:       admin.ID,
				UserType: "ADMIN",
				Name:     admin.Name,
				Email:    admin.Email,
				CenterID: admin.CenterID,
			},
		}, nil
	} else {
		teacher, err := s.teacherRepository.GetByID(ctx, claims.UserID)
		if err != nil {
			return LoginResponse{}, err
		}
		return LoginResponse{
			Token: token,
			User: IdentInfo{
				ID:         teacher.ID,
				UserType:   "TEACHER",
				Name:       teacher.Name,
				Email:      teacher.Email,
				LineUserID: teacher.LineUserID,
			},
		}, nil
	}
}

func (s *authService) Logout(ctx context.Context, token string) error {
	return nil
}
