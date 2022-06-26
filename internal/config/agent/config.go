package agent

import (
	"flag"
	"github.com/caarlos0/env/v6"
	"time"
)

var (
	DEFAULT_ADDRESS         = "localhost:8080"
	DEFAULT_REPORT_INTERVAL = time.Duration(10 * time.Second)
	DEFAULT_POLL_INTERVAL   = time.Duration(2 * time.Second)
)

type Config struct {
	Address        string        `env:"ADDRESS" envDefault:"localhost:8080"`
	ReportInterval time.Duration `env:"REPORT_INTERVAL" envDefault:"10s"`
	PollInterval   time.Duration `env:"POLL_INTERVAL" envDefault:"2s"`
}

func GetEnv() (Config, error) {
	var cfg Config
	err := env.Parse(&cfg)
	if err != nil {
		return Config{}, err
	}
	if cfg.ReportInterval == DEFAULT_REPORT_INTERVAL {
		flag.DurationVar(&cfg.ReportInterval, "r", cfg.ReportInterval, "Interval for sending the metrics to the server")
	}
	if cfg.PollInterval == DEFAULT_POLL_INTERVAL {
		flag.DurationVar(&cfg.PollInterval, "p", cfg.PollInterval, "Interval for polling the metrics")
	}
	if cfg.Address == DEFAULT_ADDRESS {
		flag.StringVar(&cfg.Address, "a", cfg.Address, "Address for the server")
	}
	flag.Parse()
	return cfg, nil
}
