package controllers

import (
	"github.com/DimKa163/go-metrics/internal/logging"
	"github.com/DimKa163/go-metrics/internal/models"
	"github.com/DimKa163/go-metrics/internal/persistence"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
)

type Metrics interface {
	Home(context *gin.Context)

	UpdateJSON(context *gin.Context)

	Update(context *gin.Context)

	Get(context *gin.Context)

	GetJSON(context *gin.Context)
}

type metrics struct {
	repository persistence.Repository
}

func (m *metrics) GetJSON(context *gin.Context) {
	var model models.Metric
	if err := context.ShouldBindJSON(&model); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	metric, err := m.repository.Find(context.Request.Context(), model.ID)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if metric == nil {
		logging.Log.Info("metric not found", zap.Any("metric", model))
		context.JSON(http.StatusNotFound, "")
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

func NewMetricController(repository persistence.Repository) Metrics {
	return &metrics{
		repository: repository,
	}
}

func (m *metrics) Home(context *gin.Context) {
	met, err := m.repository.GetAll(context)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	viewData := make([]metricView, len(met))
	for _, metric := range met {
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
	context.Writer.Header().Set("Content-Type", "text/html")
	context.HTML(http.StatusOK, "home.tmpl", gin.H{
		"metrics": viewData,
	})
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
	context.Writer.Header().Set("Content-Type", "application/json")
	existingMetric, err := m.repository.Find(context, metric.ID)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if existingMetric == nil {
		existingMetric = &metric
		logging.Log.Info("inserting metric",
			zap.Any("metric", existingMetric))

	} else {
		logging.Log.Info("updating metric", zap.Any("metric", metric))
		existingMetric.Update(&metric)
	}
	err = m.repository.Upsert(context, existingMetric)
	if err != nil {
		context.Writer.WriteHeader(http.StatusInternalServerError)
		return
	}
	context.JSON(http.StatusOK, existingMetric)
}
func (m *metrics) Update(context *gin.Context) {
	t := context.Param("type")
	if t != models.CounterType && t != models.GaugeType {
		context.Writer.WriteHeader(http.StatusBadRequest)
		return
	}
	name := context.Param("name")
	var metric *models.Metric
	context.Writer.Header().Set("Content-Type", "text/plain")
	metric, err := m.repository.Find(context, name)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if metric == nil {
		if metric, err = models.CreateMetric(t, name, context.Param("value")); err != nil {
			context.Writer.WriteHeader(http.StatusBadRequest)
			return
		}
	} else {
		err = models.Update(metric, context.Param("value"))
		if err != nil {
			context.Writer.WriteHeader(http.StatusBadRequest)
			return
		}
	}
	err = m.repository.Upsert(context, metric)
	if err != nil {
		context.Writer.WriteHeader(http.StatusBadRequest)
		return
	}
	context.Writer.WriteHeader(http.StatusOK)
}

func (m *metrics) Get(context *gin.Context) {
	t := context.Param("type")
	name := context.Param("name")
	metric, err := m.repository.Find(context, name)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if metric == nil {
		context.JSON(http.StatusNotFound, "")
		return
	}
	context.Header("Content-Type", "text/plain")
	switch t {
	case models.GaugeType:
		context.JSON(http.StatusOK, metric.Value)
	case models.CounterType:
		context.JSON(http.StatusOK, metric.Delta)
	default:
		context.JSON(http.StatusNotFound, "")
	}
}

type metricView struct {
	Name  string
	Value any
}
