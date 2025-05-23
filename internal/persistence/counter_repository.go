package persistence

type CounterRepository interface {
	Increment(key string, value int64)
}

type counterRepository struct {
	store map[string]int64
}

func NewCounterRepository() CounterRepository {
	return &counterRepository{
		store: make(map[string]int64),
	}
}

func (r *counterRepository) Increment(key string, value int64) {
	if val, ok := r.store[key]; ok {
		r.store[key] = val + value
	} else {
		r.store[key] = value
	}

}
