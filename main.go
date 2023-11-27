package main

import (
	"log"
	"pbi-rakamin-go/database"
	"pbi-rakamin-go/models"

	"github.com/joho/godotenv"
)

func main() {
    loadEnv()
    loadDatabase()
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
