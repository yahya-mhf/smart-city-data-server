// internal/repository/sensor_repository.go
package repository

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5"
	"smartcity/internal/models"
)

type SensorRepository struct {
	DB *pgx.Conn
}

func NewSensorRepository(db *pgx.Conn) *SensorRepository {
	return &SensorRepository{DB: db}
}

func (r *SensorRepository) InsertMany(req models.SensorRequest) error {

	for variable, value := range req.Data {

		if variable == "" {
			continue
		}

		_, err := r.DB.Exec(context.Background(),
			`INSERT INTO sensor_entity
			(sensor_id, time, latitude, longitude, variable, value)
			VALUES ($1, $2, $3, $4, $5, $6)`,
			req.Metadata.SensorID,
			req.Metadata.Time, // now handled in handler OR DB default
			req.Metadata.Latitude,
			req.Metadata.Longitude,
			variable,
			value,
		)

		if err != nil {
			return err
		}
	}

	return nil
}

func (r *SensorRepository) GetLatest(sensorID string) ([]models.SensorEntity, error) {
	rows, err := r.DB.Query(context.Background(), `
		SELECT sensor_id, time, latitude, longitude, variable, value
		FROM sensor_entity
		WHERE sensor_id=$1
		ORDER BY time DESC
		LIMIT 50
	`, sensorID)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []models.SensorEntity

	for rows.Next() {
		var s models.SensorEntity
		err := rows.Scan(
			&s.SensorID,
			&s.Time,
			&s.Latitude,
			&s.Longitude,
			&s.Variable,
			&s.Value,
		)
		if err != nil {
			return nil, err
		}
		result = append(result, s)
	}

	return result, nil
}

func (r *SensorRepository) GetBetween(sensorID string, from, to time.Time) ([]models.SensorEntity, error) {
	rows, err := r.DB.Query(context.Background(), `
		SELECT sensor_id, time, latitude, longitude, variable, value
		FROM sensor_entity
		WHERE sensor_id=$1
		AND time BETWEEN $2 AND $3
		ORDER BY time ASC
	`, sensorID, from, to)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []models.SensorEntity

	for rows.Next() {
		var s models.SensorEntity
		rows.Scan(&s.SensorID, &s.Time, &s.Latitude, &s.Longitude, &s.Variable, &s.Value)
		result = append(result, s)
	}

	return result, nil
}
