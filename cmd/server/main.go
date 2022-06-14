package main

import (
	"fmt"
	"github.com/dayterr/go_agent_metrics/internal/agent"
	"github.com/dayterr/go_agent_metrics/internal/server"
	"net/http"
	"time"

	"github.com/dayterr/go_agent_metrics/cmd/server/handlers"
	"github.com/dayterr/go_agent_metrics/internal/config"
)

var metrics = make(map[string]agent.Gauge)
var counters = make(map[string]agent.Counter)

var allMetrics agent.Storage = agent.Storage{
	metrics,
	counters,
}

var port = config.GetPort()

func main() {
	cfg := config.GetEnvLogger()
	ticker := time.NewTicker(cfg.StoreInterval)
	//l, _ := os.Getwd()
	go func() {
		server.LoadMetricsFromJSON(cfg, allMetrics)
		fmt.Println("working 2")
		for {
			select {
			case <- ticker.C:
				server.WriteJSON(cfg.StoreFile)
			}
		}
	}()
	r := handlers.CreateRouter()
	http.ListenAndServe(port, r)
}
