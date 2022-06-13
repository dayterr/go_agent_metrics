package main

import (
	"bufio"
	"encoding/json"
	"github.com/dayterr/go_agent_metrics/internal/agent"
	"io/ioutil"
	"log"
	"net/http"
	"os"
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
	l, _ := os.Getwd()
	go func() {
		if cfg.Restore {
			file, err := ioutil.ReadFile(l + cfg.StoreFile)
			if err != nil {
				log.Fatal(err)
			}
			err = json.Unmarshal(file, &allMetrics)
			if err != nil {
				log.Fatal(err)
			}
			agent.PostAll(allMetrics)
		}
	}()
	go func() {
		for {
			<- ticker.C

			file, err := os.OpenFile(l + cfg.StoreFile, os.O_CREATE | os.O_RDWR | os.O_SYNC, 0777)
			if err != nil {
				log.Fatal(err)
			}
			defer file.Close()
			jsn, err := handlers.MarshallMetrics()
			if err != nil {
				log.Fatal(err)
			}
			w := bufio.NewWriter(file)
			w.Write(jsn)
			w.Flush()
		}
	}()
	r := handlers.CreateRouter()
	http.ListenAndServe(port, r)
}
