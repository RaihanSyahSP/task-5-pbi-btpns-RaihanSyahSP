package controller

import (
	"net/http"
	"pbi-rakamin-go/helper"
	"pbi-rakamin-go/models"
	"strconv"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)


func DeleteUser(context *gin.Context) {

    userId, err := strconv.ParseUint(context.Param("userId"), 10, 64)
    if err != nil {
        context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    user, err := models.FindUserById(uint(userId))
    if err != nil {
        context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    if err := user.Delete(); err != nil {
        context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    context.JSON(http.StatusOK, gin.H{"message": "User deleted successfully"})
}


func UpdateUser(context *gin.Context) {

    currentUser, err := helper.CurrentUser(context)
    if err != nil {
        context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    requestedUserId := getUserIDFromRequest(context)
    if currentUser.ID != requestedUserId {
        context.JSON(http.StatusForbidden, gin.H{"error": "Not authorized to update this user"})
        return
    }

    var userInput models.User
    if err := context.ShouldBindJSON(&userInput); err != nil {
        context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    userToUpdate, err := models.FindUserById(requestedUserId)
    if err != nil {
        context.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to find the user"})
        return
    }

    userToUpdate.Username = userInput.Username
    userToUpdate.Email = userInput.Email

    hashedPassword, err := bcrypt.GenerateFromPassword([]byte(userInput.Password), 10)
    if err != nil {
        return
    }

    userToUpdate.Password = string(hashedPassword)

    if err := userToUpdate.Update(); err != nil {
        context.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update the user"})
        return
    }


    context.JSON(http.StatusOK, gin.H{"message": "User updated successfully"})
}

func getUserIDFromRequest(context *gin.Context) uint {
    userId, _ := strconv.ParseUint(context.Param("userId"), 10, 64)
    return uint(userId)
}
