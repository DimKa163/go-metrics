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

	UpdatesJSON(context *gin.Context)

	Update(context *gin.Context)

	Get(context *gin.Context)

	GetJSON(context *gin.Context)
}

type metrics struct {
	repository persistence.Repository
}

func NewMetricController(repository persistence.Repository) Metrics {
	return &metrics{
		repository: repository,
	}
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

func (m *metrics) UpdatesJSON(context *gin.Context) {
	var metricList []models.Metric
	var err error
	if err = context.ShouldBindJSON(&metricList); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	mapMetric := make(map[string]models.Metric)
	for _, metric := range metricList {
		if err = models.ValidateMetric(&metric); err != nil {
			context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		it, ok := mapMetric[metric.ID]
		if ok {
			switch it.Type {
			case models.GaugeType:
				mapMetric[metric.ID] = metric
			case models.CounterType:
				*it.Delta = *metric.Delta + *it.Delta
				mapMetric[metric.ID] = it
			}
			continue
		}
		mapMetric[metric.ID] = metric
	}
	resultList := make([]models.Metric, 0)
	for _, metric := range mapMetric {
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
		resultList = append(resultList, *existingMetric)
	}
	if err = m.repository.BatchUpsert(context, resultList); err != nil {
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
	context.Writer.Header().Set("Content-Type", "application/json")
	met, err := m.processMetric(context, metric)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	err = m.repository.Upsert(context, met)
	if err != nil {
		context.Writer.WriteHeader(http.StatusInternalServerError)
		return
	}
	context.JSON(http.StatusOK, met)
}

func (m *metrics) processMetric(context *gin.Context, metric models.Metric) (*models.Metric, error) {
	existingMetric, err := m.repository.Find(context, metric.ID)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return nil, err
	}
	if existingMetric == nil {
		existingMetric = &metric
		logging.Log.Info("inserting metric",
			zap.Any("metric", existingMetric))

	} else {
		logging.Log.Info("updating metric", zap.Any("metric", metric))
		existingMetric.Update(&metric)
	}
	return existingMetric, nil
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
	context.Writer.Header().Set("Content-Type", "text/plain")
	metric, err = m.processMetric(context, *metric)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
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
