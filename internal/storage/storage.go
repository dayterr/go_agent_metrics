package storage

import (
	"encoding/json"
	"io/ioutil"
	"os"
)

type Gauge float64
type Counter int64

type Storage struct {
	GaugeField   map[string]Gauge
	CounterField map[string]Counter
}

func New() Storage {
	return Storage{
		GaugeField: make(map[string]Gauge),
		CounterField: make(map[string]Counter),
	}
}

func (g Gauge) ToFloat() float64 {
	return float64(g)
}

func (c Counter) ToInt64() int64 {
	return int64(c)
}

func (c Counter) ToInt() int {
	return int(c)
}

func (s Storage) LoadMetricsFromJSON(filename string, isRestored bool) error {
	if isRestored {
		if _, err := os.Stat(filename); err == nil {
			file, err := ioutil.ReadFile(filename)
			if err != nil {
				return err
			}
			err = json.Unmarshal(file, &s)
			if err != nil {
				return err
			}
		}
	}
	return nil
}
