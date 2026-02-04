package configs

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
)

type Env struct {
	AppEnv      string
	AppDebug    bool
	AppTimezone string
	AppID       int

	ServerName string
	ServerHost string
	ServerPort string

	JWTSecret string

	MysqlMasterHost string
	MysqlMasterPort string
	MysqlMasterUser string
	MysqlMasterPass string
	MysqlMasterName string

	MysqlSlaveHost string
	MysqlSlavePort string
	MysqlSlaveUser string
	MysqlSlavePass string
	MysqlSlaveName string

	RedisHost     string
	RedisPort     string
	RedisPass     string
	RedisPoolSize int

	TelegramBotToken string
	TelegramChatID   string

	RabbitMqAccount  string
	RabbitMqPassword string
	RabbitMqHost     string
	RabbitMqPort     string

	WsServerPort string

	// LINE Messaging API
	LineChannelSecret      string
	LineChannelAccessToken string
	LineLiffID             string

	// Frontend
	FrontendBaseURL       string
	LineOfficialAccountID string

	// Rate Limiting
	RateLimitEnabled       bool
	RateLimitRequests      int
	RateLimitWindow        string
	RateLimitBlockDuration string

	// File Upload / Cloudflare R2
	UploadPath        string
	UploadMaxSize     int
	UploadAllowedExts []string

	// Cloudflare R2 Storage
	CloudflareR2Enabled    bool
	CloudflareR2AccountID  string
	CloudflareR2AccessKey  string
	CloudflareR2SecretKey  string
	CloudflareR2BucketName string
	CloudflareR2PublicURL  string
}

func LoadEnv() *Env {
	err := godotenv.Load("./.env")
	if err != nil {
		panic(fmt.Sprintf("讀取.env錯誤, err: %s", err.Error()))
	}

	return &Env{
		AppEnv:      os.Getenv("APP_ENV"),
		AppDebug:    getEnvAsBool("APP_DEBUG", false),
		AppTimezone: getEnvAsString("APP_TIMEZONE", "Asia/Taipei"),
		AppID:       getEnvAsInt("APP_ID", -1),

		ServerName: os.Getenv("SERVER_NAME"),
		ServerHost: os.Getenv("SERVER_HOST"),
		ServerPort: os.Getenv("SERVER_PORT"),

		JWTSecret: os.Getenv("JWT_SECRET"),

		MysqlMasterHost: os.Getenv("MYSQL_MASTER_HOST"),
		MysqlMasterPort: os.Getenv("MYSQL_MASTER_PORT"),
		MysqlMasterUser: os.Getenv("MYSQL_MASTER_USER"),
		MysqlMasterPass: os.Getenv("MYSQL_MASTER_PASS"),
		MysqlMasterName: os.Getenv("MYSQL_MASTER_NAME"),

		MysqlSlaveHost: os.Getenv("MYSQL_SLAVE_HOST"),
		MysqlSlavePort: os.Getenv("MYSQL_SLAVE_PORT"),
		MysqlSlaveUser: os.Getenv("MYSQL_SLAVE_USER"),
		MysqlSlavePass: os.Getenv("MYSQL_SLAVE_PASS"),
		MysqlSlaveName: os.Getenv("MYSQL_SLAVE_NAME"),

		RedisHost:     os.Getenv("REDIS_HOST"),
		RedisPort:     os.Getenv("REDIS_PORT"),
		RedisPass:     os.Getenv("REDIS_PASS"),
		RedisPoolSize: getEnvAsInt("REDIS_POOlSIZE", 0),

		TelegramBotToken: os.Getenv("TELEGRAM_BOT_TOKEN"),
		TelegramChatID:   os.Getenv("TELEGRAM_CHAT_ID"),

		RabbitMqAccount:  os.Getenv("RABBIT_MQ_ACCOUNT"),
		RabbitMqPassword: os.Getenv("RABBIT_MQ_PASSWORD"),
		RabbitMqHost:     os.Getenv("RABBIT_MQ_HOST"),
		RabbitMqPort:     os.Getenv("RABBIT_MQ_PORT"),

		WsServerPort: os.Getenv("WS_SERVER_PORT"),

		// LINE Messaging API
		LineChannelSecret:      os.Getenv("LINE_CHANNEL_SECRET"),
		LineChannelAccessToken: os.Getenv("LINE_CHANNEL_ACCESS_TOKEN"),
		LineLiffID:             os.Getenv("LIFF_ID"),

		// Frontend
		FrontendBaseURL:       os.Getenv("FRONTEND_BASE_URL"),
		LineOfficialAccountID: os.Getenv("LINE_OFFICIAL_ACCOUNT_ID"),

		// Rate Limiting
		RateLimitEnabled:       getEnvAsBool("RATE_LIMIT_ENABLED", true),
		RateLimitRequests:      getEnvAsInt("RATE_LIMIT_REQUESTS", 100),
		RateLimitWindow:        getEnvAsString("RATE_LIMIT_WINDOW", "1m"),
		RateLimitBlockDuration: getEnvAsString("RATE_LIMIT_BLOCK_DURATION", "5m"),

		// File Upload
		UploadPath:        getEnvAsString("UPLOAD_PATH", "./uploads"),
		UploadMaxSize:     getEnvAsInt("UPLOAD_MAX_SIZE", 10), // 10MB
		UploadAllowedExts: []string{".jpg", ".jpeg", ".png", ".pdf"},

		// Cloudflare R2 Storage
		CloudflareR2Enabled:    getEnvAsBool("CLOUDFLARE_R2_ENABLED", false),
		CloudflareR2AccountID:  os.Getenv("CLOUDFLARE_R2_ACCOUNT_ID"),
		CloudflareR2AccessKey:  os.Getenv("CLOUDFLARE_R2_ACCESS_KEY"),
		CloudflareR2SecretKey:  os.Getenv("CLOUDFLARE_R2_SECRET_KEY"),
		CloudflareR2BucketName: os.Getenv("CLOUDFLARE_R2_BUCKET_NAME"),
		CloudflareR2PublicURL:  os.Getenv("CLOUDFLARE_R2_PUBLIC_URL"),
	}
}

func getEnvAsBool(name string, defaultVal bool) bool {
	valStr := os.Getenv(name)
	if valStr == "" {
		return defaultVal
	}

	valStr = strings.ToLower(valStr)
	return valStr == "true" || valStr == "1"
}

func getEnvAsString(name string, defaultVal string) string {
	valStr := os.Getenv(name)
	if valStr == "" {
		return defaultVal
	}
	return valStr
}

func getEnvAsInt(name string, defaultVal int) int {
	valStr := os.Getenv(name)
	if valStr == "" {
		return defaultVal
	}

	num, err := strconv.Atoi(valStr)
	if err != nil {
		return 0
	}
	return num
}
