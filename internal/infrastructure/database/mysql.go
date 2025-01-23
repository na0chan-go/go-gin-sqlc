package database

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

type Config struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
}

// NewConnection はMySQLデータベースへの新しい接続を確立します
func NewConnection(config *Config) (*sql.DB, error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true&loc=Local",
		config.User,
		config.Password,
		config.Host,
		config.Port,
		config.DBName,
	)

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, fmt.Errorf("データベース接続のオープンに失敗しました: %v", err)
	}

	// 接続設定
	db.SetMaxOpenConns(25)                 // 最大接続数
	db.SetMaxIdleConns(25)                 // アイドル状態の最大接続数
	db.SetConnMaxLifetime(5 * time.Minute) // 接続の最大生存時間

	// 接続テスト
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("データベース接続のテストに失敗しました: %v", err)
	}

	log.Println("データベースへの接続に成功しました")
	return db, nil
}
