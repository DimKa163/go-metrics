package models

import (
	"errors"
)

const (
	GaugeType   = "gauge"
	CounterType = "counter"
)

var ErrUnknownMetricType = errors.New("unknown metric type")

type Gauge float64

type Counter int64

type MetricType string

type Metric struct {
	ID    string     `json:"id"`
	Type  MetricType `json:"type"`
	Delta *int64     `json:"delta,omitempty"`
	Value *float64   `json:"value,omitempty"`
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
