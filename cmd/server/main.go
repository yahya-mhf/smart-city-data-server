package main

import (
	"log"
	"net/http"

	"smartcity/internal/db"
	"smartcity/internal/handlers"
	"smartcity/internal/repository"
)

func main() {
	db.Connect()
	defer db.Close()

	repo := repository.NewSensorRepository(db.DB)
	handler := handlers.NewSensorHandler(repo)

	mux := http.NewServeMux()

	// ingestion
	mux.HandleFunc("/api/v1/sensors/data", handler.CreateSensor)
	mux.HandleFunc("/api/v1/sensors/batch", handler.CreateSensorsBatch)

	// read APIs
	mux.HandleFunc("/api/v1/sensors/latest", handler.GetLatest)
	mux.HandleFunc("/api/v1/sensors/history", handler.GetHistory)

	// health
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("OK"))
	})

	log.Println("🚀 Server running on :8080")
	log.Fatal(http.ListenAndServe(":8080", mux))
}