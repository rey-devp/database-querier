package api

import (
	"context"
	"net/http"
	"time"

	"github.com/gofiber/adaptor/v2"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"

	"database-querier-agent/internal/agent"
	"database-querier-agent/internal/config"
	"database-querier-agent/internal/logger"
	"database-querier-agent/internal/memory"
	"database-querier-agent/internal/mongodb"
	"database-querier-agent/internal/service"
)

var fiberHandler http.HandlerFunc

func init() {
	logger.Info("VERCEL", "Initializing Serverless Function...")

	cfg := config.LoadConfig()

	// Connect to MongoDB
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	mongoClient, err := mongodb.NewClient(ctx, cfg.MongoURI, cfg.DatabaseName)
	if err != nil {
		logger.Error("VERCEL", "Failed to connect to MongoDB", "error", err.Error())
		// We shouldn't os.Exit(1) in a serverless environment, just log it.
	}

	// Initialize components
	store := memory.NewStore()
	dbAgent := agent.NewAgent(mongoClient, store)
	handler := service.NewHandler(dbAgent, store)

	// Setup Fiber app
	app := fiber.New(fiber.Config{
		AppName: "Database Querier Agent (Serverless)",
	})

	// Add CORS middleware
	app.Use(cors.New(cors.Config{
		AllowOrigins: "https://jokitugas.bananaunion.web.id",
		AllowMethods: "POST,GET,OPTIONS",
		AllowHeaders: "Content-Type,Authorization",
	}))

	// Setup routes
	app.Post("/query", handler.HandleQuery)
	app.Get("/health", handler.HandleHealth)

	// Convert Fiber app to http.HandlerFunc
	fiberHandler = adaptor.FiberApp(app)
}

// Handler is the Vercel serverless entrypoint
func Handler(w http.ResponseWriter, r *http.Request) {
	fiberHandler(w, r)
}
