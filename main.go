package main

import (
	"log"
	"os"

	"example/web-service-gin/database"
	"example/web-service-gin/handlers"
	"example/web-service-gin/helpers"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	// Load environment variables from .env file if it exists
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using system environment variables")
	}
	secret := os.Getenv("JWT_TOK")

	if err := helpers.InitJWT(secret); err != nil {
		log.Fatal(err)
	}
	database.Connect_to_db()
	database.Init()

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	router := gin.Default()
	{
		auth := router.Group("/auth")
		auth.POST("/register", handlers.CreateUser)
		auth.GET("/getuser/:id", handlers.GetUserById)
		auth.POST("/login", handlers.Login)
	}
	router.POST("/all", handlers.GetAllUsers)

	router.Run(":" + port)
}
