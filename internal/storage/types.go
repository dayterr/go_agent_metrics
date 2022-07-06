package storage

type Gauge float64
type Counter int64

type Storager interface {
	LoadMetricsFromFile(filename string) error
	GetGuageByID(id string) float64
	GetCounterByID(id string) int64
	SetGuage(id string, v *float64)
	SetCounter(id string, v *int64)
	SetGaugeFromMemStats(id string, value float64)
	SetCounterFromMemStats(id string, value int64)
	ReadMetrics()
	GetGauges() map[string]Gauge
	GetCounters() map[string]Counter
	CheckGaugeByName(name string) bool
	CheckCounterByName(name string) bool
}

type InMemoryStorage struct {
	GaugeField   map[string]Gauge
	CounterField map[string]Counter
}

type DBStorage struct {
	GaugeField   map[string]Gauge
	CounterField map[string]Counter
	DSN string
}
