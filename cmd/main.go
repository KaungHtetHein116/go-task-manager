package main

import (
	"log"
	"os"

	"github.com/KaungHtetHein116/personal-task-manager/database"
	"github.com/KaungHtetHein116/personal-task-manager/repository"
	"github.com/KaungHtetHein116/personal-task-manager/routes"
	"github.com/KaungHtetHein116/personal-task-manager/service"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
)

func main() {
	// Load environment variables from .env file
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	db := database.ConnectDB()

	e := echo.New()

	authRepo := repository.NewUserRepository(db)
	authHandler := service.NewAuthService(authRepo)

	routes.SetupRoutes(e, authHandler)

	// Access the environment variable
	port := os.Getenv("APP_PORT")
	if port == "" {
		port = "8080" // Default port if not set in environment
	}

	e.Logger.Fatal(e.Start(":" + port))
}
