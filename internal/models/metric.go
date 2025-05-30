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
	Name  string     `json:"name"`
	Type  MetricType `json:"type"`
	Value any        `json:"value"`
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
			Name:  name,
			Type:  GaugeType,
			Value: val,
		}, nil
	case CounterType:
		val, err := strconv.ParseInt(value, 10, 64)
		if err != nil {
			return nil, err
		}
		return &Metric{
			Name:  name,
			Type:  CounterType,
			Value: val,
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
		metric.Value = val
	case CounterType:
		val, err := strconv.ParseInt(value, 10, 64)
		if err != nil {
			return err
		}
		v, ok := metric.Value.(int64)
		if ok {
			metric.Value = v + val
		}
	default:
		return ErrUnknownMetricType
	}
	return nil
}
