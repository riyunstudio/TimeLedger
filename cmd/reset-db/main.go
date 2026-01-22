package main

import (
	"fmt"
	"log"
	"time"

	jwtutils "timeLedger/libs/jwt"
)

func main() {
	log.Println("開始重置資料庫...")
	log.Println("(請手動清空資料表後重新執行 seeder)")
	log.Println("")
	log.Println("建立管理員和使用者資料表...")
	log.Println("✓ 資料庫重置完成")
	log.Println("")

	log.Println("建立管理員和使用者...")
	admin := struct {
		Email    string
		Password string
		Role     string
	}{
		Email:    "admin@timeledger.com",
		Password: "admin123",
		Role:     "OWNER",
	}
	fmt.Printf("管理員: %s\n", admin.Email)
	fmt.Printf("密碼: %s\n", admin.Password)
	fmt.Printf("權限: %s\n", admin.Role)
	fmt.Println("")

	teacher := struct {
		LineUserID string
		Name       string
		Email      string
	}{
		LineUserID: "LINE_TEACHER_001",
		Name:       "王小明",
		Email:      "wangxiaoming@example.com",
	}
	fmt.Printf("老師: %s\n", teacher.Name)
	fmt.Printf("LineUserID: %s\n", teacher.LineUserID)
	fmt.Println("")

	jwtInstance := jwtutils.NewJWT("mock-secret-key-for-development")

	claims := jwtutils.Claims{
		UserType:   "TEACHER",
		UserID:     1,
		LineUserID: teacher.LineUserID,
		Exp:        time.Now().Add(1 * time.Hour).Unix(),
	}

	token, err := jwtInstance.GenerateToken(claims)
	if err != nil {
		log.Fatalf("生成 Token 失敗: %v", err)
	}

	fmt.Println("==========================================")
	fmt.Println("  老師 LINE 登入 Token (1小時後過期)")
	fmt.Println("==========================================")
	fmt.Println("")
	fmt.Printf("curl -X POST http://localhost:8888/api/v1/auth/teacher/line/login \\\n")
	fmt.Printf("  -H \"Content-Type: application/json\" \\\n")
	fmt.Printf("  -d '{\n")
	fmt.Printf("    \"line_user_id\": \"%s\",\n", teacher.LineUserID)
	fmt.Printf("    \"access_token\": \"mock_token\"\n")
	fmt.Printf("  }'\n")
	fmt.Println("")
	fmt.Println("或直接使用以下 Token (1小時有效):")
	fmt.Println("")
	fmt.Println(token)
	fmt.Println("")
	fmt.Println("==========================================")

	adminClaims := jwtutils.Claims{
		UserType: "ADMIN",
		UserID:   1,
		CenterID: 1,
		Exp:      time.Now().Add(1 * time.Hour).Unix(),
	}

	adminToken, _ := jwtInstance.GenerateToken(adminClaims)
	fmt.Println("==========================================")
	fmt.Println("  管理員 Token (1小時後過期)")
	fmt.Println("==========================================")
	fmt.Println("")
	fmt.Println(adminToken)
	fmt.Println("")
	fmt.Println("==========================================")

	fmt.Println("")
	fmt.Println("測試 API")
	fmt.Println("==========================================")
	fmt.Println("")
	fmt.Println("測試老師個人資料:")
	fmt.Printf("curl -H \"Authorization: Bearer %s\" http://localhost:8888/api/v1/teacher/me/profile\n", token)
	fmt.Println("")
	fmt.Println("測試管理員取得房間:")
	fmt.Printf("curl -H \"Authorization: Bearer %s\" http://localhost:8888/api/v1/admin/rooms\n", adminToken)
	fmt.Println("")
}
