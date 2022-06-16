package main

import (
	"github.com/dayterr/go_agent_metrics/cmd/agent/flags"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/dayterr/go_agent_metrics/internal/agent"
)

func main() {
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt, syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)
	exitChan := make(chan int)
	ticker := time.NewTicker(flags.ReportInterval)
	tickerMetrics := time.NewTicker(flags.PollInterval)
	var am agent.Storage
	go func() {
		for {
			select {
			case <-tickerMetrics.C:
				am = agent.ReadMetrics()
			case <-ticker.C:
				agent.PostAll(am)
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
