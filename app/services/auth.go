package services

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
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
	app                         *app.App
	adminRepository             *repositories.AdminUserRepository
	teacherRepository           *repositories.TeacherRepository
	centerRepository            *repositories.CenterRepository
	notificationService         NotificationQueueService
	adminLoginHistoryRepository *repositories.AdminLoginHistoryRepository
	jwt                         *jwt.JWT
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

// getLineEmail 嘗試從 LINE 取得用戶 email
// 使用 https://oauth2.googleapis.com/tokeninfo 或 LINE 的 ID Token
func (s *authService) getLineEmail(accessToken string) string {
	// 嘗試從 LINE ID Token 解碼取得 email
	// LINE 的 ID Token 是 JWT 格式，可以解碼取得包含的 claims
	parts := strings.Split(accessToken, ".")
	if len(parts) != 3 {
		return ""
	}

	// 解碼 JWT payload (base64url decode)
	payload, err := base64URLDecode(parts[1])
	if err != nil {
		return ""
	}

	// 解析 claims
	var claims struct {
		Email string `json:"email"`
	}
	if err := json.Unmarshal(payload, &claims); err != nil {
		return ""
	}

	return claims.Email
}

// base64URLDecode 解碼 base64url 格式的字串
func base64URLDecode(input string) ([]byte, error) {
	// 替換 base64url 特殊字元
	decoded := strings.ReplaceAll(input, "-", "+")
	decoded = strings.ReplaceAll(decoded, "_", "/")

	// 補足 padding
	switch len(decoded) % 4 {
	case 2:
		decoded += "=="
	case 3:
		decoded += "="
	}

	return base64.StdEncoding.DecodeString(decoded)
}

func NewAuthService(app *app.App) *authService {
	baseService := NewBaseService(app, "AuthService")
	svc := &authService{
		BaseService: *baseService,
		app:         app,
	}

	if app.Env != nil {
		svc.jwt = jwt.NewJWT(app.Env.JWTSecret)
	}

	if app.MySQL != nil {
		svc.adminRepository = repositories.NewAdminUserRepository(app)
		svc.teacherRepository = repositories.NewTeacherRepository(app)
		svc.centerRepository = repositories.NewCenterRepository(app)
		svc.notificationService = NewNotificationQueueService(app)
		svc.adminLoginHistoryRepository = repositories.NewAdminLoginHistoryRepository(app)
	}

	return svc
}

func (s *authService) AdminLogin(ctx context.Context, email, password, ipAddress, userAgent string) (LoginResponse, error) {
	admin, err := s.adminRepository.GetByEmail(ctx, email)
	if err != nil {
		// 記錄失敗：帳號不存在
		s.recordLoginHistory(ctx, 0, email, models.LoginStatusFailed, ipAddress, userAgent, "admin not found")
		return LoginResponse{}, errors.New("admin not found")
	}

	if admin.Status != "ACTIVE" {
		// 記錄失敗：帳號已停用
		s.recordLoginHistory(ctx, admin.ID, email, models.LoginStatusFailed, ipAddress, userAgent, "admin account is inactive")
		return LoginResponse{}, errors.New("admin account is inactive")
	}

	// 驗證 center_id 必須有效
	if admin.CenterID == 0 {
		// 記錄失敗：未關聯中心
		s.recordLoginHistory(ctx, admin.ID, email, models.LoginStatusFailed, ipAddress, userAgent, "admin account is not associated with any center")
		return LoginResponse{}, errors.New("admin account is not associated with any center")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(admin.PasswordHash), []byte(password)); err != nil {
		// 記錄失敗：密碼錯誤
		s.recordLoginHistory(ctx, admin.ID, email, models.LoginStatusFailed, ipAddress, userAgent, "invalid password")
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

	// 記錄成功登入
	s.recordLoginHistory(ctx, admin.ID, email, models.LoginStatusSuccess, ipAddress, userAgent, "")

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

// recordLoginHistory 記錄管理員登入紀錄
func (s *authService) recordLoginHistory(ctx context.Context, adminID uint, email, status, ipAddress, userAgent, reason string) {
	loginHistory := models.AdminLoginHistory{
		AdminID:   adminID,
		Email:     email,
		Status:    status,
		IPAddress: ipAddress,
		UserAgent: userAgent,
		Reason:    reason,
	}

	// 使用 goroutine 非同步寫入，避免影響登入流程
	go func() {
		_, err := s.adminLoginHistoryRepository.Create(ctx, loginHistory)
		if err != nil {
			// 記錄錯誤，不影響主要流程
			fmt.Printf("[WARN] failed to record admin login history: %v\n", err)
		}
	}()
}

func (s *authService) TeacherLineLogin(ctx context.Context, lineUserID, accessToken string) (*LoginResponse, *errInfos.Res, error) {
	// 驗證 LINE Access Token 並取得用戶資料
	lineProfile, err := s.getLineProfile(accessToken, lineUserID)
	if err != nil {
		return nil, s.app.Err.New(errInfos.UNAUTHORIZED), err
	}

	teacher, err := s.teacherRepository.GetByLineUserID(ctx, lineUserID)
	if err != nil {
		// 老師尚未註冊，自動建立老師帳號
		// 嘗試從 LINE 取得 email
		email := s.getLineEmail(accessToken)
		if email == "" {
			email = "teacher@timeledger.com" // 如果無法取得 email，使用預設值
		}

		newTeacher := models.Teacher{
			LineUserID: lineUserID,
			Email:      email,
			Name:       lineProfile.DisplayName, // 使用 LINE 的顯示名稱
			AvatarURL:  lineProfile.PictureURL,  // 使用 LINE 的頭像
		}

		createdTeacher, err := s.teacherRepository.Create(ctx, newTeacher)
		if err != nil {
			return nil, s.app.Err.New(errInfos.SQL_ERROR), fmt.Errorf("failed to create teacher: %w", err)
		}
		teacher = createdTeacher
		s.Logger.Info("auto-created teacher account via LINE login", "teacher_id", teacher.ID, "line_user_id", lineUserID)
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
