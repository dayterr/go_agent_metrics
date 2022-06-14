package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"github.com/dayterr/go_agent_metrics/internal/agent"
	"github.com/dayterr/go_agent_metrics/internal/server"
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
	time.AfterFunc(time.Second * 2, func() {
		file, err := os.OpenFile(cfg.StoreFile, os.O_CREATE | os.O_RDWR | os.O_TRUNC, 0777)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("created file", file.Name())
		defer file.Close()
		jsn, err := handlers.MarshallMetrics()
		if err != nil {
			log.Fatal(err)
		}
		w := bufio.NewWriter(file)
		w.Write(jsn)
		w.Flush()
	})
	r := handlers.CreateRouter()
	http.ListenAndServe(port, r)
}
