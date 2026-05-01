package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"smartcity/internal/models"
	"smartcity/internal/repository"
)

type SensorHandler struct {
	Repo *repository.SensorRepository
}

func NewSensorHandler(repo *repository.SensorRepository) *SensorHandler {
	return &SensorHandler{Repo: repo}
}

func (h *SensorHandler) CreateSensor(w http.ResponseWriter, r *http.Request) {
	var req models.SensorRequest

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		log.Println("DECODE ERROR:", err)
		http.Error(w, "invalid JSON", http.StatusBadRequest)
		return
	}

	// validation
	if req.Metadata.SensorID == "" {
		http.Error(w, "missing sensor_id", http.StatusBadRequest)
		return
	}

	if len(req.Data) == 0 {
		http.Error(w, "missing data", http.StatusBadRequest)
		return
	}

	if req.Metadata.Time == "" {
		req.Metadata.Time = time.Now().Format(time.RFC3339)
	}

	err = h.Repo.InsertMany(req)
	if err != nil {
		http.Error(w, "DB insert failed", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{
		"status": "inserted",
	})
}

func (h *SensorHandler) CreateSensorsBatch(w http.ResponseWriter, r *http.Request) {
	var reqs []models.SensorRequest

	err := json.NewDecoder(r.Body).Decode(&reqs)
	if err != nil {
		http.Error(w, "invalid JSON", http.StatusBadRequest)
		return
	}

	if len(reqs) == 0 {
		http.Error(w, "empty batch", http.StatusBadRequest)
		return
	}

	for _, req := range reqs {

		// reuse existing validation
		if req.Metadata.SensorID == "" || len(req.Data) == 0 {
			continue
		}

		err := h.Repo.InsertMany(req)
		if err != nil {
			http.Error(w, "DB insert failed", http.StatusInternalServerError)
			return
		}
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{
		"status": "batch inserted",
	})
}

func (h *SensorHandler) GetLatest(w http.ResponseWriter, r *http.Request) {
	sensorID := r.URL.Query().Get("sensor_id")

	if sensorID == "" {
		http.Error(w, "missing sensor_id", http.StatusBadRequest)
		return
	}

	data, err := h.Repo.GetLatest(sensorID)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}

func (h *SensorHandler) GetHistory(w http.ResponseWriter, r *http.Request) {
	sensorID := r.URL.Query().Get("sensor_id")
	fromStr := r.URL.Query().Get("from")
	toStr := r.URL.Query().Get("to")

	if sensorID == "" {
		http.Error(w, "missing sensor_id", http.StatusBadRequest)
		return
	}

	from, err := time.Parse(time.RFC3339, fromStr)
	if err != nil {
		http.Error(w, "invalid from date", http.StatusBadRequest)
		return
	}

	to, err := time.Parse(time.RFC3339, toStr)
	if err != nil {
		http.Error(w, "invalid to date", http.StatusBadRequest)
		return
	}

	data, err := h.Repo.GetBetween(sensorID, from, to)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}

