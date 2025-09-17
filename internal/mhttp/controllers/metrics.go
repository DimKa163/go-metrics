package controllers

import (
	"errors"
	"github.com/DimKa163/go-metrics/internal/mhttp/contracts"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"github.com/DimKa163/go-metrics/internal/logging"
	"github.com/DimKa163/go-metrics/internal/models"
	"github.com/DimKa163/go-metrics/internal/usecase"
)

// @Title MetricStorage API
// @Description Metric service.
// @Version 1.0

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

// Map map all routs
func (m *metrics) Map(engine *gin.Engine) {
	engine.GET("/", m.Home)
	engine.GET("/value/:type/:name", m.Get)
	engine.POST("/value/", m.GetJSON)
	engine.POST("/update/:type/:name/:value", m.Update)
	engine.POST("/update/", m.UpdateJSON)
	engine.POST("/updates", m.UpdatesJSON)
}

// GetJSON get metric
// @Produce application/json
// @Param metric body contracts.Metric true "metric"
// @Success 200 {object} contracts.Metric "success request"
// @Failure 400 {object} contracts.ErrorModel "bad request"
// @Failure 500 {object} contracts.ErrorModel "internal server error"
// @Router /value [post]
func (m *metrics) GetJSON(context *gin.Context) {
	var model contracts.Metric
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
		context.JSON(http.StatusInternalServerError, contracts.ErrorModel{Error: err.Error()})
		return
	}
	context.Header("Content-Type", "application/json")
	switch metric.Type {
	case models.GaugeType:
		context.JSON(http.StatusOK, contracts.Metric{ID: metric.ID, Value: metric.Value})
	case models.CounterType:
		context.JSON(http.StatusOK, contracts.Metric{ID: metric.ID, Delta: metric.Delta})
	default:
		context.JSON(http.StatusNotFound, "")
	}
}

// Home all metrics
// @Produce text/html
// @Success 200 {string} []contracts.MetricView "success request"
// @Failure 400 {object} contracts.ErrorModel "bad request"
// @Failure 500 {object} contracts.ErrorModel "internal server error"
// @Router / [get]
func (m *metrics) Home(context *gin.Context) {
	met, err := m.service.GetAll(context)
	if err != nil {
		context.JSON(http.StatusInternalServerError, contracts.ErrorModel{Error: err.Error()})
		return
	}
	viewData := make([]contracts.MetricView, len(met))
	for i, metric := range met {
		switch metric.Type {
		case models.GaugeType:
			viewData[i] = contracts.MetricView{
				Name:  metric.ID,
				Value: metric.Value,
			}
		case models.CounterType:
			viewData[i] = contracts.MetricView{
				Name:  metric.ID,
				Value: metric.Delta,
			}
		}
	}
	context.Writer.Header().Set("Content-Type", "text/html")
	context.HTML(http.StatusOK, "home.tmpl", gin.H{
		"metrics": viewData,
	})
}

// UpdatesJSON update many metric
// @Produce application/json
// @Param metrics body []contracts.Metric true "metric array"
// @Success 200 {string} string "success request"
// @Failure 400 {object} contracts.ErrorModel "bad request"
// @Failure 500 {object} contracts.ErrorModel "internal server error"
// @Router /updates [post]
func (m *metrics) UpdatesJSON(context *gin.Context) {
	var metricList []contracts.Metric
	if err := context.ShouldBindJSON(&metricList); err != nil {
		context.JSON(http.StatusBadRequest, contracts.ErrorModel{Error: err.Error()})
		return
	}
	data := make([]models.Metric, len(metricList))
	for i, metric := range metricList {
		metricIt := models.Metric{
			ID:    metric.ID,
			Type:  metric.Type,
			Value: metric.Value,
			Delta: metric.Delta,
		}
		if err := models.ValidateMetric(&metricIt); err != nil {
			context.JSON(http.StatusBadRequest, contracts.ErrorModel{Error: err.Error()})
			return
		}
		data[i] = metricIt
	}
	if err := m.service.BatchUpdate(context, data); err != nil {
		context.JSON(http.StatusInternalServerError, contracts.ErrorModel{Error: err.Error()})
		return
	}
	context.Writer.Header().Set("Content-Type", "application/json")
	context.Status(http.StatusOK)
}

// UpdateJSON update metric
// @Produce application/json
// @Param metric body contracts.Metric true "metric"
// @Failure 400 {object} contracts.ErrorModel "bad request"
// @Failure 500 {object} contracts.ErrorModel "internal server error"
// @Router /update [post]
func (m *metrics) UpdateJSON(context *gin.Context) {
	var contract contracts.Metric
	if err := context.ShouldBindJSON(&contract); err != nil {
		context.JSON(http.StatusBadRequest, contracts.ErrorModel{Error: err.Error()})
		return
	}
	metric := models.Metric{
		ID:    contract.ID,
		Type:  contract.Type,
		Value: contract.Value,
		Delta: contract.Delta,
	}
	if err := models.ValidateMetric(&metric); err != nil {
		context.JSON(http.StatusBadRequest, contracts.ErrorModel{Error: err.Error()})
		return
	}
	result, err := m.service.Upsert(context, metric)
	if err != nil {
		context.JSON(http.StatusInternalServerError, contracts.ErrorModel{Error: err.Error()})
		return
	}

	context.Writer.Header().Set("Content-Type", "application/json")
	context.JSON(http.StatusOK, result)
}

// Update update metric
// @Produce plain/text
// @Produce json
// @Failure 400 {object} contracts.ErrorModel "bad request"
// @Failure 500 {object} contracts.ErrorModel "internal server error"
// @Param type path string true "Metric type"
// @Param name path string true "Metric name"
// @Param value path string true "Metric value"
// @Router /update/{type}/{name}/{value} [post]
func (m *metrics) Update(context *gin.Context) {
	t := context.Param("type")
	name := context.Param("name")
	value := context.Param("value")
	metric, err := models.CreateMetric(t, name, value)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	_, err = m.service.Upsert(context, metric)
	if err != nil {
		context.JSON(http.StatusInternalServerError, contracts.ErrorModel{Error: err.Error()})
		return
	}
	context.Writer.Header().Set("Content-Type", "text/plain")
	context.Writer.WriteHeader(http.StatusOK)
}

// Get metric
// @Produce plain/text
// @Produce json
// @Param type path string true "Metric type"
// @Param name path string true "Metric name"
// @Router /value/{type}/{name} [get]
func (m *metrics) Get(context *gin.Context) {
	t := context.Param("type")
	name := context.Param("name")
	metric, err := m.service.Get(context, name)
	if err != nil {
		if errors.Is(err, usecase.ErrMetricNotFound) {
			context.JSON(http.StatusNotFound, "")
			return
		}
		context.JSON(http.StatusInternalServerError, contracts.ErrorModel{Error: err.Error()})
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
