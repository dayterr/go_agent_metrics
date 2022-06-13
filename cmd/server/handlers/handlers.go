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

type AllMetrics struct {
	Gauge map[string]float64
	Counter map[string]int64
}

var allMetrics AllMetrics = AllMetrics{
	metrics,
	counters,
}


func MarshallMetrics() ([]byte, error){
	jsn, err := json.Marshal(metrics)
	if err != nil {
		return nil, err
	}
	return jsn, nil
}

func MarshallCounters() ([]byte, error){
	jsn, err := json.Marshal(counters)
	if err != nil {
		return nil, err
	}
	return jsn, nil
}

func UnmarshallMetrics() ([]byte, error){
	jsn, err := json.Marshal(metrics)
	if err != nil {
		return nil, err
	}
	return jsn, nil
}

func GetValue(w http.ResponseWriter, r *http.Request) {
	var m agent.Metrics
	err := json.NewDecoder(r.Body).Decode(&m)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
	}
	switch m.MType {
	case agent.GaugeType:
		m.Value = allMetrics.Gauge[m.ID]
		mJSON, err := m.MarshallJSON()
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
		}
		w.Header().Set("content-type", "application/json")
		w.Write(mJSON)
	case agent.CounterType:
		m.Delta = allMetrics.Counter[m.ID]
		mJSON, err := m.MarshallJSON()
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
		}
		w.Header().Set("content-type", "application/json")
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
		allMetrics.Gauge[m.ID] = m.Value
		w.WriteHeader(http.StatusOK)
	case agent.CounterType:
		allMetrics.Counter[m.ID] += m.Delta
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
		allMetrics.Gauge[metricName] = valFloat
		w.WriteHeader(http.StatusOK)
	case agent.CounterType:
		valInt, err := strconv.Atoi(value)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		allMetrics.Counter[metricName] = int64(valInt)
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
		if _, ok := allMetrics.Gauge[metricName]; ok {
			value := strconv.FormatFloat(allMetrics.Gauge[metricName], 'f', -1, 64)
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(value))
		} else {
			w.WriteHeader(http.StatusNotFound)
			return
		}
	case agent.CounterType:
		if _, ok := allMetrics.Counter[metricName]; ok {
			c := strconv.Itoa(int(allMetrics.Counter[metricName]))
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

	err = t.ExecuteTemplate(w, "index.html", allMetrics.Gauge)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func CreateRouter() chi.Router {
	r := chi.NewRouter()
	r.Route("/update", func(r chi.Router) {
		r.Post("/", PostJSON)
		r.Post("/{metricType}/{metricName}/{value}", PostMetric)
	})
	r.Post("/value/", GetValue)
	r.Get("/value/{metricType}/{metricName}", GetMetric)
	r.Get("/", GetIndex)
	return r
}
