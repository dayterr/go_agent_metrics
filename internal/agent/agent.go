package agent

import (
	"encoding/json"
	"fmt"
	"github.com/dayterr/go_agent_metrics/internal/hash"
	"github.com/dayterr/go_agent_metrics/internal/metric"
	"github.com/dayterr/go_agent_metrics/internal/storage"
	"github.com/levigross/grequests"
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
