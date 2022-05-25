package main

import (
	"fmt"
	"math/rand"
	"os"
	"os/signal"
	"runtime"
	"strconv"
	"syscall"
	"time"

	"github.com/levigross/grequests"
)

type Gauge float64
type Counter int64

var metrics = make(map[string]Gauge)
var counters = make(map[string]Counter)

func PostMetric (v Gauge, name string, mt string) error {
	url := fmt.Sprintf("http://localhost:8080/update/%v/%v/%v", mt, name, v)
	_, err := grequests.Post(url, &grequests.RequestOptions{Data: map[string]string {name: strconv.Itoa(int(v))},
		Headers: map[string]string{"ContentType": "text/plain"}})
	if err != nil {
		return err
	}

	return nil
}

func PostCounter (v Counter, name string, mt string) error {
	url := fmt.Sprintf("http://localhost:8080/update/%v/%v/%v", mt, name, v)
	_, err := grequests.Post(url, &grequests.RequestOptions{Data: map[string]string {name: strconv.Itoa(int(v))},
		Headers: map[string]string{"ContentType": "text/plain"}})
	if err != nil {
		return err
	}
	return nil
}

func main() {
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt, syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)
	exitChan := make(chan int)
	m := &runtime.MemStats{}
	ticker := time.NewTicker(10 * time.Second)
	go func() {
		for {
			runtime.ReadMemStats(m)
			metrics["Alloc"] = Gauge(m.Alloc)
			metrics["BuckHashSys"] = Gauge(m.BuckHashSys)
			metrics["Frees"] = Gauge(m.Frees)
			metrics["GCCPUFraction"] = Gauge(m.GCCPUFraction)
			metrics["GCSys"] = Gauge(m.GCSys)
			metrics["HeapAlloc"] = Gauge(m.HeapAlloc)
			metrics["HeapIdle"] = Gauge(m.HeapIdle)
			metrics["HeapInuse"] = Gauge(m.HeapInuse)
			metrics["HeapObjects"] = Gauge(m.HeapObjects)
			metrics["HeapReleased"] = Gauge(m.HeapReleased)
			metrics["HeapSys"] = Gauge(m.HeapSys)
			metrics["LastGC"] = Gauge(m.HeapAlloc)
			metrics["Lookups"] = Gauge(m.Lookups)
			metrics["MCacheInuse"] = Gauge(m.MCacheInuse)
			metrics["MCacheSys"] = Gauge(m.MCacheSys)
			metrics["MSpanInuse"] = Gauge(m.MSpanInuse)
			metrics["MSpanSys"] = Gauge(m.MSpanSys)
			metrics["Mallocs"] = Gauge(m.Mallocs)
			metrics["NextGC"] = Gauge(m.NextGC)
			metrics["NumForcedGC"] = Gauge(m.NumForcedGC)
			metrics["NumGC"] = Gauge(m.NumGC)
			metrics["OtherSys"] = Gauge(m.OtherSys)
			metrics["PauseTotalNs"] = Gauge(m.PauseTotalNs)
			metrics["StackInuse"] = Gauge(m.StackInuse)
			metrics["StackSys"] = Gauge(m.StackSys)
			metrics["Sys"] = Gauge(m.Sys)
			metrics["TotalAlloc"] = Gauge(m.TotalAlloc)
			metrics["RandomValue"] = Gauge(rand.Float64())
			counters["Counter"] += 1
			go func() {
				<-ticker.C
				for k, v := range metrics {
					PostMetric(v, k, "gauge")
				}
				for k, v := range counters {
					PostCounter(v, k, "counter")
				}
			}()
			time.Sleep(2 * time.Second)
			s := <- signalChan
			switch s {
			case syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT:
				exitChan <- 0
			}
		}
	}()
	
	exitCode := <-exitChan
	os.Exit(exitCode)
}
