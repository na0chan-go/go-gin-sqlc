package util

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestGenerateAndValidateToken(t *testing.T) {
	// テストケースの定義
	tests := []struct {
		name    string
		userID  int64
		wantErr bool
	}{
		{
			name:    "正常なトークンの生成と検証",
			userID:  1,
			wantErr: false,
		},
		{
			name:    "ユーザーID 0 でのトークン生成と検証",
			userID:  0,
			wantErr: false,
		},
		{
			name:    "負のユーザーID でのトークン生成と検証",
			userID:  -1,
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// トークンの生成
			token, err := GenerateToken(tt.userID)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
			assert.NotEmpty(t, token)

			// トークンの検証
			claims, err := ValidateToken(token)
			assert.NoError(t, err)
			assert.NotNil(t, claims)
			assert.Equal(t, tt.userID, claims.UserID)

			// 有効期限の検証
			assert.True(t, claims.ExpiresAt.After(time.Now()))
		})
	}
}

func TestValidateToken_Invalid(t *testing.T) {
	// テストケースの定義
	tests := []struct {
		name    string
		token   string
		wantErr string
	}{
		{
			name:    "空のトークン",
			token:   "",
			wantErr: "token contains an invalid number of segments",
		},
		{
			name:    "不正なトークン",
			token:   "invalid.token.string",
			wantErr: "signature is invalid",
		},
		{
			name:    "不正な形式のトークン",
			token:   "invalid-token",
			wantErr: "token contains an invalid number of segments",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			claims, err := ValidateToken(tt.token)
			assert.Error(t, err)
			assert.Nil(t, claims)
			assert.Contains(t, err.Error(), tt.wantErr)
		})
	}
}
