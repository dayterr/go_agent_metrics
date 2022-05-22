package main

import (
	"fmt"
	"github.com/go-chi/chi/v5"
	"net/http"
	"strconv"
)

var metrics = make(map[string]float64)
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
	switch mt {
	case "gauge":
		val, err := strconv.ParseFloat(v, 64)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		metrics[mn] = val
		w.WriteHeader(http.StatusOK)
	case "counter":
		val, err := strconv.Atoi(v)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		counters[mn] += val
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
		if _, ok := counters[mn]; ok {
			c := strconv.Itoa(counters[mn])
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(c))
		} else {
			w.WriteHeader(http.StatusNotFound)
			return
		}
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