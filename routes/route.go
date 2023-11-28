package routes

import (
	"pbi-rakamin-go/controller"
	"pbi-rakamin-go/middleware"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(router *gin.Engine) {
	publicRoutes := router.Group("/auth")
	{
		publicRoutes.POST("/register", controller.Register)
		publicRoutes.POST("/login", controller.Login)
	}

	protectedRoutes := router.Group("/api")
	protectedRoutes.Use(middleware.JWTAuthMiddleware())
	{
		protectedRoutes.POST("/photos", controller.AddPhoto)
		protectedRoutes.GET("/photos", controller.GetAllPhotos)
		protectedRoutes.PUT("/photos/:photoId", controller.UpdatePhoto)
		protectedRoutes.DELETE("/photos/:photoId", controller.DeletePhoto)

		protectedRoutes.PUT("/users/:userId", controller.UpdateUser)
		protectedRoutes.DELETE("/users/:userId", controller.DeleteUser)
	}
}
