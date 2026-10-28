package models

import (
	"time"
)

// User represents the user model in SQLite
type User struct {
	ID        uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	Username  string    `gorm:"unique;not null" json:"username"`
	Password  string    `gorm:"not null" json:"-"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
}
