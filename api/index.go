package api

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/gofiber/adaptor/v2"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"

	"database-querier-agent/pkg/agent"
	"database-querier-agent/pkg/config"
	"database-querier-agent/pkg/memory"
	"database-querier-agent/pkg/mongodb"
	"database-querier-agent/pkg/service"
)

var (
	once         sync.Once
	fiberHandler http.HandlerFunc
	initErr      error
)

func setupApp() {
	log.Println("[VERCEL] Starting lazy initialization...")

	// 1. Load config
	cfg := config.LoadConfig()
	log.Printf("[VERCEL] Config loaded: DB=%s\n", cfg.DatabaseName)

	// 2. Connect to MongoDB with longer timeout for cold start
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	log.Println("[VERCEL] Connecting to MongoDB...")
	mongoClient, err := mongodb.NewClient(ctx, cfg.MongoURI, cfg.DatabaseName)
	if err != nil {
		initErr = fmt.Errorf("MongoDB connection failed: %w", err)
		log.Printf("[VERCEL] ERROR: %s\n", initErr.Error())
		return
	}
	log.Println("[VERCEL] MongoDB connected successfully!")

	// 3. Initialize components
	store := memory.NewStore()
	dbAgent := agent.NewAgent(mongoClient, store)
	handler := service.NewHandler(dbAgent, store)

	// 4. Setup Fiber app
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

	// 5. Convert Fiber app to http.HandlerFunc
	fiberHandler = adaptor.FiberApp(app)
	log.Println("[VERCEL] Initialization complete!")
}

// Handler is the Vercel serverless entrypoint
func Handler(w http.ResponseWriter, r *http.Request) {
	// Lazy init: only connect to MongoDB on first request
	once.Do(setupApp)

	// If initialization failed, return a proper HTTP error (not a crash)
	if initErr != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"status":  "error",
			"task_id": "",
			"data":    nil,
			"message": "Agent initialization failed: " + initErr.Error(),
		})
		return
	}

	fiberHandler(w, r)
}
