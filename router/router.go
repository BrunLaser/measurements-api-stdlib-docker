package router

import (
	"measurements-api-stdlib-docker/handlers"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(r *gin.Engine, h *handlers.Handler) {
	r.GET("/measurements", h.HandleMeasurementGetAll)
	r.POST("/measurements", h.HandleMeasurementPost)
	r.GET("/measurements/:id", h.HandleMeasurementGetById)
	r.DELETE("/measurements/:id", h.HandleMeasurementDelete)
	r.PUT("/measurements/:id", h.HandleMeasurementUpdate)
	r.GET("/measurements/minmax", h.HandleMeasurementMinMax)

	r.GET("experiments/:exp/measurements", h.HandleGetMeasurementsByExperiment)
}
