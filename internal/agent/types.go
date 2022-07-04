package agent

import (
	"github.com/dayterr/go_agent_metrics/internal/storage"
	"log"
	"time"
)

type Agent struct {
	Storage storage.Storager
	Address        string
	ReportInterval time.Duration
	PollInterval   time.Duration
}

func NewAgent(address string, repInt time.Duration, pInt time.Duration) Agent {
	s := storage.NewIMS()
	log.Println("created storage for agent")
	return Agent{
		Storage: s,
		Address: address,
		ReportInterval: repInt,
		PollInterval: pInt,
	}
}

func (a Agent) Run() {

}