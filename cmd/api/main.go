package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"ApiSmart/config"
	"ApiSmart/internal/adapters/handlers"
	"ApiSmart/internal/adapters/repositories/mysql"
	"ApiSmart/internal/core/services"
	"ApiSmart/pkg/database"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	cfg := config.LoadConfig()

	db, err := database.NewMySQLConnection(cfg.DBConfig)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	userRepo := mysql.NewUserRepository(db)
	sensorRepo := mysql.NewSensorRepository(db)

	authService := services.NewAuthService(userRepo)
	alertService := services.NewAlertService()
	sensorService := services.NewSensorService(sensorRepo, alertService)

	authHandler := handlers.NewAuthHandler(authService)
	sensorHandler := handlers.NewSensorHandler(sensorService)

	router := gin.Default()

	// Configurar middleware CORS con origen específico
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000", "http://127.0.0.1:8000"}, // Origen específico del frontend
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization", "X-Requested-With"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	router.POST("/api/register", authHandler.Register)
	router.POST("/api/login", authHandler.Login)

	router.POST("/sensores", sensorHandler.CreateSensorData)

	authorized := router.Group("/api")
	{
		authorized.GET("/sensors", sensorHandler.GetAllSensorData)
		authorized.GET("/sensors/latest", sensorHandler.GetLatestSensorData)
		authorized.GET("/sensors/alerts", sensorHandler.GetAlerts)
	}

	srv := &http.Server{
		Addr:    ":" + cfg.ServerPort,
		Handler: router,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	log.Printf("Server started on port %s", cfg.ServerPort)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server exited properly")
}
