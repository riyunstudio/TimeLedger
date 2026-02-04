package services

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"timeLedger/app"
	"timeLedger/app/models"
	"timeLedger/app/repositories"
	"timeLedger/global/errInfos"
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

// LineProfile LINE API 回傳的用戶資料結構
type LineProfile struct {
	UserID        string `json:"userId"`
	DisplayName   string `json:"displayName"`
	PictureURL    string `json:"pictureUrl"`
	StatusMessage string `json:"statusMessage,omitempty"`
}

// getLineProfile 驗證 LINE Access Token 並取得用戶資料
// 使用 https://api.line.me/v2/profile API 驗證
func (s *authService) getLineProfile(accessToken string, lineUserID string) (*LineProfile, error) {
	// 建立 HTTP Client，設定超時時間
	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	// 建立 LINE Profile API 請求
	req, err := http.NewRequest("GET", "https://api.line.me/v2/profile", nil)
	if err != nil {
		return nil, errors.New("failed to create LINE API request")
	}

	// 設定 Authorization header
	req.Header.Set("Authorization", "Bearer "+accessToken)

	// 發送請求
	resp, err := client.Do(req)
	if err != nil {
		return nil, errors.New("failed to connect LINE API: " + err.Error())
	}
	defer resp.Body.Close()

	// 檢查 HTTP 狀態碼
	if resp.StatusCode != http.StatusOK {
		if resp.StatusCode == http.StatusUnauthorized {
			return nil, errors.New("LINE token 已過期或無效，請重新登入")
		}
		return nil, errors.New("LINE API 驗證失敗，狀態碼: " + fmt.Sprintf("%d", resp.StatusCode))
	}

	// 解析回應
	var profile LineProfile
	if err := json.NewDecoder(resp.Body).Decode(&profile); err != nil {
		return nil, errors.New("failed to parse LINE API response")
	}

	// 驗證用戶 ID 是否匹配
	if profile.UserID != lineUserID {
		return nil, errors.New("LINE user ID mismatch: token does not belong to the specified user")
	}

	return &profile, nil
}

// verifyLineToken 驗證 LINE Access Token
func (s *authService) verifyLineToken(accessToken string, lineUserID string) error {
	_, err := s.getLineProfile(accessToken, lineUserID)
	return err
}

func NewAuthService(app *app.App) *authService {
	baseService := NewBaseService(app, "AuthService")
	return &authService{
		BaseService:         *baseService,
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
	// 使用 context.Background() 避免 HTTP 請求結束後 context 失效
	if admin.LineUserID != "" {
		go s.triggerWelcomeAdminMessage(context.Background(), &admin, admin.CenterID)
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

func (s *authService) TeacherLineLogin(ctx context.Context, lineUserID, accessToken string) (*LoginResponse, *errInfos.Res, error) {
	// 驗證 LINE Access Token
	if err := s.verifyLineToken(accessToken, lineUserID); err != nil {
		return nil, s.app.Err.New(errInfos.UNAUTHORIZED), err
	}

	teacher, err := s.teacherRepository.GetByLineUserID(ctx, lineUserID)
	if err != nil {
		// 老師尚未註冊，返回特定錯誤碼讓前端可以引導用戶註冊
		return nil, s.app.Err.New(errInfos.TEACHER_NOT_REGISTERED), fmt.Errorf("teacher not registered")
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
		return nil, s.app.Err.New(errInfos.SYSTEM_ERROR), err
	}

	// 觸發歡迎訊息（如果尚未發送且已綁定 LINE）
	// 使用 context.Background() 避免 HTTP 請求結束後 context 失效
	if teacher.LineUserID != "" {
		go s.triggerWelcomeTeacherMessage(context.Background(), &teacher, centerID)
	}

	return &LoginResponse{
		Token: token,
		User: IdentInfo{
			ID:         teacher.ID,
			UserType:   "TEACHER",
			Name:       teacher.Name,
			Email:      teacher.Email,
			LineUserID: teacher.LineUserID,
			CenterID:   centerID,
		},
	}, nil, nil
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
	// 防禦性檢查：確保 Logger 已初始化
	if s.Logger == nil {
		fmt.Printf("[WARN] Logger not initialized, teacher welcome message skipped (teacher_id: %d)\n", teacher.ID)
		return
	}

	// centerID 為 0 表示老師沒有隸屬任何中心，跳過歡迎訊息
	if centerID == 0 {
		s.Logger.Debug("teacher has no center membership, skipping welcome message", "teacher_id", teacher.ID)
		return
	}

	s.Logger.Debug("triggering teacher welcome message", "teacher_id", teacher.ID, "center_id", centerID)

	// 檢查是否已發送過歡迎訊息
	if s.hasWelcomeMessageSent(ctx, teacher.ID, "TEACHER", models.NotificationTypeWelcomeTeacher) {
		return
	}

	// 取得中心名稱
	center, err := s.centerRepository.GetByID(ctx, centerID)
	if err != nil {
		// 防禦性檢查：Logger 可能為 nil
		if s.Logger != nil {
			s.Logger.Warn("failed to get center for welcome message", "center_id", centerID, "error", err)
		} else {
			fmt.Printf("[WARN] failed to get center for welcome message (center_id: %d, error: %v)\n", centerID, err)
		}
		return
	}

	// 發送歡迎訊息
	if err := s.notificationService.NotifyWelcomeTeacher(ctx, teacher, center.Name); err != nil {
		// 防禦性檢查：Logger 可能為 nil
		if s.Logger != nil {
			s.Logger.Error("failed to send welcome message to teacher", "teacher_id", teacher.ID, "error", err)
		}
	}
}

// triggerWelcomeAdminMessage 觸發管理員歡迎訊息
func (s *authService) triggerWelcomeAdminMessage(ctx context.Context, admin *models.AdminUser, centerID uint) {
	// 防禦性檢查：確保 Logger 已初始化
	if s.Logger == nil {
		fmt.Printf("[WARN] Logger not initialized, admin welcome message skipped (admin_id: %d)\n", admin.ID)
		return
	}

	// centerID 為 0 表示管理員沒有隸屬任何中心，跳過歡迎訊息
	if centerID == 0 {
		s.Logger.Debug("admin has no center membership, skipping welcome message", "admin_id", admin.ID)
		return
	}

	s.Logger.Debug("triggering admin welcome message", "admin_id", admin.ID, "center_id", centerID)

	// 檢查是否已發送過歡迎訊息
	if s.hasWelcomeMessageSent(ctx, admin.ID, "ADMIN", models.NotificationTypeWelcomeAdmin) {
		return
	}

	// 取得中心名稱
	center, err := s.centerRepository.GetByID(ctx, centerID)
	if err != nil {
		// 防禦性檢查：Logger 可能為 nil
		if s.Logger != nil {
			s.Logger.Warn("failed to get center for admin welcome message", "center_id", centerID, "error", err)
		} else {
			fmt.Printf("[WARN] failed to get center for admin welcome message (center_id: %d, error: %v)\n", centerID, err)
		}
		return
	}

	// 發送歡迎訊息
	if err := s.notificationService.NotifyWelcomeAdmin(ctx, admin, center.Name); err != nil {
		// 防禦性檢查：Logger 可能為 nil
		if s.Logger != nil {
			s.Logger.Error("failed to send welcome message to admin", "admin_id", admin.ID, "error", err)
		}
	}
}
