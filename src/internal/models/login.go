package models

import (
	"time"
)

// Login represents a user's login attempt
type Login struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	Username  string    `json:"username" gorm:"not null"`
	Password  string    `json:"password" gorm:"not null"`
	LoginTime time.Time `json:"login_time" gorm:"default:current_timestamp"`
}


