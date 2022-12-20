package agent

import (
	"net"
	"net/http/httptest"
	"testing"

	"github.com/dayterr/go_agent_metrics/internal/storage"

	"github.com/stretchr/testify/assert"
)

const address = "localhost:8080"

func TestPostGauge(t *testing.T) {

	tests := []struct {
		name       string
		value      storage.Gauge
		metricName string
		metricType string
		want       error
	}{
		{name: "no error for gauge metric", value: storage.Gauge(63.3), metricName: "Some_Metric", metricType: "gauge", want: nil},
		{name: "no error for gauge metric without decimal part", value: storage.Gauge(63), metricName: "Some_Metric", metricType: "gauge", want: nil},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h, err := NewAsyncHandler("", "", false)
			assert.NoError(t, err)
			r, err := CreateRouterWithAsyncHandler("", false, h)
			assert.NoError(t, err)
			ts := httptest.NewUnstartedServer(r)
			url := "127.0.0.1:8080"
			l, err := net.Listen("tcp", url)
			assert.NoError(t, err)
			ts.Listener = l
			ts.Start()
			defer ts.Close()
			v := PostGauge(tt.value, tt.metricName, address, "")
			assert.Nil(t, v)
		})
	}

}

func TestPostCounter(t *testing.T) {
	tests := []struct {
		name       string
		value      storage.Counter
		metricName string
		metricType string
		want       error
	}{
		{name: "no error for counter metric", value: storage.Counter(63), metricName: "Some_Counter", metricType: "counter", want: nil},
		{name: "no error for counter metric zero", value: storage.Counter(0), metricName: "Some_Counter", metricType: "counter", want: nil},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h, err := NewAsyncHandler("", "", false)
			assert.NoError(t, err)
			r, err := CreateRouterWithAsyncHandler("", false, h)
			assert.NoError(t, err)
			ts := httptest.NewUnstartedServer(r)
			url := "127.0.0.1:8080"
			l, err := net.Listen("tcp", url)
			assert.NoError(t, err)
			ts.Listener = l
			ts.Start()
			defer ts.Close()
			v := PostCounter(tt.value, tt.metricName, address, "", "")
			assert.Nil(t, v)
		})
	}
}
