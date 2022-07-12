package server

import (
	"bufio"
	"encoding/json"
	"github.com/dayterr/go_agent_metrics/internal/storage"
	"io/ioutil"
	"log"
	"os"
)

func WriteJSON(path string, jsn []byte) {
	file, err := os.OpenFile(path, os.O_CREATE|os.O_RDWR|os.O_TRUNC, 0777)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	w := bufio.NewWriter(file)
	w.Write(jsn)
	w.Flush()
}

func LoadMetricsFromFile(filename string) (storage.InMemoryStorage, error) {
	if _, err := os.Stat(filename); err != nil {
		file, err := os.Create(filename)
		if err != nil {
			return storage.InMemoryStorage{}, err
		}
		file.Close()
		return storage.InMemoryStorage{}, nil
	}
	file, err := ioutil.ReadFile(filename)
	if err != nil {
		return storage.InMemoryStorage{}, err
	}

	s := storage.NewIMS()
	log.Println("s is", s)

	err = json.Unmarshal(file, &s)
	if err != nil {
		return storage.InMemoryStorage{}, err
	}
	log.Println("unmarshalled", s)
	for key, value := range s.GaugeField {
		s.SetGaugeFromMemStats(key, value.ToFloat())
	}
	for key, value := range s.CounterField {
		s.SetCounterFromMemStats(key, value.ToInt64())
	}
	return s, nil
}
