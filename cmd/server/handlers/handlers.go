package handlers

import (
	"compress/gzip"
	"encoding/json"
	"github.com/dayterr/go_agent_metrics/internal/storage"
	"github.com/go-chi/chi/v5"
	"html/template"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/dayterr/go_agent_metrics/internal/agent"
)

var allMetrics = storage.New()

type gzipWriter struct {
	http.ResponseWriter
	Writer io.Writer
}

func (w gzipWriter) Write(b []byte) (int, error) {
	return w.Writer.Write(b)
}

func gzipHandle(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		if !strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
			next.ServeHTTP(w, r)
			return
		}

		gz, err := gzip.NewWriterLevel(w, gzip.BestSpeed)
		if err != nil {
			io.WriteString(w, err.Error())
			return
		}
		defer gz.Close()

		w.Header().Set("Content-Encoding", "gzip")
		next.ServeHTTP(gzipWriter{ResponseWriter: w, Writer: gz}, r)
	})
}

func MarshallMetrics() ([]byte, error) {
	jsn, err := json.Marshal(allMetrics)
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
		m.Value = allMetrics.GaugeField[m.ID].ToFloat()
		mJSON, err := json.Marshal(m)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
		}
		w.Header().Set("content-type", "application/json")
		w.Write(mJSON)
	case agent.CounterType:
		m.Delta = allMetrics.CounterField[m.ID].ToInt64()
		mJSON, err := json.Marshal(m)
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
	switch m.MType {
	case agent.GaugeType:
		allMetrics.GaugeField[m.ID] = storage.Gauge(m.Value)
		w.WriteHeader(http.StatusOK)
	case agent.CounterType:
		allMetrics.CounterField[m.ID] += storage.Counter(m.Delta)
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
	value := strings.Trim(chi.URLParam(r, "value"), "\n")
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
		allMetrics.GaugeField[metricName] = storage.Gauge(valFloat)
		w.WriteHeader(http.StatusOK)
	case agent.CounterType:
		valInt, err := strconv.Atoi(value)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		allMetrics.CounterField[metricName] += storage.Counter(int64(valInt))
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
		if _, ok := allMetrics.GaugeField[metricName]; ok {
			value := strconv.FormatFloat(allMetrics.GaugeField[metricName].ToFloat(), 'f', -1, 64)
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(value))
		} else {
			w.WriteHeader(http.StatusNotFound)
			return
		}
	case agent.CounterType:
		if _, ok := allMetrics.CounterField[metricName]; ok {
			c := strconv.Itoa(allMetrics.CounterField[metricName].ToInt())
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
	w.Header().Set("content-type", "text/html; charset=utf-8")
	t, err := template.ParseFiles("template/index.html")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	err = t.ExecuteTemplate(w, "index.html", allMetrics.GaugeField)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func CreateRouter(filename string, isRestored bool) chi.Router {
	err := allMetrics.LoadMetricsFromJSON(filename, isRestored)
	if err != nil {
		log.Fatal(err)
	}
	r := chi.NewRouter()
	r.Use(gzipHandle)
	r.Route("/update", func(r chi.Router) {
		r.Post("/", PostJSON)
		r.Post("/{metricType}/{metricName}/{value}", PostMetric)
	})
	r.Post("/value/", GetValue)
	r.Get("/value/{metricType}/{metricName}", GetMetric)
	r.Get("/", GetIndex)
	return r
}
