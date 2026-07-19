package main

import (
	"context"
	"log"
	"time"

	"database-querier-agent/pkg/config"
	"database-querier-agent/pkg/mongodb"
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
		// Mahasiswa Semester 6
		bson.M{"name": "Ahmad", "semester": 6, "gpa": 3.8, "major": "Informatika"},
		bson.M{"name": "Budi", "semester": 6, "gpa": 3.2, "major": "Sistem Informasi"},
		bson.M{"name": "Diana", "semester": 6, "gpa": 3.5, "major": "Teknik Komputer"},
		bson.M{"name": "Reza", "semester": 6, "gpa": 3.9, "major": "Informatika"},
		bson.M{"name": "Siti", "semester": 6, "gpa": 3.4, "major": "Teknik Komputer"},
		bson.M{"name": "Tono", "semester": 6, "gpa": 3.1, "major": "Sistem Informasi"},
		
		// Mahasiswa Semester Lainnya
		bson.M{"name": "Citra", "semester": 4, "gpa": 3.9, "major": "Informatika"},
		bson.M{"name": "Eko", "semester": 2, "gpa": 3.1, "major": "Informatika"},
		bson.M{"name": "Fahmi", "semester": 4, "gpa": 2.8, "major": "Teknik Komputer"},
		bson.M{"name": "Gita", "semester": 2, "gpa": 3.4, "major": "Sistem Informasi"},
		bson.M{"name": "Hadi", "semester": 8, "gpa": 3.7, "major": "Informatika"},
		bson.M{"name": "Indah", "semester": 8, "gpa": 3.9, "major": "Sistem Informasi"},
		bson.M{"name": "Joko", "semester": 4, "gpa": 2.9, "major": "Teknik Komputer"},
		bson.M{"name": "Kiki", "semester": 2, "gpa": 3.6, "major": "Informatika"},
		bson.M{"name": "Lina", "semester": 8, "gpa": 3.3, "major": "Sistem Informasi"},
		bson.M{"name": "Mira", "semester": 4, "gpa": 3.5, "major": "Informatika"},
		bson.M{"name": "Nina", "semester": 2, "gpa": 3.8, "major": "Teknik Komputer"},
		bson.M{"name": "Oki", "semester": 8, "gpa": 3.0, "major": "Informatika"},
		bson.M{"name": "Putri", "semester": 4, "gpa": 3.2, "major": "Sistem Informasi"},
		bson.M{"name": "Qori", "semester": 2, "gpa": 3.7, "major": "Teknik Komputer"},
	}

	res, err := coll.InsertMany(ctx, docs)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Successfully inserted %d documents", len(res.InsertedIDs))
}
