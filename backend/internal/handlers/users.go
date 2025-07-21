package handlers

import (
	"strconv"

	"rbac-system/backend/internal/middleware"
	"rbac-system/backend/internal/models"
	"rbac-system/backend/internal/services"
	"rbac-system/backend/internal/utils"

	"github.com/gofiber/fiber/v2"
)

type UserHandler struct {
	userService *services.UserService
	rbacService *services.RBACService
}

func NewUserHandler(userService *services.UserService, rbacService *services.RBACService) *UserHandler {
	return &UserHandler{
		userService: userService,
		rbacService: rbacService,
	}
}

func (h *UserHandler) GetUsers(c *fiber.Ctx) error {
	page := parseIntOrDefault(c.Query("page"), 1)
	limit := parseIntOrDefault(c.Query("limit"), 10)
	search := c.Query("search", "")
	sortBy := c.Query("sort_by", "")
	sortOrder := c.Query("sort_order", "")

	if limit > 100 {
		limit = 100
	}
	if page < 1 {
		page = 1
	}

	users, err := h.userService.GetUsers(page, limit, search, sortBy, sortOrder)
	if err != nil {
		return utils.SendError(c, fiber.StatusInternalServerError, "internal_error", err.Error())
	}

	return utils.SendSuccess(c, fiber.StatusOK, "Users retrieved successfully", users)
}

func (h *UserHandler) GetUser(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return utils.SendError(c, fiber.StatusBadRequest, "invalid_id", "Invalid user ID")
	}

	user, err := h.userService.GetUserByID(uint(id))
	if err != nil {
		return utils.SendError(c, fiber.StatusNotFound, "user_not_found", err.Error())
	}

	return utils.SendSuccess(c, fiber.StatusOK, "User retrieved successfully", user.ToResponse())
}

func (h *UserHandler) CreateUser(c *fiber.Ctx) error {
	var req models.UserInput
	if err := c.BodyParser(&req); err != nil {
		return utils.SendError(c, fiber.StatusBadRequest, "invalid_request", "Invalid request body")
	}

	if err := utils.ValidateStruct(&req); err != nil {
		return utils.SendValidationError(c, err)
	}

	user, err := h.userService.CreateUser(&req)
	if err != nil {
		return utils.SendError(c, fiber.StatusBadRequest, "create_failed", err.Error())
	}

	return utils.SendSuccess(c, fiber.StatusCreated, "User created successfully", user.ToResponse())
}

func (h *UserHandler) UpdateUser(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return utils.SendError(c, fiber.StatusBadRequest, "invalid_id", "Invalid user ID")
	}

	currentUserID := middleware.GetUserIDFromContext(c)
	
	if currentUserID != uint(id) {
		canManage, err := h.rbacService.CanManageUser(currentUserID, uint(id))
		if err != nil {
			return utils.SendError(c, fiber.StatusInternalServerError, "internal_error", "Error checking permissions")
		}
		if !canManage {
			return utils.SendError(c, fiber.StatusForbidden, "forbidden", "Cannot manage this user")
		}
	}

	var req models.UserUpdateInput
	if err := c.BodyParser(&req); err != nil {
		return utils.SendError(c, fiber.StatusBadRequest, "invalid_request", "Invalid request body")
	}

	if err := utils.ValidateStruct(&req); err != nil {
		return utils.SendValidationError(c, err)
	}

	user, err := h.userService.UpdateUser(uint(id), &req)
	if err != nil {
		return utils.SendError(c, fiber.StatusBadRequest, "update_failed", err.Error())
	}

	return utils.SendSuccess(c, fiber.StatusOK, "User updated successfully", user.ToResponse())
}

func (h *UserHandler) DeleteUser(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return utils.SendError(c, fiber.StatusBadRequest, "invalid_id", "Invalid user ID")
	}

	currentUserID := middleware.GetUserIDFromContext(c)
	if currentUserID == uint(id) {
		return utils.SendError(c, fiber.StatusBadRequest, "cannot_delete_self", "Cannot delete your own account")
	}

	canManage, err := h.rbacService.CanManageUser(currentUserID, uint(id))
	if err != nil {
		return utils.SendError(c, fiber.StatusInternalServerError, "internal_error", "Error checking permissions")
	}
	if !canManage {
		return utils.SendError(c, fiber.StatusForbidden, "forbidden", "Cannot delete this user")
	}

	err = h.userService.DeleteUser(uint(id))
	if err != nil {
		return utils.SendError(c, fiber.StatusBadRequest, "delete_failed", err.Error())
	}

	return utils.SendSuccess(c, fiber.StatusOK, "User deleted successfully", nil)
}

