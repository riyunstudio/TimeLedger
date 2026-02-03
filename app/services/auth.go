package services

import (
	"context"
	"errors"

	"timeLedger/app"
	"timeLedger/app/models"
	"timeLedger/app/repositories"
	jwt "timeLedger/libs/jwt"

	"golang.org/x/crypto/bcrypt"
)

type authService struct {
	BaseService
	app                 *app.App
	adminRepository     *repositories.AdminUserRepository
	teacherRepository   *repositories.TeacherRepository
	centerRepository    *repositories.CenterRepository
	notificationService NotificationQueueService
	jwt                 *jwt.JWT
}

func NewAuthService(app *app.App) *authService {
	return &authService{
		app:                 app,
		adminRepository:     repositories.NewAdminUserRepository(app),
		teacherRepository:   repositories.NewTeacherRepository(app),
		centerRepository:    repositories.NewCenterRepository(app),
		notificationService: NewNotificationQueueService(app),
		jwt:                 jwt.NewJWT(app.Env.JWTSecret),
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
		UserType: admin.Role,
		UserID:   admin.ID,
		CenterID: admin.CenterID,
	}

	token, err := s.GenerateToken(claims)
	if err != nil {
		return LoginResponse{}, err
	}

	// 觸發歡迎訊息（如果尚未發送且已綁定 LINE）
	if admin.LineUserID != "" {
		go s.triggerWelcomeAdminMessage(ctx, &admin, admin.CenterID)
	}

	return LoginResponse{
		Token: token,
		User: IdentInfo{
			ID:       admin.ID,
			UserType: admin.Role,
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

	// 觸發歡迎訊息（如果尚未發送且已綁定 LINE）
	if teacher.LineUserID != "" {
		go s.triggerWelcomeTeacherMessage(ctx, &teacher, centerID)
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

	if claims.UserType == "ADMIN" || claims.UserType == "OWNER" {
		admin, err := s.adminRepository.GetByID(ctx, claims.UserID)
		if err != nil {
			return LoginResponse{}, err
		}
		return LoginResponse{
			Token: token,
			User: IdentInfo{
				ID:       admin.ID,
				UserType: admin.Role,
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

// hasWelcomeMessageSent 檢查是否已發送過歡迎訊息
func (s *authService) hasWelcomeMessageSent(ctx context.Context, recipientID uint, recipientType string, messageType string) bool {
	queueRepo := repositories.NewNotificationQueueRepository(s.app)
	items, err := queueRepo.GetByRecipient(ctx, recipientID, recipientType, 100, 0)
	if err != nil {
		return false
	}

	for _, item := range items {
		if item.Type == messageType && item.Status == models.NotificationStatusSent {
			return true
		}
	}

	return false
}

// triggerWelcomeTeacherMessage 觸發老師歡迎訊息
func (s *authService) triggerWelcomeTeacherMessage(ctx context.Context, teacher *models.Teacher, centerID uint) {
	// 檢查是否已發送過歡迎訊息
	if s.hasWelcomeMessageSent(ctx, teacher.ID, "TEACHER", models.NotificationTypeWelcomeTeacher) {
		return
	}

	// 取得中心名稱
	center, err := s.centerRepository.GetByID(ctx, centerID)
	if err != nil {
		s.Logger.Warn("failed to get center for welcome message", "center_id", centerID, "error", err)
		return
	}

	// 發送歡迎訊息
	if err := s.notificationService.NotifyWelcomeTeacher(ctx, teacher, center.Name); err != nil {
		s.Logger.Error("failed to send welcome message to teacher", "teacher_id", teacher.ID, "error", err)
	}
}

// triggerWelcomeAdminMessage 觸發管理員歡迎訊息
func (s *authService) triggerWelcomeAdminMessage(ctx context.Context, admin *models.AdminUser, centerID uint) {
	// 檢查是否已發送過歡迎訊息
	if s.hasWelcomeMessageSent(ctx, admin.ID, "ADMIN", models.NotificationTypeWelcomeAdmin) {
		return
	}

	// 取得中心名稱
	center, err := s.centerRepository.GetByID(ctx, centerID)
	if err != nil {
		s.Logger.Warn("failed to get center for admin welcome message", "center_id", centerID, "error", err)
		return
	}

	// 發送歡迎訊息
	if err := s.notificationService.NotifyWelcomeAdmin(ctx, admin, center.Name); err != nil {
		s.Logger.Error("failed to send welcome message to admin", "admin_id", admin.ID, "error", err)
	}
}
