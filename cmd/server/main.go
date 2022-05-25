package main

import (
	"github.com/dayterr/go_agent_metrics/cmd/server/handlers"
	"net/http"
)

func main() {
	r := handlers.CreateRouter()
	http.ListenAndServe(":8080", r)
}