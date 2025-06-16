package handlers

import (
	"github.com/DimKa163/go-metrics/internal/models"
	"github.com/DimKa163/go-metrics/internal/services"
	"github.com/gin-gonic/gin"
	"net/http"
)

type metricView struct {
	Name  string
	Value any
}

func HomeHandler(container *services.ServiceContainer) func(c *gin.Context) {
	return func(c *gin.Context) {
		metrics := container.Repository.GetAll()
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
