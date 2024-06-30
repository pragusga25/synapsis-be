package models

import (
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Name     string `gorm:"type:varchar(52);not null" json:"name"`
	Email    string `gorm:"unique;not null" json:"email"`
	Password string `gorm:"not null" json:"password"`
	Role     Role   `gorm:"type:varchar(20)" json:"role"`
}
