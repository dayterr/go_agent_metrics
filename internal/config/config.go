package config

import (
	"flag"
	"log"
	"os"
	"strings"
	"time"

	"github.com/caarlos0/env/v6"
)

var (
	addr *string
	Restore *bool
	StoreFile *string
	StoreInterval time.Duration
)

func init() {
	var err error
	addr = flag.String("a", os.Getenv("ADDRESS"), "Address for the server")
	Restore = flag.Bool("r", true, "A bool flag for configuration upload")
	intervalStr := flag.String("i", os.Getenv("STORE_INTERVAL"), "Interval for saving the etrics into the file")
	StoreFile = flag.String("f", os.Getenv("STORE_FILE"), "file to store the metrics")
	StoreInterval, err = time.ParseDuration(*intervalStr)
	if err != nil {
		log.Fatal("Flag -i got an incorrect argument")
	}
	flag.Parse()
}

type Config struct {
	Address string `env:"ADDRESS" envDefault:"localhost:8080"`
	ReportInterval time.Duration `env:"REPORT_INTERVAL" envDefault:"10s"`
	PollInterval time.Duration `env:"POLL_INTERVAL" envDefault:"2s"`
}

type ConfigLogger struct {
	StoreInterval time.Duration `env:"STORE_INTERVAL" envDefault:"300s"`
	StoreFile string `env:"STORE_FILE" envDefault:"/tmp/devops-metrics-db.json"`
	Restore bool `env:"RESTORE" envDefault:"false"`
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
	port := ":" + strings.Split(*addr, ":")[1]
	return port
}