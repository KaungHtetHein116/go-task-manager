package main

import (
	"fmt"
	"os"

	"github.com/KaungHtetHein116/personal-task-manager/database"
	"github.com/KaungHtetHein116/personal-task-manager/handlers"
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
)

var logger = logrus.New()

func main() {
	// Initialize and migrate the database
	database.InitDB()
	database.Migrate()

	// Create a new Echo instance
	e := echo.New()

	// Register routes
	e.POST("/register", handlers.RegisterUser)
	e.POST("/login", handlers.LoginUser)

	// Middleware for logging requests
	e.Use(loggingMiddleware)

	// Start the server
	appPort := os.Getenv("APP_PORT")
	if appPort == "" {
		logger.Fatal("APP_PORT environment variable is not set")
	}
	e.Logger.Fatal(e.Start(fmt.Sprintf(":%s", appPort)))
}

// loggingMiddleware logs incoming requests and their responses
func loggingMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		logger.WithFields(logrus.Fields{
			"method": c.Request().Method,
			"path":   c.Request().URL.Path,
		}).Info("Incoming request")

		err := next(c)

		logger.WithFields(logrus.Fields{
			"method":     c.Request().Method,
			"path":       c.Request().URL.Path,
			"status":     c.Response().Status,
			"user_ip":    c.RealIP(),
			"user_agent": c.Request().UserAgent(),
		}).Info("Request completed")

		return err
	}
}
