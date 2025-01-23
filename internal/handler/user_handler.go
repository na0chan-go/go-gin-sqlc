package handler

import (
	"database/sql"
	"net/http"
	"strconv"
	"strings"

	db "go-gin-sqlc/db/sqlc"
	"go-gin-sqlc/internal/handler/dto"

	"golang.org/x/crypto/bcrypt"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	queries db.Querier
}

func NewUserHandler(sqlDB *sql.DB) *UserHandler {
	return &UserHandler{
		queries: db.New(sqlDB),
	}
}

// RegisterRoutes はユーザー関連のルートを登録します
func (h *UserHandler) RegisterRoutes(r gin.IRouter) {
	users := r.Group("/users")
	{
		users.POST("", h.CreateUser)
		users.GET("", h.ListUsers)
		users.GET("/:id", h.GetUser)
		users.PUT("/:id", h.UpdateUser)
		users.DELETE("/:id", h.DeleteUser)
		users.GET("/search", h.SearchUsers)
	}
}

// CreateUser は新しいユーザーを作成します
func (h *UserHandler) CreateUser(c *gin.Context) {
	var req dto.CreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// パスワードのハッシュ化
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "パスワードのハッシュ化に失敗しました"})
		return
	}

	// ユーザーの作成
	result, err := h.queries.CreateUser(c, db.CreateUserParams{
		Email:        req.Email,
		PasswordHash: string(hashedPassword),
		FirstName:    req.FirstName,
		LastName:     req.LastName,
		Status:       db.NullUsersStatus{UsersStatus: db.UsersStatusActive, Valid: true},
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 作成されたユーザーIDの取得
	id, err := result.LastInsertId()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "ユーザーIDの取得に失敗しました"})
		return
	}

	// 作成されたユーザーの取得
	user, err := h.queries.GetUser(c, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "ユーザー情報の取得に失敗しました"})
		return
	}

	c.JSON(http.StatusCreated, toUserResponse(user))
}

// ListUsers はユーザー一覧を取得します
func (h *UserHandler) ListUsers(c *gin.Context) {
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))

	users, err := h.queries.ListUsers(c, db.ListUsersParams{
		Limit:  int32(limit),
		Offset: int32(offset),
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	response := dto.UsersResponse{
		Users: make([]dto.UserResponse, len(users)),
		Total: len(users),
	}

	for i, user := range users {
		response.Users[i] = toUserResponse(user)
	}

	c.JSON(http.StatusOK, response)
}

// GetUser は指定されたIDのユーザーを取得します
func (h *UserHandler) GetUser(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "無効なユーザーID"})
		return
	}

	user, err := h.queries.GetUser(c, id)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": "ユーザーが見つかりません"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, toUserResponse(user))
}

// UpdateUser は指定されたIDのユーザー情報を更新します
func (h *UserHandler) UpdateUser(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "無効なユーザーID"})
		return
	}

	var req dto.UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 現在のユーザー情報を取得
	currentUser, err := h.queries.GetUser(c, id)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": "ユーザーが見つかりません"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 更新パラメータの準備
	params := db.UpdateUserParams{
		ID:        id,
		Email:     req.Email,
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Status:    db.NullUsersStatus{UsersStatus: db.UsersStatus(req.Status), Valid: req.Status != ""},
	}

	// 空の値は現在の値を使用
	if params.Email == "" {
		params.Email = currentUser.Email
	}
	if params.FirstName == "" {
		params.FirstName = currentUser.FirstName
	}
	if params.LastName == "" {
		params.LastName = currentUser.LastName
	}

	// ユーザー情報の更新
	err = h.queries.UpdateUser(c, params)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 更新後のユーザー情報を取得
	updatedUser, err := h.queries.GetUser(c, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "更新後のユーザー情報の取得に失敗しました"})
		return
	}

	c.JSON(http.StatusOK, toUserResponse(updatedUser))
}

// DeleteUser は指定されたIDのユーザーを削除します
func (h *UserHandler) DeleteUser(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "無効なユーザーID"})
		return
	}

	err = h.queries.DeleteUser(c, id)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": "ユーザーが見つかりません"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "ユーザーを削除しました"})
}

// SearchUsers はユーザーを検索します
func (h *UserHandler) SearchUsers(c *gin.Context) {
	query := c.Query("q")
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))

	users, err := h.queries.ListUsers(c, db.ListUsersParams{
		Limit:  int32(limit),
		Offset: int32(offset),
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// クライアントサイドでフィルタリング
	filteredUsers := users
	if query != "" {
		filteredUsers = make([]db.User, 0)
		queryLower := strings.ToLower(query)
		for _, user := range users {
			if strings.Contains(strings.ToLower(user.Email), queryLower) ||
				strings.Contains(strings.ToLower(user.FirstName), queryLower) ||
				strings.Contains(strings.ToLower(user.LastName), queryLower) {
				filteredUsers = append(filteredUsers, user)
			}
		}
	}

	response := dto.UsersResponse{
		Users: make([]dto.UserResponse, len(filteredUsers)),
		Total: len(filteredUsers),
	}

	for i, user := range filteredUsers {
		response.Users[i] = toUserResponse(user)
	}

	c.JSON(http.StatusOK, response)
}

// toUserResponse はデータベースのユーザーモデルをレスポンス用の構造体に変換します
func toUserResponse(user db.User) dto.UserResponse {
	return dto.UserResponse{
		ID:        user.ID,
		Email:     user.Email,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Status:    string(user.Status.UsersStatus),
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}
}
