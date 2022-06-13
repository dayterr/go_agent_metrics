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
		value      agent.Gauge
		metricName string
		metricType string
		want       error
	}{
		{name: "no error for gauge metric", value: agent.Gauge(63.3), metricName: "Some_Metric", metricType: "gauge", want: nil},
		{name: "no error for gauge metric without decimal part", value: agent.Gauge(63), metricName: "Some_Metric", metricType: "gauge", want: nil},
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
			v := agent.PostMetric(tt.value, tt.metricName, tt.metricType)
			assert.Nil(t, v)
		})
	}

}

func TestPostCounter(t *testing.T) {
	tests := []struct {
		name       string
		value      agent.Counter
		metricName string
		metricType string
		want       error
	}{
		{name: "no error for counter metric", value: agent.Counter(63), metricName: "Some_Counter", metricType: "counter", want: nil},
		{name: "no error for counter metric zero", value: agent.Counter(0), metricName: "Some_Counter", metricType: "counter", want: nil},
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
			v := agent.PostCounter(tt.value, tt.metricName, tt.metricType)
			assert.Nil(t, v)
		})
	}
}
