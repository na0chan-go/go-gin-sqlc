package main

import (
	"fmt"
	"log"

	"go-gin-sqlc/internal/config"
	"go-gin-sqlc/internal/handler"
	"go-gin-sqlc/internal/infrastructure/database"
	"go-gin-sqlc/internal/middleware"

	"github.com/gin-gonic/gin"
)

func main() {
	// 設定の読み込み
	cfg := config.New()

	// データベース接続
	db, err := database.Connect(cfg.DB)
	if err != nil {
		log.Fatal("データベース接続の確立に失敗しました:", err)
	}
	defer db.Close()

	// Ginルーターの初期化
	r := gin.Default()

	// ミドルウェアの適用
	r.Use(middleware.Logger())

	// パブリックルート
	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "Welcome to Go-Gin-SQLC API",
		})
	})

	r.GET("/health", func(c *gin.Context) {
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

	// 認証ハンドラーの初期化と登録
	authHandler := handler.NewAuthHandler(db)
	authHandler.RegisterRoutes(r)

	// パスワードリセットハンドラーの初期化と登録
	passwordHandler := handler.NewPasswordHandler(db, cfg)
	passwordHandler.RegisterRoutes(r)

	// 認証が必要なルート
	authorized := r.Group("/api")
	authorized.Use(middleware.AuthRequired())
	{
		// ユーザーハンドラーの初期化と登録
		userHandler := handler.NewUserHandler(db)
		userHandler.RegisterRoutes(authorized)
	}

	// サーバーの起動
	addr := fmt.Sprintf(":%s", cfg.Server.Port)
	if err := r.Run(addr); err != nil {
		log.Fatal("サーバーの起動に失敗しました:", err)
	}
}
