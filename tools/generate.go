package tools

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	mathRand "math/rand"
	"time"

	"github.com/bwmarrin/snowflake"
	"github.com/google/uuid"
)

// 產生唯一 UserId
func (tl *Tools) NewUserId() string {
	return uuid.New().String()
}

// 產生登入 SessionId
func (tl *Tools) NewSessionId(UserId string) (string, error) {
	// 生成隨機數
	randomBytes := make([]byte, 16)
	_, err := rand.Read(randomBytes)

	if err != nil {
		return "", err
	}
	// UserId + 當前時間戳 + 隨機數
	data := fmt.Sprintf("%s%d%s", UserId, tl.NowUnix(), randomBytes)
	// 使用 SHA256
	hash := sha256.Sum256([]byte(data))
	// 用 base64 加密, 得出長度約為43位數
	sessionId := base64.StdEncoding.EncodeToString(hash[:])

	return sessionId, nil
}

// 使用者密碼加密
func (tl *Tools) PwdEncode(pwd string) string {
	hash := sha256.Sum256([]byte("sdgwetegtyuxcvxfdhg*!@@#hgjv" + pwd))
	// 用 base64 加密，得出長度約為44位數
	return base64.StdEncoding.EncodeToString(hash[:])
}

// 產生帳號
func (tl *Tools) NewAccount(prefix string) string {
	return prefix + tl.RandCharset(10)
}

// 產生指定長度的隨機字符串
func (tl *Tools) RandCharset(length int) string {
	const charset = "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	r := mathRand.New(mathRand.NewSource(time.Now().UnixNano()))
	b := make([]byte, length)

	for i := range b {
		b[i] = charset[r.Intn(len(charset))]
	}
	return string(b)
}

// 產生邀請碼
func (tl *Tools) NewInvitationCode() string {
	// 邀請碼 10碼
	return tl.RandCharset(10)
}

// 產生唯一 traceId
func (tl *Tools) NewTraceId() (traceId string, err error) {
	node, err := snowflake.NewNode(1)

	if err != nil {
		return "", err
	}
	traceId = node.Generate().String()

	return traceId, nil
}

// 產生特定長度的隨機 key (可能重複)
func (tl *Tools) NewApiKey(length int) string {
	switch length {
	case 16, 24, 32:
		break
	default:
		length = 32
	}

	const charset = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789!=@#_"
	r := mathRand.New(mathRand.NewSource(time.Now().UnixNano()))
	b := make([]byte, length)

	for i := range b {
		b[i] = charset[r.Intn(len(charset))]
	}
	return string(b)
}

// 產生6位數的信箱驗證碼
func (tl *Tools) NewMailCode() string {
	code := mathRand.New(mathRand.NewSource(time.Now().UnixNano())).Intn(999999) + 100000
	return fmt.Sprintf("%06d", code)
}

// 產生隨機長度Id
func (tl *Tools) NewRandLengthId() string {
	code := mathRand.New(mathRand.NewSource(time.Now().UnixNano())).Intn(999999999) + 100000000
	return fmt.Sprintf("%09d", code)
}
