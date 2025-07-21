package services

import (
	"errors"

	"rbac-system/backend/internal/database"
	"rbac-system/backend/internal/models"

	"gorm.io/gorm"
)

type RoleService struct{}

func NewRoleService() *RoleService {
	return &RoleService{}
}

func (s *RoleService) GetRoles() ([]models.Role, error) {
	var roles []models.Role
	if err := database.DB.Preload("Permissions").Find(&roles).Error; err != nil {
		return nil, err
	}
	return roles, nil
}

func (s *RoleService) GetRoleByID(id uint) (*models.Role, error) {
	var role models.Role
	if err := database.DB.Preload("Permissions").First(&role, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("role not found")
		}
		return nil, err
	}
	return &role, nil
}

func (s *RoleService) CreateRole(req *models.RoleInput) (*models.Role, error) {
	var existingRole models.Role
	if err := database.DB.Where("name = ?", req.Name).First(&existingRole).Error; err == nil {
		return nil, errors.New("role with this name already exists")
	}

	role := models.Role{
		Name:         req.Name,
		Description:  req.Description,
		IsSystemRole: false, // Only system can create system roles
	}

	if err := database.DB.Create(&role).Error; err != nil {
		return nil, err
	}

	// Assign permissions if provided
	if len(req.PermissionIDs) > 0 {
		var permissions []models.Permission
		if err := database.DB.Where("id IN ?", req.PermissionIDs).Find(&permissions).Error; err != nil {
			return nil, err
		}
		
		if err := database.DB.Model(&role).Association("Permissions").Append(&permissions); err != nil {
			return nil, err
		}
	}

	if err := database.DB.Preload("Permissions").First(&role, role.ID).Error; err != nil {
		return nil, err
	}

	return &role, nil
}

func (s *RoleService) UpdateRole(id uint, req *models.RoleInput) (*models.Role, error) {
	var role models.Role
	if err := database.DB.First(&role, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("role not found")
		}
		return nil, err
	}

	// System roles can be updated, but cannot change IsSystemRole field
	// Only Super Admin can modify certain system roles

	if req.Name != "" && req.Name != role.Name {
		// System roles cannot change name to prevent system confusion
		if role.IsSystemRole {
			return nil, errors.New("cannot change system role name")
		}
		
		var existingRole models.Role
		if err := database.DB.Where("name = ? AND id != ?", req.Name, id).First(&existingRole).Error; err == nil {
			return nil, errors.New("role with this name already exists")
		}
		role.Name = req.Name
	}

	if req.Description != "" {
		role.Description = req.Description
	}

	if err := database.DB.Save(&role).Error; err != nil {
		return nil, err
	}

	// Update permissions if provided
	if req.PermissionIDs != nil {
		var permissions []models.Permission
		if len(req.PermissionIDs) > 0 {
			if err := database.DB.Where("id IN ?", req.PermissionIDs).Find(&permissions).Error; err != nil {
				return nil, err
			}
		}
		
		if err := database.DB.Model(&role).Association("Permissions").Replace(&permissions); err != nil {
			return nil, err
		}
	}

	if err := database.DB.Preload("Permissions").First(&role, role.ID).Error; err != nil {
		return nil, err
	}

	return &role, nil
}

func (s *RoleService) DeleteRole(id uint) error {
	var role models.Role
	if err := database.DB.First(&role, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("role not found")
		}
		return err
	}

	if role.IsSystemRole {
		return errors.New("cannot delete system role")
	}

	var userCount int64
	if err := database.DB.Model(&models.User{}).Where("role_id = ?", id).Count(&userCount).Error; err != nil {
		return err
	}

	if userCount > 0 {
		return errors.New("cannot delete role that is assigned to users")
	}

	if err := database.DB.Model(&role).Association("Permissions").Clear(); err != nil {
		return err
	}

	return database.DB.Delete(&role).Error
}

func (s *RoleService) GetPermissions() ([]models.Permission, error) {
	var permissions []models.Permission
	if err := database.DB.Find(&permissions).Error; err != nil {
		return nil, err
	}
	return permissions, nil
}

func (s *RoleService) AssignPermissions(roleID uint, permissionIDs []uint) error {
	var role models.Role
	if err := database.DB.First(&role, roleID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("role not found")
		}
		return err
	}

	var permissions []models.Permission
	if err := database.DB.Where("id IN ?", permissionIDs).Find(&permissions).Error; err != nil {
		return err
	}

	if len(permissions) != len(permissionIDs) {
		return errors.New("some permissions not found")
	}

	if err := database.DB.Model(&role).Association("Permissions").Replace(&permissions); err != nil {
		return err
	}

	return nil
}

func (s *RoleService) GetRolePermissions(roleID uint) ([]*models.Permission, error) {
	var role models.Role
	if err := database.DB.Preload("Permissions").First(&role, roleID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("role not found")
		}
		return nil, err
	}

	return role.Permissions, nil
}