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

func (m *MockQueries) CreatePasswordReset(ctx context.Context, arg db.CreatePasswordResetParams) (sql.Result, error) {
	args := m.Called(ctx, arg)
	return args.Get(0).(sql.Result), args.Error(1)
}

func (m *MockQueries) GetPasswordResetByToken(ctx context.Context, token string) (db.GetPasswordResetByTokenRow, error) {
	args := m.Called(ctx, token)
	return args.Get(0).(db.GetPasswordResetByTokenRow), args.Error(1)
}

func (m *MockQueries) DeletePasswordReset(ctx context.Context, token string) error {
	args := m.Called(ctx, token)
	return args.Error(0)
}

func (m *MockQueries) SearchUsers(ctx context.Context, arg db.SearchUsersParams) ([]db.User, error) {
	args := m.Called(ctx, arg)
	return args.Get(0).([]db.User), args.Error(1)
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

func TestRegister(t *testing.T) {
	// Ginのテストモード設定
	gin.SetMode(gin.TestMode)

	// テスト用の時間を設定
	now := time.Now()

	tests := []struct {
		name           string
		requestBody    RegisterRequest
		setupMock      func(*MockQueries)
		expectedStatus int
		expectedError  string
	}{
		{
			name: "正常な登録",
			requestBody: RegisterRequest{
				Email:     "test@example.com",
				Password:  "password123",
				FirstName: "Test",
				LastName:  "User",
			},
			setupMock: func(m *MockQueries) {
				m.On("GetUserByEmail", mock.Anything, "test@example.com").Return(db.User{}, sql.ErrNoRows)

				mockResult := new(MockSQLResult)
				mockResult.On("LastInsertId").Return(int64(1), nil)
				m.On("CreateUser", mock.Anything, mock.AnythingOfType("db.CreateUserParams")).Return(mockResult, nil)
			},
			expectedStatus: http.StatusCreated,
		},
		{
			name: "無効なメールアドレス",
			requestBody: RegisterRequest{
				Email:     "invalid-email",
				Password:  "password123",
				FirstName: "Test",
				LastName:  "User",
			},
			setupMock:      func(m *MockQueries) {},
			expectedStatus: http.StatusBadRequest,
			expectedError:  "Key: 'RegisterRequest.Email' Error:Field validation for 'Email' failed on the 'email' tag",
		},
		{
			name: "パスワードが短すぎる",
			requestBody: RegisterRequest{
				Email:     "test@example.com",
				Password:  "short",
				FirstName: "Test",
				LastName:  "User",
			},
			setupMock:      func(m *MockQueries) {},
			expectedStatus: http.StatusBadRequest,
			expectedError:  "Key: 'RegisterRequest.Password' Error:Field validation for 'Password' failed on the 'min' tag",
		},
		{
			name: "既存のメールアドレス",
			requestBody: RegisterRequest{
				Email:     "existing@example.com",
				Password:  "password123",
				FirstName: "Test",
				LastName:  "User",
			},
			setupMock: func(m *MockQueries) {
				m.On("GetUserByEmail", mock.Anything, "existing@example.com").Return(db.User{
					ID:        1,
					Email:     "existing@example.com",
					CreatedAt: now,
				}, nil)
			},
			expectedStatus: http.StatusBadRequest,
			expectedError:  "このメールアドレスは既に登録されています",
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
			c.Request = httptest.NewRequest(http.MethodPost, "/auth/register", bytes.NewBuffer(jsonData))
			c.Request.Header.Set("Content-Type", "application/json")

			// ハンドラーの実行
			handler.Register(c)

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

func TestLoginWithInvalidCredentials(t *testing.T) {
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
			name: "存在しないユーザー",
			requestBody: LoginRequest{
				Email:    "nonexistent@example.com",
				Password: "password123",
			},
			setupMock: func(m *MockQueries) {
				m.On("GetUserByEmail", mock.Anything, "nonexistent@example.com").Return(db.User{}, sql.ErrNoRows)
			},
			expectedStatus: http.StatusUnauthorized,
			expectedError:  "メールアドレスまたはパスワードが正しくありません",
		},
		{
			name: "無効なパスワード",
			requestBody: LoginRequest{
				Email:    "test@example.com",
				Password: "wrongpassword",
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
			expectedStatus: http.StatusUnauthorized,
			expectedError:  "メールアドレスまたはパスワードが正しくありません",
		},
		{
			name: "無効なユーザーステータス",
			requestBody: LoginRequest{
				Email:    "inactive@example.com",
				Password: "password123",
			},
			setupMock: func(m *MockQueries) {
				m.On("GetUserByEmail", mock.Anything, "inactive@example.com").Return(db.User{
					ID:           2,
					Email:        "inactive@example.com",
					PasswordHash: string(hashedPassword),
					FirstName:    "Inactive",
					LastName:     "User",
					Status:       db.NullUsersStatus{UsersStatus: db.UsersStatusInactive, Valid: true},
					CreatedAt:    now,
					UpdatedAt:    now,
				}, nil)
			},
			expectedStatus: http.StatusUnauthorized,
			expectedError:  "このアカウントは無効です",
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
