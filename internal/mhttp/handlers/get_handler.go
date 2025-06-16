package handlers

import (
	"github.com/DimKa163/go-metrics/internal/logging"
	"github.com/DimKa163/go-metrics/internal/models"
	"github.com/DimKa163/go-metrics/internal/services"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
)

func GetHandler(container *services.ServiceContainer) func(c *gin.Context) {
	return func(c *gin.Context) {
		t := c.Param("type")
		name := c.Param("name")
		metric := container.Repository.Find(t, name)
		if metric == nil {
			c.JSON(http.StatusNotFound, "")
			return
		}
		c.Header("Content-Type", "text/plain")
		switch t {
		case models.GaugeType:
			c.JSON(http.StatusOK, metric.Value)
		case models.CounterType:
			c.JSON(http.StatusOK, metric.Delta)
		default:
			c.JSON(http.StatusNotFound, "")
		}
	}
}

func GetHandlerJSON(container *services.ServiceContainer) func(c *gin.Context) {
	return func(c *gin.Context) {
		var model models.Metric
		if err := c.ShouldBindJSON(&model); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		metric := container.Repository.Find(model.Type, model.ID)
		if metric == nil {
			logging.Log.Info("metric not found", zap.Any("metric", model))
			c.JSON(http.StatusNotFound, "")
			return
		}
		c.Header("Content-Type", "application/json")
		switch metric.Type {
		case models.GaugeType:
			c.JSON(http.StatusOK, metric)
		case models.CounterType:
			c.JSON(http.StatusOK, metric)
		default:
			c.JSON(http.StatusNotFound, "")
		}
	}
}
