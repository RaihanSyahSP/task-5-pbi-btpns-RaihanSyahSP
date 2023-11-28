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
    // Get the user ID from the request parameters
    userId, err := strconv.ParseUint(context.Param("userId"), 10, 64)
    if err != nil {
        context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    // Find the user by ID
    user, err := models.FindUserById(uint(userId))
    if err != nil {
        context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    // Delete the user from the database
    if err := user.Delete(); err != nil {
        context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    // Send a successful response
    context.JSON(http.StatusOK, gin.H{"message": "User deleted successfully"})
}


func UpdateUser(context *gin.Context) {
    // Dapatkan user dari JWT atau sesuai dengan kebutuhan aplikasi Anda
    currentUser, err := helper.CurrentUser(context)
    if err != nil {
        context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    // Pastikan user yang melakukan request adalah pemilik akun yang akan diupdate
    requestedUserId := getUserIDFromRequest(context)
    if currentUser.ID != requestedUserId {
        context.JSON(http.StatusForbidden, gin.H{"error": "Not authorized to update this user"})
        return
    }

    // Dapatkan input dari body request
    var userInput models.User
    if err := context.ShouldBindJSON(&userInput); err != nil {
        context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    // Dapatkan user yang akan diupdate dari database
    userToUpdate, err := models.FindUserById(requestedUserId)
    if err != nil {
        context.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to find the user"})
        return
    }

    // Update informasi user sesuai dengan input
    userToUpdate.Username = userInput.Username
    userToUpdate.Email = userInput.Email

    // Hash password yang baru
    hashedPassword, err := bcrypt.GenerateFromPassword([]byte(userInput.Password), 10)
    if err != nil {
        return
    }

    userToUpdate.Password = string(hashedPassword)

    // Simpan perubahan ke database
    if err := userToUpdate.Update(); err != nil {
        context.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update the user"})
        return
    }


    // Kirim response sukses
    context.JSON(http.StatusOK, gin.H{"message": "User updated successfully"})
}

func getUserIDFromRequest(context *gin.Context) uint {
    // Fungsi ini bergantung pada bagaimana Anda mengekstrak userId dari request.
    // Sesuaikan dengan kebutuhan dan implementasi endpoint Anda.
    // Contoh: Mendapatkan userId dari path parameter
    userId, _ := strconv.ParseUint(context.Param("userId"), 10, 64)
    return uint(userId)
}
