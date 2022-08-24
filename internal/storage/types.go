package storage

import (
	"database/sql"
	"github.com/dayterr/go_agent_metrics/internal/metric"
)

type Gauge float64
type Counter int64

type Storager interface {
	GetGuageByID(id string) (float64, error)
	GetCounterByID(id string) (int64, error)
	SetGuage(id string, v *float64)
	SetCounter(id string, v *int64)
	SetGaugeFromMemStats(id string, value float64)
	SetCounterFromMemStats(id string, value int64)
	GetGauges() map[string]Gauge
	GetCounters() map[string]Counter
	CheckGaugeByName(name string) bool
	CheckCounterByName(name string) bool
	SaveMany(metricsList []metric.Metrics) error
}

type InMemoryStorage struct {
	GaugeField   map[string]Gauge   `json:"Gauge"`
	CounterField map[string]Counter `json:"Counter"`
}

type DBStorage struct {
	DB           *sql.DB
	DSN          string
	GaugeField   map[string]Gauge
	CounterField map[string]Counter
}
