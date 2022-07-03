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
	agent := agent.Agent{Storage: storage.InMemoryStorage{}}
	Cfg, err := agent2.GetEnv()
	if err != nil {
		log.Fatal(err)
	}
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt, syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)
	exitChan := make(chan int)
	tickerCollectMetrics := time.NewTicker(Cfg.PollInterval)
	tickerReportMetrics := time.NewTicker(Cfg.ReportInterval)
	go func() {
		for {
			select {
			case <-tickerCollectMetrics.C:
				agent.Storage.ReadMetrics()
			case <-tickerReportMetrics.C:
				agent.PostAll(Cfg.Address)
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
