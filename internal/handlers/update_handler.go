package handlers

import (
	"github.com/DimKa163/go-metrics/internal/models"
	"github.com/DimKa163/go-metrics/internal/persistence"
	"github.com/gin-gonic/gin"
	"net/http"
)

func Update(repository persistence.Repository) func(c *gin.Context) {
	return func(c *gin.Context) {
		var metric models.Metric
		if err := c.ShouldBindJSON(&metric); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		if err := models.ValidateMetric(&metric); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.Writer.Header().Set("Content-Type", "application/json")
		existingMetric := repository.Find(metric.Type, metric.ID)
		if existingMetric == nil {
			existingMetric = &metric
			err := repository.Create(existingMetric)
			if err != nil {
				c.Writer.WriteHeader(http.StatusBadRequest)
				return
			}
		} else {
			existingMetric.Value = metric.Value
			existingMetric.Delta = metric.Delta
		}
		c.JSON(http.StatusOK, existingMetric)
	}
}
