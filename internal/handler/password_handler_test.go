package handler

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	db "go-gin-sqlc/db/sqlc"
	"go-gin-sqlc/internal/config"
	"go-gin-sqlc/internal/util"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockSQLResult はsql.Resultのモックです
type MockSQLResult struct {
	mock.Mock
}

func (m *MockSQLResult) LastInsertId() (int64, error) {
	args := m.Called()
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockSQLResult) RowsAffected() (int64, error) {
	args := m.Called()
	return args.Get(0).(int64), args.Error(1)
}

// MockMailer はメール送信のモックです
type MockMailer struct {
	mock.Mock
}

func (m *MockMailer) SendMail(config util.MailConfig, to, subject, body string) error {
	args := m.Called(config, to, subject, body)
	return args.Error(0)
}

func TestRequestPasswordReset(t *testing.T) {
	// Ginのテストモード設定
	gin.SetMode(gin.TestMode)

	// テスト用の時間を設定
	now := time.Now()

	// オリジナルのSendMail関数を保存
	originalSendMail := util.SendMail
	defer func() {
		util.SendMail = originalSendMail
	}()

	tests := []struct {
		name           string
		requestBody    RequestPasswordResetRequest
		setupMock      func(*MockQueries)
		expectedStatus int
		expectedError  string
	}{
		{
			name: "正常なリクエスト",
			requestBody: RequestPasswordResetRequest{
				Email: "test@example.com",
			},
			setupMock: func(m *MockQueries) {
				m.On("GetUserByEmail", mock.Anything, "test@example.com").Return(db.User{
					ID:        1,
					Email:     "test@example.com",
					CreatedAt: now,
				}, nil)

				mockResult := new(MockSQLResult)
				mockResult.On("LastInsertId").Return(int64(1), nil)
				m.On("CreatePasswordReset", mock.Anything, mock.AnythingOfType("db.CreatePasswordResetParams")).Return(mockResult, nil)
			},
			expectedStatus: http.StatusOK,
		},
		{
			name: "無効なメールアドレス",
			requestBody: RequestPasswordResetRequest{
				Email: "invalid-email",
			},
			setupMock:      func(m *MockQueries) {},
			expectedStatus: http.StatusBadRequest,
			expectedError:  "Key: 'RequestPasswordResetRequest.Email' Error:Field validation for 'Email' failed on the 'email' tag",
		},
		{
			name: "存在しないユーザー",
			requestBody: RequestPasswordResetRequest{
				Email: "nonexistent@example.com",
			},
			setupMock: func(m *MockQueries) {
				m.On("GetUserByEmail", mock.Anything, "nonexistent@example.com").Return(db.User{}, sql.ErrNoRows)
			},
			expectedStatus: http.StatusOK, // セキュリティのため、成功レスポンスを返す
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// モックの準備
			mockQueries := new(MockQueries)
			tt.setupMock(mockQueries)

			// メール送信のモック
			util.SendMail = func(config util.MailConfig, to, subject, body string) error {
				return nil
			}

			// ハンドラーの準備
			handler := &PasswordHandler{
				queries: mockQueries,
				config: &config.Config{
					BaseURL: "http://localhost:8080",
					Mail: util.MailConfig{
						Host:     "smtp.example.com",
						Port:     25,
						Username: "test",
						Password: "test",
						From:     "noreply@example.com",
					},
				},
			}

			// HTTPリクエストの準備
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)

			// リクエストボディの準備
			jsonData, _ := json.Marshal(tt.requestBody)
			c.Request = httptest.NewRequest(http.MethodPost, "/passwords/reset-request", bytes.NewBuffer(jsonData))
			c.Request.Header.Set("Content-Type", "application/json")

			// ハンドラーの実行
			handler.RequestPasswordReset(c)

			// アサーション
			assert.Equal(t, tt.expectedStatus, w.Code)

			if tt.expectedError != "" {
				var response map[string]string
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Contains(t, response["error"], tt.expectedError)
			}

			// モックの検証
			mockQueries.AssertExpectations(t)
		})
	}
}

func TestResetPassword(t *testing.T) {
	// Ginのテストモード設定
	gin.SetMode(gin.TestMode)

	// テスト用の時間を設定
	now := time.Now()
	expiresAt := now.Add(24 * time.Hour)

	tests := []struct {
		name           string
		requestBody    ResetPasswordRequest
		setupMock      func(*MockQueries)
		expectedStatus int
		expectedError  string
	}{
		{
			name: "正常なリセット",
			requestBody: ResetPasswordRequest{
				Token:    "valid-token",
				Password: "newpassword123",
			},
			setupMock: func(m *MockQueries) {
				m.On("GetPasswordResetByToken", mock.Anything, "valid-token").Return(db.GetPasswordResetByTokenRow{
					UserID:    1,
					Token:     "valid-token",
					ExpiresAt: expiresAt,
				}, nil)

				m.On("UpdateUserPassword", mock.Anything, mock.AnythingOfType("db.UpdateUserPasswordParams")).Return(nil)
				m.On("DeletePasswordReset", mock.Anything, "valid-token").Return(nil)
			},
			expectedStatus: http.StatusOK,
		},
		{
			name: "無効なトークン",
			requestBody: ResetPasswordRequest{
				Token:    "invalid-token",
				Password: "newpassword123",
			},
			setupMock: func(m *MockQueries) {
				m.On("GetPasswordResetByToken", mock.Anything, "invalid-token").Return(db.GetPasswordResetByTokenRow{}, sql.ErrNoRows)
			},
			expectedStatus: http.StatusBadRequest,
			expectedError:  "無効なトークンです",
		},
		{
			name: "パスワードが短すぎる",
			requestBody: ResetPasswordRequest{
				Token:    "valid-token",
				Password: "short",
			},
			setupMock:      func(m *MockQueries) {},
			expectedStatus: http.StatusBadRequest,
			expectedError:  "Key: 'ResetPasswordRequest.Password' Error:Field validation for 'Password' failed on the 'min' tag",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// モックの準備
			mockQueries := new(MockQueries)
			tt.setupMock(mockQueries)

			// ハンドラーの準備
			handler := &PasswordHandler{
				queries: mockQueries,
				config:  &config.Config{},
			}

			// HTTPリクエストの準備
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)

			// リクエストボディの準備
			jsonData, _ := json.Marshal(tt.requestBody)
			c.Request = httptest.NewRequest(http.MethodPost, "/passwords/reset", bytes.NewBuffer(jsonData))
			c.Request.Header.Set("Content-Type", "application/json")

			// ハンドラーの実行
			handler.ResetPassword(c)

			// アサーション
			assert.Equal(t, tt.expectedStatus, w.Code)

			if tt.expectedError != "" {
				var response map[string]string
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Contains(t, response["error"], tt.expectedError)
			}

			// モックの検証
			mockQueries.AssertExpectations(t)
		})
	}
}
