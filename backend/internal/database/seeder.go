package database

import (
	"log"
	"time"

	"rbac-system/backend/internal/config"
	"rbac-system/backend/internal/models"
	"rbac-system/backend/internal/utils"
)

func SeedDatabase(cfg *config.Config) error {
	if err := seedPermissions(); err != nil {
		return err
	}

	if err := seedRoles(); err != nil {
		return err
	}

	if err := seedDefaultAdmin(cfg); err != nil {
		return err
	}

	log.Println("Database seeding completed successfully")
	return nil
}

func seedPermissions() error {
	permissions := []models.Permission{
		{Name: "users.create", Resource: "users", Action: "create", Description: "Create new users"},
		{Name: "users.read", Resource: "users", Action: "read", Description: "View users"},
		{Name: "users.update", Resource: "users", Action: "update", Description: "Update user information"},
		{Name: "users.delete", Resource: "users", Action: "delete", Description: "Delete users"},
		{Name: "roles.create", Resource: "roles", Action: "create", Description: "Create new roles"},
		{Name: "roles.read", Resource: "roles", Action: "read", Description: "View roles"},
		{Name: "roles.update", Resource: "roles", Action: "update", Description: "Update role information"},
		{Name: "roles.delete", Resource: "roles", Action: "delete", Description: "Delete roles"},
		{Name: "permissions.read", Resource: "permissions", Action: "read", Description: "View permissions"},
		{Name: "profile.read", Resource: "profile", Action: "read", Description: "View own profile"},
		{Name: "profile.update", Resource: "profile", Action: "update", Description: "Update own profile"},
		{Name: "dashboard.read", Resource: "dashboard", Action: "read", Description: "View dashboard"},
		{Name: "activity_logs.read", Resource: "activity_logs", Action: "read", Description: "View activity logs"},
	}

	for _, permission := range permissions {
		var existingPermission models.Permission
		if err := DB.Where("name = ?", permission.Name).First(&existingPermission).Error; err != nil {
			if err := DB.Create(&permission).Error; err != nil {
				return err
			}
			log.Printf("Created permission: %s", permission.Name)
		}
	}

	return nil
}

func seedRoles() error {
	roles := []models.Role{
		{Name: "Super Admin", Description: "Full system access", IsSystemRole: true},
		{Name: "Admin", Description: "Administrative access", IsSystemRole: true},
		{Name: "Manager", Description: "Manager access with limited admin privileges", IsSystemRole: true},
		{Name: "User", Description: "Basic user access", IsSystemRole: true},
	}

	rolePermissions := map[string][]string{
		"Super Admin": {
			"users.create", "users.read", "users.update", "users.delete",
			"roles.create", "roles.read", "roles.update", "roles.delete",
			"permissions.read", "profile.read", "profile.update",
			"dashboard.read", "activity_logs.read",
		},
		"Admin": {
			"users.create", "users.read", "users.update", "users.delete",
			"roles.read", "permissions.read", "profile.read", "profile.update",
			"dashboard.read", "activity_logs.read",
		},
		"Manager": {
			"users.read", "users.update", "roles.read", "permissions.read",
			"profile.read", "profile.update", "dashboard.read",
		},
		"User": {
			"profile.read", "profile.update",
		},
	}

	for _, role := range roles {
		var existingRole models.Role
		if err := DB.Where("name = ?", role.Name).First(&existingRole).Error; err != nil {
			if err := DB.Create(&role).Error; err != nil {
				return err
			}

			var permissions []models.Permission
			permissionNames := rolePermissions[role.Name]
			if err := DB.Where("name IN ?", permissionNames).Find(&permissions).Error; err != nil {
				return err
			}

			if err := DB.Model(&role).Association("Permissions").Append(&permissions); err != nil {
				return err
			}

			log.Printf("Created role: %s with %d permissions", role.Name, len(permissions))
		}
	}

	return nil
}

func seedDefaultAdmin(cfg *config.Config) error {
	if cfg.System.DefaultAdminEmail == "" || cfg.System.DefaultAdminPassword == "" {
		log.Println("Skipping default admin creation - credentials not provided")
		return nil
	}

	var existingAdmin models.User
	if err := DB.Where("email = ?", cfg.System.DefaultAdminEmail).First(&existingAdmin).Error; err == nil {
		log.Println("Default admin already exists")
		return nil
	}

	var superAdminRole models.Role
	if err := DB.Where("name = ?", "Super Admin").First(&superAdminRole).Error; err != nil {
		return err
	}

	hashedPassword, err := utils.HashPassword(cfg.System.DefaultAdminPassword)
	if err != nil {
		return err
	}

	now := time.Now()
	admin := models.User{
		Email:           cfg.System.DefaultAdminEmail,
		Username:        "admin",
		PasswordHash:    hashedPassword,
		FirstName:       "System",
		LastName:        "Administrator",
		RoleID:          superAdminRole.ID,
		IsActive:        true,
		EmailVerifiedAt: &now,
	}

	if err := DB.Create(&admin).Error; err != nil {
		return err
	}

	log.Printf("Created default admin user: %s", admin.Email)
	return nil
}