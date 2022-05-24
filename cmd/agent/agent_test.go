package main

import (
	"github.com/stretchr/testify/assert"
	"testing"

)

func TestPostGauge(t *testing.T) {
	tests := []struct {
		name  string
		v Gauge
		nm string
		tm string
		want  error
	}{
		{name: "no error for gauge metric", v: Gauge(63.3), nm: "Some_Metric", tm: "gauge", want: nil},
		{name: "no error for gauge metric without decimal part", v: Gauge(63), nm: "Some_Metric", tm: "gauge", want: nil},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := PostMetric(tt.v, tt.nm, tt.tm)
			assert.Nil(t, v)
		})
	}

}

func TestPostCounter(t *testing.T) {
	tests := []struct {
		name  string
		v Counter
		nm string
		tm string
		want  error
	}{
		{name: "no error for counter metric", v: Counter(63), nm: "Some_Counter", tm: "counter", want: nil},
		{name: "no error for counter metric zero", v: Counter(0), nm: "Some_Counter", tm: "counter", want: nil},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := PostCounter(tt.v, tt.nm, tt.tm)
			assert.Nil(t, v)
		})
	}
}