package main

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

type gauge float64

func TestPostGauge(t *testing.T) {
	tests := []struct {
		name  string
		v Gauge
		nm string
		tm string
		want  error
	}{
		{name: "no error", v: Gauge(63.3), nm: "Some Metric", tm: "gauge", want: nil},
		{name: "no error without decimal part", v: Gauge(63), nm: "Some Metric", tm: "gauge", want: nil},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := PostMetric(tt.v, tt.nm, tt.tm)
			assert.Nil(t, v)
		})
	}
}