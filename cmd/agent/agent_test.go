package main

import (
	"net"
	"net/http/httptest"
	"testing"

	"github.com/dayterr/go_agent_metrics/cmd/server/handlers"
	"github.com/stretchr/testify/assert"

	"github.com/dayterr/go_agent_metrics/internal/agent"
)

func TestPostGauge(t *testing.T) {

	tests := []struct {
		name       string
		value      Gauge
		metricName string
		metricType string
		want       error
	}{
		{name: "no error for gauge metric", value: Gauge(63.3), metricName: "Some_Metric", metricType: "gauge", want: nil},
		{name: "no error for gauge metric without decimal part", value: Gauge(63), metricName: "Some_Metric", metricType: "gauge", want: nil},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := handlers.CreateRouter()
			ts := httptest.NewUnstartedServer(r)
			url := "127.0.0.1:8080"
			l, err := net.Listen("tcp", url)
			assert.NoError(t, err)
			ts.Listener = l
			ts.Start()
			defer ts.Close()
			v := PostMetric(tt.value, tt.metricName, tt.metricType)
			assert.Nil(t, v)
		})
	}

}

func TestPostCounter(t *testing.T) {
	tests := []struct {
		name       string
		value      Counter
		metricName string
		metricType string
		want       error
	}{
		{name: "no error for counter metric", value: Counter(63), metricName: "Some_Counter", metricType: "counter", want: nil},
		{name: "no error for counter metric zero", value: Counter(0), metricName: "Some_Counter", metricType: "counter", want: nil},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := handlers.CreateRouter()
			ts := httptest.NewUnstartedServer(r)
			url := "127.0.0.1:8080"
			l, err := net.Listen("tcp", url)
			assert.NoError(t, err)
			ts.Listener = l
			ts.Start()
			defer ts.Close()
			v := PostCounter(tt.value, tt.metricName, tt.metricType)
			assert.Nil(t, v)
		})
	}
}
