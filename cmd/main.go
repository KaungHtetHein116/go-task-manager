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
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	// 1. Environment variables must be loaded first
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	// 2. Database connection must be established before creating repositories
	db := database.ConnectDB()

	// 3. Echo instance must be created before adding routes
	e := echo.New()
	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: `${time_rfc3339} ${status} ${method} ${uri} ${latency_human}` + "\n",
	}))
	// 4. Repositories must be created before handlers
	authRepo := repository.NewUserRepository(db)

	// 5. Handlers must be created before setting up routes
	authHandler := service.NewAuthService(authRepo)

	// 6. Routes must be set up before starting the server
	routes.SetupRoutes(e, authHandler)

	// 7. Port must be configured before starting the server
	port := os.Getenv("APP_PORT")
	if port == "" {
		port = "8080"
	}

	// 8. Server must be started last
	e.Logger.Fatal(e.Start(":" + port))
}
