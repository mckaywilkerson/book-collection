package handlers

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/mckaywilkerson/book-collection/internal/db"
)

// HandleGetAllBooks - returns all books from the database of given pool
func HandleGetAllBooks(dbPool *pgxpool.Pool) gin.HandlerFunc {
	return func(c *gin.Context) {
		books, err := db.GetAllBooks(dbPool)
		if err != nil {
			log.Printf("Failed to load all books from DB: %v", err)
			c.IndentedJSON(http.StatusInternalServerError, books)
			return
		}

		c.IndentedJSON(http.StatusOK, books)
	}
}
