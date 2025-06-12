package handlers

import (
	"github.com/DimKa163/go-metrics/internal/models"
	"github.com/DimKa163/go-metrics/internal/persistence"
	"github.com/gin-gonic/gin"
	"net/http"
)

func GetHandler(repository persistence.Repository) func(c *gin.Context) {
	return func(c *gin.Context) {
		var model models.Metric
		if err := c.ShouldBindJSON(&model); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		metric := repository.Find(model.Type, model.ID)
		if metric == nil {
			c.JSON(http.StatusNotFound, "")
			return
		}
		c.Header("Content-Type", "application/json")
		c.JSON(http.StatusOK, metric)
	}
}
