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
	log.Println("first line server")
	CfgLogger, err := server.GetEnvLogger()
	if err != nil {
		log.Fatal(err)
	}
	ticker := time.NewTicker(CfgLogger.StoreInterval)
	h := handlers.NewAsyncHandler()
	go func(h handlers.AsyncHandler) {
		for {
			<-ticker.C
			jsn, err := h.MarshallMetrics()
			if err != nil {
				log.Fatal(err)
			}
			server2.WriteJSON(CfgLogger.StoreFile, jsn)
		}
	}(h)
	r := handlers.CreateRouterWithAsyncHandler(CfgLogger.StoreFile, CfgLogger.Restore, h)
	err = http.ListenAndServe(CfgLogger.Address, r)
	if err != nil {
		log.Fatal(err)
	}
}
