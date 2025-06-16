package handlers

import (
	"github.com/DimKa163/go-metrics/internal/models"
	"github.com/DimKa163/go-metrics/internal/persistence"
	"github.com/gin-gonic/gin"
	"net/http"
)

func Update(repository persistence.Repository) func(c *gin.Context) {
	return func(c *gin.Context) {
		t := c.Param("type")
		name := c.Param("name")
		var metric *models.Metric
		c.Writer.Header().Set("Content-Type", "text/plain")
		metric = repository.Find(models.MetricType(t), name)
		if metric == nil {
			var err error
			if metric, err = models.CreateMetric(t, name, c.Param("value")); err != nil {
				c.Writer.WriteHeader(http.StatusBadRequest)
				return
			}
			err = repository.Create(metric)
			if err != nil {
				c.Writer.WriteHeader(http.StatusBadRequest)
				return
			}
			c.Writer.WriteHeader(http.StatusOK)
			return
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
