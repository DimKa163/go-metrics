package models

import (
	"github.com/stretchr/testify/assert"
	"strconv"
	"testing"
)

func TestCreateGauge(t *testing.T) {
	val := 42.5
	m := CreateGauge("gauge_metric", val)

	assert.NotNil(t, m)
	assert.Equal(t, "gauge_metric", m.ID)
	assert.Equal(t, val, *m.Value)
	assert.Nil(t, m.Delta)
}

func TestCreateCounter(t *testing.T) {
	val := int64(42)
	m := CreateCounter("counter_metric", val)
	assert.NotNil(t, m)
	assert.Equal(t, "counter_metric", m.ID)
	assert.Equal(t, val, *m.Delta)
	assert.Nil(t, m.Value)
}

func TestUpdateGauge(t *testing.T) {
	m := CreateGauge("gauge_metric", 42.5)

	m1 := CreateGauge("gauge_metric", 55.6)

	m.Update(*m1)

	assert.Equal(t, "gauge_metric", m.ID)
	assert.Equal(t, 55.6, *m.Value)
	assert.Nil(t, m.Delta)
}

func TestUpdateCounter(t *testing.T) {
	m := CreateCounter("counter_metric", 42)

	m1 := CreateCounter("counter_metric", 10)
	m.Update(*m1)
	assert.Equal(t, "counter_metric", m.ID)
	assert.Equal(t, int64(52), *m.Delta)
	assert.Nil(t, m.Value)
}

func TestSuccessCreateCounterMetric(t *testing.T) {
	val := int64(500)
	m, err := CreateMetric(CounterType, "Counter", strconv.FormatInt(val, 10))

	assert.NoError(t, err)
	assert.NotEqual(t, Metric{}, m)
	assert.Equal(t, "Counter", m.ID)
	assert.Equal(t, val, *m.Delta)

}

func TestFailureCreateMetric(t *testing.T) {
	val := int64(500)
	m, err := CreateMetric("ops", "Counter", strconv.FormatInt(val, 10))

	assert.Error(t, err)
	assert.Equal(t, Metric{}, m)
	assert.ErrorIs(t, err, ErrUnknownMetricType)
}

func TestSuccessCreateGaugeMetric(t *testing.T) {
	val := float64(499.99)

	m, err := CreateMetric(GaugeType, "Gauge", strconv.FormatFloat(val, 'f', -1, 64))
	assert.NoError(t, err)
	assert.NotEqual(t, Metric{}, m)
	assert.Equal(t, "Gauge", m.ID)
	assert.Equal(t, val, *m.Value)
}

func TestValidateMetric(t *testing.T) {
	g := CreateGauge("gauge", 3.14)
	if err := ValidateMetric(g); err != nil {
		t.Errorf("expected valid gauge, got error: %v", err)
	}

	c := CreateCounter("counter", 123)
	if err := ValidateMetric(c); err != nil {
		t.Errorf("expected valid counter, got error: %v", err)
	}

	invalid := &Metric{ID: "bad", Type: "unknown"}
	if err := ValidateMetric(invalid); err == nil {
		t.Errorf("expected error for invalid metric")
	}
}
