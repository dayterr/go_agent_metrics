package main

import (
	"github.com/dayterr/go_agent_metrics/cmd/server/handlers"
	"github.com/dayterr/go_agent_metrics/internal/agent"
	"github.com/dayterr/go_agent_metrics/internal/config"
	"net/http"
)

var metrics = make(map[string]agent.Gauge)
var counters = make(map[string]agent.Counter)

var allMetrics agent.Storage = agent.Storage{
	metrics,
	counters,
}

var port = config.GetPort()

func main() {
	/*cfg := config.GetEnvLogger()
	ticker := time.NewTicker(cfg.StoreInterval)
	go func() {
		for {
			select {
			case <- ticker.C:
				server.WriteJSON(cfg.StoreFile)
			}
		}
	}()*/
	r := handlers.CreateRouter()
	http.ListenAndServe(port, r)
}
