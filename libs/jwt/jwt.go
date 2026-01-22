package jwtutils

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"time"
)

type Claims struct {
	UserType   string `json:"user_type"`
	UserID     uint   `json:"user_id"`
	CenterID   uint   `json:"center_id,omitempty"`
	LineUserID string `json:"line_user_id,omitempty"`
	Exp        int64  `json:"exp"`
}

type Header struct {
	Alg string `json:"alg"`
	Typ string `json:"typ"`
}

type JWT struct {
	secretKey []byte
}

func NewJWT(secretKey string) *JWT {
	return &JWT{
		secretKey: []byte(secretKey),
	}
}

func (j *JWT) GenerateToken(claims Claims) (string, error) {
	claims.Exp = time.Now().Add(1 * time.Hour).Unix()

	header := Header{
		Alg: "HS256",
		Typ: "JWT",
	}

	headerJSON, err := json.Marshal(header)
	if err != nil {
		return "", err
	}

	claimsJSON, err := json.Marshal(claims)
	if err != nil {
		return "", err
	}

	headerEncoded := base64.URLEncoding.EncodeToString(headerJSON)
	claimsEncoded := base64.URLEncoding.EncodeToString(claimsJSON)

	signature := j.sign(headerEncoded + "." + claimsEncoded)

	return headerEncoded + "." + claimsEncoded + "." + signature, nil
}

func (j *JWT) ValidateToken(tokenString string) (*Claims, error) {
	parts := splitToken(tokenString)
	if len(parts) != 3 {
		return nil, errors.New("invalid token format")
	}

	signature := j.sign(parts[0] + "." + parts[1])
	if signature != parts[2] {
		return nil, errors.New("invalid signature")
	}

	claimsJSON, err := base64.URLEncoding.DecodeString(parts[1])
	if err != nil {
		return nil, errors.New("invalid claims encoding")
	}

	var claims Claims
	if err := json.Unmarshal(claimsJSON, &claims); err != nil {
		return nil, errors.New("invalid claims")
	}

	if time.Now().Unix() > claims.Exp {
		return nil, errors.New("token expired")
	}

	return &claims, nil
}

func (j *JWT) sign(data string) string {
	h := hmac.New(sha256.New, j.secretKey)
	h.Write([]byte(data))
	return base64.URLEncoding.EncodeToString(h.Sum(nil))
}

func splitToken(token string) []string {
	var parts []string
	start := 0
	for i, c := range token {
		if c == '.' {
			parts = append(parts, token[start:i])
			start = i + 1
		}
	}
	parts = append(parts, token[start:])
	return parts
}
