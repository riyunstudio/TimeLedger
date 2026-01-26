package main

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"time"
)

type Claims struct {
	UserType   string `json:"user_type"`
	UserID     uint   `json:"user_id"`
	CenterID   uint   `json:"center_id"`
	LineUserID string `json:"line_user_id,omitempty"`
	Exp        int64  `json:"exp"`
}

type Header struct {
	Alg string `json:"alg"`
	Typ string `json:"typ"`
}

func main() {
	// JWT Secret (從環境變量或使用預設值)
	secretKey := "timeLedger-jwt-secret-2024"

	// 生成教師 JWT
	claims := Claims{
		UserType: "TEACHER",
		UserID:   1, // 替換為實際的教師 ID
		CenterID: 1, // 替換為實際的中心 ID
		Exp:      time.Now().Add(24 * time.Hour).Unix(),
	}

	token := generateToken(secretKey, claims)
	fmt.Println("=== JWT Token ===")
	fmt.Println(token)
	fmt.Println("")
	fmt.Println("=== curl 測試指令 ===")
	fmt.Printf(`curl -X PATCH "http://localhost:8888/api/v1/teacher/me/personal-events/320260121" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer %s" \
  -d '{"title":"Test","start_at":"2026-01-22T09:00:00","end_at":"2026-01-22T10:00:00","update_mode":"SINGLE"}'
`, token)
}

func generateToken(secretKey string, claims Claims) string {
	header := Header{
		Alg: "HS256",
		Typ: "JWT",
	}

	headerJSON, _ := json.Marshal(header)
	claimsJSON, _ := json.Marshal(claims)

	headerEncoded := base64.URLEncoding.EncodeToString(headerJSON)
	claimsEncoded := base64.URLEncoding.EncodeToString(claimsJSON)

	signature := sign(secretKey, headerEncoded+"."+claimsEncoded)

	return headerEncoded + "." + claimsEncoded + "." + signature
}

func sign(secretKey, data string) string {
	h := hmac.New(sha256.New, []byte(secretKey))
	h.Write([]byte(data))
	return base64.URLEncoding.EncodeToString(h.Sum(nil))
}
