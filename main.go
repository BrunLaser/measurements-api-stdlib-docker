package main

import (
	"log"
	"measurements-api-stdlib-docker/database"
	"measurements-api-stdlib-docker/handlers"
	"measurements-api-stdlib-docker/router"
	"time"

	"github.com/gin-gonic/gin"
)

func main() {
	measurementDB, err := database.InitDB()
	if err != nil {
		log.Fatal("Database connection/creation failed:", err)
	}
	defer measurementDB.Close()

	//count all rows
	start := time.Now()
	nRows, err := measurementDB.MeasurementRows()
	if err != nil {
		log.Printf("error counting rows: %s", err.Error())
		return
	}
	log.Printf("measurement rows: %v; time %s", nRows, time.Since(start))

	//Setup API
	measurementHandler := handlers.NewHandler(measurementDB)
	r := gin.Default()
	router.SetupRoutes(r, measurementHandler)

	err = r.Run(":8080")
	if err != nil {
		measurementDB.Close()             // Ensure proper cleanup
		log.Fatal("Server failed: ", err) // Exit the program
	}
}
