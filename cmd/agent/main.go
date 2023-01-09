package main

import (
	"context"
	"fmt"
	pb "github.com/dayterr/go_agent_metrics/internal/grpc/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/dayterr/go_agent_metrics/internal/agent"
	agent2 "github.com/dayterr/go_agent_metrics/internal/config/agent"
)

var (
	buildVersion = "N/A"
	buildDate    = "N/A"
	buildCommit  = "N/A"
)

func main() {
	fmt.Printf("Build version: %v\nBuild date: %v\nBuild commit: %v\n", buildVersion, buildDate, buildCommit)
	Cfg, err := agent2.GetEnvAgent()
	if err != nil {
		log.Fatal(err)
	}
	agentInstance := agent.NewAgent(Cfg.Address, Cfg.ReportInterval, Cfg.PollInterval, Cfg.Key, Cfg.CryptoKey)
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)
	defer cancel()

	if Cfg.EnablegRPC {
		conn, err := grpc.Dial(Cfg.Address, grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			log.Fatal(err)
		}
		defer conn.Close()

		agentInstance.GRPCClient = pb.NewMetricsServiceClient(conn)

		agentInstance.RungRPC(ctx)
	} else {
		agentInstance.Run(ctx)
	}

}
