package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	_ "github.com/lib/pq"

	"github.com/renlin-code/mock-shop-api/pkg/handler"
	"github.com/renlin-code/mock-shop-api/pkg/repository"
	"github.com/renlin-code/mock-shop-api/pkg/service"
	"github.com/renlin-code/mock-shop-api/pkg/storage"
	"github.com/sirupsen/logrus"
)

// @title Mock Shop API
// @version 1.0
// @description API Service for a Mock Online Shop

// @host localhost:8020
// @BasePath /

// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
func main() {
	logrus.SetFormatter(new(logrus.JSONFormatter))
	db, err := repository.NewPostgresDB(repository.Config{
		Host:     os.Getenv("POSTGRES_HOST"),
		Port:     os.Getenv("POSTGRES_PORT"),
		Username: os.Getenv("POSTGRES_USERNAME"),
		Password: os.Getenv("POSTGRES_PASSWORD"),
		DBName:   os.Getenv("POSTGRES_NAME"),
		SSLMode:  os.Getenv("POSTGRES_SSLMODE"),
	})
	if err != nil {
		logrus.Fatalf("Failed to initialize database: %s", err.Error())
	}

	fsStorage := storage.NewFileSystemStorage(storage.Config{
		MediaBaseUrl: os.Getenv("APP_MEDIA_BASE_URL"),
	})
	storage := storage.NewStorage(fsStorage)
	repos := repository.NewRepository(db, storage)
	services := service.NewService(repos)
	handlers := handler.NewHandler(services)

	srv := new(handler.Server)

	go func() {
		if err := srv.Run(os.Getenv("APP_PORT"), handlers.InitRoutes()); err != nil {
			logrus.Fatalf("Error occurred while running http server: %s", err.Error())
		}
	}()

	logrus.Print("App started...")
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
	<-quit
	logrus.Print("App shutting down...")

	if err := srv.ShutDown(context.Background()); err != nil {
		logrus.Errorf("Error occurred on server while shutting down: %s", err.Error())
	}

	if err := db.Close(); err != nil {
		logrus.Errorf("Error occurred on db connection while closing: %s", err.Error())
	}

}
