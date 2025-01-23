package handler

import (
	"crypto/rand"
	"database/sql"
	"encoding/hex"
	"net/http"
	"time"

	db "go-gin-sqlc/db/sqlc"
	"go-gin-sqlc/internal/config"
	"go-gin-sqlc/internal/util"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

type PasswordHandler struct {
	queries db.Querier
	config  *config.Config
	mailer  util.Mailer
}

func NewPasswordHandler(sqlDB *sql.DB, cfg *config.Config) *PasswordHandler {
	return &PasswordHandler{
		queries: db.New(sqlDB),
		config:  cfg,
		mailer:  util.NewSMTPMailer(cfg.Mail),
	}
}

type RequestPasswordResetRequest struct {
	Email string `json:"email" binding:"required,email"`
}

type ResetPasswordRequest struct {
	Token    string `json:"token" binding:"required"`
	Password string `json:"password" binding:"required,min=8"`
}

// RegisterRoutes はパスワードリセット関連のルートを登録します
func (h *PasswordHandler) RegisterRoutes(r gin.IRouter) {
	passwords := r.Group("/passwords")
	{
		passwords.POST("/reset-request", h.RequestPasswordReset)
		passwords.POST("/reset", h.ResetPassword)
	}
}

// RequestPasswordReset はパスワードリセットのリクエストを処理します
func (h *PasswordHandler) RequestPasswordReset(c *gin.Context) {
	var req RequestPasswordResetRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// ユーザーの存在確認
	user, err := h.queries.GetUserByEmail(c, req.Email)
	if err != nil {
		if err == sql.ErrNoRows {
			// セキュリティのため、ユーザーが存在しない場合でも成功レスポンスを返す
			c.JSON(http.StatusOK, gin.H{"message": "パスワードリセットメールを送信しました"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// トークンの生成
	token := make([]byte, 32)
	if _, err := rand.Read(token); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "トークンの生成に失敗しました"})
		return
	}
	tokenStr := hex.EncodeToString(token)

	// パスワードリセットレコードの作成
	expiresAt := time.Now().Add(24 * time.Hour)
	_, err = h.queries.CreatePasswordReset(c, db.CreatePasswordResetParams{
		UserID:    user.ID,
		Token:     tokenStr,
		ExpiresAt: expiresAt,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// リセットURLの生成
	resetURL := h.config.BaseURL + "/reset-password?token=" + tokenStr

	// メールの送信
	mailBody := util.GeneratePasswordResetEmail(resetURL)
	err = h.mailer.SendMail(user.Email, "パスワードリセットのリクエスト", mailBody)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "メールの送信に失敗しました"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "パスワードリセットメールを送信しました"})
}

// ResetPassword はパスワードのリセットを処理します
func (h *PasswordHandler) ResetPassword(c *gin.Context) {
	var req ResetPasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// トークンの検証
	reset, err := h.queries.GetPasswordResetByToken(c, req.Token)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusBadRequest, gin.H{"error": "無効なトークンです"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// パスワードのハッシュ化
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "パスワードのハッシュ化に失敗しました"})
		return
	}

	// パスワードの更新
	err = h.queries.UpdateUserPassword(c, db.UpdateUserPasswordParams{
		ID:           reset.UserID,
		PasswordHash: string(hashedPassword),
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 使用済みトークンの削除
	err = h.queries.DeletePasswordReset(c, req.Token)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "パスワードを更新しました"})
}
