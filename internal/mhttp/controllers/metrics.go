package controllers

import (
	"errors"
	"github.com/DimKa163/go-metrics/internal/logging"
	"github.com/DimKa163/go-metrics/internal/models"
	"github.com/DimKa163/go-metrics/internal/usecase"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
)

type Metrics interface {
	Map(engine *gin.Engine)
	Home(context *gin.Context)

	UpdateJSON(context *gin.Context)

	UpdatesJSON(context *gin.Context)

	Update(context *gin.Context)

	Get(context *gin.Context)

	GetJSON(context *gin.Context)
}

type metrics struct {
	service *usecase.MetricService
}

func NewMetricController(service *usecase.MetricService) Metrics {
	return &metrics{
		service: service,
	}
}

func (m *metrics) Map(engine *gin.Engine) {
	engine.GET("/", m.Home)
	engine.GET("/value/:type/:name", m.Get)
	engine.POST("/value/", m.GetJSON)
	engine.POST("/update/:type/:name/:value", m.Update)
	engine.POST("/update/", m.UpdateJSON)
	engine.POST("/updates", m.UpdatesJSON)
}
func (m *metrics) GetJSON(context *gin.Context) {
	var model models.Metric
	if err := context.ShouldBindJSON(&model); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	metric, err := m.service.Get(context, model.ID)
	if err != nil {
		if errors.Is(err, usecase.ErrMetricNotFound) {
			logging.Log.Info("metric not found", zap.Any("metric", model))
			context.JSON(http.StatusNotFound, "")
			return
		}
		context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	context.Header("Content-Type", "application/json")
	switch metric.Type {
	case models.GaugeType:
		context.JSON(http.StatusOK, metric)
	case models.CounterType:
		context.JSON(http.StatusOK, metric)
	default:
		context.JSON(http.StatusNotFound, "")
	}
}

func (m *metrics) Home(context *gin.Context) {
	met, err := m.service.GetAll(context)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	viewData := make([]metricView, len(met))
	for i, metric := range met {
		switch metric.Type {
		case models.GaugeType:
			viewData[i] = metricView{
				Name:  metric.ID,
				Value: metric.Value,
			}
		case models.CounterType:
			viewData[i] = metricView{
				Name:  metric.ID,
				Value: metric.Value,
			}
		}
	}
	context.Writer.Header().Set("Content-Type", "text/html")
	context.HTML(http.StatusOK, "home.tmpl", gin.H{
		"metrics": viewData,
	})
}

func (m *metrics) UpdatesJSON(context *gin.Context) {
	var metricList []models.Metric
	if err := context.ShouldBindJSON(&metricList); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	for _, metric := range metricList {
		if err := models.ValidateMetric(&metric); err != nil {
			context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
	}
	if err := m.service.BatchUpdate(context, metricList); err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	context.Writer.Header().Set("Content-Type", "application/json")
	context.Status(http.StatusOK)
}

func (m *metrics) UpdateJSON(context *gin.Context) {
	var metric models.Metric
	if err := context.ShouldBindJSON(&metric); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := models.ValidateMetric(&metric); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	result, err := m.service.Update(context, metric)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	context.Writer.Header().Set("Content-Type", "application/json")
	context.JSON(http.StatusOK, result)
}
func (m *metrics) Update(context *gin.Context) {
	t := context.Param("type")
	name := context.Param("name")
	value := context.Param("value")
	metric, err := models.CreateMetric(t, name, value)
	if err != nil {
		context.Writer.WriteHeader(http.StatusBadRequest)
		return
	}
	_, err = m.service.Update(context, metric)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	context.Writer.Header().Set("Content-Type", "text/plain")
	context.Writer.WriteHeader(http.StatusOK)
}

func (m *metrics) Get(context *gin.Context) {
	t := context.Param("type")
	name := context.Param("name")
	metric, err := m.service.Get(context, name)
	if err != nil {
		if errors.Is(err, usecase.ErrMetricNotFound) {
			context.JSON(http.StatusNotFound, "")
			return
		}
		context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	context.Header("Content-Type", "text/plain")
	switch t {
	case models.GaugeType:
		context.JSON(http.StatusOK, metric.Value)
	case models.CounterType:
		context.JSON(http.StatusOK, metric.Delta)
	}
}

type metricView struct {
	Name  string
	Value any
}
