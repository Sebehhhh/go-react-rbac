package services

import (
	"fmt"

	"rbac-system/backend/internal/database"
	"rbac-system/backend/internal/models"
)

type RBACService struct{}

func NewRBACService() *RBACService {
	return &RBACService{}
}

func (s *RBACService) CheckPermission(userID uint, resource, action string) (bool, error) {
	var user models.User
	if err := database.DB.Preload("Role.Permissions").First(&user, userID).Error; err != nil {
		return false, err
	}

	requiredPermission := fmt.Sprintf("%s.%s", resource, action)
	
	for _, permission := range user.Role.Permissions {
		if permission.Name == requiredPermission {
			return true, nil
		}
		
		if permission.Resource == resource && permission.Action == action {
			return true, nil
		}
	}

	return false, nil
}

func (s *RBACService) GetUserPermissions(userID uint) ([]models.Permission, error) {
	var user models.User
	if err := database.DB.Preload("Role.Permissions").First(&user, userID).Error; err != nil {
		return nil, err
	}

	return user.Role.Permissions, nil
}

func (s *RBACService) HasRole(userID uint, roleName string) (bool, error) {
	var user models.User
	if err := database.DB.Preload("Role").First(&user, userID).Error; err != nil {
		return false, err
	}

	return user.Role.Name == roleName, nil
}

func (s *RBACService) GetUserRole(userID uint) (*models.Role, error) {
	var user models.User
	if err := database.DB.Preload("Role.Permissions").First(&user, userID).Error; err != nil {
		return nil, err
	}

	return &user.Role, nil
}

func (s *RBACService) CanManageUser(managerID, targetUserID uint) (bool, error) {
	managerRole, err := s.GetUserRole(managerID)
	if err != nil {
		return false, err
	}

	targetRole, err := s.GetUserRole(targetUserID)
	if err != nil {
		return false, err
	}

	roleHierarchy := map[string]int{
		"Super Admin": 4,
		"Admin":       3,
		"Manager":     2,
		"User":        1,
	}

	managerLevel, exists := roleHierarchy[managerRole.Name]
	if !exists {
		return false, nil
	}

	targetLevel, exists := roleHierarchy[targetRole.Name]
	if !exists {
		return false, nil
	}

	return managerLevel > targetLevel, nil
}