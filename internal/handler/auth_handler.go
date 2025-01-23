package handler

import (
	"database/sql"
	"net/http"

	db "go-gin-sqlc/db/sqlc"
	"go-gin-sqlc/internal/util"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

type AuthHandler struct {
	queries *db.Queries
}

func NewAuthHandler(sqlDB *sql.DB) *AuthHandler {
	return &AuthHandler{
		queries: db.New(sqlDB),
	}
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type LoginResponse struct {
	Token string `json:"token"`
	User  struct {
		ID        int64  `json:"id"`
		Email     string `json:"email"`
		FirstName string `json:"first_name"`
		LastName  string `json:"last_name"`
	} `json:"user"`
}

// RegisterRoutes は認証関連のルートを登録します
func (h *AuthHandler) RegisterRoutes(r gin.IRouter) {
	auth := r.Group("/auth")
	{
		auth.POST("/login", h.Login)
	}
}

// Login はユーザーログインを処理します
func (h *AuthHandler) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// メールアドレスでユーザーを検索
	user, err := h.queries.GetUserByEmail(c, req.Email)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "メールアドレスまたはパスワードが正しくありません"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// パスワードの検証
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "メールアドレスまたはパスワードが正しくありません"})
		return
	}

	// JWTトークンの生成
	token, err := util.GenerateToken(user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "トークンの生成に失敗しました"})
		return
	}

	// レスポンスの作成
	response := LoginResponse{
		Token: token,
	}
	response.User.ID = user.ID
	response.User.Email = user.Email
	response.User.FirstName = user.FirstName
	response.User.LastName = user.LastName

	c.JSON(http.StatusOK, response)
}
