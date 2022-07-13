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
	Cfg, err := server.GetEnvServer()
	if err != nil {
		log.Fatal(err)
	}
	ticker := time.NewTicker(Cfg.StoreInterval)
	var h handlers.AsyncHandler
	if Cfg.DatabaseDSN == "" {
		h = handlers.NewAsyncHandler(Cfg.Key, Cfg.DatabaseDSN,false)
	} else {
		h = handlers.NewAsyncHandler(Cfg.Key, Cfg.DatabaseDSN, true)
	}
	go func(h handlers.AsyncHandler) {
		for {
			<-ticker.C
			jsn, err := h.MarshallMetrics()
			if err != nil {
				log.Fatal(err)
			}
			server2.WriteJSON(Cfg.StoreFile, jsn)
		}
	}(h)
	log.Println("Cfg.StoreFile", Cfg.StoreFile)
	r := handlers.CreateRouterWithAsyncHandler(Cfg.StoreFile, Cfg.Restore, h)
	err = http.ListenAndServe(Cfg.Address, r)
	if err != nil {
		log.Fatal("error in server main", err)
	}
}
