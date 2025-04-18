package main

import (
	"log"
	"os"

	"pvz/internal/api/handler"
	"pvz/internal/db"
	"pvz/internal/logger"
	"pvz/internal/repository"
	"pvz/internal/service"
	"pvz/metrics"
	"pvz/server"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"

	"github.com/spf13/viper"
)

func main() {
	gin.SetMode(gin.ReleaseMode)
	metrics.Init()

	// Инициализация логгера
	if err := logger.Init(); err != nil {
		log.Fatalf("Logger initialization error: %v", err)
	}
	defer logger.Log.Sync()

	logger.Log.Infow("The application is running")

	// Загрузка конфигурации
	if err := initConfig(); err != nil {
		logger.Log.Fatalw("Error initializing configs", "error", err)
	}

	// Загрузка .env файла
	if err := godotenv.Load(); err != nil {
		logger.Log.Fatalw("Error loading env file", "error", err)
	}

	// Инициализация БД
	postgresDb, err := db.NewPostgresDB(db.Config{
		Host:     os.Getenv("DB_HOST"),
		Port:     os.Getenv("DB_PORT"),
		Username: os.Getenv("POSTGRES_USER"),
		Password: os.Getenv("POSTGRES_PASSWORD"),
		DBName:   os.Getenv("POSTGRES_DB"),
		SSLMode:  os.Getenv("SSL_MODE"),
	})
	if err != nil {
		logger.Log.Fatalw("Failed initializing DB", "error", err)
	}

	// Инициализация слоев приложения
	repos := repository.NewRepository(postgresDb, logger.Log)
	services := service.NewService(repos, logger.Log)
	handlers := handler.NewHandler(services, logger.Log)

	// Запуск сервера
	srv := new(server.Server)
	if err := srv.Run(viper.GetString("port"), handlers.InitRoutes()); err != nil {
		logger.Log.Fatalw("Error occurred while running server", "error", err)
	}
}

func initConfig() error {
	viper.AddConfigPath("config")
	viper.SetConfigName("config")
	return viper.ReadInConfig()
}
