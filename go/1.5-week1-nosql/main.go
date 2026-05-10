package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/FireChickenMP4/MuXi/go/1.5-week1-nosql/config"
	"github.com/FireChickenMP4/MuXi/go/1.5-week1-nosql/handler"
	mongorepo "github.com/FireChickenMP4/MuXi/go/1.5-week1-nosql/mongodb"
	pgrepo "github.com/FireChickenMP4/MuXi/go/1.5-week1-nosql/postgres"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	cfg := config.Load()

	mongoClient, err := connectMongo(cfg.MongoURI)
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}
	defer mongoClient.Disconnect(context.Background())
	log.Println("Connected to MongoDB")

	pgDB, err := connectPG(cfg.PgDSN)
	if err != nil {
		log.Fatalf("Failed to connect to PostgreSQL: %v", err)
	}
	log.Println("Connected to PostgreSQL")

	if err := pgrepo.AutoMigrate(pgDB); err != nil {
		log.Fatalf("Failed to migrate PostgreSQL: %v", err)
	}
	log.Println("PostgreSQL migrated")

	mux := http.NewServeMux()

	mongoRepo := mongorepo.NewRepo(mongoClient, cfg.MongoDB)
	mongoHandler := handler.New(mongoRepo)
	mux.Handle("/api/mongo/", cors(http.StripPrefix("/api/mongo", mongoHandler)))

	pgRepo := pgrepo.NewRepo(pgDB)
	pgHandler := handler.New(pgRepo)
	mux.Handle("/api/pg/", cors(http.StripPrefix("/api/pg", pgHandler)))

	mux.Handle("/", cors(http.FileServer(http.Dir("."))))

	addr := fmt.Sprintf(":%s", cfg.ServerPort)
	log.Printf("Server starting on %s", addr)
	log.Printf("  Frontend:     http://localhost%s/", addr)
	log.Printf("  MongoDB API:  http://localhost%s/api/mongo/posts", addr)
	log.Printf("  PostgreSQL API: http://localhost%s/api/pg/posts", addr)
	if err := http.ListenAndServe(addr, mux); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}

func cors(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func connectMongo(uri string) (*mongo.Client, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		return nil, err
	}
	if err := client.Ping(ctx, nil); err != nil {
		return nil, err
	}
	return client, nil
}

func connectPG(dsn string) (*gorm.DB, error) {
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}
	if err := sqlDB.Ping(); err != nil {
		return nil, err
	}
	return db, nil
}
