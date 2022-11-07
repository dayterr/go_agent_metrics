package agent

import (
	"flag"
	"log"
	"time"

	"github.com/caarlos0/env/v6"
)

const (
	defaultAddress        = "localhost:8080"
	defaultReportInterval = 10 * time.Second
	defaultPollInterval   = 2 * time.Second
	defaultKey            = ""
)

type ConfigAgent struct {
	Address        string        `env:"ADDRESS" envDefault:"localhost:8080"`
	ReportInterval time.Duration `env:"REPORT_INTERVAL" envDefault:"10s"`
	PollInterval   time.Duration `env:"POLL_INTERVAL" envDefault:"2s"`
	Key            string        `env:"KEY" envDefault:""`
}

type FlagStruct struct {
	Address        string
	ReportInterval time.Duration
	PollInterval   time.Duration
	Key            string
}

func GetEnvAgent() (ConfigAgent, error) {
	log.Println("first line config")
	var cfg ConfigAgent
	fs := FlagStruct{}
	flag.DurationVar(&fs.ReportInterval, "r", defaultReportInterval, "Interval for sending the metric to the server")
	flag.DurationVar(&fs.PollInterval, "p", defaultPollInterval, "Interval for polling the metric")
	flag.StringVar(&fs.Address, "a", defaultAddress, "Address for the server")
	flag.StringVar(&fs.Key, "k", defaultKey, "Key for encrypting")
	flag.Parse()

	err := env.Parse(&cfg)
	if err != nil {
		return ConfigAgent{}, err
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
	if cfg.Key == defaultKey && fs.Key != defaultKey {
		cfg.Key = fs.Key
	}
	log.Println("agent config", cfg)
	return cfg, nil
}
