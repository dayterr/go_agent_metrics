package storage

import (
	"context"
	"database/sql"

	"github.com/dayterr/go_agent_metrics/internal/metric"
)

type Gauge float64
type Counter int64

type Storager interface {
	GetGuageByID(ctx context.Context, id string) (float64, error)
	GetCounterByID(ctx context.Context, id string) (int64, error)
	SetGuage(ctx context.Context, id string, v *float64)
	SetCounter(ctx context.Context, id string, v *int64)
	SetGaugeFromMemStats(ctx context.Context, id string, value float64)
	SetCounterFromMemStats(ctx context.Context, id string, value int64)
	GetGauges(ctx context.Context) (map[string]Gauge, error)
	GetCounters(ctx context.Context) (map[string]Counter, error)
	CheckGaugeByName(ctx context.Context, name string) bool
	CheckCounterByName(ctx context.Context, name string) bool
	SaveMany(ctx context.Context, metricsList []metric.Metrics) error
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
