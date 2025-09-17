package models

import (
	"errors"
	"strconv"
)

const (
	GaugeType   = "gauge"
	CounterType = "counter"
)

var ErrUnknownMetricType = errors.New("unknown metric type")

type Gauge float64

type Counter int64

type Metric struct {
	ID    string   `json:"id"`
	Type  string   `json:"type"`
	Delta *int64   `json:"delta,omitempty"`
	Value *float64 `json:"value,omitempty"`
}

func (m *Metric) Update(metric Metric) {
	switch metric.Type {
	case GaugeType:
		m.Value = metric.Value
	case CounterType:
		sum := *m.Delta + *metric.Delta
		m.Delta = &sum
	}
}

func ValidateMetric(model *Metric) error {
	if model.Type != GaugeType && model.Type != CounterType {
		return ErrUnknownMetricType
	}
	switch model.Type {
	case GaugeType:
		return nil
	case CounterType:
		return nil
	default:
		return ErrUnknownMetricType
	}
}

func CreateCounter(id string, delta int64) *Metric {
	return &Metric{
		ID:    id,
		Type:  CounterType,
		Delta: &delta,
	}
}

func CreateGauge(id string, value float64) *Metric {
	return &Metric{
		ID:    id,
		Type:  GaugeType,
		Value: &value,
	}
}
func CreateMetric(tt string, name string, value string) (Metric, error) {
	switch tt {
	case GaugeType:
		val, err := strconv.ParseFloat(value, 64)
		if err != nil {
			return Metric{}, err
		}
		return Metric{
			ID:    name,
			Type:  GaugeType,
			Value: &val,
		}, nil
	case CounterType:
		val, err := strconv.ParseInt(value, 10, 64)
		if err != nil {
			return Metric{}, err
		}
		return Metric{
			ID:    name,
			Type:  CounterType,
			Delta: &val,
		}, nil
	default:
		return Metric{}, ErrUnknownMetricType
	}
}
