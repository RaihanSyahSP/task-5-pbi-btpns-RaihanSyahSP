package models

import (
	"errors"
	"pbi-rakamin-go/database"

	"gorm.io/gorm"
)

type Photo struct {
	gorm.Model
	ID 	 uint   `gorm:"primaryKey;not null"`
	Title    string `gorm:"not null"`
	Caption  string
	PhotoURL string `gorm:"not null"`
	UserID   uint   `gorm:"not null"`
	User     User   `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
}

func (photo *Photo) Save() (*Photo, error) {
	err := database.Database.Create(&photo).Error
	if err != nil {
		return nil, err
	}

	// Preload the User information
	err = database.Database.Preload("Photos").Where("ID=?", photo.UserID).Find(&photo.User).Error
	if err != nil {
		return nil, err
	}

	return photo, nil
}

func (photo *Photo) Update() error {
	// Pastikan photo memiliki ID yang valid
	if photo.ID == 0 {
		return errors.New("Invalid photo ID")
	}

	// Lakukan operasi update ke database
	err := database.Database.Model(&Photo{}).Where("id = ?", photo.ID).Updates(&photo).Error
	if err != nil {
		return err
	}

	return nil
}

func FindPhotoById(photoId string) (*Photo, error) {
	var photo Photo
	err := database.Database.Where("id = ?", photoId).First(&photo).Error
	if err != nil {
		return nil, err
	}

	return &photo, nil
}



