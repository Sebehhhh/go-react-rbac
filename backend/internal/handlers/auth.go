package handlers

import (
	"rbac-system/backend/internal/middleware"
	"rbac-system/backend/internal/models"
	"rbac-system/backend/internal/services"
	"rbac-system/backend/internal/utils"

	"github.com/gofiber/fiber/v2"
)

type AuthHandler struct {
	authService *services.AuthService
}

func NewAuthHandler(authService *services.AuthService) *AuthHandler {
	return &AuthHandler{
		authService: authService,
	}
}

func (h *AuthHandler) Register(c *fiber.Ctx) error {
	var req models.RegisterRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.SendError(c, fiber.StatusBadRequest, "invalid_request", "Invalid request body")
	}

	if err := utils.ValidateStruct(&req); err != nil {
		return utils.SendValidationError(c, err)
	}

	response, err := h.authService.Register(&req)
	if err != nil {
		return utils.SendError(c, fiber.StatusBadRequest, "registration_failed", err.Error())
	}

	return utils.SendSuccess(c, fiber.StatusCreated, "Registration successful", response)
}

func (h *AuthHandler) Login(c *fiber.Ctx) error {
	var req models.LoginRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.SendError(c, fiber.StatusBadRequest, "invalid_request", "Invalid request body")
	}

	if err := utils.ValidateStruct(&req); err != nil {
		return utils.SendValidationError(c, err)
	}

	ipAddress := c.IP()
	userAgent := c.Get("User-Agent")

	response, err := h.authService.Login(&req, ipAddress, userAgent)
	if err != nil {
		return utils.SendError(c, fiber.StatusUnauthorized, "login_failed", err.Error())
	}

	return utils.SendSuccess(c, fiber.StatusOK, "Login successful", response)
}

func (h *AuthHandler) RefreshToken(c *fiber.Ctx) error {
	var req models.RefreshTokenRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.SendError(c, fiber.StatusBadRequest, "invalid_request", "Invalid request body")
	}

	if err := utils.ValidateStruct(&req); err != nil {
		return utils.SendValidationError(c, err)
	}

	response, err := h.authService.RefreshToken(req.RefreshToken)
	if err != nil {
		return utils.SendError(c, fiber.StatusUnauthorized, "token_refresh_failed", err.Error())
	}

	return utils.SendSuccess(c, fiber.StatusOK, "Token refreshed successfully", response)
}

func (h *AuthHandler) Logout(c *fiber.Ctx) error {
	return utils.SendSuccess(c, fiber.StatusOK, "Logout successful", nil)
}

func (h *AuthHandler) Profile(c *fiber.Ctx) error {
	user := middleware.GetUserFromContext(c)
	if user == nil {
		return utils.SendError(c, fiber.StatusUnauthorized, "unauthorized", "User not found")
	}

	return utils.SendSuccess(c, fiber.StatusOK, "Profile retrieved successfully", user.ToResponse())
}

func (h *AuthHandler) UpdateProfile(c *fiber.Ctx) error {
	user := middleware.GetUserFromContext(c)
	if user == nil {
		return utils.SendError(c, fiber.StatusUnauthorized, "unauthorized", "User not found")
	}

	var req models.UserUpdateInput
	if err := c.BodyParser(&req); err != nil {
		return utils.SendError(c, fiber.StatusBadRequest, "invalid_request", "Invalid request body")
	}

	if err := utils.ValidateStruct(&req); err != nil {
		return utils.SendValidationError(c, err)
	}

	if req.Email != "" {
		user.Email = req.Email
	}
	if req.Username != "" {
		user.Username = req.Username
	}
	if req.FirstName != "" {
		user.FirstName = req.FirstName
	}
	if req.LastName != "" {
		user.LastName = req.LastName
	}

	return utils.SendSuccess(c, fiber.StatusOK, "Profile updated successfully", user.ToResponse())
}

func (h *AuthHandler) UpdatePassword(c *fiber.Ctx) error {
	user := middleware.GetUserFromContext(c)
	if user == nil {
		return utils.SendError(c, fiber.StatusUnauthorized, "unauthorized", "User not found")
	}

	var req models.PasswordUpdateInput
	if err := c.BodyParser(&req); err != nil {
		return utils.SendError(c, fiber.StatusBadRequest, "invalid_request", "Invalid request body")
	}

	if err := utils.ValidateStruct(&req); err != nil {
		return utils.SendValidationError(c, err)
	}

	if err := utils.CheckPassword(req.CurrentPassword, user.PasswordHash); err != nil {
		return utils.SendError(c, fiber.StatusBadRequest, "invalid_password", "Current password is incorrect")
	}

	hashedPassword, err := utils.HashPassword(req.NewPassword)
	if err != nil {
		return utils.SendError(c, fiber.StatusInternalServerError, "internal_error", "Failed to hash password")
	}

	user.PasswordHash = hashedPassword

	return utils.SendSuccess(c, fiber.StatusOK, "Password updated successfully", nil)
}