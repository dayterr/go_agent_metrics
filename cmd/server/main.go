package main

import (
	"context"
	"fmt"
	mygrpc "github.com/dayterr/go_agent_metrics/cmd/server/grpc"
	"net"
	"net/http"
	_ "net/http/pprof"
	"os"
	"os/signal"
	"time"


	"github.com/dayterr/go_agent_metrics/internal/encryption"
	"google.golang.org/grpc"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"github.com/dayterr/go_agent_metrics/cmd/server/handlers"
	"github.com/dayterr/go_agent_metrics/internal/config/server"
	server2 "github.com/dayterr/go_agent_metrics/internal/server"
	pb "github.com/dayterr/go_agent_metrics/internal/grpc/proto"
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
	if Cfg.EnablegRPC {
		var mServer mygrpc.GRPCServer
		if Cfg.DatabaseDSN == "" {
			mServer, err = mygrpc.NewGRPCServer(Cfg.DatabaseDSN, false)
			if err != nil {
				log.Fatal().Err(err).Msg("creating server error")
			}
		} else {
			mServer, err = mygrpc.NewGRPCServer(Cfg.DatabaseDSN, true)
			if err != nil {
				log.Fatal().Err(err).Msg("creating server error")
			}
		}

		listen, err := net.Listen("tcp", Cfg.Address)
		if err != nil {
			log.Fatal().Err(err).Msg("starting listening error")
		}

		s := grpc.NewServer()
		pb.RegisterMetricsServiceServer(s, mServer)
	} else {
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
		r, err := handlers.CreateRouterWithAsyncHandler(Cfg.StoreFile, restore, h, e, []byte(Cfg.Salt), Cfg.TrustedSubnet)
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

		err = srv.ListenAndServe()
		if err != http.ErrServerClosed {
			log.Fatal().Err(err).Msg("HTTP server ListenAndServe error")
		}
		<-idleConnsClosed
		log.Info().Msg("Server Shutdown gracefully")
	}
}
