package agent

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"log"
	"math/rand"
	"os"
	"runtime"
	"strconv"
	"time"

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

type Storage struct {
	GaugeField map[string]Gauge `json:"Gauge"`
	CounterField map[string]Counter `json"Counter"`
}

var allMetrics Storage = Storage{
	metrics,
	counters,
}

var (
	Address *string
	ReportInterval time.Duration
	PollInterval time.Duration
)

func init() {
	var err error
	Address = flag.String("a", os.Getenv("ADDRESS"), "Address for the server")
	repIntervalStr := flag.String("r", os.Getenv("REPORT_INTERVAL"), "Interval for sending the metrics to the server")
	ReportInterval, err = time.ParseDuration(*repIntervalStr)
	if err != nil {
		log.Fatal("Flag -r got an incorrect argument")
	}
	pollIntervalStr := flag.String("p", os.Getenv("POLL_INTERVAL"), "Interval for polling the metrics")
	PollInterval, err = time.ParseDuration(*pollIntervalStr)
	if err != nil {
		log.Fatal("Flag -p got an incorrect argument")
	}
	flag.Parse()
}

func ReadMetrics() Storage {
	m := &runtime.MemStats{}
	runtime.ReadMemStats(m)
	allMetrics.GaugeField["Alloc"] = Gauge(m.Alloc)
	allMetrics.GaugeField["BuckHashSys"] = Gauge(m.BuckHashSys)
	allMetrics.GaugeField["Frees"] = Gauge(m.Frees)
	allMetrics.GaugeField["GCCPUFraction"] = Gauge(m.GCCPUFraction)
	allMetrics.GaugeField["GCSys"] = Gauge(m.GCSys)
	allMetrics.GaugeField["HeapAlloc"] = Gauge(m.HeapAlloc)
	allMetrics.GaugeField["HeapIdle"] = Gauge(m.HeapIdle)
	allMetrics.GaugeField["HeapInuse"] = Gauge(m.HeapInuse)
	allMetrics.GaugeField["HeapObjects"] = Gauge(m.HeapObjects)
	allMetrics.GaugeField["HeapReleased"] = Gauge(m.HeapReleased)
	allMetrics.GaugeField["HeapSys"] = Gauge(m.HeapSys)
	allMetrics.GaugeField["LastGC"] = Gauge(m.HeapAlloc)
	allMetrics.GaugeField["Lookups"] = Gauge(m.Lookups)
	allMetrics.GaugeField["MCacheInuse"] = Gauge(m.MCacheInuse)
	allMetrics.GaugeField["MCacheSys"] = Gauge(m.MCacheSys)
	allMetrics.GaugeField["MSpanInuse"] = Gauge(m.MSpanInuse)
	allMetrics.GaugeField["MSpanSys"] = Gauge(m.MSpanSys)
	allMetrics.GaugeField["Mallocs"] = Gauge(m.Mallocs)
	allMetrics.GaugeField["NextGC"] = Gauge(m.NextGC)
	allMetrics.GaugeField["NumForcedGC"] = Gauge(m.NumForcedGC)
	allMetrics.GaugeField["NumGC"] = Gauge(m.NumGC)
	allMetrics.GaugeField["OtherSys"] = Gauge(m.OtherSys)
	allMetrics.GaugeField["PauseTotalNs"] = Gauge(m.PauseTotalNs)
	allMetrics.GaugeField["StackInuse"] = Gauge(m.StackInuse)
	allMetrics.GaugeField["StackSys"] = Gauge(m.StackSys)
	allMetrics.GaugeField["Sys"] = Gauge(m.Sys)
	allMetrics.GaugeField["TotalAlloc"] = Gauge(m.TotalAlloc)
	allMetrics.GaugeField["RandomValue"] = Gauge(rand.Float64())
	allMetrics.CounterField["PollCount"] += 1
	return allMetrics
}

func PostCounter(value Counter, metricName string, metricType string) error {
	url := fmt.Sprintf("http://%v/update/%v/%v/%v", Address, metricType, metricName, value)
	_, err := grequests.Post(url, &grequests.RequestOptions{Data: map[string]string{metricName: strconv.Itoa(int(value))},
		Headers: map[string]string{"ContentType": "text/plain"}})
	if err != nil {
		return err
	}
	url = fmt.Sprintf("http://%v/update", Address)
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
	url := fmt.Sprintf("http://%v/update/%v/%v/%v", Address, metricType, metricName, value)
	_, err := grequests.Post(url, &grequests.RequestOptions{Data: map[string]string{metricName: strconv.Itoa(int(value))},
		Headers: map[string]string{"ContentType": "text/plain"}})
	if err != nil {
		return err
	}
	url = fmt.Sprintf("http://%v/update", Address)
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


func PostAll(am Storage) {
	for k, v := range am.GaugeField {
		PostMetric(v, k, "gauge")
	}
	for k, v := range am.CounterField {
		PostCounter(v, k, "counter")
	}
}
