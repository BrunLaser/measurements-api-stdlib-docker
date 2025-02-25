package database

import "database/sql"

type Database struct {
	dbConn *sql.DB
}

type Experiment struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Date        string `json:"date"`
}

type Sensor struct {
	ID           int    `json:"id"`
	ExperimentID int    `json:"experiment_id"`
	SensorType   string `json:"sensor_type"`
}

type Measurement struct {
	ID        int64   `json:"id"`
	SensorsId int64   `json:"sensor_id"`
	Value     float64 `json:"value"`
	Unit      string  `json:"unit"`
	Timestamp string  `json:"timestamp"`
}

type MeasurementResponse struct {
	Value     float64 `json:"value"`
	Unit      string  `json:"unit"`
	Timestamp string  `json:"timestamp"`
}
