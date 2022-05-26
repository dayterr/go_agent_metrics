package main

import (
	"github.com/dayterr/go_agent_metrics/cmd/server/handlers"
	"net/http"
)

var port = ":8080"

func main() {
	r := handlers.CreateRouter()
	http.ListenAndServe(port, r)
}