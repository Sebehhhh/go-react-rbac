package services_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"rbac-system/backend/internal/models"
	"rbac-system/backend/internal/services"
)

func setupTestDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	assert.NoError(t, err)

	db.AutoMigrate(&models.User{}, &models.Role{}, &models.Permission{})

	return db
}

func TestRBACService_HasPermission(t *testing.T) {
	db := setupTestDB(t)
	service := services.NewRBACService(db)

	// Create roles
	adminRole := models.Role{Name: "Admin", Description: "Admin Role"}
	userRole := models.Role{Name: "User", Description: "User Role"}
	db.Create(&adminRole)
	db.Create(&userRole)

	// Create permissions
	userReadPerm := models.Permission{Name: "users.read", Resource: "users", Action: "read"}
	userCreatePerm := models.Permission{Name: "users.create", Resource: "users", Action: "create"}
	db.Create(&userReadPerm)
	db.Create(&userCreatePerm)

	// Assign permissions to roles
	db.Model(&adminRole).Association("Permissions").Append(&userReadPerm, &userCreatePerm)
	db.Model(&userRole).Association("Permissions").Append(&userReadPerm)

	// Create users
	adminUser := models.User{Email: "admin@example.com", Username: "adminuser", RoleID: adminRole.ID}
	regularUser := models.User{Email: "user@example.com", Username: "regularuser", RoleID: userRole.ID}
	db.Create(&adminUser)
	db.Create(&regularUser)

	// Test cases
	hasPerm, err := service.CheckPermission(adminUser.ID, "users", "read")
	assert.NoError(t, err)
	assert.True(t, hasPerm, "Admin should have users.read permission")

	hasPerm, err = service.CheckPermission(adminUser.ID, "users", "create")
	assert.NoError(t, err)
	assert.True(t, hasPerm, "Admin should have users.create permission")

	hasPerm, err = service.CheckPermission(adminUser.ID, "roles", "delete")
	assert.NoError(t, err)
	assert.False(t, hasPerm, "Admin should not have roles.delete permission")

	hasPerm, err = service.CheckPermission(regularUser.ID, "users", "read")
	assert.NoError(t, err)
	assert.True(t, hasPerm, "Regular user should have users.read permission")

	hasPerm, err = service.CheckPermission(regularUser.ID, "users", "create")
	assert.NoError(t, err)
	assert.False(t, hasPerm, "Regular user should not have users.create permission")
}

func TestRBACService_HasRole(t *testing.T) {
	db := setupTestDB(t)
	service := services.NewRBACService(db)

	// Create roles
	adminRole := models.Role{Name: "Admin", Description: "Admin Role"}
	userRole := models.Role{Name: "User", Description: "User Role"}
	db.Create(&adminRole)
	db.Create(&userRole)

	// Create users
	adminUser := models.User{Email: "admin@example.com", Username: "adminuser", RoleID: adminRole.ID}
	regularUser := models.User{Email: "user@example.com", Username: "regularuser", RoleID: userRole.ID}
	db.Create(&adminUser)
	db.Create(&regularUser)

	// Test cases
	hasRole, err := service.HasRole(adminUser.ID, "Admin")
	assert.NoError(t, err)
	assert.True(t, hasRole, "Admin user should have Admin role")

	hasRole, err = service.HasRole(adminUser.ID, "User")
	assert.NoError(t, err)
	assert.False(t, hasRole, "Admin user should not have User role")

	hasRole, err = service.HasRole(regularUser.ID, "User")
	assert.NoError(t, err)
	assert.True(t, hasRole, "Regular user should have User role")

	hasRole, err = service.HasRole(regularUser.ID, "Admin")
	assert.NoError(t, err)
	assert.False(t, hasRole, "Regular user should not have Admin role")
}
