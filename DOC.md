# Documenting create a Go server

TODO 0: CRUD
    add informations
    read informations
    delete info
    update info
    
GET url/sensors/id -> sensor with id

POST url/sensor/

DATA form
variables:
metadata
	sensorId -> string
	time -> float
	location -> string
data
	temperature -> float
	polution -> float
	...

NEXT:
	+ define the data form of endpoints & Implement them
	+ structure database Tables & store data with an ORM into just sqlite for testing
	+ refactor the code and add the Timestamp db
	...

data structure
{
	"metadata": 
		{
			"serverID" : "string",
			"time" : "string",
			"location" : "string"
		}
	"data":
		{
			"variable1" : "value",
			"variable2" : "value",
			...
		}
}

			
			
Here’s your **updated API documentation (v1)** including the new **batch ingestion endpoint**.

---

# 📡 Smart City Sensor API (v1)

Base URL:

```text id="base1"
http://localhost:8080/api/v1
```

---

# 🧠 Overview

This API allows IoT sensors to:

* Send single or batch measurements
* Retrieve latest sensor data
* Retrieve historical sensor data

Data is stored in TimescaleDB.

---

# 📥 1. Ingest Single Sensor Data

## ➤ Endpoint

```http id="ep1"
POST /sensors/data
```

---

## ➤ Description

Stores one sensor payload. Each variable in `data` becomes a row in the database.

---

## ➤ Request Body

```json id="req1"
{
  "metadata": {
    "sensor_id": "sensor-1",
    "time": "2026-05-01T12:00:00Z",
    "latitude": 31.63,
    "longitude": -8.0
  },
  "data": {
    "temperature": 25.5,
    "humidity": 60
  }
}
```

---

## ➤ Response

### Success (201)

```json id="res1"
{
  "status": "inserted"
}
```

---

# 📦 2. Batch Sensor Ingestion (NEW)

## ➤ Endpoint

```http id="ep2"
POST /sensors/batch
```

---

## ➤ Description

Allows multiple sensor payloads in a single request. Useful for IoT devices sending buffered data.

---

## ➤ Request Body

```json id="req2"
[
  {
    "metadata": {
      "sensor_id": "sensor-1",
      "time": "2026-05-01T12:00:00Z",
      "latitude": 31.63,
      "longitude": -8.0
    },
    "data": {
      "temperature": 25.5,
      "humidity": 60
    }
  },
  {
    "metadata": {
      "sensor_id": "sensor-2",
      "time": "2026-05-01T12:00:00Z",
      "latitude": 31.70,
      "longitude": -8.05
    },
    "data": {
      "temperature": 26.1
    }
  }
]
```

---

## ➤ Response

### Success (201)

```json id="res2"
{
  "status": "batch inserted"
}
```

---

## ➤ Notes

* Each object is processed independently
* Invalid entries are skipped
* Each `data` key becomes a separate database row

---

# 📊 3. Get Latest Sensor Data

## ➤ Endpoint

```http id="ep3"
GET /sensors/latest?sensor_id={id}
```

---

## ➤ Example

```http id="ex3"
GET /api/v1/sensors/latest?sensor_id=sensor-1
```

---

## ➤ Response

```json id="res3"
[
  {
    "SensorID": "sensor-1",
    "Time": "2026-05-01T12:00:00Z",
    "Latitude": 31.63,
    "Longitude": -8.0,
    "Variable": "temperature",
    "Value": 25.5
  }
]
```

---

# 📈 4. Get Historical Data

## ➤ Endpoint

```http id="ep4"
GET /sensors/history
```

---

## ➤ Query Parameters

| Parameter | Required | Description          |
| --------- | -------- | -------------------- |
| sensor_id | ✅        | Sensor identifier    |
| from      | ✅        | Start time (RFC3339) |
| to        | ✅        | End time (RFC3339)   |

---

## ➤ Example

```http id="ex4"
GET /api/v1/sensors/history?sensor_id=sensor-1&from=2026-05-01T00:00:00Z&to=2026-05-01T12:00:00Z
```

---

## ➤ Response

```json id="res4"
[
  {
    "SensorID": "sensor-1",
    "Time": "2026-05-01T10:00:00Z",
    "Latitude": 31.63,
    "Longitude": -8.0,
    "Variable": "temperature",
    "Value": 24.1
  }
]
```

---

# ❤️ 5. Health Check

## ➤ Endpoint

```http id="ep5"
GET /health
```

---

## ➤ Response

```text id="res5"
OK
```

---

# 🧱 6. Data Model

## Request Model

```go id="m1"
type SensorRequest struct {
    Metadata Metadata
    Data     map[string]float64
}
```

---

## Stored Model

```go id="m2"
type SensorEntity struct {
    SensorID  string
    Time      time.Time
    Latitude  float64
    Longitude float64
    Variable  string
    Value     float64
}
```

---

# ⚙️ 7. System Architecture

```text id="arch1"
Sensor Device
     ↓
HTTP API (Go)
     ↓
Validation Layer (Handler)
     ↓
Repository (SQL)
     ↓
TimescaleDB
```

---

# 🚀 8. Future Improvements

* streaming ingestion via RabbitMQ
* real-time dashboards
* aggregation endpoints (1min / 5min / hourly)
* anomaly detection
* WebSocket live updates

---

If you want next step, I can help you generate:

* OpenAPI/Swagger spec (auto docs)
* Postman collection
* or production-grade IoT architecture (queues + workers + scaling)
