package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/joseenriqe97/test-cabri/config"
	"github.com/joseenriqe97/test-cabri/pkg/application"
	"github.com/joseenriqe97/test-cabri/pkg/infrastructure"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/sirupsen/logrus"
)

func main() {

	err := config.ReadConfig()
	if err != nil {
		log.Fatal(err, "Error reading .env file!")
	}

	appController := infrastructure.AppController{
		User: application.NewUserApplication(),
	}

	e := echo.New()
	e = infrastructure.NewRouter(e, appController)
	e.Use(middleware.CORS())

	go func() {
		if err := e.Start(config.HTTPListener()); err != nil {
			logrus.WithError(err).Error("shutting down server")
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit
}
