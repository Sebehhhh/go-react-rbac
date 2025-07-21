package models

import (
	"time"

	"gorm.io/gorm"
)

type Permission struct {
	ID          uint           `json:"id" gorm:"primarykey"`
	Name        string         `json:"name" gorm:"type:varchar(255);uniqueIndex;not null" validate:"required,min=3,max=100"`
	Resource    string         `json:"resource" gorm:"not null" validate:"required,min=2,max=50"`
	Action      string         `json:"action" gorm:"not null" validate:"required,min=2,max=50"`
	Description string         `json:"description" gorm:"type:text"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`

	Roles []*Role `json:"roles,omitempty" gorm:"many2many:role_permissions;"`
}

type PermissionInput struct {
	Name        string `json:"name" validate:"required,min=3,max=100"`
	Resource    string `json:"resource" validate:"required,min=2,max=50"`
	Action      string `json:"action" validate:"required,min=2,max=50"`
	Description string `json:"description" validate:"max=500"`
}

func (Permission) TableName() string {
	return "permissions"
}