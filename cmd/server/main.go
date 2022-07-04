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
	log.Println(CfgLogger)
	if err != nil {
		log.Fatal(err)
	}
	ticker := time.NewTicker(CfgLogger.StoreInterval)
	h := handlers.NewAsyncHandler()
	fmt.Println(h)
	go func() {
		for {
			<-ticker.C
			jsn, err := h.MarshallMetrics()
			if err != nil {
				log.Fatal(err)
			}
			server2.WriteJSON(CfgLogger.StoreFile, jsn)
		}
	}()
	r := handlers.CreateRouterWithAsyncHandler(CfgLogger.StoreFile, CfgLogger.Restore, h)

	err = http.ListenAndServe(CfgLogger.Address, r)
	if err != nil {
		log.Fatal(err)
	}
}
