package main

import (
	"fmt"
	"github.com/dayterr/go_agent_metrics/internal/encryption"
	"net/http"
	_ "net/http/pprof"
	"os"
	"os/signal"
	"time"
	"context"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"github.com/dayterr/go_agent_metrics/cmd/server/handlers"
	"github.com/dayterr/go_agent_metrics/internal/config/server"
	server2 "github.com/dayterr/go_agent_metrics/internal/server"
)

var (
	buildVersion = "N/A"
	buildDate    = "N/A"
	buildCommit  = "N/A"
)

func main() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	fmt.Printf("Build version: %v\nBuild date: %v\nBuild commit: %v\n", buildVersion, buildDate, buildCommit)
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
	e := encryption.NewEncryptor(Cfg.CryptoKey)
	r, err := handlers.CreateRouterWithAsyncHandler(Cfg.StoreFile, restore, h, e, []byte(Cfg.Salt))
	if err != nil {
		log.Fatal().Err(err).Msg("creating router error")
	}

	srv := http.Server{Addr: Cfg.Address, Handler: r}
	idleConnsClosed := make(chan struct{})
	sigint := make(chan os.Signal, 1)
	signal.Notify(sigint, os.Interrupt)

	go func() {
		<-sigint
		if err := srv.Shutdown(context.Background()); err != nil {
			log.Info().Err(err).Msg("HTTP server Shutdown")
		}
		close(idleConnsClosed)
	}()

	if err := srv.ListenAndServe(); err != http.ErrServerClosed {
		log.Fatal().Err(err).Msg("HTTP server ListenAndServe error")
	}
	<-idleConnsClosed
	log.Info().Msg("Server Shutdown gracefully")
}
