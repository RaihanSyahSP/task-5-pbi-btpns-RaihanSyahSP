package models

type LoginInput struct {
	Email    string `json:"email" binding:"required" gorm:"unique"`
	Password string `json:"password" binding:"required,min=6"`
}
