package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
)

var metrics = make(map[string]int)

func PostGauge(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {

		b, err := ioutil.ReadAll(r.Body)
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
		args := strings.Split(r.URL.Path, "/")
		name := args[3]
		metric, err := strconv.Atoi(string(b))
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
		metrics[name] = metric
		w.WriteHeader(http.StatusCreated)
	} else {
		fmt.Println("method GET", r.URL.Path)
		args := strings.Split(r.URL.Path, "/")
		name := args[3]
		m := strconv.Itoa(metrics[name])
		w.Write([]byte(m))
	}
}

func GetUpdate(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		w.WriteHeader(http.StatusOK)
	}
}

func main() {
	http.HandleFunc("/update/gauge", PostGauge)
	http.HandleFunc("/update", GetUpdate)
	http.ListenAndServe(":8080", nil)
}