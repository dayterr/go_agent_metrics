package server

import (
	"flag"
	"github.com/caarlos0/env/v6"
	"log"
	"time"
)

var (
	defaultAddress       = "localhost:8080"
	defaultStoreInterval = 300 * time.Second
	defaultStoreFile     = "/tmp/devops-metric-db.json"
	defaultRestore       = true
	defaultKey           = ""
)

type ConfigLogger struct {
	Address       string        `env:"ADDRESS" envDefault:"localhost:8080"`
	StoreInterval time.Duration `env:"STORE_INTERVAL" envDefault:"300s"`
	StoreFile     string        `env:"STORE_FILE" envDefault:"/tmp/devops-metric-db.json"`
	Restore       bool          `env:"RESTORE" envDefault:"true"`
	Key           string        `env:"KEY" envDefault:""`
	DatabaseDSN   string        `env:"DATABASE_DSN" envDefault:""`
}

type FlagStruct struct {
	Address       string
	StoreInterval time.Duration
	StoreFile     string
	Restore       bool
	Key           string
	DatabaseDSN   string
}

func GetEnvLogger() (ConfigLogger, error) {
	log.Println("first line in GetEnvLogger")
	cfg := ConfigLogger{}
	fs := FlagStruct{}
	flag.StringVar(&fs.Address, "a", defaultAddress, "Address for the server")
	flag.BoolVar(&fs.Restore, "r", defaultRestore, "A bool flag for configuration upload")
	flag.DurationVar(&fs.StoreInterval, "i", defaultStoreInterval, "Interval for saving the metric into the file")
	flag.StringVar(&fs.StoreFile, "f", defaultStoreFile, "file to store the metric")
	flag.StringVar(&fs.Key, "k", defaultKey, "Key for encrypting")
	flag.StringVar(&fs.DatabaseDSN, "d", "", "Database DSN")
	flag.Parse()

	err := env.Parse(&cfg)
	if err != nil {
		return ConfigLogger{}, err
	}

	if cfg.Address == defaultAddress && fs.Address != defaultAddress {
		cfg.Address = fs.Address
	}
	/*if cfg.Restore == defaultRestore && fs.Restore != defaultRestore {
		cfg.Restore = fs.Restore
	}*/
	if cfg.StoreInterval == defaultStoreInterval && fs.StoreInterval != defaultStoreInterval {
		cfg.StoreInterval = fs.StoreInterval
	}
	if cfg.StoreFile == defaultStoreFile && fs.StoreFile != defaultStoreFile && cfg.DatabaseDSN == "" {
		cfg.StoreFile = fs.StoreFile
	}
	if cfg.Key == defaultKey && fs.Key != defaultKey {
		cfg.Key = fs.Key
	}
	if cfg.DatabaseDSN == "" && fs.DatabaseDSN != "" {
		cfg.DatabaseDSN = fs.DatabaseDSN
	}
	return cfg, nil
}
