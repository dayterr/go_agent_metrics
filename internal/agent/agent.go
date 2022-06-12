package agent

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/dayterr/go_agent_metrics/internal/config"
	"strconv"

	//	"fmt"
	"math/rand"
	"runtime"
//	"strconv"

	"github.com/levigross/grequests"
)

const GaugeType = "gauge"
const CounterType = "counter"

type Gauge float64
type Counter int64

type Metrics struct {
	ID    string   `json:"id"`
	MType string   `json:"type"`
	Delta int64   `json:"delta,omitempty"`
	Value float64 `json:"value,omitempty"`
}

func (m *Metrics) MarshallJSON() ([]byte, error) {
	switch m.MType {
	case GaugeType:
		aliasValue := &struct {
			ID    string `json:"id"`
			MType string `json:"type"`
			Value float64 `json:"value"`
		}{
		ID: m.ID,
		MType: m.MType,
		Value: m.Value,
		}
		return json.Marshal(aliasValue)
	case CounterType:
		aliasValue := &struct {
			ID    string `json:"id"`
			MType string `json:"type"`
			Delta int64 `json:"delta"`
		}{
			ID: m.ID,
			MType: m.MType,
			Delta: m.Delta,
		}
		return json.Marshal(aliasValue)
	default:
		return nil, errors.New("no such metric type")
	}
}

var metrics = make(map[string]Gauge)
var counters = make(map[string]Counter)

func ReadMetrics() {
	m := &runtime.MemStats{}
	runtime.ReadMemStats(m)
	metrics["Alloc"] = Gauge(m.Alloc)
	metrics["BuckHashSys"] = Gauge(m.BuckHashSys)
	metrics["Frees"] = Gauge(m.Frees)
	metrics["GCCPUFraction"] = Gauge(m.GCCPUFraction)
	metrics["GCSys"] = Gauge(m.GCSys)
	metrics["HeapAlloc"] = Gauge(m.HeapAlloc)
	metrics["HeapIdle"] = Gauge(m.HeapIdle)
	metrics["HeapInuse"] = Gauge(m.HeapInuse)
	metrics["HeapObjects"] = Gauge(m.HeapObjects)
	metrics["HeapReleased"] = Gauge(m.HeapReleased)
	metrics["HeapSys"] = Gauge(m.HeapSys)
	metrics["LastGC"] = Gauge(m.HeapAlloc)
	metrics["Lookups"] = Gauge(m.Lookups)
	metrics["MCacheInuse"] = Gauge(m.MCacheInuse)
	metrics["MCacheSys"] = Gauge(m.MCacheSys)
	metrics["MSpanInuse"] = Gauge(m.MSpanInuse)
	metrics["MSpanSys"] = Gauge(m.MSpanSys)
	metrics["Mallocs"] = Gauge(m.Mallocs)
	metrics["NextGC"] = Gauge(m.NextGC)
	metrics["NumForcedGC"] = Gauge(m.NumForcedGC)
	metrics["NumGC"] = Gauge(m.NumGC)
	metrics["OtherSys"] = Gauge(m.OtherSys)
	metrics["PauseTotalNs"] = Gauge(m.PauseTotalNs)
	metrics["StackInuse"] = Gauge(m.StackInuse)
	metrics["StackSys"] = Gauge(m.StackSys)
	metrics["Sys"] = Gauge(m.Sys)
	metrics["TotalAlloc"] = Gauge(m.TotalAlloc)
	metrics["RandomValue"] = Gauge(rand.Float64())
	counters["PollCount"] += 1
}

func PostCounter(value Counter, metricName string, metricType string) error {
	cfg := config.GetEnv()
	url := fmt.Sprintf("http://%v/update/%v/%v/%v", cfg.Address, metricType, metricName, value)
	_, err := grequests.Post(url, &grequests.RequestOptions{Data: map[string]string{metricName: strconv.Itoa(int(value))},
		Headers: map[string]string{"ContentType": "text/plain"}})
	if err != nil {
		return err
	}
	url = fmt.Sprintf("http://%v/update", cfg.Address)
	metric := Metrics{ID: metricName, MType: metricType, Delta: int64(value)}
	mJSON, err := metric.MarshallJSON()
	if err != nil {
		return err
	}
	_, err = grequests.Post(url, &grequests.RequestOptions{JSON: mJSON,
		Headers: map[string]string{"ContentType": "application/json"}})
	if err != nil {
		return err
	}
	return nil
}

func PostMetric(value Gauge, metricName string, metricType string) error {
	cfg := config.GetEnv()
	url := fmt.Sprintf("http://%v/update/%v/%v/%v", cfg.Address, metricType, metricName, value)
	_, err := grequests.Post(url, &grequests.RequestOptions{Data: map[string]string{metricName: strconv.Itoa(int(value))},
		Headers: map[string]string{"ContentType": "text/plain"}})
	if err != nil {
		return err
	}
	url = fmt.Sprintf("http://%v/update", cfg.Address)
	metric := Metrics{ID: metricName, MType: metricType, Value: float64(value)}
	mJSON, err := metric.MarshallJSON()
	if err != nil {
		return err
	}
   	_, err = grequests.Post(url, &grequests.RequestOptions{JSON: mJSON,
   		Headers: map[string]string{"ContentType": "application/json"}})
	if err != nil {
		return err
	}
	return nil
}


func PostAll() {
	for k, v := range metrics {
		PostMetric(v, k, "gauge")
	}
	for k, v := range counters {
		PostCounter(v, k, "counter")
	}
}
