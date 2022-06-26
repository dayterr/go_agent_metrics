package server

import (
	"flag"
	"github.com/caarlos0/env/v6"
	"time"
)

var (
	DEFAULT_ADDRESS         = "localhost:8080"
	DEFAULT_STORE_INTERVAL  = time.Duration(300 * time.Second)
	DEFAULT_STORE_FILE      = "/tmp/devops-metrics-db.json"
	DEFAULT_RESTORE         = true
)

type ConfigLogger struct {
	Address        string        `env:"ADDRESS" envDefault:"localhost:8080"`
	StoreInterval time.Duration `env:"STORE_INTERVAL" envDefault:"300s"`
	StoreFile     string        `env:"STORE_FILE" envDefault:"/tmp/devops-metrics-db.json"`
	Restore       bool          `env:"RESTORE" envDefault:"true"`
}

func GetEnvLogger() (ConfigLogger, error) {
	var cfg ConfigLogger
	err := env.Parse(&cfg)
	if err != nil {
		return ConfigLogger{}, err
	}
	if cfg.Address == DEFAULT_ADDRESS {
		flag.StringVar(&cfg.Address, "a", cfg.Address, "Address for the server")
	}
	if cfg.Restore == DEFAULT_RESTORE {
		flag.BoolVar(&cfg.Restore, "r", cfg.Restore, "A bool flag for configuration upload")
	}
	if cfg.StoreInterval == DEFAULT_STORE_INTERVAL {
		flag.DurationVar(&cfg.StoreInterval, "i", cfg.StoreInterval, "Interval for saving the metrics into the file")
	}
	if cfg.StoreFile == DEFAULT_STORE_FILE {
		flag.StringVar(&cfg.StoreFile, "f", cfg.StoreFile, "file to store the metrics")
	}
	flag.Parse()
	return cfg, nil
}

