package models

import (
	"time"

	"gorm.io/gorm"
)

type Role struct {
	ID           uint           `json:"id" gorm:"primarykey"`
	Name         string         `json:"name" gorm:"type:varchar(50);uniqueIndex;not null" validate:"required,min=2,max=50"`
	Description  string         `json:"description" gorm:"type:text"`
	IsSystemRole bool           `json:"is_system_role" gorm:"default:false"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `json:"-" gorm:"index"`

	Users       []User        `json:"users,omitempty" gorm:"foreignKey:RoleID"`
	Permissions []*Permission `json:"permissions,omitempty" gorm:"many2many:role_permissions;"`
}

type RoleInput struct {
	Name          string  `json:"name" validate:"required,min=2,max=50"`
	Description   string  `json:"description" validate:"max=500"`
	IsSystemRole  bool    `json:"is_system_role"`
	PermissionIDs []uint  `json:"permission_ids"`
}

type RoleWithPermissions struct {
	Role
	PermissionIDs []uint `json:"permission_ids"`
}

type RolePermissionInput struct {
	PermissionIDs []uint `json:"permission_ids" validate:"required"`
}

func (Role) TableName() string {
	return "roles"
}