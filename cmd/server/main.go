package main

import (
	"flag"
	"github.com/dayterr/go_agent_metrics/cmd/server/handlers"
	"github.com/dayterr/go_agent_metrics/internal/agent"
	"github.com/dayterr/go_agent_metrics/internal/config"
	"github.com/dayterr/go_agent_metrics/internal/server"
	"net/http"
	"time"
)

var metrics = make(map[string]agent.Gauge)
var counters = make(map[string]agent.Counter)

var Cfg = config.GetEnv()
var CfgLogger = config.GetEnvLogger()


func main() {
	flag.StringVar(&Cfg.Address, "a", Cfg.Address, "Address for the server")
	flag.BoolVar(&CfgLogger.Restore, "r", CfgLogger.Restore, "A bool flag for configuration upload")
	flag.DurationVar(&CfgLogger.StoreInterval, "i", CfgLogger.StoreInterval, "Interval for saving the metrics into the file")
	flag.StringVar(&CfgLogger.StoreFile, "f", CfgLogger.StoreFile, "file to store the metrics")
	flag.Parse()
	var port = handlers.GetPort(Cfg.Address)
	ticker := time.NewTicker(CfgLogger.StoreInterval)
	go func() {
		for {
			select {
			case <- ticker.C:
				server.WriteJSON(CfgLogger.StoreFile)
			}
		}
	}()
	r := handlers.CreateRouter(CfgLogger.StoreFile, CfgLogger.Restore)
	http.ListenAndServe(port, r)
}
