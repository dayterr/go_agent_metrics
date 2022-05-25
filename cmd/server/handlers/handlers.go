package handlers

import (
	"github.com/go-chi/chi/v5"
	"html/template"
	"log"
	"net/http"
	"strconv"
)

var metrics = make(map[string]float64)
var counters = make(map[string]int)

func PostMetric(w http.ResponseWriter, r *http.Request) {
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
		if _, ok := metrics[mn]; ok {
			v := strconv.FormatFloat(metrics[mn], 'f', -1, 64)
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(v))
		} else {
			w.WriteHeader(http.StatusNotFound)
			return
		}
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

func GetIndex(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("cmd/server/index.html")
	if err != nil {
		log.Fatalln(err)
	}

	err = t.ExecuteTemplate(w, "index.html", metrics)
	if err != nil {
		log.Fatalln(err)
	}
}

func CreateRouter() chi.Router {
	r := chi.NewRouter()
	r.Route("/update", func(r chi.Router) {
		r.Post("/{mt}/{mn}/{v}", PostMetric)
	})
	r.Get("/value/{mt}/{mn}", GetMetric)
	r.Get("/", GetIndex)
	return r
}
