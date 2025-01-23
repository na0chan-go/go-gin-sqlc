package main

import (
	"fmt"
	"log"

	"go-gin-sqlc/internal/config"
	"go-gin-sqlc/internal/handler"
	"go-gin-sqlc/internal/infrastructure/database"

	"github.com/gin-gonic/gin"
)

func main() {
	// 設定の読み込み
	cfg := config.New()

	// データベース接続
	db, err := database.NewConnection(cfg.DB)
	if err != nil {
		log.Fatal("データベース接続の確立に失敗しました:", err)
	}
	defer db.Close()

	// Ginルーターの初期化
	r := gin.Default()

	// ユーザーハンドラーの初期化と登録
	userHandler := handler.NewUserHandler(db)
	userHandler.RegisterRoutes(r)

	// ルートエンドポイント
	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "Welcome to Go-Gin-SQLC API",
		})
	})

	// ヘルスチェックエンドポイント
	r.GET("/health", func(c *gin.Context) {
		// データベース接続の確認
		err := db.Ping()
		if err != nil {
			c.JSON(500, gin.H{
				"status":  "error",
				"message": "データベース接続エラー",
			})
			return
		}

		c.JSON(200, gin.H{
			"status":  "ok",
			"message": "サービスは正常に動作しています",
		})
	})

	// サーバーの起動
	addr := fmt.Sprintf(":%s", cfg.Server.Port)
	if err := r.Run(addr); err != nil {
		log.Fatal("サーバーの起動に失敗しました:", err)
	}
}
