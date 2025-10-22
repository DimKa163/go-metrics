package contracts

type (
	MetricView struct {
		Name  string
		Value any
	}
	// Metric info
	Metric struct {
		ID    string  `json:"id"`
		Type  string  `json:"type"`
		Value float64 `json:"value,omitempty"`
		Delta int64   `json:"delta,omitempty"`
	}
)

func CreateCounter(id string, delta int64) *Metric {
	return &Metric{
		ID:    id,
		Type:  "counter",
		Delta: delta,
	}
}

func CreateGauge(id string, value float64) *Metric {
	return &Metric{
		ID:    id,
		Type:  "gauge",
		Value: value,
	}
}
