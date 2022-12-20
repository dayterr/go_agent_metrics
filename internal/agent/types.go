package agent

import (
	"os"
	"os/signal"
	"syscall"
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

func (a Agent) Run() {
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt, syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)
	exitChan := make(chan int)
	tickerCollectMetrics := time.NewTicker(a.PollInterval)
	tickerReportMetrics := time.NewTicker(a.ReportInterval)
	go func() {
		for {
			select {
			case <-tickerCollectMetrics.C:
				a.ReadMetrics()
			case <-tickerReportMetrics.C:
				err := a.PostMany()
				if err != nil {
					a.PostAll()
				}
			case s := <-signalChan:
				switch s {
				case syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT:
					exitChan <- 0
				}
			}
		}
	}()
	exitCode := <-exitChan
	os.Exit(exitCode)
}
