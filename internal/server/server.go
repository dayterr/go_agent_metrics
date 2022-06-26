package server

import (
	"bufio"
	"github.com/dayterr/go_agent_metrics/cmd/server/handlers"
	"log"
	"os"
)

func WriteJSON(path string) {
	file, err := os.OpenFile(path, os.O_CREATE|os.O_RDWR|os.O_TRUNC, 0777)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	jsn, err := handlers.MarshallMetrics()
	if err != nil {
		log.Fatal(err)
	}
	w := bufio.NewWriter(file)
	w.Write(jsn)
	w.Flush()
}
