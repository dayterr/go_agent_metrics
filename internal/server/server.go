package server

import (
	"bufio"
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
	log.Println("wrote to file")
	w.Flush()
}
