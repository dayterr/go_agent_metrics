package agent

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	pb "github.com/dayterr/go_agent_metrics/internal/grpc/proto"
	"log"
	"math/rand"
	"net"
	"net/http"
	"runtime"
	"strconv"
	"time"

	"github.com/levigross/grequests"
	"github.com/shirou/gopsutil/v3/mem"

	"github.com/dayterr/go_agent_metrics/internal/encryption"
	"github.com/dayterr/go_agent_metrics/internal/hash"
	"github.com/dayterr/go_agent_metrics/internal/metric"
	"github.com/dayterr/go_agent_metrics/internal/storage"
)

const GaugeType = "gauge"
const CounterType = "counter"

func PostCounter(value storage.Counter, metricName string, address string, key, cryptoKey string) error {
	url := fmt.Sprintf("http://%v/update/", address)
	delta := value.ToInt64()
	metric := metric.Metrics{ID: metricName, MType: CounterType, Delta: &delta}
	if key != "" {
		metric.Hash = hash.EncryptMetric(metric, key)
	}
	mJSON, err := json.Marshal(metric)
	if err != nil {
		return err
	}

	client := &http.Client{}

	if cryptoKey != "" {
		enc := encryption.NewEncryptor(cryptoKey)
		client.Transport = encryption.NewRoundTripperWithEncryption(enc)
	}

	addr, err := GetMyIP()
	if err != nil {
		return nil
	}

	_, err = grequests.Post(url, &grequests.RequestOptions{JSON: mJSON,
		Headers: map[string]string{"ContentType": "application/json", "X-Real-IP": addr}, HTTPClient: client})
	if err != nil {
		return err
	}
	return nil
}

func PostGauge(value storage.Gauge, metricName string, address string, key, cryptoKey string) error {
	url := fmt.Sprintf("http://%v/update/", address)
	v := value.ToFloat()
	metric := metric.Metrics{ID: metricName, MType: GaugeType, Value: &v}
	if key != "" {
		metric.Hash = hash.EncryptMetric(metric, key)
	}
	mJSON, err := json.Marshal(metric)
	if err != nil {
		return err
	}

	client := &http.Client{}

	if cryptoKey != "" {
		enc := encryption.NewEncryptor(cryptoKey)
		client.Transport = encryption.NewRoundTripperWithEncryption(enc)
	}

	addr, err := GetMyIP()
	if err != nil {
		return nil
	}

	_, err = grequests.Post(url, &grequests.RequestOptions{JSON: mJSON,
		Headers: map[string]string{"ContentType": "application/json", "X-Real-IP": addr}, HTTPClient: client})
	if err != nil {
		return err
	}
	return nil
}

func (a Agent) PostAll() {
	ctx := context.Background()
	gauges, err := a.Storage.GetGauges(ctx)
	if err != nil {
		log.Println("getting gauges error", err)
	}
	counters, err := a.Storage.GetCounters(ctx)
	if err != nil {
		log.Println("getting counters error", err)
	}
	for k, v := range gauges {
		PostGauge(v, k, a.Address, a.Key, a.CryptoKey)
	}
	for k, v := range counters {
		PostCounter(v, k, a.Address, a.Key, a.CryptoKey)
	}
}

func (a Agent) ReadMetrics() {
	m := &runtime.MemStats{}
	v, _ := mem.VirtualMemory()
	runtime.ReadMemStats(m)
	ctx := context.Background()
	a.Storage.SetGaugeFromMemStats(ctx, "Alloc", float64(m.Alloc))
	a.Storage.SetGaugeFromMemStats(ctx, "BuckHashSys", float64(m.BuckHashSys))
	a.Storage.SetGaugeFromMemStats(ctx, "Frees", float64(m.Frees))
	a.Storage.SetGaugeFromMemStats(ctx, "GCCPUFraction", m.GCCPUFraction)
	a.Storage.SetGaugeFromMemStats(ctx, "GCSys", float64(m.GCSys))
	a.Storage.SetGaugeFromMemStats(ctx, "HeapAlloc", float64(m.HeapAlloc))
	a.Storage.SetGaugeFromMemStats(ctx, "HeapIdle", float64(m.HeapIdle))
	a.Storage.SetGaugeFromMemStats(ctx, "HeapInuse", float64(m.HeapInuse))
	a.Storage.SetGaugeFromMemStats(ctx, "HeapObjects", float64(m.HeapObjects))
	a.Storage.SetGaugeFromMemStats(ctx, "HeapReleased", float64(m.HeapReleased))
	a.Storage.SetGaugeFromMemStats(ctx, "HeapSys", float64(m.HeapSys))
	a.Storage.SetGaugeFromMemStats(ctx, "LastGC", float64(m.HeapAlloc))
	a.Storage.SetGaugeFromMemStats(ctx, "Lookups", float64(m.Lookups))
	a.Storage.SetGaugeFromMemStats(ctx, "MCacheInuse", float64(m.MCacheInuse))
	a.Storage.SetGaugeFromMemStats(ctx, "MCacheSys", float64(m.MCacheSys))
	a.Storage.SetGaugeFromMemStats(ctx, "MSpanInuse", float64(m.MSpanInuse))
	a.Storage.SetGaugeFromMemStats(ctx, "MSpanSys", float64(m.MSpanSys))
	a.Storage.SetGaugeFromMemStats(ctx, "Mallocs", float64(m.Mallocs))
	a.Storage.SetGaugeFromMemStats(ctx, "NextGC", float64(m.NextGC))
	a.Storage.SetGaugeFromMemStats(ctx, "NumForcedGC", float64(m.NumForcedGC))
	a.Storage.SetGaugeFromMemStats(ctx, "NumGC", float64(m.NumGC))
	a.Storage.SetGaugeFromMemStats(ctx, "OtherSys", float64(m.OtherSys))
	a.Storage.SetGaugeFromMemStats(ctx, "PauseTotalNs", float64(m.PauseTotalNs))
	a.Storage.SetGaugeFromMemStats(ctx, "StackInuse", float64(m.StackInuse))
	a.Storage.SetGaugeFromMemStats(ctx, "StackSys", float64(m.StackSys))
	a.Storage.SetGaugeFromMemStats(ctx, "Sys", float64(m.Sys))
	a.Storage.SetGaugeFromMemStats(ctx, "TotalAlloc", float64(m.TotalAlloc))
	a.Storage.SetGaugeFromMemStats(ctx, "RandomValue", rand.Float64())
	a.Storage.SetGaugeFromMemStats(ctx, "TotalMemory", float64(v.Total))
	a.Storage.SetGaugeFromMemStats(ctx, "FreeMemory", float64(v.Free))
	a.Storage.SetGaugeFromMemStats(ctx, "CPUutilization1", float64(v.Used))
	a.Storage.SetCounterFromMemStats(ctx, "PollCount", 1)
}

