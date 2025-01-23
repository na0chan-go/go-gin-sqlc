package middleware

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"go-gin-sqlc/internal/util"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestAuthRequired(t *testing.T) {
	// Ginのテストモード設定
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name           string
		setupAuth      func(*http.Request)
		expectedStatus int
		expectedError  string
	}{
		{
			name: "有効なトークン",
			setupAuth: func(r *http.Request) {
				token, _ := util.GenerateToken(1)
				r.Header.Set("Authorization", "Bearer "+token)
			},
			expectedStatus: http.StatusOK,
		},
		{
			name: "認証ヘッダーなし",
			setupAuth: func(r *http.Request) {
			},
			expectedStatus: http.StatusUnauthorized,
			expectedError:  "認証ヘッダーがありません",
		},
		{
			name: "不正な認証形式",
			setupAuth: func(r *http.Request) {
				r.Header.Set("Authorization", "Invalid token")
			},
			expectedStatus: http.StatusUnauthorized,
			expectedError:  "無効な認証形式です",
		},
		{
			name: "不正なトークン",
			setupAuth: func(r *http.Request) {
				r.Header.Set("Authorization", "Bearer invalid.token.here")
			},
			expectedStatus: http.StatusUnauthorized,
			expectedError:  "無効なトークンです",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// テスト用のルーターとレスポンスレコーダーの設定
			w := httptest.NewRecorder()
			_, r := gin.CreateTestContext(w)

			// ミドルウェアとハンドラーの設定
			r.Use(AuthRequired())
			r.GET("/test", func(c *gin.Context) {
				c.Status(http.StatusOK)
			})

			// リクエストの準備
			req := httptest.NewRequest(http.MethodGet, "/test", nil)
			tt.setupAuth(req)

			// リクエストの実行
			r.ServeHTTP(w, req)

			// アサーション
			assert.Equal(t, tt.expectedStatus, w.Code)

			if tt.expectedError != "" {
				var response map[string]string
				err := json.NewDecoder(w.Body).Decode(&response)
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedError, response["error"])
			}
		})
	}
}
