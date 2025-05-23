package persistence

type GaugeRepository interface {
	Update(key string, value float64)
}

type gaugeRepository struct {
	store map[string]float64
}

func NewGaugeRepository() GaugeRepository {
	return &gaugeRepository{
		store: make(map[string]float64),
	}
}

func (g gaugeRepository) Update(key string, value float64) {
	g.store[key] = value
}
