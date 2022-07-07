package storage

import "github.com/dayterr/go_agent_metrics/internal/metric"

type Gauge float64
type Counter int64

type Storager interface {
	LoadMetricsFromFile(filename string) error
	GetGuageByID(id string) (float64, error)
	GetCounterByID(id string) (int64, error)
	SetGuage(id string, v *float64)
	SetCounter(id string, v *int64)
	SetGaugeFromMemStats(id string, value float64)
	SetCounterFromMemStats(id string, value int64)
	ReadMetrics()
	GetGauges() map[string]Gauge
	GetCounters() map[string]Counter
	CheckGaugeByName(name string) bool
	CheckCounterByName(name string) bool
	SaveMany(metricsList []metric.Metrics) error
}

type InMemoryStorage struct {
	GaugeField   map[string]Gauge `json:"Gauge"`
	CounterField map[string]Counter `json:"Counter"`
}

type DBStorage struct {
	GaugeField   map[string]Gauge `json:"Gauge"`
	CounterField map[string]Counter `json:"Counter"`
	DSN string
}
