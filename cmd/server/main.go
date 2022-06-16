package main

import (
	"flag"
	"fmt"
	"github.com/dayterr/go_agent_metrics/cmd/server/handlers"
	"github.com/dayterr/go_agent_metrics/internal/agent"
	"github.com/dayterr/go_agent_metrics/internal/config"
	"github.com/dayterr/go_agent_metrics/internal/server"
	"log"
	"net/http"
	"os"
	"time"
)

var metrics = make(map[string]agent.Gauge)
var counters = make(map[string]agent.Counter)

var (
	Addr *string
	Restore *bool
	StoreFile *string
	StoreInterval time.Duration
)

func init() {
	var err error
	cfg := config.GetEnv()
	cfgLogger := config.GetEnvLogger()
	serverFlags := flag.NewFlagSet("", flag.ExitOnError)
	Addr = serverFlags.String("a", cfg.Address, "Address for the server")
	Restore = serverFlags.Bool("r", true, "A bool flag for configuration upload")
	intervalStr := serverFlags.String("i", "300s", "Interval for saving the metrics into the file")
	StoreFile = serverFlags.String("f", cfgLogger.StoreFile, "file to store the metrics")
	StoreInterval, err = time.ParseDuration(*intervalStr)
	if err != nil {
		log.Fatal("Flag -i got an incorrect argument")
	}
	if len(os.Args) >= 2 {
		fmt.Println(Addr)
		serverFlags.Parse(os.Args[1:])
	}
	fmt.Println()
	fmt.Println(*Addr)
}

func main() {
	var port = handlers.GetPort(*Addr)
	ticker := time.NewTicker(StoreInterval)
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
