package persistence

type MemStorage struct {
	Gauge   map[string]float64
	Counter map[string]int64
}
