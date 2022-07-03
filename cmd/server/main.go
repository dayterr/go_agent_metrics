package main

import (
	"fmt"
	"github.com/dayterr/go_agent_metrics/cmd/server/handlers"
	"github.com/dayterr/go_agent_metrics/internal/config/server"
	server2 "github.com/dayterr/go_agent_metrics/internal/server"
	"log"
	"net/http"
	"time"
)

func main() {
	CfgLogger, err := server.GetEnvLogger()
	if err != nil {
		log.Fatal(err)
	}
	ticker := time.NewTicker(CfgLogger.StoreInterval)
	h := handlers.NewAsyncHandler()
	go func() {
		for {
			<-ticker.C
			jsn, _ := h.MarshallMetrics()
			server2.WriteJSON(CfgLogger.StoreFile, jsn)
		}
	}()
	r := handlers.CreateRouterWithAsyncHandler(CfgLogger.StoreFile, CfgLogger.Restore, h)
	http.ListenAndServe(CfgLogger.Address, r)
}
