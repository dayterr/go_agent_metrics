package main

import (
	"fmt"
	"github.com/dayterr/go_agent_metrics/cmd/server/handlers"
	"github.com/dayterr/go_agent_metrics/internal/agent"
	"github.com/dayterr/go_agent_metrics/internal/server"
	"net/http"
	"time"
)

var metrics = make(map[string]agent.Gauge)
var counters = make(map[string]agent.Counter)

var port = handlers.GetPort()

func main() {
	ticker := time.NewTicker(handlers.StoreInterval)
	go func() {
		for {
			select {
			case <- ticker.C:
				server.WriteJSON(*handlers.StoreFile)
			}
		}
	}()
	fmt.Println(port)
	r := handlers.CreateRouter()
	http.ListenAndServe(port, r)
}
