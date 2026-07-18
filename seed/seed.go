package main

import (
	"context"
	"log"
	"time"

	"database-querier-agent/internal/config"
	"database-querier-agent/internal/mongodb"
	"go.mongodb.org/mongo-driver/v2/bson"
)

func main() {
	cfg := config.LoadConfig()
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongodb.NewClient(ctx, cfg.MongoURI, cfg.DatabaseName)
	if err != nil {
		log.Fatal(err)
	}
	defer client.Close(context.Background())

	coll := client.GetDatabase().Collection("students")

	// Clean up existing data
	coll.DeleteMany(ctx, bson.M{})

	docs := []interface{}{
		// Mahasiswa Semester 6 (seperti di contoh payload)
		bson.M{"name": "Ahmad", "semester": 6, "gpa": 3.8, "major": "Informatika"},
		bson.M{"name": "Budi", "semester": 6, "gpa": 3.2, "major": "Sistem Informasi"},
		bson.M{"name": "Diana", "semester": 6, "gpa": 3.5, "major": "Teknik Komputer"},
		
		// Mahasiswa Semester Lainnya
		bson.M{"name": "Citra", "semester": 4, "gpa": 3.9, "major": "Informatika"},
		bson.M{"name": "Eko", "semester": 2, "gpa": 3.1, "major": "Informatika"},
		bson.M{"name": "Fahmi", "semester": 4, "gpa": 2.8, "major": "Teknik Komputer"},
		bson.M{"name": "Gita", "semester": 2, "gpa": 3.4, "major": "Sistem Informasi"},
		bson.M{"name": "Hadi", "semester": 8, "gpa": 3.7, "major": "Informatika"},
		bson.M{"name": "Indah", "semester": 8, "gpa": 3.9, "major": "Sistem Informasi"},
		bson.M{"name": "Joko", "semester": 4, "gpa": 2.9, "major": "Teknik Komputer"},
	}

	res, err := coll.InsertMany(ctx, docs)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Successfully inserted %d documents", len(res.InsertedIDs))
}
