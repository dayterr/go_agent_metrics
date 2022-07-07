package storage

import (
	"encoding/json"
	"fmt"
	"github.com/dayterr/go_agent_metrics/internal/metric"
	"io/ioutil"
	"math/rand"
	"os"
	"runtime"
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

func (s InMemoryStorage) LoadMetricsFromFile(filename string) error {
	if _, err := os.Stat(filename); err != nil {
		file, err := os.Create(filename)
		if err != nil {
			return err
		}
		file.Close()
		return nil
	}
	file, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}
	err = json.Unmarshal(file, &s)
	fmt.Println(s)
	if err != nil {
		return err
	}
	return nil
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

func (s InMemoryStorage) ReadMetrics() {
	m := &runtime.MemStats{}
	runtime.ReadMemStats(m)
	s.SetGaugeFromMemStats("Alloc", float64(m.Alloc))
	s.SetGaugeFromMemStats("BuckHashSys", float64(m.BuckHashSys))
	s.SetGaugeFromMemStats("Frees", float64(m.Frees))
	s.SetGaugeFromMemStats("GCCPUFraction", m.GCCPUFraction)
	s.SetGaugeFromMemStats("GCSys", float64(m.GCSys))
	s.SetGaugeFromMemStats("HeapAlloc", float64(m.HeapAlloc))
	s.SetGaugeFromMemStats("HeapIdle", float64(m.HeapIdle))
	s.SetGaugeFromMemStats("HeapInuse", float64(m.HeapInuse))
	s.SetGaugeFromMemStats("HeapObjects", float64(m.HeapObjects))
	s.SetGaugeFromMemStats("HeapReleased", float64(m.HeapReleased))
	s.SetGaugeFromMemStats("HeapSys", float64(m.HeapSys))
	s.SetGaugeFromMemStats("LastGC", float64(m.HeapAlloc))
	s.SetGaugeFromMemStats("Lookups", float64(m.Lookups))
	s.SetGaugeFromMemStats("MCacheInuse", float64(m.MCacheInuse))
	s.SetGaugeFromMemStats("MCacheSys", float64(m.MCacheSys))
	s.SetGaugeFromMemStats("MSpanInuse", float64(m.MSpanInuse))
	s.SetGaugeFromMemStats("MSpanSys", float64(m.MSpanSys))
	s.SetGaugeFromMemStats("Mallocs", float64(m.Mallocs))
	s.SetGaugeFromMemStats("NextGC", float64(m.NextGC))
	s.SetGaugeFromMemStats("NumForcedGC", float64(m.NumForcedGC))
	s.SetGaugeFromMemStats("NumGC", float64(m.NumGC))
	s.SetGaugeFromMemStats("OtherSys", float64(m.OtherSys))
	s.SetGaugeFromMemStats("PauseTotalNs", float64(m.PauseTotalNs))
	s.SetGaugeFromMemStats("StackInuse", float64(m.StackInuse))
	s.SetGaugeFromMemStats("StackSys", float64(m.StackSys))
	s.SetGaugeFromMemStats("Sys", float64(m.Sys))
	s.SetGaugeFromMemStats("TotalAlloc", float64(m.TotalAlloc))
	s.SetGaugeFromMemStats("RandomValue", rand.Float64())
	s.SetCounterFromMemStats("PollCount", 1)
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