package agent

import (
	"flag"
	"github.com/caarlos0/env/v6"
	"log"
	"time"
)

const (
	defaultAddress         = "localhost:8080"
	defaultReportInterval = 10 * time.Second
	defaultPollInterval   = 2 * time.Second
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
	flag.DurationVar(&fs.ReportInterval, "r", defaultReportInterval, "Interval for sending the metrics to the server")
	flag.DurationVar(&fs.PollInterval, "p", defaultPollInterval, "Interval for polling the metrics")
	flag.StringVar(&fs.Address, "a", defaultAddress, "Address for the server")
	flag.Parse()

	err := env.Parse(&cfg)
	log.Println(cfg)
	if err != nil {
		return Config{}, err
	}
	if cfg.ReportInterval == defaultReportInterval && fs.ReportInterval != defaultReportInterval {
		cfg.ReportInterval = fs.ReportInterval
	}
	if cfg.PollInterval == defaultPollInterval && fs.PollInterval != defaultPollInterval {
		cfg.PollInterval = fs.PollInterval
	}
	if cfg.Address == defaultAddress && fs.Address != defaultAddress {
		cfg.Address = fs.Address
	}
	return cfg, nil
}
