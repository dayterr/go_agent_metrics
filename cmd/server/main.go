package main

import (
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
	r, h := handlers.CreateRouterWithAsyncHandler(CfgLogger.StoreFile, CfgLogger.Restore)
	go func() {
		for {
			<-ticker.C
			jsn, _ := h.MarshallMetrics()
			server2.WriteJSON(CfgLogger.StoreFile, jsn)
		}
	}()
	http.ListenAndServe(CfgLogger.Address, r)
}
