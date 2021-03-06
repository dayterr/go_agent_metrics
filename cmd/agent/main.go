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
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt, syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)
	exitChan := make(chan int)
	tickerCollectMetrics := time.NewTicker(Cfg.PollInterval)
	tickerReportMetrics := time.NewTicker(Cfg.ReportInterval)
	var am = storage.New()
	go func() {
		for {
			select {
			case <-tickerCollectMetrics.C:
				am = agent.ReadMetrics()
			case <-tickerReportMetrics.C:
				agent.PostAll(am, Cfg.Address)
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
