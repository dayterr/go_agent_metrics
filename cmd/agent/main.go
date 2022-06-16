package main

import (
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/dayterr/go_agent_metrics/internal/agent"
	"github.com/dayterr/go_agent_metrics/internal/config"
)

var (
	Address *string
	ReportInterval time.Duration
	PollInterval time.Duration
)

func init() {
	var err error
	cfg := config.GetEnv()
	Address = flag.String("a", cfg.Address, "Address for the server")
	repIntervalStr := flag.String("r", "10s", "Interval for sending the metrics to the server")
	ReportInterval, err = time.ParseDuration(*repIntervalStr)
	if err != nil {
		log.Fatal("Flag -r for REPORT_INTERVAL got an incorrect argument")
	}
	pollIntervalStr := flag.String("p", "2s", "Interval for polling the metrics")
	PollInterval, err = time.ParseDuration(*pollIntervalStr)
	if err != nil {
		log.Fatal("Flag -p got an incorrect argument")
	}
	flag.Parse()
}

func main() {
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt, syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)
	exitChan := make(chan int)
	ticker := time.NewTicker(ReportInterval)
	tickerMetrics := time.NewTicker(PollInterval)
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
