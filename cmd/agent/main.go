package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/dayterr/go_agent_metrics/internal/agent"
	"github.com/dayterr/go_agent_metrics/internal/config"
)

func main() {
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt, syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)
	exitChan := make(chan int)
	conf := config.GetEnv()
	ticker := time.NewTicker(conf.ReportInterval)
	tickerMetrics := time.NewTicker(conf.PollInterval)
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
