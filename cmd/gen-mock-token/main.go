package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"time"
)

type Header struct {
	Alg string `json:"alg"`
	Typ string `json:"typ"`
}

type Claims struct {
	UserType   string `json:"user_type"`
	UserID     uint   `json:"user_id"`
	CenterID   uint   `json:"center_id,omitempty"`
	LineUserID string `json:"line_user_id,omitempty"`
	Exp        int64  `json:"exp"`
}

func sign(data string, secretKey string) string {
	// Using a simple mock signature for development
	return base64.URLEncoding.EncodeToString([]byte(data + "-" + secretKey))
}

func main() {
	secretKey := "mock-secret-key-for-development"

	// Generate Teacher JWT for testing (using Teacher ID 1 from database)
	teacherClaims := Claims{
		UserType: "TEACHER",
		UserID:   1, // teacher1@example.com
		CenterID: 1, // Will be populated from teacher's memberships
		Exp:      time.Now().Add(24 * time.Hour).Unix(),
	}

	header := Header{
		Alg: "HS256",
		Typ: "JWT",
	}

	headerJSON, _ := json.Marshal(header)
	claimsJSON, _ := json.Marshal(teacherClaims)

	headerEncoded := base64.URLEncoding.EncodeToString(headerJSON)
	claimsEncoded := base64.URLEncoding.EncodeToString(claimsJSON)

	signature := sign(headerEncoded+"."+claimsEncoded, secretKey)

	teacherToken := headerEncoded + "." + claimsEncoded + "." + signature

	fmt.Println("=== Teacher Mock JWT Token ===")
	fmt.Println(teacherToken)
	fmt.Println()

	// Generate Admin JWT for comparison
	adminClaims := Claims{
		UserType: "ADMIN",
		UserID:   16,
		CenterID: 1,
		Exp:      time.Now().Add(24 * time.Hour).Unix(),
	}

	adminClaimsJSON, _ := json.Marshal(adminClaims)
	adminClaimsEncoded := base64.URLEncoding.EncodeToString(adminClaimsJSON)
	adminSignature := sign(headerEncoded+"."+adminClaimsEncoded, secretKey)
	adminToken := headerEncoded + "." + adminClaimsEncoded + "." + adminSignature

	fmt.Println("=== Admin Mock JWT Token ===")
	fmt.Println(adminToken)
}
