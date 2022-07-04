package agent

import (
	"encoding/json"
	"fmt"
	"github.com/dayterr/go_agent_metrics/internal/storage"
	"github.com/levigross/grequests"
	"log"
)

const GaugeType = "gauge"
const CounterType = "counter"

type Metrics struct {
	ID    string  `json:"id"`
	MType string  `json:"type"`
	Delta *int64   `json:"delta,omitempty"`
	Value *float64 `json:"value,omitempty"`
}

func PostCounter(value storage.Counter, metricName string, address string) error {
	url := fmt.Sprintf("http://%v/update", address)
	delta := value.ToInt64()
	metric := Metrics{ID: metricName, MType: CounterType, Delta: &delta}
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
	url := fmt.Sprintf("http://%v/update", address)
	v := value.ToFloat()
	metric := Metrics{ID: metricName, MType: GaugeType, Value: &v}
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
	log.Println("got gauges", gauges)
	log.Println("got counters", counters)
	for k, v := range gauges {
		PostGauge(v, k, a.Address)
	}
	for k, v := range counters {
		PostCounter(v, k, a.Address)
	}
}
