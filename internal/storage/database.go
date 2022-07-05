package storage

import (
	"encoding/json"
	"io/ioutil"
	"math/rand"
	"os"
	"runtime"
)

func (s DBStorage) LoadMetricsFromFile(filename string) error {
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
	if err != nil {
		return err
	}
	return nil
}

func (s DBStorage) GetGuageByID(id string) float64 {
	return 0
}

func (s DBStorage) GetCounterByID(id string) int64 {
	return 0
}

func (s DBStorage) SetGuage(id string, v *float64) {
	s.GaugeField[id] = Gauge(*v)
}

func (s DBStorage) SetCounter(id string, v *int64) {
	s.CounterField[id] += Counter(*v)
}

func (s DBStorage) SetGaugeFromMemStats(id string, value float64) {
	s.GaugeField[id] = Gauge(value)
}

func (s DBStorage) SetCounterFromMemStats(id string, value int64) {
	s.CounterField[id] += Counter(value)
}

func (s DBStorage) ReadMetrics() {
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

func (s DBStorage) GetGauges() map[string]Gauge {
	return s.GaugeField
}

func (s DBStorage) GetCounters() map[string]Counter {
	return s.CounterField
}

func (s DBStorage) CheckGaugeByName(name string) bool {
	_, ok := s.GaugeField[name]
	return ok
}

func (s DBStorage) CheckCounterByName(name string) bool {
	_, ok := s.CounterField[name]
	return ok
}
