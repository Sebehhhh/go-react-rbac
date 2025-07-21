package utils

import (
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

type ErrorResponse struct {
	Error   string                 `json:"error"`
	Message string                 `json:"message,omitempty"`
	Details map[string]interface{} `json:"details,omitempty"`
}

type SuccessResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
}

func SendError(c *fiber.Ctx, status int, err string, message ...string) error {
	response := ErrorResponse{
		Error: err,
	}
	
	if len(message) > 0 {
		response.Message = message[0]
	}
	
	return c.Status(status).JSON(response)
}

func SendSuccess(c *fiber.Ctx, status int, message string, data interface{}) error {
	response := SuccessResponse{
		Success: true,
		Message: message,
		Data:    data,
	}
	
	return c.Status(status).JSON(response)
}

func SendValidationError(c *fiber.Ctx, err error) error {
	if validationErrors, ok := err.(validator.ValidationErrors); ok {
		details := make(map[string]interface{})
		
		for _, validationError := range validationErrors {
			field := strings.ToLower(validationError.Field())
			switch validationError.Tag() {
			case "required":
				details[field] = "This field is required"
			case "email":
				details[field] = "Invalid email format"
			case "min":
				details[field] = "Value is too short"
			case "max":
				details[field] = "Value is too long"
			case "eqfield":
				details[field] = "Field must match with " + validationError.Param()
			default:
				details[field] = "Invalid value"
			}
		}
		
		return c.Status(fiber.StatusBadRequest).JSON(ErrorResponse{
			Error:   "validation_failed",
			Message: "Validation failed",
			Details: details,
		})
	}
	
	return SendError(c, fiber.StatusBadRequest, "validation_failed", "Invalid input data")
}