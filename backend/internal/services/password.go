package services

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"time"

	"golang.org/x/crypto/bcrypt"

	"rbac-system/backend/internal/database"
	"rbac-system/backend/internal/models"
)

type PasswordService struct{}

func NewPasswordService() *PasswordService {
	return &PasswordService{}
}

// CreateResetToken generates a new password reset token for a user.
func (s *PasswordService) CreateResetToken(email string) (string, error) {
	var user models.User
	if err := database.DB.Where("email = ?", email).First(&user).Error; err != nil {
		return "", errors.New("user not found")
	}

	token := make([]byte, 32)
	if _, err := rand.Read(token); err != nil {
		return "", err
	}
	tokenStr := hex.EncodeToString(token)

	resetToken := models.PasswordResetToken{
		UserID:    user.ID,
		Token:     tokenStr,
		ExpiresAt: time.Now().Add(1 * time.Hour), // Token valid for 1 hour
	}

	if err := database.DB.Create(&resetToken).Error; err != nil {
		return "", err
	}

	// Here you would typically send an email with the reset link.
	// For this example, we'll just return the token.
	// log.Printf("Password reset token for %s: %s", email, tokenStr)

	return tokenStr, nil
}

// ResetPassword validates the token and resets the user's password.
func (s *PasswordService) ResetPassword(token, newPassword string) error {
	var resetToken models.PasswordResetToken
	if err := database.DB.Where("token = ? AND expires_at > ?", token, time.Now()).First(&resetToken).Error; err != nil {
		return errors.New("invalid or expired token")
	}

	var user models.User
	if err := database.DB.First(&user, resetToken.UserID).Error; err != nil {
		return errors.New("user not found")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	user.PasswordHash = string(hashedPassword)
	if err := database.DB.Save(&user).Error; err != nil {
		return err
	}

	// Invalidate the token after use
	if err := database.DB.Delete(&resetToken).Error; err != nil {
		// Log the error, but don't block the user
		// log.Printf("Failed to delete used reset token: %v", err)
	}

	return nil
}
