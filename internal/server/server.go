package server

import (
	"bufio"
	"encoding/json"
	"github.com/dayterr/go_agent_metrics/cmd/server/handlers"
	"github.com/dayterr/go_agent_metrics/internal/agent"
	"github.com/dayterr/go_agent_metrics/internal/config"
	"io/ioutil"
	"log"
	"os"
)

func LoadMetricsFromJSON(cfg config.ConfigLogger, allMetrics agent.Storage) {
	if cfg.Restore {
		file, err := ioutil.ReadFile(cfg.StoreFile)
		if err != nil {
			log.Fatal(err)
		}
		err = json.Unmarshal(file, &allMetrics)
		if err != nil {
			log.Fatal(err)
		}
		agent.PostAll(allMetrics)
	}
}

func WriteJSON(path string) {
	file, err := os.OpenFile(path, os.O_CREATE | os.O_RDWR | os.O_TRUNC, 0777)
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