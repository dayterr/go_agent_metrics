package agent

import (
	"encoding/json"
	"flag"
	"io/ioutil"
	"log"
	"os"
	"time"

	"github.com/caarlos0/env/v6"
)

const (
	defaultAddress        = "localhost:8080"
	defaultReportInterval = 10 * time.Second
	defaultPollInterval   = 2 * time.Second
	defaultKey            = ""
)

type ConfigAgent struct {
	Address        string        `env:"ADDRESS" envDefault:"localhost:8080"`
	ReportInterval time.Duration `env:"REPORT_INTERVAL" envDefault:"10s"`
	PollInterval   time.Duration `env:"POLL_INTERVAL" envDefault:"2s"`
	Key            string        `env:"KEY" envDefault:""`
	Salt           string
	CryptoKey      string
	File           string `env:"CONFIG" envDefault:""`
}

type FlagStruct struct {
	Address        string
	ReportInterval time.Duration
	PollInterval   time.Duration
	Key            string
	Salt           string
	CryptoKey      string
	File           string
}

type FileStruct struct {
	Address        string
	ReportInterval time.Duration
	PollInterval   time.Duration
	Key            string
	Salt           string
	CryptoKey      string
}

func GetEnvAgent() (ConfigAgent, error) {
	log.Println("first line config")
	var cfg ConfigAgent
	fs := FlagStruct{}
	flag.DurationVar(&fs.ReportInterval, "r", defaultReportInterval, "Interval for sending the metric to the server")
	flag.DurationVar(&fs.PollInterval, "p", defaultPollInterval, "Interval for polling the metric")
	flag.StringVar(&fs.Address, "a", defaultAddress, "Address for the server")
	flag.StringVar(&fs.Key, "k", defaultKey, "Key for encrypting")
	flag.StringVar(&fs.Salt, "salt", "", "salt for crypto key")
	flag.StringVar(&fs.CryptoKey, "cryptokey", "", "crypto key")
	flag.StringVar(&fs.File, "c", "", "config file")
	flag.Parse()

	err := env.Parse(&cfg)
	if err != nil {
		return ConfigAgent{}, err
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
			log.Println("config file error")
		}
	}

	if cfg.ReportInterval == defaultReportInterval && fileCfg.ReportInterval != 0 {
		cfg.ReportInterval = fileCfg.ReportInterval
	}
	if cfg.PollInterval == defaultPollInterval && fileCfg.PollInterval != 0 {
		cfg.PollInterval = fileCfg.PollInterval
	}
	if cfg.Address == defaultAddress && fileCfg.Address != "" {
		cfg.Address = fileCfg.Address
	}
	if cfg.Key == defaultKey && fileCfg.Key != "" {
		cfg.Key = fileCfg.Key
	}
	if cfg.Salt == "" && fileCfg.Salt != "" {
		cfg.Salt = fileCfg.Salt
	}
	if cfg.CryptoKey == "" && fileCfg.CryptoKey != "" {
		cfg.CryptoKey = fileCfg.CryptoKey
	}

	log.Println("agent config", cfg)
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