func (h *UserHandler) ActivateUser(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return utils.SendError(c, fiber.StatusBadRequest, "invalid_id", "Invalid user ID")
	}

	user, err := h.userService.ActivateUser(uint(id))
	if err != nil {
		return utils.SendError(c, fiber.StatusBadRequest, "activate_failed", err.Error())
	}

	return utils.SendSuccess(c, fiber.StatusOK, "User activated successfully", user.ToResponse())
}

func (h *UserHandler) DeactivateUser(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return utils.SendError(c, fiber.StatusBadRequest, "invalid_id", "Invalid user ID")
	}

	currentUserID := middleware.GetUserIDFromContext(c)
	if currentUserID == uint(id) {
		return utils.SendError(c, fiber.StatusBadRequest, "cannot_deactivate_self", "Cannot deactivate your own account")
	}

	user, err := h.userService.DeactivateUser(uint(id))
	if err != nil {
		return utils.SendError(c, fiber.StatusBadRequest, "deactivate_failed", err.Error())
	}

	return utils.SendSuccess(c, fiber.StatusOK, "User deactivated successfully", user.ToResponse())
}

func (h *UserHandler) GetUserActivity(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return utils.SendError(c, fiber.StatusBadRequest, "invalid_id", "Invalid user ID")
	}

	page := parseIntOrDefault(c.Query("page"), 1)
	limit := parseIntOrDefault(c.Query("limit"), 20)

	if limit > 100 {
		limit = 100
	}
	if page < 1 {
		page = 1
	}

	activities, err := h.userService.GetUserActivityLogs(uint(id), page, limit)
	if err != nil {
		return utils.SendError(c, fiber.StatusInternalServerError, "internal_error", err.Error())
	}

	return utils.SendSuccess(c, fiber.StatusOK, "User activity retrieved successfully", activities)
}

func (h *UserHandler) BulkActions(c *fiber.Ctx) error {
	var req services.BulkActionRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.SendError(c, fiber.StatusBadRequest, "invalid_request", "Invalid request body")
	}

	if err := utils.ValidateStruct(&req); err != nil {
		return utils.SendValidationError(c, err)
	}

	err := h.userService.BulkUserActions(&req)
	if err != nil {
		return utils.SendError(c, fiber.StatusBadRequest, "bulk_action_failed", err.Error())
	}

	return utils.SendSuccess(c, fiber.StatusOK, "Bulk action completed successfully", nil)
}

func (h *UserHandler) UpdatePassword(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return utils.SendError(c, fiber.StatusBadRequest, "invalid_id", "Invalid user ID")
	}

	currentUserID := middleware.GetUserIDFromContext(c)
	
	// Only allow users to update their own password or admin can update any user's password
	if currentUserID != uint(id) {
		canManage, err := h.rbacService.CanManageUser(currentUserID, uint(id))
		if err != nil {
			return utils.SendError(c, fiber.StatusInternalServerError, "internal_error", "Error checking permissions")
		}
		if !canManage {
			return utils.SendError(c, fiber.StatusForbidden, "forbidden", "Cannot update this user's password")
		}
	}

	var req models.PasswordUpdateInput
	if err := c.BodyParser(&req); err != nil {
		return utils.SendError(c, fiber.StatusBadRequest, "invalid_request", "Invalid request body")
	}

	if err := utils.ValidateStruct(&req); err != nil {
		return utils.SendValidationError(c, err)
	}

	err = h.userService.UpdatePassword(uint(id), &req)
	if err != nil {
		return utils.SendError(c, fiber.StatusBadRequest, "update_password_failed", err.Error())
	}

	return utils.SendSuccess(c, fiber.StatusOK, "Password updated successfully", nil)
}

func parseIntOrDefault(s string, defaultVal int) int {
	if val, err := strconv.Atoi(s); err == nil {
		return val
	}
	return defaultVal
}