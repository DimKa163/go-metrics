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
func CreateMetric(tt string, name string, value string) (*Metric, error) {
	if tt != GaugeType && tt != CounterType {
		return nil, ErrUnknownMetricType
	}
	switch tt {
	case GaugeType:
		val, err := strconv.ParseFloat(value, 64)
		if err != nil {
			return nil, err
		}
		return &Metric{
			ID:    name,
			Type:  GaugeType,
			Value: &val,
		}, nil
	case CounterType:
		val, err := strconv.ParseInt(value, 10, 64)
		if err != nil {
			return nil, err
		}
		return &Metric{
			ID:    name,
			Type:  CounterType,
			Delta: &val,
		}, nil
	default:
		return nil, ErrUnknownMetricType
	}
}
func Update(metric *Metric, value string) error {
	switch metric.Type {
	case GaugeType:
		val, err := strconv.ParseFloat(value, 64)
		if err != nil {
			return err
		}
		metric.Value = &val
	case CounterType:
		val, err := strconv.ParseInt(value, 10, 64)
		if err != nil {
			return err
		}
		newVal := *metric.Delta + val
		metric.Delta = &newVal
	default:
		return ErrUnknownMetricType
	}
	return nil
}
