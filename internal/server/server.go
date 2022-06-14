package server

import (
	"bufio"
	"encoding/json"
	"fmt"
	"github.com/dayterr/go_agent_metrics/cmd/server/handlers"
	"github.com/dayterr/go_agent_metrics/internal/agent"
	"github.com/dayterr/go_agent_metrics/internal/config"
	"io/ioutil"
	"log"
	"os"
)

func LoadMetricsFromJSON(cfg config.ConfigLogger, allMetrics agent.Storage) {
	fmt.Println("I'm working")
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
	file, err := os.OpenFile(path, os.O_CREATE | os.O_RDWR , 0777)
	if err != nil {
		fmt.Println("create error", err)
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