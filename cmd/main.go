package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"pvz/internal/api/handler"
	"pvz/internal/db"
	"pvz/internal/logger"
	"pvz/internal/repository"
	"pvz/internal/service"
	"pvz/server"

	"github.com/spf13/viper"
)

func main() {
	gin.SetMode(gin.ReleaseMode)

	if err := logger.Init(); err != nil {
		log.Fatalf("Logger initialization error: %v", err)
	}
	defer logger.Logger.Sync()
	logger.SugaredLogger.Info("The application is running")

	if err := initConfig(); err != nil {
		logger.SugaredLogger.Fatalw("Error initializing configs", "error", err)
	}

	if err := godotenv.Load(); err != nil {
		logger.SugaredLogger.Fatalw("Error loading env file", "error", err)
	}

	postgresDb, err := db.NewPostgresDB(db.Config{
		Host:     viper.GetString("db.host"),
		Port:     viper.GetString("db.port"),
		Username: viper.GetString("db.username"),
		Password: os.Getenv("PostgresPassword"),
		DBName:   viper.GetString("db.dbname"),
		SSLMode:  viper.GetString("db.sslmode"),
	})
	if err != nil {
		logger.SugaredLogger.Fatalw("Failed initializing DB", "error", err)
	}

	repos := repository.NewRepositore(postgresDb)
	services := service.NewService(repos)
	handlers := handler.NewHandler(services)

	srv := new(server.Server)
	if err := srv.Run(viper.GetString("port"), handlers.InitRoutes()); err != nil {
		logger.SugaredLogger.Fatalw("Error occurred while running server", "error", err)
	}
}

func initConfig() error {
	viper.AddConfigPath("config")
	viper.SetConfigName("config")
	return viper.ReadInConfig()
}
