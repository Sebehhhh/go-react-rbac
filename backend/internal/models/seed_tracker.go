package models

import "time"

type SeedTracker struct {
	ID          uint      `json:"id" gorm:"primarykey"`
	SeedName    string    `json:"seed_name" gorm:"type:varchar(255);uniqueIndex;not null"`
	IsCompleted bool      `json:"is_completed" gorm:"default:false"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func (SeedTracker) TableName() string {
	return "seed_trackers"
}