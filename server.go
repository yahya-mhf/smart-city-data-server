package main

import (
	"fmt"
	"encoding/json"
	"net/http"
)

type Sensor struct{
	SensorId string `json:"sensorId"`
	Time string `json:"time"`
	Location string `json:"location"`
}

func handler(w http.ResponseWriter, r *http.Request){
	fmt.Println("Request Method: " + r.Method)
	if r.Method == "GET" {
		w.Header().Set("Content-Type", "application/json")

		sensor := Sensor{
			SensorId: "s1",
			Time: "4-23-2026 00:18",
			Location: "Marrakech, Gueliz",
		}

		json.NewEncoder(w).Encode(sensor)
	} else if r.Method == "POST" {		
		var sensor Sensor

		err := json.NewDecoder(r.Body).Decode(&sensor)
		if err != nil {
			http.Error(w, "Invalid JSON", 400)
			return
		}

		fmt.Println("SensorID: " + sensor.SensorId)
		fmt.Println("Time: " + sensor.Time)
		fmt.Println("Location: " + sensor.Location)
	} else {
		http.Error(w, "Method not allowed", 403)
	}
}

func main(){
	http.HandleFunc("/sensors", handler)
	fmt.Printf("Server Listening On: http://localhost:8080\n")
	http.ListenAndServe(":8080", nil)
}

