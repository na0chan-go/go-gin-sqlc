package config

import (
	"os"

	"go-gin-sqlc/internal/infrastructure/database"
)

// Config はアプリケーション全体の設定を保持します
type Config struct {
	DB     *database.Config
	Server *ServerConfig
}

// ServerConfig はサーバーの設定を保持します
type ServerConfig struct {
	Port string
}

// New は新しい設定インスタンスを作成します
func New() *Config {
	return &Config{
		DB: &database.Config{
			Host:     getEnv("DB_HOST", "localhost"),
			Port:     getEnv("DB_PORT", "3306"),
			User:     getEnv("DB_USER", "user"),
			Password: getEnv("DB_PASSWORD", "password"),
			DBName:   getEnv("DB_NAME", "go_gin_db"),
		},
		Server: &ServerConfig{
			Port: getEnv("SERVER_PORT", "8080"),
		},
	}
}

// getEnv は環境変数を取得し、設定されていない場合はデフォルト値を返します
func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}
