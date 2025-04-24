package api

import (
	"os"

	"github.com/KaungHtetHein116/personal-task-manager/api/middleware"
	v1 "github.com/KaungHtetHein116/personal-task-manager/api/v1"
	"github.com/KaungHtetHein116/personal-task-manager/config"
	redisdb "github.com/KaungHtetHein116/personal-task-manager/internal/redis"
	"github.com/KaungHtetHein116/personal-task-manager/utils"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	log "github.com/sirupsen/logrus"
)

func init() {
	log.SetFormatter(&log.TextFormatter{
		FullTimestamp: true,
	})
	log.SetLevel(log.InfoLevel)
}

func StartServer() {
	db := config.ConnectDB()
	redisdb.InitRedis()

	e := echo.New()
	e.Validator = &utils.CustomValidator{Validator: validator.New()}
	e.HTTPErrorHandler = utils.CustomHTTPErrorHandler

	middleware.RegisterBasicMiddlewares(e)
	middleware.AuthMiddleware(e)

	v1.RegisterUserRoute(e, db)
	v1.RegisterProjectRoute(e, db)

	port := ":" + os.Getenv("APP_PORT")

	log.Fatal(e.Start(port))
}
