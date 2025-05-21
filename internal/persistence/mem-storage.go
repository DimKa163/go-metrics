package persistence

var Store MemStorage

func init() {
	Store = MemStorage{
		Gauge:   make(map[string]float64),
		Counter: make(map[string]int64),
	}
}

type MemStorage struct {
	Gauge   map[string]float64
	Counter map[string]int64
}
