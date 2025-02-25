package handlers

import (
	"fmt"
	"log"
	"measurements-api-stdlib-docker/database"
	"measurements-api-stdlib-docker/util"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	db *database.Database
}

func NewHandler(db *database.Database) *Handler {
	return &Handler{db: db}
}

func (h *Handler) HandleMeasurementPost(c *gin.Context) {
	newPoint := &database.Measurement{}
	//json to struct
	if err := c.ShouldBindJSON(newPoint); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON Data"})
		return
	}
	//struct to database
	if err := h.db.InsertMeasurement(newPoint); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database INSERT error"})
		return
	}
	fmt.Println(*newPoint)
	location := fmt.Sprintf("/measurements/%d", newPoint.ID) //The ID is set to the actual ID in the database
	c.Header("Location", location)
	c.JSON(http.StatusCreated, gin.H{
		"message":  "Point created",
		"location": location,
		"data":     *newPoint,
	})
}

func (h *Handler) HandleMeasurementGetAll(c *gin.Context) {
	//w.Header().Set("Content-Type", "application/json") //gin does the header when I do json stuff
	measurements, err := h.db.GetAllMeasurements()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, measurements) //write json data back
}

func (h *Handler) HandleMeasurementGetById(c *gin.Context) {
	id, err := util.GetParamInt(c, "id")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	point, err := h.db.GetMeasurementById(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, point)
}

func (h *Handler) HandleMeasurementDelete(c *gin.Context) {
	id, err := util.GetParamInt(c, "id")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	err = h.db.DeleteMeasurement(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": fmt.Sprintf("succesfully deleted measurement %v", id)})
}

func (h *Handler) HandleMeasurementUpdate(c *gin.Context) {
	id, err := util.GetParamInt(c, "id")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var updateData map[string]any
	if err := c.BindJSON(&updateData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.db.UpdateMeasurement(id, updateData); err != nil {
		if err.Error() == "record not found" {
			// Return 404 if record not found
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		} else {
			//Internal Server Error
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{"message": "Successfully updated measurement"})
}

func (h *Handler) HandleGetMeasurementsByExperiment(c *gin.Context) {
	expName := c.Param("exp")
	if expName == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "empty parameter"})
		return
	}
	startTime := c.Query("startTime")
	endTime := c.Query("endTime")
	fmt.Println(startTime)
	fmt.Println(endTime)
	measurements, err := h.db.GetMeasurementsByExperiment(expName, startTime, endTime)
	if err != nil {
		//we could check for different errors here
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, measurements)
}

func (h *Handler) HandleMeasurementMinMax(c *gin.Context) {
	start := time.Now()
	Measurements, err := h.db.GetMeasurementMinMax()
	log.Printf("Time to get MinMax: %s", time.Since(start))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	log.Printf("Number of MinMax: %v", len(Measurements))
	c.JSON(http.StatusOK, Measurements)
}
