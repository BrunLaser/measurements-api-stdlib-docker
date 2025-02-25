package database

import (
	"fmt"
	"time"
)

func (d *Database) experimentExists(name string) (bool, error) {
	queryDB := `SELECT EXISTS (
		SELECT 1
		FROM experiments
		WHERE name = ?
	) AS order_exists;`

	exists := false
	if err := d.dbConn.QueryRow(queryDB, name).Scan(&exists); err != nil {
		return false, fmt.Errorf("error checking if %s exists: %w", name, err)
	}
	return exists, nil
}

func (d *Database) constructTimeRangeSQL(name, startTime, endTime string) (string, []any, error) {
	//Basic query
	queryDB := `SELECT value, unit, timestamp 
	FROM measurements 
	INNER JOIN sensors 		ON measurements.sensors_id 	= sensors.id
	INNER JOIN experiments 	ON sensors.experiment_id 	= experiments.id
	WHERE experiments.name = ?`

	queryParams := make([]any, 0, 2)
	queryParams = append(queryParams, name)

	if startTime != "" {
		parsedStartTime, err := time.Parse(time.DateTime, startTime)
		if err != nil {
			return "", nil, fmt.Errorf("invalid format of startingtime(%w)", err)
		}

		startTime = parsedStartTime.Format(time.DateTime)
		queryParams = append(queryParams, startTime)
		queryDB += " AND measurements.timestamp >= ?"
	}

	if endTime != "" {
		parsedEndTime, err := time.Parse(time.DateTime, endTime)
		if err != nil {
			return "", nil, fmt.Errorf("invalid format of endingtime(%w)", err)
		}

		endTime = parsedEndTime.Format(time.DateTime)
		queryParams = append(queryParams, endTime)
		queryDB += " AND measurements.timestamp <= ?"
	}
	queryDB += ";"
	//fmt.Println(queryDB)
	return queryDB, queryParams, nil
}

func (d *Database) GetMeasurementsByExperiment(expName, startTime, endTime string) (*[]MeasurementResponse, error) {
	//check for an experiment with the submitted name
	exists, err := d.experimentExists(expName)
	if err != nil {
		return nil, fmt.Errorf("failed to check experiment %s exists: %w", expName, err)
	} else if !exists {
		return nil, fmt.Errorf("experiment %s does not exist", expName)
	}

	//build the query accordingt to submitted params
	queryDB, params, err := d.constructTimeRangeSQL(expName, startTime, endTime)
	if err != nil {
		return nil, fmt.Errorf("error constructing sql query: %w", err)
	}

	rows, err := d.dbConn.Query(queryDB, params...)
	if err != nil {
		return nil, fmt.Errorf("error querying measurements: %w", err)
	}
	defer rows.Close() // Ensure rows are closed after processing

	var measurements []MeasurementResponse
	for rows.Next() {
		var m MeasurementResponse
		if err := rows.Scan(&m.Value, &m.Unit, &m.Timestamp); err != nil {
			return nil, fmt.Errorf("error scanning row: %w", err)
		}
		measurements = append(measurements, m)
	}

	// Check for errors from iterating over rows
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating over rows: %w", err)
	}

	return &measurements, nil
}
