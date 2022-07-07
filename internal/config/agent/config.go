package agent

import (
	"flag"
	"github.com/caarlos0/env/v6"
	"log"
	"time"
)

const (
	defaultAddress        = "localhost:8080"
	defaultReportInterval = 10 * time.Second
	defaultPollInterval   = 2 * time.Second
	defaultKey            = ""
)

type Config struct {
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

func GetEnv() (Config, error) {
	var cfg Config
	fs := FlagStruct{}
	flag.DurationVar(&fs.ReportInterval, "r", defaultReportInterval, "Interval for sending the metric to the server")
	flag.DurationVar(&fs.PollInterval, "p", defaultPollInterval, "Interval for polling the metric")
	flag.StringVar(&fs.Address, "a", defaultAddress, "Address for the server")
	flag.StringVar(&fs.Key, "k", defaultKey, "Key for encrypting")
	flag.Parse()

	err := env.Parse(&cfg)
	if err != nil {
		//return Config{}, err
		log.Println("agent config error", err)
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
