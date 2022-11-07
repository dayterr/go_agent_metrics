package storage

import (
	"context"
	"errors"
	"github.com/dayterr/go_agent_metrics/internal/metric"
	"log"
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

func (s InMemoryStorage) GetGuageByID(ctx context.Context, id string) (float64, error) {
	select {
	case <-ctx.Done():
		return 0, errors.New("the request was cancelled")
	default:
		v := s.GaugeField[id].ToFloat()
		return v, nil
	}
}

func (s InMemoryStorage) GetCounterByID(ctx context.Context, id string) (int64, error) {
	select {
	case <-ctx.Done():
		return 0, errors.New("the request was cancelled")
	default:
		v := s.CounterField[id].ToInt64()
		return v, nil
	}
}

func (s InMemoryStorage) SetGuage(ctx context.Context, id string, v *float64) {
	select {
	case <-ctx.Done():
		log.Println(errors.New("the request was cancelled"))
	default:
		s.GaugeField[id] = Gauge(*v)
	}
}

func (s InMemoryStorage) SetCounter(ctx context.Context, id string, v *int64) {
	select {
	case <-ctx.Done():
		log.Println(errors.New("the request was cancelled"))
	default:
		s.CounterField[id] += Counter(*v)
	}
}

func (s InMemoryStorage) SetGaugeFromMemStats(ctx context.Context, id string, value float64) {
	select {
	case <-ctx.Done():
		log.Println(errors.New("the request was cancelled"))
	default:
		s.GaugeField[id] = Gauge(value)
	}
}

func (s InMemoryStorage) SetCounterFromMemStats(ctx context.Context, id string, value int64) {
	select {
	case <-ctx.Done():
		log.Println(errors.New("the request was cancelled"))
	default:
		s.CounterField[id] += Counter(value)
	}
}

func (s InMemoryStorage) GetGauges(ctx context.Context) (map[string]Gauge, error) {
	select {
	case <-ctx.Done():
		return map[string]Gauge{}, errors.New("the request was cancelled")
	default:
		return s.GaugeField, nil
	}
}

func (s InMemoryStorage) GetCounters(ctx context.Context) (map[string]Counter, error) {
	select {
	case <-ctx.Done():
		return map[string]Counter{}, errors.New("the request was cancelled")
	default:
		return s.CounterField, nil
	}
}

func (s InMemoryStorage) CheckGaugeByName(ctx context.Context, name string) bool {
	select {
	case <-ctx.Done():
		return false
	default:
		_, ok := s.GaugeField[name]
		return ok
	}
}

func (s InMemoryStorage) CheckCounterByName(ctx context.Context, name string) bool {
	select {
	case <-ctx.Done():
		return false
	default:
		_, ok := s.CounterField[name]
		return ok
	}
}

func (s InMemoryStorage) SaveMany(ctx context.Context, metricsList []metric.Metrics) error {
	select {
	case <-ctx.Done():
		return errors.New("the request was cancelled")
	default:
		for _, metric := range metricsList {
			if metric.MType == "gauge" {
				s.GaugeField[metric.ID] = Gauge(*metric.Value)
			} else {
				s.CounterField[metric.ID] = Counter(*metric.Delta)
			}
		}
		return nil
	}
}
