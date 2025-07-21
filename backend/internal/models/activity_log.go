package models

import (
	"time"

	"gorm.io/gorm"
)

type ActivityLog struct {
	ID        uint           `json:"id" gorm:"primarykey"`
	UserID    uint           `json:"user_id" gorm:"not null"`
	Action    string         `json:"action" gorm:"not null" validate:"required,max=100"`
	Resource  string         `json:"resource" gorm:"not null" validate:"required,max=100"`
	Details   string         `json:"details" gorm:"type:text"`
	IPAddress string         `json:"ip_address" gorm:"max=45"`
	UserAgent string         `json:"user_agent" gorm:"type:text"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`

	User User `json:"user" gorm:"foreignKey:UserID"`
}

type ActivityLogInput struct {
	UserID    uint   `json:"user_id" validate:"required"`
	Action    string `json:"action" validate:"required,max=100"`
	Resource  string `json:"resource" validate:"required,max=100"`
	Details   string `json:"details" validate:"max=1000"`
	IPAddress string `json:"ip_address" validate:"max=45"`
	UserAgent string `json:"user_agent" validate:"max=500"`
}

type ActivityLogResponse struct {
	ID        uint      `json:"id"`
	UserID    uint      `json:"user_id"`
	Action    string    `json:"action"`
	Resource  string    `json:"resource"`
	Details   string    `json:"details"`
	IPAddress string    `json:"ip_address"`
	UserAgent string    `json:"user_agent"`
	CreatedAt time.Time `json:"created_at"`
	User      struct {
		ID        uint   `json:"id"`
		Username  string `json:"username"`
		FirstName string `json:"first_name"`
		LastName  string `json:"last_name"`
	} `json:"user"`
}

type ActivityLogListResponse struct {
	Activities []ActivityLogResponse `json:"activities"`
	Total      int64                 `json:"total"`
	Page       int                   `json:"page"`
	Limit      int                   `json:"limit"`
	TotalPages int                   `json:"total_pages"`
}

func (ActivityLog) TableName() string {
	return "activity_logs"
}

func (a *ActivityLog) ToResponse() *ActivityLogResponse {
	response := &ActivityLogResponse{
		ID:        a.ID,
		UserID:    a.UserID,
		Action:    a.Action,
		Resource:  a.Resource,
		Details:   a.Details,
		IPAddress: a.IPAddress,
		UserAgent: a.UserAgent,
		CreatedAt: a.CreatedAt,
	}
	
	response.User.ID = a.User.ID
	response.User.Username = a.User.Username
	response.User.FirstName = a.User.FirstName
	response.User.LastName = a.User.LastName

	return response
}