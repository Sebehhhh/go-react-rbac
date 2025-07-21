package middleware

import (
	"strings"

	"rbac-system/backend/internal/models"
	"rbac-system/backend/internal/services"
	"rbac-system/backend/internal/utils"

	"github.com/gofiber/fiber/v2"
)

func AuthMiddleware(authService *services.AuthService, jwtService *utils.JWTService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return utils.SendError(c, fiber.StatusUnauthorized, "unauthorized", "Authorization header is required")
		}

		tokenParts := strings.SplitN(authHeader, " ", 2)
		if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
			return utils.SendError(c, fiber.StatusUnauthorized, "unauthorized", "Invalid authorization header format")
		}

		tokenString := tokenParts[1]
		claims, err := jwtService.ValidateAccessToken(tokenString)
		if err != nil {
			return utils.SendError(c, fiber.StatusUnauthorized, "unauthorized", "Invalid or expired token")
		}

		user, err := authService.GetUserByID(claims.UserID)
		if err != nil {
			return utils.SendError(c, fiber.StatusUnauthorized, "unauthorized", "User not found")
		}

		if !user.IsActive {
			return utils.SendError(c, fiber.StatusUnauthorized, "unauthorized", "Account is deactivated")
		}

		c.Locals("user", user)
		c.Locals("user_id", user.ID)
		c.Locals("role_id", user.RoleID)
		
		return c.Next()
	}
}

func GetUserFromContext(c *fiber.Ctx) *models.User {
	user, ok := c.Locals("user").(*models.User)
	if !ok {
		return nil
	}
	return user
}

func GetUserIDFromContext(c *fiber.Ctx) uint {
	userID, ok := c.Locals("user_id").(uint)
	if !ok {
		return 0
	}
	return userID
}

func GetRoleIDFromContext(c *fiber.Ctx) uint {
	roleID, ok := c.Locals("role_id").(uint)
	if !ok {
		return 0
	}
	return roleID
}