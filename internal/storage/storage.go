package storage

import (
	"github.com/dayterr/go_agent_metrics/internal/metric"
)

func NewIMS() InMemoryStorage {
	return InMemoryStorage{
		GaugeField:   make(map[string]Gauge),
		CounterField: make(map[string]Counter),
	}
}

func (g Gauge) ToFloat() float64 {
	return float64(g)
}

func (c Counter) ToInt64() int64 {
	return int64(c)
}

func (c Counter) ToInt() int {
	return int(c)
}

func (s InMemoryStorage) GetGuageByID(id string) (float64, error) {
	v := s.GaugeField[id].ToFloat()
	return v, nil
}

func (s InMemoryStorage) GetCounterByID(id string) (int64, error) {
	v := s.CounterField[id].ToInt64()
	return v, nil
}

func (s InMemoryStorage) SetGuage(id string, v *float64) {
	s.GaugeField[id] = Gauge(*v)
}

func (s InMemoryStorage) SetCounter(id string, v *int64) {
	s.CounterField[id] += Counter(*v)
}

func (s InMemoryStorage) SetGaugeFromMemStats(id string, value float64) {
	s.GaugeField[id] = Gauge(value)
}

func (s InMemoryStorage) SetCounterFromMemStats(id string, value int64) {
	s.CounterField[id] += Counter(value)
}

func (s InMemoryStorage) GetGauges() map[string]Gauge {
	return s.GaugeField
}

func (s InMemoryStorage) GetCounters() map[string]Counter {
	return s.CounterField
}

func (s InMemoryStorage) CheckGaugeByName(name string) bool {
	_, ok := s.GaugeField[name]
	return ok
}

func (s InMemoryStorage) CheckCounterByName(name string) bool {
	_, ok := s.CounterField[name]
	return ok
}

func (s InMemoryStorage) SaveMany(metricsList []metric.Metrics) error {
	for _, metric := range metricsList {
		if metric.MType == "gauge" {
			s.GaugeField[metric.ID] = Gauge(*metric.Value)
		} else {
			s.CounterField[metric.ID] = Counter(*metric.Delta)
		}
	}
	return nil
}
