package services

import (
	"errors"
	"time"

	"rbac-system/backend/internal/database"
	"rbac-system/backend/internal/models"
	"rbac-system/backend/internal/utils"

	"gorm.io/gorm"
)

type AuthService struct {
	jwtService *utils.JWTService
}

func NewAuthService(jwtService *utils.JWTService) *AuthService {
	return &AuthService{
		jwtService: jwtService,
	}
}

func (s *AuthService) Register(req *models.RegisterRequest) (*models.TokenResponse, error) {
	var existingUser models.User
	if err := database.DB.Where("email = ? OR username = ?", req.Email, req.Username).First(&existingUser).Error; err == nil {
		return nil, errors.New("user with this email or username already exists")
	}

	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		return nil, err
	}

	var defaultRole models.Role
	if err := database.DB.Where("name = ?", "User").First(&defaultRole).Error; err != nil {
		return nil, errors.New("default role not found")
	}

	user := models.User{
		Email:        req.Email,
		Username:     req.Username,
		PasswordHash: hashedPassword,
		FirstName:    req.FirstName,
		LastName:     req.LastName,
		RoleID:       defaultRole.ID,
		IsActive:     true,
	}

	if err := database.DB.Create(&user).Error; err != nil {
		return nil, err
	}

	if err := database.DB.Preload("Role.Permissions").First(&user, user.ID).Error; err != nil {
		return nil, err
	}

	return s.generateTokenResponse(&user)
}

func (s *AuthService) Login(req *models.LoginRequest, ipAddress, userAgent string) (*models.TokenResponse, error) {
	var user models.User
	if err := database.DB.Preload("Role.Permissions").Where("email = ?", req.Email).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("invalid credentials")
		}
		return nil, err
	}

	if !user.IsActive {
		return nil, errors.New("account is deactivated")
	}

	if err := utils.CheckPassword(req.Password, user.PasswordHash); err != nil {
		return nil, errors.New("invalid credentials")
	}

	now := time.Now()
	user.LastLoginAt = &now
	database.DB.Save(&user)

	activityLog := models.ActivityLog{
		UserID:    user.ID,
		Action:    "login",
		Resource:  "auth",
		Details:   "User logged in successfully",
		IPAddress: ipAddress,
		UserAgent: userAgent,
	}
	database.DB.Create(&activityLog)

	return s.generateTokenResponse(&user)
}

func (s *AuthService) RefreshToken(refreshToken string) (*models.TokenResponse, error) {
	claims, err := s.jwtService.ValidateRefreshToken(refreshToken)
	if err != nil {
		return nil, errors.New("invalid refresh token")
	}

	var user models.User
	if err := database.DB.Preload("Role.Permissions").First(&user, claims.UserID).Error; err != nil {
		return nil, errors.New("user not found")
	}

	if !user.IsActive {
		return nil, errors.New("account is deactivated")
	}

	return s.generateTokenResponse(&user)
}

func (s *AuthService) generateTokenResponse(user *models.User) (*models.TokenResponse, error) {
	accessToken, expiresAt, err := s.jwtService.GenerateAccessToken(user)
	if err != nil {
		return nil, err
	}

	refreshToken, _, err := s.jwtService.GenerateRefreshToken(user)
	if err != nil {
		return nil, err
	}

	return &models.TokenResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresAt:    expiresAt,
		User:         *user.ToResponse(),
	}, nil
}

func (s *AuthService) GetUserByID(userID uint) (*models.User, error) {
	var user models.User
	if err := database.DB.Preload("Role.Permissions").First(&user, userID).Error; err != nil {
		return nil, err
	}
	return &user, nil
}