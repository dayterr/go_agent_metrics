package agent

import (
	"encoding/json"
	"fmt"
	"github.com/levigross/grequests"
	"math/rand"
	"runtime"

	"github.com/dayterr/go_agent_metrics/internal/storage"
)

const GaugeType = "gauge"
const CounterType = "counter"

type Metrics struct {
	ID    string  `json:"id"`
	MType string  `json:"type"`
	Delta int64   `json:"delta,omitempty"`
	Value float64 `json:"value,omitempty"`
}

func ReadMetrics() storage.Storage {
	var allMetrics = storage.New()
	m := &runtime.MemStats{}
	runtime.ReadMemStats(m)
	allMetrics.GaugeField["Alloc"] = storage.Gauge(m.Alloc)
	allMetrics.GaugeField["BuckHashSys"] = storage.Gauge(m.BuckHashSys)
	allMetrics.GaugeField["Frees"] = storage.Gauge(m.Frees)
	allMetrics.GaugeField["GCCPUFraction"] = storage.Gauge(m.GCCPUFraction)
	allMetrics.GaugeField["GCSys"] = storage.Gauge(m.GCSys)
	allMetrics.GaugeField["HeapAlloc"] = storage.Gauge(m.HeapAlloc)
	allMetrics.GaugeField["HeapIdle"] = storage.Gauge(m.HeapIdle)
	allMetrics.GaugeField["HeapInuse"] = storage.Gauge(m.HeapInuse)
	allMetrics.GaugeField["HeapObjects"] = storage.Gauge(m.HeapObjects)
	allMetrics.GaugeField["HeapReleased"] = storage.Gauge(m.HeapReleased)
	allMetrics.GaugeField["HeapSys"] = storage.Gauge(m.HeapSys)
	allMetrics.GaugeField["LastGC"] = storage.Gauge(m.HeapAlloc)
	allMetrics.GaugeField["Lookups"] = storage.Gauge(m.Lookups)
	allMetrics.GaugeField["MCacheInuse"] = storage.Gauge(m.MCacheInuse)
	allMetrics.GaugeField["MCacheSys"] = storage.Gauge(m.MCacheSys)
	allMetrics.GaugeField["MSpanInuse"] = storage.Gauge(m.MSpanInuse)
	allMetrics.GaugeField["MSpanSys"] = storage.Gauge(m.MSpanSys)
	allMetrics.GaugeField["Mallocs"] = storage.Gauge(m.Mallocs)
	allMetrics.GaugeField["NextGC"] = storage.Gauge(m.NextGC)
	allMetrics.GaugeField["NumForcedGC"] = storage.Gauge(m.NumForcedGC)
	allMetrics.GaugeField["NumGC"] = storage.Gauge(m.NumGC)
	allMetrics.GaugeField["OtherSys"] = storage.Gauge(m.OtherSys)
	allMetrics.GaugeField["PauseTotalNs"] = storage.Gauge(m.PauseTotalNs)
	allMetrics.GaugeField["StackInuse"] = storage.Gauge(m.StackInuse)
	allMetrics.GaugeField["StackSys"] = storage.Gauge(m.StackSys)
	allMetrics.GaugeField["Sys"] = storage.Gauge(m.Sys)
	allMetrics.GaugeField["TotalAlloc"] = storage.Gauge(m.TotalAlloc)
	allMetrics.GaugeField["RandomValue"] = storage.Gauge(rand.Float64())
	allMetrics.CounterField["PollCount"] += 1
	return allMetrics
}

func PostCounter(value storage.Counter, metricName string, address string) error {
	/*url := fmt.Sprintf("http://%v/update/%v/%v/%v", address, CounterType, metricName, value)
	_, err := grequests.Post(url, &grequests.RequestOptions{Data: map[string]string{metricName: strconv.Itoa(int(value))},
		Headers: map[string]string{"ContentType": "text/plain"}})
	if err != nil {
		return err
	}*/
	url := fmt.Sprintf("http://%v/update", address)
	//url := fmt.Sprintf("http://%v/update/%v/%v/%v", address, CounterType, metricName, value)
	metric := Metrics{ID: metricName, MType: CounterType, Delta: int64(value)}
	mJSON, err := json.Marshal(metric)
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

func PostGauge(value storage.Gauge, metricName string, address string) error {
	/*url := fmt.Sprintf("http://%v/update/%v/%v/%v", address, storage.GaugeType, metricName, value)
	_, err := grequests.Post(url, &grequests.RequestOptions{Data: map[string]string{metricName: strconv.Itoa(int(value))},
		Headers: map[string]string{"ContentType": "text/plain"}})
	if err != nil {
		return err
	}*/
	url := fmt.Sprintf("http://%v/update", address)
	//url := fmt.Sprintf("http://%v/update/%v/%v/%v", address, GaugeType, metricName, value)
	metric := Metrics{ID: metricName, MType: GaugeType, Value: float64(value)}
	mJSON, err := json.Marshal(metric)
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

func PostAll(am storage.Storage, address string) {
	for k, v := range am.GaugeField {
		PostGauge(v, k, address)
	}
	for k, v := range am.CounterField {
		PostCounter(v, k, address)
	}
}
