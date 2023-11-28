package controller

import (
	"errors"
	"net/http"
	"pbi-rakamin-go/helper"
	"pbi-rakamin-go/models"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func Register(context *gin.Context) {
    var input models.RegitstrationInput

    if err := context.ShouldBindJSON(&input); err != nil {
        context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    user := models.User{
        Username: input.Username,
		Email:    input.Email,
        Password: input.Password,
    }

    savedUser, err := user.Save()

    if err != nil {
        context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    context.JSON(http.StatusCreated, gin.H{"user": savedUser})
}

// func Login(context *gin.Context) {
//     var input models.LoginInput

//     if err := context.ShouldBindJSON(&input); err != nil {
//         context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
//         return
//     }

//     user, err := models.FindUserByEmail(input.Email)

//     if err != nil {
//         context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
//         return
//     }

//     err = user.ValidatePassword(input.Password)

//     if err != nil {
//         context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
//         return
//     }

//     jwt, err := helper.GenerateJWT(user)
//     if err != nil {
//         context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
//         return
//     }

//     context.JSON(http.StatusOK, gin.H{
//         "id": user.ID,
//         "username": user.Username,
//         "email": user.Email,
//         "password": user.Password,
//         "jwt": jwt,
//     })
// } 

func Login(context *gin.Context) {
    var input models.LoginInput

    if err := context.ShouldBindJSON(&input); err != nil {
        context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    user, err := models.FindUserByEmail(input.Email)

    if err != nil {
        // Check if the error is because the user doesn't exist
        if errors.Is(err, gorm.ErrRecordNotFound) {
            context.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
            return
        } else {
            context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
            return
        }
    }

    // Validate the password using bcrypt's CompareHashAndPassword function
    err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password))

    if err != nil {
        // Handle invalid password error
        context.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
        return
    }

    // Generate and return JWT token
    jwt, err := helper.GenerateJWT(user)
    if err != nil {
        context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    context.JSON(http.StatusOK, gin.H{
        "id": user.ID,
        "username": user.Username,
        "email": user.Email,
        "password": user.Password,
        "jwt": jwt,
    })
}

