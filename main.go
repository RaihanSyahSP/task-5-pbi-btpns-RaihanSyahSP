package main

import (
	"fmt"
	"log"
	"pbi-rakamin-go/controller"
	"pbi-rakamin-go/database"
	"pbi-rakamin-go/models"

	"github.com/gin-gonic/gin"

	"github.com/joho/godotenv"
)

func main() {
    loadEnv()
    loadDatabase()
    serveApplication()
}

func loadDatabase() {
    database.Connect()
    database.Database.AutoMigrate(&models.User{})
    database.Database.AutoMigrate(&models.Photo{})
}

func loadEnv() {
    err := godotenv.Load(".env.local")
    if err != nil {
        log.Fatal("Error loading .env file")
    }
}

func serveApplication() {
    router := gin.Default()

    publicRoutes := router.Group("/auth")
    publicRoutes.POST("/register", controller.Register)
    publicRoutes.POST("/login", controller.Login)

    router.Run(":8000")
    fmt.Println("Server running on port 8000")
}