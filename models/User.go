package models

import (
	"gorm.io/gorm"
)

type User struct {
    gorm.Model
    ID       uint   `gorm:"primaryKey;not null"`
    Username string `gorm:"not null"`
    Email    string `gorm:"unique;not null"`
    Password string `gorm:"not null;min:6"`
    Photos   []Photo
}

