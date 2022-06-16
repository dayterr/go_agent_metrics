package server

import (
	"bufio"
	"fmt"
	"github.com/dayterr/go_agent_metrics/cmd/server/handlers"
	"log"
	"os"
)

func WriteJSON(path string) {
	l, _ := os.Getwd()
	fmt.Println("writing to a file")
	file, err := os.OpenFile(l + path, os.O_CREATE | os.O_RDWR | os.O_TRUNC, 0777)
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