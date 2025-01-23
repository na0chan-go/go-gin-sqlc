package handler

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	db "go-gin-sqlc/db/sqlc"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"golang.org/x/crypto/bcrypt"
)

// MockQueries はデータベースクエリのモックです
type MockQueries struct {
	mock.Mock
}

// GetUserByEmail はdb.Queriesインターフェースの実装です
func (m *MockQueries) GetUserByEmail(ctx context.Context, email string) (db.User, error) {
	args := m.Called(ctx, email)
	return args.Get(0).(db.User), args.Error(1)
}

// CreateUser はdb.Queriesインターフェースの実装です
func (m *MockQueries) CreateUser(ctx context.Context, arg db.CreateUserParams) (sql.Result, error) {
	args := m.Called(ctx, arg)
	return args.Get(0).(sql.Result), args.Error(1)
}

// GetUser はdb.Queriesインターフェースの実装です
func (m *MockQueries) GetUser(ctx context.Context, id int64) (db.User, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(db.User), args.Error(1)
}

// DeleteUser はdb.Queriesインターフェースの実装です
func (m *MockQueries) DeleteUser(ctx context.Context, id int64) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

// GetUsersByStatus はdb.Queriesインターフェースの実装です
func (m *MockQueries) GetUsersByStatus(ctx context.Context, arg db.GetUsersByStatusParams) ([]db.User, error) {
	args := m.Called(ctx, arg)
	return args.Get(0).([]db.User), args.Error(1)
}

// ListUsers はdb.Queriesインターフェースの実装です
func (m *MockQueries) ListUsers(ctx context.Context, arg db.ListUsersParams) ([]db.User, error) {
	args := m.Called(ctx, arg)
	return args.Get(0).([]db.User), args.Error(1)
}

// UpdateUser はdb.Queriesインターフェースの実装です
func (m *MockQueries) UpdateUser(ctx context.Context, arg db.UpdateUserParams) error {
	args := m.Called(ctx, arg)
	return args.Error(0)
}

// UpdateUserPassword はdb.Queriesインターフェースの実装です
func (m *MockQueries) UpdateUserPassword(ctx context.Context, arg db.UpdateUserPasswordParams) error {
	args := m.Called(ctx, arg)
	return args.Error(0)
}

func TestLogin(t *testing.T) {
	// Ginのテストモード設定
	gin.SetMode(gin.TestMode)

	// テスト用の時間を設定
	now := time.Now()

	// テスト用のパスワードハッシュを生成
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)
	if err != nil {
		t.Fatal("パスワードのハッシュ化に失敗しました:", err)
	}

	tests := []struct {
		name           string
		requestBody    LoginRequest
		setupMock      func(*MockQueries)
		expectedStatus int
		expectedError  string
	}{
		{
			name: "正常なログイン",
			requestBody: LoginRequest{
				Email:    "test@example.com",
				Password: "password123",
			},
			setupMock: func(m *MockQueries) {
				m.On("GetUserByEmail", mock.Anything, "test@example.com").Return(db.User{
					ID:           1,
					Email:        "test@example.com",
					PasswordHash: string(hashedPassword),
					FirstName:    "Test",
					LastName:     "User",
					Status:       db.NullUsersStatus{UsersStatus: db.UsersStatusActive, Valid: true},
					CreatedAt:    now,
					UpdatedAt:    now,
				}, nil)
			},
			expectedStatus: http.StatusOK,
		},
		{
			name: "無効なメールアドレス",
			requestBody: LoginRequest{
				Email:    "invalid-email",
				Password: "password123",
			},
			setupMock:      func(m *MockQueries) {},
			expectedStatus: http.StatusBadRequest,
			expectedError:  "Key: 'LoginRequest.Email' Error:Field validation for 'Email' failed on the 'email' tag",
		},
		{
			name: "パスワードなし",
			requestBody: LoginRequest{
				Email: "test@example.com",
			},
			setupMock:      func(m *MockQueries) {},
			expectedStatus: http.StatusBadRequest,
			expectedError:  "Key: 'LoginRequest.Password' Error:Field validation for 'Password' failed on the 'required' tag",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// モックの準備
			mockQueries := new(MockQueries)
			tt.setupMock(mockQueries)

			// ハンドラーの準備
			handler := &AuthHandler{
				queries: mockQueries,
			}

			// HTTPリクエストの準備
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)

			// リクエストボディの準備
			jsonData, _ := json.Marshal(tt.requestBody)
			c.Request = httptest.NewRequest(http.MethodPost, "/auth/login", bytes.NewBuffer(jsonData))
			c.Request.Header.Set("Content-Type", "application/json")

			// ハンドラーの実行
			handler.Login(c)

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
