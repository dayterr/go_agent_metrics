package main

import (
	"flag"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/dayterr/go_agent_metrics/internal/agent"
	"github.com/dayterr/go_agent_metrics/internal/config"
)

var cfg = config.GetEnv()

func main() {
	flag.StringVar(&cfg.Address, "a", cfg.Address, "Address for the server")
	flag.DurationVar(&cfg.ReportInterval, "r", cfg.ReportInterval, "Interval for sending the metrics to the server")
	flag.DurationVar(&cfg.PollInterval, "p", cfg.PollInterval, "Interval for polling the metrics")
	flag.Parse()
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt, syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)
	exitChan := make(chan int)
	ticker := time.NewTicker(cfg.ReportInterval)
	tickerMetrics := time.NewTicker(cfg.PollInterval)
	var am agent.Storage
	go func() {
		for {
			select {
			case <-tickerMetrics.C:
				am = agent.ReadMetrics()
			case <-ticker.C:
				agent.PostAll(am, cfg.Address)
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
