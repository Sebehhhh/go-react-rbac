package handlers

import (
	"strconv"

	"rbac-system/backend/internal/services"
	"rbac-system/backend/internal/utils"

	"github.com/gofiber/fiber/v2"
)

type DashboardHandler struct {
	dashboardService *services.DashboardService
}

func NewDashboardHandler(dashboardService *services.DashboardService) *DashboardHandler {
	return &DashboardHandler{
		dashboardService: dashboardService,
	}
}

func (h *DashboardHandler) GetStats(c *fiber.Ctx) error {
	stats, err := h.dashboardService.GetDashboardStats()
	if err != nil {
		return utils.SendError(c, fiber.StatusInternalServerError, "internal_error", err.Error())
	}

	return utils.SendSuccess(c, fiber.StatusOK, "Dashboard stats retrieved successfully", stats)
}

func (h *DashboardHandler) GetRoleDistribution(c *fiber.Ctx) error {
	distribution, err := h.dashboardService.GetRoleDistribution()
	if err != nil {
		return utils.SendError(c, fiber.StatusInternalServerError, "internal_error", err.Error())
	}

	return utils.SendSuccess(c, fiber.StatusOK, "Role distribution retrieved successfully", distribution)
}

func (h *DashboardHandler) GetRecentActivity(c *fiber.Ctx) error {
	limit := 20
	if limitParam := c.Query("limit"); limitParam != "" {
		if parsedLimit, err := strconv.Atoi(limitParam); err == nil && parsedLimit > 0 && parsedLimit <= 100 {
			limit = parsedLimit
		}
	}

	activity, err := h.dashboardService.GetRecentActivity(limit)
	if err != nil {
		return utils.SendError(c, fiber.StatusInternalServerError, "internal_error", err.Error())
	}

	return utils.SendSuccess(c, fiber.StatusOK, "Recent activity retrieved successfully", activity)
}

func (h *DashboardHandler) GetUserAnalytics(c *fiber.Ctx) error {
	days := 30
	if daysParam := c.Query("days"); daysParam != "" {
		if parsedDays, err := strconv.Atoi(daysParam); err == nil && parsedDays > 0 && parsedDays <= 365 {
			days = parsedDays
		}
	}

	analytics, err := h.dashboardService.GetUserAnalytics(days)
	if err != nil {
		return utils.SendError(c, fiber.StatusInternalServerError, "internal_error", err.Error())
	}

	return utils.SendSuccess(c, fiber.StatusOK, "User analytics retrieved successfully", analytics)
}

func (h *DashboardHandler) GetSystemHealth(c *fiber.Ctx) error {
	health, err := h.dashboardService.GetSystemHealth()
	if err != nil {
		return utils.SendError(c, fiber.StatusInternalServerError, "internal_error", err.Error())
	}

	return utils.SendSuccess(c, fiber.StatusOK, "System health retrieved successfully", health)
}