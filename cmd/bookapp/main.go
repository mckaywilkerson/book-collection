package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/mckaywilkerson/book-collection/internal/db"
	"github.com/mckaywilkerson/book-collection/internal/handlers"
)

func main() {
	// connect to DB
	dsn := getDSN()
	pgxPool, err := db.ConnectDB(dsn)
	if err != nil {
		log.Fatalf("Could not connect to DB: %f", err)
	}
	defer pgxPool.Close()

	// setup router
	router := gin.Default()

	router.GET("/books", handlers.HandleGetAllBooks(pgxPool))

	log.Println("Starting server on :8081")
	router.Run("localhost:8081")

}

// getDSN - read environment variables orr fallback to deefaults
func getDSN() string {
	dbUser := os.Getenv("DB_USER")
	if dbUser == "" {
		dbUser = "postgres"
	}

	dbPassword := os.Getenv("DB_PASSWORD")
	if dbPassword == "" {
		dbPassword = "postgres"
	}

	dbHost := os.Getenv("DB_HOST")
	if dbHost == "" {
		dbHost = "localhost"
	}

	dbName := os.Getenv("DB_NAME")
	if dbName == "" {
		dbName = "books"
	}

	return "postgres://" + dbUser + ":" + dbPassword + "@" + dbHost + ":5432/" + dbName + "?sslmode=disable"
}
