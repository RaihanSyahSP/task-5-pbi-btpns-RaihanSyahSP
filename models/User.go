package models

import (
	"errors"
	"html"
	"pbi-rakamin-go/database"
	"strings"

	"golang.org/x/crypto/bcrypt"
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

func (user *User) Save() (*User, error) {
    err := database.Database.Create(&user).Error
    if err != nil {
        return &User{}, err
    }
    return user, nil
}

func (user *User) Update() error {
    // Pastikan user memiliki ID yang valid
    if user.ID == 0 {
        return errors.New("Invalid user ID")
    }

    // Lakukan operasi update ke database
    err := database.Database.Model(&User{}).Where("id = ?", user.ID).Updates(&user).Error
    if err != nil {
        return err
    }

    return nil
}

func (user *User) Delete() error {
    if user.ID == 0 {
        return errors.New("Invalid user ID")
    }

    err:= database.Database.Model(&User{}).Where("id = ?", user.ID).Delete(&user).Error
    if err != nil {
        return err
    }

    return nil
}



func (user *User) BeforeSave(*gorm.DB) error {
    passwordHash, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
    if err != nil {
        return err
    }
    user.Password = string(passwordHash)
    user.Username = html.EscapeString(strings.TrimSpace(user.Username))
    return nil
}

func (user *User) ValidatePassword(password string) error {
    return bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
}

func FindUserByEmail(email string) (User, error) {
    var user User
    err := database.Database.Where("email=?", email).Find(&user).Error
    if err != nil {
        return User{}, err
    }
    return user, nil
}

func FindUserById(id uint) (User, error) {
    var user User
    err := database.Database.Preload("Photos").Where("ID=?", id).Find(&user).Error
    if err != nil {
        return User{}, err
    }
    return user, nil
}



