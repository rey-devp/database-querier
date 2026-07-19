package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	MongoURI     string
	DatabaseName string
	ServerPort   string
}

func LoadConfig() *Config {
	// Try loading from .env file (will not exist in Vercel, that's OK)
	err := godotenv.Load()
	if err != nil {
		_ = godotenv.Load("../.env")
	}

	mongoURI := os.Getenv("MONGO_DBQ")
	if mongoURI == "" {
		log.Println("[CONFIG] MONGO_DBQ is empty! Available env vars check...")
		panic("MONGO_DBQ environment variable must be set")
	}

	databaseName := os.Getenv("DATABASE_NAME")
	if databaseName == "" {
		log.Println("[CONFIG] DATABASE_NAME is empty!")
		panic("DATABASE_NAME environment variable must be set")
	}

	serverPort := os.Getenv("SERVER_PORT")
	if serverPort == "" {
		serverPort = "8080" // Default port
	}

	log.Printf("[CONFIG] Loaded: DB=%s, Port=%s\n", databaseName, serverPort)

	return &Config{
		MongoURI:     mongoURI,
		DatabaseName: databaseName,
		ServerPort:   serverPort,
	}
}

