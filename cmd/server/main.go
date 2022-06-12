package main

import (
	"bufio"
	"github.com/dayterr/go_agent_metrics/internal/config"
	"log"
	"os"
	"time"
	"net/http"


	"github.com/dayterr/go_agent_metrics/cmd/server/handlers"
)

var port = config.GetPort()

func main() {
	cfg := config.GetEnvLogger()
	ticker := time.NewTicker(cfg.StoreInterval)
	go func() {
		for {
			<- ticker.C
			file, err := os.OpenFile(cfg.StoreFile, os.O_RDWR | os.O_CREATE | os.O_APPEND | os.O_SYNC, 0777)
			if err != nil {
				log.Fatal(err)
			}
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
