package agent

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/dayterr/go_agent_metrics/internal/config"
	"reflect"
	"strconv"
	"strings"

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

type MetricsJSON struct {
	Alloc Gauge
	BuckHashSys Gauge
	Frees Gauge
	GCCPUFraction Gauge
	GCSys Gauge
	HeapAlloc Gauge
	HeapIdle Gauge
	HeapInuse Gauge
	HeapObjects Gauge
	HeapReleased Gauge
	HeapSys Gauge
	LastGC Gauge
	Lookups Gauge
	MCacheInuse Gauge
	MCacheSys Gauge
	MSpanInuse Gauge
	MSpanSys Gauge
	Mallocs Gauge
	NextGC Gauge
	NumForcedGC Gauge
	NumGC Gauge
	OtherSys Gauge
	PauseTotalNs Gauge
	StackInuse Gauge
	StackSys Gauge
	Sys Gauge
	TotalAlloc Gauge
	RandomValue Gauge
	testGauge Gauge
	PollCount Counter
	testCounter Counter
}


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

func ReadMetrics() MetricsJSON {
	m := &runtime.MemStats{}
	mj := MetricsJSON{}
	runtime.ReadMemStats(m)
	mj.Alloc = Gauge(m.Alloc)
	mj.BuckHashSys = Gauge(m.BuckHashSys)
	mj.Frees = Gauge(m.Frees)
	mj.GCCPUFraction = Gauge(m.GCCPUFraction)
	mj.GCSys = Gauge(m.GCSys)
	mj.HeapAlloc = Gauge(m.HeapAlloc)
	mj.HeapIdle = Gauge(m.HeapIdle)
	mj.HeapInuse = Gauge(m.HeapInuse)
	mj.HeapObjects = Gauge(m.HeapObjects)
	mj.HeapReleased = Gauge(m.HeapReleased)
	mj.HeapSys = Gauge(m.HeapSys)
	mj.LastGC = Gauge(m.HeapAlloc)
	mj.Lookups = Gauge(m.Lookups)
	mj.MCacheInuse = Gauge(m.MCacheInuse)
	mj.MCacheSys = Gauge(m.MCacheSys)
	mj.MSpanInuse = Gauge(m.MSpanInuse)
	mj.MSpanSys = Gauge(m.MSpanSys)
	mj.Mallocs = Gauge(m.Mallocs)
	mj.NextGC = Gauge(m.NextGC)
	mj.NumForcedGC = Gauge(m.NumForcedGC)
	mj.NumGC = Gauge(m.NumGC)
	mj.OtherSys = Gauge(m.OtherSys)
	mj.PauseTotalNs = Gauge(m.PauseTotalNs)
	mj.StackInuse = Gauge(m.StackInuse)
	mj.StackSys = Gauge(m.StackSys)
	mj.Sys = Gauge(m.Sys)
	mj.TotalAlloc = Gauge(m.TotalAlloc)
	mj.RandomValue = Gauge(rand.Float64())
	mj.PollCount += 1
	return mj
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


func PostAll(mj MetricsJSON) {
	v := reflect.ValueOf(&mj).Elem()
	for i := 0; i < v.NumField(); i++ {
		metricName := v.Type().Field(i).Name
		metricType := strings.ToLower(v.Type().Field(i).Type.Name())
		switch metricType {
		case GaugeType:
			metricGauge := v.Field(i).Float()
			PostMetric(Gauge(metricGauge), metricName, metricType)
		case CounterType:
			metricCounter := v.Field(i).Int()
			PostCounter(Counter(metricCounter), metricName, metricType)
		}

	}
}
