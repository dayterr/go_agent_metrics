package main

import (
	"fmt"
	"github.com/go-chi/chi/v5"
	"net/http"
	"strconv"
)

var metrics = make(map[string]int)
var counters = make(map[string]int)

func PostGauge(w http.ResponseWriter, r *http.Request) {
	mt := chi.URLParam(r, "mt")
	mn := chi.URLParam(r,"mn")
	if mn == "" {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	v := chi.URLParam(r,"v")
	fmt.Println("v", v)
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
		fmt.Println(val)
		metrics[mn] = val
		w.WriteHeader(http.StatusOK)
	case "counter":
		counters[mn] = val
		w.WriteHeader(http.StatusOK)
	default:
		w.WriteHeader(http.StatusNotImplemented)
	}
}

func GetMetric(w http.ResponseWriter, r *http.Request) {
	mt := chi.URLParam(r, "mt")
	mn := chi.URLParam(r, "mn")
	if mn == "" {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	switch mt {
	case "gauge":
		fmt.Println(metrics[mn])
		v := strconv.Itoa(metrics[mn])
		fmt.Println(v)
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(v))
	case "counter":
		c := strconv.Itoa(counters[mn])
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(c))
	default:
		w.WriteHeader(http.StatusNotFound)
		return
	}
}

func main() {
	r := chi.NewRouter()
	r.Route("/update", func(r chi.Router) {
		r.Post("/{mt}/{mn}/{v}", PostGauge)
	})
	r.Get("/value/{mt}/{mn}", GetMetric)
	http.ListenAndServe(":8080", r)
}