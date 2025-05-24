package handlers

import (
	"github.com/DimKa163/go-metrics/internal/models"
	"github.com/DimKa163/go-metrics/internal/persistence"
	"github.com/gin-gonic/gin"
	"net/http"
)

func GetHandler(repository persistence.Repository) func(c *gin.Context) {
	return func(c *gin.Context) {
		t := c.Param("type")
		name := c.Param("name")
		var metric *models.Metric
		metric = repository.Find(models.MetricType(t), name)
		if metric == nil {
			c.JSON(http.StatusNotFound, "")
			return
		}
		c.Header("Content-Type", "text/plain")
		c.JSON(http.StatusOK, metric.Value)
	}
}
