package main

import (
	"net/http"

	"github.com/dayterr/go_agent_metrics/cmd/server/handlers"
)

var port = ":8080"

func main() {
	r := handlers.CreateRouter()
	http.ListenAndServe(port, r)
}
