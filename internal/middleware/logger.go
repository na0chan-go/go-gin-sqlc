package middleware

import (
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
)

// Logger はリクエストのロギングを行うミドルウェアです
func Logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		// リクエスト開始時刻
		startTime := time.Now()

		// リクエストパス
		path := c.Request.URL.Path

		// リクエストメソッド
		method := c.Request.Method

		// ハンドラの処理を実行
		c.Next()

		// レスポンスステータス
		statusCode := c.Writer.Status()

		// 処理時間
		duration := time.Since(startTime)

		// ログ出力
		fmt.Printf("[%s] %s %s %d %v\n",
			startTime.Format("2006/01/02 - 15:04:05"),
			method,
			path,
			statusCode,
			duration,
		)
	}
}
