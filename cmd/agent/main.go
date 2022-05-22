package main

import (
	"fmt"
	"runtime"
	"strconv"
	"time"

	"github.com/levigross/grequests"
)

type gauge int64
type counter int64

var metrics = make(map[string]gauge)

func PostMetric (v gauge, name string, mt string) error {
	fmt.Println("Sending", v)
	url := fmt.Sprintf("http://localhost:8080/update/%v/%v/%v/", mt, name, v)
	//req, err := http.NewRequest("POST", url, bytes.NewBuffer([]byte(strconv.Itoa(int(v)))))
	_, err := grequests.Post(url, &grequests.RequestOptions{Data: map[string]string {name: strconv.Itoa(int(v))},
		Headers: map[string]string{"ContentType": "text/plain"}})
	if err != nil {
		return err
	}

	return nil
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
			PostMetric(Alloc, "alloc", "gauge")
			PostMetric(BuckHashSys, "buckhashsys", "gauge")
		}()
		time.Sleep(2 * time.Second)
	}
}
