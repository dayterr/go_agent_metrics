package main

import (
	"flag"
	"fmt"
	"github.com/dayterr/go_agent_metrics/cmd/server/handlers"
	"github.com/dayterr/go_agent_metrics/internal/agent"
	"github.com/dayterr/go_agent_metrics/internal/config"
	"github.com/dayterr/go_agent_metrics/internal/server"
	"net/http"
	"time"
)

var metrics = make(map[string]agent.Gauge)
var counters = make(map[string]agent.Counter)

var (
	Addr *string
	Restore *bool
	StoreFile *string
	StoreInterval *time.Duration
)

func init() {
	cfg := config.GetEnv()
	cfgLogger := config.GetEnvLogger()
	Addr = flag.String("a", cfg.Address, "Address for the server")
	Restore = flag.Bool("r", cfgLogger.Restore, "A bool flag for configuration upload")
	StoreInterval = flag.Duration("i", cfgLogger.StoreInterval, "Interval for saving the metrics into the file")
	StoreFile = flag.String("f", cfgLogger.StoreFile, "file to store the metrics")
	flag.Parse()
}

func main() {
	var port = handlers.GetPort(*Addr)
	ticker := time.NewTicker(*StoreInterval)
	go func() {
		for {
			select {
			case <- ticker.C:
				server.WriteJSON(*StoreFile)
			}
		}
	}()
	r := handlers.CreateRouter(*StoreFile, *Restore)
	http.ListenAndServe(port, r)
}
