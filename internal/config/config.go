package config

import (
	"os"

	"go-gin-sqlc/internal/util"
)

// Config はアプリケーション全体の設定を保持します
type Config struct {
	DB      DBConfig
	Server  *ServerConfig
	Mail    util.MailConfig
	BaseURL string
}

// ServerConfig はサーバーの設定を保持します
type ServerConfig struct {
	Port string
}

type DBConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
}

// New は新しい設定インスタンスを作成します
func New() *Config {
	return &Config{
		DB: DBConfig{
			Host:     getEnv("DB_HOST", "localhost"),
			Port:     getEnv("DB_PORT", "3306"),
			User:     getEnv("DB_USER", "root"),
			Password: getEnv("DB_PASSWORD", "password"),
			DBName:   getEnv("DB_NAME", "go_gin_sqlc"),
		},
		Server: &ServerConfig{
			Port: getEnv("SERVER_PORT", "8080"),
		},
		Mail: util.MailConfig{
			Host:     getEnv("MAIL_HOST", "smtp.gmail.com"),
			Port:     25,
			Username: getEnv("MAIL_USERNAME", ""),
			Password: getEnv("MAIL_PASSWORD", ""),
			From:     getEnv("MAIL_FROM", "noreply@example.com"),
		},
		BaseURL: getEnv("BASE_URL", "http://localhost:8080"),
	}
}

// getEnv は環境変数を取得し、設定されていない場合はデフォルト値を返します
func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}
