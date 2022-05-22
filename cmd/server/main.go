package main

import (
	"fmt"
	"github.com/go-chi/chi/v5"
	"net/http"
	"strconv"
)

var metrics = make(map[string]int)
var counter int

func PostGauge(w http.ResponseWriter, r *http.Request) {
	fmt.Println("workig")
	mt := chi.URLParam(r, "mt")
	mn := chi.URLParam(r,"mn")
	if mn == "" {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	v := chi.URLParam(r,"v")
	if v == "" {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	val, err := strconv.Atoi(v)
		if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	switch mt {
	case "gauge":
		metrics[mn] = val
		w.WriteHeader(http.StatusOK)
	case "counter":
		counter = val
		w.WriteHeader(http.StatusOK)
	default:
		w.WriteHeader(http.StatusNotImplemented)
	}
	/*else {
		fmt.Println("method GET", r.URL.Path)
		args := strings.Split(r.URL.Path, "/")
		name := args[3]
		m := strconv.Itoa(metrics[name])
		w.Write([]byte(m))
	}*/
}

func GetMetric(w http.ResponseWriter, r *http.Request) {
	mn := chi.URLParam(r, "mn")
	fmt.Println("cho works", mn)
	w.Write([]byte(mn))
}

func main() {
	r := chi.NewRouter()
	r.Route("/update", func(r chi.Router) {
		r.Post("/{mt}/{mn}/{v}", PostGauge)
	})
	r.Get("/value/gauge/{mn}", GetMetric)

	//http.HandleFunc("/update/", PostGauge)
	//http.HandleFunc("/update", GetUpdate)
	http.ListenAndServe(":8080", r)
}