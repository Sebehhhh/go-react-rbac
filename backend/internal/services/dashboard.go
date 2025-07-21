package services

import (
	"time"

	"rbac-system/backend/internal/database"
	"rbac-system/backend/internal/models"
)

type DashboardService struct{}

func NewDashboardService() *DashboardService {
	return &DashboardService{}
}

type DashboardStats struct {
	TotalUsers       int64 `json:"total_users"`
	ActiveUsers      int64 `json:"active_users"`
	InactiveUsers    int64 `json:"inactive_users"`
	NewUsersToday    int64 `json:"new_users_today"`
	NewUsersThisWeek int64 `json:"new_users_this_week"`
	TotalRoles       int64 `json:"total_roles"`
}

type RoleDistribution struct {
	RoleName  string `json:"role_name"`
	UserCount int64  `json:"user_count"`
}

type UserAnalytics struct {
	Date      string `json:"date"`
	UserCount int64  `json:"user_count"`
}

func (s *DashboardService) GetDashboardStats() (*DashboardStats, error) {
	var stats DashboardStats

	if err := database.DB.Model(&models.User{}).Count(&stats.TotalUsers).Error; err != nil {
		return nil, err
	}

	if err := database.DB.Model(&models.User{}).Where("is_active = ?", true).Count(&stats.ActiveUsers).Error; err != nil {
		return nil, err
	}

	if err := database.DB.Model(&models.User{}).Where("is_active = ?", false).Count(&stats.InactiveUsers).Error; err != nil {
		return nil, err
	}

	today := time.Now().Format("2006-01-02")
	if err := database.DB.Model(&models.User{}).Where("DATE(created_at) = ?", today).Count(&stats.NewUsersToday).Error; err != nil {
		return nil, err
	}

	weekAgo := time.Now().AddDate(0, 0, -7)
	if err := database.DB.Model(&models.User{}).Where("created_at >= ?", weekAgo).Count(&stats.NewUsersThisWeek).Error; err != nil {
		return nil, err
	}

	if err := database.DB.Model(&models.Role{}).Count(&stats.TotalRoles).Error; err != nil {
		return nil, err
	}

	return &stats, nil
}

func (s *DashboardService) GetRoleDistribution() ([]RoleDistribution, error) {
	var distributions []RoleDistribution

	err := database.DB.Model(&models.User{}).
		Select("roles.name as role_name, COUNT(users.id) as user_count").
		Joins("JOIN roles ON users.role_id = roles.id").
		Group("roles.id, roles.name").
		Scan(&distributions).Error

	if err != nil {
		return nil, err
	}

	return distributions, nil
}

func (s *DashboardService) GetRecentActivity(limit int) ([]models.ActivityLogResponse, error) {
	var activityLogs []models.ActivityLog

	err := database.DB.Preload("User").
		Order("created_at DESC").
		Limit(limit).
		Find(&activityLogs).Error

	if err != nil {
		return nil, err
	}

	activities := make([]models.ActivityLogResponse, len(activityLogs))
	for i, log := range activityLogs {
		activities[i] = *log.ToResponse()
	}

	return activities, nil
}

func (s *DashboardService) GetUserAnalytics(days int) ([]UserAnalytics, error) {
	var analytics []UserAnalytics

	startDate := time.Now().AddDate(0, 0, -days)

	rows, err := database.DB.Raw(`
		SELECT DATE(created_at) as date, COUNT(*) as user_count
		FROM users
		WHERE created_at >= ?
		GROUP BY DATE(created_at)
		ORDER BY DATE(created_at)
	`, startDate).Rows()

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var analytic UserAnalytics
		if err := rows.Scan(&analytic.Date, &analytic.UserCount); err != nil {
			return nil, err
		}
		analytics = append(analytics, analytic)
	}

	return analytics, nil
}

type SystemHealth struct {
	DatabaseStatus string `json:"database_status"`
	TotalRequests  int64  `json:"total_requests"`
	ActiveSessions int64  `json:"active_sessions"`
	SystemUptime   string `json:"system_uptime"`
}

func (s *DashboardService) GetSystemHealth() (*SystemHealth, error) {
	health := &SystemHealth{
		DatabaseStatus: "healthy",
		TotalRequests:  0, 
		ActiveSessions: 0,
		SystemUptime:   "N/A",
	}

	if err := database.DB.Raw("SELECT 1").Error; err != nil {
		health.DatabaseStatus = "unhealthy"
	}

	var activityCount int64
	if err := database.DB.Model(&models.ActivityLog{}).Count(&activityCount).Error; err == nil {
		health.TotalRequests = activityCount
	}

	var activeUserCount int64
	if err := database.DB.Model(&models.User{}).Where("is_active = ? AND last_login_at > ?", true, time.Now().Add(-24*time.Hour)).Count(&activeUserCount).Error; err == nil {
		health.ActiveSessions = activeUserCount
	}

	return health, nil
}