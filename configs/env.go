package configs

import (
	"os"
	"strconv"
	"strings"
)

type Env struct {
	AppEnv      string
	AppDebug    bool
	AppTimezone string

	ServerName string
	ServerHost string
	ServerPort string

	PrefixCode int

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
}

func LoadEnv() *Env {
	return &Env{
		AppEnv:      os.Getenv("APP_ENV"),
		AppDebug:    getEnvAsBool("APP_DEBUG", false),
		AppTimezone: os.Getenv("APP_TIMEZONE"),

		ServerName: os.Getenv("SERVER_NAME"),
		ServerHost: os.Getenv("SERVER_HOST"),
		ServerPort: os.Getenv("SERVER_PORT"),

		PrefixCode: getEnvAsInt(os.Getenv("PREFIX_CODE"), -1),

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
