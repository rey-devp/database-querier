package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	fiberLogger "github.com/gofiber/fiber/v2/middleware/logger"

	"database-querier-agent/internal/agent"
	"database-querier-agent/internal/config"
	"database-querier-agent/internal/logger"
	"database-querier-agent/internal/memory"
	"database-querier-agent/internal/mongodb"
	"database-querier-agent/internal/service"
)

func main() {
	logger.Info("STARTUP", "Loading configuration...")

	// 1. Load config
	cfg := config.LoadConfig()
	logger.Info("STARTUP", "Config loaded", "port", cfg.ServerPort, "database", cfg.DatabaseName)

	// 2. Connect to MongoDB Atlas
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	mongoClient, err := mongodb.NewClient(ctx, cfg.MongoURI, cfg.DatabaseName)
	if err != nil {
		logger.Error("STARTUP", "Failed to connect to MongoDB", "error", err.Error())
		os.Exit(1)
	}
	defer func() {
		if err := mongoClient.Close(context.Background()); err != nil {
			logger.Error("SHUTDOWN", "Error disconnecting MongoDB", "error", err.Error())
		}
	}()
	logger.Info("STARTUP", "Connected to MongoDB Atlas", "database", cfg.DatabaseName)

	// 3. Initialize components
	store := memory.NewStore()
	dbAgent := agent.NewAgent(mongoClient, store)
	handler := service.NewHandler(dbAgent, store)

	// 4. Setup Fiber app
	app := fiber.New(fiber.Config{
		AppName: "Database Querier Agent",
	})

	// Add CORS middleware
	app.Use(cors.New(cors.Config{
		AllowOrigins: "https://jokitugas.bananaunion.web.id",
		AllowMethods: "POST,GET,OPTIONS",
		AllowHeaders: "Content-Type,Authorization",
	}))

	// Add Fiber Logger middleware
	app.Use(fiberLogger.New(fiberLogger.Config{
		Format: "| ${status} | ${method} | ${path} | ${latency} |\n",
	}))

	// Setup routes
	app.Post("/query", handler.HandleQuery)
	app.Get("/health", handler.HandleHealth)

	// 5. Start Server
	serverAddr := ":" + cfg.ServerPort
	
	go func() {
		logger.Info("STARTUP", "Server listening", "port", cfg.ServerPort)
		if err := app.Listen(serverAddr); err != nil {
			logger.Error("STARTUP", "Server error", "error", err.Error())
		}
	}()

	// 6. Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	logger.Info("SHUTDOWN", "Shutting down server...")

	if err := app.ShutdownWithTimeout(5 * time.Second); err != nil {
		logger.Error("SHUTDOWN", "Server forced to shutdown", "error", err.Error())
	}

	logger.Info("SHUTDOWN", "Server exiting")
}
