package main

import (
	"github.com/dayterr/go_agent_metrics/internal/agent"
	agent2 "github.com/dayterr/go_agent_metrics/internal/config/agent"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	Cfg, err := agent2.GetEnv()
	if err != nil {
		log.Fatal(err)
	}
	agentInstance := agent.NewAgent(Cfg.Address, Cfg.ReportInterval, Cfg.PollInterval, Cfg.Key)
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt, syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)
	exitChan := make(chan int)
	tickerCollectMetrics := time.NewTicker(agentInstance.PollInterval)
	tickerReportMetrics := time.NewTicker(agentInstance.ReportInterval)
	go func() {
		for {
			select {
			case <-tickerCollectMetrics.C:
				agentInstance.Storage.ReadMetrics()
			case <-tickerReportMetrics.C:
				agentInstance.PostAll()
			case s := <-signalChan:
				switch s {
				case syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT:
					exitChan <- 0
				}
			}
		}
	}()
	exitCode := <-exitChan
	os.Exit(exitCode)
}
