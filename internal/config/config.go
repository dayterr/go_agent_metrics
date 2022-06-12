package config

import (
	"fmt"
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

func GetEnv() Config {
	var cfg Config
	err := env.Parse(&cfg)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(cfg)
	return cfg
}

func GetPort() string {
	cfg := GetEnv()
	port := ":" + strings.Split(cfg.Address, ":")[1]
	fmt.Println(port)
	return port
}