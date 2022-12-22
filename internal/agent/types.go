package agent

import (
	"time"

	"github.com/dayterr/go_agent_metrics/internal/storage"
)

type Agent struct {
	Storage        storage.Storager
	Address        string
	ReportInterval time.Duration
	PollInterval   time.Duration
	Key            string
	CryptoKey string
}

func NewAgent(address string, repInt time.Duration, pInt time.Duration, key, cryptoKey string) Agent {
	s := storage.NewIMS()
	return Agent{
		Storage:        s,
		Address:        address,
		ReportInterval: repInt,
		PollInterval:   pInt,
		Key:            key,
		CryptoKey: cryptoKey,
	}
}
