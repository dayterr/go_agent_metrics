package handlers

import (
	"html/template"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

var metrics = make(map[string]float64)
var counters = make(map[string]int)

const gauge = "gauge"
const counter = "counter"

func PostMetric(w http.ResponseWriter, r *http.Request) {
	metricType := chi.URLParam(r, "metricType")
	metricName := chi.URLParam(r, "metricName")
	if metricName == "" {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	value := chi.URLParam(r, "value")
	if value == "" {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	switch metricType {
	case gauge:
		valFloat, err := strconv.ParseFloat(value, 64)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		metrics[metricName] = valFloat
		w.WriteHeader(http.StatusOK)
	case counter:
		valInt, err := strconv.Atoi(value)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		counters[metricName] += valInt
		w.WriteHeader(http.StatusOK)
	default:
		w.WriteHeader(http.StatusNotImplemented)
	}
}

func GetMetric(w http.ResponseWriter, r *http.Request) {
	metricType := chi.URLParam(r, "metricType")
	metricName := chi.URLParam(r, "metricName")
	if metricName == "" {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	switch metricType {
	case gauge:
		if _, ok := metrics[metricName]; ok {
			value := strconv.FormatFloat(metrics[metricName], 'f', -1, 64)
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(value))
		} else {
			w.WriteHeader(http.StatusNotFound)
			return
		}
	case counter:
		if _, ok := counters[metricName]; ok {
			c := strconv.Itoa(counters[metricName])
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

func GetIndex(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("cmd/server/index.html")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	err = t.ExecuteTemplate(w, "index.html", metrics)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func CreateRouter() chi.Router {
	r := chi.NewRouter()
	r.Route("/update", func(r chi.Router) {
		r.Post("/{metricType}/{metricName}/{value}", PostMetric)
	})
	r.Get("/value/{metricType}/{metricName}", GetMetric)
	r.Get("/", GetIndex)
	return r
}
