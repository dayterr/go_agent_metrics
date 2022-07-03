package main

import (
	"github.com/dayterr/go_agent_metrics/internal/storage"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/dayterr/go_agent_metrics/internal/agent"
	agent2 "github.com/dayterr/go_agent_metrics/internal/config/agent"
)

func main() {
	Cfg, err := agent2.GetEnv()
	if err != nil {
		log.Fatal(err)
	}
	agentInstance := agent.NewAgent(Cfg.Address, Cfg.ReportInterval, Cfg.PollInterval)
	agentInstance.Run()
}
