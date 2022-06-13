package handlers

import (
	"encoding/json"
	"html/template"
	"net/http"
	"reflect"
	"strconv"

	"github.com/go-chi/chi/v5"

	"github.com/dayterr/go_agent_metrics/internal/agent"
)

var metrics = make(map[string]float64)
var counters = make(map[string]int64)

var metricJSON agent.MetricsJSON

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
	v := reflect.ValueOf(&metricJSON).Elem()
	switch m.MType {
	case agent.GaugeType:
		m.Value = v.Elem().FieldByName(m.ID).Float()
		mJSON, err := m.MarshallJSON()
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
		}
		w.Header().Set("content-type", "application/json")
		w.Write(mJSON)
	case agent.CounterType:
		m.Delta = v.Elem().FieldByName(m.ID).Int()
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
	v := reflect.ValueOf(&metricJSON)
	switch m.MType{
	case agent.GaugeType:
		v.Elem().FieldByName(m.ID).SetFloat(m.Value)
		w.WriteHeader(http.StatusOK)
	case agent.CounterType:
		uptValue := v.Elem().FieldByName(m.ID).Int() + m.Delta
		v.Elem().FieldByName(m.ID).SetInt(uptValue)
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
	v := reflect.ValueOf(&metricJSON)
	switch metricType {
	case agent.GaugeType:
		valFloat, err := strconv.ParseFloat(value, 64)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		v.Elem().FieldByName(metricName).SetFloat(valFloat)
		w.WriteHeader(http.StatusOK)
	case agent.CounterType:
		valInt, err := strconv.Atoi(value)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		uptValue := v.Elem().FieldByName(metricName).Int() + int64(valInt)
		v.Elem().FieldByName(metricName).SetInt(uptValue)
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
	var metrics map[string]interface{}
	data, _ := json.Marshal(metricJSON)
	json.Unmarshal(data, &metrics)
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
		r.Post("/", PostJSON)
		r.Post("/{metricType}/{metricName}/{value}", PostMetric)
	})
	r.Post("/value/", GetValue)
	r.Get("/value/{metricType}/{metricName}", GetMetric)
	r.Get("/", GetIndex)
	return r
}
