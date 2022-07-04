package server

import (
	"flag"
	"github.com/caarlos0/env/v6"
	"log"
	"os"
	"time"
)

var (
	defaultAddress        = "localhost:8080"
	defaultStoreInterval  = 300 * time.Second
	defaultStoreFile      = "/tmp/devops-metrics-db.json"
	defaultRestore        = true
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
	log.Println("first line in GetEnvLogger")
	cfg := ConfigLogger{}
	fs := FlagStruct{}
	flag.StringVar(&fs.Address, "a", defaultAddress, "Address for the server")
	flag.BoolVar(&fs.Restore, "r", defaultRestore, "A bool flag for configuration upload")
	flag.DurationVar(&fs.StoreInterval, "i", defaultStoreInterval, "Interval for saving the metrics into the file")
	flag.StringVar(&fs.StoreFile, "f", defaultStoreFile, "file to store the metrics")
	flag.Parse()

	err := env.Parse(&cfg)
	log.Println("err server config", err)
	if err != nil {
		//return ConfigLogger{}, err
		log.Println("server config error", err)
	}
	log.Println("server config", cfg)

	if cfg.Address == defaultAddress && fs.Address != defaultAddress {
		cfg.Address = fs.Address
	}
	if cfg.Restore == defaultRestore && fs.Restore != defaultRestore && os.Getenv("RESTORE") == "" {
		cfg.Restore = fs.Restore
	}
	if cfg.StoreInterval == defaultStoreInterval && fs.StoreInterval != defaultStoreInterval {
		cfg.StoreInterval = fs.StoreInterval
	}
	if cfg.StoreFile == defaultStoreFile && fs.StoreFile != defaultStoreFile {
		cfg.StoreFile = fs.StoreFile
	}
	log.Println("read flags")
	log.Println("final config", cfg)
	return cfg, nil
}

