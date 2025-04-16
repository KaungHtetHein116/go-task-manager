package main

import (
	"os"

	"github.com/KaungHtetHein116/personal-task-manager/database"
	redisdb "github.com/KaungHtetHein116/personal-task-manager/internal/redis"
	"github.com/KaungHtetHein116/personal-task-manager/repository"
	"github.com/KaungHtetHein116/personal-task-manager/routes"
	"github.com/KaungHtetHein116/personal-task-manager/service"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	log "github.com/sirupsen/logrus"
)

func init() {
	// Configure logrus format and level
	log.SetFormatter(&log.TextFormatter{
		FullTimestamp: true,
	})
	log.SetLevel(log.InfoLevel)
}

func main() {
	// 1. Environment variables must be loaded first
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	// 2. Database connection must be established before creating repositories
	db := database.ConnectDB()
	redisdb.InitRedis()

	// 3. Echo instance must be created before adding routes
	e := echo.New()
	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: `${time_rfc3339} ${status} ${method} ${uri} ${latency_human}` + "\n",
	}))

	// 4. Repositories must be created before handlers
	authRepo := repository.NewUserRepository(db)
	projectRepo := repository.NewProjectRepository(db)

	// 5. Handlers must be created before setting up routes
	authHandler := service.NewAuthService(authRepo)
	projectHandler := service.NewProjectService(projectRepo)

	// 6. Routes must be set up before starting the server
	routes.SetupRoutes(e, authHandler, projectHandler)

	// 7. Port must be configured before starting the server
	port := os.Getenv("APP_PORT")
	if port == "" {
		log.Fatal("APP_PORT is not set in the environment variables")
	}
	log.Infof("Starting server on port %s...", port)

	// 8. Server must be started last
	if err := e.Start(":" + port); err != nil {
		log.WithError(err).Fatal("Failed to start server")
	}
}
