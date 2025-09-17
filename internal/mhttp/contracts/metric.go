package contracts

type (
	MetricView struct {
		Name  string
		Value any
	}
	// Metric info
	Metric struct {
		ID    string   `json:"id"`
		Type  string   `json:"type"`
		Value *float64 `json:"value,omitempty"`
		Delta *int64   `json:"delta,omitempty"`
	}
)
