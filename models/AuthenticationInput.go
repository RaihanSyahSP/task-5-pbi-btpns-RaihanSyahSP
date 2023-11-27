package models

type AuthenticationInput struct {
	Username string `json:"username" binding:"required"`
	Email    string `json:"email" binding:"required" gorm:"unique"`
	Password string `json:"password" binding:"required,min=6"`
}
