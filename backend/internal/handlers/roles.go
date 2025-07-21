package handlers

import (
	"strconv"

	"rbac-system/backend/internal/models"
	"rbac-system/backend/internal/services"
	"rbac-system/backend/internal/utils"

	"github.com/gofiber/fiber/v2"
)

type RoleHandler struct {
	roleService *services.RoleService
}

func NewRoleHandler(roleService *services.RoleService) *RoleHandler {
	return &RoleHandler{
		roleService: roleService,
	}
}

func (h *RoleHandler) GetRoles(c *fiber.Ctx) error {
	roles, err := h.roleService.GetRoles()
	if err != nil {
		return utils.SendError(c, fiber.StatusInternalServerError, "internal_error", err.Error())
	}

	return utils.SendSuccess(c, fiber.StatusOK, "Roles retrieved successfully", roles)
}

func (h *RoleHandler) GetRole(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return utils.SendError(c, fiber.StatusBadRequest, "invalid_id", "Invalid role ID")
	}

	role, err := h.roleService.GetRoleByID(uint(id))
	if err != nil {
		return utils.SendError(c, fiber.StatusNotFound, "role_not_found", err.Error())
	}

	return utils.SendSuccess(c, fiber.StatusOK, "Role retrieved successfully", role)
}

func (h *RoleHandler) CreateRole(c *fiber.Ctx) error {
	var req models.RoleInput
	if err := c.BodyParser(&req); err != nil {
		return utils.SendError(c, fiber.StatusBadRequest, "invalid_request", "Invalid request body")
	}

	if err := utils.ValidateStruct(&req); err != nil {
		return utils.SendValidationError(c, err)
	}

	role, err := h.roleService.CreateRole(&req)
	if err != nil {
		return utils.SendError(c, fiber.StatusBadRequest, "create_failed", err.Error())
	}

	return utils.SendSuccess(c, fiber.StatusCreated, "Role created successfully", role)
}

func (h *RoleHandler) UpdateRole(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return utils.SendError(c, fiber.StatusBadRequest, "invalid_id", "Invalid role ID")
	}

	var req models.RoleInput
	if err := c.BodyParser(&req); err != nil {
		return utils.SendError(c, fiber.StatusBadRequest, "invalid_request", "Invalid request body")
	}

	if err := utils.ValidateStruct(&req); err != nil {
		return utils.SendValidationError(c, err)
	}

	role, err := h.roleService.UpdateRole(uint(id), &req)
	if err != nil {
		return utils.SendError(c, fiber.StatusBadRequest, "update_failed", err.Error())
	}

	return utils.SendSuccess(c, fiber.StatusOK, "Role updated successfully", role)
}

func (h *RoleHandler) DeleteRole(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return utils.SendError(c, fiber.StatusBadRequest, "invalid_id", "Invalid role ID")
	}

	err = h.roleService.DeleteRole(uint(id))
	if err != nil {
		return utils.SendError(c, fiber.StatusBadRequest, "delete_failed", err.Error())
	}

	return utils.SendSuccess(c, fiber.StatusOK, "Role deleted successfully", nil)
}

func (h *RoleHandler) GetPermissions(c *fiber.Ctx) error {
	permissions, err := h.roleService.GetPermissions()
	if err != nil {
		return utils.SendError(c, fiber.StatusInternalServerError, "internal_error", err.Error())
	}

	return utils.SendSuccess(c, fiber.StatusOK, "Permissions retrieved successfully", permissions)
}

func (h *RoleHandler) AssignPermissions(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return utils.SendError(c, fiber.StatusBadRequest, "invalid_id", "Invalid role ID")
	}

	var req models.RolePermissionInput
	if err := c.BodyParser(&req); err != nil {
		return utils.SendError(c, fiber.StatusBadRequest, "invalid_request", "Invalid request body")
	}

	if err := utils.ValidateStruct(&req); err != nil {
		return utils.SendValidationError(c, err)
	}

	err = h.roleService.AssignPermissions(uint(id), req.PermissionIDs)
	if err != nil {
		return utils.SendError(c, fiber.StatusBadRequest, "assign_failed", err.Error())
	}

	return utils.SendSuccess(c, fiber.StatusOK, "Permissions assigned successfully", nil)
}

func (h *RoleHandler) GetRolePermissions(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return utils.SendError(c, fiber.StatusBadRequest, "invalid_id", "Invalid role ID")
	}

	permissions, err := h.roleService.GetRolePermissions(uint(id))
	if err != nil {
		return utils.SendError(c, fiber.StatusBadRequest, "get_permissions_failed", err.Error())
	}

	return utils.SendSuccess(c, fiber.StatusOK, "Role permissions retrieved successfully", permissions)
}