package handlers

import (
	"github.com/DimKa163/go-metrics/internal/models"
	"github.com/DimKa163/go-metrics/internal/persistence"
	"github.com/gin-gonic/gin"
	"net/http"
)

type metricView struct {
	Name  string
	Value any
}

func HomeHandler(repository persistence.Repository) func(c *gin.Context) {
	return func(c *gin.Context) {
		metrics := repository.GetAll()
		viewData := make([]metricView, len(metrics))
		for _, metric := range metrics {
			switch metric.Type {
			case models.GaugeType:
				viewData = append(viewData, metricView{
					Name:  metric.ID,
					Value: metric.Value,
				})
			case models.CounterType:
				viewData = append(viewData, metricView{
					Name:  metric.ID,
					Value: metric.Delta,
				})
			}
		}
		c.Writer.Header().Set("Content-Type", "text/html")
		c.HTML(http.StatusOK, "home.tmpl", gin.H{
			"metrics": viewData,
		})
	}
}
