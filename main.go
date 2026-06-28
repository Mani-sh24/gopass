package main

import (
	"log"
	"os"

	"example/web-service-gin/database"
	_ "example/web-service-gin/docs" // Import swagger docs
	"example/web-service-gin/handlers"
	"example/web-service-gin/helpers"
	"example/web-service-gin/middleware"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title           Gin Password Manager API
// @version         1.0
// @description     API for managing user registration, logins, and credential storage.
// @host            localhost:8080
// @BasePath        /
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

	// Swagger route
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	{
		auth := router.Group("/auth")
		auth.POST("/register", handlers.CreateUser)
		auth.POST("/login", handlers.Login)
	}
	{
		protected := router.Group("/protected")
		protected.Use(middleware.JWTMiddleware())
		protected.GET("/getuser", handlers.GetUserById)

	}
	{
		posts := router.Group("passwordcrud")
		posts.Use(middleware.JWTMiddleware())
		posts.POST("/create", handlers.CreatePassword)
		posts.GET("/all", handlers.GetAllPasswords)
		posts.DELETE("/delete/:id", handlers.DeletePassword)
		posts.PATCH("/update/:id", handlers.UpdatePassword)

	}

	router.Run(":" + port)
}
