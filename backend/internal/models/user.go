package models

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	ID              uint           `json:"id" gorm:"primarykey"`
	Email           string         `json:"email" gorm:"type:varchar(255);uniqueIndex;not null" validate:"required,email"`
	Username        string         `json:"username" gorm:"type:varchar(255);uniqueIndex;not null" validate:"required,min=3,max=50"`
	PasswordHash    string         `json:"-" gorm:"not null"`
	FirstName       string         `json:"first_name" gorm:"not null" validate:"required,min=1,max=50"`
	LastName        string         `json:"last_name" gorm:"not null" validate:"required,min=1,max=50"`
	RoleID          uint           `json:"role_id" gorm:"not null"`
	IsActive        bool           `json:"is_active" gorm:"default:true"`
	EmailVerifiedAt *time.Time     `json:"email_verified_at"`
	LastLoginAt     *time.Time     `json:"last_login_at"`
	CreatedAt       time.Time      `json:"created_at"`
	UpdatedAt       time.Time      `json:"updated_at"`
	DeletedAt       gorm.DeletedAt `json:"-" gorm:"index"`

	Role         Role          `json:"role" gorm:"foreignKey:RoleID"`
	ActivityLogs []ActivityLog `json:"activity_logs,omitempty" gorm:"foreignKey:UserID"`
}

type UserInput struct {
	Email     string `json:"email" validate:"required,email"`
	Username  string `json:"username" validate:"required,min=3,max=50"`
	Password  string `json:"password" validate:"required,min=8,max=100"`
	FirstName string `json:"first_name" validate:"required,min=1,max=50"`
	LastName  string `json:"last_name" validate:"required,min=1,max=50"`
	RoleID    uint   `json:"role_id" validate:"required,min=1"`
}

type UserUpdateInput struct {
	Email     string `json:"email" validate:"omitempty,email"`
	Username  string `json:"username" validate:"omitempty,min=3,max=50"`
	FirstName string `json:"first_name" validate:"omitempty,min=1,max=50"`
	LastName  string `json:"last_name" validate:"omitempty,min=1,max=50"`
	RoleID    uint   `json:"role_id" validate:"omitempty,min=1"`
	IsActive  *bool  `json:"is_active"`
}

type PasswordUpdateInput struct {
	CurrentPassword string `json:"current_password" validate:"required,min=8"`
	NewPassword     string `json:"new_password" validate:"required,min=8,max=100"`
	ConfirmPassword string `json:"confirm_password" validate:"required,eqfield=NewPassword"`
}

type UserResponse struct {
	ID              uint       `json:"id"`
	Email           string     `json:"email"`
	Username        string     `json:"username"`
	FirstName       string     `json:"first_name"`
	LastName        string     `json:"last_name"`
	RoleID          uint       `json:"role_id"`
	IsActive        bool       `json:"is_active"`
	EmailVerifiedAt *time.Time `json:"email_verified_at"`
	LastLoginAt     *time.Time `json:"last_login_at"`
	CreatedAt       time.Time  `json:"created_at"`
	UpdatedAt       time.Time  `json:"updated_at"`
	Role            Role       `json:"role"`
	Permissions     []string   `json:"permissions"`
}

type UserListResponse struct {
	Users      []UserResponse `json:"users"`
	Total      int64          `json:"total"`
	Page       int            `json:"page"`
	Limit      int            `json:"limit"`
	TotalPages int            `json:"total_pages"`
}

func (User) TableName() string {
	return "users"
}

func (u *User) ToResponse() *UserResponse {
	permissions := []string{}
	// Ensure Role and Permissions are preloaded
	if u.Role.Permissions != nil {
		for _, p := range u.Role.Permissions {
			permissions = append(permissions, p.Name)
		}
	}

	return &UserResponse{
		ID:              u.ID,
		Email:           u.Email,
		Username:        u.Username,
		FirstName:       u.FirstName,
		LastName:        u.LastName,
		RoleID:          u.RoleID,
		IsActive:        u.IsActive,
		EmailVerifiedAt: u.EmailVerifiedAt,
		LastLoginAt:     u.LastLoginAt,
		CreatedAt:       u.CreatedAt,
		UpdatedAt:       u.UpdatedAt,
		Role:            u.Role,
		Permissions:     permissions,
	}
}