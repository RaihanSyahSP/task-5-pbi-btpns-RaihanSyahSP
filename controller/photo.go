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

	userWithPhotos, err := models.FindUserById(user.ID)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

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

func UpdatePhoto (context *gin.Context) {
	user, err := helper.CurrentUser(context)

	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	photoId := getPhotoIDFromRequest(context)

	photo, err := models.FindPhotoById(photoId)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if photo.UserID != user.ID {
		context.JSON(http.StatusForbidden, gin.H{"error": "Not authorized to update this photo"})
		return
	}

	var input models.Photo
	if err := context.ShouldBindJSON(&input); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	photo.Title = input.Title
	photo.Caption = input.Caption
	photo.PhotoURL = input.PhotoURL

	if err := photo.Update(); err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	context.JSON(http.StatusOK, gin.H{"data": photo})
}

func DeletePhoto(context *gin.Context) {
	user, err := helper.CurrentUser(context)

	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	photoId := getPhotoIDFromRequest(context)

	photo, err := models.FindPhotoById(photoId)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if photo.UserID != user.ID {
		context.JSON(http.StatusForbidden, gin.H{"error": "Not authorized to update this photo"})
		return
	}

	if err := photo.Delete(); err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	context.JSON(http.StatusOK, gin.H{"message": "Photo deleted successfully"})
}

func getPhotoIDFromRequest(context *gin.Context) string {
    photoId := context.Param("photoId")
    return photoId
}

