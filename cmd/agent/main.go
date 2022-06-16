package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/dayterr/go_agent_metrics/internal/agent"
	"github.com/dayterr/go_agent_metrics/internal/config"
)

var (
	Address *string
	ReportInterval *time.Duration
	PollInterval *time.Duration
)

func init() {
	cfg := config.GetEnv()
	Address = flag.String("a", cfg.Address, "Address for the server")
	ReportInterval = flag.Duration("r", cfg.ReportInterval, "Interval for sending the metrics to the server")
	PollInterval = flag.Duration("p", cfg.PollInterval, "Interval for polling the metrics")
	flag.Parse()
	fmt.Println(*Address, *ReportInterval, *PollInterval)
}

func main() {
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt, syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)
	exitChan := make(chan int)
	ticker := time.NewTicker(*ReportInterval)
	tickerMetrics := time.NewTicker(*PollInterval)
	var am agent.Storage
	go func() {
		for {
			select {
			case <-tickerMetrics.C:
				am = agent.ReadMetrics()
			case <-ticker.C:
				agent.PostAll(am, *Address)
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
