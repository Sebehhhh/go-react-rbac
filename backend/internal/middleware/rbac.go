package middleware

import (
	"fmt"

	"rbac-system/backend/internal/services"
	"rbac-system/backend/internal/utils"

	"github.com/gofiber/fiber/v2"
)

func RequirePermission(rbacService *services.RBACService, resource, action string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		userID := GetUserIDFromContext(c)
		if userID == 0 {
			return utils.SendError(c, fiber.StatusUnauthorized, "unauthorized", "User not authenticated")
		}

		hasPermission, err := rbacService.CheckPermission(userID, resource, action)
		if err != nil {
			return utils.SendError(c, fiber.StatusInternalServerError, "internal_error", "Error checking permissions")
		}

		if !hasPermission {
			return utils.SendError(c, fiber.StatusForbidden, "forbidden", "Insufficient permissions")
		}

		return c.Next()
	}
}

func RequireRole(rbacService *services.RBACService, roleName string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		userID := GetUserIDFromContext(c)
		if userID == 0 {
			return utils.SendError(c, fiber.StatusUnauthorized, "unauthorized", "User not authenticated")
		}

		hasRole, err := rbacService.HasRole(userID, roleName)
		if err != nil {
			return utils.SendError(c, fiber.StatusInternalServerError, "internal_error", "Error checking role")
		}

		if !hasRole {
			return utils.SendError(c, fiber.StatusForbidden, "forbidden", "Insufficient role privileges")
		}

		return c.Next()
	}
}

func RequireRoleAny(rbacService *services.RBACService, roleNames ...string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		userID := GetUserIDFromContext(c)
		if userID == 0 {
			return utils.SendError(c, fiber.StatusUnauthorized, "unauthorized", "User not authenticated")
		}

		userRole, err := rbacService.GetUserRole(userID)
		if err != nil {
			return utils.SendError(c, fiber.StatusInternalServerError, "internal_error", "Error getting user role")
		}

		for _, roleName := range roleNames {
			if userRole.Name == roleName {
				return c.Next()
			}
		}

		return utils.SendError(c, fiber.StatusForbidden, "forbidden", "Insufficient role privileges")
	}
}

func SelfOrPermission(rbacService *services.RBACService, resource, action string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		userID := GetUserIDFromContext(c)
		if userID == 0 {
			return utils.SendError(c, fiber.StatusUnauthorized, "unauthorized", "User not authenticated")
		}

		targetUserID := c.Params("id")
		if targetUserID == "" {
			return utils.SendError(c, fiber.StatusBadRequest, "bad_request", "User ID parameter is required")
		}

		if targetUserID == fmt.Sprintf("%d", userID) {
			return c.Next()
		}

		hasPermission, err := rbacService.CheckPermission(userID, resource, action)
		if err != nil {
			return utils.SendError(c, fiber.StatusInternalServerError, "internal_error", "Error checking permissions")
		}

		if !hasPermission {
			return utils.SendError(c, fiber.StatusForbidden, "forbidden", "Insufficient permissions")
		}

		return c.Next()
	}
}