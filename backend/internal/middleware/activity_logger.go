package middleware

import (
	"rbac-system/backend/internal/database"
	"rbac-system/backend/internal/models"

	"github.com/gofiber/fiber/v2"
)

func ActivityLogger() fiber.Handler {
	return func(c *fiber.Ctx) error {
		userID := GetUserIDFromContext(c)
		if userID == 0 {
			return c.Next()
		}

		err := c.Next()

		if c.Response().StatusCode() >= 200 && c.Response().StatusCode() < 300 {
			action := getActionFromMethod(c.Method())
			resource := getResourceFromPath(c.Path())
			
			activityLog := models.ActivityLog{
				UserID:    userID,
				Action:    action,
				Resource:  resource,
				Details:   getDetailsFromContext(c),
				IPAddress: c.IP(),
				UserAgent: c.Get("User-Agent"),
			}
			
			go func() {
				database.DB.Create(&activityLog)
			}()
		}

		return err
	}
}

func getActionFromMethod(method string) string {
	switch method {
	case "GET":
		return "read"
	case "POST":
		return "create"
	case "PUT", "PATCH":
		return "update"
	case "DELETE":
		return "delete"
	default:
		return "unknown"
	}
}

func getResourceFromPath(path string) string {
	if len(path) > 0 && path[0] == '/' {
		path = path[1:]
	}
	
	parts := []string{}
	current := ""
	
	for _, char := range path {
		if char == '/' {
			if current != "" {
				parts = append(parts, current)
				current = ""
			}
		} else {
			current += string(char)
		}
	}
	
	if current != "" {
		parts = append(parts, current)
	}
	
	if len(parts) >= 2 && parts[0] == "api" {
		return parts[1]
	}
	
	if len(parts) >= 1 {
		return parts[0]
	}
	
	return "unknown"
}

func getDetailsFromContext(c *fiber.Ctx) string {
	method := c.Method()
	path := c.Path()
	return method + " " + path
}