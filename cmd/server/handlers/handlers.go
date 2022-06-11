package handlers

import (
	"encoding/json"
	"html/template"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"

	"github.com/dayterr/go_agent_metrics/internal/agent"
)

var metrics = make(map[string]float64)
var counters = make(map[string]int64)

func GetValue(w http.ResponseWriter, r *http.Request) {
	var m agent.Metrics
	err := json.NewDecoder(r.Body).Decode(&m)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
	}
	switch m.MType {
	case agent.GaugeType:
		m.Value = metrics[m.ID]
		mJSON, err := m.MarshallJSON()
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
		}
		w.Write(mJSON)
	case agent.CounterType:
		m.Delta = counters[m.ID]
		mJSON, err := m.MarshallJSON()
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
		}
		w.Write(mJSON)
	default:
		w.WriteHeader(http.StatusBadRequest)
	}
}


func PostJSON(w http.ResponseWriter, r *http.Request) {
	var m agent.Metrics
	err := json.NewDecoder(r.Body).Decode(&m)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
	}
	switch m.MType{
	case agent.GaugeType:
		metrics[m.ID] = m.Value
		w.WriteHeader(http.StatusOK)
	case agent.CounterType:
		counters[m.ID] = m.Delta
		w.WriteHeader(http.StatusOK)
	default:
		w.WriteHeader(http.StatusNotFound)
	}
}

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
	case agent.GaugeType:
		valFloat, err := strconv.ParseFloat(value, 64)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		metrics[metricName] = valFloat
		w.WriteHeader(http.StatusOK)
	case agent.CounterType:
		valInt, err := strconv.Atoi(value)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		counters[metricName] += int64(valInt)
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
	case agent.GaugeType:
		if _, ok := metrics[metricName]; ok {
			value := strconv.FormatFloat(metrics[metricName], 'f', -1, 64)
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(value))
		} else {
			w.WriteHeader(http.StatusNotFound)
			return
		}
	case agent.CounterType:
		if _, ok := counters[metricName]; ok {
			c := strconv.Itoa(int(counters[metricName]))
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
	r.Post("/update", PostJSON)
	r.Route("/update", func(r chi.Router) {
		r.Post("/{metricType}/{metricName}/{value}", PostMetric)
	})
	r.Post("/value", GetValue)
	r.Get("/value/{metricType}/{metricName}", GetMetric)
	r.Get("/", GetIndex)
	return r
}
