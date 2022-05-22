package main

import (
	"bytes"
	"fmt"
	"net/http"
	"runtime"
	"strconv"
	"time"
)

type gauge int64
type counter int64

func postGauge (v gauge, name string) {
	fmt.Println("Sending", v)
	url := fmt.Sprintf("http://localhost:8080/update/gauge/%v/%v/", name, v)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer([]byte(strconv.Itoa(int(v)))))
	if err != nil {
		panic(err)
	}
	req.Header.Set("Content-Type", "text/plain")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
}

func main() {
	m := &runtime.MemStats{}
	ticker := time.NewTicker(10 * time.Second)
	for {
		runtime.ReadMemStats(m)
		Alloc := gauge(m.Alloc)
		BuckHashSys := gauge(m.BuckHashSys)
		go func() {
			<-ticker.C
			postGauge(Alloc, "alloc")
			postGauge(BuckHashSys, "buckhashsys")
		}()
		time.Sleep(2 * time.Second)
	}
}
