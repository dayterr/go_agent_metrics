package agent

import (
	"flag"
	"fmt"
	"github.com/caarlos0/env/v6"
	"time"
)

const (
	DEFAULT_ADDRESS         = "localhost:8080"
	DEFAULT_REPORT_INTERVAL = 10 * time.Second
	DEFAULT_POLL_INTERVAL   = 2 * time.Second
)

type Config struct {
	Address        string        `env:"ADDRESS" envDefault:"localhost:8080"`
	ReportInterval time.Duration `env:"REPORT_INTERVAL" envDefault:"10s"`
	PollInterval   time.Duration `env:"POLL_INTERVAL" envDefault:"2s"`
}

type FlagStruct struct {
	Address        string
	ReportInterval time.Duration
	PollInterval   time.Duration
}

func GetEnv() (Config, error) {
	var cfg Config
	fs := FlagStruct{}
	flag.DurationVar(&cfg.ReportInterval, "r", cfg.ReportInterval, "Interval for sending the metrics to the server")
	flag.DurationVar(&cfg.PollInterval, "p", cfg.PollInterval, "Interval for polling the metrics")
	flag.StringVar(&cfg.Address, "a", cfg.Address, "Address for the server")
	flag.Parse()

	err := env.Parse(&cfg)
	if err != nil {
		return Config{}, err
	}
	if cfg.ReportInterval == DEFAULT_REPORT_INTERVAL {
		cfg.ReportInterval = fs.ReportInterval
	}
	if cfg.PollInterval == DEFAULT_POLL_INTERVAL {
		cfg.PollInterval = fs.PollInterval
	}
	if cfg.Address == DEFAULT_ADDRESS {
		cfg.Address = fs.Address
	}
	return cfg, nil
}
