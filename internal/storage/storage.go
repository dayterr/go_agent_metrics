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

type TempStorage struct {
	GaugeField   map[string]Gauge   `json:"Gauge"`
	CounterField map[string]Counter `json:"Counter"`
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
			var ts TempStorage
			err = json.Unmarshal(file, &ts)
			if err != nil {
				return err
			}
			for k, v := range ts.GaugeField {
				s.GaugeField[k] = v
			}
			for k, v := range ts.CounterField {
				s.CounterField[k] = v
			}
		}
	}
	return nil
}
