package server

import (
	"flag"
	"time"

	"github.com/caarlos0/env/v6"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

var (
	defaultAddress       = "localhost:8080"
	defaultStoreInterval = 300 * time.Second
	defaultStoreFile     = "/tmp/devops-metric-db.json"
	defaultRestore       = true
	defaultKey           = ""
)

type ConfigServer struct {
	Address       string        `env:"ADDRESS" envDefault:"localhost:8080"`
	StoreInterval time.Duration `env:"STORE_INTERVAL" envDefault:"300s"`
	StoreFile     string        `env:"STORE_FILE" envDefault:"/tmp/devops-metric-db.json"`
	Restore       bool          `env:"RESTORE" envDefault:"true"`
	Key           string        `env:"KEY" envDefault:""`
	DatabaseDSN   string        `env:"DATABASE_DSN" envDefault:""`
	Salt string `env:"SALT" envDefault:""`
	CryptoKey string `env:"CRYPTO_KEY" envDefault:""`
}

type FlagStruct struct {
	Address       string
	StoreInterval time.Duration
	StoreFile     string
	Restore       bool
	Key           string
	DatabaseDSN   string
	Salt string
	CryptoKey string
}

func GetEnvServer() (ConfigServer, error) {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	log.Print("first line in GetEnvServer")
	cfg := ConfigServer{}
	fs := FlagStruct{}
	flag.StringVar(&fs.Address, "a", defaultAddress, "Address for the server")
	flag.BoolVar(&fs.Restore, "r", defaultRestore, "A bool flag for configuration upload")
	flag.DurationVar(&fs.StoreInterval, "i", defaultStoreInterval, "Interval for saving the metric into the file")
	flag.StringVar(&fs.StoreFile, "f", defaultStoreFile, "file to store the metric")
	flag.StringVar(&fs.Key, "k", defaultKey, "Key for encrypting")
	flag.StringVar(&fs.DatabaseDSN, "d", "", "Database DSN")
	flag.StringVar(&fs.Salt, "salt", "", "salt for crypto key")
	flag.StringVar(&fs.CryptoKey, "cryptokey", "", "crypto key")
	flag.Parse()

	err := env.Parse(&cfg)
	if err != nil {
		return ConfigServer{}, err
	}

	if cfg.Address == defaultAddress && fs.Address != defaultAddress {
		cfg.Address = fs.Address
	}
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
	log.Print("server config", cfg)
	return cfg, nil
}
