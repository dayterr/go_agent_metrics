package server

import (
	"flag"
	"io/ioutil"
	"os"
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
	File           string        `env:"CONFIG" envDefault:""`
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
	File string
}

type FileStruct struct {
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
	flag.StringVar(&fs.File, "c", "", "config file")
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
	if cfg.Salt == "" && fs.Salt != "" {
		cfg.Salt = fs.Salt
	}
	if cfg.CryptoKey == "" && fs.CryptoKey != "" {
		cfg.CryptoKey = fs.CryptoKey
	}
	if cfg.File == "" && fs.File != "" {
		cfg.File = fs.File
	}

	var fileCfg FileStruct
	if cfg.File != "" {
		fileCfg, err = readConfigFile(cfg.File)
		if err != nil {
			log.Info().Msg("config file error")
		}
	}

	if cfg.Address == defaultAddress && fileCfg.Address != "" {
		cfg.Address = fileCfg.Address
	}
	if cfg.StoreInterval == defaultStoreInterval && fileCfg.StoreInterval != 0 {
		cfg.StoreInterval = fileCfg.StoreInterval
	}
	if cfg.Key == defaultKey && fileCfg.Key != "" {
		cfg.Key = fileCfg.Key
	}
	if cfg.DatabaseDSN == "" && fileCfg.DatabaseDSN != "" {
		cfg.DatabaseDSN = fileCfg.DatabaseDSN
	}
	if cfg.Salt == "" && fileCfg.Salt != "" {
		cfg.Salt = fileCfg.Salt
	}
	if cfg.CryptoKey == "" && fileCfg.CryptoKey != "" {
		cfg.CryptoKey = fileCfg.CryptoKey
	}

	log.Print("server config", cfg)
	return cfg, nil
}

func readConfigFile(filepath string) (FileStruct, error) {
	jsonFile, err := os.Open(filepath)

	if err != nil {
		return FileStruct{}, err
	}

	defer jsonFile.Close()

	byteValue, err := ioutil.ReadAll(jsonFile)

	if err != nil {
		return FileStruct{}, err
	}

	var fs FileStruct

	err = json.Unmarshal(byteValue, &fs)

	if err != nil {
		return FileStruct{}, err
	}

	return fs, nil
}