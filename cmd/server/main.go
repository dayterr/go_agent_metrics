package main

import (
	"encoding/json"
	"fmt"
	"github.com/dayterr/go_agent_metrics/internal/agent"
	"github.com/dayterr/go_agent_metrics/internal/server"
	"io/ioutil"
	"log"
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
	time.AfterFunc(time.Second, func() {
		if cfg.Restore {
			file, err := ioutil.ReadFile(cfg.StoreFile)
			if err != nil {
				fmt.Println("trying to read the file")
				log.Fatal(err)
			}
			err = json.Unmarshal(file, &allMetrics)
			if err != nil {
				log.Fatal(err)
			}
			agent.PostAll(allMetrics)
		}
	})
	go func() {
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
