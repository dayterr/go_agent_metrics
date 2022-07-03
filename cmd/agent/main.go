package main

import (
	"fmt"
	"github.com/dayterr/go_agent_metrics/internal/agent"
	agent2 "github.com/dayterr/go_agent_metrics/internal/config/agent"
	"log"
)

func main() {
	Cfg, err := agent2.GetEnv()
	fmt.Println(Cfg)
	if err != nil {
		log.Fatal(err)
	}
	agentInstance := agent.NewAgent(Cfg.Address, Cfg.ReportInterval, Cfg.PollInterval)
	agentInstance.Run()
}
