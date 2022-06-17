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

func main() {
	var Cfg = config.GetEnv()
	flag.StringVar(&Cfg.Address, "a", Cfg.Address, "Address for the server")
	flag.DurationVar(&Cfg.ReportInterval, "r", Cfg.ReportInterval, "Interval for sending the metrics to the server")
	flag.DurationVar(&Cfg.PollInterval, "p", Cfg.PollInterval, "Interval for polling the metrics")
	flag.CommandLine.Parse(os.Args[1:])
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt, syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)
	exitChan := make(chan int)
	ticker := time.NewTicker(Cfg.ReportInterval)
	tickerMetrics := time.NewTicker(Cfg.PollInterval)
	var am agent.Storage
	go func() {
		for {
			select {
			case <-tickerMetrics.C:
				am = agent.ReadMetrics()
			case <-ticker.C:
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
