package db

import (
	"context"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
)

// ConnectDB - connects to a given postgres database, via pgxpool (for simultaneous connections)
func ConnectDB(dsn string) (*pgxpool.Pool, error) {
	dbPool, err := pgxpool.New(context.Background(), dsn)
	if err != nil {
		return nil, err
	}

	// Test connection
	err = dbPool.Ping(context.Background())
	if err != nil {
		return nil, err
	}
	log.Println("Connected to Postgres DB!")

	return dbPool, nil
}

func GetBooks(conn *pgxpool.Conn) ([]Book, error) {

}
