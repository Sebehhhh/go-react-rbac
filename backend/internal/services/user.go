package services

import (
	"errors"
	"math"
	"strconv"
	"strings"

	"rbac-system/backend/internal/database"
	"rbac-system/backend/internal/models"
	"rbac-system/backend/internal/utils"

	"gorm.io/gorm"
)

type UserService struct {
	rbacService *RBACService
}

func NewUserService(rbacService *RBACService) *UserService {
	return &UserService{
		rbacService: rbacService,
	}
}

func (s *UserService) GetUsers(page, limit int, search, sortBy, sortOrder string) (*models.UserListResponse, error) {
	offset := (page - 1) * limit
	
	query := database.DB.Model(&models.User{}).Preload("Role")
	
	if search != "" {
		searchTerm := "%" + strings.ToLower(search) + "%"
		query = query.Where(
			"LOWER(first_name) LIKE ? OR LOWER(last_name) LIKE ? OR LOWER(email) LIKE ? OR LOWER(username) LIKE ?",
			searchTerm, searchTerm, searchTerm, searchTerm,
		)
	}
	
	if sortBy == "" {
		sortBy = "created_at"
	}
	if sortOrder == "" {
		sortOrder = "desc"
	}
	
	orderClause := sortBy + " " + sortOrder
	query = query.Order(orderClause)
	
	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, err
	}
	
	var users []models.User
	if err := query.Offset(offset).Limit(limit).Find(&users).Error; err != nil {
		return nil, err
	}
	
	userResponses := make([]models.UserResponse, len(users))
	for i, user := range users {
		userResponses[i] = *user.ToResponse()
	}
	
	totalPages := int(math.Ceil(float64(total) / float64(limit)))
	
	return &models.UserListResponse{
		Users:      userResponses,
		Total:      total,
		Page:       page,
		Limit:      limit,
		TotalPages: totalPages,
	}, nil
}

func (s *UserService) GetUserByID(id uint) (*models.User, error) {
	var user models.User
	if err := database.DB.Preload("Role.Permissions").First(&user, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("user not found")
		}
		return nil, err
	}
	return &user, nil
}

func (s *UserService) CreateUser(req *models.UserInput) (*models.User, error) {
	var existingUser models.User
	if err := database.DB.Where("email = ? OR username = ?", req.Email, req.Username).First(&existingUser).Error; err == nil {
		return nil, errors.New("user with this email or username already exists")
	}

	var role models.Role
	if err := database.DB.First(&role, req.RoleID).Error; err != nil {
		return nil, errors.New("role not found")
	}

	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		return nil, err
	}

	user := models.User{
		Email:        req.Email,
		Username:     req.Username,
		PasswordHash: hashedPassword,
		FirstName:    req.FirstName,
		LastName:     req.LastName,
		RoleID:       req.RoleID,
		IsActive:     true,
	}

	if err := database.DB.Create(&user).Error; err != nil {
		return nil, err
	}

	if err := database.DB.Preload("Role.Permissions").First(&user, user.ID).Error; err != nil {
		return nil, err
	}

	return &user, nil
}

func (s *UserService) UpdateUser(id uint, req *models.UserUpdateInput) (*models.User, error) {
	var user models.User
	if err := database.DB.First(&user, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("user not found")
		}
		return nil, err
	}

	if req.Email != "" && req.Email != user.Email {
		var existingUser models.User
		if err := database.DB.Where("email = ? AND id != ?", req.Email, id).First(&existingUser).Error; err == nil {
			return nil, errors.New("user with this email already exists")
		}
		user.Email = req.Email
	}

	if req.Username != "" && req.Username != user.Username {
		var existingUser models.User
		if err := database.DB.Where("username = ? AND id != ?", req.Username, id).First(&existingUser).Error; err == nil {
			return nil, errors.New("user with this username already exists")
		}
		user.Username = req.Username
	}

	if req.FirstName != "" {
		user.FirstName = req.FirstName
	}
	if req.LastName != "" {
		user.LastName = req.LastName
	}
	if req.RoleID != 0 {
		var role models.Role
		if err := database.DB.First(&role, req.RoleID).Error; err != nil {
			return nil, errors.New("role not found")
		}
		user.RoleID = req.RoleID
	}
	if req.IsActive != nil {
		user.IsActive = *req.IsActive
	}

	if err := database.DB.Save(&user).Error; err != nil {
		return nil, err
	}

	if err := database.DB.Preload("Role.Permissions").First(&user, user.ID).Error; err != nil {
		return nil, err
	}

	return &user, nil
}

func (s *UserService) DeleteUser(id uint) error {
	var user models.User
	if err := database.DB.First(&user, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("user not found")
		}
		return err
	}

	return database.DB.Delete(&user).Error
}

func (s *UserService) ActivateUser(id uint) (*models.User, error) {
	var user models.User
	if err := database.DB.First(&user, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("user not found")
		}
		return nil, err
	}

	user.IsActive = true
	if err := database.DB.Save(&user).Error; err != nil {
		return nil, err
	}

	if err := database.DB.Preload("Role.Permissions").First(&user, user.ID).Error; err != nil {
		return nil, err
	}

	return &user, nil
}

func (s *UserService) DeactivateUser(id uint) (*models.User, error) {
	var user models.User
	if err := database.DB.First(&user, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("user not found")
		}
		return nil, err
	}

	user.IsActive = false
	if err := database.DB.Save(&user).Error; err != nil {
		return nil, err
	}

	if err := database.DB.Preload("Role.Permissions").First(&user, user.ID).Error; err != nil {
		return nil, err
	}

	return &user, nil
}

func (s *UserService) GetUserActivityLogs(userID uint, page, limit int) (*models.ActivityLogListResponse, error) {
	offset := (page - 1) * limit
	
	var total int64
	if err := database.DB.Model(&models.ActivityLog{}).Where("user_id = ?", userID).Count(&total).Error; err != nil {
		return nil, err
	}
	
	var activityLogs []models.ActivityLog
	if err := database.DB.Preload("User").Where("user_id = ?", userID).
		Order("created_at desc").Offset(offset).Limit(limit).Find(&activityLogs).Error; err != nil {
		return nil, err
	}
	
	activities := make([]models.ActivityLogResponse, len(activityLogs))
	for i, log := range activityLogs {
		activities[i] = *log.ToResponse()
	}
	
	totalPages := int(math.Ceil(float64(total) / float64(limit)))
	
	return &models.ActivityLogListResponse{
		Activities: activities,
		Total:      total,
		Page:       page,
		Limit:      limit,
		TotalPages: totalPages,
	}, nil
}

type BulkActionRequest struct {
	UserIDs []uint `json:"user_ids" validate:"required"`
	Action  string `json:"action" validate:"required,oneof=activate deactivate delete"`
}

func (s *UserService) BulkUserActions(req *BulkActionRequest) error {
	if len(req.UserIDs) == 0 {
		return errors.New("no user IDs provided")
	}

	switch req.Action {
	case "activate":
		return database.DB.Model(&models.User{}).Where("id IN ?", req.UserIDs).Update("is_active", true).Error
	case "deactivate":
		return database.DB.Model(&models.User{}).Where("id IN ?", req.UserIDs).Update("is_active", false).Error
	case "delete":
		return database.DB.Where("id IN ?", req.UserIDs).Delete(&models.User{}).Error
	default:
		return errors.New("invalid action")
	}
}

func parseIntOrDefault(s string, defaultVal int) int {
	if val, err := strconv.Atoi(s); err == nil {
		return val
	}
	return defaultVal
}