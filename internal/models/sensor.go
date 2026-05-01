// internal/models/sensor.go
package models

import "time"

// type SensorRequest struct {
// 	Metadata struct {
// 		SensorID  string  `json:"sensor_id"`
// 		Time      string  `json:"time"`
// 		Latitude  float64 `json:"latitude"`
// 		Longitude float64 `json:"longitude"`
// 	} `json:"metadata"`

// 	Data map[string]float64 `json:"data"`
// }

type SensorRequest struct {
	Metadata Metadata         `json:"metadata"`
	Data     map[string]float64 `json:"data"`
}

type Metadata struct {
	SensorID  string  `json:"sensor_id"`
	Time      string  `json:"time"`
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

type SensorEntity struct {
	SensorID  string
	Time      time.Time
	Latitude  float64
	Longitude float64
	Variable  string
	Value     float64
}
