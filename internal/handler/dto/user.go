package dto

import "time"

// CreateUserRequest はユーザー作成リクエストの構造体です
type CreateUserRequest struct {
	Email     string `json:"email" binding:"required,email"`
	Password  string `json:"password" binding:"required,min=8"`
	FirstName string `json:"first_name" binding:"required"`
	LastName  string `json:"last_name" binding:"required"`
}

// UpdateUserRequest はユーザー更新リクエストの構造体です
type UpdateUserRequest struct {
	Email     string `json:"email" binding:"omitempty,email"`
	FirstName string `json:"first_name" binding:"omitempty"`
	LastName  string `json:"last_name" binding:"omitempty"`
	Status    string `json:"status" binding:"omitempty,oneof=active inactive suspended"`
}

// UserResponse はユーザー情報レスポンスの構造体です
type UserResponse struct {
	ID        int64     `json:"id"`
	Email     string    `json:"email"`
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// UsersResponse は複数ユーザー情報のレスポンス構造体です
type UsersResponse struct {
	Users []UserResponse `json:"users"`
	Total int            `json:"total"`
}
