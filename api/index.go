package api

import (
	"context"
	"net/http"
	"time"

	"github.com/gofiber/adaptor/v2"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"

	"database-querier-agent/pkg/agent"
	"database-querier-agent/pkg/config"
	"database-querier-agent/pkg/logger"
	"database-querier-agent/pkg/memory"
	"database-querier-agent/pkg/mongodb"
	"database-querier-agent/pkg/service"
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
		panic("Gagal terhubung ke MongoDB saat inisialisasi: " + err.Error())
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
