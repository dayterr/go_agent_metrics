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

var Cfg config.Config
var CfgLogger config.ConfigLogger

func main() {
	Cfg = config.GetEnv()
	CfgLogger = config.GetEnvLogger()
	fmt.Println("before", Cfg, CfgLogger)
	if Cfg.Address == "localhost:8080" {
		flag.StringVar(&Cfg.Address, "a", Cfg.Address, "Address for the server")
	}
	_ = flag.Bool("r", CfgLogger.Restore, "A bool flag for configuration upload")
	if CfgLogger.StoreInterval == 300 * time.Second {
		flag.DurationVar(&CfgLogger.StoreInterval, "i", CfgLogger.StoreInterval, "Interval for saving the metrics into the file")
	} else {
		_ = flag.Duration("i", CfgLogger.StoreInterval, "Interval for saving the metrics into the file")
	}
	if CfgLogger.StoreFile == "/tmp/devops-metrics-db.json" {
		flag.StringVar(&CfgLogger.StoreFile, "f", CfgLogger.StoreFile, "file to store the metrics")
	}
	flag.Parse()
	fmt.Println("after", Cfg, CfgLogger)
	//fmt.Println(CfgLogger.Restore, Cfg.Address)
	var port = handlers.GetPort(Cfg.Address)
	fmt.Println(CfgLogger, Cfg)
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
