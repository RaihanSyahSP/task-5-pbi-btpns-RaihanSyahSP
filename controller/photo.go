package controller

import (
	"net/http"
	"pbi-rakamin-go/helper"
	"pbi-rakamin-go/models"

	"github.com/gin-gonic/gin"
)

func AddPhoto(context *gin.Context) {
	var input models.Photo
	if err := context.ShouldBindJSON(&input); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := helper.CurrentUser(context)

	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	input.UserID = user.ID

	savedPhoto, err := input.Save()

	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Fetch the User information to include in the response
	userWithPhotos, err := models.FindUserById(user.ID)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Exclude the Photos array from the User information in the response
	response := gin.H{
		"data": gin.H{
			"CreatedAt": savedPhoto.CreatedAt,
			"UpdatedAt": savedPhoto.UpdatedAt,
			"DeletedAt": savedPhoto.DeletedAt,
			"ID":        savedPhoto.ID,
			"Title":     savedPhoto.Title,
			"Caption":   savedPhoto.Caption,
			"PhotoURL":  savedPhoto.PhotoURL,
			"User": gin.H{
				"CreatedAt": userWithPhotos.CreatedAt,
				"UpdatedAt": userWithPhotos.UpdatedAt,
				"DeletedAt": userWithPhotos.DeletedAt,
				"Username":  userWithPhotos.Username,
				"Email":     userWithPhotos.Email,
				"Password":  userWithPhotos.Password,
                "ID":        userWithPhotos.ID,
			},
		},
	}

	context.JSON(http.StatusCreated, response)
}

func GetAllPhotos(context *gin.Context) {
    user, err := helper.CurrentUser(context)

    if err != nil {
        context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    // Exclude the "User" attribute in the response
    var photoResponse []gin.H
    for _, photo := range user.Photos {
        photoResponse = append(photoResponse, gin.H{
            "CreatedAt": photo.CreatedAt,
            "UpdatedAt": photo.UpdatedAt,
            "DeletedAt": photo.DeletedAt,
            "ID":        photo.ID,
            "Title":     photo.Title,
            "Caption":   photo.Caption,
            "PhotoURL":  photo.PhotoURL,
            "UserID":    photo.UserID,
        })
    }

    context.JSON(http.StatusOK, gin.H{"data": photoResponse})
}

