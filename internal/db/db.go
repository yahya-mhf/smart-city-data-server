package db

import (
	"context"
	"log"

	"github.com/jackc/pgx/v5"
)

var DB *pgx.Conn

func Connect() {
	connStr := "host=localhost port=5432 user=postgres password=postgres dbname=mydb sslmode=disable"
	conn, err := pgx.Connect(context.Background(), connStr)
	if err != nil {
		log.Fatal("DB connection failed:", err)
	}

	DB = conn
	log.Println("✅ Connected to TimescaleDB")
}

func Close() {
	if DB != nil {
		DB.Close(context.Background())
	}
}