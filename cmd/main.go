package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	_ "github.com/lib/pq"

	"github.com/joho/godotenv"
	"github.com/renlin-code/mock-shop-api/pkg/handler"
	"github.com/renlin-code/mock-shop-api/pkg/repository"
	"github.com/renlin-code/mock-shop-api/pkg/service"
	"github.com/renlin-code/mock-shop-api/pkg/storage"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func main() {
	logrus.SetFormatter(new(logrus.JSONFormatter))
	if err := initConfig(); err != nil {
		logrus.Fatalf("Error initializating configs: %s", err.Error())
	}
	if err := godotenv.Load(); err != nil {
		logrus.Fatalf("Error loading env variables: %s", err.Error())
	}

	db, err := repository.NewPostgresDB(repository.Config{
		Host:     viper.GetString("db.host"),
		Port:     viper.GetString("db.port"),
		Username: viper.GetString("db.username"),
		Password: os.Getenv("DB_PASSWORD"),
		DBName:   viper.GetString("db.dbname"),
		SSLMode:  viper.GetString("db.sslmode"),
	})
	if err != nil {
		logrus.Fatalf("Failed to initialize database: %s", err.Error())
	}

	fsStorage := storage.NewFileSystemStorage(storage.Config{
		BaseUrl: viper.GetString("server.base_url"),
	})
	storage := storage.NewStorage(fsStorage)
	repos := repository.NewRepository(db, storage)
	services := service.NewService(repos)
	handlers := handler.NewHandler(services)

	srv := new(handler.Server)

	go func() {
		if err := srv.Run(viper.GetString("port"), handlers.InitRoutes()); err != nil {
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

func initConfig() error {
	viper.AddConfigPath("configs")
	viper.SetConfigName("config")
	return viper.ReadInConfig()
}
