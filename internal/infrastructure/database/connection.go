package database

import (
	"database/sql"
	"fmt"

	"go-gin-sqlc/internal/config"

	_ "github.com/go-sql-driver/mysql"
)

// Connect はデータベース接続を確立します
func Connect(cfg config.DBConfig) (*sql.DB, error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true",
		cfg.User,
		cfg.Password,
		cfg.Host,
		cfg.Port,
		cfg.DBName,
	)

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, fmt.Errorf("データベース接続のオープンに失敗しました: %w", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("データベース接続の確認に失敗しました: %w", err)
	}

	return db, nil
}
