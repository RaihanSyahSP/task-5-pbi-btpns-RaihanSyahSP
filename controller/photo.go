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

func UpdatePhoto (context *gin.Context) {
	user, err := helper.CurrentUser(context)

	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Dapatkan ID photo dari URL
	photoId := getPhotoIDFromRequest(context)
	// Dapatkan photo dari database

	photo, err := models.FindPhotoById(photoId)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Pastikan photo yang akan diupdate adalah milik user yang sedang login
	if photo.UserID != user.ID {
		context.JSON(http.StatusForbidden, gin.H{"error": "Not authorized to update this photo"})
		return
	}

	// Dapatkan input dari body request
	var input models.Photo
	if err := context.ShouldBindJSON(&input); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Update informasi photo sesuai dengan input
	photo.Title = input.Title
	photo.Caption = input.Caption
	photo.PhotoURL = input.PhotoURL

	// Lakukan operasi update ke database
	if err := photo.Update(); err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Kirim response ke client
	context.JSON(http.StatusOK, gin.H{"data": photo})
}

func getPhotoIDFromRequest(context *gin.Context) string {
    // Fungsi ini bergantung pada bagaimana Anda mengekstrak userId dari request.
    // Sesuaikan dengan kebutuhan dan implementasi endpoint Anda.
    // Contoh: Mendapatkan userId dari path parameter
    photoId := context.Param("photoId")
    return photoId
}

