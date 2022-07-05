package agent

import (
	"github.com/dayterr/go_agent_metrics/internal/storage"
	"time"
)

type Agent struct {
	Storage        storage.Storager
	Address        string
	ReportInterval time.Duration
	PollInterval   time.Duration
	Key            string
}

func NewAgent(address string, repInt time.Duration, pInt time.Duration, key string) Agent {
	s := storage.NewIMS()
	return Agent{
		Storage:        s,
		Address:        address,
		ReportInterval: repInt,
		PollInterval:   pInt,
		Key:            key,
	}
}

func (a Agent) Run() {

}
