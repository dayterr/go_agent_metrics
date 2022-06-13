package main

import (
	"bufio"
	"encoding/json"
	"fmt"
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

type AllMetrics struct {
	GaugeField map[string]agent.Gauge
	CounterField map[string]agent.Counter
}

var allMetrics AllMetrics = AllMetrics{
	metrics,
	counters,
}

var port = config.GetPort()

func main() {
	cfg := config.GetEnvLogger()
	ticker := time.NewTicker(cfg.StoreInterval)
	if cfg.Restore {
		file, err := ioutil.ReadFile(cfg.StoreFile)
		if err != nil {
			log.Fatal(err)
		}
		err = json.Unmarshal(file, &allMetrics)
		if err != nil {
			log.Fatal(err)
		}
		agent.PostAll()

	}
	go func() {
		for {
			<- ticker.C
			file, err := os.OpenFile(cfg.StoreFile, os.O_CREATE | os.O_APPEND | os.O_RDWR | os.O_SYNC, 0777)
			if err != nil {
				log.Fatal(err)
			}
			fmt.Println(cfg.StoreFile)
			defer file.Close()
			jsn, err := handlers.MarshallMetrics()
			if err != nil {
				log.Fatal(err)
			}
			jsn2, err := handlers.MarshallCounters()
			_ = jsn2
			if err != nil {
				log.Fatal(err)
			}
			w := bufio.NewWriter(file)
			w.Write(jsn)
		}
	}()
	r := handlers.CreateRouter()
	http.ListenAndServe(port, r)
}
