package utils

import (
	"errors"
	"time"

	"rbac-system/backend/internal/config"
	"rbac-system/backend/internal/models"

	"github.com/golang-jwt/jwt/v5"
)

type JWTService struct {
	config *config.Config
}

func NewJWTService(cfg *config.Config) *JWTService {
	return &JWTService{config: cfg}
}

func (s *JWTService) GenerateAccessToken(user *models.User) (string, time.Time, error) {
	expiresAt := time.Now().Add(s.config.JWT.AccessTokenExpiry)
	
	claims := jwt.MapClaims{
		"user_id": user.ID,
		"email":   user.Email,
		"role_id": user.RoleID,
		"type":    "access",
		"exp":     expiresAt.Unix(),
		"iat":     time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(s.config.JWT.Secret))
	if err != nil {
		return "", time.Time{}, err
	}

	return tokenString, expiresAt, nil
}

func (s *JWTService) GenerateRefreshToken(user *models.User) (string, time.Time, error) {
	expiresAt := time.Now().Add(s.config.JWT.RefreshTokenExpiry)
	
	claims := jwt.MapClaims{
		"user_id": user.ID,
		"email":   user.Email,
		"role_id": user.RoleID,
		"type":    "refresh",
		"exp":     expiresAt.Unix(),
		"iat":     time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(s.config.JWT.RefreshSecret))
	if err != nil {
		return "", time.Time{}, err
	}

	return tokenString, expiresAt, nil
}

func (s *JWTService) ValidateAccessToken(tokenString string) (*models.JWTClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(s.config.JWT.Secret), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		tokenType, ok := claims["type"].(string)
		if !ok || tokenType != "access" {
			return nil, errors.New("invalid token type")
		}

		userID, ok := claims["user_id"].(float64)
		if !ok {
			return nil, errors.New("invalid user_id claim")
		}

		email, ok := claims["email"].(string)
		if !ok {
			return nil, errors.New("invalid email claim")
		}

		roleID, ok := claims["role_id"].(float64)
		if !ok {
			return nil, errors.New("invalid role_id claim")
		}

		return &models.JWTClaims{
			UserID: uint(userID),
			Email:  email,
			RoleID: uint(roleID),
			Type:   tokenType,
		}, nil
	}

	return nil, errors.New("invalid token")
}

func (s *JWTService) ValidateRefreshToken(tokenString string) (*models.JWTClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(s.config.JWT.RefreshSecret), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		tokenType, ok := claims["type"].(string)
		if !ok || tokenType != "refresh" {
			return nil, errors.New("invalid token type")
		}

		userID, ok := claims["user_id"].(float64)
		if !ok {
			return nil, errors.New("invalid user_id claim")
		}

		email, ok := claims["email"].(string)
		if !ok {
			return nil, errors.New("invalid email claim")
		}

		roleID, ok := claims["role_id"].(float64)
		if !ok {
			return nil, errors.New("invalid role_id claim")
		}

		return &models.JWTClaims{
			UserID: uint(userID),
			Email:  email,
			RoleID: uint(roleID),
			Type:   tokenType,
		}, nil
	}

	return nil, errors.New("invalid token")
}