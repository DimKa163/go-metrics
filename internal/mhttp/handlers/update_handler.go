package handlers

import (
	"github.com/DimKa163/go-metrics/internal/logging"
	"github.com/DimKa163/go-metrics/internal/models"
	"github.com/DimKa163/go-metrics/internal/services"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
)

func UpdateJSON(container *services.ServiceContainer) func(c *gin.Context) {
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
		existingMetric := container.Repository.Find(metric.Type, metric.ID)
		if existingMetric == nil {
			existingMetric = &metric
			logging.Log.Info("inserting metric",
				zap.Any("metric", metric))
			err := container.Repository.Create(existingMetric)
			if err != nil {
				c.Writer.WriteHeader(http.StatusBadRequest)
				return
			}
		} else {
			logging.Log.Info("updating metric", zap.Any("metric", metric))
			existingMetric.Add(&metric)
		}
		c.JSON(http.StatusOK, existingMetric)
	}
}

func Update(container *services.ServiceContainer) func(c *gin.Context) {
	return func(c *gin.Context) {
		t := c.Param("type")
		name := c.Param("name")
		var metric *models.Metric
		c.Writer.Header().Set("Content-Type", "text/plain")
		metric = container.Repository.Find(t, name)
		if metric == nil {
			var err error

			if metric, err = models.CreateMetric(t, name, c.Param("value")); err != nil {
				c.Writer.WriteHeader(http.StatusBadRequest)
				return
			}
			err = container.Repository.Create(metric)
			if err != nil {
				c.Writer.WriteHeader(http.StatusBadRequest)
				return
			}
		} else {
			err := models.Update(metric, c.Param("value"))
			if err != nil {
				c.Writer.WriteHeader(http.StatusBadRequest)
				return
			}
		}

		c.Writer.WriteHeader(http.StatusOK)
	}
}
