package main

import (
	"log"

	"github.com/dayterr/go_agent_metrics/internal/agent"
	agent2 "github.com/dayterr/go_agent_metrics/internal/config/agent"
)

func main() {
	Cfg, err := agent2.GetEnvAgent()
	if err != nil {
		log.Fatal(err)
	}
	agentInstance := agent.NewAgent(Cfg.Address, Cfg.ReportInterval, Cfg.PollInterval, Cfg.Key)
	agentInstance.Run()
}
