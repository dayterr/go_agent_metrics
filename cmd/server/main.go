package main

import (
	"github.com/dayterr/go_agent_metrics/internal/config"
	"net/http"

	"github.com/dayterr/go_agent_metrics/cmd/server/handlers"
)

var port = config.GetPort()

func main() {
	r := handlers.CreateRouter()
	http.ListenAndServe(port, r)
}
