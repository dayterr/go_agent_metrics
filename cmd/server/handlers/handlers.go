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

type AsyncHandler struct {
	storage storage.Storager
}

type SyncHandler struct {
	storage storage.Storager
}

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

func (ah AsyncHandler) MarshallMetrics() ([]byte, error) {
	jsn, err := json.Marshal(ah.storage)
	if err != nil {
		return nil, err
	}
	log.Println("marshall metrics")
	return jsn, nil
}

func (ah AsyncHandler) GetValue(w http.ResponseWriter, r *http.Request) {
	var m agent.Metrics
	err := json.NewDecoder(r.Body).Decode(&m)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
	}
	switch m.MType {
	case agent.GaugeType:
		v := ah.storage.GetGuageByID(m.ID)
		m.Value = &v
		mJSON, err := json.Marshal(m)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
		}
		w.Header().Set("content-type", "application/json")
		w.Write(mJSON)
	case agent.CounterType:
		d := ah.storage.GetCounterByID(m.ID)
		m.Delta = &d
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

func (as AsyncHandler) PostJSON(w http.ResponseWriter, r *http.Request) {
	var m agent.Metrics
	err := json.NewDecoder(r.Body).Decode(&m)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
	}
	switch m.MType {
	case agent.GaugeType:
		as.storage.SetGuage(m.ID, m.Value)
		w.WriteHeader(http.StatusOK)
	case agent.CounterType:
		as.storage.SetCounter(m.ID, m.Delta)
		w.WriteHeader(http.StatusOK)
	default:
		w.WriteHeader(http.StatusNotFound)
	}
}

func (ah AsyncHandler) PostMetric(w http.ResponseWriter, r *http.Request) {
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
		ah.storage.SetGuage(metricName, &valFloat)
		w.WriteHeader(http.StatusOK)
	case agent.CounterType:
		valInt, err := strconv.Atoi(value)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		ah.storage.SetCounterFromMemStats(metricName, int64(valInt))
		w.WriteHeader(http.StatusOK)
	default:
		w.WriteHeader(http.StatusNotImplemented)
	}
}

func (ah AsyncHandler) GetMetric(w http.ResponseWriter, r *http.Request) {
	metricType := chi.URLParam(r, "metricType")
	metricName := chi.URLParam(r, "metricName")
	if metricName == "" {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	switch metricType {
	case agent.GaugeType:
		if ah.storage.CheckGaugeByName(metricName) {
			value := strconv.FormatFloat(ah.storage.GetGuageByID(metricName), 'f', -1, 64)
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(value))
		} else {
			w.WriteHeader(http.StatusNotFound)
			return
		}
	case agent.CounterType:
		if ah.storage.CheckCounterByName(metricName) {
			c := strconv.Itoa(int(ah.storage.GetCounterByID(metricName)))
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

func (ah AsyncHandler) GetIndex(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", "text/html; charset=utf-8")
	t, err := template.ParseFiles("template/index.html")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	err = t.ExecuteTemplate(w, "index.html", ah.storage.GetGauges())
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

