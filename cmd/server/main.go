package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
)
// HelloWorld — обработчик запроса.
var metrics = make(map[string][]int)
func PostGauge(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		_, err := ioutil.ReadAll(r.Body)
		// обрабатываем ошибку
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
		args := strings.Split(r.URL.Path, "/")
		name := args[3]
		metric, err := strconv.Atoi(args[4])
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
		metrics[name] = append(metrics[name], metric)
	} else {
		fmt.Println("method GET", r.URL.Path)
		args := strings.Split(r.URL.Path, "/")
		name := args[3]
		l := len(metrics[name]) - 1
		m := strconv.Itoa(metrics[name][l])
		w.Write([]byte(m))
	}
}

func main() {
	http.HandleFunc("/update/gauge/", PostGauge)
	http.ListenAndServe(":8080", nil)
}