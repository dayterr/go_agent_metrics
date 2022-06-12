package server

import "github.com/dayterr/go_agent_metrics/internal/agent"

type MetricsJSON struct {
	Alloc agent.Gauge
	BuckHashSys agent.Gauge
	Frees agent.Gauge
	GCCPUFraction agent.Gauge
	GCSys agent.Gauge
	HeapAlloc agent.Gauge
	HeapIdle agent.Gauge
	HeapInuse agent.Gauge
	HeapObjects agent.Gauge
	HeapReleased agent.Gauge
	HeapSys agent.Gauge
	LastGC agent.Gauge
	Lookups agent.Gauge
	MCacheInuse agent.Gauge
	MCacheSys agent.Gauge
	MSpanInuse agent.Gauge
	MSpanSys agent.Gauge
	Mallocs agent.Gauge
	NextGC agent.Gauge
	NumForcedGC agent.Gauge
	NumGC agent.Gauge
	OtherSys agent.Gauge
	PauseTotalNs agent.Gauge
	StackInuse agent.Gauge
	StackSys agent.Gauge
	Sys agent.Gauge
	TotalAlloc agent.Gauge
	RandomValue agent.Gauge
	PollCount agent.Counter
}
