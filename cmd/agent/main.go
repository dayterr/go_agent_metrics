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

var Cfg config.Config

func main() {
	Cfg = config.GetEnv()
	fmt.Println("before", Cfg)
	flag.DurationVar(&Cfg.ReportInterval, "r", Cfg.ReportInterval, "Interval for sending the metrics to the server")
	flag.DurationVar(&Cfg.PollInterval, "p", Cfg.PollInterval, "Interval for polling the metrics")
	Cfg = config.GetEnv()
	if Cfg.Address == "localhost:8080" {
		flag.StringVar(&Cfg.Address, "a", Cfg.Address, "Address for the server")
	}
	flag.CommandLine.Parse(os.Args[1:])
	fmt.Println("after", Cfg)
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
