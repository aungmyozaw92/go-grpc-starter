package entity

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	Username  string         `gorm:"uniqueIndex;not null;size:30" json:"username"`
	Name      string         `gorm:"not null;size:100" json:"name"`
	Email     *string        `gorm:"uniqueIndex;size:100" json:"email"`
	Phone     string         `gorm:"size:20" json:"phone"`
	Mobile    string         `gorm:"size:20" json:"mobile"`
	ImageURL  string         `gorm:"size:255" json:"image_url"`
	Password  string         `gorm:"not null;size:255" json:"-"`
	IsActive  *bool          `gorm:"default:true" json:"is_active"`
	RoleID    int            `gorm:"not null;default:1" json:"role_id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}
