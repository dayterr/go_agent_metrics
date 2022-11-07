package main

import (
	"net/http"
	_ "net/http/pprof"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"github.com/dayterr/go_agent_metrics/cmd/server/handlers"
	"github.com/dayterr/go_agent_metrics/internal/config/server"
	server2 "github.com/dayterr/go_agent_metrics/internal/server"
)

func main() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	Cfg, err := server.GetEnvServer()
	if err != nil {
		log.Fatal().Err(err).Msg("getting config error")
	}
	ticker := time.NewTicker(Cfg.StoreInterval)
	var h handlers.AsyncHandler
	if Cfg.DatabaseDSN == "" {
		h, err = handlers.NewAsyncHandler(Cfg.Key, Cfg.DatabaseDSN, false)
		if err != nil {
			log.Fatal().Err(err).Msg("creating handler error")
		}
	} else {
		h, err = handlers.NewAsyncHandler(Cfg.Key, Cfg.DatabaseDSN, true)
		if err != nil {
			log.Fatal().Err(err).Msg("creating handler error")
		}
	}
	if Cfg.DatabaseDSN == "" {
		go func(h handlers.AsyncHandler) {
			for {
				<-ticker.C
				jsn, err := h.MarshallMetrics()
				if err != nil {
					log.Fatal().Err(err).Msg("marshalling metrics error")
				}
				server2.WriteJSON(Cfg.StoreFile, jsn)
			}
		}(h)
	}
	restore := Cfg.Restore && Cfg.DatabaseDSN == ""
	r, err := handlers.CreateRouterWithAsyncHandler(Cfg.StoreFile, restore, h)
	if err != nil {
		log.Fatal().Err(err).Msg("creating router error")
	}
	err = http.ListenAndServe(Cfg.Address, r)
	if err != nil {
		log.Fatal().Err(err).Msg("error in server main")
	}
}
