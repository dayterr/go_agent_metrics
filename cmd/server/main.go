package main

import (
	"net/http"
	_ "net/http/pprof"
	"time"

	"github.com/dayterr/go_agent_metrics/cmd/server/handlers"
	"github.com/dayterr/go_agent_metrics/internal/config/server"
	server2 "github.com/dayterr/go_agent_metrics/internal/server"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func main() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	Cfg, err := server.GetEnvServer()
	if err != nil {
		log.Print(err)
	}
	ticker := time.NewTicker(Cfg.StoreInterval)
	var h handlers.AsyncHandler
	if Cfg.DatabaseDSN == "" {
		h = handlers.NewAsyncHandler(Cfg.Key, Cfg.DatabaseDSN, false)
	} else {
		h = handlers.NewAsyncHandler(Cfg.Key, Cfg.DatabaseDSN, true)
	}
	if Cfg.DatabaseDSN == "" {
		go func(h handlers.AsyncHandler) {
			for {
				<-ticker.C
				jsn, err := h.MarshallMetrics()
				if err != nil {
					log.Print(err)
				}
				server2.WriteJSON(Cfg.StoreFile, jsn)
			}
		}(h)
	}
	restore := Cfg.Restore && Cfg.DatabaseDSN == ""
	r := handlers.CreateRouterWithAsyncHandler(Cfg.StoreFile, restore, h)
	err = http.ListenAndServe(Cfg.Address, r)
	if err != nil {
		log.Print("error in server main", err)
	}
}
