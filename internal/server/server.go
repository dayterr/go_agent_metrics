package server

import (
	"bufio"
	"fmt"
	"github.com/dayterr/go_agent_metrics/cmd/server/handlers"
	"log"
	"os"
)

func WriteJSON(path string) {
	file, err := os.OpenFile(path, os.O_CREATE | os.O_RDWR , 0777)
	if err != nil {
		fmt.Println("create error", err)
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