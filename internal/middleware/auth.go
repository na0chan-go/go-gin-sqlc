package middleware

import (
	"net/http"
	"strings"

	"go-gin-sqlc/internal/util"

	"github.com/gin-gonic/gin"
)

// AuthRequired は認証を必要とするエンドポイントに使用するミドルウェアです
func AuthRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "認証ヘッダーがありません"})
			c.Abort()
			return
		}

		// Bearer トークンの形式を確認
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "無効な認証形式です"})
			c.Abort()
			return
		}

		// トークンの検証
		claims, err := util.ValidateToken(parts[1])
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "無効なトークンです"})
			c.Abort()
			return
		}

		// ユーザーIDをコンテキストに設定
		c.Set("userID", claims.UserID)
		c.Next()
	}
}
