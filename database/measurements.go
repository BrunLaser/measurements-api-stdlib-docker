// methods for the measurement table
package database

import (
	"database/sql"
	"fmt"
	"log"
	"math/rand/v2"
)

func (d *Database) InsertMeasurement(m *Measurement) error {
	insertSQL := `INSERT INTO measurements (
		sensors_id,
		value,
		unit) VALUES (?, ?, ?);`
	result, err := d.dbConn.Exec(insertSQL, m.SensorsId, m.Value, m.Unit)
	if err != nil {
		log.Println("Error inserting point: ", err)
		return err
	}
	// Retrieve the last inserted ID
	lastInsertId, err := result.LastInsertId()
	if err != nil {
		log.Println("Error retrieving last insert ID: ", err)
		return err
	}
	m.ID = lastInsertId //change the ID to the asserted one
	return nil
}

func (d *Database) GetAllMeasurements() ([]Measurement, error) {
	rows, err := d.dbConn.Query(`SELECT * FROM measurements;`)
	if err != nil {
		log.Println("Error getting all points: ", err)
		return nil, err
	}
	defer rows.Close()

	var points []Measurement
	for rows.Next() {
		var p Measurement
		if err := rows.Scan(&p.ID, &p.SensorsId, &p.Value, &p.Unit, &p.Timestamp); err != nil {
			log.Println("Error scanning row: ", err)
			return nil, err
		}
		points = append(points, p) //Append points to the Measurement slice
	}

	// Check for errors from the row iteration
	if err := rows.Err(); err != nil {
		log.Println("Error during row iteration: ", err)
		return nil, err
	}

	return points, nil
}

func (d *Database) GetMeasurementById(queryId int) (*Measurement, error) {
	row := d.dbConn.QueryRow(`SELECT * FROM measurements WHERE id = ? LIMIT 1;`, queryId)
	p := &Measurement{}

	if err := row.Scan(&p.ID, &p.SensorsId, &p.Value, &p.Unit, &p.Timestamp); err != nil {
		log.Println("Error getting single row: ", err)
		return nil, err
	}
	return p, nil
}

func (d *Database) DeleteMeasurement(id int) error {
	res, err := d.dbConn.Exec(`DELETE FROM measurements WHERE id = ?;`, id)
	if err != nil {
		return fmt.Errorf("error deleting measurement(id=%v): %w", id, err)
	}
	//Check if there was actually a delete
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("error retrieving rows affected: %w", err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("measurement(id=%v) not found", id)
	}
	return nil
}

func (d *Database) UpdateMeasurement(id int, updateData map[string]any) error {
	//should check here if correct types are passed in json request

	// Build SQL query dynamically
	query := "UPDATE measurements SET "
	args := make([]any, 0, len(updateData)+1) //max cap is one more than updateData (+id)
	i := 0
	for key, value := range updateData {
		if i > 0 {
			//no comma in first iteration
			query += ", "
		}
		query += key + " = ?"
		args = append(args, value)
		i++
	}
	query += " WHERE id = ?"
	fmt.Printf("final update query: %s", query)

	args = append(args, id)
	// Execute the query
	res, err := d.dbConn.Exec(query, args...)
	if err != nil {
		return fmt.Errorf("failed to update: %w", err)
	}

	// Check if the row exists
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("error retrieving rows affected: %w", err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("record not found")
	}

	return nil
}

func (d *Database) BulkInsertRandMeasurementSlow(amount, sensorId int, unit string) error {
	sqlInsert := `INSERT INTO measurements
		(sensors_id,
		value,
		unit)
		VALUES (?, ?, ?);`
	for i := 0; i < amount; i++ {
		_, err := d.dbConn.Exec(sqlInsert, sensorId, rand.Float64()*100, unit)
		if err != nil {
			return fmt.Errorf("error inserting measurement: %w", err)
		}
	}
	return nil
}

func (d *Database) BulkInsertRandMeasurementFast(tx *sql.Tx, amount, sensorId int, unit string) error {
	sqlInsert := `INSERT INTO measurements
	(sensors_id,
	value,
	unit)
	VALUES (?, ?, ?);`
	sqlStmt, err := tx.Prepare(sqlInsert)
	if err != nil {
		return fmt.Errorf("error preparing sql stmt: %w", err)
	}
	defer sqlStmt.Close()

	for i := 0; i < amount; i++ {
		_, err := sqlStmt.Exec(sensorId, rand.Float64()*100, unit)
		if err != nil {
			return fmt.Errorf("error inserting measurement: %w", err)
		}
	}
	return nil
}

func (db *Database) MeasurementRows() (int64, error) {
	sqlQuery := `SELECT COUNT(*) AS total_rows
	FROM measurements;`

	var totalRows sql.NullInt64
	err := db.dbConn.QueryRow(sqlQuery).Scan(&totalRows)
	if err != nil {
		return -1, fmt.Errorf("error executing row count query: %w", err)
	}
	if !totalRows.Valid {
		return 0, fmt.Errorf("no rows found")
	}
	return int64(totalRows.Int64), nil
}

func (db *Database) GetMeasurementMinMax() ([]Measurement, error) {
	//Query could be improved but for now it's good
	sqlQuery := `SELECT * 
	FROM measurements
	WHERE value = (SELECT MAX(value) FROM measurements)
	OR value = (SELECT MIN(value) FROM measurements);`

	rows, err := db.dbConn.Query(sqlQuery)
	if err != nil {
		return nil, fmt.Errorf("error querying min/max: %w", err)
	}
	defer rows.Close()

	var MinMaxMeasurements []Measurement
	for rows.Next() {
		var m Measurement
		if err := rows.Scan(&m.ID, &m.SensorsId, &m.Value, &m.Unit, &m.Timestamp); err != nil {
			return nil, fmt.Errorf("error minmax scanning row: %w", err)
		}
		MinMaxMeasurements = append(MinMaxMeasurements, m)
	}

	// Check for errors from the row iteration
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error during row iteration: %w", err)
	}

	return MinMaxMeasurements, nil
}

// this should fix the problem with updating not supported types, but not finished
/*func (d *Database) measurementToMap(updateData Measurement) map[string]any {
	var updateDataMap map[string]any
	val := reflect.ValueOf(updateData)
	for i := 0; i < val.NumField(); i++ {
		field := val.Field(i)
		if !val.IsZero() {
			// Get the field name
			fieldName := val.Type().Field(i).Name
			// Add to the map
			updateDataMap[fieldName] = field.Interface()
		}
	}
	return updateDataMap
}*/
