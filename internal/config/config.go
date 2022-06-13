package config

import (
	"log"
	"strings"
	"time"

	"github.com/caarlos0/env/v6"
)

type Config struct {
	Address string `env:"ADDRESS" envDefault:"localhost:8080"`
	ReportInterval time.Duration `env:"REPORT_INTERVAL" envDefault:"10s"`
	PollInterval time.Duration `env:"POLL_INTERVAL" envDefault:"2s"`
}

type ConfigLogger struct {
	StoreInterval time.Duration `env:"STORE_INTERVAL" envDefault:"10s"`
	StoreFile string `env:"STORE_FILE" envDefault:"/tmp/devops-metrics-db.json"`
	Restore bool `env:"RESTORE" envDefault:"true"`
}

func GetEnv() Config {
	var cfg Config
	err := env.Parse(&cfg)
	if err != nil {
		log.Fatal(err)
	}
	return cfg
}

func GetEnvLogger() ConfigLogger {
	var cfg ConfigLogger
	err := env.Parse(&cfg)
	if err != nil {
		log.Fatal(err)
	}
	return cfg
}

func GetPort() string {
	cfg := GetEnv()
	port := ":" + strings.Split(cfg.Address, ":")[1]
	return port
}