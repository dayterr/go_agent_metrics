package flags

import (
	"flag"
	"github.com/dayterr/go_agent_metrics/internal/config"
	"log"
	"os"
	"time"
)

var (
	Address *string
	ReportInterval time.Duration
	PollInterval time.Duration
)

func init() {
	var err error
	cfg := config.GetEnv()
	agentFlags := flag.NewFlagSet("", flag.ExitOnError)
	Address = agentFlags.String("a", cfg.Address, "Address for the server")
	repIntervalStr := agentFlags.String("r", "10s", "Interval for sending the metrics to the server")
	ReportInterval, err = time.ParseDuration(*repIntervalStr)
	if err != nil {
		log.Fatal("Flag -r for REPORT_INTERVAL got an incorrect argument")
	}
	pollIntervalStr := agentFlags.String("p", "2s", "Interval for polling the metrics")
	PollInterval, err = time.ParseDuration(*pollIntervalStr)
	if err != nil {
		log.Fatal("Flag -p got an incorrect argument")
	}
	if len(os.Args) >= 2 {
		agentFlags.Parse(os.Args[1:])
	}
}