package server

import (
	"bufio"
	"context"
	"encoding/json"
	"io/ioutil"
	"log"
	"os"

	"github.com/dayterr/go_agent_metrics/internal/storage"
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
	ctx := context.Background()
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

	log.Println("file is", string(file))
	s := storage.NewIMS()
	log.Println("s is", s)

	err = json.Unmarshal(file, &s)
	if err != nil {
		return storage.InMemoryStorage{}, err
	}
	log.Println("unmarshalled", s)
	for key, value := range s.GaugeField {
		s.SetGaugeFromMemStats(ctx, key, value.ToFloat())
	}
	for key, value := range s.CounterField {
		s.SetCounterFromMemStats(ctx, key, value.ToInt64())
	}
	return s, nil
}