func (a Agent) PostMany() error {
	var listMetrics []metric.Metrics

	ctx := context.Background()
	gs, err := a.Storage.GetGauges(ctx)
	if err != nil {
		return err
	}
	cs, err := a.Storage.GetCounters(ctx)
	if err != nil {
		return err
	}
	if len(gs) == 0 && len(cs) == 0 {
		return errors.New("the batch is empty")
	}

	gs, err = a.Storage.GetGauges(ctx)
	if err != nil {
		return err
	}
	for key, value := range gs {
		var m metric.Metrics
		m.ID = key
		v := value.ToFloat()
		m.Value = &v
		m.MType = GaugeType
		if a.Key != "" {
			m.Hash = hash.EncryptMetric(m, key)
		}
		listMetrics = append(listMetrics, m)
	}
	cs, err = a.Storage.GetCounters(ctx)
	if err != nil {
		return err
	}
	for key, value := range cs {
		var m metric.Metrics
		m.ID = key
		d := value.ToInt64()
		m.Delta = &d
		m.MType = CounterType
		if a.Key != "" {
			m.Hash = hash.EncryptMetric(m, key)
		}
		listMetrics = append(listMetrics, m)
	}

	jsn, err := json.Marshal(listMetrics)
	if err != nil {

		return err
	}

	url := fmt.Sprintf("http://%v/updates/", a.Address)
	_, err = grequests.Post(url, &grequests.RequestOptions{JSON: jsn,
		Headers: map[string]string{"ContentType": "application/json"}, DisableCompression: false})
	if err != nil {
		return err
	}

	return nil
}

func (a Agent) Run(ctx context.Context) {
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
			case <-ctx.Done():
				return
			}
		}
	}()
}

func (a Agent) RungRPC(ctx context.Context) {
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
					a.PostAllgRPC()
				}
			case <-ctx.Done():
				return
			}
		}
	}()
}

func (a Agent) PostMetricgRPC(value string, metricName string, metricType string) error {
	var metric pb.Metric
	switch metricType {
	case GaugeType:
		valFloat, err := strconv.ParseFloat(value, 64)
		if err != nil {
			return err
		}
		metric.Value = valFloat
		metric.Type = pb.Metric_GAUGE
		metric.Id = metricName

	case CounterType:
		valInt, err := strconv.Atoi(value)
		if err != nil {
			return err
		}
		metric.Delta = int64(valInt)
		metric.Type = pb.Metric_GAUGE
		metric.Id = metricName

	default:
		return errors.New("unknown metric type")
	}

	ctx := context.Background()
	_, err := a.GRPCClient.PostMetric(ctx, &pb.PostMetricRequest{Metrics: &metric})
	if err != nil {
		return err
	}

	return nil
}

func (a Agent) PostAllgRPC() {
	ctx := context.Background()
	gauges, err := a.Storage.GetGauges(ctx)
	if err != nil {
		log.Println("getting gauges error", err)
	}
	counters, err := a.Storage.GetCounters(ctx)
	if err != nil {
		log.Println("getting counters error", err)
	}
	for k, v := range gauges {
		valStr := fmt.Sprintf("%f", v)
		a.PostMetricgRPC(valStr, k, GaugeType)
	}
	for k, v := range counters {
		valStr := fmt.Sprintf("%d", v)
		a.PostMetricgRPC(valStr, k, CounterType)
	}
}

func GetMyIP() (string, error) {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return "", err
	}

	for _, addr := range addrs {
		ipnet, ok := addr.(*net.IPNet)
		if ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				return ipnet.IP.String(), nil
			}
		}
	}

	return "", nil
}
