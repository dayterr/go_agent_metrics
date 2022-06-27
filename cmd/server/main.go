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
	go func() {
		for {
			<-ticker.C
			server2.WriteJSON(CfgLogger.StoreFile)
		}
	}()
	r := handlers.CreateRouter(CfgLogger.StoreFile, CfgLogger.Restore)
	http.ListenAndServe(CfgLogger.Address, r)
}
