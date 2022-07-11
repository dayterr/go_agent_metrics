package agent

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/dayterr/go_agent_metrics/internal/hash"
	"github.com/dayterr/go_agent_metrics/internal/metric"
	"github.com/dayterr/go_agent_metrics/internal/storage"
	"github.com/levigross/grequests"
	"math/rand"
	"runtime"
)

const GaugeType = "gauge"
const CounterType = "counter"

func PostCounter(value storage.Counter, metricName string, address string, key string) error {
	url := fmt.Sprintf("http://%v/update/", address)
	delta := value.ToInt64()
	metric := metric.Metrics{ID: metricName, MType: CounterType, Delta: &delta}
	if key != "" {
		metric.Hash = hash.EncryptMetric(metric, key)
	}
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

func PostGauge(value storage.Gauge, metricName string, address string, key string) error {
	url := fmt.Sprintf("http://%v/update/", address)
	v := value.ToFloat()
	metric := metric.Metrics{ID: metricName, MType: GaugeType, Value: &v}
	if key != "" {
		metric.Hash = hash.EncryptMetric(metric, key)
	}
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

func (a Agent) PostAll() {
	gauges := a.Storage.GetGauges()
	counters := a.Storage.GetCounters()

	for k, v := range gauges {
		PostGauge(v, k, a.Address, a.Key)
	}
	for k, v := range counters {
		PostCounter(v, k, a.Address, a.Key)
	}
}

func (a Agent) ReadMetrics() {
	m := &runtime.MemStats{}
	runtime.ReadMemStats(m)
	a.Storage.SetGaugeFromMemStats("Alloc", float64(m.Alloc))
	a.Storage.SetGaugeFromMemStats("BuckHashSys", float64(m.BuckHashSys))
	a.Storage.SetGaugeFromMemStats("Frees", float64(m.Frees))
	a.Storage.SetGaugeFromMemStats("GCCPUFraction", m.GCCPUFraction)
	a.Storage.SetGaugeFromMemStats("GCSys", float64(m.GCSys))
	a.Storage.SetGaugeFromMemStats("HeapAlloc", float64(m.HeapAlloc))
	a.Storage.SetGaugeFromMemStats("HeapIdle", float64(m.HeapIdle))
	a.Storage.SetGaugeFromMemStats("HeapInuse", float64(m.HeapInuse))
	a.Storage.SetGaugeFromMemStats("HeapObjects", float64(m.HeapObjects))
	a.Storage.SetGaugeFromMemStats("HeapReleased", float64(m.HeapReleased))
	a.Storage.SetGaugeFromMemStats("HeapSys", float64(m.HeapSys))
	a.Storage.SetGaugeFromMemStats("LastGC", float64(m.HeapAlloc))
	a.Storage.SetGaugeFromMemStats("Lookups", float64(m.Lookups))
	a.Storage.SetGaugeFromMemStats("MCacheInuse", float64(m.MCacheInuse))
	a.Storage.SetGaugeFromMemStats("MCacheSys", float64(m.MCacheSys))
	a.Storage.SetGaugeFromMemStats("MSpanInuse", float64(m.MSpanInuse))
	a.Storage.SetGaugeFromMemStats("MSpanSys", float64(m.MSpanSys))
	a.Storage.SetGaugeFromMemStats("Mallocs", float64(m.Mallocs))
	a.Storage.SetGaugeFromMemStats("NextGC", float64(m.NextGC))
	a.Storage.SetGaugeFromMemStats("NumForcedGC", float64(m.NumForcedGC))
	a.Storage.SetGaugeFromMemStats("NumGC", float64(m.NumGC))
	a.Storage.SetGaugeFromMemStats("OtherSys", float64(m.OtherSys))
	a.Storage.SetGaugeFromMemStats("PauseTotalNs", float64(m.PauseTotalNs))
	a.Storage.SetGaugeFromMemStats("StackInuse", float64(m.StackInuse))
	a.Storage.SetGaugeFromMemStats("StackSys", float64(m.StackSys))
	a.Storage.SetGaugeFromMemStats("Sys", float64(m.Sys))
	a.Storage.SetGaugeFromMemStats("TotalAlloc", float64(m.TotalAlloc))
	a.Storage.SetGaugeFromMemStats("RandomValue", rand.Float64())
	a.Storage.SetCounterFromMemStats("PollCount", 1)
}

func (a Agent) PostMany() error {
	var listMetrics []metric.Metrics

	if len(a.Storage.GetGauges()) == 0 && len(a.Storage.GetCounters()) == 0 {
		return errors.New("the batch is empty")
	}

	for key, value := range a.Storage.GetGauges() {
		var m metric.Metrics
		m.ID = key
		v := value.ToFloat()
		m.Value = &v
		m.MType = GaugeType
		if a.Key != "" {
			m.Hash = hash.EncryptMetric(m, key)
		}
		listMetrics = append(listMetrics, m)
	}
	for key, value := range a.Storage.GetCounters() {
		var m metric.Metrics
		m.ID = key
		d := value.ToInt64()
		m.Delta = &d
		m.MType = CounterType
		if a.Key != "" {
			m.Hash = hash.EncryptMetric(m, key)
		}
		listMetrics = append(listMetrics, m)
	}

	jsn, err := json.Marshal(listMetrics)
	if err != nil {

		return err
	}

	url := fmt.Sprintf("http://%v/updates/", a.Address)
	_, err = grequests.Post(url, &grequests.RequestOptions{JSON: jsn,
		Headers: map[string]string{"ContentType": "application/json"}, DisableCompression: false})
	if err != nil {
		return err
	}

	return nil
}
