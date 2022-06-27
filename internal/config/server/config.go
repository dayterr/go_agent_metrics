package server

import (
	"flag"
	"github.com/caarlos0/env/v6"
	"time"
)

var (
	defaultAddress        = "localhost:8080"
	defaultStoreInterval  = 300 * time.Second
	DEFAULT_STORE_FILE      = "/tmp/devops-metrics-db.json"
	DEFAULT_RESTORE         = true
)

type ConfigLogger struct {
	Address       string        `env:"ADDRESS" envDefault:"localhost:8080"`
	StoreInterval time.Duration `env:"STORE_INTERVAL" envDefault:"300s"`
	StoreFile     string        `env:"STORE_FILE" envDefault:"/tmp/devops-metrics-db.json"`
	Restore       bool          `env:"RESTORE" envDefault:"true"`
}

type FlagStruct struct {
	Address       string
	StoreInterval time.Duration
	StoreFile     string
	Restore       bool
}

func GetEnvLogger() (ConfigLogger, error) {
	cfg := ConfigLogger{}
	fs := FlagStruct{}
	flag.StringVar(&fs.Address, "a", defaultAddress, "Address for the server")
	flag.BoolVar(&fs.Restore, "r", DEFAULT_RESTORE, "A bool flag for configuration upload")
	flag.DurationVar(&fs.StoreInterval, "i", defaultStoreInterval, "Interval for saving the metrics into the file")
	flag.StringVar(&fs.StoreFile, "f", DEFAULT_STORE_FILE, "file to store the metrics")
	flag.Parse()

	err := env.Parse(&cfg)
	if err != nil {
		return ConfigLogger{}, err
	}

	if cfg.Address == defaultAddress{
		cfg.Address = fs.Address
	}
	if cfg.Restore == DEFAULT_RESTORE {
		cfg.Restore = fs.Restore
	}
	if cfg.StoreInterval == defaultStoreInterval {
		cfg.StoreInterval = fs.StoreInterval
	}
	if cfg.StoreFile == DEFAULT_STORE_FILE {
		cfg.StoreFile = fs.StoreFile
	}
	return cfg, nil
}

