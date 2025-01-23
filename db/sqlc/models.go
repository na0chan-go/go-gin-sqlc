// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.28.0

package db

import (
	"database/sql"
	"database/sql/driver"
	"fmt"
	"time"
)

type UsersStatus string

const (
	UsersStatusActive    UsersStatus = "active"
	UsersStatusInactive  UsersStatus = "inactive"
	UsersStatusSuspended UsersStatus = "suspended"
)

func (e *UsersStatus) Scan(src interface{}) error {
	switch s := src.(type) {
	case []byte:
		*e = UsersStatus(s)
	case string:
		*e = UsersStatus(s)
	default:
		return fmt.Errorf("unsupported scan type for UsersStatus: %T", src)
	}
	return nil
}

type NullUsersStatus struct {
	UsersStatus UsersStatus `json:"users_status"`
	Valid       bool        `json:"valid"` // Valid is true if UsersStatus is not NULL
}

// Scan implements the Scanner interface.
func (ns *NullUsersStatus) Scan(value interface{}) error {
	if value == nil {
		ns.UsersStatus, ns.Valid = "", false
		return nil
	}
	ns.Valid = true
	return ns.UsersStatus.Scan(value)
}

// Value implements the driver Valuer interface.
func (ns NullUsersStatus) Value() (driver.Value, error) {
	if !ns.Valid {
		return nil, nil
	}
	return string(ns.UsersStatus), nil
}

type PasswordReset struct {
	ID        int64        `json:"id"`
	UserID    int64        `json:"user_id"`
	Token     string       `json:"token"`
	ExpiresAt time.Time    `json:"expires_at"`
	CreatedAt sql.NullTime `json:"created_at"`
}

type User struct {
	ID           int64           `json:"id"`
	Email        string          `json:"email"`
	PasswordHash string          `json:"password_hash"`
	FirstName    string          `json:"first_name"`
	LastName     string          `json:"last_name"`
	Status       NullUsersStatus `json:"status"`
	CreatedAt    time.Time       `json:"created_at"`
	UpdatedAt    time.Time       `json:"updated_at"`
}
